package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/lightec-xyz/btc_provers/circuits/blockchain/baselevel"
	"github.com/lightec-xyz/btc_provers/circuits/blockchain/midlevel"
	"github.com/lightec-xyz/btc_provers/circuits/blockchain/recursiveduper"
	"github.com/lightec-xyz/btc_provers/circuits/blockchain/upperlevel"
	"github.com/lightec-xyz/btc_provers/circuits/common"
	baselevelUtil "github.com/lightec-xyz/btc_provers/utils/blockchain"
	midlevelUtil "github.com/lightec-xyz/btc_provers/utils/blockchain"
	recursiveduperUtil "github.com/lightec-xyz/btc_provers/utils/blockchain"
	upperlevelUtil "github.com/lightec-xyz/btc_provers/utils/blockchain"
	"github.com/lightec-xyz/btc_provers/utils/client"
	"github.com/lightec-xyz/daemon/logger"
	reLight_common "github.com/lightec-xyz/reLight/circuits/common"
)

type BtcSetup struct {
	cfg              *RunConfig
	client           *client.Client
	exit             chan os.Signal
	fileStore        *FileStorage
	startBlockheight uint32
}

func NewBtcSetup(cfg *RunConfig) (*BtcSetup, error) {
	err := logger.InitLogger(nil)
	if err != nil {
		return nil, err
	}

	err = cfg.check()
	if err != nil {
		return nil, err
	}

	btcClient, err := client.NewClient(cfg.BtcHost, cfg.BtcUser, cfg.BtcPwd)
	if err != nil {
		logger.Error("new btcClient error: %v", err)
		return nil, err
	}

	proofPath := fmt.Sprintf("%v/proof", cfg.DataDir)
	fileStorage, err := NewFileStorage(proofPath)
	if err != nil {
		logger.Error("new file storage error: %v %v", proofPath, err)
		return nil, err
	}

	startBlockheight := uint32(0)
	if !cfg.IsFromGenesis {
		lastestBh, err := btcClient.GetBlockCount()
		if err != nil {
			logger.Error("get block height error: %v", err)
			return nil, err
		}
		startBlockheight = (uint32(lastestBh)/common.CapacityDifficultyBlock - 3) * common.CapacityDifficultyBlock
	}
	logger.Info("start block height: %v", startBlockheight)

	return &BtcSetup{
		cfg:              cfg,
		client:           btcClient,
		exit:             make(chan os.Signal, 1),
		fileStore:        fileStorage,
		startBlockheight: startBlockheight,
	}, nil
}

func (bs *BtcSetup) Run() error {
	if bs.cfg.Setup {
		err := bs.Setup()
		if err != nil {
			logger.Error("bs setup error: %v", err)
			return err
		}
	}
	err := bs.Prove()
	if err != nil {
		logger.Error("bs prove error: %v", err)
		return err
	}
	return nil
}

func (bs *BtcSetup) Close() error {
	return nil
}

func (bs *BtcSetup) Setup() error {
	logger.Debug("start baselevel setup ...")

	err := baselevel.Setup(bs.cfg.SetupDir, bs.cfg.SrsDir)
	if err != nil {
		logger.Error("setup baselevel error: %v", err)
		return err
	}

	logger.Debug("start midlevel setup ...")
	err = midlevel.Setup(bs.cfg.SetupDir, bs.cfg.SrsDir)
	if err != nil {
		logger.Error("setup midlevel error: %v", err)
		return err
	}

	logger.Debug("start upperlevel setup ...")
	err = upperlevel.Setup(bs.cfg.SetupDir, bs.cfg.SrsDir)
	if err != nil {
		logger.Error("setup upperlevel error: %v", err)
		return err
	}

	logger.Debug("start recursiveduper setup ...")
	err = recursiveduper.Setup(bs.cfg.SetupDir, bs.cfg.SrsDir)
	if err != nil {
		logger.Error("setup recursiveduper error: %v", err)
		return err
	}

	return nil
}

