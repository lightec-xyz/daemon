package proof

import (
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/store"
)

var _ rpc.ProofAPI = (*Handler)(nil)

type Handler struct {
	store       store.IStore
	memoryStore store.IStore
}

func NewHandler(store store.IStore, memoryStore store.IStore) *Handler {
	return &Handler{
		store:       store,
		memoryStore: memoryStore,
	}
}

func (h *Handler) Info() (rpc.ProofInfo, error) {
	info := rpc.ProofInfo{
		Version: "1.2.3",
	}
	return info, nil
}

func (h *Handler) GenBtcProof(request rpc.ProofRequest) (rpc.BtcProofResponse, error) {
	//todo
	response := rpc.BtcProofResponse{
		TxId:   request.TxId,
		Status: 0,
		Msg:    "ok",
		Proof:  "test proof",
	}
	return response, nil
}

func (h *Handler) GenEthProof(request rpc.EthProofRequest) (rpc.EthProofResponse, error) {
	//todo
	response := rpc.EthProofResponse{
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
