package rpc

import "testing"

func TestProofClient_GenZkProof(t *testing.T) {
	client, err := NewWsProofClient("ws://127.0.0.1:30001")
	if err != nil {
		t.Fatal(err)
	}
	proof, err := client.GenZkProof(ProofRequest{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(proof)
	proofInfo, err := client.ProofInfo("")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(proofInfo)
}
