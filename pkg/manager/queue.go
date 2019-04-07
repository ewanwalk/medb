package manager

import "encoder-backend/pkg/watcher/events"

type queue struct {
	data []events.Event
}

func (q *queue) Enqueue(ev events.Event) {
	q.data = append(q.data, ev)
}

func (q *queue) Dequeue() (batch []events.Event) {
	q.data, batch = q.data[:0], q.data[0:]
	return
}
