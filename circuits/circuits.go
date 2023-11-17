package circuits

import (
	"fmt"
	nativeplonk "github.com/consensys/gnark/backend/plonk"
	"github.com/consensys/gnark/backend/witness"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/lightec-xyz/btc_provers/circuits/blockchain"
	"github.com/lightec-xyz/btc_provers/circuits/blockchain/baselevel"
	"github.com/lightec-xyz/btc_provers/circuits/blockchain/midlevel"
	"github.com/lightec-xyz/btc_provers/circuits/blockchain/upperlevel"
	"github.com/lightec-xyz/btc_provers/circuits/blockdepth/blockbulk"
	depthCommon "github.com/lightec-xyz/btc_provers/circuits/blockdepth/common"
	"github.com/lightec-xyz/btc_provers/circuits/blockdepth/recursivebulks"
	"github.com/lightec-xyz/btc_provers/circuits/blockdepth/timestamp"
	"github.com/lightec-xyz/btc_provers/circuits/txinchain"
	blockCu "github.com/lightec-xyz/btc_provers/utils/blockchain"
	btcbase "github.com/lightec-xyz/btc_provers/utils/blockchain"
	btcmiddle "github.com/lightec-xyz/btc_provers/utils/blockchain"
	btcupper "github.com/lightec-xyz/btc_provers/utils/blockchain"
	recursiveUtil "github.com/lightec-xyz/btc_provers/utils/blockchain"
	blockDepthUtil "github.com/lightec-xyz/btc_provers/utils/blockdepth"
	grUtil "github.com/lightec-xyz/btc_provers/utils/txinchain"
	"github.com/lightec-xyz/common/operations"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	beaconheader "github.com/lightec-xyz/provers/circuits/beacon-header"
	beaconheaderfinality "github.com/lightec-xyz/provers/circuits/beacon-header-finality"
	"github.com/lightec-xyz/provers/circuits/redeem"
	syncCommittee "github.com/lightec-xyz/provers/circuits/sync-committee"
	txineth2 "github.com/lightec-xyz/provers/circuits/tx-in-eth2"
	proverType "github.com/lightec-xyz/provers/circuits/types"
	proversCommon "github.com/lightec-xyz/provers/common"
	txineth2Utils "github.com/lightec-xyz/provers/utils/tx-in-eth2"
)

var _ ICircuit = (*Circuit)(nil)

type Circuit struct {
	cfg   *CircuitConfig
	debug bool
}

func (c *Circuit) GetBtcChainVerifyKey(step uint64) nativeplonk.VerifyingKey {
	var lastVk nativeplonk.VerifyingKey
	if step == common.BtcUpperDistance || step == common.BtcBaseDistance {
		lastVk = blockchain.GetVKey(c.cfg.BtcSetupDir)
	} else {
		lastVk = blockchain.GetHybridVKey(c.cfg.BtcSetupDir, int(step))
	}
	return lastVk
}

func (c *Circuit) BtcChainHybridProve(firstType string, firstStep uint64, data *blockCu.HybridProofData, first *operations.Proof) (*operations.Proof, error) {
	logger.Debug("current zk circuit btcChainHybridProve....")
	if c.debug {
		logger.Warn("current zk circuit btcChainHybridProve prove is debug,skip prove")
		return debugProof()
	}
	verifyKey, err := c.getBtcChainVerifyKey(firstType, firstStep)
	if err != nil {
		logger.Error("getBtcChainVerifyKey error: %v", err)
		return nil, err
	}
	proof, err := blockchain.ProveHybrid(c.cfg.BtcSetupDir, first, verifyKey, data)
	if err != nil {
		logger.Error("btcChainHybrid prove error: %v", err)
		return nil, err
	}
	return proof, nil

}

func (c *Circuit) getBtcChainVerifyKey(proveType string, step uint64) (nativeplonk.VerifyingKey, error) {
	if proveType == BtcChainUpper {
		return upperlevel.GetVKey(c.cfg.BtcSetupDir), nil
	} else if proveType == BtcBlockChain {
		return blockchain.GetVKey(c.cfg.BtcSetupDir), nil
	} else if proveType == BtcChainBase {
		return baselevel.GetVKey(c.cfg.BtcSetupDir), nil
	} else if proveType == BtcChainHybrid {
		return blockchain.GetHybridVKey(c.cfg.BtcSetupDir, int(step)), nil
	} else {
		return nil, fmt.Errorf("invalid proveType:%s", proveType)
	}

}

