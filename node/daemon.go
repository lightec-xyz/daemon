package node

import (
	"fmt"
	"github.com/aviate-labs/agent-go/identity"
	btcproverClient "github.com/lightec-xyz/btc_provers/utils/client"
	"github.com/lightec-xyz/daemon/node/p2p"
	"github.com/lightec-xyz/daemon/rpc/sgx"
	"github.com/ybbus/jsonrpc/v3"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lightec-xyz/daemon/rpc/dfinity"

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
	ReScan(height uint64) error
	ProofResponse(resp *common.ProofResponse) error
	Init() error
	CheckState() error
	Close() error
	Name() string
}

type IManager interface {
	IScheduler
	Init() error
	PendingProofRequest() []*common.ProofRequest
	CheckState() error
	GetProofRequest(proofTypes []common.ProofType) (*common.ProofRequest, bool, error)
	ReceiveProofs(submitProof *common.SubmitProof) error
	LibP2pMessage(msg *p2p.Msg) error
	MinerPower() error
	AddP2pPeer(addr string) error
	ChainFork(signal *ChainFork) error
	EthNotify() chan *Notify
	BtcNotify() chan *Notify
	BeaconNotify() chan *Notify
	Close() error
}

type IFetch interface {
	Init() error
	Bootstrap()
	FinalityUpdate() error
	LightClientUpdate() error
	Close() error
}

