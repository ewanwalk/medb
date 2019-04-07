package events

type scan struct {
	generic
}

func scanned(src generic) *scan {
	return &scan{
		generic: src,
	}
}

func (c *scan) Type() Type {
	return Scan
}
