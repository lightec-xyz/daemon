package node

import (
	"fmt"
	"github.com/lightec-xyz/daemon/store"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestDb_Mock(t *testing.T) {
	dbPath := "testdb"
	file, err := os.Stat(dbPath)
	assert.Nil(t, err)
	if file.IsDir() {
		err = os.RemoveAll(dbPath)
		assert.Nil(t, err)
	}
	db, err := store.NewStore(dbPath, 0, 0, "zkbtc", false)
	assert.Nil(t, err)
	before, err := ReadInitBitcoinHeight(db)
	assert.Nil(t, err)
	assert.Equal(t, false, before)
	err = WriteBitcoinHeight(db, 100)
	assert.Nil(t, err)
	after, err := ReadInitBitcoinHeight(db)
	assert.Nil(t, err)
	assert.Equal(t, true, after)
	var bitcoinHeight int64
	bitcoinHeight, err = ReadBitcoinHeight(db)
	assert.Nil(t, err)
	assert.Equal(t, int64(100), bitcoinHeight)
	var txes []DbTx
	var proofs []DbProof
	for i := 0; i < 100; i++ {
		txes = append(txes, DbTx{
			TxHash: fmt.Sprintf("%v", i),
			Height: 100,
		})
		proofs = append(proofs, DbProof{
			TxId: fmt.Sprintf("%v", i),
		})
	}
	err = WriteBitcoinTxIds(db, 100, txes)
	assert.Nil(t, err)
	err = WriteBitcoinTx(db, txes)
	assert.Nil(t, err)
	err = WriteProof(db, proofs)
	assert.Nil(t, err)
	txIds, err := ReadBitcoinTxIds(db, 100)
	assert.Nil(t, err)
	for _, txId := range txIds {
		t.Log(txId)
		tx, err := ReadTransaction(db, txId)
		assert.Nil(t, err)
		t.Log(tx)
		proof, err := ReadProof(db, txId)
		assert.Nil(t, err)
		t.Log(proof)
	}

}
