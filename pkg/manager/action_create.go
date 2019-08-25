package manager

import (
	"encoder-backend/pkg/models"
	"encoder-backend/pkg/watcher/events"
	"errors"
	log "github.com/sirupsen/logrus"
	"time"
)

// TODO figure out if we should use the path id
type findKey struct {
	Name string
	//Checksum string
	//PathID   int64
}

type update struct {
	old models.File
	new models.File
}

var (
	ErrDuplicateFile = errors.New("create.updateandcompare: duplicate file detected")
)

func (c *Client) createFunc() func() error {
	return func() error {
		return c.create()
	}
}

func (c *Client) create(list ...events.Event) error {

	// flag should on be trigger on initial scan not sub-sequent runtime scans
	fromScan := false

	if len(list) == 0 {
		list = c.queues[events.Scan].Dequeue()
		fromScan = true
	}

	if len(list) == 0 {
		return nil
	}

	var (
		measure = time.Now()
		// the files to be created (do not currently exist)
		creates = make([]interface{}, 0)
		// the files to be updated to "exists" from "does not exist"
		updates       = make([]update, 0)
		created int64 = 0 // count: files created
		updated       = 0 // count: files updated
	)

	defer func() {
		log.WithFields(log.Fields{
			"duration": time.Since(measure),
			"created":  created,
			"updated":  updated,
		}).Debug("manager.client.create: completed")
	}()

	var err error

	if len(list) < 25 {
		for _, ev := range list {
			file := ev.Get()
			found := models.File{}

			file.Checksum, err = file.CurrentChecksum()
			// when we cannot compute the checksum we should send it to be re-evaluated
			if err != nil {
				c.queues[events.Scan].Enqueue(ev)
				continue

			}

			c.db.Where("name = ?", file.Name /*file.Checksum*/).First(&found)
			if found.ID == 0 {
				creates = append(creates, file)
				continue
			}

			err = updateAndCompare(fromScan, file, &found)
			if err != nil {
				continue
			}

			updates = append(updates, update{
				old: found,
				new: *file,
			})

		}
	} else {

		finds := make([]models.File, 0)

		c.db.Find(&finds)

		mappedFinds := make(map[findKey]models.File, len(finds))
		for _, file := range finds {
			mappedFinds[findKey{
				file.Name, //file.Checksum, //file.PathID,
			}] = file
		}

		for _, ev := range list {
			file := ev.Get()

			file.Checksum, err = file.CurrentChecksum()
			// when we cannot compute the checksum we should send it to be re-evaluated
			if err != nil {
				c.queues[events.Scan].Enqueue(ev)
				continue

			}

			found, ok := mappedFinds[findKey{
				file.Name, //file.Checksum, //file.PathID,
			}]
			if !ok {
				creates = append(creates, file)
				continue
			}

			err = updateAndCompare(fromScan, file, &found)
			if err != nil {
				continue
			}

			updates = append(updates, update{
				old: found,
				new: *file,
			})
		}
	}

	if len(creates) != 0 {

		qry := c.db.Model(creates[0]).CreateBatch(creates...)
		if qry.Error != nil {
			return qry.Error
		}

		created += qry.RowsAffected
	}

	// there should never be more than a couple hundred updates so
	// lets update individually
	if len(updates) != 0 {

		for _, update := range updates {

			err := c.db.Model(update.new).Updates(map[string]interface{}{
				"path_id":        update.new.PathID,
				"source":         update.new.Source,
				"name":           update.new.Name,
				"size":           update.new.Size,
				"checksum":       update.new.Checksum,
				"status":         models.FileStatusEnabled,
				"status_encoder": update.new.StatusEncoder,
			}).Error
			if err != nil {
				return err
			}

			if update.new.Checksum != update.old.Checksum {
				// we need to insert a revision if there is an applicable update
				err = c.db.Model(models.Revision{}).Create(&models.Revision{
					FileID:   update.old.ID,
					PathID:   update.old.PathID,
					Checksum: update.old.Checksum,
					Size:     update.old.Size,
					Encoded:  update.old.StatusEncoder,
				}).Error
				if err != nil {
					return err
				}
			}
		}

		updated = len(updates)

	}

	return nil
}

// updateAndCompare
// checks the file and compares it with the "found" file from our records
// returning nil when an update is due on the file
func updateAndCompare(fromScan bool, file, found *models.File) error {

	needsClear := fromScan && file.StatusEncoder == models.FileEncodeStatusRunning || file.StatusEncoder == models.FileEncodeStatusPending

	if found.Checksum == file.Checksum && found.Status == models.FileStatusEnabled && found.Source == file.Source && !needsClear {
		return errors.New("create.updateandcompare: update not needed")
	}

	// When a files paths do not match but the destination
	if found.PathID != file.PathID && found.ExistsShallow() {
		log.WithFields(log.Fields{
			"file": file.Name,
			"path": file.Source,
		}).Warn("duplicate file detected, ignoring...")
		return ErrDuplicateFile
	}

	file.StatusEncoder = found.StatusEncoder

	// TODO check revisions if we've seen this file in the past
	if found.Checksum != file.Checksum {
		file.StatusEncoder = models.FileEncodeStatusNotDone
	}

	// clear any potential issues with encode status
	if needsClear {
		file.StatusEncoder = models.FileEncodeStatusNotDone
	}

	file.ID = found.ID

	return nil
}
