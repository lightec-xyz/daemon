package proof

import (
	"encoding/json"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/store"
	"sync"
	"sync/atomic"
	"time"
)

var _ rpc.IProof = (*Handler)(nil)

type Handler struct {
	memoryStore store.IStore
	maxNums     int // The maximum number of proofs that can be generated at the same time
	proofs      *sync.Map
	currentNums atomic.Int64
	lock        sync.Mutex
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
		Body: []byte("deposit proof"),
	}, nil
}

func (h *Handler) GenRedeemProof(req rpc.RedeemRequest) (rpc.RedeemResponse, error) {
	logger.Debug("gen redeem proof")
	time.Sleep(10 * time.Second)
	return rpc.RedeemResponse{
		Body: []byte("redeem proof"),
	}, nil
}

func (h *Handler) GenVerifyProof(req rpc.VerifyRequest) (rpc.VerifyResponse, error) {
	logger.Debug("gen verify proof")
	time.Sleep(10 * time.Second)
	return rpc.VerifyResponse{
		Body: []byte("verify proof"),
	}, nil
}

func (h *Handler) GenSyncCommGenesisProof(req rpc.SyncCommGenesisRequest) (rpc.SyncCommGenesisResponse, error) {
	logger.Debug("gen sync comm genesis proof")
	time.Sleep(10 * time.Second)
	return rpc.SyncCommGenesisResponse{
		Body: []byte("genesis proof"),
	}, nil
}

func (h *Handler) GenSyncCommitUnitProof(req rpc.SyncCommUnitsRequest) (rpc.SyncCommUnitsResponse, error) {
	logger.Debug("gen sync comm units proof")
	time.Sleep(10 * time.Second)
	return rpc.SyncCommUnitsResponse{
		Body: []byte("units proof"),
	}, nil
}

func (h *Handler) GenSyncCommRecursiveProof(req rpc.SyncCommRecursiveRequest) (rpc.SyncCommRecursiveResponse, error) {
	logger.Debug("gen sync comm recursive proof")
	time.Sleep(10 * time.Second)
	return rpc.SyncCommRecursiveResponse{
		Body: []byte("recursive proof"),
	}, nil
}

func (h *Handler) AddReqNum() {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) DelReqNum() {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) MaxNums() int {
	return h.maxNums
}

func (h *Handler) CurrentNums() int {
	return int(h.currentNums.Load())
}

func NewHandler(memoryStore store.IStore, max int) *Handler {
	return &Handler{
		memoryStore: memoryStore,
		maxNums:     max,
	}
}

func objParse(src, dest interface{}) error {
	marshal, err := json.Marshal(src)
	if err != nil {
		return err
	}
	err = json.Unmarshal(marshal, dest)
	if err != nil {
		return err
	}
	return nil
}
