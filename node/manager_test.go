package node

import (
	"testing"
)

func TestSafeList(t *testing.T) {
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
		queue.Remove(front)
	}

}
