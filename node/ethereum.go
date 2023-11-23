package node

import (
	"github.com/lightec-xyz/daemon/rpc/ethereum"
	"github.com/lightec-xyz/daemon/store"
	"time"
)

type EthereumAgent struct {
	client     *ethereum.Client
	store      *store.Store
	name       string
	exitSignal chan struct{}
	blockTime  time.Duration
}

func NewEthereumAgent(cfg EthConfig, store *store.Store) (IAgent, error) {
	return &EthereumAgent{}, nil
}

func (e *EthereumAgent) Init() error {
	return nil
}

func (e *EthereumAgent) Run() error {
	ticker := time.NewTicker(e.blockTime)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			// TODO
		case <-e.exitSignal:
			return nil
		}
	}
}

func (e *EthereumAgent) Close() error {
	panic(e)
}
func (e *EthereumAgent) Name() string {
	return e.name
}
