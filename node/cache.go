package node

import (
	"sync"
)

type Cache struct {
	requests *sync.Map
}

func NewCacheState() *Cache {
	return &Cache{
		requests: new(sync.Map),
	}
}

func (cs *Cache) Store(key, value interface{}) {
	cs.requests.Store(key, value)
}

func (cs *Cache) Check(key interface{}) bool {
	_, ok := cs.requests.Load(key)
	return ok
}
func (cs *Cache) Get(key interface{}) (interface{}, bool) {
	value, ok := cs.requests.Load(key)
	if !ok {
		return nil, false
	}
	return value, true
}

func (cs *Cache) Delete(key interface{}) {
	cs.requests.Delete(key)
}
