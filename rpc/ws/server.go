package ws

import (
	"github.com/gorilla/websocket"
)

type Server struct {
	conn *websocket.Conn
	exit chan struct{}
	fn   func([]byte) []byte
}

func NewServer(conn *websocket.Conn) *Server {
	return &Server{
		conn: conn,
		exit: make(chan struct{}, 1),
	}
}

func (s *Server) Run() {
	go s.read()
}

func (s *Server) read() {
	for {
		select {
		case <-s.exit:
			return
		default:
			msgType, data, err := s.conn.ReadMessage()
			if err != nil {
				continue
			}
			switch msgType {
			case websocket.TextMessage:
				resp := s.fn(data)
				err := s.conn.WriteMessage(websocket.TextMessage, resp)
				if err != nil {
					continue
				}
			}
		}

	}
}

func (s *Server) Close() error {
	close(s.exit)
	return s.conn.Close()
}
