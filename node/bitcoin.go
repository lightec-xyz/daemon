package node

import (
	"encoding/hex"
	"fmt"
	"strings"

	ethcommon "github.com/ethereum/go-ethereum/common"
	btcproverClient "github.com/lightec-xyz/btc_provers/utils/client"
	grUtil "github.com/lightec-xyz/btc_provers/utils/txinchain"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/rpc/bitcoin"
	"github.com/lightec-xyz/daemon/rpc/bitcoin/types"
	"github.com/lightec-xyz/daemon/rpc/ethereum"
	"github.com/lightec-xyz/daemon/store"
)

type BitcoinAgent struct {
	btcClient       *bitcoin.Client
	ethClient       *ethereum.Client
	btcProverClient *btcproverClient.Client
	store           store.IStore
	memoryStore     store.IStore
	cache           *Cache
	proofRequest    chan<- []*common.ZkProofRequest
	operatorAddr    string
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
		logger.Info(" btc agent success send btc proof request: %v", request.RequestId())
		b.cache.Store(request.RequestId(), nil)
	}
}

func (b *BitcoinAgent) updateRedeemInfo(height int64, txList []*Transaction) error {
	//todo
	return nil
}

func (b *BitcoinAgent) saveData(height int64, txes []*Transaction) error {
	err := WriteBitcoinTxIdsByHeight(b.store, height, txesToTxIds(txes))
	if err != nil {
		logger.Error("write bitcoin tx ids error: %v %v", height, err)
		return err
	}
	err = WriteTxes(b.store, txesToDbTxes(txes))
	if err != nil {
		logger.Error("put redeem tx error: %v %v", height, err)
		return err
	}
	err = WriteDbProof(b.store, txesToDbProofs(txes))
	if err != nil {
		logger.Error("write Proof error: %v", err)
		return err
	}
	err = WriteUnGenProof(b.store, common.BitcoinChain, txesToUnGenProofs(txes))
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
	logger.Info("bitcoinAgent receive Proof resp: %v %x", resp.RespId(), resp.Proof)
	b.cache.Delete(resp.RespId())
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
	return nil
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
		TxType:    common.DepositTx,
		ChainType: common.BitcoinChain,
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
		TxType:    common.RedeemTx,
		ChainType: common.BitcoinChain,
		ProofType: common.VerifyTxType,
		BlockHash: blockHash,
	}
}
