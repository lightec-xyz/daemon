package cmd

import (
	"fmt"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/spf13/cobra"
)

var setup = &cobra.Command{
	Use:     "run",
	Long:    "./setup  --datadir ./test --srsdir ./srs run --type beaconOuter --group all",
	Example: "./setup  --datadir ./test --srsdir ./srs run --type beaconOuter --group all",
	Run: func(cmd *cobra.Command, args []string) {
		err := check()
		if err != nil {
			fmt.Printf("check error: %v \n", err)
			return
		}
		err = initLogger()
		if err != nil {
			fmt.Printf("init logger error: %v \n", err)
			return
		}
		circuitType, err := cmd.Flags().GetString("type")
		if err != nil {
			fmt.Printf("get type error: %v \n", err)
			return
		}
		group, err := cmd.Flags().GetString("group")
		if err != nil {
			fmt.Printf("get type error: %v \n", err)
			return
		}
		circuitSetup := NewCircuitSetup(datadir, srsddir)
		if group != "" {
			if err = circuitSetup.SetupGroup(Group(group)); err != nil {
				fmt.Printf("setup group error: %v \n", err)
				return
			}
		}
		if circuitType != "" {
			if err = circuitSetup.Setup(CircuitType(circuitType)); err != nil {
				fmt.Printf("setup circuit type error: %v \n", err)
				return
			}
		}

	},
}

func init() {
	setup.Flags().String("type", "", "setup circuit type value: beaconOuter, beaconInner, beaconUnit, beaconGenesis,...")
	setup.Flags().String("group", "", "batch setup circuit group value: all, bitcoin, beacon, ethereum")
	rootCmd.AddCommand(setup)
}

func initLogger() error {
	err := logger.InitLogger(&logger.LogCfg{
		IsStdout: true,
		File:     false,
	})
	return err
}
func check() error {
	if datadir == "" {
		return fmt.Errorf("datadir can not be empty")
	}
	if srsddir == "" {
		return fmt.Errorf("srsddir can not be empty")
	}
	return nil
}
