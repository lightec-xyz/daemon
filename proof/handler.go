package proof

import (
	"github.com/lightec-xyz/daemon/common"
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

func (h *Handler) BtcDuperRecursiveProve(req *rpc.BtcDuperRecursiveRequest) (*rpc.ProofResponse, error) {
	response, err := h.worker.BtcDuperRecursiveProve(req)
	if err != nil {
		logger.Error("btc duper recursive prove error:%v", err)
		return nil, err
	}
	return response, nil
}

func (h *Handler) BtcDepthRecursiveProve(req *rpc.BtcDepthRecursiveRequest) (*rpc.ProofResponse, error) {
	response, err := h.worker.BtcDepthRecursiveProve(req)
	if err != nil {
		logger.Error("btc depth recursive prove error:%v", err)
		return nil, err
	}
	return response, nil
}

func (h *Handler) BtcChainProve(req *rpc.BtcChainRequest) (*rpc.ProofResponse, error) {
	response, err := h.worker.BtcChainProve(req)
	if err != nil {
		logger.Error("btc chain prove error:%v", err)
		return nil, err
	}
	return response, nil
}

func (h *Handler) BtcDepositProve(req *rpc.BtcDepositRequest) (*rpc.ProofResponse, error) {
	response, err := h.worker.BtcDepositProve(req)
	if err != nil {
		logger.Error("btc deposit prove error:%v", err)
		return nil, err
	}
	return response, nil
}

func (h *Handler) BtcChangeProve(req *rpc.BtcChangeRequest) (*rpc.ProofResponse, error) {
	response, err := h.worker.BtcChangeProve(req)
	if err != nil {
		logger.Error("btc change prove error:%v", err)
		return nil, err
	}
	return response, nil
}

func (h *Handler) BtcBaseProve(req *rpc.BtcBaseRequest) (*rpc.ProofResponse, error) {
	response, err := h.worker.BtcBaseProve(req)
	if err != nil {
		logger.Error("btc base prove error:%v", err)
		return nil, err
	}
	return response, nil
}

func (h *Handler) BtcMiddleProve(req *rpc.BtcMiddleRequest) (*rpc.ProofResponse, error) {
	response, err := h.worker.BtcMiddleProve(req)
	if err != nil {
		logger.Error("btc middle prove error:%v", err)
		return nil, err
	}
	return response, nil
}

func (h *Handler) BtcUpperProve(req *rpc.BtcUpperRequest) (*rpc.ProofResponse, error) {
	response, err := h.worker.BtcUpperProve(req)
	if err != nil {
		logger.Error("btc up prove error:%v", err)
		return nil, err
	}
	return response, nil
}

func (h *Handler) SupportProofType() []common.ZkProofType {
	return nil
}

func (h *Handler) Close() error {
	err := h.worker.Close()
	return err
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

func (h *Handler) GenRedeemProof(req *rpc.RedeemRequest) (*rpc.RedeemResponse, error) {
	response, err := h.worker.GenRedeemProof(req)
	if err != nil {
		logger.Error("gen redeem proof error:%v", err)
		return nil, err
	}
	return response, nil
}

func (h *Handler) GenSyncCommGenesisProof(req rpc.SyncCommGenesisRequest) (*rpc.SyncCommGenesisResponse, error) {
	response, err := h.worker.GenSyncCommGenesisProof(req)
	if err != nil {
		logger.Error("gen sync comm genesis proof error:%v", err)
		return nil, err
	}
	return response, nil
}

func (h *Handler) GenSyncCommitUnitProof(req rpc.SyncCommUnitsRequest) (*rpc.SyncCommUnitsResponse, error) {
	response, err := h.worker.GenSyncCommitUnitProof(req)
	if err != nil {
		logger.Error("gen sync comm unit proof error:%v", err)
		return nil, err
	}
	return response, nil
}

func (h *Handler) GenSyncCommRecursiveProof(req rpc.SyncCommRecursiveRequest) (*rpc.SyncCommRecursiveResponse, error) {
	response, err := h.worker.GenSyncCommRecursiveProof(req)
	if err != nil {
		logger.Error("gen sync comm recursive proof error:%v", err)
		return nil, err
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
