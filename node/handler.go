package node

import (
	"fmt"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/store"
	"os"
	"syscall"
)

var _ rpc.NodeAPI = (*Handler)(nil)

type Handler struct {
	store    store.IStore
	memoryDb store.IStore
	exitCh   chan os.Signal
	schedule *Schedule
}

func (h *Handler) Transaction(txHash string) (rpc.Transaction, error) {
	panic(h)
}

func (h *Handler) ProofInfo(txId string) (rpc.ProofInfo, error) {
	proofId := TxIdToProofId(txId)
	var txProof TxProof
	err := h.store.GetObj(proofId, &txProof)
	if err != nil {
		logger.Error("get proof error: %v %v", proofId, err)
		return rpc.ProofInfo{}, err
	}
	result := rpc.ProofInfo{
		Status: txProof.Status,
		Msg:    txProof.Msg,
		Proof:  txProof.Proof,
	}
	return result, nil
}

func (h *Handler) Stop() error {
	logger.Debug("rpc handler receive stop signal")
	h.exitCh <- syscall.SIGQUIT
	return nil
}

func (h *Handler) AddWorker(endpoint string, max int) (string, error) {
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

func NewHandler(store, memoryDb store.IStore, schedule *Schedule, exitCh chan os.Signal) *Handler {
	return &Handler{
		store:    store,
		memoryDb: memoryDb,
		exitCh:   exitCh,
		schedule: schedule,
	}
}

func (h *Handler) HelloWorld(name *string, age *int) (string, error) {
	fmt.Printf("req: %v %v", name, age)
	return fmt.Sprintf(" name %v age %v", *name, *age), nil
}
