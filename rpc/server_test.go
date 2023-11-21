package rpc

import (
	"testing"
	"time"
)

func TestServer_server(t *testing.T) {
	server, err := NewServer(":8445")
	if err != nil {
		t.Fatal(err)
	}
	defer server.Shutdown()
	time.Sleep(10 * time.Minute)
}
