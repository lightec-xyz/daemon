package store

import (
	"github.com/lightec-xyz/daemon/codec"
	"github.com/lightec-xyz/daemon/logger"
)

type MemoryStore struct {
	memoryDb *MemoryDb
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

func (m *MemoryStore) BatchPut(key []byte, value []byte) error {
	return m.memoryDb.batch.Put(key, value)
}

func (m *MemoryStore) BatchDelete(key []byte) error {
	return m.memoryDb.batch.Delete(key)
}

func (m *MemoryStore) BatchWrite() error {
	return m.memoryDb.batch.Write()
}

func (m *MemoryStore) HasObj(key interface{}) (bool, error) {
	bytesKey, err := keyParse(key)
	if err != nil {
		logger.Error("key parse bytes error:%v", err)
		return false, err
	}
	return m.Has(bytesKey)
}

func (m *MemoryStore) GetObj(key interface{}, value interface{}) error {
	bytesKey, err := keyParse(key)
	if err != nil {
		logger.Error("key parse bytes error:%v", err)
		return err
	}
	valueBytes, err := m.Get(bytesKey)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return codec.Unmarshal(valueBytes, value)
}

func (m *MemoryStore) DeleteObj(key interface{}) error {
	bytesKey, err := keyParse(key)
	if err != nil {
		logger.Error("key parse bytes error:%v", err)
		return err
	}
	return m.Delete(bytesKey)
}
func (m *MemoryStore) PutObj(key interface{}, value interface{}) error {
	bytes, err := codec.Marshal(value)
	if err != nil {
		logger.Error("value can't Marshal error:%v", err)
		return err
	}
	bytesKey, err := keyParse(key)
	if err != nil {
		logger.Error("key parse bytes error:%v", err)
		return err
	}
	return m.Put(bytesKey, bytes)
}

func (m *MemoryStore) BatchPutObj(key interface{}, value interface{}) error {
	bytes, err := codec.Marshal(value)
	if err != nil {
		logger.Error("value can't Marshal error:%v", err)
		return err
	}
	bytesKey, err := keyParse(key)
	if err != nil {
		logger.Error("key parse bytes error:%v", err)
		return err
	}
	return m.BatchPut(bytesKey, bytes)
}

func (m *MemoryStore) BatchDeleteObj(key interface{}) error {
	bytesKey, err := keyParse(key)
	if err != nil {
		logger.Error("key parse bytes error:%v", err)
		return err
	}
	return m.BatchDelete(bytesKey)
}
