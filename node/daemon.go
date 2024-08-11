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

type IAgent interface {
	ScanBlock() error
	ProofResponse(resp *common.ZkProofResponse) error
	Init() error
	CheckState() error
	Close() error
	Name() string
	FetchDataResponse(resp *FetchResponse) error
}

type IManager interface {
	IState
	Init() error
	PendingProofList() []*common.ZkProofRequest
	ReceiveRequest(requests []*common.ZkProofRequest) error
	CheckState() error
	GetProofRequest(proofTypes []common.ZkProofType) (*common.ZkProofRequest, bool, error)
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
	fetch             IFetch
	rpcServer         *rpc.Server
	wsServer          *rpc.Server
	nodeConfig        Config
	exitSignal        chan os.Signal
	manager           *WrapperManger
	txManager         *TxManager // todo
	enableLocalWorker bool
	debug             bool
}

func NewDaemon(cfg Config) (*Daemon, error) {
	err := logger.InitLogger(&logger.LogCfg{
		File:           true,
		IsStdout:       true,
		DiscordHookUrl: cfg.DiscordHookUrl,
	})
	if err != nil {
		return nil, err
	}
	logger.Info("current DebugMode :%v", cfg.Debug)

	// todo need to refactor
	logger.Info("beacon genesis Index: %v, slot:%v", cfg.GenesisSyncPeriod, cfg.BeaconInitSlot)
	btcClient, err := bitcoin.NewClient(cfg.BtcUrl, cfg.BtcUser, cfg.BtcPwd)
	if err != nil {
		logger.Error("new btc btcClient error:%v", err)
		return nil, err
	}
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
	ethClient, err := ethereum.NewClient(cfg.EthUrl, cfg.ZkBridgeAddr, cfg.ZkBtcAddr, cfg.UtxoManagerAddr)
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
	params.UseHoleskyNetworkConfig()
	params.OverrideBeaconConfig(params.HoleskyConfig())
	beaClient, err := apiclient.NewClient(cfg.BeaconUrl)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	// todo
	if cfg.Debug {
		count, err := btcClient.GetBlockCount()
		if err != nil {
			logger.Error("get btc block count error:%v", err)
			return nil, err
		}
		cfg.BtcGenesisHeight = (count/2016 - 1) * 2016
		cfg.BtcGenesisHeight = 2871700
		cfg.BeaconInitSlot = 273 * 8192
		cfg.GenesisSyncPeriod = 273
		//fileStore, err := NewFileStorage(cfg.Datadir, cfg.BeaconInitSlot, uint64(cfg.BtcGenesisHeight))
		//if err != nil {
		//	logger.Error("new fileStorage error: %v", err)
		//	return nil, err
		//}
		//logger.Debug("clear fileStore: %v", fileStore.RootPath)
		//err = fileStore.Clear()
		//if err != nil {
		//	logger.Error("clear fileStore error: %v", err)
		//	return nil, err
		//}
	}
	logger.Debug("beaconGenesisPeriod: %v,ethInitHeight: %v,btcGenesisHeight: %v,btcInitHeight: %v", cfg.GenesisSyncPeriod,
		cfg.EthInitHeight, cfg.BtcGenesisHeight, cfg.BtcInitHeight)

	dbPath := fmt.Sprintf("%s/%s", cfg.Datadir, cfg.Network)
	logger.Info("levelDbPath: %s", dbPath)
	storeDb, err := store.NewStore(dbPath, 0, 0, common.DbNameSpace, false)
	if err != nil {
		logger.Error("new store error:%v,dbPath:%s", err, dbPath)
		return nil, err
	}
	memoryStore := store.NewMemoryStore()
	proofRequest := make(chan []*common.ZkProofRequest, 1000)
	btcProofResp := make(chan *common.ZkProofResponse, 1000)
	ethProofResp := make(chan *common.ZkProofResponse, 1000)
	syncCommitResp := make(chan *common.ZkProofResponse, 1000)

	ethFetchDataResp := make(chan *FetchResponse, 1000)
	beaconFetchDataResp := make(chan *FetchResponse, 1000)
	btcFetchDataResp := make(chan *FetchResponse, 1000)

	fileStore, err := NewFileStorage(cfg.Datadir, cfg.BeaconInitSlot, uint64(cfg.BtcGenesisHeight))
	if err != nil {
		logger.Error("new fileStorage error: %v", err)
		return nil, err
	}

	keyStore, err := NewKeyStore(cfg.EthPrivateKey)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	logger.Debug("ethereum submit address: %v", keyStore.address)
	taskManager, err := NewTxManager(storeDb, keyStore, ethClient, btcClient, oasisClient)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	cache := NewCacheState()

	var agents []*WrapperAgent

	if false {
		beaconAgent, err := NewBeaconAgent(storeDb, beaconClient, beaClient, proofRequest, fileStore, cache, cfg.BeaconInitSlot, cfg.GenesisSyncPeriod)
		if err != nil {
			logger.Error("new node btcClient error:%v", err)
			return nil, err
		}
		agents = append(agents, NewWrapperAgent(beaconAgent, 15*time.Second, 17*time.Second, syncCommitResp, beaconFetchDataResp))
	}
	if true {
		btcAgent, err := NewBitcoinAgent(cfg, storeDb, memoryStore, fileStore, btcClient, ethClient, btcProverClient,
			proofRequest, keyStore, taskManager, cache)
		if err != nil {
			logger.Error(err.Error())
			return nil, err
		}
		agents = append(agents, NewWrapperAgent(btcAgent, cfg.BtcScanTime, 1*time.Minute, btcProofResp, btcFetchDataResp))

	}
	if false {
		ethAgent, err := NewEthereumAgent(cfg, cfg.BeaconInitSlot, fileStore, storeDb, memoryStore, beaClient, btcClient, ethClient,
			beaconClient, oasisClient, proofRequest, taskManager, cache)
		if err != nil {
			logger.Error(err.Error())
			return nil, err
		}
		agents = append(agents, NewWrapperAgent(ethAgent, cfg.EthScanTime, 1*time.Minute, ethProofResp, ethFetchDataResp))
	}
	workers := make([]rpc.IWorker, 0)
	if cfg.EnableLocalWorker {
		logger.Info("local worker enabled")
		zkParamDir := common.GetEnvZkParameterDir() // todo
		if !cfg.Debug {
			if zkParamDir == "" {
				logger.Error("zkParamDir is empty,please config  ZkParameterDir env")
				return nil, fmt.Errorf("zkParamDir is empty,please config  ZkParameterDir env")
			}
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
	scheduler := NewSchedule(workers)
	preparedData, err := NewPreparedData(fileStore, storeDb, cfg.GenesisSyncPeriod, uint64(cfg.BtcGenesisHeight), btcProverClient, btcClient, ethClient, beaClient, beaconClient)
	if err != nil {
		logger.Error("new PreparedData error: %v", err)
		return nil, err
	}

	msgManager, err := NewManager(btcClient, ethClient, preparedData, btcProofResp, ethProofResp, syncCommitResp,
		storeDb, memoryStore, scheduler, fileStore, uint64(cfg.BtcGenesisHeight), cfg.GenesisSyncPeriod, cfg.BeaconInitSlot, cache)
	if err != nil {
		logger.Error("new msgManager error: %v", err)
		return nil, err
	}
	exitSignal := make(chan os.Signal, 1)

	rpcHandler := NewHandler(msgManager, storeDb, memoryStore, scheduler, fileStore, exitSignal)
	server, err := rpc.NewServer(RpcRegisterName, fmt.Sprintf("%s:%s", cfg.Rpcbind, cfg.Rpcport), rpcHandler, keyStore, nil)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	logger.Debug("rpcServer: listen on %v,port  %v", cfg.Rpcbind, cfg.Rpcport)
	wsServer, err := rpc.NewWsServer(RpcRegisterName, fmt.Sprintf("%s:%s", cfg.Rpcbind, cfg.WsPort), rpcHandler)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	logger.Debug("wsServer: listen on %v,port  %v", cfg.Rpcbind, cfg.WsPort)
	//fetch, err := NewFetch(beaconClient, storeDb, fileStore, cfg.BeaconInitSlot, ethFetchDataResp)
	//if err != nil {
	//	logger.Error("new fetch error: %v", err)
	//	return nil, err
	//}
	daemon := &Daemon{
		agents:            agents,
		rpcServer:         server,
		wsServer:          wsServer,
		nodeConfig:        cfg,
		enableLocalWorker: cfg.EnableLocalWorker,
		exitSignal:        exitSignal,
		txManager:         taskManager,
		fetch:             nil, // todo
		debug:             common.GetEnvDebugMode(),
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
	if d.txManager != nil {
		err = d.txManager.init()
		if err != nil {
			logger.Error("txManager init error %v", err)
			return err
		}
	}
	if d.fetch != nil {
		err = d.fetch.Init()
		if err != nil {
			logger.Error("fetch init error %v", err)
			return err
		}
	}
	for _, agent := range d.agents {
		if err := agent.node.Init(); err != nil {
			logger.Error("%v:init agent error %v", agent.node.Name(), err)
			return err
		}
	}
	return nil
}

func (d *Daemon) Run() error {
	logger.Info("start daemon")
	// rpc rpcServer
	go d.rpcServer.Run()
	// ws server
	go d.wsServer.Run()

	// fetch
	if d.fetch != nil {
		go DoTimerTask("fetch-finality-update", 40*time.Second, d.fetch.FinalityUpdate, d.exitSignal)
		go DoTimerTask("fetch-update", 1*time.Minute, d.fetch.LightClientUpdate, d.exitSignal)
	}

	if !d.debug {
		go DoTimerTask("txManager-check", 30*time.Second, d.txManager.Check, d.exitSignal) // todo
	}

	// proof request manager
	go doProofRequestTask("manager-proofRequest", d.manager.proofRequest, d.manager.manager.ReceiveRequest, d.exitSignal)
	if d.enableLocalWorker {
		go DoTask("manager-generateProof:", d.manager.manager.DistributeRequest, d.exitSignal) // todo
	}
	go DoTimerTask("manager-checkState", d.manager.checkTime, d.manager.manager.CheckState, d.exitSignal)
	go DoTimerTask("manager-checkBtcState", d.manager.checkTime, d.manager.manager.CheckBtcState, d.exitSignal)
	//go DoTimerTask("manager-checkEthState", d.manager.checkTime, d.manager.manager.CheckEthState, d.exitSignal)
	//go DoTimerTask("manager-checkBeaconState", d.manager.checkTime, d.manager.manager.CheckBeaconState, d.exitSignal)

	for _, agent := range d.agents {
		proofReplyName := fmt.Sprintf("%s-proofResponse", agent.node.Name())
		go doProofResponseTask(proofReplyName, agent.proofResp, agent.node.ProofResponse, d.exitSignal)
		//scanName := fmt.Sprintf("%s-scanBlock", agent.node.Name())
		//go DoTimerTask(scanName, agent.scanTime, agent.node.ScanBlock, d.exitSignal)
		//fetchName := fmt.Sprintf("%s-fetchResponse", agent.node.Name())
		//go doFetchRespTask(fetchName, agent.fetchResp, agent.node.FetchDataResponse, d.exitSignal)
		//checkStateName := fmt.Sprintf("%s-checkState", agent.node.Name())
		//go DoTimerTask(checkStateName, agent.checkStateTime, agent.node.CheckState, d.exitSignal)

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
	err := d.rpcServer.Shutdown()
	if err != nil {
		logger.Error("rpc rpcServer shutdown error:%v", err)
	}
	err = d.wsServer.Shutdown()
	if err != nil {
		logger.Error("ws server shutdown error:%v", err)
	}
	for _, agent := range d.agents {
		if err := agent.node.Close(); err != nil {
			logger.Error("%v:close agent error %v", agent.node.Name(), err)
		}
	}
	err = d.manager.manager.Close()
	if err != nil {
		logger.Error("manager close error:%v", err)
	}
	if d.exitSignal != nil {
		close(d.exitSignal)
		//todo waitGroup
		time.Sleep(5 * time.Second)
	}
	err = logger.Close()
	if err != nil {
		fmt.Printf("logger close error: %v \n", err)
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
	fetchResp      chan *FetchResponse
	checkStateTime time.Duration
}

func NewWrapperAgent(agent IAgent, scanTime, checkState time.Duration, proofResp chan *common.ZkProofResponse, fetch chan *FetchResponse) *WrapperAgent {
	return &WrapperAgent{
		node:           agent,
		scanTime:       scanTime,
		proofResp:      proofResp,
		fetchResp:      fetch,
		checkStateTime: checkState,
	}
}
