package node

import "testing"

func TestFileStorage(t *testing.T) {
	fileStore, err := NewFileStorage("test", 90*8192)
	if err != nil {
		t.Error(err)
	}
	err = fileStore.StorePeriod(100)
	if err != nil {
		t.Error(err)
	}
	indexes, err := fileStore.NeedUpdateIndexes()
	if err != nil {
		t.Error(err)
	}
	t.Log(indexes)
}
