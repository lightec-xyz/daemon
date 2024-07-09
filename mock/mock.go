package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/node"
	btcrpc "github.com/lightec-xyz/daemon/rpc/bitcoin"
	"github.com/lightec-xyz/daemon/rpc/bitcoin/types"
	ethrpc "github.com/lightec-xyz/daemon/rpc/ethereum"
	btctx "github.com/lightec-xyz/daemon/transaction/bitcoin"
)

type Mock struct {
	cfg       Config
	btcClient *btcrpc.Client
	ethClient *ethrpc.Client
}

func (m *Mock) DepositBtc(value int64) error {
	utxoSet, err := m.btcClient.Scantxoutset(m.cfg.BtcDepositAddr)
	if err != nil {
		logger.Error("scan btc utxo error:%v", err)
		return err
	}
	if len(utxoSet.Unspents) == 0 {
		logger.Error("no utxo found")
		return err
	}

	var utxo types.Unspents
	for _, item := range utxoSet.Unspents {
		if item.Amount > 0.0001 {
			utxo = item
			break
		}
	}
	if utxo.Amount == 0 {
		logger.Error("no utxo found")
		return err
	}

	secretBytes, err := hex.DecodeString(m.cfg.BtcPrivateKey)
	if err != nil {
		logger.Error("decode private key error:%v", err)
		return err
	}
	ethAddrBytes, err := hex.DecodeString(m.cfg.EthAddr)
	if err != nil {
		logger.Error("decode eth addr error:%v", err)
		return err
	}
	var btcTxIn []btctx.TxIn
	inputValue := floatToInt(utxo.Amount)
	btcTxIn = append(btcTxIn, btctx.TxIn{
		Hash:     utxo.Txid,
		VOut:     uint32(utxo.Vout),
		Amount:   inputValue,
		PkScript: utxo.ScriptPubKey,
	})
	networkInfo, err := m.btcClient.GetNetworkInfo()
	if err != nil {
		logger.Error("get btc network info error:%v", err)
		return err
	}

	//todo
	outpuValue := inputValue - floatToInt(networkInfo.Relayfee) - 50
	fmt.Println(outpuValue)
	var btcTxOuts []btctx.TxOut
	btcTxOuts = append(btcTxOuts, btctx.TxOut{
		Amount:  outpuValue,
		Address: m.cfg.BtcOperatorAddr,
	})
	//btcTxOuts = append(btcTxOuts, btctx.TxOut{
	//	Amount:  outpuValue,
	//	EthAddress: m.cfg.BtcDepositAddr,
	//})

	transaction, err := btctx.CreateDepositTransaction(secretBytes, ethAddrBytes, btcTxIn, btcTxOuts, m.cfg.Network)
	if err != nil {
		logger.Error("create btc tx error:%v", err)
	}
	txHexBytes := hex.EncodeToString(transaction)
	logger.Info("send btc tx:%v", txHexBytes)
	txHash, err := m.btcClient.Sendrawtransaction(txHexBytes)
	if err != nil {
		logger.Error("send btc tx error:%v", err)
		return err
	}
	logger.Info("success send btc tx hash:%v index: %v", txHash, 1)

	return nil
}

func (m *Mock) DepositBtcToEth(txId, receiverAddr string, index uint32, amount *big.Int) error {
	nonce, err := m.ethClient.GetNonce(m.cfg.EthAddr)
	if err != nil {
		logger.Error("get nonce error:%v", err)
		return err
	}
	chainID, err := m.ethClient.ChainID(context.Background())
	if err != nil {
		logger.Error("get chain id error:%v", err)
		return err
	}
	gasPrice, err := m.ethClient.GetGasPrice()
	if err != nil {
		logger.Error("get gas price error:%v", err)
		return err
	}
	gasPrice = big.NewInt(0).Mul(gasPrice, big.NewInt(3))
	gasLimit := uint64(500000)
	ethTxHash, err := m.ethClient.Deposit(m.cfg.EthPrivateKey, txId, receiverAddr, index, nonce, gasLimit, chainID, gasPrice,
		amount, common.ZkProof([]byte("test proof")))
	if err != nil {
		logger.Error(" deposit eth error:%v", err)
		return err
	}
	logger.Info("success deposit eth tx hash:%v", ethTxHash)
	return nil
}

