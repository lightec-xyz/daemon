package node

import (
	"encoding/hex"
	"fmt"
	btcprovercom "github.com/lightec-xyz/btc_provers/circuits/common"
	btcprovertypes "github.com/lightec-xyz/btc_provers/circuits/types"
	baselevelUtil "github.com/lightec-xyz/btc_provers/utils/baselevel"
	btcproverClient "github.com/lightec-xyz/btc_provers/utils/client"
	grUtil "github.com/lightec-xyz/btc_provers/utils/grandrollup"
	midlevelUtil "github.com/lightec-xyz/btc_provers/utils/midlevel"
	recursiveduperUtil "github.com/lightec-xyz/btc_provers/utils/recursiveduper"
	upperlevelUtil "github.com/lightec-xyz/btc_provers/utils/upperlevel"
	"github.com/lightec-xyz/daemon/circuits"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/rpc/beacon"
	btcrpc "github.com/lightec-xyz/daemon/rpc/bitcoin"
	ethrpc "github.com/lightec-xyz/daemon/rpc/ethereum"
	"github.com/lightec-xyz/daemon/store"
	ethblock "github.com/lightec-xyz/provers/circuits/fabric/tx-in-eth2"
	txineth2 "github.com/lightec-xyz/provers/circuits/tx-in-eth2"
	proverType "github.com/lightec-xyz/provers/circuits/types"
	apiclient "github.com/lightec-xyz/provers/utils/api-client"
	"github.com/lightec-xyz/reLight/circuits/utils"
	"github.com/prysmaticlabs/prysm/v5/api/server/structs"
	"strconv"
)

type PreparedData struct {
	filestore        *FileStorage
	store            store.IStore
	proverClient     *btcproverClient.Client
	btcClient        *btcrpc.Client
	ethClient        *ethrpc.Client
	apiClient        *apiclient.Client
	beaconClient     *beacon.Client
	genesisPeriod    uint64
	btcGenesisHeight uint64 // startIndex
}

