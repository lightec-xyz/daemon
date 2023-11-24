package node

import (
	"fmt"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/store"
	"os"
	"os/signal"
	"syscall"
)

type IAgent interface {
	Run() error
	Init() error
	Close() error
	Name() string
}

type Daemon struct {
	agents []IAgent
	server *rpc.Server
}

func NewDaemon(cfg Config) (*Daemon, error) {
	//todo
	err := logger.InitLogger()
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	storeDb, err := store.NewStore(cfg.DbConfig.Path, cfg.DbConfig.Cache, cfg.DbConfig.Handler, "zkbtc", false)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	memoryStore := store.NewMemoryStore()
	btcAgent, err := NewBitcoinAgent(cfg.Bitcoin, storeDb, memoryStore)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	ethAgent, err := NewEthereumAgent(cfg.Ethereum, storeDb, memoryStore)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	rpcHandler := NewHandler(storeDb)
	server, err := rpc.NewServer(fmt.Sprintf("%s:%s", cfg.SeverConfig.IP, cfg.SeverConfig.Port), rpcHandler)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return &Daemon{
		agents: []IAgent{btcAgent, ethAgent},
		server: server,
	}, nil
}

func (d *Daemon) Init() error {
	for _, node := range d.agents {
		if err := node.Init(); err != nil {
			logger.Error("%v:init node error %v", node.Name(), err)
			return err
		}
	}
	return nil
}

func (d *Daemon) Close() error {
	for _, node := range d.agents {
		if err := node.Close(); err != nil {
			logger.Error("%v:close node error %v", node.Name(), err)
			//need continue,close next node
		}
	}
	err := d.server.Shutdown()
	if err != nil {
		logger.Error("server shutdown error:%v", err)
	}
	return nil
}

func (d *Daemon) Run() error {
	err := d.Init()
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	ch := make(chan os.Signal, 1)

	for _, node := range d.agents {
		go func(cNode IAgent) {
			err := cNode.Run()
			if err != nil {
				logger.Error("%v run error %v", cNode.Name(), err)
				ch <- syscall.SIGQUIT
			}
		}(node)
	}
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
			err := d.Close()
			if err != nil {
				logger.Error(err.Error())
			}
		}
	}
}
