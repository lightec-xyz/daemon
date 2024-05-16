package beacon

import (
	"fmt"
	apiclient "github.com/lightec-xyz/provers/utils/api-client"
	"github.com/prysmaticlabs/prysm/v5/api/client/beacon"
	"github.com/prysmaticlabs/prysm/v5/config/params"
	"github.com/prysmaticlabs/prysm/v5/testing/require"
	"strconv"
	"strings"
	"testing"
)

func TestGetHeadSlot(t *testing.T) {
	cl, err := apiclient.NewClient("http://58.41.9.129:8970")
	require.NoError(t, err)
	params.UseHoleskyNetworkConfig()
	params.OverrideBeaconConfig(params.HoleskyConfig())
	headerResp, err := cl.GetBlockHeader(beacon.IdFinalized)
	require.NoError(t, err)

	require.NotNil(t, headerResp)
	finalizedSlot, err := strconv.Atoi(headerResp.Data.Header.Message.Slot)
	t.Log(finalizedSlot)
}

func TestBeaconSlot(t *testing.T) {
	cl, err := apiclient.NewClient("http://58.41.9.129:8970")
	require.NoError(t, err)

	params.UseHoleskyNetworkConfig()
	params.OverrideBeaconConfig(params.HoleskyConfig())
	headerResp, err := cl.GetBlockHeader(beacon.IdFinalized)
	require.NoError(t, err)

	require.NotNil(t, headerResp)
	finalizedSlot, err := strconv.Atoi(headerResp.Data.Header.Message.Slot)
	t.Log(finalizedSlot)
	require.NoError(t, err)
	for index := 1482752; index < finalizedSlot; index++ {
		eth1MapToEth2, err := GetEth1MapToEth2(cl, uint64(index))
		if err != nil {
			if strings.Contains(err.Error(), "404 NotFound response") {
				t.Logf("slot not found %v \n", index)
				continue
			}
			t.Fatal(err)
		}
		fmt.Printf("slot:%v,number:%v\n", eth1MapToEth2.BlockSlot, eth1MapToEth2.BlockNumber)
	}
}
