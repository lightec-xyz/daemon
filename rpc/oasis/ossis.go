package oasis

import (
	"encoding/hex"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	signer "github.com/lightec-xyz/daemon/rpc/oasis/singer"
)

type Client struct {
	signerCall *signer.ZkbridgeSigner
	timout     time.Duration
	network    string
	address    string
}

func NewClient(url string, address string) (*Client, error) {
	rpcDial, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	client := ethclient.NewClient(rpcDial)
	testnetSignerCall, err := signer.NewZkbridgeSigner(common.HexToAddress(address), client)
	if err != nil {
		return nil, err
	}
	return &Client{
		signerCall: testnetSignerCall,
		timout:     60 * time.Second,
	}, nil
}

func (c *Client) PublicKey() ([][]byte, error) {
	publicKeys, err := c.signerCall.GetPublicKeys(nil)
	if err != nil {
		return nil, err
	}
	return publicKeys, nil
}

func (c *Client) SignBtcTx(txId, currentScRoot, proof string, sigHashes []string, minerReward *big.Int) ([][][]byte, error) {
	txIdBytes, err := hex.DecodeString(txId)
	if err != nil {
		return nil, err
	}
	scRootBytes, err := hex.DecodeString(currentScRoot)
	if err != nil {
		return nil, err
	}
	proofBytes, err := hex.DecodeString(proof)
	if err != nil {
		return nil, err
	}
	var sigHashBytes [][32]byte
	for _, v := range sigHashes {
		sigHash, err := hex.DecodeString(v)
		if err != nil {
			return nil, err
		}
		sigHashBytes = append(sigHashBytes, [32]byte(sigHash))
	}
	signature1, err := c.signerCall.SignBtcTx(nil, [32]byte(txIdBytes), minerReward, sigHashBytes, [32]byte(scRootBytes), proofBytes)
	if err != nil {
		return nil, err
	}
	return signature1, nil
}
