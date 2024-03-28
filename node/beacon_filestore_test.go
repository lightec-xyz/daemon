package node

import (
	"fmt"
	"regexp"
	"testing"
)

var fileStore *FileStore
var err error

func init() {
	fileStore, err = NewFileStore("test")
	if err != nil {
		panic(err)
	}
}
func TestFileStoreGenesis(t *testing.T) {
	err := fileStore.StoreLatestPeriod(123)
	if err != nil {
		t.Fatal(err)
	}
	checkLatestPeriod, err := fileStore.CheckLatestPeriod()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(checkLatestPeriod)
	period, ok, err := fileStore.GetLatestPeriod()
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal(err)
	}
	t.Log(period)
	err = fileStore.StoreGenesisUpdate("update")
	if err != nil {
		t.Fatal(err)
	}
	genesisUpdate, err := fileStore.CheckGenesisUpdate()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(genesisUpdate)

	err = fileStore.StoreUpdate(1, "update")
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

}

func TestTraverseFile(t *testing.T) {
	files, err := traverseFile("/Users/red/lworkspace/lightec/daemon/node")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(files)
}

func TestDemo001(t *testing.T) {
	// 定义一个只匹配数字的正则表达式
	pattern := regexp.MustCompile(`^\d+$`)

	// 测试字符串
	testStrings := []string{"09123", "abc123", "456xyz", "789_!@#"}

	// 遍历测试字符串
	for _, str := range testStrings {
		if pattern.MatchString(str) {
			fmt.Printf("%s 匹配\n", str)
		} else {
			fmt.Printf("%s 不匹配\n", str)
		}
	}
}
