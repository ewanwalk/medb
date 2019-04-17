package message

import "time"

type Base struct {
	Name    string    `json:"name"`
	Created time.Time `json:"created"`
}

func fromBase(name string) Base {
	return Base{
		name,
		time.Now().UTC(),
	}
}

func (b Base) Type() string {
	return b.Name
}
