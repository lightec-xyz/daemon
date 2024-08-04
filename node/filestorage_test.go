package node

import "testing"

func TestFileStorage(t *testing.T) {
	fileStore, err := NewFileStorage("/Users/red/lworkspace/lightec/daemon/node/test", 157*8192, 0)
	if err != nil {
		t.Error(err)
	}
	finalizedSlot, ok, err := fileStore.GetNearTxSlotFinalizedSlot(1315329)
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Fatal("")
	}
	t.Log(finalizedSlot)
}

func TestParseKey(t *testing.T) {

}
