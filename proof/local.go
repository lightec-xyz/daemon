package proof

import (
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/node"
	"github.com/lightec-xyz/daemon/rpc"
	"time"
)

type Local struct {
	Id     string
	client *rpc.NodeClient
	worker rpc.IWorker
	exit   chan struct{}
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
	return &Local{
		client: client,
		worker: worker,
		Id:     worker.Id(),
		exit:   make(chan struct{}, 1),
	}, nil
}

func (l *Local) Run() error {
	for {
		select {
		case <-l.exit:
			return nil
		default:
			time.Sleep(10 * time.Second)
		}
		if l.worker.CurrentNums() >= l.worker.MaxNums() {
			logger.Warn("wait proof generated")
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
			logger.Info("daemon server no task need to generate proof now ....")
			continue
		}
		l.worker.AddReqNum()
		go func(request common.ZkProofRequest) {
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
				_, err = l.client.SubmitProof(common.SubmitProof{Data: proof})
				if err != nil {
					logger.Error("submit proof error:%v", err)
					continue // Todo ,retry should in queue
				}
				logger.Info("submit proof to daemon type: %v Period: %v,txHash: %v %v", request.ReqType.String(), request.Period, request.TxHash)
				return
			}
		}(requestResp.Request)

	}
}

func (l *Local) Close() error {
	close(l.exit)
	return nil
}