func (bs *BtcSetup) Prove() error {
	duperProofs := make([]reLight_common.Proof, 3)
	beginHeight := bs.startBlockheight

	for i := uint32(0); i < 3; i++ {
		duperBeginHeight := beginHeight + i*common.CapacityDifficultyBlock

		duperProof, err := bs.DuperProve(duperBeginHeight)
		if err != nil {
			logger.Error("DuperProve(begin Height: %v) error: %v", duperBeginHeight, err)
			return err
		}

		duperProofs[i] = *duperProof
	}

	endHeight1 := beginHeight + common.CapacityDifficultyBlock*2 - 1
	logger.Info("start genesis recursiveduper prove: %v~%v", beginHeight, endHeight1)

	proofData1, err := recursiveduperUtil.GetRecursiveProofData(bs.client, endHeight1, beginHeight)
	if err != nil {
		logger.Error("get recursiveduper data error: %v %v", endHeight1, err)
		return err
	}

	genesisProof, err := recursiveduper.ProveGenesis(bs.cfg.SetupDir, &duperProofs[0], &duperProofs[1], proofData1)
	if err != nil {
		logger.Error("genesis recursiveduper prove error: %v %v", endHeight1, err)
		return err
	}

	err = bs.fileStore.StoreRecursive(genKey(string(recursiveTable), beginHeight, endHeight1), genesisProof)
	if err != nil {
		logger.Error("store recursiveduper proof error %v %v", endHeight1, err)
		return err
	}

	recursiveduper.SaveProof(bs.cfg.DataDir, genesisProof, endHeight1, beginHeight)
	logger.Info("complete genesis recursiveduper prove: %v~%v", beginHeight, endHeight1)

	endHeight2 := endHeight1 + common.CapacityDifficultyBlock
	logger.Info("start recursiveduper prove: %v~%v", beginHeight, endHeight2)

	proofData2, err := recursiveduperUtil.GetRecursiveProofData(bs.client, endHeight2, beginHeight)
	if err != nil {
		logger.Error("get recursiveduper data error: %v %v", endHeight2, err)
		return err
	}

	recursiveProof, err := recursiveduper.ProveRecursive(
		bs.cfg.SetupDir, genesisProof, &duperProofs[2], proofData2)
	if err != nil {
		logger.Error("recursiveduper prove error: %v %v", endHeight2, err)
		return err
	}

	err = bs.fileStore.StoreRecursive(genKey(string(recursiveTable), beginHeight, endHeight2), recursiveProof)
	if err != nil {
		logger.Error("store recursiveduper proof error %v %v", endHeight2, err)
		return err
	}

	recursiveduper.SaveProof(bs.cfg.DataDir, recursiveProof, endHeight2, beginHeight)
	logger.Info("complete recursiveduper prove: %v~%v", beginHeight, endHeight2)
	return nil

}

func (bs *BtcSetup) BatchProve(beginHeight uint32) (*reLight_common.Proof, error) {
	endHeight := beginHeight + common.CapacityBaseLevel - 1
	logger.Info("start baseLevel prove: %v~%v", beginHeight, endHeight)

	baseData, err := baselevelUtil.GetBaseLevelProofData(bs.client, endHeight)
	if err != nil {
		logger.Error("get baseLevel proof data error: %v~%v %v", beginHeight, endHeight, err)
		return nil, err
	}

	baseProof, err := baselevel.Prove(bs.cfg.SetupDir, baseData)
	if err != nil {
		logger.Error("baseLevel prove error: %v~%v %v", beginHeight, endHeight, err)
		return nil, err
	}

	err = bs.fileStore.StoreBase(genKey(string(baseTable), beginHeight, endHeight), baseProof)
	if err != nil {
		logger.Error("store baseLevel proof error %v~%v %v", beginHeight, endHeight, err)
		return nil, err
	}

	baselevel.SaveProof(bs.cfg.DataDir, baseProof, endHeight)
	logger.Info("complete baseLevel prove: %v~%v", beginHeight, endHeight)
	return baseProof, nil
}

