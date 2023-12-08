package store

import (
	"github.com/lightec-xyz/daemon/logger"
	"testing"
)

func init() {
	logger.InitLogger()
}
func TestStore_Demo(t *testing.T) {
	store, err := NewStore("~/.daemon/testnet", 0, 0, "zkbtc", false)
	if err != nil {
		t.Fatal(err)
	}
	var result interface{}
	err = store.GetObj([]byte("ethCurHeight"), &result)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}
