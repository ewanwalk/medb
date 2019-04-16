package dispatcher

import (
	"encoder-backend/pkg/config"
	"encoder-backend/pkg/database"
	"encoder-backend/pkg/encoder/dispatcher/worker"
	"github.com/ewanwalk/gorm"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"sync"
)

type Client struct {
	db *gorm.DB

	mtx     sync.Mutex
	workers []*worker.Worker
}

func New() *Client {

	c := &Client{
		mtx:     sync.Mutex{},
		workers: make([]*worker.Worker, 0),
	}

	db, err := database.Connect()
	if err != nil {
		log.WithError(err).Fatal("watcher: failed to connect to database")
	}

	c.db = db

	c.spawn()

	return c
}

func (c *Client) Close() {

	c.mtx.Lock()
	for _, w := range c.workers {
		w.Stop()
	}
	c.mtx.Unlock()

}

// spawn new workers
func (c *Client) spawn() {

	concurrency := 1
	if env := os.Getenv(config.EnvEncoderConcurrency); len(env) != 0 {
		con, err := strconv.Atoi(env)
		if err == nil {
			concurrency = con
		}
	}

	for i := 0; i < concurrency; i++ {

		w := worker.New(i, c.db)

		c.mtx.Lock()
		c.workers = append(c.workers, w)
		c.mtx.Unlock()

		go w.Start()
	}

}

// TODO get un-encoded files from the database

// TODO create workers based on the allowed concurrency
