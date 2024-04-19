package common

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/lightec-xyz/daemon/rpc/ethereum"
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

func GetSlotByHash(client *ethereum.Client, hash string) (uint64, error) {
	txHash := common.HexToHash(hash)
	receipt, err := client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		return 0, err
	}
	slot, err := GetSlot(receipt.BlockNumber.Int64())
	if err != nil {
		return 0, err
	}
	return slot, nil
}
