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

		temp.PathID = file.PathID
		temp.Source = file.Source
		temp.Name = file.Name
		temp.Status = models.FileStatusEnabled

		err := c.db.Save(&temp).Error
		if err != nil {
			return err
		}

	}

	renames = len(list)

	return nil
}
