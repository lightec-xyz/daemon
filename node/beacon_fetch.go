package node

import (
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc/beacon"
	"sync"
	"sync/atomic"
)

const MaxReqNums = 10

type BeaconFetch struct {
	beaconClient      *beacon.Client
	currentReqNums    *atomic.Int64
	downloadResponse  chan downloadResponse
	downloadRequest   chan downloadRequest
	fileStore         *FileStore
	exit              chan struct{}
	fetchProofRequest chan FetchDataResponse
	genesisSyncPeriod uint64
	download          *sync.Map
}

func NewBeaconFetch(client *beacon.Client, fileStore *FileStore, unitRequest chan FetchDataResponse) (*BeaconFetch, error) {
	maxReqNums := &atomic.Int64{}
	maxReqNums.Store(0)
	return &BeaconFetch{
		beaconClient:      client,
		currentReqNums:    maxReqNums,
		fileStore:         fileStore,
		exit:              make(chan struct{}, 1),
		fetchProofRequest: unitRequest,
		download:          new(sync.Map),
	}, nil
}

func (bf *BeaconFetch) canNewRequest() bool {
	return bf.currentReqNums.Load() < MaxReqNums
}

func (bf *BeaconFetch) NewUpdateRequest(period uint64) {
	bf.downloadRequest <- downloadRequest{
		period:     period,
		UpdateType: SyncComUnitType,
	}
}

func (bf *BeaconFetch) GenesisUpdateRequest() {
	bf.downloadRequest <- downloadRequest{
		period:     bf.genesisSyncPeriod,
		UpdateType: SyncComGenesisType,
	}
}

func (bf *BeaconFetch) FetchUpdate() {
	for {
		select {
		case <-bf.exit:
			logger.Info("beacon Run fetch goroutine exit now ...")
			return
		case request := <-bf.downloadRequest:
			if _, exists := bf.download.Load(request.period); exists {
				continue
			}
			logger.Info("get update request %v", request.period)
			bf.currentReqNums.Add(1)
			if request.UpdateType == SyncComGenesisType {
				go bf.getGenesisData(bf.genesisSyncPeriod)
			} else {
				go bf.getUpdateData(request.period)
			}
		}
	}
}

func (bf *BeaconFetch) FetchRepose() {
	for {
		select {
		case <-bf.exit:
			logger.Info("beacon FetchRepose fetch goroutine exit now ...")
			return
		case response := <-bf.downloadResponse:
			if response.status == Done {
				bf.download.Delete(response.period)
				logger.Info("success update response %v", response.period)
				bf.currentReqNums.Add(-1)
				if response.reqType == SyncComGenesisType {
					bf.fetchProofRequest <- FetchDataResponse{
						period:  response.period,
						reqType: SyncComGenesisType,
					}
				} else {
					bf.fetchProofRequest <- FetchDataResponse{
						period:  response.period,
						reqType: SyncComUnitType,
					}
				}

			} else {
				// retry until success
				logger.Warn("fail update response %v", response.period)
				bf.downloadRequest <- downloadRequest{
					period:     response.period,
					UpdateType: response.reqType,
				}
			}
		}
	}
}

func (bf *BeaconFetch) Close() error {
	close(bf.exit)
	return nil
}

func (bf *BeaconFetch) getGenesisData(period uint64) {
	updateResponse := downloadResponse{
		reqType: SyncComGenesisType,
		period:  period,
	}
	bootStrap, err := bf.beaconClient.GetBootstrap(uint64(bf.genesisSyncPeriod) * 32)
	if err != nil {
		updateResponse.status = Fail
		logger.Error("get bootstrap error:%v", err)
		return
	}
	err = bf.fileStore.StoreGenesisUpdate(bootStrap)
	if err != nil {
		// todo
		updateResponse.status = Fail
		logger.Error(err.Error())
		return
	}
	updateResponse.status = Done
	bf.downloadResponse <- updateResponse
}

func (bf *BeaconFetch) getUpdateData(index uint64) {
	updateResponse := downloadResponse{
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
	bf.downloadResponse <- updateResponse
}
