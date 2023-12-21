package store

type ILevelDb interface {
	Has(key []byte) (bool, error)
	Get(key []byte) ([]byte, error)
	Put(key []byte, value []byte) error
	Delete(key []byte) error
	Iterator(prefix []byte, start []byte) Iterator
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

type Iterator interface {
	// Next moves the iterator to the next key/value pair. It returns whether the
	// iterator is exhausted.
	Next() bool

	// Error returns any accumulated error. Exhausting all the key/value pairs
	// is not considered to be an error.
	Error() error

	// Key returns the key of the current key/value pair, or nil if done. The caller
	// should not modify the contents of the returned slice, and its contents may
	// change on the next call to Next.
	Key() []byte

	// Value returns the value of the current key/value pair, or nil if done. The
	// caller should not modify the contents of the returned slice, and its contents
	// may change on the next call to Next.
	Value() []byte

	// Release releases associated resources. Release should always succeed and can
	// be called multiple times without causing error.
	Release()
}