func (bs *BtcSetup) SuperProve(beginHeight uint32) (*reLight_common.Proof, error) {
	batchProofs := make([]reLight_common.Proof, common.CapacityMidLevel)

	for i := uint32(0); i < common.CapacityMidLevel; i++ {
		batchedBeginHeight := beginHeight + i*common.CapacityBaseLevel

		batchProof, err := bs.BatchProve(batchedBeginHeight)
		if err != nil {
			logger.Error("BatchProve(begin Height: %v) error: %v", batchedBeginHeight, err)
			return nil, err
		}

		batchProofs[i] = *batchProof
	}

	endHeight := beginHeight + common.CapacitySuperBatch - 1
	logger.Info("start middLevel prove: %v~%v", beginHeight, endHeight)

	middleData, err := midlevelUtil.GetBatchedProofData(bs.client, endHeight)
	if err != nil {
		logger.Error("get middLevel data error: %v~%v %v", beginHeight, endHeight, err)
		return nil, err
	}

	middleProof, err := midlevel.Prove(bs.cfg.SetupDir, middleData, batchProofs)
	if err != nil {
		logger.Error("middLevel prove error: %v~%v %v", beginHeight, endHeight, err)
		return nil, err
	}

	err = bs.fileStore.StoreMiddle(genKey(string(middleTable), beginHeight, endHeight), middleProof)
	if err != nil {
		logger.Error("store middLevel level proof error: %v~%v %v", beginHeight, endHeight, err)
		return nil, err
	}

	midlevel.SaveProof(bs.cfg.DataDir, middleProof, endHeight)
	logger.Info("complete middLevel level proof: %v~%v", beginHeight, endHeight)
	return middleProof, nil

}

func (bs *BtcSetup) DuperProve(beginHeight uint32) (*reLight_common.Proof, error) {
	superProofs := make([]reLight_common.Proof, common.CapacityUpperLevel)

	for i := uint32(0); i < common.CapacityUpperLevel; i++ {
		superBeginHeight := beginHeight + i*common.CapacitySuperBatch

		superProof, err := bs.SuperProve(superBeginHeight)
		if err != nil {
			logger.Error("SuperProve(begin Height: %v) error: %v", superBeginHeight, err)
			return nil, err
		}

		superProofs[i] = *superProof
	}

	endHeight := beginHeight + common.CapacityDifficultyBlock - 1
	logger.Info("start upperlevel prove: %v~%v", beginHeight, endHeight)

	upData, err := upperlevelUtil.GetBatchedProofData(bs.client, endHeight)
	if err != nil {
		logger.Error("get upperlevel data error: %v~%v %v", beginHeight, endHeight, err)
		return nil, err
	}

	upProof, err := upperlevel.Prove(bs.cfg.SetupDir, upData, superProofs)
	if err != nil {
		logger.Error("upperlevel prove error: %v~%v %v", beginHeight, endHeight, err)
		return nil, err
	}

	err = bs.fileStore.StoreUp(genKey(string(upTable), beginHeight, endHeight), upProof)
	if err != nil {
		logger.Error("store upperlevel proof error: %v~%v %v", beginHeight, endHeight, err)
		return nil, err
	}

	upperlevel.SaveProof(bs.cfg.DataDir, upProof, endHeight)
	logger.Info("complete upperlevel prove: %v~%v", beginHeight, endHeight)
	return upProof, nil
}

type RunConfig struct {
	IsSetup            bool   `json:"isSetup"`
	PubKeyInDfinity    string `json:"pubKeyInDfinity"`
	SrsDir             string `json:"srsDir"`
	RedeemSetupDir     string `json:"redeemSetupDir"`
	SetupDir           string `json:"setupDir"`
	ProveDir           string `json:"proveDir"`
	BtcHost            string `json:"btcHost"`
	BtcUser            string `json:"btcUser"`
	BtcPwd             string `json:"btcPwd"`
	GenesisBlockHeight int    `json:"genesisBlockHeight"`
	CpBlockHeight      int    `json:"cpBlockHeight"`
	EndBlockHeight     int    `json:"endBlockHeight"`
}

func (rc *RunConfig) check() error {
	if rc.BtcHost == "" {
		return fmt.Errorf("host is empty")
	}
	if rc.DataDir == "" {
		return fmt.Errorf("dataDir can not be empty")
	}
	if rc.SetupDir == "" {
		return fmt.Errorf("setupdir is empty")
	}
	if rc.SrsDir == "" {
		return fmt.Errorf("srsDir can not be empty")
	}
	return nil
}

func readRunConfig(path string) (*RunConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	fmt.Printf("config data: %v\n", string(data))
	var runConfig RunConfig
	err = json.Unmarshal(data, &runConfig)
	if err != nil {
		return nil, err
	}
	return &runConfig, nil
}