func (c *Circuit) BtcChainRecursiveProve(firstType, secondType string, firstStep, secondStep uint64, data *recursiveUtil.BlockChainProofData, first, second *operations.Proof) (*operations.Proof, error) {
	logger.Debug("current zk circuit btcChainRecursive....")
	if c.debug {
		logger.Warn("current zk circuit btcChainRecursive prove is debug,skip prove")
		return debugProof()
	}
	firstVk, err := c.getBtcChainVerifyKey(firstType, firstStep)
	if err != nil {
		logger.Error("get verify vk error: %v", err)
		return nil, err
	}
	secondVk, err := c.getBtcChainVerifyKey(secondType, secondStep)
	if err != nil {
		logger.Error("get verify vk error: %v", err)
		return nil, err
	}
	proof, err := blockchain.ProveChain(c.cfg.BtcSetupDir, first, second, firstVk, secondVk, data)
	if err != nil {
		logger.Error("btcChainRecursive prove error: %v", err)
		return nil, err
	}
	return proof, nil
}

func (c *Circuit) BtcDepthRecursiveProve(recursive bool, step uint64, data *blockDepthUtil.RecursiveBulksProofData, first *operations.Proof) (*operations.Proof, error) {
	logger.Debug("current zk circuit btcDepthRecursive....")
	if c.debug {
		logger.Warn("current zk circuit btcDepthRecursive prove is debug,skip prove")
		return debugProof()
	}
	depthProof := &depthCommon.DepthProof{
		Proof:             first,
		IsRecursive:       recursive,
		LastAbsorbedDepth: uint32(step),
	}
	proof, err := recursivebulks.Prove(c.cfg.BtcSetupDir, depthProof, data)
	if err != nil {
		logger.Error("btcDepthRecursive prove error: %v", err)
		return nil, err
	}
	return proof, nil
}

func (c *Circuit) BtcDepositProve(chainType string, chainStep, txStep, cpStep uint64, txRecursive, cpRecursive bool, data *grUtil.TxInChainProofData, blockChain, txDepth, cpDepth, sigVerify *operations.Proof,
	proverAddr ethCommon.Address, smoothedTimestamp uint32, cpFlag uint8, sigVerifyData *blockDepthUtil.SigVerifProofData) (*operations.Proof, error) {
	logger.Debug("chainType: %v,chainStep: %v,txStep: %v,cpStep: %v,txRecursive: %v,cpRecursive: %v,data: %v,blockChain: %v,txDepth: %v,cpDepth: %v,proverAddr: %x",
		chainType, chainStep, txStep, cpStep, txRecursive, cpRecursive, data, blockChain, txDepth, cpDepth, proverAddr)
	logger.Debug("current zk circuit DepositProve....")
	if c.debug {
		logger.Warn("current zk circuit DepositProve prove is debug,skip prove")
		return debugProof()
	}
	printProof(blockChain, "blockChain")
	printProof(txDepth, "txDepth")
	printProof(cpDepth, "cpDepth")
	verifyKey, err := c.getBtcChainVerifyKey(chainType, chainStep)
	if err != nil {
		logger.Error("getBtcChainVerifyKey error: %v", err)
		return nil, err
	}
	txDepthParam := &depthCommon.DepthProof{
		Proof:             txDepth,
		IsRecursive:       txRecursive,
		LastAbsorbedDepth: uint32(txStep),
	}
	cpDepthParam := &depthCommon.DepthProof{
		Proof:             cpDepth,
		IsRecursive:       cpRecursive,
		LastAbsorbedDepth: uint32(cpStep),
	}
	proof, err := txinchain.DepositProve(c.cfg.BtcSetupDir, verifyKey, blockChain, txDepthParam, cpDepthParam, sigVerify, *sigVerifyData,
		data, proverAddr[:], smoothedTimestamp, cpFlag)
	if err != nil {
		logger.Error("DepositProve prove error: %v", err)
		return nil, err
	}
	return proof, nil
}

