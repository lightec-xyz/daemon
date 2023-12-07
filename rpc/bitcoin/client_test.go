package bitcoin

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/lightec-xyz/daemon/transaction/bitcoin"
	"math/big"
	"strconv"
	"testing"
)

var client *Client
var err error

func init() {
	//url := "https://go.getblock.io/d54c59f635654cc082de1f3fd14e5d02"
	url := "http://127.0.0.1:8332"
	user := "lightec"
	pwd := "abcd1234"
	network := "regtest"
	client, err = NewClient(url, user, pwd, network)
	if err != nil {
		panic(err)
	}
}

func TestClient_GetBlockHeader(t *testing.T) {
	header, err := client.GetBlockHeader("0ca10b19b94eedc77da894b15c6a62407e73a8eb312c8d9befd1769448f92ce7")
	if err != nil {
		panic(err)
	}
	fmt.Println(header)
}

func TestClient_GetBlockCount1(t *testing.T) {
	blockCount, err := client.GetBlockCount()
	if err != nil {
		panic(err)
	}
	fmt.Println(blockCount)
}

func TestClient_GetBlockHash(t *testing.T) {
	hash, err := client.GetBlockHash(200)
	if err != nil {
		panic(err)
	}
	fmt.Println(hash)
}

func TestClient_GetBlockTx(t *testing.T) {
	hash, err := client.GetBlockHash(2540940)
	if err != nil {
		panic(err)
	}
	blockWithTx, err := client.GetBlock(hash)
	if err != nil {
		panic(err)
	}
	fmt.Println(blockWithTx)
}

