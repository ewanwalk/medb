package events

import (
	"encoder-backend/pkg/config"
	"encoder-backend/pkg/models"
	"github.com/Ewan-Walker/watcher"
	"runtime"
	"strings"
)

type Type int64

const (
	Create Type = iota
	Move
	Rename
	Delete
)

// TODO categorize events
type Event interface {
	Type() Type
	Get() *models.File
}

type generic struct {
	watcher.Event
}

func New(src watcher.Event) Event {

	// TODO possibly initialize File here
	g := generic{
		Event: src,
	}

	switch src.Op {
	case watcher.Create:
		return create(g)
	case watcher.Move:
		return move(g)
	case watcher.Remove:
		return delete(g)
	case watcher.Rename:
		return rename(g)
	}

	return nil
}

// Get
// generic `get file` function for all events
// to be overridden dependant on the parent event type
func (g *generic) Get() *models.File {

	path := strings.Split(g.Abs(), config.Separator)

	return &models.File{
		Name:     g.FileInfo.Name(),
		Size:     g.FileInfo.Size(),
		Checksum: "",
		Source:   strings.Join(path[0:len(path)-1], config.Separator),
	}
}

// Abs
// obtain the absolute path of the file related to the events
func (g *generic) Abs() string {
	if g.Op == watcher.Move || g.Op == watcher.Rename {
		if split := strings.Split(g.Path, " -> "); len(split) == 2 {
			return split[1]
		}
	}

	return g.Path
}
