package cmd

import (
	"fmt"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/spf13/cobra"
)

var setup = &cobra.Command{
	Use:  "run",
	Long: "./setup  --datadir ./test --srsdir ./srs run --type beaconOuter --group all",
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
		if circuitType == "" && group == "" {
			fmt.Printf("type or group is empty \n")
			return
		}

		chainId, err := cmd.Flags().GetInt("chainId")
		if err != nil {
			fmt.Printf("get type error: %v \n", err)
			return
		}
		zkbtcBridgeAddr, err := cmd.Flags().GetString("zkbtcBridgeAddr")
		if err != nil {
			fmt.Printf("get type error: %v \n", err)
			return
		}
		icpPublickey, err := cmd.Flags().GetString("icpPublickey")
		if err != nil {
			fmt.Printf("get type error: %v \n", err)
			return
		}
		if group == Ethereum.String() || circuitType == ethTxInEth2.String() {
			if chainId == 0 || zkbtcBridgeAddr == "" {
				fmt.Printf("chainId or zkbtcBridgeAddr is empty \n")
				return
			}
		} else if group == Bitcoin.String() || circuitType == btcTxInChain.String() {
			if icpPublickey == "" {
				fmt.Printf("icpPublickey is empty \n")
				return
			}
		} else if group == All.String() {
			if chainId == 0 || zkbtcBridgeAddr == "" || icpPublickey == "" {
				fmt.Printf("chainId or zkbtcBridgeAddr or icpPublickey is empty \n")
				return
			}
		}
		fmt.Printf("group: %v ,circuitType: %v,chainId: %v,zkbtcBridgeAddr: %v,icpPublickey: %v \n",
			group, circuitType, chainId, zkbtcBridgeAddr, icpPublickey)
		circuitSetup := NewCircuitSetup(datadir, srsddir, zkbtcBridgeAddr, icpPublickey, chainId)
		if group != "" {
			if err = circuitSetup.SetupGroup(Group(group)); err != nil {
				fmt.Printf("setup group error: %v \n", err)
				return
			}
		} else {
			if circuitType != "" {
				if err = circuitSetup.Setup(CircuitType(circuitType)); err != nil {
					fmt.Printf("setup circuit type error: %v \n", err)
					return
				}
			}
		}
		fileMd5s, err := GetFilesMd5(datadir)
		if err != nil {
			fmt.Printf("get files md5 error: %v \n", err)
			return
		}
		err = SaveMd5sToJson(fileMd5s, datadir+"/md5.json")
		if err != nil {
			fmt.Printf("save md5 error: %v \n", err)
			return
		}
	},
}

func init() {
	setup.Flags().String("type", "", "setup circuit type value: beaconInner, beaconOuter, beaconSyncCommittee, beaconGenesis, beaconRecursive,ethTxInEth2,ethBeaconHeader,ethFinalityHeader,ethRedeem, btcBase, btcMiddle, btcUpper")
	setup.Flags().String("group", "", "batch setup circuit group value: all, bitcoin, beacon, ethereum")
	setup.Flags().Int("chainId", 0, "batch setup circuit group value: all, bitcoin, beacon, ethereum")
	setup.Flags().String("zkbtcBridgeAddr", "", "batch setup circuit group value: all, bitcoin, beacon, ethereum")
	setup.Flags().String("icpPublickey", "", "batch setup circuit group value: all, bitcoin, beacon, ethereum")

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
