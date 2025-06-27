package node

import (
	"encoding/json"
	"fmt"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/store"
	"os"
	"sort"
	"strings"
	"syscall"
	"time"
)

var _ rpc.INode = (*Handler)(nil)

type Handler struct {
	fileStore  *FileStorage
	exitCh     chan os.Signal
	ethReScan  chan *ReScnSignal
	btcReScan  chan *ReScnSignal
	chainFork  chan *ChainFork
	manager    IManager
	chainStore *ChainStore
}

func (h *Handler) Eth2Slot(height uint64) (uint64, error) {
	slot, _, err := h.chainStore.ReadSlotByHeight(height)
	if err != nil {
		return 0, err
	}
	return slot, nil
}

func (h *Handler) Eth1Height(slot uint64) (uint64, error) {
	height, err := h.chainStore.ReadEthNumberBySlot(slot)
	if err != nil {
		return 0, err
	}
	return height, nil
}

func (h *Handler) ReScan(height uint64, chain string) error {
	logger.Debug("re scan height: %v, chain: %v", height, chain)
	if chain == common.EthereumChain.String() {
		h.ethReScan <- &ReScnSignal{Height: height}
	} else if chain == common.BitcoinChain.String() {
		h.btcReScan <- &ReScnSignal{Height: height}
	} else {
		logger.Error("reScan unsupported chain: %v", chain)
	}
	return nil
}

func (h *Handler) MinerInfo() ([]*rpc.MinerInfo, error) {
	miners, err := h.chainStore.ReadAllMiners()
	if err != nil {
		logger.Error("read miner powers error: %v", err)
		return nil, err
	}
	var exportRpcMiners []*rpc.MinerInfo
	for _, miner := range miners {
		power, err := h.chainStore.ReadMinerPower(miner)
		if err != nil {
			logger.Error("read miner powers error: %v", err)
			return nil, err
		}
		// only export miner power in last 1 hour
		if time.Now().Sub(time.Unix(int64(power.Timestamp), 0)) < 1*time.Hour {
			exportRpcMiners = append(exportRpcMiners, &rpc.MinerInfo{
				Address:   miner,
				Power:     power.Power,
				Timestamp: power.Timestamp,
			})
		}

	}
	return exportRpcMiners, nil
}

func (h *Handler) AddP2pPeer(addr string) (string, error) {
	err := h.manager.AddP2pPeer(addr)
	if err != nil {
		logger.Error("add p2p peer error:%v", err)
		return "", err
	}
	return "ok", nil
}

func (h *Handler) RemoveUnGenProof(hash string) (string, error) {
	err := h.chainStore.DeleteUnGenProof(common.BitcoinChain, hash)
	if err != nil {
		//logger.Error("remove ungen proof error: %v %v", hash, err)

	}
	err = h.chainStore.DeleteUnGenProof(common.EthereumChain, hash)
	if err != nil {
		//logger.Error("remove ungen proof error: %v %v", hash, err)
	}
	logger.Debug("remove ungen proof: %v", hash)
	return "ok", err
}

func (h *Handler) RemoveUnSubmitTx(hash string) (string, error) {
	err := h.chainStore.DeleteUnSubmitTx(hash)
	if err != nil {
		logger.Error("remove unsubmit tx error: %v %v", hash, err)
		return "", err
	}
	logger.Debug("remove unsubmit tx: %v", hash)
	return "ok", err
}

func (h *Handler) ProofTask(id string) (*rpc.ProofTaskInfo, error) {
	taskInfo, err := h.chainStore.ReadAllTaskTime(id)
	if err != nil {
		logger.Error("read queue time error: %v %v", id, err)
		return nil, err
	}
	logger.Info("proof task: %v, queue time: %v, generating time: %v, proof time: %v", id,
		taskInfo.QueueTime, taskInfo.GeneratingTime, taskInfo.EndTime)
	return &rpc.ProofTaskInfo{
		Id:             id,
		QueueTime:      taskInfo.QueueTime,
		GeneratingTime: taskInfo.GeneratingTime,
		EndTime:        taskInfo.EndTime,
	}, nil
}

func (h *Handler) PendingTask() ([]*rpc.ProofTaskInfo, error) {
	proofList := h.manager.PendingProofRequest()
	var proofInfos []*rpc.ProofTaskInfo
	for _, proof := range proofList {
		taskInfo, err := h.ProofTask(proof.FileKey.String())
		if err != nil {
			logger.Error("read proof task error: %v %v", proof.FileKey, err)
			return nil, err
		}
		proofInfos = append(proofInfos, taskInfo)
	}
	return proofInfos, nil
}

