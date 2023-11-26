package node

import (
	"fmt"
	"github.com/lightec-xyz/daemon/store"
)

var _ API = (*Handler)(nil)

type Handler struct {
	//todo
	store    store.IStore
	memoryDb store.IStore
}

func NewHandler(store store.IStore, memoryDb store.IStore) *Handler {
	return &Handler{
		store:    store,
		memoryDb: memoryDb,
	}
}

func (h *Handler) Version() (DaemonInfo, error) {
	return DaemonInfo{}, nil
}

func (h *Handler) HelloWorld(name *string, age *int) (string, error) {
	fmt.Printf("req: %v %v", name, age)
	return fmt.Sprintf(" name %v age %v", *name, *age), nil
}
