package node

import (
	"testing"
)

func TestStoreProofToFile(t *testing.T) {
	fileStorage, err := NewFileStorage("/Users/red/lworkspace/lightec/audit/daemon/testdata/.daemon", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	err = fileStorage.RemoveBtcProof(83000)
	if err != nil {
		t.Fatal(err)
	}
}
