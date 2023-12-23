package main

import (
	"fmt"
	"github.com/lightec-xyz/daemon/rpc"
	"strings"
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

func TestDemo(t *testing.T) {
	src := "dddd_aaaaa"
	s := src[strings.Index(src, "_")+1:]
	fmt.Printf(s)
}
