package common

import (
	"github.com/lightec-xyz/provers/utils"
	"math/big"
)

// todo

func GetSlot(blockNumber int64) (uint64, error) {
	//url := "https://holesky.beaconcha.in"
	slot, err := utils.GetSlotOfEth1Block("https://holesky.beaconcha.in", blockNumber)
	if err != nil {
		return 0, err
	}
	slotBig, ok := big.NewInt(0).SetString(slot, 10)
	if !ok {
		return 0, err
	}
	return slotBig.Uint64(), nil
}
