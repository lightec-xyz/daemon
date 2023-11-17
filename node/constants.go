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
	TestnetBtcOperatorAddress = "tb1qwpcasq5qj5hmaz9wtulvjxt79w8um6hzz2f977xz8q69h39ya2rq7e2n5y"
	TestnetBtcLockScript      = "0x00207071d80280952fbe88ae5f3ec9197e2b8fcdeae212925f78c238345bc4a4ea86"
	TestnetMultiSig           = "0x5221038b6ce7e785f30c0eee59deda56c132a291d81800d2040bfb1e7b367c0e01f16d2102971351ad0a4e80b4d61003a152c746bde6d7ac5cba52466727c611fdc8c20f5b2103727dd022322bf061c92e06de4ed3608dc5f1bde0f6e13b5f7e9f160ae2f6517b53ae"
	TestnetEthZkBridgeAddress = "0x1974004A05516799E123628b8926b51e970f313f"
	TestnetEthZkBtcAddress    = "0xC3e425523878eA044178C4feaaFeB7Abd631e9e5"
	TestEthUtxoManagerAddress = "0xde2b14fF56c692a45c86D33D1E8CE7da883E5A60"
	TestEthBtcTxVerifyAddress = "0x03D1a417B139203E74443E73C78BB3c6CF038F99"
	TestnetFeePoolAddr        = "0xCfA0BF8729ff82Bb4249B2d495f9bf2d2B040c0A"
	TestnetOasisSignerAddr    = "0x5286d6EF240C09c0Fbd85DA530dcAB712f8Fa5C8"
	TestnetMigrateProto       = "6a141234560000000000000000000000000000000000"
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
	GeneratorVersion  = "1.0.0"
	NodeVersion       = "1.0.0"
	RpcRegisterName   = "zkbtc"
	BitcoinNetwork    = "bitcoin"
	EthereumNetwork   = "ethereum"
	BitcoinAgentName  = "bitcoinAgent"
	EthereumAgentName = "ethereumAgent"
	BeaconAgentName   = "beaconAgent"
)
