package node

import (
	"fmt"
	"github.com/lightec-xyz/daemon/common"
	"sort"
	"sync"
)

type QueueManager struct {
	requests *ArrayQueue
	pending  *PendingQueue
	cache    *cache
}

func NewQueueManager() *QueueManager {
	return &QueueManager{
		requests: NewArrayQueue(sortRequest),
		pending:  NewPendingQueue(),
		cache:    NewCacheState(),
	}
}

func (q *QueueManager) CheckId(proofId string) bool {
	return q.cache.Check(proofId)
}

func (q *QueueManager) StoreId(proofId string) {
	q.cache.Store(proofId, nil)
}

func (q *QueueManager) DeleteId(proofId string) {
	q.cache.Delete(proofId)
}

func (q *QueueManager) PushRequest(req *common.ProofRequest) {
	q.requests.Push(req)
}

func (q *QueueManager) PopRequest() (*common.ProofRequest, bool) {
	return q.requests.Pop()
}

func (q *QueueManager) PopFnRequest(fn func(req *common.ProofRequest) bool) (*common.ProofRequest, bool) {
	return q.requests.PopFn(fn)
}

func (q *QueueManager) FilterRequest(fn func(value *common.ProofRequest) bool) []*common.ProofRequest {
	return q.requests.Filter(fn)
}
func (q *QueueManager) RequestLen() int {
	return q.requests.Len()
}

func (q *QueueManager) RemoveRequest(fn func(value *common.ProofRequest) bool) []*common.ProofRequest {
	return q.requests.Remove(fn)
}

func (q *QueueManager) ListRequest() []*common.ProofRequest {
	return q.requests.List()
}

func (q *QueueManager) AddPending(key string, value *common.ProofRequest) {
	q.pending.Add(key, value)
}

func (q *QueueManager) DeletePending(key string) {
	q.pending.Delete(key)
}

func (q *QueueManager) FilterPending(fn func(value *common.ProofRequest) bool) []*common.ProofRequest {
	return q.pending.Iterator(fn)
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

func (q *PendingQueue) Iterator(fn func(value *common.ProofRequest) bool) []*common.ProofRequest {
	var filters []*common.ProofRequest
	q.list.Range(func(key, value interface{}) bool {
		req, ok := value.(*common.ProofRequest)
		if !ok {
			return true
		}
		if fn != nil {
			match := fn(req)
			if match {
				filters = append(filters, req)
			}
		}

		return true
	})
	return filters
}

func sortRequest(a, b *common.ProofRequest) bool {
	if a.BlockTime != 0 && b.BlockTime != 0 {
		return a.BlockTime < b.BlockTime
	}
	if a.Weight == b.Weight {
		if a.ProofType == b.ProofType {
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

func (q *ArrayQueue) Filter(fn func(value *common.ProofRequest) bool) []*common.ProofRequest {
	q.lock.Lock()
	defer q.lock.Unlock()
	var filtersReq []*common.ProofRequest
	for _, value := range q.list {
		if fn == nil {
			return filtersReq
		}
		ok := fn(value)
		if ok {
			filtersReq = append(filtersReq, value)
		}
	}
	return filtersReq
}

func (q *ArrayQueue) Len() int {
	q.lock.Lock()
	defer q.lock.Unlock()
	return len(q.list)
}

func (q *ArrayQueue) Remove(fn func(value *common.ProofRequest) bool) []*common.ProofRequest {
	q.lock.Lock()
	defer q.lock.Unlock()
	var filters []*common.ProofRequest
	var newList []*common.ProofRequest
	for _, value := range q.list {
		if fn != nil {
			return filters
		}
		ok := fn(value)
		if ok {
			filters = append(filters, value)
		} else {
			newList = append(newList, value)
		}
	}
	q.list = newList
	return filters
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
