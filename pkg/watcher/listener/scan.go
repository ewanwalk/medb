package listener

import (
	"encoder-backend/pkg/watcher/events"
	"github.com/ewanwalk/watcher"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"time"
)

// scan
// completes a full scan of the path and its subdirectories
func (l *Listener) scan() error {

	measure := time.Now()

	count := 0

	log.WithFields(log.Fields{
		"path": l.path.Directory,
	}).Info("listener.scan: starting")

	defer func() {
		log.WithFields(log.Fields{
			"duration": time.Since(measure),
			"files":    count,
			"path":     l.path.Directory,
		}).Info("listener.scan: completed")
	}()

	return filepath.Walk(l.path.Directory, func(current string, info os.FileInfo, err error) error {

		if info.IsDir() {
			return nil
		}

		count++

		if !l.IsAllowedExtension(info.Name()) {
			return nil
		}

		l.emit(events.New(l.path.ID, watcher.Event{
			Op:       watcher.Chmod,
			Path:     current,
			FileInfo: info,
		}))

		return nil
	})
}

func (l *Listener) periodicScan(interval time.Duration) {
	for {
		select {
		case <-time.After(interval):
			err := l.scan()
			if err != nil {
				log.WithError(err).Warn("listener.scan: error while scanning library")
			}
		case <-l.quit:
			return
		}
	}
}

func (l *Listener) dummyScan() {
	initialCount := 0

	log.WithFields(log.Fields{
		"path": l.path.Directory,
	}).Infof("listener.dummyScan: starting")

	err := filepath.Walk(l.path.Directory, func(current string, info os.FileInfo, err error) error {

		if info.IsDir() {
			return nil
		}

		if !l.IsAllowedExtension(info.Name()) {
			return nil
		}

		initialCount++

		return nil
	})
	if err != nil {
		log.WithError(err).Warn("listener.dummyScan")
	}

	log.WithFields(log.Fields{
		"path":  l.path.Directory,
		"count": initialCount,
	}).Infof("listener.dummyScan: completed")

	for {
		select {
		case <-time.After(15 * time.Minute):
			count := 0
			err := filepath.Walk(l.path.Directory, func(current string, info os.FileInfo, err error) error {

				if info.IsDir() {
					return nil
				}

				if !l.IsAllowedExtension(info.Name()) {
					return nil
				}

				count++

				return nil
			})
			if err != nil {
				log.WithError(err).Warn("listener.dummyScan")
				continue
			}

			if count < initialCount {
				log.Error("listener.dummyScan: failed to correctly find files")
			}

			log.WithFields(log.Fields{
				"path":    l.path.Directory,
				"count":   initialCount,
				"compare": count,
			}).Infof("listener.dummyScan.debug: completed")
		}
	}
}
