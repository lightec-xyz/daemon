package cmd

import (
	"fmt"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/lightec-xyz/btc_provers/circuits/blockdepth/timestamp"
	"github.com/lightec-xyz/daemon/common"
	"math/big"
	"path/filepath"

	"github.com/lightec-xyz/btc_provers/circuits/blockchain"
	"github.com/lightec-xyz/btc_provers/circuits/blockchain/baselevel"
	"github.com/lightec-xyz/btc_provers/circuits/blockchain/midlevel"
	"github.com/lightec-xyz/btc_provers/circuits/blockchain/upperlevel"
	"github.com/lightec-xyz/btc_provers/circuits/blockdepth"
	"github.com/lightec-xyz/btc_provers/circuits/txinchain"
	"github.com/lightec-xyz/daemon/logger"
	beacon_header "github.com/lightec-xyz/provers/circuits/beacon-header"
	beacon_header_finality "github.com/lightec-xyz/provers/circuits/beacon-header-finality"
	"github.com/lightec-xyz/provers/circuits/redeem"
	syncCommittee "github.com/lightec-xyz/provers/circuits/sync-committee"
	txineth2 "github.com/lightec-xyz/provers/circuits/tx-in-eth2"
	proversCom "github.com/lightec-xyz/provers/common"
)

type Group string

func (g Group) String() string {
	return string(g)
}

const (
	Bitcoin  Group = "bitcoin"
	Beacon   Group = "beacon"
	Ethereum Group = "ethereum"
	TxAll    Group = "txes"
	All      Group = "all"
)

type CircuitType string

func (c CircuitType) String() string {
	return string(c)
}

const (
	btcBase       CircuitType = "btcBase"
	btcMiddle     CircuitType = "btcMiddle"
	btcUpper      CircuitType = "btcUpper"
	btcBlockChain CircuitType = "btcBlockChain"
	btcBlockDepth CircuitType = "btcBlockDepth"
	btcTxInChain  CircuitType = "btcTxInChain"
	btcTxInBlock  CircuitType = "btcTxInBlock"
	btcTimestamp  CircuitType = "btcTimestamp"

	beaconSyncCommittee CircuitType = "beaconSyncCommittee"

	ethTxInEth2       CircuitType = "ethTxInEth2"
	ethBeaconHeader   CircuitType = "ethBeaconHeader"
	ethFinalityHeader CircuitType = "ethFinalityHeader"
	ethRedeem         CircuitType = "ethRedeem"
)

var btcGroups = []CircuitType{btcTimestamp, btcBlockDepth, btcBlockChain, btcTxInChain}
var beaconGroups = []CircuitType{beaconSyncCommittee}
var ethGroups = []CircuitType{ethTxInEth2, ethBeaconHeader, ethFinalityHeader, ethRedeem}
var txesGroups = []CircuitType{ethTxInEth2, ethRedeem, btcTxInChain}

type CircuitSetup struct {
	datadir         string
	srsdir          string
	chainId         *big.Int
	zkbtcBridgeAddr string
	icpPublickey    string
}

func NewCircuitSetup(datadir, srsdir, zkbtcBridgeAddr, icpPublickey string, chainId int) *CircuitSetup {
	return &CircuitSetup{
		datadir:         datadir,
		srsdir:          srsdir,
		chainId:         big.NewInt(int64(chainId)),
		zkbtcBridgeAddr: zkbtcBridgeAddr,
		icpPublickey:    icpPublickey,
	}
}

func (cs *CircuitSetup) SetupGroup(group Group) error {
	logger.Info("start setup group: %s", group)
	circuitTypes, err := cs.CircuitTypes(group)
	if err != nil {
		logger.Error("get circuit types error: %v", err)
		return err
	}
	for _, circuitType := range circuitTypes {
		if err = cs.Setup(circuitType); err != nil {
			logger.Error("setup circuit error: %v", err)
			return err
		}
		logger.Info("finish setup circuit: %s", circuitType)
	}
	logger.Info("finish setup group: %s", group)
	return nil
}

