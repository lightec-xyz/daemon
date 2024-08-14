package node

import (
	"encoding/json"
	"fmt"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/rpc/ws"
	"os"
	"testing"
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
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	server.Run()

}
