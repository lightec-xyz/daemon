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

func (h *Handler) ProofInfo(proofId string) (rpc.ProofResponse, error) {
	//todo
	status := rpc.ProofResponse{
		TxId:   proofId,
		Status: 0,
		Msg:    "ok",
		Proof:  "test proof",
	}
	return status, nil
}
