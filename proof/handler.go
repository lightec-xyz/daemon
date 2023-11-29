package proof

import (
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/store"
)

var _ rpc.ProofAPI = (*Handler)(nil)

type Handler struct {
	memoryStore     store.IStore
	store           store.IStore
	maxParallelNums int
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Info() (rpc.ProofInfo, error) {
	info := rpc.ProofInfo{
		Version: "1.2.3",
	}
	return info, nil
}

func (h *Handler) GenZkProof(request rpc.ProofRequest) (rpc.ProofResponse, error) {
	//todo
	response := rpc.ProofResponse{
		TxId:   request.TxId,
		Status: 0,
		Msg:    "ok",
		Proof:  "test proof",
	}
	return response, nil
}

func (h *Handler) ProofStatus(proofId string) (rpc.ProofStatus, error) {
	//todo
	status := rpc.ProofStatus{
		State: 0,
		Msg:   "ok",
	}
	return status, nil
}