type Daemon struct {
	agents     []*WrapperAgent
	fetch      IFetch
	rpcServer  *rpc.Server
	wsServer   *rpc.Server
	cfg        Config
	exitSignal chan os.Signal
	manager    *WrapperManger
	txManager  *TxManager
	libp2p     *p2p.LibP2p
	debug      bool
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
	btcClient, err := bitcoin.NewClient(cfg.BtcUrl, cfg.BtcUser, cfg.BtcPwd)
	if err != nil {
		logger.Error("new btc btcClient error:%v", err)
		return nil, err
	}
	dbPath := fmt.Sprintf("%s/%s", cfg.Datadir, cfg.Network)
	logger.Info("levelDbPath: %s", dbPath)
	storeDb, err := store.NewStore(dbPath, 0, 0, common.DbNameSpace, false)
	if err != nil {
		logger.Error("new store error:%v,dbPath:%s", err, dbPath)
		return nil, err
	}
	proverClient := btcproverClient.NewJsonRpcClient(cfg.BtcUrl, cfg.BtcUser, cfg.BtcPwd, &jsonrpc.RPCClientOpts{
		HTTPClient: &http.Client{
			Timeout: 1 * time.Minute,
		},
	})
	btcProverClient := NewBtcClient(proverClient, storeDb, btcClient, int64(cfg.BtcInitHeight))

	beaconClient, err := beacon.NewClient(cfg.BeaconUrl)
	if err != nil {
		logger.Error("new beacon client error:%v", err)
		return nil, err
	}
	ethClient, err := ethereum.NewClient(cfg.EthUrl, cfg.ZkBridgeAddr, cfg.UtxoManagerAddr, cfg.BtcTxVerifyAddr, cfg.ZkBtcAddr)
	if err != nil {
		logger.Error("new eth client error:%v", err)
		return nil, err
	}
	sgxClient := sgx.NewClient(cfg.SgxUrl)
	oasisClient, err := oasis.NewClient(cfg.OasisUrl, cfg.OasisSignerAddress)
	if err != nil {
		logger.Error("new oasis client error:%v", err)
		return nil, err
	}
	params.UseHoleskyNetworkConfig()
	params.OverrideBeaconConfig(params.HoleskyConfig())
	beaClient, err := apiclient.NewClient(cfg.BeaconUrl)
	if err != nil {
		logger.Error("new provers api client error: %v", err)
		return nil, err
	}
	var dfinityClient *dfinity.Client
	if cfg.IcpPrivateKey != "" {
		icpIdentity, err := identity.NewSecp256k1IdentityFromPEMWithoutParameters([]byte(cfg.IcpPrivateKey))
		if err != nil {
			logger.Error("new icp identity error: %v", err)
			return nil, err
		}
		dfinityClient, err = dfinity.NewClientWithIdentity(cfg.IcpSingerAddress, cfg.IcpWalletAddress,
			cfg.IcpSingerUrl, icpIdentity)
		if err != nil {
			logger.Error("new dfinity client error: %v", err)
			return nil, err
		}
	} else {
		dfinityClient, err = dfinity.NewClient(cfg.IcpSingerAddress)
		if err != nil {
			logger.Error("new dfinity client error: %v", err)
			return nil, err
		}
	}

	logger.Debug("beaconGenesisSlot: %v, beaconGenesisPeriod: %v,ethInitHeight: %v,btcGenesisHeight: %v, btcInitHeight: %v",
		cfg.GenesisBeaconSlot, cfg.GenesisSyncPeriod, cfg.EthInitHeight, cfg.BtcGenesisHeight, cfg.BtcInitHeight)

	btcProofResp := make(chan *common.ProofResponse, 64)
	ethProofResp := make(chan *common.ProofResponse, 64)
	syncCommitResp := make(chan *common.ProofResponse, 64)

	btcNotify := make(chan *Notify, 32)
	ethNotify := make(chan *Notify, 32)
	beaconNotify := make(chan *Notify, 32)

	ethReScan := make(chan *ReScnSignal, 16)
	btcReScan := make(chan *ReScnSignal, 16)
	chainFork := make(chan *ChainFork, 16)

	fileStore, err := NewFileStorage(cfg.Datadir, cfg.GenesisBeaconSlot, cfg.BtcGenesisHeight)
	if err != nil {
		logger.Error("new fileStorage error: %v", err)
		return nil, err
	}
	keyStore, err := NewKeyStore(cfg.EthPrivateKey)
	if err != nil {
		logger.Error("new keyStore error: %v", err)
		return nil, err
	}
	logger.Debug("ethereum  address: %v", keyStore.address)
	if cfg.MinerAddr == "" {
		logger.Warn("miner address is empty, %v as miner address", keyStore.address)
		cfg.MinerAddr = keyStore.address
	}
	preparedData, err := NewPreparedData(fileStore, storeDb, cfg.GenesisBeaconSlot, cfg.BtcGenesisHeight,
		btcProverClient, btcClient, ethClient, beaClient, beaconClient, cfg.MinerAddr)
	if err != nil {
		logger.Error("new proof Prepared data error: %v", err)
		return nil, err
	}
	taskManager, err := NewTxManager(storeDb, fileStore, preparedData, keyStore, ethClient, btcClient, oasisClient, dfinityClient,
		sgxClient, btcProverClient, cfg.MinerAddr)
	if err != nil {
		logger.Error("new tx manager error: %v", err)
		return nil, err
	}
	var agents []*WrapperAgent
	if !cfg.DisableBeaconAgent {
		beaconAgent, err := NewBeaconAgent(cfg.BeaconReScan, storeDb, beaconClient, beaClient, fileStore, cfg.BeaconInitSlot)
		if err != nil {
			logger.Error("new beacon agent error:%v", err)
			return nil, err
		}
		agents = append(agents, NewWrapperAgent(beaconAgent, 15*time.Second, nil, syncCommitResp))
	}
	if !cfg.DisableBtcAgent {
		btcAgent, err := NewBitcoinAgent(cfg, storeDb, btcClient, ethClient, dfinityClient, taskManager, chainFork, fileStore)
		if err != nil {
			logger.Error("new bitcoin agent error:%v", err)
			return nil, err
		}
		agents = append(agents, NewWrapperAgent(btcAgent, cfg.BtcScanTime, btcReScan, btcProofResp))

	}
	if !cfg.DisableEthAgent {
		ethAgent, err := NewEthereumAgent(cfg, fileStore, storeDb, btcClient, ethClient, taskManager, chainFork)
		if err != nil {
			logger.Error("new ethereum agent error:%v", err)
			return nil, err
		}
		agents = append(agents, NewWrapperAgent(ethAgent, cfg.EthScanTime, ethReScan, ethProofResp))
	}
	var libp2p *p2p.LibP2p
	if !cfg.DisableLipP2p {
		libp2p, err = p2p.NewLibP2p(p2p.NewP2pConfig(cfg.EthPrivateKey, cfg.P2pPort, cfg.P2pBootstraps))
		if err != nil {
			logger.Error("new libp2p error: %v", err)
			return nil, err
		}
	}
	workers := make([]rpc.IWorker, 0)
	if cfg.EnableLocalWorker {
		logger.Info("local worker enabled")
		if cfg.BtcSetupDir == "" || cfg.EthSetupDir == "" {
			return nil, fmt.Errorf("btcSetupDir or ethSetupDir is empty, please config it")
		}
		localWorkerDir := fmt.Sprintf("%s/localWorker", cfg.Datadir)
		localWorker, err := NewLocalWorker(cfg.BtcSetupDir, cfg.EthSetupDir, localWorkerDir, UUID(), 1, 0)
		if err != nil {
			logger.Error("new local worker error:%v", err)
			return nil, err
		}
		workers = append(workers, localWorker)
	} else {
		logger.Warn("no local worker to generate proof")
	}

	manager, err := NewManager(cfg.MinerAddr, libp2p, dfinityClient, btcClient, ethClient, preparedData, btcProofResp, ethProofResp, syncCommitResp,
		storeDb, fileStore, btcNotify, ethNotify, beaconNotify)
	if err != nil {
		logger.Error("new manager error: %v", err)
		return nil, err
	}
	exitSignal := make(chan os.Signal, 1)

	rpcHandler := NewHandler(manager, ethReScan, btcReScan, storeDb, fileStore, exitSignal)
	server, err := rpc.NewServer(RpcRegisterName, fmt.Sprintf("%s:%s", cfg.Rpcbind, cfg.Rpcport), rpcHandler, keyStore, nil)
	if err != nil {
		logger.Error("new rpc server error: %v", err)
		return nil, err
	}
	logger.Debug("rpcServer: listen on %v,port  %v", cfg.Rpcbind, cfg.Rpcport)
	var fetch IFetch
	if !cfg.DisableFetch {
		fetch, err = NewFetch(beaconClient, storeDb, fileStore, cfg.GenesisBeaconSlot, beaconNotify, ethNotify)
		if err != nil {
			logger.Error("new fetch error: %v", err)
			return nil, err
		}
	}
	daemon := &Daemon{
		agents:     agents,
		rpcServer:  server,
		wsServer:   nil,
		cfg:        cfg,
		exitSignal: exitSignal,
		txManager:  taskManager,
		libp2p:     libp2p,
		fetch:      fetch,
		debug:      common.GetEnvDebugMode(),
		manager:    NewWrapperManger(manager, chainFork),
	}
	return daemon, nil
}

