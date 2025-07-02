package node

import "testing"

func TestBtcToSat(t *testing.T) {
	sat := BtcToSat(0.00000011)
	t.Log(sat)
	
}
