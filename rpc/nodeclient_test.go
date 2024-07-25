package rpc

import (
	"github.com/lightec-xyz/daemon/common"
	"testing"
)

var nodeClient *NodeClient
var err error

func init() {
	url := "https://testnet.zkbtc.money/api"
	//url:="http://127.0.0.1:9780"
	nodeClient, err = NewNodeClient(url)
	if err != nil {
		panic(err)
	}
}

func TestNodeClient_ProofTask(t *testing.T) {
	task, err := nodeClient.ProofTask("")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(task)
}

func TestNodeClient_PendingTask(t *testing.T) {
	task, err := nodeClient.PendingTask()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(task)
}

func TestNodeClient_TxesByAddr(t *testing.T) {
	txes, err := nodeClient.TxesByAddr("", "")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(txes)
}

func TestNodeClient_GetTask(t *testing.T) {
	request := common.TaskRequest{
		Id:        "test_id",
		ProofType: nil,
	}
	task, err := nodeClient.GetZkProofTask(request)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(task)
}

func TestNodeClient_SubmitProof(t *testing.T) {
	result, err := nodeClient.SubmitProof(&common.SubmitProof{
		Data: []*common.ZkProofResponse{
			{},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
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
	transaction, err := nodeClient.Transaction("0x98f008ce0140ee5e1b90fe1840da40fed2c83eafb5dfc557ecea97f24f1be31d")
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
	proofInfo, err := nodeClient.ProofInfo([]string{"0x4438c9e843b35e549173658a1409c4577ad78dae5b2cda70008cb31a541c4458"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(proofInfo)

}
