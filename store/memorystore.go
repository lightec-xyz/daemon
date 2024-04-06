package store

import (
	"github.com/lightec-xyz/daemon/codec"
	"github.com/lightec-xyz/daemon/logger"
)

var _ IStore = (*MemoryStore)(nil)

type MemoryStore struct {
	memoryDb *MemoryDb
}

func (m *MemoryStore) Batch() IBatch {
	//TODO implement me
	panic("implement me")
}

func (m *MemoryStore) Iterator(prefix []byte, start []byte) Iterator {
	//TODO implement me
	panic("implement me")
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
	keyBytes, err := KeyEncode(key)
	if err != nil {
		logger.Error("key parse bytes error:%v", err)
		return false, err
	}
	return m.Has(keyBytes)
}

func (m *MemoryStore) GetObj(key interface{}, value interface{}) error {
	keyBytes, err := KeyEncode(key)
	if err != nil {
		logger.Error("key parse bytes error:%v", err)
		return err
	}
	valueBytes, err := m.Get(keyBytes)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return codec.Unmarshal(valueBytes, value)
}

func (m *MemoryStore) DeleteObj(key interface{}) error {
	keyBytes, err := KeyEncode(key)
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
	keyBytes, err := KeyEncode(key)
	if err != nil {
		logger.Error("key parse bytes error:%v", err)
		return err
	}
	return m.Put(keyBytes, bytes)
}
