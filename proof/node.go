package proof

import (
	"fmt"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	dnode "github.com/lightec-xyz/daemon/node"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/store"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Node Todo
type Node struct {
	rpcServer *rpc.Server
	mode      common.Mode
	local     *Local
	store     store.IStore
	exit      chan os.Signal
	Id        string
}

func NewNode(cfg Config) (*Node, error) {
	err := logger.InitLogger(nil)
	if err != nil {
		logger.Error("init logger error:%v", err)
		return nil, err
	}
	err = cfg.Check()
	if err != nil {
		logger.Error("config check error:%v", err)
		return nil, err
	}
	dbPath := fmt.Sprintf("%s/%s", cfg.DataDir, cfg.Network)
	logger.Info("dbPath:%s", dbPath)

	fileStorage, err := dnode.NewFileStorage(cfg.DataDir, 0)
	if err != nil {
		logger.Error("new fileStorage error:%v", err)
		return nil, err
	}

	storeDb, err := store.NewStore(dbPath, 0, 0, "zkbtc", false)
	if err != nil {
		logger.Error("new store error:%v,dbPath:%s", err, cfg.DataDir)
		return nil, err
	}
	workerId, exists, err := ReadWorkerId(storeDb)
	if err != nil {
		logger.Error("read worker id error:%v", err)
		return nil, err
	}
	if !exists {
		workerId = common.MustUUID()
		err := WriteWorkerId(storeDb, workerId)
		if err != nil {
			logger.Error("write worker id error:%v", err)
			return nil, err
		}
	}
	if cfg.Mode == common.Client {
		local, err := NewLocal(cfg.Url, cfg.DataDir, workerId, cfg.MaxNums, storeDb, fileStorage)
		if err != nil {
			logger.Error("new local error:%v", err)
			return nil, err
		}
		return &Node{
			local: local,
			mode:  cfg.Mode,
			exit:  make(chan os.Signal, 1),
			store: storeDb,
			Id:    workerId,
		}, nil
	} else if cfg.Mode == common.Cluster {
		host := fmt.Sprintf("%v:%v", cfg.RpcBind, cfg.RpcPort)
		memoryStore := store.NewMemoryStore()
		handler := NewHandler(storeDb, memoryStore, cfg.MaxNums)
		logger.Info("proof worker info: %v", cfg.Info())
		server, err := rpc.NewWsServer(RpcRegisterName, host, handler)
		if err != nil {
			logger.Error("new rpc rpcServer error:%v", err)
			return nil, err
		}
		return &Node{
			rpcServer: server,
			mode:      cfg.Mode,
			exit:      make(chan os.Signal, 1),
			store:     storeDb,
			Id:        workerId,
		}, nil
	}
	return nil, fmt.Errorf("new node error: unknown model:%v", cfg.Mode)

}

func (node *Node) Start() error {
	if node.mode == common.Client {
		go dnode.DoTimerTask("local-generator", 1*time.Minute, node.local.Run, node.exit)
		go dnode.DoTimerTask("local-checkState", 1*time.Minute, node.local.CheckState, node.exit)
	} else if node.mode == common.Cluster {
		go node.rpcServer.Run()
	}
	logger.Info("proof worker node start now ....")
	signal.Notify(node.exit, syscall.SIGTERM, syscall.SIGQUIT)
	for {
		msg := <-node.exit
		switch msg {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGTSTP:
			logger.Info("get shutdown signal ...")
			err := node.Close()
			if err != nil {
				logger.Error(" node close info error:%v", err)
			}
			return err
		}
	}
}

func (node *Node) Close() error {
	logger.Warn("proof worker node exit now ....")
	if node.rpcServer != nil {
		err := node.rpcServer.Shutdown()
		if err != nil {
			logger.Error(" proof worker node exit now: %v", err)
		}
	}
	if node.local != nil {
		err := node.local.Close()
		if err != nil {
			logger.Error(" proof worker node exit now: %v", err)
		}
	}
	return nil
}
