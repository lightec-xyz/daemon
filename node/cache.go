package node

import (
	"sync"
)

type CacheState struct {
	requests *sync.Map
}

func NewCacheState() *CacheState {
	return &CacheState{
		requests: new(sync.Map),
	}
}

func (cs *CacheState) Store(key, value interface{}) {
	cs.requests.Store(key, value)
}

func (cs *CacheState) Check(key interface{}) bool {
	_, ok := cs.requests.Load(key)
	return ok
}
func (cs *CacheState) Get(key interface{}) (interface{}, bool) {
	value, ok := cs.requests.Load(key)
	if !ok {
		return nil, false
	}
	return value, true
}

func (cs *CacheState) Delete(key interface{}) {
	cs.requests.Delete(key)
}

// todo

type State struct {
	txes         map[string]string            // txHash -> redeem
	txSlotMap    map[uint64]map[string]string // txSlot -> map
	finalizeSlot map[uint64]map[string]string // finalized -> map
	lock         sync.Mutex                   // todo
	proofs       *sync.Map
}

func NewState() *State {
	return &State{
		txes:         make(map[string]string),
		txSlotMap:    make(map[uint64]map[string]string),
		finalizeSlot: make(map[uint64]map[string]string),
		proofs:       new(sync.Map),
	}
}

func (rs *State) DeleteId(id string) {
	rs.proofs.Delete(id)
}

func (rs *State) CheckId(id string) bool {
	_, ok := rs.proofs.Load(id)
	return ok
}

func (rs *State) StoreId(id string) {
	rs.proofs.Store(id, "")
}

func (rs *State) GetFinalizeSlot(slot uint64) map[string]string {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	return rs.finalizeSlot[slot]
}

func (rs *State) CheckFinalizeSlot(slot uint64) bool {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	_, ok := rs.finalizeSlot[slot]
	return ok
}

func (rs *State) AddFinalizeSlot(slot uint64, hash string) {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	list, ok := rs.finalizeSlot[slot]
	if ok {
		list[hash] = ""
		rs.finalizeSlot[slot] = list
	} else {
		finalizeSlotMap := make(map[string]string)
		finalizeSlotMap[hash] = ""
		rs.finalizeSlot[slot] = finalizeSlotMap
	}
}
func (rs *State) DeleteFinalizeSlot(slot uint64) {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	delete(rs.finalizeSlot, slot)
}

func (rs *State) GetTxSlot(txSlot uint64) map[string]string {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	return rs.txSlotMap[txSlot]
}

func (rs *State) CheckTxSlot(txSlot uint64) bool {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	_, ok := rs.txSlotMap[txSlot]
	return ok

}
func (rs *State) DeleteTxSlot(txSlot uint64) {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	delete(rs.txSlotMap, txSlot)
}

func (rs *State) AddTxSlot(txSlot uint64, hash string) {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	list, ok := rs.txSlotMap[txSlot]
	if ok {
		list[hash] = ""
		rs.txSlotMap[txSlot] = list
	} else {
		txSlotMap := make(map[string]string)
		txSlotMap[hash] = ""
		rs.txSlotMap[txSlot] = txSlotMap
	}

}

func (rs *State) CheckTx(hash string) bool {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	_, ok := rs.txes[hash]
	return ok
}

func (rs *State) DeleteTx(hash string) {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	delete(rs.txes, hash)
}

func (rs *State) AddTx(hash string) {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	rs.txes[hash] = ""
}
