package bitcoin

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TXPrevOutFetcher(tx *wire.MsgTx, prevPkScripts [][]byte,
	inputValues []btcutil.Amount) (*txscript.MultiPrevOutFetcher, error) {

	if len(tx.TxIn) != len(prevPkScripts) {
		return nil, errors.New("tx.TxIn and prevPkScripts slices " +
			"must have equal length")
	}
	if len(tx.TxIn) != len(inputValues) {
		return nil, errors.New("tx.TxIn and inputValues slices " +
			"must have equal length")
	}

	fetcher := txscript.NewMultiPrevOutFetcher(nil)
	for idx, txin := range tx.TxIn {
		fetcher.AddPrevOut(txin.PreviousOutPoint, &wire.TxOut{
			Value:    int64(inputValues[idx]),
			PkScript: prevPkScripts[idx],
		})
	}

	return fetcher, nil
}

func spendNestedWitnessPubKeyHash(txIn *wire.TxIn, pkScript []byte,
	inputValue int64, chainParams *chaincfg.Params, privKey *btcec.PrivateKey,
	tx *wire.MsgTx, hashCache *txscript.TxSigHashes, idx int) error {

	// First we need to obtain the key pair related to this p2sh output.
	//_, addrs, _, err := txscript.ExtractPkScriptAddrs(pkScript,
	//	chainParams)
	//if err != nil {
	//	return err
	//}
	//privKey, compressed := privKey, true
	//if err != nil {
	//	return err
	//}
	pubKey := privKey.PubKey()

	var pubKeyHash []byte
	compressed := true
	if compressed {
		pubKeyHash = btcutil.Hash160(pubKey.SerializeCompressed())
	} else {
		pubKeyHash = btcutil.Hash160(pubKey.SerializeUncompressed())
	}

	// Next, we'll generate a valid sigScript that'll allow us to spend
	// the p2sh output. The sigScript will contain only a single push of
	// the p2wkh witness program corresponding to the matching public key
	// of this address.
	p2wkhAddr, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, chainParams)
	if err != nil {
		return err
	}
	witnessProgram, err := txscript.PayToAddrScript(p2wkhAddr)
	if err != nil {
		return err
	}
	bldr := txscript.NewScriptBuilder()
	bldr.AddData(witnessProgram)
	sigScript, err := bldr.Script()
	if err != nil {
		return err
	}
	txIn.SignatureScript = sigScript

	// With the sigScript in place, we'll next generate the proper witness
	// that'll allow us to spend the p2wkh output.
	witnessScript, err := txscript.WitnessSignature(tx, hashCache, idx,
		inputValue, witnessProgram, txscript.SigHashAll, privKey, compressed)
	if err != nil {
		return err
	}

	txIn.Witness = witnessScript

	return nil
}

// spendWitnessKeyHash generates, and sets a valid witness for spending the
// passed pkScript with the specified input amount. The input amount *must*
// correspond to the output value of the previous pkScript, or else verification
// will fail since the new sighash digest algorithm defined in BIP0143 includes
// the input value in the sighash.
func spendWitnessKeyHash(txIn *wire.TxIn, pkScript []byte,
	inputValue int64, chainParams *chaincfg.Params, privKey *btcec.PrivateKey,
	tx *wire.MsgTx, hashCache *txscript.TxSigHashes, idx int) error {

	// First obtain the key pair associated with this p2wkh address.
	//_, addrs, _, err := txscript.ExtractPkScriptAddrs(pkScript,
	//	chainParams)
	//if err != nil {
	//	return err
	//}
	//privKey, compressed, err := secrets.GetKey(addrs[0])
	//if err != nil {
	//	return err
	//}
	pubKey := privKey.PubKey()
	compressed := true
	// Once we have the key pair, generate a p2wkh address type, respecting
	// the compression type of the generated key.
	var pubKeyHash []byte
	if compressed {
		pubKeyHash = btcutil.Hash160(pubKey.SerializeCompressed())
	} else {
		pubKeyHash = btcutil.Hash160(pubKey.SerializeUncompressed())
	}
	p2wkhAddr, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, chainParams)
	if err != nil {
		return err
	}

	// With the concrete address type, we can now generate the
	// corresponding witness program to be used to generate a valid witness
	// which will allow us to spend this output.
	witnessProgram, err := txscript.PayToAddrScript(p2wkhAddr)
	if err != nil {
		return err
	}
	witnessScript, err := txscript.WitnessSignature(tx, hashCache, idx,
		inputValue, witnessProgram, txscript.SigHashAll, privKey, true)
	if err != nil {
		return err
	}

	txIn.Witness = witnessScript

	return nil
}

// spendTaprootKey generates, and sets a valid witness for spending the passed
// pkScript with the specified input amount. The input amount *must*
// correspond to the output value of the previous pkScript, or else verification
// will fail since the new sighash digest algorithm defined in BIP0341 includes
// the input value in the sighash.
func spendTaprootKey(txIn *wire.TxIn, pkScript []byte,
	inputValue int64, chainParams *chaincfg.Params, privKey *btcec.PrivateKey,
	tx *wire.MsgTx, hashCache *txscript.TxSigHashes, idx int) error {

	// First obtain the key pair associated with this p2tr address. If the
	// pkScript is incorrect or derived from a different internal key or
	// with a script root, we simply won't find a corresponding private key
	// here.
	//_, addrs, _, err := txscript.ExtractPkScriptAddrs(pkScript, chainParams)
	//if err != nil {
	//	return err
	//}
	//privKey, _, err := secrets.GetKey(addrs[0])
	//if err != nil {
	//	return err
	//}

	// We can now generate a valid witness which will allow us to spend this
	// output.
	witnessScript, err := txscript.TaprootWitnessSignature(
		tx, hashCache, idx, inputValue, pkScript,
		txscript.SigHashDefault, privKey,
	)
	if err != nil {
		return err
	}

	txIn.Witness = witnessScript

	return nil
}