func (m *Mock) RedeemTx(amount int64) error {
	redeemAmount := big.NewInt(amount)
	minerFee := big.NewInt(300)
	fromAddr := m.cfg.EthAddr
	//balance, err := m.ethClient.GetZkBtcBalance(fromAddr)
	//if err != nil {
	//	logger.Error("get balance error:%v", err)
	//	return err
	//}
	//if balance.Cmp(big.NewInt(10000)) < 0 {
	//	logger.Error("balance less than 10000")
	//	return err
	//}
	payToAddrScript, err := btctx.GenPayToAddrScript(m.cfg.BtcDepositAddr, m.cfg.Network)
	if err != nil {
		logger.Error("gen pay to addr script error:%v", err)
		return err
	}
	redeemLockScript, err := hex.DecodeString(payToAddrScript)
	if err != nil {
		logger.Error("decode redeem lock script error:%v", err)
		return err
	}
	from := ethcommon.HexToAddress(fromAddr)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	gasLimit := 500000
	gasPrice, err := m.ethClient.SuggestGasPrice(ctx)
	if err != nil {
		logger.Error("get gas price error:%v", err)
		return err
	}
	gasPrice = big.NewInt(0).Mul(big.NewInt(2), gasPrice)
	chainID, err := m.ethClient.ChainID(ctx)
	if err != nil {
		logger.Error("get chain id error:%v", err)
		return err
	}
	nonce, err := m.ethClient.NonceAt(ctx, from, nil)
	if err != nil {
		logger.Error("get nonce error:%v", err)
		return err
	}
	txhash, err := m.ethClient.Redeem(m.cfg.EthPrivateKey, uint64(gasLimit), chainID, big.NewInt(int64(nonce)), gasPrice, redeemAmount, minerFee, redeemLockScript)
	if err != nil {
		logger.Error("redeem error:%v", err)
		return err
	}
	logger.Info("success redeem hash:%v", txhash)

	return nil
}

func (m *Mock) MergeBtcTx() error {
	utxoSet, err := m.btcClient.Scantxoutset(m.cfg.BtcDepositAddr)
	if err != nil {
		logger.Error("scan btc utxo error:%v", err)
		return err
	}
	if len(utxoSet.Unspents) == 0 {
		logger.Error("no utxo found")
		return err
	}
	secretBytes, err := hex.DecodeString(m.cfg.BtcPrivateKey)
	if err != nil {
		logger.Error("decode private key error:%v", err)
		return err
	}
	var btcTxIn []btctx.TxIn
	total := 0.0
	for _, utxo := range utxoSet.Unspents {
		logger.Info("utxo: txid:%v, index:%v, amount:%0.8f ,scriptPubKey:%v", utxo.Txid, utxo.Vout, utxo.Amount, utxo.ScriptPubKey)
		btcTxIn = append(btcTxIn, btctx.TxIn{
			Hash:     utxo.Txid,
			VOut:     uint32(utxo.Vout),
			Amount:   floatToInt(utxo.Amount),
			PkScript: utxo.ScriptPubKey,
		})
		total = total + utxo.Amount
	}
	logger.Info("total: %0.8f", total)
	networkInfo, err := m.btcClient.GetNetworkInfo()
	if err != nil {
		logger.Error("get btc network info error:%v", err)
		return err
	}

	//todo
	outpuValue := total - networkInfo.Relayfee - 0.00000150
	var btcTxOuts []btctx.TxOut
	btcTxOuts = append(btcTxOuts, btctx.TxOut{
		Amount:  floatToInt(outpuValue),
		Address: m.cfg.BtcDepositAddr,
	})
	logger.Info("outpu: value:%v", outpuValue)
	if outpuValue <= 0 {
		logger.Error("output value less than 0")
		return fmt.Errorf("output value less than 0")
	}

	transaction, err := btctx.CreateTransaction(secretBytes, btcTxIn, btcTxOuts, m.cfg.Network)
	if err != nil {
		logger.Error("create btc tx error:%v", err)
		return err
	}
	txHexBytes := hex.EncodeToString(transaction)
	logger.Info("send btc tx:%v", txHexBytes)
	txHash, err := m.btcClient.Sendrawtransaction(txHexBytes)
	if err != nil {
		logger.Error("send btc tx error:%v", err)
		return err
	}
	logger.Info("success send btc tx hash:%v", txHash)
	return nil
}

