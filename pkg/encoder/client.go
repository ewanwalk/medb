package encoder

import "encoder-backend/pkg/encoder/dispatcher"

type Client struct {
	dispatch *dispatcher.Client
}

func New() *Client {

	c := &Client{
		dispatch: dispatcher.New(),
	}

	return c
}

func (c *Client) Close() {
	c.dispatch.Close()
}
