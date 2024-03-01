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
