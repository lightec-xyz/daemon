package bitcoin

import (
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/btcutil/base58"
)

func TestClient_GetaddressInfo(t *testing.T) {
	getaddressinfo, err := client.Getaddressinfo("bcrt1qvx3puzwjlr3ttqay5g3pfaj973trgzryl0mce0")
	if err != nil {
		panic(err)
	}
	t.Log(getaddressinfo.Pubkey)
	t.Log(getaddressinfo.ScriptPubKey)
}

func TestClient_DumpPrivteKey(t *testing.T) {
	privkey, err := client.DumpPrivkey("bcrt1qvx3puzwjlr3ttqay5g3pfaj973trgzryl0mce0")
	if err != nil {
		panic(err)
	}
	t.Log(privkey)
	t.Log(fmt.Sprintf("%x", base58.Decode(privkey)))
}

func TestGetrawChangeAddress(t *testing.T) {
	address, err := client.GetRawChangeAddress()
	if err != nil {
		panic(err)
	}
	t.Log(address)
}

func TestGenTestAddress(t *testing.T) {
	for index := 0; index < 3; index++ {
		address, err := client.GetRawChangeAddress()
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("address: %v\n", address)
		privkey, err := client.DumpPrivkey(address)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("privateKey: %v | %x \n", privkey, base58.Decode(privkey))
		addressInfo, err := client.Getaddressinfo(address)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("publicKey: %v \n", addressInfo.Pubkey)
		t.Logf("scriptPubKey: %v \n", addressInfo.ScriptPubKey)
		t.Log("------------------------------------")
	}

}
