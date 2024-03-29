package node

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestLocalDevDaemon(t *testing.T) {
	config := LocalDevDaemonConfig()
	daemon, err := NewDaemon(config)
	if err != nil {
		t.Fatal(err)
	}
	defer daemon.Close()
	err = daemon.Init()
	if err != nil {
		t.Fatal(err)
	}
	err = daemon.Run()
	if err != nil {
		t.Fatal(err)
	}
}
func TestTestnetDaemon(t *testing.T) {
	config := TestnetDaemonConfig()
	daemon, err := NewDaemon(config)
	if err != nil {
		t.Fatal(err)
	}
	defer daemon.Close()
	err = daemon.Init()
	if err != nil {
		t.Fatal(err)
	}
	err = daemon.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestLightDaemon(t *testing.T) {
	dataDir := "/Users/red/lworkspace/lightec/daemon/node/test"
	cfg, err := NewLightLocalDaemonConfig(true, dataDir, "testnet",
		"127.0.0.1", "9780", "http://127.0.0.1:8970")
	if err != nil {
		t.Fatal(err)
	}
	daemon, err := NewRecursiveLightDaemon(cfg)
	if err != nil {
		t.Fatal(err)
	}
	err = daemon.Init()
	if err != nil {
		t.Fatal(err)
	}
	err = daemon.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestConfig(t *testing.T) {
	config := LocalDevDaemonConfig()
	data, _ := json.Marshal(config)
	fmt.Printf("%v \n", string(data))
}

func TestParseUnitUpdateData(t *testing.T) {

	//dataDir := "/Users/red/lworkspace/lightec/daemon/node/test"
	//fileStore, err := NewFileStore(dataDir)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//unitReqParamDir := fmt.Sprintf("%s/%s", dataDir, "unitReqParam")
	//ok, err := dirNotExistsAndCreate(unitReqParamDir)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//if !ok {
	//	t.Fatal("create dir error")
	//}
	//for index := 0; index <= 146; index++ {
	//	var unitPram UnitProofParam
	//	if index == 0 {
	//		var genesisData structs.LightClientBootstrapResponse
	//		err = fileStore.GetBootstrap(&genesisData)
	//		if err != nil {
	//			t.Fatal(err)
	//		}
	//		unitPram.CurrentSyncCommittee = genesisData.Data.CurrentSyncCommittee
	//	} else {
	//		var preUnitUpdateData []structs.LightClientUpdateWithVersion
	//		perPeriod := uint64(index - 1)
	//		err = fileStore.GetUpdate(perPeriod, &preUnitUpdateData)
	//		if err != nil {
	//			t.Fatal(err)
	//		}
	//		unitPram.CurrentSyncCommittee = preUnitUpdateData[0].Data.NextSyncCommittee
	//	}
	//
	//	var unitUpdateData []structs.LightClientUpdateWithVersion
	//	err = fileStore.GetUpdate(uint64(index), &unitUpdateData)
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//	unitPram.Version = unitUpdateData[0].Version
	//	unitPram.AttestedHeader = unitUpdateData[0].Data.AttestedHeader
	//	unitPram.SyncAggregate = unitUpdateData[0].Data.SyncAggregate
	//	unitPram.FinalizedHeader = unitUpdateData[0].Data.FinalizedHeader
	//	unitPram.FinalityBranch = unitUpdateData[0].Data.FinalityBranch
	//	unitPram.NextSyncCommittee = unitUpdateData[0].Data.NextSyncCommittee
	//	unitPram.NextSyncCommitteeBranch = unitUpdateData[0].Data.NextSyncCommitteeBranch
	//	unitPram.SignatureSlot = unitUpdateData[0].Data.SignatureSlot
	//
	//	dataBytes, err := json.Marshal(unitPram)
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//	filePath := fmt.Sprintf("%v/holesky_sync_committee_update_%d", unitReqParamDir, index)
	//	err = WriteFile(filePath, dataBytes)
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//	t.Logf("complete index: %v \n", index)
	//}

}
