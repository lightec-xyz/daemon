package proof

import (
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/node"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/store"
)

type Local struct {
	workerId    string
	client      *rpc.NodeClient
	worker      rpc.IWorker
	fileStore   *node.FileStorage
	exit        chan struct{}
	proofTypes  []common.ProofType
	cacheProofs *node.ProofRespQueue
	chainStore  *node.ChainStore
}

func NewLocal(url string, worker rpc.IWorker, proofTypes []common.ProofType, store store.IStore, fileStore *node.FileStorage) (*Local, error) {
	client, err := rpc.NewNodeClient(url, "")
	if err != nil {
		logger.Error("new node client error:%v", err)
		return nil, err
	}
	return &Local{
		client:      client,
		fileStore:   fileStore,
		worker:      worker,
		workerId:    worker.Id(),
		chainStore:  node.NewChainStore(store),
		proofTypes:  proofTypes,
		exit:        make(chan struct{}, 1),
		cacheProofs: node.NewProofRespQueue(),
	}, nil
}
func (l *Local) Init() error {
	submitProofs, err := l.chainStore.ReadAllProofResponse()
	if err != nil {
		logger.Error("read all proof response error:%v", err)
		return err
	}
	for _, resp := range submitProofs {
		for _, item := range resp.Responses {
			logger.Debug("load pending proof response:%v %v", resp.Id, item.ProofId())
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
		err = l.chainStore.DeleteProofResponse(resp.Id)
		if err != nil {
			logger.Error("delete proof response error:%v", err)
			return err
		}

	}
	return nil
}

func (l *Local) polling() error {
	if l.worker.CurrentNums() >= l.worker.MaxNums() {
		//logger.Warn("maxNums limit reached, wait proof generated")
		return nil
	}
	request := common.TaskRequest{
		Id:        l.workerId,
		ProofType: l.proofTypes,
		Version:   node.GeneratorVersion,
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
	zkRequest, err := requestResp.ToRequest()
	if err != nil {
		logger.Error("unmarshal request error:%v", err)
		return nil
	}
	err = l.fileStore.StoreRequest(zkRequest)
	if err != nil {
		logger.Error("store request error:%v %v", zkRequest.ProofId(), err)
		//return nil
	}
	l.worker.AddReqNum()
	go func(request *common.ProofRequest) {
		defer l.worker.DelReqNum()
		logger.Info("worker %v start generate Proof type: %v", l.workerId, request.ProofId())
		proofs, err := node.WorkerGenProof(l.worker, request)
		if err != nil {
			logger.Error("worker gen proof error:%v %v", request.ProofId(), err)
			_, err := l.client.SubmitProof(NewSubmitProof(l.workerId, false, nil, []*common.ProofRequest{request}))
			if err != nil {
				logger.Error("submit proof error:%v %v", request.ProofId(), err)
			}
			return
		}
		logger.Info("complete generate Proof type: %v", request.ProofId())
		submitProof := NewSubmitProof(l.workerId, true, proofs, nil)
		_, err = l.client.SubmitProof(submitProof)
		if err != nil {
			for _, proof := range proofs {
				logger.Error("submit proof %v %v error cache now, %v", submitProof.Id, proof.ProofId(), err)
			}
			l.cacheProofs.Push(submitProof)
			return
		}
		for _, proof := range proofs {
			logger.Info("success submit proof to server: workerId %v,proofId %v", submitProof.WorkerId, proof.ProofId())
			storeKey := node.NewStoreKey(proof.ProofType, proof.Hash, proof.Prefix, proof.FIndex, proof.SIndex)
			err := l.fileStore.StoreProof(storeKey, proof.Proof, proof.Witness)
			if err != nil {
				logger.Error("store proof error:%v %v", proof.ProofId(), err)
			}
		}
	}(zkRequest)
	return nil
}

func (l *Local) CheckState() error {
	l.cacheProofs.Iterator(func(value *common.SubmitProof) error {
		_, err := l.client.SubmitProof(value)
		if err == nil {
			for _, resp := range value.Responses {
				logger.Info("success submit proof again:%v,proofId: %v", value.Id, resp.ProofId())
			}
		}
		l.cacheProofs.Delete(value.Id)
		return nil
	})
	return nil
}

func (l *Local) Close() error {
	logger.Debug("store cache data to db now ...")
	l.cacheProofs.Iterator(func(value *common.SubmitProof) error {
		err := l.chainStore.WriteProofResponse(value)
		if err != nil {
			logger.Error("write proof response error:%v", err)
			return err
		}
		for _, item := range value.Responses {
			logger.Debug("write pending proof response: %v %v", value.Id, item.ProofId())
		}
		return nil
	})
	return nil

}

func NewSubmitProof(workerId string, status bool, proofs []*common.ProofResponse, requests []*common.ProofRequest) *common.SubmitProof {
	return &common.SubmitProof{
		Id:        common.MustUUID(),
		WorkerId:  workerId,
		Status:    status,
		Version:   node.GeneratorVersion,
		Responses: proofs,
		Requests:  requests,
	}
}
