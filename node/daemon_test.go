package node

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestLocalDevDaemon(t *testing.T) {
	config := LocalDevDaemonConfig()
	daemon, err := NewDaemon(config)
	if err != nil {
		t.Fatal(err)
	}
	defer daemon.Close()
	err = daemon.Init()
	if err != nil {
		t.Fatal(err)
	}
	err = daemon.Run()
	if err != nil {
		t.Fatal(err)
	}
}
func TestTestnetDaemon(t *testing.T) {
	config := TestnetDaemonConfig()
	daemon, err := NewDaemon(config)
	if err != nil {
		t.Fatal(err)
	}
	defer daemon.Close()
	err = daemon.Init()
	if err != nil {
		t.Fatal(err)
	}
	err = daemon.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestConfig(t *testing.T) {
	config := LocalDevDaemonConfig()
	data, _ := json.Marshal(config)
	fmt.Printf("%v \n", string(data))
}
