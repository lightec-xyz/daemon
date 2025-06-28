package node

import (
	btctx "github.com/lightec-xyz/daemon/rpc/bitcoin/common"
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

// ********************* testnet ************************
const (
	LightecTestnet            = "testnet"
	TestnetBtcOperatorAddress = "tb1qxv89x4mmu0e64e7g5zu3f4h9ar0zd7xjfw8fk3c7nd4vxh3swc6spz645v"
	TestnetBtcLockScript      = "0x0020330e53577be3f3aae7c8a0b914d6e5e8de26f8d24b8e9b471e9b6ac35e307635"
	TestnetMultiSig           = "0x52210327716110843c703f59f78c6d5e1b9c634307b87f95d161251a97f722c8bb9aa62102971351ad0a4e80b4d61003a152c746bde6d7ac5cba52466727c611fdc8c20f5b21037a0e87290a962bad95efe3edbb61d70a7ee6cf7d5d5afbbdd3490ce098cf562953ae"
	TestnetEthZkBridgeAddress = "0x430C22DA4251D18f710347e316cB54303e6bA155"
	TestnetEthZkBtcAddress    = "0x3540141758cf3824c00d9CF07143f435b18D25c3"
	TestEthUtxoManagerAddress = "0xa6B0B23eCC6fcfa3bC8Bb62d176EFCc70E7e09d9"
	TestEthBtcTxVerifyAddress = "0x2B7DFb385C81582FE151E85Db0E160A17E4971BF"
	TestnetFeePoolAddr        = "0x4996dCcA6fEe37aF47c3073c8D27fFC71eFc7a41"
	TestnetOasisSignerAddr    = "0xbAa42115A13d62B4e99a227bd650991B5DB1a6Bb"
	TestnetIcpPublicKey       = "0x02971351ad0a4e80b4d61003a152c746bde6d7ac5cba52466727c611fdc8c20f5b"
	TestnetBtcScanTime        = 1 * time.Minute
	TestnetEthScanTime        = 30 * time.Second

	TestnetIcpSingerId  = "6ybqh-oiaaa-aaaak-quffa-cai"
	TestnetIcpUrl       = "https://icp0.io"
	TestnetSgxServerUrl = ""

	TestnetInitBitcoinHeight  = 2812015
	TestnetInitBeaconHeight   = 153
	TestnetInitEthereumHeight = 1489369

	TestnetDepositTopic       = "0xd063609fea0cb9b8a1b53a4fbf0e659c270b3bc99eab08dcc7f4433b4937e074"
	TestnetRedeemTopic        = "0x379299efe6911678ce0f23cfce13a7c61a5b2c1723f583f9217b6ee0887b3ef4"
	TestnetUpdateUtxoTopic    = "0x5d6cc5c33e60ae274f09159956ae8fd20271c63c95b9004445271437b335ed6e"
	TestnetDepositRewardTopic = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
	TestnetRedeemRewardTopic  = "0xaaa8a5dd1203745b5ddf449a9bc6cd2b6ec919f7b839ef4326133ccf6fbe7bd0"
)

// *********************Other**************************************

const (
	// MigrateProto operator migrate address
	MigrateProto = "6a141234560000000000000000000000000000000000"
	// MinDepositValue deposit min value
	MinDepositValue   = float64(0.00002100) // 2100 sats
	GeneratorVersion  = 1
	NodeVersion       = "1.0.0"
	RpcRegisterName   = "zkbtc"
	BitcoinNetwork    = "bitcoin"
	EthereumNetwork   = "ethereum"
	BitcoinAgentName  = "bitcoinAgent"
	EthereumAgentName = "ethereumAgent"
	BeaconAgentName   = "beaconAgent"
)
