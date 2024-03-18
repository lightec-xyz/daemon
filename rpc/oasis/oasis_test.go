package oasis

import "testing"

var err error
var client *Client

func init() {
	client, err = NewClient("https://testnet.sapphire.oasis.dev/",
		[]string{
			"0x8Cf463B54C5481E772841870E01c8c8d2671e66f",
			"0xB806e81B33FDD9b457aE88DEEA25258C688Ee470",
			"0x5ADDC1A7E0bd5b05Bb1fd454a57D494Cfb62443F",
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
		t.Logf("%x\n", item)
	}
}
