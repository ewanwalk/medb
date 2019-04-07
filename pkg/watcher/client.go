package watcher

import (
	"encoder-backend/pkg/database"
	"encoder-backend/pkg/models"
	"encoder-backend/pkg/repeat"
	"encoder-backend/pkg/watcher/events"
	"encoder-backend/pkg/watcher/listener"
	"github.com/Ewan-Walker/gorm"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Client struct {
	db    *gorm.DB
	paths []models.Path
	wait  *sync.WaitGroup

	mtx       sync.Mutex
	listeners map[int64]*listener.Listener

	stream      chan events.Event
	btx         sync.Mutex
	subscribers map[string]int
	streams     []chan<- events.Event
}

var (
	instance *Client
)

func New() *Client {

	if instance != nil {
		return instance
	}

	c := &Client{
		wait:        &sync.WaitGroup{},
		mtx:         sync.Mutex{},
		listeners:   make(map[int64]*listener.Listener, 0),
		stream:      make(chan events.Event, 1024),
		btx:         sync.Mutex{},
		subscribers: make(map[string]int),
		streams:     make([]chan<- events.Event, 0),
	}

	db, err := database.Connect()
	if err != nil {
		log.WithError(err).Fatal("watcher: failed to connect to database")
	}

	c.db = db

	go c.process()

	c.load()

	go repeat.Every(15*time.Second, c.load)

	instance = c

	return c
}

// Close
// shuts down the watcher
func (c *Client) Close() {

	c.mtx.Lock()
	for _, l := range c.listeners {
		l.Close()
	}
	c.mtx.Unlock()

	close(c.stream)
	c.wait.Wait()
}

// process all events
func (c *Client) process() {
	c.wait.Add(1)
	defer c.wait.Done()

	for ev := range c.stream {
		for _, stream := range c.streams {
			select {
			case stream <- ev:
			default:
				// send on closed channel (?)
				continue
			}
		}
	}
}

// load
// attempts to load and unload any relevant paths
func (c *Client) load() error {

	paths := make([]models.Path, 0)

	err := c.db.Scopes(models.PathEnabled).Preload("QualityProfile").Find(&paths).Error
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

		log.WithFields(log.Fields{
			"path": path.Directory,
		}).Debug("watcher.client.load: removed path")

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

		measure := time.Now()

		l := listener.New(
			path,
			listener.WithScanInterval(path.EventScanInterval),
			listener.WithExtensionWhitelist("mkv", "mp4", "avi"),
		)

		l.PublishTo(c.stream)

		log.WithFields(log.Fields{
			"path":     path.Directory,
			"duration": time.Since(measure),
		}).Debug("watcher.client.load: added path")

		c.mtx.Lock()
		c.listeners[path.ID] = l
		c.mtx.Unlock()
	}

	return nil
}

// Subscribe
// add a new subscriber
func (c *Client) Subscribe(name string, channel chan<- events.Event) {

	c.btx.Lock()
	_, ok := c.subscribers[name]
	c.btx.Unlock()

	if ok {
		log.Warnf("watcher.client.subscribe: subscriber [%s] already subscribed", name)
		return
	}

	c.streams = append(c.streams, channel)

	c.btx.Lock()
	c.subscribers[name] = len(c.streams) - 1
	c.btx.Unlock()
}
