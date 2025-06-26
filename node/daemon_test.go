package node

import (
	"encoding/json"
	"fmt"
	"github.com/lightec-xyz/daemon/logger"
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
	cfgBytes, err := os.ReadFile("/Users/red/lworkspace/lightec/audit/daemon/cmd/node/devnet.json")
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

func Test_depositScript(t *testing.T) {
	ethAddr, err := getEthAddrFromScript("6a14e96af29bb5bb124c705c69034262fbc9fbb2d5f3ddd")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(ethAddr)
}
