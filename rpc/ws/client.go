package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

type WsClient struct {
	conn      *websocket.Conn
	writeByte chan []byte
	exit      chan struct{}
	cache     *sync.Map
	timeout   time.Duration
}

func NewWsClient(endpoint string) (*WsClient, error) {
	//url := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	conn, _, err := websocket.DefaultDialer.Dial(endpoint, nil)
	if err != nil {
		return nil, err
	}
	return &WsClient{
		conn:      conn,
		writeByte: make(chan []byte, 10),
		exit:      make(chan struct{}, 1),
		cache:     new(sync.Map),
		timeout:   20 * time.Second,
	}, nil
}

func (w *WsClient) Run() {
	go w.read()
	go w.write()
}

func (w *WsClient) read() {
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

func (w *WsClient) write() {
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

func (w *WsClient) Execute(req Message) ([]byte, error) {
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

func (w *WsClient) Close() error {
	close(w.exit)
	return w.conn.Close()
}
