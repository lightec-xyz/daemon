package cmd

import (
	"fmt"
	"github.com/lightec-xyz/daemon/node"
	"github.com/lightec-xyz/reLight/circuits/utils"
	"github.com/spf13/cobra"
)

var unitCmd = &cobra.Command{
	Use:   "unit",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		index, err := cmd.Flags().GetUint64("index")
		if err != nil {
			fmt.Printf("get index error: %v \n", err)
			return
		}
		fileStore, err := node.NewFileStore(datadir, genesisSlot)
		if err != nil {
			fmt.Printf("new file store error: %v \n", err)
			return
		}
		update, ok, err := node.GetSyncCommitUpdate(fileStore, index)
		if err != nil {
			fmt.Printf("get update error: %v \n", err)
			return
		}
		if !ok {
			fmt.Printf("update not found \n")
			return
		}
		err = UnitProve(update)
		if err != nil {
			fmt.Printf("prove error: %v \n", err)
			return
		}
	},
}

func init() {
	unitCmd.Flags().Uint64("index", 0, "proof period")
	rootCmd.AddCommand(unitCmd)
}

func UnitProve(update *utils.LightClientUpdateInfo) error {
	panic(update)
}
