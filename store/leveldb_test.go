package store

import (
	"github.com/ethereum/go-ethereum/ethdb/leveldb"
	"log"
	"testing"
)

func TestLevelDb(t *testing.T) {
	// 指定数据库路径
	db, err := leveldb.New("./testdb", 0, 0, "zkbtc", false)
	if err != nil {
		log.Fatal(err)
	}
	t.Log(db)
}
