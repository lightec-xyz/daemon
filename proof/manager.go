package proof

import "github.com/lightec-xyz/daemon/store"

type Manager struct {
	store       *store.Store
	memoryStore *store.MemoryStore
}

func NewManager(store *store.Store, memoryStore *store.MemoryStore) *Manager {
	return &Manager{
		store:       store,
		memoryStore: memoryStore,
	}
}