func (c *Circuit) BtcRedeemProve(chainType string, chainStep, txStep, cpStep uint64, txRecursive, cpRecursive bool, data *grUtil.TxInChainProofData, chain, txDepth, cpDepth, redeem, sigVerify *operations.Proof,
	minerReward [32]byte, proverAddr ethCommon.Address, smoothedTimestamp uint32, cpFlag uint8, sigVerifyData *blockDepthUtil.SigVerifProofData) (*operations.Proof, error) {
	logger.Debug("current zk circuit btcRedeemProve....")
	if c.debug {
		logger.Warn("current zk circuit btcRedeemProve prove is debug,skip prove")
		return debugProof()
	}
	//logger.Debug("chainStep:%v,txStep:%v,cpStep:%v,txRecursive:%v,cpRecursive:%v,data:%v,chain:%v,txDepth:%v,cpDepth:%v,redeem:%v,signature:%x,proverAddr:%x,minerReward:%x",
	//	chainStep, txStep, cpStep, txRecursive, cpRecursive, data, chain, txDepth, cpDepth, redeem, signature, proverAddr, minerReward)
	verifyKey, err := c.getBtcChainVerifyKey(chainType, chainStep)
	if err != nil {
		logger.Error("getBtcChainVerifyKey error: %v", err)
		return nil, err
	}
	txDepthParam := &depthCommon.DepthProof{
		Proof:             txDepth,
		IsRecursive:       txRecursive,
		LastAbsorbedDepth: uint32(txStep),
	}
	cpDepthParam := &depthCommon.DepthProof{
		Proof:             cpDepth,
		IsRecursive:       cpRecursive,
		LastAbsorbedDepth: uint32(cpStep),
	}
	proof, err := txinchain.RedeemProve(c.cfg.BtcSetupDir, c.cfg.EthSetupDir, chain, redeem, txDepthParam,
		cpDepthParam, sigVerify, verifyKey, *sigVerifyData, data, proverAddr[:], smoothedTimestamp, cpFlag, minerReward)
	if err != nil {
		logger.Error("btcRedeemProve prove error: %v", err)
		return nil, err
	}
	return proof, nil
}

func (c *Circuit) BtcTimestamp(cpTime *blockDepthUtil.CptimestampProofData, smoothData *blockDepthUtil.SmoothedTimestampProofData) (*operations.Proof, error) {
	logger.Debug("current zk circuit btcTimestamp....")
	if c.debug {
		logger.Warn("current zk circuit btcTimestamp prove is debug,skip prove")
		return debugProof()
	}
	proof, err := timestamp.Prove(c.cfg.BtcSetupDir, cpTime, smoothData)
	if err != nil {
		logger.Error("btcTimestamp prove error: %v", err)
		return nil, err
	}
	return proof, nil
}

func (c *Circuit) BtcBaseProve(data *btcbase.BaseLevelProofData) (*operations.Proof, error) {
	logger.Debug("current zk circuit btcBaseProve....")
	if c.debug {
		logger.Warn("current zk circuit btcBase prove is debug,skip prove")
		return debugProof()
	}
	proof, err := baselevel.Prove(c.cfg.BtcSetupDir, data)
	if err != nil {
		logger.Error("btcBase prove error: %v", err)
		return nil, err
	}
	return proof, nil

}

func (c *Circuit) BtcMiddleProve(data *btcmiddle.BatchedProofData, proofs []operations.Proof) (*operations.Proof, error) {
	logger.Debug("current zk circuit btcMiddle prove....")
	if c.debug {
		logger.Warn("current zk circuit btcMiddle prove is debug,skip prove")
		return debugProof()
	}
	proof, err := midlevel.Prove(c.cfg.BtcSetupDir, data, proofs)
	if err != nil {
		logger.Error("btcMiddle prove error: %v", err)
		return nil, err
	}
	return proof, nil

}

func (c *Circuit) BtcUpperProve(data *btcupper.BatchedProofData, proofs []operations.Proof) (*operations.Proof, error) {
	logger.Debug("current zk circuit btcUpper prove....")
	if c.debug {
		logger.Warn("current zk circuit btcUpper prove is debug,skip prove")
		return debugProof()
	}
	proof, err := upperlevel.Prove(c.cfg.BtcSetupDir, data, proofs)
	if err != nil {
		logger.Error("btcUpper prove error: %v", err)
		return nil, err
	}
	return proof, nil

}

func (c *Circuit) BtcBulkProve(data *blockDepthUtil.BlockBulkProofData) (*operations.Proof, error) {
	logger.Debug("current zk circuit btcBlockBulk prove....")
	if c.debug {
		logger.Warn("current zk circuit btcBlockBulk prove is debug,skip prove")
		return debugProof()
	}
	proof, err := blockbulk.Prove(c.cfg.BtcSetupDir, data)
	if err != nil {
		logger.Error("btcBlockBulk prove error:%v", err)
		return nil, err
	}
	return proof, nil
}

