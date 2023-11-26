package proof

import "github.com/lightec-xyz/daemon/store"

type Manager struct {
	store       store.IStore
	memoryStore store.IStore
}

func NewManager(store store.IStore, memoryStore store.IStore) *Manager {
	return &Manager{
		store:       store,
		memoryStore: memoryStore,
	}
}
