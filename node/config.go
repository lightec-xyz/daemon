package node

import (
	"encoding/hex"
	"fmt"
	"github.com/lightec-xyz/daemon/logger"
	"os/user"
	"time"
)

type NodeConfig struct {
	DataDir string `json:"datadir"`
	Network string `json:"network"`
	Rpcbind string `json:"rpcbind"`
	RpcPort string `json:"rpcport"`

	BtcUrl           string           `json:"btcUrl"`
	BtcUser          string           `json:"btcUser"`
	BtcPwd           string           `json:"btcPwd"`
	BtcNetwork       string           `json:"btcNetwork"`
	BtcScanBlockTime time.Duration    `json:"btcBlockTime"`
	BtcOperatorAddr  string           `json:"btcOperatorAddr"`
	BtcPrivateKeys   []string         `json:"btcPrivateKeys"`
	BtcWhiteList     []string         `json:"btcWhiteList"`
	BtcInitHeight    int64            `json:"btcInitHeight"`
	MultiAddressInfo MultiAddressInfo `json:"multiAddressInfo"`

	EthInitHeight    int64          `json:"ethInitHeight"`
	EthWhiteList     []string       `json:"ethWhiteList"`
	EthUrl           string         `json:"ethUrl"`
	ZkBridgeAddr     string         `json:"zkBridgeAddr"`
	ZkBtcAddr        string         `json:"zkBtcAddr"`
	EthScanBlockTime time.Duration  `json:"ethBlockTime"`
	EthPrivateKey    string         `json:"ethPrivateKey"`
	LogAddr          []string       `json:"logAddr"`
	LogTopic         []string       `json:"logTopic"`
	Workers          []WorkerConfig `json:"workers"`
}

func NewNodeConfig(dataDir, network, rpcbind, rpcport, btcUrl, btcUser, btcPwd, btcNetwork, ethUrl, ethPrivateKey string) (NodeConfig, error) {
	var config NodeConfig
	current, err := user.Current()
	if err != nil {
		logger.Error("get current user error", err)
		return config, err
	}
	if dataDir == "" {
		dataDir = fmt.Sprintf("%v/.daemon", current.HomeDir)
	}
	if network == "" {
		network = "mainnet"
	}
	if rpcbind == "" {
		rpcbind = "127.0.0.1"
	}
	if rpcport == "" {
		rpcport = "30000"
	}
	if btcUrl == "" {
		return config, fmt.Errorf("btcUrl is empty")
	}
	if ethUrl == "" {
		return config, fmt.Errorf("ethUrl is empty")
	}
	if ethPrivateKey == "" {
		return config, fmt.Errorf("ethPrivateKey is empty")
	}

	switch network {
	case "mainnet":
		return newMainnetConfig(dataDir, network, rpcbind, rpcport, btcUrl, btcUser, btcPwd, btcNetwork, ethUrl, ethPrivateKey)
	case "testnet":
		return newTestConfig(dataDir, network, rpcbind, rpcport, btcUrl, btcUser, btcPwd, btcNetwork, ethUrl, ethPrivateKey)
	case "local":
		return newLocalConfig(dataDir, network, rpcbind, rpcport, btcUrl, btcUser, btcPwd, btcNetwork, ethUrl, ethPrivateKey)
	default:
		return config, fmt.Errorf("unsupport network now: %v", network)
	}

}

