package node

import (
	"encoding/hex"
	"fmt"
	ethcommon "github.com/ethereum/go-ethereum/common"
	btcproverUtils "github.com/lightec-xyz/btc_provers/utils"
	btcproverClient "github.com/lightec-xyz/btc_provers/utils/client"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/rpc/bitcoin"
	"github.com/lightec-xyz/daemon/rpc/bitcoin/types"
	"github.com/lightec-xyz/daemon/rpc/ethereum"
	"github.com/lightec-xyz/daemon/store"
	"math/big"
	"strings"
)

type BitcoinAgent struct {
	btcClient       *bitcoin.Client
	ethClient       *ethereum.Client
	btcProverClient *btcproverClient.Client
	store           store.IStore
	memoryStore     store.IStore
	fileStore       *FileStorage
	cache           *CacheState
	proofRequest    chan<- []*common.ZkProofRequest
	operatorAddr    string
	submitTxEthAddr string
	keyStore        *KeyStore
	minDepositValue float64
	initHeight      int64
	txManager       *TxManager
}

func NewBitcoinAgent(cfg Config, submitTxEthAddr string, store, memoryStore store.IStore, fileStore *FileStorage, btcClient *bitcoin.Client,
	ethClient *ethereum.Client, btcProverClient *btcproverClient.Client, requests chan []*common.ZkProofRequest, keyStore *KeyStore, task *TxManager) (IAgent, error) {
	return &BitcoinAgent{
		btcClient:       btcClient,
		ethClient:       ethClient,
		store:           store,
		memoryStore:     memoryStore,
		operatorAddr:    cfg.BtcOperatorAddr,
		proofRequest:    requests,
		minDepositValue: 0, // todo
		btcProverClient: btcProverClient,
		keyStore:        keyStore,
		submitTxEthAddr: submitTxEthAddr,
		fileStore:       fileStore,
		initHeight:      cfg.BtcInitHeight,
		txManager:       task,
		cache:           NewCacheState(),
	}, nil
}

func (b *BitcoinAgent) Init() error {
	logger.Info("bitcoin agent init now")
	height, exists, err := ReadBitcoinHeight(b.store)
	if err != nil {
		logger.Error("get btc current height error:%v", err)
		return err
	}
	if !exists || height < b.initHeight {
		logger.Debug("init btc current height: %v", b.initHeight)
		err := WriteBitcoinHeight(b.store, b.initHeight)
		if err != nil {
			logger.Error("put init btc current height error:%v", err)
			return err
		}
	}
	// test rpc
	_, err = b.btcClient.GetBlockCount()
	if err != nil {
		logger.Error(" bitcoin json rpc get block count error:%v", err)
		return err
	}
	logger.Info("init bitcoin agent completed")
	return nil
}

// checkUnGenerateProof check uncompleted generate Proof tx,resend again
func (b *BitcoinAgent) checkUnGenerateProof() error {
	// todo
	return nil
}

func (b *BitcoinAgent) ScanBlock() error {
	logger.Debug("bitcoin scan block ...")
	curHeight, ok, err := ReadBitcoinHeight(b.store)
	if err != nil {
		logger.Error("get btc current height error:%v", err)
		return err
	}
	if !ok {
		logger.Error("never should happen")
		return fmt.Errorf("no btc current height")
	}
	blockCount, err := b.btcClient.GetBlockCount()
	if err != nil {
		logger.Error("bitcoin client get block count error:%v", err)
		return err
	}
	forked, err := b.CheckChainFork(blockCount)
	if err != nil {
		logger.Error("bitcoin chain fork error:%v %v", blockCount, err)
		return err
	}
	if forked {
		logger.Error("bitcoin chain forked,need to rollback %v ", blockCount, err)
		// todo
		//return nil
	}
	blockCount = blockCount - 1
	if curHeight >= blockCount {
		logger.Debug("btc current height:%d,node block count:%d", curHeight, blockCount)
		return nil
	}
	for index := curHeight + 1; index <= blockCount; index++ {
		if index%10 == 0 {
			logger.Debug("bitcoin parse block height:%d", index)
		}
		depositTxes, redeemTxes, proofRequests, err := b.parseBlock(uint64(index))
		if err != nil {
			logger.Error("bitcoin agent parse block error: %v %v", index, err)
			return err
		}
		err = b.updateRedeemInfo(index, redeemTxes)
		if err != nil {
			logger.Error("bitcoin agent update redeem info error: %v %v", index, err)
			return err
		}
		allTxes := append(depositTxes, redeemTxes...)
		err = b.saveData(index, allTxes)
		if err != nil {
			logger.Error("bitcoin agent save transaction error: %v %v", index, err)
			return err
		}
		err = WriteBitcoinHeight(b.store, index)
		if err != nil {
			logger.Error("write btc height error: %v %v", index, err)
			return err
		}
		if len(proofRequests) > 0 {
			b.SendProofRequest(proofRequests...)
		}

	}
	return nil
}

