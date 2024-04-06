package cmd

import (
	"encoding/hex"
	"github.com/lightec-xyz/reLight/circuits/unit"
	"github.com/spf13/cobra"
)

// outerCmd represents the outer command
var outerCmd = &cobra.Command{
	Use:   "outer",
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
		innerFp, err := cmd.Flags().GetString("innerFp")
		if err != nil {
			panic(err)
		}
		innerFpBytes, err := hex.DecodeString(innerFp)
		if err != nil {
			panic(err)
		}
		outerConfig := unit.NewOuterConfig(dataDir, srsDir, subDir, innerFpBytes)
		inner := unit.NewOuter(&outerConfig)
		err = inner.Setup()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(outerCmd)
}
