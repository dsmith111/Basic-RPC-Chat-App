package shared

type MessageController interface {
	Send()
}

type Message struct {
	Data         string
	User         string
	IpAddress    string
	Ack          int
	MessagesSent int
}
