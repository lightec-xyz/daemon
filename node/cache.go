package node

import (
	"github.com/lightec-xyz/daemon/common"
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

func (cs *CacheState) StoreZkRequest(key, value interface{}) {
	cs.requests.Store(key, value)
}

func (cs *CacheState) CheckZkRequest(key interface{}) bool {
	_, ok := cs.requests.Load(key)
	return ok
}
func (cs *CacheState) GetZkRequest(key interface{}) (*common.ZkProofRequest, bool) {
	value, ok := cs.requests.Load(key)
	if !ok {
		return nil, false
	}
	req, ok := value.(*common.ZkProofRequest)
	if !ok {
		return nil, false
	}
	return req, true
}

func (cs *CacheState) DeleteZkRequest(key interface{}) {
	cs.requests.Delete(key)
}
