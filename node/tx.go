package node

import (
	"context"
	"encoding/hex"
	"fmt"
	blockdepthUtil "github.com/lightec-xyz/btc_provers/utils/blockdepth"
	btcproverClient "github.com/lightec-xyz/btc_provers/utils/client"
	"github.com/lightec-xyz/daemon/rpc/ethereum/zkbridge"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/lightec-xyz/daemon/circuits"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc/bitcoin"
	btctx "github.com/lightec-xyz/daemon/rpc/bitcoin/common"
	"github.com/lightec-xyz/daemon/rpc/dfinity"
	ethrpc "github.com/lightec-xyz/daemon/rpc/ethereum"
	"github.com/lightec-xyz/daemon/rpc/oasis"
	"github.com/lightec-xyz/daemon/rpc/sgx"
	"github.com/lightec-xyz/daemon/store"
	redeemUtils "github.com/lightec-xyz/provers/utils/redeem-tx"
)

//todo record tx hash to check status

type TxManager struct {
	ethClient    *ethrpc.Client
	btcClient    *bitcoin.Client
	oasisClient  *oasis.Client
	icpClient    *dfinity.Client
	sgxClient    *sgx.Client
	proverClient btcproverClient.IClient
	minerAddr    string
	submitAddr   string
	chainStore   *ChainStore
	fileStore    *FileStorage
	prepared     *Prepared
	keyStore     *KeyStore
	lock         sync.Mutex
	icpSigMap    map[string][][]byte
}

func NewTxManager(store store.IStore, fileStore *FileStorage, prepared *Prepared, keyStore *KeyStore, ethClient *ethrpc.Client, btcClient *bitcoin.Client,
	oasisClient *oasis.Client, dfinityClient *dfinity.Client, sgxClient *sgx.Client, proverClient btcproverClient.IClient, minerAddr string) (*TxManager, error) {
	return &TxManager{
		ethClient:    ethClient,
		btcClient:    btcClient,
		oasisClient:  oasisClient,
		icpClient:    dfinityClient,
		sgxClient:    sgxClient,
		proverClient: proverClient,
		keyStore:     keyStore,
		prepared:     prepared,
		chainStore:   NewChainStore(store),
		minerAddr:    minerAddr,
		submitAddr:   keyStore.EthAddress(),
		icpSigMap:    make(map[string][][]byte),
		fileStore:    fileStore,
	}, nil
}

func (t *TxManager) init() error {
	return nil
}

func (t *TxManager) AddTask(resp *common.ProofResponse) {
	logger.Info("txManager add retry task: %v ,type:%v", resp.ProofId(), resp.ProofType.Name())
	unSubmitTx := NewDbUnSubmitTx(resp.Hash, hex.EncodeToString(resp.Proof), resp.ProofType)
	err := t.chainStore.WriteUnSubmitTx(unSubmitTx)
	if err != nil {
		logger.Error("write unsubmit tx error: %v %v", resp.ProofId(), err)
		return
	}
}