func (b *BitcoinAgent) CheckChainFork(height int64) (bool, error) {
	// todo
	return false, nil
}

func (b *BitcoinAgent) SendProofRequest(requests ...*common.ZkProofRequest) {
	if len(requests) > 0 {
		b.proofRequest <- requests
	}
	for _, request := range requests {
		logger.Info(" btc agent success send btc proof request: %v", request.Id())
		b.cache.Store(request.Id(), nil)
	}
}

func (b *BitcoinAgent) updateRedeemInfo(height int64, txList []*Transaction) error {
	//todo
	return nil
}

func (b *BitcoinAgent) saveData(height int64, txes []*Transaction) error {
	err := WriteTxes(b.store, txesToDbTxes(txes))
	if err != nil {
		logger.Error("put redeem tx error: %v %v", height, err)
		return err
	}
	err = WriteDbProof(b.store, txesToDbProofs(txes))
	if err != nil {
		logger.Error("write Proof error: %v", err)
		return err
	}
	err = WriteUnGenProof(b.store, Bitcoin, txesToUnGenProofs(Bitcoin, txes))
	if err != nil {
		logger.Error("write ungen Proof error:%v", err)
		return err
	}
	return nil
}

func (b *BitcoinAgent) parseBlock(height uint64) ([]*Transaction, []*Transaction, []*common.ZkProofRequest, error) {
	blockHash, err := b.btcClient.GetBlockHash(int64(height))
	if err != nil {
		logger.Error("btcClient get block hash error: %v %v", height, err)
		return nil, nil, nil, err
	}
	blockWithTx, err := b.btcClient.GetBlock(blockHash)
	if err != nil {
		logger.Error("btcClient get block error: %v %v", blockHash, err)
		return nil, nil, nil, err
	}
	var depositTxes []*Transaction
	var redeemTxes []*Transaction
	var requests []*common.ZkProofRequest
	for _, tx := range blockWithTx.Tx {
		redeemTx, isRedeem := b.isRedeemTx(tx, height, blockHash)
		if isRedeem {
			proofed, err := b.CheckChainProof(common.VerifyTxType, tx.Txid)
			if err != nil {
				logger.Error("bitcoin check chain proof error: %v %v", tx.Txid, err)
				return nil, nil, nil, err
			}
			if proofed {
				redeemTx.Proofed = true
			} else {
				proofData, err := btcproverUtils.GetDefaultGrandRollupProofData(b.btcProverClient, tx.Txid, blockHash)
				if err != nil {
					logger.Error("get verify proof data error: %v %v", tx.Txid, err)
					return nil, nil, nil, err
				}
				data := rpc.VerifyRequest{
					TxHash:    tx.Txid,
					BlockHash: blockHash,
					Data:      proofData,
				}
				requests = append(requests, common.NewZkProofRequest(common.VerifyTxType, data, 0, tx.Txid))
			}
			redeemTxes = append(redeemTxes, redeemTx)
			continue
		}
		depositTx, isDeposit, err := parseDepositTx(tx, b.operatorAddr, height, b.minDepositValue)
		if err != nil {
			logger.Error("check deposit tx error: %v %v", tx.Txid, err)
			return nil, nil, nil, err
		}
		if isDeposit {
			proofed, err := b.CheckChainProof(common.DepositTxType, depositTx.TxHash)
			if err != nil {
				logger.Error("check deposit Proof error: %v %v", tx.Txid, err)
				return nil, nil, nil, err
			}
			if proofed {
				depositTx.Proofed = true
			} else {
				proofData, err := btcproverUtils.GetDefaultGrandRollupProofData(b.btcProverClient, tx.Txid, blockHash)
				if err != nil {
					logger.Error("get deposit proof data error: %v %v", tx.Txid, err)
					return nil, nil, nil, err
				}
				data := rpc.DepositRequest{
					TxHash:    tx.Txid,
					BlockHash: blockHash,
					Data:      proofData,
				}
				requests = append(requests, common.NewZkProofRequest(common.DepositTxType, data, 0, tx.Txid))
			}
			depositTxes = append(depositTxes, depositTx)
		}
	}
	return depositTxes, redeemTxes, requests, nil
}

