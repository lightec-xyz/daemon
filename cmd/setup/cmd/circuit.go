package cmd

import (
	"fmt"
	"github.com/lightec-xyz/btc_provers/circuits/baselevel"
	"github.com/lightec-xyz/btc_provers/circuits/midlevel"
	"github.com/lightec-xyz/btc_provers/circuits/upperlevel"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/reLight/circuits/genesis"
	"github.com/lightec-xyz/reLight/circuits/recursive"
	"github.com/lightec-xyz/reLight/circuits/unit"
)

type Group string

const (
	Bitcoin  Group = "bitcoin"
	Beacon   Group = "beacon"
	Ethereum Group = "ethereum"
	All      Group = "all"
)

type CircuitType string

const (
	btcBase         CircuitType = "btcBase"
	btcMiddle       CircuitType = "btcMiddle"
	btcUpper        CircuitType = "btcUpper"
	beaconInner     CircuitType = "beaconInner"
	beaconOuter     CircuitType = "beaconOuter"
	beaconUnit      CircuitType = "beaconUnit"
	beaconGenesis   CircuitType = "beaconGenesis"
	beaconRecursive CircuitType = "beaconRecursive"
)

var btcGroups = []CircuitType{btcBase, btcMiddle, btcUpper}
var beaconGroups = []CircuitType{beaconInner, beaconOuter, beaconUnit, beaconGenesis, beaconRecursive}
var ethGroups = []CircuitType{}

type CircuitSetup struct {
	datadir string
	srsdir  string
}

func NewCircuitSetup(datadir, srsdir string) *CircuitSetup {
	return &CircuitSetup{
		datadir: datadir,
		srsdir:  srsdir,
	}
}

func (cs *CircuitSetup) SetupGroup(group Group) error {
	circuitTypes, err := cs.CircuitTypes(group)
	if err != nil {
		return err
	}
	for _, circuitType := range circuitTypes {
		if err = cs.Setup(circuitType); err != nil {
			return err
		}
		logger.Info("finish setup circuit: %s", circuitType)
	}
	return nil
}

func (cs *CircuitSetup) Setup(circuitType CircuitType) error {
	logger.Info("start setup circuit: %s", circuitType)
	switch circuitType {
	case beaconOuter:
		return cs.SyncCommOuter()
	case beaconInner:
		return cs.SyncCommInner()
	case beaconUnit:
		return cs.SyncCommUnit()
	case beaconGenesis:
		return cs.SyncCommGenesis()
	case beaconRecursive:
		return cs.SyncCommRecursive()
	case btcBase:
		return cs.BtcBase()
	case btcMiddle:
		return cs.BtcMiddle()
	case btcUpper:
		return cs.BtcUpleve()
	default:
		return fmt.Errorf("invalid circuitType: %s", circuitType)
	}
}

func (cs *CircuitSetup) CircuitTypes(group Group) ([]CircuitType, error) {
	switch group {
	case Bitcoin:
		return btcGroups, nil
	case Beacon:
		return beaconGroups, nil
	case Ethereum:
		return ethGroups, nil
	case All:
		return append(beaconGroups, append(ethGroups, btcGroups...)...), nil
	default:
		return nil, fmt.Errorf("invalid group: %s", group)
	}
}

func (cs *CircuitSetup) SyncCommInner() error {
	config := unit.NewInnerConfig(cs.datadir, cs.srsdir, "")
	inner := unit.NewInner(&config)
	err := inner.Setup()
	if err != nil {
		return err
	}
	return nil
}

func (cs *CircuitSetup) SyncCommOuter() error {
	config := unit.NewOuterConfig(cs.datadir, cs.srsdir, "")
	outer := unit.NewOuter(&config)
	err := outer.Setup()
	if err != nil {
		return err
	}
	return nil
}

func (cs *CircuitSetup) SyncCommUnit() error {
	config := unit.NewUnitConfig(cs.datadir, cs.srsdir, "")
	unitCir := unit.NewUnit(config)
	err := unitCir.Setup()
	if err != nil {
		return err
	}
	return nil
}
func (cs *CircuitSetup) SyncCommGenesis() error {
	genesisConfig := genesis.NewGenesisConfig(cs.datadir, cs.srsdir, "")
	genesis := genesis.NewGenesis(genesisConfig)
	err := genesis.Setup()
	if err != nil {
		return err
	}
	return nil
}

func (cs *CircuitSetup) SyncCommRecursive() error {
	recursiveConfig := recursive.NewRecursiveConfig(cs.datadir, cs.srsdir, "")
	recursive := recursive.NewRecursive(recursiveConfig)
	err := recursive.Setup()
	if err != nil {
		return err
	}
	return nil
}
func (cs *CircuitSetup) BtcBase() error {
	err := baselevel.Setup(cs.datadir, cs.srsdir)
	if err != nil {
		return err
	}
	return nil
}

func (cs *CircuitSetup) BtcMiddle() error {
	err := midlevel.Setup(cs.datadir, cs.srsdir)
	if err != nil {
		return err
	}
	return nil
}

func (cs *CircuitSetup) BtcUpleve() error {
	err := upperlevel.Setup(cs.datadir, cs.srsdir)
	if err != nil {
		return err
	}
	return nil
}
