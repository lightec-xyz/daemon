package node

import "time"

// ********************* mainnet ************************
const (
	BtcOperatorAddress = "bcrt1qalv7aduqdpz9wc4fut3nt44tsf42anleed76yj3el3rgd4rgldvq2aw6ze"
	BtcMultiSigPublic1 = "03bd96c4d06aa773e5d282f0b6bccd1fb91268484918648ccda1ae768209edb050"
	BtcMultiSigPublic2 = "03aa9c4245340a02864c903f7f9e7bc9ef1cc374093aacbf72b614002f6d8c8c22"
	BtcMultiSigPublic3 = "03351a7971bf7ed886fca99aebdc3b195fc79ffe93b499e2309a4e69ab115405e0"
	BtcMultiNRequired  = 2

	EthZkBridgeAddress = "0xbdfb7b89e9c77fe647ac1628416773c143ca4b51"
	EthZkBtcAddress    = "0x5898953ff9c1c11a8a6bc578bd6c93aabcd1f083"

	BtcScanTime = 1 * time.Minute
	EthScanTime = 5 * time.Second

	InitBitcoinHeight  = 2540942
	InitEthereumHeight = 10127532
)

var RedeemLogTopices = []string{"0xb28ad0403b0a341130002b9eef334c5daa3c1002a73dd90d4626f7079d0a804a"}

// ********************* testnet ************************
const (
	TestnetBtcOperatorAddress = "tb1qalv7aduqdpz9wc4fut3nt44tsf42anleed76yj3el3rgd4rgldvq8yyuhr"
	TestnetBtcMultiSigPublic1 = "03bd96c4d06aa773e5d282f0b6bccd1fb91268484918648ccda1ae768209edb050"
	TestnetBtcMultiSigPublic2 = "03aa9c4245340a02864c903f7f9e7bc9ef1cc374093aacbf72b614002f6d8c8c22"
	TestnetBtcMultiSigPublic3 = "03351a7971bf7ed886fca99aebdc3b195fc79ffe93b499e2309a4e69ab115405e0"
	TestnetBtcMultiNRequired  = 2

	TestnetEthZkBridgeAddress = "0xbdfb7b89e9c77fe647ac1628416773c143ca4b51"
	TestnetEthZkBtcAddress    = "0x5898953ff9c1c11a8a6bc578bd6c93aabcd1f083"

	TestnetBtcScanTime = 1 * time.Minute
	TestnetEthScanTime = 5 * time.Second

	TestnetInitBitcoinHeight  = 2540942
	TestnetInitEthereumHeight = 10127532
)

var TestnetRedeemLogTopices = []string{"0xb28ad0403b0a341130002b9eef334c5daa3c1002a73dd90d4626f7079d0a804a"}

// ********************* local ************************
const (
	LocalBtcOperatorAddress = "bcrt1qalv7aduqdpz9wc4fut3nt44tsf42anleed76yj3el3rgd4rgldvq2aw6ze"
	LocalBtcMultiSigPublic1 = "03bd96c4d06aa773e5d282f0b6bccd1fb91268484918648ccda1ae768209edb050"
	LocalBtcMultiSigPublic2 = "03aa9c4245340a02864c903f7f9e7bc9ef1cc374093aacbf72b614002f6d8c8c22"
	LocalBtcMultiSigPublic3 = "03351a7971bf7ed886fca99aebdc3b195fc79ffe93b499e2309a4e69ab115405e0"
	LocalBtcMultiNRequired  = 2

	LocalEthZkBridgeAddress = "0xbdfb7b89e9c77fe647ac1628416773c143ca4b51"
	LocalEthZkBtcAddress    = "0x5898953ff9c1c11a8a6bc578bd6c93aabcd1f083"

	LocalBtcScanTime = 1 * time.Minute
	LocalEthScanTime = 5 * time.Second

	LocalInitBitcoinHeight  = 2540942
	LocalInitEthereumHeight = 10127532
)

var LocalRedeemLogTopices = []string{"0xb28ad0403b0a341130002b9eef334c5daa3c1002a73dd90d4626f7079d0a804a"}

// ***********************************************************
const (
	Deposit = "deposit"
	Redeem  = "redeem"
)

const RpcRegisterName = "zkbtc"

const ProofPrefix = "p"

var (
	btcCurHeightKey = []byte("btcCurHeight")
	ethCurHeightKey = []byte("ethCurHeight")
)
