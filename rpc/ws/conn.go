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
	ModeClient = "client"
	ModeServer = "server"
)

type Conn struct {
	conn      *websocket.Conn
	writeByte chan []byte
	exit      chan struct{}
	notifySig *sync.Map
	fn        func(body []byte)
	close     func()
	waitReply bool
	timeout   time.Duration
	mode      string
	lock      sync.Mutex
	closing   bool
}

func NewWsConn(endpoint string, fn func(body []byte), close func(), waitReply bool) (*Conn, error) {
	//url := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	conn, _, err := websocket.DefaultDialer.Dial(endpoint, nil)
	if err != nil {
		return nil, err
	}
	timeout := 60 * time.Second
	conn.SetPingHandler(nil)
	conn.SetPongHandler(func(appData string) error {
		err := conn.SetReadDeadline(time.Now().Add(timeout))
		return err
	})
	return &Conn{
		conn:      conn,
		writeByte: make(chan []byte, 10),
		exit:      make(chan struct{}, 1),
		notifySig: new(sync.Map),
		timeout:   timeout,
		fn:        fn,
		close:     close,
		waitReply: waitReply,
		mode:      ModeClient,
	}, nil
}

func NewConn(conn *websocket.Conn, fn func(body []byte), close func(), waitReply bool) *Conn {
	timeout := 60 * time.Second
	conn.SetPingHandler(nil)
	conn.SetPongHandler(func(appData string) error {
		err := conn.SetReadDeadline(time.Now().Add(timeout))
		return err
	})
	return &Conn{
		conn:      conn,
		writeByte: make(chan []byte, 10),
		exit:      make(chan struct{}, 1),
		notifySig: new(sync.Map),
		timeout:   timeout,
		fn:        fn,
		close:     close,
		waitReply: waitReply,
		mode:      ModeServer,
	}
}

func (w *Conn) Run() {
	go w.read()
	go w.write()
	go w.heart()

}

func (w *Conn) heart() {
	ticker := time.NewTicker(w.timeout / 2)
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
				if w.fn != nil {
					w.fn(data)
				}
				if w.waitReply {
					var req Message
					err := json.Unmarshal(data, &req)
					if err != nil {
						continue
					}
					if msg, ok := w.notifySig.Load(req.Id); ok {
						if value, ok := msg.(chan []byte); ok {
							value <- req.Data
						}
						w.notifySig.Delete(req.Id)
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
	resp := make(chan []byte, 1)
	w.notifySig.Store(req.Id, resp)
	w.writeByte <- bytes
	select {
	case data := <-resp:
		if _, ok := w.notifySig.Load(req.Id); ok {
			w.notifySig.Delete(req.Id)
		}
		return data, nil
	case <-ctx.Done():
		if _, ok := w.notifySig.Load(req.Id); ok {
			w.notifySig.Delete(req.Id)
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
	if w.close != nil {
		w.close()
	}
	return nil
}
