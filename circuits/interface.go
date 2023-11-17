package circuits

import (
	ethCommon "github.com/ethereum/go-ethereum/common"

	blockChainUtil "github.com/lightec-xyz/btc_provers/utils/blockchain"
	blockDepthUtil "github.com/lightec-xyz/btc_provers/utils/blockdepth"
	txChainUtil "github.com/lightec-xyz/btc_provers/utils/txinchain"
	"github.com/lightec-xyz/common/operations"
	proverType "github.com/lightec-xyz/provers/circuits/types"
	txineth2Utils "github.com/lightec-xyz/provers/utils/tx-in-eth2"
)

const (
	SyncCommitteeGenesis   = "genesis"
	SyncCommitteeRecursive = "recursive"
	BtcChainUpper          = "upper"
	BtcBlockChain          = "blockchain"
	BtcChainHybrid         = "hybrid"
	BtcChainBase           = "base"
)

type ICircuit interface {
	SyncInnerProve(index uint64, update *proverType.SyncCommittee) (*operations.Proof, error)
	SyncOutProve(index uint64, update *proverType.SyncCommittee, innersProof []operations.Proof) (*operations.Proof, error)
	SyncCommitteeUnitProve(period uint64, outerProof operations.Proof, update *proverType.SyncCommitteeUpdate) (*operations.Proof, error)
	SyncCommitteeDutyProve(choice string, first, unit, nextOuter *operations.Proof,
		beginId, relayId, endId []byte, scIndex int, update *proverType.SyncCommitteeUpdate) (*operations.Proof, *operations.Proof, error)
	TxInEth2Prove(param *txineth2Utils.TxInEth2ProofData) (*operations.Proof, error)
	BeaconHeaderProve(header *proverType.BeaconHeaderChain) (*operations.Proof, error)
	BeaconHeaderFinalityUpdateProve(finalityUpdate *proverType.FinalityUpdate,
		scUpdate *proverType.SyncCommittee) (*operations.Proof, error)
	RedeemProve(tx, bh, bhf, duty *operations.Proof, genesisScRoot, currentSCSSZRoot []byte,
		btcTxId, minerReward [32]byte, sigHashs [][32]byte, nBeaconHeaders int, isFront bool) (*operations.Proof, error)

	BtcBulkProve(data *blockDepthUtil.BlockBulkProofData) (*operations.Proof, error)
	BtcDepthRecursiveProve(recursive bool, step uint64, data *blockDepthUtil.RecursiveBulksProofData, first *operations.Proof) (*operations.Proof, error)
	BtcTimestamp(cpTime *blockDepthUtil.CptimestampProofData, smoothData *blockDepthUtil.SmoothedTimestampProofData) (*operations.Proof, error)
	BtcBaseProve(data *blockChainUtil.BaseLevelProofData) (*operations.Proof, error)
	BtcMiddleProve(data *blockChainUtil.BatchedProofData, batch []operations.Proof) (*operations.Proof, error)
	BtcUpperProve(data *blockChainUtil.BatchedProofData, super []operations.Proof) (*operations.Proof, error)
	BtcChainRecursiveProve(firstType, secondType string, firstStep, secondStep uint64, data *blockChainUtil.BlockChainProofData, first, second *operations.Proof) (*operations.Proof, error)
	BtcChainHybridProve(firstType string, firstStep uint64, data *blockChainUtil.HybridProofData, first *operations.Proof) (*operations.Proof, error)
	BtcDepositProve(chainType string, chainStep, txStep, cpStep uint64, txRecursive, cpRecursive bool, data *txChainUtil.TxInChainProofData, blockChain, txDepth, cpDepth, sigVerify *operations.Proof,
		addr ethCommon.Address, smoothedTimestamp uint32, cpFlag uint8, sigVerifyData *blockDepthUtil.SigVerifProofData) (*operations.Proof, error)
	BtcRedeemProve(chainType string, chainStep, txStep, cpStep uint64, txRecursive, cpRecursive bool, data *txChainUtil.TxInChainProofData, blockChain, txDepth, cpDepth, redeem, sigVerify *operations.Proof,
		minerReward [32]byte, addr ethCommon.Address, smoothedTimestamp uint32, cpFlag uint8, sigVerifyData *blockDepthUtil.SigVerifProofData) (*operations.Proof, error)
}
