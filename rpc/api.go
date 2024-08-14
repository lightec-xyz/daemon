package rpc

import (
	"github.com/lightec-xyz/daemon/common"
)

type INode interface {
	IAdmin
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
	GenSyncCommitUnitProof(req SyncCommUnitsRequest) (*SyncCommUnitsResponse, error)
	GenSyncCommGenesisProof(req SyncCommGenesisRequest) (*SyncCommGenesisResponse, error)
	GenSyncCommRecursiveProof(req SyncCommRecursiveRequest) (*SyncCommRecursiveResponse, error)
	TxInEth2Prove(req *TxInEth2ProveRequest) (*TxInEth2ProveResponse, error)
	BlockHeaderProve(req *BlockHeaderRequest) (*BlockHeaderResponse, error)
	BlockHeaderFinalityProve(req *BlockHeaderFinalityRequest) (*BlockHeaderFinalityResponse, error)
	GenRedeemProof(req *RedeemRequest) (*RedeemResponse, error)

	BtcBulkProve(req *BtcBulkRequest) (*BtcBulkResponse, error)
	BtcPackedRequest(req *BtcPackedRequest) (*BtcPackResponse, error)
	BtcBaseProve(req *BtcBaseRequest) (*ProofResponse, error)
	BtcMiddleProve(req *BtcMiddleRequest) (*ProofResponse, error)
	BtcUpperProve(req *BtcUpperRequest) (*ProofResponse, error)
	BtcDuperRecursiveProve(req *BtcDuperRecursiveRequest) (*ProofResponse, error)
	BtcDepthRecursiveProve(req *BtcDepthRequest) (*ProofResponse, error)
	BtcChainProve(req *BtcChainRequest) (*ProofResponse, error)
	BtcDepositProve(req *BtcDepositRequest) (*ProofResponse, error)
	BtcChangeProve(req *BtcChangeRequest) (*ProofResponse, error)

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

type IAdmin interface {
	RemoveRequest(id string) error
}

type IVerify interface {
	VerifyJwt(token string) (*CustomClaims, error)
	CheckPermission(method string) (Permission, error)
}
