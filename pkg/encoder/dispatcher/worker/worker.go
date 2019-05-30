package worker

import (
	"context"
	"encoder-backend/pkg/bus"
	"encoder-backend/pkg/bus/message"
	"encoder-backend/pkg/encoder/handbrake"
	"encoder-backend/pkg/encoder/job"
	"encoder-backend/pkg/models"
	"encoding/json"
	"github.com/ewanwalk/gorm"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	Waiting = iota
	Running
)

const (
	MessageStart  = "worker.start"
	MessageStop   = "worker.stop"
	MessageTick   = "worker.tick"
	MessageStatus = "worker.status"
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
	if file.ID == 0 || file.Path == nil || file.Path.QualityProfile == nil {
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

	go w.tick(ctx)

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

			bus.Broadcast(message.Obj(MessageStart, map[string]interface{}{
				"id": w.file.ID,
			}))

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

		bus.Broadcast(message.Obj(MessageStop, map[string]interface{}{
			"id": w.file.ID,
		}))

		w.mtx.Lock()
		w.file = nil
		w.mtx.Unlock()

		w.status = Waiting
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

	w.status = Running

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

func (w *Worker) tick(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(2 * time.Second):
			if w.status != Running {
				bus.Broadcast(message.Obj(MessageStatus, map[string]interface{}{
					"status": "ok",
					"id":     w.id,
				}))
				break
			}

			bus.Broadcast(message.Obj(MessageTick, w.report()))
		}
	}
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

// report
func (w Worker) report() map[string]interface{} {
	data := map[string]interface{}{
		"id":     w.id,
		"status": w.status,
	}

	w.mtx.Lock()
	if w.file != nil {
		file := *w.file
		data["report"] = w.job.Report()
		data["file"] = map[string]interface{}{
			"id":   file.ID,
			"name": file.Name,
			"encodes": []map[string]interface{}{
				{"id": file.Encodes[0].ID},
			},
		}
	}
	w.mtx.Unlock()

	return data
}

// MarshalJSON
// implement the json interface
func (w Worker) MarshalJSON() ([]byte, error) {
	return json.Marshal(w.report())
}