func NewMock(network string) (*Mock, error) {
	cfg, err := NewConfig(network)
	if err != nil {
		return nil, err
	}
	btcClient, err := btcrpc.NewClient(cfg.BtcUrl, cfg.BtcUser, cfg.BtcPwd, cfg.BtcNetwork)
	if err != nil {
		return nil, err
	}
	ethClient, err := ethrpc.NewClient(cfg.EthUrl, cfg.ZkBridgeAddr, cfg.ZkBtcAddr)
	if err != nil {
		return nil, err
	}
	mock := Mock{
		cfg:       cfg,
		btcClient: btcClient,
		ethClient: ethClient,
	}
	return &mock, nil
}

func NewConfig(network string) (Config, error) {
	var cfg Config
	if network == "testnet" {
		cfg = NewTestnetConfig()
	} else {
		cfg = NewDevConfig()
	}
	return cfg, nil
}

func NewDevConfig() Config {
	config := node.LocalDevDaemonConfig()
	return Config{
		NodeConfig:     config,
		Network:        btctx.RegTest,
		BtcPrivateKey:  "3c5579d538347d56ed5ef6d56e7e36ae453dcbd6ff6586783f82c17c8190d71601",
		BtcDepositAddr: "bcrt1q6lawf77u30mvs6sgcuthchgxdqm4f6n3kvx4z5",
		EthPrivateKey:  "c0781e4ca498e0ad693751bac014c0ab00c2841f28903e59cdfe1ab212438e49",
		EthAddr:        "771815eFD58e8D6e66773DB0bc002899c00d5b0c",
	}
}

func NewTestnetConfig() Config {
	config := node.TestnetDaemonConfig()
	return Config{
		NodeConfig:     config,
		Network:        btctx.TestNet,
		BtcPrivateKey:  "3c5579d538347d56ed5ef6d56e7e36ae453dcbd6ff6586783f82c17c8190d71601",
		BtcDepositAddr: "tb1q6lawf77u30mvs6sgcuthchgxdqm4f6n359lc4a",
		EthPrivateKey:  "c0781e4ca498e0ad693751bac014c0ab00c2841f28903e59cdfe1ab212438e49",
		EthAddr:        "771815eFD58e8D6e66773DB0bc002899c00d5b0c",
	}
}

type Config struct {
	BtcDepositAddr string
	BtcPrivateKey  string
	EthPrivateKey  string
	EthAddr        string
	Network        btctx.NetWork
	node.Config
}

func floatToInt(value float64) int64 {
	valueRat := NewRat().Mul(NewRat().SetFloat64(value), NewRat().SetUint64(100000000))
	floatStr := valueRat.FloatString(1)
	intValue := strings.Split(floatStr, ".")
	amountBig, ok := big.NewInt(0).SetString(intValue[0], 10)
	if !ok {
		panic("set big int error")
	}
	return amountBig.Int64()
}

func NewRat() *big.Rat {
	rat := new(big.Rat)
	return rat
}
