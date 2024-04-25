package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/node"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/spf13/cobra"
)

var paramPath string
var index uint64
var proofType string
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("datadir: %v,setupdir: %v,genesisPeriod: %v \n", datadir, setupDir, genesisSlot)
		localProof, err := NewLocalProof(genesisSlot, datadir, setupDir)
		if err != nil {
			fmt.Printf("new local proof error: %v \n", err)
			return
		}
		err = localProof.GenProof(proofType, index)
		if err != nil {
			fmt.Printf("gen proof error: %v \n", err)
			return
		}
	},
}

func init() {
	err := logger.InitLogger()
	if err != nil {
		panic(err)
	}
	genCmd.Flags().StringVar(&paramPath, "paramPath", "", "param file path")
	genCmd.Flags().Uint64Var(&index, "index", 0, "proof index")
	genCmd.Flags().StringVar(&proofType, "proofType", "", "proof type")
	rootCmd.AddCommand(genCmd)
}

type LocalProof struct {
	genesisPeriod uint64
	genesisSlot   uint64
	fileStore     *node.FileStore
	dataDir       string
	worker        rpc.IWorker
}

func NewLocalProof(genesisSlot uint64, datadir, setupDir string) (*LocalProof, error) {
	fileStore, err := node.NewFileStore(datadir, genesisSlot)
	if err != nil {
		logger.Error("new file store error:%v", err)
		return nil, err
	}
	worker, err := node.NewLocalWorker(setupDir, datadir, 1)
	if err != nil {
		logger.Error("new local worker error:%v", err)
		return nil, err
	}
	return &LocalProof{
		genesisPeriod: fileStore.GetGenesisPeriod(),
		genesisSlot:   genesisSlot,
		fileStore:     fileStore,
		worker:        worker,
		dataDir:       datadir,
	}, nil

}

func (lp *LocalProof) GenProof(proofType string, index uint64) error {
	zkProofType, err := getZkProofType(proofType)
	if err != nil {
		logger.Error("get zk proof type error:%v", err)
		return err
	}
	data, ok, err := GetProofRequestData(lp.fileStore, zkProofType, index)
	if err != nil {
		logger.Error("get proof request data error:%v", err)
		return err
	}
	if !ok {
		logger.Error("proof request data not found")
		return nil
	}
	logger.Info("start gen proof: %v %v", zkProofType.String(), index)
	zkProofRequest := common.NewZkProofRequest(zkProofType, data, index, "")
	err = lp.SaveRequest(zkProofRequest)
	if err != nil {
		logger.Error("save request error:%v", err)
		return err
	}
	proofResponses, err := node.WorkerGenProof(lp.worker, zkProofRequest)
	if err != nil {
		logger.Error("worker gen proof error:%v", err)
		return err
	}
	for _, resp := range proofResponses {
		err := node.StoreZkProof(lp.fileStore, resp.ZkProofType, resp.Period, resp.TxHash, resp.Proof, resp.Witness)
		if err != nil {
			logger.Error("store zk proof error:%v", err)
			return err
		}
		logger.Info("success store proof: %v %v", resp.ZkProofType.String(), resp.Period)
	}
	return nil
}

func (lp *LocalProof) SaveRequest(req *common.ZkProofRequest) error {
	path := fmt.Sprintf("%s/reqData/%v.json", lp.fileStore.RootPath(), req.Id())
	reqBytes, err := json.Marshal(req)
	if err != nil {
		logger.Error("json marshal error:%v", err)
		return err
	}
	err = common.WriteFile(path, reqBytes)
	if err != nil {
		logger.Error("write file error:%v", err)
		return err
	}
	return nil
}

func GetProofRequestData(fileStore *node.FileStore, proofType common.ZkProofType, index uint64) (interface{}, bool, error) {
	genesisPeriod := fileStore.GetGenesisPeriod()
	switch proofType {
	case common.SyncComUnitType:
		update, ok, err := node.GetSyncCommitUpdate(fileStore, index)
		if err != nil {
			logger.Error("get sync commit update error:%v", err)
			return nil, false, err
		}
		return update, ok, nil
	case common.SyncComGenesisType:
		data, ok, err := node.GetGenesisData(fileStore)
		if err != nil {
			logger.Error("get genesis data error:%v", err)
			return nil, false, err
		}
		return data, ok, nil
	case common.SyncComRecursiveType:
		if index == genesisPeriod+2 {
			data, ok, err := node.GetRecursiveGenesisData(fileStore, index)
			if err != nil {
				logger.Error("get recursive genesis data error:%v", err)
				return nil, false, err
			}
			return data, ok, nil
		} else if index > genesisPeriod+2 {
			data, ok, err := node.GetRecursiveData(fileStore, index)
			if err != nil {
				logger.Error("get recursive data error:%v", err)
				return nil, false, err
			}
			return data, ok, nil
		}
	case common.BeaconHeaderFinalityType:
		data, ok, err := node.GetBhfUpdateData(fileStore, index)
		if err != nil {
			logger.Error("get bhf update data error:%v", err)
			return nil, false, err
		}
		return data, ok, nil
	default:
		return nil, false, fmt.Errorf("unSupport now  proof type: %v", proofType)
	}
	return nil, false, fmt.Errorf("never reach here")
}

func getZkProofType(proofType string) (common.ZkProofType, error) {
	switch proofType {
	case "SyncComUnitType":
		return common.SyncComUnitType, nil
	case "SyncComGenesisType":
		return common.SyncComGenesisType, nil
	case "SyncComRecursiveType":
		return common.SyncComRecursiveType, nil
	case "BeaconHeaderFinalityType":
		return common.BeaconHeaderFinalityType, nil
	default:
		return 0, fmt.Errorf("unSupport now  proof type: %v", proofType)
	}
}
