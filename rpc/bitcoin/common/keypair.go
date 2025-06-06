package bitcoin

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type IKeyPair interface {
	Sign(msg []byte) []byte
	PublicKey() PublicKey
	PrivateKey() PrivateKey
	Address(addrType AddrType, network NetWork) (string, string, error)
	Verify(message, signature []byte) (bool, error)
}

type PrivateKey []byte
type PublicKey []byte
type AddrType string
type NetWork string

type KeyPair struct {
	privateKey *btcec.PrivateKey
	publicKey  *secp256k1.PublicKey
}

func (k *KeyPair) PrivateKey() PrivateKey {
	return k.privateKey.Serialize()
}

func (k *KeyPair) PublicKey() PublicKey {
	return k.publicKey.SerializeCompressed()
}

func (k *KeyPair) Address(addrType AddrType, network NetWork) (string, string, error) {
	pubKey := k.publicKey
	netParams, err := getNetworkParams(network)
	if err != nil {
		return "", "", err
	}
	switch addrType {
	case P2PKH:
		pkhAddr, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(pubKey.SerializeCompressed()), netParams)
		if err != nil {
			return "", "", err
		}
		lockingScript, err := txscript.PayToAddrScript(pkhAddr)
		if err != nil {
			return "", "", err
		}
		return pkhAddr.EncodeAddress(), fmt.Sprintf("%x", lockingScript), nil

	case P2WPKH:
		wpkhAddr, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(pubKey.SerializeCompressed()), netParams)
		if err != nil {
			return "", "", err
		}
		lockingScript, err := txscript.PayToAddrScript(wpkhAddr)
		if err != nil {
			return "", "", err
		}
		return wpkhAddr.EncodeAddress(), fmt.Sprintf("%x", lockingScript), nil

	default:
		return "", "", fmt.Errorf("unSupport addrType:%v", addrType)
	}

}

func (k *KeyPair) Sign(msg []byte) []byte {
	signature := ecdsa.Sign(k.privateKey, msg)
	return signature.Serialize()
}

func (k *KeyPair) Verify(message, signature []byte) (bool, error) {
	sig, err := ecdsa.ParseSignature(signature)
	if err != nil {
		return false, err
	}
	return sig.Verify(message, k.publicKey), nil
}

func NewKeyPairFromSecret(seed string) (IKeyPair, error) {
	secret, err := hex.DecodeString(seed)
	if err != nil {
		return nil, err
	}
	privateKey, pubKey := btcec.PrivKeyFromBytes(secret)
	keyPari := &KeyPair{
		privateKey: privateKey,
		publicKey:  pubKey,
	}
	return keyPari, nil
}

func NewRandSeed() (IKeyPair, error) {
	privateKey, err := btcec.NewPrivateKey()
	if err != nil {
		return nil, err
	}
	publicKey := privateKey.PubKey()
	keyPari := &KeyPair{
		privateKey: privateKey,
		publicKey:  publicKey,
	}
	return keyPari, nil
}

func MultiScriptAddress(required int, network NetWork, publicKeyList [][]byte) (string, error) {
	networkParams, err := getNetworkParams(network)
	if err != nil {
		return "", err
	}
	var addrPubKeys []*btcutil.AddressPubKey
	for _, pubKey := range publicKeyList {
		addressPubKey, err := btcutil.NewAddressPubKey(pubKey, networkParams)
		if err != nil {
			return "", err
		}
		addrPubKeys = append(addrPubKeys, addressPubKey)
	}
	multiSigScript, err := txscript.MultiSigScript(addrPubKeys, required)
	if err != nil {
		return "", err
	}
	scriptHash := sha256.Sum256(multiSigScript)
	wshAddr, err := btcutil.NewAddressWitnessScriptHash(scriptHash[:], networkParams)
	if err != nil {
		return "", err
	}
	return wshAddr.EncodeAddress(), nil
}

func GetMultiSigScriptRelateds(required int, net *chaincfg.Params, pubBytesList [][]byte) (
	multiSigScript []byte, addr btcutil.Address, lockScript []byte, err error) {

	var addrPubKeys []*btcutil.AddressPubKey
	for _, pubKey := range pubBytesList {
		addressPubKey, err := btcutil.NewAddressPubKey(pubKey, net)
		if err != nil {
			return nil, nil, nil, err
		}
		addrPubKeys = append(addrPubKeys, addressPubKey)
	}

	multiSigScript, err = txscript.MultiSigScript(addrPubKeys, required)
	if err != nil {
		return nil, nil, nil, err
	}

	scriptHash := sha256.Sum256(multiSigScript)
	addr, err = btcutil.NewAddressWitnessScriptHash(scriptHash[:], net)
	if err != nil {
		return nil, nil, nil, err
	}

	lockScript, err = txscript.PayToAddrScript(addr)
	if err != nil {
		return nil, nil, nil, err
	}

	return multiSigScript, addr, lockScript, nil
}

func getNetworkParams(network NetWork) (*chaincfg.Params, error) {
	var netParams *chaincfg.Params
	switch network {
	case MainNet:
		netParams = &chaincfg.MainNetParams
	case TestNet:
		netParams = &chaincfg.TestNet3Params
	case RegTest:
		netParams = &chaincfg.RegressionNetParams
	default:
		return nil, fmt.Errorf("unKnown network:%v", network)
	}
	return netParams, nil
}

func GenPayToAddrScript(address string, network NetWork) (string, error) {
	networkParams, err := getNetworkParams(network)
	if err != nil {
		return "", err
	}
	addr, err := btcutil.DecodeAddress(address, networkParams)
	if err != nil {
		return "", err
	}
	txOutScript, err := txscript.PayToAddrScript(addr)
	if err != nil {
		return "", err
	}
	hexLockScript := hexutil.Encode(txOutScript)
	return hexLockScript, nil
}
