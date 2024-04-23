package node

import (
	"github.com/lightec-xyz/daemon/common"
	"testing"
)

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
	arrayQueue.Push(&common.ZkProofRequest{Weight: 10})
	arrayQueue.Push(&common.ZkProofRequest{Weight: 5})
	arrayQueue.Push(&common.ZkProofRequest{Weight: 1})
	arrayQueue.Push(&common.ZkProofRequest{Weight: 7})
	arrayQueue.Push(&common.ZkProofRequest{Weight: 8})
	t.Logf("lenghth: %v \n", arrayQueue.Len())
	request, ok := arrayQueue.Pop()
	if !ok {
		t.Fatal(err)
	}
	t.Logf("lenghth: %v \n", arrayQueue.Len())
	t.Logf("pop result: %v \n", request.Weight)
	arrayQueue.Iterator(func(index int, value *common.ZkProofRequest) error {
		t.Log(value.Weight)
		return nil
	})
	t.Logf("lenghth: %v \n", arrayQueue.Len())
}