func ValidateMsgTx(tx *wire.MsgTx, prevScripts [][]byte,
	inputValues []btcutil.Amount) error {

	inputFetcher, err := TXPrevOutFetcher(
		tx, prevScripts, inputValues,
	)
	if err != nil {
		return err
	}

	hashCache := txscript.NewTxSigHashes(tx, inputFetcher)
	for i, prevScript := range prevScripts {
		vm, err := txscript.NewEngine(
			prevScript, tx, i, txscript.StandardVerifyFlags, nil,
			hashCache, int64(inputValues[i]), inputFetcher,
		)
		if err != nil {
			return fmt.Errorf("cannot create script engine: %s", err)
		}
		err = vm.Execute()
		if err != nil {
			return fmt.Errorf("cannot validate transaction: %s", err)
		}
	}
	return nil
}

func CreateTx(tx *wire.MsgTx, prevPkScripts [][]byte,
	inputValues []btcutil.Amount, privKeys []*btcec.PrivateKey, chainParams *chaincfg.Params) error {

	inputFetcher, err := TXPrevOutFetcher(tx, prevPkScripts, inputValues)
	if err != nil {
		return err
	}

	inputs := tx.TxIn
	hashCache := txscript.NewTxSigHashes(tx, inputFetcher)

	if len(inputs) != len(prevPkScripts) {
		return errors.New("tx.TxIn and prevPkScripts slices must " +
			"have equal length")
	}

	if len(inputs) != len(privKeys) {
		return errors.New("tx.TxIn and privKeys slices must " +
			"have equal length")
	}

	for i := range inputs {
		pkScript := prevPkScripts[i]

		switch {
		// If this is a p2sh output, who's script hash pre-image is a
		// witness program, then we'll need to use a modified signing
		// function which generates both the sigScript, and the witness
		// 如果前一个pkScript是P2SH
		case txscript.IsPayToScriptHash(pkScript):
			err := spendNestedWitnessPubKeyHash(
				inputs[i], pkScript, int64(inputValues[i]),
				chainParams, privKeys[i], tx, hashCache, i,
			)
			if err != nil {
				return err
			}
			// 如果前一个pkScript是 P2WPKH
		case txscript.IsPayToWitnessPubKeyHash(pkScript):
			err := spendWitnessKeyHash(
				inputs[i], pkScript, int64(inputValues[i]),
				chainParams, privKeys[i], tx, hashCache, i,
			)
			if err != nil {
				return err
			}

		case txscript.IsPayToTaproot(pkScript):
			err := spendTaprootKey(
				inputs[i], pkScript, int64(inputValues[i]),
				chainParams, privKeys[i], tx, hashCache, i,
			)
			if err != nil {
				return err
			}

		default:
			sigScript := inputs[i].SignatureScript //
			lookupKey := func(a btcutil.Address) (*btcec.PrivateKey, bool, error) {
				return privKeys[i], true, nil
			}
			script, err := txscript.SignTxOutput(chainParams, tx, i,
				pkScript, txscript.SigHashAll, txscript.KeyClosure(lookupKey), nil,
				sigScript)
			if err != nil {
				return err
			}
			inputs[i].SignatureScript = script
		}
	}
	return nil
}

// extractWitnessPubKeyHash extracts the witness public key hash from the passed
// script if it is a standard pay-to-witness-pubkey-hash script. It will return
// nil otherwise.
const (
	witnessV0PubKeyHashLen = 22
	witnessV0ScriptHashLen = 34
	sigHashMask            = 0x1f
)

func extractWitnessPubKeyHash(script []byte) []byte {
	// A pay-to-witness-pubkey-hash script is of the form:
	//   OP_0 OP_DATA_20 <20-byte-hash>
	if len(script) == witnessV0PubKeyHashLen &&
		script[0] == txscript.OP_0 &&
		script[1] == txscript.OP_DATA_20 {

		return script[2:witnessV0PubKeyHashLen]
	}

	return nil
}

func extractWitnessV0ScriptHash(script []byte) []byte {
	// A pay-to-witness-script-hash script is of the form:
	//   OP_0 OP_DATA_32 <32-byte-hash>
	if len(script) == witnessV0ScriptHashLen &&
		script[0] == txscript.OP_0 &&
		script[1] == txscript.OP_DATA_32 {

		return script[2:34]
	}

	return nil
}

