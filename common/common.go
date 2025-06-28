package common

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	native_plonk "github.com/consensys/gnark/backend/plonk"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/google/uuid"
)

func Md5(data []byte) []byte {
	ret := md5.Sum(data)
	return ret[:]
}
func HexMd5(data []byte) string {
	return hex.EncodeToString(Md5(data))
}

func ReadObj(path string, obj interface{}) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, obj)
	if err != nil {
		return err
	}
	return nil
}

func StrEqual(a, b string) bool {
	if strings.ToLower(TrimOx(a)) == strings.ToLower(TrimOx(b)) {
		return true
	}
	return false
}

func TrimOx(value string) string {
	if strings.HasPrefix(value, "0x") || strings.HasPrefix(value, "0X") {
		return value[2:]
	}
	return value
}

func ReverseU32(input []uint32) []uint32 {
	b := make([]uint32, len(input))
	copy(b, input)
	for i := 0; i < len(b)/2; i++ {
		b[i], b[len(b)-i-1] = b[len(b)-i-1], b[i]
	}
	return b
}

func ReverseBytes(data []byte) []byte {
	res := make([]byte, len(data))
	copy(res, data)
	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}
	return res
}

func MustUUID() string {
	newV7, err := uuid.NewV7()
	if err != nil {
		panic("gen uuid error,should never happen")
	}
	return newV7.String()
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
func ParseWithNumber(src, dst interface{}) error {
	if reflect.ValueOf(dst).Kind() != reflect.Ptr {
		return fmt.Errorf("dst must be a pointer")
	}
	srcBytes, err := json.Marshal(src)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewReader(srcBytes))
	decoder.UseNumber()
	err = decoder.Decode(dst)
	if err != nil {
		return err
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

func GetEnvDebugMode() bool {
	zkDebugEnv := os.Getenv(ZkDebugEnv)
	if zkDebugEnv == "1" {
		return true
	}
	return false
}

func GetEnvZkProofTypes() string {
	zkProofTypes := os.Getenv(ZkProofTypes)
	return zkProofTypes
}

func BytesArrayToHex(arr [][32]byte) []string {
	hexArr := make([]string, 0)
	for _, bytes := range arr {
		hexArr = append(hexArr, hex.EncodeToString(bytes[:]))
	}
	return hexArr
}

type CircuitFile struct {
	File string `json:"file"`
	Md5  string `json:"md5"`
}

func GetCircuitMd5(url string) ([]CircuitFile, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var circuitFiles []CircuitFile
	err = json.Unmarshal(body, &circuitFiles)
	if err != nil {
		return nil, err
	}
	return circuitFiles, nil
}

func ExportSolidity(solFile string, vk native_plonk.VerifyingKey) error {
	openFile, err := OverwriteFile(solFile)
	if err != nil {
		return err
	}
	defer openFile.Close()
	err = vk.ExportSolidity(openFile)
	if err != nil {
		return err
	}
	return nil
}

func OverwriteFile(file string) (*os.File, error) {
	dir := filepath.Dir(file)
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	exists, err := FileExists(file)
	if err != nil {
		return nil, err
	}
	if exists {
		err := os.Remove(file)
		if err != nil {
			return nil, err
		}
	}
	fFile, err := os.Create(file)
	if err != nil {
		return nil, err
	}
	return fFile, nil
}

func GetBtcPwd() string {
	return os.Getenv("btcPwd")
}

func GetBtcUser() string {
	return os.Getenv("btcUser")
}

func GetBtcUrl() string {
	return os.Getenv("btcUrl")
}

func SetUpDir() string {
	return os.Getenv("setupDir")
}
