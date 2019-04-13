package worker

import (
	"context"
	"encoder-backend/pkg/encoder/handbrake"
	"encoder-backend/pkg/encoder/job"
	"encoder-backend/pkg/models"
	"encoding/json"
	"github.com/Ewan-Walker/gorm"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	WorkerWaiting = iota
	WorkerRunning
)

type Worker struct {
	id     int
	status int
	db     *gorm.DB
	logger *log.Entry

	// job sepcific
	cancel context.CancelFunc
	grp    *sync.WaitGroup

	mtx  *sync.Mutex
	job  *job.Encode
	file *models.File
}

func New(id int, db *gorm.DB) *Worker {

	w := &Worker{
		id:  id,
		db:  db,
		grp: &sync.WaitGroup{},
		logger: log.WithFields(log.Fields{
			"worker": id,
		}),
		mtx: &sync.Mutex{},
	}

	return w
}

func (w *Worker) Stop() {

	w.logger.Info("worker.stop: stopping")

	if w.cancel != nil {
		w.cancel()
	}

	w.grp.Wait()

	w.logger.Info("worker.stop: stopped")
}

// Next
// obtain the next available encode (if any are present)
// this will update the file and set its encode status to "pending"
func (w *Worker) next() *models.File {

	file := models.File{}

	txn := w.db.Begin()

	txn.Scopes(models.FileNeedsEncode).Preload("Path.QualityProfile").First(&file)

	// no file was found
	if file.ID == 0 {
		txn.Rollback()
		return &models.File{}
	}

	err := txn.Model(&file).Updates(map[string]interface{}{
		"status_encoder": models.FileEncodeStatusPending,
	}).Error
	if err != nil {
		txn.Rollback()
		return &models.File{}
	}

	txn.Commit()

	return &file
}

// Start
// initialize a workers main runtime loop, having it search for new work
func (w *Worker) Start() {

	w.grp.Add(1)
	defer w.grp.Done()

	ctx, cancel := context.WithCancel(context.Background())

	w.cancel = cancel

	backoff := w.backoff(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		default:

			w.mtx.Lock()
			w.file = w.next()
			w.mtx.Unlock()
			if w.file.ID == 0 {
				backoff(false)
				break
			}

			backoff(true)

			// TODO more advanced validation of eligibility on the file
			// check to ensure the file still exists
			if !w.file.Exists() {
				w.file.Status = models.FileStatusDeleted
				w.file.StatusEncoder = models.FileEncodeStatusNotDone
				w.db.Save(w.file)
				break
			}

			err := w.run(ctx)
			if err != nil {

				if err == handbrake.ErrBinaryNotFound {
					w.file.StatusEncoder = models.FileEncodeStatusNotDone
					w.db.Save(w.file)
				}

				w.logger.WithError(err).Warn("dispatcher.worker: failed to complete job correctly")
				break
			}

		}

		w.mtx.Lock()
		w.file = nil
		w.mtx.Unlock()

		w.status = WorkerWaiting
	}

}

// run
// attempts to run the current job
func (w *Worker) run(ctx context.Context) error {

	encode, err := job.New(w.file)
	if err != nil {
		return err
	}

	w.mtx.Lock()
	w.job = encode
	w.mtx.Unlock()

	w.status = WorkerRunning

	// flag job as started
	if err := w.onJobStart(); err != nil {

		w.file.StatusEncoder = models.FileEncodeStatusNotDone
		w.db.Save(w.file)

		return err
	}

	// execute job
	err = w.job.Run(ctx)
	if err != nil {

		if err == job.ErrCancelled || err == handbrake.ErrStagingDirectory {
			return w.onJobCancel(err)
		}

		return w.onJobError(err)
	}

	// move encoded file back
	err = os.Rename(w.job.Output(), filepath.Join(w.file.Source, w.file.Name))
	if err != nil {
		return w.onJobError(err)
	}

	// flag job as ended
	return w.onJobEnd()
}

// backoff
// used to prevent a worker from spamming the database for more work
func (w *Worker) backoff(ctx context.Context) func(bool) {

	current, start := 0*time.Second, 0*time.Second

	max := 32 * time.Second

	return func(reset bool) {

		if reset {
			current = start
			return
		}

		if current == start {
			current = 1 * time.Second
		} else if current*2 <= max {
			current *= 2
		}

		w.logger.WithFields(log.Fields{
			"waiting": current,
		}).Debug("worker.backoff: backing off for next job")

		// whichever comes first
		select {
		case <-time.After(current):
		case <-ctx.Done():
		}
	}
}

// MarshalJSON
// implement the json interface
func (w Worker) MarshalJSON() ([]byte, error) {

	data := map[string]interface{}{
		"id":     w.id,
		"status": w.status,
	}

	w.mtx.Lock()
	if w.file != nil {
		file := *w.file
		data["report"] = w.job.Report()
		data["file"] = file
	}
	w.mtx.Unlock()

	return json.Marshal(data)
}
