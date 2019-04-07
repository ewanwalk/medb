package listener

import (
	"encoder-backend/pkg/models"
	"encoder-backend/pkg/watcher/events"
	"github.com/Ewan-Walker/watcher"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

type Listener struct {
	*watcher.Watcher
	*options
	path models.Path

	subscribers []chan<- events.Event
}

func New(path models.Path, opts ...Option) *Listener {

	l := &Listener{
		path: path,
		options: &options{
			ScanInterval: 500 * time.Millisecond,
		},
		Watcher: watcher.New(),
	}

	info, err := os.Stat(path.Directory)
	if err != nil || !info.IsDir() {
		log.WithField("path", path.Directory).
			WithError(err).
			Warn("listener.new: directory provided may not exist")
		return l
	}

	l.FilterOps(
		watcher.Create, watcher.Remove, watcher.Rename, watcher.Move,
	)
	l.IgnoreHiddenFiles(true)

	// apply custom options
	for _, opt := range opts {
		opt(l.options)
	}

	go l.listen()

	log.WithField("path", path.Directory).Info("listener.new: adding path")

	err = l.AddRecursive(path.Directory)
	if err != nil {
		log.WithError(err).
			WithField("path", path.Directory).
			Warn("listener.new: failed to add path")
	}

	// TODO determine if we want to disable "real-time" events and only run periodic scans
	// TODO determine if we need to do a full scan periodically

	go l.run()

	go func() {
		err := l.scan()
		if err != nil {
			log.WithError(err).Warn("listener: full directory scan failed")
		}
	}()

	return l
}

func (l *Listener) run() {
	err := l.Start(l.ScanInterval)
	if err != nil {
		log.WithError(err).
			WithField("path", l.path.Directory).
			Warn("listener.new: failed to start watcher")
	}
}

// Close
// shuts down the internal event listener
func (l *Listener) Close() {
	l.Watcher.Close()
}

// listen
// listens for events coming off the file watcher
func (l *Listener) listen() {
	for {
		select {
		case ev := <-l.Event:
			// skip directories
			if ev.FileInfo.IsDir() {
				continue
			}

			if !l.IsAllowedExtension(ev.Name()) {
				continue
			}

			l.emit(events.New(l.path.ID, ev))
		case <-l.Closed:
			return
		}
	}
}

// emit
// submits an event to all current subscribers
func (l *Listener) emit(event events.Event) {
	// TODO this may become an issue due to blocking
	for _, sub := range l.subscribers {
		select {
		case sub <- event:
		default:
			log.Warn("listener.emit: subscriber channel blocked, skipping event")
		}
	}
}

// PublishTo
// provide a one-way channel that wants to subscribe to this event listeners events
func (l *Listener) PublishTo(sub chan<- events.Event) {
	l.subscribers = append(l.subscribers, sub)
}