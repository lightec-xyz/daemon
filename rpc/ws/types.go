package ws

import "time"

type Message struct {
	Id     int64
	Data   []byte
	Method string
	Error  string
}

func NewReqMessage(method string, data []byte) Message {
	return Message{
		Id:     time.Now().UnixNano(),
		Data:   data,
		Method: method,
		Error:  "",
	}
}

func NewErrorMsg(id int64, method string, error string) Message {
	return Message{
		Id:     id,
		Error:  error,
		Method: method,
	}
}

func NewRespMessage(id int64, method string, data []byte) Message {
	return Message{
		Id:     id,
		Data:   data,
		Method: method,
	}
}
