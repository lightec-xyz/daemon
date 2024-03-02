package node

import "testing"

var fileStore *FileStore
var err error

func init() {
	fileStore, err = NewFileStore("test")
	if err != nil {
		panic(err)
	}
}
func TestFileStoreGenesis(t *testing.T) {
	err := fileStore.StoreLatestPeriod(1)
	if err != nil {
		t.Fatal(err)
	}
	err = fileStore.StoreGenesisUpdate("update")
	if err != nil {
		t.Fatal(err)
	}
	err = fileStore.StoreUpdate(1, "update")
	if err != nil {
		t.Fatal(err)
	}

	err = fileStore.StoreUnitProof(1, "unit")
	if err != nil {
		t.Fatal(err)
	}
	err = fileStore.StoreRecursiveProof(1, "recursive")
	if err != nil {
		t.Fatal(err)
	}
}

func TestFileLatestPeriod(t *testing.T) {

	existsPeriod, err := fileStore.CheckLatestPeriod()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(existsPeriod)
	err = fileStore.StoreLatestPeriod(100)
	if err != nil {
		t.Fatal(err)
	}
	existsPeriod, err = fileStore.CheckLatestPeriod()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(existsPeriod)
	period, err := fileStore.GetLatestPeriod()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(period)
	err = fileStore.StoreLatestPeriod(123)
	if err != nil {
		t.Fatal(err)
	}
	period, err = fileStore.GetLatestPeriod()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(period)
}
