package node

import (
	"bytes"
	"encoding/json"
	blockdepthUtil "github.com/lightec-xyz/btc_provers/utils/blockdepth"
	btcproverClient "github.com/lightec-xyz/btc_provers/utils/client"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/ybbus/jsonrpc/v3"
	"net/http"
	"testing"
	"time"
)

func TestBlockDepthIndex(t *testing.T) {
	start := uint64(common.BtcCpMinDepth)
	end := start + 3*common.BtcCpMinDepth + 32 - 1
	t.Log(start, end)
	indexes := BlockDepthPlan(60003, start, end, false)
	for _, index := range indexes {
		t.Logf("step%v: %v_%v\n", index.Step, index.Start, index.End)
	}
}

func TestBlockChainPlan(t *testing.T) {
	start := uint64(common.BtcUpperDistance)
	end := start + 4*common.BtcUpperDistance - 12 //19,14,9
	t.Log(start, end)
	indexes := BlockChainPlan(start, end)
	for _, index := range indexes {
		t.Logf("step%v: %v_%v\n", index.Step, index.Start, index.End)
	}
}

func TestData(t *testing.T) {
	proverClient := btcproverClient.NewJsonRpcClient("http://127.0.0.1:9935", "", "", &jsonrpc.RPCClientOpts{
		HTTPClient: &http.Client{
			Timeout: 1 * time.Minute,
		},
	})
	cpData, err := blockdepthUtil.GetCpTimestampProofData(proverClient, uint32(79802))
	if err != nil {
		t.Fatal(err)
	}
	request := common.ProofRequest{Data: &rpc.BtcTimestampRequest{
		CpTime: cpData,
	}}

	data, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))

	var newReq common.ProofRequest
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	err = decoder.Decode(&newReq)
	if err != nil {
		t.Fatal(err)
	}

	if v, ok := newReq.Data.(rpc.ICheck); ok {
		err := v.Check()
		if err != nil {
			t.Fatal(err)
		}
	}
	t.Log(newReq)

}
