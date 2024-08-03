package node

import (
	btcproverClient "github.com/lightec-xyz/btc_provers/utils/client"
	"github.com/lightec-xyz/daemon/rpc/beacon"
	btcrpc "github.com/lightec-xyz/daemon/rpc/bitcoin"
	ethrpc "github.com/lightec-xyz/daemon/rpc/ethereum"
	"github.com/lightec-xyz/daemon/store"
	apiclient "github.com/lightec-xyz/provers/utils/api-client"
	"testing"
)

func TestState(t *testing.T) {

	state, err := initState()
	if err != nil {
		t.Fatal(err)
	}
	go state.CheckBtcState()
	ch := make(chan struct{}, 1)
	<-ch

}

func initState() (*State, error) {
	genesisSlot := uint64(0)
	genesisPeriod := genesisSlot / 8192
	fileStorage, err := NewFileStorage("/Users/red/lworkspace/lightec/daemon/daemon/node/test/testnet", genesisSlot, 0)
	if err != nil {
		panic(err)
	}
	store, err := store.NewStore("/Users/red/lworkspace/lightec/daemon/daemon/node/test/testnet", 0, 0, "zkbtc", false)
	if err != nil {
		panic(err)
	}
	client, err := btcproverClient.NewClient("http://18.116.118.39:18332", "Lightec", "Abcd1234")
	if err != nil {
		panic(err)
	}
	btcClient, err := btcrpc.NewClient("http://18.116.118.39:18332", "Lightec", "Abcd1234")
	if err != nil {
		panic(err)
	}
	ethClient, err := ethrpc.NewClient("http://3.15.40.243:8545", "", "", "")
	if err != nil {
		panic(err)
	}
	apiClient, err := apiclient.NewClient("http://58.41.9.129:8970")
	if err != nil {
		panic(err)
	}
	beacon, err := beacon.NewClient("http://58.41.9.129:8970")
	if err != nil {
		panic(err)
	}
	preparedData, err := NewPreparedData(fileStorage, store, genesisPeriod, client, btcClient, ethClient, apiClient, beacon)
	if err != nil {
		panic(err)
	}
	state, err := NewState(NewArrayQueue(), fileStorage, store, NewCacheState(), preparedData,
		0, 0, nil, nil)
	if err != nil {
		return nil, err
	}
	return state, nil
}
