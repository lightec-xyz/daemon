package oasis

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	localSigner "github.com/lightec-xyz/daemon/rpc/oasis/local"
	testnetSigner "github.com/lightec-xyz/daemon/rpc/oasis/testnet"
)

const (
	Testnet = "testnet"
	Local   = "local"
	Mainnet = "mainnet"
)

type Option struct {
	Address        string
	TestnetAddress string
	LocalAddress   string
}

type Client struct {
	localSignerCall   *localSigner.ZkbridgeSigner // todo
	testnetSignerCall *testnetSigner.ZkbridgeSigner
	timout            time.Duration
	network           string
}

func NewClient(url, network string, option *Option) (*Client, error) {
	rpcDial, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	client := ethclient.NewClient(rpcDial)
	localSignerCall, err := localSigner.NewZkbridgeSigner(common.HexToAddress(option.LocalAddress), client)
	if err != nil {
		return nil, err
	}
	testnetSignerCall, err := testnetSigner.NewZkbridgeSigner(common.HexToAddress(option.TestnetAddress), client)
	if err != nil {
		return nil, err
	}
	return &Client{
		localSignerCall:   localSignerCall,
		testnetSignerCall: testnetSignerCall,
		timout:            60 * time.Second,
		network:           network,
	}, nil
}

func (c *Client) PublicKey() ([][]byte, error) {
	switch c.network {
	case Testnet:
		return c.TestnetPublicKey()
	case Local:
		return c.LocalPublicKey()
	default:
		return nil, fmt.Errorf("unknown network: %s", c.network)
	}
}

func (c *Client) SignBtcTx(rawTx, receiptTx, proof []byte) ([][][]byte, error) {
	switch c.network {
	case Testnet:
		return c.TestnetSignBtcTx(rawTx, receiptTx, proof)
	case Local:
		return c.LocalSignBtcTx(rawTx, receiptTx, proof)
	default:
		return nil, fmt.Errorf("unknown network: %s", c.network)
	}
}

func (c *Client) TestnetPublicKey() ([][]byte, error) {
	publicKeys, err := c.testnetSignerCall.GetPublicKeys(nil)
	if err != nil {
		return nil, err
	}
	return publicKeys, nil

}

func (c *Client) LocalSignBtcTx(rawTx, receiptTx, proof []byte) ([][][]byte, error) {
	signature1, err := c.localSignerCall.SignBtcTx(nil, rawTx, receiptTx, proof)
	if err != nil {
		return nil, err
	}
	return signature1, nil
}

func (c *Client) LocalPublicKey() ([][]byte, error) {
	publicKeys, err := c.localSignerCall.GetPublicKeys(nil)
	if err != nil {
		return nil, err
	}
	return publicKeys, nil

}

func (c *Client) TestnetSignBtcTx(rawTx, receiptTx, proof []byte) ([][][]byte, error) {
	signature1, err := c.testnetSignerCall.SignBtcTx(nil, rawTx, receiptTx, proof)
	if err != nil {
		return nil, err
	}
	return signature1, nil
}
