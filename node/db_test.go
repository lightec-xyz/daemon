package node

import (
	"fmt"
	"github.com/lightec-xyz/daemon/store"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func initStore() store.IStore {
	dbPath := "testdb"
	db, err := store.NewStore(dbPath, 0, 0, "zkbtc", false)
	if err != nil {
		panic(err)
	}
	return db
}

func TestWriteUnGenProof(t *testing.T) {
	store := initStore()
	err := WriteUnGenProof(store, Bitcoin, []string{"sdfsdfsdf", "2sdfsdfsd"})
	assert.Nil(t, err)
	ids, err := ReadAllUnGenProofIds(store, Bitcoin)
	assert.Nil(t, err)
	t.Log(ids)
	for _, item := range ids {
		err := DeleteUnGenProof(store, Bitcoin, item.TxId)
		assert.Nil(t, err)
	}
	ids, err = ReadAllUnGenProofIds(store, Bitcoin)
	assert.Nil(t, err)
	t.Log(ids)
}

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
	before, err := CheckBitcoinHeight(db)
	assert.Nil(t, err)
	assert.Equal(t, false, before)
	err = WriteBitcoinHeight(db, 100)
	assert.Nil(t, err)
	after, err := CheckBitcoinHeight(db)
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
			TxHash: fmt.Sprintf("%v", i),
		})
	}
	assert.Nil(t, err)
	err = WriteBitcoinTx(db, txes)
	assert.Nil(t, err)
	err = WriteDbProof(db, proofs)
	assert.Nil(t, err)
	txIds, err := ReadBitcoinTxIds(db, 100)
	assert.Nil(t, err)
	for _, txId := range txIds {
		t.Log(txId)
		tx, err := ReadDbTx(db, txId)
		assert.Nil(t, err)
		t.Log(tx)
		proof, err := ReadDbProof(db, txId)
		assert.Nil(t, err)
		t.Log(proof)
	}

}
