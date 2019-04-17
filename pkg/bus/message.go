package bus

type Message interface {
	Type() string
}
