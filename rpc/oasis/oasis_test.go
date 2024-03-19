package oasis

import (
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

var err error
var client *Client

func init() {
	client, err = NewClient("https://testnet.sapphire.oasis.io",
		[]string{
			"0xBb8b61bD363221281A105b6a37ad4CF7DDf24BAc",
			"0x5ee2C3FABED0780abB5905fCD6DEbf1C3C42C729",
			"0x7e0d35F36a1103Fe0Ad91911b2798Cb24A6beC7f",
		},
	)
	if err != nil {
		panic(err)
	}
}

func TestClient_PublicKey(t *testing.T) {
	publicKey, err := client.PublicKey()
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range publicKey {
		t.Logf("%v\n", hexutil.Encode(item))
	}
}

func TestClient_SignBtcTx(t *testing.T) {
	txRaw, err := hexutil.Decode("0xBb8b61bD363221281A105b6a37ad4CF7DDf24BAc")
	if err != nil {
		t.Fatal(err)
	}
	receiptRaw, err := hexutil.Decode("0xBb8b61bD363221281A105b6a37ad4CF7DDf24BAc")
	if err != nil {
		t.Fatal(err)
	}
	proofData, err := hexutil.Decode("0xBb8b61bD363221281A105b6a37ad4CF7DDf24BAc")
	if err != nil {
		t.Fatal(err)
	}
	sig1, sig2, sig3, err := client.SignBtcTx(txRaw, receiptRaw, proofData)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(sig1)
	t.Log(sig2)
	t.Log(sig3)
}
