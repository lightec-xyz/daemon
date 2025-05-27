package node

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc/ethereum"
	apiclient "github.com/lightec-xyz/provers/utils/api-client"
	txineth2Utils "github.com/lightec-xyz/provers/utils/tx-in-eth2"
	prysmClient "github.com/prysmaticlabs/prysm/v5/api/client"
	"testing"
)

func TestPrepared_GetTxInEth2Request(t *testing.T) {
	ethClient, err := ethereum.NewClient("http://127.0.0.1:8002",
		"0xB86E9A8391d3df83F53D3f39E3b5Fce4D7da405d",
		"0x2635Dc72706478F4bD784A8D04B3e0af8AB053dc",
		"0xB86E9A8391d3df83F53D3f39E3b5Fce4D7da405d",
		"0x199CC8f0ac008Bdc8cF0B1CCd5187F84E168C4D2")
	if err != nil {
		t.Fatal(err)
	}
	url := "http://127.0.0.1:8003"
	apiClient, err := apiclient.NewClient(url, prysmClient.WithAuthenticationToken(getUrlToken(url)))
	if err != nil {
		logger.Error("new provers api client error: %v", err)
		t.Fatal(err)
		//params.UseSepoliaNetworkConfig()
		//params.OverrideBeaconConfig(params.SepoliaConfig())
		txData, err := txineth2Utils.GetTxInEth2ProofData(ethClient.Client, apiClient, func(blockNumber uint64) (uint64, error) {
			return 12089032, nil
		}, ethcommon.HexToHash("08be2acd2b07cda893c33be6790939500c49f61ace481a3507a272de383e5927"))
		if err != nil {
			t.Fatal(err)
		}
		t.Log(txData)
	}
}
