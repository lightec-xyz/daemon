package node

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/mempool"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/lightec-xyz/daemon/rpc/bitcoin"

	"github.com/lightec-xyz/daemon/rpc/dfinity"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc/ethereum"
	"github.com/lightec-xyz/daemon/store"
)

type bitcoinAgent struct {
	btcClient       *bitcoin.Client
	ethClient       *ethereum.Client
	dfinityClient   *dfinity.Client
	proverClient    *BtcClient
	btcFilter       *BtcFilter
	initHeight      uint64
	curHeight       uint64
	txManager       *TxManager
	chainStore      *ChainStore
	fileStore       *FileStorage
	mode            Mode
	chainForkSignal chan<- *ChainFork
	reScan          bool
	check           bool
}

func NewBitcoinAgent(cfg Config, store store.IStore, btcProverClient *BtcClient, btcClient *bitcoin.Client, ethClient *ethereum.Client,
	dfinityClient *dfinity.Client, txManager *TxManager, chainFork chan *ChainFork, fileStore *FileStorage) (IAgent, error) {
	return &bitcoinAgent{
		btcClient:       btcClient,
		ethClient:       ethClient,
		dfinityClient:   dfinityClient,
		btcFilter:       cfg.BtcFilter,
		initHeight:      cfg.BtcInitHeight,
		txManager:       txManager,
		chainStore:      NewChainStore(store),
		chainForkSignal: chainFork,
		reScan:          cfg.BtcReScan,
		fileStore:       fileStore,
		check:           true,
		proverClient:    btcProverClient,
		mode:            cfg.Mode,
	}, nil
}

func (b *bitcoinAgent) Init() error {
	logger.Info("start init bitcoin agent")
	_, err := b.btcClient.GetBlockCount()
	if err != nil {
		logger.Error(" btcClient get block count error:%v", err)
		return err
	}
	height, exists, err := b.chainStore.ReadBtcHeight()
	if err != nil {
		logger.Error("get btc current height error:%v", err)
		return err
	}

	if !exists && b.initHeight-24 > 0 {
		b.initHeight = b.initHeight - 24
	}
	if !exists || height < b.initHeight || b.reScan {
		logger.Info("use btc init height: %v", b.initHeight)
		err := b.chainStore.WriteBtcHeight(b.initHeight)
		if err != nil {
			logger.Error("put init btc current height error:%v", err)
			return err
		}
	}
	if exists && height-BtcClientCacheHeight > b.initHeight {
		b.proverClient.SetInitHeight(int64(height - BtcClientCacheHeight))
	}

	////todo
	//for index := b.initHeight; index < height; index++ {
	//	logger.Debug("set btc client cache: %v", index)
	//	err := b.setBtcClientCache(index)
	//	if err != nil {
	//		logger.Error("set btc client cache error: %v %v", index, err)
	//		return err
	//	}
	//}
	return nil
}

func (b *bitcoinAgent) ScanBlock() error {
	logger.Debug("bitcoin scan block ...")
	currentHeight, ok, err := b.chainStore.ReadBtcHeight()
	if err != nil {
		logger.Error("get btc current height error:%v", err)
		return err
	}
	if !ok {
		logger.Warn("no find btc current height")
		return fmt.Errorf("no btc current height")
	}
	latestHeight, err := b.btcClient.GetBlockCount()
	if err != nil {
		logger.Error("get block count error:%v", err)
		return err
	}
	blockCount := uint64(latestHeight)
	if currentHeight >= blockCount {
		logger.Debug("btc current height:%d,node block count:%d", currentHeight, blockCount)
		return nil
	}

	for index := currentHeight + 1; index <= blockCount; index++ {
		preHeight := index - 1
		chainFork, err := b.checkChainFork(preHeight)
		if err != nil {
			logger.Error("check chain fork error: %v %v", index, err)
			return err
		}
		if chainFork {
			err := b.rollback(preHeight)
			if err != nil {
				logger.Error("rollback error: %v %v", index, err)
				return err
			}
			return nil
		}
		err = b.scan(index)
		if err != nil {
			logger.Error("scan error: %v %v", index, err)
			return err
		}
		err = b.chainStore.WriteBtcHeight(index)
		if err != nil {
			logger.Error("write btc height error: %v %v", index, err)
			return err
		}
		err = b.setBtcClientCache(index)
		if err != nil {
			logger.Error("set cache btc client error: %v %v", index, err)
		}

	}
	//err = b.cropData(currentHeight)
	//if err != nil {
	//	logger.Error("crop data error: %v %v", currentHeight, err)
	//	return err
	//}

	return nil
}

