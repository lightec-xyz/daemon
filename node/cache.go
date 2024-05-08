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
