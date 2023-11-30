package bitcoin

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/lightec-xyz/daemon/transaction/bitcoin"
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
	header, err := client.GetBlockHeader("")
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
	privateKey := "cPbyxXLqYqAjAHDtbvKq7ETd6BsQBbS643RLHH4u3k1YeVAXkAqR"
	secret, _, err := base58.CheckDecode(privateKey)
	if err != nil {
		t.Fatal(err)
	}
	ethAddr := "e8c84a631D71E1Bb7083D3a82a3a74870a286B97"
	ethAddrBytes, err := hex.DecodeString(ethAddr)
	if err != nil {
		t.Fatal(err)
	}
	inputs := []bitcoin.TxIn{
		{
			Hash:     "7d8f46b43caebfc8f5940b3bbab189aa96d6569580e7328f19d5542de2a51467",
			VOut:     0,
			PkScript: "0014d7fae4fbdc8bf6c86a08c7177c5d06683754ea71",
			Amount:   1200000000,
		},
	}

	outputs := []bitcoin.TxOut{
		{
			Address: "bcrt1qalv7aduqdpz9wc4fut3nt44tsf42anleed76yj3el3rgd4rgldvq2aw6ze",
			Amount:  1199999500,
		},
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
	//t.Logf("%x\n%x\n%x", secret1, secret2, secret3)
	secrets := [][]byte{secret1, secret2, secret3}
	inputs := []bitcoin.TxIn{
		{
			Hash:     "99bc9c085120370e3f89aeeb2dcb11657f36d4411f832dd6d71cdaad65e2512b",
			VOut:     0,
			PkScript: "0020efd9eeb78068445762a9e2e335d6ab826aaecff9cb7da24a39fc4686d468fb58",
			Amount:   1200000000,
		},
		{
			Hash:     "98ea87268149e34d7c818e2cb8fd0b68044beeb0527d0baac83e8619c93450ae",
			VOut:     0,
			PkScript: "0020efd9eeb78068445762a9e2e335d6ab826aaecff9cb7da24a39fc4686d468fb58",
			Amount:   12000000000,
		},
		{
			Hash:     "ee3cd919cef42d0ca63f62410fccff3a6c1c8754b9ce564dd3c24f551a6911c4",
			VOut:     1,
			PkScript: "0020efd9eeb78068445762a9e2e335d6ab826aaecff9cb7da24a39fc4686d468fb58",
			Amount:   1200000000,
		},
	}

	outputs := []bitcoin.TxOut{
		{
			Address: "bcrt1q6lawf77u30mvs6sgcuthchgxdqm4f6n3kvx4z5",
			Amount:  14399999600,
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

func TestOpReturnTransaction(t *testing.T) {

}
