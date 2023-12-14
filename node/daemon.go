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
	ScanBlock() error
	Transfer()
	Init() error
	Close() error
	Name() string
	BlockTime() time.Duration
}

type Daemon struct {
	agents         []IAgent
	server         *rpc.Server
	nodeConfig     NodeConfig
	exitScanSignal chan struct{}
	manager        *Manager
	exitSignal     chan os.Signal
}

func NewDaemon(cfg NodeConfig) (*Daemon, error) {
	err := logger.InitLogger()
	if err != nil {
		logger.Error("init logger error:%v", err)
		return nil, err
	}
	btcClient, err := bitcoin.NewClient(cfg.BtcUrl, cfg.BtcUser, cfg.BtcPwd, cfg.BtcNetwork)
	if err != nil {
		logger.Error("new btc btcClient error:%v", err)
		return nil, err
	}
	ethClient, err := ethereum.NewClient(cfg.EthUrl, cfg.ZkBridgeAddr, cfg.ZkBtcAddr)
	if err != nil {
		logger.Error("new eth btcClient error:%v", err)
		return nil, err
	}
	dbPath := fmt.Sprintf("%s/%s", cfg.DataDir, cfg.Network)
	logger.Info("dbPath:%s", dbPath)
	storeDb, err := store.NewStore(dbPath, 0, 0, "zkbtc", false)
	if err != nil {
		logger.Error("new store error:%v,dbPath:%s", err, dbPath)
		return nil, err
	}
	memoryStore := store.NewMemoryStore()
	proofRequest := make(chan []ProofRequest, 10000)
	btcProofResp := make(chan ProofResponse, 1000)
	ethProofResp := make(chan ProofResponse, 1000)
	var agents []IAgent
	btcAgent, err := NewBitcoinAgent(cfg, storeDb, memoryStore, btcClient, ethClient, proofRequest, btcProofResp)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	agents = append(agents, btcAgent)
	ethAgent, err := NewEthereumAgent(cfg, storeDb, memoryStore, btcClient, ethClient, proofRequest, ethProofResp)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	agents = append(agents, ethAgent)


	workers := make([]IWorker, 1)
	if cfg.EnableLocalWorker {
		logger.Info("local worker enable")
		workers = append(workers, NewLocalWorker(1))
	}
	schedule := NewSchedule(workers...)
	manager := NewManager(proofRequest, btcProofResp, ethProofResp, storeDb, memoryStore, schedule)
	exitSignal := make(chan os.Signal, 1)

	// todo new store
	rpcHandler := NewHandler(storeDb, memoryStore, schedule, exitSignal)
	server, err := rpc.NewServer(RpcRegisterName, fmt.Sprintf("%s:%s", cfg.Rpcbind, cfg.RpcPort), rpcHandler)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	daemon := &Daemon{
		agents:         agents,
		server:         server,
		nodeConfig:     cfg,
		exitScanSignal: make(chan struct{}, 1),
		manager:        manager,
		exitSignal:     exitSignal,
	}
	return daemon, nil
}

func (d *Daemon) Init() error {
	for _, node := range d.agents {
		if err := node.Init(); err != nil {
			logger.Error("%v:init node error %v", node.Name(), err)
			return err
		}
		go node.Transfer()
	}
	go d.manager.run()
	go d.manager.genProof()
	return nil
}

func (d *Daemon) Run() error {
	for _, node := range d.agents {
		go func(tNode IAgent) {
			//todo
			ticker := time.NewTicker(tNode.BlockTime())
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					err := tNode.ScanBlock()
					if err != nil {
						logger.Error("%v run error %v", tNode.Name(), err)
					}
				case <-d.exitScanSignal:
					logger.Info("exit scan block goroutine: %v", tNode.Name())
					return
				}
			}
		}(node)
	}
	signal.Notify(d.exitSignal, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT, syscall.SIGTSTP, syscall.SIGQUIT)
	for {
		msg := <-d.exitSignal
		switch msg {
		case syscall.SIGHUP:
			logger.Info("daemon get SIGHUP")
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGTSTP:
			logger.Info("get shutdown signal ,exit now ...")
			err := d.Close()
			if err != nil {
				logger.Error(err.Error())
			}
			return nil
		}
	}
}

func (d *Daemon) Close() error {
	if d.exitScanSignal != nil {
		close(d.exitScanSignal)
	}
	time.Sleep(2 * time.Second)
	for _, node := range d.agents {
		if err := node.Close(); err != nil {
			logger.Error("%v:close node error %v", node.Name(), err)
		}
	}
	d.manager.Close()
	err := d.server.Shutdown()
	if err != nil {
		logger.Error("rpc server shutdown error:%v", err)
	}
	return nil
}

//todo local worker ?

func NewWorkers(workers []WorkerConfig) ([]IWorker, error) {
	workersList := make([]IWorker, 0)
	for _, cfg := range workers {
		client, err := rpc.NewProofClient(cfg.ProofUrl)
		if err != nil {
			logger.Error("new worker error:%v", err)
			return nil, err
		}
		worker := NewWorker(client, cfg.ParallelNums)
		workersList = append(workersList, worker)
	}
	return workersList, nil
}
