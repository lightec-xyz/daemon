package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

const (
	client = "client"
	server = "server"
)

type Fn func(req Message) (Message, error)

// todo
type Conn struct {
	Name      string
	conn      *websocket.Conn
	writeByte chan []byte
	exit      chan struct{}
	notify    *sync.Map
	fn        Fn
	close     func()
	waitReply bool
	lock      sync.Mutex
	closing   bool
	autoConn  bool
	mode      string
}

func NewClientConn(endpoint string, fn Fn, close func(), waitReply bool) (*Conn, error) {
	//url := url.Url{Scheme: "ws", Host: "localhost:8970", Path: "/ws"}
	conn, _, err := websocket.DefaultDialer.Dial(endpoint, nil)
	if err != nil {
		return nil, err
	}
	conn.SetPingHandler(nil)
	conn.SetPongHandler(func(appData string) error {
		err := conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return err
	})
	return &Conn{
		conn:      conn,
		writeByte: make(chan []byte, 10),
		exit:      make(chan struct{}, 1),
		notify:    new(sync.Map),
		fn:        fn,
		close:     close,
		waitReply: waitReply,
		mode:      client,
	}, nil
}

func NewConn(conn *websocket.Conn, fn Fn, close func(), waitReply bool) *Conn {
	conn.SetPingHandler(nil)
	conn.SetPongHandler(func(appData string) error {
		err := conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return err
	})
	return &Conn{
		conn:      conn,
		writeByte: make(chan []byte, 10),
		exit:      make(chan struct{}, 1),
		notify:    new(sync.Map),
		fn:        fn,
		close:     close,
		waitReply: waitReply,
		mode:      server,
	}
}

func (w *Conn) Run() {
	go w.read()
	go w.write()
	go w.heart()

}

func (w *Conn) heart() {
	ticker := time.NewTicker(60 * time.Second)
	for {
		select {
		case <-w.exit:
			return
		case <-ticker.C:
			err := w.conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				err := w.Close()
				if err != nil {
					return
				}
			}
		}
	}
}

func (w *Conn) read() {
	for {
		select {
		case <-w.exit:
			return
		default:
			messageType, data, err := w.conn.ReadMessage()
			if err != nil {
				err := w.Close()
				if err != nil {
					return
				}
			}
			switch messageType {
			case websocket.TextMessage:
				//fmt.Printf("read: %v \n", string(data))
				var req Message
				err := json.Unmarshal(data, &req)
				if err != nil {
					continue
				}
				// todo
				if w.fn != nil {
					go func(req Message) {
						reply, _ := w.fn(req)
						w.Write(reply)
					}(req)
				}
				if w.waitReply {
					if msg, ok := w.notify.Load(req.Id); ok {
						if value, ok := msg.(chan []byte); ok {
							value <- data
						}
						w.notify.Delete(req.Id)
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
					continue
				}
			}
		}
	}
}

func (w *Conn) Write(req Message) error {
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	w.writeByte <- data
	return nil
}

func (w *Conn) Call(ctx context.Context, method string, args ...interface{}) ([]byte, error) {
	if !w.waitReply {
		return nil, fmt.Errorf("need set waitReploy to true")
	}
	params, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}
	req := NewReqMessage(method, params)
	bytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp := make(chan []byte, 1)
	w.notify.Store(req.Id, resp)
	w.writeByte <- bytes
	select {
	case data := <-resp:
		if _, ok := w.notify.Load(req.Id); ok {
			w.notify.Delete(req.Id)
		}
		var reply Message
		err := json.Unmarshal(data, &reply)
		if err != nil {
			return nil, err
		}
		if reply.Error != "" {
			return nil, fmt.Errorf("%v", reply.Error)
		}
		return reply.Data, nil
	case <-ctx.Done():
		if _, ok := w.notify.Load(req.Id); ok {
			w.notify.Delete(req.Id)
		}
		return nil, fmt.Errorf("ws execute timeout")
	}
}

func (w *Conn) Close() error {
	if w.closing {
		return nil
	}
	w.lock.Lock()
	w.closing = true
	defer w.lock.Unlock()
	close(w.exit)
	err := w.conn.Close()
	if err != nil {
	}
	return nil
}
