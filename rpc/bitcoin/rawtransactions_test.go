package bitcoin

import (
	"fmt"
	"github.com/lightec-xyz/daemon/rpc/bitcoin/types"
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
	txHash, err := client.Sendrawtransaction("0200000000010150227b3619e77a3777b0c8aaf0095f1c6d8496c762544c007500bf99e6bc16860000000000000000000200ca9a3b000000001600143aa3b51760517ebf5f39739a3246a1bad1cf6bdf3046725300000000220020129e8693a89563af3524bacc46cf7f8cbaeb09383868f71d13045ebb3373fdc002473044022026191f0e5f7c573e648e1a60189b3a490d43ef6da6a53d84ab9be71e13efc8830220062ead551e295c79bd7e1baace675d6b655748e2b41aa1c2fc5843307dcad1d7012103c3b882486d88778ebebae212c03dc332566d28964293c4bd45f7b8506223160400000000")
	if err != nil {
		panic(err)
	}
	t.Log(txHash)
}

func TestClient_Getrawtransaction(t *testing.T) {
	txHash, err := client.GetRawtransaction("47810c242dc3c5df7dad90e68b247afa4484ff121bbcea8908cf7d1fe0391593")
	if err != nil {
		panic(err)
	}
	t.Log(txHash)
}
