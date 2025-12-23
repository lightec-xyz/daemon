package node

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	btcproverClient "github.com/lightec-xyz/btc_provers/utils/client"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/rpc/bitcoin"
	"github.com/lightec-xyz/daemon/store"
)

var _ rpc.INode = (*Handler)(nil)

type Handler struct {
	fileStore    *FileStorage
	exitCh       chan os.Signal
	ethReScan    chan *ReScnSignal
	btcReScan    chan *ReScnSignal
	chainFork    chan *ChainFork
	manager      IManager
	txManager    *TxManager
	chainStore   *ChainStore
	btcClient    *bitcoin.Client         // todo
	proverClient btcproverClient.IClient // todo
	miner        string
	network      string
}

func (h *Handler) AutoSubmitMaxValue(max uint64) (string, error) {
	logger.Debug("set auto submit max value: %v", max)
	if max == 0 {
		return "", fmt.Errorf("max value is 0")
	}
	err := h.chainStore.WriteSubmitMaxValue(max)
	if err != nil {
		return "", err
	}
	h.txManager.setSubmitMax(max)
	return "ok", nil
}

func (h *Handler) AutoSubmitMinValue(min uint64) (string, error) {
	logger.Debug("set auto submit min value: %v", min)
	if min == 0 {
		return "", fmt.Errorf("min value is 0")
	}
	err := h.chainStore.WriteSubmitMinValue(min)
	if err != nil {
		return "", err
	}
	h.txManager.setSubmitMin(min)
	return "ok", nil
}

func (h *Handler) SetGasPrice(gasPrice uint64) (string, error) {
	logger.Warn("set gas price: %v", gasPrice)
	err := h.chainStore.WriteMaxGasPrice(gasPrice)
	if err != nil {
		return "", err
	}
	h.txManager.setMaxGasPrice(gasPrice)
	return "ok", nil
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

func (h *Handler) PendingTask() ([]*rpc.ProofTask, error) {
	proofList := h.manager.PendingProofRequest()
	var proofInfos []*rpc.ProofTask
	for _, proof := range proofList {
		proofInfos = append(proofInfos, &rpc.ProofTask{
			Id: proof.ProofId(),
		})
	}
	return proofInfos, nil
}

func (h *Handler) GetZkProofTask(request common.TaskRequest) (*common.TaskResponse, error) {
	if request.Version < GeneratorVersion {
		return nil, fmt.Errorf("generator version %v, less than node version %v,please upgrade generator", request.Version, GeneratorVersion)
	}
	zkReq, ok, err := h.manager.GetProofRequest(request.ProofType)
	if err != nil {
		logger.Error("get proof request error: %v %v", request.Id, err)
		return nil, err
	}
	response := common.TaskResponse{}
	if !ok {
		response.CanGen = false
		return &response, nil
	}
	reqData, err := json.Marshal(zkReq)
	if err != nil {
		logger.Error("marshal zk proof request error: workerId:%v,zkRequest:%v ,error:%v", request.Id, zkReq, err)
		return nil, err
	}
	response.CanGen = true
	response.Data = string(reqData)
	err = h.fileStore.StoreRequest(zkReq)
	if err != nil {
		logger.Error("store zk proof request error: %v %v", request.Id, err)
	}
	logger.Info("worker: %v get zk proof task proofId: %v, timestamp :%v txIndex:%v", request.Id, zkReq.ProofId(), zkReq.BlockTime, zkReq.TxIndex)
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
			logger.Error("read ethereum tx ids error: %v %v", height, err)
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
		//todo refactor
		params, err := getProofParams(txId, h.miner, h.network, h.chainStore, h.btcClient, h.proverClient)
		if err != nil {
			//logger.Warn("get proof params error: %v %v", txId, err)
			continue
		}
		results = append(results, rpc.ProofInfo{
			Status: int(proof.Status),
			Proof:  proof.Proof,
			TxId:   proof.TxHash,
			Params: &rpc.ProofParams{
				Checkpoint:        hex.EncodeToString(params.Checkpoint[:]),
				CpDepth:           params.CpDepth,
				TxDepth:           params.TxDepth,
				TxBlockHash:       hex.EncodeToString(params.TxBlockHash[:]),
				TxTimestamp:       params.TxTimestamp,
				ZkpMiner:          params.ZkpMiner.String(),
				Flag:              uint32(params.Flag.Int64()),
				SmoothedTimestamp: params.SmoothedTimestamp,
			},
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

func NewHandler(txManager *TxManager, manager IManager, ethReScan, btcReScan chan *ReScnSignal, store store.IStore, fileStore *FileStorage, exitCh chan os.Signal,
	btcClient *bitcoin.Client, proverClient btcproverClient.IClient, miner, network string) *Handler {
	return &Handler{
		txManager:    txManager,
		chainStore:   NewChainStore(store),
		exitCh:       exitCh,
		ethReScan:    ethReScan,
		btcReScan:    btcReScan,
		manager:      manager,
		fileStore:    fileStore,
		btcClient:    btcClient,
		proverClient: proverClient,
		miner:        miner,
		network:      network,
	}
}