func (p *PreparedData) GetBtcGenesisData(endHeight uint64) (*rpc.BtcGenesisRequest, bool, error) {
	data, err := recursiveduperUtil.GetRecursiveProofData(p.proverClient, uint32(endHeight-1), uint32(p.btcGenesisHeight))
	if err != nil {
		logger.Error("get base level proof data error: %v %v", 0, err)
		return nil, false, err
	}
	firstProof, ok, err := p.filestore.GetBtcUpperProof(p.btcGenesisHeight, p.btcGenesisHeight+common.BtcUpperDistance)
	if err != nil {
		logger.Error("get btc genesis proof data error: %v %v", 0, err)
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	secondProof, ok, err := p.filestore.GetBtcUpperProof(p.btcGenesisHeight+common.BtcUpperDistance, p.btcGenesisHeight+common.BtcUpperDistance*2)
	if err != nil {
		logger.Error("get base level proof data error: %v %v", 0, err)
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	genesisRequest := rpc.BtcGenesisRequest{
		Data: data,
		First: rpc.Proof{
			Proof:   firstProof.Proof,
			Witness: firstProof.Witness,
		},
		Second: rpc.Proof{
			Proof:   secondProof.Proof,
			Witness: secondProof.Witness,
		},
	}
	return &genesisRequest, true, nil
}

func (p *PreparedData) GetBtcRecursiveData(endHeight uint64) (*rpc.BtcRecursiveRequest, bool, error) {
	data, err := recursiveduperUtil.GetRecursiveProofData(p.proverClient, uint32(endHeight-1), uint32(p.btcGenesisHeight))
	if err != nil {
		logger.Error("get base level proof data error: %v %v", 0, err)
		return nil, false, err
	}

	var fistProof rpc.Proof
	/*
			example:
			up1: 0~2, up2 2~4, up3 4~6  up4 6~8
			genesis: 0~4(up1,up2)
			recursive1: 0~6(genesis,up3)
		    recursive2: 0~8(recursive1,up4)
		    ....
	*/
	if endHeight == p.btcGenesisHeight+common.BtcUpperDistance*3 {
		genesisProof, ok, err := p.filestore.GetBtcGenesisProof()
		if err != nil {
			logger.Error("get base level proof data error: %v %v", 0, err)
			return nil, false, err
		}
		if !ok {
			return nil, false, nil
		}
		fistProof = rpc.Proof{
			Proof:   genesisProof.Proof,
			Witness: genesisProof.Witness,
		}

	} else if endHeight > p.btcGenesisHeight+common.BtcUpperDistance*3 {
		recursiveProof, ok, err := p.filestore.GetBtcRecursiveProof(endHeight-2*common.BtcUpperDistance, endHeight-common.BtcUpperDistance)
		if err != nil {
			logger.Error("get base level proof data error: %v %v", 0, err)
			return nil, false, err
		}
		if !ok {
			return nil, false, nil
		}
		fistProof = rpc.Proof{
			Proof:   recursiveProof.Proof,
			Witness: recursiveProof.Witness,
		}

	}
	secondProof, ok, err := p.filestore.GetBtcUpperProof(endHeight-common.BtcUpperDistance, endHeight)
	if err != nil {
		logger.Error("get base level proof data error: %v %v", 0, err)
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	btcRecursiveRequest := rpc.BtcRecursiveRequest{
		Data:  data,
		First: fistProof,
		Second: rpc.Proof{
			Proof:   secondProof.Proof,
			Witness: secondProof.Witness,
		},
	}
	return &btcRecursiveRequest, true, nil
}

func (p *PreparedData) GetBtcBaseData(endHeight uint64) (*rpc.BtcBaseRequest, bool, error) {
	data, err := baselevelUtil.GetBaseLevelProofData(p.proverClient, uint32(endHeight-1))
	if err != nil {
		logger.Error("get base level proof data error: %v %v", endHeight, err)
		return nil, false, err
	}
	baseRequest := rpc.BtcBaseRequest{
		Data: data,
	}
	return &baseRequest, true, nil
}

func (p *PreparedData) GetBtcMiddleData(endHeight uint64) (*rpc.BtcMiddleRequest, bool, error) {
	data, err := midlevelUtil.GetMidLevelProofData(p.proverClient, uint32(endHeight-1))
	if err != nil {
		logger.Error("get base level proof data error: %v %v", endHeight, err)
		return nil, false, err
	}
	var proofs []rpc.Proof
	for index := endHeight - common.BtcMiddleDistance; index < endHeight; index = index + common.BtcBaseDistance {
		startIndex := index
		endIndex := index + common.BtcBaseDistance
		baseProof, ok, err := p.filestore.GetBtcBaseProof(startIndex, endIndex)
		if err != nil {
			logger.Error("get base level proof data error: %v~%v %v", startIndex, endIndex, err)
			return nil, false, err
		}
		if ok {
			proofs = append(proofs, rpc.Proof{
				Proof:   baseProof.Proof,
				Witness: baseProof.Witness,
			})
		}
	}
	// todo
	if len(proofs) != 2 {
		logger.Error("get base level proof data error: %v %v", endHeight, err)
		return nil, false, fmt.Errorf("get base level proof data error: %v", endHeight)
	}
	baseRequest := rpc.BtcMiddleRequest{
		Data:   data,
		Proofs: proofs,
	}
	return &baseRequest, true, nil
}

func (p *PreparedData) GetBtcUpperData(endHeight uint64) (*rpc.BtcUpperRequest, bool, error) {
	data, err := upperlevelUtil.GetUpperLevelProofData(p.proverClient, uint32(endHeight-1))
	if err != nil {
		logger.Error("get base level proof data error: %v %v", endHeight, err)
		return nil, false, err
	}
	var proofs []rpc.Proof
	for index := endHeight - common.BtcUpperDistance; index < endHeight; index = index + common.BtcMiddleDistance {
		startIndex := index
		endIndex := index + common.BtcMiddleDistance
		middleProof, ok, err := p.filestore.GetBtcMiddleProof(startIndex, endIndex)
		if err != nil {
			logger.Error("get base level proof data error: %v~%v %v", startIndex, endIndex, err)
			return nil, false, err
		}
		if ok {
			proofs = append(proofs, rpc.Proof{
				Proof:   middleProof.Proof,
				Witness: middleProof.Witness,
			})
		}
	}
	baseRequest := rpc.BtcUpperRequest{
		Data:   data,
		Proofs: proofs,
	}
	return &baseRequest, true, nil
}

func (p *PreparedData) GetVerifyData(txHash string) (interface{}, bool, error) {
	tx, err := p.btcClient.GetTransaction(txHash)
	if err != nil {
		logger.Error("get verify tx error: %v %v", txHash, err)
		return nil, false, err
	}
	proofData, err := grUtil.GetDefaultGrandRollupProofData(p.proverClient, txHash, tx.Blockhash)
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

func (p *PreparedData) GetDepositData(txHash string) (*rpc.DepositRequest, bool, error) {
	tx, err := p.btcClient.GetTransaction(txHash)
	if err != nil {
		logger.Error("get deposit tx error: %v %v", txHash, err)
		return nil, false, err
	}
	proofData, err := grUtil.GetDefaultGrandRollupProofData(p.proverClient, txHash, tx.Blockhash)
	if err != nil {
		logger.Error("get deposit proof data error: %v %v", txHash, err)
		return nil, false, err
	}
	depositRequest := rpc.DepositRequest{
		TxHash:    txHash,
		BlockHash: tx.Blockhash,
		Data:      proofData,
	}
	return &depositRequest, true, nil
}

func (p *PreparedData) GetTxInEth2Data(txHash string, getSlotByNumber func(uint64) (uint64, error)) (*rpc.TxInEth2ProveRequest, bool, error) {
	txData, err := ethblock.GenerateTxInEth2Proof(p.ethClient.Client, p.apiClient, getSlotByNumber, txHash)
	if err != nil {
		logger.Error("get tx data error: %v", err)
		return nil, false, err
	}
	return &rpc.TxInEth2ProveRequest{
		TxHash: txHash,
		TxData: txData,
	}, true, nil
}

func (p *PreparedData) GetBlockHeaderRequestData(index uint64) (*rpc.BlockHeaderRequest, bool, error) {
	// todo
	finalizedSlot, ok, err := p.filestore.GetNearTxSlotFinalizedSlot(index)
	if err != nil {
		logger.Error("get finalized slot error: %v", err)
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	logger.Debug("get beaconHeader %v ~ %v", index, finalizedSlot)
	beaconBlockHeaders, err := p.beaconClient.RetrieveBeaconHeaders(index, finalizedSlot)
	if err != nil {
		logger.Error("get beacon block headers error: %v", err)
		return nil, false, err
	}
	if len(beaconBlockHeaders) == 0 {
		return nil, false, fmt.Errorf("never should happen %v", index)
	}
	beginSlot, beginRoot, err := BeaconBlockHeaderToSlotAndRoot(beaconBlockHeaders[0])
	if err != nil {
		logger.Error("get beacon block headers error: %v", err)
		return nil, false, err
	}
	endSlot, endRoot, err := BeaconBlockHeaderToSlotAndRoot(beaconBlockHeaders[len(beaconBlockHeaders)-1])
	if err != nil {
		logger.Error("get beacon block headers error: %v", err)
		return nil, false, err
	}
	return &rpc.BlockHeaderRequest{
		Index:     index,
		BeginSlot: beginSlot,
		EndSlot:   endSlot,
		BeginRoot: hex.EncodeToString(beginRoot),
		EndRoot:   hex.EncodeToString(endRoot),
		Headers:   beaconBlockHeaders[1:],
	}, true, nil
}

func (p *PreparedData) GetRecursiveData(period uint64) (interface{}, bool, error) {
	//todo
	genesisPeriod := p.filestore.GetGenesisPeriod()
	genesisId, ok, err := p.GetSyncCommitRootId(genesisPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v Index genesis commitId no find", genesisPeriod)
		return nil, false, nil
	}
	relayId, ok, err := p.GetSyncCommitRootId(period)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v Index relay commitId no find", period)
		return nil, false, nil
	}
	endPeriod := period + 1
	endId, ok, err := p.GetSyncCommitRootId(endPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v Index end commitId no find", endPeriod)
		return nil, false, nil
	}
	secondProof, exists, err := p.filestore.GetUnitProof(period)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !exists {
		logger.Warn("no find %v unit proof Data, send new proof request", period)
		return nil, false, nil
	}

	prePeriod := period - 1
	firstProof, exists, err := p.filestore.GetRecursiveProof(prePeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !exists {
		logger.Warn("no find %v Index recursive Data, send new proof request", prePeriod)
		return nil, false, nil
	}
	return &rpc.SyncCommRecursiveRequest{
		Choice:        circuits.SyncRecursive,
		FirstProof:    firstProof.Proof,
		FirstWitness:  firstProof.Witness,
		SecondProof:   secondProof.Proof,
		SecondWitness: secondProof.Witness,
		BeginId:       genesisId,
		RelayId:       relayId,
		EndId:         endId,
	}, true, nil
}

func (p *PreparedData) getSlotByNumber(number uint64) (uint64, error) {
	slot, ok, err := ReadBeaconSlot(p.store, number)
	if err != nil {
		logger.Error("get slot error: %v %v", number, err)
		return 0, err
	}
	if !ok {
		return 0, fmt.Errorf("no find %v slot", number)
	}
	return slot, nil
}

func (p *PreparedData) GetRecursiveGenesisData(period uint64) (interface{}, bool, error) {
	genesisPeriod := p.filestore.GetGenesisPeriod()
	genesisId, ok, err := p.GetSyncCommitRootId(genesisPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v Index genesis commitId no find", genesisPeriod)
		return nil, false, nil
	}
	relayPeriod := period
	relayId, ok, err := p.GetSyncCommitRootId(relayPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v  Index relay commitId no find ", relayPeriod)
		return nil, false, nil
	}
	endPeriod := relayPeriod + 1
	endId, ok, err := p.GetSyncCommitRootId(endPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v Index end commitId no find", endPeriod)
		return nil, false, nil
	}

	fistProof, firstExists, err := p.filestore.GetGenesisProof()
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !firstExists {
		logger.Warn("no find genesis proof ,start new proof request")
		return nil, false, nil
	}
	secondProof, secondExists, err := p.filestore.GetUnitProof(relayPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !secondExists {
		logger.Warn("no find %v Index unit proof , send new proof request", relayPeriod)
		return nil, false, nil
	}
	return &rpc.SyncCommRecursiveRequest{
		Choice:        circuits.SyncGenesis,
		FirstProof:    fistProof.Proof,
		FirstWitness:  fistProof.Witness,
		SecondProof:   secondProof.Proof,
		SecondWitness: secondProof.Witness,
		BeginId:       genesisId,
		RelayId:       relayId,
		EndId:         endId,
	}, true, nil

}

func (p *PreparedData) GetSyncComGenesisData() (*rpc.SyncCommGenesisRequest, bool, error) {
	genesisPeriod := p.filestore.GetGenesisPeriod()
	genesisId, ok, err := p.GetSyncCommitRootId(genesisPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v Index genesis commitId  no find", genesisPeriod)
		return nil, false, nil
	}

	nextPeriod := genesisPeriod + 1
	firstId, ok, err := p.GetSyncCommitRootId(nextPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v Index first commitId no find", nextPeriod)
		return nil, false, nil
	}
	secondPeriod := nextPeriod + 1
	secondId, ok, err := p.GetSyncCommitRootId(secondPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v Index second commitId no find", secondPeriod)
		return nil, false, nil
	}

	firstProof, exists, err := p.filestore.GetUnitProof(genesisPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !exists {
		logger.Warn("get genesis Data,first proof not exists: %v Index", genesisPeriod)
		return nil, false, nil
	}
	logger.Info("get genesis first proof: %v", genesisPeriod)

	secondProof, exists, err := p.filestore.GetUnitProof(nextPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}

	if !exists {
		logger.Warn("get genesis Data,second proof not exists: %v Index", nextPeriod)
		return nil, false, nil
	}
	logger.Info("get genesis second proof: %v", nextPeriod)
	genesisProofParam := &rpc.SyncCommGenesisRequest{
		FirstProof:    firstProof.Proof,
		FirstWitness:  firstProof.Witness,
		SecondProof:   secondProof.Proof,
		SecondWitness: secondProof.Witness,
		GenesisID:     genesisId,
		FirstID:       firstId,
		SecondID:      secondId,
	}
	return genesisProofParam, true, nil

}

func (p *PreparedData) GetSyncCommitRootId(period uint64) ([]byte, bool, error) {
	update, ok, err := p.GetSyncCommitUpdate(period)
	if err != nil {
		logger.Error("get %v Index update error: %v", period, err)
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	syncCommitRoot, err := circuits.SyncCommitRoot(update)
	if err != nil {
		logger.Error("get commit root: %v error: %v", period, err)
		return nil, false, err
	}
	return syncCommitRoot, true, nil
}

func (p *PreparedData) GetSyncComUnitData(period uint64) (*rpc.SyncCommUnitsRequest, bool, error) {
	update, ok, err := p.GetSyncCommitUpdate(period)
	if err != nil {
		logger.Error("get %v Index update error: %v", period, err)
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	commUnitsRequest := rpc.SyncCommUnitsRequest{
		Data:    update,
		Index:   period,
		Version: update.Version,
	}
	return &commUnitsRequest, true, nil

}

func (p *PreparedData) GetReverseHash(height uint64) (string, error) {
	hash, err := p.btcClient.GetBlockHash(int64(height))
	if err != nil {
		logger.Error("get block header error: %v %v", height, err)
		return "", err
	}
	reverseHash, err := common.ReverseHex(hash)
	if err != nil {
		logger.Error("reverse hex error: %v", err)
		return "", err
	}
	return reverseHash, nil
}
func (p *PreparedData) GetBtcMidBlockHeader(start, end uint64) (*btcprovertypes.BlockHeaderChain, error) {
	startHash, err := p.GetReverseHash(start)
	if err != nil {
		logger.Error("get block header error: %v %v", start, err)
		return nil, err
	}

	endHash, err := p.GetReverseHash(end)
	if err != nil {
		logger.Error("get block header error: %v %v", end, err)
		return nil, err
	}
	var middleHeaders []string
	for index := start + 1; index <= end; index++ {
		header, err := p.btcClient.GetHexBlockHeader(int64(index))
		if err != nil {
			logger.Error("get block header error: %v %v", index, err)
			return nil, err
		}
		middleHeaders = append(middleHeaders, header)
	}
	data := &btcprovertypes.BlockHeaderChain{
		BeginHeight:        start,
		BeginHash:          startHash,
		EndHeight:          end,
		EndHash:            endHash,
		MiddleBlockHeaders: middleHeaders,
	}
	err = data.Verify()
	if err != nil {
		logger.Error("verify block header error: %v", err)
		return nil, err
	}
	return data, nil

}

func (p *PreparedData) GetBtcWrapData(start, end uint64) (*rpc.BtcWrapRequest, error) {
	startHash, err := p.GetReverseHash(start)
	if err != nil {
		logger.Error("get block header error: %v %v", start, err)
		return nil, err
	}
	endHash, err := p.GetReverseHash(end)
	if err != nil {
		logger.Error("get block header error: %v %v", end, err)
		return nil, err
	}
	nRequired := end - start
	var proof *StoreProof
	var ok bool
	var flag string
	if nRequired <= btcprovercom.MaxNbBlockPerBulk { // todo
		proof, ok, err = p.filestore.GetBtcBulkProof(start, end)
		if err != nil {
			logger.Error("get btc bulk proof error: %v", err)
			return nil, err
		}
		if !ok {
			return nil, err
		}
		flag = circuits.BtcBulk
	} else {
		proof, ok, err = p.filestore.GetBtcPackedProof(start)
		if err != nil {
			logger.Error("get btc bulk proof error: %v", err)
			return nil, err
		}
		if !ok {
			return nil, err
		}
		flag = circuits.BtcPacked
	}
	data := &rpc.BtcWrapRequest{
		Flag:      flag,
		Proof:     proof.Proof,
		Witness:   proof.Witness,
		BeginHash: startHash,
		EndHash:   endHash,
		NbBlocks:  nRequired,
	}
	return data, nil
}

func (p *PreparedData) GetSyncCommitUpdate(period uint64) (*utils.SyncCommitteeUpdate, bool, error) {
	var currentPeriodUpdate structs.LightClientUpdateWithVersion
	exists, err := p.filestore.GetUpdate(period, &currentPeriodUpdate)
	if err != nil {
		logger.Error("get %v index update error: %v", period, err)
		return nil, false, err
	}
	if !exists {
		logger.Warn("no find %v index update Data", period)
		return nil, false, nil
	}
	var update utils.SyncCommitteeUpdate
	err = common.ParseObj(currentPeriodUpdate.Data, &update)
	if err != nil {
		logger.Error("parse obj error: %v %v", period, err)
		return nil, false, err
	}
	update.Version = currentPeriodUpdate.Version
	if p.filestore.GetGenesisPeriod() == period {
		var genesisData structs.LightClientBootstrapResponse
		genesisExists, err := p.filestore.GetBootstrap(&genesisData)
		if err != nil {
			logger.Error("get genesis update error: %v %v", period, err)
			return nil, false, err
		}
		if !genesisExists {
			logger.Warn("no find genesis update Data,%v", period)
			return nil, false, nil
		}
		var genesisCommittee utils.SyncCommittee
		err = common.ParseObj(genesisData.Data.CurrentSyncCommittee, &genesisCommittee)
		if err != nil {
			logger.Error(err.Error())
			return nil, false, err
		}
		update.CurrentSyncCommittee = &genesisCommittee
	} else {
		prePeriod := period - 1
		if prePeriod < p.filestore.GetGenesisPeriod() {
			logger.Error("should never happen: %v", prePeriod)
			return nil, false, nil
		}
		var preUpdateData structs.LightClientUpdateWithVersion
		preUpdateExists, err := p.filestore.GetUpdate(prePeriod, &preUpdateData)
		if err != nil {
			logger.Error("get %v index update error: %v", prePeriod, err)
			return nil, false, err
		}
		if !preUpdateExists {
			logger.Warn("get unit Data,no find %v Index update Data", prePeriod)
			return nil, false, nil
		}
		var currentSyncCommittee utils.SyncCommittee
		err = common.ParseObj(preUpdateData.Data.NextSyncCommittee, &currentSyncCommittee)
		if err != nil {
			logger.Error(err.Error())
			return nil, false, err
		}
		update.CurrentSyncCommittee = &currentSyncCommittee
	}
	ok, err := common.VerifyLightClientUpdate(update)
	if err != nil {
		logger.Error("verify light client update error: %v %v", period, err)
		return nil, false, err
	}
	if !ok {
		logger.Error("update data verify false: %v", period)
		return nil, false, fmt.Errorf("update data verify false: %v", period)
	}
	return &update, true, nil

}

func (p *PreparedData) GetBhfUpdateData(finalizedSlot, genesisPeriod uint64) (*rpc.BlockHeaderFinalityRequest, bool, error) {
	logger.Debug("get bhf update data: %v", finalizedSlot)
	var currentFinalityUpdate structs.LightClientFinalityUpdateEvent
	exists, err := p.filestore.GetFinalityUpdate(finalizedSlot, &currentFinalityUpdate)
	if err != nil {
		logger.Error("get finality update error: %v %v", finalizedSlot, err)
		return nil, false, err
	}
	if !exists {
		logger.Warn("no find finality update: %v", finalizedSlot)
		return nil, false, nil
	}

	genesisId, ok, err := p.GetSyncCommitRootId(genesisPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v Index genesis commitId no find", genesisPeriod)
		return nil, false, nil
	}
	attestedSlot, err := strconv.ParseUint(currentFinalityUpdate.Data.AttestedHeader.Slot, 10, 64)
	if err != nil {
		logger.Error("parse big error %v %v", currentFinalityUpdate.Data.AttestedHeader.Slot, err)
		return nil, false, err
	}
	// todo
	period := (attestedSlot / 8192) - 1
	recursiveProof, ok, err := p.filestore.GetRecursiveProof(period)
	if err != nil {
		logger.Error("get recursive proof error: %v %v", period, err)
		return nil, false, err
	}
	if !ok {
		logger.Warn("no find recursive proof: %v", period)
		return nil, false, nil
	}
	outerPeriod := period + 1
	logger.Debug("get bhf update data finalizedSlot: %v,recPeriod: %v,outPeriod %v", finalizedSlot, period, outerPeriod)
	outerProof, ok, err := p.filestore.GetOuterProof(outerPeriod)
	if err != nil {
		logger.Error("get outer proof error: %v %v", outerPeriod, err)
		return nil, false, err
	}
	if !ok {
		logger.Warn("no find outer proof: %v", outerPeriod)
		return nil, false, nil
	}

	var finalUpdate proverType.FinalityUpdate
	err = common.ParseObj(currentFinalityUpdate.Data, &finalUpdate)
	if err != nil {
		logger.Error("parse finality update error: %v", err)
		return nil, false, err
	}
	finalUpdate.Version = currentFinalityUpdate.Version

	currentSyncCommitUpdate, ok, err := p.GetSyncCommitUpdate(outerPeriod)
	if err != nil {
		logger.Error("get sync committee update error: %v %v", outerPeriod, err)
		return nil, false, err
	}
	if !ok {
		logger.Error("no find sync committee update: %v", period)
		return nil, false, nil
	}
	var scUpdate proverType.SyncCommitteeUpdate
	err = common.ParseObj(currentSyncCommitUpdate, &scUpdate)
	if err != nil {
		logger.Error("parse sync committee update error: %v", err)
		return nil, false, err
	}
	request := rpc.BlockHeaderFinalityRequest{
		Index:            finalizedSlot,
		GenesisSCSSZRoot: fmt.Sprintf("%x", genesisId),
		RecursiveProof:   recursiveProof.Proof,
		RecursiveWitness: recursiveProof.Witness,
		OuterProof:       outerProof.Proof,
		OuterWitness:     outerProof.Witness,
		FinalityUpdate:   &finalUpdate,
		ScUpdate:         &scUpdate,
	}
	return &request, true, nil
}

func (p *PreparedData) GetRedeemRequestData(genesisPeriod, txSlot uint64, txHash string) (*rpc.RedeemRequest, bool, error) {
	txProof, ok, err := p.filestore.GetTxProof(txHash)
	if err != nil {
		logger.Error("get tx proof error: %v", err)
		return nil, false, err
	}
	if !ok {
		logger.Debug("proof request data not prepared: %v", txHash)
		return nil, false, nil
	}
	finalizedSlot, ok, err := p.filestore.GetNearTxSlotFinalizedSlot(txSlot)
	if err != nil {
		logger.Error("get bhf update proof error: %v", err)
		return nil, false, err
	}
	if !ok {
		logger.Debug("proof request data not prepared: %v", txSlot)
		return nil, false, nil
	}
	blockHeaderProof, ok, err := p.filestore.GetBeaconHeaderProof(txSlot, finalizedSlot)
	if err != nil {
		logger.Error("get block header proof error: %v", err)
		return nil, false, err
	}
	if !ok {
		logger.Debug("proof request data not prepared: %v", txSlot)
		return nil, false, nil
	}
	bhfProof, ok, err := p.filestore.GetBhfProof(finalizedSlot)
	if err != nil {
		logger.Error("get bhf update proof error: %v", err)
		return nil, false, err
	}
	if !ok {
		logger.Warn("no find bhf update %v", finalizedSlot)
		return nil, false, nil
	}
	genesisRoot, ok, err := p.GetSyncCommitRootId(genesisPeriod)
	if err != nil {
		logger.Error("get genesis root error: %v", err)
		return nil, false, err
	}
	if !ok {
		logger.Warn("no find genesis root %v", genesisPeriod)
		return nil, false, nil
	}

	var finalityUpdate *structs.LightClientFinalityUpdateEvent
	ok, err = p.filestore.GetFinalityUpdate(finalizedSlot, &finalityUpdate)
	if err != nil {
		logger.Error("get finality update error: %v", err)
		return nil, false, err
	}
	if !ok {
		logger.Warn("no find finality update %v", finalizedSlot)
		return nil, false, nil
	}

	attestedSlot, err := strconv.ParseUint(finalityUpdate.Data.AttestedHeader.Slot, 10, 64)
	if err != nil {
		logger.Error("parse slot error: %v", err)
		return nil, false, err
	}
	period := attestedSlot / common.SlotPerPeriod
	currentRoot, ok, err := p.GetSyncCommitRootId(period)
	if err != nil {
		logger.Error("get current root error: %v", err)
		return nil, false, err
	}
	beginID, endId, err := p.GetBeaconHeaderId(txSlot, finalizedSlot)
	if err != nil {
		logger.Error("get begin and end id error: %v", err)
		return nil, false, err
	}
	txVar, receiptVar, err := txineth2.GenerateTxAndReceiptU128Padded(p.ethClient.Client, txHash)
	if err != nil {
		logger.Error("get tx and receipt error: %v", err)
		return nil, false, err
	}
	txVarHex, err := common.TxVarToHex(txVar)
	if err != nil {
		logger.Error("tx var to bytes error: %v", err)
		return nil, false, err
	}
	receiptVarHex, err := common.ReceiptVarToHex(receiptVar)
	if err != nil {
		logger.Error("receipt var to bytes error: %v", err)
		return nil, false, err
	}

	redeemRequest := rpc.RedeemRequest{
		TxHash:           txHash,
		TxProof:          txProof.Proof,
		TxWitness:        txProof.Witness,
		BhProof:          blockHeaderProof.Proof,
		BhWitness:        blockHeaderProof.Witness,
		BhfProof:         bhfProof.Proof,
		BhfWitness:       bhfProof.Witness,
		GenesisScRoot:    hex.EncodeToString(genesisRoot),
		BeginId:          hex.EncodeToString(beginID),
		EndId:            hex.EncodeToString(endId),
		CurrentSCSSZRoot: hex.EncodeToString(currentRoot),
		TxVar:            txVarHex,
		ReceiptVar:       receiptVarHex,
	}
	return &redeemRequest, true, nil

}

func (p *PreparedData) GetBeaconHeaderId(start, end uint64) ([]byte, []byte, error) {
	beaconBlockHeaders, err := p.beaconClient.RetrieveBeaconHeaders(start, end)
	if err != nil {
		logger.Error("get beacon block headers error: %v", err)
		return nil, nil, err
	}
	if len(beaconBlockHeaders) == 0 {
		return nil, nil, fmt.Errorf("never should happen %v", start)
	}
	_, beginRoot, err := BeaconBlockHeaderToSlotAndRoot(beaconBlockHeaders[0])
	if err != nil {
		logger.Error("get beacon block headers error: %v", err)
		return nil, nil, err
	}
	_, endRoot, err := BeaconBlockHeaderToSlotAndRoot(beaconBlockHeaders[len(beaconBlockHeaders)-1])
	if err != nil {
		logger.Error("get beacon block headers error: %v", err)
		return nil, nil, err
	}
	return beginRoot, endRoot, nil
}

func NewPreparedData(filestore *FileStorage, store store.IStore, genesisPeriod, btcGenesisHeight uint64, proverClient *btcproverClient.Client, btcClient *btcrpc.Client,
	ethClient *ethrpc.Client, apiClient *apiclient.Client, beaconClient *beacon.Client) (*PreparedData, error) {
	return &PreparedData{
		filestore:        filestore,
		store:            store,
		proverClient:     proverClient,
		btcClient:        btcClient,
		ethClient:        ethClient,
		apiClient:        apiClient,
		beaconClient:     beaconClient,
		genesisPeriod:    genesisPeriod,
		btcGenesisHeight: btcGenesisHeight,
	}, nil
}