func (c *Circuit) RedeemProve(tx, bh, bhf, duty *operations.Proof, genesisScRoot,
	currentSCSSZRoot []byte, txid, minerReward [32]byte, sigHashs [][32]byte, nbBeaconHeaders int, isFront bool) (*operations.Proof, error) {
	logger.Debug("current zk circuit redeemProve")
	if c.debug {
		logger.Warn("current zk circuit redeemProve prove is debug,skip prove")
		return debugProof()
	}
	//logger.Debug("redeem prove genesisScSszRoot: %x, currentScSszRoot: %x,txid: %x, minerReward: %x,sigHashs: %x,nbBeaconHeaders: %v,isFront: %v",
	//	genesisScRoot, currentSCSSZRoot, txid, minerReward, sigHashs, nbBeaconHeaders, isFront)

	proof, err := redeem.Prove(c.cfg.EthSetupDir, tx.Proof, tx.Witness, bh.Proof, bh.Witness, bhf.Proof, bhf.Witness, duty.Proof, duty.Witness,
		genesisScRoot, currentSCSSZRoot, txid, minerReward, sigHashs, nbBeaconHeaders, isFront)
	if err != nil {
		logger.Error("redeem prove error:%v", err)
		return nil, err
	}
	return proof, nil
}

func (c *Circuit) BeaconHeaderFinalityUpdateProve(finalityUpdate *proverType.FinalityUpdate, scUpdate *proverType.SyncCommittee) (*operations.Proof, error) {
	logger.Debug("current zk circuit beaconHeaderFinalityUpdateProve")
	if c.debug {
		logger.Warn("current zk circuit beaconHeaderFinalityUpdateProve prove is debug,skip prove")
		return debugProof()
	}
	proof, err := beaconheaderfinality.Prove(c.cfg.EthSetupDir, finalityUpdate, scUpdate)
	if err != nil {
		logger.Error("beacon header finality update prove error:%v", err)
		return nil, err
	}
	return proof, nil
}

func (c *Circuit) BeaconHeaderProve(header *proverType.BeaconHeaderChain) (*operations.Proof, error) {
	logger.Debug("current zk circuit BeaconHeaderProve")
	if c.debug {
		logger.Warn("current zk circuit BeaconHeaderProve prove is debug,skip prove")
		return debugProof()
	}
	proof, err := beaconheader.Prove(c.cfg.EthSetupDir, c.cfg.DataDir, header)
	if err != nil {
		logger.Error("beacon header prove error:%v %v %v", header.BeginSlot, header.EndSlot, err)
		return nil, err
	}
	return proof, nil
}

func (c *Circuit) TxInEth2Prove(param *txineth2Utils.TxInEth2ProofData) (*operations.Proof, error) {
	logger.Debug("current zk circuit TxInEth2Prove")
	if c.debug {
		logger.Warn("current zk circuit TxInEth2Prove prove is debug,skip prove")
		return debugProof()
	}
	proof, err := txineth2.Prove(c.cfg.EthSetupDir, param)
	if err != nil {
		logger.Error("txInEth2 prove error:%v", err)
		return nil, err
	}
	return proof, nil
}

func (c *Circuit) SyncInnerProve(index uint64, update *proverType.SyncCommittee) (*operations.Proof, error) {
	logger.Debug("current zk circuit syncUnitInner....")
	if c.debug {
		logger.Warn("current zk circuit syncUnitInner  prove is debug,skip prove")
		return debugProof()
	}
	instance := syncCommittee.NewSyncCommitteeInner(syncCommittee.NewSyncCommitteeInnerConfig(c.cfg.EthSetupDir, c.cfg.SrsDir, c.cfg.DataDir))
	err := instance.Load()
	if err != nil {
		logger.Error("load circuit setup file error: %v", err)
		return nil, err
	}
	proof, err := instance.ProveSingle(update, int(index))
	if err != nil {
		logger.Error("syncUnitInner prove error: %v", err)
		return nil, err
	}
	return proof, nil
}

func (c *Circuit) SyncOutProve(index uint64, update *proverType.SyncCommittee, innerProofs []operations.Proof) (*operations.Proof, error) {
	logger.Debug("current zk circuit syncOuter prove")
	if c.debug {
		logger.Warn("current zk circuit syncOuter prove is debug mode,skip prove")
		proof, err := debugProof()
		if err != nil {
			logger.Error("debug proof error:%v", err)
			return nil, err
		}
		return proof, nil
	}
	outerProof, err := outerProve(c.cfg.EthSetupDir, c.cfg.SrsDir, c.cfg.DataDir, update, innerProofs)
	if err != nil {
		logger.Error("outer prove error:%v", err)
		return nil, err
	}
	return outerProof, nil
}

