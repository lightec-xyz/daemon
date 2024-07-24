package dfinity

import "testing"

var client *Client
var err error

func init() {
	client, err = NewClient()
	if err != nil {
		panic(err)
	}
}

func TestClient_VerifyPublicKey(t *testing.T) {
	publicKey, err := client.PublicKey("lgj7u-yaaaa-aaaak-qipsa-cai")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(publicKey)
}

func TestClient_BtcBalance(t *testing.T) {
	balance, err := client.BtcBalance("lgj7u-yaaaa-aaaak-qipsa-cai", "tb1qn9fpljh5ggp407z02jx8x76pemzclgd6rla0qp")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(balance)
}

func TestClient_Sign(t *testing.T) {
	resp, err := client.Sign("lgj7u-yaaaa-aaaak-qipsa-cai", "hello")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}

func TestClient_Verify(t *testing.T) {
	resp, err := client.Verify("lgj7u-yaaaa-aaaak-qipsa-cai", "hello", "hello", "tb1qn9fpljh5ggp407z02jx8x76pemzclgd6rla0qp")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}

func TestClient_BtcUtxo(t *testing.T) {
	resp, err := client.BtcUtxo("lgj7u-yaaaa-aaaak-qipsa-cai", "tb1qn9fpljh5ggp407z02jx8x76pemzclgd6rla0qp")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}

func TestClient_CanisterInfo(t *testing.T) {

}
