package node

import "testing"

func TestBtcToSat(t *testing.T) {
	sat := BtcToSat(21.00000011645645645)
	t.Log(sat)

}
