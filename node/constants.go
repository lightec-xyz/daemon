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
	TestnetBtcOperatorAddress = "tb1qtysxx7zkmm5nwy0hv2mjxfrermsry2vjsygg0eqawwwp6gy4hl4s2tudtw"
	TestnetBtcLockScript      = "00205920637856dee93711f762b72324791ee0322992811087e41d739c1d2095bfeb"
	TestnetBtcMultiSigPublic1 = "0363f549d250342df02ee8b51ad6c9148dabc587c6569761ab58aa68488bd2e2c5"
	TestnetBtcMultiSigPublic2 = "031cbb294f9955d80f65d9499feaeb5cb29d44c070adddd75cd48a40791d39b971"
	TestnetBtcMultiSigPublic3 = "035c54e8287a7f7ba31886249fc89f295a4cb74cebf0d925f1eafe87f22fba57f9"
	TestnetBtcMultiNRequired  = 2
	TestnetEthZkBridgeAddress = "0x8e4f5a8f3e24a279d8ed39e868f698130777fded"
	TestnetEthZkBtcAddress    = "0xbf3041e37be70a58920a6fd776662b50323021c9"

	TestEthUtxoManagerAddress = "0x9d2aaea60dee441981edf44300c26f1946411548"

	TestnetBtcScanTime = 1 * time.Minute
	TestnetEthScanTime = 30 * time.Second

	TestnetInitBitcoinHeight  = 2812015
	TestnetInitBeaconHeight   = 153
	TestnetInitEthereumHeight = 1489369
	TestnetOasisSignerAddr    = "0x99e514Dc90f4Dd36850C893bec2AdC9521caF8BB"

	// utxo manager
	TestLogDepositAddr = "0x9d2aaea60dee441981edf44300c26f1946411548"

	TestLogRedeemAddr = "0x8e4f5a8f3e24a279d8ed39e868f698130777fded"

	TestTopicDepositAddr = "0xbfb6a0aa850eff6109c854ffb48321dcf37f02d6c7a44c46987a5ddf3419fc07"
	TestTopicRedeemAddr  = "0x1e5e2baa6d11cc5bcae8c0d1187d7b9ebf13d6d9b932f7dbbf4e396438845fb8"
)

// ********************* local ************************
const (
	Lighteclocal            = "local"
	LocalBtcOperatorAddress = "tb1qj4atskm3pj6m2achrc3mkdyj2dkgq2wvn9g37wfq60xp8pc6yvnsnnpxj4"
	LocalBtcMultiSigPublic1 = "034def276d763bfb937a4f356d26b58cb0428bc198d000b91630db5d04bb7f35dc"
	LocalBtcMultiSigPublic2 = "03183ee062dafa5a0e536ba497c6375a996364682bf22cd5de989df1b0b9d23621"
	LocalBtcMultiSigPublic3 = "03a868050ec7b61b6956d6c1ca722f4d2a32671902486980d5fd6ebf9b4c64dd93"
	LocalBtcMultiNRequired  = 2
	LocalBtcScanTime        = 15 * time.Second
	LocalEthScanTime        = 5 * time.Second

	LocalInitBitcoinHeight  = 2812015
	LocalInitBeaconHeight   = 153 //slot of 596751 in holesky
	LocalInitEthereumHeight = 1489369

	LocalOasisSignerAddr = "0x7ccCc552F55C05FD33d9827070E1dB3D28322622"

	LocalEthZkBridgeAddress = "0xb2631368c8c8151875ea67cb5faf8f1377ec02a0"
	LocalEthZkBtcAddress    = "0xbf3041e37be70a58920a6fd776662b50323021c9"
	LocalEthUtxoManager     = "0x9d2aaea60dee441981edf44300c26f1946411548"
	// utxo manager contract
	LocalLogDepositAddr = "0xe8965848879eb831e3c8f47d2256eff883d9a0d9"
	LocalLogRedeemAddr  = "0xb2631368c8c8151875ea67cb5faf8f1377ec02a0"

	LocalTopicDepositAddr = "0x975dbbd59299029fdfc12db336ede29e2e2b2d117effa1a45be55f0b4f9cfbce"
	LocalTopicRedeemAddr  = "0x1e5e2baa6d11cc5bcae8c0d1187d7b9ebf13d6d9b932f7dbbf4e396438845fb8"
)

// *********************Other**************************************

const RpcRegisterName = "zkbtc"

const BitcoinNetwork = "bitcoin"
const EthereumNetwork = "ethereum"

const BitcoinAgentName = "bitcoinAgent"
const EthereumAgentName = "ethereumAgent"
const BeaconAgentName = "beaconAgent"

const (
	GeneratorVersion = "1.0.0"
)
