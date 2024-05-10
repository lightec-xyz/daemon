package rpc

import (
	"github.com/lightec-xyz/daemon/common"
)

type INode interface {
	Version() (NodeInfo, error)
	AddWorker(endpoint string, max int) (string, error)
	ProofInfo(txIds []string) ([]ProofInfo, error)
	Transaction(txHash string) (Transaction, error)
	TransactionsByHeight(height uint64, network string) ([]string, error)
	Transactions(txId []string) ([]Transaction, error)
	GetZkProofTask(request common.TaskRequest) (*common.TaskResponse, error)
	SubmitProof(req *common.SubmitProof) (string, error)
	TxesByAddr(addr, txType string) ([]Transaction, error)
	Stop() error
}

// IProof api between node and proof
type IProof interface {
	GenDepositProof(req DepositRequest) (DepositResponse, error)
	GenVerifyProof(req VerifyRequest) (VerifyResponse, error)
	GenSyncCommGenesisProof(req SyncCommGenesisRequest) (SyncCommGenesisResponse, error)
	GenSyncCommitUnitProof(req SyncCommUnitsRequest) (SyncCommUnitsResponse, error)
	GenSyncCommRecursiveProof(req SyncCommRecursiveRequest) (SyncCommRecursiveResponse, error)

	GenRedeemProof(req *RedeemRequest) (*RedeemResponse, error)
	TxInEth2Prove(req *TxInEth2ProveRequest) (*TxInEth2ProveResponse, error)
	BlockHeaderProve(req *BlockHeaderRequest) (*BlockHeaderResponse, error)
	BlockHeaderFinalityProve(req *BlockHeaderFinalityRequest) (*BlockHeaderFinalityResponse, error)

	ProofInfo(proofId string) (ProofInfo, error)
}

type IWorker interface {
	IProof
	AddReqNum()
	DelReqNum()
	MaxNums() int
	CurrentNums() int
	Id() string
}

// ICRequest Todo
type ICRequest interface {
	DepositProofReq() (IDepositRequest, error)
	DepositResponse(data DepositResponse) error
}
