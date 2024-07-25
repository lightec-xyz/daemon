package circuits

import (
	"fmt"
	"path/filepath"

	"github.com/lightec-xyz/btc_provers/circuits/baselevel"
	btcprovercom "github.com/lightec-xyz/btc_provers/circuits/common"
	"github.com/lightec-xyz/btc_provers/circuits/grandrollup"
	"github.com/lightec-xyz/btc_provers/circuits/header-to-latest/bulk"
	"github.com/lightec-xyz/btc_provers/circuits/header-to-latest/packed"
	"github.com/lightec-xyz/btc_provers/circuits/header-to-latest/wrap"
	"github.com/lightec-xyz/btc_provers/circuits/midlevel"
	btcprovertypes "github.com/lightec-xyz/btc_provers/circuits/types"
	"github.com/lightec-xyz/btc_provers/circuits/upperlevel"
	btcbase "github.com/lightec-xyz/btc_provers/utils/baselevel"
	grUtil "github.com/lightec-xyz/btc_provers/utils/grandrollup"
	btcmiddle "github.com/lightec-xyz/btc_provers/utils/midlevel"
	btcupper "github.com/lightec-xyz/btc_provers/utils/upperlevel"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	beacon_header "github.com/lightec-xyz/provers/circuits/beacon-header"
	beacon_header_finality "github.com/lightec-xyz/provers/circuits/beacon-header-finality"
	ethblock "github.com/lightec-xyz/provers/circuits/fabric/tx-in-eth2"
	"github.com/lightec-xyz/provers/circuits/redeem"
	txineth2 "github.com/lightec-xyz/provers/circuits/tx-in-eth2"
	proverType "github.com/lightec-xyz/provers/circuits/types"
	reLightCommon "github.com/lightec-xyz/reLight/circuits/common"
	"github.com/lightec-xyz/reLight/circuits/genesis"
	"github.com/lightec-xyz/reLight/circuits/recursive"
	"github.com/lightec-xyz/reLight/circuits/unit"
	"github.com/lightec-xyz/reLight/circuits/utils"
)

var _ ICircuit = (*Circuit)(nil)

type Circuit struct {
	Cfg   *CircuitConfig
	debug bool
}

func (c *Circuit) BtcBaseProve(req *btcbase.BaseLevelProofData) (*reLightCommon.Proof, error) {
	logger.Debug("current zk circuit BtcBaseProve....")
	if c.debug {
		logger.Warn("current zk circuit btcBase prove is debug,skip prove")
		return debugProof()
	}
	baseProof, err := baselevel.Prove(c.Cfg.SetupDir, req)
	if err != nil {
		logger.Error("btcBase prove error: %v %v", req.FirstBlockHash, err)
		return nil, err
	}
	return baseProof, nil
}

func (c *Circuit) BtcMiddleProve(req *btcmiddle.MidLevelProofData, proofList []reLightCommon.Proof) (*reLightCommon.Proof, error) {
	logger.Debug("current zk circuit btcMiddle prove....")
	if c.debug {
		logger.Warn("current zk circuit btcMiddle prove is debug,skip prove")
		return debugProof()
	}
	middleProof, err := midlevel.Prove(c.Cfg.SetupDir, req, proofList)
	if err != nil {
		logger.Error("btcMiddle prove error: %v %v", req.FirstBlockHash, err)
		return nil, err
	}
	return middleProof, nil
}

func (c *Circuit) BtcUpperProve(req *btcupper.UpperLevelProofData, proofList []reLightCommon.Proof) (*reLightCommon.Proof, error) {
	logger.Debug("current zk circuit btcUpper prove....")
	if c.debug {
		logger.Warn("current zk circuit btcUpper prove is debug,skip prove")
		return debugProof()
	}
	upProof, err := upperlevel.Prove(c.Cfg.SetupDir, req, proofList)
	if err != nil {
		logger.Error("btcUpper prove error: %v %v", req.FirstBlockHash, err)
		return nil, err
	}
	return upProof, nil
}

