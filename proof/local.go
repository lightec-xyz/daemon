package proof

import (
	"fmt"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/node"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/store"
)

type Local struct {
	Id          string
	client      *rpc.NodeClient
	worker      rpc.IWorker
	store       store.IStore
	fileStore   *node.FileStorage
	exit        chan struct{}
	cacheProofs *node.ProofRespQueue
}

func NewLocal(url, datadir, id string, num int, store store.IStore, fileStore *node.FileStorage) (*Local, error) {
	client, err := rpc.NewNodeClient(url)
	if err != nil {
		logger.Error("new node client error:%v", err)
		return nil, err
	}
	zkParamDir := common.GetEnvZkParameterDir()
	debugMode := common.GetEnvDebugMode()
	logger.Debug("DebugMode: %v", debugMode)

	if !debugMode {
		if zkParamDir == "" {
			logger.Error("zkParamDir is empty,please config  ZkParameterDir env")
			return nil, fmt.Errorf("zkParamDir is empty,please config  ZkParameterDir env")
		}
	}
	logger.Info("zkParamDir: %v", zkParamDir)
	worker, err := node.NewLocalWorker(zkParamDir, datadir, num)
	if err != nil {
		logger.Error("new local worker error:%v", err)
		return nil, err
	}
	logger.Info("workerId: %v", id)
	return &Local{
		client:      client,
		fileStore:   fileStore,
		worker:      worker,
		Id:          id,
		store:       store,
		exit:        make(chan struct{}, 1),
		cacheProofs: node.NewProofRespQueue(),
	}, nil
}

func (l *Local) Run() error {
	logger.Debug("generator proof run")
	if l.worker.CurrentNums() >= l.worker.MaxNums() {
		logger.Warn("maxNums limit reached, wait proof generated")
		return nil
	}
	request := common.TaskRequest{
		Id:        l.Id,
		ProofType: []common.ZkProofType{}, // Todo worker support which proof type
	}
	requestResp, err := l.client.GetZkProofTask(request)
	if err != nil {
		logger.Error("get task error:%v", err)
		return nil
	}
	if !requestResp.CanGen {
		logger.Debug("no new proof request, wait  request coming now ....")
		return nil
	}
	l.worker.AddReqNum()
	err = l.fileStore.StoreRequest(requestResp.Request)
	if err != nil {
		logger.Error("store request error:%v %v", requestResp.Request.Id(), err)
		return nil
	}
	go func(request *common.ZkProofRequest) {
		count := 0
		for {
			defer l.worker.DelReqNum()
			if count >= 1 { // todo
				logger.Error("retry gen proof too much time,stop generate this proof now: %v", request.Id())
				return
			}
			count = count + 1
			logger.Info("worker %v start generate Proof type: %v", l.Id, request.Id())
			proofs, err := node.WorkerGenProof(l.worker, request)
			if err != nil {
				logger.Error("worker gen proof error:%v %v", request.Id(), err)
				continue
			}
			logger.Info("complete generate Proof type: %v", request.Id())
			submitProof := common.SubmitProof{Id: common.MustUUID(), WorkerId: l.Id, Data: proofs}
			_, err = l.client.SubmitProof(&submitProof)
			if err != nil {
				for _, proof := range proofs {
					logger.Error("submit proof %v_%v error cache now, %v", submitProof.Id, proof.Id(), err)
				}
				l.cacheProofs.Push(&submitProof)
				return
			}
			for _, proof := range proofs {
				logger.Info("success submit proof to server: %v %v", submitProof.Id, proof.Id())
			}
			return
		}
	}(requestResp.Request)
	return nil
}

func (l *Local) CheckState() error {
	logger.Debug("generator proof check state")
	l.cacheProofs.Iterator(func(value *common.SubmitProof) error {
		_, err := l.client.SubmitProof(value)
		if err != nil {
			logger.Error("submit proof error again:%v %v ", value.Id, err)
			return err
		}
		logger.Info("success submit proof again:%v", value.Id)
		l.cacheProofs.Delete(value.Id)
		return nil
	})
	return nil
}

func (l *Local) Close() error {
	return nil
}