func (b *bitcoinAgent) setBtcClientCache(height uint64) error {
	cStartHeight := height - BtcClientCacheHeight
	if cStartHeight > b.initHeight {
		b.proverClient.SetInitHeight(int64(cStartHeight))
		err := b.chainStore.DelBtcClientCache(cStartHeight)
		if err != nil {
			logger.Error("delete btc client cache error: %v %v", cStartHeight, err)
			return err
		}
	}
	return nil
}

func (b *bitcoinAgent) cropData(height uint64) error {
	// todo
	if !(b.mode == LiteMode && height%10 == 0 && height-b.initHeight > BtcLiteCacheHeight) {
		return nil
	}
	eHeight := height - BtcLiteCacheHeight
	sHeight := eHeight - BtcLiteCacheHeight
	if sHeight < b.initHeight {
		sHeight = b.initHeight
	}
	for index := eHeight; index >= sHeight; index-- {
		err := b.chainStore.BtcDeleteData(index)
		if err != nil {
			logger.Warn("delete btc data error: %v %v", index, err)
			//return err
		}
	}
	return nil

}

func (b *bitcoinAgent) ReScan(height uint64) error {
	logger.Debug("bitcoin rescan block height:%d", height)
	err := b.scan(height)
	if err != nil {
		logger.Error("scan error: %v %v", height, err)
		return err
	}
	txIds, err := b.chainStore.ReadBtcTxHeight(height)
	if err != nil {
		logger.Error("get btc tx height error: %v %v", height, err)
		return err
	}
	for _, txId := range txIds {
		logger.Debug("delete proof: %v", txId)
		_ = b.fileStore.DelProof(NewHashStoreKey(common.BtcDepositType, DbValue(txId)))
		_ = b.fileStore.DelProof(NewHashStoreKey(common.BtcChangeType, DbValue(txId)))
	}
	return nil
}

func (b *bitcoinAgent) scan(index uint64) error {
	logger.Debug("bitcoin parse block height:%d", index)
	blockHash, err := b.btcClient.GetBlockHash(int64(index))
	if err != nil {
		logger.Error("btcClient get block hash error: %v %v", index, err)
		return err
	}
	blockHeader, err := b.btcClient.GetHexBlockHeader(blockHash)
	if err != nil {
		logger.Error("btcClient get block header error: %v %v", index, err)
		return err
	}
	err = b.chainStore.WriteBitcoinHash(index, blockHash)
	if err != nil {
		logger.Error("write bitcoin hash error: %v %v", index, err)
		return err
	}
	err = b.chainStore.WriteBlockHeader(blockHash, blockHeader)
	if err != nil {
		logger.Error("write block header error: %v %v", index, err)
		return err
	}
	depositTxes, redeemTxes, err := b.parseBlock(blockHash, index)
	if err != nil {
		logger.Error("bitcoin agent parse block error: %v %v", index, err)
		return err
	}
	err = b.chainStore.BtcSaveData(index, depositTxes, redeemTxes)
	if err != nil {
		logger.Error("")
		return err
	}
	if b.reScan {
		checkTxes := append(depositTxes, redeemTxes...)
		for _, tx := range checkTxes {
			if !tx.Proved {
				logger.Debug("delete proof: %v %v", tx.ProofType.Name(), tx.Hash)
				_ = b.fileStore.DelProof(NewHashStoreKey(common.BtcDepositType, DbValue(tx.Hash)))
				_ = b.fileStore.DelProof(NewHashStoreKey(common.BtcChangeType, DbValue(tx.Hash)))
			}
		}
	}
	return nil

}

