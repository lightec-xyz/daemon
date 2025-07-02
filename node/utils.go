package node

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	btccdEcdsa "github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/davecgh/go-spew/spew"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"math/big"
	"runtime"
	"strings"
)

func txSkipCheck(txes []*DbTx) []*DbTx {
	for _, item := range txes {
		item.Proved = false
	}
	return txes
}

func LowercaseFirstLetter(s string) string {
	if s == "" {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func RsToSignature(sig string) ([]byte, error) {
	sigBytes, err := hex.DecodeString(sig)
	if err != nil {
		return nil, err
	}
	r := secp256k1.ModNScalar{}
	s := secp256k1.ModNScalar{}
	r.SetByteSlice(sigBytes[:32])
	s.SetByteSlice(sigBytes[32:])
	signature := btccdEcdsa.NewSignature(&r, &s)
	return signature.Serialize(), nil
}

func TxIdIsEmpty(txId [32]byte) bool {
	for _, b := range txId {
		if b != 0 {
			return false
		}
	}
	return true
}

func UUID() string {
	newV7, err := uuid.NewV7()
	if err != nil {
		panic("should never happen")
	}
	return newV7.String()
}
func BtcToSat(value float64) int64 {
	return int64(value * 100000000)
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
	if strings.HasPrefix(strings.ToLower(id), "0x") {
		return id[2:]
	}
	return id
}
func NewRat() *big.Rat {
	rat := new(big.Rat)
	return rat
}

func ethTxesToBtcIds(txes ...[]*DbTx) []string {
	var btcIds []string
	for _, item := range txes {
		for _, tx := range item {
			btcIds = append(btcIds, tx.UtxoId)
		}
	}
	return btcIds
}

func txesToDbTxIds(txes []*DbTx) []string {
	uniqueList := NewUniqueList()
	for _, tx := range txes {
		uniqueList.Add(tx.Hash)
	}
	return uniqueList.List()
}

func txesToDbProofs(txes []*DbTx) []DbProof {
	var dbProofs []DbProof
	for _, tx := range txes {
		dbProof := DbProof{
			TxHash:    tx.Hash,
			ProofType: tx.ProofType,
		}
		if tx.Proved {
			dbProof.Status = int(common.ProofSuccess)
		}
		dbProofs = append(dbProofs, dbProof)
	}
	return dbProofs
}

func txesToUnGenProofs(txes []*DbTx) []*DbUnGenProof {
	var proofs []*DbUnGenProof
	for _, tx := range txes {
		if !tx.Proved {
			proofs = append(proofs, &DbUnGenProof{
				ChainType: tx.ChainType,
				ProofType: tx.ProofType,
				Hash:      tx.Hash,
				Height:    tx.Height,
				TxIndex:   tx.TxIndex,
				Amount:    uint64(tx.Amount),
			})
		}
	}
	return proofs
}

func mergeDbTxes(txes ...[]*DbTx) []*DbTx {
	var dbtxes []*DbTx
	for _, item := range txes {
		dbtxes = append(dbtxes, item...)
	}
	return dbtxes
}

func txesByAddrGroup(txes []*DbTx) map[string][]DbTx {
	txMap := make(map[string][]DbTx)
	for _, tx := range txes {
		if tx.Sender == "" {
			continue
		}
		list, ok := txMap[tx.Sender]
		if ok {
			list = append(list, DbTx{
				Hash:      tx.Hash,
				Height:    tx.Height,
				TxType:    tx.TxType,
				ChainType: tx.ChainType,
				Amount:    tx.Amount,
				TxIndex:   tx.TxIndex,
			})
			txMap[tx.Sender] = list
		} else {
			txMap[tx.Sender] = []DbTx{
				{
					Hash:      tx.Hash,
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
func PrintPanicStack(extras ...interface{}) {
	defer func() {
		if x := recover(); x != nil {
			logger.Error("caught panic in PrintPanicStack() %v", x)
		}
	}()
	if x := recover(); x != nil {
		logger.Error("%v", x)
		i := 0
		funcName, file, line, ok := runtime.Caller(i)
		for ok {
			logger.Error("frame %v:[func:%v,file:%v,line:%v]\n", i, runtime.FuncForPC(funcName).Name(), file, line)
			i++
			funcName, file, line, ok = runtime.Caller(i)
		}

		for k := range extras {
			logger.Error("EXRAS#%v DATA:%v\n", k, spew.Sdump(extras[k]))
		}
	}
}
