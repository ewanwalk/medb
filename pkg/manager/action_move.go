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

		temp.PathID = file.PathID
		temp.Source = file.Source
		temp.Name = file.Name
		temp.Status = models.FileStatusEnabled

		err := c.db.Save(&temp).Error
		if err != nil {
			return err
		}

	}

	moves = len(list)

	return nil
}
