package rpc

import (
	"testing"
)

var nodeClient *NodeClient
var err error

func init() {
	nodeClient, err = NewNodeClient("http://127.0.0.1:8545")
	if err != nil {
		panic(err)
	}
}

func TestProofClient(t *testing.T) {
	proofInfo, err := nodeClient.ProofInfo("ce214c17648920ea947bed3e0ad5b62837e80614b28107c8e58f8e4d7172bb32")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(proofInfo)

}
