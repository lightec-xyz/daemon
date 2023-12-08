package store

import (
	"github.com/lightec-xyz/daemon/codec"
	"github.com/lightec-xyz/daemon/logger"
)

var var_ IStore = (*Store)(nil)

type Store struct {
	levelDb *LevelDb
}

func NewStore(file string, cache int, handles int, namespace string, readonly bool) (*Store, error) {
	levelDb, err := NewLevelDb(file, cache, handles, namespace, readonly)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return &Store{
		levelDb: levelDb,
	}, nil
}

func (s *Store) Close() error {
	return s.levelDb.Close()
}

func (s *Store) Has(key []byte) (bool, error) {
	return s.levelDb.Has(key)
}

func (s *Store) Get(key []byte) ([]byte, error) {
	return s.levelDb.Get(key)
}

func (s *Store) Put(key []byte, value []byte) error {
	return s.levelDb.Put(key, value)
}

func (s *Store) Delete(key []byte) error {
	return s.levelDb.Delete(key)
}

func (s *Store) BatchPut(key []byte, value []byte) error {
	return s.levelDb.batch.Put(key, value)
}

func (s *Store) BatchDelete(key []byte) error {
	return s.levelDb.batch.Delete(key)
}

func (s *Store) BatchWrite() error {
	return s.levelDb.batch.Write()
}

func (s *Store) HasObj(key interface{}) (bool, error) {
	keyBytes, err := objKeyEncode(key)
	if err != nil {
		logger.Error("key parse bytes error:%v", err)
		return false, err
	}
	return s.Has(keyBytes)
}

func (s *Store) GetObj(key interface{}, value interface{}) error {
	keyBytes, err := objKeyEncode(key)
	if err != nil {
		logger.Error("key parse bytes error:%v", err)
		return err
	}
	valueBytes, err := s.Get(keyBytes)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return codec.Unmarshal(valueBytes, value)
}

func (s *Store) DeleteObj(key interface{}) error {
	keyBytes, err := objKeyEncode(key)
	if err != nil {
		logger.Error("key parse bytes error:%v", err)
		return err
	}
	return s.Delete(keyBytes)
}
func (s *Store) PutObj(key interface{}, value interface{}) error {
	bytes, err := codec.Marshal(value)
	if err != nil {
		logger.Error("value can't Marshal error:%v", err)
		return err
	}
	keyBytes, err := objKeyEncode(key)
	if err != nil {
		logger.Error("key parse bytes error:%v", err)
		return err
	}
	return s.Put(keyBytes, bytes)
}

func (s *Store) BatchPutObj(key interface{}, value interface{}) error {
	bytes, err := codec.Marshal(value)
	if err != nil {
		logger.Error("value can't Marshal error:%v", err)
		return err
	}
	keyBytes, err := objKeyEncode(key)
	if err != nil {
		logger.Error("key parse bytes error:%v", err)
		return err
	}
	return s.BatchPut(keyBytes, bytes)
}

func (s *Store) BatchDeleteObj(key interface{}) error {
	keyBytes, err := objKeyEncode(key)
	if err != nil {
		logger.Error("key parse bytes error:%v", err)
		return err
	}
	return s.BatchDelete(keyBytes)
}

func (s *Store) BatchWriteObj() error {
	return s.BatchWrite()
}
func objKeyEncode(key interface{}) ([]byte, error) {
	//todo
	switch key.(type) {
	case []byte:
		keyBytes := key.([]byte)
		return keyBytes, nil
	case string:
		keyBytes := []byte(key.(string))
		return keyBytes, nil
	default:
		keyBytes, err := codec.Marshal(key)
		if err != nil {
			return nil, err
		}
		return keyBytes, nil
	}
}
