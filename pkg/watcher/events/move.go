package events

import (
	"encoder-backend/pkg/models"
)

type moved struct {
	generic
}

func move(src generic) *moved {
	return &moved{
		generic: src,
	}
}

func (m *moved) Type() Type {
	return Move
}

func (m *moved) Get() *models.File {
	file := m.generic.Get()

	if m.New == nil {
		return file
	}

	file.Name = m.New.Name()
	file.Size = m.New.Size()
	file.Checksum, _ = file.CurrentChecksum()

	return file
}
