package node

// todo

type KeyStore struct {
	privateKey string
}

func NewKeyStore(privateKey string) *KeyStore {
	return &KeyStore{
		privateKey: privateKey,
	}
}

func (k *KeyStore) GetPrivateKey() string {
	return k.privateKey
}
