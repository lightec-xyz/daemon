package circuits

import (
	btcprovertypes "github.com/lightec-xyz/btc_provers/circuits/types"
	btcbase "github.com/lightec-xyz/btc_provers/utils/baselevel"
	grUtil "github.com/lightec-xyz/btc_provers/utils/grandrollup"
	btcmiddle "github.com/lightec-xyz/btc_provers/utils/midlevel"
	btcupper "github.com/lightec-xyz/btc_provers/utils/upperlevel"
	ethblock "github.com/lightec-xyz/provers/circuits/fabric/tx-in-eth2"
	proverType "github.com/lightec-xyz/provers/circuits/types"
	"github.com/lightec-xyz/reLight/circuits/common"
	"github.com/lightec-xyz/reLight/circuits/utils"
)

const (
	BtcBulk   = "bulk"
	BtcPacked = "packed"
)

type ICircuit interface {
	BeaconHeaderFinalityUpdateProve(genesisSCSSZRoot string, recursiveProof, recursiveWitness, outerProof,
		outerWitness string, finalityUpdate *proverType.FinalityUpdate, scUpdate *proverType.SyncCommitteeUpdate) (*common.Proof, error)
	TxInEth2Prove(param *ethblock.TxInEth2ProofData) (*common.Proof, error)
	BeaconHeaderProve(header proverType.BeaconHeaderChain) (*common.Proof, error)
	RedeemProve(txProof, txWitness, bhProof, bhWitness, bhfProof, bhfWitness string, beginId, endId, genesisScRoot, currentSCSSZRoot string,
		txVar, receiptVar []string) (*common.Proof, error)
	DepositProve(data *grUtil.GrandRollupProofData) (*common.Proof, error)
	GenesisProve(firstProof, secondProof, firstWitness, secondWitness string,
		genesisId, firstId, secondId []byte) (*common.Proof, error)
	UnitProve(period uint64, update *utils.SyncCommitteeUpdate) (*common.Proof, *common.Proof, error)
	RecursiveProve(choice string, firstProof, secondProof, firstWitness, secondWitness string,
		beginId, relayId, endId []byte) (*common.Proof, error)
	UpdateChangeProve(data *grUtil.GrandRollupProofData) (*common.Proof, error)
	BtcBulkProve(data *btcprovertypes.BlockHeaderChain) (*common.Proof, error)
	BtcPackProve(data *btcprovertypes.BlockHeaderChain) (*common.Proof, error)
	BtcWrapProve(flag, proof, witness, beginHash, endHash string, nbBlocks uint64) (*common.Proof, error)
	BtcBaseProve(req *btcbase.BaseLevelProofData) (*common.Proof, error)
	BtcMiddleProve(req *btcmiddle.MidLevelProofData, proofList *btcprovertypes.WitnessFile) (*common.Proof, error)
	BtcUpperProve(req *btcupper.UpperLevelProofData, proofList *btcprovertypes.WitnessFile) (*common.Proof, error)
}
