package node

import (
	"container/heap"
	"container/list"
	"fmt"
	"github.com/lightec-xyz/daemon/common"
	"sort"
	"sync"
)

//Todo abstract queue

type HeapQueue struct {
	queue *PriorityQueue
	lock  sync.Mutex
}

func NewHeapQueue() *HeapQueue {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	return &HeapQueue{
		queue: &pq,
	}
}

func (h *HeapQueue) Push(value *common.ZkProofRequest) {
	h.lock.Lock()
	defer h.lock.Unlock()
	heap.Push(h.queue, &Item{
		value:    value,
		priority: int(value.Weight),
	})
}

func (h *HeapQueue) Len() int {
	h.lock.Lock()
	defer h.lock.Unlock()
	return h.queue.Len()
}

func (h *HeapQueue) Pop() (*common.ZkProofRequest, bool) {
	h.lock.Lock()
	defer h.lock.Unlock()
	if h.queue.Len() == 0 {
		return nil, false
	}
	return heap.Pop(h.queue).(*Item).value, true
}

// https://pkg.go.dev/container/heap
type Item struct {
	value    *common.ZkProofRequest
	priority int
	index    int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	if pq[i].priority != pq[j].priority {
		return pq[i].priority > pq[j].priority
	}
	return pq[i].value.Index < pq[j].value.Index
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func (pq *PriorityQueue) update(item *Item, value *common.ZkProofRequest, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}

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

func (sl *Queue) Iterator(fn func(value *list.Element) error) {
	for element := sl.list.Front(); element != nil; element = element.Next() {
		err := fn(element)
		if err != nil {
			return
		}
	}
}

// PendingQueue
type PendingQueue struct {
	list *sync.Map
}

func NewPendingQueue() *PendingQueue {
	return &PendingQueue{
		list: new(sync.Map),
	}
}

func (q *PendingQueue) Push(key string, value *common.ZkProofRequest) {
	q.list.Store(key, value)
}

func (q *PendingQueue) Delete(key string) {
	q.list.Delete(key)
}

func (q *PendingQueue) Get(key string) (*common.ZkProofRequest, error) {
	value, ok := q.list.Load(key)
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	req, ok := value.(*common.ZkProofRequest)
	if !ok {
		return nil, fmt.Errorf("parse error")
	}
	return req, nil
}

func (q *PendingQueue) Iterator(fn func(value *common.ZkProofRequest) error) {
	q.list.Range(func(key, value interface{}) bool {
		req, ok := value.(*common.ZkProofRequest)
		if !ok {
			return false
		}
		err := fn(req)
		if err != nil {
			return false
		}
		return true
	})
}

type ArrayQueue struct {
	list []*common.ZkProofRequest
	lock sync.Mutex
}

func NewArrayQueue() *ArrayQueue {
	return &ArrayQueue{
		list: make([]*common.ZkProofRequest, 0),
	}
}

func (aq *ArrayQueue) Push(value *common.ZkProofRequest) {
	aq.lock.Lock()
	defer aq.lock.Unlock()
	aq.list = append(aq.list, value)
}
func (aq *ArrayQueue) Pop() (*common.ZkProofRequest, bool) {
	aq.lock.Lock()
	defer aq.lock.Unlock()
	if len(aq.list) == 0 {
		return nil, false
	}
	aq.sortList()
	value := aq.list[0]
	aq.list = aq.list[1:]
	return value, true
}

func (aq *ArrayQueue) sortList() {
	sort.Slice(aq.list, func(i, j int) bool {
		if aq.list[i].Weight != aq.list[j].Weight {
			return aq.list[i].Weight > aq.list[j].Weight
		}
		return aq.list[i].Index < aq.list[j].Index
	})
}

func (aq *ArrayQueue) Iterator(fn func(index int, value *common.ZkProofRequest) error) {
	// todo
	aq.lock.Lock()
	defer aq.lock.Unlock()
	for index, value := range aq.list {
		err := fn(index, value)
		if err != nil {
			return
		}
	}
}
func (aq *ArrayQueue) Len() int {
	aq.lock.Lock()
	defer aq.lock.Unlock()
	return len(aq.list)
}

// ProofRespQueue todo
type SubmitQueue struct {
	list *sync.Map
}

func NewSubmitQueue() *SubmitQueue {
	return &SubmitQueue{
		list: new(sync.Map),
	}
}

func (q *SubmitQueue) Push(value *common.ZkProofResponse) {
	q.list.Store(value.Id(), value)
}

func (q *SubmitQueue) Delete(key string) {
	q.list.Delete(key)
}

func (q *SubmitQueue) Get(key string) (*common.ZkProofResponse, error) {
	value, ok := q.list.Load(key)
	if !ok {
		return nil, fmt.Errorf("not found: %v", key)
	}
	req, ok := value.(*common.ZkProofResponse)
	if !ok {
		return nil, fmt.Errorf("parse error:%v", key)
	}
	return req, nil
}

func (q *SubmitQueue) Iterator(fn func(value *common.ZkProofResponse) error) {
	q.list.Range(func(key, value interface{}) bool {
		req, ok := value.(*common.ZkProofResponse)
		if !ok {
			return false
		}
		err := fn(req)
		if err != nil {
			return false
		}
		return true
	})
}

// ProofRespQueue todo
type ProofRespQueue struct {
	list *sync.Map
}

func NewProofRespQueue() *ProofRespQueue {
	return &ProofRespQueue{
		list: new(sync.Map),
	}
}

func (q *ProofRespQueue) Push(value *common.SubmitProof) {
	q.list.Store(value.Id, value)
}

func (q *ProofRespQueue) Delete(key string) {
	q.list.Delete(key)
}

func (q *ProofRespQueue) Get(key string) (*common.SubmitProof, error) {
	value, ok := q.list.Load(key)
	if !ok {
		return nil, fmt.Errorf("not found: %v", key)
	}
	req, ok := value.(*common.SubmitProof)
	if !ok {
		return nil, fmt.Errorf("parse error:%v", key)
	}
	return req, nil
}

func (q *ProofRespQueue) Iterator(fn func(value *common.SubmitProof) error) {
	q.list.Range(func(key, value interface{}) bool {
		req, ok := value.(*common.SubmitProof)
		if !ok {
			return false
		}
		err := fn(req)
		if err != nil {
			return false
		}
		return true
	})
}
