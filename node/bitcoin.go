package node

import (
	"encoding/hex"
	"fmt"
	ethcommon "github.com/ethereum/go-ethereum/common"
	btcproverClient "github.com/lightec-xyz/btc_provers/utils/client"
	grUtil "github.com/lightec-xyz/btc_provers/utils/grandrollup"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/rpc/bitcoin"
	"github.com/lightec-xyz/daemon/rpc/bitcoin/types"
	"github.com/lightec-xyz/daemon/rpc/ethereum"
	"github.com/lightec-xyz/daemon/store"
	"strings"
)

type BitcoinAgent struct {
	btcClient       *bitcoin.Client
	ethClient       *ethereum.Client
	btcProverClient *btcproverClient.Client
	store           store.IStore
	memoryStore     store.IStore
	fileStore       *FileStorage
	cache           *Cache
	proofRequest    chan<- []*common.ZkProofRequest
	operatorAddr    string
	keyStore        *KeyStore
	minDepositValue float64
	initHeight      int64
	txManager       *TxManager
	debug           bool
	force           bool
}

func (b *BitcoinAgent) FetchDataResponse(resp *FetchResponse) error {
	// todo
	return nil
}

func NewBitcoinAgent(cfg Config, store, memoryStore store.IStore, fileStore *FileStorage, btcClient *bitcoin.Client,
	ethClient *ethereum.Client, btcProverClient *btcproverClient.Client, requests chan []*common.ZkProofRequest, keyStore *KeyStore,
	task *TxManager, state *Cache) (IAgent, error) {
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
		fileStore:       fileStore,
		initHeight:      cfg.BtcInitHeight,
		txManager:       task,
		cache:           state,
		debug:           true, //common.GetEnvDebugMode(),
	}, nil
}

func (b *BitcoinAgent) Init() error {
	logger.Info("bitcoin agent init now")
	if b.force {
		err := WriteBitcoinHeight(b.store, b.initHeight)
		if err != nil {
			logger.Error("write btc height error: %v %v", b.initHeight, err)
			return err
		}
	} else {
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
	}
	// test rpc
	_, err := b.btcClient.GetBlockCount()
	if err != nil {
		logger.Error(" bitcoin json rpc get block count error:%v", err)
		return err
	}
	if b.debug {
		err = WriteBitcoinHeight(b.store, 2835035)
		if err != nil {
			logger.Error("write btc height error: %v %v", 2835035, err)
			return err
		}
	}

	logger.Info("init bitcoin agent completed")
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
		logger.Error("bitcoin chain forked,need to rollback %v %v ", blockCount, err)
		// todo
		//return nil
	}
	if b.debug {
		blockCount = 2835038
		if curHeight >= 2835037 {
			return nil
		}
	}
	//todo
	blockCount = blockCount - 1
	if curHeight >= blockCount {
		logger.Debug("btc current height:%d,node block count:%d", curHeight, blockCount)
		return nil
	}
	for index := curHeight + 1; index <= blockCount; index++ {
		logger.Debug("bitcoin parse block height:%d", index)
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
				logger.Debug("btc redeem tx proofed: %v", tx.Txid)
			} else {
				proofData, err := grUtil.GetDefaultGrandRollupProofData(b.btcProverClient, tx.Txid, blockHash)
				if err != nil {
					logger.Error("get verify proof data error: %v %v", tx.Txid, err)
					return nil, nil, nil, err
				}
				data := rpc.VerifyRequest{
					TxHash:    tx.Txid,
					BlockHash: blockHash,
					Data:      proofData,
				}
				requests = append(requests, common.NewZkProofRequest(common.VerifyTxType, data, 0, 0, tx.Txid))
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
				logger.Debug("btc deposit tx proofed: %v", tx.Txid)
			} else {
				//proofData, err := btcproverUtils.GetDefaultGrandRollupProofData(b.btcProverClient, tx.Txid, blockHash)
				//if err != nil {
				//	logger.Error("get deposit proof data error: %v %v", tx.Txid, err)
				//	return nil, nil, nil, err
				//}
				//data := rpc.DepositRequest{
				//	TxHash:    tx.Txid,
				//	BlockHash: blockHash,
				//	Data:      proofData,
				//}
				//requests = append(requests, common.NewZkProofRequest(common.DepositTxType, data, 0, 0, tx.Txid))
			}
			depositTxes = append(depositTxes, depositTx)
		}
	}
	return depositTxes, redeemTxes, requests, nil
}

