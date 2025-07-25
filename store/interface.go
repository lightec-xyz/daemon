package store

type ILevelDb interface {
	Has(key []byte) (bool, error)
	Get(key []byte) ([]byte, error)
	Put(key []byte, value []byte) error
	Delete(key []byte) error
	Iterator(prefix []byte, start []byte) Iterator // todo
	Iter(prefix []byte, start []byte, fn func(key, value []byte) error) error
	Compact(start, limit []byte) error
	Batch() IBatch
}

type IBatch interface {
	BatchPut(key []byte, value []byte) error
	BatchDelete(key []byte) error
	BatchWrite() error
	BatchPutObj(key interface{}, value interface{}) error
	BatchDeleteObj(key interface{}) error
	BatchWriteObj() error
}

type IStoreObj interface {
	HasObj(key interface{}) (bool, error)
	GetObj(key interface{}, value interface{}) error
	GetValue(key interface{}, value interface{}) (bool, error)
	DeleteObj(key interface{}) error
	PutObj(key interface{}, value interface{}) error
}

type IStore interface {
	ILevelDb
	IStoreObj
	WrapBatch(fn func(batch IBatch) error) error
}

type Iterator interface {
	Next() bool
	Error() error
	Key() []byte
	Value() []byte
	Release()
}
