package watcher

import (
	"encoder-backend/pkg/watcher/events"
	log "github.com/sirupsen/logrus"
)

// Subscribe
// add a new subscriber
func (c *Client) Subscribe(name string, channel chan<- events.Event) {

	c.btx.Lock()
	_, ok := c.subscribers[name]
	c.btx.Unlock()

	if ok {
		log.Warnf("watcher.client.subscribe: subscriber [%s] already subscribed", name)
		return
	}

	c.streams = append(c.streams, channel)

	c.btx.Lock()
	c.subscribers[name] = len(c.streams) - 1
	c.btx.Unlock()
}
