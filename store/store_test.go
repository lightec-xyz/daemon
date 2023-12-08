package store

import (
	"github.com/lightec-xyz/daemon/logger"
	"testing"
)

func init() {
	logger.InitLogger()
}
func TestStore_Demo(t *testing.T) {
	storeDb, err := NewStore("~/.daemon/testnet", 0, 0, "zkbtc", false)
	if err != nil {
		t.Fatal(err)
	}
	has, err := storeDb.Has([]byte("ethCurHeight"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(has)
	var result interface{}
	err = storeDb.GetObj([]byte("ethCurHeight"), &result)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}