func (c *Circuit) BtcBulkProve(data *btcprovertypes.BlockHeaderChain) (*reLightCommon.Proof, error) {
	logger.Debug("current zk circuit BtcBulkProve")
	err := data.Verify()
	if err != nil {
		logger.Error("verify error:%v", err)
		return nil, err
	}
	if c.debug {
		logger.Warn("current zk circuit BtcBulkProve prove is debug,skip prove")
		return debugProof()
	}
	proof, err := bulk.Prove(c.Cfg.DataDir, data)
	if err != nil {
		logger.Error("btc bulk prove error:%v", err)
		return nil, err
	}
	return proof, nil
}

func (c *Circuit) BtcPackProve(data *btcprovertypes.BlockHeaderChain) (*reLightCommon.Proof, error) {
	logger.Debug("current zk circuit BtcPackedRequest")
	if c.debug {
		logger.Warn("current zk circuit btcPack prove is debug,skip prove")
		return debugProof()
	}
	proof, err := packed.Prove(c.Cfg.DataDir, data)
	if err != nil {
		logger.Error("btc pack prove error:%v", err)
		return nil, err
	}
	return proof, nil
}

func (c *Circuit) BtcWrapProve(flag, hexProof, hexWitness, beginHash, endHash string, nbBlocks uint64) (*reLightCommon.Proof, error) {
	logger.Debug("current zk circuit BtcWrapProve")
	if c.debug {
		logger.Warn("current zk circuit BtcWrapProve prove is debug,skip prove")
		return debugProof()
	}
	proof, err := HexToProof(hexProof)
	if err != nil {
		logger.Error("parse proof error:%v", err)
		return nil, err
	}
	witness, err := HexToWitness(hexWitness)
	if err != nil {
		logger.Error("parse witness error:%v", err)
		return nil, err
	}
	beginId, err := HashToLinkageId(beginHash)
	if err != nil {
		logger.Error("parse beginId error:%v", err)
		return nil, err
	}
	endId, err := HashToLinkageId(endHash)
	if err != nil {
		logger.Error("parse endId error:%v", err)
		return nil, err
	}
	var vkPath string
	if flag == BtcBulk {
		vkPath = filepath.Join(c.Cfg.SetupDir, btcprovercom.BlockHeaderBulkOuterVkFile)
	} else if flag == BtcPacked {
		vkPath = filepath.Join(c.Cfg.SetupDir, btcprovercom.BlockHeaderPackedVkFile)
	} else {
		return nil, fmt.Errorf("unknown flag")
	}
	verifyingKey, err := utils.ReadVk(vkPath)
	if err != nil {
		logger.Error("read vk error:%v", err)
		return nil, err
	}
	result, err := wrap.Prove(c.Cfg.DataDir, verifyingKey, proof, witness, beginId, endId, nbBlocks)
	if err != nil {
		logger.Error("btc wrap prove error:%v", err)
		return nil, err
	}
	return result, nil
}

