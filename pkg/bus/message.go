package bus

type Message interface {
	Type() string
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(b []byte) error
}