func (d *Daemon) Init() error {
	logger.Debug("init zkbtc daemon ...")
	err := d.manager.Init()
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
	if !d.cfg.DisableFetch {
		err = d.fetch.Init()
		if err != nil {
			logger.Error("fetch init error %v", err)
			return err
		}
	}
	for _, agent := range d.agents {
		if err := agent.Init(); err != nil {
			logger.Error("%v:init agent error %v", agent.Name(), err)
			return err
		}
	}
	return nil
}

func (d *Daemon) Run() error {
	logger.Info("start zkbtc daemon ...")
	if d.rpcServer != nil {
		go d.rpcServer.Run()
	}
	if d.wsServer != nil {
		go d.wsServer.Run()
	}
	if d.cfg.DiscordHookUrl != "" {
		go DoTimerTask("heartbeat", 1*time.Hour, d.heartBeat, d.exitSignal)
	}
	// fetch
	if !d.cfg.DisableFetch {
		go DoTimerTask("fetch-finality-update", 30*time.Second, d.fetch.FinalityUpdate, d.exitSignal)
		go DoTimerTask("fetch-update", 2*time.Minute, d.fetch.LightClientUpdate, d.exitSignal)

	}
	if !d.cfg.DisableLipP2p {
		d.libp2p.Run()
		d.libp2p.SayHello(d.cfg.MinerAddr)
		go doLibP2pMsgTask("manager-libp2p", d.libp2p.MsgChan(), d.manager.LibP2pMessage, d.exitSignal)
	}
	if !d.debug {
		go DoTimerTask("txManager-check", 10*time.Minute, d.txManager.Check, d.exitSignal)
	}
	go doChainForkTask("manager-chainFork", d.manager.chainFork, d.manager.ChainFork, d.exitSignal)
	go DoTimerTask("manager-minerPower", 1*time.Minute, d.manager.MinerPower, d.exitSignal)
	if d.cfg.EnableLocalWorker {
		//go DoTask("manager-generateProof:", d.manager.manager.DistributeRequest, d.exitSignal)
	}
	if !d.cfg.DisableBtcAgent {
		go DoTimerTask("manager-checkBtcState", 2*time.Minute, d.manager.CheckBtcState, d.exitSignal, d.manager.BtcNotify())
		go DoTimerTask("manager-checkPreBtcState", 3*time.Minute, d.manager.CheckPreBtcState, d.exitSignal)
		go DoTimerTask("manager-icpSignature", 4*time.Minute, d.manager.BlockSignature, d.exitSignal)
		go DoTimerTask("manager-updateBtcCp", 24*time.Hour, d.manager.UpdateBtcCp, d.exitSignal)
	}
	if !d.cfg.DisableEthAgent {
		go DoTimerTask("manager-checkEthState", 1*time.Minute, d.manager.CheckEthState, d.exitSignal, d.manager.EthNotify())
	}
	if !d.cfg.DisableBeaconAgent {
		go DoTimerTask("manager-checkBeaconState", 1*time.Minute, d.manager.CheckBeaconState, d.exitSignal, d.manager.BeaconNotify())
	}
	go DoTimerTask("manager-checkState", 2*time.Minute, d.manager.CheckState, d.exitSignal)
	for _, agent := range d.agents {
		proofReplyName := fmt.Sprintf("%s-proofResponse", agent.Name())
		go doProofResponseTask(proofReplyName, agent.proofResp, agent.ProofResponse, d.exitSignal)
		reScanName := fmt.Sprintf("%s-reScan", agent.Name())
		go doReScanTask(reScanName, agent.reScan, agent.ReScan, d.exitSignal)
		scanName := fmt.Sprintf("%s-scanBlock", agent.Name())
		go DoTimerTask(scanName, agent.scanTime, agent.ScanBlock, d.exitSignal)
		checkStateName := fmt.Sprintf("%s-checkState", agent.Name())
		go DoTimerTask(checkStateName, 30*time.Minute, agent.CheckState, d.exitSignal)

	}
	signal.Notify(d.exitSignal, syscall.SIGTERM, syscall.SIGQUIT)
	for {
		msg := <-d.exitSignal
		switch msg {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGTSTP:
			logger.Info("get shutdown signal ,waiting exit now ...")
			d.Close()
			time.Sleep(2 * time.Second)
			return nil
		}
	}
}

