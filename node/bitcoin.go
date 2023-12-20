package node

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc/bitcoin"
	"github.com/lightec-xyz/daemon/rpc/bitcoin/types"
	"github.com/lightec-xyz/daemon/rpc/ethereum"
	"github.com/lightec-xyz/daemon/store"
	"math/big"
	"strings"
	"time"
)

type BitcoinAgent struct {
	btcClient            *bitcoin.Client
	ethClient            *ethereum.Client
	store                store.IStore
	memoryStore          store.IStore
	blockTime            time.Duration
	proofResponse        <-chan ProofResponse
	proofRequest         chan<- []ProofRequest
	nonceManager         *NonceManager
	checkProofHeightNums int64
	whiteList            map[string]bool // todo
	operatorAddr         string
	submitTxEthAddr      string
	keyStore             *KeyStore
	minDepositValue      float64
	initStartHeight      int64
	autoSubmit           bool
	exitSign             chan struct{}
}

func NewBitcoinAgent(cfg NodeConfig, store, memoryStore store.IStore, btcClient *bitcoin.Client, ethClient *ethereum.Client,
	request chan []ProofRequest, response <-chan ProofResponse, nonceManager *NonceManager, keyStore *KeyStore) (IAgent, error) {
	submitTxEthAddr, err := privateKeyToEthAddr(cfg.EthPrivateKey)
	if err != nil {
		logger.Error("privateKeyToEthAddr error:%v", err)
		return nil, err
	}
	return &BitcoinAgent{
		btcClient:            btcClient,
		ethClient:            ethClient,
		store:                store,
		memoryStore:          memoryStore,
		blockTime:            cfg.BtcScanBlockTime,
		operatorAddr:         cfg.BtcOperatorAddr,
		proofRequest:         request,
		proofResponse:        response,
		checkProofHeightNums: 100, // todo
		minDepositValue:      0,   // todo
		nonceManager:         nonceManager,
		keyStore:             keyStore,
		submitTxEthAddr:      submitTxEthAddr,
		exitSign:             make(chan struct{}, 1),
		initStartHeight:      cfg.BtcInitHeight,
		autoSubmit:           cfg.AutoSubmit,
	}, nil
}

