package node

import (
	"sync"
)

type cache struct {
	requests *sync.Map
}

func NewCacheState() *cache {
	return &cache{
		requests: new(sync.Map),
	}
}

func (cs *cache) Store(key, value interface{}) {
	cs.requests.Store(key, value)
}

func (cs *cache) Check(key interface{}) bool {
	_, ok := cs.requests.Load(key)
	return ok
}
func (cs *cache) Get(key interface{}) (interface{}, bool) {
	value, ok := cs.requests.Load(key)
	if !ok {
		return nil, false
	}
	return value, true
}

func (cs *cache) Delete(key interface{}) {
	cs.requests.Delete(key)
}
