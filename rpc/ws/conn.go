package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

type Conn struct {
	conn      *websocket.Conn
	writeByte chan []byte
	exit      chan struct{}
	cache     *sync.Map
	fn        func(body []byte)
	waitReply bool
	timeout   time.Duration
}

func NewWsConn(endpoint string, fn func(body []byte), waitReply bool) (*Conn, error) {
	//url := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	conn, _, err := websocket.DefaultDialer.Dial(endpoint, nil)
	if err != nil {
		return nil, err
	}
	return &Conn{
		conn:      conn,
		writeByte: make(chan []byte, 10),
		exit:      make(chan struct{}, 1),
		cache:     new(sync.Map),
		timeout:   20 * time.Second,
		fn:        fn,
		waitReply: waitReply,
	}, nil
}

func NewConn(conn *websocket.Conn, fn func(body []byte), waitReply bool) *Conn {
	return &Conn{
		conn:      conn,
		writeByte: make(chan []byte, 10),
		exit:      make(chan struct{}, 1),
		cache:     new(sync.Map),
		timeout:   20 * time.Second,
		fn:        fn,
		waitReply: waitReply,
	}
}

func (w *Conn) Run() {
	go w.read()
	go w.write()
}

func (w *Conn) read() {
	for {
		select {
		case <-w.exit:
			return
		default:
			messageType, data, err := w.conn.ReadMessage()
			if err != nil {
				time.Sleep(200 * time.Millisecond) //todo
				continue
			}
			switch messageType {
			case websocket.TextMessage:
				if w.fn != nil {
					w.fn(data)
				}
				if w.waitReply {
					var req Message
					err := json.Unmarshal(data, &req)
					if err != nil {
						continue
					}
					if msg, ok := w.cache.Load(req.Id); ok {
						if value, ok := msg.(chan []byte); ok {
							value <- req.Data
						}
						w.cache.Delete(req.Id)
					}
				}
			}
		}
	}
}

func (w *Conn) write() {
	for {
		select {
		case <-w.exit:
			return
		case data, ok := <-w.writeByte:
			if ok {
				err := w.conn.WriteMessage(websocket.TextMessage, data)
				if err != nil {
					time.Sleep(500 * time.Millisecond) //todo
					continue
				}
			}
		}
	}
}

func (w *Conn) Write(data []byte) error {
	w.writeByte <- data
	return nil
}

func (w *Conn) Call(req Message) ([]byte, error) {
	if !w.waitReply {
		return nil, fmt.Errorf("need set waitReploy to true")
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), w.timeout)
	defer cancelFunc()
	bytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	msgSig := make(chan []byte, 1)
	w.cache.Store(req.Id, msgSig)
	w.writeByte <- bytes
	select {
	case data := <-msgSig:
		if _, ok := w.cache.Load(req.Id); ok {
			w.cache.Delete(req.Id)
		}
		return data, nil
	case <-ctx.Done():
		if _, ok := w.cache.Load(req.Id); ok {
			w.cache.Delete(req.Id)
		}
		return nil, fmt.Errorf("ws execute timeout")
	}

}

func (w *Conn) Close() error {
	close(w.exit)
	err := w.conn.Close()
	if err != nil {
		return err
	}
	return nil
}
