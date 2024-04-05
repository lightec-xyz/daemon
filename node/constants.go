package node

import (
	btctx "github.com/lightec-xyz/daemon/transaction/bitcoin"
	"time"
)

// ********************* mainnet ************************
const (
	LightecMainnet     = "mainnet"
	BtcOperatorAddress = "bcrt1qalv7aduqdpz9wc4fut3nt44tsf42anleed76yj3el3rgd4rgldvq2aw6ze"
	BtcMultiSigPublic1 = "03bd96c4d06aa773e5d282f0b6bccd1fb91268484918648ccda1ae768209edb050"
	BtcMultiSigPublic2 = "03aa9c4245340a02864c903f7f9e7bc9ef1cc374093aacbf72b614002f6d8c8c22"
	BtcMultiSigPublic3 = "03351a7971bf7ed886fca99aebdc3b195fc79ffe93b499e2309a4e69ab115405e0"
	BtcMultiNRequired  = 2
	BtcNetwork         = btctx.MainNet
	EthZkBridgeAddress = "0x07417a531e376ce150493ffa98cd5516b544441d"
	EthZkBtcAddress    = "0xd2a00777a7e5b6afaa5c053a6425619653541c82"

	BtcScanTime = 1 * time.Minute
	EthScanTime = 5 * time.Second

	InitBitcoinHeight  = 2540942
	InitBeaconHeight   = 1024256
	InitEthereumHeight = 10127532

	LogDepositAddr   = "0x07417a531e376ce150493ffa98cd5516b544441d"
	LogRedeemAddr    = "0xa7becea4ce9040336d7d4aad84e684d1daeabea1"
	TopicDepositAddr = "0xb28ad0403b0a341130002b9eef334c5daa3c1002a73dd90d4626f7079d0a804a"
	TopicRedeemAddr  = "0x975dbbd59299029fdfc12db336ede29e2e2b2d117effa1a45be55f0b4f9cfbce"
)

var LogAddrs = []string{LogDepositAddr, LogRedeemAddr}
var LogTopics = []string{TopicDepositAddr, TopicRedeemAddr}

// ********************* testnet ************************
const (
	LightecTestnet            = "testnet"
	TestnetBtcOperatorAddress = "tb1qalv7aduqdpz9wc4fut3nt44tsf42anleed76yj3el3rgd4rgldvq8yyuhr"
	TestnetBtcMultiSigPublic1 = "03bd96c4d06aa773e5d282f0b6bccd1fb91268484918648ccda1ae768209edb050"
	TestnetBtcMultiSigPublic2 = "03aa9c4245340a02864c903f7f9e7bc9ef1cc374093aacbf72b614002f6d8c8c22"
	TestnetBtcMultiSigPublic3 = "03351a7971bf7ed886fca99aebdc3b195fc79ffe93b499e2309a4e69ab115405e0"
	TestnetBtcMultiNRequired  = 2
	TestnetBtcNetwork         = btctx.TestNet
	TestnetEthZkBridgeAddress = "0x3ca427befe5b8b821c09a8d6425fbcee20f952f6"
	TestnetEthZkBtcAddress    = "0x3528594509fcf7b06f70976a9fae1c3b0ab92e22"

	TestnetBtcScanTime = 1 * time.Minute
	TestnetEthScanTime = 5 * time.Second

	TestnetInitBitcoinHeight  = 2544083
	TestnetInitBeaconHeight   = 1024256
	TestnetInitEthereumHeight = 598020

	TestLogRedeemAddr  = "0x3ca427befe5b8b821c09a8d6425fbcee20f952f6"
	TestLogDepositAddr = "0x96ffb80f74a646940569b599039e0fbd0b3a4711"

	TestTopicDepositAddr = "0x975dbbd59299029fdfc12db336ede29e2e2b2d117effa1a45be55f0b4f9cfbce"
	TestTopicRedeemAddr  = "0xb28ad0403b0a341130002b9eef334c5daa3c1002a73dd90d4626f7079d0a804a"
)

var TestLogAddrs = []string{TestLogDepositAddr, TestLogRedeemAddr}
var TestLogTopics = []string{TestTopicDepositAddr, TestTopicRedeemAddr}

// ********************* local ************************
const (
	Lighteclocal            = "local"
	LocalBtcOperatorAddress = "tb1qtysxx7zkmm5nwy0hv2mjxfrermsry2vjsygg0eqawwwp6gy4hl4s2tudtw"
	LocalBtcMultiSigPublic1 = "03bd96c4d06aa773e5d282f0b6bccd1fb91268484918648ccda1ae768209edb050"
	LocalBtcMultiSigPublic2 = "03aa9c4245340a02864c903f7f9e7bc9ef1cc374093aacbf72b614002f6d8c8c22"
	LocalBtcMultiSigPublic3 = "03351a7971bf7ed886fca99aebdc3b195fc79ffe93b499e2309a4e69ab115405e0"
	LocalBtcMultiNRequired  = 2
	LocalBtcNetwork         = btctx.RegTest
	LocalBtcScanTime        = 15 * time.Second
	LocalEthScanTime        = 5 * time.Second

	LocalInitBitcoinHeight  = 2585316
	LocalInitBeaconHeight   = 153 //slot of 596751 in holesky
	LocalInitEthereumHeight = 1286750

	LocalEthZkBridgeAddress = "0x8e4f5a8f3e24a279d8ed39e868f698130777fded"
	LocalEthZkBtcAddress    = "0xbf3041e37be70a58920a6fd776662b50323021c9"

	// utxo manager contract
	LocalLogDepositAddr = "0xab5146a46e90c497b3d23afab7ddaedf3ff61eaf"
	LocalLogRedeemAddr  = "0x19d376e6a10aad92e787288464d4c738de97d135"

	LocalTopicDepositAddr = "0x975dbbd59299029fdfc12db336ede29e2e2b2d117effa1a45be55f0b4f9cfbce"
	LocalTopicRedeemAddr  = "0x1e5e2baa6d11cc5bcae8c0d1187d7b9ebf13d6d9b932f7dbbf4e396438845fb8"
)

var LocalLogAddrs = []string{LocalLogDepositAddr, LocalLogRedeemAddr}
var LocalLogTopics = []string{LocalTopicDepositAddr, LocalTopicRedeemAddr}

// ***********************************************************

const RpcRegisterName = "zkbtc"

const BitcoinNetwork = "bitcoin"
const EthereumNetwork = "ethereum"
