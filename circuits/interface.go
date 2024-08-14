package circuits

import (
	ethCommon "github.com/ethereum/go-ethereum/common"
	blockCu "github.com/lightec-xyz/btc_provers/utils/blockchain"
	blockDu "github.com/lightec-xyz/btc_provers/utils/blockdepth"
	txChainU "github.com/lightec-xyz/btc_provers/utils/txinchain"
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
	UnitProve(period uint64, update *utils.SyncCommitteeUpdate) (*common.Proof, *common.Proof, error)
	GenesisProve(first, second *common.Proof,
		genesisId, firstId, secondId []byte) (*common.Proof, error)
	RecursiveProve(choice string, firstProof, secondProof *common.Proof,
		beginId, relayId, endId []byte) (*common.Proof, error)
	TxInEth2Prove(param *ethblock.TxInEth2ProofData) (*common.Proof, error)
	BeaconHeaderProve(header proverType.BeaconHeaderChain) (*common.Proof, error)
	BeaconHeaderFinalityUpdateProve(genesisSCSSZRoot string, recursive, outer *common.Proof, finalityUpdate *proverType.FinalityUpdate,
		scUpdate *proverType.SyncCommitteeUpdate) (*common.Proof, error)
	RedeemProve(tx, bh, bhf *common.Proof, beginId, endId []byte, genesisScRoot, currentSCSSZRoot string,
		txVar, receiptVar []string) (*common.Proof, error)

	BtcBulkProve(data *blockDu.BlockBulkProofData) (*common.Proof, error)
	BtcPackProve(data *blockDu.BulksProofData, recursive, bulk *common.Proof) (*common.Proof, error)
	BtcDepthRecursiveProve(data *blockDu.BulksProofData, recursive, unit *common.Proof) (*common.Proof, error)
	BtcBaseProve(data *blockCu.BaseLevelProofData) (*common.Proof, error)
	BtcMiddleProve(data *blockCu.BatchedProofData, batch []common.Proof) (*common.Proof, error)
	BtcUpperProve(data *blockCu.BatchedProofData, super []common.Proof) (*common.Proof, error)
	BtcDuperRecursiveProve(data *blockCu.RecursiveProofData, first, second *common.Proof) (*common.Proof, error)
	BtcChainProve(data *blockCu.BlockChainProofData, recursive, base, middle, upper *common.Proof) (*common.Proof, error)
	BtcDepositProve(data *txChainU.TxInChainProofData, blockChain, txDepth, cpDepth *common.Proof,
		r, s ethCommon.Hash, addr ethCommon.Address) (*common.Proof, error)
	BtcChangeProve(data *txChainU.TxInChainProofData, blockChain, txDepth, cpDepth, redeem *common.Proof,
		r, s ethCommon.Hash, addr ethCommon.Address) (*common.Proof, error)
}
