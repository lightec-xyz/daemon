package node

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/lightec-xyz/daemon/circuits"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/rpc/beacon"
	txineth2 "github.com/lightec-xyz/provers/circuits/tx-in-eth2"
	proverType "github.com/lightec-xyz/provers/circuits/types"
	"github.com/lightec-xyz/reLight/circuits/utils"
	"github.com/prysmaticlabs/prysm/v5/api/server/structs"
	"strconv"
)

func GetRecursiveData(fileStore *FileStorage, period uint64) (interface{}, bool, error) {
	//todo
	genesisPeriod := fileStore.GetGenesisPeriod()
	genesisId, ok, err := GetSyncCommitRootId(fileStore, genesisPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v Index genesis commitId no find", genesisPeriod)
		return nil, false, nil
	}
	relayId, ok, err := GetSyncCommitRootId(fileStore, period)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v Index relay commitId no find", period)
		return nil, false, nil
	}
	endPeriod := period + 1
	endId, ok, err := GetSyncCommitRootId(fileStore, endPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v Index end commitId no find", endPeriod)
		return nil, false, nil
	}
	secondProof, exists, err := fileStore.GetUnitProof(period)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !exists {
		logger.Warn("no find %v unit proof Data, send new proof request", period)
		return nil, false, nil
	}

	prePeriod := period - 1
	firstProof, exists, err := fileStore.GetRecursiveProof(prePeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !exists {
		logger.Warn("no find %v Index recursive Data, send new proof request", prePeriod)
		return nil, false, nil
	}
	return &rpc.SyncCommRecursiveRequest{
		Choice:        "recursive",
		FirstProof:    firstProof.Proof,
		FirstWitness:  firstProof.Witness,
		SecondProof:   secondProof.Proof,
		SecondWitness: secondProof.Witness,
		BeginId:       genesisId,
		RelayId:       relayId,
		EndId:         endId,
	}, true, nil
}

