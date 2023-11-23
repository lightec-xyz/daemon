package node

import "github.com/lightec-xyz/daemon/store"

func getCurrentHeight(store *store.Store, key string, value interface{}) error {
	return store.GetObj(key, value)
}
