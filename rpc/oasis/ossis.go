package oasis

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	zkbridge_verify "github.com/lightec-xyz/daemon/rpc/oasis/contract"
	"time"
)

type Client struct {
	zkBridgeVerifyCall1 *zkbridge_verify.ZkbridgeVerify // todo
	zkBridgeVerifyCall2 *zkbridge_verify.ZkbridgeVerify
	zkBridgeVerifyCall3 *zkbridge_verify.ZkbridgeVerify
	timout              time.Duration
}

func NewClient(url string, address []string) (*Client, error) {
	// todo
	rpcDial, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	client := ethclient.NewClient(rpcDial)
	zkBridgeVerifyCall1, err := zkbridge_verify.NewZkbridgeVerify(common.HexToAddress(address[0]), client)
	if err != nil {
		return nil, err
	}
	zkBridgeVerifyCall2, err := zkbridge_verify.NewZkbridgeVerify(common.HexToAddress(address[1]), client)
	if err != nil {
		return nil, err
	}
	zkBridgeVerifyCall3, err := zkbridge_verify.NewZkbridgeVerify(common.HexToAddress(address[2]), client)
	if err != nil {
		return nil, err
	}
	return &Client{
		zkBridgeVerifyCall1: zkBridgeVerifyCall1,
		zkBridgeVerifyCall2: zkBridgeVerifyCall2,
		zkBridgeVerifyCall3: zkBridgeVerifyCall3,
		timout:              60 * time.Second,
	}, nil
}

func (c *Client) PublicKey() ([][]byte, error) {
	// todo
	publicKey1, err := c.zkBridgeVerifyCall1.SignerPublicKey(nil)
	if err != nil {
		return nil, err
	}
	publicKey2, err := c.zkBridgeVerifyCall2.SignerPublicKey(nil)
	if err != nil {
		return nil, err
	}
	publicKey3, err := c.zkBridgeVerifyCall3.SignerPublicKey(nil)
	if err != nil {
		return nil, err
	}
	return [][]byte{publicKey1, publicKey2, publicKey3}, nil
}

func (c *Client) SignBtcTx(rawTx, receiptTx, proof []byte) ([][]byte, error) {
	// todo
	signature1, err := c.zkBridgeVerifyCall1.SignBtcTx(nil, rawTx, receiptTx, proof)
	if err != nil {
		return nil, err
	}
	return signature1, nil
}
