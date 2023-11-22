package store

import (
	"github.com/ethereum/go-ethereum/ethdb/leveldb"
	"github.com/lightec-xyz/daemon/logger"
)

type LevelDb struct {
	*leveldb.Database
}

func NewLevelDb(file string, cache int, handles int, namespace string, readonly bool) (*LevelDb, error) {
	db, err := leveldb.New(file, cache, handles, namespace, readonly)
	if err != nil {
		logger.Error("new leveldb error:%v", err)
		return nil, err
	}
	return &LevelDb{db}, nil
}
