package node

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/wire"
	"github.com/lightec-xyz/daemon/common"
	btctypes "github.com/lightec-xyz/daemon/rpc/bitcoin"
	"os"
	"time"
)

type RunConfig struct {
	Datadir        string `json:"datadir"`
	Rpcbind        string `json:"rpcbind"`
	Rpcport        string `json:"rpcport"`
	WsPort         string `json:"wsport"`
	Network        string `json:"network"`
	BtcUser        string `json:"btcUser"`
	BtcPwd         string `json:"btcPwd"`
	BtcToken       string `json:"btcToken"`
	BtcUrl         string `json:"btcUrl"`
	EthUrl         string `json:"ethUrl"`
	BeaconUrl      string `json:"beaconUrl"`
	OasisUrl       string `json:"oasisUrl"`
	SgxUrl         string `json:"sgxUrl"`
	DiscordHookUrl string `json:"discordHookUrl"`

	IcpWalletAddress string `json:"icpWalletAddress"`
	IcpPrivateKey    string `json:"icpPrivateKey"`

	MinerAddr          string        `json:"minerAddr"`
	BtcReScan          bool          `json:"btcReScan"`
	EthReScan          bool          `json:"ethReScan"`
	TxMode             common.TxMode `json:"txMode"`
	BeaconReScan       bool          `json:"beaconReScan"`
	EthPrivateKey      string        `json:"ethPrivateKey"`
	EnableLocalWorker  bool          `json:"enableLocalWorker"`
	BtcInitHeight      uint64        `json:"btcInitHeight"`
	EthInitHeight      uint64        `json:"ethInitHeight"`
	BeaconInitSlot     uint64        `json:"beaconInitSlot"`
	GenesisBeaconSlot  uint64        `json:"genesisBeaconSlot"`
	BtcGenesisHeight   uint64        `json:"btcGenesisHeight"`
	BtcCpBlockHeight   int64         `json:"btcCpBlockHeight"`
	DisableBtcAgent    bool          `json:"disableBtcAgent"`
	DisableEthAgent    bool          `json:"disableEthAgent"`
	DisableBeaconAgent bool          `json:"disableBeaconAgent"`
	BtcMainnetPath     string        `json:"btcMainnetPath"`

	DisableLipP2p bool     `json:"disableLipP2p"`
	BtcSetupDir   string   `json:"btcSetupDir"`
	EthSetupDir   string   `json:"ethSetupDir"`
	DisableFetch  bool     `json:"disableFetch"`
	P2pPort       int      `json:"p2pPort"`
	P2pBootstraps []string `json:"p2pBootstraps"`
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
	if rc.WsPort == "" {
		rc.WsPort = "9880"
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
	//if rc.MinerAddr == "" {
	//	return fmt.Errorf("minerAddr is empty")
	//}
	if rc.SgxUrl == "" {
		rc.SgxUrl = TestnetSgxServerUrl
	}
	if rc.OasisUrl == "" {
		return fmt.Errorf("oasisUrl is empty")
	}
	if rc.EthPrivateKey == "" {
		return fmt.Errorf("ethPrivateKey is empty")
	}
	return nil
}

type Config struct {
	RunConfig
	BtcOperatorAddr       string        `json:"btcOperatorAddr"`
	BtcLockScript         string        `json:"btcLockScript"`
	GenesisSyncPeriod     uint64        `json:"genesisPeriod"`
	ZkBridgeAddr          string        `json:"zkBridgeAddr"`
	ZkBtcAddr             string        `json:"zkBtcAddr"`
	UtxoManagerAddr       string        `json:"utxoManagerAddr"`
	OasisSignerAddress    string        `json:"oasisSignerAddress"`
	IcpTxSingerAddress    string        `json:"icpTxSingerAddress"`
	IcpBlockSignerAddress string        `json:"icpBlockSignerAddress"`
	BtcTxVerifyAddr       string        `json:"txVerifyAddr"`
	EthAddrFilter         *EthFilter    `json:"ethAddrFilter"`
	BtcFilter             *BtcFilter    `json:"btcFilter"`
	EthScanTime           time.Duration `json:"ethScanTime"`
	BtcScanTime           time.Duration `json:"btcScanTime"`
	Debug                 bool
}

func NewConfig(cfg RunConfig) (Config, error) {
	err := cfg.Check()
	if err != nil {
		return Config{}, err
	}
	switch cfg.Network {
	case LightecTestnet:
		return getTestnetConfig(cfg)
	default:
		return Config{}, fmt.Errorf("unsupport network now: %v", cfg.Network)
	}
}

func getTestnetConfig(cfg RunConfig) (Config, error) {
	if cfg.BtcInitHeight == 0 {
		cfg.BtcInitHeight = TestnetInitBitcoinHeight
	}
	if cfg.EthInitHeight == 0 {
		cfg.EthInitHeight = TestnetInitEthereumHeight
	}
	if cfg.BeaconInitSlot == 0 {
		cfg.BeaconInitSlot = TestnetInitBeaconHeight
	}
	return Config{
		RunConfig: cfg,

		IcpBlockSignerAddress: TestnetBlockSingerId,
		IcpTxSingerAddress:    TestnetIcpTxSingerId,
		GenesisSyncPeriod:     cfg.GenesisBeaconSlot / common.SlotPerPeriod,
		BtcOperatorAddr:       TestnetBtcOperatorAddress,
		BtcLockScript:         TestnetBtcLockScript,
		ZkBridgeAddr:          TestnetEthZkBridgeAddress,
		ZkBtcAddr:             TestnetEthZkBtcAddress,
		UtxoManagerAddr:       TestEthUtxoManagerAddress,
		BtcScanTime:           TestnetBtcScanTime,
		EthScanTime:           TestnetEthScanTime,
		BtcTxVerifyAddr:       TestEthBtcTxVerifyAddress,
		OasisSignerAddress:    TestnetOasisSignerAddr,
		BtcFilter:             NewBtcAddrFilter(TestnetBtcOperatorAddress, TestnetBtcLockScript, MinDepositValue, cfg.TxMode),
		EthAddrFilter: NewEthAddrFilter(TestnetBtcLockScript, TestEthUtxoManagerAddress, TestnetEthZkBridgeAddress, TestnetFeePoolAddr,
			TestnetDepositTopic, TestnetRedeemTopic, TestnetUpdateUtxoTopic, TestnetDepositRewardTopic, TestnetRedeemRewardTopic,
			cfg.TxMode),
		Debug: common.GetEnvDebugMode(),
	}, nil
}

type BtcFilter struct {
	OperatorAddr    string
	LockScript      string
	minDepositValue float64
	txMode          common.TxMode
}

func NewBtcAddrFilter(operator, lockScript string, minDepositValue float64, txMode common.TxMode) *BtcFilter {
	return &BtcFilter{
		OperatorAddr:    operator,
		LockScript:      lockScript,
		minDepositValue: minDepositValue,
		txMode:          txMode,
	}
}

func (b *BtcFilter) GetMinDepositValue() float64 {
	return b.minDepositValue
}

func (b *BtcFilter) Redeem(inputs []btctypes.TxVin) bool {
	if b.txMode != common.NormalTx {
		return false
	}
	for _, vin := range inputs {
		if vin.Prevout.ScriptPubKey.Address == b.OperatorAddr {
			return true
		}
	}
	return false
}

func (b *BtcFilter) Migrate(outputs []btctypes.TxVout) bool {
	var migrate bool
	for _, out := range outputs {
		if out.ScriptPubKey.Type == "nulldata" && common.StrEqual(out.ScriptPubKey.Hex, MigrateProto) {
			migrate = true
		}
	}
	return migrate && b.Deposit(outputs)
}

func (b *BtcFilter) Deposit(outputs []btctypes.TxVout) bool {
	for _, out := range outputs {
		if out.ScriptPubKey.Address == b.OperatorAddr {
			return true
		}
	}
	return false
}

type EthFilter struct {
	UtxoManagerAddr string `json:"depositAddr"`
	ZkbtcBridgeAddr string `json:"redeemAddr"`
	FeePoolAddr     string `json:"feePoolAddr"`

	DepositTxTopic     string `json:"depositTxTopic"`
	RedeemTxTopic      string `json:"redeemTxTopic"`
	UpdateUtxoTopic    string `json:"updateUtxoTopic"`
	DepositRewardTopic string `json:"depositRewardTopic"`
	RedeemRewardTopic  string `json:"redeemRewardTopic"`
	BtcLockScript      string `json:"btcLockScript"`
	txMode             common.TxMode
}

func (e *EthFilter) DepositTx(addr, topic string) bool {
	return common.StrEqual(e.UtxoManagerAddr, addr) && common.StrEqual(e.DepositTxTopic, topic)
}
func (e *EthFilter) RedeemTx(addr, topic string) bool {
	return common.StrEqual(e.ZkbtcBridgeAddr, addr) && common.StrEqual(e.RedeemTxTopic, topic)
}
func (e *EthFilter) DepositReward(addr, topic string) bool {
	return common.StrEqual(e.FeePoolAddr, addr) && common.StrEqual(e.DepositRewardTopic, topic)
}

func (e *EthFilter) RedeemReward(addr, topic string) bool {
	return common.StrEqual(e.FeePoolAddr, addr) && common.StrEqual(e.RedeemRewardTopic, topic)
}
func (e *EthFilter) UpdateUtxo(addr, topic string) bool {
	return common.StrEqual(e.UtxoManagerAddr, addr) && common.StrEqual(e.UpdateUtxoTopic, topic)
}

func (e *EthFilter) MigrateTx(outs []*wire.TxOut) bool {
	migrate := false
	for _, out := range outs {
		if common.StrEqual(hex.EncodeToString(out.PkScript), MigrateProto) {
			migrate = true
		}
	}
	if (migrate && e.txMode == common.OnlyMigrateTx) ||
		(!migrate && e.txMode == common.NormalTx) {
		return true
	}

	return false
}

func NewEthAddrFilter(btcLockScript, utxoManagerAddr, zkbtcBridgeAddr, feePoolAddr string, depositTxTopic, redeemTxTopic, updateUtxoTopic,
	depositRewardTopic, redeemRewardTopic string, txMode common.TxMode) *EthFilter {
	return &EthFilter{
		BtcLockScript:      btcLockScript,
		UtxoManagerAddr:    utxoManagerAddr,
		ZkbtcBridgeAddr:    zkbtcBridgeAddr,
		FeePoolAddr:        feePoolAddr,
		DepositTxTopic:     depositTxTopic,
		RedeemTxTopic:      redeemTxTopic,
		UpdateUtxoTopic:    updateUtxoTopic,
		DepositRewardTopic: depositRewardTopic,
		RedeemRewardTopic:  redeemRewardTopic,
		txMode:             txMode,
	}
}

func (e *EthFilter) FilterLogs() (logFilters []string, topicFilters []string) {
	logFilters = []string{e.UtxoManagerAddr, e.ZkbtcBridgeAddr, e.FeePoolAddr}
	topicFilters = []string{e.DepositTxTopic, e.RedeemTxTopic, e.UpdateUtxoTopic, e.DepositRewardTopic, e.RedeemRewardTopic}
	return logFilters, topicFilters
}
