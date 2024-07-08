package cmd

import (
	"fmt"
	"github.com/lightec-xyz/btc_provers/circuits/baselevel"
	"github.com/lightec-xyz/btc_provers/circuits/midlevel"
	"github.com/lightec-xyz/btc_provers/circuits/upperlevel"
	"github.com/spf13/cobra"
)

const (
	btcBase   = "base"
	btcMiddle = "middle"
	btcUp     = "up"
)

var btcBaseCmd = &cobra.Command{
	Use:   "btcBase",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		flag, err := cmd.Flags().GetString("flag")
		if err != nil {
			fmt.Printf("get flag error: %v \n", err)
			return
		}
		switch flag {
		case btcBase:
			err := newBaseSetup(datadir, srsdata)
			if err != nil {
				fmt.Printf("newBaseSetup error: %v \n", err)
				return
			}
		case btcMiddle:
			err := newMiddleSetup(datadir, srsdata)
			if err != nil {
				fmt.Printf("newMiddleSetup error: %v \n", err)
				return
			}
		case btcUp:
			err := newUplevelSetup(datadir, srsdata)
			if err != nil {
				fmt.Printf("upperlevel.Setup error: %v \n", err)
				return
			}
		default:
			err := newBaseSetup(datadir, srsdata)
			if err != nil {
				fmt.Printf("newBaseSetup error: %v \n", err)
				return
			}
			err = newMiddleSetup(datadir, srsdata)
			if err != nil {
				fmt.Printf("newMiddleSetup error: %v \n", err)
				return
			}
			err = newUplevelSetup(datadir, srsdata)
			if err != nil {
				fmt.Printf("upperlevel.Setup error: %v \n", err)
				return
			}
		}
	},
}

func init() {
	btcBaseCmd.Flags().String("flag", "", "")
	rootCmd.AddCommand(btcBaseCmd)
}

func newBaseSetup(datadir, srsdir string) error {
	err := baselevel.Setup(datadir, srsdir)
	if err != nil {
		return err
	}
	return nil
}

func newMiddleSetup(datadir, srsdir string) error {
	err := midlevel.Setup(datadir, srsdir)
	if err != nil {
		return err
	}
	return nil
}

func newUplevelSetup(datadir, srsdir string) error {
	err := upperlevel.Setup(datadir, srsdir)
	if err != nil {
		return err
	}
	return nil
}
