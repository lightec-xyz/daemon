package cmd

import (
	"fmt"
	"github.com/lightec-xyz/daemon/node"
	"github.com/spf13/cobra"
)

var filststoreCmd = &cobra.Command{
	Use:   "filestore",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		InitLogger()
		height, err := cmd.Flags().GetInt64("height")
		if err != nil {
			fmt.Printf("get height error: %v \n", err)
			return
		}
		if height == 0 {
			fmt.Printf("height is empty \n")
			return
		}
		cfg, err := readCfg(cfgFile)
		if err != nil {
			fmt.Printf("read config error: %v %v \n", cfgFile, err)
			return
		}
		fileStorage, err := node.NewFileStorage(cfg.Datadir, 0, 0)
		if err != nil {
			fmt.Printf("new fileStorage error: %v \n", err)
			return
		}
		fmt.Printf("start remove btc proof  >=%v \n", height)
		err = fileStorage.RemoveBtcProof(uint64(height))
		if err != nil {
			fmt.Printf("remove btc proof error: %v \n", err)
			return
		}
		fmt.Printf("remove btc proof success \n")

	},
}

func init() {
	rootCmd.AddCommand(filststoreCmd)
	filststoreCmd.Flags().Int64("height", 0, "btc height")
}
