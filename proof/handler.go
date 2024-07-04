package proof

import (
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/store"
)

var _ rpc.IProof = (*Handler)(nil)

type Handler struct {
	memoryStore store.IStore
	store       store.IStore
	worker      rpc.IWorker
}

func (h *Handler) BtcBulkProve(req *rpc.BtcBulkRequest) (*rpc.BtcBulkResponse, error) {
	response, err := h.worker.BtcBulkProve(req)
	if err != nil {
		logger.Error("btc bulk prove error:%v", err)
		return nil, err
	}
	return response, nil
}

func (h *Handler) BtcPackedRequest(req *rpc.BtcPackedRequest) (*rpc.BtcPackResponse, error) {
	response, err := h.worker.BtcPackedRequest(req)
	if err != nil {
		logger.Error("btc pack request error:%v", err)
		return nil, err
	}
	return response, nil
}

func (h *Handler) BtcWrapProve(req *rpc.BtcWrapRequest) (*rpc.BtcWrapResponse, error) {
	response, err := h.worker.BtcWrapProve(req)
	if err != nil {
		logger.Error("btc wrap prove error:%v", err)
		return nil, err
	}
	return response, nil
}

func (h *Handler) TxInEth2Prove(req *rpc.TxInEth2ProveRequest) (*rpc.TxInEth2ProveResponse, error) {
	response, err := h.worker.TxInEth2Prove(req)
	if err != nil {
		logger.Error("tx in eth2 prove error:%v", err)
		return nil, err
	}
	return response, nil
}

func (h *Handler) BlockHeaderProve(req *rpc.BlockHeaderRequest) (*rpc.BlockHeaderResponse, error) {
	response, err := h.worker.BlockHeaderProve(req)
	if err != nil {
		logger.Error("block header prove error:%v", err)
		return nil, err
	}
	return response, nil
}

func (h *Handler) BlockHeaderFinalityProve(req *rpc.BlockHeaderFinalityRequest) (*rpc.BlockHeaderFinalityResponse, error) {
	response, err := h.worker.BlockHeaderFinalityProve(req)
	if err != nil {
		logger.Error("block header finality prove error:%v", err)
		return nil, err
	}
	return response, nil
}

func (h *Handler) GenDepositProof(req rpc.DepositRequest) (rpc.DepositResponse, error) {
	response, err := h.worker.GenDepositProof(req)
	if err != nil {
		logger.Error("gen deposit proof error:%v", err)
		return rpc.DepositResponse{}, err
	}
	return response, nil
}

func (h *Handler) GenRedeemProof(req *rpc.RedeemRequest) (*rpc.RedeemResponse, error) {
	response, err := h.worker.GenRedeemProof(req)
	if err != nil {
		logger.Error("gen redeem proof error:%v", err)
		return nil, err
	}
	return response, nil
}

func (h *Handler) GenVerifyProof(req rpc.VerifyRequest) (rpc.VerifyResponse, error) {
	logger.Debug("gen verify proof ")
	response, err := h.worker.GenVerifyProof(req)
	if err != nil {
		logger.Error("gen verify proof error:%v", err)
		return rpc.VerifyResponse{}, err
	}
	return response, nil
}

func (h *Handler) GenSyncCommGenesisProof(req rpc.SyncCommGenesisRequest) (rpc.SyncCommGenesisResponse, error) {
	response, err := h.worker.GenSyncCommGenesisProof(req)
	if err != nil {
		logger.Error("gen sync comm genesis proof error:%v", err)
		return rpc.SyncCommGenesisResponse{}, err
	}
	return response, nil
}

func (h *Handler) GenSyncCommitUnitProof(req rpc.SyncCommUnitsRequest) (rpc.SyncCommUnitsResponse, error) {
	response, err := h.worker.GenSyncCommitUnitProof(req)
	if err != nil {
		logger.Error("gen sync comm unit proof error:%v", err)
		return rpc.SyncCommUnitsResponse{}, err
	}
	return response, nil
}

func (h *Handler) GenSyncCommRecursiveProof(req rpc.SyncCommRecursiveRequest) (rpc.SyncCommRecursiveResponse, error) {
	response, err := h.worker.GenSyncCommRecursiveProof(req)
	if err != nil {
		logger.Error("gen sync comm recursive proof error:%v", err)
		return rpc.SyncCommRecursiveResponse{}, err
	}
	return response, nil
}

func (h *Handler) ProofInfo(proofId string) (rpc.ProofInfo, error) {
	return rpc.ProofInfo{}, nil
}

func (h *Handler) MaxNums() (int, error) {
	maxNums := h.worker.MaxNums()
	return maxNums, nil
}

func (h *Handler) CurrentNums() (int, error) {
	currentNums := h.worker.CurrentNums()
	return currentNums, nil
}

func (h *Handler) AddReqNum() {
	h.worker.AddReqNum()
}

func (h *Handler) DelReqNum() {
	h.worker.DelReqNum()
}

func NewHandler(store, memoryStore store.IStore, worker rpc.IWorker) *Handler {
	return &Handler{
		memoryStore: memoryStore,
		store:       store,
		worker:      worker,
	}
}
