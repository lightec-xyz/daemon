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
	TestnetBtcOperatorAddress = "tb1qzurfm8wattd63nk9ytvent8u869k0krq0087qcrupuneeacgsj4qudvlln"
	TestnetBtcLockScript      = "0x002017069d9ddd5adba8cec522d999acfc3e8b67d8607bcfe0607c0f279cf70884aa"
	TestnetMultiSig           = "0x522103af9e8678ba9d10de6df594766b3ecf4890c5d0306656c331eff5b746d8d449f52102971351ad0a4e80b4d61003a152c746bde6d7ac5cba52466727c611fdc8c20f5b21037a0e87290a962bad95efe3edbb61d70a7ee6cf7d5d5afbbdd3490ce098cf562953ae"
	TestnetEthZkBridgeAddress = "0x95E728c449ffF859CfbC100835e53064E6134f33"
	TestnetEthZkBtcAddress    = "0xACFB20CF49C542e8a79BFB9f3ed5b70818fa6407"
	TestEthUtxoManagerAddress = "0xb2677028b91b3BC34cADdbF23b3824585C781C4e"
	TestEthBtcTxVerifyAddress = "0x2EAda3FE674DaF7968FA95E8ea83971b322F86A5"
	TestnetFeePoolAddr        = "0xE9507Ecfc12b13213aA4094FbDE9A49F7031244F"
	TestnetOasisSignerAddr    = "0xbaa175e4286bB735265FecAD25c60e329359D8bF"
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
	GeneratorVersion  = "1.0.0"
	NodeVersion       = "1.0.0"
	RpcRegisterName   = "zkbtc"
	BitcoinNetwork    = "bitcoin"
	EthereumNetwork   = "ethereum"
	BitcoinAgentName  = "bitcoinAgent"
	EthereumAgentName = "ethereumAgent"
	BeaconAgentName   = "beaconAgent"
)