func (b *BitcoinAgent) CheckChainProof(proofType common.ZkProofType, txHash string) (bool, error) {
	// todo
	return false, nil
}

func (b *BitcoinAgent) ProofResponse(resp *common.ZkProofResponse) error {
	logger.Info("bitcoinAgent receive Proof resp: %v %v %v %x",
		resp.ZkProofType.String(), resp.Period, resp.TxHash, resp.Proof)
	err := StoreZkProof(b.fileStore, resp.ZkProofType, resp.Period, resp.TxHash, resp.Proof, resp.Witness)
	if err != nil {
		logger.Error("store Proof error: %v %v", resp.TxHash, err)
		return err
	}
	proofId := resp.TxHash
	err = b.updateDepositProof(proofId, hex.EncodeToString(resp.Proof), resp.Status)
	if err != nil {
		logger.Error("update Proof error: %v %v", proofId, err)
		return err
	}
	switch resp.ZkProofType {
	case common.VerifyTxType:
		logger.Info("start update utxo change: %v", proofId)
		err := updateContractUtxoChange(b.ethClient, b.submitTxEthAddr, b.keyStore.GetPrivateKey(), []string{resp.TxHash}, resp.Proof)
		if err != nil {
			logger.Error("update utxo error: %v %v", proofId, err)
			b.txManager.AddTask(resp)
			return err
		}
	default:
	}
	return nil
}

func updateContractUtxoChange(ethClient *ethereum.Client, address, privateKey string, txIds []string, proof []byte) error {
	// todo
	nonce, err := ethClient.GetNonce(address)
	if err != nil {
		logger.Error("get  nonce error:%v", err)
		return err
	}
	chainId, err := ethClient.GetChainId()
	if err != nil {
		logger.Error("get chain id error:%v", err)
		return err
	}
	gasPrice, err := ethClient.GetGasPrice()
	if err != nil {
		logger.Error("get gas price error:%v", err)
		return err
	}
	gasLimit := uint64(500000)
	gasPrice = big.NewInt(0).Mul(gasPrice, big.NewInt(2))
	txHash, err := ethClient.UpdateUtxoChange(privateKey, txIds, nonce, gasLimit, chainId, gasPrice, proof)
	if err != nil {
		logger.Error("update utxo change error:%v", err)
		return err
	}
	logger.Info("success send update utxo change  hash:%v", txHash)
	return nil
}

func (b *BitcoinAgent) isRedeemTx(tx types.Tx, height uint64, blockHash string) (*Transaction, bool) {
	// todo more check
	var inputs []Utxo
	isRedeemTx := false
	for _, vin := range tx.Vin {
		if vin.Prevout.ScriptPubKey.Address == b.operatorAddr {
			isRedeemTx = true
		}
		inputs = append(inputs, Utxo{
			TxId:  vin.TxId,
			Index: uint32(vin.Vout),
		})
	}
	var outputs []TxOut
	for _, out := range tx.Vout {
		scriptHex, err := hex.DecodeString(out.ScriptPubKey.Hex)
		if err != nil {
			logger.Error("decode hex error:%v %v", tx.Txid, err)
			return nil, false
		}
		outputs = append(outputs, TxOut{
			Value:    BtcToSat(out.Value),
			PkScript: scriptHex,
		})
	}
	if isRedeemTx {
		logger.Info("bitcoin agent find redeem tx: %v,inputs:%v ,outputs:%v", tx.Txid, formatUtxo(inputs), formatOut(outputs))
	}
	redeemBtcTx := NewRedeemBtcTx(height, tx.Txid, blockHash, inputs, outputs)
	return redeemBtcTx, isRedeemTx
}

