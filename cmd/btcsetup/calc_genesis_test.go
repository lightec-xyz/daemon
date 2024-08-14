package main

import (
	"log"
	"testing"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/lightec-xyz/btc_provers/circuits/common"
	"github.com/lightec-xyz/btc_provers/utils/client"
	"github.com/stretchr/testify/assert"
)

var connCfg = &rpcclient.ConnConfig{
	Host: "localhost:18332",
	User: "test",
	Pass: "123456",
}

func TestCalcGenesis(t *testing.T) {
	cl, err := client.NewClient(connCfg.Host, connCfg.User, connCfg.Pass)
	assert.NoError(t, err)

	lastestBh, err := cl.GetBlockCount()
	assert.NoError(t, err)
	log.Printf("lastest block height: %d", lastestBh)

	genesisBlockheight := (uint32(lastestBh)/common.CapacityDifficultyBlock - 2) * common.CapacityDifficultyBlock

	log.Printf("genesis block height: %d", genesisBlockheight)
}
