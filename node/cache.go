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

type RedeemState struct {
	txes         map[string]string            // txHash -> redeem
	txSlotMap    map[uint64]map[string]string // txSlot -> map
	finalizeSlot map[uint64]map[string]string // finalized -> map
	lock         sync.Mutex                   // todo
}

func NewRedeemState() *RedeemState {
	return &RedeemState{
		txes:         make(map[string]string),
		txSlotMap:    make(map[uint64]map[string]string),
		finalizeSlot: make(map[uint64]map[string]string),
	}
}

func (rs *RedeemState) GetFinalizeSlot(slot uint64) map[string]string {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	return rs.finalizeSlot[slot]
}

func (rs *RedeemState) CheckFinalizeSlot(slot uint64) bool {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	_, ok := rs.finalizeSlot[slot]
	return ok
}

func (rs *RedeemState) AddFinalizeSlot(slot uint64, hash string) {
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
func (rs *RedeemState) DeleteFinalizeSlot(slot uint64) {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	delete(rs.finalizeSlot, slot)
}

func (rs *RedeemState) GetTxSlot(txSlot uint64) map[string]string {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	return rs.txSlotMap[txSlot]
}

func (rs *RedeemState) CheckTxSlot(txSlot uint64) bool {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	_, ok := rs.txSlotMap[txSlot]
	return ok

}
func (rs *RedeemState) DeleteTxSlot(txSlot uint64) {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	delete(rs.txSlotMap, txSlot)
}

func (rs *RedeemState) AddTxSlot(txSlot uint64, hash string) {
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

func (rs *RedeemState) CheckTx(hash string) bool {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	_, ok := rs.txes[hash]
	return ok
}

func (rs *RedeemState) DeleteTx(hash string) {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	delete(rs.txes, hash)
}

func (rs *RedeemState) AddTx(hash string) {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	rs.txes[hash] = ""
}
