package node

import (
	"github.com/lightec-xyz/daemon/rpc/beacon"
	"testing"
)

func TestBeaconFetch_Fetch(t *testing.T) {
	beacon.NewClient("")
	beaconFetch, err := NewBeaconFetch(nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		beaconFetch.Fetch()
	}()
	ch := make(chan int)
	<-ch
}