func (c *Circuit) RedeemProve(txProof, txWitness, bhProof, bhWitness, bhfProof, bhfWitness string, beginId, endId, genesisScRoot,
	currentSCSSZRoot string, txVar, receiptVar []string) (*reLightCommon.Proof, error) {
	logger.Debug("current zk circuit redeemProve")
	if c.debug {
		logger.Warn("current zk circuit redeemProve prove is debug,skip prove")
		return debugProof()
	}
	txInEth2Proof, err := HexToProof(txProof)
	if err != nil {
		logger.Error("parse proof error:%v", err)
		return nil, err
	}
	txInEth2Witness, err := HexToWitness(txWitness)
	if err != nil {
		logger.Error("parse witness error:%v", err)
		return nil, err
	}

	blockHeaderProof, err := HexToProof(bhProof)
	if err != nil {
		logger.Error("parse proof error:%v", err)
		return nil, err
	}
	blockHeaderWitness, err := HexToWitness(bhWitness)
	if err != nil {
		logger.Error("parse witness error:%v", err)
		return nil, err
	}
	blockHeaderFinalityProof, err := HexToProof(bhfProof)
	if err != nil {
		logger.Error("parse proof error:%v", err)
		return nil, err
	}

	blockHeaderFinalityWitness, err := HexToWitness(bhfWitness)
	if err != nil {
		logger.Error("parse witness error:%v", err)
		return nil, err
	}
	genesisSCSSZRoot, err := GetGenesisSCSSZRoot(genesisScRoot)
	if err != nil {
		logger.Error("get genesis scssz root error:%v", err)
		return nil, err
	}
	beginIdBytes, err := HexToBytes(beginId)
	if err != nil {
		logger.Error("decode begin id error:%v", err)
		return nil, err
	}
	endIdBytes, err := HexToBytes(endId)
	if err != nil {
		logger.Error("decode begin id error:%v", err)
		return nil, err
	}
	curScRootBytes, err := HexToBytes(currentSCSSZRoot)
	if err != nil {
		logger.Error("decode current scssz root error:%v", err)
		return nil, err
	}
	txValue, err := common.HexToTxVar(txVar)
	if err != nil {
		logger.Error("decode tx value error:%v", err)
		return nil, err
	}
	hexToReceiptVar, err := common.HexToReceiptVar(receiptVar)
	if err != nil {
		logger.Error("decode receipt value error:%v", err)
		return nil, err
	}
	proof, err := redeem.Prove(c.Cfg.DataDir, txInEth2Proof, txInEth2Witness, blockHeaderProof, blockHeaderWitness,
		blockHeaderFinalityProof, blockHeaderFinalityWitness, beginIdBytes, endIdBytes, genesisSCSSZRoot, curScRootBytes, *txValue,
		*hexToReceiptVar)
	if err != nil {
		logger.Error("redeem prove error:%v", err)
		return nil, err
	}
	return &reLightCommon.Proof{
		Proof: proof.Proof,
		Wit:   proof.Wit,
	}, nil
}

func (c *Circuit) BeaconHeaderFinalityUpdateProve(genesisSCSSZRoot string, recursiveProof, recursiveWitness, outerProof,
	outerWitness string, finalityUpdate *proverType.FinalityUpdate, scUpdate *proverType.SyncCommitteeUpdate) (*reLightCommon.Proof, error) {
	ok, err := common.VerifyLightClientUpdate(scUpdate)
	if err != nil {
		logger.Error("verify light client update error:%v", err)
		return nil, err
	}
	if !ok {
		logger.Error("verify light client update error")
		return nil, fmt.Errorf("verify light client update error")
	}
	logger.Debug("current zk circuit BeaconHeaderFinalityUpdateProve")
	if c.debug {
		logger.Warn("current zk circuit BeaconHeaderFinalityUpdateProve prove is debug,skip prove")
		return debugProof()
	}
	scRecursiveProof, err := HexToProof(recursiveProof)
	if err != nil {
		logger.Error("parse proof error:%v", err)
		return nil, err
	}
	scRecursiveWitness, err := HexToWitness(recursiveWitness)
	if err != nil {
		logger.Error("parse witness error:%v", err)
		return nil, err
	}
	scOuterProof, err := HexToProof(outerProof)
	if err != nil {
		logger.Error("parse proof error:%v", err)
		return nil, err
	}
	scOuterWitness, err := HexToWitness(outerWitness)
	if err != nil {
		logger.Error("parse witness error:%v", err)
		return nil, err
	}

	proof, err := beacon_header_finality.Prove(c.Cfg.DataDir, genesisSCSSZRoot, scRecursiveProof,
		scRecursiveWitness, scOuterProof, scOuterWitness, finalityUpdate, scUpdate)
	if err != nil {
		logger.Error("beacon header finality update prove error:%v", err)
		return nil, err
	}
	return &reLightCommon.Proof{
		Proof: proof.Proof,
		Wit:   proof.Wit,
	}, nil
}

func (c *Circuit) BeaconHeaderProve(header proverType.BeaconHeaderChain) (*reLightCommon.Proof, error) {
	logger.Debug("current zk circuit BeaconHeaderProve")
	if c.debug {
		logger.Warn("current zk circuit BeaconHeaderProve prove is debug,skip prove")
		return debugProof()
	}
	proof, err := beacon_header.Prove(c.Cfg.DataDir, header)
	if err != nil {
		logger.Error("beacon header prove error:%v %v %v", header.BeginSlot, header.EndSlot, err)
		return nil, err
	}
	return &reLightCommon.Proof{
		Proof: proof.Proof,
		Wit:   proof.Wit,
	}, nil
}

