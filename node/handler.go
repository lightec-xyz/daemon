package node

import (
	"fmt"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/store"
	"os"
	"sort"
	"strings"
	"syscall"
)

var _ rpc.INode = (*Handler)(nil)

type Handler struct {
	store     store.IStore
	memoryDb  store.IStore
	fileStore *FileStorage
	exitCh    chan os.Signal
	schedule  *Schedule
	manager   IManager
}

func (h *Handler) ProofTask(id string) (*rpc.ProofTaskInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) PendingTask() ([]*rpc.ProofTaskInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) TxesByAddr(addr, txType string) ([]*rpc.Transaction, error) {
	if addr == "" || txType == "" {
		return nil, fmt.Errorf("addr or txType is empty")
	}
	dbTxes, err := ReadTxesByAddr(h.store, addr)
	if err != nil {
		logger.Error("read addr txes error: %v %v %v", addr, txType, err)
		return nil, err
	}
	var rpcTxes []*rpc.Transaction
	for _, tx := range dbTxes {
		transaction, err := h.Transaction(tx.TxHash)
		if err != nil {
			logger.Error("read transaction error: %v %v", tx.TxHash, err)
			return nil, err
		}
		rpcTxes = append(rpcTxes, transaction)
	}
	sort.SliceStable(rpcTxes, func(i, j int) bool {
		if rpcTxes[i].Height == rpcTxes[j].Height {
			return rpcTxes[i].TxIndex < rpcTxes[j].TxIndex
		}
		return rpcTxes[i].Height < rpcTxes[j].Height
	})
	return rpcTxes, nil

}

func (h *Handler) GetZkProofTask(request common.TaskRequest) (*common.TaskResponse, error) {
	// Todo
	if request.Version != GeneratorVersion {
		logger.Error("id: %v current version: %v, unsupported version: %v", request.Id, GeneratorVersion, request.Version)
		return nil, fmt.Errorf("current version: %v, unsupported version: %v", GeneratorVersion, request.Version)
	}

	zkProofRequest, ok, err := h.manager.GetProofRequest(request.ProofType)
	if err != nil {
		logger.Error("get proof request error: %v %v", request.Id, err)
		return nil, err
	}
	var response common.TaskResponse
	if !ok {
		logger.Warn("workerId: %v ,rpcServer maybe no new proof task", request.Id)
		response.CanGen = false
		return &response, nil
	}
	response.CanGen = true
	response.Request = zkProofRequest
	logger.Info("worker: %v get zk proof task: %v", request.Id, zkProofRequest.Id())
	return &response, nil
}

func (h *Handler) SubmitProof(req *common.SubmitProof) (string, error) {
	//if req.Version != GeneratorVersion {
	//	logger.Error("id: %v current version: %v, unsupported version: %v", req.Id, GeneratorVersion, req.Version)
	//	return "", fmt.Errorf("current version: %v, unsupported version: %v", GeneratorVersion, req.Version)
	//}
	for _, item := range req.Data {
		logger.Info("workerId %v,submit proof %v", req.WorkerId, item.Id())
		// todo
		err := StoreZkProof(h.fileStore, item.ZkProofType, item.Index, item.End, item.TxHash, item.Proof, item.Witness)
		if err != nil {
			logger.Error("store zk proof error: %v %v", item.Id(), err)
			return "", err
		}
	}
	err := h.manager.SendProofResponse(req.Data)
	if err != nil {
		logger.Error("worker %v send proof to manager error: %v", req.WorkerId, err)
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

func (h *Handler) Transactions(txIds []string) ([]*rpc.Transaction, error) {
	var txList []*rpc.Transaction
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

func (h *Handler) Transaction(txHash string) (*rpc.Transaction, error) {
	tx, err := ReadDbTx(h.store, txHash)
	if err != nil {
		logger.Error("read transaction error: %v %v", txHash, err)
		return nil, err
	}
	destChainHash, err := ReadDestHash(h.store, txHash)
	if err != nil {
		//logger.Error("read dest chain hash error: %v %v", txHash, err)
		//return nil, err
	}
	dbProof, err := ReadDbProof(h.store, txHash)
	if err != nil {
		//logger.Error("read dbProof error: %v %v", txHash, err)
		//return nil, err
	}
	transaction := rpc.Transaction{
		Height:    tx.Height,
		Hash:      txHash,
		TxIndex:   tx.TxIndex,
		TxType:    tx.TxType.String(),
		ChainType: tx.ChainType.String(),
		Amount:    tx.Amount,
		DestChain: rpc.DestChainInfo{
			Hash: destChainHash,
		},
		Proof: rpc.ProofInfo{
			TxId:   txHash,
			Proof:  dbProof.Proof,
			Status: dbProof.Status,
		},
	}
	return &transaction, err
}

func (h *Handler) ProofInfo(txIds []string) ([]rpc.ProofInfo, error) {
	var results []rpc.ProofInfo
	for _, txId := range txIds {
		if txId == "" {
			continue
		}
		proof, err := ReadDbProof(h.store, txId)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				continue
			}
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

func NewHandler(manager IManager, store, memoryDb store.IStore, schedule *Schedule, fileStore *FileStorage, exitCh chan os.Signal) *Handler {
	return &Handler{
		store:     store,
		memoryDb:  memoryDb,
		exitCh:    exitCh,
		schedule:  schedule,
		manager:   manager,
		fileStore: fileStore,
	}
}

func (h *Handler) HelloWorld(name *string, age *int) (string, error) {
	fmt.Printf("req: %v %v", name, age)
	return fmt.Sprintf(" name %v age %v", *name, *age), nil
}
