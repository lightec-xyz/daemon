package node

import (
	"fmt"
	"github.com/lightec-xyz/daemon/common"
	"testing"
)

func TestHeapQueue(t *testing.T) {
	heapQueue := NewHeapQueue()
	heapQueue.Push(&common.ZkProofRequest{Weight: 6, Index: 10})
	heapQueue.Push(&common.ZkProofRequest{Weight: 7, Index: 4})
	heapQueue.Push(&common.ZkProofRequest{Weight: 10, Index: 9})
	heapQueue.Push(&common.ZkProofRequest{Weight: 5, Index: 6})
	heapQueue.Push(&common.ZkProofRequest{Weight: 1, Index: 10})
	heapQueue.Push(&common.ZkProofRequest{Weight: 1, Index: 4})
	heapQueue.Push(&common.ZkProofRequest{Weight: 1, Index: 9})
	heapQueue.Push(&common.ZkProofRequest{Weight: 1, Index: 6})
	for heapQueue.Len() > 0 {
		request, _ := heapQueue.Pop()
		t.Log(request.Weight, request.Index, heapQueue.Len())
	}
}

func TestQueue(t *testing.T) {
	queue := NewQueue()
	queue.PushFront(1)
	queue.PushFront(2)
	queue.PushFront(3)
	queue.PushFront(4)
	queue.PushFront(5)
	queue.PushBack(9)
	//queue.PushBack(6)
	for i := 0; i < 5; i++ {
		front := queue.Back()
		t.Log(front.Value)
		//queue.Remove(front)
	}
}

func TestArrayQueue(t *testing.T) {
	arrayQueue := NewArrayQueue()
	arrayQueue.Push(common.NewZkProofRequest(common.TxInEth2, nil, 0, 0, "1"))
	arrayQueue.Push(common.NewZkProofRequest(common.BeaconHeaderType, nil, 100, 0, ""))
	arrayQueue.Push(common.NewZkProofRequest(common.BeaconHeaderFinalityType, nil, 110, 0, ""))
	arrayQueue.Push(common.NewZkProofRequest(common.TxInEth2, nil, 0, 0, "2"))
	arrayQueue.Push(common.NewZkProofRequest(common.BeaconHeaderType, nil, 200, 0, ""))
	arrayQueue.Push(common.NewZkProofRequest(common.TxInEth2, nil, 0, 0, "3"))
	arrayQueue.Push(common.NewZkProofRequest(common.TxInEth2, nil, 0, 0, "4"))
	arrayQueue.Push(common.NewZkProofRequest(common.BeaconHeaderType, nil, 300, 0, ""))
	arrayQueue.Push(common.NewZkProofRequest(common.TxInEth2, nil, 0, 0, "5"))
	arrayQueue.Push(common.NewZkProofRequest(common.BeaconHeaderType, nil, 400, 0, ""))

	arrayQueue.Iterator(func(index int, value *common.ZkProofRequest) error {
		t.Log(value.RequestId())
		return nil
	})
	fmt.Printf("*********** lenghth: %v  **************** \n", arrayQueue.Len())
	for arrayQueue.Len() > 0 {
		request, ok := arrayQueue.PopFn(func(req *common.ZkProofRequest) bool {
			if req.ProofType == common.BeaconHeaderType {
				return true
			}
			return false

		})
		if ok {
			t.Log(request.RequestId())
		}

	}
}
