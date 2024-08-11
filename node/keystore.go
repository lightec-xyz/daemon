package node

import (
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
	memguard.Store(SecretKeyId, []byte(privateKey))
	return &KeyStore{
		memguard: memguard,
		address:  address,
	}, nil
}

func (k *KeyStore) EthAddress() string {
	return k.address
}

func (k *KeyStore) GetPrivateKey() (string, error) {
	bytes, err := k.memguard.Load(SecretKeyId)
	if err != nil {
		logger.Error("get private key error:%v", err)
		return "", err
	}
	return string(bytes), nil
}

func (k *KeyStore) VerifyJwt(token string) (*rpc.CustomClaims, error) {
	secret, err := k.GetPrivateKey()
	if err != nil {
		logger.Error("get private key error:%v", err)
		return nil, err
	}
	jwt, err := rpc.VerifyJWT([]byte(secret), token)
	if err != nil {
		logger.Error("verify jwt error:%v", err)
		return nil, err
	}
	return jwt, nil
}
