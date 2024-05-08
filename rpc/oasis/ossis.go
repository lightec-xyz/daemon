package oasis

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	alphaSigner "github.com/lightec-xyz/daemon/rpc/oasis/alpha"
	signer "github.com/lightec-xyz/daemon/rpc/oasis/contract"
)

type Option struct {
	Address      string
	AlphaAddress string
}

type Client struct {
	signerCall      *signer.ZkbridgeSigner // todo
	alphaSignerCall *alphaSigner.ZkbridgeSigner
	timout          time.Duration
}

func NewClient(url string, option *Option) (*Client, error) {
	rpcDial, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	client := ethclient.NewClient(rpcDial)
	signerCall, err := signer.NewZkbridgeSigner(common.HexToAddress(option.Address), client)
	if err != nil {
		return nil, err
	}
	alphaSignerCall, err := alphaSigner.NewZkbridgeSigner(common.HexToAddress(option.AlphaAddress), client)
	if err != nil {
		return nil, err
	}
	return &Client{
		signerCall:      signerCall,
		alphaSignerCall: alphaSignerCall,
		timout:          60 * time.Second,
	}, nil
}

func (c *Client) AlphaPublicKey() ([][]byte, error) {
	publicKeys, err := c.alphaSignerCall.GetPublicKeys(nil)
	if err != nil {
		return nil, err
	}
	return publicKeys, nil

}

func (c *Client) AlphaSignBtcTx(rawTx, receiptTx, proof []byte) ([][][]byte, error) {
	signature1, err := c.alphaSignerCall.SignBtcTx(nil, rawTx, receiptTx, proof)
	if err != nil {
		return nil, err
	}
	return signature1, nil
}

func (c *Client) PublicKey() ([][]byte, error) {
	publicKeys, err := c.signerCall.GetPublicKeys(nil)
	if err != nil {
		return nil, err
	}
	return publicKeys, nil
}

func (c *Client) SignBtcTx(rawTx, receiptTx, proof []byte) ([][][]byte, error) {
	signature1, err := c.signerCall.SignBtcTx(nil, rawTx, receiptTx, proof)
	if err != nil {
		return nil, err
	}
	return signature1, nil
}
