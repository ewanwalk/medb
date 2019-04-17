package message

import (
	"encoder-backend/pkg/bus"
)

type Object struct {
	Base
	Body interface{} `json:"body"`
}

func Obj(name string, body interface{}) bus.Message {
	return &Object{
		Base: fromBase(name),
		Body: body,
	}
}
