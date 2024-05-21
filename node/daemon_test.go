package node

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestLocalDevDaemon(t *testing.T) {
	cfgBytes, err := os.ReadFile("/Users/red/lworkspace/lightec/daemon/cmd/node/node_local.json")
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
