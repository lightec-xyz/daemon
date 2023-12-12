package proof

import (
	"encoding/json"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/store"
)

var _ rpc.ProofAPI = (*Handler)(nil)

type Handler struct {
	memoryStore     store.IStore
	maxParallelNums int // The maximum number of proofs that can be generated at the same time
}

func NewHandler(memoryStore store.IStore, max int) *Handler {
	return &Handler{
		memoryStore:     memoryStore,
		maxParallelNums: max,
	}
}

func (h *Handler) GenZkProof(request rpc.ProofRequest) (rpc.ProofResponse, error) {
	//todo ffi
	logger.Debug("new proof request: %v", request)
	response := rpc.ProofResponse{}
	err := objParse(request, &response)
	if err != nil {
		logger.Error("parse proof request error:%v", err)
		return response, nil
	}
	logger.Debug("proof generated: %v", response)
	return response, nil
}

func (h *Handler) ProofInfo(proofId string) (rpc.ProofInfo, error) {
	logger.Debug("proof info: %v", proofId)
	status := rpc.ProofInfo{
		Status: 1,
		Msg:    "success",
	}
	return status, nil
}

func objParse(src, dest interface{}) error {
	marshal, err := json.Marshal(src)
	if err != nil {
		return err
	}
	err = json.Unmarshal(marshal, dest)
	if err != nil {
		return err
	}
	return nil
}
