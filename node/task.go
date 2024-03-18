package node

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc/bitcoin"
	btcTypes "github.com/lightec-xyz/daemon/rpc/bitcoin/types"
	"github.com/lightec-xyz/daemon/rpc/ethereum"
	"github.com/lightec-xyz/daemon/rpc/oasis"
	"math/big"
	"sync"
	"time"
)

type TaskManager struct {
	ethNonce    uint64
	oasisNonce  uint64
	ethClient   *ethereum.Client
	btcClient   *bitcoin.Client
	oasisClient *oasis.Client
	address     string
	lock        sync.Mutex
	exit        chan struct{}
	queue       *sync.Map
	keyStore    *KeyStore
	timeout     time.Duration
}

func NewTaskManager(address, privateKey string, ethClient *ethereum.Client, btcClient *bitcoin.Client, oasisClient *oasis.Client) (*TaskManager, error) {
	ethNonce, err := ethClient.GetNonce(address)
	if err != nil {
		logger.Error("get ethNonce error:%v", err)
		return nil, err
	}
	return &TaskManager{
		queue:       new(sync.Map),
		ethNonce:    ethNonce,
		oasisNonce:  0,
		address:     address,
		ethClient:   ethClient,
		btcClient:   btcClient,
		oasisClient: oasisClient,
		keyStore:    NewKeyStore(privateKey),
		exit:        make(chan struct{}, 1),
		timeout:     5 * time.Minute,
	}, nil
}

func (t *TaskManager) GetEthNewNonce() (uint64, error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	chainNonce, err := t.ethClient.GetNonce(t.address)
	if err != nil {
		logger.Error("get ethNonce error:%v", err)
		return 0, err
	}
	if t.ethNonce >= chainNonce {
		return t.ethNonce + 1, nil
	}
	return chainNonce, nil
}

func (t *TaskManager) GetOasisNewNonce() (uint64, error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	chainNonce, err := t.ethClient.GetNonce(t.address)
	if err != nil {
		logger.Error("get ethNonce error:%v", err)
		return 0, err
	}
	if t.ethNonce >= chainNonce {
		return t.ethNonce + 1, nil
	}
	return chainNonce, nil
}

func (t *TaskManager) execute() error {
	t.queue.Range(func(key, value interface{}) bool {
		task, ok := value.(*Task)
		if !ok {
			logger.Error("never should happen innerTask type: %v", value)
			return false
		}
		switch task.Type {
		case DepositTask:
			_, err := t.submitDepositTx(task)
			if err != nil {
				logger.Error(err.Error())
			}
		case UpdateTask:
			_, err := t.submitUpdateTx(task)
			if err != nil {
				logger.Error(err.Error())
			}
		case RedeemTask:
			_, err := t.submitOasisTx(task)
			if err != nil {
				logger.Error(err.Error())
			}
		default:
			logger.Error("never should happen network: %v", task)
			return false
		}
		return true

	})
	return nil
}

func (t *TaskManager) MintZkBtcRequest(proofId, proof string) (string, error) {
	// todo
	btcTx, err := t.btcClient.GetTransaction(proofId)
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}
	param, err := parseBtcTx(btcTx)
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}
	param.Proof = proof
	nonce, err := t.GetEthNewNonce()
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}
	task := NewDepositTask(nonce, param)
	t.addTask(task)
	txHash, err := t.submitDepositTx(task)
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}
	return txHash, nil
}

func (t *TaskManager) UpdateUtxoRequest(txIds []string, proof string) (string, error) {
	// todo
	nonce, err := t.GetEthNewNonce()
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}
	param := UpdateParam{
		Proof: proof,
		TxIds: txIds,
	}
	task := NewUpdateTask(nonce, param)
	t.addTask(task)
	txHash, err := t.submitUpdateTx(task)
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}
	return txHash, nil
}