func (b *bitcoinAgent) rollback(height uint64) error {
	startForkHeight, err := b.findForkHeight(height)
	if err != nil {
		logger.Error("find fork height error:%v", err)
		return err
	}
	logger.Warn("bitcoin found start fork height: %v", startForkHeight)
	for index := height; index >= startForkHeight; index-- {
		logger.Debug("bitcoin start rollback height: %v", index)
		err := b.chainStore.BtcDeleteData(index)
		if err != nil {
			return err
		}
		err = b.chainStore.WriteBtcHeight(index - 1)
		if err != nil {
			logger.Error("write btc height error: %v %v", index, err)
			return err
		}
	}
	chainFork := ChainFork{
		ForkHeight: startForkHeight,
		Chain:      common.BitcoinChain,
		Timestamp:  time.Now().UnixNano(),
	}
	b.chainForkSignal <- &chainFork
	return nil
}

func (b *bitcoinAgent) parseBlock(hash string, height uint64) ([]*DbTx, []*DbTx, error) {
	blockV0Hex, err := b.btcClient.GetBlockStr(hash)
	if err != nil {
		logger.Error("btcClient get block error: %v %v", hash, err)
		return nil, nil, err
	}
	blockBytes, err := hex.DecodeString(string(blockV0Hex))
	if err != nil {
		logger.Error("btcClient hex-decode block data error: %v %v", hash, err)
		return nil, nil, err
	}
	rawBlock, err := btcutil.NewBlockFromBytes(blockBytes)
	if err != nil {
		logger.Error("btcClient deserialize block data error: %v %v", hash, err)
		return nil, nil, err
	}

	var params chaincfg.Params
	if height >= 500000 { // only a hack TODO FIXME
		params = chaincfg.MainNetParams
	} else {
		params = chaincfg.TestNet4Params
	}

	// begin codes adapted from btcd handleGetBlock, including the getDifficultyRatio function, and createTxRawResult and its depencencies
	blockHeader := &rawBlock.MsgBlock().Header
	block := btcjson.GetBlockVerboseResult{
		Hash:          hash,
		Version:       blockHeader.Version,
		VersionHex:    fmt.Sprintf("%08x", blockHeader.Version),
		MerkleRoot:    blockHeader.MerkleRoot.String(),
		PreviousHash:  blockHeader.PrevBlock.String(),
		Nonce:         blockHeader.Nonce,
		Time:          blockHeader.Timestamp.Unix(),
		Confirmations: int64(1), // we don't know yet and not important, so ...
		Height:        int64(height),
		Size:          int32(len(blockBytes)),
		StrippedSize:  int32(rawBlock.MsgBlock().SerializeSizeStripped()),
		Weight:        int32(blockchain.GetBlockWeight(rawBlock)),
		Bits:          strconv.FormatInt(int64(blockHeader.Bits), 16),
		Difficulty:    getDifficultyRatio(blockHeader.Bits, &params),
		NextHash:      "", // we don't care
	}
	txns := rawBlock.Transactions()
	rawTxns := make([]btcjson.TxRawResult, len(txns))
	for i, tx := range txns {
		rawTxn, err := createTxRawResult(&params, tx.MsgTx(),
			tx.Hash().String(), blockHeader, hash,
			int32(height), int32(height)+1) // we don't care the chain height for now
		if err != nil {
			return nil, nil, err
		}
		rawTxns[i] = *rawTxn
	}
	block.RawTx = rawTxns
	// end codes adapted from btcd handleGetBlock

	blockVerboseStr, err := json.Marshal(&block)
	if err != nil {
		logger.Error("marshal btc block error: %v %v", hash, err)
		return nil, nil, err
	}
	err = b.chainStore.WriteBtcBlock(hash, string(blockVerboseStr))
	if err != nil {
		logger.Error("write btc block error: %v %v", hash, err)
		return nil, nil, err
	}

	var depositTxes []*DbTx
	var redeemTxes []*DbTx
	for txIndex, tx := range block.RawTx {

		migrateTx, isMigrate, err := b.migrateTx(tx, height, uint64(txIndex), uint64(block.Time))
		if err != nil {
			logger.Error("check btc migrate tx error: %v %v", tx.Txid, err)
			return nil, nil, err
		}
		if isMigrate {
			depositTxes = append(depositTxes, migrateTx)
			continue
		}
		redeemTx, isRedeem, err := b.redeemTx(tx, height, uint64(txIndex), uint64(block.Time))
		if err != nil {
			logger.Error("check btc Redeem tx error: %v %v", tx.Txid, err)
			return nil, nil, err
		}
		if isRedeem {
			redeemTxes = append(redeemTxes, redeemTx)
			continue
		}
		depositTx, isDeposit, err := b.depositTx(tx, height, uint64(txIndex), uint64(block.Time))
		if err != nil {
			logger.Error("check deposit tx error: %v %v", tx.Txid, err)
			return nil, nil, err
		}
		if isDeposit {
			depositTxes = append(depositTxes, depositTx)
		}
	}
	return depositTxes, redeemTxes, nil
}

