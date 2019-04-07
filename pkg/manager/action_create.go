package manager

import (
	"encoder-backend/pkg/models"
	"encoder-backend/pkg/watcher/events"
	log "github.com/sirupsen/logrus"
	"time"
)

// TODO figure out if we should use the path id
type findKey struct {
	Name     string
	Checksum string
	//PathID   int64
}

func (c *Client) createFunc() func() error {
	return func() error {
		return c.create()
	}
}

func (c *Client) create(list ...events.Event) error {

	if len(list) == 0 {
		list = c.queues[events.Scan].Dequeue()
	}

	if len(list) == 0 {
		return nil
	}

	var (
		measure = time.Now()
		// the files to be created (do not currently exist)
		creates = make([]interface{}, 0)
		// the files to be updated to "exists" from "does not exist"
		updates       = make([]models.File, 0)
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

	// TODO what happens if a file has moved paths, we must re-assign the path
	// TODO do we allow for duplicate files (?)

	if len(list) < 25 {
		for _, ev := range list {
			file := ev.Get()
			found := models.File{}

			file.Checksum, _ = file.CurrentChecksum()

			c.db.Where("name = ? AND checksum = ?", file.Name, file.Checksum).First(&found)
			if found.ID == 0 {
				creates = append(creates, file)
				continue
			}

			if found.Status == models.FileStatusDeleted || found.Source != file.Source {
				file.ID = found.ID
				updates = append(updates, *file)
				continue
			}
		}
	} else {

		finds := make([]models.File, 0)

		c.db.Find(&finds)

		mappedFinds := make(map[findKey]models.File, len(finds))
		for _, file := range finds {
			mappedFinds[findKey{
				file.Name, file.Checksum, //file.PathID,
			}] = file
		}

		for _, ev := range list {
			file := ev.Get()

			found, ok := mappedFinds[findKey{
				file.Name, file.Checksum, //file.PathID,
			}]
			if !ok {
				file.Checksum, _ = file.CurrentChecksum()
				creates = append(creates, file)
				continue
			}

			if found.Status == models.FileStatusDeleted || found.Source != file.Source {
				file.ID = found.ID
				updates = append(updates, *file)
			}
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

		for _, file := range updates {

			err := c.db.Model(file).Updates(map[string]interface{}{
				"path_id": file.PathID,
				"source":  file.Source,
				"name":    file.Name,
				"status":  models.FileStatusEnabled,
			}).Error
			if err != nil {
				return err
			}
		}

		updated = len(updates)

	}

	return nil
}
