package proof

import "github.com/lightec-xyz/daemon/store"

type Handler struct {
	store       *store.Store
	memoryStore *store.MemoryStore
}

func NewHandler(store *store.Store, memoryStore *store.MemoryStore) *Handler {
	return &Handler{
		store:       store,
		memoryStore: memoryStore,
	}
}

func (h *Handler) Info() (string, error) {
	panic("implement me")
}

func (h *Handler) GenBtcProof(request BtcProofRequest) (string, error) {
	panic("implement me")
}

func (h *Handler) GenEthProof(request EthProofRequest) (string, error) {
	panic("implement me")
}

func (h *Handler) ProofStatus(proofId string) (string, error) {
	panic("implement me")
}
