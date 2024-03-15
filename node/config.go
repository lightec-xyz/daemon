package node

import (
	"encoding/hex"
	"fmt"
	"github.com/lightec-xyz/daemon/logger"
	beaconconfig "github.com/prysmaticlabs/prysm/v5/config/params"
	"os/user"
	"time"
)

type NodeConfig struct {
	DataDir           string `json:"datadir"`
	Network           string `json:"network"`
	Rpcbind           string `json:"rpcbind"`
	RpcPort           string `json:"rpcport"`
	EnableLocalWorker bool   `json:"enableLocalWorker"`

	BtcUrl            string           `json:"btcUrl"`
	BtcUser           string           `json:"btcUser"`
	BtcPwd            string           `json:"btcPwd"`
	BtcNetwork        string           `json:"btcNetwork"`
	BtcScanBlockTime  time.Duration    `json:"btcBlockTime"`
	BtcOperatorAddr   string           `json:"btcOperatorAddr"`
	BtcPrivateKeys    []string         `json:"btcPrivateKeys"`
	BtcWhiteList      []string         `json:"btcWhiteList"`
	BtcInitHeight     int64            `json:"btcInitHeight"`
	MultiAddressInfo  MultiAddressInfo `json:"multiAddressInfo"`
	GenesisSyncPeriod uint64           `json:"genesisSyncPeriod"`
	AutoSubmit        bool             `json:"autoSubmit"`

	//Beacon config
	BeaconInitHeight uint64                          `json:"beaconInitHeight"`
	BeaconUrl        string                          `json:"beaconUrl"`
	BeaconConfig     *beaconconfig.BeaconChainConfig `json:"beaconConfig"`

	//Eth1 config
	EthInitHeight    int64          `json:"ethInitHeight"`
	EthWhiteList     []string       `json:"ethWhiteList"`
	EthUrl           string         `json:"ethUrl"` //eth1 url
	ZkBridgeAddr     string         `json:"zkBridgeAddr"`
	ZkBtcAddr        string         `json:"zkBtcAddr"`
	EthScanBlockTime time.Duration  `json:"ethBlockTime"`
	EthPrivateKey    string         `json:"ethPrivateKey"`
	LogAddr          []string       `json:"logAddr"`
	LogTopic         []string       `json:"logTopic"`
	Workers          []WorkerConfig `json:"workers"`
	EthAddrFilter    EthAddrFilter
}

type EthAddrFilter struct {
	LogDepositAddr      string
	LogRedeemAddr       string
	LogTopicDepositAddr string
	LogTopicRedeemAddr  string
}

func (f *EthAddrFilter) FilterLogs() (logFilters []string, topicFilters []string) {
	logFilters = []string{f.LogDepositAddr, f.LogRedeemAddr}
	topicFilters = []string{f.LogTopicDepositAddr, f.LogTopicRedeemAddr}
	return logFilters, topicFilters
}

func NewNodeConfig(enableLocalWorker, autoSubmit bool, dataDir, network, rpcbind, rpcport, btcUrl, btcUser, btcPwd, beaconUrl, ethUrl, ethPrivateKey string) (NodeConfig, error) {
	var config NodeConfig
	if network == "" {
		network = LightecMainnet
	}
	logger.Info("current network: %v", network)
	if dataDir == "" {
		current, err := user.Current()
		if err != nil {
			logger.Error("get current user error: %v", err)
			return config, err
		}
		dataDir = fmt.Sprintf("%v/.daemon", current.HomeDir)
		logger.Info("datadir: %v", dataDir)
	}
	if rpcbind == "" {
		rpcbind = "127.0.0.1"
	}
	if rpcport == "" {
		rpcport = "9780"
	}
	if btcUrl == "" {
		return config, fmt.Errorf("btcUrl is empty")
	}
	if beaconUrl == "" {
		return config, fmt.Errorf("beaconUrl is empty")
	}
	if ethUrl == "" {
		return config, fmt.Errorf("ethUrl is empty")
	}

	if autoSubmit {
		if ethPrivateKey == "" {
			return config, fmt.Errorf("ethPrivateKey is empty")
		}
	}

	switch network {
	case LightecMainnet:
		return newMainnetConfig(enableLocalWorker, autoSubmit, dataDir, network, rpcbind, rpcport, btcUrl, btcUser, btcPwd, beaconUrl, ethUrl, ethPrivateKey, beaconconfig.MainnetConfig())
	case LightecTestnet:
		return newTestConfig(enableLocalWorker, autoSubmit, dataDir, network, rpcbind, rpcport, btcUrl, btcUser, btcPwd, beaconUrl, ethUrl, ethPrivateKey, beaconconfig.HoleskyConfig())
	case Lighteclocal:
		return newLocalConfig(enableLocalWorker, autoSubmit, dataDir, network, rpcbind, rpcport, btcUrl, btcUser, btcPwd, beaconUrl, ethUrl, ethPrivateKey, beaconconfig.HoleskyConfig())
	default:
		return config, fmt.Errorf("unsupport network now: %v", network)
	}

}

