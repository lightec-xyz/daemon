package store

import (
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
)

type MemoryDb struct {
	*memorydb.Database
	batch ethdb.Batch
}

func NewMemoryDb() *MemoryDb {
	database := memorydb.New()
	batch := database.NewBatch()
	return &MemoryDb{Database: database, batch: batch}
}
