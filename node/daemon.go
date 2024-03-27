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
	ProofResponse(resp ZkProofResponse) error
	Init() error
	Close() error
	Name() string
}

type IBeaconAgent interface {
	ScanSyncPeriod() error
	ProofResponse(resp ZkProofResponse) error
	FetchDataResponse(resp FetchDataResponse) error
	CheckData() error
	Init() error
	Close() error
	Name() string
}

type Daemon struct {
	agents        []*WrapperAgent
	beaconAgent   *WrapperBeacon
	server        *rpc.Server
	nodeConfig    NodeConfig
	exitSignal    chan os.Signal
	manager       *WrapperManger
	onlyRecursive bool // true ,Only enable the function of generating recursive proofs
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
		logger.Error("new node btcClient error:%v", err)
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
	proofRequest := make(chan []ZkProofRequest, 1000)
	btcProofResp := make(chan ZkProofResponse, 1000)
	ethProofResp := make(chan ZkProofResponse, 1000)
	syncCommitResp := make(chan ZkProofResponse, 1000)
	fetchDataResp := make(chan FetchDataResponse, 1000)

	beaconAgent, err := NewBeaconAgent(cfg, beaconClient, proofRequest, fetchDataResp)
	if err != nil {
		logger.Error("new node btcClient error:%v", err)
		return nil, err
	}

	keyStore := NewKeyStore(cfg.EthPrivateKey)
	var agents []*WrapperAgent
	btcAgent, err := NewBitcoinAgent(cfg, submitTxEthAddr, storeDb, memoryStore, btcClient, ethClient, proofRequest, keyStore)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	agents = append(agents, NewWrapperAgent(btcAgent, cfg.BtcScanBlockTime, btcProofResp))

	ethAgent, err := NewEthereumAgent(cfg, submitTxEthAddr, storeDb, memoryStore, btcClient, ethClient, proofRequest)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	agents = append(agents, NewWrapperAgent(ethAgent, cfg.EthScanBlockTime, ethProofResp))

	workers := make([]rpc.IWorker, 0)
	if cfg.EnableLocalWorker {
		logger.Info("local worker enable")
		localWorker, err := NewLocalWorker("", "", 1)
		if err != nil {
			logger.Error("new local worker error:%v", err)
			return nil, err
		}
		workers = append(workers, localWorker)
	}
	schedule := NewSchedule(workers)
	manager, err := NewManager(cfg, btcProofResp, ethProofResp, syncCommitResp, storeDb, memoryStore, schedule)
	if err != nil {
		logger.Error("new manager error: %v", err)
		return nil, err
	}
	exitSignal := make(chan os.Signal, 1)

	rpcHandler := NewHandler(storeDb, memoryStore, schedule, exitSignal)
	server, err := rpc.NewServer(RpcRegisterName, fmt.Sprintf("%s:%s", cfg.Rpcbind, cfg.RpcPort), rpcHandler)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	daemon := &Daemon{
		agents:      agents,
		server:      server,
		nodeConfig:  cfg,
		exitSignal:  make(chan os.Signal, 1),
		beaconAgent: NewWrapperBeacon(beaconAgent, 1*time.Minute, 1*time.Minute, syncCommitResp, fetchDataResp),
		manager:     NewWrapperManger(manager, proofRequest),
	}
	return daemon, nil
}

func (d *Daemon) Init() error {

	if !d.onlyRecursive {
		// true ,skip init node agent
		for _, agent := range d.agents {
			if err := agent.node.Init(); err != nil {
				logger.Error("%v:init agent error %v", agent.node.Name(), err)
				return err
			}
		}
	}
	err := d.manager.manager.init()
	if err != nil {
		logger.Error("manager init error %v", err)
		return err
	}
	err = d.beaconAgent.node.Init()
	if err != nil {
		logger.Error("node agent init error %v", err)
		return err
	}
	return nil
}

