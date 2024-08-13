package cmd

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/lightec-xyz/btc_provers/circuits/blockchain"
	"github.com/lightec-xyz/btc_provers/circuits/blockchain/baselevel"
	"github.com/lightec-xyz/btc_provers/circuits/blockchain/midlevel"
	"github.com/lightec-xyz/btc_provers/circuits/blockchain/recursiveduper"
	"github.com/lightec-xyz/btc_provers/circuits/blockchain/upperlevel"
	"github.com/lightec-xyz/btc_provers/circuits/blockdepth"
	"github.com/lightec-xyz/btc_provers/circuits/common"
	"github.com/lightec-xyz/btc_provers/circuits/txinchain"
	blockchainlUtil "github.com/lightec-xyz/btc_provers/utils/blockchain"
	"github.com/lightec-xyz/btc_provers/utils/client"
	"github.com/lightec-xyz/daemon/logger"
	reLight_common "github.com/lightec-xyz/reLight/circuits/common"
)

type BtcSetup struct {
	cfg       *RunConfig
	client    *client.Client
	exit      chan os.Signal
	fileStore *FileStorage
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

	fileStorage, err := NewFileStorage(cfg.ProveDir)
	if err != nil {
		logger.Error("new file storage error: %v %v", cfg.ProveDir, err)
		return nil, err
	}

	return &BtcSetup{
		cfg:       cfg,
		client:    btcClient,
		exit:      make(chan os.Signal, 1),
		fileStore: fileStorage,
	}, nil
}

func (bs *BtcSetup) Run() error {
	if bs.cfg.IsSetup {
		err := bs.Setup()
		if err != nil {
			logger.Error("setup error: %v", err)
			return err
		}
	}

	err := bs.BlockChainProve()
	if err != nil {
		logger.Error("blockchain prove error: %v", err)
		return err
	}

	return nil
}

func (bs *BtcSetup) Close() error {
	return nil
}

func (bs *BtcSetup) Setup() error {
	logger.Debug("start blockchain setup ...")
	err := blockchain.BlockChainSetup(bs.cfg.SetupDir, bs.cfg.SrsDir)
	if err != nil {
		logger.Error("setup blockchain error: %v", err)
		return err
	}

	logger.Debug("start blockdepth setup ...")
	err = blockdepth.BlockDepthSetup(bs.cfg.SetupDir, bs.cfg.SrsDir)
	if err != nil {
		logger.Error("setup blockdepth error: %v", err)
		return err
	}

	logger.Debug("start txinchain setup ...")
	pubKey, err := hex.DecodeString(bs.cfg.PubKeyInDfinity)
	if err != nil {
		logger.Error("decode pubkey error: %v", err)
		return err
	}
	if len(pubKey) != 33 {
		logger.Error("pubkey length is not 33")
		return err
	}

	err = txinchain.Setup(bs.cfg.SetupDir, bs.cfg.SrsDir, bs.cfg.RedeemSetupDir, [33]byte(pubKey))
	if err != nil {
		logger.Error("setup txinchain error: %v", err)
		return err
	}

	return nil
}

func (bs *BtcSetup) BlockChainProve() error {
	duperProofs := make([]reLight_common.Proof, 2)
	genesisHeight := uint32(bs.cfg.GenesisBlockHeight)

	for i := uint32(0); i < 2; i++ {
		duperBeginHeight := genesisHeight + i*common.CapacityDifficultyBlock

		duperProof, err := bs.DuperProve(duperBeginHeight)
		if err != nil {
			logger.Error("DuperProve(begin Height: %v) error: %v", duperBeginHeight, err)
			return err
		}

		duperProofs[i] = *duperProof
	}

	chainEndHeight := genesisHeight + common.CapacityDifficultyBlock*2 - 1
	logger.Info("start genesis recursiveduper prove: %v~%v", genesisHeight, chainEndHeight)

	proofData, err := blockchainlUtil.GetRecursiveProofData(bs.client, chainEndHeight, genesisHeight)
	if err != nil {
		logger.Error("get recursiveduper data error: %v~%v %v", genesisHeight, chainEndHeight, err)
		return err
	}

	genesisProof, err := recursiveduper.Prove(bs.cfg.SetupDir, &duperProofs[0], &duperProofs[1], proofData)
	if err != nil {
		logger.Error("genesis recursiveduper prove error: %v~%v %v", genesisHeight, chainEndHeight, err)
		return err
	}

	err = bs.fileStore.StoreRecursiveDuper(genKey(string(recursiveDuperTable), genesisHeight, chainEndHeight), genesisProof)
	if err != nil {
		logger.Error("store recursiveduper proof error %v~%v %v", genesisHeight, chainEndHeight, err)
		return err
	}

	recursiveduper.SaveProof(bs.cfg.ProveDir, genesisProof, chainEndHeight, genesisHeight)
	logger.Info("complete genesis recursiveduper prove: %v~%v", genesisHeight, chainEndHeight)

	return nil
}

func (bs *BtcSetup) BatchProve(beginHeight uint32) (*reLight_common.Proof, error) {
	endHeight := beginHeight + common.CapacityBaseLevel - 1
	logger.Info("start baseLevel prove: %v~%v", beginHeight, endHeight)

	baseData, err := blockchainlUtil.GetBaseLevelProofData(bs.client, endHeight)
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

	baselevel.SaveProof(bs.cfg.ProveDir, baseProof, endHeight)
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

	middleData, err := blockchainlUtil.GetMidLevelProofData(bs.client, endHeight)
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

	midlevel.SaveProof(bs.cfg.ProveDir, middleProof, endHeight)
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

	upData, err := blockchainlUtil.GetUpperLevelProofData(bs.client, endHeight)
	if err != nil {
		logger.Error("get upperlevel data error: %v~%v %v", beginHeight, endHeight, err)
		return nil, err
	}

	upProof, err := upperlevel.Prove(bs.cfg.SetupDir, upData, superProofs)
	if err != nil {
		logger.Error("upperlevel prove error: %v~%v %v", beginHeight, endHeight, err)
		return nil, err
	}

	err = bs.fileStore.StoreUpper(genKey(string(upperTable), beginHeight, endHeight), upProof)
	if err != nil {
		logger.Error("store upperlevel proof error: %v~%v %v", beginHeight, endHeight, err)
		return nil, err
	}

	upperlevel.SaveProof(bs.cfg.ProveDir, upProof, endHeight)
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
	if rc.PubKeyInDfinity == "" {
		return fmt.Errorf("pubKeyInDfinity is empty")
	}

	if rc.SetupDir == "" || rc.ProveDir == "" || rc.SrsDir == "" || rc.RedeemSetupDir == "" {
		return fmt.Errorf("dir is empty")
	}

	if rc.BtcHost == "" || rc.BtcUser == "" || rc.BtcPwd == "" {
		return fmt.Errorf("btc config is empty")
	}

	if rc.GenesisBlockHeight > rc.CpBlockHeight || rc.GenesisBlockHeight >= rc.EndBlockHeight {
		return fmt.Errorf("invalid genesisBlockHeight")
	}

	if rc.CpBlockHeight >= rc.EndBlockHeight {
		return fmt.Errorf("invalid cpBlockHeight")
	}

	if rc.EndBlockHeight-rc.CpBlockHeight < common.CapacityBulkUint*2 ||
		rc.EndBlockHeight-rc.GenesisBlockHeight+1 < common.CapacityDifficultyBlock*2 {
		return fmt.Errorf("invalid endBlockHeight")
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
