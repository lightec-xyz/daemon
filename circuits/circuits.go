package circuits

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/consensys/gnark-crypto/ecc"
	native_plonk "github.com/consensys/gnark/backend/plonk"
	plonk_bn254 "github.com/consensys/gnark/backend/plonk/bn254"
	"github.com/consensys/gnark/backend/witness"
	dCom "github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	ethblock "github.com/lightec-xyz/provers/circuits/fabric/tx-in-eth2"
	txineth2 "github.com/lightec-xyz/provers/circuits/tx-in-eth2"
	"github.com/lightec-xyz/reLight/circuits/common"
	"github.com/lightec-xyz/reLight/circuits/genesis"
	"github.com/lightec-xyz/reLight/circuits/recursive"
	"github.com/lightec-xyz/reLight/circuits/unit"
	"github.com/lightec-xyz/reLight/circuits/utils"
)

type Circuit struct {
	unit      *unit.Unit
	recursive *recursive.Recursive
	genesis   *genesis.Genesis
	Cfg       *CircuitConfig
	debug     bool
}

func NewCircuit(cfg *CircuitConfig) (*Circuit, error) {
	unitConfig := unit.NewUnitConfig(cfg.DataDir, cfg.SrsDir, cfg.SubDir)
	genesisConfig := genesis.NewGenesisConfig(cfg.DataDir, cfg.SrsDir, cfg.SubDir)
	recursiveConfig := recursive.NewRecursiveConfig(cfg.DataDir, cfg.SrsDir, cfg.SubDir)
	var zkDebug bool
	var err error
	zkDebugEnv := os.Getenv(dCom.ZkDebugEnv)
	if zkDebugEnv != "" {
		zkDebug, err = strconv.ParseBool(zkDebugEnv)
		if err != nil {
			return nil, err
		}
	}
	return &Circuit{
		unit:      unit.NewUnit(unitConfig),
		recursive: recursive.NewRecursive(recursiveConfig),
		genesis:   genesis.NewGenesis(genesisConfig),
		Cfg:       cfg,
		debug:     zkDebug, // todo
	}, nil
}

func (c *Circuit) Load() error {
	if c.debug {
		logger.Warn("current zk circuit is debug mode,skip load")
		return nil
	}
	// todo
	err := c.genesis.Load()
	if err != nil {
		logger.Error("genesis load error:%v", err)
		return err
	}
	err = c.unit.Load()
	if err != nil {
		logger.Error("unit load error:%v", err)
		return err
	}
	err = c.recursive.Load()
	if err != nil {
		logger.Error("recursive load error:%v", err)
		return err
	}
	return nil
}

func (c *Circuit) TxInEth2Prove(param *ethblock.TxInEth2ProofData) (*common.Proof, error) {
	proof, wit, err := txineth2.Prove(c.Cfg.DataDir, param)
	if err != nil {
		return nil, err
	}
	return &common.Proof{
		Proof: proof,
		Wit:   wit,
	}, nil
}

func (c *Circuit) TxBlockIsParentOfCheckPointProve() (*common.Proof, error) {

	return nil, nil
}

func (c *Circuit) CheckPointFinalityProve() (*common.Proof, error) {

	return nil, nil
}

func (c *Circuit) RedeemProve() (*common.Proof, error) {
	panic(c)
	return nil, nil
}

func (c *Circuit) DepositProve() (*common.Proof, error) {
	panic(c)
	return nil, nil
}

func (c *Circuit) UnitProve(period uint64, update *utils.LightClientUpdateInfo) (*common.Proof, error) {
	if c.debug {
		logger.Warn("current zk circuit unit prove is debug mode,skip prove")
		return debugProof()
	}
	// todo
	logger.Warn("really do unit prove now: %v", period)
	//proof, err := unitProve(c.Cfg.DataDir, c.Cfg.SrsDir, fmt.Sprintf("sc%d", period), update)
	//proof, err := c.unit.Prove(update)
	subDir := fmt.Sprintf("sc%d", period)
	err := innerProve(c.Cfg.DataDir, subDir, update)
	if err != nil {
		logger.Error("inner prove error:%v", err)
		return nil, err
	}
	err = outerProve(c.Cfg.DataDir, subDir, update)
	if err != nil {
		logger.Error("outer prove error:%v", err)
		return nil, err
	}
	proof, err := innerUnitProv(c.Cfg.DataDir, subDir, update)
	if err != nil {
		logger.Error("unit prove error:%v", err)
		return nil, err
	}
	return proof, nil
}

