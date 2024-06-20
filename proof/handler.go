package proof

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/store"
)

var _ rpc.IProof = (*Handler)(nil)

type Handler struct {
	memoryStore store.IStore
	store       store.IStore
	maxNums     int // The maximum number of proofs that can be generated at the same time
	proofs      *sync.Map
	currentNums atomic.Int64
	lock        sync.Mutex
}

func (h *Handler) BtcBulkProve(data *rpc.BtcBulkRequest) (*rpc.BtcBulkResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) BtcPackedRequest(data *rpc.BtcPackedRequest) (*rpc.BtcPackResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) BtcWrapProve(data *rpc.BtcWrapRequest) (*rpc.BtcWrapResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) TxInEth2Prove(req *rpc.TxInEth2ProveRequest) (*rpc.TxInEth2ProveResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) BlockHeaderProve(req *rpc.BlockHeaderRequest) (*rpc.BlockHeaderResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) BlockHeaderFinalityProve(req *rpc.BlockHeaderFinalityRequest) (*rpc.BlockHeaderFinalityResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) ProofInfo(proofId string) (rpc.ProofInfo, error) {
	logger.Debug("proof info: %v", proofId)
	time.Sleep(10 * time.Second)
	return rpc.ProofInfo{}, nil
}

func (h *Handler) GenDepositProof(req rpc.DepositRequest) (rpc.DepositResponse, error) {
	logger.Debug("gen deposit proof")
	time.Sleep(10 * time.Second)
	return rpc.DepositResponse{
		Proof: common.ZkProof([]byte("deposit proof")),
	}, nil
}

func (h *Handler) GenRedeemProof(req *rpc.RedeemRequest) (*rpc.RedeemResponse, error) {
	logger.Debug("gen redeem proof")
	time.Sleep(10 * time.Second)
	return &rpc.RedeemResponse{
		Proof: common.ZkProof([]byte("redeem proof")),
	}, nil
}

func (h *Handler) GenVerifyProof(req rpc.VerifyRequest) (rpc.VerifyResponse, error) {
	logger.Debug("gen verify proof")
	time.Sleep(10 * time.Second)
	return rpc.VerifyResponse{
		Proof: []byte("verify proof"),
	}, nil
}

func (h *Handler) GenSyncCommGenesisProof(req rpc.SyncCommGenesisRequest) (rpc.SyncCommGenesisResponse, error) {
	logger.Debug("gen sync comm genesis proof")
	time.Sleep(10 * time.Second)
	return rpc.SyncCommGenesisResponse{
		Proof: common.ZkProof([]byte("genesis proof")),
	}, nil
}

func (h *Handler) GenSyncCommitUnitProof(req rpc.SyncCommUnitsRequest) (rpc.SyncCommUnitsResponse, error) {
	logger.Debug("gen sync comm units proof")
	time.Sleep(10 * time.Second)
	return rpc.SyncCommUnitsResponse{
		Proof: common.ZkProof([]byte("units proof")),
	}, nil
}

func (h *Handler) GenSyncCommRecursiveProof(req rpc.SyncCommRecursiveRequest) (rpc.SyncCommRecursiveResponse, error) {
	logger.Debug("gen sync comm recursive proof")
	time.Sleep(10 * time.Second)
	return rpc.SyncCommRecursiveResponse{
		Proof: common.ZkProof([]byte("recursive proof")),
	}, nil
}

func (h *Handler) MaxNums() (int, error) {
	logger.Debug("max nums: %v", h.maxNums)
	time.Sleep(2 * time.Second)
	return h.maxNums, nil
}

func (h *Handler) CurrentNums() (int, error) {
	logger.Debug("current nums: %v", h.currentNums.Load())
	time.Sleep(2 * time.Second)
	return int(h.currentNums.Load()), nil
}

func (h *Handler) AddReqNum() {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) DelReqNum() {
	//TODO implement me
	panic("implement me")
}

func NewHandler(store, memoryStore store.IStore, max int) *Handler {
	return &Handler{
		memoryStore: memoryStore,
		store:       store,
		maxNums:     max,
	}
}
