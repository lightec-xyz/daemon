package node

import (
	btcproverClient "github.com/lightec-xyz/btc_provers/utils/client"
	"github.com/ybbus/jsonrpc/v3"
	"net/http"
	"testing"
	"time"
)

func TestNewBtcClient(t *testing.T) {
	proverClient := btcproverClient.NewJsonRpcClient("http://127.0.0.1:9935", "", "", &jsonrpc.RPCClientOpts{
		HTTPClient: &http.Client{
			Timeout: 1 * time.Minute,
		},
	})
	block, err := proverClient.GetBlock("0000000000000000a524ba8922066e55acd50836e909072508acb999c1c5ce72")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(block)

}
