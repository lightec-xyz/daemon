package ethereum

import "testing"

var err error
var client *Client

func init() {
	client, err = NewClient("https://endpoints.omniatech.io/v1/eth/goerli/public")
	if err != nil {
		panic(err)
	}
}

func TestClient_TestEth(t *testing.T) {
	result, err := client.EthRPC.EthGetBlockByNumber(10127442, false)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}
