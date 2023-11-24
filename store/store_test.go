package store

import (
	"github.com/lightec-xyz/daemon/logger"
	"testing"
)

func init() {
	logger.InitLogger()
}
func TestStore_Demo(t *testing.T) {
	store, err := NewStore("/Users/red/.daemon/testnet/data", 0, 0, "zkbtc", false)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(store)
	err = store.PutObj(1000, 10000)
	if err != nil {
		t.Fatal(err)
	}
	var value int64
	err = store.GetObj(1000, &value)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(value)
}
