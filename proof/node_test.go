package proof

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/node"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/rpc/ws"
)

func init() {
	err := logger.InitLogger(&logger.LogCfg{
		File:     false,
		IsStdout: true,
	})
	if err != nil {
		panic(err)
	}
}

func TestClientModeProof(t *testing.T) {
	config := NewTestClientModeConfig()
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
	config := NewTestClusterModeConfig()

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

func TestCustomModeProof(t *testing.T) {
	config := NewTestCustomModeConfig()
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
	server, err := rpc.NewCustomWsServer("zkbtc", "127.0.0.1:8080", func(opt *rpc.WsOpt) error {
		t.Log("new connection")
		wsConn := ws.NewConn(opt.Conn, nil, nil, true)
		wsConn.Run()
		proofClient, err = rpc.NewCustomWsProofClient(wsConn)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		for {
			time.Sleep(10 * time.Second)
			if proofClient != nil {
				t.Log("send new req")
				proof, err := proofClient.TxInEth2Prove(&rpc.TxInEth2ProveRequest{})
				if err != nil {
					t.Log(err)
					continue
				}
				t.Logf("response: %v \n", proof)
			}
		}
	}()
	err = server.Run()
	if err != nil {
		t.Fatal(err)
	}

}

func TestWsClient(t *testing.T) {
	conn, err := ws.NewClientConn("ws://127.0.0.1:8970/ws", func(req ws.Message) (ws.Message, error) {
		t.Logf("clinet receive new req: %v \n", req)
		response := rpc.TxInEth2ProveRequest{
			TxHash: "testVerifyResp",
		}
		data, err := json.Marshal(response)
		if err != nil {
			return ws.Message{}, err
		}
		t.Logf("clinet response: %v \n", string(data))
		return ws.NewRespMessage(req.Id, req.Method, data), nil
	}, nil, false)

	if err != nil {
		t.Fatal(err)
	}
	go conn.Run()
	exit := make(chan struct{})
	<-exit
}

func TestRpcHandler(t *testing.T) {
	worker, err := node.NewLocalWorker("", "", "", "", 1, 0)
	if err != nil {
		t.Fatal(err)
	}
	handler := NewHandler(nil, nil, worker)
	service := ws.NewService(handler)
	t.Log(service)
	request := rpc.TxInEth2ProveRequest{
		TxHash: "testhash",
	}
	param, err := json.Marshal([]interface{}{request})
	if err != nil {
		t.Fatal(err)
	}
	result, err := service.Call("genVerifyProof", param)
	if err != nil {
		t.Fatal(err)
	}
	bytes, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(bytes))
}
