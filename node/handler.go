package node

import (
	"fmt"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/store"
)

var _ rpc.NodeAPI = (*Handler)(nil)

type Handler struct {
	//todo
	store       store.IStore
	memoryDb    store.IStore
	proofClient rpc.ProofAPI
}

func (h *Handler) Version() (rpc.NodeInfo, error) {
	daemonInfo := rpc.NodeInfo{}
	return daemonInfo, nil
}

func NewHandler(store, memoryDb store.IStore, client rpc.ProofAPI) *Handler {
	return &Handler{
		store:       store,
		memoryDb:    memoryDb,
		proofClient: client,
	}
}

func (h *Handler) HelloWorld(name *string, age *int) (string, error) {
	fmt.Printf("req: %v %v", name, age)
	return fmt.Sprintf(" name %v age %v", *name, *age), nil
}
