package listener

import (
	"encoder-backend/pkg/watcher/events"
	"path/filepath"
	"strings"
	"time"
)

type Option func(*options)

type options struct {
	ScanInterval       time.Duration
	ExtensionWhitelist []string
	Publisher          chan<- events.Event
}

// IsAllowedExtension
// whether or not a filename contains a whitelisted extension or not
func (o options) IsAllowedExtension(name string) bool {
	if len(o.ExtensionWhitelist) == 0 {
		return true
	}

	match := strings.ToLower(strings.Replace(filepath.Ext(name), ".", "", 1))
	for _, ext := range o.ExtensionWhitelist {
		if match == ext {
			return true
		}
	}

	return false
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
