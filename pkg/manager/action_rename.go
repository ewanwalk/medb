package manager

import (
	"encoder-backend/pkg/models"
	"encoder-backend/pkg/watcher/events"
	log "github.com/sirupsen/logrus"
	"time"
)

func (c *Client) rename(list ...events.Event) error {

	//list := c.queues[events.Rename].Dequeue()

	if len(list) == 0 {
		return nil
	}

	var (
		measure = time.Now()
		renames = 0 // count: files deleted
	)

	defer func() {
		log.WithFields(log.Fields{
			"duration": time.Since(measure),
			"renames":  renames,
		}).Debug("manager.client.rename: completed")
	}()

	for _, ev := range list {

		file := ev.Get()
		temp := models.File{}

		c.db.Select("*").Where(&models.File{
			Checksum: file.Checksum, PathID: file.PathID,
		}).First(&temp)

		err := c.db.Model(temp).Updates(map[string]interface{}{
			"path_id": file.PathID,
			"source":  file.Source,
			"name":    file.Name,
			"status":  models.FileStatusEnabled,
		}).Error
		if err != nil {
			return err
		}

	}

	renames = len(list)

	return nil
}
