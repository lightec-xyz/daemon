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
	MaxReqNums   = 3
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
	genesisSlot        uint64
	fetchQueue         *Queue
	cache              *sync.Map
	lock               *sync.Mutex
}

func NewBeaconFetch(client *beacon.Client, fileStore *FileStore, genesisSlot uint64, fetchDataResp chan FetchDataResponse) (*BeaconFetch, error) {
	maxReqNums := &atomic.Int64{}
	maxReqNums.Store(0)
	return &BeaconFetch{
		beaconClient:       client,
		currentReqNums:     maxReqNums,
		fileStore:          fileStore,
		exit:               make(chan struct{}, 1),
		fetchProofResponse: fetchDataResp,
		genesisSyncPeriod:  genesisSlot / 8192,
		genesisSlot:        genesisSlot,
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
	exists, err := bf.fileStore.CheckUpdate(period)
	if err != nil {
		logger.Error("check update error:%v", err)
		return
	}
	if exists {
		return
	}
	bf.innerNewUpdateRequest(false, period)

}

func (bf *BeaconFetch) BootStrapRequest() {
	bf.lock.Lock()
	defer bf.lock.Unlock()
	logger.Debug("add  genesis request to queue: %v", bf.genesisSyncPeriod)
	fetchCacheKey := cacheKey(GenesisUpdateType, bf.genesisSyncPeriod)
	if _, exists := bf.cache.Load(fetchCacheKey); exists {
		return
	}
	// todo
	bf.cache.Store(fetchCacheKey, struct{}{})
	exists, err := bf.fileStore.CheckBootstrap()
	if err != nil {
		logger.Error("check genesis error:%v", err)
		return
	}
	if exists {
		return
	}
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
		logger.Error("should never happen,parse Proof request error")
		time.Sleep(2 * time.Second)
		return nil
	}
	bf.fetchQueue.Remove(element)
	bf.currentReqNums.Add(1)
	logger.Debug("get fetch request Period:%v,type:%v,fetch data now", request.period, request.UpdateType.String())
	if request.UpdateType == GenesisUpdateType {
		go func() {
			err := bf.getBootStrap()
			if err != nil {
				logger.Error("get genesis data error:%v", err)
				bf.innerNewGenesisRequest(true)
				return
			} else {
				bf.currentReqNums.Add(-1)
				bf.cache.Delete(cacheKey(GenesisUpdateType, bf.genesisSyncPeriod))
			}

		}()
	} else if request.UpdateType == PeriodUpdateType {
		go func() {
			err := bf.getUpdateData(request.period)
			if err != nil {
				logger.Error("get update data error:%v %v", err, request.period)
				bf.innerNewUpdateRequest(true, request.period)
			} else {
				bf.currentReqNums.Add(-1)
				bf.cache.Delete(cacheKey(PeriodUpdateType, request.period))
			}
		}()
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

func (bf *BeaconFetch) getBootStrap() error {
	// todo
	bootStrap, err := bf.beaconClient.Bootstrap(bf.genesisSlot)
	if err != nil {
		logger.Error("get bootstrap error:%v %v", bf.genesisSyncPeriod, err)
		return err
	}
	err = bf.fileStore.StoreBootstrap(bootStrap)
	if err != nil {
		// todo
		logger.Error("store genesis update error:%v %v", bf.genesisSyncPeriod, err)
		return err
	}
	updateResponse := FetchDataResponse{
		UpdateType: GenesisUpdateType,
		period:     bf.genesisSlot / 8192,
	}
	logger.Debug("success get genesis update data:%v", bf.genesisSlot/8192)
	bf.fetchProofResponse <- updateResponse
	return nil
}

func (bf *BeaconFetch) getUpdateData(period uint64) error {
	updates, err := bf.beaconClient.GetLightClientUpdates(period, 1)
	if err != nil {
		logger.Error("get light client updates error:%v %v", period, err)
		return err
	}
	if len(updates) != 1 {
		logger.Error("get light client updates error:%v %v", period, err)
		return nil
	}
	err = bf.fileStore.StoreUpdate(period, updates[0])
	if err != nil {
		// todo
		logger.Error("store update error:%v %v", period, err)
		return err
	}
	updateResponse := FetchDataResponse{
		period:     period,
		UpdateType: PeriodUpdateType,
	}
	logger.Debug("success get update data:%v", period)
	bf.fetchProofResponse <- updateResponse
	return nil
}

func (bf *BeaconFetch) Close() error {
	close(bf.exit)
	return nil
}

func cacheKey(fetchType FetchType, period uint64) string {
	return fmt.Sprintf("%v-%v", fetchType, period)
}
