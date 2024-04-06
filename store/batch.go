package store

import (
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/lightec-xyz/daemon/codec"
)

var _ IBatch = (*Batch)(nil)

type Batch struct {
	batch ethdb.Batch
}

func NewBatch(db ethdb.Batch) *Batch {
	return &Batch{
		batch: db,
	}
}

func (b *Batch) BatchPut(key []byte, value []byte) error {
	return b.batch.Put(key, value)
}

func (b *Batch) BatchDelete(key []byte) error {
	return b.batch.Delete(key)
}

func (b *Batch) BatchWrite() error {
	return b.batch.Write()
}

func (b *Batch) BatchPutObj(key interface{}, value interface{}) error {
	keyEncode, err := KeyEncode(key)
	if err != nil {
		return err
	}
	valueBytes, err := codec.Marshal(value)
	if err != nil {
		return err
	}
	return b.batch.Put(keyEncode, valueBytes)
}

func (b *Batch) BatchDeleteObj(key interface{}) error {
	keyEncode, err := KeyEncode(key)
	if err != nil {
		return err
	}
	return b.batch.Delete(keyEncode)
}

func (b *Batch) BatchWriteObj() error {
	return b.batch.Write()
}
