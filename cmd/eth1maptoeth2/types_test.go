package main

import (
	"fmt"
	"strconv"
	"testing"

	apiclient "github.com/lightec-xyz/provers/utils/api-client"
	"github.com/prysmaticlabs/prysm/v5/api/client"
	"github.com/prysmaticlabs/prysm/v5/api/client/beacon"
	"github.com/prysmaticlabs/prysm/v5/config/params"
	"github.com/stretchr/testify/require"
)

func Test_GetHeadSlot(t *testing.T) {
	cl, err := apiclient.NewClient("http://58.41.9.129:8970")
	require.NoError(t, err)

	headSlot, err := GetHeadSlot(cl)
	require.NoError(t, err)

	t.Log(headSlot)
}

func Test_GetEth1MapToEth2(t *testing.T) {
	tokenOpt := client.WithAuthenticationToken("3ac3d8d70361a628192b6fd7cd71b88a0b17638d")

	cl, err := apiclient.NewClient("https://young-morning-meadow.ethereum-holesky.quiknode.pro", tokenOpt)
	require.NoError(t, err)

	params.UseHoleskyNetworkConfig()
	params.OverrideBeaconConfig(params.HoleskyConfig())

	headerResp, err := cl.GetBlockHeader(beacon.IdFinalized)
	require.NoError(t, err)

	require.NotNil(t, headerResp)
	finalizedSlot, err := strconv.Atoi(headerResp.Data.Header.Message.Slot)
	require.NoError(t, err)

	eth1MapToEth2, err := GetEth1MapToEth2(cl, finalizedSlot)
	require.NoError(t, err)

	fmt.Printf("eth1MapToEth2: %+v\n", eth1MapToEth2)
}
