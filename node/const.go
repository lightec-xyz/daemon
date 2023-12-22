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
	EthZkBridgeAddress = "0x8b404b735afe5bcdce85a1ce753c79715f86062c"
	EthZkBtcAddress    = "0xdf68798c22c683f72e3a1359f9de8bbedb7ab920"

	BtcScanTime = 1 * time.Minute
	EthScanTime = 5 * time.Second

	InitBitcoinHeight  = 2540942
	InitEthereumHeight = 10127532

	LogDepositAddr   = "0x8b404b735afe5bcdce85a1ce753c79715f86062c"
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
	TestnetEthZkBridgeAddress = "0x8b404b735afe5bcdce85a1ce753c79715f86062c"
	TestnetEthZkBtcAddress    = "0xdf68798c22c683f72e3a1359f9de8bbedb7ab920"

	TestnetBtcScanTime = 1 * time.Minute
	TestnetEthScanTime = 5 * time.Second

	TestnetInitBitcoinHeight  = 2540942
	TestnetInitEthereumHeight = 10127532

	TestLogDepositAddr   = "0x8b404b735afe5bcdce85a1ce753c79715f86062c"
	TestLogRedeemAddr    = "0xa7becea4ce9040336d7d4aad84e684d1daeabea1"
	TestTopicDepositAddr = "0xb28ad0403b0a341130002b9eef334c5daa3c1002a73dd90d4626f7079d0a804a"
	TestTopicRedeemAddr  = "0x975dbbd59299029fdfc12db336ede29e2e2b2d117effa1a45be55f0b4f9cfbce"
)

var TestLogAddrs = []string{TestLogDepositAddr, TestLogRedeemAddr}
var TestLogTopics = []string{TestTopicDepositAddr, TestTopicRedeemAddr}

// ********************* local ************************
const (
	Lighteclocal            = "local"
	LocalBtcOperatorAddress = "bcrt1qalv7aduqdpz9wc4fut3nt44tsf42anleed76yj3el3rgd4rgldvq2aw6ze"
	LocalBtcMultiSigPublic1 = "03bd96c4d06aa773e5d282f0b6bccd1fb91268484918648ccda1ae768209edb050"
	LocalBtcMultiSigPublic2 = "03aa9c4245340a02864c903f7f9e7bc9ef1cc374093aacbf72b614002f6d8c8c22"
	LocalBtcMultiSigPublic3 = "03351a7971bf7ed886fca99aebdc3b195fc79ffe93b499e2309a4e69ab115405e0"
	LocalBtcMultiNRequired  = 2
	LocalBtcNetwork         = btctx.RegTest
	LocalBtcScanTime        = 15 * time.Second
	LocalEthScanTime        = 5 * time.Second

	LocalInitBitcoinHeight  = 12980
	LocalInitEthereumHeight = 576019

	LocalEthZkBridgeAddress = "0x8b404b735afe5bcdce85a1ce753c79715f86062c"
	LocalEthZkBtcAddress    = "0xdf68798c22c683f72e3a1359f9de8bbedb7ab920"

	LocalLogDepositAddr = "0x52ebc075616195cc7deb79d5c21bd9b04acc33ee"
	LocalLogRedeemAddr  = "0x8b404b735afe5bcdce85a1ce753c79715f86062c"

	LocalTopicDepositAddr = "0x975dbbd59299029fdfc12db336ede29e2e2b2d117effa1a45be55f0b4f9cfbce"
	LocalTopicRedeemAddr  = "0xb28ad0403b0a341130002b9eef334c5daa3c1002a73dd90d4626f7079d0a804a"
)

var LocalLogAddrs = []string{LocalLogDepositAddr, LocalLogRedeemAddr}
var LocalLogTopics = []string{LocalTopicDepositAddr, LocalTopicRedeemAddr}

// ***********************************************************

const BtcDeposit = 0
const BtcRedeem = 1

const RpcRegisterName = "zkbtc"
