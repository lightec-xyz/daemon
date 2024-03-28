package circuits

import (
	"bytes"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	native_plonk "github.com/consensys/gnark/backend/plonk"
	plonk_bn254 "github.com/consensys/gnark/backend/plonk/bn254"
	"github.com/consensys/gnark/backend/witness"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/reLight/circuits/common"
	"github.com/lightec-xyz/reLight/circuits/genesis"
	"github.com/lightec-xyz/reLight/circuits/recursive"
	"github.com/lightec-xyz/reLight/circuits/unit"
	"github.com/lightec-xyz/reLight/circuits/utils"
	"time"
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
	return &Circuit{
		unit:      unit.NewUnit(unitConfig),
		recursive: recursive.NewRecursive(recursiveConfig),
		genesis:   genesis.NewGenesis(genesisConfig),
		Cfg:       cfg,
		debug:     true, // todo
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

func (c *Circuit) UnitProve(update *utils.LightClientUpdateInfo) (*common.Proof, error) {
	if c.debug {
		logger.Warn("current zk circuit unit prove is debug mode,skip prove")
		return debugProof()
	}

	proof, err := c.unit.Prove(update)
	if err != nil {
		logger.Error("unit prove error:%v", err)
		return nil, err
	}
	return proof, nil
}

func (c *Circuit) RecursiveProve(choice string, firstProof, secondProof, firstWitness, secondWitness []byte,
	beginId, relayId, endId []byte) (*common.Proof, error) {
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

func WriteProof(proofFile string, proof native_plonk.Proof) error {
	return utils.WriteProof(proofFile, proof)
}

func WriteWitness(witnessFile string, witness witness.Witness) error {
	return utils.WriteWintess(witnessFile, witness)
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
