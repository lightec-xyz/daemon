package node

import (
	"context"
	"encoding/hex"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc/bitcoin"
	btctx "github.com/lightec-xyz/daemon/rpc/bitcoin/common"
	ethrpc "github.com/lightec-xyz/daemon/rpc/ethereum"
	"github.com/lightec-xyz/daemon/rpc/oasis"
	"github.com/lightec-xyz/daemon/store"
	"math/big"
	"sync"
	"time"
)

/*
todo
*/

type TxManager struct {
	ethClient   *ethrpc.Client
	btcClient   *bitcoin.Client
	oasisClient *oasis.Client
	store       store.IStore
	keyStore    *KeyStore
	lock        sync.Mutex
}

func NewTxManager(store store.IStore, keyStore *KeyStore, ethClient *ethrpc.Client, btcClient *bitcoin.Client, oasisClient *oasis.Client) (*TxManager, error) {
	return &TxManager{
		ethClient:   ethClient,
		btcClient:   btcClient,
		oasisClient: oasisClient,
		keyStore:    keyStore,
		store:       store,
	}, nil
}

func (t *TxManager) init() error {
	//allUnSubmitTxs, err := ReadAllUnSubmitTxs(t.store)
	//if err != nil {
	//	logger.Error("get unsubmit tx error:%v", err)
	//	return err
	//}
	return nil
}

func (t *TxManager) AddTask(resp *common.ZkProofResponse) {
	logger.Info("txManager add retry task: %v", resp.RespId())
	unSubmitTx := NewDbUnSubmitTx(resp.Hash, hex.EncodeToString(resp.Proof), resp.ProofType)
	err := WriteUnSubmitTx(t.store, []DbUnSubmitTx{unSubmitTx})
	if err != nil {
		logger.Error("write unsubmit tx error: %v %v", resp.RespId(), err)
		return
	}
}

func (t *TxManager) Check() error {
	unSubmitTxs, err := ReadAllUnSubmitTxs(t.store)
	if err != nil {
		logger.Error("read unsubmit tx error:%v", err)
		return err
	}
	for _, tx := range unSubmitTxs {
		switch tx.ProofType {
		case common.VerifyTxType:
			hash, err := t.UpdateUtxoChange(tx.Hash, tx.Proof)
			if err != nil {
				logger.Error("update utxo error: %v %v", tx.ProofType.String(), tx.Hash)
				return err
			}
			logger.Info("success update utxo txId: %v,hash: %v", tx.Hash, hash)
		case common.RedeemTxType:
			txHash, err := t.RedeemZkbtc(tx.Hash, tx.Proof)
			if err != nil {
				logger.Error("update utxo error: %v %v", tx.ProofType.String(), tx.Hash)
				continue
			}
			logger.Debug("success redeem btc ethHash: %v,btcHash: %v", tx.Hash, txHash)
		default:
			logger.Warn("never should happen: %v %v", tx.ProofType.String(), tx.Hash)
		}
	}
	return nil
}

func (t *TxManager) getEthAddrNonce(addr string) (uint64, error) {
	chainNonce, err := t.ethClient.GetNonce(addr)
	if err != nil {
		logger.Error("get nonce error: %v %v", addr, err)
		return 0, err
	}
	dbNonce, exists, err := ReadNonce(t.store, common.ETH.String(), addr)
	if err != nil {
		logger.Error("read nonce error: %v %v", addr, err)
		return 0, err
	}
	if !exists {
		return chainNonce, nil
	}
	if chainNonce <= dbNonce {
		return dbNonce + 1, nil
	}
	return chainNonce, nil
}