func (cs *CircuitSetup) Setup(circuitType CircuitType) error {
	logger.Info("start setup circuit: %s", circuitType)
	switch circuitType {
	case beaconSyncCommittee:
		return cs.SyncCommittee()
	case ethTxInEth2:
		return cs.EthTxInEth2()
	case ethBeaconHeader:
		return cs.EthBeaconHeader()
	case ethFinalityHeader:
		return cs.EthFinalityHeader()
	case ethRedeem:
		return cs.EthRedeem()
	case btcBase:
		return cs.BtcBase()
	case btcMiddle:
		return cs.BtcMiddle()
	case btcUpper:
		return cs.BtcUpleve()
	case btcBlockChain:
		return cs.BtcBlockChain()
	case btcBlockDepth:
		return cs.BtcBlockDepth()
	case btcTxInChain:
		return cs.BtcTxInChain()
	case btcTxInBlock:
		return cs.BtcTxInBlock()
	case btcTimestamp:
		return cs.BtcTimestamp()
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
	case TxAll:
		return txesGroups, nil
	case All:
		return append(beaconGroups, append(ethGroups, btcGroups...)...), nil
	default:
		return nil, fmt.Errorf("invalid group: %s", group)
	}
}

func (cs *CircuitSetup) SyncCommittee() error {
	err := syncCommittee.Setup(cs.datadir, cs.srsdir)
	if err != nil {
		logger.Error("setup unit circuit error: %v", err)
		return err
	}
	return nil
}

func (cs *CircuitSetup) EthTxInEth2() error {
	err := txineth2.Setup(cs.datadir, cs.srsdir, cs.chainId, ethCommon.HexToAddress(cs.zkbtcBridgeAddr))
	if err != nil {
		logger.Error("setup txineth2 circuit error: %v", err)
		return err
	}
	return nil
}

func (cs *CircuitSetup) EthBeaconHeader() error {
	err := beacon_header.Setup(cs.datadir, cs.srsdir)
	if err != nil {
		logger.Error("setup inner circuit error: %v", err)
		return err
	}
	return nil
}

func (cs *CircuitSetup) EthFinalityHeader() error {
	err := beacon_header_finality.Setup(cs.datadir, cs.srsdir)
	if err != nil {
		logger.Error("setup finality circuit error: %v", err)
		return err
	}
	return nil
}

func (cs *CircuitSetup) EthRedeem() error {
	err := redeem.Setup(cs.datadir, cs.srsdir)
	if err != nil {
		logger.Error("setup redeem circuit error: %v", err)
		return err
	}
	return nil
}

func (cs *CircuitSetup) BtcBase() error {
	err := baselevel.Setup(cs.datadir, cs.srsdir)
	if err != nil {
		logger.Error("setup baselevel error: %v", err)
		return err
	}
	return nil
}

func (cs *CircuitSetup) BtcMiddle() error {
	err := midlevel.Setup(cs.datadir, cs.srsdir)
	if err != nil {
		logger.Error("setup midlevel error: %v", err)
		return err
	}
	return nil
}

func (cs *CircuitSetup) BtcUpleve() error {
	err := upperlevel.Setup(cs.datadir, cs.srsdir)
	if err != nil {
		logger.Error("setup upperlevel error: %v", err)
		return err
	}
	return nil
}

func (cs *CircuitSetup) BtcTxInBlock() error {
	//err := txinblock.Setup(cs.srsdir, cs.datadir)
	//if err != nil {
	//	logger.Error("setup txinblock error: %v", err)
	//	return err
	//}
	return nil
}

func (cs *CircuitSetup) BtcTxInChain() error {
	publicKey := ethCommon.FromHex(cs.icpPublickey)
	redeemFile := filepath.Join(cs.datadir, proversCom.RedeemCcsFile)
	exists, err := common.FileExists(redeemFile)
	if err != nil {
		logger.Error("check redeem file error: %v", err)
		return err
	}
	if !exists {
		err := cs.EthRedeem()
		if err != nil {
			logger.Error("setup redeem error: %v", err)
			return err
		}
	}
	err = txinchain.Setup(cs.datadir, cs.srsdir, cs.datadir, [33]byte(publicKey))
	if err != nil {
		logger.Error("setup txinchain error: %v", err)
		return err
	}
	return err
}

func (cs *CircuitSetup) BtcBlockChain() error {
	err := blockchain.BlockChainSetup(cs.datadir, cs.srsdir)
	if err != nil {
		logger.Error("setup blockchain error: %v", err)
		return err
	}
	return err
}

func (cs *CircuitSetup) BtcBlockDepth() error {
	err := blockdepth.BlockDepthSetup(cs.datadir, cs.srsdir)
	if err != nil {
		logger.Error("setup blockdepth error: %v", err)
		return err
	}
	return err
}

func (cs *CircuitSetup) BtcTimestamp() error {
	err := timestamp.Setup(cs.datadir, cs.srsdir)
	if err != nil {
		logger.Error("setup timestamp error: %v", err)
		return err
	}
	return nil
}