func (d *Daemon) Run() error {
	logger.Info("start daemon")
	// rpc server
	go d.server.Run()

	// syncCommit
	go doTimerTask("beacon-scanSyncPeriod", d.beaconAgent.scanPeriodTime, d.beaconAgent.node.ScanSyncPeriod, d.exitSignal)
	go doProofResponseTask("beacon-proofResponse", d.beaconAgent.proofResponse, d.beaconAgent.node.ProofResponse, d.exitSignal)
	go doFetchRespTask("beacon-fetchDataResponse", d.beaconAgent.fetchDataResponse, d.beaconAgent.node.FetchDataResponse, d.exitSignal)
	go doTimerTask("beacon-checkData", d.beaconAgent.checkDataTime, d.beaconAgent.node.CheckData, d.exitSignal)

	// generate Proof manager
	go doProofRequestTask("manager-proofRequest", d.manager.proofRequest, d.manager.manager.run, d.exitSignal)
	go doTask("manager-generateProof:", d.manager.manager.genProof, d.exitSignal)

	if !d.onlyRecursive {
		// true skip this step
		//tx Proof
		for _, agent := range d.agents {
			name := fmt.Sprintf("%s-submitProof", agent.node.Name())
			go doProofResponseTask(name, agent.proofResp, agent.node.ProofResponse, d.exitSignal)
		}
		//scan block with tx
		for _, agent := range d.agents {
			name := fmt.Sprintf("%s-scanBlock", agent.node.Name())
			go doTimerTask(name, agent.scanTime, agent.node.ScanBlock, d.exitSignal)
		}
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
	if !d.onlyRecursive {
		for _, agent := range d.agents {
			if err := agent.node.Close(); err != nil {
				logger.Error("%v:close agent error %v", agent.node.Name(), err)
			}
		}
	}
	err := d.server.Shutdown()
	if err != nil {
		logger.Error("rpc server shutdown error:%v", err)
	}
	d.manager.manager.Close()
	err = d.beaconAgent.node.Close()
	if err != nil {
		logger.Error("node agent close error:%v", err)
	}
	if d.exitSignal != nil {
		close(d.exitSignal)
		//todo waitGroup
		time.Sleep(5 * time.Second)
	}
	return nil
}

func NewWorkers(workers []WorkerConfig) ([]rpc.IWorker, error) {
	workersList := make([]rpc.IWorker, 0)
	for _, cfg := range workers {
		client, err := rpc.NewProofClient(cfg.Url)
		if err != nil {
			logger.Error("new worker error:%v", err)
			return nil, err
		}
		worker := NewWorker(client, cfg.MaxNums)
		workersList = append(workersList, worker)
	}
	return workersList, nil
}

type WrapperBeacon struct {
	node              IBeaconAgent
	scanPeriodTime    time.Duration // get node Period
	checkDataTime     time.Duration
	proofResponse     chan ZkProofResponse
	fetchDataResponse chan FetchDataResponse
}

func NewWrapperBeacon(beacon IBeaconAgent, scanPeriodTime, checkDataTime time.Duration, proofResponse chan ZkProofResponse, fetchDataResp chan FetchDataResponse) *WrapperBeacon {
	return &WrapperBeacon{
		node:              beacon,
		scanPeriodTime:    scanPeriodTime,
		proofResponse:     proofResponse,
		fetchDataResponse: fetchDataResp,
		checkDataTime:     checkDataTime, // todo
	}
}

type WrapperManger struct {
	manager      *manager
	proofRequest chan []ZkProofRequest
}

func NewWrapperManger(manager *manager, request chan []ZkProofRequest) *WrapperManger {
	return &WrapperManger{
		manager:      manager,
		proofRequest: request,
	}
}

type WrapperAgent struct {
	node      IAgent
	scanTime  time.Duration
	proofResp chan ZkProofResponse
}

func NewWrapperAgent(agent IAgent, scanTime time.Duration, proofResp chan ZkProofResponse) *WrapperAgent {
	return &WrapperAgent{
		node:      agent,
		scanTime:  scanTime,
		proofResp: proofResp,
	}
}

func doTask(name string, fn func() error, exit chan os.Signal) {
	logger.Info("%v goroutine start ...", name)
	for {
		select {
		case <-exit:
			logger.Info("%v goroutine exit now ...", name)
			return
		default:
			err := fn()
			if err != nil {
				logger.Error("%v error %v", name, err.Error())
			}
		}
	}
}

func doTimerTask(name string, interval time.Duration, fn func() error, exit chan os.Signal) {
	logger.Info("%v ticker goroutine start ...", name)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-exit:
			logger.Info("%v goroutine exit now ...", name)
			return
		case <-ticker.C:
			err := fn()
			if err != nil {
				logger.Error("%v error %v", name, err.Error())
			}
		}
	}
}

func doProofRequestTask(name string, req chan []ZkProofRequest, fn func(req []ZkProofRequest) error, exit chan os.Signal) {
	logger.Info("%v goroutine start ...", name)
	for {
		select {
		case <-exit:
			logger.Info("%v goroutine exit now ...", name)
			return
		case request := <-req:
			err := fn(request)
			if err != nil {
				logger.Error("%v error %v", name, err.Error())
			}
		}

	}
}

func doFetchRespTask(name string, resp chan FetchDataResponse, fn func(resp FetchDataResponse) error, exit chan os.Signal) {
	logger.Info("%v goroutine start ...", name)
	for {
		select {
		case <-exit:
			logger.Info("%v goroutine exit now ...", name)
			return
		case response := <-resp:
			err := fn(response)
			if err != nil {
				logger.Error("%v error %v", name, err.Error())
			}
		}
	}
}

func doProofResponseTask(name string, resp chan ZkProofResponse, fn func(resp ZkProofResponse) error, exit chan os.Signal) {
	logger.Info("%v goroutine start ...", name)
	for {
		select {
		case <-exit:
			logger.Info("%v goroutine exit now ...", name)
			return
		case response := <-resp:
			err := fn(response)
			if err != nil {
				logger.Error("%v error %v", name, err.Error())
			}
		}
	}
}
