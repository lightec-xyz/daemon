package cmd

import (
	"github.com/lightec-xyz/daemon/circuits"
	"github.com/lightec-xyz/daemon/node"
	"testing"
)

func TestBhfProve(t *testing.T) {
	fileStore, err := node.NewFileStore("/Users/red/lworkspace/lightec/daemon/node/test", 176*8192)
	if err != nil {
		t.Fatal(err)
	}
	syncCommitUpdate, ok, err := node.GetSyncCommitUpdate(fileStore, 177)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("no find sync commit update")
	}
	rootId, err := circuits.SyncCommitRoot(syncCommitUpdate)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%x", rootId)
}
