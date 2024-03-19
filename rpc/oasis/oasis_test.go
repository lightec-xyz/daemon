package oasis

import (
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

var err error
var client *Client

func init() {
	client, err = NewClient("https://testnet.sapphire.oasis.dev/",
		[]string{
			"0xA3D5838913497AD0fcdE036F128a446289EBaD03",
			"0x907e95F678e4D746D5b532B33D6dC17705a71aB6",
			"0x8DDa72eE36AB9c91e92298823D3C0d4D73894081",
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