// function getDifficultyRatio is adapted from btcd
// getDifficultyRatio returns the proof-of-work difficulty as a multiple of the
// minimum difficulty using the passed bits field from the header of a block.
func getDifficultyRatio(bits uint32, params *chaincfg.Params) float64 {
	// The minimum difficulty is the max possible proof-of-work limit bits
	// converted back to a number.  Note this is not the same as the proof of
	// work limit directly because the block difficulty is encoded in a block
	// with the compact form which loses precision.
	max := blockchain.CompactToBig(params.PowLimitBits)
	target := blockchain.CompactToBig(bits)

	difficulty := new(big.Rat).SetFrac(max, target)
	outString := difficulty.FloatString(8)
	diff, err := strconv.ParseFloat(outString, 64)
	if err != nil {
		logger.Error("Cannot get difficulty: %v", err)
		return 0
	}
	return diff
}

// function createTxRawResult is adapted from btcd
// createTxRawResult converts the passed transaction and associated parameters
// to a raw transaction JSON object.
func createTxRawResult(chainParams *chaincfg.Params, mtx *wire.MsgTx,
	txHash string, blkHeader *wire.BlockHeader, blkHash string,
	blkHeight int32, chainHeight int32) (*btcjson.TxRawResult, error) {

	mtxHex, err := messageToHex(mtx)
	if err != nil {
		return nil, err
	}

	txReply := &btcjson.TxRawResult{
		Hex:      mtxHex,
		Txid:     txHash,
		Hash:     mtx.WitnessHash().String(),
		Size:     int32(mtx.SerializeSize()),
		Vsize:    int32(mempool.GetTxVirtualSize(btcutil.NewTx(mtx))),
		Weight:   int32(blockchain.GetTransactionWeight(btcutil.NewTx(mtx))),
		Vin:      createVinList(mtx),
		Vout:     createVoutList(mtx, chainParams, nil),
		Version:  uint32(mtx.Version),
		LockTime: mtx.LockTime,
	}

	if blkHeader != nil {
		// This is not a typo, they are identical in bitcoind as well.
		txReply.Time = blkHeader.Timestamp.Unix()
		txReply.Blocktime = blkHeader.Timestamp.Unix()
		txReply.BlockHash = blkHash
		txReply.Confirmations = uint64(1 + chainHeight - blkHeight)
	}

	return txReply, nil
}

// from btcd
func messageToHex(msg wire.Message) (string, error) {
	var buf bytes.Buffer
	if err := msg.BtcEncode(&buf, 70002, wire.WitnessEncoding); err != nil {
		context := fmt.Sprintf("Failed to encode msg of type %T", msg)
		return "", internalRPCError(err.Error(), context)
	}

	return hex.EncodeToString(buf.Bytes()), nil
}

// from btcd
func internalRPCError(errStr, context string) *btcjson.RPCError {
	return btcjson.NewRPCError(btcjson.ErrRPCInternal.Code, errStr)
}

