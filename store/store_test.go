package store

import (
	"fmt"
	"github.com/lightec-xyz/daemon/logger"
	"testing"
)

var db IStore
var err error

func init() {
	logger.InitLogger()
	db, err = NewStore("/Users/red/.daemon/testnet", 0, 0, "zkbtc", false)
	if err != nil {
		panic(err)
	}
}

func TestStore_Iterator(t *testing.T) {
	err := db.Put([]byte("test_01"), []byte("test01"))
	if err != nil {
		t.Fatal(err)
	}
	err = db.Put([]byte("test_02"), []byte("test02"))
	if err != nil {
		t.Fatal(err)
	}
	err = db.Put([]byte("test_03"), []byte("test03"))
	if err != nil {
		t.Fatal(err)
	}
	err = db.Put([]byte("test_04"), []byte("test04"))
	if err != nil {
		t.Fatal(err)
	}

	iterator := db.Iterator([]byte("test"), []byte("_03"))
	defer iterator.Release()
	for iterator.Next() {
		key := iterator.Key()
		value := iterator.Value()
		fmt.Printf("Key: %s, Value: %s\n", key, value)
	}
	err = iterator.Error()
	if err != nil {
		t.Fatal(err)
	}

}

func TestStore_Demo(t *testing.T) {
	has, err := db.Has([]byte("ethCurHeight"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(has)
	var result interface{}
	err = db.GetObj([]byte("ethCurHeight"), &result)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}
