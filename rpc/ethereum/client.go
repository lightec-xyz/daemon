package ethereum

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
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
	return &Client{Client: client, EthRPC: ethRPC, zkBridgeCall: zkBridgeCall}, nil
}

func (c *Client) Deposit(secret, txId, ethAddr string, index uint32,
	gasLimit uint64, nonce, chainID, gasPrice, amount *big.Int, proof []byte) (string, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 15*time.Second)
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
	auth.GasPrice = gasPrice
	auth.GasLimit = gasLimit
	fixedTxId := [32]byte{}
	copy(fixedTxId[:], common.FromHex(txId))
	transaction, err := c.zkBridgeCall.Deposit(auth, fixedTxId, index, amount, common.HexToAddress(ethAddr), proof)
	if err != nil {
		return "", err
	}
	return transaction.Hash().Hex(), nil

}

func (c *Client) Redeem(secret string, gasLimit uint64, chainID, nonce, gasPrice,
	amount, btcMinerFee *big.Int, receiveLockScript []byte) (string, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 15*time.Second)
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
	auth.GasPrice = gasPrice
	auth.GasLimit = gasLimit
	transaction, err := c.zkBridgeCall.Redeem(auth, amount, btcMinerFee, receiveLockScript)
	if err != nil {
		return "", err
	}
	return transaction.Hash().Hex(), nil
}
