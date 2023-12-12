package main

import (
	"github.com/lightec-xyz/daemon/rpc"
	"testing"
)

func TestNodeVersion(t *testing.T) {
	client, err := rpc.NewNodeClient("http://127.0.0.1:8899")
	if err != nil {
		t.Fatal(err)
	}
	version, err := client.Version()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(version)
}
