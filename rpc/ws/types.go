package ws

type Message struct {
	Id     int64
	Data   []byte
	Method string
	error  error
}
