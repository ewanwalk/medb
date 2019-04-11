package events

import (
	"encoder-backend/pkg/config"
	"encoder-backend/pkg/models"
	"github.com/Ewan-Walker/watcher"
	"strings"
)

type Type int64

const (
	Create Type = iota
	Move
	Rename
	Delete
	Scan
)

type Event interface {
	Type() Type
	Get() *models.File
}

type generic struct {
	PathID int64
	watcher.Event
}

func (t Type) String() string {
	switch t {
	case Create:
		return "create"
	case Move:
		return "move"
	case Rename:
		return "rename"
	case Delete:
		return "delete"
	}

	return ""
}

func New(id int64, src watcher.Event) Event {

	g := generic{
		PathID: id,
		Event:  src,
	}

	switch src.Op {
	case watcher.Chmod: // placeholder for a create due to initial scan
		return scanned(g)
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
		PathID:   g.PathID,
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
