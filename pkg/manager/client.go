package manager

import (
	"encoder-backend/pkg/database"
	"encoder-backend/pkg/repeat"
	"encoder-backend/pkg/watcher"
	"encoder-backend/pkg/watcher/events"
	"github.com/Ewan-Walker/gorm"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

// TODO [integrity] periodic scan of files based on what we have in the database

type Client struct {
	db      *gorm.DB
	wait    *sync.WaitGroup
	watcher *watcher.Client

	events chan events.Event
	queues map[events.Type]*queue
}

func New() *Client {

	c := &Client{
		watcher: watcher.New(),
		wait:    &sync.WaitGroup{},
		events:  make(chan events.Event, 1024),
		queues: map[events.Type]*queue{
			events.Scan:   {},
			events.Rename: {},
			events.Move:   {},
			events.Create: {},
			events.Delete: {},
		},
	}

	db, err := database.Connect()
	if err != nil {
		log.WithError(err).Fatal("watcher: failed to connect to database")
	}

	c.db = db

	c.watcher.Subscribe("manager", c.events)

	go c.listen()

	// TODO this can be racy if multiple actions happen to a single file

	interval := 5 * time.Second
	go repeat.Every(interval, c.createFunc())
	//go repeat.Every(interval, c.delete)
	//go repeat.Every(interval, c.rename)
	//go repeat.Every(interval, c.move)

	return c
}

// Close
// shutdown routine
func (c *Client) Close() {

	c.watcher.Close()

	close(c.events)
	c.wait.Wait()

	repeat.Close()
}

// listen
// waits for events
func (c *Client) listen() {
	c.wait.Add(1)
	defer c.wait.Done()

	for ev := range c.events {

		var err error

		switch ev.Type() {
		case events.Scan:
			c.queues[ev.Type()].Enqueue(ev)
		case events.Create:
			err = c.create(ev)
		case events.Delete:
			err = c.delete(ev)
		case events.Rename:
			err = c.rename(ev)
		case events.Move:
			err = c.move(ev)
		}

		if err != nil {
			log.WithError(err).Warn("manager.client.listener: failed to process event")
		}
	}
}
