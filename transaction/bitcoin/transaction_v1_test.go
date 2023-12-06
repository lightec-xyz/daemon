package bitcoin

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"log"
	"strconv"
	"testing"
)

func TestDemo(t *testing.T) {
	//src:="e8c84a631D71E1Bb7083D3a82a3a74870a286B97"
	data, err := hex.DecodeString("6a14e8c84a631d71e1bb7083d3a82a3a74870a286b97")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(data)
	fmt.Printf("%v : %x\n", len(data[2:]), data[2:])

	result, err := hex.DecodeString("cbee12cf5411935db7ba6311a16c2e5b1aa7ac7d7562593312707fb343551117")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)
}

func TestTX(t *testing.T) {
	var msgTx wire.MsgTx
	hexData, err := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000007102000000011aae5c5a37f9003aaa12c63dcebdfcd0e5cb6d753c4265ec055d0697e5e0d6100100000000ffffffff026e86010000000000160014d7fae4fbdc8bf6c86a08c7177c5d06683754ea71ecdc7ee202000000160014fb5defb676e7f0a6711e3bc385849572a57fbe7e00000000000000000000000000000000000000")
	if err != nil {
		log.Fatal(err)
	}
	err = msgTx.Deserialize(bytes.NewReader(hexData))
	if err != nil {
		log.Fatal(err)
	}
	t.Log(msgTx.TxHash())

}

func TestTxDemo(t *testing.T) {
	//0000000000000000000000000000000000000000000000000000000000000020
	//0000000000000000000000000000000000000000000000000000000000000071
	//02000000011aae5c5a37f9003aaa12c63dcebdfcd0e5cb6d753c4265ec055d06
	//97e5e0d6100100000000ffffffff026e86010000000000160014d7fae4fbdc8b
	//f6c86a08c7177c5d06683754ea71ecdc7ee202000000160014fb5defb676e7f0
	//a6711e3bc385849572a57fbe7e00000000000000000000000000000000000000
	data := "0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000007102000000011aae5c5a37f9003aaa12c63dcebdfcd0e5cb6d753c4265ec055d0697e5e0d6100100000000ffffffff026e86010000000000160014d7fae4fbdc8bf6c86a08c7177c5d06683754ea71ecdc7ee202000000160014fb5defb676e7f0a6711e3bc385849572a57fbe7e00000000000000000000000000000000000000"
	hexData, err := hex.DecodeString(data)
	if err != nil {
		t.Fatal(err)
	}
	ver, err := strconv.ParseInt("0000000000000000000000000000000000000000000000000000000000000020", 16, 32)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(ver)
	l, err := strconv.ParseInt("0000000000000000000000000000000000000000000000000000000000000071", 16, 32)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(l)
	version := binary.LittleEndian.Uint32(hexData[0:32])
	length := binary.LittleEndian.Uint32(hexData[32:64])

	fmt.Println(version, length)
	txData := hexData[64 : 64+113]
	fmt.Printf("%v : %x\n", len(txData), txData)
	msgTx := wire.NewMsgTx(2)
	err = msgTx.Deserialize(bytes.NewReader(txData))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(msgTx.TxHash())

}

func TestLockScript(t *testing.T) {
	//address := "bcrt1qalv7aduqdpz9wc4fut3nt44tsf42anleed76yj3el3rgd4rgldvq2aw6ze"
	address := "bcrt1qalv7aduqdpz9wc4fut3nt44tsf42anleed76yj3el3rgd4rgldvq2aw6ze"
	addr, err := btcutil.DecodeAddress(address, &chaincfg.RegressionNetParams)
	lockingScript, err := txscript.PayToAddrScript(addr)
	if err != nil {
		fmt.Println("Error creating locking script:", err)
		return
	}
	fmt.Printf("%x\n", lockingScript)
}

func Test0001(t *testing.T) {

}
