package node

import (
	"encoding/json"
	"fmt"
	"github.com/lightec-xyz/reLight/circuits/utils"
	"github.com/prysmaticlabs/prysm/v5/api/server/structs"
	"os"
	"testing"
)

var fileStore *FileStore
var err error

func init() {
	fileStore, err = NewFileStore("test", 100)
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
	err = fileStore.StoreBootstrap("update")
	if err != nil {
		t.Fatal(err)
	}
	genesisUpdate, err := fileStore.CheckBootstrap()
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
	files, err := traverseFile("//Users/red/lworkspace/lightec/daemon/node/test/proofData/update")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(files)
}

func TestDemo001(t *testing.T) {

	for index := uint64(0); index < 162; index++ {
		var update structs.LightClientUpdateWithVersion
		exists, err := fileStore.GetUpdate(index, &update)
		if err != nil {
			t.Fatal(err)
		}
		if !exists {
			t.Fatal(err)
		}
		var reUpdate utils.LightClientUpdateInfo
		err = deepCopy(update.Data, &reUpdate)
		if err != nil {
			t.Fatal(err)
		}
		reUpdate.Version = update.Version
		if index == 0 {
			var bootstrap structs.LightClientBootstrapResponse
			exists, err := fileStore.GetBootstrap(&bootstrap)
			if err != nil {
				t.Fatal(err)
			}
			if !exists {
				t.Fatal(err)
			}
			var genesisCommittee utils.SyncCommittee
			err = deepCopy(bootstrap.Data.CurrentSyncCommittee, &genesisCommittee)
			if err != nil {
				t.Fatal(err)
			}
			reUpdate.CurrentSyncCommittee = &genesisCommittee
		} else {
			preIndex := index - 1
			var preUpdate structs.LightClientUpdateWithVersion
			exists, err := fileStore.GetUpdate(preIndex, &preUpdate)
			if err != nil {
				t.Fatal(err)
			}
			if !exists {
				panic("not exists")
			}
			var currentSyncCommittee utils.SyncCommittee
			err = deepCopy(preUpdate.Data.NextSyncCommittee, &currentSyncCommittee)
			if err != nil {
				t.Fatal(err)
			}
			reUpdate.CurrentSyncCommittee = &currentSyncCommittee
		}
		t.Log(reUpdate)
		reData, err := json.Marshal(reUpdate)
		if err != nil {
			t.Fatal(err)
		}
		dir := "/Users/red/lworkspace/lightec/daemon/node/test/parseUpdate"
		ok, err := fileExists(dir)
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			os.MkdirAll(dir, os.ModePerm)
		}
		fileName := fmt.Sprintf("%s/holesky_sync_committee_update_%d.json", dir, index)
		err = WriteFile(fileName, reData)
		if err != nil {
			t.Fatal(err)
		}

	}

}
