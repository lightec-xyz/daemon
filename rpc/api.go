package rpc

import (
	"github.com/lightec-xyz/daemon/common"
)

type INode interface {
	IAdmin
	Version() (NodeInfo, error)
	AddWorker(endpoint string, max int) (string, error)
	ProofInfo(txIds []string) ([]ProofInfo, error)
	Transaction(txHash string) ([]*Transaction, error)
	TransactionsByHeight(height uint64, network string) ([]string, error)
	Transactions(txId []string) ([]*Transaction, error)
	GetZkProofTask(request common.TaskRequest) (*common.TaskResponse, error)
	SubmitProof(req *common.SubmitProof) (string, error)
	ProofTask(id string) (*ProofTaskInfo, error)
	PendingTask() ([]*ProofTaskInfo, error)
	Eth2Slot(height uint64) (uint64, error)
	Eth1Height(slot uint64) (uint64, error)
	AddP2pPeer(endpoint string) (string, error)
	ReScan(height uint64, chain string) error
	MinerInfo() ([]*MinerInfo, error)
	Stop() error
}

// IProof api between node and proof
type IProof interface {
	SyncCommInner(req *SyncCommInnerRequest) (*ProofResponse, error)
	SyncCommOuter(req *SyncCommOuterRequest) (*ProofResponse, error)
	SyncCommitUnitProve(req SyncCommUnitsRequest) (*SyncCommUnitsResponse, error)
	SyncCommDutyProve(req SyncCommDutyRequest) (*SyncCommDutyResponse, error)
	TxInEth2Prove(req *TxInEth2ProveRequest) (*TxInEth2ProveResponse, error)
	BlockHeaderProve(req *BlockHeaderRequest) (*BlockHeaderResponse, error)
	BlockHeaderFinalityProve(req *BlockHeaderFinalityRequest) (*BlockHeaderFinalityResponse, error)
	RedeemProof(req *RedeemRequest) (*RedeemResponse, error)
	BackendRedeemProof(req *RedeemRequest) (*RedeemResponse, error)

	BtcBulkProve(req *BtcBulkRequest) (*BtcBulkResponse, error)
	BtcBaseProve(req *BtcBaseRequest) (*ProofResponse, error)
	BtcMiddleProve(req *BtcMiddleRequest) (*ProofResponse, error)
	BtcUpperProve(req *BtcUpperRequest) (*ProofResponse, error)
	BtcDepthRecursiveProve(req *BtcDepthRecursiveRequest) (*ProofResponse, error)
	BtcDuperRecursiveProve(req *BtcDuperRecursiveRequest) (*ProofResponse, error)
	BtcDepositProve(req *BtcDepositRequest) (*ProofResponse, error)
	BtcChangeProve(req *BtcChangeRequest) (*ProofResponse, error)
	BtcTimestamp(req *BtcTimestampRequest) (*ProofResponse, error)
	ProofInfo(proofId string) (ProofInfo, error)
	SupportProofType() []common.ProofType
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
	RemoveUnSubmitTx(hash string) (string, error)
	RemoveUnGenProof(hash string) (string, error)
	SetGasPrice(gasPrice uint64) (string, error)
}

type IVerify interface {
	VerifyJwt(token string) (*CustomClaims, error)
	CheckPermission(method string) (Permission, error)
}

// ICheck when encode or deconde obj, need to check data
type ICheck interface {
	Check() error
}
