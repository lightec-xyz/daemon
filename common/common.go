package common

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/consensys/gnark/frontend"
	"github.com/google/uuid"
	"github.com/lightec-xyz/provers/circuits/fabric/receipt-proof"
	"github.com/lightec-xyz/provers/circuits/fabric/tx-proof"
	"os"
	"path/filepath"
	"reflect"
)

func HexToBytes(data string) ([]byte, error) {
	if data[0:2] == "0x" {
		data = data[2:]
	}
	bytes, err := hex.DecodeString(data)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func Uuid() (string, error) {
	newV7, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return newV7.String(), nil
}
func MustUUID() string {
	newV7, err := uuid.NewV7()
	if err != nil {
		panic("gen uuid error,should never happen")
	}
	return newV7.String()
}

func objToJson(obj interface{}) string {
	ojbBytes, err := json.Marshal(obj)
	if err != nil {
		return "error obj to josn"
	}
	return string(ojbBytes)
}

func ParseObj(src, dst interface{}) error {
	if reflect.ValueOf(dst).Kind() != reflect.Ptr {
		return fmt.Errorf("dst must be a pointer")
	}
	srcBytes, err := json.Marshal(src)
	if err != nil {
		return err
	}
	err = json.Unmarshal(srcBytes, dst)
	if err != nil {
		return err
	}
	return nil
}

func WriteFile(path string, data []byte) error {
	err := CheckOrCreateDir(path)
	if err != nil {
		return err
	}
	exists, err := FileExists(path)
	if err != nil {
		return err
	}
	if exists {
		err := os.Remove(path)
		if err != nil {
			return err
		}
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	err = file.Sync()
	if err != nil {
		return err
	}
	return nil
}
func CheckOrCreateDir(path string) error {
	dir := filepath.Dir(path)
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, fmt.Errorf("stat error: %v", err)
}

func TxVarToHex(data *[tx.MaxTxUint128Len]frontend.Variable) (string, error) {
	var bytes []byte
	for _, v := range data {
		bytes = append(bytes, v.(byte))
	}
	return hex.EncodeToString(bytes), nil

}

func HexToTxVar(data string) (*[tx.MaxTxUint128Len]frontend.Variable, error) {
	bytes, err := hex.DecodeString(data)
	if err != nil {
		return nil, err
	}
	var res [tx.MaxTxUint128Len]frontend.Variable
	for k, v := range bytes {
		res[k] = v
	}
	return &res, nil
}

func ReceiptVarToHex(data *[receipt.MaxReceiptUint128Len]frontend.Variable) (string, error) {
	var bytes []byte
	for _, v := range data {
		bytes = append(bytes, v.(byte))
	}
	return hex.EncodeToString(bytes), nil
}

func HexToReceiptVar(data string) (*[receipt.MaxReceiptUint128Len]frontend.Variable, error) {
	bytes, err := hex.DecodeString(data)
	if err != nil {
		return nil, err
	}
	var res [receipt.MaxReceiptUint128Len]frontend.Variable
	for k, v := range bytes {
		res[k] = v
	}
	return &res, nil
}
