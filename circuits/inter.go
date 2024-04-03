package circuits

import (
	"github.com/ethereum/go-ethereum/ethclient"
	apiclient "github.com/lightec-xyz/provers/utils/api-client"
	"github.com/lightec-xyz/reLight/circuits/common"
	"github.com/lightec-xyz/reLight/circuits/utils"
)

type ICircuit interface {
	CheckPointFinalityProve() (*common.Proof, error)
	TxInEth2Prove() (*common.Proof, error)
	
	TxBlockIsParentOfCpsProve() (*common.Proof, error)
	RedeemProve() (*common.Proof, error)

	DepositProve(ec *ethclient.Client, cl *apiclient.Client, txHash string) (*common.Proof, error)

	GenesisProve(firstProof, secondProof, firstWitness, secondWitness []byte,
		genesisId, firstId, secondId []byte) (*common.Proof, error)
	UnitProve(period uint64, update *utils.LightClientUpdateInfo) (*common.Proof, error)

	RecursiveProve(choice string, firstProof, secondProof, firstWitness, secondWitness []byte,
		beginId, relayId, endId []byte) (*common.Proof, error)
}
