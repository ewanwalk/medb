package repeat

import (
	log "github.com/sirupsen/logrus"
	"math/rand"
	"sync"
	"time"
)

var (
	wait = &sync.WaitGroup{}
	run  = make(chan struct{})
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Close() {
	close(run)
	wait.Wait()
}

// Every
// repeats a task every X provided interval
func Every(interval time.Duration, task Task) {

	wait.Add(1)
	defer wait.Done()

	for {

		val := int(interval / time.Millisecond)
		// this jitter helps prevent many tasks from overlapping
		jitter := time.Duration(rand.Intn(val)) * time.Millisecond

		select {
		case <-time.After(interval + jitter):
			err := task()
			if err != nil {
				log.WithError(err).Warn("repeat.every: task completed with error")
			}
		case <-run:
			return
		}
	}

}
