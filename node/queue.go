package node

import (
	"fmt"
	"github.com/lightec-xyz/daemon/common"
	"sort"
	"sync"
)

// todo
type QueueManager struct {
}

func NewQueueManager() *QueueManager {
	return &QueueManager{}
}

type PendingQueue struct {
	list *sync.Map
}

func NewPendingQueue() *PendingQueue {
	return &PendingQueue{
		list: new(sync.Map),
	}
}

func (q *PendingQueue) Add(key string, value *common.ProofRequest) {
	q.list.Store(key, value)
}

func (q *PendingQueue) Delete(key string) {
	q.list.Delete(key)
}

func (q *PendingQueue) Get(key string) (*common.ProofRequest, error) {
	value, ok := q.list.Load(key)
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	req, ok := value.(*common.ProofRequest)
	if !ok {
		return nil, fmt.Errorf("parse error")
	}
	return req, nil
}

func (q *PendingQueue) Check(key string) bool {
	_, ok := q.list.Load(key)
	return ok
}

func (q *PendingQueue) Iterator(fn func(value *common.ProofRequest) error) {
	q.list.Range(func(key, value interface{}) bool {
		req, ok := value.(*common.ProofRequest)
		if !ok {
			return true
		}
		err := fn(req)
		if err != nil {
			return true //continue to check next item
		}
		return true
	})
}

func sortRequest(a, b *common.ProofRequest) bool {
	if a.Weight == b.Weight {
		if a.ProofType == b.ProofType { // todo more rule
			if a.ProofType == common.SyncComInnerType {
				if a.Prefix == b.Prefix {
					return a.FIndex < b.FIndex
				}
				return a.Prefix < b.Prefix
			}
			return a.FIndex < b.FIndex
		}
		return false
	}
	return a.Weight > b.Weight
}

// todo
type ArrayQueue struct {
	list   []*common.ProofRequest
	lock   sync.Mutex
	sortFn func(a, b *common.ProofRequest) bool
}

func NewArrayQueue(sortFn func(a, b *common.ProofRequest) bool) *ArrayQueue {
	return &ArrayQueue{
		list:   make([]*common.ProofRequest, 0),
		sortFn: sortFn,
	}
}

func (q *ArrayQueue) AddFirst(value *common.ProofRequest) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.list = append([]*common.ProofRequest{value}, q.list...)
}
func (q *ArrayQueue) Push(value *common.ProofRequest) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.list = append(q.list, value)
	q.sortList()
}
func (q *ArrayQueue) Pop() (*common.ProofRequest, bool) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if len(q.list) == 0 {
		return nil, false
	}
	value := q.list[0]
	q.list = q.list[1:]
	return value, true
}

func (q *ArrayQueue) PopFn(fn func(req *common.ProofRequest) bool) (*common.ProofRequest, bool) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if len(q.list) == 0 {
		return nil, false
	}
	var newList []*common.ProofRequest
	for index, value := range q.list {
		ok := fn(value)
		if ok {
			newList = append(newList, q.list[0:index]...)
			if index+1 < len(q.list) {
				newList = append(newList, q.list[index+1:]...)
			}
			q.list = newList
			return value, true
		}
	}
	return nil, false

}

func (q *ArrayQueue) sortList() {
	sort.SliceStable(q.list, func(i, j int) bool {
		if q.sortFn != nil {
			return q.sortFn(q.list[i], q.list[j])
		}
		//todo
		if q.list[i].Weight == q.list[j].Weight {
			if q.list[i].ProofType == q.list[j].ProofType {
				return q.list[i].FIndex < q.list[j].FIndex
			}
			return false
		}
		return q.list[i].Weight > q.list[j].Weight
	})
}

func (q *ArrayQueue) Iterator(fn func(index int, value *common.ProofRequest) error) {
	// todo
	q.lock.Lock()
	defer q.lock.Unlock()
	for index, value := range q.list {
		err := fn(index, value)
		if err != nil {
			return
		}
	}
}

func (q *ArrayQueue) Filter(fn func(value *common.ProofRequest) (bool, error)) error {
	q.lock.Lock()
	defer q.lock.Unlock()
	var newList []*common.ProofRequest
	for _, value := range q.list {
		ok, err := fn(value)
		if err != nil {
			return err
		}
		if ok {
			newList = append(newList, value)
		}
	}
	q.list = newList
	return nil
}

func (q *ArrayQueue) Len() int {
	q.lock.Lock()
	defer q.lock.Unlock()
	return len(q.list)
}

func (q *ArrayQueue) Remove(fn func(value *common.ProofRequest) bool) []*common.ProofRequest {
	q.lock.Lock()
	defer q.lock.Unlock()
	var expired []*common.ProofRequest
	var newList []*common.ProofRequest
	for _, value := range q.list {
		ok := fn(value)
		if ok {
			expired = append(expired, value)
		} else {
			newList = append(newList, value)
		}
	}
	q.list = newList
	return expired
}

func (q *ArrayQueue) List() []*common.ProofRequest {
	return q.list
}

type SubmitQueue struct {
	list *sync.Map
}

func NewSubmitQueue() *SubmitQueue {
	return &SubmitQueue{
		list: new(sync.Map),
	}
}

func (q *SubmitQueue) Push(value *common.ProofResponse) {
	q.list.Store(value.ProofId(), value)
}

func (q *SubmitQueue) Delete(key string) {
	q.list.Delete(key)
}

func (q *SubmitQueue) Get(key string) (*common.ProofResponse, error) {
	value, ok := q.list.Load(key)
	if !ok {
		return nil, fmt.Errorf("not found: %v", key)
	}
	req, ok := value.(*common.ProofResponse)
	if !ok {
		return nil, fmt.Errorf("parse error:%v", key)
	}
	return req, nil
}

func (q *SubmitQueue) Iterator(fn func(value *common.ProofResponse) error) {
	q.list.Range(func(key, value interface{}) bool {
		req, ok := value.(*common.ProofResponse)
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
