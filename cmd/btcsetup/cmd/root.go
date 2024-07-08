package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var datadir string
var srsdir string

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

func check() error {
	if datadir == "" {
		return fmt.Errorf("datadir can not be empty")
	}
	if srsdir == "" {
		return fmt.Errorf("srsdir can not be empty")
	}
	return nil
}

func init() {
	rootCmd.PersistentFlags().StringVar(&datadir, "datadir", "", "data directory")
	rootCmd.PersistentFlags().StringVar(&srsdir, "srsdir", "", "srs directory")
	err := check()
	if err != nil {
		fmt.Printf("error: %v \n", err)
		os.Exit(1)
	}
}
