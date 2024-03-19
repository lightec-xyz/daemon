package cmd

import (
	"fmt"
	"github.com/lightec-xyz/daemon/circuits"
	"github.com/spf13/cobra"
)

var paramFile string

var unitCmd = &cobra.Command{
	Use:   "unit",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		unit := circuits.NewUnit()
		optUnit := circuits.OptUnit{
			DataDir:    dataDir,
			SrsDataDir: srcDataDir,
			SubDir:     fmt.Sprintf("%s/sc", dataDir),
			ParamFile:  paramFile,
		}
		err := unit.GenerateProof(&optUnit)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(unitCmd)
	unitCmd.PersistentFlags().StringVar(&paramFile, "file", "", "update data file path")
	if paramFile == "" {
		panic("param file can not be empty")
	}

}
