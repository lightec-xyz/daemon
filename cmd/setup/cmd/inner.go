package cmd

import (
	"github.com/lightec-xyz/reLight/circuits/unit"
	"github.com/spf13/cobra"
)

var innerCmd = &cobra.Command{
	Use:   "inner",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		dataDir, err := cmd.Flags().GetString("datadir")
		if err != nil {
			panic(err)
		}
		srsDir, err := cmd.Flags().GetString("srsDir")
		if err != nil {
			panic(err)
		}
		subDir, err := cmd.Flags().GetString("subDir")
		if err != nil {
			panic(err)
		}
		innerConfig := unit.NewInnerConfig(dataDir, srsDir, subDir)
		inner := unit.NewInner(&innerConfig)
		err = inner.Setup()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(innerCmd)
}