// calcWitnessSignatureHashRaw 从btcd 拷贝而来，增加了对pwsh的支持
func calcWitnessSignatureHashRaw(subScript []byte, sigHashes *txscript.TxSigHashes,
	hashType txscript.SigHashType, tx *wire.MsgTx, idx int, amt int64, origScript []byte) ([]byte, error) {

	// As a sanity check, ensure the passed input index for the transaction
	// is valid.
	//
	// TODO(roasbeef): check needs to be lifted elsewhere?
	if idx > len(tx.TxIn)-1 {
		return nil, fmt.Errorf("idx %d but %d txins", idx, len(tx.TxIn))
	}

	// We'll utilize this buffer throughout to incrementally calculate
	// the signature hash for this transaction.
	var sigHash bytes.Buffer

	// First write out, then encode the transaction's version number.
	var bVersion [4]byte
	binary.LittleEndian.PutUint32(bVersion[:], uint32(tx.Version))
	sigHash.Write(bVersion[:])

	// Next write out the possibly pre-calculated hashes for the sequence
	// numbers of all inputs, and the hashes of the previous outs for all
	// outputs.
	var zeroHash chainhash.Hash

	// If anyone can pay isn't active, then we can use the cached
	// hashPrevOuts, otherwise we just write zeroes for the prev outs.
	if hashType&txscript.SigHashAnyOneCanPay == 0 {
		sigHash.Write(sigHashes.HashPrevOutsV0[:])
	} else {
		sigHash.Write(zeroHash[:])
	}

	// If the sighash isn't anyone can pay, single, or none, the use the
	// cached hash sequences, otherwise write all zeroes for the
	// hashSequence.
	if hashType&txscript.SigHashAnyOneCanPay == 0 &&
		hashType&sigHashMask != txscript.SigHashSingle &&
		hashType&sigHashMask != txscript.SigHashNone {
		sigHash.Write(sigHashes.HashSequenceV0[:])
	} else {
		sigHash.Write(zeroHash[:])
	}

	txIn := tx.TxIn[idx]

	// Next, write the outpoint being spent.
	sigHash.Write(txIn.PreviousOutPoint.Hash[:])
	var bIndex [4]byte
	binary.LittleEndian.PutUint32(bIndex[:], txIn.PreviousOutPoint.Index)
	sigHash.Write(bIndex[:])

	switch {
	case txscript.IsPayToWitnessPubKeyHash(subScript):
		sigHash.Write([]byte{0x19})
		sigHash.Write([]byte{txscript.OP_DUP})
		sigHash.Write([]byte{txscript.OP_HASH160})
		sigHash.Write([]byte{txscript.OP_DATA_20})
		sigHash.Write(extractWitnessPubKeyHash(subScript))
		sigHash.Write([]byte{txscript.OP_EQUALVERIFY})
		sigHash.Write([]byte{txscript.OP_CHECKSIG})

	case txscript.IsPayToWitnessScriptHash(subScript):
		if origScript == nil {
			panic("origScript is nil")
		}
		// TODO(keep), 检查subScript != hash(original script), encode original script
		calpkScript := sha256.Sum256(origScript)

		if bytes.Compare(calpkScript[:], extractWitnessV0ScriptHash(subScript)) != 0 {
			panic("sha256.Sum256(origScript) != subScript")
		}
		wire.WriteVarBytes(&sigHash, 0, origScript)

	}
	// Next, add the input amount, and sequence number of the input being
	// signed.
	var bAmount [8]byte
	binary.LittleEndian.PutUint64(bAmount[:], uint64(amt))
	sigHash.Write(bAmount[:])
	var bSequence [4]byte
	binary.LittleEndian.PutUint32(bSequence[:], txIn.Sequence)
	sigHash.Write(bSequence[:])

	// If the current signature mode isn't single, or none, then we can
	// re-use the pre-generated hashoutputs sighash fragment. Otherwise,
	// we'll serialize and add only the target output index to the signature
	// pre-image.
	if hashType&sigHashMask != txscript.SigHashSingle &&
		hashType&sigHashMask != txscript.SigHashNone {
		sigHash.Write(sigHashes.HashOutputsV0[:])
	} else if hashType&sigHashMask == txscript.SigHashSingle && idx < len(tx.TxOut) {
		var b bytes.Buffer
		wire.WriteTxOut(&b, 0, 0, tx.TxOut[idx])
		sigHash.Write(chainhash.DoubleHashB(b.Bytes()))
	} else {
		sigHash.Write(zeroHash[:])
	}

	// Finally, write out the transaction's locktime, and the sig hash
	// type.
	var bLockTime [4]byte
	binary.LittleEndian.PutUint32(bLockTime[:], tx.LockTime)
	sigHash.Write(bLockTime[:])
	var bHashType [4]byte
	binary.LittleEndian.PutUint32(bHashType[:], uint32(hashType))
	sigHash.Write(bHashType[:])
	return chainhash.DoubleHashB(sigHash.Bytes()), nil
}

// checkScriptParses returns an error if the provided script fails to parse.
func checkScriptParses(scriptVersion uint16, script []byte) error {
	tokenizer := txscript.MakeScriptTokenizer(scriptVersion, script)
	for tokenizer.Next() {
		// Nothing to do.
	}
	return tokenizer.Err()
}

func CalcWitnessSigHash(txInLockingScript []byte, sigHashes *txscript.TxSigHashes, hType txscript.SigHashType,
	tx *wire.MsgTx, idx int, amt int64, txInOriginalScript []byte) ([]byte, error) {

	const scriptVersion = 0
	if err := checkScriptParses(scriptVersion, txInLockingScript); err != nil {
		return nil, err
	}

	return calcWitnessSignatureHashRaw(txInLockingScript, sigHashes, hType, tx, idx, amt, txInOriginalScript)
}

