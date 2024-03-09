package node

import (
	"fmt"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc/beacon"
	"sync"
	"sync/atomic"
	"time"
)

const (
	MaxReqNums   = 2
	MaxQueueSize = 2
)

type BeaconFetch struct {
	beaconClient       *beacon.Client
	currentReqNums     *atomic.Int64
	fetchRequest       chan FetchRequest
	fileStore          *FileStore
	exit               chan struct{}
	fetchProofResponse chan FetchDataResponse
	genesisSyncPeriod  uint64
	fetchQueue         *Queue
	cache              *sync.Map
	lock               *sync.Mutex
}

func NewBeaconFetch(client *beacon.Client, fileStore *FileStore, genesisPeriod uint64, fetchDataResp chan FetchDataResponse) (*BeaconFetch, error) {
	maxReqNums := &atomic.Int64{}
	maxReqNums.Store(0)
	return &BeaconFetch{
		beaconClient:       client,
		currentReqNums:     maxReqNums,
		fileStore:          fileStore,
		exit:               make(chan struct{}, 1),
		fetchProofResponse: fetchDataResp,
		genesisSyncPeriod:  genesisPeriod,
		fetchQueue:         NewQueueWithCapacity(MaxQueueSize),
		cache:              new(sync.Map),
		lock:               new(sync.Mutex),
	}, nil
}

func (bf *BeaconFetch) canNewRequest() bool {
	bf.lock.Lock()
	defer bf.lock.Unlock()
	return bf.fetchQueue.CanPush()
}

func (bf *BeaconFetch) NewUpdateRequest(period uint64) {
	bf.lock.Lock()
	defer bf.lock.Unlock()
	logger.Debug("add new update request to queue: %v", period)
	fetchCacheKey := cacheKey(PeriodUpdateType, period)
	if _, exists := bf.cache.Load(fetchCacheKey); exists {
		return
	}
	bf.cache.Store(fetchCacheKey, struct{}{})
	bf.innerNewUpdateRequest(false, period)

}

func (bf *BeaconFetch) GenesisUpdateRequest() {
	bf.lock.Lock()
	defer bf.lock.Unlock()
	logger.Debug("add  genesis request to queue: %v", bf.genesisSyncPeriod)
	fetchCacheKey := cacheKey(GenesisUpdateType, bf.genesisSyncPeriod)
	if _, exists := bf.cache.Load(fetchCacheKey); exists {
		return
	}
	bf.cache.Store(fetchCacheKey, struct{}{})
	bf.innerNewGenesisRequest(false)
}

func (bf *BeaconFetch) fetch() error {
	if bf.fetchQueue.Len() == 0 {
		time.Sleep(2 * time.Second)
		return nil
	}
	if bf.currentReqNums.Load() > MaxReqNums {
		//logger.Warn("fetch too many request now: maxReqNums:%v", MaxReqNums)
		time.Sleep(2 * time.Second)
		return nil
	}
	element := bf.fetchQueue.Back()
	request, ok := element.Value.(FetchRequest)
	if !ok {
		logger.Error("should never happen,parse proof request error")
		time.Sleep(2 * time.Second)
		return nil
	}
	bf.fetchQueue.Remove(element)
	bf.currentReqNums.Add(1)
	time.Sleep(1 * time.Second)
	logger.Debug("get fetch request period:%v,type:%v,fetch data now", request.period, request.UpdateType.String())
	if request.UpdateType == GenesisUpdateType {
		go bf.getGenesisData(bf.genesisSyncPeriod)
	} else if request.UpdateType == PeriodUpdateType {
		go bf.getUpdateData(request.period)
	}
	return nil
}

func (bf *BeaconFetch) Fetch() {
	for {
		select {
		case <-bf.exit:
			logger.Info("beacon Fetch fetch goroutine exit now ...")
			return
		default:
			err := bf.fetch()
			if err != nil {
				logger.Error(err.Error())
			}
		}
	}
}

func (bf *BeaconFetch) innerNewUpdateRequest(highPriority bool, period uint64) {
	if highPriority {
		bf.fetchQueue.PushBack(FetchRequest{
			period:     period,
			UpdateType: PeriodUpdateType,
		})
	} else {
		bf.fetchQueue.PushFront(FetchRequest{
			period:     period,
			UpdateType: PeriodUpdateType,
		})
	}

}

func (bf *BeaconFetch) innerNewGenesisRequest(highPriority bool) {
	if highPriority {
		bf.fetchQueue.PushBack(FetchRequest{
			period:     bf.genesisSyncPeriod,
			UpdateType: GenesisUpdateType,
		})
	} else {
		bf.fetchQueue.PushFront(FetchRequest{
			period:     bf.genesisSyncPeriod,
			UpdateType: GenesisUpdateType,
		})
	}
}

func (bf *BeaconFetch) getGenesisData(period uint64) {
	defer func() {
		bf.currentReqNums.Add(-1)
		bf.cache.Delete(cacheKey(GenesisUpdateType, period))
	}()
	bootStrap, err := bf.beaconClient.Bootstrap(bf.genesisSyncPeriod)
	if err != nil {
		logger.Error("get bootstrap error:%v %v", bf.genesisSyncPeriod, err)
		// todo retry have higher priority
		bf.innerNewGenesisRequest(true)
		return
	}
	err = bf.fileStore.StoreGenesisUpdate(bootStrap)
	if err != nil {
		// todo
		logger.Error("store genesis update error:%v %v", bf.genesisSyncPeriod, err)
		return
	}
	updateResponse := FetchDataResponse{
		UpdateType: GenesisUpdateType,
		period:     period,
	}
	logger.Debug("success get genesis update data:%v", period)
	bf.fetchProofResponse <- updateResponse
}

func (bf *BeaconFetch) getUpdateData(period uint64) {
	defer func() {
		bf.currentReqNums.Add(-1)
		bf.cache.Delete(cacheKey(PeriodUpdateType, period))
	}()
	updates, err := bf.beaconClient.GetLightClientUpdates(period, 1)
	if err != nil {
		logger.Error("get light client updates error:%v %v", period, err)
		// todo
		bf.innerNewUpdateRequest(true, period)
		return
	}
	err = bf.fileStore.StoreUpdate(period, updates)
	if err != nil {
		// todo
		logger.Error("store update error:%v %v", period, err)
		return
	}
	updateResponse := FetchDataResponse{
		period:     period,
		UpdateType: PeriodUpdateType,
	}
	logger.Debug("success get update data:%v", period)
	bf.fetchProofResponse <- updateResponse
}

func (bf *BeaconFetch) Close() error {
	close(bf.exit)
	return nil
}

func cacheKey(fetchType FetchType, period uint64) string {
	return fmt.Sprintf("%v-%v", fetchType, period)
}
