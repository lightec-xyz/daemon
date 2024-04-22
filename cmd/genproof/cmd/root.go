package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var paramFile string
var datadir string
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
	rootCmd.PersistentFlags().StringVar(&paramFile, "param", "", "param file path")
	rootCmd.PersistentFlags().StringVar(&datadir, "datadir", "", "data dir")
}
