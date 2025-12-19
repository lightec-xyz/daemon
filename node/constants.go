package node

import (
	"time"
)

const (
	LightecNetwork        = "mainnet"
	BtcOperatorAddress    = "bc1qcwc08898aseahq2e5920m9395py3jjtm2xnt0s90slqra7cry2dq2xv9uq"
	BtcLockScript         = "0020c3b0f39ca7ec33db8159a154fd9625a04919497b51a6b7c0af87c03efb03229a"
	BtcMultiSig           = "522102e510903d107b5594a5ee854eecb2858aeed5e84838c36fb12041676f71a17eba2103beeb82e07222ca3b22057028311c052c9acd9b844310463086d033dabcb8c3fb210218e65b47da2b63543f5b40c4b98532a97c737fa39c7e18d117bd3351eabbdc6c53ae"
	EthZkBridgeAddress    = "0xF90966fd006a5B18Cb0E3A0568226010CED426FD"
	EthZkBtcAddress       = "0x199CC8f0ac008Bdc8cF0B1CCd5187F84E168C4D2"
	EthUtxoManagerAddress = "0x205a1E85C7d4d0fcd4344335120181aB5e796562"
	EthBtcTxVerifyAddress = "0x1F0f891fB88287091DFc6225038336207374ec79"
	FeePoolAddr           = "0x7be6F1ECac63c8562Da8fF769347c45fc4590bFb"
	OasisSignerAddr       = "0xA81Fc99DBC654D68513B8C1475aFeC3B5d76496e"
	IcpPublicKey          = "03183007b9afcfa519871885380d4dfd1144269d8050ec2a51992065af2a87d3df"
	GenesisRoot           = "52bbd8287d0e455ce6cd732fa8a5f003e2ad82fd0ed3a59516f9ae1642f1b182"
	BlockSingerId         = "xdqo6-dqaaa-aaaal-qsqva-cai"
	IcpTxSingerId         = "wlkxr-hqaaa-aaaad-aaxaa-cai"

	// we don't need those default values, every config file should have these three values
	// InitBitcoinHeight  = 2812015
	// InitBeaconHeight   = 153
	// InitEthereumHeight = 1489369

	DepositTopic    = "0xd063609fea0cb9b8a1b53a4fbf0e659c270b3bc99eab08dcc7f4433b4937e074" // the UTXOAdded event from the UTXO Manager contract
	RedeemTopic     = "0x379299efe6911678ce0f23cfce13a7c61a5b2c1723f583f9217b6ee0887b3ef4" // the CreateRedeemUnsignedTx event from the zkBTC bridge contract
	UpdateUtxoTopic = "0x5d6cc5c33e60ae274f09159956ae8fd20271c63c95b9004445271437b335ed6e" // the ChangeUTXOUpdated event from the UTXO Manager contract

	// DepositRewardTopic = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef" // the transfer call of minted zkBTC, the zkBTC bridge contract
	// should have used "0x9d186086b05c87611fc259bfefbd7b1261b646a6715ce64039050613604e9b41", // the MinerDepositReward event from the Fee Pool Contract
	// see zkBtc Explorer common.Constants

	// RedeemRewardTopic  = "0xaaa8a5dd1203745b5ddf449a9bc6cd2b6ec919f7b839ef4326133ccf6fbe7bd0" // the MinerRedeemReward event to the Fee Pool contract
)

var SgxServerUrl []string = []string{""} // todo

// *********************Other**************************************

const (
	BtcScanTime = 5 * time.Minute
	EthScanTime = 3 * time.Minute
	// MigrateProto operator migrate address
	MigrateProto = "6a141234560000000000000000000000000000000000"
	// MinDepositValue deposit min value
	MinDepositValue      = float64(0.00021000) // 21000 sats
	GeneratorVersion     = 1
	NodeVersion          = "1.0.0"
	RpcRegisterName      = "zkbtc"
	BitcoinNetwork       = "bitcoin"
	EthereumNetwork      = "ethereum"
	BitcoinAgentName     = "bitcoinAgent"
	EthereumAgentName    = "ethereumAgent"
	BeaconAgentName      = "beaconAgent"
	BtcLiteCacheHeight   = 24 * 6 * 45
	ProofExpired         = 5 * time.Hour
	BtcClientCacheHeight = 100
)

type Mode string

const (
	LiteMode Mode = "lite"
	FullMode Mode = "full"
)
