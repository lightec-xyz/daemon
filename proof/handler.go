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
	worker      rpc.IWorker
}

func (h *Handler) BtcTimestamp(req *rpc.BtcTimestampRequest) (*rpc.ProofResponse, error) {
	response, err := h.worker.BtcTimestamp(req)
	if err != nil {
		logger.Error("btc timestamp error:%v", err)
		return nil, err
	}
	return response, nil
}

func (h *Handler) SyncCommOuter(req *rpc.SyncCommOuterRequest) (*rpc.ProofResponse, error) {
	response, err := h.worker.SyncCommOuter(req)
	if err != nil {
		logger.Error("sync comm outer error:%v", err)
		return nil, err
	}
	return response, nil
}

func (h *Handler) SyncCommDutyProve(req rpc.SyncCommDutyRequest) (*rpc.SyncCommDutyResponse, error) {
	response, err := h.worker.SyncCommDutyProve(req)
	if err != nil {
		logger.Error("gen sync comm duty proof error:%v", err)
		return nil, err
	}
	return response, nil
}

func (h *Handler) BtcDuperRecursiveProve(req *rpc.BtcDuperRecursiveRequest) (*rpc.ProofResponse, error) {
	response, err := h.worker.BtcDuperRecursiveProve(req)
	if err != nil {
		logger.Error("btc duper recursive prove error:%v", err)
		return nil, err
	}
	return response, nil
}

func (h *Handler) SyncCommInner(req *rpc.SyncCommInnerRequest) (*rpc.ProofResponse, error) {
	response, err := h.worker.SyncCommInner(req)
	if err != nil {
		logger.Error("sync comm inner error:%v", err)
		return nil, err
	}
	return response, nil
}

func (h *Handler) BackendRedeemProof(req *rpc.RedeemRequest) (*rpc.RedeemResponse, error) {
	response, err := h.worker.BackendRedeemProof(req)
	if err != nil {
		logger.Error("backend redeem proof error:%v", err)
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

func (h *Handler) SupportProofType() []common.ProofType {
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

func (h *Handler) RedeemProof(req *rpc.RedeemRequest) (*rpc.RedeemResponse, error) {
	response, err := h.worker.RedeemProof(req)
	if err != nil {
		logger.Error("gen redeem proof error:%v", err)
		return nil, err
	}
	return response, nil
}

func (h *Handler) SyncCommitUnitProve(req rpc.SyncCommUnitsRequest) (*rpc.SyncCommUnitsResponse, error) {
	response, err := h.worker.SyncCommitUnitProve(req)
	if err != nil {
		logger.Error("gen sync comm unit proof error:%v", err)
		return nil, err
	}
	return response, nil
}

func (h *Handler) SyncCommRecursiveProve(req rpc.SyncCommDutyRequest) (*rpc.SyncCommDutyResponse, error) {
	response, err := h.worker.SyncCommDutyProve(req)
	if err != nil {
		logger.Error("gen sync comm recursive proof error:%v", err)
		return nil, err
	}
	return response, nil
}

func (h *Handler) ProofInfo(proofId string) (rpc.ProofInfo, error) {
	return rpc.ProofInfo{}, nil
}

func NewHandler(store, memoryStore store.IStore, worker rpc.IWorker) *Handler {
	return &Handler{
		memoryStore: memoryStore,
		worker:      worker,
	}
}
