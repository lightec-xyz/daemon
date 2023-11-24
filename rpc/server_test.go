package rpc

import (
	"testing"
)

func TestServer_server(t *testing.T) {
	ch := make(chan int, 1)
	handler := &TestHandler{}
	server, err := NewServer(":8445", handler)
	if err != nil {
		t.Fatal(err)
	}
	defer server.Shutdown()
	<-ch
}

type TestHandler struct {
}

func (t *TestHandler) Test() {

}
