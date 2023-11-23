package node

import "github.com/lightec-xyz/daemon/logger"

type IAgent interface {
	ScanBlock(height int64) (int64, error)
	Init() error
	Close() error
	Name() string
}

type Daemon struct {
	nodeList []IAgent
}

func NewNode(nodeList []IAgent) *Daemon {
	return &Daemon{
		nodeList: nodeList,
	}
}

func (d *Daemon) Init() error {
	for _, node := range d.nodeList {
		if err := node.Init(); err != nil {
			logger.Error("%v:init node error %v", node.Name(), err)
			return err
		}
	}
	return nil
}

func (d *Daemon) Close() error {
	for _, node := range d.nodeList {
		if err := node.Close(); err != nil {
			logger.Error("%v:close node error %v", node.Name(), err)
			//
		}
	}
	return nil
}