func newMainnetConfig(dataDir, testnet, rpcbind, rpcport, btcUrl, btcUser, btcPwd, btcNetwork, ethUrl, ethPrivateKey string) (NodeConfig, error) {
	multiSigPub1, err := hex.DecodeString(BtcMultiSigPublic1)
	if err != nil {
		logger.Error("hex decode string error", err)
		return NodeConfig{}, err
	}
	multiSigPub2, err := hex.DecodeString(BtcMultiSigPublic2)
	if err != nil {
		logger.Error("hex decode string error", err)
		return NodeConfig{}, err
	}
	multiSigPub3, err := hex.DecodeString(BtcMultiSigPublic3)
	if err != nil {
		logger.Error("hex decode string error", err)
		return NodeConfig{}, err
	}
	multiSigAddressInfo := MultiAddressInfo{
		PublicKeyList: [][]byte{
			multiSigPub1, multiSigPub2, multiSigPub3,
		},
		NRequired: BtcMultiNRequired,
	}
	return NodeConfig{
		DataDir: dataDir,
		Network: testnet,
		Rpcbind: rpcbind,
		RpcPort: rpcport,

		BtcUrl:           btcUrl,
		BtcUser:          btcUser,
		BtcPwd:           btcPwd,
		BtcNetwork:       btcNetwork,
		BtcScanBlockTime: BtcScanTime,
		BtcOperatorAddr:  BtcOperatorAddress,
		BtcPrivateKeys: []string{
			"b26dbaab82d9ebd8f37c88bbe56e22bf9cb21150c96dfb35ece4b787d3710d3301",
			"62dd5835dc2ce7f4f40eea1b88c816043d288532c8bb91964adef9bc0f0b4b7201",
			"9ff573d948c80fa1a50da6f66229b4bede9ec3fb482dd126f58d3acfb4b2979801",
		},
		BtcInitHeight: InitBitcoinHeight,

		EthInitHeight:    InitEthereumHeight,
		EthUrl:           ethUrl,
		ZkBridgeAddr:     EthZkBridgeAddress,
		ZkBtcAddr:        EthZkBtcAddress,
		EthScanBlockTime: EthScanTime,
		EthPrivateKey:    ethPrivateKey,
		LogAddr:          []string{EthZkBridgeAddress},
		LogTopic:         RedeemLogTopices,
		MultiAddressInfo: multiSigAddressInfo,
	}, nil
}

func newTestConfig(dataDir, testnet, rpcbind, rpcport, btcUrl, btcUser, btcPwd, btcNetwork, ethUrl, ethPrivateKey string) (NodeConfig, error) {
	multiSigPub1, err := hex.DecodeString(TestnetBtcMultiSigPublic1)
	if err != nil {
		logger.Error("hex decode string error", err)
		return NodeConfig{}, err
	}
	multiSigPub2, err := hex.DecodeString(TestnetBtcMultiSigPublic2)
	if err != nil {
		logger.Error("hex decode string error", err)
		return NodeConfig{}, err
	}
	multiSigPub3, err := hex.DecodeString(TestnetBtcMultiSigPublic3)
	if err != nil {
		logger.Error("hex decode string error", err)
		return NodeConfig{}, err
	}
	multiSigAddressInfo := MultiAddressInfo{
		PublicKeyList: [][]byte{
			multiSigPub1, multiSigPub2, multiSigPub3,
		},
		NRequired: TestnetBtcMultiNRequired,
	}
	return NodeConfig{
		DataDir:          dataDir,
		Network:          testnet,
		Rpcbind:          rpcbind,
		RpcPort:          rpcport,
		BtcUrl:           btcUrl,
		BtcUser:          btcUser,
		BtcPwd:           btcPwd,
		BtcNetwork:       btcNetwork,
		BtcScanBlockTime: TestnetBtcScanTime,
		BtcOperatorAddr:  TestnetBtcOperatorAddress,
		BtcPrivateKeys: []string{
			"b26dbaab82d9ebd8f37c88bbe56e22bf9cb21150c96dfb35ece4b787d3710d3301",
			"62dd5835dc2ce7f4f40eea1b88c816043d288532c8bb91964adef9bc0f0b4b7201",
			"9ff573d948c80fa1a50da6f66229b4bede9ec3fb482dd126f58d3acfb4b2979801",
		},
		BtcInitHeight:    TestnetInitBitcoinHeight,
		EthInitHeight:    TestnetInitEthereumHeight,
		EthUrl:           ethUrl,
		ZkBridgeAddr:     TestnetEthZkBridgeAddress,
		ZkBtcAddr:        TestnetEthZkBtcAddress,
		EthScanBlockTime: TestnetEthScanTime,
		EthPrivateKey:    ethPrivateKey,
		LogAddr:          []string{TestnetEthZkBridgeAddress},
		LogTopic:         TestnetRedeemLogTopices,
		MultiAddressInfo: multiSigAddressInfo,
	}, nil
}

