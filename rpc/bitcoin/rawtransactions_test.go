package bitcoin

import (
	"fmt"
	"github.com/lightec-xyz/daemon/rpc/bitcoin/types"
	"math/big"
	"testing"
)

func TestClient_CreateTransaction(t *testing.T) {
	var inputs []types.TxIn
	inputs = append(inputs, types.TxIn{
		TxId:     "8616bce699bf0075004c5462c796846d1c5f09f0aac8b077377ae719367b2250",
		Vout:     0,
		Sequence: 0,
	})
	var outs []types.TxOut
	outs = append(outs, types.TxOut{
		Address: "bcrt1q823m29mq29lt7heewwdry34phtgu767lgv5mx6",
		Amount:  10,
	},
		types.TxOut{
			Address: "bcrt1qz20gdyagj4367dfyhtxydnml3jawkzfc8p50w8gnq30tkvmnlhqq9wudl2",
			Amount:  13.99998000,
		})
	createrawtransaction, err := client.Createrawtransaction(inputs, outs)
	if err != nil {
		panic(err)
	}
	fmt.Println(createrawtransaction)
}

func TestClient_SigRawTransaction(t *testing.T) {
	hexData := "020000000150227b3619e77a3777b0c8aaf0095f1c6d8496c762544c007500bf99e6bc16860000000000000000000200ca9a3b000000001600143aa3b51760517ebf5f39739a3246a1bad1cf6bdf3046725300000000220020129e8693a89563af3524bacc46cf7f8cbaeb09383868f71d13045ebb3373fdc000000000"
	privateKys := []string{"cS6QFVqv9v9XREMTUaYQ38GfkLcqXcTYDV5oMt4kw76MGJvfF5pV"}
	var inputs []types.TxIn
	inputs = append(inputs, types.TxIn{
		TxId:         "8616bce699bf0075004c5462c796846d1c5f09f0aac8b077377ae719367b2250",
		Vout:         0,
		ScriptPubKey: "00143aa3b51760517ebf5f39739a3246a1bad1cf6bdf",
		Sequence:     0,
		Amount:       23.99998500,
	})
	signrawtransaction, err := client.Signrawtransactionwithkey(hexData, privateKys, inputs)
	if err != nil {
		panic(err)
	}
	fmt.Println(signrawtransaction)
}
func TestClient_SendRawTransaction(t *testing.T) {
	txHash, err := client.Sendrawtransaction("02000000000102740cd4d6f30fc3d7a65e1a112782c42744b47d7791d630634d56d5a716d6c5720100000000fffffffff4aecbb6a943cb8a8ae8124330e70cfe38dae0c5fd5830e8bee70a8d678b9c760000000000ffffffff0130150d8f00000000160014d7fae4fbdc8bf6c86a08c7177c5d06683754ea71040047304402204ef8f9d9c6230d11126aa5c1c3bcd8e7d5d6d22937c03bdea343641eef5084df0220312d94fca8aa7878a988ab632b15189922ce7e7d5b9564b15579ec7b3037a1f701473044022045c61947c5f2d37e9f13f71aa1e50fc1f705b10849f750e8d910acfbffbd0b0602205dfce8ec4b06c0b13dcc1059a99b64e60a82b8b792b4ce93bce5ab238e7dbd260169522103bd96c4d06aa773e5d282f0b6bccd1fb91268484918648ccda1ae768209edb0502103aa9c4245340a02864c903f7f9e7bc9ef1cc374093aacbf72b614002f6d8c8c222103351a7971bf7ed886fca99aebdc3b195fc79ffe93b499e2309a4e69ab115405e053ae0400483045022100fd6bfa0e45beadc51b11cf4e72cc4dc95d7160e72a509f766f130a1a56b138330220181d15010661d0fbdb8fa80012df0067c972acef2cc36b57bedd5f46166acd6901483045022100c2924fdaba8494aecd8cef6b7cd1a57261cf297f5762a8d1873d0d7f238f7ad8022043b21a46eef1f26da76a121f05e81f2a900b2e69c98661113fb4400ea279327a0169522103bd96c4d06aa773e5d282f0b6bccd1fb91268484918648ccda1ae768209edb0502103aa9c4245340a02864c903f7f9e7bc9ef1cc374093aacbf72b614002f6d8c8c222103351a7971bf7ed886fca99aebdc3b195fc79ffe93b499e2309a4e69ab115405e053ae00000000")
	if err != nil {
		panic(err)
	}
	t.Log(txHash)
}

func TestClient_GetUtxoByTxId(t *testing.T) {
	utxo, err := client.GetUtxoByTxId("b5a43b150d8f9a305b9f19b11411f6a166f57080c8347b258ba48ebce77a2bc9", 4)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(utxo)
}

func TestClient_Getrawtransaction(t *testing.T) {
	//aabd19556b19981ae317a26d47bb6f69bb0daa230612dd9a857a11aae5c22cab
	//09f15de0a93ead808978d31e43644ef2d38e7f6f9cdf70a195ae79f93046eaf3
	tx, err := client.GetRawTransaction("f1600eaff05f75671978bb86c27a2de68efe70dceaefda73f8d5a1108bc3660f")
	if err != nil {
		panic(err)
	}
	t.Log(tx)
}

func TestClient_GetTransaction(t *testing.T) {
	tx, err := client.GetTransaction("6108b2003208d310a4afe49ec963dbb62a17f2407af9aeec0ef4ded5ab95d8bd")
	if err != nil {
		panic(err)
	}
	t.Log(tx)
}

func TestDmeo(t *testing.T) {
	amount, _ := big.NewFloat(0.0009975 * 100000000.0).Int64()
	t.Log(amount)
	fmt.Printf("%0.8f", 1.12313123123123)
}
