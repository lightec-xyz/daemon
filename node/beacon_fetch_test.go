package node

import (
	"github.com/lightec-xyz/daemon/rpc/beacon"
	"testing"
	"time"
)

func TestBeaconFetch_Fetch(t *testing.T) {
	genesisPeriod := uint64(60)
	client, err := beacon.NewClient("http://127.0.0.1:8970")
	if err != nil {
		t.Fatal(err)
	}
	fileStore, err := NewFileStore("test")
	if err != nil {
		t.Fatal(err)
	}
	fetchResp := make(chan FetchDataResponse, 1)
	beaconFetch, err := NewBeaconFetch(client, fileStore, fetchResp)
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		beaconFetch.Fetch()
	}()
	go func() {
		for {
			select {
			case resp := <-fetchResp:
				t.Logf("receive resp: %v %v \n", resp.period, resp.reqType)
			}
		}
	}()
	go func() {
		beaconFetch.GenesisUpdateRequest()
		for {
			if beaconFetch.canNewRequest() {
				beaconFetch.NewUpdateRequest(genesisPeriod)
				genesisPeriod = genesisPeriod + 1
			} else {
				time.Sleep(10 * time.Second)
			}
		}
	}()

	ch := make(chan int)
	<-ch
}
