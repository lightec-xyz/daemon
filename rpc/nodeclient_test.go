package rpc

import (
	"testing"
)

var nodeClient *NodeClient
var err error

func init() {
	nodeClient, err = NewNodeClient("http://127.0.0.1:9780")
	if err != nil {
		panic(err)
	}
}

func TestProofClient(t *testing.T) {
	proofInfo, err := nodeClient.ProofInfo("0020efd9eeb78068445762a9e2e335d6ab826aaecff9cb7da24a39fc4686d468fb58")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(proofInfo)

}
