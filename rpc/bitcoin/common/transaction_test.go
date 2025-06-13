package bitcoin

import (
	"encoding/hex"
	"testing"
)

func TestBuildTx(t *testing.T) {
	txRaw, _ := hex.DecodeString("02000000012A2553D1444B10ED6677CB860A1D9F44D34D2C5B45D6012A016F59E774E5AD1D0000000000FFFFFFFF02EE22000000000000160014C3CCAFD4C7D930E4A5539882A55737039303227736000000000000002200201CFF3AE65030961B7CC68683638F689AA1110F6C2D537237D2D9516C8207CEA600000000")
	builder := NewMultiTransactionBuilder()
	err := builder.Deserialize(txRaw)
	if err != nil {
		panic(err)
	}
	for _, in := range builder.MsgTx.TxIn {
		t.Logf("txId:%v,index:%v,script:%x,secquence:%v\n", in.PreviousOutPoint.Hash.String(), in.PreviousOutPoint.Index,
			in.SignatureScript, in.Sequence)
	}
	for _, out := range builder.MsgTx.TxOut {
		t.Logf("value:%v,pkScript:%x\n", out.Value, out.PkScript)
	}
	t.Log(builder.TxHash())
}
