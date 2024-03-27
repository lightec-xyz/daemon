package circuits

import (
	"bytes"
	"encoding/hex"
	"fmt"
	native_plonk "github.com/consensys/gnark/backend/plonk"
	plonk_bn254 "github.com/consensys/gnark/backend/plonk/bn254"
	"github.com/consensys/gnark/backend/witness"
	"github.com/lightec-xyz/daemon/logger"
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
}

func NewCircuit(cfg *CircuitConfig) (*Circuit, error) {
	outerFpBytes, err := hex.DecodeString(cfg.OuterFp)
	if err != nil {
		return nil, err
	}
	unitFpBytes, err := hex.DecodeString(cfg.UnitFp)
	if err != nil {
		return nil, err
	}
	genesisFpBytes, err := hex.DecodeString(cfg.GenesisFp)
	if err != nil {
		return nil, err
	}
	unitConfig := unit.NewUnitConfig(cfg.DataDir, cfg.SrsDir, cfg.SubDir, outerFpBytes)
	genesisConfig := genesis.NewGenesisConfig(cfg.DataDir, cfg.SrsDir, cfg.SubDir, unitFpBytes)
	recursiveConfig := recursive.NewRecursiveConfig(cfg.DataDir, cfg.SrsDir, cfg.SubDir, unitFpBytes, genesisFpBytes)
	return &Circuit{
		unit:      unit.NewUnit(&unitConfig),
		recursive: recursive.NewRecursive(recursiveConfig),
		genesis:   genesis.NewGenesis(genesisConfig),
		Cfg:       cfg,
	}, nil
}

func (c *Circuit) Load() error {
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
	proof, err := c.unit.Prove(update)
	if err != nil {
		logger.Error("unit prove error:%v", err)
		return nil, err
	}
	return proof, nil
}

func (c *Circuit) RecursiveProve(choice string, firstProof, secondProof, firstWitness, secondWitness []byte,
	beginId, relayId, endId, recursiveFp []byte) (*common.Proof, error) {
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
	return c.recursive.Prove(choice, firstPr, secondPr, firstWit, secondWit, beginId, relayId, endId, recursiveFp)
}

func (c *Circuit) GenesisProve(firstProof, secondProof, firstWitness, secondWitness []byte,
	genesisId, firstId, secondId, recursiveFp []byte) (*common.Proof, error) {
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
	return c.genesis.Prove(firstPf, secondPf, firstWit, secondWit, genesisId, firstId, secondId, recursiveFp)
}

func WriteProof(proofFile string, proof native_plonk.Proof) error {
	return utils.WriteProof(proofFile, proof)
}

func WriteWitness(witnessFile string, witness witness.Witness) error {
	return utils.WriteWintess(witnessFile, witness)
}

func SyncCommitRoot(update *utils.LightClientUpdateInfo) ([]byte, error) {
	return unit.SyncCommitRoot(update)
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
	DataDir   string
	SrsDir    string
	SubDir    string
	OuterFp   string
	UnitFp    string
	RecFp     string
	GenesisFp string
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