func (b *BitcoinAgent) updateDepositProof(txId string, proof string, status common.ProofStatus) error {
	logger.Debug("bitcoin update Proof status: %v %v %v", txId, proof, status)
	err := UpdateProof(b.store, txId, proof, common.DepositTxType, status)
	if err != nil {
		logger.Error("update Proof error: %v %v", txId, err)
		return err
	}
	return nil

}

func (b *BitcoinAgent) CheckState() error {
	unGenProofs, err := ReadAllUnGenProofs(b.store, Bitcoin)
	if err != nil {
		logger.Error("read unGen proof error:%v", err)
		return err
	}
	for _, proof := range unGenProofs {
		logger.Debug("bitcoin check ungen proof: %v %v", proof.ProofType.String(), proof.TxHash)
		if proof.ProofType == 0 || proof.TxHash == "" {
			logger.Warn("unGenProof error:%v %v", proof.ProofType.String(), proof.TxHash)
			err := DeleteUnGenProof(b.store, Bitcoin, proof.TxHash)
			if err != nil {
				logger.Error("delete ungen proof error:%v %v", proof.TxHash, err)
			}
			continue
		}
		exists, err := CheckProof(b.fileStore, proof.ProofType, 0, proof.TxHash)
		if err != nil {
			logger.Error("check proof error:%v %v", proof.TxHash, err)
			return nil
		}
		if exists {
			logger.Debug("%v %v proof exists ,delete ungen proof now", proof.ProofType.String(), proof.TxHash)
			err = DeleteUnGenProof(b.store, Bitcoin, proof.TxHash)
			if err != nil {
				logger.Error("delete ungen proof error:%v %v", proof.TxHash, err)
				return nil
			}
			continue
		}
		err = b.tryProofRequest(proof.ProofType, proof.TxHash)
		if err != nil {
			logger.Error("try proof request error:%v %v", proof.TxHash, err)
			return nil
		}
	}
	return nil
}

func (b *BitcoinAgent) tryProofRequest(proofType common.ZkProofType, txHash string) error {
	proofId := common.NewProofId(proofType, 0, txHash)
	exists := b.cache.Check(proofId)
	if exists {
		logger.Debug("proof request exists: %v", proofId)
		return nil
	}
	exists, err := CheckProof(b.fileStore, proofType, 0, txHash)
	if err != nil {
		logger.Error("check proof error:%v %v", txHash, err)
		return err
	}
	if exists {
		return nil
	}
	data, ok, err := b.getRequestData(proofType, txHash)
	if err != nil {
		logger.Error("get request data error:%v %v", txHash, err)
		return err
	}
	if !ok {
		return nil
	}
	zkProofRequest := common.NewZkProofRequest(proofType, data, 0, txHash)
	b.SendProofRequest(zkProofRequest)
	return nil
}

func (b *BitcoinAgent) getRequestData(proofType common.ZkProofType, txHash string) (interface{}, bool, error) {
	switch proofType {
	case common.DepositTxType:
		data, ok, err := b.getDepositData(txHash)
		if err != nil {
			logger.Error("get deposit data error:%v %v", txHash, err)
			return nil, false, err
		}
		return data, ok, nil
	case common.VerifyTxType:
		data, ok, err := b.getVerifyData(txHash)
		if err != nil {
			logger.Error("get verify data error:%v %v", txHash, err)
			return nil, false, err
		}
		return data, ok, nil
	default:
		return nil, false, fmt.Errorf("unknown proof type: %v", proofType)
	}

}

func (b *BitcoinAgent) getDepositData(txHash string) (interface{}, bool, error) {
	tx, err := b.btcClient.GetTransaction(txHash)
	if err != nil {
		logger.Error("get deposit tx error: %v %v", txHash, err)
		return nil, false, err
	}
	proofData, err := btcproverUtils.GetDefaultGrandRollupProofData(b.btcProverClient, txHash, tx.Blockhash)
	if err != nil {
		logger.Error("get deposit proof data error: %v %v", txHash, err)
		return nil, false, err
	}
	depositRequest := rpc.DepositRequest{
		TxHash:    txHash,
		BlockHash: tx.Blockhash,
		Data:      proofData,
	}
	return depositRequest, true, nil
}

