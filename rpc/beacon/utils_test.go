package beacon

import (
	"github.com/lightec-xyz/daemon/logger"
	apiclient "github.com/lightec-xyz/provers/utils/api-client"
	prysmClient "github.com/prysmaticlabs/prysm/v5/api/client"
	"github.com/prysmaticlabs/prysm/v5/config/params"
	"testing"
)

func TestGetEth1MapToEth2(t *testing.T) {
	//params.UseMainnetNetworkConfig()
	params.OverrideBeaconConfig(params.MainnetConfig().Copy())
	url := "https://rpc.ankr.com/premium-http/eth_beacon/8c933202fbe8dbe6d63377a319b6020f4a4c35bb4424f6368f630b676b4fcc2e"
	beaClient, err := apiclient.NewClient(url, prysmClient.WithAuthenticationToken("8c933202fbe8dbe6d63377a319b6020f4a4c35bb4424f6368f630b676b4fcc2e"))
	if err != nil {
		logger.Error("new provers api client error: %v", err)
		t.Fatal(err)
	}
	result, err := GetEth1MapToEth2(beaClient, 12086551)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}
