package beacon

//
//import (
//	"encoding/json"
//	"fmt"
//	"github.com/lightec-xyz/reLight/circuits/utils"
//	"github.com/prysmaticlabs/prysm/v5/api/server/structs"
//	"github.com/prysmaticlabs/prysm/v5/testing/require"
//	"math/big"
//	"os"
//	"path/filepath"
//	"testing"
//)
//
//func TestLightClientUpdateInfo(t *testing.T) {
//	for i := 157; i <= 160; i++ {
//		file := filepath.Join("testdata", fmt.Sprintf("sc%v", i), fmt.Sprintf("holesky_sync_committee_update_%v.json", i))
//		data, err := os.ReadFile(file)
//		require.NoError(t, err)
//		update := &utils.LightClientUpdateInfo{}
//		err = json.Unmarshal(data, update)
//		require.NoError(t, err)
//		valid, err := verifyLightClientUpdateInfo(update)
//		require.Equal(t, true, valid)
//		require.NoError(t, err)
//		fmt.Printf("verify pass:%v, version:%v\n", i, update.Version)
//	}
//}
//
//func TestFileExist(t *testing.T) {
//	_, err := os.Stat("test.json")
//	exist := os.IsExist(err)
//	fmt.Printf("file exist: %v\n", exist)
//}
//
//func buildAndVerifyLightClientUpdateInfos(bootStrapFile string, updateFiles []string) ([]utils.LightClientUpdateInfo, error) {
//	data, err := os.ReadFile(bootStrapFile)
//	if err != nil {
//		return nil, err
//	}
//
//	bootStrap := structs.LightClientBootstrapResponse{}
//	err = json.Unmarshal(data, &bootStrap)
//	if err != nil {
//		return nil, err
//	}
//
//	updates := make([]structs.LightClientUpdateWithVersion, 0)
//	for _, f := range updateFiles {
//		data, err = os.ReadFile(f)
//		if err != nil {
//			return nil, err
//		}
//		update := structs.LightClientUpdateWithVersion{}
//		err = json.Unmarshal(data, &update)
//		if err != nil {
//			return nil, err
//		}
//		updates = append(updates, update)
//	}
//
//	infos := make([]utils.LightClientUpdateInfo, 0)
//	info := utils.LightClientUpdateInfo{
//		Version:                 updates[0].Version,
//		AttestedHeader:          updates[0].Data.AttestedHeader,
//		CurrentSyncCommittee:    bootStrap.Data.CurrentSyncCommittee,
//		SyncAggregate:           updates[0].Data.SyncAggregate,
//		NextSyncCommittee:       updates[0].Data.NextSyncCommittee,
//		NextSyncCommitteeBranch: updates[0].Data.NextSyncCommitteeBranch,
//		FinalizedHeader:         updates[0].Data.FinalizedHeader,
//		FinalityBranch:          updates[0].Data.FinalityBranch,
//		SignatureSlot:           updates[0].Data.SignatureSlot,
//	}
//
//	infos = append(infos, info)
//	for i := 1; i < len(updates); i++ {
//		info = LightClientUpdateInfo{
//			Version:                 updates[i].Version,
//			AttestedHeader:          updates[i].Data.AttestedHeader,
//			CurrentSyncCommittee:    updates[i-1].Data.NextSyncCommittee,
//			SyncAggregate:           updates[i].Data.SyncAggregate,
//			NextSyncCommittee:       updates[i].Data.NextSyncCommittee,
//			NextSyncCommitteeBranch: updates[i].Data.NextSyncCommitteeBranch,
//			FinalizedHeader:         updates[i].Data.FinalizedHeader,
//			FinalityBranch:          updates[i].Data.FinalityBranch,
//			SignatureSlot:           updates[i].Data.SignatureSlot,
//		}
//		infos = append(infos, info)
//	}
//
//	for i := 0; i < len(infos); i++ {
//		valid, err := verifyLightClientUpdateInfo(&infos[i])
//		if err != nil {
//			return nil, err
//		}
//		if !valid {
//			return nil, fmt.Errorf("%v verifyLightClientUpdateInfo failed", i)
//		}
//		fmt.Printf("%v verify pass, version:%v\n", infos[i].AttestedHeader.Slot, infos[i].Version)
//	}
//	return infos, nil
//}
//
//// build light client update info from bootstrap file + followed update files
//func TestBuildLightClientUpdateInfo(t *testing.T) {
//	bootStrapFile := "testdata/bootstrap_157.json"
//	updateFiles := []string{"testdata/light_client_update_157.json", "testdata/light_client_update_158.json", "testdata/light_client_update_159.json", "testdata/light_client_update_160.json", "testdata/light_client_update_161.json"}
//
//	infos, err := buildAndVerifyLightClientUpdateInfos(bootStrapFile, updateFiles)
//	require.NoError(t, err)
//
//	for i := 0; i < len(infos); i++ {
//		slot, ok := big.NewInt(0).SetString(infos[i].AttestedHeader.Slot, 10)
//		require.Equal(t, true, ok)
//		period := slot.Uint64() / 8192
//		err = os.MkdirAll(fmt.Sprintf("testdata/sc%d", period), 0775)
//		require.NoError(t, err)
//
//		fn := fmt.Sprintf("holesky_sync_committee_update_%d.json", period)
//
//		f, err := os.Create(fn)
//		require.NoError(t, err)
//		defer f.Close()
//
//		data, err := json.Marshal(infos[i])
//		require.NoError(t, err)
//		_, err = f.Write(data)
//	}
//}
