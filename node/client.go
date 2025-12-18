package node

import (
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	btcproverClient "github.com/lightec-xyz/btc_provers/utils/client"
	"github.com/lightec-xyz/daemon/rpc/bitcoin"
	btctypes "github.com/lightec-xyz/daemon/rpc/bitcoin"
	"github.com/lightec-xyz/daemon/store"
)

// BtcClient why exists this client, because btc_provers use btcproverClient.IClient to get block header,it`s maybe get forked chain data
type BtcClient struct {
	IClient    btcproverClient.IClient
	btcClient  *bitcoin.Client
	chainStore *ChainStore
	initHeight int64
}

func (c *BtcClient) SetInitHeight(height int64) {
	//logger.Debug("set init height %v", height)
	c.initHeight = height
}

func (c *BtcClient) GetHeaderByHashStr(hash string) (string, error) {
	header, err := c.headerByHash(hash)
	if err != nil {
		bHeader, iErr := c.btcClient.GetBlockHeader(hash)
		if iErr != nil {
			return "", iErr
		}
		if int64(bHeader.Height) <= c.initHeight {
			//logger.Warn("headerByHash error:%v %v %v ", bHeader.Height, hash, err)
			return c.IClient.GetHeaderByHashStr(hash)
		}
		return "", err
	}
	return header, nil
}

func (c *BtcClient) GetHeaderByHash(hash *chainhash.Hash) (string, error) {
	header, err := c.GetHeaderByHashStr(hash.String())
	if err != nil {
		return "", err
	}
	return header, nil
}

func (c *BtcClient) GetBlockHash(height int64) (*chainhash.Hash, error) {
	hash, err := c.blockHash(height)
	if err != nil {
		if height <= c.initHeight {
			//logger.Warn("blockHash error:%v %v", height, err)
			return c.IClient.GetBlockHash(height)
		}
		return nil, err
	}
	return hash, nil
}

func (c *BtcClient) GetHeaderByHeight(height int64) (string, error) {
	header, err := c.headerByHeight(height)
	if err != nil {
		if height <= c.initHeight {
			//logger.Warn("headerByHeight error: %v %v %v ", height, c.initHeight, err)
			return c.IClient.GetHeaderByHeight(height)
		}
		return "", err
	}
	return header, nil
}

func (c *BtcClient) GetBlock(hash string) (*btcjson.GetBlockVerboseResult, error) {
	block, err := c.readBlockFromDb(hash)
	if err != nil {
		bHeader, bErr := c.btcClient.GetBlockHeader(hash)
		if bErr != nil {
			return nil, bErr
		}
		if int64(bHeader.Height) <= c.initHeight {
			return c.IClient.GetBlock(hash)
		}
		return nil, err
	}
	return block, nil

}

func (c *BtcClient) headerByHash(hash string) (string, error) {
	header, ok, err := c.chainStore.ReadBlockHeader(hash)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", fmt.Errorf("db no find header %v", hash)
	}
	return header, nil
}

func (c *BtcClient) headerByHeight(height int64) (string, error) {
	hash, err := c.blockHash(height)
	if err != nil {
		return "", err
	}
	return c.headerByHash(hash.String())
}
func (c *BtcClient) blockHash(height int64) (*chainhash.Hash, error) {
	hash, ok, err := c.chainStore.ReadBitcoinHash(uint64(height))
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("db no find hash %v", height)
	}
	return chainhash.NewHashFromStr(hash)
}

func (c *BtcClient) readBlockFromDb(hash string) (*btcjson.GetBlockVerboseResult, error) {
	blockData, exists, err := c.chainStore.ReadBtcBlock(hash)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("db no find block %v", hash)
	}
	var block btcjson.GetBlockVerboseResult
	err = json.Unmarshal([]byte(blockData), &block)
	if err != nil {
		return nil, err
	}
	txes := make([]string, len(block.RawTx))
	for i := 0; i < len(txes); i++ {
		txes[i] = block.RawTx[i].Txid
	}
	block.Tx = txes
	return &block, nil
}

func NewBtcClient(client btcproverClient.IClient, store store.IStore, btClient *bitcoin.Client, initHeight int64) *BtcClient {
	return &BtcClient{
		IClient:    client,
		chainStore: NewChainStore(store),
		btcClient:  btClient,
		initHeight: initHeight,
	}
}

func toTxIds(txes []btctypes.Tx) []string {
	var txIds []string
	for _, tx := range txes {
		txIds = append(txIds, tx.Txid)
	}
	return txIds
}
