package node

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	btcproverClient "github.com/lightec-xyz/btc_provers/utils/client"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/rpc/beacon"
	"github.com/lightec-xyz/daemon/rpc/bitcoin"
	"github.com/lightec-xyz/daemon/rpc/ethereum"
	"github.com/lightec-xyz/daemon/rpc/oasis"
	"github.com/lightec-xyz/daemon/store"
	apiclient "github.com/lightec-xyz/provers/utils/api-client"

	"github.com/prysmaticlabs/prysm/v5/config/params"
)

func init() {
	err := logger.InitLogger()
	if err != nil {
		panic(err)
	}
}

type IAgent interface {
	ScanBlock() error
	ProofResponse(resp *common.ZkProofResponse) error
	Init() error
	CheckState() error
	Close() error
	Name() string
}

type IBeaconAgent interface {
	ProofResponse(resp *common.ZkProofResponse) error
	FetchDataResponse(resp FetchDataResponse) error
	CheckState() error
	Init() error
	Close() error
	Name() string
}

type IManager interface {
	Init() error
	ReceiveRequest(requests []*common.ZkProofRequest) error
	CheckPendingRequest() error
	GetProofRequest() (*common.ZkProofRequest, bool, error)
	SendProofResponse(response []*common.ZkProofResponse) error
	DistributeRequest() error
	Close() error
}

type IFetch interface {
	Init() error
	Bootstrap() error
	FinalityUpdate() error
	LightClientUpdate() error
	Close() error
}

type Daemon struct {
	agents            []*WrapperAgent
	beaconAgent       *WrapperBeacon
	fetch             IFetch
	rpcServer         *rpc.Server
	nodeConfig        Config
	exitSignal        chan os.Signal
	manager           *WrapperManger
	taskManager       *TaskManager // todo
	disableRecurAgent bool         // true ,Only enable the function of generating recursive proofs
	disableTxAgent    bool
	enableLocal       bool
}

