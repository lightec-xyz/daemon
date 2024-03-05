package node

import "testing"

func TestBeaconFetch_Fetch(t *testing.T) {
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
