package events

type created struct {
	generic
}

func create(src generic) *created {
	return &created{
		generic: src,
	}
}

func (c *created) Type() Type {
	return Create
}
