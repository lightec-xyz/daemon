package node

import "sync"

type NonceManager struct {
	sync.Mutex
}

func NewNonceManager() *NonceManager {
	return &NonceManager{}
}

type DepositTask struct {
	Nonce int
	Id    string
}

type RedeemTask struct {
	Nonce int
	Id    string
}