func (d *Daemon) Close() {
	logger.Debug("daemon start close now ...")
	if d.exitSignal != nil {
		close(d.exitSignal)
	}
	if d.manager != nil {
		err := d.manager.Close()
		if err != nil {
			logger.Warn("manager close error:%v", err)
		}
	}
	for _, agent := range d.agents {
		if err := agent.Close(); err != nil {
			logger.Warn("%v agent close error %v", agent.Name(), err)
		}
	}
	if d.fetch != nil {
		err := d.fetch.Close()
		if err != nil {
			logger.Warn("fetch close error:%v", err)
		}
	}
	if d.txManager != nil {
		err := d.txManager.Close()
		if err != nil {
			logger.Warn("txManager close error:%v", err)
		}
	}
	if d.libp2p != nil {
		d.libp2p.Close()
	}
	if d.rpcServer != nil {
		err := d.rpcServer.Shutdown()
		if err != nil {
			logger.Warn("rpc rpcServer shutdown error:%v", err)
		}
	}
	if d.wsServer != nil {
		err := d.wsServer.Shutdown()
		if err != nil {
			logger.Warn("ws rpcServer shutdown error:%v", err)
		}
	}
	logger.Debug("all goroutine closed now ...")

	err := logger.Flush()
	if err != nil {
		fmt.Printf("logger close error: %v \n", err)
	}
}

func (d *Daemon) heartBeat() error {
	logger.Error("heartBeat,i am zkbtc daemon living ...")
	return nil
}

type WrapperManger struct {
	IManager
	chainFork chan *ChainFork
}

func NewWrapperManger(manager IManager, request chan *ChainFork) *WrapperManger {
	return &WrapperManger{
		IManager:  manager,
		chainFork: request,
	}
}

type WrapperAgent struct {
	IAgent
	scanTime  time.Duration
	proofResp chan *common.ProofResponse
	reScan    chan *ReScnSignal
}

func NewWrapperAgent(agent IAgent, scanTime time.Duration, reScan chan *ReScnSignal, proofResp chan *common.ProofResponse) *WrapperAgent {
	return &WrapperAgent{
		IAgent:    agent,
		scanTime:  scanTime,
		proofResp: proofResp,
		reScan:    reScan,
	}
}

type WorkerConfig struct {
	MaxNums int    `json:"maxNums"`
	Url     string `json:"url"`
}
