package node

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/rpc"
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

func proofsToDbProofs(proofs []Proof) []DbProof {
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

func depositToTxHash(txs []DepositProofParam) []string {
	var txHashList []string
	for _, tx := range txs {
		txHashList = append(txHashList, tx.TxHash)
	}
	return txHashList
}

func redeemToTxHashList(txs []RedeemProofParam) []string {
	var txHashList []string
	for _, tx := range txs {
		txHashList = append(txHashList, tx.TxHash)
	}
	return txHashList
}
func toDepositZkProofRequest(list []DepositProofParam) ([]*common.ZkProofRequest, error) {
	var result []*common.ZkProofRequest
	for _, item := range list {
		data := rpc.DepositRequest{
			TxHash:    item.TxHash,
			BlockHash: item.BlockHash,
		}
		result = append(result, common.NewZkProofRequest(common.DepositTxType, data, 0, item.TxHash))
	}
	return result, nil
}

func toUpdateZkProofRequest(redeemTxes []Transaction) ([]*common.ZkProofRequest, error) {
	var result []*common.ZkProofRequest
	for _, item := range redeemTxes {
		data := rpc.VerifyRequest{TxHash: item.TxHash, BlockHash: item.BlockHash}
		result = append(result, common.NewZkProofRequest(common.VerifyTxType, data, 0, item.TxHash))
	}
	return result, nil
}

func toTxInEth2Request(list []RedeemProofParam) ([]*common.ZkProofRequest, error) {
	var result []*common.ZkProofRequest
	for _, item := range list {
		result = append(result, common.NewZkProofRequest(common.TxInEth2, item, 0, item.TxHash))
	}
	return result, nil
}

func NewRat() *big.Rat {
	rat := new(big.Rat)
	return rat
}
