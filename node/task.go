package node

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc/bitcoin"
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
			err := t.submitEthTx(task)
			if err != nil {
				logger.Error(err.Error())
			}
		case VerifyTask:
			err := t.submitEthTx(task)
			if err != nil {
				logger.Error(err.Error())
			}
		case RedeemTask:
			err := t.SubmitDFinityTx(task)
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

func (t *TaskManager) DepositRequest() error {
	nonce, err := t.GetEthNewNonce()
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	task := NewDepositTask(nonce)
	err = t.submitEthTx(task)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}

func (t *TaskManager) VerifyRequest() error {
	nonce, err := t.GetEthNewNonce()
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	task := NewVerifyTask(nonce)
	err = t.submitEthTx(task)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}

func (t *TaskManager) SubmitOasisTx(task *Task, highPriority ...bool) error {
	panic(t)
}

func (t *TaskManager) submitEthTx(task *Task, highPriority ...bool) error {
	highPrio := false
	if len(highPriority) > 0 {
		highPrio = highPriority[0]
	}
	if len(task.tasks) != 1 {
		return fmt.Errorf("never should happen innerTask length: %v", len(task.tasks))
	}
	depositTask := task.tasks[0]
	switch depositTask.Status {
	case Default:
	case Pending:
		currentTime := time.Now()
		if currentTime.Sub(depositTask.StartTime) <= t.timeout {
			return nil
		} else {
			highPrio = true
		}
	case Success:
	case Failed:
	default:
		return fmt.Errorf("never should happen innerTask status: %v", depositTask.Status)
	}

	if depositTask.TxHash != "" {
		transaction, err := t.ethClient.TransactionReceipt(context.Background(), common.HexToHash(depositTask.TxHash))
		if err != nil {
			logger.Error(err.Error())
			return err
		}
		if transaction.Status == types.ReceiptStatusSuccessful {
			t.RemoveTask(task)
			return nil
		}
	}
	gasPrice, err := t.ethClient.GetGasPrice()
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	if highPrio {
		gasPrice = big.NewInt(0).Add(gasPrice, big.NewInt(2))
	}
	// todo submit tx
	switch task.Type {
	case DepositTask:

	case VerifyTask:

	default:
		logger.Error("never should happen network: %v", task)
		return fmt.Errorf("never should happen network: %v", task)
	}
	depositTask.Status = Pending
	return nil
}

func (t *TaskManager) SubmitDFinityTx(task *Task) error {
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
	data  []string
}

func NewDepositTask(nonce uint64) *Task {
	return &Task{
		Type: DepositTask,
		tasks: []*innerTask{
			{
				Id:     UUID(),
				Nonce:  nonce,
				Status: Default,
			},
		},
	}
}

func NewVerifyTask(nonce uint64) *Task {
	return &Task{
		Type: VerifyTask,
		tasks: []*innerTask{
			{
				Id:     UUID(),
				Nonce:  nonce,
				Status: Default,
			},
		},
	}
}

func NewRedeemTask() *Task {
	return &Task{
		Type: RedeemTask,
		tasks: []*innerTask{
			{
				Id:     UUID(),
				Nonce:  0,
				Status: Default,
			},
		},
	}
}

type innerTask struct {
	Id        string
	Nonce     uint64
	TxHash    string
	StartTime time.Time
	EndTime   time.Time
	Status    TaskStatus
	Network   Network
}

type TaskType int

const (
	DepositTask TaskType = iota + 1
	VerifyTask
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
	DFinityChain
)