func (c *Circuit) RecursiveProve(choice string, firstProof, secondProof, firstWitness, secondWitness []byte,
	beginId, relayId, endId []byte) (*common.Proof, error) {
	logger.Debug("recursive prove request data choice:%v", choice)
	if c.debug {
		logger.Warn("current zk circuit recursive prove is debug mode,skip prove")
		return debugProof()
	}
	if !(choice == "genesis" || choice == "recursive") {
		return nil, fmt.Errorf("invalid choice: %s", choice)
	}
	firstPr, err := ParseProof(firstProof)
	if err != nil {
		logger.Error("parse proof error:%v", err)
		return nil, err
	}
	secondPr, err := ParseProof(secondProof)
	if err != nil {
		logger.Error("parse proof error:%v", err)
		return nil, err
	}
	firstWit, err := ParseWitness(firstWitness)
	if err != nil {
		logger.Error("parse witness error:%v", err)
		return nil, err
	}
	secondWit, err := ParseWitness(secondWitness)
	if err != nil {
		logger.Error("parse witness error:%v", err)
		return nil, err
	}
	return c.recursive.Prove(choice, firstPr, secondPr, firstWit, secondWit, beginId, relayId, endId)
}

func (c *Circuit) GenesisProve(firstProof, secondProof, firstWitness, secondWitness []byte,
	genesisId, firstId, secondId []byte) (*common.Proof, error) {
	//logger.Debug("genesis prove request data firstProof:%x secondProof:%x firstWitness:%x secondWitness:%x,genesisId:%x firstId:%x secondId:%x",
	//	firstProof, secondProof, firstWitness, secondWitness, genesisId, firstId, secondId)
	if c.debug {
		logger.Warn("current zk circuit genesis prove is debug mode,skip prove")
		return debugProof()
	}
	firstPf, err := ParseProof(firstProof)
	if err != nil {
		logger.Error("parse proof error:%v", err)
		return nil, err
	}
	secondPf, err := ParseProof(secondProof)
	if err != nil {
		logger.Error("parse proof error:%v", err)
		return nil, err
	}
	firstWit, err := ParseWitness(firstWitness)
	if err != nil {
		logger.Error("parse witness error:%v", err)
		return nil, err
	}
	secondWit, err := ParseWitness(secondWitness)
	if err != nil {
		logger.Error("parse witness error:%v", err)
		return nil, err
	}
	return c.genesis.Prove(firstPf, secondPf, firstWit, secondWit, genesisId, firstId, secondId)
}
func (c *Circuit) VerifyProve() (*common.Proof, error) {

	return nil, nil
}

func SyncCommitRoot(update *utils.LightClientUpdateInfo) ([]byte, error) {
	return utils.SyncCommitRoot(update)
}

func ParseProof(proof []byte) (native_plonk.Proof, error) {
	reader := bytes.NewReader(proof)
	var bn254Proof plonk_bn254.Proof
	_, err := bn254Proof.ReadFrom(reader)
	if err != nil {
		return nil, err
	}
	return &bn254Proof, nil
}

func innerUnitProv(dataDir string, subDir string, update *utils.LightClientUpdateInfo) (*common.Proof, error) {
	unitCfg := unit.NewUnitConfig(dataDir, "", subDir)
	unit := unit.NewUnit(unitCfg)
	err := unit.Load()
	if err != nil {
		return nil, err
	}
	proofs, err := unit.Prove(update)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return proofs, nil
}

func outerProve(dataDir string, subDir string, update *utils.LightClientUpdateInfo) error {
	outerCfg := unit.NewOuterConfig(dataDir, "", subDir)
	outer := unit.NewOuter(&outerCfg)
	err := outer.Load()
	if err != nil {
		return err
	}
	proofs, err := outer.Prove(update)
	if err != nil {
		return err
	}
	err = outer.Save(proofs)
	if err != nil {
		return err
	}
	return nil
}

func innerProve(dataDir string, subDir string, update *utils.LightClientUpdateInfo) error {
	innerCfg := unit.NewInnerConfig(dataDir, "", subDir)
	inner := unit.NewInner(&innerCfg)
	err := inner.Load()
	if err != nil {
		return err
	}
	assignments, err := inner.GetCircuitAssignments(update)
	if err != nil {
		return err
	}
	for index, assignment := range assignments {
		proof, err := inner.Prove(assignment)
		if err != nil {
			return err
		}
		err = inner.Save(index, proof)
		if err != nil {
			return err
		}
	}
	return nil
}

func ParseWitness(body []byte) (witness.Witness, error) {
	reader := bytes.NewReader(body)
	var wit witness.Witness
	_, err := wit.ReadFrom(reader)
	if err != nil {
		return nil, err
	}
	return wit, nil
}

type CircuitConfig struct {
	DataDir string
	SrsDir  string
	SubDir  string
	Debug   bool
}

func debugProof() (*common.Proof, error) {
	// todo
	time.Sleep(15 * time.Second)
	field := ecc.BN254.ScalarField()
	w, err := witness.New(field)
	if err != nil {
		return nil, err
	}
	return &common.Proof{
		Proof: &plonk_bn254.Proof{},
		Wit:   w,
	}, nil
}

func ProofToBytes(proof native_plonk.Proof) []byte {
	var buf bytes.Buffer
	_, err := proof.WriteTo(&buf)
	if err != nil {
		return nil
	}
	return buf.Bytes()
}

func WitnessToBytes(witness witness.Witness) []byte {
	var buf bytes.Buffer
	_, err := witness.WriteTo(&buf)
	if err != nil {
		return nil
	}
	return buf.Bytes()
}
