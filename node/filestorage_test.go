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
	key := parseKey(1)
	t.Log(key)
	k1 := parseKey(100, 200)
	t.Log(k1)
	k2 := parseKey(100, 200, 300)
	t.Log(k2)
}
