package proof

import (
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/node"
	"github.com/lightec-xyz/daemon/rpc"
	"sync"
	"time"
)

type Local struct {
	Id                string
	client            *rpc.NodeClient
	worker            rpc.IWorker
	exit              chan struct{}
	pendingProofsList *sync.Map
}

func NewLocal(url, datadir string, num int) (*Local, error) {
	client, err := rpc.NewNodeClient(url)
	if err != nil {
		logger.Error("new node client error:%v", err)
		return nil, err
	}
	worker, err := node.NewLocalWorker(datadir, datadir, int(num))
	if err != nil {
		logger.Error("new local worker error:%v", err)
		return nil, err
	}
	logger.Info("workerId: %v", worker.Id())
	return &Local{
		client:            client,
		worker:            worker,
		Id:                worker.Id(),
		exit:              make(chan struct{}, 1),
		pendingProofsList: new(sync.Map),
	}, nil
}

func (l *Local) Run() error {
	for {
		select {
		case <-l.exit:
			return nil
		default:
			time.Sleep(1 * time.Minute)
		}
		if l.worker.CurrentNums() >= l.worker.MaxNums() {
			logger.Warn("maxNums limit reached, wait proof generated")
			continue
		}
		request := common.TaskRequest{
			Id:        l.Id,
			ProofType: []common.ZkProofType{}, // Todo worker support which proof type
		}
		requestResp, err := l.client.GetTask(request)
		if err != nil {
			logger.Error("get task error:%v", err)
			continue
		}
		if !requestResp.CanGen {
			logger.Debug("no new proof request, wait  request coming now ....")
			continue
		}
		l.worker.AddReqNum()
		go func(request *common.ZkProofRequest) {
			count := 0
			for {
				count = count + 1
				if count > 10 {
					logger.Error("retry gen proof too much time,stop generate this proof now: %v %v %v", request.Period, request.TxHash, request.ReqType.String())
					return
				}
				logger.Info("worker %v start generate Proof type: %v Period: %v", l.worker.Id(), request.ReqType.String(), request.Period)
				proof, err := node.WorkerGenProof(l.worker, request)
				if err != nil {
					logger.Error("worker gen proof error:%v", err)
					continue
				}
				logger.Info("complete generate Proof type: %v Period: %v", request.ReqType.String(), request.Period)
				submitProof := common.SubmitProof{Id: common.MustUUID(), WorkerId: l.Id, Data: proof}
				_, err = l.client.SubmitProof(submitProof)
				if err != nil {
					logger.Error("submit proof error:%v,store proof: %v", err, submitProof.Id)
					l.pendingProofsList.Store(submitProof.Id, &submitProof)
					// todo check again
					return
				}
				logger.Info("submit proof to daemon type: %v Period: %v,txHash: %v", request.ReqType.String(), request.Period, request.TxHash)
				return
			}
		}(requestResp.Request)

	}
}

func (l *Local) CheckState() error {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-l.exit:
			logger.Info("%v goroutine exit now ...", "CheckState")
			return nil
		case <-ticker.C:
			err := l.checkPendingProof()
			if err != nil {
				logger.Error("check pending proof error:%v", err)
			}
		}
	}
}

func (l *Local) checkPendingProof() error {
	l.pendingProofsList.Range(func(key, value any) bool {
		proof, ok := value.(*common.SubmitProof)
		if !ok {
			logger.Error("value is not SubmitProof")
			return false
		}
		_, err := l.client.SubmitProof(*proof)
		if err != nil {
			logger.Error("submit proof error again:%v %v ", key, err)
			return false
		}
		logger.Info("success submit proof again:%v", key)
		l.pendingProofsList.Delete(key)
		return true
	})
	return nil
}

func (l *Local) Close() error {
	close(l.exit)
	return nil
}
