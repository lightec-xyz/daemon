package ethereum

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/lightec-xyz/daemon/transaction/ethereum/zkbridge"
	"github.com/onrik/ethrpc"
	"math/big"
	"time"
)

// todo
type Client struct {
	*ethclient.Client
	*ethrpc.EthRPC
	zkBridgeCall *zkbridge.Zkbridge
	timeout      time.Duration
}

func NewClient(endpoint string, zkBridgeAddr string) (*Client, error) {
	rpcDial, err := rpc.Dial(endpoint)
	if err != nil {
		return nil, err
	}
	ethRPC := ethrpc.New(endpoint)
	client := ethclient.NewClient(rpcDial)
	zkBridgeCall, err := zkbridge.NewZkbridge(common.HexToAddress(zkBridgeAddr), client)
	if err != nil {
		return nil, err
	}
	return &Client{Client: client, EthRPC: ethRPC, zkBridgeCall: zkBridgeCall, timeout: 15 * time.Second}, nil
}

func (c *Client) GetLogs(hash string, addrList []string, topicList []string) ([]types.Log, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), c.timeout)
	defer cancelFunc()
	blockHash := common.HexToHash(hash)
	var addresses []common.Address
	for _, addr := range addrList {
		address := common.HexToAddress(addr)
		addresses = append(addresses, address)
	}
	var topics [][]common.Hash
	for _, topic := range topicList {
		topicHash := common.HexToHash(topic)
		topics = append(topics, []common.Hash{topicHash})
	}

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
	nonce, err := c.NonceAt(ctx, common.HexToAddress(addr), nil)
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

func (c *Client) GetBalance(addr string) (*big.Int, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), c.timeout)
	defer cancelFunc()
	balance, err := c.BalanceAt(ctx, common.HexToAddress(addr), nil)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func (c *Client) Deposit(secret, txId string, index uint32,
	nonce, gasLimit uint64, chainID, gasPrice, amount *big.Int, proof []byte) (string, error) {
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
	fixedTxId := [32]byte{}
	copy(fixedTxId[:], common.FromHex(txId))
	fmt.Println("txIndex:", index)
	transaction, err := c.zkBridgeCall.Deposit(auth, fixedTxId, index, amount, proof)
	if err != nil {
		return "", err
	}
	return transaction.Hash().Hex(), nil

}

func (c *Client) UpdateUtxoChange(secret, txId string, nonce, gasLimit uint64, chainID, gasPrice *big.Int, proof []byte) (string, error) {
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
	fixedTxId := [32]byte{}
	copy(fixedTxId[:], common.FromHex(txId))
	transaction, err := c.zkBridgeCall.UpdateChange(auth, fixedTxId, proof)
	if err != nil {
		return "", err
	}
	return transaction.Hash().Hex(), nil

}

func (c *Client) Redeem(secret string, gasLimit uint64, chainID, nonce, gasPrice,
	amount, btcMinerFee *big.Int, receiveLockScript []byte) (string, error) {
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