func (b *BitcoinAgent) CheckChainProof(proofType common.ZkProofType, txHash string) (bool, error) {
	switch proofType {
	case common.VerifyTxType:
		utxo, err := b.ethClient.GetUtxo(txHash)
		if err != nil {
			logger.Error("check utxo error: %v %v", txHash, err)
			return false, err
		}
		return utxo.IsChangeConfirmed, nil
	case common.DepositTxType:
		utxo, err := b.ethClient.GetUtxo(txHash)
		if err != nil {
			logger.Error("check tx error: %v %v", txHash, err)
			return false, err
		}
		if TxIdIsEmpty(utxo.Txid) {
			return false, nil
		}
		return true, nil
	default:
		return false, fmt.Errorf("unsupported proof type: %v", proofType)
	}
}

func (b *BitcoinAgent) ProofResponse(resp *common.ZkProofResponse) error {
	logger.Info("bitcoinAgent receive Proof resp: %v %x", resp.Id(), resp.Proof)
	b.cache.Delete(resp.Id())
	switch resp.ZkProofType {
	case common.DepositTxType:
		err := b.updateDepositProof(resp.TxHash, hex.EncodeToString(resp.Proof), resp.Status)
		if err != nil {
			logger.Error("update Proof error: %v %v", resp.TxHash, err)
			return err
		}
	case common.VerifyTxType:
		logger.Info("start update utxo change: %v", resp.TxHash)
		hash, err := b.txManager.UpdateUtxoChange(resp.TxHash, hex.EncodeToString(resp.Proof))
		if err != nil {
			logger.Error("update utxo fail: %v %v,save to db", resp.TxHash, err)
			b.txManager.AddTask(resp)
			return err
		}
		logger.Debug("success update utxo: txId:%v hash:%v", resp.TxHash, hash)

	default:
	}
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
	for _, tx := range unGenProofs {
		logger.Debug("bitcoin check ungen proof: %v %v", tx.ProofType.String(), tx.TxHash)
		if tx.ProofType == 0 || tx.TxHash == "" {
			logger.Warn("unGenProof error:%v %v", tx.ProofType.String(), tx.TxHash)
			err := DeleteUnGenProof(b.store, Bitcoin, tx.TxHash)
			if err != nil {
				logger.Error("delete ungen proof error:%v %v", tx.TxHash, err)
			}
			continue
		}
		switch tx.ProofType {
		case common.DepositTxType:
			err := b.checkDepositRequest(tx)
			if err != nil {
				logger.Error("check deposit request error:%v %v", tx.TxHash, err)
				continue
			}
		case common.VerifyTxType:
			err := b.tryProofRequest(common.VerifyTxType, 0, 0, tx.TxHash)
			if err != nil {
				logger.Error("try proof request error:%v %v", tx.TxHash, err)
				continue
			}
		default:
			logger.Error("unknown proof type: %v", tx.ProofType.String())
		}
	}
	return nil
}

