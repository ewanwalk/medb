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
