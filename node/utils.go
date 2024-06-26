package node

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/lightec-xyz/daemon/common"
	"math/big"
	"strings"
	"time"
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

func trimOx(id string) string {
	return strings.TrimPrefix(id, "0x")
}

func txesToTxIds(txes []*Transaction) []string {
	var txHashes []string
	for _, tx := range txes {
		txHashes = append(txHashes, tx.TxHash)
	}
	return txHashes
}

func txesToDbProofs(txes []*Transaction) []DbProof {
	var dbProofs []DbProof
	for _, tx := range txes {
		dbProof := DbProof{
			TxHash: tx.TxHash,
		}
		if tx.Proofed {
			dbProof.Status = int(common.ProofSuccess)
		}
		dbProofs = append(dbProofs, dbProof)
	}
	return dbProofs
}

func proofsToDbProofs(proofs []*common.ZkProofRequest) []DbProof {
	var dbProofs []DbProof
	for _, proof := range proofs {
		dbProofs = append(dbProofs, DbProof{ // todo
			TxHash: proof.TxHash,
		})
	}
	return dbProofs
}

func proofToUnSubmitTx(resp *common.ZkProofResponse) DbUnSubmitTx {
	return DbUnSubmitTx{
		Hash:      resp.TxHash,
		ProofType: resp.ZkProofType,
		Proof:     hex.EncodeToString(resp.Proof),
		Timestamp: time.Now().UnixNano(),
	}
}

func txesToDbTxes(txes []*Transaction) []DbTx {
	var dbtxes []DbTx
	for _, tx := range txes {
		dbtxes = append(dbtxes, DbTx{
			TxHash:    tx.TxHash,
			Height:    tx.Height,
			TxType:    tx.TxType,
			ChainType: tx.ChainType,
			Amount:    tx.Amount,
		})
	}
	return dbtxes
}

func txesToUnGenProofs(chainType ChainType, txes []*Transaction) []*DbUnGenProof {
	var proofs []*DbUnGenProof
	for _, tx := range txes {
		if !tx.Proofed {
			proofs = append(proofs, &DbUnGenProof{
				ChainType: chainType,
				ProofType: tx.ProofType,
				TxHash:    tx.TxHash,
				Height:    tx.Height,
				TxIndex:   tx.TxIndex,
				Amount:    uint64(tx.Amount),
			})
		}
	}
	return proofs
}

func requestsToUnGenProofs(chainType ChainType, requests []*common.ZkProofRequest) []*DbUnGenProof {
	var proofs []*DbUnGenProof
	for _, req := range requests {
		proofs = append(proofs, &DbUnGenProof{
			TxHash:    req.TxHash,
			ProofType: req.ReqType,
			ChainType: chainType,
		})
	}
	return proofs
}

func txesByAddrGroup(txes []*Transaction) map[string][]DbTx {
	txMap := make(map[string][]DbTx)
	for _, tx := range txes {
		if tx.From == "" {
			continue
		}
		list, ok := txMap[tx.From]
		if ok {
			list = append(list, DbTx{
				TxHash:    tx.TxHash,
				Height:    tx.Height,
				TxType:    tx.TxType,
				ChainType: tx.ChainType,
				Amount:    tx.Amount,
				TxIndex:   tx.TxIndex,
			})
			txMap[tx.From] = list
		} else {
			txMap[tx.From] = []DbTx{
				{
					TxHash:    tx.TxHash,
					Height:    tx.Height,
					TxType:    tx.TxType,
					ChainType: tx.ChainType,
					Amount:    tx.Amount,
					TxIndex:   tx.TxIndex,
				},
			}
		}
	}
	return txMap
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
