package manager

import (
	"encoder-backend/pkg/models"
	"encoder-backend/pkg/watcher/events"
	log "github.com/sirupsen/logrus"
	"time"
)

func (c *Client) move(list ...events.Event) error {

	//list := c.queues[events.Move].Dequeue()

	if len(list) == 0 {
		return nil
	}

	var (
		measure = time.Now()
		moves   = 0 // count: files deleted
	)

	defer func() {
		log.WithFields(log.Fields{
			"duration": time.Since(measure),
			"moves":    moves,
		}).Debug("manager.client.move: completed")
	}()

	for _, ev := range list {

		file := ev.Get()
		temp := models.File{}

		c.db.Select("*").Where(&models.File{
			Name: file.Name, Checksum: file.Checksum,
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

	moves = len(list)

	return nil
}
