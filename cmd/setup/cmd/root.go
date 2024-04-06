package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cmd",
	Short: "A brief description of your application",
	Long:  ``,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().String("datadir", "./data", "data directory")
	rootCmd.Flags().String("srsdir", "./data", "data directory")
	rootCmd.Flags().String("subdir", "./data", "data directory")
}
