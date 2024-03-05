package proof

import (
	"encoding/json"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/store"
	"sync"
)

var _ rpc.IProof = (*Handler)(nil)

type Handler struct {
	memoryStore store.IStore
	maxNums     int // The maximum number of proofs that can be generated at the same time
	proofs      *sync.Map
	currentNums int
	lock        sync.Mutex
}

func (h *Handler) ProofInfo(proofId string) (rpc.ProofInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) GenDepositProof(req rpc.DepositRequest) (rpc.DepositResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) GenRedeemProof(req rpc.RedeemRequest) (rpc.RedeemResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) GenVerifyProof(req rpc.VerifyRequest) (rpc.VerifyResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) GenSyncCommGenesisProof(req rpc.SyncCommGenesisRequest) (rpc.SyncCommGenesisResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) GenSyncCommitUnitProof(req rpc.SyncCommUnitsRequest) (rpc.SyncCommUnitsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) GenSyncCommRecursiveProof(req rpc.SyncCommRecursiveRequest) (rpc.SyncCommRecursiveResponse, error) {
	//TODO implement me
	panic("implement me")
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
	//TODO implement me
	panic("implement me")
}

func (h *Handler) CurrentNums() int {
	//TODO implement me
	panic("implement me")
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
