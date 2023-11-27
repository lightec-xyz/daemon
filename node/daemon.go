package node

import (
	"fmt"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/rpc/bitcoin"
	"github.com/lightec-xyz/daemon/rpc/ethereum"
	"github.com/lightec-xyz/daemon/store"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type IAgent interface {
	Run() error
	Init() error
	Close() error
	Name() string
	BlockTime() time.Duration
}

type Daemon struct {
	agents     []IAgent
	server     *rpc.Server
	nodeConfig NodeConfig
	exitSignal chan struct{}
}

func NewDaemon(cfg Config) (*Daemon, error) {
	err := logger.InitLogger()
	if err != nil {
		logger.Error("init logger error:%v", err)
		return nil, err
	}
	btcClient, err := bitcoin.NewClient(cfg.NodeConfig.BtcUrl, cfg.NodeConfig.BtcUser, cfg.NodeConfig.BtcPwd, cfg.NodeConfig.BtcNetwork)
	if err != nil {
		logger.Error("new btc btcClient error:%v", err)
		return nil, err
	}
	ethClient, err := ethereum.NewClient(cfg.NodeConfig.EthUrl)
	if err != nil {
		logger.Error("new eth btcClient error:%v", err)
		return nil, err
	}
	proorClient, err := rpc.NewProofClient(cfg.NodeConfig.ProofUrl)
	if err != nil {
		logger.Error("new proofClient error:%v", err)
		return nil, err
	}
	//todo
	dbPath := fmt.Sprintf("%s/%s", cfg.NodeConfig.DataDir, cfg.NodeConfig.Network)
	storeDb, err := store.NewStore(dbPath, cfg.DbConfig.Cache, cfg.DbConfig.Handler, "zkbtc", false)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	memoryStore := store.NewMemoryStore()
	btcAgent, err := NewBitcoinAgent(cfg.NodeConfig, storeDb, memoryStore, btcClient, ethClient, proorClient)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	ethAgent, err := NewEthereumAgent(cfg.NodeConfig, storeDb, memoryStore, btcClient, ethClient, proorClient)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	rpcHandler := NewHandler(storeDb, memoryStore, proorClient)
	server, err := rpc.NewServer(fmt.Sprintf("%s:%s", cfg.SeverConfig.IP, cfg.SeverConfig.Port), rpcHandler)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	daemon := &Daemon{
		agents:     []IAgent{btcAgent, ethAgent},
		server:     server,
		nodeConfig: cfg.NodeConfig,
		exitSignal: make(chan struct{}, 1),
	}
	err = daemon.Init()
	if err != nil {
		logger.Error("daemon init error:%v", err)
		return nil, err
	}
	return daemon, nil
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

func (d *Daemon) Run() error {
	ch := make(chan os.Signal, 1)
	for _, node := range d.agents {
		go func(tNode IAgent) {
			//todo
			ticker := time.NewTicker(tNode.BlockTime())
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					err := tNode.Run()
					if err != nil {
						logger.Error("% run error %v", tNode.Name(), err)
					}
				case <-d.exitSignal:
					return
				}
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

func (d *Daemon) Close() error {
	for _, node := range d.agents {
		if err := node.Close(); err != nil {
			logger.Error("%v:close node error %v", node.Name(), err)
			//need continue,close next node
		}
	}
	close(d.exitSignal)
	err := d.server.Shutdown()
	if err != nil {
		logger.Error("server shutdown error:%v", err)
	}
	return nil
}