func (t *TxManager) RedeemZkbtc(hash, proof string) (string, error) {
	if common.GetEnvDebugMode() {
		logger.Debug("current is debug,skip verify and sign: %v", hash)
		return "", nil
	}
	proofBytes, err := hex.DecodeString(proof)
	if err != nil {
		logger.Error("decode proof error: %v %v", hash, err)
		return "", err
	}
	ethTxHash := ethcommon.HexToHash(hash)
	ethTx, _, err := t.ethClient.TransactionByHash(context.Background(), ethTxHash)
	if err != nil {
		logger.Error("get eth tx error:%v", err)
		return "", err
	}
	receipt, err := t.ethClient.TransactionReceipt(context.Background(), ethTxHash)
	if err != nil {
		logger.Error("get eth tx receipt error:%v", err)
		return "", err
	}

	btcRawTx, _, err := ethrpc.DecodeRedeemLog(receipt.Logs[3].Data)
	if err != nil {
		logger.Error("decode redeem log error:%v", err)
		return "", err
	}
	logger.Info("btcRawTx: %v\n", hexutil.Encode(btcRawTx))
	rawTx, rawReceipt := ethrpc.GetRawTxAndReceipt(ethTx, receipt)
	logger.Info("rawTx: %v\n", hexutil.Encode(rawTx))
	logger.Info("rawReceipt: %v\n", hexutil.Encode(rawReceipt))

	sigs, err := t.oasisClient.SignBtcTx(rawTx, rawReceipt, proofBytes)
	if err != nil {
		logger.Error("sign btc tx error:%v", err)
		return "", nil
	}
	transaction := btctx.NewMultiTransactionBuilder()
	err = transaction.Deserialize(btcRawTx)
	if err != nil {
		logger.Error("deserialize btc tx error:%v", err)
		return "", err
	}
	multiSigScript, err := t.ethClient.GetMultiSigScript()
	if err != nil {
		logger.Error("get multi sig script error:%v", err)
		return "", err
	}
	nTotal, nRequred := 3, 2
	transaction.AddMultiScript(multiSigScript, nRequred, nTotal)
	err = transaction.MergeSignature(sigs[:nRequred])
	if err != nil {
		logger.Error("merge signature error:%v", err)
		return "", err
	}
	btxTx, err := transaction.Serialize()
	if err != nil {
		logger.Error("serialize btc tx error:%v", err)
		return "", err
	}
	btcTxHash := transaction.TxHash()
	_, err = t.btcClient.GetTransaction(btcTxHash) // todo
	if err == nil {
		logger.Warn("btc tx already exist: %v", btcTxHash)
		return "", nil
	}
	txHex := hex.EncodeToString(btxTx)
	logger.Info("btc Tx: %v\n", txHex)
	txHash, err := t.btcClient.Sendrawtransaction(txHex)
	if err != nil {
		logger.Error("send btc tx error:%v %v", btcTxHash, err)
		// todo  just test
		_, err = bitcoin.BroadcastTx(txHex)
		if err != nil {
			logger.Error("broadcast btc tx error %v:%v", btcTxHash, err)
			return "", err
		}
	}
	logger.Info("send redeem btc tx: %v", btcTxHash)
	err = DeleteUnSubmitTx(t.store, hash)
	if err != nil {
		logger.Error("delete unSubmit tx error: %v %v", hash, err)
		return "", err
	}
	return txHash, err
}

func (t *TxManager) UpdateUtxoChange(txId, proof string) (string, error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	proofBytes, err := hex.DecodeString(proof)
	if err != nil {
		logger.Error("decode proof error: %v %v", txId, err)
		return "", err
	}
	privateKey, err := t.keyStore.GetPrivateKey()
	if err != nil {
		logger.Error("get private key error: %v", err)
		return "", err
	}
	ethAddress := t.keyStore.EthAddress()
	nonce, err := t.getEthAddrNonce(ethAddress)
	if err != nil {
		logger.Error("get nonce error: %v %v", ethAddress, err)
		return "", err
	}
	chainId, err := t.ethClient.GetChainId()
	if err != nil {
		logger.Error("get chain id error:%v", err)
		return "", err
	}
	gasPrice, err := t.ethClient.GetGasPrice()
	if err != nil {
		logger.Error("get gas price error:%v", err)
		return "", err
	}
	gasLimit := uint64(500000)
	gasPrice = big.NewInt(0).Mul(gasPrice, big.NewInt(2)) // todo
	txHash, err := t.ethClient.UpdateUtxoChange(privateKey, []string{txId}, nonce, gasLimit, chainId, gasPrice, proofBytes)
	if err != nil {
		logger.Error("update utxo change error:%v", err)
		return "", err
	}
	logger.Debug("update utxo info address: %v txHash: %v,txId: %v,nonce: %v", ethAddress, txHash, txId, nonce)
	err = WriteNonce(t.store, common.ETH.String(), ethAddress, nonce)
	if err != nil {
		logger.Error("write nonce error: %v %v", ethAddress, err)
		return "", err
	}
	err = DeleteUnSubmitTx(t.store, txId)
	if err != nil {
		logger.Error("delete unsubmit tx error: %v %v", txId, err)
		return "", err
	}
	return txHash, nil
}

func (t *TxManager) Close() error {
	return nil
}

func NewDbUnSubmitTx(hash, proof string, proofType common.ZkProofType) DbUnSubmitTx {
	return DbUnSubmitTx{
		Hash:      hash,
		Proof:     proof,
		ProofType: proofType,
		Timestamp: time.Now().UnixNano(),
	}
}