func (c *Circuit) TxInEth2Prove(param *ethblock.TxInEth2ProofData) (*reLightCommon.Proof, error) {
	logger.Debug("current zk circuit TxInEth2Prove")
	if c.debug {
		logger.Warn("current zk circuit TxInEth2Prove prove is debug,skip prove")
		return debugProof()
	}
	proof, err := txineth2.Prove(c.Cfg.DataDir, param)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return &reLightCommon.Proof{
		Proof: proof.Proof,
		Wit:   proof.Wit,
	}, err
}
func (c *Circuit) DepositProve(data *grUtil.GrandRollupProofData) (*reLightCommon.Proof, error) {
	logger.Debug("current zk circuit DepositProve")
	if c.debug {
		logger.Warn("current zk circuit DepositProve is debug,skip prove ")
		return debugProof()
	}
	proof, _, err := grandrollup.Prove(c.Cfg.DataDir, data)
	if err != nil {
		logger.Error("deposit prove error:%v", err)
		return nil, err
	}
	return proof, nil

}
func (c *Circuit) UnitProve(period uint64, update *utils.SyncCommitteeUpdate) (*reLightCommon.Proof, *reLightCommon.Proof, error) {
	logger.Debug("current zk circuit unit prove")
	ok, err := common.VerifyLightClientUpdate(update)
	if err != nil {
		logger.Error("verify light client update error:%v", err)
		return nil, nil, err
	}
	if !ok {
		logger.Error("verify light client update error")
		return nil, nil, fmt.Errorf("verify light client update error")
	}
	if c.debug {
		logger.Warn("current zk circuit unit prove is debug mode,skip prove")
		proof, err := debugProof()
		if err != nil {
			logger.Error("debug proof error:%v", err)
			return nil, nil, err
		}
		return proof, proof, err
	}
	subDir := fmt.Sprintf("sc%d", period) // todo need remove
	err = innerProve(c.Cfg.DataDir, subDir, update)
	if err != nil {
		logger.Error("inner prove error:%v", err)
		return nil, nil, err
	}
	outerProof, err := outerProve(c.Cfg.DataDir, subDir, update)
	if err != nil {
		logger.Error("outer prove error:%v", err)
		return nil, nil, err
	}
	unitProof, err := unitProv(c.Cfg.DataDir, subDir, update)
	if err != nil {
		logger.Error("unit prove error:%v", err)
		return nil, nil, err
	}
	return unitProof, outerProof, nil
}

func (c *Circuit) RecursiveProve(choice string, firstProof, secondProof, firstWitness, secondWitness string,
	beginId, relayId, endId []byte) (*reLightCommon.Proof, error) {
	logger.Debug("current zk circuit recursive prove,choice: %v", choice)
	if c.debug {
		logger.Warn("current zk circuit recursive prove is debug mode,skip prove")
		return debugProof()
	}
	if !(choice == "genesis" || choice == "recursive") { // todo
		return nil, fmt.Errorf("invalid choice: %s", choice)
	}
	firstPr, err := HexToProof(firstProof)
	if err != nil {
		logger.Error("parse proof error:%v", err)
		return nil, err
	}
	secondPr, err := HexToProof(secondProof)
	if err != nil {
		logger.Error("parse proof error:%v", err)
		return nil, err
	}
	firstWit, err := HexToWitness(firstWitness)
	if err != nil {
		logger.Error("parse witness error:%v", err)
		return nil, err
	}
	secondWit, err := HexToWitness(secondWitness)
	if err != nil {
		logger.Error("parse witness error:%v", err)
		return nil, err
	}
	config := recursive.NewRecursiveConfig(c.Cfg.DataDir, c.Cfg.SrsDir, "")
	recursiveCir := recursive.NewRecursive(config)
	err = recursiveCir.Load()
	if err != nil {
		logger.Error("recursive load error:%v", err)
		return nil, err
	}
	proof, err := recursiveCir.Prove(choice, firstPr, secondPr, firstWit, secondWit, beginId, relayId, endId)
	if err != nil {
		logger.Error("recursive prove error:%v", err)
		return nil, err
	}
	return proof, err
}

