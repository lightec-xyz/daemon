package node

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	btctx "github.com/lightec-xyz/daemon/rpc/bitcoin/common"
	"testing"
)

func TestClient_NewBtcSignatures(t *testing.T) {
	btcRawTx := ethcommon.FromHex("0x020000000127d9cb550399f553ef5e18588716a4945296e0ef19030322fac195daa535bf9a0000000000ffffffff026261000000000000160014ca12a4e423ae13a604b20de38d49c7b9abc2504dfa0a000000000000220020e0381fbba457e811b16a296457be1486e21d09c51f9b4f98f573c08c09b5806900000000")
	transaction := btctx.NewMultiTransactionBuilder()
	err := transaction.Deserialize(btcRawTx)
	if err != nil {
		t.Fatal(err)
	}
	var signatues [][][]byte
	sgx := []string{"3045022100a8d1dde756ab5f5749c07552054d4793595861969549bae9808ad3e595abe6c402203c8618064bcb3478f346812b41ba7543b1230cfc2b8b54d550ab44837cb008a2"}
	oasis := []string{"3045022100839a63e1f1e9957ce3da5e9b85359927076f1269e697d94be9d24d1208ffbe4802206a61e662a343fefc29830c71407198a60500744e0c9f8c71123583c1315b4ce701"}

	sgxBytes, err := sgxSigToBytes(sgx)
	if err != nil {
		t.Fatal(err)
	}
	signatues = append(signatues, sgxBytes)
	var oasisByts [][]byte
	for _, s := range oasis {
		oasisByts = append(oasisByts, ethcommon.FromHex(s))
	}
	signatues = append(signatues, oasisByts)
	multiSigScriptBytes := ethcommon.FromHex(BtcMultiSig)
	err = transaction.AddMultiScript(multiSigScriptBytes, 2, 3)
	if err != nil {
		t.Fatal(err)
	}
	err = transaction.MergeSignature(signatues)
	if err != nil {
		t.Fatal(err)
	}
	btcTxBytes, err := transaction.Serialize()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%x", btcTxBytes)

}
