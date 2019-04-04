package events

type deleted struct {
	generic
}

func delete(src generic) *deleted {
	return &deleted{
		generic: src,
	}
}

func (c *deleted) Type() Type {
	return Delete
}
