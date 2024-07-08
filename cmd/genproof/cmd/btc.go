package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var btcCmd = &cobra.Command{
	Use:   "btc",
	Short: "A brief description of your command",
	Long:  `A`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("btc called")
	},
}

func init() {
	rootCmd.AddCommand(btcCmd)

}
