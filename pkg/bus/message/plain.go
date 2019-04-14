package message

import (
	"encoder-backend/pkg/bus"
	"encoding/json"
	"time"
)

type Plain struct {
	Name    string
	Body    string
	Created time.Time
}

func Text(name, body string) bus.Message {
	return &Plain{
		Name:    name,
		Body:    body,
		Created: time.Now().UTC(),
	}
}

func (t *Plain) Type() string {
	return t.Name
}

func (t *Plain) MarshalJSON() ([]byte, error) {
	return json.Marshal(t)
}

func (t *Plain) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, t)
}
