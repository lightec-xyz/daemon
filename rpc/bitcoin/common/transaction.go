package bitcoin

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

const TxVersion = 2

type MultiTransactionBuilder struct {
	MsgTx            *wire.MsgTx
	txInPkScripts    [][]byte
	txInValues       []btcutil.Amount
	multiScriptsList [][]byte
	multiSigScript   []byte
	netParams        *chaincfg.Params
	nRequired        int
	nTotal           int
}

func NewMultiTransactionBuilder() *MultiTransactionBuilder {
	return NewMultiTxBuilder(&chaincfg.MainNetParams)
}

func NewMultiTxBuilder(netParams *chaincfg.Params) *MultiTransactionBuilder {
	return &MultiTransactionBuilder{
		MsgTx:            wire.NewMsgTx(TxVersion),
		txInPkScripts:    make([][]byte, 0),
		txInValues:       make([]btcutil.Amount, 0),
		multiScriptsList: make([][]byte, 0),
		multiSigScript:   make([]byte, 0),
		netParams:        netParams,
	}
}

func (mb *MultiTransactionBuilder) NetParams(network NetWork) error {
	params, err := getNetworkParams(network)
	if err != nil {
		return err
	}
	mb.netParams = params
	return nil
}

func (mb *MultiTransactionBuilder) AddMultiScript(multiSigScript []byte, nRequired, nTotal int) error {
	if nRequired > nTotal {
		return errors.New("the number of required signatures is too much")
	}

	mb.multiSigScript = multiSigScript
	mb.nRequired = nRequired
	mb.nTotal = nTotal

	return nil
}

func (mb *MultiTransactionBuilder) AddMultiPublicKey(publicList [][]byte, nRequired int) error {
	var publicKeys []*btcutil.AddressPubKey
	for _, pub := range publicList {
		addrPubKey, err := btcutil.NewAddressPubKey(pub, mb.netParams)
		if err != nil {
			return err
		}
		publicKeys = append(publicKeys, addrPubKey)
	}
	multiSigScript, err := txscript.MultiSigScript(publicKeys, nRequired)
	if err != nil {
		return err
	}
	return mb.AddMultiScript(multiSigScript, nRequired, len(publicKeys))
}

func (mb *MultiTransactionBuilder) AddTxIn(inputs []TxIn) error {
	for _, input := range inputs {
		hash, err := chainhash.NewHashFromStr(input.Hash)
		if err != nil {
			return err
		}
		txIn := wire.NewTxIn(wire.NewOutPoint(hash, input.VOut), nil, nil)
		mb.MsgTx.AddTxIn(txIn)
		pkScriptBytes, err := hex.DecodeString(input.PkScript)
		if err != nil {
			return err
		}
		mb.txInPkScripts = append(mb.txInPkScripts, pkScriptBytes)
		mb.txInValues = append(mb.txInValues, btcutil.Amount(input.Amount))
		mb.multiScriptsList = append(mb.multiScriptsList, mb.multiSigScript)
	}
	return nil
}

func (mb *MultiTransactionBuilder) AddTxOutScript(outputs []TxOut) error {
	for _, output := range outputs {
		txOut := wire.NewTxOut(output.Amount, output.PayScript)
		mb.MsgTx.AddTxOut(txOut)
	}
	return nil
}

func (mb *MultiTransactionBuilder) AddTxOut(outputs []TxOut) error {
	for _, output := range outputs {
		address, err := btcutil.DecodeAddress(output.Address, mb.netParams)
		if err != nil {
			return err
		}
		pkScript, err := txscript.PayToAddrScript(address)
		if err != nil {
			return err
		}
		txOut := wire.NewTxOut(output.Amount, pkScript)
		mb.MsgTx.AddTxOut(txOut)
	}
	return nil
}

func (mb *MultiTransactionBuilder) TxHash() string {
	return mb.MsgTx.TxHash().String()
}

func (mb *MultiTransactionBuilder) SigHashes() ([]string, error) {
	scriptHash := sha256.Sum256(mb.multiSigScript)
	from, err := btcutil.NewAddressWitnessScriptHash(scriptHash[:], mb.netParams)
	if err != nil {
		return nil, err
	}
	hashes, err := CalWitnessSigHash(mb.MsgTx, mb.txInPkScripts, mb.txInValues, from, mb.netParams, mb.multiScriptsList)
	if err != nil {
		return nil, err
	}
	var hashList []string
	for _, hash := range hashes {
		hashList = append(hashList, hex.EncodeToString(hash))
		//mb.MsgTx.TxIn[index].Witness = witnessScript
	}
	return hashList, nil
}

func (mb *MultiTransactionBuilder) Sign(signFn func(hash []byte) ([][]byte, error)) error {
	scriptHash := sha256.Sum256(mb.multiSigScript)
	from, err := btcutil.NewAddressWitnessScriptHash(scriptHash[:], mb.netParams)
	if err != nil {
		return err
	}
	hashes, err := CalWitnessSigHash(mb.MsgTx, mb.txInPkScripts, mb.txInValues, from, mb.netParams, mb.multiScriptsList)
	if err != nil {
		return err
	}
	for index, hash := range hashes {
		sigs, err := signFn(hash)
		if err != nil {
			return err
		}
		witnessScript, err := MergeMultiSignatures(mb.nRequired, mb.multiSigScript, sigs)
		if err != nil {
			return err
		}
		mb.MsgTx.TxIn[index].Witness = witnessScript
	}
	return nil
}

