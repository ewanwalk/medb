package message

import (
	"encoder-backend/pkg/bus"
)

type Plain struct {
	Base
	Body string `json:"body"`
}

func Text(name, body string) bus.Message {
	return &Plain{
		Base: fromBase(name),
		Body: body,
	}
}
