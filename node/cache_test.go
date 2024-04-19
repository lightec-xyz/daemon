package node

import "testing"

func TestCache_test(t *testing.T) {
	cacheState := NewCacheState()
	cacheState.Store(1, "dsdfsdfsd")
	value, ok := cacheState.Get(1)
	t.Log(value, ok)
}
