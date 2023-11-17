package store

import (
	"fmt"
	"github.com/emirpasic/gods/maps/treemap"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestDemo01(t *testing.T) {
	comparator := treemap.NewWithIntComparator()
	comparator.Put(4, "ddd_4_4")
	comparator.Put(3, "ddd_3_3")
	comparator.Put(1, "ddd_1_1")
	comparator.Put(2, "ddd_2_2")
	keys := comparator.Keys()
	for _, item := range keys {
		t.Log(item)
	}
	t.Log(comparator.Max())

}

func TestFileStore(t *testing.T) {
	err := filepath.WalkDir("/Users/red/lworkspace/lightec/audit/daemon/node/test/data/bakdaemon/proofData/finalityUpdate",
		func(path string, d fs.DirEntry, err error) error {
			info, err := d.Info()
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println(info.Name())
			return nil
		})
	if err != nil {
		t.Fatal(err)
	}
}

func TestDemo(t *testing.T) {
	dir, err := os.ReadDir("/Users/red/lworkspace/lightec/audit/daemon/node/test/data/bakdaemon/proofData/finalityUpdate")
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range dir {
		t.Log(item)
	}
}
