package node

import (
	"github.com/lightec-xyz/daemon/codec"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/store"
)

type MemoryStore struct {
	store store.IStore
}

func NewMemoryStore(store store.IStore) *MemoryStore {
	return &MemoryStore{store: store}
}

func (s *MemoryStore) WriteFinalityUpdateSlot(slot uint64) error {
	return s.store.PutObj(dbFinalityUpdateSlotId(slot), slot)
}

func (s *MemoryStore) FindFinalityUpdateNearestSlot(txSlot uint64) (uint64, bool, error) {
	iterator := s.store.Iterator([]byte(finalityUpdateSlotPrefix), nil)
	defer iterator.Release()
	defer iterator.Release()
	for iterator.Next() {
		var slot uint64
		err := codec.Unmarshal(iterator.Value(), &slot)
		if err != nil {
			return 0, false, err
		}
		if slot >= txSlot {
			return slot, slot-txSlot <= common.MaxDiffTxFinalitySlot, nil
		}
	}
	return 0, false, nil
}
