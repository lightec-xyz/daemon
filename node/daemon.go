package node

type IAgent interface {
	ScanBlock(height int64) (int64, error)
	Init() error
	Close() error
}

type Daemon struct {
	nodeList []IAgent
}

func NewNode(nodeList []IAgent) *Daemon {
	return &Daemon{
		nodeList: nodeList,
	}
}
