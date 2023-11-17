package ethereum

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc/ethereum/zkbridge"
)

type Client struct {
	*ethclient.Client
	zkBridgeCall    *zkbridge.Zkbridge
	utxoCall        *zkbridge.Utxo
	btcTxVerifyCall *zkbridge.BtcTxVerify
	zkbtcBridgeAbi  abi.ABI
	zkbtcBridgeAddr string
	timeout         time.Duration
}

func NewClient(endpoint string, zkBridgeAddr, utxoManager, txVerify string) (*Client, error) {
	client, err := ethclient.Dial(endpoint)
	if err != nil {
		return nil, err
	}
	zkBridgeCall, err := zkbridge.NewZkbridge(ethcommon.HexToAddress(zkBridgeAddr), client)
	if err != nil {
		return nil, err
	}
	utxo, err := zkbridge.NewUtxo(ethcommon.HexToAddress(utxoManager), client)
	if err != nil {
		return nil, err
	}
	btcTxVerify, err := zkbridge.NewBtcTxVerify(ethcommon.HexToAddress(txVerify), client)
	if err != nil {
		return nil, err
	}
	zkbtcBridgeAbi, err := abi.JSON(bytes.NewReader([]byte(zkbtcBridgeAbiConst)))
	if err != nil {
		return nil, err
	}
	return &Client{
		Client:          client,
		zkBridgeCall:    zkBridgeCall,
		utxoCall:        utxo,
		btcTxVerifyCall: btcTxVerify,
		timeout:         15 * time.Second,
		zkbtcBridgeAbi:  zkbtcBridgeAbi,
		zkbtcBridgeAddr: zkBridgeAddr,
	}, nil
}

func (c *Client) GetRaised(hash string) (bool, error) {
	return false, nil
}

func (c *Client) EthBalance(addr string) (*big.Int, error) {
	balance, err := c.Client.BalanceAt(context.Background(), ethcommon.HexToAddress(addr), nil)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func (c *Client) GetCpLatestAddedTime() (uint64, error) {
	cpLatestAddedTime, err := c.btcTxVerifyCall.CpLatestAddedTime(nil)
	if err != nil {
		return 0, err
	}
	return cpLatestAddedTime, nil
}

func (c *Client) CheckUtxo(txId string) (bool, error) {
	txIdBytes := ethcommon.FromHex(txId)
	utxoOf, err := c.utxoCall.UtxoOf(nil, [32]byte(txIdBytes))
	if err != nil {
		return false, err
	}
	if bytes.Equal(utxoOf.Txid[:], txIdBytes) {
		return true, nil
	}
	return false, nil
}

func (c *Client) UtxoConfirm(txId string) (bool, error) {
	txIdBytes, err := hex.DecodeString(txId)
	if err != nil {
		return false, err
	}
	utxoOf, err := c.utxoCall.UtxoOf(nil, [32]byte(txIdBytes))
	if err != nil {
		return false, err
	}
	if bytes.Equal(utxoOf.Txid[:], txIdBytes) && utxoOf.IsChangeConfirmed {
		return true, nil
	}
	return false, nil
}

func (c *Client) GetDepthByAmount(amount uint64, raise, blockSig bool) (uint32, error) {
	depth, _, err := c.btcTxVerifyCall.GetDepthByAmount(nil, amount, raise, blockSig)
	if err != nil {
		return 0, nil
	}
	return depth, nil
}

func (c *Client) SuggestedCP() ([]byte, error) {
	hash, err := c.btcTxVerifyCall.SuggestedCP(nil)
	if err != nil {
		return nil, err
	}
	return hash[:], nil

}

func (c *Client) GetUtxo(hash string) (zkbridge.UTXOManagerUTXO, error) {
	txId := [32]byte{}
	hexBytes, err := hex.DecodeString(hash)
	if err != nil {
		return zkbridge.UTXOManagerUTXO{}, err
	}
	copy(txId[:], hexBytes)
	result, err := c.utxoCall.UtxoOf(nil, txId)
	if err != nil {
		return zkbridge.UTXOManagerUTXO{}, err
	}
	return result, nil

}

func (c *Client) CheckTx(txHash string) (bool, error) {
	tx, err := c.Client.TransactionReceipt(context.Background(), ethcommon.HexToHash(txHash))
	if err != nil {
		return false, err
	}
	if tx.Status == types.ReceiptStatusSuccessful {
		return true, nil
	}
	return false, nil
}

func (c *Client) GetTxSender(txHash, blockHash string, index uint) (string, error) {
	tx, pending, err := c.Client.TransactionByHash(context.Background(), ethcommon.HexToHash(txHash))
	if err != nil {
		logger.Error("get eth tx error:%v %v", txHash, err)
		return "", err
	}
	if pending {
		return "", fmt.Errorf("tx %v is pending", txHash)
	}
	sender, err := c.Client.TransactionSender(context.Background(), tx, ethcommon.HexToHash(blockHash), index)
	if err != nil {
		logger.Error("get eth tx sender error:%v %v", txHash, err)
		return "", err
	}
	return sender.Hex(), nil

}

func (c *Client) CheckDepositProof(txId string) (bool, error) {
	//todo
	return false, nil
}

func (c *Client) GetLogs(hash string, addrList []string, topicList []string) ([]types.Log, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), c.timeout)
	defer cancelFunc()
	blockHash := ethcommon.HexToHash(hash)
	var addresses []ethcommon.Address
	for _, addr := range addrList {
		address := ethcommon.HexToAddress(addr)
		addresses = append(addresses, address)
	}
	var topics [][]ethcommon.Hash
	var matchTopic []ethcommon.Hash
	for _, topic := range topicList {
		topicHash := ethcommon.HexToHash(topic)
		matchTopic = append(matchTopic, topicHash)
	}
	topics = append(topics, matchTopic)

	filterQuery := ethereum.FilterQuery{
		BlockHash: &blockHash,
		Addresses: addresses,
		Topics:    topics,
	}
	logs, err := c.FilterLogs(ctx, filterQuery)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (c *Client) GetBlock(height int64) (*types.Block, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), c.timeout)
	defer cancelFunc()
	block, err := c.BlockByNumber(ctx, big.NewInt(height))
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (c *Client) GetNonce(addr string) (uint64, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), c.timeout)
	defer cancelFunc()
	nonce, err := c.NonceAt(ctx, ethcommon.HexToAddress(addr), nil)
	if err != nil {
		return 0, err
	}
	return nonce, nil
}

