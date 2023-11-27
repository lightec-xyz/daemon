package proof

import (
	"github.com/lightec-xyz/daemon/rpc"
	"testing"
)

func TestProofServer(t *testing.T) {
	ch := make(chan int, 1)
	server, err := rpc.NewServer("127.0.0.1:8545", &Handler{})
	if err != nil {
		t.Fatal(err.Error())
	}
	defer server.Shutdown()
	<-ch
}
