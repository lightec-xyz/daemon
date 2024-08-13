package circuits

import (
	ethCommon "github.com/ethereum/go-ethereum/common"
	blockchainUtil "github.com/lightec-xyz/btc_provers/utils/blockchain"
	blockdepthUtil "github.com/lightec-xyz/btc_provers/utils/blockdepth"
	txinchainUtil "github.com/lightec-xyz/btc_provers/utils/txinchain"
	ethblock "github.com/lightec-xyz/provers/circuits/fabric/tx-in-eth2"
	proverType "github.com/lightec-xyz/provers/circuits/types"
	"github.com/lightec-xyz/reLight/circuits/common"
	"github.com/lightec-xyz/reLight/circuits/utils"
)

const (
	BtcBulk       = "bulk"
	BtcPacked     = "packed"
	SyncGenesis   = "genesis"
	SyncRecursive = "recursive"
)

type ICircuit interface {
	BeaconHeaderFinalityUpdateProve(genesisSCSSZRoot string, recursiveProof, recursiveWitness, outerProof,
		outerWitness string, finalityUpdate *proverType.FinalityUpdate, scUpdate *proverType.SyncCommitteeUpdate) (*common.Proof, error)
	TxInEth2Prove(param *ethblock.TxInEth2ProofData) (*common.Proof, error)
	BeaconHeaderProve(header proverType.BeaconHeaderChain) (*common.Proof, error)
	RedeemProve(txProof, txWitness, bhProof, bhWitness, bhfProof, bhfWitness string, beginId, endId, genesisScRoot, currentSCSSZRoot string,
		txVar, receiptVar []string) (*common.Proof, error)
	GenesisProve(firstProof, secondProof, firstWitness, secondWitness string,
		genesisId, firstId, secondId []byte) (*common.Proof, error)
	UnitProve(period uint64, update *utils.SyncCommitteeUpdate) (*common.Proof, *common.Proof, error)
	RecursiveProve(choice string, firstProof, secondProof, firstWitness, secondWitness string,
		beginId, relayId, endId []byte) (*common.Proof, error)

	// for btc
	BtcBulkProve(proofData *blockdepthUtil.BlockBulkProofData) (*common.Proof, error)
	BtcDepthRecursiveProve(proofData *blockdepthUtil.BulksProofData, recursiveProofFile, unitProofFile *common.Proof) (*common.Proof, error)
	BtcPackProve(proofData *blockdepthUtil.BulksProofData, depthRecursiveProofFile, bulkProofFile *common.Proof) (*common.Proof, error)
	BtcBaseProve(proofData *blockchainUtil.BaseLevelProofData) (*common.Proof, error)
	BtcMiddleProve(proofData *blockchainUtil.BatchedProofData, batchProofList []common.Proof) (*common.Proof, error)
	BtcUpperProve(proofData *blockchainUtil.BatchedProofData, superProofList []common.Proof) (*common.Proof, error)
	BtcDuperRecursiveProve(proofData *blockchainUtil.RecursiveProofData, firstProofFile, duperProofFile *common.Proof) (*common.Proof, error)
	BtcChainProve(
		proofData *blockchainUtil.BlockChainProofData,
		duperRecursiveProofFile, baseLevelProofFile, midLevelProofFile, upperLevelProofFile *common.Proof,
	) (*common.Proof, error)
	BtcDepositProve(
		proofData *txinchainUtil.TxInChainProofData,
		blockChainProofFile, txDepthProofFile, cpDepthProofFile *common.Proof,
		r, s ethCommon.Hash,
		proverAddr ethCommon.Address,
	) (*common.Proof, error)
	BtcChangeProve(
		proofData *txinchainUtil.TxInChainProofData,
		blockChainProofFile, txDepthProofFile, cpDepthProofFile, redeemInEthProofFile *common.Proof,
		r, s ethCommon.Hash,
		proverAddr ethCommon.Address,
	) (*common.Proof, error)
}
