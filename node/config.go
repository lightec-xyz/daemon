package node

import (
	"encoding/hex"
	"fmt"
	"os"
	"time"
)

type RunConfig struct {
	Datadir               string `json:"datadir"`
	Rpcbind               string `json:"rpcbind"`
	Rpcport               string `json:"rpcport"`
	Network               string `json:"network"`
	BtcUser               string `json:"btcUser"`
	BtcPwd                string `json:"btcPwd"`
	BtcUrl                string `json:"btcUrl"`
	EthUrl                string `json:"ethUrl"`
	BeaconUrl             string `json:"beaconUrl"`
	EthPrivateKey         string `json:"ethPrivateKey"`
	EnableLocalWorker     bool   `json:"enableLocalWorker"`
	DisableRecursiveAgent bool   `json:"disableRecursiveAgent"`
	DisableTxAgent        bool   `json:"disableTxAgent"`
	BtcInitHeight         int64  `json:"btcInitHeight"`
	EthInitHeight         int64  `json:"ethInitHeight"`
	BeaconInitSlot        uint64 `json:"beaconInitSlot"`
}

func (rc *RunConfig) Check() error {
	if rc.Datadir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		rc.Datadir = fmt.Sprintf("%s/.lightec", homeDir)
	}
	if rc.Rpcbind == "" {
		rc.Rpcbind = "127.0.0.1"
	}
	if rc.Rpcport == "" {
		rc.Rpcport = "9870"
	}
	if rc.Network == "" {
		rc.Network = "testnet" // todo
	}
	if rc.BtcUrl == "" {
		return fmt.Errorf("btcUrl is empty")
	}
	if rc.EthUrl == "" {
		return fmt.Errorf("ethUrl is empty")
	}
	if rc.BeaconUrl == "" {
		return fmt.Errorf("beaconUrl is empty")
	}
	if rc.EthPrivateKey == "" {
		return fmt.Errorf("ethPrivateKey is empty")
	}
	return nil
}

type Config struct {
	RunConfig
	BtcOperatorAddr   string           `json:"btcOperatorAddr"`
	MultiAddressInfo  MultiAddressInfo `json:"multiAddressInfo"`
	GenesisSyncPeriod uint64           `json:"genesisPeriod"`
	ZkBridgeAddr      string           `json:"zkBridgeAddr"`
	ZkBtcAddr         string           `json:"zkBtcAddr"`
	EthAddrFilter     EthAddrFilter    `json:"ethAddrFilter"`
	EthScanTime       time.Duration    `json:"ethScanTime"`
	BtcScanTime       time.Duration    `json:"btcScanTime"`
}

func NewConfig(cfg RunConfig) (Config, error) {
	err := cfg.Check()
	if err != nil {
		return Config{}, err
	}
	switch cfg.Network {
	case LightecTestnet:
		return getTestnetConfig(cfg)
	case Lighteclocal:
		return getLocalConfig(cfg)
	default:
		return Config{}, fmt.Errorf("unsupport network now: %v", cfg.Network)
	}
}

func getLocalConfig(runCnfg RunConfig) (Config, error) {
	multiAddressInfo, err := NewMultiAddressInfo([]string{LocalBtcMultiSigPublic1, LocalBtcMultiSigPublic2,
		LocalBtcMultiSigPublic3}, LocalBtcMultiNRequired)
	if err != nil {
		return Config{}, err
	}
	if runCnfg.BtcInitHeight == 0 {
		runCnfg.BtcInitHeight = LocalInitBitcoinHeight
	}
	if runCnfg.EthInitHeight == 0 {
		runCnfg.EthInitHeight = LocalInitEthereumHeight
	}
	if runCnfg.BeaconInitSlot == 0 {
		runCnfg.BeaconInitSlot = LocalInitBeaconHeight
	}
	return Config{
		RunConfig:         runCnfg,
		GenesisSyncPeriod: runCnfg.BeaconInitSlot / 8192, // todo
		BtcOperatorAddr:   LocalBtcOperatorAddress,
		MultiAddressInfo:  multiAddressInfo,
		ZkBridgeAddr:      LocalEthZkBridgeAddress,
		ZkBtcAddr:         LocalEthZkBtcAddress,
		BtcScanTime:       LocalBtcScanTime,
		EthScanTime:       LocalEthScanTime,
		EthAddrFilter:     NewEthAddrFilter(LocalLogDepositAddr, LocalLogRedeemAddr, LocalTopicDepositAddr, LocalTopicRedeemAddr),
	}, nil
}