func (t *TxManager) Check() error {
	unSubmitTxs, err := t.chainStore.ReadUnSubmitTxs()
	if err != nil {
		logger.Error("read unsubmit tx error:%v", err)
		return err
	}
	for _, tx := range unSubmitTxs {
		switch tx.ProofType {
		case common.BtcDepositType, common.BtcUpdateCpType:
			hash, err := t.DepositBtc(tx.ProofType, tx.Hash, tx.Proof)
			if err != nil {
				logger.Error("update deposit error: %v %v", tx.ProofType.Name(), tx.Hash)
				continue
			}
			logger.Info("success update deposit txId: %v,hash: %v", tx.Hash, hash)
		case common.BtcChangeType:
			hash, err := t.UpdateUtxoChange(tx.Hash, tx.Proof)
			if err != nil {
				logger.Error("update utxo error: %v %v", tx.ProofType.Name(), tx.Hash)
				continue
			}
			logger.Info("success update utxo txId: %v,hash: %v", tx.Hash, hash)
		case common.RedeemTxType:
			hash, err := t.RedeemZkbtc(tx.Hash, tx.Proof)
			if err != nil {
				logger.Error("Redeem btx tx error: %v %v", tx.ProofType.Name(), tx.Hash)
				continue
			}
			logger.Debug("success Redeem btc ethHash: %v,btcHash: %v", tx.Hash, hash)
		default:
			logger.Warn("never should happen: %v %v", tx.ProofType.Name(), tx.Hash)
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
	dbNonce, exists, err := t.chainStore.ReadNonce(common.ETH.String(), addr)
	if err != nil {
		logger.Error("read nonce error: %v %v", addr, err)
		return 0, err
	}
	if !exists {
		return chainNonce, nil
	}
	if chainNonce <= dbNonce+1 {
		return dbNonce + 1, nil
	}
	return chainNonce, nil
}

func (t *TxManager) DepositBtc(proofType common.ProofType, txId, proof string) (string, error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	exists, err := t.ethClient.CheckUtxo(txId)
	if err != nil {
		logger.Error("check utxo error: %v %v", txId, err)
		return "", err
	}
	if exists {
		logger.Warn("deposit tx already exists: %v,skip it", txId)
		return "", nil
	}
	proofBytes := ethcommon.FromHex(proof)
	transaction, err := t.btcClient.GetHexRawTransaction(txId)
	if err != nil {
		logger.Error("get transaction error: %v %v", txId, err)
		return "", err
	}

	btcRawTx := ethcommon.FromHex(transaction)
	privateKey, err := t.keyStore.GetPrivateKey()
	if err != nil {
		logger.Error("get private key error: %v", err)
		return "", err
	}
	ethAddress := t.submitAddr
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
	gasPrice = getSuggestGasPrice(gasPrice)
	params, err := t.getParams(txId)
	if err != nil {
		logger.Error("get params error: %v %v", txId, err)
		return "", err
	}
	logger.Debug("submit deposit tx:%v cpDepth:%v,txDepth:%v,checkPoint:%x,blockHash:%x,blockTime:%v,flag:%v,smoothedTimestamp: %v,minerAddr:%v,gasPrice:%v,btcTxRaw:%x,proof:%v",
		txId, params.CpDepth, params.TxDepth, params.Checkpoint, params.TxBlockHash, params.TxTimestamp, params.Flag, params.SmoothedTimestamp, t.minerAddr, gasPrice, btcRawTx, proof)
	gasLimit, mockErr := t.ethClient.EstimateDepositGasLimit(t.submitAddr, params, gasPrice, btcRawTx, proofBytes)
	if mockErr != nil {
		switch proofType {
		case common.BtcUpdateCpType:
			logger.Warn("mock updateCp error:%v %v", txId, mockErr)
			err := t.chainStore.DeleteUnSubmitTx(txId)
			if err != nil {
				logger.Error("delete unSubmit tx error: %v", err)
				return "", err
			}
			return "", nil
		case common.BtcDepositType:
			logger.Error("deposit zkbtc error:%v %v", txId, mockErr)
			if strings.Contains(mockErr.Error(), "execution reverted") {
				err := t.chainStore.DeleteUnSubmitTx(txId)
				if err != nil {
					logger.Error("delete unSubmit tx error: %v", err)
					return "", err
				}
				if !proofFailed(mockErr) {
					logger.Warn("deposit tx expired now try again: %v", txId)
					err = t.addBtcUnGenProof(txId)
					if err != nil {
						logger.Error("add btc ungen proof error: %v", err)
						return "", err
					}
				}
				return "", nil
			}
			return "", mockErr
		default:
			logger.Warn("never should happen: %v %v", txId, mockErr)
			return "", mockErr
		}
	}
	gasLimit = getSuggestGasLimit(gasLimit)
	if err != nil {
		logger.Error("deposit zkbtc error:%v %v", txId, err)
		return "", err
	}
	balOk, err := t.CheckEthBalance(t.submitAddr, gasPrice, gasLimit)
	if err != nil {
		logger.Error("check balance error:%v", err)
		return "", err
	}
	if !balOk {
		return "", fmt.Errorf("balace check error")
	}

	txHash, err := t.ethClient.Deposit(privateKey, params, nonce, gasLimit, chainId, gasPrice, btcRawTx, proofBytes)
	if err != nil {
		logger.Error("deposit zkbtc error:%v %v", txId, err)
		return "", err
	}
	logger.Debug("deposit zkbtc info address: %v ethTxHash: %v, btcTxId: %v,nonce: %v", ethAddress, txHash, txId, nonce)
	err = t.chainStore.WriteNonce(common.ETH.String(), ethAddress, nonce)
	if err != nil {
		logger.Error("write nonce error: %v %v", ethAddress, err)
		return "", err
	}
	err = t.chainStore.DeleteUnSubmitTx(txId)
	if err != nil {
		logger.Error("delete unsubmit tx error: %v %v", txId, err)
		return "", err
	}
	return txHash, nil
}

func (t *TxManager) RedeemZkbtc(hash, proof string) (string, error) {
	ethTxHash := ethcommon.HexToHash(hash)
	ethTx, _, err := t.ethClient.TransactionByHash(context.Background(), ethTxHash)
	if err != nil {
		logger.Error("get eth tx error:%v %v", hash, err)
		return "", err
	}
	receipt, err := t.ethClient.TransactionReceipt(context.Background(), ethTxHash)
	if err != nil {
		logger.Error("get eth tx receipt error:%v %v", hash, err)
		return "", err
	}
	if len(receipt.Logs) < 4 {
		logger.Error("decode Redeem log error:%v %v", hash, err)
		return "", err
	}
	btcTxId, rewardBytes, btcRawTx, sigHashes, _, err := redeemUtils.DecodeRedeemReceipt(receipt)
	if err != nil {
		logger.Error("decode Redeem log error:%v %v", hash, err)
		return "", err
	}

	logger.Info("txId:%x,minerReward:%x,btxRawTx:%x,sigHashes:%x", btcTxId, rewardBytes[:], btcRawTx, sigHashes)

	minerReward := big.NewInt(0).SetBytes(rewardBytes[:])
	transaction := btctx.NewMultiTransactionBuilder()
	err = transaction.Deserialize(btcRawTx)
	if err != nil {
		logger.Error("deserialize btc tx error: %v %v", hash, err)
		return "", err
	}
	logger.Info("btcRawTx: %v\n", hexutil.Encode(btcRawTx))
	rawTx, rawReceipt := ethrpc.GetRawTxAndReceipt(ethTx, receipt)
	logger.Info("rawTx: %v\n", hexutil.Encode(rawTx))
	logger.Info("rawReceipt: %v\n", hexutil.Encode(rawReceipt))
	for index, vin := range transaction.MsgTx.TxIn {
		logger.Debug("%v vin: %v", index, vin)
	}

	scRoot, err := t.getTxScRoot(hash)
	if err != nil {
		logger.Error("get tx sc root error: %v %v", hash, err)
		return "", err
	}

	//currentScRoot, ethTxHash, ethUrl, btcTxId, proof string, sigHashes []string, minerReward *big.Int
	btcTxSignatures, err := t.SignerBtc(scRoot, hash, hex.EncodeToString(btcTxId[:]), proof, common.BytesArrayToHex(sigHashes), minerReward)
	if err != nil {
		logger.Error("sign btc tx error: %v %v", hash, err)
		return "", err
	}
	multiSigScriptBytes := ethcommon.FromHex(TestnetMultiSig)
	err = transaction.AddMultiScript(multiSigScriptBytes, 2, 3)
	if err != nil {
		logger.Error("add multi script error: %v %v", hash, err)
		return "", err
	}
	err = transaction.MergeSignature(btcTxSignatures)
	if err != nil {
		logger.Error("merge signature error: %v %v", hash, err)
		return "", err
	}
	btcTxBytes, err := transaction.Serialize()
	if err != nil {
		logger.Error("serialize btc tx error:%v %v", hash, err)
		return "", err
	}
	btcTxHash := transaction.TxHash()
	exists, err := t.btcTxOnChain(btcTxHash)
	if err != nil {
		logger.Error("get btc tx %v on chain error: %v %v", hash, btcTxHash, err)
		return "", err
	}
	if exists {
		logger.Warn("btc tx already exist: %v %v", btcTxHash, hash)
		err := t.chainStore.DeleteUnSubmitTx(hash)
		if err != nil {
			logger.Error("delete unSubmit tx error: %v %v", hash, err)
			return "", err
		}
		return "", nil
	}
	btcTxHex := hex.EncodeToString(btcTxBytes)
	logger.Info("Redeem btc: %v %v", btcTxHash, btcTxHex)
	txHash, err := t.btcClient.Sendrawtransaction(btcTxHex)
	if err != nil {
		logger.Error("send btc tx btcHash:%v,hash:%v %v", btcTxHash, hash, err)
		return "", err
	}
	logger.Info("send Redeem btc tx btcHash:%v  txHash:%v", btcTxHash, hash)
	err = t.chainStore.DeleteUnSubmitTx(hash)
	if err != nil {
		logger.Error("delete unSubmit tx error: %v %v", hash, err)
		return "", err
	}
	delete(t.icpSigMap, hash)
	return txHash, err
}

func (t *TxManager) btcTxOnChain(hash string) (bool, error) {
	_, err := t.btcClient.Getmempoolentry(hash)
	if err == nil {
		return true, nil
	}
	_, err = t.btcClient.GetRawTransaction(hash)
	if err == nil {
		return true, nil
	}
	return false, nil
}

func (t *TxManager) UpdateUtxoChange(txId, proof string) (string, error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	confirmed, err := t.ethClient.UtxoConfirm(txId)
	if err != nil {
		logger.Error("utxo confirm error: %v %v", txId, err)
		return "", err
	}
	if confirmed {
		logger.Warn("btc Redeem tx: %v confirmed,skip it", txId)
		return "", nil
	}
	txIdBytes := common.ReverseBytes(ethcommon.FromHex(txId))
	proofBytes := ethcommon.FromHex(proof)

	privateKey, err := t.keyStore.GetPrivateKey()
	if err != nil {
		logger.Error("get private key error: %v", err)
		return "", err
	}
	fromAddress := t.submitAddr
	nonce, err := t.getEthAddrNonce(fromAddress)
	if err != nil {
		logger.Error("get nonce error: %v %v", fromAddress, err)
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
	destHash, err := t.chainStore.ReadDestHash(txId)
	if err != nil {
		logger.Error("read dest hash error:%v", err)
		return "", err
	}
	receipt, err := t.ethClient.TransactionReceipt(context.Background(), ethcommon.HexToHash(destHash))
	if err != nil {
		logger.Error("get eth tx receipt error:%v %v", destHash, err)
		return "", err
	}
	if len(receipt.Logs) < 4 {
		logger.Error("decode Redeem log error:%v %v", destHash, err)
		return "", err
	}
	_, rewardBytes, _, _, _, err := redeemUtils.DecodeRedeemReceipt(receipt)
	if err != nil {
		logger.Error("decode Redeem log error:%v %v", destHash, err)
		return "", err
	}
	minerReward := big.NewInt(0).SetBytes(rewardBytes[:])
	gasPrice = getSuggestGasPrice(gasPrice)

	params, err := t.getParams(txId)
	if err != nil {
		logger.Error("get params %v error %v", txId, err)
		return "", err
	}
	logger.Debug("submit updateUtxo txId:%v, cpDepth:%v, txDepth:%v,blochHash:%x,cpHash:%x, blocktime:%v,flag:%v,smoothedTimestamp: %v,minerReward:%v,proof:%x",
		txId, params.CpDepth, params.TxDepth, params.TxBlockHash, params.Checkpoint, params.TxTimestamp, params.Flag, params.SmoothedTimestamp, minerReward.String(), proofBytes)

	gasLimit, mockErr := t.ethClient.EstimateUpdateUtxoGasLimit(t.submitAddr, params, gasPrice, minerReward, txIdBytes, proofBytes)
	if mockErr != nil {
		logger.Error("estimate update utxo gas limit error:%v %v", txId, mockErr)
		if strings.Contains(mockErr.Error(), "execution reverted") {
			err := t.chainStore.DeleteUnSubmitTx(txId)
			if err != nil {
				logger.Error("delete unSubmit tx error: %v", err)
				return "", err
			}
			if !proofFailed(mockErr) {
				logger.Warn("update utxo tx expired now,try again: %v", txId)
				err = t.addBtcUnGenProof(txId)
				if err != nil {
					logger.Error("add btc ungen proof error: %v", mockErr)
					return "", err
				}
			}
			return "", nil
		}
		return "", mockErr
	}
	gasLimit = getSuggestGasLimit(gasLimit)
	balOk, err := t.CheckEthBalance(t.submitAddr, gasPrice, gasLimit)
	if err != nil {
		logger.Error("check balance error:%v", err)
		return "", err
	}
	if !balOk {
		return "", fmt.Errorf("balace check error")
	}
	txHash, err := t.ethClient.UpdateUtxoChange(privateKey, params, nonce, gasLimit, chainId, gasPrice, minerReward, txIdBytes, proofBytes)
	if err != nil {
		logger.Error("update utxo change error:%v", err)
		return "", err
	}
	logger.Debug("update utxo info address: %v txHash: %v,txId: %v,nonce: %v", fromAddress, txHash, txId, nonce)
	err = t.chainStore.WriteNonce(common.ETH.String(), fromAddress, nonce)
	if err != nil {
		logger.Error("write nonce error: %v %v", fromAddress, err)
		return "", err
	}
	err = t.chainStore.DeleteUnSubmitTx(txId)
	if err != nil {
		logger.Error("delete unsubmit tx error: %v %v", txId, err)
		return "", err
	}
	return txHash, nil
}

func (t *TxManager) SignerBtc(currentScRoot, ethTxHash, btcTxId, proof string, sigHashes []string, minerReward *big.Int) ([][][]byte, error) {
	logger.Debug("signer btc tx currentScRoot:%v,ethTxHash:%v,btcTxId:%v,proof:%v,minerReward:%v,sigHashes:%v",
		currentScRoot, ethTxHash, btcTxId, proof, minerReward.Uint64(), sigHashes)
	var signatures [][][]byte
	oasisSignature, err := t.oasisClient.SignBtcTx(btcTxId, currentScRoot, proof, sigHashes, minerReward)
	if err != nil {
		logger.Error("oasis sign btc tx error: %v %v", btcTxId, err)
	} else {
		logger.Debug("txId:%v,oasis signature:%x", btcTxId, oasisSignature)
		signatures = append(signatures, oasisSignature...)
	}
	if icpSignatures, ok := t.icpSigMap[btcTxId]; ok {
		logger.Debug("txId:%v,use cache icp signature:%x", btcTxId, icpSignatures)
		signatures = append(signatures, icpSignatures)
	} else {
		//currentScRoot, ethTxHash, ethUrl, btcTxId, proof string, minerReward uint64, sigHashes []string
		icpTxSignatures, err := t.icpClient.BtcTxSign(currentScRoot, ethTxHash, btcTxId, proof, minerReward.String(), sigHashes)
		if err != nil {
			logger.Error("sign btc tx error: %v %v", btcTxId, err)
		} else {
			logger.Debug("txId:%v,icp signature:%v", btcTxId, icpTxSignatures)
			if icpTxSignatures.Signed {
				logger.Debug("txId:%v,icp signature:%v", btcTxId, icpTxSignatures.Signature)
				icpSignaturesBytes, err := icpSigToBytes(icpTxSignatures.Signature)
				if err != nil {
					logger.Error("sign btc tx error: %v %v", btcTxId, err)
					return nil, err
				}
				t.icpSigMap[btcTxId] = icpSignaturesBytes
				signatures = append(signatures, icpSignaturesBytes)
			}
		}
	}
	if len(signatures) == 0 {
		return nil, fmt.Errorf("signer faile: %v", btcTxId)
	}
	if len(signatures) == 2 {
		return signatures, nil
	}
	sgxRedeemProof, exists, err := t.fileStore.GetSgxRedeemProof(ethTxHash)
	if err != nil {
		logger.Error("get sgx Redeem proof error: %v %v", ethTxHash, err)
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("get sgx Redeem proof error: %v", ethTxHash)
	}
	sgxSignatures, err := t.sgxClient.BtcTxSignature(currentScRoot, minerReward.String(), btcTxId, sgxRedeemProof.Proof, sigHashes)
	if err != nil {
		logger.Error("sign btc tx error: %v %v", btcTxId, err)
		return nil, err
	}
	logger.Debug("txId:%v,sgx signature:%v", btcTxId, sgxSignatures.Signatures)
	sgxSignaturesBytes, err := sgxSigToBytes(sgxSignatures.Signatures)
	if err != nil {
		logger.Error("sign btc tx error: %v %v", btcTxId, err)
		return nil, err
	}
	signatures = append(signatures, sgxSignaturesBytes)
	return signatures, nil
}

func (t *TxManager) getTxScRoot(hash string) (string, error) {
	txes, err := t.chainStore.ReadDbTxes(hash)
	if err != nil {
		logger.Error("read db tx error:%v", err)
		return "", err
	}
	if len(txes) != 1 {
		logger.Warn("read db tx error:%v", err)
		return "", fmt.Errorf("read db tx error:%v", err)
	}
	slot, ok, err := t.chainStore.ReadSlotByHeight(txes[0].Height)
	if err != nil {
		logger.Error("read beacon slot error:%v", err)
		return "", err
	}
	if !ok {
		logger.Warn("read beacon slot error:%v", err)
		return "", fmt.Errorf("read beacon slot error:%v", err)
	}
	finalizedSlot, ok, err := t.fileStore.GetTxFinalizedSlot(slot)
	if err != nil {
		logger.Error("get tx finalized slot error: %v %v", slot, err)
		return "", err
	}
	if !ok {
		logger.Warn("no find tx finalized slot: %v", slot)
		return "", fmt.Errorf("no find tx finalized slot: %v", slot)
	}

	var currentFinalityUpdate common.LightClientFinalityUpdateEvent
	exists, err := t.fileStore.GetFinalityUpdate(finalizedSlot, &currentFinalityUpdate)
	if err != nil {
		logger.Error("get finality update error: %v %v", finalizedSlot, err)
		return "", err
	}
	if !exists {
		logger.Warn("no find finality update: %v", finalizedSlot)
		return "", fmt.Errorf("no find finality update: %v", finalizedSlot)
	}
	attestedSlot, err := strconv.ParseUint(currentFinalityUpdate.Data.AttestedHeader.Slot, 10, 64)
	if err != nil {
		logger.Error("parse big error %v %v", currentFinalityUpdate.Data.AttestedHeader.Slot, err)
		return "", err
	}
	period := attestedSlot / common.SlotPerPeriod
	update, ok, err := t.prepared.GetSyncCommitUpdate(period)
	if err != nil {
		logger.Error("read update error:%v", err)
		return "", err
	}
	if !ok {
		logger.Warn("read update error:%v", err)
		return "", fmt.Errorf("read update error:%v", err)
	}
	syncCommitRoot, err := circuits.SyncCommitRoot(update.CurrentSyncCommittee)
	if err != nil {
		logger.Error("get syncCommitRoot error:%v", err)
		return "", err
	}
	return hex.EncodeToString(syncCommitRoot), nil

}

func (t *TxManager) getParams(txId string) (*zkbridge.IBtcTxVerifierPublicWitnessParams, error) {
	dbTx, ok, err := t.chainStore.ReadBtcTx(txId)
	if err != nil {
		logger.Error("read btc tx error: %v %v", txId, err)
		return nil, err
	}
	if !ok {
		logger.Warn("no find btc tx: %v", txId)
		return nil, fmt.Errorf("no find btc tx:%v", txId)
	}
	cpDepth := dbTx.LatestHeight - dbTx.CheckPointHeight
	txDepth := dbTx.LatestHeight - dbTx.Height
	cpHash, ok, err := t.chainStore.ReadCheckpoint(dbTx.CheckPointHeight)
	if err != nil {
		logger.Error("%v", err.Error())
		return nil, err
	}
	btcTx, err := t.btcClient.GetRawTransaction(dbTx.Hash)
	if err != nil {
		logger.Error("get transaction error: %v %v", dbTx.Hash, err)
		return nil, err
	}
	blockHash := common.ReverseBytes(ethcommon.FromHex(btcTx.Blockhash))

	icpSignature, ok, err := t.chainStore.ReadIcpSignature(dbTx.LatestHeight)
	if err != nil {
		logger.Error("read dfinity sign error: %v", err)
		return nil, err
	}
	if !ok {
		logger.Warn("no find: %v icp %v signature", dbTx.Hash, dbTx.LatestHeight)
		// no work,just placeholder
		icpSignature.Hash = "6aeb6ec6f0fbc707b91a3bec690ae6536fe0abaa1994ef24c3463eb20494785d"
		icpSignature.Signature = "3f8e02c743e76a4bd655873a428db4fa2c46ac658854ba38f8be0fbbf9af9b2b6b377aaaaf231b6b890a5ee3c15a558f1ccc18dae0c844b6f06343b88a8d12e3"
	}
	smoothedTimestamp, err := blockdepthUtil.GetSmoothedTimestampProofData(t.proverClient, uint32(dbTx.LatestHeight))
	if err != nil {
		logger.Error("%v", err.Error())
		return nil, err
	}
	cptimeData, err := blockdepthUtil.GetCpTimestampProofData(t.proverClient, uint32(dbTx.Height))
	if err != nil {
		logger.Error("%v", err)
		return nil, err
	}
	sigVerif, err := blockdepthUtil.GetSigVerifProofData(
		ethcommon.FromHex(icpSignature.Hash),
		ethcommon.FromHex(icpSignature.Signature),
		ethcommon.FromHex(TestnetIcpPublicKey))
	if err != nil {
		logger.Error("%v", err.Error())
		return nil, err
	}
	flag := cptimeData.Flag<<1 | sigVerif.Flag
	params := &zkbridge.IBtcTxVerifierPublicWitnessParams{
		Checkpoint:        [32]byte(ethcommon.FromHex(cpHash)),
		CpDepth:           uint32(cpDepth),
		TxDepth:           uint32(txDepth),
		TxBlockHash:       [32]byte(blockHash),
		TxTimestamp:       uint32(btcTx.Blocktime),
		ZkpMiner:          ethcommon.HexToAddress(t.minerAddr),
		Flag:              big.NewInt(int64(flag)),
		SmoothedTimestamp: smoothedTimestamp.Timestamp,
	}
	return params, nil
}

func (t *TxManager) CheckEthBalance(addr string, gasPrice *big.Int, gasLimit uint64) (bool, error) {
	balance, err := t.ethClient.EthBalance(addr)
	if err != nil {
		return false, nil
	}
	gasFee := big.NewInt(0).Mul(gasPrice, big.NewInt(int64(gasLimit)))
	fixGasFee := big.NewInt(0).Mul(gasFee, big.NewInt(2)) // todo
	if balance.Cmp(fixGasFee) <= 0 {
		logger.Error("not enough gas fee to submit proof,please deposit gasFee:%v", addr)
		return false, nil
	}
	return true, nil

}

func (t *TxManager) addBtcUnGenProof(txId string) error {
	tx, ok, err := t.chainStore.ReadBtcTx(txId)
	if err != nil {
		logger.Error("read btc tx error: %v %v", txId, err)
		return err
	}
	if !ok {
		logger.Warn("no find btc tx: %v", txId)
		return fmt.Errorf("no find btc tx:%v", txId)
	}
	err = t.fileStore.DelProof(NewHashStoreKey(tx.ProofType, txId))
	if err != nil {
		logger.Error("del proof error: %v", err)
		//return err
	}
	// re select latest height to gen proof
	tx.LatestHeight = 0
	tx.CheckPointHeight = 0
	tx.GenProofNums = tx.GenProofNums + 1
	err = t.chainStore.WriteDbTxes(tx)
	if err != nil {
		logger.Error("write db tx error: %v", err)
		return err
	}
	err = t.chainStore.WriteUnGenProof(common.BitcoinChain, &DbUnGenProof{
		ChainType: tx.ChainType,
		ProofType: tx.ProofType,
		Hash:      tx.Hash,
		Height:    tx.Height,
		TxIndex:   tx.TxIndex,
		Amount:    uint64(tx.Amount),
	})
	if err != nil {
		logger.Error("write ungen proof error: %v", err)
		return err
	}
	logger.Debug("add gen btc proof again:%v %v %v", tx.ChainType.String(), tx.ProofType.Name(), txId)
	return nil
}

func (t *TxManager) Close() error {
	return nil
}

func NewDbUnSubmitTx(hash, proof string, proofType common.ProofType) DbUnSubmitTx {
	return DbUnSubmitTx{
		Hash:      hash,
		Proof:     proof,
		ProofType: proofType,
		Timestamp: time.Now().UnixNano(),
	}
}

func icpSigToBytes(signatures []string) ([][]byte, error) {
	var signaturesBytes [][]byte
	for _, sig := range signatures {
		sigBytes, err := RsToSignature(sig)
		if err != nil {
			return nil, err
		}
		signaturesBytes = append(signaturesBytes, append(sigBytes, 0x01)) //todo
	}
	return signaturesBytes, nil
}

func sgxSigToBytes(signatures []string) ([][]byte, error) {
	var signaturesBytes [][]byte
	for _, sig := range signatures {
		sigBytes, err := hex.DecodeString(sig)
		if err != nil {
			return nil, err
		}
		signaturesBytes = append(signaturesBytes, append(sigBytes, 0x01)) //todo
	}
	return signaturesBytes, nil
}

func getSuggestGasLimit(value uint64) uint64 {
	gasLimit := big.NewInt(0).Div(
		big.NewInt(0).Mul(big.NewInt(int64(value)), big.NewInt(3)),
		big.NewInt(2))
	return gasLimit.Uint64()

}
func getSuggestGasPrice(value *big.Int) *big.Int {
	gasPrice := big.NewInt(0).Div(
		big.NewInt(0).Mul(value, big.NewInt(3)),
		big.NewInt(2))
	return gasPrice
}

func proofFailed(err error) bool {
	//todo
	switch err.Error() {
	case "execution reverted: no practical use operation":
		return true
	default:
		return false
	}
}

func proofExpired(err error) bool {
	// todo
	switch err.Error() {
	case "execution reverted: txDepth check failed":
		return true
	case "execution reverted: cpDepth check failed":
		return true
	case "execution reverted: deposit proof verification failed":
		return true
	case "execution reverted: deposit to previous need role":
		return true
	}
	return false

}
