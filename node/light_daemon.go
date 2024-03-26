package node

import (
	"fmt"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/rpc/beacon"
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
	daemon := &Daemon{
		nodeConfig:  cfg,
		exitSignal:  make(chan os.Signal, 1),
		beaconAgent: NewWrapperBeacon(beaconAgent, 1*time.Minute, 1*time.Minute, syncCommitResp, fetchDataResp),
		manager:     NewWrapperManger(manager, proofRequest),
	}
	return daemon, nil
}
