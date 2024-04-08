package proof

import (
	"fmt"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/store"
	"os"
	"os/signal"
	"syscall"
)

type Node struct {
	server  *rpc.Server
	handler *Handler
}

func NewNode(cfg Config) (*Node, error) {
	err := logger.InitLogger()
	if err != nil {
		logger.Error("init logger error:%v", err)
		return nil, err
	}
	host := fmt.Sprintf("%v:%v", cfg.RpcBind, cfg.RpcPort)
	memoryStore := store.NewMemoryStore()
	handler := NewHandler(memoryStore, cfg.MaxNums)
	logger.Info("proof worker info: %v", cfg.Info())
	server, err := rpc.NewWsServer(RpcRegisterName, host, handler)
	if err != nil {
		logger.Error("new rpc server error:%v", err)
		return nil, err
	}
	return &Node{
		server: server,
	}, nil
}

func (node *Node) Start() error {
	go node.server.Run()
	logger.Info("proof worker node start now ....")
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT, syscall.SIGTSTP, syscall.SIGQUIT)
	for {
		msg := <-ch
		switch msg {
		case syscall.SIGHUP:
			logger.Info("node get SIGHUP")
		case syscall.SIGQUIT, syscall.SIGTERM:
			logger.Info("get shutdown signal ...")
			err := node.Close()
			if err != nil {
				logger.Error(" node close info error:%v", err)
			}
			return err
		}
	}
}

func (node *Node) Close() error {
	if node.server != nil {
		err := node.server.Shutdown()
		if err != nil {
			logger.Error(" proof worker node exit now: %v", err)
		}
		return err
	}
	return nil
}
