package cmd

import (
	"fmt"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/node"
	"github.com/lightec-xyz/daemon/rpc/beacon"
	"github.com/lightec-xyz/daemon/store"
	"github.com/spf13/cobra"
	"math/big"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

var dataDir string
var beaconUrl string
var genesisPeriod uint64
var runCmd = &cobra.Command{
	Use:     "run",
	Long:    ``,
	Example: `./fetch run --genesisPeriod 157 --beaconUrl "http://127.0.0.1:9870"  --datadir "./test"`,
	Run: func(cmd *cobra.Command, args []string) {
		if dataDir == "" || beaconUrl == "" {
			fmt.Printf("datadir or beaconUrl can not be empty\n")
			return
		}
		fu, err := NewFetch(beaconUrl, dataDir, genesisPeriod)
		if err != nil {
			fmt.Printf("new Fetch error: %v \n", err)
			return
		}
		err = fu.Run()
		if err != nil {
			fmt.Printf("run Fetch error: %v \n", err)
		}
	},
}

func init() {
	runCmd.Flags().StringVar(&beaconUrl, "beaconUrl", "", "rpc beaconUrl")
	runCmd.Flags().StringVar(&dataDir, "datadir", "", "storage dir")
	runCmd.Flags().Uint64Var(&genesisPeriod, "genesisPeriod", 0, "genesisPeriod value")
	rootCmd.AddCommand(runCmd)
}

type Fetch struct {
	client              *beacon.Client
	datadir             string
	updateStore         *store.FileStore
	finalityUpdateStore *store.FileStore
	cache               *node.CacheState
	genesisPeriod       uint64
	maxReqs             *atomic.Int64
	exitSignal          chan os.Signal
}

func NewFetch(url, datadir string, genesisPeriod uint64) (*Fetch, error) {
	fmt.Printf("datadir: %v,beaconUrl: %v,genesisPeriod: %v \n", datadir, url, genesisPeriod)
	err := logger.InitLogger()
	if err != nil {
		return nil, err
	}
	client, err := beacon.NewClient(url)
	if err != nil {
		logger.Error("new beacon client error: %v", err)
		return nil, err
	}
	finalityUpdateStore, err := store.NewFileStore(fmt.Sprintf("%s/finalityUpdate", datadir))
	if err != nil {
		logger.Error("new finalityUpdateStore error: %v", err)
		return nil, err
	}
	updateStore, err := store.NewFileStore(fmt.Sprintf("%s/update", datadir))
	if err != nil {
		logger.Error("new updateStore error: %v", err)
		return nil, err
	}
	maxReqs := atomic.Int64{}
	maxReqs.Store(0)
	return &Fetch{
		client:              client,
		finalityUpdateStore: finalityUpdateStore,
		updateStore:         updateStore,
		datadir:             datadir,
		genesisPeriod:       genesisPeriod,
		cache:               node.NewCacheState(),
		exitSignal:          make(chan os.Signal, 1),
		maxReqs:             &maxReqs,
	}, nil
}

func (f *Fetch) Run() error {
	logger.Info("start Fetch now...")
	go node.DoTimerTask("fetch-finality-update", 1*time.Minute, f.fetchFinalityUpdate, f.exitSignal)
	go node.DoTimerTask("fetch-update", 1*time.Minute, f.fetchUpdate, f.exitSignal)
	signal.Notify(f.exitSignal, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT, syscall.SIGTSTP, syscall.SIGQUIT)
	for {
		msg := <-f.exitSignal
		switch msg {
		case syscall.SIGHUP:
			logger.Info("get exit sign")
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGTSTP:
			logger.Info("get shutdown signal ,waiting exit now ...")
			err := f.Close()
			if err != nil {
				logger.Error(err.Error())
			}
			return nil
		}
	}

}

func (f *Fetch) fetchUpdate() error {
	if f.maxReqs.Load() >= 5 {
		return nil
	}
	latestPeriod, err := f.client.GetFinalizedSyncPeriod()
	if err != nil {
		logger.Error("get finalized sync period error: %v", err)
		return err
	}
	indexes, err := f.updateStore.AllIndexes()
	if err != nil {
		logger.Error("all indexes error: %v", err)
		return err
	}
	for index := f.genesisPeriod; index <= latestPeriod; index++ {
		if f.maxReqs.Load() >= 5 {
			return nil
		}
		if _, ok := indexes[index]; !ok {
			if !f.cache.Check(index) {
				f.cache.Store(index, nil)
				f.maxReqs.Add(1)
				go f.getUpdate(index)
			}
		}
	}
	return nil
}

func (f *Fetch) getUpdate(index uint64) error {
	defer func() {
		f.cache.Delete(index)
		f.maxReqs.Add(-1)
	}()
	logger.Debug("fetch update index: %v", index)
	updates, err := f.client.GetLightClientUpdates(index, 1)
	if err != nil {
		logger.Error("get light client updates error: %v", err)
		return err
	}
	key := fmt.Sprintf("%v", index)
	err = f.updateStore.Store(key, updates)
	if err != nil {
		logger.Error("store error: %v", err)
		return err
	}
	logger.Info("success store update: %v", key)
	return nil
}

func (f *Fetch) fetchFinalityUpdate() error {
	logger.Debug("check finality update now")
	finalityUpdate, err := f.client.GetFinalityUpdate()
	if err != nil {
		logger.Error("get finality proof error: %v", err)
		return err
	}
	slotBig, ok := big.NewInt(0).SetString(finalityUpdate.Data.FinalizedHeader.Slot, 10)
	if !ok {
		logger.Error("parse slot error: %v", finalityUpdate.Data.FinalizedHeader.Slot)
		return fmt.Errorf("parse slot error: %v", finalityUpdate.Data.FinalizedHeader.Slot)
	}
	slot := slotBig.Uint64()
	keyName := fmt.Sprintf("%v", slot)
	exists, err := f.finalityUpdateStore.CheckExists(keyName)
	if err != nil {
		logger.Error("check exists error: %v", err)
		return err
	}
	if exists {
		return nil
	}
	err = f.finalityUpdateStore.Store(keyName, finalityUpdate)
	if err != nil {
		logger.Error("store error: %v", err)
		return err
	}
	logger.Info("success store finality update file: %v", keyName)
	return nil
}

func (f *Fetch) Close() error {
	close(f.exitSignal)
	return nil
}