func newMainnetConfig(enableLocalWorker, autoSubmit bool, dataDir, testnet, rpcbind, rpcport, btcUrl, btcUser, btcPwd, beaconUrl, ethUrl, ethPrivateKey string, beaconConfig *beaconconfig.BeaconChainConfig) (NodeConfig, error) {
	multiSigPub1, err := hex.DecodeString(BtcMultiSigPublic1)
	if err != nil {
		logger.Error("hex decode string error: %v", err)
		return NodeConfig{}, err
	}
	multiSigPub2, err := hex.DecodeString(BtcMultiSigPublic2)
	if err != nil {
		logger.Error("hex decode string error:%v", err)
		return NodeConfig{}, err
	}
	multiSigPub3, err := hex.DecodeString(BtcMultiSigPublic3)
	if err != nil {
		logger.Error("hex decode string error:%v", err)
		return NodeConfig{}, err
	}
	multiSigAddressInfo := MultiAddressInfo{
		PublicKeyList: [][]byte{
			multiSigPub1, multiSigPub2, multiSigPub3,
		},
		NRequired: BtcMultiNRequired,
	}
	return NodeConfig{
		DataDir:           dataDir,
		Network:           testnet,
		Rpcbind:           rpcbind,
		RpcPort:           rpcport,
		EnableLocalWorker: enableLocalWorker,

		BtcUrl:           btcUrl,
		BtcUser:          btcUser,
		BtcPwd:           btcPwd,
		BtcNetwork:       string(BtcNetwork),
		BtcScanBlockTime: BtcScanTime,
		BtcOperatorAddr:  BtcOperatorAddress,
		BtcPrivateKeys: []string{
			"b26dbaab82d9ebd8f37c88bbe56e22bf9cb21150c96dfb35ece4b787d3710d3301",
			"62dd5835dc2ce7f4f40eea1b88c816043d288532c8bb91964adef9bc0f0b4b7201",
			"9ff573d948c80fa1a50da6f66229b4bede9ec3fb482dd126f58d3acfb4b2979801",
		},
		BtcInitHeight:    InitBitcoinHeight,
		BeaconInitHeight: InitBeaconHeight,
		BeaconUrl:        beaconUrl,
		BeaconConfig:     beaconConfig,
		EthInitHeight:    InitEthereumHeight,
		EthUrl:           ethUrl,
		ZkBridgeAddr:     EthZkBridgeAddress,
		ZkBtcAddr:        EthZkBtcAddress,
		EthScanBlockTime: EthScanTime,
		EthPrivateKey:    ethPrivateKey,
		LogAddr:          LogAddrs,
		LogTopic:         LogTopics,
		AutoSubmit:       autoSubmit,
		MultiAddressInfo: multiSigAddressInfo,
		EthAddrFilter: EthAddrFilter{
			LogDepositAddr:      LogDepositAddr,
			LogRedeemAddr:       LogRedeemAddr,
			LogTopicDepositAddr: TopicDepositAddr,
			LogTopicRedeemAddr:  TopicRedeemAddr,
		},
	}, nil
}

