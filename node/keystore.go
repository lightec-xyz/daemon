package node

import (
	"encoding/hex"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
)

// todo

const (
	SecretKeyId = "secretKey"
)

type KeyStore struct {
	memguard *Memguard
	address  string
}

func NewKeyStore(privateKey string) (*KeyStore, error) {
	memguard := NewMemguard()
	address, err := privateKeyToEthAddr(privateKey)
	if err != nil {
		logger.Error("privateKeyToEthAddr error:%v", err)
		return nil, err
	}
	logger.Debug("keystore address: %v", address)
	hexSecret, err := hex.DecodeString(privateKey)
	if err != nil {
		logger.Error("decode private key error:%v", err)
		return nil, err
	}
	memguard.Store(SecretKeyId, hexSecret)
	return &KeyStore{
		memguard: memguard,
		address:  address,
	}, nil
}

func (k *KeyStore) EthAddress() string {
	return k.address
}

func (k *KeyStore) GetPrivateKey() ([]byte, error) {
	bytes, err := k.memguard.Load(SecretKeyId)
	if err != nil {
		logger.Error("get private key error:%v", err)
		return nil, err
	}
	return bytes, nil
}

func (k *KeyStore) VerifyJwt(token string) (*rpc.CustomClaims, error) {
	secret, err := k.GetPrivateKey()
	if err != nil {
		logger.Error("get private key error:%v", err)
		return nil, err
	}
	jwt, err := rpc.VerifyJWT(secret, token)
	if err != nil {
		logger.Error("verify jwt error:%v", err)
		return nil, err
	}
	return jwt, nil
}
