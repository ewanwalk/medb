package listener

import (
	"github.com/radovskyb/watcher"
)

type listener struct {
	watcher.Watcher
}

func newListener(path string) *listener {

	l := &listener{}

	go l.listen()

	return l
}

func (l *listener) Close() {
	l.Watcher.Close()
}

func (l *listener) listen() {
	for {
		select {
		case ev := <-l.Event:
			// skip directories
			if ev.FileInfo.IsDir() {
				continue
			}

		case <-l.Closed:
			return
		}
	}
}
