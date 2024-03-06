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
	Submit(resp ZkProofResponse) error
	Init() error
	Close() error
	Name() string
}

type IBeaconAgent interface {
	ScanSyncPeriod() error
	ProofResponse(resp ZkProofResponse) error
	FetchResponse(resp FetchDataResponse) error
	CheckData() error
	Init() error
	Close() error
	Name() string
}

type Daemon struct {
	agents      []*WrapperAgent
	beaconAgent *WrapperBeacon
	server      *rpc.Server
	nodeConfig  NodeConfig
	exitSignal  chan os.Signal
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

	beaconAgent, err := NewBeaconAgent(cfg, beaconClient, fetchDataResp)
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

	workers := make([]rpc.IProof, 0)
	if cfg.EnableLocalWorker {
		logger.Info("local worker enable")
		workers = append(workers, NewLocalWorker(1))
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
		beaconAgent: NewWrapperBeacon(beaconAgent, 1*time.Hour, syncCommitResp, fetchDataResp),
		manager:     NewWrapperManger(manager, proofRequest),
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
	err = d.beaconAgent.node.Init()
	if err != nil {
		logger.Error("node agent init error %v", err)
		return err
	}
	return nil
}

func (d *Daemon) Run() error {
	// syncCommit
	//go doTimerTask("beacon-ScanSyncPeriod", d.beaconAgent.time, d.beaconAgent.node.ScanSyncPeriod, d.exitSignal)
	//go doProofResponseTask("beacon-DepositResponse", d.beaconAgent.proofResponse, d.beaconAgent.node.ProofResponse, d.exitSignal)
	//go doFetchRespTask("beacon-DepositResponse", d.beaconAgent.fetchDataResponse, d.beaconAgent.node.FetchResponse, d.exitSignal)
	//go doTask("beacon-CheckData", d.beaconAgent.node.CheckData, d.exitSignal)

	// proof manager
	go doProofRequestTask("manager-DepositRequest", d.manager.proofRequest, d.manager.manager.run, d.exitSignal)
	go doTask("manager-GenerateProof:", d.manager.manager.genProof, d.exitSignal)

	// tx proof
	for _, agent := range d.agents {
		name := fmt.Sprintf("%s-SubmitProof", agent.node.Name())
		go doProofResponseTask(name, agent.proofResp, agent.node.Submit, d.exitSignal)
	}
	// scan block with tx
	for _, agent := range d.agents {
		name := fmt.Sprintf("%s-ScanBlock", agent.node.Name())
		go doTimerTask(name, agent.scanTime, agent.node.ScanBlock, d.exitSignal)
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

func NewWorkers(workers []WorkerConfig) ([]rpc.IProof, error) {
	workersList := make([]rpc.IProof, 0)
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

type WrapperBeacon struct {
	node              IBeaconAgent
	time              time.Duration // get node period
	proofResponse     chan ZkProofResponse
	fetchDataResponse chan FetchDataResponse
}

func NewWrapperBeacon(beacon IBeaconAgent, time time.Duration, proofResponse chan ZkProofResponse, fetchDataResp chan FetchDataResponse) *WrapperBeacon {
	return &WrapperBeacon{
		node:              beacon,
		time:              time,
		proofResponse:     proofResponse,
		fetchDataResponse: fetchDataResp,
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