func TestMultiTransactionBuilder(t *testing.T) {
	secrerts := []string{
		"b26dbaab82d9ebd8f37c88bbe56e22bf9cb21150c96dfb35ece4b787d3710d3301",
		"62dd5835dc2ce7f4f40eea1b88c816043d288532c8bb91964adef9bc0f0b4b7201",
		"9ff573d948c80fa1a50da6f66229b4bede9ec3fb482dd126f58d3acfb4b2979801",
	}
	var privateKeys []*btcec.PrivateKey
	var pubkeylist [][]byte
	for _, secret := range secrerts {
		hexPriv, err := hex.DecodeString(secret)
		if err != nil {
			t.Fatal(err)
		}
		privateKey, publikey := btcec.PrivKeyFromBytes(hexPriv)
		if err != nil {
			t.Fatal(err)
		}
		privateKeys = append(privateKeys, privateKey)
		pubkeylist = append(pubkeylist, publikey.SerializeCompressed())

	}
	txInputs := []bitcoin.TxIn{
		{
			Hash:     "1960288483d27a24c33b148a3e3d75a64ea99531f77406a1a719c0dad34670e9",
			VOut:     1,
			PkScript: "0020efd9eeb78068445762a9e2e335d6ab826aaecff9cb7da24a39fc4686d468fb58",
			Amount:   9975,
		},
	}
	pk1, _ := hex.DecodeString("0014d7fae4fbdc8bf6c86a08c7177c5d06683754ea71")
	pk2, _ := hex.DecodeString("0020efd9eeb78068445762a9e2e335d6ab826aaecff9cb7da24a39fc4686d468fb58")
	txOutputs := []bitcoin.TxOut{
		{
			PayScript: pk1,
			Amount:    1199899300,
		},
		{
			PayScript: pk2,
			Amount:    1199899300,
		},
	}

	builder := bitcoin.NewMultiTransactionBuilder()
	err = builder.NetParams(bitcoin.RegTest)
	if err != nil {
		t.Fatal(err)
	}
	err = builder.AddMultiPublicKey(pubkeylist, 2)
	if err != nil {
		t.Fatal(err)
	}
	err = builder.AddTxIn(txInputs)
	if err != nil {
		t.Fatal(err)
	}
	err = builder.AddTxOut(txOutputs)
	if err != nil {
		t.Fatal(err)
	}
	hash := builder.TxHash()
	t.Logf("before: %v \n", hash)
	err = builder.Sign(func(hash []byte) ([][]byte, error) {
		var sigs [][]byte
		for _, privateKey := range privateKeys {
			sig := ecdsa.Sign(privateKey, hash)
			sigWithType := append(sig.Serialize(), byte(txscript.SigHashAll))
			sigs = append(sigs, sigWithType)
		}
		return sigs, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	txBytes, err := builder.Build()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("after: %v \n", hash)
	hexTxBytes := fmt.Sprintf("%x", txBytes)
	txHash, err := client.Sendrawtransaction(hexTxBytes)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("txHash: %v\n", txHash)

}

func TestMockDemo(t *testing.T) {
	for i := 0; i < 3; i++ {
		address, err := client.GetRawChangeAddress()
		if err != nil {
			t.Fatal(err)
		}
		dumpPrivkey, err := client.DumpPrivkey(address)
		if err != nil {
			t.Fatal(err)
		}
		addressInfo, err := client.Getaddressinfo(address)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("address:%v,privatekey:%v,publickey:%v\n", address, dumpPrivkey, addressInfo.Pubkey)
	}
}

func TestSimpleTx(t *testing.T) {
	privateKey := "cPbyxXLqYqAjAHDtbvKq7ETd6BsQBbS643RLHH4u3k1YeVAXkAqR"
	secret, _, err := base58.CheckDecode(privateKey)
	if err != nil {
		t.Fatal(err)
	}
	inputs := []bitcoin.TxIn{
		{
			Hash:     "820556b04b6955cc0646bd424792c5940753a82d55140d9af3cc2f8060576f21",
			VOut:     1,
			PkScript: "0014d7fae4fbdc8bf6c86a08c7177c5d06683754ea71",
			Amount:   1200000000,
		},
	}

	outputs := []bitcoin.TxOut{
		{
			Address: "bcrt1qalv7aduqdpz9wc4fut3nt44tsf42anleed76yj3el3rgd4rgldvq2aw6ze",
			Amount:  1199999800,
		},
	}

	result, err := bitcoin.CreateTransaction(secret, inputs, outputs, bitcoin.RegTest)
	if err != nil {
		t.Fatal(err)
	}
	hexTxData := hex.EncodeToString(result)
	t.Log(hexTxData)
	txHash, err := client.Sendrawtransaction(hexTxData)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(txHash)
}

func TestDepositTransaction(t *testing.T) {

	utxoSet, err := client.Scantxoutset("bcrt1q6lawf77u30mvs6sgcuthchgxdqm4f6n3kvx4z5")
	if err != nil {
		t.Fatal(err)
	}
	if len(utxoSet.Unspents) == 0 {
		t.Fatal("no utxo found")
	}
	utxo := utxoSet.Unspents[0]
	floatBig := big.NewFloat(0).Mul(big.NewFloat(utxo.Amount), big.NewFloat(100000000))
	amount, _ := floatBig.Int64()
	inputs := []bitcoin.TxIn{
		{
			Hash:     utxo.Txid,
			VOut:     uint32(utxo.Vout),
			PkScript: utxo.ScriptPubKey,
			Amount:   amount,
		},
	}
	outputs := []bitcoin.TxOut{
		{
			Address: "bcrt1qalv7aduqdpz9wc4fut3nt44tsf42anleed76yj3el3rgd4rgldvq2aw6ze",
			Amount:  amount - 200,
		},
	}

	//bcrt1q6lawf77u30mvs6sgcuthchgxdqm4f6n3kvx4z5
	privateKey := "cPbyxXLqYqAjAHDtbvKq7ETd6BsQBbS643RLHH4u3k1YeVAXkAqR"
	secret, _, err := base58.CheckDecode(privateKey)
	if err != nil {
		t.Fatal(err)
	}
	ethAddr := "771815eFD58e8D6e66773DB0bc002899c00d5b0c"
	ethAddrBytes, err := hex.DecodeString(ethAddr)
	if err != nil {
		t.Fatal(err)
	}

	result, err := bitcoin.CreateDepositTransaction(secret, ethAddrBytes, inputs, outputs, bitcoin.RegTest)
	if err != nil {
		t.Fatal(err)
	}
	hexTxData := hex.EncodeToString(result)
	t.Log(hexTxData)
	txHash, err := client.Sendrawtransaction(hexTxData)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(txHash)
}

func TestMultiTransaction(t *testing.T) {
	privateKey1 := "cTZYUrhkQmpKGf6QjNhiLrw7g52VL449Tgwzo9UmKbnzphfZUUZp"
	secret1, _, err := base58.CheckDecode(privateKey1)
	if err != nil {
		t.Fatal(err)
	}
	privateKey2 := "cQtt5owpfyv79fGpnv2yncKtQTU42rzohXGCSD3Mcn4CGgEkPVVE"
	secret2, _, err := base58.CheckDecode(privateKey2)
	if err != nil {
		t.Fatal(err)
	}
	privateKey3 := "cSwe7Np3o7eCec6hgKFrwqGs9bb6x2dubKBucLYqQ6mJu5JH1aCn"
	secret3, _, err := base58.CheckDecode(privateKey3)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%x\n%x\n%x\n", secret1, secret2, secret3)
	//t.Logf("%x\n%x\n%x", secret1, secret2, secret3)
	secrets := [][]byte{secret1, secret2, secret3}
	inputs := []bitcoin.TxIn{
		{
			Hash:     "d84498c08fadb2276b1010e8a572e351a4969b8bbb58ac9ff71d8fd3748e2faa",
			VOut:     1,
			PkScript: "0020efd9eeb78068445762a9e2e335d6ab826aaecff9cb7da24a39fc4686d468fb58",
			Amount:   11999998,
		},
		//{
		//	Hash:     "98ea87268149e34d7c818e2cb8fd0b68044beeb0527d0baac83e8619c93450ae",
		//	VOut:     0,
		//	PkScript: "0020efd9eeb78068445762a9e2e335d6ab826aaecff9cb7da24a39fc4686d468fb58",
		//	Amount:   12000000000,
		//},
		//{
		//	Hash:     "ee3cd919cef42d0ca63f62410fccff3a6c1c8754b9ce564dd3c24f551a6911c4",
		//	VOut:     1,
		//	PkScript: "0020efd9eeb78068445762a9e2e335d6ab826aaecff9cb7da24a39fc4686d468fb58",
		//	Amount:   1200000000,
		//},
	}

	outputs := []bitcoin.TxOut{
		{
			Address: "bcrt1q6lawf77u30mvs6sgcuthchgxdqm4f6n3kvx4z5",
			Amount:  1199999300,
		},
	}

	txBytes, err := bitcoin.CreateMultiSigTransaction(2, secrets, inputs, outputs, bitcoin.RegTest)
	if err != nil {
		t.Fatal(err)
	}
	hexTxData := hex.EncodeToString(txBytes)
	t.Log(hexTxData)
	txHash, err := client.Sendrawtransaction(hexTxData)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(txHash)
}

func TestTxHash(t *testing.T) {
	msgTx := wire.NewMsgTx(2)

	hash, _ := chainhash.NewHashFromStr("d84498c08fadb2276b1010e8a572e351a4969b8bbb58ac9ff71d8fd3748e2faa")
	txIn := wire.NewTxIn(wire.NewOutPoint(hash, 1), nil, nil)
	msgTx.AddTxIn(txIn)

	pk1, _ := hex.DecodeString("0014d7fae4fbdc8bf6c86a08c7177c5d06683754ea71")
	pk2, _ := hex.DecodeString("0020efd9eeb78068445762a9e2e335d6ab826aaecff9cb7da24a39fc4686d468fb58")
	out := &wire.TxOut{
		Value:    99950,
		PkScript: pk1,
	}
	msgTx.AddTxOut(out)
	out1 := &wire.TxOut{
		Value:    1199899500,
		PkScript: pk2,
	}
	msgTx.AddTxOut(out1)
	txId := msgTx.TxHash().String()
	t.Log(txId)

}

func TestDeserializeTransaction(t *testing.T) {
	data := "0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000007d0200000001aa2f8e74d38f1df79fac58bb8b9b96a451e372a5e810106b27b2ad8fc09844d80000000000ffffffff026e86010000000000160014d7fae4fbdc8bf6c86a08c7177c5d06683754ea716c03854700000000220020efd9eeb78068445762a9e2e335d6ab826aaecff9cb7da24a39fc4686d468fb5800000000000000"
	hexData, err := hex.DecodeString(data)
	if err != nil {
		t.Fatal(err)
	}
	//version := hexData[0:32]
	l := hexData[32:64]
	length, err := strconv.ParseInt(fmt.Sprintf("%x", l), 16, 64)
	if err != nil {
		t.Fatal(err)
	}
	tx := hexData[64 : 64+length]
	t.Logf("%x\n", tx)
	msgTx := wire.MsgTx{}
	err = msgTx.Deserialize(bytes.NewReader(tx))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(msgTx.TxHash().String())
	t.Log(msgTx)

}
