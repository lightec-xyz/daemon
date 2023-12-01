package node

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
)

func TxIdToProofId(txId string) string {
	pTxID := fmt.Sprintf("%s%s", ProofPrefix, txId)
	return pTxID
}

func getEthAddrFromScript(script string) (string, error) {
	//6a14e8c84a631d71e1bb7083d3a82a3a74870a286b97
	//todo
	data, err := hex.DecodeString(script)
	if err != nil {
		return "", err
	}
	data = data[2:int(data[1])]
	return hex.EncodeToString(data), nil
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