func newTestConfig(enableLocalWorker, autoSubmit bool, dataDir, testnet, rpcbind, rpcport, btcUrl, btcUser, btcPwd, beaconUrl, ethUrl, ethPrivateKey string, beaconConfig *beaconconfig.BeaconChainConfig) (NodeConfig, error) {
	multiSigPub1, err := hex.DecodeString(TestnetBtcMultiSigPublic1)
	if err != nil {
		logger.Error("hex decode string error:%v", err)
		return NodeConfig{}, err
	}
	multiSigPub2, err := hex.DecodeString(TestnetBtcMultiSigPublic2)
	if err != nil {
		logger.Error("hex decode string error:%v", err)
		return NodeConfig{}, err
	}
	multiSigPub3, err := hex.DecodeString(TestnetBtcMultiSigPublic3)
	if err != nil {
		logger.Error("hex decode string error:%v", err)
		return NodeConfig{}, err
	}
	multiSigAddressInfo := MultiAddressInfo{
		PublicKeyList: [][]byte{
			multiSigPub1, multiSigPub2, multiSigPub3,
		},
		NRequired: TestnetBtcMultiNRequired,
	}
	return NodeConfig{
		DataDir:           dataDir,
		Network:           testnet,
		Rpcbind:           rpcbind,
		RpcPort:           rpcport,
		EnableLocalWorker: enableLocalWorker,

		BtcUrl:           btcUrl,
		BtcUser:          btcUser,
		BtcPwd:           btcPwd,
		BtcNetwork:       string(TestnetBtcNetwork),
		BtcScanBlockTime: TestnetBtcScanTime,
		BtcOperatorAddr:  TestnetBtcOperatorAddress,
		BtcPrivateKeys: []string{
			"b26dbaab82d9ebd8f37c88bbe56e22bf9cb21150c96dfb35ece4b787d3710d3301",
			"62dd5835dc2ce7f4f40eea1b88c816043d288532c8bb91964adef9bc0f0b4b7201",
			"9ff573d948c80fa1a50da6f66229b4bede9ec3fb482dd126f58d3acfb4b2979801",
		},
		BtcInitHeight: TestnetInitBitcoinHeight,
		AutoSubmit:    autoSubmit,

		//BeaconConfig
		BeaconInitHeight: TestnetInitBeaconHeight,
		BeaconUrl:        beaconUrl,
		BeaconConfig:     beaconConfig,

		EthInitHeight:    TestnetInitEthereumHeight,
		EthUrl:           ethUrl,
		ZkBridgeAddr:     TestnetEthZkBridgeAddress,
		ZkBtcAddr:        TestnetEthZkBtcAddress,
		EthScanBlockTime: TestnetEthScanTime,
		EthPrivateKey:    ethPrivateKey,
		LogAddr:          TestLogAddrs,
		LogTopic:         TestLogTopics,
		MultiAddressInfo: multiSigAddressInfo,
		EthAddrFilter: EthAddrFilter{
			LogDepositAddr:      TestLogDepositAddr,
			LogRedeemAddr:       TestLogRedeemAddr,
			LogTopicDepositAddr: TestTopicDepositAddr,
			LogTopicRedeemAddr:  TestTopicRedeemAddr,
		},
	}, nil
}

