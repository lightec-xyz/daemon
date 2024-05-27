package rpc

import (
	"github.com/gorilla/websocket"
	"sync"
)

type WsManager struct {
	wsConn *sync.Map // id -> *WsConn
}

func NewWsManager() *WsManager {
	return &WsManager{
		wsConn: &sync.Map{},
	}
}

type WsConn struct {
	conn *websocket.Conn
}

func NewWsConn() *WsConn {
	return &WsConn{}
}
