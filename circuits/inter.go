package circuits

import (
	ethblock "github.com/lightec-xyz/provers/circuits/fabric/tx-in-eth2"
	"github.com/lightec-xyz/reLight/circuits/common"
	"github.com/lightec-xyz/reLight/circuits/utils"
)

type ICircuit interface {
	CheckPointFinalityProve() (*common.Proof, error)
	TxInEth2Prove(param *ethblock.TxInEth2ProofData) (*common.Proof, error)
	TxBlockIsParentOfCheckPointProve() (*common.Proof, error)
	RedeemProve() (*common.Proof, error)
	DepositProve() (*common.Proof, error)
	GenesisProve(firstProof, secondProof, firstWitness, secondWitness []byte,
		genesisId, firstId, secondId []byte) (*common.Proof, error)
	UnitProve(period uint64, update *utils.LightClientUpdateInfo) (*common.Proof, error)
	RecursiveProve(choice string, firstProof, secondProof, firstWitness, secondWitness []byte,
		beginId, relayId, endId []byte) (*common.Proof, error)
}
