package oasis

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	zkbridgeSigner "github.com/lightec-xyz/daemon/rpc/oasis/contract"
)

type Client struct {
	zkBridgeVerifyCall1 *zkbridgeSigner.ZkbridgeSigner // todo
	timout              time.Duration
}

func NewClient(url string, signerAddress string) (*Client, error) {
	// todo
	rpcDial, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	client := ethclient.NewClient(rpcDial)
	zkBridgeVerifyCall1, err := zkbridgeSigner.NewZkbridgeSigner(common.HexToAddress(signerAddress), client)
	if err != nil {
		return nil, err
	}
	return &Client{
		zkBridgeVerifyCall1: zkBridgeVerifyCall1,
		timout:              60 * time.Second,
	}, nil
}

func (c *Client) PublicKey() ([][]byte, error) {
	publicKeys, err := c.zkBridgeVerifyCall1.GetPublicKeys(nil)
	if err != nil {
		return nil, err
	}
	return publicKeys, nil
}

func (c *Client) SignBtcTx(rawTx, receiptTx, proof []byte) ([][][]byte, error) {
	signature1, err := c.zkBridgeVerifyCall1.SignBtcTx(nil, rawTx, receiptTx, proof)
	if err != nil {
		return nil, err
	}
	return signature1, nil
}
