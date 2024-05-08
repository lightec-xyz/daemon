package proof

import "github.com/lightec-xyz/daemon/store"

const workerIdKey = "workerIdKey"

func ReadWorkerId(store store.IStore) (string, bool, error) {
	exists, err := store.HasObj(workerIdKey)
	if err != nil {
		return "", false, err
	}
	if !exists {
		return "", false, nil
	}
	var id string
	err = store.GetObj(workerIdKey, &id)
	if err != nil {
		return "", false, err
	}
	return id, true, nil
}

func WriteWorkerId(store store.IStore, id string) error {
	return store.PutObj(workerIdKey, id)
}
