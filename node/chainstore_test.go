package node

import (
	"github.com/lightec-xyz/daemon/store"
	"testing"
)

func initStore() store.IStore {
	dbPath := "/Users/red/lworkspace/lightec/daemon/daemon/node/test/dbtest"
	db, err := store.NewStore(dbPath, 0, 0, "zkbtc", false)
	if err != nil {
		panic(err)
	}
	return db
}

func TestDemo0001(t *testing.T) {
	id := dbTxSlotId(1, "dsdsfs")
	t.Log(string(id))
}
