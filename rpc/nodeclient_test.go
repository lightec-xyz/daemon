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

func TestNodeClient_TransactionsByHeight(t *testing.T) {
	txes, err := nodeClient.TransactionsByHeight(585719, "ethereum")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(txes)
}

func TestNodeClient_TransactionsByHeight01(t *testing.T) {
	txes, err := nodeClient.TransactionsByHeight(16743, "bitcoin")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(txes)
}

func TestNodeClient_Transaction(t *testing.T) {
	transaction, err := nodeClient.Transaction("0x6deff065bbaf2c9e9c12faf1d841d1f0b96502a20e6e5a864cc398cf6d54d6e4")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(transaction)
}

func TestNodeClient_Transactions(t *testing.T) {
	transaction, err := nodeClient.Transactions([]string{"0x6deff065bbaf2c9e9c12faf1d841d1f0b96502a20e6e5a864cc398cf6d54d6e4"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(transaction)
}

func TestNodeClient_ProofInfo(t *testing.T) {
	proofInfo, err := nodeClient.ProofInfo("0x6deff065bbaf2c9e9c12faf1d841d1f0b96502a20e6e5a864cc398cf6d54d6e4")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(proofInfo)

}
