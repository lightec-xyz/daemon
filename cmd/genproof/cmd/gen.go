package cmd

import (
	"fmt"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/spf13/cobra"
)

var paramPath string
var index uint64
var proofType string
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("datadir: %v,setupdir: %v,genesisPeriod: %v \n", datadir, setupDir, genesisSlot)

	},
}

func init() {
	err := logger.InitLogger(nil)
	if err != nil {
		panic(err)
	}
	genCmd.Flags().StringVar(&paramPath, "paramPath", "", "param file path")
	genCmd.Flags().Uint64Var(&index, "index", 0, "proof index")
	genCmd.Flags().StringVar(&proofType, "proofType", "", "proof type")
	rootCmd.AddCommand(genCmd)
}
