package watcher

import (
	"encoder-backend/pkg/database"
	"encoder-backend/pkg/models"
	"encoder-backend/pkg/watcher/events"
	"encoder-backend/pkg/watcher/listener"
	"github.com/Ewan-Walker/gorm"
	log "github.com/sirupsen/logrus"
	"sync"
)

type Client struct {
	db    *gorm.DB
	paths []models.Path

	mtx       sync.Mutex
	listeners map[int64]*listener.Listener

	stream chan events.Event
}

func New() *Client {

	c := &Client{
		mtx:       sync.Mutex{},
		listeners: make(map[int64]*listener.Listener, 0),
		stream:    make(chan events.Event, 1024),
	}

	db, err := database.Connect()
	if err != nil {
		log.WithError(err).Fatal("watcher: failed to connect to database")
	}

	c.db = db

	return c
}

func (c *Client) Close() {

	c.mtx.Lock()
	for _, l := range c.listeners {
		l.Close()
	}
	c.mtx.Unlock()

	close(c.stream)
}

// load
// attempts to load and unload any relevant paths
func (c *Client) load() error {

	paths := make([]models.Path, 0)

	err := c.db.Scopes(models.PathEnabled).Preload("QualityProfile").Find(paths).Error
	if err != nil {
		return err
	}

	// check to see if we need to remove any paths
	remove := make([]models.Path, 0)

	for _, path := range c.paths {
		exists := false
		for _, match := range paths {
			if match.ID != path.ID {
				continue
			}

			exists = true
			break
		}

		if !exists {
			remove = append(remove, path)
		}
	}

	c.paths = paths

	// shutdown non-existent paths
	for _, path := range remove {
		c.mtx.Lock()
		l, ok := c.listeners[path.ID]
		c.mtx.Unlock()
		if !ok {
			continue
		}

		l.Close()

		c.mtx.Lock()
		delete(c.listeners, path.ID)
		c.mtx.Unlock()
	}

	// start any new paths
	for _, path := range c.paths {

		c.mtx.Lock()
		_, ok := c.listeners[path.ID]
		c.mtx.Unlock()
		if ok {
			continue
		}

		l := listener.New(
			path,
			listener.WithScanInterval(path.EventScanInterval),
			listener.WithExtensionWhitelist("mkv", "mp4", "avi"),
		)

		l.PublishTo(c.stream)

		c.mtx.Lock()
		c.listeners[path.ID] = l
		c.mtx.Unlock()
	}

	return nil
}

func (c *Client) Subscribe() chan events.Event {
	return c.stream
}