func (c *Circuit) SyncCommitteeUnitProve(period uint64, outerProof operations.Proof, update *proverType.SyncCommitteeUpdate) (*operations.Proof, error) {
	logger.Debug("current zk circuit syncUnit prove")
	if c.debug {
		logger.Warn("current zk circuit syncUnit prove is debug mode,skip prove")
		proof, err := debugProof()
		if err != nil {
			logger.Error("debug proof error:%v", err)
			return nil, err
		}
		return proof, err
	}
	unitProof, err := unitProve(c.cfg.EthSetupDir, c.cfg.SrsDir, c.cfg.DataDir, update, &outerProof)
	if err != nil {
		logger.Error("syncUnit prove error:%v", err)
		return nil, err
	}
	return unitProof, nil
}

func (c *Circuit) SyncCommitteeDutyProve(choice string, first, unit, nextOuter *operations.Proof,
	beginId, relayId, endId []byte, scIndex int, update *proverType.SyncCommitteeUpdate) (*operations.Proof, *operations.Proof, error) {
	logger.Debug("current zk circuit syncCommitDuty prove,choice: %v", choice)
	if c.debug {
		logger.Warn("current zk circuit syncCommitDuty prove is debug mode,skip prove")
		proof, err := debugProof()
		if err != nil {
			logger.Error("debug syncCommitDuty error:%v", err)
			return nil, nil, err
		}
		return proof, proof, err
	}
	firstGenesis := choice == SyncCommitteeGenesis
	recursiveProof, err := syncCommittee.ProveRecursive(c.cfg.EthSetupDir, c.cfg.DataDir, first.Proof, unit.Proof, first.Witness, unit.Witness,
		beginId, relayId, endId, scIndex, firstGenesis, false)
	if err != nil {
		logger.Error("syncRecursive prove error:%v", err)
		return nil, nil, err
	}
	if recursiveProof == nil {
		return nil, nil, fmt.Errorf("syncRecursive proof is nil")
	}
	dutyProof, err := syncCommittee.ProveDuty(c.cfg.EthSetupDir, c.cfg.DataDir, [32]byte(beginId), recursiveProof.Proof, nextOuter.Proof, recursiveProof.Witness,
		nextOuter.Witness, update.CurrentSyncCommittee, false)
	if err != nil {
		logger.Error("syncCommitDuty prove error:%v", err)
		return nil, nil, err
	}
	if dutyProof == nil {
		return nil, nil, fmt.Errorf("syncDuty proof is nil")
	}
	return dutyProof, recursiveProof, err
}

func NewCircuit(cfg *CircuitConfig) (*Circuit, error) {
	if cfg.CacheCap > 0 {
		logger.Debug("current zk circuit init lru cap: %v", cfg.CacheCap)
		operations.InitLru(cfg.CacheCap)
	}
	return &Circuit{
		cfg:   cfg,
		debug: cfg.Debug,
	}, nil
}

func unitProve(setupDir, srsDir, datadir string, update *proverType.SyncCommitteeUpdate, outerProof *operations.Proof) (*operations.Proof, error) {
	config := syncCommittee.NewSyncCommitteeUnitConfig(setupDir, srsDir, datadir)
	instance := syncCommittee.NewSyncCommitteeUnit(config)
	err := instance.Load()
	if err != nil {
		return nil, err
	}
	proofs, err := instance.Prove(
		&outerProof.Proof,
		&outerProof.Witness,
		update)
	if err != nil {
		return nil, err
	}
	return proofs, nil
}

func outerProve(setupDir string, srsDir, datadir string, update *proverType.SyncCommittee, innerProofs []operations.Proof) (*operations.Proof, error) {
	config := syncCommittee.NewSyncCommitteeOuterConfig(setupDir, srsDir, datadir)
	instance := syncCommittee.NewSyncCommitteeOuter(config)
	err := instance.Load()
	if err != nil {
		return nil, err
	}
	if len(innerProofs) != proversCommon.NbBatches {
		return nil, fmt.Errorf("invalid number of inner proofs: %v", len(innerProofs))
	}
	tmpProofs := [proversCommon.NbBatches]*nativeplonk.Proof{}
	tmpWitnesses := [proversCommon.NbBatches]*witness.Witness{}
	for i := 0; i < proversCommon.NbBatches; i++ {
		tmpProofs[i] = &innerProofs[i].Proof
		tmpWitnesses[i] = &innerProofs[i].Witness
	}
	proofs, err := instance.Prove(
		tmpProofs,
		tmpWitnesses,
		update)
	if err != nil {
		return nil, err
	}
	return proofs, nil
}
