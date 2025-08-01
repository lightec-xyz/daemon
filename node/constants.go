package node

import (
	"time"
)

const (
	LightecNetwork        = "testnet"
	BtcOperatorAddress    = "tb1quxl2g39x50vlc4ycne2hrrq8dl3cqv84prdnvzlrshm2a4dkt24snxa5nz"
	BtcLockScript         = "0x0020e1bea444a6a3d9fc54989e55718c076fe38030f508db360be385f6aed5b65aab"
	BtcMultiSig           = "0x522102fdb41469ab1536cbfadd8a659b2e2667795d6a98b7cb696c1626a70e937bbe142102971351ad0a4e80b4d61003a152c746bde6d7ac5cba52466727c611fdc8c20f5b21037a0e87290a962bad95efe3edbb61d70a7ee6cf7d5d5afbbdd3490ce098cf562953ae"
	EthZkBridgeAddress    = "0x8a6415Ff6865f0194ef2742fFfCe1Ab861212cb1"
	EthZkBtcAddress       = "0xeCed2a18fB48B34671eE63E9370da5A9bE7493FB"
	EthUtxoManagerAddress = "0x579D5Bc54629E9F2C1974A8058B3399fA3B8b7e2"
	EthBtcTxVerifyAddress = "0x4991bD7f47221513E8c8e289CE0879c0C0C8bAC0"
	FeePoolAddr           = "0x285d14dB804c41A03C50e6e11dbBe12B8eC4AfDc"
	OasisSignerAddr       = "0x10622d51ABF42860111DA3329156f5ac56c135aF"
	IcpPublicKey          = "0x02971351ad0a4e80b4d61003a152c746bde6d7ac5cba52466727c611fdc8c20f5b"
	GenesisRoot           = "0x095016b5bb894d9076ed599e2ef0c54d88048d0d8192143b5008dac531a6a43a"
	BtcScanTime           = 3 * time.Minute
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
