package node

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/lightec-xyz/daemon/common"
	"math/big"
	"strings"
)

func UUID() string {
	newV7, err := uuid.NewV7()
	if err != nil {
		panic("should never happen")
	}
	return newV7.String()
}
func BtcToSat(value float64) int64 {
	valueRat := NewRat().Mul(NewRat().SetFloat64(value), NewRat().SetUint64(100000000))
	floatStr := valueRat.FloatString(1)
	valuesStr := strings.Split(floatStr, ".")
	amountBig, ok := big.NewInt(0).SetString(valuesStr[0], 10)
	if !ok {
		panic(fmt.Sprintf("never should happen:%v", value))
	}
	return amountBig.Int64()
}

func privateKeyToEthAddr(secret string) (string, error) {
	privateKey, err := crypto.HexToECDSA(secret)
	if err != nil {
		return "", err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	return address, nil
}

func trimOx(hash string) string {
	return strings.TrimPrefix(hash, "0x")
}

func txesToTxIds(txes []Transaction) []string {
	var txHashes []string
	for _, tx := range txes {
		txHashes = append(txHashes, tx.TxHash)
	}
	return txHashes
}

func proofsToDbProofs(proofs []*common.ZkProofRequest) []DbProof {
	var dbProofs []DbProof
	for _, proof := range proofs {
		dbProofs = append(dbProofs, DbProof{
			TxHash: proof.TxHash,
		})
	}
	return dbProofs
}

func txesToDbTxes(txes []Transaction) []DbTx {
	var dbtxes []DbTx
	for _, tx := range txes {
		dbtxes = append(dbtxes, DbTx{
			TxHash: tx.TxHash,
		})
	}
	return dbtxes
}

func requestsToProofUnGenId(requests []*common.ZkProofRequest) []string {
	var ids []string
	for _, req := range requests {
		ids = append(ids, req.Id())
	}
	return ids
}

func redeemToTxHashList(txs []RedeemProofParam) []string {
	var txHashList []string
	for _, tx := range txs {
		txHashList = append(txHashList, tx.TxHash)
	}
	return txHashList
}

func NewRat() *big.Rat {
	rat := new(big.Rat)
	return rat
}
