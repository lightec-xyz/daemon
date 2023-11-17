package ethereum

import "math/big"

type ZkbtcInfo struct {
	TotalCrossFound *big.Int
	DepositL2Reward *big.Int
	RedeemL2Reward  *big.Int
	NextHaving      *big.Int
}
