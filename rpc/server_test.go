package rpc

import (
	"github.com/gorilla/websocket"
	"github.com/lightec-xyz/daemon/rpc/ws"
	"testing"
	"time"
)

func TestWsSimpleServer(t *testing.T) {
	server, err := NewCustomWsServer("wsSimpleSever", "localhost:8080", func(conn *websocket.Conn) {
		c := ws.NewConn(conn, func(body []byte) {
			t.Log("server receive", string(body))

		}, func() {
			t.Log("server close")
		}, false)
		c.Run()

		go func() {
			for {
				time.Sleep(3 * time.Second)
				c.Write([]byte("hello"))
			}
		}()
	})
	if err != nil {
		t.Fatal(err)
	}
	err = server.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestWsClient(t *testing.T) {
	conn, err := ws.NewWsConn("ws://localhost:8080", func(body []byte) {
		t.Log("client receive", string(body))
	}, func() {
		t.Log("client close")
	}, false)
	if err != nil {
		t.Fatal(err)
	}
	conn.Run()
	go func() {
		for {
			time.Sleep(3 * time.Second)
			conn.Write([]byte("hello"))
		}
	}()
	sig := make(chan struct{})
	<-sig

}
