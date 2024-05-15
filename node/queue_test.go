package node

import (
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
		t.Log(request.Weight, request.Index)
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
	arrayQueue.Push(&common.ZkProofRequest{Weight: 10, Index: 6})
	arrayQueue.Push(&common.ZkProofRequest{Weight: 10, Index: 4})
	arrayQueue.Push(&common.ZkProofRequest{Weight: 10, Index: 9})
	arrayQueue.Push(&common.ZkProofRequest{Weight: 5, Index: 10})
	arrayQueue.Push(&common.ZkProofRequest{Weight: 5, Index: 5})
	arrayQueue.Push(&common.ZkProofRequest{Weight: 5, Index: 7})
	arrayQueue.Push(&common.ZkProofRequest{Weight: 1, Index: 7})
	arrayQueue.Push(&common.ZkProofRequest{Weight: 7, Index: 88})
	arrayQueue.Push(&common.ZkProofRequest{Weight: 7, Index: 10})
	arrayQueue.Push(&common.ZkProofRequest{Weight: 7, Index: 8})
	arrayQueue.Pop()
	arrayQueue.Iterator(func(index int, value *common.ZkProofRequest) error {
		t.Log(value.Weight, value.Index)
		return nil
	})
	t.Logf("lenghth: %v \n", arrayQueue.Len())
}