func (c *Client) GetGasPrice() (*big.Int, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), c.timeout)
	defer cancelFunc()
	gasPrice, err := c.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}
	return gasPrice, nil
}

func (c *Client) GetChainId() (*big.Int, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), c.timeout)
	defer cancelFunc()
	chainId, err := c.ChainID(ctx)
	if err != nil {
		return nil, err
	}
	return chainId, nil
}

func (c *Client) EstimateDepositGasLimit(from string, params *zkbridge.IBtcTxVerifierPublicWitnessParams, gasPrice *big.Int,
	btcRawTx, proof []byte) (uint64, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), c.timeout)
	defer cancelFunc()
	payload, err := c.zkbtcBridgeAbi.Pack("deposit", btcRawTx, params, proof)
	if err != nil {
		return 0, err
	}
	toAddress := ethcommon.HexToAddress(c.zkbtcBridgeAddr)
	msg := ethereum.CallMsg{
		From:     ethcommon.HexToAddress(from),
		To:       &toAddress,
		GasPrice: gasPrice,
		Value:    big.NewInt(0),
		Data:     payload,
	}
	gasLimit, err := c.EstimateGas(ctx, msg)
	if err != nil {
		return 0, err
	}
	return gasLimit, nil
}

func (c *Client) Deposit(secret []byte, params *zkbridge.IBtcTxVerifierPublicWitnessParams, nonce, gasLimit uint64, chainID, gasPrice *big.Int,
	btcRawTx, proof []byte) (string, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), c.timeout)
	defer cancelFunc()
	privateKey, err := crypto.ToECDSA(secret)
	if err != nil {
		return "", err
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return "", err
	}
	auth.Context = ctx
	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasFeeCap = gasPrice
	auth.GasLimit = gasLimit
	transaction, err := c.zkBridgeCall.Deposit(auth, btcRawTx, *params, proof)
	if err != nil {
		return "", err
	}
	return transaction.Hash().Hex(), nil

}

func (c *Client) EstimateUpdateUtxoGasLimit(from string, param *zkbridge.IBtcTxVerifierPublicWitnessParams, gasPrice,
	minerReward *big.Int, txId, proof []byte) (uint64, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), c.timeout)
	defer cancelFunc()
	payload, err := c.zkbtcBridgeAbi.Pack("updateRedeem", param, [32]byte(txId), minerReward, proof)
	if err != nil {
		return 0, err
	}
	toAddress := ethcommon.HexToAddress(c.zkbtcBridgeAddr)
	msg := ethereum.CallMsg{
		From:     ethcommon.HexToAddress(from),
		To:       &toAddress,
		GasPrice: gasPrice,
		Value:    big.NewInt(0),
		Data:     payload,
	}
	gasLimit, err := c.EstimateGas(ctx, msg)
	if err != nil {
		return 0, err
	}
	return gasLimit, nil
}

func (c *Client) UpdateUtxoChange(secret []byte, params *zkbridge.IBtcTxVerifierPublicWitnessParams, nonce, gasLimit uint64,
	chainID, gasPrice, minerReward *big.Int, txId, proof []byte) (string, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), c.timeout)
	defer cancelFunc()
	privateKey, err := crypto.ToECDSA(secret)
	if err != nil {
		return "", err
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return "", err
	}
	auth.Context = ctx
	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasFeeCap = gasPrice
	auth.GasLimit = gasLimit
	transaction, err := c.zkBridgeCall.UpdateRedeem(auth, *params, [32]byte(txId), minerReward, proof)
	if err != nil {
		return "", err
	}
	return transaction.Hash().Hex(), nil

}

func TxIdsToFixedIds(txIds []string) [][32]byte {
	fixedTxIds := make([][32]byte, 0)
	for _, txId := range txIds {
		fixedTxId := [32]byte{}
		copy(fixedTxId[:], ethcommon.FromHex(txId))
		fixedTxIds = append(fixedTxIds, fixedTxId)
	}
	return fixedTxIds
}

func (c *Client) Redeem(secret string, gasLimit uint64, chainID, nonce, gasPrice *big.Int,
	amount, btcMinerFee uint64, receiveLockScript []byte) (string, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), c.timeout)
	defer cancelFunc()
	privateKey, err := crypto.HexToECDSA(secret)
	if err != nil {
		return "", err
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return "", err
	}
	auth.Context = ctx
	auth.Nonce = nonce
	auth.GasFeeCap = gasPrice
	auth.GasLimit = gasLimit
	transaction, err := c.zkBridgeCall.Redeem(auth, amount, btcMinerFee, receiveLockScript)
	if err != nil {
		return "", err
	}
	return transaction.Hash().Hex(), nil
}