// from btcd
// createVinList returns a slice of JSON objects for the inputs of the passed
// transaction.
func createVinList(mtx *wire.MsgTx) []btcjson.Vin {
	// Coinbase transactions only have a single txin by definition.
	vinList := make([]btcjson.Vin, len(mtx.TxIn))
	if blockchain.IsCoinBaseTx(mtx) {
		txIn := mtx.TxIn[0]
		vinList[0].Coinbase = hex.EncodeToString(txIn.SignatureScript)
		vinList[0].Sequence = txIn.Sequence
		vinList[0].Witness = txIn.Witness.ToHexStrings()
		return vinList
	}

	for i, txIn := range mtx.TxIn {
		// The disassembled string will contain [error] inline
		// if the script doesn't fully parse, so ignore the
		// error here.
		disbuf, _ := txscript.DisasmString(txIn.SignatureScript)

		vinEntry := &vinList[i]
		vinEntry.Txid = txIn.PreviousOutPoint.Hash.String()
		vinEntry.Vout = txIn.PreviousOutPoint.Index
		vinEntry.Sequence = txIn.Sequence
		vinEntry.ScriptSig = &btcjson.ScriptSig{
			Asm: disbuf,
			Hex: hex.EncodeToString(txIn.SignatureScript),
		}

		if mtx.HasWitness() {
			vinEntry.Witness = txIn.Witness.ToHexStrings()
		}
	}

	return vinList
}

// from btcd
// createVoutList returns a slice of JSON objects for the outputs of the passed
// transaction.
func createVoutList(mtx *wire.MsgTx, chainParams *chaincfg.Params, filterAddrMap map[string]struct{}) []btcjson.Vout {
	voutList := make([]btcjson.Vout, 0, len(mtx.TxOut))
	for i, v := range mtx.TxOut {
		// The disassembled string will contain [error] inline if the
		// script doesn't fully parse, so ignore the error here.
		disbuf, _ := txscript.DisasmString(v.PkScript)

		// Ignore the error here since an error means the script
		// couldn't parse and there is no additional information about
		// it anyways.
		scriptClass, addrs, reqSigs, _ := txscript.ExtractPkScriptAddrs(
			v.PkScript, chainParams)

		// Encode the addresses while checking if the address passes the
		// filter when needed.
		passesFilter := len(filterAddrMap) == 0
		encodedAddrs := make([]string, len(addrs))
		for j, addr := range addrs {
			encodedAddr := addr.EncodeAddress()
			encodedAddrs[j] = encodedAddr

			// No need to check the map again if the filter already
			// passes.
			if passesFilter {
				continue
			}
			if _, exists := filterAddrMap[encodedAddr]; exists {
				passesFilter = true
			}
		}

		if !passesFilter {
			continue
		}

		var vout btcjson.Vout
		vout.N = uint32(i)
		vout.Value = btcutil.Amount(v.Value).ToBTC()
		vout.ScriptPubKey.Addresses = encodedAddrs
		vout.ScriptPubKey.Asm = disbuf
		vout.ScriptPubKey.Hex = hex.EncodeToString(v.PkScript)
		vout.ScriptPubKey.Type = scriptClass.String()
		vout.ScriptPubKey.ReqSigs = int32(reqSigs)

		// Address is defined when there's a single well-defined
		// receiver address. To spend the output a signature for this,
		// and only this, address is required.
		if len(encodedAddrs) == 1 && reqSigs <= 1 {
			vout.ScriptPubKey.Address = encodedAddrs[0]
		}

		voutList = append(voutList, vout)
	}

	return voutList
}

func (b *bitcoinAgent) checkChainFork(height uint64) (bool, error) {
	if height <= b.initHeight {
		return false, nil
	}
	preHash, exists, err := b.chainStore.ReadBitcoinHash(height)
	if err != nil {
		logger.Error("get btc hash error: %v %v", height, err)
		return false, err
	}
	if !exists {
		logger.Warn("local btc hash not exist: %v", height)
		return false, nil
	}
	blockHash, err := b.btcClient.GetBlockHash(int64(height))
	if err != nil {
		logger.Error("btcClient get block hash error: %v %v", height, err)
		return false, err
	}
	equal := common.StrEqual(blockHash, preHash)
	if !equal {
		logger.Error("bitcoin chain forked: found height %v,chainHash %v,localHash %v", height, blockHash, preHash)
		return true, nil
	}
	return false, nil
}

