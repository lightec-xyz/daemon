package store

import (
	"errors"
	"github.com/lightec-xyz/daemon/codec"
	"github.com/lightec-xyz/daemon/logger"
)

var _ IStore = (*MemoryStore)(nil)

type MemoryStore struct {
	memoryDb *MemoryDb
}

func (m *MemoryStore) WrapBatch(fn func(batch IBatch) error) error {
	b := m.Batch()
	err := fn(b)
	if err != nil {
		return err
	}
	return b.BatchWrite()
}

func (m *MemoryStore) Iter(prefix []byte, start []byte, fn func(key []byte, value []byte) error) error {
	if fn == nil {
		return errors.New("fn is nil")
	}
	iter := m.memoryDb.NewIterator(prefix, start)
	defer iter.Release()
	for iter.Next() {
		err := fn(iter.Key(), iter.Value())
		if err != nil {
			return err
		}
	}
	err := iter.Error()
	if err != nil {
		return err
	}
	return nil
}

func (m *MemoryStore) Batch() IBatch {
	return NewBatch(m.memoryDb.NewBatch())
}

func (m *MemoryStore) Iterator(prefix []byte, start []byte) Iterator {
	return m.memoryDb.NewIterator(prefix, start)
}

func NewMemoryStore() *MemoryStore {
	memoryDb := NewMemoryDb()
	return &MemoryStore{memoryDb: memoryDb}
}

func (m *MemoryStore) Has(key []byte) (bool, error) {
	return m.memoryDb.Has(key)
}

func (m *MemoryStore) Get(key []byte) ([]byte, error) {
	return m.memoryDb.Get(key)
}

func (m *MemoryStore) Put(key []byte, value []byte) error {
	return m.memoryDb.Put(key, value)
}

func (m *MemoryStore) Delete(key []byte) error {
	return m.memoryDb.Delete(key)
}

func (m *MemoryStore) HasObj(key interface{}) (bool, error) {
	keyBytes, err := keyEncode(key)
	if err != nil {
		logger.Error("key parse bytes error:%v", err)
		return false, err
	}
	return m.Has(keyBytes)
}

func (m *MemoryStore) GetObj(key interface{}, value interface{}) error {
	keyBytes, err := keyEncode(key)
	if err != nil {
		logger.Error("key parse bytes error:%v", err)
		return err
	}
	valueBytes, err := m.Get(keyBytes)
	if err != nil {
		return err
	}
	return codec.Unmarshal(valueBytes, value)
}
func (m *MemoryStore) GetValue(key interface{}, value interface{}) (bool, error) {
	keyBytes, err := keyEncode(key)
	if err != nil {
		return false, err
	}
	has, err := m.Has(keyBytes)
	if err != nil {
		return false, err
	}
	if !has {
		return false, nil
	}
	valueBytes, err := m.Get(keyBytes)
	if err != nil {
		return false, err
	}
	return true, codec.Unmarshal(valueBytes, value)
}

func (m *MemoryStore) DeleteObj(key interface{}) error {
	keyBytes, err := keyEncode(key)
	if err != nil {
		logger.Error("key parse bytes error:%v", err)
		return err
	}
	return m.Delete(keyBytes)
}
func (m *MemoryStore) PutObj(key interface{}, value interface{}) error {
	bytes, err := codec.Marshal(value)
	if err != nil {
		logger.Error("value can't Marshal error:%v", err)
		return err
	}
	keyBytes, err := keyEncode(key)
	if err != nil {
		logger.Error("key parse bytes error:%v", err)
		return err
	}
	return m.Put(keyBytes, bytes)
}
