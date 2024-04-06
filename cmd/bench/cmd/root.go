package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var dataDir string
var srcDataDir string

var rootCmd = &cobra.Command{
	Use:   "bench",
	Short: "test generate zkb proof",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&dataDir, "datadir", "", "datadir default current running directory")
	rootCmd.PersistentFlags().StringVar(&srcDataDir, "srsDatadir", "", "srs data directory")
	if srcDataDir == "" {
		panic("srs data directory can not be empty")
	}
	if dataDir == "" {
		currentPath, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		dataDir = fmt.Sprintf("%s/unitProof-%s", currentPath, getTimeStr())
	}
}

func getTimeStr() string {
	currentTime := time.Now()
	year := currentTime.Format("2006")
	month := currentTime.Format("01")
	day := currentTime.Format("02")
	hour := currentTime.Format("15")
	minute := currentTime.Format("04")
	second := currentTime.Format("05")
	return fmt.Sprintf("%s-%s-%s-%s-%s-%s", year, month, day, hour, minute, second)
}
