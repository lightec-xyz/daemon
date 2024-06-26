package node

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/lightec-xyz/daemon/common"
	"os"
	"testing"
)

func TestLocalDevDaemon(t *testing.T) {
	cfgBytes, err := os.ReadFile("/Users/red/lworkspace/lightec/daemon/cmd/node/local.json")
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

func TestDemo(t *testing.T) {
	startHeader, err := hex.DecodeString("0000000000005fcfce1ba578d04676b167575b4890c4e75ce86d0b92b7078a50")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(startHeader)
	common.ReverseBytes(startHeader)

	middleHeader, err := hex.DecodeString("00409135508a07b7920b6de85ce7c490485b5767b17646d078a51bcecf5f00000000000056ad3b883149e96c495487cf1f5fe39296cf92b7ac12d8cb0a9dbc5cc0066b8a588c7b66f0ff0f1bdfa65698")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(middleHeader[4:36])
	t.Log(startHeader)

	t.Log("---------------------------")
	startHeader, err = hex.DecodeString("b5fbf970bf362cc3203d71022d0764ce966a9d5cee7615354e27362400000000")
	if err != nil {
		t.Fatal(err)
	}
	middleHeader, err = hex.DecodeString("01000000b5fbf970bf362cc3203d71022d0764ce966a9d5cee7615354e273624000000008c209cca50575be7aad6faf11c26af9d91fc91f9bf953c1e7d4fca44e44be3fa3d286f49ffff001d2e18e5ed")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(startHeader)
	t.Log(middleHeader[4:36])

}
