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
	proofInfo, err := nodeClient.ProofInfo("0x6192c3a62383898cf4368638c718b9343099f64adb654347a55163e886a43758")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(proofInfo)

}
