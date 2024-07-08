package cmd

import (
	"encoding/json"
	"fmt"
	native_plonk "github.com/consensys/gnark/backend/plonk"
	"github.com/consensys/gnark/backend/witness"
	"github.com/lightec-xyz/btc_provers/circuits"
	"github.com/lightec-xyz/btc_provers/circuits/baselevel"
	"github.com/lightec-xyz/btc_provers/circuits/midlevel"
	"github.com/lightec-xyz/btc_provers/circuits/upperlevel"
	baselevelUtil "github.com/lightec-xyz/btc_provers/utils/baselevel"
	"github.com/lightec-xyz/btc_provers/utils/client"
	midlevelUtil "github.com/lightec-xyz/btc_provers/utils/midlevel"
	upperlevelUtil "github.com/lightec-xyz/btc_provers/utils/upperlevel"
	"github.com/lightec-xyz/daemon/logger"
	reLightCommon "github.com/lightec-xyz/reLight/circuits/common"
	"os"
)

const (
	upDistance   = 1008
	baseDistance = 56
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
	client, err := client.NewClient(cfg.Url, cfg.BtcUser, cfg.BtcPwd)
	if err != nil {
		logger.Error("new client error: %v", err)
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
		client:    client,
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
	err := baselevel.Setup(bs.cfg.DataDir, bs.cfg.SrsDir)
	if err != nil {
		logger.Error("setup baseLevel error: %v", err)
		return err
	}
	err = midlevel.Setup(bs.cfg.DataDir, bs.cfg.SrsDir)
	if err != nil {
		logger.Error("setup middleLevel error: %v", err)
		return err
	}
	err = upperlevel.Setup(bs.cfg.DataDir, bs.cfg.SrsDir)
	if err != nil {
		logger.Error("setup upLevel error: %v", err)
		return err
	}
	return nil
}

func (bs *BtcSetup) Prove() error {
	endBlockHeight := bs.cfg.EndHeight
	startIndex := endBlockHeight - 3*upDistance
	if startIndex < 0 {
		return fmt.Errorf("end height is too small: %v", bs.cfg.EndHeight)
	}
	for index := startIndex; index <= endBlockHeight; index = index + upDistance {
		baseProof, err := bs.baseProve(index)
		if err != nil {
			logger.Error("baseLevel prove error: %v %v", index, err)
			return err
		}
		middleProof, err := bs.middleProve(index, baseProof)
		if err != nil {
			logger.Error("middle prove error: %v %v", index, err)
			return err
		}
		upProof, err := bs.uplevelProve(index, middleProof)
		if err != nil {
			logger.Error("upLever prove error: %v %v", index, err)
			return err
		}
		err = bs.fileStore.StoreUp(genKey("up", index, index+upDistance), upProof.Proof, upProof.Wit)
		if err != nil {
			logger.Error("store upLevel proof error %v %v", index, err)
			return err
		}
		logger.Info("upLevel prove proof complete: %v", index)
	}
	return nil

}

func (bs *BtcSetup) baseProve(height int64) (*circuits.WitnessFile, error) {
	var proofs []native_plonk.Proof
	var witnesses []witness.Witness
	startIndex := height - 6*baseDistance
	if startIndex < 0 {
		return nil, fmt.Errorf("height less than 0")
	}
	for index := startIndex; index <= height; index = index + baseDistance {
		baseData, err := baselevelUtil.GetBaseLevelProofData(bs.client, uint32(index))
		if err != nil {
			logger.Error("get baseLevel proof data error: %v %v", index, err)
			return nil, err
		}
		baseProof, err := baselevel.Prove(bs.cfg.SetupDir, bs.cfg.DataDir, baseData)
		if err != nil {
			logger.Error("baseLevel prove error: %v %v", index, err)
			return nil, err
		}
		err = bs.fileStore.StoreBase(genKey("base", index, index+baseDistance), baseProof.Proof, baseProof.Wit)
		if err != nil {
			logger.Error("store base level proof error %v %v", index, err)
			return nil, err
		}
		proofs = append(proofs, baseProof.Proof)
		witnesses = append(witnesses, baseProof.Wit)
	}
	return &circuits.WitnessFile{
		Proofs:    proofs,
		Witnesses: witnesses,
	}, nil
}

func (bs *BtcSetup) middleProve(height int64, baseProof *circuits.WitnessFile) (*circuits.WitnessFile, error) {
	startIndex := height - 6*baseDistance
	if startIndex < 0 {
		return nil, fmt.Errorf("middle height less than 0")
	}
	var proofs []native_plonk.Proof
	var witnesses []witness.Witness
	for index := startIndex; index <= height; index++ {
		middleData, err := midlevelUtil.GetMidLevelProofData(bs.client, uint32(index))
		if err != nil {
			logger.Error("get middle data error: %v %v", index, err)
			return nil, err
		}
		middleProof, err := midlevel.Prove(bs.cfg.SetupDir, bs.cfg.DataDir, middleData, baseProof)
		if err != nil {
			logger.Error("middLevel prove error: %v %v", index, err)
			return nil, err
		}
		err = bs.fileStore.StoreMiddle(genKey("middle", index, index+baseDistance), middleProof.Proof, middleProof.Wit)
		if err != nil {
			logger.Error("store middle level proof error %v %v", index, err)
			return nil, err
		}
		proofs = append(proofs, middleProof.Proof)
		witnesses = append(witnesses, middleProof.Wit)
	}
	return &circuits.WitnessFile{
		Proofs:    proofs,
		Witnesses: witnesses,
	}, nil
}

func (bs *BtcSetup) uplevelProve(height int64, middleProof *circuits.WitnessFile) (*reLightCommon.Proof, error) {
	upData, err := upperlevelUtil.GetUpperLevelProofData(bs.client, uint32(height))
	if err != nil {
		logger.Error("get upLevel data : %v %v", height, err)
		return nil, err
	}
	upProof, err := upperlevel.Prove(bs.cfg.SetupDir, bs.cfg.DataDir, upData, middleProof)
	if err != nil {
		logger.Error("upLevel prove error: %v %v", height, err)
		return nil, err
	}
	return upProof, nil
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
	var runConfig RunConfig
	err = json.Unmarshal(data, &runConfig)
	if err != nil {
		return nil, err
	}
	return &runConfig, nil
}
