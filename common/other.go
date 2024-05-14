package common

import (
	"strconv"

	"github.com/lightec-xyz/provers/utils"
)

// todo

func GetSlot(blockNumber int64) (uint64, error) {
	//url := "https://holesky.beaconcha.in"
	slot, err := utils.GetSlotOfEth1Block("https://holesky.beaconcha.in", blockNumber)
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(slot, 10, 64)
}

func GetNearTxSlot(slot uint64) uint64 {
	tmp := slot % 32
	if tmp == 0 {
		return slot
	}
	c := slot - tmp
	return c + 32
}
