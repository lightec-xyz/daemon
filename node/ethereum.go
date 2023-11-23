package node

type EthereumNode struct {
}

func NewEthereumNode() *EthereumNode {
	return &EthereumNode{}
}

func (e *EthereumNode) Init() error {

}

func (e *EthereumNode) ScanBlock(height int64) (int64, error) {

}

func (e *EthereumNode) Close() error {

}
