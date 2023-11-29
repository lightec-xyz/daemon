package proof

import (
	"fmt"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"os"
	"os/signal"
	"syscall"
)

type Node struct {
	server *rpc.Server
}

func NewNode(cfg Config) (*Node, error) {
	host := fmt.Sprintf("%v:%v", cfg.RpcBind, cfg.RpcPort)
	handler := NewHandler()
	server, err := rpc.NewWsServer(host, handler)
	if err != nil {
		logger.Error("new server error:%v", err)
		return nil, err
	}
	return &Node{
		server: server,
	}, nil
}

func (node *Node) Start() error {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)
	for {
		msg := <-ch
		switch msg {
		case syscall.SIGHUP:
			logger.Info("node get SIGHUP")

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
