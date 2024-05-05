package circuits

import (
	btcproverUtils "github.com/lightec-xyz/btc_provers/utils"
	ethblock "github.com/lightec-xyz/provers/circuits/fabric/tx-in-eth2"
	proverType "github.com/lightec-xyz/provers/circuits/types"
	"github.com/lightec-xyz/reLight/circuits/common"
	"github.com/lightec-xyz/reLight/circuits/utils"
)

type ICircuit interface {
	// todo per 32 slot to generate proof
	BeaconHeaderFinalityUpdateProve(genesisSCSSZRoot string, recursiveProof, recursiveWitness, outerProof,
		outerWitness string, finalityUpdate *proverType.FinalityUpdate, scUpdate *proverType.SyncCommitteeUpdate) (*common.Proof, error)
	// todo find redeem tx,immediately to generate proof
	TxInEth2Prove(param *ethblock.TxInEth2ProofData) (*common.Proof, error)
	// todo only container rededm tx, we should need generate prof
	BeaconHeaderProve(header proverType.BeaconHeaderChain) (*common.Proof, error)
	// todo submit to eth contract proof
	RedeemProve(txProof, txWitness, bhProof, bhWitness, bhfProof, bhfWitness string, beginId, endId, genesisScRoot, currentSCSSZRoot string,
		txVar, receiptVar []string) (*common.Proof, error)
	DepositProve(data *btcproverUtils.GrandRollupProofData) (*common.Proof, error)
	GenesisProve(firstProof, secondProof, firstWitness, secondWitness string,
		genesisId, firstId, secondId []byte) (*common.Proof, error)
	UnitProve(period uint64, update *utils.SyncCommitteeUpdate) (*common.Proof, *common.Proof, error)
	RecursiveProve(choice string, firstProof, secondProof, firstWitness, secondWitness string,
		beginId, relayId, endId []byte) (*common.Proof, error)
	UpdateChangeProve(data *btcproverUtils.GrandRollupProofData) (*common.Proof, error)
}