func NewDaemon(cfg Config) (*Daemon, error) {
	var submitTxEthAddr string
	var err error
	if cfg.EthPrivateKey != "" {
		submitTxEthAddr, err = privateKeyToEthAddr(cfg.EthPrivateKey)
		if err != nil {
			logger.Error("privateKeyToEthAddr error:%v", err)
			return nil, err
		}
		logger.Info("ethereum submit address:%v", submitTxEthAddr)
	}
	logger.Info("beacon genesis Index: %v, slot:%v", cfg.GenesisSyncPeriod, cfg.BeaconInitSlot)

	btcClient, err := bitcoin.NewClient(cfg.BtcUrl, cfg.BtcUser, cfg.BtcPwd)
	if err != nil {
		logger.Error("new btc btcClient error:%v", err)
		return nil, err
	}
	// todo
	url := strings.Replace(cfg.BtcUrl, "http://", "", 1)
	btcProverClient, err := btcproverClient.NewClient(url, cfg.BtcUser, cfg.BtcPwd)
	if err != nil {
		logger.Error("new btc btcProverClient error:%v", err)
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
	oasisClient, err := oasis.NewClient(cfg.OasisUrl, cfg.Network, &oasis.Option{
		LocalAddress:   LocalOasisSignerAddr,
		TestnetAddress: TestnetOasisSignerAddr,
	})
	if err != nil {
		logger.Error("new eth btcClient error:%v", err)
		return nil, err
	}

	dbPath := fmt.Sprintf("%s/%s", cfg.Datadir, cfg.Network)
	logger.Info("levelDbPath: %s", dbPath)
	storeDb, err := store.NewStore(dbPath, 0, 0, "zkbtc", false)
	if err != nil {
		logger.Error("new store error:%v,dbPath:%s", err, dbPath)
		return nil, err
	}
	// todo
	memoryStore := store.NewMemoryStore()
	proofRequest := make(chan []*common.ZkProofRequest, 10)
	btcProofResp := make(chan *common.ZkProofResponse, 10)
	ethProofResp := make(chan *common.ZkProofResponse, 10)
	syncCommitResp := make(chan *common.ZkProofResponse, 10)
	fetchDataResp := make(chan FetchDataResponse, 10)

	fileStore, err := NewFileStorage(cfg.Datadir, cfg.BeaconInitSlot)
	if err != nil {
		logger.Error("new fileStorage error: %v", err)
		return nil, err
	}
	beaconAgent, err := NewBeaconAgent(cfg, beaconClient, proofRequest, fileStore, cfg.BeaconInitSlot, cfg.GenesisSyncPeriod, fetchDataResp)
	if err != nil {
		logger.Error("new node btcClient error:%v", err)
		return nil, err
	}
	keyStore := NewKeyStore(cfg.EthPrivateKey)

	taskManager, err := NewTaskManager(keyStore, ethClient, btcClient, oasisClient)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	var agents []*WrapperAgent
	btcAgent, err := NewBitcoinAgent(cfg, submitTxEthAddr, storeDb, memoryStore, fileStore, btcClient, ethClient, btcProverClient, proofRequest, keyStore, taskManager)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	agents = append(agents, NewWrapperAgent(btcAgent, cfg.BtcScanTime, 1*time.Minute, btcProofResp))

	// todo find a better way
	params.UseHoleskyNetworkConfig()
	params.OverrideBeaconConfig(params.HoleskyConfig())

	// tokenOpt := client.WithAuthenticationToken("3ac3d8d70361a628192b6fd7cd71b88a0b17638d")
	// beaClient, err := apiclient.NewClient("https://young-morning-meadow.ethereum-holesky.quiknode.pro", tokenOpt)

	beaClient, err := apiclient.NewClient(cfg.BeaconUrl)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	ethAgent, err := NewEthereumAgent(cfg, cfg.BeaconInitSlot, fileStore, storeDb, memoryStore, beaClient, btcClient, ethClient, beaconClient, oasisClient, proofRequest, taskManager)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	agents = append(agents, NewWrapperAgent(ethAgent, cfg.EthScanTime, 1*time.Minute, ethProofResp))

	workers := make([]rpc.IWorker, 0)
	if cfg.EnableLocalWorker {
		logger.Info("local worker enabled")
		zkParamDir := os.Getenv(common.ZkParameterDir) // todo
		if zkParamDir == "" {
			logger.Error("zkParamDir is empty,please config  ZkParameterDir env")
			return nil, fmt.Errorf("zkParamDir is empty,please config  ZkParameterDir env")
		}
		logger.Info("zkParamDir: %v", zkParamDir)
		localWorker, err := NewLocalWorker(zkParamDir, "", 1)
		if err != nil {
			logger.Error("new local worker error:%v", err)
			return nil, err
		}
		workers = append(workers, localWorker)
	} else {
		logger.Warn("no local worker to generate proof")
	}
	schedule := NewSchedule(workers)
	msgManager, err := NewManager(btcClient, ethClient, btcProofResp, ethProofResp, syncCommitResp,
		storeDb, memoryStore, schedule, fileStore)
	if err != nil {
		logger.Error("new msgManager error: %v", err)
		return nil, err
	}
	exitSignal := make(chan os.Signal, 1)

	rpcHandler := NewHandler(msgManager, storeDb, memoryStore, schedule, exitSignal)
	server, err := rpc.NewServer(RpcRegisterName, fmt.Sprintf("%s:%s", cfg.Rpcbind, cfg.Rpcport), rpcHandler)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	fetch, err := NewFetch(beaconClient, fileStore, cfg.BeaconInitSlot)
	if err != nil {
		logger.Error("new fetch error: %v", err)
		return nil, err
	}

	daemon := &Daemon{
		agents:            agents,
		rpcServer:         server,
		nodeConfig:        cfg,
		disableRecurAgent: cfg.DisableRecursiveAgent,
		disableTxAgent:    cfg.DisableTxAgent,
		enableLocal:       cfg.EnableLocalWorker,
		exitSignal:        exitSignal,
		taskManager:       taskManager,
		fetch:             fetch,
		beaconAgent:       NewWrapperBeacon(beaconAgent, 1*time.Minute, 1*time.Minute, syncCommitResp, fetchDataResp),
		manager:           NewWrapperManger(msgManager, proofRequest, 1*time.Minute),
	}
	return daemon, nil
}

func (d *Daemon) Init() error {
	err := d.manager.manager.Init()
	if err != nil {
		logger.Error("manager init error %v", err)
		return err
	}
	if !d.disableTxAgent {
		err := d.fetch.Init()
		if err != nil {
			logger.Error("fetch init error %v", err)
			return err
		}
		for _, agent := range d.agents {
			if err := agent.node.Init(); err != nil {
				logger.Error("%v:init agent error %v", agent.node.Name(), err)
				return err
			}
		}
	}
	if !d.disableRecurAgent {
		err = d.beaconAgent.node.Init()
		if err != nil {
			logger.Error("node agent init error %v", err)
			return err
		}
	}
	return nil
}

func (d *Daemon) Run() error {
	logger.Info("start daemon")
	// rpc rpcServer
	go d.rpcServer.Run()

	if !d.disableRecurAgent {
		// syncCommit proof
		go doProofResponseTask("beacon-proofResponse", d.beaconAgent.proofResponse, d.beaconAgent.node.ProofResponse, d.exitSignal)
		go doFetchRespTask("beacon-fetchDataResponse", d.beaconAgent.fetchDataResponse, d.beaconAgent.node.FetchDataResponse, d.exitSignal)
		go DoTimerTask("beacon-checkData", d.beaconAgent.checkDataTime, d.beaconAgent.node.CheckState, d.exitSignal)
		go DoTimerTask("fetch-finality-update", 1*time.Minute, d.fetch.FinalityUpdate, d.exitSignal)
		go DoTimerTask("fetch-update", 1*time.Minute, d.fetch.LightClientUpdate, d.exitSignal)
	}

	// task manager
	go DoTimerTask("task-manager", 1*time.Minute, d.taskManager.Check, d.exitSignal) // todo

	// proof request manager
	go doProofRequestTask("manager-proofRequest", d.manager.proofRequest, d.manager.manager.ReceiveRequest, d.exitSignal)
	if d.enableLocal {
		go DoTask("manager-generateProof:", d.manager.manager.DistributeRequest, d.exitSignal) // todo
	}
	go DoTimerTask("manager-checkPending", d.manager.checkTime, d.manager.manager.CheckPendingRequest, d.exitSignal)

	if !d.disableTxAgent {
		//tx Proof
		for _, agent := range d.agents {
			name := fmt.Sprintf("%s-submitProof", agent.node.Name())
			go doProofResponseTask(name, agent.proofResp, agent.node.ProofResponse, d.exitSignal)
		}
		//scan block with tx
		for _, agent := range d.agents {
			scanName := fmt.Sprintf("%s-scanBlock", agent.node.Name())
			go DoTimerTask(scanName, agent.scanTime, agent.node.ScanBlock, d.exitSignal)
			checkStateName := fmt.Sprintf("%s-checkState", agent.node.Name())
			go DoTimerTask(checkStateName, agent.checkStateTime, agent.node.CheckState, d.exitSignal)
		}

	}

	signal.Notify(d.exitSignal, syscall.SIGTERM, syscall.SIGQUIT)
	for {
		msg := <-d.exitSignal
		switch msg {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGTSTP:
			logger.Info("get shutdown signal ,waiting exit now ...")
			err := d.Close()
			if err != nil {
				logger.Error(err.Error())
			}
			return nil
		}
	}
}

func (d *Daemon) Close() error {
	if d.disableTxAgent {
		for _, agent := range d.agents {
			if err := agent.node.Close(); err != nil {
				logger.Error("%v:close agent error %v", agent.node.Name(), err)
			}
		}
	}
	err := d.rpcServer.Shutdown()
	if err != nil {
		logger.Error("rpc rpcServer shutdown error:%v", err)
	}
	d.manager.manager.Close()
	if d.disableRecurAgent {
		err = d.beaconAgent.node.Close()
		if err != nil {
			logger.Error("node agent close error:%v", err)
		}
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
	scanPeriodTime    time.Duration // get node Index
	checkDataTime     time.Duration
	proofResponse     chan *common.ZkProofResponse
	fetchDataResponse chan FetchDataResponse
}

func NewWrapperBeacon(beacon IBeaconAgent, scanPeriodTime, checkDataTime time.Duration, proofResponse chan *common.ZkProofResponse, fetchDataResp chan FetchDataResponse) *WrapperBeacon {
	return &WrapperBeacon{
		node:              beacon,
		scanPeriodTime:    scanPeriodTime,
		proofResponse:     proofResponse,
		fetchDataResponse: fetchDataResp,
		checkDataTime:     checkDataTime, // todo
	}
}

type WrapperManger struct {
	manager      IManager
	proofRequest chan []*common.ZkProofRequest
	checkTime    time.Duration
}

func NewWrapperManger(manager IManager, request chan []*common.ZkProofRequest, checkTime time.Duration) *WrapperManger {
	return &WrapperManger{
		manager:      manager,
		proofRequest: request,
		checkTime:    checkTime,
	}
}

type WrapperAgent struct {
	node           IAgent
	scanTime       time.Duration
	proofResp      chan *common.ZkProofResponse
	checkStateTime time.Duration
}

func NewWrapperAgent(agent IAgent, scanTime, checkState time.Duration, proofResp chan *common.ZkProofResponse) *WrapperAgent {
	return &WrapperAgent{
		node:           agent,
		scanTime:       scanTime,
		proofResp:      proofResp,
		checkStateTime: checkState,
	}
}
