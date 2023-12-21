package rpc

import (
	"testing"
)

var nodeClient *NodeClient
var err error

func init() {
	nodeClient, err = NewNodeClient("http://127.0.0.1:30001")
	if err != nil {
		panic(err)
	}
}

func TestProofClient(t *testing.T) {
	proofInfo, err := nodeClient.ProofInfo("")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(proofInfo)

}
