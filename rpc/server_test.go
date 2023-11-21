package rpc

import (
	"testing"
)

func TestServer_server(t *testing.T) {
	ch := make(chan int, 1)
	server, err := NewServer(":8445")
	if err != nil {
		t.Fatal(err)
	}
	defer server.Shutdown()
	<-ch
}
