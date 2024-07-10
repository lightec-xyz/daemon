package proof

import (
	"fmt"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	dnode "github.com/lightec-xyz/daemon/node"
	"github.com/lightec-xyz/daemon/store"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Node Todo
type Node struct {
	server *Server
	local  *Local
	mode   common.Mode
	exit   chan os.Signal
}

func NewNode(cfg Config) (*Node, error) {
	err := logger.InitLogger(&logger.LogCfg{
		File:           true,
		IsStdout:       true,
		DiscordHookUrl: cfg.DiscordHookUrl,
	})
	if err != nil {
		return nil, err
	}
	err = cfg.Check()
	if err != nil {
		logger.Error("config check error:%v", err)
		return nil, err
	}
	dbPath := fmt.Sprintf("%s/%s", cfg.DataDir, cfg.Network)
	logger.Info("dbPath:%s", dbPath)
	logger.Info("mode:%v", cfg.Mode)
	zkProofTypes, err := toZkProofType(cfg.ProofType)
	if err != nil {
		logger.Error("convert proof type error:%v", err)
		return nil, err
	}
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
	debugMode := common.GetEnvDebugMode()
	logger.Debug("DebugMode: %v", debugMode)
	zkParameterDir := common.GetEnvZkParameterDir()
	if !debugMode {
		if zkParameterDir == "" {
			return nil, fmt.Errorf("zkParameterDir is empty,please config ZkParameterDir env")
		}
	}
	workerId, exists, err := dnode.ReadWorkerId(storeDb)
	if err != nil {
		logger.Error("read worker id error:%v", err)
		return nil, err
	}
	if !exists {
		workerId = common.MustUUID()
		err := dnode.WriteWorkerId(storeDb, workerId)
		if err != nil {
			logger.Error("write worker id error:%v", err)
			return nil, err
		}
	}
	verified, err := dnode.ReadZkParamVerify(storeDb)
	if err != nil {
		logger.Error("read zkParamVerify error:%v", err)
		return nil, err
	}
	if !verified && !debugMode {
		logger.Debug("**** start check zk parameters md5 ****")
		parameters, err := ReadParameters([]byte(common.ParametersStr)) // todo get from http server
		if err != nil {
			logger.Error("read parameters error:%v", err)
			return nil, err
		}
		for index, item := range parameters {
			path := zkParameterDir + "/" + item.FileName
			logger.Debug("start verify zk file:%v %v ", index, path)
			fileBytes, err := os.ReadFile(path)
			if err != nil {
				logger.Error("read zk file error: %v %v", path, err)
				return nil, fmt.Errorf("read zk file error: %v %v", path, err)
			}
			fileMd5 := common.HexMd5(fileBytes)
			if fileMd5 != item.Hash {
				logger.Error("check zk md5 not match: %v %v %v", path, fileMd5, item.Hash)
				return nil, fmt.Errorf("check zk md5 not match: %v %v %v", path, fileMd5, item.Hash)
			}
		}
		logger.Debug("**** check zk parameters md5 end ****")
		err = dnode.WriteZkParamVerify(storeDb, true)
		if err != nil {
			logger.Error("write zkParamVerify error:%v", err)
			return nil, err
		}
	}
	worker, err := dnode.NewLocalWorker(zkParameterDir, cfg.DataDir, cfg.MaxNums)
	if err != nil {
		logger.Error("new local worker error:%v", err)
		return nil, err
	}
	memoryStore := store.NewMemoryStore()
	if cfg.Mode == common.Client {
		local, err := NewLocal(cfg.Url, worker, zkProofTypes, storeDb, fileStorage)
		if err != nil {
			logger.Error("new local error:%v", err)
			return nil, err
		}
		return &Node{
			local: local,
			mode:  cfg.Mode,
			exit:  make(chan os.Signal, 1),
		}, nil
	} else if cfg.Mode == common.Custom {
		handler := NewHandler(storeDb, memoryStore, worker)
		url := fmt.Sprintf("%v?id=%v", cfg.Url, workerId)
		server, err := NewServer(url, cfg.Mode, handler)
		if err != nil {
			logger.Error("new server error:%v", err)
			return nil, err
		}
		return &Node{
			server: server,
			mode:   cfg.Mode,
			exit:   make(chan os.Signal, 1),
		}, nil
	} else if cfg.Mode == common.Cluster {
		handler := NewHandler(storeDb, memoryStore, worker)
		server, err := NewServer(fmt.Sprintf("%v:%v", cfg.RpcBind, cfg.RpcPort), cfg.Mode, handler)
		if err != nil {
			logger.Error("new server error:%v", err)
			return nil, err
		}
		return &Node{
			server: server,
			mode:   cfg.Mode,
			exit:   make(chan os.Signal, 1),
		}, nil
	}
	return nil, fmt.Errorf("new node error: unknown model:%v", cfg.Mode)

}

func (node *Node) Init() error {
	if node.local != nil {
		err := node.local.Init()
		if err != nil {
			logger.Error("local init error:%v", err)
			return err
		}
	}
	if node.server != nil {
		err := node.server.init()
		if err != nil {
			logger.Error("server init error:%v", err)
			return err
		}
	}
	return nil
}

func (node *Node) Start() error {
	if node.mode == common.Client {
		go dnode.DoTimerTask("generator", 30*time.Second, node.local.polling, node.exit)
		go dnode.DoTimerTask("checkState", 30*time.Second, node.local.CheckState, node.exit)
	} else if node.mode == common.Cluster || node.mode == common.Custom {
		node.server.Run()
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
	if node.server != nil {
		err := node.server.Close()
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
	err := logger.Close()
	if err != nil {
		fmt.Printf("logger close error: %v", err)
	}
	return nil
}
