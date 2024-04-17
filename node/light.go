package node

import (
	"fmt"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/rpc/beacon"
	"github.com/lightec-xyz/daemon/rpc/bitcoin"
	"github.com/lightec-xyz/daemon/rpc/ethereum"
	"github.com/lightec-xyz/daemon/store"
	"os"
	"time"
)

// only generate recursive proof

func NewRecursiveLightDaemon(cfg NodeConfig) (*Daemon, error) {
	beaconClient, err := beacon.NewClient(cfg.BeaconUrl)
	if err != nil {
		logger.Error("new node btcClient error:%v", err)
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
	proofRequest := make(chan []*common.ZkProofRequest)
	btcProofResp := make(chan *common.ZkProofResponse)
	ethProofResp := make(chan *common.ZkProofResponse)
	syncCommitResp := make(chan *common.ZkProofResponse)
	fetchDataResp := make(chan FetchDataResponse)

	genesisPeriod := uint64(cfg.BeaconSlotHeight) / common.SlotPerPeriod
	fileStore, err := NewFileStore(cfg.DataDir, cfg.BeaconSlotHeight)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	beaconAgent, err := NewBeaconAgent(cfg, beaconClient, proofRequest, fileStore, cfg.BeaconSlotHeight, genesisPeriod, fetchDataResp)
	if err != nil {
		logger.Error("new node btcClient error:%v", err)
		return nil, err
	}
	workers := make([]rpc.IWorker, 0)
	if cfg.EnableLocalWorker {
		logger.Info("local worker enable")
		// todo
		zkParamDir := os.Getenv(common.ZkParameterDir)
		if zkParamDir == "" {
			zkParamDir = fmt.Sprintf("%s/setup", cfg.DataDir)
		}
		localWorkerDir := fmt.Sprintf("%s/localWork", cfg.DataDir)
		localWorker, err := NewLocalWorker(zkParamDir, localWorkerDir, 1)
		if err != nil {
			logger.Error("new local worker error:%v", err)
			return nil, err
		}
		workers = append(workers, localWorker)
	}
	schedule := NewSchedule(workers)
	manager, err := NewManager(nil, nil, btcProofResp, ethProofResp, syncCommitResp, storeDb, memoryStore, schedule)
	if err != nil {
		logger.Error("new manager error: %v", err)
		return nil, err
	}
	exitSignal := make(chan os.Signal, 1)

	rpcHandler := NewHandler(manager, storeDb, memoryStore, schedule, exitSignal)
	server, err := rpc.NewServer(RpcRegisterName, fmt.Sprintf("%s:%s", cfg.Rpcbind, cfg.RpcPort), rpcHandler)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	keyStore := NewKeyStore(cfg.EthPrivateKey)
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
	taskManager, err := NewTaskManager(keyStore, ethClient, btcClient)
	if err != nil {
		logger.Error("new taskManager error:%v", err)
		return nil, err
	}
	daemon := &Daemon{
		nodeConfig:    cfg,
		server:        server,
		exitSignal:    exitSignal,
		enableSyncCom: true,
		enableTx:      false,
		taskManager:   taskManager,
		beaconAgent:   NewWrapperBeacon(beaconAgent, 1*time.Minute, 1*time.Minute, 1*time.Minute, syncCommitResp, fetchDataResp),
		manager:       NewWrapperManger(manager, proofRequest, 1*time.Minute),
	}
	return daemon, nil
}
