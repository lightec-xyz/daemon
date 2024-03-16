package node

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc/ethereum"
	"math/big"
	"sync"
	"time"
)

type TaskManager struct {
	nonce     uint64
	ethClient *ethereum.Client
	address   string
	lock      sync.Mutex
	exit      chan struct{}
	queue     *sync.Map
	keyStore  *KeyStore
	timeout   time.Duration
}

func NewTaskManager(address, privateKey string, ethClient ethereum.Client) (*TaskManager, error) {
	nonce, err := ethClient.GetNonce(address)
	if err != nil {
		logger.Error("get nonce error:%v", err)
		return nil, err
	}
	return &TaskManager{
		queue:     new(sync.Map),
		nonce:     nonce,
		address:   address,
		ethClient: &ethClient,
		keyStore:  NewKeyStore(privateKey),
		exit:      make(chan struct{}, 1),
		timeout:   5 * time.Minute,
	}, nil
}

func (t *TaskManager) GetNewNonce() (uint64, error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	chainNonce, err := t.ethClient.GetNonce(t.address)
	if err != nil {
		logger.Error("get nonce error:%v", err)
		return 0, err
	}
	if t.nonce >= chainNonce {
		return t.nonce + 1, nil
	}
	return chainNonce, nil
}

func (t *TaskManager) execute() error {
	t.queue.Range(func(key, value interface{}) bool {
		task, ok := value.(*Task)
		if !ok {
			logger.Error("never should happen task type: %v", value)
			return false
		}
		switch task.Status {
		case Default:
			err := t.SendTask(task)
			if err != nil {
				logger.Error(err.Error())
			}
		case Pending:
			// too long pending tx,retry again
			currentTime := time.Now()
			if currentTime.Sub(task.StartTime) >= t.timeout {
				err := t.SendTask(task, true)
				if err != nil {
					logger.Error(err.Error())
				}
			}
		case Success:
			t.RemoveTask(task)
		case Failed:
			err := t.SendTask(task)
			if err != nil {
				logger.Error(err.Error())
			}
		default:
			logger.Error("never should happen task status: %v", task.Status)
			return false
		}
		return true

	})
	return nil
}

func (t *TaskManager) Submit() error {
	task := &Task{}
	err := t.SendTask(task)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	t.AddTask(task)
	return nil
}

func (t *TaskManager) AddTask(task *Task) {
	t.queue.Store(task.Id, task)
}

func (t *TaskManager) SendTask(task *Task, highPriority ...bool) error {
	highPrio := false
	if len(highPriority) > 0 {
		highPrio = highPriority[0]
	}
	if task.TxHash != "" {
		transaction, err := t.ethClient.TransactionReceipt(context.Background(), common.HexToHash(task.TxHash))
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
	switch task.Type {
	case DepositTask:

	case RedeemTask:

	case VerifyTask:

	default:
		logger.Error("never should happen task type: %v", task.Type)
		return fmt.Errorf("never should happen task type: %v", task.Type)
	}
	return nil

}

func (t *TaskManager) RemoveTask(task *Task) {
	t.queue.Delete(task.Id)
}

func (t *TaskManager) Execute() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-t.exit:
			logger.Info("task manager goroutine exit now ...")
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
	Nonce     uint64
	Id        string
	TxHash    string
	StartTime time.Time
	EndTime   time.Time
	Status    TaskStatus
	Type      TaskType
}

func NewTask(nonce uint64) (*Task, error) {
	return &Task{
		Nonce:     nonce,
		Id:        UUID(),
		StartTime: time.Now(),
		EndTime:   time.Now(),
		Status:    Default,
	}, nil
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
