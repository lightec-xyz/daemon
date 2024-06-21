package node

import "github.com/lightec-xyz/daemon/logger"

// todo

const (
	SecretKey = "secretKey"
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
	memguard.Store(SecretKey, []byte(privateKey))
	return &KeyStore{
		memguard: memguard,
		address:  address,
	}, nil
}

func (k *KeyStore) Address() (string, error) {
	return k.address, nil
}

func (k *KeyStore) GetPrivateKey() string {
	bytes, err := k.memguard.Load(SecretKey)
	if err != nil {
		logger.Error("get private key error:%v", err)
		return ""
	}
	return string(bytes)
}
