package node

import "testing"

func TestFileStorage(t *testing.T) {
	fileStore, err := NewFileStorage("/Users/red/lworkspace/lightec/daemon/node/test", 157*8192)
	if err != nil {
		t.Error(err)
	}
	finalizedSlot, err := fileStore.GetNearTxSlotFinalizedSlot(130)
	if err != nil {
		t.Error(err)
	}
	t.Log(finalizedSlot)
}
