package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/lightec-xyz/btc_provers/circuits/baselevel"
	"github.com/lightec-xyz/btc_provers/circuits/common"
	"github.com/lightec-xyz/btc_provers/circuits/midlevel"
	"github.com/lightec-xyz/btc_provers/circuits/recursiveduper"
	"github.com/lightec-xyz/btc_provers/circuits/upperlevel"
	baselevelUtil "github.com/lightec-xyz/btc_provers/utils/baselevel"
	"github.com/lightec-xyz/btc_provers/utils/client"
	midlevelUtil "github.com/lightec-xyz/btc_provers/utils/midlevel"
	upperlevelUtil "github.com/lightec-xyz/btc_provers/utils/upperlevel"
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

	btcClient, err := client.NewClient(cfg.Url, cfg.BtcUser, cfg.BtcPwd)
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
		startBlockheight = (lastestBh/common.CapacityDifficultyBlock - 3) * common.CapacityDifficultyBlock
	}

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
	endBlockHeight := bs.cfg.EndHeight
	_, err := bs.uplevelProve(endBlockHeight)
	if err != nil {
		logger.Error("uplevel prove error: %v", err)
		return err
	}
	return nil

}

func (bs *BtcSetup) baseProve(beginHeight uint32) (*reLight_common.Proof, error) {
	endHeight := beginHeight + common.CapacityBaseLevel - 1
	logger.Info("start baseLevel prove: %v~%v", beginHeight, endHeight)

	baseData, err := baselevelUtil.GetBaseLevelProofData(bs.client, endHeight)
	if err != nil {
		logger.Error("get baseLevel proof data error: %v %v", beginHeight, err)
		return nil, err
	}

	baseProof, err := baselevel.Prove(bs.cfg.SetupDir, bs.cfg.DataDir, baseData)
	if err != nil {
		logger.Error("baseLevel prove error: %v %v", beginHeight, err)
		return nil, err
	}

	err = bs.fileStore.StoreBase(genKey(string(baseTable), beginHeight, endHeight), baseProof)
	if err != nil {
		logger.Error("store base level proof error %v %v", beginHeight, err)
		return nil, err
	}

	logger.Info("complete baseLevel prove: %v~%v", beginHeight, endHeight)
	return baseProof, nil
}

func (bs *BtcSetup) middleProve(beginHeight uint32) (*reLight_common.Proof, error) {
	baseProofs := make([]reLight_common.Proof, common.CapacityMidLevel)

	for i := uint32(0); i < common.CapacityMidLevel; i++ {
		baseBeginHeight := beginHeight + i*common.CapacityBaseLevel

		baseProof, err := bs.baseProve(baseBeginHeight)
		if err != nil {
			logger.Error("base prove error: %v %v", baseBeginHeight, err)
			return nil, err
		}

		baseProofs[i] = *baseProof
	}

	endHeight := beginHeight + common.CapacitySuperBatch - 1
	logger.Info("start middle prove: %v~%v", beginHeight, endHeight)

	middleData, err := midlevelUtil.GetMidLevelProofData(bs.client, endHeight)
	if err != nil {
		logger.Error("get middle data error: %v %v", beginHeight, err)
		return nil, err
	}

	middleProof, err := midlevel.Prove(bs.cfg.SetupDir, bs.cfg.DataDir, middleData, baseProofs)
	if err != nil {
		logger.Error("middLevel prove error: %v %v", beginHeight, err)
		return nil, err
	}

	err = bs.fileStore.StoreMiddle(genKey(string(middleTable), beginHeight, endHeight), middleProof)
	if err != nil {
		logger.Error("store middle level proof error %v %v", beginHeight, err)
		return nil, err
	}

	logger.Info("complete middle level proof: %v~%v", beginHeight, endHeight)
	return middleProof, nil

}

func (bs *BtcSetup) uplevelProve(beginHeight uint32) (*reLight_common.Proof, error) {
	midProofs := make([]reLight_common.Proof, common.CapacityUpperLevel)

	for i := uint32(0); i < common.CapacityUpperLevel; i++ {
		midBeginHeight := beginHeight + i*common.CapacitySuperBatch

		midProof, err := bs.middleProve(midBeginHeight)
		if err != nil {
			logger.Error("base prove error: %v %v", midBeginHeight, err)
			return nil, err
		}

		midProofs[i] = *midProof
	}

	endHeight := beginHeight + common.CapacityDifficultyBlock - 1
	logger.Info("start upLevel prove: %v~%v", beginHeight, endHeight)

	upData, err := upperlevelUtil.GetUpperLevelProofData(bs.client, endHeight)
	if err != nil {
		logger.Error("get upLevel data : %v %v", beginHeight, err)
		return nil, err
	}

	upProof, err := upperlevel.Prove(bs.cfg.SetupDir, bs.cfg.DataDir, upData, midProofs)
	if err != nil {
		logger.Error("upLevel prove error: %v %v", beginHeight, err)
		return nil, err
	}

	err = bs.fileStore.StoreUp(genKey(string(upTable), beginHeight, endHeight), upProof)
	if err != nil {
		logger.Error("store upLevel proof error %v %v", beginHeight, err)
		return nil, err
	}

	logger.Info("complete upLevel prove: %v~%v", beginHeight, endHeight)
	return upProof, nil
}

type RunConfig struct {
	DataDir       string `json:"datadir"`
	SetupDir      string `json:"setupdir"`
	SrsDir        string `json:"srsdir"`
	Setup         bool   `json:"setup"`
	IsFromGenesis bool   `json:"isFromGenesis"`
	Url           string `json:"url"`
	BtcUser       string `json:"btcUser"`
	BtcPwd        string `json:"btcPwd"`
}

func (rc *RunConfig) check() error {
	if rc.Url == "" {
		return fmt.Errorf("url is empty")
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
