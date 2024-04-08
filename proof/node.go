package proof

import (
	"fmt"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/store"
	"os"
	"os/signal"
	"syscall"
)

// Node Todo
type Node struct {
	server *rpc.Server
	mode   common.Mode
	local  *Local
	exit   chan os.Signal
}

func NewNode(cfg Config) (*Node, error) {
	err := logger.InitLogger()
	if err != nil {
		logger.Error("init logger error:%v", err)
		return nil, err
	}
	err = cfg.Check()
	if err != nil {
		logger.Error("config check error:%v", err)
		return nil, err
	}
	if cfg.Mode == common.Client {
		local, err := NewLocal(cfg.Url, cfg.DataDir, cfg.MaxNums)
		if err != nil {
			return nil, err
		}
		return &Node{
			local: local,
			mode:  cfg.Mode,
			exit:  make(chan os.Signal, 1),
		}, nil
	} else if cfg.Mode == common.Cluster {
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
			mode:   cfg.Mode,
			exit:   make(chan os.Signal, 1),
		}, nil
	}
	return nil, fmt.Errorf("new node error: unknown model:%v", cfg.Mode)

}

func (node *Node) Start() error {
	if node.mode == common.Client {
		go node.local.Run()
	} else if node.mode == common.Cluster {
		go node.server.Run()
	}
	logger.Info("proof worker node start now ....")

	signal.Notify(node.exit, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT, syscall.SIGTSTP, syscall.SIGQUIT)
	for {
		msg := <-node.exit
		switch msg {
		case syscall.SIGHUP:
			logger.Info("node get SIGHUP")
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGTSTP:
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
	logger.Warn("proof worker node exit now ....")
	if node.server != nil {
		err := node.server.Shutdown()
		if err != nil {
			logger.Error(" proof worker node exit now: %v", err)
		}
	}
	if node.local != nil {
		err := node.local.Close()
		if err != nil {
			logger.Error(" proof worker node exit now: %v", err)
		}
	}
	return nil
}
