package encoder

import "encoder-backend/pkg/encoder/dispatcher"

type Client struct {
	dispatch *dispatcher.Client
}

var (
	global *Client
)

func New() *Client {

	if global != nil {
		return global
	}

	c := &Client{
		dispatch: dispatcher.New(),
	}

	global = c

	return c
}

func (c *Client) Close() {
	c.dispatch.Close()
}
