package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc/beacon"
	"github.com/spf13/cobra"
	"math/big"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var dataDir string
var beaconUrl string
var runCmd = &cobra.Command{
	Use:     "run",
	Long:    ``,
	Example: `./finalityUpdate run --beaconUrl "http://127.0.0.1:9870"  --datadir "./test"`,
	Run: func(cmd *cobra.Command, args []string) {
		if dataDir == "" || beaconUrl == "" {
			fmt.Printf("datadir or beaconUrl can not be empty\n")
			return
		}
		fu, err := NewFinalityUpdate(beaconUrl, dataDir)
		if err != nil {
			fmt.Printf("new FinalityUpdate error: %v \n", err)
			return
		}
		err = fu.Run()
		if err != nil {
			fmt.Printf("run FinalityUpdate error: %v \n", err)
		}
	},
}

func init() {
	runCmd.Flags().StringVar(&beaconUrl, "beaconUrl", "", "rpc beaconUrl")
	runCmd.Flags().StringVar(&dataDir, "datadir", "", "storage dir")
	rootCmd.AddCommand(runCmd)
}

type FinalityUpdate struct {
	client     *beacon.Client
	datadir    string
	exitSignal chan os.Signal
}

func NewFinalityUpdate(url, datadir string) (*FinalityUpdate, error) {
	err := logger.InitLogger()
	if err != nil {
		return nil, err
	}
	client, err := beacon.NewClient(url)
	if err != nil {
		return nil, err
	}
	return &FinalityUpdate{
		client:     client,
		datadir:    datadir,
		exitSignal: make(chan os.Signal, 1),
	}, nil
}

func (f *FinalityUpdate) Run() error {
	logger.Info("start FinalityUpdate now...")
	go f.Fetch()
	signal.Notify(f.exitSignal, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT, syscall.SIGTSTP, syscall.SIGQUIT)
	for {
		msg := <-f.exitSignal
		switch msg {
		case syscall.SIGHUP:
			logger.Info("daemon get SIGHUP")
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

func (f *FinalityUpdate) Fetch() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-f.exitSignal:
			logger.Info("fetch exit now...")
			return
		case <-ticker.C:
			err := f.fetch()
			if err != nil {
				logger.Error("fetch error: %v", err)
			}
		}
	}
}

func (f *FinalityUpdate) fetch() error {
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
	path := fmt.Sprintf("%s/%d", f.datadir, slot)
	exists, err := common.FileExists(path)
	if err != nil {
		logger.Error("file exists error: %v", err)
		return err
	}
	if exists {
		return nil
	}
	data, err := json.Marshal(finalityUpdate)
	if err != nil {
		logger.Error("json marshal error: %v", err)
		return err
	}
	err = common.WriteFile(path, data)
	if err != nil {
		logger.Error("write file error: %v", err)
		return err
	}
	logger.Info("success write file: %s", slot)
	return nil
}

func (f *FinalityUpdate) Close() error {
	close(f.exitSignal)
	return nil
}
