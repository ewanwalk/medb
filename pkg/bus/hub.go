package bus

type Hub struct {
	messages  chan Message
	registers chan subscribe

	subscribers map[string][]chan Message
}

type subscribe struct {
	name string
	send chan Message
}

func New() *Hub {

	h := &Hub{
		messages:    make(chan Message, 8<<16),
		registers:   make(chan subscribe),
		subscribers: make(map[string][]chan Message, 0),
	}

	go h.processor()

	return h
}

func (h *Hub) Broadcast(message Message) {
	h.messages <- message
}

func (h *Hub) processor() {

	for {
		select {
		case sub := <-h.registers:
			subs, ok := h.subscribers[sub.name]
			if !ok {
				subs = []chan Message{}
			}

			h.subscribers[sub.name] = append(subs, sub.send)
		case msg := <-h.messages:
			if subs, ok := h.subscribers[msg.Type()]; ok {
				for _, sub := range subs {
					sub <- msg
				}
			}

		}

	}

}

// Subscribe
// sends all messages with the provided name to the provided subscriber
func (h *Hub) Subscribe(name string, sub chan Message) {
	h.registers <- subscribe{
		name, sub,
	}
}

func (h *Hub) Unsubscribe(name string) {

}