func newLocalConfig(enableLocalWorker, autoSubmit bool, dataDir, testnet, rpcbind, rpcport, btcUrl, btcUser, btcPwd, beaconUrl, ethUrl, ethPrivateKey string, beaconConfig *beaconconfig.BeaconChainConfig) (NodeConfig, error) {
	multiSigPub1, err := hex.DecodeString(LocalBtcMultiSigPublic1)
	if err != nil {
		logger.Error("hex decode string error: %v", err)
		return NodeConfig{}, err
	}
	multiSigPub2, err := hex.DecodeString(LocalBtcMultiSigPublic2)
	if err != nil {
		logger.Error("hex decode string error: %v", err)
		return NodeConfig{}, err
	}
	multiSigPub3, err := hex.DecodeString(LocalBtcMultiSigPublic3)
	if err != nil {
		logger.Error("hex decode string error: %v", err)
		return NodeConfig{}, err
	}
	multiSigAddressInfo := MultiAddressInfo{
		PublicKeyList: [][]byte{
			multiSigPub1, multiSigPub2, multiSigPub3,
		},
		NRequired: LocalBtcMultiNRequired,
	}
	return NodeConfig{
		DataDir:           dataDir,
		Network:           testnet,
		Rpcbind:           rpcbind,
		RpcPort:           rpcport,
		EnableLocalWorker: enableLocalWorker,
		BtcUrl:            btcUrl,
		BtcUser:           btcUser,
		BtcPwd:            btcPwd,
		BtcNetwork:        string(LocalBtcNetwork),
		BtcScanBlockTime:  LocalBtcScanTime,
		BtcOperatorAddr:   LocalBtcOperatorAddress,
		BtcPrivateKeys: []string{
			"b26dbaab82d9ebd8f37c88bbe56e22bf9cb21150c96dfb35ece4b787d3710d3301",
			"62dd5835dc2ce7f4f40eea1b88c816043d288532c8bb91964adef9bc0f0b4b7201",
			"9ff573d948c80fa1a50da6f66229b4bede9ec3fb482dd126f58d3acfb4b2979801",
		},
		BtcInitHeight:    LocalInitBitcoinHeight,
		AutoSubmit:       autoSubmit,
		BeaconInitHeight: LocalInitBeaconHeight,
		BeaconUrl:        beaconUrl,
		BeaconConfig:     beaconConfig,
		EthInitHeight:    LocalInitEthereumHeight,
		EthUrl:           ethUrl,
		ZkBridgeAddr:     LocalEthZkBridgeAddress,
		ZkBtcAddr:        LocalEthZkBtcAddress,
		EthScanBlockTime: LocalEthScanTime,
		EthPrivateKey:    ethPrivateKey,
		LogAddr:          LocalLogAddrs,
		LogTopic:         LocalLogTopics,
		MultiAddressInfo: multiSigAddressInfo,
		EthAddrFilter: EthAddrFilter{
			LogDepositAddr:      LocalLogDepositAddr,
			LogRedeemAddr:       LocalLogRedeemAddr,
			LogTopicDepositAddr: LocalTopicDepositAddr,
			LogTopicRedeemAddr:  LocalTopicRedeemAddr,
		},
	}, nil
}

type MultiAddressInfo struct {
	PublicKeyList [][]byte
	NRequired     int
}

type WorkerConfig struct {
	MaxNums int    `json:"maxNums"`
	Url     string `json:"url"`
}

func TestnetDaemonConfig() NodeConfig {
	user, err := user.Current()
	config, err := NewNodeConfig(
		true,
		true,
		fmt.Sprintf("%v/.daemon", user.HomeDir),
		"testnet",
		"127.0.0.1",
		"9870",
		"https://go.getblock.io/d54c59f635654cc082de1f3fd14e5d02",
		"lightec",
		"abcd1234",
		"http://127.0.0.1:8970",
		"https://go.getblock.io/0d372517498b419a97613e2bbf882a30",
		"c0781e4ca498e0ad693751bac014c0ab00c2841f28903e59cdfe1ab212438e49",
	)
	if err != nil {
		panic(err)
	}
	return config
}
func LocalDevDaemonConfig() NodeConfig {
	//user, err := user.Current()
	//if err != nil {
	//	panic(err)
	//}
	//dataDir := fmt.Sprintf("%v/.daemon", user.HomeDir)
	dataDir := "/Users/red/lworkspace/lightec/daemon/node/test"
	config, err := NewNodeConfig(
		true,
		false,
		dataDir,
		"local",
		"127.0.0.1",
		"9870",
		"http://127.0.0.1:8332",
		"lightec",
		"abcd1234",
		"http://127.0.0.1:8970",
		"https://go.getblock.io/0d372517498b419a97613e2bbf882a30",
		"c0781e4ca498e0ad693751bac014c0ab00c2841f28903e59cdfe1ab212438e49",
	)
	if err != nil {
		panic(err)
	}
	return config
}
