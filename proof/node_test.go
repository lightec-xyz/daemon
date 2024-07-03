package proof

import (
	"github.com/gorilla/websocket"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/rpc/ws"
	"testing"
	"time"
)

func TestClientModeProof(t *testing.T) {
	config := NewClientModeConfig()
	node, err := NewNode(config)
	if err != nil {
		t.Fatal(err)
	}
	err = node.Init()
	if err != nil {
		t.Fatal(err)
	}
	err = node.Start()
	if err != nil {
		t.Fatal(err)
	}
}

func TestClusterModeProof(t *testing.T) {
	config := NewClusterModeConfig()

	node, err := NewNode(config)
	if err != nil {
		t.Fatal(err)
	}
	err = node.Init()
	if err != nil {
		t.Fatal(err)
	}
	err = node.Start()
	if err != nil {
		t.Fatal(err)
	}
}

func TestWsServer(t *testing.T) {
	var proofClient *rpc.ProofClient
	var err error
	server, err := rpc.NewCustomWsServer("test", "ws://127.0.0.1:8545", func(conn *websocket.Conn) {
		wsConn := ws.NewConn(conn, func(body []byte) {
		}, func() {

		}, true)
		wsConn.Run()
		proofClient, err = rpc.NewCustomWsProofClient(wsConn)
		if err != nil {
			return
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		for {
			time.Sleep(10 * time.Second)
			if proofClient != nil {

			}
		}
	}()

	err = server.Run()
	if err != nil {
		t.Fatal(err)
	}

}