func (b *BitcoinAgent) Init() error {
	logger.Info("bitcoin agent init now")
	exists, err := ReadInitBitcoinHeight(b.store)
	if err != nil {
		logger.Error("get btc current height error:%v", err)
		return err
	}
	if exists {
		logger.Debug("bitcoin agent check uncompleted generate proof tx")
		err := b.checkUnCompleteGenerateProofTx()
		if err != nil {
			logger.Error("check uncompleted generate proof tx error:%v", err)
			return err
		}
	} else {
		logger.Debug("init btc current height: %v", b.initStartHeight)
		err := WriteBitcoinHeight(b.store, b.initStartHeight)
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

// checkUnCompleteGenerateProofTx check uncompleted generate proof tx,resend again
func (b *BitcoinAgent) checkUnCompleteGenerateProofTx() error {
	return nil
	//currentHeight, err := b.getCurrentHeight()
	//if err != nil {
	//	logger.Error("get btc current height error:%v", err)
	//	return err
	//}
	//start := currentHeight - b.checkProofHeightNums
	//var proofList []ProofRequest
	//for index := start; index < currentHeight; index++ {
	//	var txIdList []string
	//	hasObj, err := b.store.HasObj(index)
	//	if err != nil {
	//		logger.Error("get txIdList error:%v", err)
	//		return err
	//	}
	//	if !hasObj {
	//		continue
	//	}
	//	err = b.store.GetObj(index, &txIdList)
	//	if err != nil {
	//		logger.Error("get txIdList error:%v", err)
	//		return err
	//	}
	//	for _, txId := range txIdList {
	//		var proof Proof
	//		err := b.store.GetObj(TxIdToProofId(txId), &proof)
	//		if err != nil {
	//			logger.Error("get proof error:%v", err)
	//			return err
	//		}
	//		//todo
	//		proofList = append(proofList, ProofRequest{
	//			TxHash:      proof.TxHash,
	//			ProofType: Deposit,
	//			Msg:       proof.Msg,
	//		})
	//	}
	//}
	//b.proofRequest <- proofList
	//return nil
}

func (b *BitcoinAgent) getCurrentHeight() (int64, error) {
	return ReadBitcoinHeight(b.store)

}

func (b *BitcoinAgent) ScanBlock() error {
	logger.Debug("bitcoin scan block ...")
	curHeight, err := b.getCurrentHeight()
	if err != nil {
		logger.Error("get btc current height error:%v", err)
		return err
	}
	if curHeight < b.initStartHeight {
		curHeight = b.initStartHeight
	}
	blockCount, err := b.btcClient.GetBlockCount()
	if err != nil {
		logger.Error("bitcoin client get block count error:%v", err)
		return err
	}
	//todo
	blockCount = blockCount - 0
	if curHeight >= blockCount {
		logger.Debug("btc current height:%d,node block count:%d", curHeight, blockCount)
		return nil
	}
	for index := curHeight + 1; index <= blockCount; index++ {
		logger.Debug("bitcoin parse block height:%d", index)
		depositTxes, redeemTxes, proofRequests, proofs, err := b.parseBlock(index)
		if err != nil {
			logger.Error("bitcoin agent parse block error: %v %v", index, err)
			return err
		}
		err = b.saveDataToDb(index, depositTxes, proofs)
		if err != nil {
			logger.Error("bitcoin agent save data to db error: %v %v", index, err)
			return err
		}
		b.proofRequest <- proofRequests

		if len(redeemTxes) > 0 {
			err := b.updateRedeemTxInfo(index, redeemTxes)
			if err != nil {
				logger.Error("update tx info error: %v %v", index, err)
				return err
			}
			err = b.updateContractUtxoChange(redeemTxes)
			if err != nil {
				logger.Error("update utxo error: %v %v", index, err)
				return err
			}
		}
	}
	return nil
}

func (b *BitcoinAgent) updateRedeemTxInfo(height int64, txList []*BitcoinTx) error {
	// todo
	err := UpdateRedeemInfo(b.store, txList)
	if err != nil {
		logger.Error("update redeem info error: %v %v", height, err)
		return err
	}
	return nil

}

func (b *BitcoinAgent) saveDataToDb(height int64, depositTxes []*BitcoinTx, proofs []Proof) error {
	//todo
	err := WriteBitcoinTx(b.store, height, depositTxes)
	if err != nil {
		logger.Error("write deposit tx error: %v %v", height, err)
		return err
	}
	err = WriteProof(b.store, proofs)
	if err != nil {
		logger.Error("write proof error: %v %v", height, err)
		return err
	}
	err = WriteBitcoinHeight(b.store, height)
	if err != nil {
		logger.Error("write btc height error: %v %v", height, err)
		return err
	}
	return nil
}

func (b *BitcoinAgent) parseBlock(height int64) ([]*BitcoinTx, []*BitcoinTx, []ProofRequest, []Proof, error) {
	blockHash, err := b.btcClient.GetBlockHash(height)
	if err != nil {
		logger.Error("btcClient get block hash error: %v %v", height, err)
		return nil, nil, nil, nil, err
	}
	blockWithTx, err := b.btcClient.GetBlock(blockHash)
	if err != nil {
		logger.Error("btcClient get block error: %v %v", blockHash, err)
		return nil, nil, nil, nil, err
	}
	var requests []ProofRequest
	var depositTxes []*BitcoinTx
	var redeemTxes []*BitcoinTx
	var proofs []Proof
	for _, tx := range blockWithTx.Tx {
		redeemTx, isRedeem := b.isRedeemTx(tx)
		if isRedeem {
			logger.Info("find bitcoin redeem tx: %v", tx.Txid)
			redeemTxes = append(redeemTxes, redeemTx)
			continue
		}
		depositTx, isDeposit, err := b.isDepositTx(tx)
		if err != nil {
			logger.Error("check deposit tx error: %v %v", tx.Txid, err)
			return nil, nil, nil, nil, err
		}
		if isDeposit {
			logger.Info("find bitcoin deposit tx: %v", tx.Txid)
			proofs = append(proofs, NewDepositTxProof(tx.Txid))
			requests = append(requests, NewDepositProofRequest(depositTx.TxId, depositTx.EthAddr, depositTx.Amount, depositTx.Utxos))
			depositTxes = append(depositTxes, depositTx)
		}
	}
	return depositTxes, redeemTxes, requests, proofs, nil
}

func (b *BitcoinAgent) Transfer() {
	// todo queue ?
	logger.Debug("start bitcoin transfer goroutine")
	for {
		select {
		case <-b.exitSign:
			logger.Info("bitcoin transfer goroutine exit ...")
			return
		case response := <-b.proofResponse:
			logger.Info("bitcoinAgent receive deposit proof response: %v", response.String())
			err := b.updateProof(response)
			if err != nil {
				logger.Error("update proof error: %v %v", response.TxId, err)
				continue
			}
			exists, err := b.ethClient.CheckDepositProof(response.TxId)
			if err != nil {
				// todo retry  add queue?
				logger.Error("check deposit proof error: %v %v", response.TxId, err)
				continue
			}
			if exists {
				logger.Warn("deposit utxo already exists: %v", response.TxId)
				continue
			}
			if b.autoSubmit {
				if ProofStatus(response.Status) == ProofSuccess {
					txHash, err := b.MintZKBtcTx(response)
					if err != nil {
						//todo add queue or cli retry ?
						logger.Error("mint btc tx error:%v", err)
						continue
					}
					err = b.updateDestChainHash(response.TxId, txHash)
					if err != nil {
						logger.Error("update deposit info error: %v %v", response.TxId, err)
						continue
					}
				}

			}
		}

	}
}

func (b *BitcoinAgent) updateDestChainHash(txId, ethTxHash string) error {
	err := WriteDestChainHash(b.store, txId, ethTxHash)
	if err != nil {
		logger.Error("write dest hash error: %v %v", txId, err)
		return err
	}
	return nil

}

func (b *BitcoinAgent) updateContractUtxoChange(utxoList []*BitcoinTx) error {
	// todo
	var txIds []string
	for _, tx := range utxoList {
		txIds = append(txIds, tx.TxId)
	}
	nonce, err := b.nonceManager.GetNonce(b.submitTxEthAddr, b.ethClient, b.store)
	if err != nil {
		logger.Error("get  nonce error:%v", err)
		return err
	}
	chainId, err := b.ethClient.GetChainId()
	if err != nil {
		logger.Error("get chain id error:%v", err)
		return err
	}
	gasPrice, err := b.ethClient.GetGasPrice()
	if err != nil {
		logger.Error("get gas price error:%v", err)
		return err
	}
	gasLimit := uint64(500000)
	proofBytes := []byte("test ok")
	txHash, err := b.ethClient.UpdateUtxoChange(b.keyStore.GetPrivateKey(), txIds, nonce, gasLimit, chainId, gasPrice, proofBytes)
	if err != nil {
		logger.Error("update utxo change error:%v", err)
		return err
	}
	logger.Info("success send update utxo change  hash:%v", txHash)
	return nil
}

func (b *BitcoinAgent) MintZKBtcTx(resp ProofResponse) (string, error) {
	//todo
	nonce, err := b.nonceManager.GetNonce(b.submitTxEthAddr, b.ethClient, b.store)
	if err != nil {
		logger.Error("get nonce error:%v", err)
		return "", err
	}
	chainId, err := b.ethClient.GetChainId()
	if err != nil {
		logger.Error("get chain id error:%v %v", b.submitTxEthAddr, err)
		return "", err
	}
	gasPrice, err := b.ethClient.GetGasPrice()
	if err != nil {
		logger.Error("get gas price error:%v", err)
		return "", err
	}
	//todo
	gasLimit := uint64(500000)
	amountBig := big.NewInt(resp.Amount)
	proofBytes := []byte(resp.Proof)
	index := resp.Utxos[0].Index
	txHash, err := b.ethClient.Deposit(b.keyStore.GetPrivateKey(), resp.TxId, index, nonce, gasLimit, chainId, gasPrice,
		amountBig, proofBytes)
	if err != nil {
		logger.Error("mint btc tx error:%v", err)
		return "", err
	}
	logger.Info("success send mint zkbtctx hash:%v, amount: %v", txHash, amountBig.String())
	return txHash, nil
}

func (b *BitcoinAgent) isRedeemTx(tx types.Tx) (*BitcoinTx, bool) {
	// todo more check
	for _, vin := range tx.Vin {
		if vin.Prevout.ScriptPubKey.Address == b.operatorAddr {
			bitcoinTx := &BitcoinTx{
				TxId:   tx.Txid,
				TxType: BtcRedeem,
			}
			return bitcoinTx, true
		}
	}
	return nil, false
}

func (b *BitcoinAgent) isDepositTx(tx types.Tx) (*BitcoinTx, bool, error) {
	// todo more rule
	txOuts := tx.Vout
	if len(txOuts) < 2 {
		return nil, false, nil
	}
	if txOuts[1].ScriptPubKey.Address != b.operatorAddr {
		return nil, false, nil
	}
	if txOuts[1].Value <= b.minDepositValue {
		logger.Warn("deposit tx less than min value: %v %v", b.minDepositValue, tx.Txid)
		return nil, false, nil
	}
	if !(txOuts[0].ScriptPubKey.Type == "nulldata" && strings.HasPrefix(txOuts[0].ScriptPubKey.Hex, "6a")) {
		logger.Warn("find deposit tx but check rule fail: %v", tx.Txid)
		return nil, false, nil
	}
	ethAddr, err := getEthAddrFromScript(txOuts[0].ScriptPubKey.Hex)
	if err != nil {
		logger.Error("get eth addr from script error:%v %v", txOuts[0].ScriptPubKey.Hex, err)
		return nil, false, err
	}
	depositTx := &BitcoinTx{}
	utxoList := []Utxo{
		{
			TxId:  tx.Txid,
			Index: 1,
		},
	}
	depositTx.Utxos = utxoList
	depositTx.TxId = tx.Txid
	depositTx.EthAddr = ethAddr
	depositTx.TxType = BtcDeposit
	depositTx.Amount = BtcToSat(txOuts[1].Value)
	return depositTx, true, nil
}

func (b *BitcoinAgent) updateProof(resp ProofResponse) error {
	err := UpdateProof(b.store, resp.TxId, resp.Proof, Deposit, ProofStatus(resp.Status))
	if err != nil {
		logger.Error("update proof error: %v %v", resp.TxId, err)
		return err
	}
	return nil

}

func (b *BitcoinAgent) Close() error {
	close(b.exitSign)
	return nil
}
func (b *BitcoinAgent) Name() string {
	return "Bitcoin Agent"
}

func (b *BitcoinAgent) BlockTime() time.Duration {
	return b.blockTime
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
	isHexAddress := common.IsHexAddress(script[4:])
	if !isHexAddress {
		return "", fmt.Errorf("script is not hex address")
	}
	return script[4:], nil
}

func NewDepositProofRequest(txId, ethAddr string, amount int64, utxo []Utxo) ProofRequest {
	return ProofRequest{
		TxId:    txId,
		EthAddr: ethAddr,
		Amount:  amount,
		Utxos:   utxo,
	}
}

func NewDepositTxProof(txId string) Proof {
	return Proof{
		TxId:      txId,
		ProofType: Deposit,
		Status:    ProofDefault,
	}
}
