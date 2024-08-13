package node

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	grUtil "github.com/lightec-xyz/btc_provers/utils/txinchain"
	"github.com/lightec-xyz/daemon/logger"
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

func TestLocalDevDaemon(t *testing.T) {
	cfgBytes, err := os.ReadFile("/Users/red/lworkspace/lightec/daemon/daemon/cmd/node/local.json")
	if err != nil {
		t.Fatal(err)
	}
	var runCfg RunConfig
	err = json.Unmarshal(cfgBytes, &runCfg)
	if err != nil {
		t.Fatal(err)
	}
	config, err := NewConfig(runCfg)
	if err != nil {
		t.Fatal(err)
	}
	marshal, err := json.Marshal(config)
	fmt.Println(string(marshal))
	daemon, err := NewDaemon(config)
	if err != nil {
		t.Fatal(err)
	}
	err = daemon.Init()
	if err != nil {
		t.Fatal(err)
	}
	err = daemon.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestServer(t *testing.T) {
	handler := NewHandler(nil, nil, nil, nil, nil, nil)
	server, err := rpc.NewServer(RpcRegisterName, "127.0.0.1:8970", handler, nil, func(opt *rpc.WsOpt) error {
		t.Logf("new connection: %v \n", opt.Id)
		newConn := ws.NewConn(opt.Conn, nil, nil, true)
		go newConn.Run()
		go func() {
			for {
				time.Sleep(5 * time.Second)
				client, err := rpc.NewCustomWsProofClient(newConn)
				if err != nil {
					t.Error(err)
					continue
				}
				result, err := client.GenVerifyProof(rpc.VerifyRequest{
					TxHash:    "testhash",
					BlockHash: "blockHash",
					Data: &grUtil.GrandRollupProofData{
						GenesisHash: &chainhash.Hash{},
					}})
				if err != nil {
					t.Error(err)
					continue
				}
				t.Log(result)
			}
		}()
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	server.Run()

}
