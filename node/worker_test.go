package node

import (
	"testing"

	"github.com/lightec-xyz/daemon/rpc"
)

func TestLocalWorker_GenProof(t *testing.T) {
	client, err := rpc.NewProofClient("http://127.0.0.1:30001")
	if err != nil {
		t.Fatal(err)
	}
	worker := NewWorker(client, 1)
	proof, err := worker.ProofInfo("")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(proof)
}