func (b *bitcoinAgent) findForkHeight(height uint64) (uint64, error) {
	for index := height; index >= b.initHeight; index = index - 1 {
		localBlockHash, exists, err := b.chainStore.ReadBitcoinHash(index)
		if err != nil {
			logger.Error("get btc hash error: %v %v", index, err)
			return 0, err
		}
		if !exists {
			logger.Error("btc hash not exist: %v", index)
			return 0, fmt.Errorf("btc hash not exist: %v", index)
		}
		chainBlockHash, err := b.btcClient.GetBlockHash(int64(index))
		if err != nil {
			logger.Error("btcClient get block hash error: %v %v", index, err)
			return 0, err
		}
		if common.StrEqual(localBlockHash, chainBlockHash) {
			logger.Info("find rollback start height: %v", index)
			return index + 1, nil
		}
	}
	return b.initHeight, nil
}

func (b *bitcoinAgent) checkTxProved(proofType common.ProofType, hash string) (bool, error) {
	switch proofType {
	case common.BtcChangeType:
		_, exists, err := b.chainStore.ReadUpdateUtxoDest(hash)
		if err != nil {
			logger.Error("check utxo error: %v %v", hash, err)
			return false, err
		}
		if exists {
			return true, nil
		}
		utxo, err := b.ethClient.GetUtxo(hash)
		if err != nil {
			logger.Error("check utxo error: %v %v", hash, err)
			return false, nil
		}
		return utxo.IsChangeConfirmed, nil
	case common.BtcDepositType:
		_, exists, err := b.chainStore.GetDestHash(hash)
		if err != nil {
			logger.Error("check deposit tx utxo error: %v %v", hash, err)
			return false, err
		}
		if exists {
			return true, nil
		}
		utxo, err := b.ethClient.GetUtxo(hash)
		if err != nil {
			logger.Warn("check deposit tx utxo error: %v %v", hash, err)
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
func (b *bitcoinAgent) ProofResponse(resp *common.ProofResponse) error {
	switch resp.ProofType {
	case common.BtcDepositType:
		logger.Info("find deposit proof: %v %v %v", resp.ProofId(), resp.Hash, hex.EncodeToString(resp.Proof))
		err := b.chainStore.UpdateProof(resp.Hash, hex.EncodeToString(resp.Proof), common.BtcDepositType, common.ProofSuccess)
		if err != nil {
			logger.Error("update Proof error: %v %v", resp.Hash, err)
			return err
		}
		submitTx := NewDbUnSubmitTx(resp.Hash, hex.EncodeToString(resp.Proof), common.BtcDepositType)
		hash, err := b.txManager.DepositBtc(submitTx)
		if err != nil {
			logger.Error("update deposit error: %v %v,save to db", resp.Hash, err)
			b.txManager.AddTask(resp)
			return err
		}
		logger.Debug("success  deposit: btcTxId:%v  ethHash:%v", resp.Hash, hash)

	case common.BtcUpdateCpType:
		logger.Info("find deposit proof: %v %v %v", resp.ProofId(), resp.Hash, hex.EncodeToString(resp.Proof))
		submitTx := NewDbUnSubmitTx(resp.Hash, hex.EncodeToString(resp.Proof), common.BtcUpdateCpType)
		hash, err := b.txManager.DepositBtc(submitTx)
		if err != nil {
			logger.Error("update deposit error: %v %v,save to db", resp.Hash, err)
			b.txManager.AddTask(resp)
			return err
		}
		logger.Debug("success  deposit: btcTxId:%v  ethHash:%v", resp.Hash, hash)
	case common.BtcChangeType:
		logger.Info("find change proof: %v %v", resp.Hash, hex.EncodeToString(resp.Proof))
		submitTx := NewDbUnSubmitTx(resp.Hash, hex.EncodeToString(resp.Proof), common.BtcChangeType)
		hash, err := b.txManager.UpdateUtxoChange(submitTx)
		if err != nil {
			logger.Error("update utxo fail: %v %v,save to db", resp.Hash, err)
			b.txManager.AddTask(resp)
			return err
		}
		logger.Debug("success update utxo: btcTxId:%v ,ethHash:%v", resp.Hash, hash)

	default:
	}
	return nil
}

func (b *bitcoinAgent) redeemTx(tx btcjson.TxRawResult, height, txIndex, blockTime uint64) (*DbTx, bool, error) {
	isRedeem := b.btcFilter.Redeem(tx.Vin)
	if !isRedeem {
		return nil, false, nil
	}
	redeemAmount := getRedeemAmount(tx.Vout, b.btcFilter.OperatorAddr)
	proved, err := b.checkTxProved(common.BtcChangeType, tx.Txid)
	if err != nil {
		logger.Error("check btc change proved error: %v,%v", tx.Txid, err)
		return nil, false, err
	}
	logger.Info("bitcoin agent find Redeem tx height: %v,hash: %v,amount: %.8f,proved:%v", height, tx.Txid, redeemAmount, proved)
	redeemBtcTx := NewRedeemBtcTx(height, txIndex, blockTime, tx.Txid, BtcToSat(redeemAmount), proved)
	return redeemBtcTx, true, nil
}

func (b *bitcoinAgent) migrateTx(tx btcjson.TxRawResult, height, txIndex, blockTime uint64, skipCheck ...bool) (*DbTx, bool, error) {
	txOuts := tx.Vout
	if len(txOuts) < 2 {
		return nil, false, nil
	}
	isMigrate := b.btcFilter.Migrate(tx.Vout)
	if !isMigrate {
		return nil, false, nil
	}
	amount := getDepositAmount(tx.Vout, b.btcFilter.OperatorAddr)
	minDepositValue := b.btcFilter.GetMinDepositValue()
	if amount < minDepositValue {
		logger.Debug("migrate amount tool low %v ,less than minDepositValue %v", amount, minDepositValue)
		return nil, false, nil
	}

	proved, err := b.checkTxProved(common.BtcDepositType, tx.Txid)
	if err != nil {
		logger.Error("bitcoin check chain proof error: %v %v", tx.Txid, err)
		return nil, false, err
	}
	logger.Info("bitcoin agent find  migrate tx height: %v, hash: %v, amount:%v ,proved:%v",
		height, tx.Txid, BtcToSat(amount), proved)
	depositTx := NewMigrateBtcTx(height, txIndex, blockTime, tx.Txid, BtcToSat(amount), proved)
	return depositTx, true, nil

}

func (b *bitcoinAgent) depositTx(tx btcjson.TxRawResult, height, txIndex, blockTime uint64) (*DbTx, bool, error) {
	txOuts := tx.Vout
	if len(txOuts) < 2 {
		return nil, false, nil
	}
	isDeposit := b.btcFilter.Deposit(tx.Vout)
	if !isDeposit {
		return nil, false, nil
	}
	ethAddr, err := isOpZkpProto(tx.Vout)
	if err != nil {
		logger.Error("get deposit info error: %v %v", tx.Txid, err)
		return nil, false, nil
	}
	amount := getDepositAmount(tx.Vout, b.btcFilter.OperatorAddr)
	minDepositValue := b.btcFilter.GetMinDepositValue()
	if amount < minDepositValue {
		logger.Debug("deposit amount tool low %f ,less than minDepositValue %f", amount, minDepositValue)
		return nil, false, nil
	}
	proved, err := b.checkTxProved(common.BtcDepositType, tx.Txid)
	if err != nil {
		logger.Error("bitcoin check chain proof error: %v %v", tx.Txid, err)
		return nil, false, err
	}
	logger.Info("bitcoin agent find  deposit tx height: %v, hash: %v, ethAddr:%v, amount:%v ,proved:%v",
		height, tx.Txid, ethAddr, BtcToSat(amount), proved)
	depositTx := NewDepositBtcTx(height, txIndex, blockTime, tx.Txid, ethAddr, BtcToSat(amount), proved)
	return depositTx, true, nil
}

func (b *bitcoinAgent) CheckState() error {
	// bitcoin sync blocks per half hour
	height, ok, err := b.chainStore.ReadBtcHeight()
	if err != nil {
		logger.Error("read beacon latest height error: %v", err)
		return err
	}
	if ok {
		diff := height - b.curHeight
		if diff < 2 { //normal 3
			logger.Error("bitcoin sync too slow, node maybe offline: diff %v prevHeight:%v curHeight:%v", diff, b.curHeight, height)
		}

	}
	b.curHeight = height
	return nil
}

func (b *bitcoinAgent) Close() error {
	return nil
}
func (b *bitcoinAgent) Name() string {
	return BitcoinAgentName
}

func NewMigrateBtcTx(height, txIndex, blockTime uint64, txId string, amount int64, proofed bool) *DbTx {
	return &DbTx{
		Hash:      DbValue(txId),
		TxIndex:   uint(txIndex),
		Height:    height,
		TxType:    common.DepositTx,
		ChainType: common.BitcoinChain,
		ProofType: common.BtcDepositType,
		BlockTime: blockTime,
		Amount:    amount,
		Proved:    proofed,
	}
}

func NewDepositBtcTx(height, txIndex, blockTime uint64, txId, ethAddr string, amount int64, proofed bool) *DbTx {
	return &DbTx{
		Hash:      DbValue(txId),
		TxIndex:   uint(txIndex),
		Height:    height,
		TxType:    common.DepositTx,
		ChainType: common.BitcoinChain,
		EthAddr:   DbValue(ethAddr),
		ProofType: common.BtcDepositType,
		BlockTime: blockTime,
		Amount:    amount,
		Proved:    proofed,
	}
}

func NewRedeemBtcTx(height, txIndex, blockTime uint64, txId string, amount int64, proofed bool) *DbTx {
	return &DbTx{
		Height:    height,
		TxIndex:   uint(txIndex),
		Hash:      DbValue(txId),
		TxType:    common.RedeemTx,
		ChainType: common.BitcoinChain,
		ProofType: common.BtcChangeType,
		BlockTime: blockTime,
		Proved:    proofed,
		Amount:    amount,
	}
}

func isOpZkpProto(outputs []btcjson.Vout) (string, error) {
	//https://mempool.space/zh/testnet4/tx/923d9f0fcb3654a343fd3e23d53f729c227f3ae77619e795e20c8b11a34bd358
	//op_return + length + ethAddr （20 byte） + extra （0+byte）
	//6a14e96af29bb5bb124c705c69034262fbc9fbb2d5f3
	for _, out := range outputs {
		if out.ScriptPubKey.Type == "nulldata" && strings.HasPrefix(out.ScriptPubKey.Hex, "6a") {
			ethAddr, err := getEthAddrFromScript(out.ScriptPubKey.Hex)
			if err != nil {
				return "", err
			}
			return ethAddr, nil
		}
	}
	return "", fmt.Errorf("no find zkbtc opreturn")
}

func getDepositAmount(txOuts []btcjson.Vout, addr string) float64 {
	var total float64
	for _, out := range txOuts {
		if out.ScriptPubKey.Address == addr {
			total = total + out.Value
		}
	}
	return total
}

func getRedeemAmount(txOuts []btcjson.Vout, addr string) float64 {
	var total float64
	for _, out := range txOuts {
		if out.ScriptPubKey.Address != addr {
			total = total + out.Value
		}
	}
	return total
}

func getEthAddrFromScript(script string) (string, error) {
	//6a 14 e96af29bb5bb124c705c69034262fbc9fbb2d5f3
	if len(script) < 44 {
		return "", fmt.Errorf("scritp lenght is less than 44")
	}
	if !strings.HasPrefix(script, "6a") {
		return "", fmt.Errorf("script is not start with 6a")
	}
	isHexAddress := ethcommon.IsHexAddress(script[4:44])
	if !isHexAddress {
		return "", fmt.Errorf("script is not hex address")
	}
	return script[4:44], nil
}
