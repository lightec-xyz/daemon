package store

type IStore interface {
	Has(key []byte) (bool, error)
	Get(key []byte) ([]byte, error)
	Put(key []byte, value []byte) error
	Delete(key []byte) error
	BatchPut(key []byte, value []byte) error
	BatchDelete(key []byte) error
	BatchWrite() error
	HasObj(key interface{}) (bool, error)
	GetObj(key interface{}, value interface{}) error
	DeleteObj(key interface{}) error
	PutObj(key interface{}, value interface{}) error
	BatchPutObj(key interface{}, value interface{}) error
	BatchDeleteObj(key interface{}) error
	BatchWriteObj() error
}
