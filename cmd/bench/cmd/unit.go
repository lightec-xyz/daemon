package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/lightec-xyz/daemon/circuits"
	"github.com/lightec-xyz/reLight/circuits/utils"
	"github.com/spf13/cobra"
	"os"
	"time"
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
		wrapTime(func() {
			err := unit.Prove(&optUnit)
			if err != nil {
				panic(err)
			}
		}, "generate unit proof")
		wrapTime(func() {
			paramBytes, err := os.ReadFile(paramFile)
			if err != nil {
				panic(err)
			}
			lightClientUpdateInfo := &utils.LightClientUpdateInfo{}
			err = json.Unmarshal(paramBytes, lightClientUpdateInfo)
			if err != nil {
				panic(err)
			}
			verify, err := unit.Verify(&optUnit, lightClientUpdateInfo)
			if err != nil {
				panic(err)
			}
			if verify {
				fmt.Println("verify success")
			} else {
				fmt.Println("verify failed")
			}
		}, "verify unit proof")
	},
}

func init() {
	rootCmd.AddCommand(unitCmd)
	unitCmd.PersistentFlags().StringVar(&paramFile, "file", "", "update data file path")
	if paramFile == "" {
		panic("param file can not be empty")
	}

}

func wrapTime(fn func(), desc ...string) {
	var name string
	if len(desc) != 0 {
		name = desc[0]
	}
	start := time.Now()
	fn()
	end := time.Now()
	t := end.Sub(start)
	fmt.Printf("%s task time  %02d:%02d:%02d\n", name, int(t.Hours()), int(t.Minutes())%60, int(t.Seconds())%60)
}