func newLocalConfig(dataDir, testnet, rpcbind, rpcport, btcUrl, btcUser, btcPwd, btcNetwork, ethUrl, ethPrivateKey string) (NodeConfig, error) {
	multiSigPub1, err := hex.DecodeString(LocalBtcMultiSigPublic1)
	if err != nil {
		logger.Error("hex decode string error", err)
		return NodeConfig{}, err
	}
	multiSigPub2, err := hex.DecodeString(LocalBtcMultiSigPublic2)
	if err != nil {
		logger.Error("hex decode string error", err)
		return NodeConfig{}, err
	}
	multiSigPub3, err := hex.DecodeString(LocalBtcMultiSigPublic3)
	if err != nil {
		logger.Error("hex decode string error", err)
		return NodeConfig{}, err
	}
	multiSigAddressInfo := MultiAddressInfo{
		PublicKeyList: [][]byte{
			multiSigPub1, multiSigPub2, multiSigPub3,
		},
		NRequired: LocalBtcMultiNRequired,
	}
	return NodeConfig{
		DataDir:          dataDir,
		Network:          testnet,
		Rpcbind:          rpcbind,
		RpcPort:          rpcport,
		BtcUrl:           btcUrl,
		BtcUser:          btcUser,
		BtcPwd:           btcPwd,
		BtcNetwork:       btcNetwork,
		BtcScanBlockTime: LocalBtcScanTime,
		BtcOperatorAddr:  LocalBtcOperatorAddress,
		BtcPrivateKeys: []string{
			"b26dbaab82d9ebd8f37c88bbe56e22bf9cb21150c96dfb35ece4b787d3710d3301",
			"62dd5835dc2ce7f4f40eea1b88c816043d288532c8bb91964adef9bc0f0b4b7201",
			"9ff573d948c80fa1a50da6f66229b4bede9ec3fb482dd126f58d3acfb4b2979801",
		},
		BtcInitHeight:    LocalInitBitcoinHeight,
		EthInitHeight:    LocalInitEthereumHeight,
		EthUrl:           ethUrl,
		ZkBridgeAddr:     LocalEthZkBridgeAddress,
		ZkBtcAddr:        LocalEthZkBtcAddress,
		EthScanBlockTime: LocalEthScanTime,
		EthPrivateKey:    ethPrivateKey,
		LogAddr:          []string{LocalEthZkBridgeAddress},
		LogTopic:         LocalRedeemLogTopices,
		MultiAddressInfo: multiSigAddressInfo,
	}, nil
}

type MultiAddressInfo struct {
	PublicKeyList [][]byte
	NRequired     int
}

type WorkerConfig struct {
	ParallelNums int    `json:"parallelNums"`
	ProofUrl     string `json:"proofUrl"`
}

func TestnetDaemonConfig() NodeConfig {
	user, err := user.Current()
	config, err := NewNodeConfig(
		fmt.Sprintf("%v/.daemon", user.HomeDir),
		"testnet",
		"127.0.0.1",
		"8545",
		"https://go.getblock.io/d54c59f635654cc082de1f3fd14e5d02",
		"",
		"",
		"Testnet",
		"https://ethereum-holesky.publicnode.com",
		"c0781e4ca498e0ad693751bac014c0ab00c2841f28903e59cdfe1ab212438e49",
	)
	if err != nil {
		panic(err)
	}
	return config
}
func LocalDevDaemonConfig() NodeConfig {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	config, err := NewNodeConfig(
		fmt.Sprintf("%v/.daemon", user.HomeDir),
		"local",
		"127.0.0.1",
		"8545",
		"http://127.0.0.1:8332",
		"lightec",
		"Abcd1234",
		"Regtest",
		"https://ethereum-holesky.publicnode.com",
		"c0781e4ca498e0ad693751bac014c0ab00c2841f28903e59cdfe1ab212438e49",
	)
	if err != nil {
		panic(err)
	}
	return config
}
