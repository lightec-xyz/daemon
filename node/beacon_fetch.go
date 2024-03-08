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
	beaconClient      *beacon.Client
	currentReqNums    *atomic.Int64
	downloadRequest   chan downloadRequest
	fileStore         *FileStore
	exit              chan struct{}
	fetchProofRequest chan FetchDataResponse
	genesisSyncPeriod uint64
	fetchQueue        *Queue
	cache             *sync.Map
}

func NewBeaconFetch(client *beacon.Client, fileStore *FileStore, genesisPeriod uint64, fetchDataResp chan FetchDataResponse) (*BeaconFetch, error) {
	maxReqNums := &atomic.Int64{}
	maxReqNums.Store(0)
	return &BeaconFetch{
		beaconClient:      client,
		currentReqNums:    maxReqNums,
		fileStore:         fileStore,
		exit:              make(chan struct{}, 1),
		fetchProofRequest: fetchDataResp,
		genesisSyncPeriod: genesisPeriod,
		fetchQueue:        NewQueue(),
		cache:             new(sync.Map),
	}, nil
}

func (bf *BeaconFetch) canNewRequest() bool {
	return bf.fetchQueue.Len() < MaxQueueSize
}

func (bf *BeaconFetch) NewUpdateRequest(period uint64) {
	logger.Debug("add new update request to queue: %v", period)
	fetchCacheKey := cacheKey(PeriodUpdateType, period)
	if _, exists := bf.cache.Load(fetchCacheKey); exists {
		return
	}
	bf.cache.Store(fetchCacheKey, struct{}{})
	bf.fetchQueue.PushFront(downloadRequest{
		period:     period,
		UpdateType: PeriodUpdateType,
	})
}

func (bf *BeaconFetch) GenesisUpdateRequest() {
	logger.Debug("add  genesis request to queue: %v", bf.genesisSyncPeriod)
	fetchCacheKey := cacheKey(GenesisUpdateType, bf.genesisSyncPeriod)
	if _, exists := bf.cache.Load(fetchCacheKey); exists {
		return
	}
	bf.cache.Store(fetchCacheKey, struct{}{})
	bf.fetchQueue.PushBack(downloadRequest{
		period:     bf.genesisSyncPeriod,
		UpdateType: GenesisUpdateType,
	})
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
	request, ok := element.Value.(downloadRequest)
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

func (bf *BeaconFetch) getGenesisData(period uint64) {
	defer func() {
		bf.currentReqNums.Add(-1)
		bf.cache.Delete(cacheKey(GenesisUpdateType, period))
	}()
	bootStrap, err := bf.beaconClient.Bootstrap(bf.genesisSyncPeriod)
	if err != nil {
		logger.Error("get bootstrap error:%v %v", bf.genesisSyncPeriod, err)
		// todo
		// retry again
		bf.GenesisUpdateRequest()
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
	bf.fetchProofRequest <- updateResponse
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
		// retry again
		bf.NewUpdateRequest(period)
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
	bf.fetchProofRequest <- updateResponse
}

func (bf *BeaconFetch) Close() error {
	close(bf.exit)
	return nil
}

func cacheKey(fetchType FetchType, period uint64) string {
	return fmt.Sprintf("%v-%v", fetchType, period)
}
