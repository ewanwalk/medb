package socket

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	// time allowed to write message to peer
	writeWait = 10 * time.Second
	// time allowed to read  next pong from peer
	pongWait = 60 * time.Second
	// send pings to peer with this period / must be less than pongWait
	pingPeriod = (pongWait * 9) / 10
	// max message size allowed from peer
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	//space = []byte{' '}
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// TODO add validation
			return true
		},
	}
)

type Client struct {
	hub  *hub
	conn *websocket.Conn
	send chan []byte
}

func NewClient(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.WithError(err).Warn("http.web.sockets: failed to upgrade connection")
		return
	}

	client := &Client{
		hub:  std,
		conn: conn,
		send: make(chan []byte, 256),
	}

	client.hub.register <- client

	go client.write()
	//go client.read()
}

/*func (c *Client) read() {
	defer func() {
		c.hub.unregister <- c
		_ = c.conn.Close()
	} ()

	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(data string) error {
		_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.WithError(err).Warn("http.web.sockets: unexpected connection closure")
			}
			break
		}
		// TODO maybe accept messages
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.hub.broadcast <- message
	}
}*/

func (c *Client) write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			_, _ = w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				_, _ = w.Write(newline)
				_, _ = w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}

		}
	}
}
