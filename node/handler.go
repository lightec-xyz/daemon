package node

import (
	"fmt"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/store"
	"os"
	"syscall"
)

var _ rpc.INode = (*Handler)(nil)

type Handler struct {
	store    store.IStore
	memoryDb store.IStore
	exitCh   chan os.Signal
	schedule *Schedule
	manager  *manager
}

func (h *Handler) GetTask(request *common.TaskRequest) (*common.TaskResponse, error) {
	// Todo
	zkProofRequest, ok, err := h.manager.GetProofRequest()
	if err != nil {
		logger.Error("get proof request error: %v", err)
		return nil, err
	}
	var response common.TaskResponse
	if !ok {
		logger.Warn("maybe no new proof task")
		response.CanGen = false
		return &response, nil
	}
	response.CanGen = true
	response.Request = zkProofRequest
	logger.Info("worker: %v get task: type:%v hash:%v %v:period:%v", request.Id, zkProofRequest.ReqType.String(),
		zkProofRequest.TxHash, zkProofRequest.Period)
	return &response, nil
}

func (h *Handler) SubmitProof(req *common.SubmitProof) (string, error) {
	//todo check
	err := h.manager.SendProofResponse(req.Data)
	if err != nil {
		logger.Error("send proof to manager error: %v", err)
		return "", err
	}
	return "ok", nil
}

func (h *Handler) TransactionsByHeight(height uint64, network string) ([]string, error) {
	if network == BitcoinNetwork {
		txIds, err := ReadBitcoinTxIds(h.store, int64(height))
		if err != nil {
			logger.Error("read bitcoin tx ids error: %v %v", height, err)
			return nil, err
		}
		return txIds, nil

	} else if network == EthereumNetwork {
		txIds, err := ReadEthereumTxIds(h.store, int64(height))
		if err != nil {
			logger.Error("read bitcoin tx ids error: %v %v", height, err)
			return nil, err
		}
		return txIds, nil
	} else {
		return nil, fmt.Errorf("unsupported network: %v", network)
	}

}

func (h *Handler) Transactions(txIds []string) ([]rpc.Transaction, error) {
	var txList []rpc.Transaction
	for _, txId := range txIds {
		transaction, err := h.Transaction(txId)
		if err != nil {
			logger.Error("read transaction error: %v %v", txId, err)
			return nil, err
		}
		txList = append(txList, transaction)
	}
	return txList, nil

}

func (h *Handler) Transaction(txHash string) (rpc.Transaction, error) {
	tx, err := ReadDbTx(h.store, txHash)
	if err != nil {
		logger.Error("read transaction error: %v %v", txHash, err)
		return rpc.Transaction{}, err
	}
	transaction := rpc.Transaction{}
	err = objParse(tx, &transaction)
	if err != nil {
		logger.Error("parse transaction error: %v %v", txHash, err)
		return rpc.Transaction{}, err
	}
	return transaction, err
}

func (h *Handler) ProofInfo(txIds []string) ([]rpc.ProofInfo, error) {
	var results []rpc.ProofInfo
	for _, txId := range txIds {
		proof, err := ReadDbProof(h.store, txId)
		if err != nil {
			logger.Error("read Proof error: %v %v", txId, err)
			return nil, err
		}
		results = append(results, rpc.ProofInfo{
			Status: int(proof.Status),
			Proof:  proof.Proof,
			TxId:   proof.TxHash,
		})
	}
	return results, nil
}

func (h *Handler) Stop() error {
	logger.Debug("rpc handler receive stop signal")
	h.exitCh <- syscall.SIGQUIT
	return nil
}

func (h *Handler) AddWorker(endpoint string, max int) (string, error) {
	logger.Info("add new worker now: %v %v", endpoint, max)
	err := h.schedule.AddWorker(endpoint, max)
	if err != nil {
		return "", err
	}
	return "success", err
}

func (h *Handler) Version() (rpc.NodeInfo, error) {
	daemonInfo := rpc.NodeInfo{
		Version: "0.0.1",
	}
	return daemonInfo, nil
}

func NewHandler(manager *manager, store, memoryDb store.IStore, schedule *Schedule, exitCh chan os.Signal) *Handler {
	return &Handler{
		store:    store,
		memoryDb: memoryDb,
		exitCh:   exitCh,
		schedule: schedule,
		manager:  manager,
	}
}

func (h *Handler) HelloWorld(name *string, age *int) (string, error) {
	fmt.Printf("req: %v %v", name, age)
	return fmt.Sprintf(" name %v age %v", *name, *age), nil
}