func (t *TaskManager) RedeemBtcRequest(proof string) (string, error) {
	var innerTasks []*innerTask
	for index := 0; index < 3; index++ {
		newNonce, err := t.GetOasisNewNonce()
		if err != nil {
			logger.Error(err.Error())
			return "", err
		}
		task := NewRedeemInnerTask(newNonce, proof)
		innerTasks = append(innerTasks, task)
	}
	task := &Task{
		Id:    UUID(),
		Type:  RedeemTask,
		tasks: innerTasks,
	}
	t.addTask(task)
	txHash, err := t.submitOasisTx(task)
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}
	return txHash, nil
}

func (t *TaskManager) submitOasisTx(oasisTask *Task) (string, error) {
	//todo
	if oasisTask.Type != RedeemTask {
		return "", fmt.Errorf("never should happen network: %v", oasisTask)
	}
	for _, task := range oasisTask.tasks {
		if task.TxHash != "" {
			_, err := t.doOasisTx(task)
			if err != nil {
				logger.Error(err.Error())
				return "", err
			}
		}
	}
	return "", nil

}

func (t *TaskManager) submitUpdateTx(ethTask *Task, highPriority ...bool) (string, error) {
	highPrio := false
	if len(highPriority) > 0 {
		highPrio = highPriority[0]
	}
	if len(ethTask.tasks) != 1 {
		return "", fmt.Errorf("never should happen innerTask length: %v", len(ethTask.tasks))
	}
	task := ethTask.tasks[0]
	if task.TxHash != "" {
		transaction, err := t.ethClient.TransactionReceipt(context.Background(), common.HexToHash(task.TxHash))
		if err != nil {
			logger.Error(err.Error())
			return "", err
		}
		if transaction.Status == types.ReceiptStatusSuccessful {
			t.RemoveTask(ethTask)
			return "", nil
		}
	}
	// todo only focus on pending
	if task.Status == Pending {
		currentTime := time.Now()
		if currentTime.Sub(task.StartTime) < t.timeout {
			return "", nil
		}
	}
	gasPrice, err := t.ethClient.GetGasPrice()
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}
	chainId, err := t.ethClient.GetChainId()
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}
	if highPrio {
		gasPrice = big.NewInt(0).Add(gasPrice, big.NewInt(2))
	}
	param, ok := task.data.(UpdateParam)
	if !ok {
		logger.Error("never should happen innerTask type: %v", task)
		return "", fmt.Errorf("never should happen innerTask type: %v", task)
	}
	txHash, err := t.ethClient.UpdateUtxoChange(t.keyStore.GetPrivateKey(), param.TxIds,
		task.Nonce, 0, chainId, gasPrice, []byte(param.Proof))
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}
	logger.Info("submit deposit tx: %v", txHash)
	task.Status = Pending
	task.StartTime = time.Now()
	return txHash, nil

}

func (t *TaskManager) submitDepositTx(ethTask *Task, highPriority ...bool) (string, error) {
	highPrio := false
	if len(highPriority) > 0 {
		highPrio = highPriority[0]
	}
	if len(ethTask.tasks) != 1 {
		return "", fmt.Errorf("never should happen innerTask length: %v", len(ethTask.tasks))
	}
	task := ethTask.tasks[0]
	if task.TxHash != "" {
		transaction, err := t.ethClient.TransactionReceipt(context.Background(), common.HexToHash(task.TxHash))
		if err != nil {
			logger.Error(err.Error())
			return "", err
		}
		if transaction.Status == types.ReceiptStatusSuccessful {
			t.RemoveTask(ethTask)
			return "", nil
		}
	}
	// todo only focus on pending
	if task.Status == Pending {
		currentTime := time.Now()
		if currentTime.Sub(task.StartTime) < t.timeout {
			return "", nil
		}
	}
	gasPrice, err := t.ethClient.GetGasPrice()
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}
	chainId, err := t.ethClient.GetChainId()
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}
	if highPrio {
		gasPrice = big.NewInt(0).Add(gasPrice, big.NewInt(2))
	}
	param, ok := task.data.(DepositParam)
	if !ok {
		logger.Error("never should happen innerTask type: %v", task)
		return "", fmt.Errorf("never should happen innerTask type: %v", task)
	}
	txHash, err := t.ethClient.Deposit(t.keyStore.GetPrivateKey(), param.TxId, param.EthAddr, param.TxIndex,
		task.Nonce, 0, chainId, gasPrice, big.NewInt(param.Amount), []byte(param.Proof))
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}
	task.Status = Pending
	task.StartTime = time.Now()
	logger.Info("submit deposit tx: %v", txHash)
	return txHash, nil
}

