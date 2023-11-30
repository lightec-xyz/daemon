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

func (c *Client) Deposit(secret, txId, ethAddr string, index uint32, amount *big.Int, proof []byte) (string, error) {
	ctx := context.Background()
	chainID, err := c.ChainID(ctx)
	if err != nil {
		return "", err
	}
	nonce, err := c.NonceAt(ctx, common.HexToAddress(ethAddr), nil)
	if err != nil {
		return "", err
	}
	gasPrice, err := c.EthGasPrice()
	if err != nil {
		return "", err
	}
	privateKey, err := crypto.HexToECDSA(secret)
	if err != nil {
		return "", err
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return "", err
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasPrice = &gasPrice
	auth.GasLimit = 500000
	fixedTxId := [32]byte{}
	copy(fixedTxId[:], common.FromHex(txId))
	transaction, err := c.zkBridgeCall.Deposit(auth, fixedTxId, index, amount, common.HexToAddress(ethAddr), proof)
	if err != nil {
		return "", err
	}
	return transaction.Hash().Hex(), nil

}

func (c *Client) Redeem(from, secret string, amount, btcMinerFee *big.Int, receiveLockScript []byte) (string, error) {
	ctx := context.Background()
	chainID, err := c.ChainID(ctx)
	if err != nil {
		return "", err
	}
	nonce, err := c.NonceAt(ctx, common.HexToAddress(from), nil)
	if err != nil {
		return "", err
	}
	gasPrice, err := c.EthGasPrice()
	if err != nil {
		return "", err
	}
	privateKey, err := crypto.HexToECDSA(secret)
	if err != nil {
		return "", err
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return "", err
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasPrice = &gasPrice
	auth.GasLimit = 500000
	transaction, err := c.zkBridgeCall.Redeem(auth, amount, btcMinerFee, receiveLockScript)
	if err != nil {
		return "", err
	}
	return transaction.Hash().Hex(), nil
}
