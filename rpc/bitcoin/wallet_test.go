package bitcoin

import (
	"fmt"
	"testing"
)

func TestClient_GetaddressInfo(t *testing.T) {
	getaddressinfo, err := client.Getaddressinfo("bcrt1q6lawf77u30mvs6sgcuthchgxdqm4f6n3kvx4z5")
	if err != nil {
		panic(err)
	}
	fmt.Println(getaddressinfo.Pubkey)
}

func TestClient_DumpPrivteKey(t *testing.T) {
	privkey, err := client.DumpPrivkey("bcrt1q2v3adhw34kc2am22w6rw88mryufmv9dtg5rwd2")
	if err != nil {
		panic(err)
	}
	fmt.Println(privkey)
}

func TestGetrawChangeAddress(t *testing.T) {
	address, err := client.GetRawChangeAddress()
	if err != nil {
		panic(err)
	}
	fmt.Print(address)
}