func (h *Handler) TxesByAddr(addr, txType string) ([]*rpc.Transaction, error) {
	if addr == "" || txType == "" {
		return nil, fmt.Errorf("addr or txType is empty")
	}
	tType, err := common.ToTxType(txType)
	if err != nil {
		logger.Error("to tx type error: %v %v %v", addr, txType, err)
		return nil, err
	}
	dbTxes, err := h.chainStore.ReadTxIdsByAddr(tType, addr)
	if err != nil {
		logger.Error("read addr txes error: %v %v %v", addr, txType, err)
		return nil, err
	}
	var rpcTxes []*rpc.Transaction
	for _, txId := range dbTxes {
		transaction, err := h.Transaction(txId)
		if err != nil {
			logger.Error("read transaction error: %v %v", txId, err)
			return nil, err
		}
		rpcTxes = append(rpcTxes, transaction...)
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
	if request.Version < GeneratorVersion {
		return nil, fmt.Errorf("generator version %v, less than node version %v,please upgrade generator", request.Version, GeneratorVersion)
	}
	zkProofRequest, ok, err := h.manager.GetProofRequest(request.ProofType)
	if err != nil {
		logger.Error("get proof request error: %v %v", request.Id, err)
		return nil, err
	}
	response := common.TaskResponse{}
	if !ok {
		response.CanGen = false
		return &response, nil
	}
	reqData, err := json.Marshal(zkProofRequest)
	if err != nil {
		logger.Error("marshal zk proof request error: workerId:%v,zkRequest:%v ,error:%v", request.Id, zkProofRequest, err)
		return nil, err
	}
	response.CanGen = true
	response.Data = string(reqData)
	err = h.fileStore.StoreRequest(zkProofRequest)
	if err != nil {
		logger.Error("store zk proof request error: %v %v", request.Id, err)
	}
	logger.Info("worker: %v get zk proof task proofId: %v,blockTime:%v", request.Id, zkProofRequest.ProofId(), zkProofRequest.BlockTime)
	return &response, nil
}

func (h *Handler) SubmitProof(req *common.SubmitProof) (string, error) {
	if req.Version < GeneratorVersion {
		return "", fmt.Errorf("generator version %v, less than node version %v,please upgrade generator", req.Version, GeneratorVersion)
	}
	for _, item := range req.Responses {
		logger.Debug("worker: %v submit proof: %v", req.WorkerId, item.ProofId())
	}
	for _, item := range req.Requests {
		logger.Warn("worker: %v submit fail request: %v", req.WorkerId, item.ProofId())
	}
	err := h.manager.ReceiveProofs(req)
	if err != nil {
		logger.Error("worker %v send proof to manager error: %v", req.WorkerId, err)
		return "", err
	}
	return "ok", nil
}

func (h *Handler) TransactionsByHeight(height uint64, network string) ([]string, error) {
	if network == BitcoinNetwork {
		txIds, err := h.chainStore.ReadBtcTxHeight(height)
		if err != nil {
			logger.Error("read bitcoin tx ids error: %v %v", height, err)
			return nil, err
		}
		return txIds, nil

	} else if network == EthereumNetwork {
		txIds, err := h.chainStore.ReadEthTxHeight(height)
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
		txList = append(txList, transaction...)
	}
	return txList, nil

}

func (h *Handler) Transaction(txHash string) ([]*rpc.Transaction, error) {
	txes, err := h.chainStore.ReadDbTxes(txHash)
	if err != nil {
		logger.Error("read transaction error: %v %v", txHash, err)
		return nil, err
	}
	var transactions []*rpc.Transaction
	for _, tx := range txes {
		destChainHash, _ := h.chainStore.ReadDestHash(txHash)
		dbProof, _ := h.chainStore.ReadDbProof(txHash)
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
		transactions = append(transactions, &transaction)
	}
	return transactions, err
}

func (h *Handler) ProofInfo(txIds []string) ([]rpc.ProofInfo, error) {
	var results []rpc.ProofInfo
	for _, txId := range txIds {
		if txId == "" {
			continue
		}
		proof, err := h.chainStore.ReadDbProof(txId)
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
	logger.Error("unsupported add worker")
	return "", nil
}

func (h *Handler) Version() (rpc.NodeInfo, error) {
	daemonInfo := rpc.NodeInfo{
		Version: "0.0.1",
		Network: "devnet", // todo
	}
	return daemonInfo, nil
}

func NewHandler(manager IManager, ethReScan, btcReScan chan *ReScnSignal, store store.IStore, fileStore *FileStorage, exitCh chan os.Signal) *Handler {
	return &Handler{
		chainStore: NewChainStore(store),
		exitCh:     exitCh,
		ethReScan:  ethReScan,
		btcReScan:  btcReScan,
		manager:    manager,
		fileStore:  fileStore,
	}
}
