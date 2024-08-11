package node

import (
	"fmt"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/store"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func initStore() store.IStore {
	dbPath := "/Users/red/lworkspace/lightec/daemon/daemon/node/test/dbtest"
	db, err := store.NewStore(dbPath, 0, 0, "zkbtc", false)
	if err != nil {
		panic(err)
	}
	return db
}

func TestFindFinalityUpdateNearestSlot(t *testing.T) {

	iStore := initStore()
	err := WriteFinalityUpdateSlot(iStore, 900)
	assert.Nil(t, err)
	err = WriteFinalityUpdateSlot(iStore, 200)
	assert.Nil(t, err)
	err = WriteFinalityUpdateSlot(iStore, 600)
	assert.Nil(t, err)
	err = WriteFinalityUpdateSlot(iStore, 500)
	assert.Nil(t, err)
	err = WriteFinalityUpdateSlot(iStore, 100)
	assert.Nil(t, err)
	err = WriteFinalityUpdateSlot(iStore, 400)
	assert.Nil(t, err)
	err = WriteFinalityUpdateSlot(iStore, 700)
	assert.Nil(t, err)
	slot, ok, err := FindFinalityUpdateNearestSlot(iStore, 450)
	if !ok {
		t.Fatal("no find slot")
	}
	assert.Nil(t, err)
	t.Log(slot)
}

func TestPendingRequest(t *testing.T) {
	iStore := initStore()
	err := WritePendingRequest(iStore, "test01", &common.ZkProofRequest{
		Id:        "test01",
		StartTime: time.Now(),
	})
	assert.Nil(t, err)

	ids, err := ReadAllPendingRequests(iStore)
	assert.Nil(t, err)
	t.Log(ids)
	for _, id := range ids {
		err := DeletePendingRequest(iStore, id.RequestId())
		assert.Nil(t, err)
	}
	requests, err := ReadAllPendingRequests(iStore)
	assert.Nil(t, err)
	t.Log(requests)

}

func TestWriteUnGenProof(t *testing.T) {

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
	bitcoinHeight, ok, err := ReadBitcoinHeight(db)
	if !ok {
		panic("")
	}
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
	err = WriteDbProof(db, proofs)
	assert.Nil(t, err)
	txIds, err := ReadBitcoinTxIdsByHeight(db, 100)
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
