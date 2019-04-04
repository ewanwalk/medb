package listener

import (
	"encoder-backend/pkg/watcher/events"
	"time"
)

type Option func(*options)

type options struct {
	ScanInterval       time.Duration
	ExtensionWhitelist []string
	Publisher          chan<- events.Event
}

// WithScanInterval
// Interval in milliseconds
// how often to poll the file system smaller equals faster
func WithScanInterval(interval int64) Option {
	return func(o *options) {
		o.ScanInterval = time.Duration(interval) * time.Millisecond
	}
}

// WithExtensionWhitelist
// Restricts the extension types we wish to watch to the ones provided
func WithExtensionWhitelist(extensions ...string) Option {
	return func(o *options) {
		o.ExtensionWhitelist = append(o.ExtensionWhitelist, extensions...)
	}
}

func WithPublisher() Option {
	return func(o *options) {

	}
}