func (mb *MultiTransactionBuilder) MergeSignature(signatures [][][]byte) error {
	if len(signatures) < mb.nRequired || len(signatures) > mb.nTotal {
		return fmt.Errorf("the number of signatures is too much or too little")
	}
	nTxin := len(mb.MsgTx.TxIn)
	nNil := mb.nTotal - mb.nRequired
	for i := 0; i < nTxin; i++ {
		witnessElements := make(wire.TxWitness, 0, mb.nTotal+1)
		for j := 0; j < nNil; j++ {
			witnessElements = append(witnessElements, nil)
		}
		for j := 0; j < mb.nRequired; j++ {
			witnessElements = append(witnessElements, signatures[j][i])
		}
		witnessElements = append(witnessElements, mb.multiSigScript)
		mb.MsgTx.TxIn[i].Witness = witnessElements
	}
	return nil
}

func (mb *MultiTransactionBuilder) Deserialize(txData []byte) error {
	return mb.MsgTx.Deserialize(bytes.NewReader(txData))
}

func (mb *MultiTransactionBuilder) Serialize() ([]byte, error) {
	var buf bytes.Buffer
	err := mb.MsgTx.Serialize(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (mb *MultiTransactionBuilder) Build() ([]byte, error) {
	err := validateMsgTx(mb.MsgTx, mb.txInPkScripts, mb.txInValues)
	if err != nil {
		return nil, err
	}
	return mb.Serialize()
}

type TxIn struct {
	Hash     string
	VOut     uint32
	PkScript string
	Amount   int64
}

type TxOut struct {
	Address   string
	PayScript []byte
	Amount    int64
}

type Transaction struct {
	*wire.MsgTx
}

func NewTransaction() *Transaction {
	return &Transaction{
		MsgTx: wire.NewMsgTx(2),
	}
}

func CreateDepositTransaction(secret, ethAddr []byte, inputs []TxIn, outputs []TxOut, network NetWork) ([]byte, error) {
	networkParams, err := getNetworkParams(network)
	if err != nil {
		return nil, err
	}
	privateKey, _ := btcec.PrivKeyFromBytes(secret)
	var txInPkScripts [][]byte
	var txInValues []btcutil.Amount
	msgTx := wire.NewMsgTx(wire.TxVersion)
	for _, input := range inputs {
		hash, err := chainhash.NewHashFromStr(input.Hash)
		if err != nil {
			return nil, err
		}
		txIn := wire.NewTxIn(wire.NewOutPoint(hash, input.VOut), nil, nil)
		msgTx.AddTxIn(txIn)
		pkScriptBytes, err := hex.DecodeString(input.PkScript)
		if err != nil {
			return nil, err
		}
		txInPkScripts = append(txInPkScripts, pkScriptBytes)
		txInValues = append(txInValues, btcutil.Amount(input.Amount))
	}
	//opReturnScript, err := txscript.NewScriptBuilder().AddOp(txscript.OP_RETURN).AddData(ethAddr).Script()
	//if err != nil {
	//	return nil, err
	//}
	//msgTx.AddTxOut(wire.NewTxOut(0, opReturnScript))
	for _, output := range outputs {
		address, err := btcutil.DecodeAddress(output.Address, networkParams)
		if err != nil {
			return nil, err
		}
		txOutScript, err := txscript.PayToAddrScript(address)
		if err != nil {
			return nil, err
		}
		txOut := wire.NewTxOut(output.Amount, txOutScript)
		msgTx.AddTxOut(txOut)
	}

	err = createTx(msgTx, txInPkScripts, txInValues, []*btcec.PrivateKey{privateKey}, networkParams)
	if err != nil {
		return nil, err
	}
	err = validateMsgTx(msgTx, txInPkScripts, txInValues)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = msgTx.Serialize(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func txPrevOutFetcher(tx *wire.MsgTx, prevPkScripts [][]byte,
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

func validateMsgTx(tx *wire.MsgTx, prevScripts [][]byte,
	inputValues []btcutil.Amount) error {

	inputFetcher, err := txPrevOutFetcher(
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

func createTx(tx *wire.MsgTx, prevPkScripts [][]byte,
	inputValues []btcutil.Amount, privKeys []*btcec.PrivateKey, chainParams *chaincfg.Params) error {

	inputFetcher, err := txPrevOutFetcher(tx, prevPkScripts, inputValues)
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

// calcWitnessSignatureHashRaw 从btcd 拷贝而来，增加了对p2wsh的支持
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

// calcWitnessSignatureHashRaw 从btcd 拷贝而来，增加了对p2wsh的支持
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

	inputFetcher, err := txPrevOutFetcher(tx, txInLockingScripts, inputValues)
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
