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
	ethClient, err := ethereum.NewClient(cfg.EthUrl, cfg.ZkBridgeAddr)
	if err != nil {
		logger.Error("new eth btcClient error:%v", err)
		return nil, err
	}
	dbPath := fmt.Sprintf("%s/%s", cfg.DataDir, cfg.Network)
	storeDb, err := store.NewStore(dbPath, 1000, 10000, "zkbtc", false)
	if err != nil {
		logger.Error("new store error:%v,dbPath:%s", err, dbPath)
		return nil, err
	}
	memoryStore := store.NewMemoryStore()

	proofRequest := make(chan []ProofRequest, 10000)
	btcProofResp := make(chan ProofResponse, 1000)
	ethProofResp := make(chan ProofResponse, 1000)
	//btcAgent, err := NewBitcoinAgent(cfg, storeDb, memoryStore, btcClient, ethClient, proofRequest, btcProofResp)
	//if err != nil {
	//	logger.Error(err.Error())
	//	return nil, err
	//}
	ethAgent, err := NewEthereumAgent(cfg, storeDb, memoryStore, btcClient, ethClient, proofRequest, ethProofResp)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	//todo
	//workers, err := NewWorkers(cfg.Workers)
	//if err != nil {
	//	logger.Error("new workers error:%v", err)
	//	return nil, err
	//}
	workers := []IWorker{NewLocalWorker(1)}
	manager := NewManager(proofRequest, btcProofResp, ethProofResp, storeDb, memoryStore, NewSchedule(workers...))
	rpcHandler := NewHandler(storeDb, memoryStore)
	server, err := rpc.NewServer(fmt.Sprintf("%s:%s", cfg.Rpcbind, cfg.RpcPort), rpcHandler)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	daemon := &Daemon{
		agents:         []IAgent{ethAgent},
		server:         server,
		nodeConfig:     cfg,
		exitScanSignal: make(chan struct{}, 1),
		manager:        manager,
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
	ch := make(chan os.Signal, 1)
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