func GetRecursiveGenesisData(fileStore *FileStorage, period uint64) (interface{}, bool, error) {
	genesisPeriod := fileStore.GetGenesisPeriod()
	genesisId, ok, err := GetSyncCommitRootId(fileStore, genesisPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v Index genesis commitId no find", genesisPeriod)
		return nil, false, nil
	}
	relayPeriod := period
	relayId, ok, err := GetSyncCommitRootId(fileStore, relayPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v  Index relay commitId no find ", relayPeriod)
		return nil, false, nil
	}
	endPeriod := relayPeriod + 1
	endId, ok, err := GetSyncCommitRootId(fileStore, endPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v Index end commitId no find", endPeriod)
		return nil, false, nil
	}

	fistProof, firstExists, err := fileStore.GetGenesisProof()
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !firstExists {
		logger.Warn("no find genesis proof ,start new proof request")
		return nil, false, nil
	}
	secondProof, secondExists, err := fileStore.GetUnitProof(relayPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !secondExists {
		logger.Warn("no find %v Index unit proof , send new proof request", relayPeriod)
		return nil, false, nil
	}
	return &rpc.SyncCommRecursiveRequest{
		Choice:        "genesis",
		FirstProof:    fistProof.Proof,
		FirstWitness:  fistProof.Witness,
		SecondProof:   secondProof.Proof,
		SecondWitness: secondProof.Witness,
		BeginId:       genesisId,
		RelayId:       relayId,
		EndId:         endId,
	}, true, nil

}

func GetGenesisData(fileStore *FileStorage) (*rpc.SyncCommGenesisRequest, bool, error) {
	genesisPeriod := fileStore.GetGenesisPeriod()
	genesisId, ok, err := GetSyncCommitRootId(fileStore, genesisPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v Index genesis commitId  no find", genesisPeriod)
		return nil, false, nil
	}

	nextPeriod := genesisPeriod + 1
	firstId, ok, err := GetSyncCommitRootId(fileStore, nextPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v Index first commitId no find", nextPeriod)
		return nil, false, nil
	}
	secondPeriod := nextPeriod + 1
	secondId, ok, err := GetSyncCommitRootId(fileStore, secondPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v Index second commitId no find", secondPeriod)
		return nil, false, nil
	}

	firstProof, exists, err := fileStore.GetUnitProof(genesisPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !exists {
		logger.Warn("get genesis Data,first proof not exists: %v Index", genesisPeriod)
		return nil, false, nil
	}
	logger.Info("get genesis first proof: %v", genesisPeriod)

	secondProof, exists, err := fileStore.GetUnitProof(nextPeriod)
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

func GetSyncCommitRootId(fileStore *FileStorage, period uint64) ([]byte, bool, error) {
	update, ok, err := GetSyncCommitUpdate(fileStore, period)
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

func GetSyncCommitUpdate(fileStore *FileStorage, period uint64) (*utils.SyncCommitteeUpdate, bool, error) {
	var currentPeriodUpdate structs.LightClientUpdateWithVersion
	exists, err := fileStore.GetUpdate(period, &currentPeriodUpdate)
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
	if fileStore.GetGenesisPeriod() == period {
		var genesisData structs.LightClientBootstrapResponse
		genesisExists, err := fileStore.GetBootstrap(&genesisData)
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
		if prePeriod < fileStore.GetGenesisPeriod() {
			logger.Error("should never happen: %v", prePeriod)
			return nil, false, nil
		}
		var preUpdateData structs.LightClientUpdateWithVersion
		preUpdateExists, err := fileStore.GetUpdate(prePeriod, &preUpdateData)
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
	return &update, true, nil

}

func GetBhfUpdateData(fileStore *FileStorage, finalizedSlot, genesisPeriod uint64) (*rpc.BlockHeaderFinalityRequest, bool, error) {
	logger.Debug("get bhf update data: %v", finalizedSlot)
	var currentFinalityUpdate structs.LightClientUpdateWithVersion
	exists, err := fileStore.GetFinalityUpdate(finalizedSlot, &currentFinalityUpdate)
	if err != nil {
		logger.Error("get finality update error: %v %v", finalizedSlot, err)
		return nil, false, err
	}
	if !exists {
		logger.Warn("no find finality update: %v", finalizedSlot)
		return nil, false, nil
	}

	genesisId, ok, err := GetSyncCommitRootId(fileStore, genesisPeriod)
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
	recursiveProof, ok, err := fileStore.GetRecursiveProof(period)
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
	outerProof, ok, err := fileStore.GetOuterProof(outerPeriod)
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

	currentSyncCommitUpdate, ok, err := GetSyncCommitUpdate(fileStore, outerPeriod)
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

func GetRedeemRequestData(fileStore *FileStorage, genesisPeriod, txSlot uint64, txHash string,
	beaconClient *beacon.Client, ethClient *ethclient.Client) (rpc.RedeemRequest, bool, error) {
	txProof, ok, err := fileStore.GetTxProof(txHash)
	if err != nil {
		logger.Error("get tx proof error: %v", err)
		return rpc.RedeemRequest{}, false, err
	}
	if !ok {
		logger.Debug("proof request data not prepared: %v", txHash)
		return rpc.RedeemRequest{}, false, nil
	}
	blockHeaderProof, ok, err := fileStore.GetBeaconHeaderProof(txSlot)
	if err != nil {
		logger.Error("get block header proof error: %v", err)
		return rpc.RedeemRequest{}, false, err
	}
	if !ok {
		logger.Debug("proof request data not prepared: %v", txSlot)
		return rpc.RedeemRequest{}, false, nil
	}
	finalizedSlot, ok, err := fileStore.GetNearTxSlotFinalizedSlot(txSlot)
	if err != nil {
		logger.Error("get bhf update proof error: %v", err)
		return rpc.RedeemRequest{}, false, err
	}
	if !ok {
		logger.Debug("proof request data not prepared: %v", txSlot)
		return rpc.RedeemRequest{}, false, nil
	}
	bhfProof, ok, err := fileStore.GetBhfProof(finalizedSlot)
	if err != nil {
		logger.Error("get bhf update proof error: %v", err)
		return rpc.RedeemRequest{}, false, err
	}
	if !ok {
		logger.Warn("no find bhf update %v", finalizedSlot)
		return rpc.RedeemRequest{}, false, nil
	}
	genesisRoot, ok, err := GetSyncCommitRootId(fileStore, genesisPeriod)
	if err != nil {
		logger.Error("get genesis root error: %v", err)
		return rpc.RedeemRequest{}, false, err
	}
	if !ok {
		logger.Warn("no find genesis root %v", genesisPeriod)
		return rpc.RedeemRequest{}, false, nil
	}

	var finalityUpdate *structs.LightClientUpdateWithVersion
	ok, err = fileStore.GetFinalityUpdate(finalizedSlot, &finalityUpdate)
	if err != nil {
		logger.Error("get finality update error: %v", err)
		return rpc.RedeemRequest{}, false, err
	}
	if !ok {
		logger.Warn("no find finality update %v", finalizedSlot)
		return rpc.RedeemRequest{}, false, nil
	}

	attestedSlot, err := strconv.ParseUint(finalityUpdate.Data.AttestedHeader.Slot, 10, 64)
	if err != nil {
		logger.Error("parse slot error: %v", err)
		return rpc.RedeemRequest{}, false, err
	}
	period := attestedSlot / common.SlotPerPeriod
	currentRoot, ok, err := GetSyncCommitRootId(fileStore, period)
	if err != nil {
		logger.Error("get current root error: %v", err)
		return rpc.RedeemRequest{}, false, err
	}
	beginID, endId, err := GetBeaconHeaderId(beaconClient, txSlot, finalizedSlot)
	if err != nil {
		logger.Error("get begin and end id error: %v", err)
		return rpc.RedeemRequest{}, false, err
	}
	// todo need cache
	txVar, receiptVar, err := txineth2.GenerateTxAndReceiptU128Padded(ethClient, txHash)
	if err != nil {
		logger.Error("get tx and receipt error: %v", err)
		return rpc.RedeemRequest{}, false, err
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
		TxVar:            txVar,
		ReceiptVar:       receiptVar,
	}
	return redeemRequest, true, nil

}

func GetBeaconHeaderId(beaconClient *beacon.Client, start, end uint64) ([]byte, []byte, error) {
	beaconBlockHeaders, err := beaconClient.RetrieveBeaconHeaders(start, end)
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
