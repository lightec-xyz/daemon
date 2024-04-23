package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var datadir string
var setupDir string
var genesisSlot uint64
var rootCmd = &cobra.Command{
	Use:   "cmd",
	Short: "",
	Long:  ``,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&datadir, "datadir", "", "data dir")
	rootCmd.PersistentFlags().StringVar(&setupDir, "setupDir", "", "setup dir")
	rootCmd.PersistentFlags().Uint64Var(&genesisSlot, "genesisSlot", 0, "genesis slot")
}
