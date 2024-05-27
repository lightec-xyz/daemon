package proof

import (
	"encoding/json"
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
	proofTypes  []common.ZkProofType
	cacheProofs *node.ProofRespQueue
}

func NewLocal(zkParamDir, url, datadir, id string, proofTypes []common.ZkProofType, num int, store store.IStore, fileStore *node.FileStorage) (*Local, error) {
	client, err := rpc.NewNodeClient(url)
	if err != nil {
		logger.Error("new node client error:%v", err)
		return nil, err
	}
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
		proofTypes:  proofTypes,
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
		ProofType: l.proofTypes,
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

func (l *Local) Init() error {
	submitProofs, err := node.ReadAllProofResponse(l.store)
	if err != nil {
		logger.Error("read all proof response error:%v", err)
		return err
	}
	for _, resp := range submitProofs {
		for _, item := range resp.Data {
			logger.Debug("load pending proof response:%v %v", resp.Id, item.Id())
		}
		_, err := l.client.SubmitProof(resp)
		if err != nil {
			logger.Error("submit proof error again:%v %v", resp.Id, err)
			l.cacheProofs.Push(resp)
			logger.Debug("add proof response to pending queue:%v", resp.Id)
		} else {
			logger.Debug("success submit proof again:%v", resp.Id)
		}
		// todo
		err = node.DeleteProofResponse(l.store, resp.Id)
		if err != nil {
			logger.Error("delete proof response error:%v", err)
			return err
		}

	}
	return nil
}

func (l *Local) Close() error {
	logger.Debug("store cache data to db now ...")
	l.cacheProofs.Iterator(func(value *common.SubmitProof) error {
		err := node.WriteProofResponse(l.store, value)
		if err != nil {
			logger.Error("write proof response error:%v", err)
			return err
		}
		for _, item := range value.Data {
			logger.Debug("write pending proof response: %v %v", value.Id, item.Id())
		}
		return nil
	})

	return nil
}

func ReadParameters(data []byte) ([]*common.Parameters, error) {
	var result []*common.Parameters
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
