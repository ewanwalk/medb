package events

import "encoder-backend/pkg/models"

type renamed struct {
	generic
}

func rename(src generic) *renamed {
	return &renamed{
		generic: src,
	}
}

func (r *renamed) Type() Type {
	return Rename
}

func (r *renamed) Get() *models.File {
	file := r.generic.Get()

	if r.New == nil {
		return file
	}

	file.Name = r.New.Name()
	file.Size = r.New.Size()
	file.Checksum, _ = file.CurrentChecksum()

	return file
}