func (t *TaskManager) doOasisTx(task *innerTask) (string, error) {
	// todo
	param, ok := task.data.(RedeemParam)
	if !ok {
		logger.Error("never should happen innerTask type: %v", task)
		return "", fmt.Errorf("never should happen innerTask type: %v", task)
	}

	if task.TxHash != "" {

	}
	txHash, err := t.oasisClient.Redeem(param.Proof)
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}
	logger.Info("submit redeem tx: %v", txHash)
	return "", nil
}

func (t *TaskManager) deposit(nonce uint64, gasPrice *big.Int) (string, error) {
	panic(t)
}

func (t *TaskManager) verify(nonce uint64, gasPrice *big.Int) (string, error) {
	panic(t)
}

func (t *TaskManager) SubmitDfinityTx(task *Task) error {
	panic(task)
}

func (t *TaskManager) RemoveTask(task *Task) {
	t.queue.Delete(task.Id)
}

func (t *TaskManager) addTask(task *Task) {
	t.queue.Store(task.Id, task)
}

func (t *TaskManager) Execute() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-t.exit:
			logger.Info("innerTask manager goroutine exit now ...")
			return
		case <-ticker.C:
			err := t.execute()
			if err != nil {
				logger.Error(err.Error())
			}
		}
	}
}

func (t *TaskManager) Exit() {
	t.exit <- struct{}{}
}

type Task struct {
	Id    string
	Type  TaskType
	tasks []*innerTask
}

type innerTask struct {
	Id        string
	Nonce     uint64
	TxHash    string
	StartTime time.Time
	EndTime   time.Time
	Status    TaskStatus
	Network   Network
	data      interface{}
}

type DepositParam struct {
	EthAddr string
	TxId    string
	TxIndex uint32
	Amount  int64
	Proof   string
}

type UpdateParam struct {
	Proof string
	TxIds []string
}

type RedeemParam struct {
	Proof string
}

type TaskType int

const (
	DepositTask TaskType = iota + 1
	UpdateTask
	RedeemTask
)

type TaskStatus int

const (
	Default TaskStatus = iota
	Success
	Failed
	Pending
)

type Network int

const (
	EthereumChain Network = iota + 1
	OasisChain
	DfinityChain
)

func parseBtcTx(tx btcTypes.RawTransaction) (DepositParam, error) {
	panic(tx)
}

func NewDepositTask(nonce uint64, data interface{}) *Task {
	return &Task{
		Type: DepositTask,
		tasks: []*innerTask{
			{
				Id:      UUID(),
				Nonce:   nonce,
				Status:  Default,
				data:    data,
				Network: EthereumChain,
			},
		},
	}
}

func NewUpdateTask(nonce uint64, data interface{}) *Task {
	return &Task{
		Type: UpdateTask,
		tasks: []*innerTask{
			{
				Id:      UUID(),
				Nonce:   nonce,
				Status:  Default,
				data:    data,
				Network: EthereumChain,
			},
		},
	}
}

func NewRedeemInnerTask(nonce uint64, data interface{}) *innerTask {
	return &innerTask{
		Id:      UUID(),
		Nonce:   nonce,
		Status:  Default,
		data:    data,
		Network: OasisChain,
	}

}
