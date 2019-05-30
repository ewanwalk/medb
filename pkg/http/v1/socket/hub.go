package socket

import (
	"encoder-backend/pkg/bus"
	"encoder-backend/pkg/encoder/dispatcher/worker"
	"encoding/json"
)

type hub struct {
	// registered clients
	clients map[*Client]bool
	// inbound messages from clients
	broadcast chan []byte
	internal  chan bus.Message
	// register requests from clients
	register chan *Client
	// unregister requests from clients
	unregister chan *Client
}

var (
	std = newHub()
)

func newHub() *hub {
	h := &hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		internal:   make(chan bus.Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	go h.run()

	return h
}

func (h *hub) run() {

	// subscribe to internal hub
	for _, name := range []string{
		worker.MessageStart,
		worker.MessageStop,
		worker.MessageTick,
		//worker.MessageStatus,
	} {
		bus.Subscribe(name, h.internal)
	}

	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		case message := <-h.internal:
			payload, err := json.Marshal(map[string]interface{}{"data": message})
			if err != nil {
				continue
			}

			for client := range h.clients {
				select {
				case client.send <- payload:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
