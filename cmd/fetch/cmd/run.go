package cmd

import (
	"fmt"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/node"
	"github.com/lightec-xyz/daemon/rpc/beacon"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var dataDir string
var beaconUrl string
var genesisSlot uint64
var runCmd = &cobra.Command{
	Use:     "run",
	Long:    ``,
	Example: `./fetch run --genesisSlot 157 --beaconUrl "http://127.0.0.1:9870"  --datadir "./test"`,
	Run: func(cmd *cobra.Command, args []string) {
		if dataDir == "" || beaconUrl == "" {
			fmt.Printf("datadir or beaconUrl can not be empty\n")
			return
		}
		fu, err := NewFetch(beaconUrl, dataDir, genesisSlot)
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
	runCmd.Flags().Uint64Var(&genesisSlot, "genesisSlot", 0, "genesisSlot value")
	rootCmd.AddCommand(runCmd)
}

type Fetch struct {
	fetch      node.IFetch
	exitSignal chan os.Signal
}

func NewFetch(url, datadir string, genesisSlot uint64) (*Fetch, error) {
	fmt.Printf("datadir: %v,beaconUrl: %v,genesisSlot: %v \n", datadir, url, genesisSlot)
	err := logger.InitLogger(&logger.LogCfg{
		LogDir: "logs",
		File:   true,
	})
	if err != nil {
		return nil, err
	}
	client, err := beacon.NewClient(url)
	if err != nil {
		logger.Error("new beacon client error: %v", err)
		return nil, err
	}
	tables := []node.Table{node.GenesisTable, node.IndexTable,
		node.FinalityTable, node.UpdateTable}
	fileStorage, err := node.NewFileStorage(datadir, genesisSlot, tables...)
	if err != nil {
		logger.Error("new fileStorage error: %v", err)
		return nil, err
	}
	fetch, err := node.NewFetch(client, fileStorage, genesisSlot, nil)
	if err != nil {
		logger.Error("new fetch error: %v", err)
		return nil, err
	}
	return &Fetch{
		fetch:      fetch,
		exitSignal: make(chan os.Signal, 1),
	}, nil
}

func (f *Fetch) Run() error {
	logger.Info("start Fetch now...")
	err := f.fetch.Init()
	if err != nil {
		logger.Error("init fetch error: %v", err)
		return err
	}
	go node.DoTimerTask("fetch-finality-update", 1*time.Minute, f.fetch.FinalityUpdate, f.exitSignal)
	go node.DoTimerTask("fetch-update", 1*time.Minute, f.fetch.LightClientUpdate, f.exitSignal)
	signal.Notify(f.exitSignal, syscall.SIGTERM, syscall.SIGQUIT)
	for {
		msg := <-f.exitSignal
		switch msg {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGTSTP:
			logger.Info("get shutdown signal ,waiting exit now ...")
			err := f.Close()
			if err != nil {
				logger.Error(err.Error())
			}
			return nil
		}
	}

}

func (f *Fetch) Close() error {
	logger.Close()
	return f.fetch.Close()
}
