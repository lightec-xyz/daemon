package circuits

import (
	ethblock "github.com/lightec-xyz/provers/circuits/fabric/tx-in-eth2"
	"github.com/lightec-xyz/reLight/circuits/common"
	"github.com/lightec-xyz/reLight/circuits/utils"
)

type ICircuit interface {
	// todo per 32 slot to generate proof
	BeaconHeaderFinalityUpdateProve() (*common.Proof, error)
	// todo find redeem tx,immediately to generate proof
	TxInEth2Prove(param *ethblock.TxInEth2ProofData) (*common.Proof, error)
	// todo only container rededm tx, we should need generate prof
	BeaconHeaderProve() (*common.Proof, error)
	// todo submit to eth contract proof
	RedeemProve() (*common.Proof, error)
	DepositProve(txId, blockHash string) (*common.Proof, error)
	GenesisProve(firstProof, secondProof, firstWitness, secondWitness []byte,
		genesisId, firstId, secondId []byte) (*common.Proof, error)
	UnitProve(period uint64, update *utils.LightClientUpdateInfo) (*common.Proof, error)
	RecursiveProve(choice string, firstProof, secondProof, firstWitness, secondWitness []byte,
		beginId, relayId, endId []byte) (*common.Proof, error)
}
