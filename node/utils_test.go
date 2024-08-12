package node

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/lightec-xyz/daemon/rpc"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestVerifySignature(t *testing.T) {
	publicKey, err := hex.DecodeString("02a3f345daecb93c23528849817a00762d611f385e84549974aec7c62ee8cae987")
	if err != nil {
		t.Fatal(err)
	}
	signature, err := hex.DecodeString("50eae8431e89a15bca4d0e72e92183dd2b4530d3e5efa6ba27938d771aa0ac5947ef5a1cbe9ea30e1da0f291465fa0369f4f8d108e8387946332adced8dbc9ae")
	if err != nil {
		t.Fatal(err)
	}
	verify := VerifySha256Signature(publicKey, signature, []byte("0000000000000007fc451dcb79d027d52c0cdd3f5a3e19f230a435dc72ad6a65"))
	t.Log(verify)
}

func VerifySha256Signature(publicKey []byte, signature []byte, data []byte) bool {
	hash := sha256.Sum256(data)
	fmt.Println(hexutil.Encode(hash[:]))
	verified := crypto.VerifySignature(publicKey, hash[:], signature)
	return verified
}

func TestEIP55Addr(t *testing.T) {

	addr := EIP55Addr("0x0ef907a4cfd17202d4035ff07f06030e2203da5a")
	t.Log(addr)
	addr = EIP55Addr("0x2A6443B5838f9524970c471289AB22f399395Ff6")
	t.Log(addr)
	addr = EIP55Addr("0x0ef907a4cfd172")
	t.Log(addr)
	addr = EIP55Addr("0x0ef907a4cfd17202d4035ff07f06030e20ef907a4cfd17202d4035ff07f06030e2")
	t.Log(addr)
	addr = EIP55Addr("0x")
	t.Log(addr)
	addr = EIP55Addr("")
	t.Log(addr)
}

func TestHandler(t *testing.T) {
	typeValue := reflect.TypeOf((*rpc.IAdmin)(nil)).Elem()
	for i := 0; i < typeValue.NumMethod(); i++ {
		fmt.Println(typeValue.Method(i).Name)
	}
}
