package node

import (
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc/bitcoin"
	ethrpc "github.com/lightec-xyz/daemon/rpc/ethereum"
)

// Todo

type TaskManager struct {
	ethClient *ethrpc.Client
	btcClient *bitcoin.Client
	queue     *SubmitQueue
	address   string
	keyStore  *KeyStore
}

func NewTaskManager(keyStore *KeyStore, ethClient *ethrpc.Client, btcClient *bitcoin.Client) (*TaskManager, error) {
	address, err := keyStore.Address()
	if err != nil {
		return nil, err
	}
	return &TaskManager{
		ethClient: ethClient,
		btcClient: btcClient,
		keyStore:  keyStore,
		address:   address,
		queue:     NewSubmitQueue(),
	}, nil
}

func (t *TaskManager) AddTask(resp *common.ZkProofResponse) {
	logger.Info("add retry task: %v", resp.Id())
	t.queue.Push(resp.Id(), resp)
}

func (t *TaskManager) Check() error {
	t.queue.Iterator(func(value *common.ZkProofResponse) error {
		logger.Info("task check", value.Id())
		switch value.ZkProofType {
		case common.TxInEth2:
			_ = t.MintZkBtcRequest(value)
		case common.VerifyTxType:
			_ = t.UpdateUtxoChange(value)

		}
		return nil
	})
	return nil
}

func (t *TaskManager) MintZkBtcRequest(resp *common.ZkProofResponse) error {
	_, err := RedeemBtcTx(t.btcClient, resp.TxHash, resp.Proof)
	if err != nil {
		logger.Error("mint zk btc error: %v", err)
		return err
	}
	t.queue.Delete(resp.Id())
	return nil
}

func (t *TaskManager) UpdateUtxoChange(resp *common.ZkProofResponse) error {
	err := updateContractUtxoChange(t.ethClient, t.address, t.keyStore.GetPrivateKey(), []string{resp.TxHash}, resp.Proof)
	if err != nil {
		logger.Error("update utxo error: %v", err)
		return err
	}
	t.queue.Delete(resp.Id())
	return nil
}

func (t *TaskManager) Close() error {
	return nil
}

type Task struct {
	Id   string
	Data *common.ZkProofResponse
}
