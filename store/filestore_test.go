package store

import (
	"testing"
)

func TestNewFileStore(t *testing.T) {
	fileStore, err := NewFileStore("test", "unit", "json")
	if err != nil {
		t.Error(err)
	}
	fileName := "10"
	fileData := "{}"

	exists, err := fileStore.CheckExists(fileName)
	if err != nil {
		t.Error(err)
	}
	t.Log(exists)
	err = fileStore.Store(fileName, fileData)
	if err != nil {
		t.Error(err)
	}
	exists, err = fileStore.CheckExists(fileName)
	if err != nil {
		t.Error(err)
	}
	t.Log(exists)
	data, err := fileStore.GetData(fileName)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(data))
	var result string
	err = fileStore.GetObj(fileName, &result)
	if err != nil {
		t.Error(err)
	}
	t.Log(result)
	var result2 string
	exists, err = fileStore.Get(fileName, &result2)
	if err != nil {
		t.Error(err)
	}
	t.Log(exists, result2)
	files, err := fileStore.AllFiles()
	if err != nil {
		t.Error(err)
	}
	t.Log(files)
	indexes, err := fileStore.AllIndexes()
	if err != nil {
		t.Error(err)
	}
	t.Log(indexes)
}

func TestDemo(t *testing.T) {

}
