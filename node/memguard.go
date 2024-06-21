package node

import (
	"fmt"
	"github.com/awnumar/memguard"
	"sync"
)

type Memguard struct {
	lock        sync.Mutex
	memguardMap map[string]*memguard.Enclave
}

func NewMemguard() *Memguard {
	// todo
	//memguard.CatchInterrupt()
	return &Memguard{
		memguardMap: make(map[string]*memguard.Enclave),
	}
}

func (m *Memguard) Store(key string, value []byte) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.memguardMap[key] = memguard.NewEnclave(value)
}

func (m *Memguard) Load(key string) ([]byte, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	value, ok := m.memguardMap[key]
	if !ok {
		return nil, fmt.Errorf("key %s not found", key)
	}
	lockedBuffer, err := value.Open()
	if err != nil {
		return nil, err
	}
	defer lockedBuffer.Destroy()
	dest := make([]byte, len(lockedBuffer.Bytes()))
	copy(dest, lockedBuffer.Bytes())
	return dest, nil
}

func (m *Memguard) Check(key string) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	_, ok := m.memguardMap[key]
	return ok
}

func (m *Memguard) Delete(key string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.memguardMap, key)
}

func (m *Memguard) Close() {
	m.memguardMap = nil
	memguard.Purge()
}
