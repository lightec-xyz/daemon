package proof

import "github.com/lightec-xyz/daemon/store"

var _ API = (*Handler)(nil)

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

func (h *Handler) Info() (ProofInfo, error) {
	panic("implement me")
}

func (h *Handler) GenBtcProof(request BtcProofRequest) (BtcProofResponse, error) {
	panic("implement me")
}

func (h *Handler) GenEthProof(request EthProofRequest) (EthProofResponse, error) {
	panic("implement me")
}

func (h *Handler) ProofStatus(proofId string) (ProofStatus, error) {
	panic("implement me")
}
