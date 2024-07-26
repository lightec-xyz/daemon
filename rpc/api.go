package rpc

import (
	"github.com/lightec-xyz/daemon/common"
)

type INode interface {
	Version() (NodeInfo, error)
	AddWorker(endpoint string, max int) (string, error)
	ProofInfo(txIds []string) ([]ProofInfo, error)
	Transaction(txHash string) (*Transaction, error)
	TransactionsByHeight(height uint64, network string) ([]string, error)
	Transactions(txId []string) ([]*Transaction, error)
	GetZkProofTask(request common.TaskRequest) (*common.TaskResponse, error)
	SubmitProof(req *common.SubmitProof) (string, error)
	TxesByAddr(addr, txType string) ([]*Transaction, error)
	ProofTask(id string) (*ProofTaskInfo, error)
	PendingTask() ([]*ProofTaskInfo, error)
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

	BtcBulkProve(req *BtcBulkRequest) (*BtcBulkResponse, error)
	BtcPackedRequest(req *BtcPackedRequest) (*BtcPackResponse, error)
	BtcWrapProve(req *BtcWrapRequest) (*BtcWrapResponse, error)
	BtcBaseProve(req *BtcBaseRequest) (*ProofResponse, error)
	BtcMiddleProve(req *BtcMiddleRequest) (*ProofResponse, error)
	BtcUpProve(req *BtcUpperRequest) (*ProofResponse, error)
	ProofInfo(proofId string) (ProofInfo, error)
	SupportProofType() []common.ZkProofType
	Close() error
}

type IWorker interface {
	IProof
	AddReqNum()
	DelReqNum()
	MaxNums() int
	CurrentNums() int
	Id() string
}
