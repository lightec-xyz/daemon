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