func (b *BitcoinAgent) checkDepositRequest(tx *DbUnGenProof) error {
	exists, err := CheckProof(b.fileStore, common.DepositTxType, 0, 0, tx.TxHash)
	if err != nil {
		logger.Error("check proof error:%v %v", tx.TxHash, err)
		return err
	}
	if exists {
		logger.Debug("%v %v proof exists ,delete ungen proof now", tx.ProofType.String(), tx.TxHash)
		err = DeleteUnGenProof(b.store, Bitcoin, tx.TxHash)
		if err != nil {
			logger.Error("delete ungen proof error:%v %v", tx.TxHash, err)
			return err
		}
		return nil
	}

	ok, confirms, err := b.CheckTxConfirms(tx.TxHash, tx.Amount)
	if err != nil {
		logger.Error("check tx confirms error: %v %v", tx.TxHash, err)
		return err
	}
	if !ok {
		logger.Warn("wait tx %v confirm: %v %v", tx.TxHash, tx.Amount, confirms)
		return nil
	}
	endHeight := tx.Height + uint64(confirms)
	if confirms <= 48 {
		exists, err := CheckProof(b.fileStore, common.BtcBulkType, tx.Height, endHeight, "")
		if err != nil {
			logger.Error(err.Error())
			return err
		}
		if !exists {
			err := b.tryProofRequest(common.BtcBulkType, tx.Height, endHeight, "")
			if err != nil {
				logger.Error("try proof request error:%v %v", tx.TxHash, err)
				return err
			}
			return nil
		}

	} else {
		exists, err := CheckProof(b.fileStore, common.BtcPackedType, tx.Height, endHeight, "")
		if err != nil {
			logger.Error(err.Error())
			return err
		}
		if !exists {
			err := b.tryProofRequest(common.BtcPackedType, tx.Height, endHeight, "")
			if err != nil {
				logger.Error(err.Error())
				return err
			}
			return nil
		}
	}
	wrapExists, err := CheckProof(b.fileStore, common.BtcWrapType, tx.Height, endHeight, "")
	if err != nil {
		logger.Error("check proof error:%v %v", tx.TxHash, err)
		return err
	}
	if !wrapExists {
		err := b.tryProofRequest(common.BtcWrapType, tx.Height, endHeight, "")
		if err != nil {
			logger.Error("try proof request error:%v %v", tx.TxHash, err)
			return err
		}
		return nil
	}
	err = b.tryProofRequest(common.DepositTxType, tx.Height, endHeight, tx.TxHash)
	if err != nil {
		logger.Error("try proof request error:%v %v", tx.TxHash, err)
		return err
	}
	return nil
}

func (b *BitcoinAgent) CheckTxConfirms(hash string, amount uint64) (bool, int, error) {
	needConfirms := 0
	if amount < 100000000 {
		needConfirms = 1
	} else if amount < 200000000 {
		needConfirms = 2
	} else {
		needConfirms = 3
	}
	tx, err := b.btcClient.GetTransaction(hash)
	if err != nil {
		logger.Error("get tx error:%v %v", hash, err)
		return false, 0, err
	}
	if tx.Confirmations >= needConfirms {
		return true, needConfirms, nil
	}
	return false, 0, nil
}

func (b *BitcoinAgent) tryProofRequest(proofType common.ZkProofType, index, end uint64, txHash string) error {
	proofId := common.NewProofId(proofType, index, end, txHash)
	exists := b.cache.Check(proofId)
	if exists {
		logger.Debug("proof request exists: %v", proofId)
		return nil
	}
	exists, err := CheckProof(b.fileStore, proofType, index, end, txHash)
	if err != nil {
		logger.Error("check proof error:%v %v", txHash, err)
		return err
	}
	if exists {
		return nil
	}
	data, ok, err := b.getRequestData(proofType, index, end, txHash)
	if err != nil {
		logger.Error("get request data error:%v %v", txHash, err)
		return err
	}
	if !ok {
		return nil
	}
	zkProofRequest := common.NewZkProofRequest(proofType, data, index, end, txHash)
	b.SendProofRequest(zkProofRequest)
	return nil
}

func (b *BitcoinAgent) getRequestData(proofType common.ZkProofType, index, end uint64, txHash string) (interface{}, bool, error) {
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

	case common.BtcBulkType:
		data, err := GetBtcMidBlockHeader(b.btcClient, index, end)
		if err != nil {
			logger.Error("get mid block header error:%v %v", txHash, err)
			return nil, false, err
		}
		return rpc.BtcBulkRequest{
			Data: data,
		}, true, nil
	case common.BtcPackedType:
		data, err := GetBtcMidBlockHeader(b.btcClient, index, end)
		if err != nil {
			logger.Error("get mid block header error:%v %v", txHash, err)
			return nil, false, err
		}
		return rpc.BtcPackedRequest{
			Data: data,
		}, true, nil
	case common.BtcWrapType:
		data, err := GetBtcWrapData(b.fileStore, b.btcClient, index, end)
		if err != nil {
			logger.Error("get btc wrap data error:%v %v", txHash, err)
			return nil, false, err
		}
		return data, true, nil

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
	proofData, err := grUtil.GetDefaultGrandRollupProofData(b.btcProverClient, txHash, tx.Blockhash)
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
	proofData, err := grUtil.GetDefaultGrandRollupProofData(b.btcProverClient, txHash, tx.Blockhash)
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
	return BitcoinAgentName
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
