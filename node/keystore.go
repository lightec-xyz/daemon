package node

import "github.com/lightec-xyz/daemon/logger"

// todo

type KeyStore struct {
	privateKey string
}

func NewKeyStore(privateKey string) *KeyStore {
	return &KeyStore{
		privateKey: privateKey,
	}
}

func (k *KeyStore) Address() (string, error) {
	address, err := privateKeyToEthAddr(k.privateKey)
	if err != nil {
		logger.Error("privateKeyToEthAddr error:%v", err)
		return "", err
	}
	return address, nil
}

func (k *KeyStore) GetPrivateKey() string {
	return k.privateKey
}
