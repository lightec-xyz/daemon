package node

import (
	"time"
)

const (
	LightecNetwork        = "mainnet"
	BtcOperatorAddress    = "bc1ququplway2l5prvt299j900s5sm3p6zw9r7d5lx84w0qgczd4sp5scxcpls"
	BtcLockScript         = "0020e0381fbba457e811b16a296457be1486e21d09c51f9b4f98f573c08c09b58069"
	BtcMultiSig           = "522102e510903d107b5594a5ee854eecb2858aeed5e84838c36fb12041676f71a17eba2103ea366ed0cfa0f48ee1e40ae321dab882f017fc8e7cb6a2159ced6fc42c6746da210218e65b47da2b63543f5b40c4b98532a97c737fa39c7e18d117bd3351eabbdc6c53ae"
	EthZkBridgeAddress    = "0xB86E9A8391d3df83F53D3f39E3b5Fce4D7da405d"
	EthZkBtcAddress       = "0x199CC8f0ac008Bdc8cF0B1CCd5187F84E168C4D2"
	EthUtxoManagerAddress = "0x2635Dc72706478F4bD784A8D04B3e0af8AB053dc"
	EthBtcTxVerifyAddress = "0x45bc26C34f9e6BA7704cBBf8dE090f2D1fEfA097"
	FeePoolAddr           = "0x7be6F1ECac63c8562Da8fF769347c45fc4590bFb"
	OasisSignerAddr       = "0xA81Fc99DBC654D68513B8C1475aFeC3B5d76496e"
	IcpPublicKey          = "03183007b9afcfa519871885380d4dfd1144269d8050ec2a51992065af2a87d3df"
	GenesisRoot           = "52bbd8287d0e455ce6cd732fa8a5f003e2ad82fd0ed3a59516f9ae1642f1b182"
	BtcScanTime           = 1 * time.Minute
	EthScanTime           = 30 * time.Second

	BlockSingerId = "xdqo6-dqaaa-aaaal-qsqva-cai"
	IcpTxSingerId = "wlkxr-hqaaa-aaaad-aaxaa-cai"

	InitBitcoinHeight  = 2812015
	InitBeaconHeight   = 153
	InitEthereumHeight = 1489369

	DepositTopic       = "0xd063609fea0cb9b8a1b53a4fbf0e659c270b3bc99eab08dcc7f4433b4937e074"
	RedeemTopic        = "0x379299efe6911678ce0f23cfce13a7c61a5b2c1723f583f9217b6ee0887b3ef4"
	UpdateUtxoTopic    = "0x5d6cc5c33e60ae274f09159956ae8fd20271c63c95b9004445271437b335ed6e"
	DepositRewardTopic = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
	RedeemRewardTopic  = "0xaaa8a5dd1203745b5ddf449a9bc6cd2b6ec919f7b839ef4326133ccf6fbe7bd0"
)

var SgxServerUrl []string = []string{""} // todo

// *********************Other**************************************

const (
	// MigrateProto operator migrate address
	MigrateProto = "6a141234560000000000000000000000000000000000"
	// MinDepositValue deposit min value
	MinDepositValue    = float64(0.00021000) // 21000 sats
	GeneratorVersion   = 1
	NodeVersion        = "1.0.0"
	RpcRegisterName    = "zkbtc"
	BitcoinNetwork     = "bitcoin"
	EthereumNetwork    = "ethereum"
	BitcoinAgentName   = "bitcoinAgent"
	EthereumAgentName  = "ethereumAgent"
	BeaconAgentName    = "beaconAgent"
	BtcLiteCacheHeight = 24 * 6 * 45
)

type Mode string

const (
	LiteMode Mode = "lite"
	FullMode Mode = "full"
)
