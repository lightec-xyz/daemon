package node

import (
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
func TestNetDevDaemon(t *testing.T) {
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
