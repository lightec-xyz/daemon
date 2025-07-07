package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/lightec-xyz/daemon/circuits"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/node"
	"github.com/spf13/cobra"
	"os"
)

var scRootCmd = &cobra.Command{
	Use:   "scRoot",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := logger.InitLogger(nil)
		if err != nil {
			fmt.Printf("init logger error: %v \n", err)
			return
		}
		period, err := cmd.Flags().GetInt("period")
		if err != nil {
			fmt.Printf("get period error: %v \n", err)
			return
		}
		genesisSlot, err := cmd.Flags().GetInt("genesisSlot")
		if err != nil {
			fmt.Printf("get period error: %v \n", err)
			return
		}
		fmt.Printf("period: %v genesisSlot: %v \n", period, genesisSlot)
		cfgBytes, err := os.ReadFile(cfgFile)
		if err != nil {
			fmt.Printf("read config error: %v %v \n", cfgFile, err)
			return
		}
		fmt.Printf("confg data: %v \n", string(cfgBytes))
		var runCfg node.RunConfig
		err = json.Unmarshal(cfgBytes, &runCfg)
		if err != nil {
			fmt.Printf("unmarshal config error: %v %v \n", cfgFile, err)
			return
		}
		config, err := node.NewConfig(runCfg)
		if err != nil {
			fmt.Printf("new config error: %v \n", err)
			return
		}
		fileStorage, err := node.NewFileStorage(config.Datadir, 0, 0)
		if err != nil {
			fmt.Printf("new fileStorage error: %v \n", err)
			return
		}
		prepared, err := node.NewPreparedData(fileStorage, nil, uint64(genesisSlot), 0, nil, nil,
			nil, nil, nil, "", false)
		if err != nil {
			fmt.Printf("new preparedData error: %v \n", err)
			return
		}
		update, ok, err := prepared.GetSyncCommitUpdate(uint64(period))
		if err != nil {
			fmt.Printf("get syncCommitUpdate error: %v \n", err)
			return
		}
		if !ok {
			fmt.Printf("get syncCommitUpdate error: %v \n", err)
			return
		}
		syncCommitRoot, err := circuits.SyncCommitRoot(update.CurrentSyncCommittee)
		if err != nil {
			fmt.Printf("get syncCommitRoot error: %v \n", err)
			return
		}
		fmt.Printf(" %v:syncCommitRoot: %x \n", period, syncCommitRoot)
	},
}

func init() {
	scRootCmd.Flags().Int("period", 0, "the period of the sc root")
	scRootCmd.Flags().Int("genesisSlot", 0, "the period of the sc root")
	rootCmd.AddCommand(scRootCmd)
}
