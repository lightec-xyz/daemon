package node

type BitcoinNode struct {
}

func NewBitcoin() *BitcoinNode {
	return &BitcoinNode{}
}

func (b *BitcoinNode) Init() error {
	return nil
}

func (b *BitcoinNode) ScanBlock(height int64) (int64, error) {

}

func (b *BitcoinNode) Close() error {

}
