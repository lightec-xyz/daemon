package node

type IAgent interface {
	ScanBlock(height int64) (int64, error)
	Init() error
}

type Node struct {
}

func NewNode() *Node {
	return &Node{}
}
