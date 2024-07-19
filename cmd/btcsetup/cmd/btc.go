package cmd

import (
	"encoding/json"
	"fmt"
	native_plonk "github.com/consensys/gnark/backend/plonk"
	"github.com/consensys/gnark/backend/witness"
	"github.com/lightec-xyz/btc_provers/circuits/baselevel"
	"github.com/lightec-xyz/btc_provers/circuits/midlevel"
	btcprovertypes "github.com/lightec-xyz/btc_provers/circuits/types"
	"github.com/lightec-xyz/btc_provers/circuits/upperlevel"
	baselevelUtil "github.com/lightec-xyz/btc_provers/utils/baselevel"
	"github.com/lightec-xyz/btc_provers/utils/client"
	midlevelUtil "github.com/lightec-xyz/btc_provers/utils/midlevel"
	upperlevelUtil "github.com/lightec-xyz/btc_provers/utils/upperlevel"
	"github.com/lightec-xyz/daemon/logger"
	"os"
)

const (
	upDistance     = 6 * middleDistance
	middleDistance = 6 * baseDistance
	baseDistance   = 56
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
	return &BtcSetup{
		cfg:       cfg,
		client:    btcClient,
		exit:      make(chan os.Signal, 1),
		fileStore: fileStorage,
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
	logger.Debug("start base setup ...")
	err := baselevel.Setup(bs.cfg.SetupDir, bs.cfg.SrsDir)
	if err != nil {
		logger.Error("setup baseLevel error: %v", err)
		return err
	}
	logger.Debug("start midlevel setup ...")
	err = midlevel.Setup(bs.cfg.SetupDir, bs.cfg.SrsDir)
	if err != nil {
		logger.Error("setup middleLevel error: %v", err)
		return err
	}
	logger.Debug("start up setup ...")
	err = upperlevel.Setup(bs.cfg.SetupDir, bs.cfg.SrsDir)
	if err != nil {
		logger.Error("setup upLevel error: %v", err)
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

func (bs *BtcSetup) baseProve(endHeight int64) (*btcprovertypes.WitnessFile, error) {
	var proofs []native_plonk.Proof
	var witnesses []witness.Witness
	start := endHeight - 6*baseDistance
	if start < 0 {
		return nil, fmt.Errorf("endHeight less than 0")
	}
	for index := start; index < endHeight; index = index + baseDistance {
		eHeight := index + baseDistance
		logger.Info("start baseLevel prove: %v~%v", index, eHeight)
		baseData, err := baselevelUtil.GetBaseLevelProofData(bs.client, uint32(eHeight))
		if err != nil {
			logger.Error("get baseLevel proof data error: %v %v", index, err)
			return nil, err
		}
		baseProof, err := baselevel.Prove(bs.cfg.SetupDir, bs.cfg.DataDir, baseData)
		if err != nil {
			logger.Error("baseLevel prove error: %v %v", index, err)
			return nil, err
		}
		err = bs.fileStore.StoreBase(genKey("base", index, eHeight), baseProof.Proof, baseProof.Wit)
		if err != nil {
			logger.Error("store base level proof error %v %v", index, err)
			return nil, err
		}
		proofs = append(proofs, baseProof.Proof)
		witnesses = append(witnesses, baseProof.Wit)
		logger.Info("complete baseLevel prove: %v~%v", index, eHeight)
	}
	return &btcprovertypes.WitnessFile{
		Proofs:    proofs,
		Witnesses: witnesses,
	}, nil
}

func (bs *BtcSetup) middleProve(endHeight int64) (*btcprovertypes.WitnessFile, error) {
	start := endHeight - 6*middleDistance
	if start < 0 {
		return nil, fmt.Errorf("middle endHeight less than 0")
	}
	var proofs []native_plonk.Proof
	var witnesses []witness.Witness
	for index := start; index < endHeight; index = index + middleDistance {
		eHeight := index + middleDistance
		logger.Info("start middle prove: %v~%v", index, eHeight)
		baseProof, err := bs.baseProve(eHeight)
		if err != nil {
			logger.Error("base prove error: %v %v", index, err)
			return nil, err
		}
		middleData, err := midlevelUtil.GetMidLevelProofData(bs.client, uint32(eHeight-1))
		if err != nil {
			logger.Error("get middle data error: %v %v", index, err)
			return nil, err
		}
		middleProof, err := midlevel.Prove(bs.cfg.SetupDir, bs.cfg.DataDir, middleData, baseProof)
		if err != nil {
			logger.Error("middLevel prove error: %v %v", index, err)
			return nil, err
		}
		err = bs.fileStore.StoreMiddle(genKey("middle", index, eHeight), middleProof.Proof, middleProof.Wit)
		if err != nil {
			logger.Error("store middle level proof error %v %v", index, err)
			return nil, err
		}
		proofs = append(proofs, middleProof.Proof)
		witnesses = append(witnesses, middleProof.Wit)
		logger.Info("complete middle level proof: %v~%v", index, eHeight)
	}
	return &btcprovertypes.WitnessFile{
		Proofs:    proofs,
		Witnesses: witnesses,
	}, nil
}

func (bs *BtcSetup) uplevelProve(endHeight int64) (*btcprovertypes.WitnessFile, error) {
	start := endHeight - 3*upDistance // todo
	if start < 0 {
		return nil, fmt.Errorf("up height less than 0")
	}
	for index := start; index < endHeight; index = index + upDistance {
		eHeight := index + upDistance
		logger.Info("start upLevel prove: %v~%v", index, eHeight)

		middleProof, err := bs.middleProve(eHeight)
		if err != nil {
			logger.Error("middle prove error: %v %v", index, err)
			return nil, err
		}
		upData, err := upperlevelUtil.GetUpperLevelProofData(bs.client, uint32(eHeight-1))
		if err != nil {
			logger.Error("get upLevel data : %v %v", endHeight, err)
			return nil, err
		}
		upProof, err := upperlevel.Prove(bs.cfg.SetupDir, bs.cfg.DataDir, upData, middleProof)
		if err != nil {
			logger.Error("upLevel prove error: %v %v", endHeight, err)
			return nil, err
		}
		err = bs.fileStore.StoreUp(genKey("up", index, endHeight), upProof.Proof, upProof.Wit)
		if err != nil {
			logger.Error("store upLevel proof error %v %v", index, err)
			return nil, err
		}
		logger.Info("complete upLevel prove: %v~%v", index, eHeight)
	}
	return nil, nil
}

type RunConfig struct {
	DataDir   string `json:"datadir"`
	SetupDir  string `json:"setupdir"`
	SrsDir    string `json:"srsdir"`
	Setup     bool   `json:"setup"`
	EndHeight int64  `json:"endHeight"`
	Url       string `json:"url"`
	BtcUser   string `json:"btcUser"`
	BtcPwd    string `json:"btcPwd"`
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
