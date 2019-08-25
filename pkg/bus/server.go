package bus

var (
	std = New()
)

// Broadcast
// send a message on the standard hub
func Broadcast(message Message) {
	std.Broadcast(message)
}

// Subscribe
// listen for the message type (name) on the standard hub
func Subscribe(name string, sub chan Message) {
	std.Subscribe(name, sub)
}
