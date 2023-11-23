package node

import (
	"github.com/lightec-xyz/daemon/rpc/bitcoin"
	"github.com/lightec-xyz/daemon/store"
	"time"
)

type BitcoinAgent struct {
	client     *bitcoin.Client
	store      *store.Store
	blockTime  time.Duration
	exitSignal chan struct{}
}

func NewBitcoinAgent(cfg BtcConfig) *BitcoinAgent {

	return &BitcoinAgent{}
}

func (b *BitcoinAgent) Init() error {
	return nil
}

func (b *BitcoinAgent) Run() error {
	ticker := time.NewTicker(b.blockTime)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			// TODO
		case <-b.exitSignal:
			return nil
		}
	}

}

func (b *BitcoinAgent) Close() error {
	panic(b)
}
func (b *BitcoinAgent) Name() string {
	return b.Name()
}
