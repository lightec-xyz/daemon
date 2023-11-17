package store

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/lightec-xyz/daemon/codec"
	"github.com/lightec-xyz/daemon/logger"
	"strings"
)

const ProtocolSeparator = "_"

var var_ IStore = (*Store)(nil)

type Store struct {
	levelDb *LevelDb
}

func NewStore(file string, cache int, handles int, namespace string, readonly bool) (*Store, error) {
	levelDb, err := NewLevelDb(file, cache, handles, namespace, readonly)
	if err != nil {
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

func (s *Store) Batch() IBatch {
	return NewBatch(s.levelDb.NewBatch())
}

func (s *Store) WrapBatch(fn func(batch IBatch) error) error {
	b := s.Batch()
	err := fn(b)
	if err != nil {
		return err
	}
	return b.BatchWrite()
}

type Key string

func (k Key) String() string {
	return string(k)
}

func (s *Store) HasObj(key interface{}) (bool, error) {
	keyBytes, err := keyEncode(key)
	if err != nil {
		logger.Error("key parse bytes error:%v", err)
		return false, err
	}
	return s.Has(keyBytes)
}

func (s *Store) GetObj(key interface{}, value interface{}) error {
	keyBytes, err := keyEncode(key)
	if err != nil {
		logger.Error("key parse bytes error:%v", err)
		return err
	}
	valueBytes, err := s.Get(keyBytes)
	if err != nil {
		return err
	}
	return codec.Unmarshal(valueBytes, value)
}

func (s *Store) GetValue(key interface{}, value interface{}) (bool, error) {
	keyBytes, err := keyEncode(key)
	if err != nil {
		return false, err
	}
	has, err := s.Has(keyBytes)
	if err != nil {
		return false, err
	}
	if !has {
		return false, nil
	}
	valueBytes, err := s.Get(keyBytes)
	if err != nil {
		return false, err
	}
	return true, codec.Unmarshal(valueBytes, value)
}

func (s *Store) DeleteObj(key interface{}) error {
	keyBytes, err := keyEncode(key)
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
	keyBytes, err := keyEncode(key)
	if err != nil {
		logger.Error("key parse bytes error:%v", err)
		return err
	}
	return s.Put(keyBytes, bytes)
}

func (s *Store) Iter(prefix, start []byte, fn func(key, value []byte) error) error {
	if fn == nil {
		return errors.New("fn is nil")
	}
	iter := s.levelDb.NewIterator(prefix, start)
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

func (s *Store) Iterator(prefix []byte, start []byte) Iterator {
	return s.levelDb.NewIterator(prefix, start)
}

func keyEncode(key interface{}) ([]byte, error) {
	//todo
	switch key.(type) {
	case []byte:
		keyBytes := key.([]byte)
		return keyBytes, nil
	case string:
		keyBytes := []byte(key.(string))
		return keyBytes, nil
	case uint64:
		keyBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(keyBytes, key.(uint64))
		return keyBytes, nil
	case int64:
		keyBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(keyBytes, key.(uint64))
		return keyBytes, nil
	case int:
		keyBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(keyBytes, key.(uint64))
		return keyBytes, nil
	default:
		keyBytes, err := codec.Marshal(key)
		if err != nil {
			return nil, err
		}
		return keyBytes, nil
	}
}

func GenDbKey(args ...interface{}) Key {
	var key string
	for _, arg := range args {
		key = key + fmt.Sprintf("_%v", arg)
	}
	return Key(strings.ToLower(strings.TrimPrefix(key, "_")))
}
