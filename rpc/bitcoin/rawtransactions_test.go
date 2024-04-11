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
	txHash, err := client.Sendrawtransaction("0200000000010150ec60185108cfb91a92278f4dea0489d44d980fe2ed6e008668d60318e1f7240000000000ffffffff022a0100000000000016001499521fcaf4420357f84f548c737b41cec58fa1bac6000000000000002200205920637856dee93711f762b72324791ee0322992811087e41d739c1d2095bfeb0400483045022100d88fa38c82beab0ccda5605dafdb02dab286d9188c6d01ca1f6fd86f69d14a4602206e1bfe6bbade010e5c8f69c71861929a96c400d45607a77bb65b8391a6237a9e014830450221009cfafb07b51ccf32392f29dafd925d7e8b3df05d1705d4f21f3de206a51a190102205e8522f12b30f2fff7d355f474c8217ae07edc16e7dbc8e812100f008900eea9016952210363f549d250342df02ee8b51ad6c9148dabc587c6569761ab58aa68488bd2e2c521031cbb294f9955d80f65d9499feaeb5cb29d44c070adddd75cd48a40791d39b97121035c54e8287a7f7ba31886249fc89f295a4cb74cebf0d925f1eafe87f22fba57f953ae00000000")
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
	tx, err := client.GetRawTransaction("7a9321ddd1f9ed7433fa6fddf03601d5e4541669319ee76a48280b65f74510e9")
	if err != nil {
		panic(err)
	}
	t.Log(tx.Blockhash)
}

func TestClient_GetTransaction(t *testing.T) {
	tx, err := client.GetTransaction("f3558552478bbc873759c4dac9655a19d41efb10cedf66f36be37d97c57155cf")
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
