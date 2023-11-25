package proof

import (
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"os"
	"os/signal"
	"syscall"
)

type Node struct {
	server *rpc.Server
}

func NewNode() (*Node, error) {
	return &Node{}, nil
}

func (node *Node) Start() error {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)
	for {
		msg := <-ch
		switch msg {
		case syscall.SIGHUP:
			logger.Info("daemon get SIGHUP")

		case syscall.SIGQUIT:
			fallthrough
		case syscall.SIGTERM:
			logger.Info("get shutdown sigterm...")
			err := node.Close()
			if err != nil {
				logger.Error(err.Error())
			}
		}
	}
}

func (node *Node) Close() error {
	if node.server != nil {
		err := node.server.Shutdown()
		if err != nil {
			logger.Error("proof server shutdown error:%v", err)
		}
	}
	return nil
}
