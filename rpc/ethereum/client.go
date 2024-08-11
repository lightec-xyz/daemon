package ethereum

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/lightec-xyz/daemon/logger"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/lightec-xyz/daemon/rpc/ethereum/zkbridge"
)

// todo
type Client struct {
	*ethclient.Client
	zkBridgeCall *zkbridge.Zkbridge
	zkBtcCall    *zkbridge.Zkbtc
	utxoCall     *zkbridge.Utxo
	timeout      time.Duration
}

func NewClient(endpoint string, zkBridgeAddr, zkBtcAddr, utxoManager string) (*Client, error) {
	client, err := ethclient.Dial(endpoint)
	if err != nil {
		return nil, err
	}
	zkBridgeCall, err := zkbridge.NewZkbridge(ethcommon.HexToAddress(zkBridgeAddr), client)
	if err != nil {
		return nil, err
	}
	zkBtcCall, err := zkbridge.NewZkbtc(ethcommon.HexToAddress(zkBtcAddr), client)
	if err != nil {
		return nil, err
	}
	utxo, err := zkbridge.NewUtxo(ethcommon.HexToAddress(utxoManager), client)
	if err != nil {
		return nil, err
	}
	return &Client{
		Client:       client,
		zkBridgeCall: zkBridgeCall,
		zkBtcCall:    zkBtcCall,
		utxoCall:     utxo,
		timeout:      15 * time.Second,
	}, nil
}

func (e *Client) GetUtxo(hash string) (zkbridge.UTXOManagerUTXO, error) {
	txId := [32]byte{}
	hexBytes, err := hex.DecodeString(hash)
	if err != nil {
		return zkbridge.UTXOManagerUTXO{}, err
	}
	copy(txId[:], hexBytes)
	result, err := e.utxoCall.UtxoOf(nil, txId)
	if err != nil {
		return zkbridge.UTXOManagerUTXO{}, err
	}
	return result, nil

}

func (e *Client) ChainFork(height uint64) (bool, error) {
	block, err := e.Client.BlockByNumber(context.Background(), big.NewInt(int64(height)))
	if err != nil {
		return false, err
	}
	preHeight := height - 1
	if preHeight < 0 {
		return false, nil
	}
	preBlock, err := e.Client.BlockByNumber(context.Background(), big.NewInt(int64(preHeight)))
	if err != nil {
		return false, err
	}
	if block.ParentHash().Hex() != preBlock.Hash().Hex() {
		return true, nil
	}
	return false, nil
}

func (e *Client) CheckTx(txHash string) (bool, error) {
	tx, err := e.Client.TransactionReceipt(context.Background(), ethcommon.HexToHash(txHash))
	if err != nil {
		return false, err
	}
	if tx.Status == types.ReceiptStatusSuccessful {
		return true, nil
	}
	return false, nil
}

func (e *Client) GetTxSender(txHash, blockHash string, index uint) (string, error) {
	tx, pending, err := e.Client.TransactionByHash(context.Background(), ethcommon.HexToHash(txHash))
	if err != nil {
		logger.Error("get eth tx error:%v %v", txHash, err)
		return "", err
	}
	if pending {
		return "", fmt.Errorf("tx %v is pending", txHash)
	}
	sender, err := e.Client.TransactionSender(context.Background(), tx, ethcommon.HexToHash(blockHash), index)
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

func (c *Client) GetPendingNonce(addr string) (uint64, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), c.timeout)
	defer cancelFunc()
	nonce, err := c.PendingNonceAt(ctx, ethcommon.HexToAddress(addr))
	if err != nil {
		return 0, err
	}
	return nonce, nil
}

func (c *Client) GetEstimateGasLimit(from, to, txId string, index uint32, amount *big.Int, proofData []byte) (uint64, error) {
	//todo
	ctx, cancelFunc := context.WithTimeout(context.Background(), c.timeout)
	defer cancelFunc()
	zkBridgeAbi, err := abi.JSON(strings.NewReader(zkbridge.ZkbridgeMetaData.ABI))
	if err != nil {
		return 0, err
	}
	fixedTxId := [32]byte{}
	copy(fixedTxId[:], ethcommon.FromHex(txId))
	txData, err := zkBridgeAbi.Pack("deposit", fixedTxId, index, amount, proofData)
	if err != nil {
		return 0, err
	}
	toAddress := ethcommon.HexToAddress(to)
	gas, err := c.EstimateGas(ctx, ethereum.CallMsg{
		From:      ethcommon.HexToAddress(from),
		To:        &toAddress,
		Gas:       0,
		Value:     big.NewInt(0),
		GasFeeCap: big.NewInt(200000000000),
		Data:      txData,
	})
	if err != nil {
		return 0, err
	}
	return gas, nil
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

func (c *Client) GetZkBtcBalance(addr string) (*big.Int, error) {
	balance, err := c.zkBtcCall.BalanceOf(nil, ethcommon.HexToAddress(addr))
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func (c *Client) Deposit(secret string, nonce, gasLimit uint64, chainID, gasPrice *big.Int, rawBtcTx, proof []byte) (string, error) {
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
	auth.Nonce = big.NewInt(int64(nonce))
	//auth.GasPrice = gasPrice todo
	auth.GasFeeCap = gasPrice
	auth.GasLimit = gasLimit
	transaction, err := c.zkBridgeCall.Deposit(auth, rawBtcTx, proof[:])
	if err != nil {
		return "", err
	}
	return transaction.Hash().Hex(), nil

}

func (c *Client) UpdateUtxoChange(secret []byte, txIds []string, nonce, gasLimit uint64, chainID, gasPrice *big.Int, proof []byte) (string, error) {
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
	// todo
	transaction, err := c.zkBridgeCall.UpdateChange(auth, TxIdsToFixedIds(txIds)[0], proof)
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
	//auth.GasPrice = gasPrice  todo
	auth.GasLimit = gasLimit
	auth.GasFeeCap = gasPrice
	transaction, err := c.zkBridgeCall.Redeem(auth, amount, btcMinerFee, receiveLockScript)
	if err != nil {
		return "", err
	}
	return transaction.Hash().Hex(), nil
}

func (c *Client) GetMultiSigScript() ([]byte, error) {
	multiSigScript, err := c.zkBridgeCall.MultiSigScript(nil)
	if err != nil {
		return nil, err
	}
	return multiSigScript, nil
}
