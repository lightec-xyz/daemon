package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/node"
	btcrpc "github.com/lightec-xyz/daemon/rpc/bitcoin"
	btctx "github.com/lightec-xyz/daemon/rpc/bitcoin/common"
	"github.com/lightec-xyz/daemon/rpc/bitcoin/types"
	ethrpc "github.com/lightec-xyz/daemon/rpc/ethereum"
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
	ethTxHash, err := m.ethClient.Deposit("", nonce, gasLimit, chainID, gasPrice, nil, nil)
	if err != nil {
		logger.Error(" deposit eth error:%v", err)
		return err
	}
	logger.Info("success deposit eth tx hash:%v", ethTxHash)
	return nil
}

func (m *Mock) RedeemTx(amount int64) error {
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
	txhash, err := m.ethClient.Redeem(m.cfg.EthPrivateKey, uint64(gasLimit), chainID, big.NewInt(int64(nonce)), gasPrice, 0, 0, redeemLockScript)
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