func (c *Circuit) GenesisProve(firstProof, secondProof, firstWitness, secondWitness string,
	genesisId, firstId, secondId []byte) (*reLightCommon.Proof, error) {
	logger.Debug("current zk circuit syncCommittee genesis prove")
	if c.debug {
		logger.Warn("current zk circuit genesis prove is debug mode,skip prove")
		return debugProof()
	}
	firstPf, err := HexToProof(firstProof)
	if err != nil {
		logger.Error("parse proof error:%v", err)
		return nil, err
	}
	secondPf, err := HexToProof(secondProof)
	if err != nil {
		logger.Error("parse proof error:%v", err)
		return nil, err
	}
	firstWit, err := HexToWitness(firstWitness)
	if err != nil {
		logger.Error("parse witness error:%v", err)
		return nil, err
	}
	secondWit, err := HexToWitness(secondWitness)
	if err != nil {
		logger.Error("parse witness error:%v", err)
		return nil, err
	}
	config := genesis.NewGenesisConfig(c.Cfg.DataDir, "", "")
	genesisCir := genesis.NewGenesis(config)
	err = genesisCir.Load()
	if err != nil {
		logger.Error("genesis load error:%v", err)
		return nil, err
	}
	proof, err := genesisCir.Prove(firstPf, secondPf, firstWit, secondWit, genesisId, firstId, secondId)
	if err != nil {
		logger.Error("genesis prove error:%v", err)
		return nil, err
	}
	return proof, err
}

func (c *Circuit) UpdateChangeProve(data *grUtil.GrandRollupProofData) (*reLightCommon.Proof, error) {
	logger.Debug("current zk circuit UpdateChangeProve")
	if c.debug {
		logger.Warn("current zk circuit DepositProve is debug,skip prove ")
		return debugProof()
	}
	proof, _, err := grandrollup.Prove(c.Cfg.DataDir, data)
	if err != nil {
		logger.Error("update change prove error:%v", err)
		return nil, err
	}
	return proof, nil
}

func NewCircuit(cfg *CircuitConfig) (*Circuit, error) {
	return &Circuit{
		Cfg:   cfg,
		debug: cfg.Debug,
	}, nil
}

func unitProv(dataDir string, subDir string, update *utils.SyncCommitteeUpdate) (*reLightCommon.Proof, error) {
	unitCfg := unit.NewUnitConfig(dataDir, "", subDir)
	unit := unit.NewUnit(unitCfg)
	err := unit.Load()
	if err != nil {
		logger.Error("load unit error:%v", err)
		return nil, err
	}
	proofs, err := unit.Prove(update)
	if err != nil {
		logger.Error("unit prove error:%v", err)
		return nil, err
	}
	return proofs, nil
}

func outerProve(dataDir string, subDir string, update *utils.SyncCommitteeUpdate) (*reLightCommon.Proof, error) {
	outerCfg := unit.NewOuterConfig(dataDir, "", subDir)
	outer := unit.NewOuter(&outerCfg)
	err := outer.Load()
	if err != nil {
		logger.Error("load outer error:%v", err)
		return nil, err
	}
	proofs, err := outer.Prove(update)
	if err != nil {
		logger.Error("outer prove error:%v", err)
		return nil, err
	}
	err = outer.Save(proofs)
	if err != nil {
		logger.Error("outer save error:%v", err)
		return nil, err
	}
	return proofs, nil
}

func innerProve(dataDir string, subDir string, update *utils.SyncCommitteeUpdate) error {
	innerCfg := unit.NewInnerConfig(dataDir, "", subDir)
	inner := unit.NewInner(&innerCfg)
	err := inner.Load()
	if err != nil {
		logger.Error("load inner error:%v", err)
		return err
	}
	assignments, err := inner.GetCircuitAssignments(update)
	if err != nil {
		logger.Error("get circuit assignments error:%v", err)
		return err
	}
	for index, assignment := range assignments {
		proof, err := inner.Prove(assignment)
		if err != nil {
			logger.Error("prove error:%v", err)
			return err
		}
		err = inner.Save(proof, index)
		if err != nil {
			logger.Error("save error:%v", err)
			return err
		}
	}
	return nil
}
