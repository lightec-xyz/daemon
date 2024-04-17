package node

import (
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"os"
	"time"
)

func doTask(name string, fn func() error, exit chan os.Signal) {
	logger.Info("%v goroutine start ...", name)
	for {
		select {
		case <-exit:
			logger.Info("%v goroutine exit now ...", name)
			return
		default:
			err := fn()
			if err != nil {
				logger.Error("%v error %v", name, err.Error())
			}
		}
	}
}

func doTimerTask(name string, interval time.Duration, fn func() error, exit chan os.Signal) {
	logger.Info("%v ticker goroutine start ...", name)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-exit:
			logger.Info("%v goroutine exit now ...", name)
			return
		case <-ticker.C:
			err := fn()
			if err != nil {
				logger.Error("%v error %v", name, err.Error())
			}
		}
	}
}

func doProofRequestTask(name string, req chan []*common.ZkProofRequest, fn func(req []*common.ZkProofRequest) error, exit chan os.Signal) {
	logger.Info("%v goroutine start ...", name)
	for {
		select {
		case <-exit:
			logger.Info("%v goroutine exit now ...", name)
			return
		case request := <-req:
			err := fn(request)
			if err != nil {
				logger.Error("%v error %v", name, err.Error())
			}
		}

	}
}

func doFetchRespTask(name string, resp chan FetchDataResponse, fn func(resp FetchDataResponse) error, exit chan os.Signal) {
	logger.Info("%v goroutine start ...", name)
	for {
		select {
		case <-exit:
			logger.Info("%v goroutine exit now ...", name)
			return
		case response := <-resp:
			err := fn(response)
			if err != nil {
				logger.Error("%v error %v", name, err.Error())
			}
		}
	}
}

func doProofResponseTask(name string, resp chan *common.ZkProofResponse, fn func(resp *common.ZkProofResponse) error, exit chan os.Signal) {
	logger.Info("%v goroutine start ...", name)
	for {
		select {
		case <-exit:
			logger.Info("%v goroutine exit now ...", name)
			return
		case response := <-resp:
			err := fn(response)
			if err != nil {
				logger.Error("%v error %v", name, err.Error())
			}
		}
	}
}
