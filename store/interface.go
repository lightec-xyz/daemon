package store

type ILevelDb interface {
	Has(key []byte) (bool, error)
	Get(key []byte) ([]byte, error)
	Put(key []byte, value []byte) error
	Delete(key []byte) error
}

type IBatch interface {
	BatchPut(key []byte, value []byte) error
	BatchDelete(key []byte) error
	BatchWrite() error
}

type IStoreObj interface {
	HasObj(key interface{}) (bool, error)
	GetObj(key interface{}, value interface{}) error
	DeleteObj(key interface{}) error
	PutObj(key interface{}, value interface{}) error
	BatchPutObj(key interface{}, value interface{}) error
	BatchDeleteObj(key interface{}) error
	BatchWriteObj() error
}

type IStore interface {
	ILevelDb
	IBatch
	IStoreObj
}
