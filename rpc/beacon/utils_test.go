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
	url := ""
	beaClient, err := apiclient.NewClient(url, prysmClient.WithAuthenticationToken(""))
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
