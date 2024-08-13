package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var config string
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := readRunConfig(config)
		if err != nil {
			fmt.Printf("read config: %v \n", err)
			return
		}
		btcSetup, err := NewBtcSetup(config)
		if err != nil {
			fmt.Printf("new btcsetup error: %v \n", err)
			return
		}
		err = btcSetup.Run()
		if err != nil {
			fmt.Printf("run btcsetup error: %v \n", err)
			return
		}
	},
}

func init() {
	runCmd.Flags().StringVar(&config, "config", "config.json", "config file")
	rootCmd.AddCommand(runCmd)
}
