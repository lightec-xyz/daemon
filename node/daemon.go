package node

import (
	"fmt"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/rpc/beacon"
	"github.com/lightec-xyz/daemon/rpc/bitcoin"
	"github.com/lightec-xyz/daemon/rpc/ethereum"
	"github.com/lightec-xyz/daemon/store"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	err := logger.InitLogger()
	if err != nil {
		panic(err)
	}
}

type IAgent interface {
	ScanBlock() error
	Transfer(resp ProofResponse) error
	Init() error
	Close() error
	Name() string
}

type IBeaconAgent interface {
	ScanSycnPeriod() error
	SaveSyncCommitteeProof(resp rpc.SyncCommitteeProofResponse) error
	Init() error
	Close() error
	Name() string
}

type Daemon struct {
	agents      []*Agent
	beaconAgent *BeaconAgent
	server      *rpc.Server
	nodeConfig  NodeConfig
	exitSignal  chan struct{}
	manager     *WrapperManger
}

func NewDaemon(cfg NodeConfig) (*Daemon, error) {
	var submitTxEthAddr string
	var err error
	if cfg.AutoSubmit {
		submitTxEthAddr, err = privateKeyToEthAddr(cfg.EthPrivateKey)
		if err != nil {
			logger.Error("privateKeyToEthAddr error:%v", err)
			return nil, err
		}
		logger.Info("ethereum submit address:%v", submitTxEthAddr)
	}

	btcClient, err := bitcoin.NewClient(cfg.BtcUrl, cfg.BtcUser, cfg.BtcPwd, cfg.BtcNetwork)
	if err != nil {
		logger.Error("new btc btcClient error:%v", err)
		return nil, err
	}

	beaconClient, err := beacon.NewClient(cfg.BeaconUrl)
	if err != nil {
		logger.Error("new beacon btcClient error:%v", err)
		return nil, err
	}
	//TODO(keep), should be replaced with actual url
	proofClient, err := rpc.NewSyncCommitteeProofClient("http://127.0.0.1:8980")
	if err != nil {
		logger.Error("new proofClient error:%v", err)
		return nil, err
	}

	beaconAgent, err := NewBeaconAgent(cfg, beaconClient, proofClient)
	if err != nil{
		logger.Error("new beacon btcClient error:%v", err)
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
	proofRequest := make(chan []ProofRequest, 1000)
	btcProofResp := make(chan ProofResponse, 1000)
	ethProofResp := make(chan ProofResponse, 1000)
	keyStore := NewKeyStore(cfg.EthPrivateKey)
	var agents []*Agent
	btcAgent, err := NewBitcoinAgent(cfg, submitTxEthAddr, storeDb, memoryStore, btcClient, ethClient, proofRequest, keyStore)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	agents = append(agents, NewAgent(btcAgent, cfg.BtcScanBlockTime, btcProofResp))
	ethAgent, err := NewEthereumAgent(cfg, submitTxEthAddr, storeDb, memoryStore, btcClient, ethClient, proofRequest)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	agents = append(agents, NewAgent(ethAgent, cfg.EthScanBlockTime, ethProofResp))

	workers := make([]IWorker, 0)
	if cfg.EnableLocalWorker {
		logger.Info("local worker enable")
		workers = append(workers, NewLocalWorker(1))
	}
	schedule := NewSchedule(workers...)
	manager, err := NewManager(cfg, btcProofResp, ethProofResp, storeDb, memoryStore, schedule)
	if err != nil {
		logger.Error("new manager error: %v", err)
		return nil, err
	}
	exitSignal := make(chan os.Signal, 1)
	// todo new store
	rpcHandler := NewHandler(storeDb, memoryStore, schedule, exitSignal)
	server, err := rpc.NewServer(RpcRegisterName, fmt.Sprintf("%s:%s", cfg.Rpcbind, cfg.RpcPort), rpcHandler)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	daemon := &Daemon{
		agents: agents,
		beaconAgent:
		server:     server,
		nodeConfig: cfg,
		exitSignal: make(chan struct{}, 1),
		manager:    NewWrapperManger(manager, proofRequest),
	}
	return daemon, nil
}

func (d *Daemon) Init() error {
	for _, agent := range d.agents {
		if err := agent.node.Init(); err != nil {
			logger.Error("%v:init agent error %v", agent.node.Name(), err)
			return err
		}
	}
	err := d.manager.manager.init()
	if err != nil {
		logger.Error("manager init error %v", err)
		return err
	}

	return nil
}

func (d *Daemon) Run() error {

	go func() {
		proofRequest := d.manager.proofRequest
		for {
			select {
			case <-d.exitSignal:
				logger.Info("manager proof queue exit ...")
				return
			case requests := <-proofRequest:
				err := d.manager.manager.run(requests)
				if err != nil {
					logger.Error("manager run error %v", err)
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case <-d.exitSignal:
				logger.Info("manager generate proof exit ...")
				return
			default:
				err := d.manager.manager.genProof()
				if err != nil {
					logger.Error("manager gen proof error %v", err)

				}
			}
		}
	}()

	for _, agent := range d.agents {
		go func(tAgent *Agent) {
			proofResponses := tAgent.proofResp
			for {
				select {
				case response := <-proofResponses:
					err := tAgent.node.Transfer(response)
					if err != nil {
						logger.Error("transfer error %v", err)
					}
				case <-d.exitSignal:
					logger.Error("%v transfer goroutine exit", tAgent.node.Name())
				}
			}
		}(agent)
	}
	for _, agent := range d.agents {
		go func(tAgent *Agent) {
			ticker := time.NewTicker(tAgent.scanTime)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					err := tAgent.node.ScanBlock()
					if err != nil {
						logger.Error("%v run error %v", tAgent.node.Name(), err)
					}
				case <-d.exitSignal:
					logger.Info("%v scan block goroutine exit ", tAgent.node.Name())
					return
				}
			}
		}(agent)
	}
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT, syscall.SIGTSTP, syscall.SIGQUIT)
	for {
		msg := <-ch
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
	for _, agent := range d.agents {
		if err := agent.node.Close(); err != nil {
			logger.Error("%v:close agent error %v", agent.node.Name(), err)
		}
	}
	d.manager.manager.Close()
	err := d.server.Shutdown()
	if err != nil {
		logger.Error("rpc server shutdown error:%v", err)
	}
	if d.exitSignal != nil {
		close(d.exitSignal)
	}
	return nil
}

//todo local worker ?

func NewWorkers(workers []WorkerConfig) ([]IWorker, error) {
	workersList := make([]IWorker, 0)
	for _, cfg := range workers {
		client, err := rpc.NewProofClient(cfg.Url)
		if err != nil {
			logger.Error("new worker error:%v", err)
			return nil, err
		}
		worker := NewWorker(client, cfg.ParallelNums)
		workersList = append(workersList, worker)
	}
	return workersList, nil
}

func NewAgent(agent IAgent, scanTime time.Duration, proofResp chan ProofResponse) *Agent {
	return &Agent{
		node:      agent,
		scanTime:  scanTime,
		proofResp: proofResp,
	}
}

type WrapperManger struct {
	manager      *manager
	proofRequest chan []ProofRequest
}

func NewWrapperManger(manager *manager, request chan []ProofRequest) *WrapperManger {
	return &WrapperManger{
		manager:      manager,
		proofRequest: request,
	}
}

type Agent struct {
	node      IAgent
	scanTime  time.Duration
	proofResp chan ProofResponse
}


