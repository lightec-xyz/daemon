package node

import (
	"container/list"
	"sync"
)

type Queue struct {
	list     *list.List
	lock     sync.Mutex
	capacity uint64
}

func NewQueue() *Queue {
	return &Queue{
		list: list.New(),
	}
}

func NewQueueWithCapacity(capacity uint64) *Queue {
	return &Queue{
		list:     list.New(),
		lock:     sync.Mutex{},
		capacity: capacity,
	}
}

func (sl *Queue) CanPush() bool {
	if sl.capacity == 0 {
		return true
	}
	sl.lock.Lock()
	defer sl.lock.Unlock()
	return sl.list.Len() < int(sl.capacity)
}

func (sl *Queue) PushBack(value interface{}) {
	sl.lock.Lock()
	defer sl.lock.Unlock()
	sl.list.PushBack(value)
}

func (sl *Queue) PushFront(value interface{}) {
	sl.lock.Lock()
	defer sl.lock.Unlock()
	sl.list.PushFront(value)
}

func (sl *Queue) Front() *list.Element {
	sl.lock.Lock()
	defer sl.lock.Unlock()
	return sl.list.Front()

}
func (sl *Queue) Back() *list.Element {
	sl.lock.Lock()
	defer sl.lock.Unlock()
	return sl.list.Back()

}

func (sl *Queue) Len() int {
	sl.lock.Lock()
	defer sl.lock.Unlock()
	return sl.list.Len()
}
func (sl *Queue) Remove(e *list.Element) {
	sl.lock.Lock()
	defer sl.lock.Unlock()
	sl.list.Remove(e)
}
