package cmd

import (
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
		localProof, err := NewLocalProof(genesisSlot, datadir, setupDir)
		if err != nil {
			fmt.Printf("new local proof error: %v \n", err)
			return
		}
		err = localProof.GenProof(proofType, paramPath, index)
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
	}, nil

}

func (lp *LocalProof) GenProof(proofType, paramPath string, index uint64) error {
	panic("unSupport now")
	proofResponses, err := node.WorkerGenProof(lp.worker, nil)
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
	}
	return nil
}

func (lp *LocalProof) GetProofRequestData(fileStore *node.FileStore, proofType common.ZkProofType, index uint64) (interface{}, bool, error) {
	switch proofType {
	case common.SyncComUnitType:
		update, ok, err := node.GetSyncCommitUpdate(fileStore, index)
		if err != nil {
			logger.Error("get sync commit update error:%v", err)
			return nil, false, err
		}
		return update, ok, nil
	case common.SyncComGenesisType:
		return nil, false, fmt.Errorf("unSupport now  proof type: %v", proofType)
	case common.SyncComRecursiveType:
		return nil, false, fmt.Errorf("unSupport now  proof type: %v", proofType)
	default:
		return nil, false, fmt.Errorf("unSupport now  proof type: %v", proofType)
	}
}
