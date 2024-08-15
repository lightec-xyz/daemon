package node

import (
	"encoding/hex"
	"fmt"
	grUtil "github.com/lightec-xyz/btc_provers/utils/txinchain"
	"github.com/lightec-xyz/daemon/rpc/dfinity"
	"strconv"

	baselevelUtil "github.com/lightec-xyz/btc_provers/utils/blockchain"
	blockCu "github.com/lightec-xyz/btc_provers/utils/blockchain"
	blockDu "github.com/lightec-xyz/btc_provers/utils/blockdepth"
	btcproverClient "github.com/lightec-xyz/btc_provers/utils/client"
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
)

type Prepared struct {
	filestore        *FileStorage
	store            store.IStore
	proverClient     *btcproverClient.Client
	btcClient        *btcrpc.Client
	ethClient        *ethrpc.Client
	apiClient        *apiclient.Client
	beaconClient     *beacon.Client
	icpClient        *dfinity.Client
	genesisPeriod    uint64
	minerAddr        string
	btcGenesisHeight uint64 // startIndex
}

func (p *Prepared) GetBtcBaseRequest(endHeight uint64) (*rpc.BtcBaseRequest, bool, error) {
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

func (p *Prepared) GetBtcMiddleRequest(endHeight uint64) (*rpc.BtcMiddleRequest, bool, error) {
	data, err := blockCu.GetMidLevelProofData(p.proverClient, uint32(endHeight-1))
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
	baseRequest := rpc.BtcMiddleRequest{
		Data:   data,
		Proofs: proofs,
	}
	return &baseRequest, true, nil
}

func (p *Prepared) GetBtcUpperRequest(endHeight uint64) (*rpc.BtcUpperRequest, bool, error) {
	data, err := blockCu.GetUpperLevelProofData(p.proverClient, uint32(endHeight-1))
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

func (p *Prepared) GetTxInEth2Request(txHash string, getSlotByNumber func(uint64) (uint64, error)) (*rpc.TxInEth2ProveRequest, bool, error) {
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

func (p *Prepared) GetBlockHeaderRequest(index uint64) (*rpc.BlockHeaderRequest, bool, error) {
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

func (p *Prepared) GetRecursiveRequest(period uint64) (interface{}, bool, error) {
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
		Choice: circuits.SyncRecursive,
		FirstProof: rpc.Proof{
			Proof:   firstProof.Proof,
			Witness: firstProof.Witness},
		SecondProof: rpc.Proof{
			Proof:   secondProof.Proof,
			Witness: secondProof.Witness,
		},
		BeginId: genesisId,
		RelayId: relayId,
		EndId:   endId,
	}, true, nil
}

func (p *Prepared) getSlotByNumber(number uint64) (uint64, error) {
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

func (p *Prepared) GetRecursiveGenesisRequest(period uint64) (interface{}, bool, error) {
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
		Choice: circuits.SyncGenesis,
		FirstProof: rpc.Proof{
			Proof:   fistProof.Proof,
			Witness: fistProof.Witness,
		},
		SecondProof: rpc.Proof{
			Proof:   secondProof.Proof,
			Witness: secondProof.Witness,
		},
		BeginId: genesisId,
		RelayId: relayId,
		EndId:   endId,
	}, true, nil

}

func (p *Prepared) GetSyncComGenesisRequest() (*rpc.SyncCommGenesisRequest, bool, error) {
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
		FirstProof: rpc.Proof{
			Proof:   firstProof.Proof,
			Witness: firstProof.Witness,
		},
		SecondProof: rpc.Proof{
			Proof:   secondProof.Proof,
			Witness: secondProof.Witness,
		},
		GenesisID: genesisId,
		FirstID:   firstId,
		SecondID:  secondId,
	}
	return genesisProofParam, true, nil

}

func (p *Prepared) GetSyncCommitRootId(period uint64) ([]byte, bool, error) {
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

func (p *Prepared) GetSyncComUnitRequest(period uint64) (*rpc.SyncCommUnitsRequest, bool, error) {
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

func (p *Prepared) GetReverseHash(height uint64) (string, error) {
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
func (p *Prepared) GetBtcBulkRequest(start, end uint64) (*rpc.BtcBulkRequest, error) {
	proofData, err := blockDu.GetBlockBulkProofData(p.proverClient, uint32(start), uint32(end))
	if err != nil {
		logger.Error("btc bulk data error: %v", err)
		return nil, err
	}
	return &rpc.BtcBulkRequest{
		Data: proofData,
	}, nil

}

func (p *Prepared) GetBtcPackRequest(start, end uint64) (*rpc.BtcPackedRequest, bool, error) {
	proofData, err := blockDu.GetPackedProofData(p.proverClient, uint32(start), uint32(end))
	if err != nil {
		logger.Error("btc bulk data error: %v", err)
		return nil, false, err
	}
	recursiveProof, ok, err := p.filestore.GetBtcDepthRecursiveProof(start, end)
	if err != nil {
		logger.Error("get btc depth recursive proof data error: %v %v", 0, err)
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	bulkProof, ok, err := p.filestore.GetBtcBulkProof(start, end)
	if err != nil {
		logger.Error("get btc bulk proof data error: %v %v", 0, err)
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	return &rpc.BtcPackedRequest{
		Data: proofData,
		Recursive: rpc.Proof{
			Proof:   recursiveProof.Proof,
			Witness: recursiveProof.Witness,
		},
		Bulk: rpc.Proof{
			Proof:   bulkProof.Proof,
			Witness: bulkProof.Witness,
		},
	}, true, nil

}

func (p *Prepared) GetSyncCommitUpdate(period uint64) (*utils.SyncCommitteeUpdate, bool, error) {
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

func (p *Prepared) GetBhfUpdateRequest(finalizedSlot, genesisPeriod uint64) (*rpc.BlockHeaderFinalityRequest, bool, error) {
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
		RecursiveProof: rpc.Proof{
			Proof:   recursiveProof.Proof,
			Witness: recursiveProof.Witness,
		},
		OuterProof: rpc.Proof{
			Proof:   outerProof.Proof,
			Witness: outerProof.Witness,
		},
		FinalityUpdate: &finalUpdate,
		ScUpdate:       &scUpdate,
	}
	return &request, true, nil
}

func (p *Prepared) GetRedeemRequest(genesisPeriod, txSlot uint64, txHash string) (*rpc.RedeemRequest, bool, error) {
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
		TxHash: txHash,
		TxProof: rpc.Proof{
			Proof:   txProof.Proof,
			Witness: txProof.Witness,
		},
		BhProof: rpc.Proof{
			Proof:   blockHeaderProof.Proof,
			Witness: blockHeaderProof.Witness,
		},
		BhfProof: rpc.Proof{
			Proof:   bhfProof.Proof,
			Witness: bhfProof.Witness,
		},
		GenesisScRoot:    hex.EncodeToString(genesisRoot),
		BeginId:          beginID,
		EndId:            endId,
		CurrentSCSSZRoot: hex.EncodeToString(currentRoot),
		TxVar:            txVarHex,
		ReceiptVar:       receiptVarHex,
	}
	return &redeemRequest, true, nil

}

func (p *Prepared) GetBeaconHeaderId(start, end uint64) ([]byte, []byte, error) {
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

func (p *Prepared) GetBtcDuperRecursiveRequest(index uint64) (*rpc.BtcDuperRecursiveRequest, bool, error) {
	data, err := blockCu.GetRecursiveProofData(p.proverClient, uint32(index), uint32(p.btcGenesisHeight))
	if err != nil {
		logger.Error("get recursive proof data error: %v", err)
		return nil, false, err
	}
	var firstProof *StoreProof
	var ok bool
	// todo
	if index == p.btcGenesisHeight+2*common.BtcUpperDistance {
		firstProof, ok, err = p.filestore.GetBtcDuperGenesisProof()
		if err != nil {
			logger.Error("get btc genesis proof data error: %v %v", 0, err)
			return nil, false, err
		}
		if !ok {
			return nil, false, nil
		}
	} else {
		firstProof, ok, err = p.filestore.GetBtcDuperRecursiveProof(index, index+common.BtcUpperDistance)
		if err != nil {
			logger.Error("get btc genesis proof data error: %v %v", 0, err)
			return nil, false, err
		}
		if !ok {
			return nil, false, nil
		}
	}
	secondProof, ok, err := p.filestore.GetBtcUpperProof(index+common.BtcUpperDistance, index+common.BtcUpperDistance*2)
	if err != nil {
		logger.Error("get base level proof data error: %v %v", 0, err)
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	request := rpc.BtcDuperRecursiveRequest{
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
	return &request, false, nil
}

func (p *Prepared) GetBtcDepthRecursiveRequest(start, end uint64) (*rpc.BtcDepthRecursiveRequest, bool, error) {
	proofData, err := blockDu.GetPackedProofData(p.proverClient, uint32(start), uint32(end))
	if err != nil {
		logger.Error("btc bulk data error: %v", err)
		return nil, false, err
	}
	var firstProof *StoreProof
	var ok bool
	// todo
	if p.btcGenesisHeight == start+common.CapacityBulkUint {
		firstProof, ok, err = p.filestore.GetBtcDepthGenesisProof()
		if err != nil {
			logger.Error("get btc genesis proof data error: %v %v", 0, err)
			return nil, false, err
		}
		if !ok {
			return nil, false, nil
		}
	} else {
		firstProof, ok, err = p.filestore.GetBtcDepthRecursiveProof(start, end)
		if err != nil {
			logger.Error("get btc genesis proof data error: %v %v", 0, err)
			return nil, false, err
		}
		if !ok {
			return nil, false, nil
		}
	}
	secondProof, ok, err := p.filestore.GetBtcBulkProof(start, end)
	if err != nil {
		logger.Error("get btc bulk proof data error: %v %v", 0, err)
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	request := rpc.BtcDepthRecursiveRequest{
		Data: proofData,
		Recursive: rpc.Proof{
			Proof:   firstProof.Proof,
			Witness: firstProof.Witness,
		},
		Unit: rpc.Proof{
			Proof:   secondProof.Proof,
			Witness: secondProof.Witness,
		},
	}
	return &request, false, nil
}

func (p *Prepared) GetBtcChainRequest(start, end uint64) (*rpc.BtcChainRequest, bool, error) {
	data, err := blockCu.GetBlockChainProofData(p.proverClient, uint32(start), uint32(end))
	if err != nil {
		logger.Error("get block chain proof data error: %v", err)
		return nil, false, err
	}
	recursiveProof, ok, err := p.filestore.GetBtcDuperRecursiveProof(start, end)
	if err != nil {
		logger.Error("get btc duper recursive proof data error: %v %v", 0, err)
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	baseProof, ok, err := p.filestore.GetBtcBaseProof(start, end)
	if err != nil {
		logger.Error("get btc base proof data error: %v %v", 0, err)
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	midLevelProof, ok, err := p.filestore.GetBtcMiddleProof(start, end)
	if err != nil {
		logger.Error("get btc mid level proof data error: %v %v", 0, err)
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	upperProof, ok, err := p.filestore.GetBtcUpperProof(start, end)
	if err != nil {
		logger.Error("get btc upper proof data error: %v %v", 0, err)
		return nil, false, err
	}

	request := rpc.BtcChainRequest{
		Data: data,
		Recursive: rpc.Proof{
			Proof:   recursiveProof.Proof,
			Witness: recursiveProof.Witness,
		},
		Base: rpc.Proof{
			Proof:   baseProof.Proof,
			Witness: baseProof.Witness,
		},
		MidLevel: rpc.Proof{
			Proof:   midLevelProof.Proof,
			Witness: midLevelProof.Witness,
		},
		Upper: rpc.Proof{
			Proof:   upperProof.Proof,
			Witness: upperProof.Witness,
		},
	}
	return &request, false, nil
}

func (p *Prepared) GetBtcDepositRequest(hash string) (*rpc.BtcDepositRequest, bool, error) {
	data, err := grUtil.GetTxInChainProofData(p.proverClient, hash, "", 0, 0, 0)
	if err != nil {
		logger.Error("get tx in chain proof data error: %v", err)
		return nil, false, err
	}
	blockChainProof, ok, err := p.filestore.GetBtcBlockChainProof(hash)
	if err != nil {
		logger.Error("get btc block chain proof data error: %v %v", 0, err)
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	txDepthProof, ok, err := p.filestore.GetBtcDuperRecursiveProof()
	if err != nil {
		logger.Error("get btc tx depth proof data error: %v %v", 0, err)
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	cpDepthProof, ok, err := p.filestore.GetBtcDepthRecursiveProof(hash)
	if err != nil {
		logger.Error("get btc cp depth proof data error: %v %v", 0, err)
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	blockSignature, err := p.icpClient.BlockSignature()
	if err != nil {
		logger.Error("get block signature error: %v", err)
		return nil, false, err
	}
	r, s, err := blockSignature.ToRS()
	if err != nil {
		logger.Error("get block signature error: %v", err)
		return nil, false, err
	}
	request := rpc.BtcDepositRequest{
		Data: data,
		BlockChain: rpc.Proof{
			Proof:   blockChainProof.Proof,
			Witness: blockChainProof.Witness,
		},
		TxDepth: rpc.Proof{
			Proof:   txDepthProof.Proof,
			Witness: txDepthProof.Witness,
		},
		CpDepth: rpc.Proof{
			Proof:   cpDepthProof.Proof,
			Witness: cpDepthProof.Witness,
		},
		R:          hex.EncodeToString(r),
		S:          hex.EncodeToString(s),
		ProverAddr: p.minerAddr,
	}
	return &request, false, nil
}

func (p *Prepared) GetBtcChangeRequest(hash string) (*rpc.BtcChangeRequest, bool, error) {
	data, err := grUtil.GetTxInChainProofData(p.proverClient, hash, "", 0, 0, 0)
	if err != nil {
		logger.Error("get tx in chain proof data error: %v", err)
		return nil, false, err
	}
	blockChainProof, ok, err := p.filestore.GetBtcBlockChainProof(hash)
	if err != nil {
		logger.Error("get btc block chain proof data error: %v %v", 0, err)
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	txDepthProof, ok, err := p.filestore.GetBtcDuperRecursiveProof()
	if err != nil {
		logger.Error("get btc tx depth proof data error: %v %v", 0, err)
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	cpDepthProof, ok, err := p.filestore.GetBtcDepthRecursiveProof(hash)
	if err != nil {
		logger.Error("get btc cp depth proof data error: %v %v", 0, err)
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	redeemProof, ok, err := p.filestore.GetTxProof("")
	if err != nil {
		logger.Error("get btc redeem proof data error: %v %v", 0, err)
		return nil, false, err
	}

	blockSignature, err := p.icpClient.BlockSignature()
	if err != nil {
		logger.Error("get block signature error: %v", err)
		return nil, false, err
	}
	r, s, err := blockSignature.ToRS()
	if err != nil {
		logger.Error("get block signature error: %v", err)
		return nil, false, err
	}
	request := rpc.BtcChangeRequest{
		Data: data,
		BlockChain: rpc.Proof{
			Proof:   blockChainProof.Proof,
			Witness: blockChainProof.Witness,
		},
		TxDepth: rpc.Proof{
			Proof:   txDepthProof.Proof,
			Witness: txDepthProof.Witness,
		},
		CpDepth: rpc.Proof{
			Proof:   cpDepthProof.Proof,
			Witness: cpDepthProof.Witness,
		},
		Redeem: rpc.Proof{
			Proof:   redeemProof.Proof,
			Witness: redeemProof.Witness,
		},
		R:          hex.EncodeToString(r),
		S:          hex.EncodeToString(s),
		ProverAddr: p.minerAddr,
	}
	return &request, false, nil
}

func (p *Prepared) GetBtcDuperGenesisRequest() (*rpc.BtcDuperRecursiveRequest, bool, error) {
	data, err := blockCu.GetRecursiveProofData(p.proverClient, uint32(p.btcGenesisHeight+common.BtcUpperDistance-1), uint32(p.btcGenesisHeight))
	if err != nil {
		logger.Error("get recursive proof data error: %v", err)
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
	request := rpc.BtcDuperRecursiveRequest{
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
	return &request, false, nil
}

func (p *Prepared) BtcDepthGenesisRequest() (*rpc.BtcDepthRecursiveRequest, bool, error) {
	data, err := blockDu.GetPackedProofData(p.proverClient, uint32(0), uint32(common.CapacityBulkUint))
	if err != nil {
		logger.Error("btc bulk data error: %v", err)
		return nil, false, err
	}
	cpHeight := uint64(0)
	firstProof, ok, err := p.filestore.GetBtcBulkProof(cpHeight, cpHeight+common.CapacityBulkUint)
	if err != nil {
		logger.Error("get btc bulk proof data error: %v %v", 0, err)
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	secondProof, ok, err := p.filestore.GetBtcBulkProof(cpHeight+common.CapacityBulkUint, cpHeight+common.CapacityBulkUint*2)
	if err != nil {
		logger.Error("get btc bulk proof data error: %v %v", 0, err)
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	request := rpc.BtcDepthRecursiveRequest{
		Data: data,
		Recursive: rpc.Proof{
			Proof:   firstProof.Proof,
			Witness: firstProof.Witness,
		},
		Unit: rpc.Proof{
			Proof:   secondProof.Proof,
			Witness: secondProof.Witness,
		},
	}
	return &request, true, nil
}

func NewPreparedData(filestore *FileStorage, store store.IStore, genesisPeriod, btcGenesisHeight uint64, proverClient *btcproverClient.Client, btcClient *btcrpc.Client,
	ethClient *ethrpc.Client, apiClient *apiclient.Client, beaconClient *beacon.Client) (*Prepared, error) {
	return &Prepared{
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
