package circuits

import (
	ethblock "github.com/lightec-xyz/provers/circuits/fabric/tx-in-eth2"
	proverType "github.com/lightec-xyz/provers/circuits/types"
	"github.com/lightec-xyz/reLight/circuits/common"
	"github.com/lightec-xyz/reLight/circuits/utils"
)

type ICircuit interface {
	// todo per 32 slot to generate proof
	BeaconHeaderFinalityUpdateProve(genesisSCSSZRoot string, recursiveProof, recursiveWitness, outerProof,
		outerWitness []byte, finalityUpdate *proverType.FinalityUpdate, scUpdate *proverType.SyncCommitteeUpdate) (*common.Proof, error)
	// todo find redeem tx,immediately to generate proof
	TxInEth2Prove(param *ethblock.TxInEth2ProofData) (*common.Proof, error)
	// todo only container rededm tx, we should need generate prof
	BeaconHeaderProve(header proverType.BeaconHeaderChain) (*common.Proof, error)
	// todo submit to eth contract proof
	RedeemProve(txProof, txWitness, bhProof, bhWitness, bhfProof, bhfWitness, beginId, endId, genesisScRoot,
		currentSCSSZRoot, txVarBytes, receiptVarBytes []byte) (*common.Proof, error)
	DepositProve(txId, blockHash string) (*common.Proof, error)
	GenesisProve(firstProof, secondProof, firstWitness, secondWitness []byte,
		genesisId, firstId, secondId []byte) (*common.Proof, error)
	UnitProve(period uint64, update *utils.LightClientUpdateInfo) (*common.Proof, *common.Proof, error)
	RecursiveProve(choice string, firstProof, secondProof, firstWitness, secondWitness []byte,
		beginId, relayId, endId []byte) (*common.Proof, error)
	UpdateChangeProve(txId, blockHash string) (*common.Proof, error)
}