func (b *BitcoinAgent) getVerifyData(txHash string) (interface{}, bool, error) {
	tx, err := b.btcClient.GetTransaction(txHash)
	if err != nil {
		logger.Error("get verify tx error: %v %v", txHash, err)
		return nil, false, err
	}
	proofData, err := btcproverUtils.GetDefaultGrandRollupProofData(b.btcProverClient, txHash, tx.Blockhash)
	if err != nil {
		logger.Error("get verify proof data error: %v %v", txHash, err)
		return nil, false, err
	}
	verifyRequest := rpc.VerifyRequest{
		TxHash:    txHash,
		BlockHash: tx.Blockhash,
		Data:      proofData,
	}
	return verifyRequest, true, nil
}

func (b *BitcoinAgent) Close() error {
	return nil
}
func (b *BitcoinAgent) Name() string {
	return "bitcoinAgent"
}

func parseDepositTx(tx types.Tx, operatorAddr string, height uint64, minDepositValue float64) (*Transaction, bool, error) {
	// todo more rule
	txOuts := tx.Vout
	if len(txOuts) < 2 {
		return nil, false, nil
	}
	amount, isDeposit, err := isContainOperator(tx.Vout, operatorAddr)
	if err != nil {
		return nil, false, err
	}
	if !isDeposit {
		return nil, false, nil
	}

	ethAddr, ok, err := getOPReturn(tx.Vout)
	if !ok {
		return nil, false, nil
	}
	utxoList := []Utxo{
		{
			TxId:  tx.Txid,
			Index: 1,
		},
	}
	logger.Info("bitcoin agent find  deposit tx: %v, ethAddr:%v,amount:%v,utxo:%v", tx.Txid, ethAddr, amount, formatUtxo(utxoList))
	depositTx := NewDepositBtcTx(height, tx.Txid, ethAddr, utxoList, BtcToSat(amount))
	return depositTx, true, nil
}

func isContainOperator(txOuts []types.TxVout, operatorAddr string) (float64, bool, error) {
	var isDeposit bool
	var total float64
	for _, out := range txOuts {
		if out.ScriptPubKey.Address == operatorAddr {
			isDeposit = true
			total = total + out.Value
		}
	}
	return total, isDeposit, nil
}

func getOPReturn(txOuts []types.TxVout) (string, bool, error) {
	for _, out := range txOuts {
		if out.ScriptPubKey.Type == "nulldata" && strings.HasPrefix(out.ScriptPubKey.Hex, "6a") {
			ethAddr, err := getEthAddrFromScript(out.ScriptPubKey.Hex)
			if err != nil {
				return "", false, err
			}
			return ethAddr, true, nil
		}
	}
	return "", false, nil
}

func getEthAddrFromScript(script string) (string, error) {
	// todo
	// example https://live.blockcypher.com/btc-testnet/tx/fa1bee4165f1720b33047792e47743aeb406940f4b2527874929db9cdbb9da42/
	if len(script) < 5 {
		return "", fmt.Errorf("scritp lenght is less than 4")
	}
	if !strings.HasPrefix(script, "6a") {
		return "", fmt.Errorf("script is not start with 6a")
	}
	isHexAddress := ethcommon.IsHexAddress(script[4:])
	if !isHexAddress {
		return "", fmt.Errorf("script is not hex address")
	}
	return script[4:], nil
}

func NewDepositBtcTx(height uint64, txId, ethAddr string, utxo []Utxo, amount int64) *Transaction {
	return &Transaction{
		TxHash:    txId,
		Height:    height,
		TxType:    DepositTx,
		ChainType: Bitcoin,
		EthAddr:   ethAddr,
		ProofType: common.DepositTxType,
		Utxo:      utxo,
		Amount:    amount,
	}
}

func NewRedeemBtcTx(height uint64, txId, blockHash string, inputs []Utxo, outputs []TxOut) *Transaction {
	return &Transaction{
		TxHash:    txId,
		Height:    height,
		TxType:    RedeemTx,
		ChainType: Bitcoin,
		ProofType: common.VerifyTxType,
		BlockHash: blockHash,
	}
}
