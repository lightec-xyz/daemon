package store

import (
	"fmt"
	"github.com/lightec-xyz/daemon/codec"
	"github.com/lightec-xyz/daemon/logger"
)

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

func (s *Store) HasObj(key interface{}) (bool, error) {
	bytesKey, err := KeyParse(key)
	if err != nil {
		logger.Error("key parse bytes error:%v", err)
		return false, err
	}
	return s.Has(bytesKey), nil
}

func (s *Store) GetObj(key interface{}, value interface{}) error {
	bytesKey, err := KeyParse(key)
	if err != nil {
		logger.Error("key parse bytes error:%v", err)
		return err
	}
	valueBytes, err := s.Get(bytesKey)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return codec.Unmarshal(valueBytes, value)
}

func (s *Store) DeleteObj(key interface{}) error {
	bytesKey, err := KeyParse(key)
	if err != nil {
		logger.Error("key parse bytes error:%v", err)
		return err
	}
	return s.Delete(bytesKey)
}
func (s *Store) PutObj(key interface{}, value interface{}) error {
	bytes, err := codec.Marshal(value)
	if err != nil {
		logger.Error("value can't Marshal error:%v", err)
		return err
	}
	bytesKey, err := KeyParse(key)
	if err != nil {
		logger.Error("key parse bytes error:%v", err)
		return err
	}
	return s.Put(bytesKey, bytes)
}

func KeyParse(key interface{}) ([]byte, error) {
	//todo
	switch key.(type) {
	case string:
		return []byte(key.(string)), nil
	}
	return nil, fmt.Errorf("unsupported key type:%v", key)
}
