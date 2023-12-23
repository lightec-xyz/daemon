package store

import (
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
)

type MemoryDb struct {
	*memorydb.Database
}

func NewMemoryDb() *MemoryDb {
	database := memorydb.New()
	return &MemoryDb{Database: database}
}
