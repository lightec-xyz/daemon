package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var datadir string
var srsddir string

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
	rootCmd.PersistentFlags().StringVar(&datadir, "datadir", "", "data directory")
	rootCmd.PersistentFlags().StringVar(&srsddir, "srsdir", "", "srs directory")
}
