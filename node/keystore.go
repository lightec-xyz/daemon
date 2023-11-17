package node

import (
	"encoding/hex"
	"fmt"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"reflect"
	"sync"
)

const (
	SecretKeyId = "secretKey"
)

type KeyStore struct {
	memguard  *Memguard
	address   string
	adminPath *sync.Map
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
	adminPath := new(sync.Map)
	typeValue := reflect.TypeOf((*rpc.IAdmin)(nil)).Elem()
	for i := 0; i < typeValue.NumMethod(); i++ {
		path := fmt.Sprintf("%s_%s", RpcRegisterName, LowercaseFirstLetter(typeValue.Method(i).Name))
		logger.Debug("admin path: %v", path)
		adminPath.Store(path, rpc.JwtPermission)
	}
	memguard.Store(SecretKeyId, hexSecret)
	return &KeyStore{
		memguard:  memguard,
		address:   address,
		adminPath: adminPath,
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

func (k *KeyStore) CheckPermission(method string) (rpc.Permission, error) {
	if value, ok := k.adminPath.Load(method); ok {
		perm, ok := value.(rpc.Permission)
		if ok {
			return perm, nil
		}
		return rpc.NonePermission, fmt.Errorf("invalid permission value ")
	}
	return rpc.NonePermission, nil
}
