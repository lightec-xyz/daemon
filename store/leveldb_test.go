package store

import (
	"log"
	"testing"
)

func TestLevelDb(t *testing.T) {
	// 指定数据库路径
	db, err := NewLevelDb("/Users/red/.daemon/testnet/data", 0, 0, "zkbtc", false)
	if err != nil {
		log.Fatal(err)
	}
	t.Log(db)
}
