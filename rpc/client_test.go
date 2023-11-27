package rpc

import (
	"testing"
)

func TestProofClient(t *testing.T) {
	client, err := NewProofClient("http://127.0.0.1:8545")
	if err != nil {
		t.Fatal(err)
	}
	btcProofResponse, err := client.GenZkProof(ProofRequest{
		TxId:    "test",
		EthAddr: "test",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(btcProofResponse)
}
