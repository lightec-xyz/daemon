package node

import (
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc/beacon"
	"sync/atomic"
)

const MaxReqNums = 10

type BeaconFetch struct {
	beaconClient   *beacon.Client
	currentReqNums *atomic.Int64
	updateResponse chan UpdateResponse
	updateRequest  chan UpdateRequest
	fileStore      *FileStore
	exit           chan struct{}
}

func NewBeaconFetch(client *beacon.Client, fileStore *FileStore) (*BeaconFetch, error) {
	maxReqNums := &atomic.Int64{}
	maxReqNums.Store(0)
	return &BeaconFetch{
		beaconClient:   client,
		currentReqNums: maxReqNums,
		fileStore:      fileStore,
		exit:           make(chan struct{}, 1),
	}, nil
}

func (bf *BeaconFetch) CanAddRequest() bool {
	return bf.currentReqNums.Load() < MaxReqNums
}

func (bf *BeaconFetch) SendRequest(period uint64) {
	bf.updateRequest <- UpdateRequest{
		period: period,
	}
}

func (bf *BeaconFetch) run() {
	for {
		select {
		case <-bf.exit:
			logger.Info("beacon Run fetch goroutine exit now ...")
			return
		case request := <-bf.updateRequest:
			logger.Info("get update request %v", request.period)
			bf.currentReqNums.Add(1)
			go bf.getUpdate(request.period)
		}
	}
}

func (bf *BeaconFetch) Response() {
	for {
		select {
		case <-bf.exit:
			logger.Info("beacon Response fetch goroutine exit now ...")
			return
		case response := <-bf.updateResponse:
			if response.status == Done {
				logger.Info("success update response %v", response.period)
				bf.currentReqNums.Add(-1)
			} else {
				// retry until success
				logger.Warn("fail update response %v", response.period)
				bf.updateRequest <- UpdateRequest{
					period: response.period,
				}
			}
		}
	}
}

func (bf *BeaconFetch) Close() error {
	close(bf.exit)
	return nil
}

func (bf *BeaconFetch) getUpdate(index uint64) {
	updateResponse := UpdateResponse{
		period: index,
	}
	updates, err := bf.beaconClient.GetLightClientUpdates(index, 1)
	if err != nil {
		updateResponse.status = Fail
		logger.Error("get light client updates error:%v", err)
		return
	}
	err = bf.fileStore.StoreUpdate(index, updates)
	if err != nil {
		// todo
		updateResponse.status = Fail
		logger.Error(err.Error())
		return
	}
	updateResponse.status = Done
	bf.updateResponse <- updateResponse
}