func getTestnetConfig(runCnfg RunConfig) (Config, error) {
	multiAddressInfo, err := NewMultiAddressInfo([]string{TestnetBtcMultiSigPublic1, TestnetBtcMultiSigPublic2,
		TestnetBtcMultiSigPublic3}, TestnetBtcMultiNRequired)
	if err != nil {
		return Config{}, err
	}
	if runCnfg.BtcInitHeight == 0 {
		runCnfg.BtcInitHeight = TestnetInitBitcoinHeight
	}
	if runCnfg.EthInitHeight == 0 {
		runCnfg.EthInitHeight = TestnetInitEthereumHeight
	}
	if runCnfg.BeaconInitSlot == 0 {
		runCnfg.BeaconInitSlot = TestnetInitBeaconHeight
	}
	return Config{
		RunConfig:         runCnfg,
		GenesisSyncPeriod: runCnfg.BeaconInitSlot / 8192, // todo
		BtcOperatorAddr:   TestnetBtcOperatorAddress,
		MultiAddressInfo:  multiAddressInfo,
		ZkBridgeAddr:      TestnetEthZkBridgeAddress,
		ZkBtcAddr:         TestnetEthZkBtcAddress,
		BtcScanTime:       TestnetBtcScanTime,
		EthScanTime:       TestnetEthScanTime,
		EthAddrFilter:     NewEthAddrFilter(TestLogDepositAddr, TestLogRedeemAddr, TestTopicDepositAddr, TestTopicRedeemAddr),
	}, nil
}

type MultiAddressInfo struct {
	PublicKeyList [][]byte `json:"publicKeyList"`
	NRequired     int      `json:"nRequired"`
}

type WorkerConfig struct {
	MaxNums int    `json:"maxNums"`
	Url     string `json:"url"`
	DataDir string `json:"dataDir"`
}

func NewMultiAddressInfo(publicKeys []string, nRequired int) (MultiAddressInfo, error) {
	var publicKeyListByte [][]byte
	for _, pub := range publicKeys {
		pubBytes, err := hex.DecodeString(pub)
		if err != nil {
			return MultiAddressInfo{}, err
		}
		publicKeyListByte = append(publicKeyListByte, pubBytes)
	}
	return MultiAddressInfo{
		PublicKeyList: publicKeyListByte,
		NRequired:     nRequired,
	}, nil
}

type EthAddrFilter struct {
	LogDepositAddr      string `json:"logDepositAddr"`
	LogRedeemAddr       string `json:"logRedeemAddr"`
	LogTopicDepositAddr string `json:"logTopicDepositAddr"`
	LogTopicRedeemAddr  string `json:"logTopicRedeemAddr"`
}

func NewEthAddrFilter(depositAddr, redeemAddr, topicDepositAddr, topicRedeemAddr string) EthAddrFilter {
	return EthAddrFilter{
		LogDepositAddr:      depositAddr,
		LogRedeemAddr:       redeemAddr,
		LogTopicDepositAddr: topicDepositAddr,
		LogTopicRedeemAddr:  topicRedeemAddr,
	}
}

func (f *EthAddrFilter) FilterLogs() (logFilters []string, topicFilters []string) {
	logFilters = []string{f.LogDepositAddr, f.LogRedeemAddr}
	topicFilters = []string{f.LogTopicDepositAddr, f.LogTopicRedeemAddr}
	return logFilters, topicFilters
}