func calWitnessSigHashForNestedWitnessPubKeyHash(txIn *wire.TxIn, pkScript []byte,
	inputValue int64, chainParams *chaincfg.Params, from btcutil.Address,
	tx *wire.MsgTx, hashCache *txscript.TxSigHashes, idx int) ([]byte, error) {

	//pubKey := privKey.PubKey()
	//var pubKeyHash []byte
	//compressed := true
	//if compressed {
	//	pubKeyHash = btcutil.Hash160(pubKey.SerializeCompressed())
	//} else {
	//	pubKeyHash = btcutil.Hash160(pubKey.SerializeUncompressed())
	//}
	//
	//// Next, we'll generate a valid sigScript that'll allow us to spend
	//// the p2sh output. The sigScript will contain only a single push of
	//// the p2wkh witness program corresponding to the matching public key
	//// of this address.
	//p2wkhAddr, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, chainParams)
	//if err != nil {
	//	return nil, err
	//}
	witnessProgram, err := txscript.PayToAddrScript(from)
	if err != nil {
		return nil, err
	}
	//bldr := txscript.NewScriptBuilder()
	//bldr.AddData(witnessProgram)
	//sigScript, err := bldr.Script()
	//if err != nil {
	//	return nil, err
	//}
	//txIn.SignatureScript = sigScript
	////
	// With the sigScript in place, we'll next generate the proper witness
	// that'll allow us to spend the p2wkh output.

	hash, err := CalcWitnessSigHash(witnessProgram, hashCache, txscript.SigHashAll, tx, idx, inputValue, nil)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func calWitnessSigHashForWitnessKeyHash(txIn *wire.TxIn, pkScript []byte,
	inputValue int64, chainParams *chaincfg.Params, from btcutil.Address,
	tx *wire.MsgTx, hashCache *txscript.TxSigHashes, idx int) ([]byte, error) {

	// First obtain the key pair associated with this p2wkh address.
	//_, addrs, _, err := txscript.ExtractPkScriptAddrs(pkScript,
	//	chainParams)
	//if err != nil {
	//	return err
	//}
	//privKey, compressed, err := secrets.GetKey(addrs[0])
	//if err != nil {
	//	return err
	//}
	//pubKey := privKey.PubKey()
	//compressed := true
	//// Once we have the key pair, generate a p2wkh address type, respecting
	//// the compression type of the generated key.
	//var pubKeyHash []byte
	//if compressed {
	//	pubKeyHash = btcutil.Hash160(pubKey.SerializeCompressed())
	//} else {
	//	pubKeyHash = btcutil.Hash160(pubKey.SerializeUncompressed())
	//}
	//p2wkhAddr, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, chainParams)
	//if err != nil {
	//	return nil, err
	//}

	// With the concrete address type, we can now generate the
	// corresponding witness program to be used to generate a valid witness
	// which will allow us to spend this output.
	txInLockScript, err := txscript.PayToAddrScript(from)
	if err != nil {
		return nil, err
	}
	hash, err := CalcWitnessSigHash(txInLockScript, hashCache, txscript.SigHashAll, tx, idx, inputValue, nil)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func calWitnessSigHashForWitnessScriptHash(txIn *wire.TxIn, pkScript []byte,
	inputValue int64, chainParams *chaincfg.Params, from btcutil.Address,
	tx *wire.MsgTx, hashCache *txscript.TxSigHashes, idx int, txInOriginalScript []byte) ([]byte, error) {

	// First obtain the key pair associated with this p2wkh address.
	//_, addrs, _, err := txscript.ExtractPkScriptAddrs(pkScript,
	//	chainParams)
	//if err != nil {
	//	return err
	//}
	//privKey, compressed, err := secrets.GetKey(addrs[0])
	//if err != nil {
	//	return err
	//}
	//pubKey := privKey.PubKey()
	//compressed := true
	//// Once we have the key pair, generate a p2wkh address type, respecting
	//// the compression type of the generated key.
	//var pubKeyHash []byte
	//if compressed {
	//	pubKeyHash = btcutil.Hash160(pubKey.SerializeCompressed())
	//} else {
	//	pubKeyHash = btcutil.Hash160(pubKey.SerializeUncompressed())
	//}
	//p2wkhAddr, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, chainParams)
	//if err != nil {
	//	return nil, err
	//}

	// With the concrete address type, we can now generate the
	// corresponding witness program to be used to generate a valid witness
	// which will allow us to spend this output.
	txInLockingScript, err := txscript.PayToAddrScript(from)
	if err != nil {
		return nil, err
	}
	hash, err := CalcWitnessSigHash(txInLockingScript, hashCache, txscript.SigHashAll, tx, idx, inputValue, txInOriginalScript)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

// sigs must be sorted according pubkeys sequence
func MergeMultiSignatures(required int, multiSigScript []byte, sigs [][]byte) (wire.TxWitness, error) {
	if len(sigs) < required {
		return nil, fmt.Errorf("not enough signatures")
	}

	witnessElements := make(wire.TxWitness, 0, required+2)
	witnessElements = append(witnessElements, nil)
	for i := 0; i < required; i++ {
		witnessElements = append(witnessElements, sigs[i])
	}
	witnessElements = append(witnessElements, multiSigScript)

	return witnessElements, nil
}

// func CalcWitnessSigHash(script []byte, sigHashes *TxSigHashes, hType SigHashType, tx *wire.MsgTx, idx int, amt int64) ([]byte, error) {
func CalWitnessSigHash(tx *wire.MsgTx, txInLockingScripts [][]byte,
	inputValues []btcutil.Amount, from btcutil.Address, chainParams *chaincfg.Params, txInOriginalScripts [][]byte) ([][]byte, error) {

	inputFetcher, err := TXPrevOutFetcher(tx, txInLockingScripts, inputValues)
	if err != nil {
		return nil, err
	}

	inputs := tx.TxIn
	hashCache := txscript.NewTxSigHashes(tx, inputFetcher)

	if len(inputs) != len(txInLockingScripts) {
		return nil, errors.New("tx.TxIn and txInLockingScripts slices must " +
			"have equal length")
	}

	if len(inputs) != len(txInOriginalScripts) {
		return nil, errors.New("tx.TxIn and txInOriginalScripts slices must " +
			"have equal length")
	}

	//if len(inputs) != len(privKeys) {
	//	return nil, errors.New("tx.TxIn and privKeys slices must " +
	//		"have equal length")
	//}

	var hashes [][]byte
	for i := range inputs {
		pkScript := txInLockingScripts[i]
		switch {
		// If this is a p2sh output, who's script hash pre-image is a
		// witness program, then we'll need to use a modified signing
		// function which generates both the sigScript, and the witness
		// 如果前一个pkScript是P2SH
		case txscript.IsPayToScriptHash(pkScript):
			hash, err := calWitnessSigHashForNestedWitnessPubKeyHash(
				inputs[i], pkScript, int64(inputValues[i]),
				chainParams, from, tx, hashCache, i,
			)
			if err != nil {
				return nil, err
			}
			hashes = append(hashes, hash)
		case txscript.IsPayToWitnessPubKeyHash(pkScript):
			hash, err := calWitnessSigHashForWitnessKeyHash(
				inputs[i], pkScript, int64(inputValues[i]),
				chainParams, from, tx, hashCache, i,
			)
			if err != nil {
				return nil, err
			}
			hashes = append(hashes, hash)

		case txscript.IsPayToWitnessScriptHash(pkScript):
			hash, err := calWitnessSigHashForWitnessScriptHash(
				inputs[i], pkScript, int64(inputValues[i]),
				chainParams, from, tx, hashCache, i, txInOriginalScripts[i],
			)
			if err != nil {
				return nil, err
			}
			hashes = append(hashes, hash)

		case txscript.IsPayToTaproot(pkScript):
			panic("unsupported")

		default:
			panic("unsupported")
		}
	}
	return hashes, nil
}

func TestCreateDepositTx_1(t *testing.T) {
	secret, _ := hex.DecodeString("b29c3157e4a68b240ec821515fc77181c7a828259efbb3c1ab1df9b67d03c645")
	privKey, pubKey := btcec.PrivKeyFromBytes(secret)
	//fmt.Printf("priv: %v\n", privKey)
	//fmt.Printf("pub: %v\n", hex.EncodeToString(pubKey.SerializeCompressed()))
	from, _ := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(pubKey.SerializeCompressed()), &chaincfg.TestNet3Params)
	fmt.Printf("from addr:%v, from sriptAddress:%v\n", from.EncodeAddress(), hex.EncodeToString(from.ScriptAddress()))

	// Create the transaction to redeem the fake transaction.
	to, _ := btcutil.DecodeAddress("tb1q7yc8ncrxy6wsdlhvhd6gglpfatg07835uses5mpsc2rfv7zulhcqy0m979", &chaincfg.TestNet3Params)
	fmt.Printf("to address:%v, to sriptAddress:%v\n", to.String(), hex.EncodeToString(to.ScriptAddress()))

	depositTx := wire.NewMsgTx(wire.TxVersion)

	hash, err := chainhash.NewHashFromStr("6658fcd6da67b838a7405c1a6269423e3c0b09787ff96fabde57c2e84c8b5c48")
	txIn := wire.NewTxIn(wire.NewOutPoint(hash, 0), nil, nil)
	depositTx.AddTxIn(txIn)

	//TxIn's pkScript and value
	txInPkScript, err := hex.DecodeString("0014f97a2ead90717062357c8c1ee15d3ed0a5324efd")
	txInValue := btcutil.Amount(10000)

	//
	txOutScript, err := txscript.PayToAddrScript(to)
	if err != nil {
		fmt.Println(err)
		return
	}
	txOut := wire.NewTxOut(9000, txOutScript)
	depositTx.AddTxOut(txOut)

	var buf bytes.Buffer
	depositTx.Serialize(&buf)
	fmt.Printf("before sign deposit TxHash: %v\n", depositTx.TxHash())
	fmt.Printf("before sign deposit WitnessHash: %v\n", depositTx.WitnessHash())
	fmt.Printf("before sign deposit: %v\n", hex.EncodeToString(buf.Bytes()))

	err = CreateTx(depositTx, [][]byte{txInPkScript}, []btcutil.Amount{txInValue}, []*btcec.PrivateKey{privKey}, &chaincfg.TestNet3Params)
	assert.NoError(t, err)

	err = ValidateMsgTx(depositTx, [][]byte{txInPkScript}, []btcutil.Amount{txInValue})
	assert.NoError(t, err)

	buf.Reset()
	depositTx.Serialize(&buf)
	fmt.Printf("after sign deposit TxHash: %v\n", depositTx.TxHash())
	fmt.Printf("after sign deposit WitnessHash: %v\n", depositTx.WitnessHash())
	fmt.Printf("after sign deposit: %v\n", hex.EncodeToString(buf.Bytes()))
}

func TestCreateDepositTx_2(t *testing.T) {
	secret, _ := hex.DecodeString("b29c3157e4a68b240ec821515fc77181c7a828259efbb3c1ab1df9b67d03c645")
	privKey, pubKey := btcec.PrivKeyFromBytes(secret)
	from, _ := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(pubKey.SerializeCompressed()), &chaincfg.TestNet3Params)
	fmt.Printf("from addr:%v, from sriptAddress:%v\n", from.EncodeAddress(), hex.EncodeToString(from.ScriptAddress()))

	// Create the transaction to redeem the fake transaction.
	to, _ := btcutil.DecodeAddress("tb1q7yc8ncrxy6wsdlhvhd6gglpfatg07835uses5mpsc2rfv7zulhcqy0m979", &chaincfg.TestNet3Params)
	fmt.Printf("to address:%v, to sriptAddress:%v\n", to.String(), hex.EncodeToString(to.ScriptAddress()))

	depositTx := wire.NewMsgTx(2)

	hash, err := chainhash.NewHashFromStr("77014d81ac23bdb3f4646c29afe3e8803e291ea9f69e585aee331832c0a62581")
	txIn := wire.NewTxIn(wire.NewOutPoint(hash, 0), nil, nil)
	depositTx.AddTxIn(txIn)

	//TxIn's pkScript and value
	txInPkScript, err := hex.DecodeString("0014f97a2ead90717062357c8c1ee15d3ed0a5324efd")
	txInValue := btcutil.Amount(7000)

	//
	txOutScript, err := txscript.PayToAddrScript(to)
	if err != nil {
		fmt.Println(err)
		return
	}
	txOut := wire.NewTxOut(6500, txOutScript)
	depositTx.AddTxOut(txOut)

	err = CreateTx(depositTx, [][]byte{txInPkScript}, []btcutil.Amount{txInValue}, []*btcec.PrivateKey{privKey}, &chaincfg.TestNet3Params)
	assert.NoError(t, err)

	err = ValidateMsgTx(depositTx, [][]byte{txInPkScript}, []btcutil.Amount{txInValue})
	assert.NoError(t, err)

	var buf bytes.Buffer
	depositTx.Serialize(&buf)
	fmt.Printf("after sign deposit: %v\n", hex.EncodeToString(buf.Bytes()))
}

func TestCreateRedeemTx_1(t *testing.T) {
	expectedHash, _ := hex.DecodeString("f889a5df00f886ba0c932ab668a28ab8b9f60ef8332a628065a61561de515585")
	scretes := []string{
		"23c9cdb2685d0905c0969dbbbfd27fdc1791e16e43b0352d9f11a89053d268ac",
		"47b38c30407286330562e228a73bf84f0c6d5d9593bd16b2dfc66ca1654ab83d",
		"968b40431da7f3aba9dfea20f0c9790ca38117d884ce47ef03d36829cfc48f49",
	}

	privKeys := []*btcec.PrivateKey{}
	pubKeys := []*btcec.PublicKey{}
	addrPubKeys := []*btcutil.AddressPubKey{}
	for _, secret := range scretes {
		s, _ := hex.DecodeString(secret)
		privKey, pubKey := btcec.PrivKeyFromBytes(s)
		privKeys = append(privKeys, privKey)
		pubKeys = append(pubKeys, pubKey)
		addrPubKey, _ := btcutil.NewAddressPubKey(pubKey.SerializeCompressed(), &chaincfg.TestNet3Params)
		addrPubKeys = append(addrPubKeys, addrPubKey)
	}

	multiSigScript, _ := txscript.MultiSigScript(addrPubKeys, 2)

	scriptHash := sha256.Sum256(multiSigScript)
	from, _ := btcutil.NewAddressWitnessScriptHash(scriptHash[:], &chaincfg.RegressionNetParams)

	to, _ := btcutil.DecodeAddress("tb1ql9azatvsw9cxydtu3s0wzhf76zjnynhasuy4zy", &chaincfg.TestNet3Params)

	reedeemTx := wire.NewMsgTx(2)
	hash, err := chainhash.NewHashFromStr("a403d4cdd7d1cbc5839fa7c3f842175013d47c295f50acb8a0d2319e30eb7002")
	txIn := wire.NewTxIn(wire.NewOutPoint(hash, 0), nil, nil)
	reedeemTx.AddTxIn(txIn)

	//TxIn's pkScript and value
	txInPkScript, err := hex.DecodeString("0020f13079e066269d06feecbb74847c29ead0ff1e34e4330a6c30c28696785cfdf0")
	txInValue := btcutil.Amount(9000)

	//TxOut
	txOutScript, err := txscript.PayToAddrScript(to)
	if err != nil {
		fmt.Println(err)
		return
	}
	txOut := wire.NewTxOut(8000, txOutScript)
	reedeemTx.AddTxOut(txOut)

	hashes, err := CalWitnessSigHash(reedeemTx, [][]byte{txInPkScript}, []btcutil.Amount{txInValue}, from, &chaincfg.TestNet3Params, [][]byte{multiSigScript})
	assert.NoError(t, err)
	fmt.Printf("hash:%v\n", hex.EncodeToString(hashes[0]))
	assert.Equal(t, expectedHash, hashes[0])

	var sigs [][]byte
	for _, priv := range privKeys {
		sig := ecdsa.Sign(priv, hashes[0])

		sigWithType := append(sig.Serialize(), byte(txscript.SigHashAll))
		sigs = append(sigs, sigWithType)
	}

	witnessScript, err := MergeMultiSignatures(2, multiSigScript, sigs)
	assert.NoError(t, err)

	reedeemTx.TxIn[0].Witness = witnessScript
	var buf bytes.Buffer
	err = reedeemTx.Serialize(&buf)
	fmt.Printf("signed reedeem: %v\n", hex.EncodeToString(buf.Bytes()))
	//correct: 020000000001010270eb309e31d2a0b8ac505f297cd413501742f8c3a79f83c5cbd1d7cdd403a40000000000ffffffff01401f000000000000160014f97a2ead90717062357c8c1ee15d3ed0a5324efd04004830450221009a2ccd91d89bf37c556863f13ed939aed04694e34dc97e0ea9f1c35018e46d23022055e657a3d93ceb693a4983773d6907ffdc8325798ce977546e1c87f43a67bf5b014730440220765f46fcb6bc52d24ee6fe593661d414c26242aac6ec8c17e7b61c9e1d8fbacc02202c0ddd31048508f2b9ed73d81540055c868b2e78df14f93c137c0bf2baaa39e001695221028fa190883221d93c3ecd3d9a7c7afa130393d56826acc811b3d27834b4986f3221033e8d41a47d121a6a4ac4e05db8967b47ff3036507e7d95a6b912483bea9ab7162103d78e3a9b9b1b966b930e13acf2eb90eb9b9c87c044e6f05a49b6bc0c3d5c5a2b53ae00000000
	err = ValidateMsgTx(reedeemTx, [][]byte{txInPkScript}, []btcutil.Amount{txInValue})
	assert.NoError(t, err)
}

func TestCreateRedeemTx_2(t *testing.T) {
	expectedHash, _ := hex.DecodeString("00a165032df678da63f0e024bff6f6593a3d3af2bb318a992fc5c88fb2bcf613")
	scretes := []string{
		"23c9cdb2685d0905c0969dbbbfd27fdc1791e16e43b0352d9f11a89053d268ac",
		"47b38c30407286330562e228a73bf84f0c6d5d9593bd16b2dfc66ca1654ab83d",
		"968b40431da7f3aba9dfea20f0c9790ca38117d884ce47ef03d36829cfc48f49",
	}

	privKeys := []*btcec.PrivateKey{}
	pubKeys := []*btcec.PublicKey{}
	addrPubKeys := []*btcutil.AddressPubKey{}
	for _, secret := range scretes {
		s, _ := hex.DecodeString(secret)
		privKey, pubKey := btcec.PrivKeyFromBytes(s)
		privKeys = append(privKeys, privKey)
		pubKeys = append(pubKeys, pubKey)
		addrPubKey, _ := btcutil.NewAddressPubKey(pubKey.SerializeCompressed(), &chaincfg.TestNet3Params)
		addrPubKeys = append(addrPubKeys, addrPubKey)
	}

	multiSigScript, _ := txscript.MultiSigScript(addrPubKeys, 2)

	scriptHash := sha256.Sum256(multiSigScript)
	from, _ := btcutil.NewAddressWitnessScriptHash(scriptHash[:], &chaincfg.RegressionNetParams)

	to, _ := btcutil.DecodeAddress("tb1ql9azatvsw9cxydtu3s0wzhf76zjnynhasuy4zy", &chaincfg.TestNet3Params)

	reedeemTx := wire.NewMsgTx(2)
	hash, err := chainhash.NewHashFromStr("a0b391b03d17c3a07a65652b5807931bcbb31d63894b8fd46538fc50602948c3")
	txIn := wire.NewTxIn(wire.NewOutPoint(hash, 0), nil, nil)
	reedeemTx.AddTxIn(txIn)

	//TxIn's pkScript and value
	txInPkScript, err := hex.DecodeString("0020f13079e066269d06feecbb74847c29ead0ff1e34e4330a6c30c28696785cfdf0")
	txInValue := btcutil.Amount(6500)

	//TxOut
	txOutScript, err := txscript.PayToAddrScript(to)
	if err != nil {
		fmt.Println(err)
		return
	}
	txOut := wire.NewTxOut(6000, txOutScript)
	reedeemTx.AddTxOut(txOut)

	hashes, err := CalWitnessSigHash(reedeemTx, [][]byte{txInPkScript}, []btcutil.Amount{txInValue}, from, &chaincfg.TestNet3Params, [][]byte{multiSigScript})
	assert.NoError(t, err)
	fmt.Printf("hash:%v\n", hex.EncodeToString(hashes[0]))
	assert.Equal(t, expectedHash, hashes[0])

	var sigs [][]byte
	for _, priv := range privKeys {
		sig := ecdsa.Sign(priv, hashes[0])

		sigWithType := append(sig.Serialize(), byte(txscript.SigHashAll))
		sigs = append(sigs, sigWithType)
	}

	witnessScript, err := MergeMultiSignatures(2, multiSigScript, sigs)
	assert.NoError(t, err)

	reedeemTx.TxIn[0].Witness = witnessScript
	var buf bytes.Buffer
	err = reedeemTx.Serialize(&buf)
	//fmt.Printf("signed reedeem TxHash: %v\n", reedeemTx.TxHash())
	//fmt.Printf("signed reedeem WitnessHash: %v\n", reedeemTx.WitnessHash())
	fmt.Printf("signed reedeem: %v\n", hex.EncodeToString(buf.Bytes()))
	//correct: 020000000001010270eb309e31d2a0b8ac505f297cd413501742f8c3a79f83c5cbd1d7cdd403a40000000000ffffffff01401f000000000000160014f97a2ead90717062357c8c1ee15d3ed0a5324efd04004830450221009a2ccd91d89bf37c556863f13ed939aed04694e34dc97e0ea9f1c35018e46d23022055e657a3d93ceb693a4983773d6907ffdc8325798ce977546e1c87f43a67bf5b014730440220765f46fcb6bc52d24ee6fe593661d414c26242aac6ec8c17e7b61c9e1d8fbacc02202c0ddd31048508f2b9ed73d81540055c868b2e78df14f93c137c0bf2baaa39e001695221028fa190883221d93c3ecd3d9a7c7afa130393d56826acc811b3d27834b4986f3221033e8d41a47d121a6a4ac4e05db8967b47ff3036507e7d95a6b912483bea9ab7162103d78e3a9b9b1b966b930e13acf2eb90eb9b9c87c044e6f05a49b6bc0c3d5c5a2b53ae00000000
	err = ValidateMsgTx(reedeemTx, [][]byte{txInPkScript}, []btcutil.Amount{txInValue})
	assert.NoError(t, err)
}

func TestCalculateHash_1(t *testing.T) {
	reedeemTx := wire.NewMsgTx(2)
	expectedHash, _ := hex.DecodeString("f889a5df00f886ba0c932ab668a28ab8b9f60ef8332a628065a61561de515585")

	from, _ := btcutil.DecodeAddress("tb1q7yc8ncrxy6wsdlhvhd6gglpfatg07835uses5mpsc2rfv7zulhcqy0m979", &chaincfg.TestNet3Params)
	hash, err := chainhash.NewHashFromStr("a403d4cdd7d1cbc5839fa7c3f842175013d47c295f50acb8a0d2319e30eb7002")
	assert.NoError(t, err)
	txIn := wire.NewTxIn(wire.NewOutPoint(hash, 0), nil, nil)
	reedeemTx.AddTxIn(txIn)

	//TxIn's pkScript and value. get from
	txInOriginalScript, err := hex.DecodeString("5221028fa190883221d93c3ecd3d9a7c7afa130393d56826acc811b3d27834b4986f3221033e8d41a47d121a6a4ac4e05db8967b47ff3036507e7d95a6b912483bea9ab7162103d78e3a9b9b1b966b930e13acf2eb90eb9b9c87c044e6f05a49b6bc0c3d5c5a2b53ae")
	txInPkScript, err := hex.DecodeString("0020f13079e066269d06feecbb74847c29ead0ff1e34e4330a6c30c28696785cfdf0")
	txInValue := btcutil.Amount(9000)

	to, _ := btcutil.DecodeAddress("tb1ql9azatvsw9cxydtu3s0wzhf76zjnynhasuy4zy", &chaincfg.TestNet3Params)
	txOutScript, err := txscript.PayToAddrScript(to)
	if err != nil {
		fmt.Println(err)
		return
	}
	txOut := wire.NewTxOut(8000, txOutScript)
	reedeemTx.AddTxOut(txOut)

	hashes, err := CalWitnessSigHash(reedeemTx, [][]byte{txInPkScript}, []btcutil.Amount{txInValue}, from, &chaincfg.TestNet3Params, [][]byte{txInOriginalScript})
	assert.Equal(t, expectedHash, hashes[0])
}

func TestCalculateHash_2(t *testing.T) {
	expectedHash, err := hex.DecodeString("f889a5df00f886ba0c932ab668a28ab8b9f60ef8332a628065a61561de515585")
	assert.NoError(t, err)
	raw, _ := hex.DecodeString("02000000081aba67bc0e257c733f0c00dab79853b37e19f617e9a63b271e9e617a16377a3bb13029ce7b1f559ef5e747fcac439f1455a2ec7c5f09b72290795e706650440270eb309e31d2a0b8ac505f297cd413501742f8c3a79f83c5cbd1d7cdd403a400000000695221028fa190883221d93c3ecd3d9a7c7afa130393d56826acc811b3d27834b4986f3221033e8d41a47d121a6a4ac4e05db8967b47ff3036507e7d95a6b912483bea9ab7162103d78e3a9b9b1b966b930e13acf2eb90eb9b9c87c044e6f05a49b6bc0c3d5c5a2b53ae2823000000000000ffffffffedfe2b7591d80ff513a432a6fdb4ae8a410469dd3a700e3749bb05d25251207a0000000001000000")
	hash := chainhash.DoubleHashB(raw)
	assert.Equal(t, expectedHash, hash)
}
