package node

import (
	"github.com/lightec-xyz/daemon/rpc/beacon"
	"testing"
	"time"
)

func TestBeaconFetch_Fetch(t *testing.T) {
	genesisPeriod := uint64(50)
	client, err := beacon.NewClient("http://127.0.0.1:8970")
	if err != nil {
		t.Fatal(err)
	}
	fileStore, err := NewFileStore("test", genesisPeriod)
	if err != nil {
		t.Fatal(err)
	}
	fetchResp := make(chan FetchDataResponse, 1)
	beaconFetch, err := NewBeaconFetch(client, fileStore, genesisPeriod, fetchResp)
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
				t.Logf("receive resp: %v %v \n", resp.period, resp.UpdateType)
			}
		}
	}()
	go func() {
		beaconFetch.BootStrapRequest()
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
