package proof

import (
	"encoding/json"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/gorilla/websocket"
	btcproverUtils "github.com/lightec-xyz/btc_provers/utils"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/node"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/rpc/ws"
	"testing"
	"time"
)

func init() {
	err := logger.InitLogger(&logger.LogCfg{
		File: false,
	})
	if err != nil {
		panic(err)
	}
}

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
	server, err := rpc.NewCustomWsServer("test", "127.0.0.1:8080", func(conn *websocket.Conn) {
		t.Log("new connection")
		wsConn := ws.NewConn(conn, func(req ws.Message) (ws.Message, error) {
			return ws.Message{}, nil
		}, nil, true)
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
				t.Log("send new req")
				proof, err := proofClient.GenVerifyProof(rpc.VerifyRequest{})
				if err != nil {
					continue
				}
				t.Log(proof)
			}
		}
	}()

	err = server.Run()
	if err != nil {
		t.Fatal(err)
	}

}

func TestWsClient(t *testing.T) {
	_, err := rpc.NewCustomWsProofClientByUrl("ws://127.0.0.1:8080", func(req ws.Message) (ws.Message, error) {
		return ws.Message{}, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	exit := make(chan struct{})
	<-exit
}

func TestRpcHandler(t *testing.T) {
	worker, err := node.NewLocalWorker("", "", 1)
	if err != nil {
		t.Fatal(err)
	}
	handler := NewHandler(nil, nil, worker)
	service := ws.NewService(handler)
	t.Log(service)
	result, err := service.Call("genVerifyProof", rpc.VerifyRequest{
		TxHash:    "testhash",
		BlockHash: "blockHash",
		Data: &btcproverUtils.GrandRollupProofData{
			GenesisHash: &chainhash.Hash{},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	bytes, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(bytes))
}
