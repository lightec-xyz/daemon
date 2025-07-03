package proof

import (
	"fmt"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/node"
	"github.com/spf13/cobra"
)

var removeBtcProofCmd = &cobra.Command{
	Use:   "filestore",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cfgFile, err := cmd.Root().PersistentFlags().GetString("config")
		if err != nil {
			fmt.Printf("get config error: %v \n", err)
			return
		}
		fmt.Printf("config file: %v \n", cfgFile)
		var cfg node.RunConfig
		err = common.ReadObj(cfgFile, &cfg)
		if err != nil {
			fmt.Printf("read config error: %v %v \n", cfgFile, err)
			return
		}
		height, err := cmd.Flags().GetInt64("height")
		if err != nil {
			fmt.Printf("get height error: %v \n", err)
			return
		}
		if height == 0 {
			fmt.Printf("height is empty \n")
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
	removeBtcProofCmd.Flags().Int64("height", 0, "btc height")
}
