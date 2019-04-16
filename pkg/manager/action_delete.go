package manager

import (
	"encoder-backend/pkg/models"
	"encoder-backend/pkg/watcher/events"
	log "github.com/sirupsen/logrus"
	"time"
)

func (c *Client) delete(list ...events.Event) error {

	//list := c.queues[events.Delete].Dequeue()

	if len(list) == 0 {
		return nil
	}

	var (
		measure = time.Now()
		deleted = 0 // count: files deleted
	)

	defer func() {
		log.WithFields(log.Fields{
			"duration": time.Since(measure),
			"deletes":  deleted,
		}).Debug("manager.client.delete: completed")
	}()

	for _, ev := range list {

		file := ev.Get()
		temp := models.File{}

		c.db.Select("*").Where(&models.File{
			Name: file.Name, Size: file.Size, PathID: file.PathID,
		}).First(&temp)

		if file.Source != temp.Source {
			continue
		}

		temp.Status = models.FileStatusDeleted

		err := c.db.Save(&temp).Error
		if err != nil {
			return err
		}

	}

	deleted = len(list)

	return nil
}
