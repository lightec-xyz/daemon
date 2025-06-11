package node

import (
	"context"
	"encoding/hex"
	"fmt"
	"strconv"

	ethcommon "github.com/ethereum/go-ethereum/common"

	baselevelUtil "github.com/lightec-xyz/btc_provers/utils/blockchain"
	blockchainUtil "github.com/lightec-xyz/btc_provers/utils/blockchain"
	blockdepthUtil "github.com/lightec-xyz/btc_provers/utils/blockdepth"
	btcproverClient "github.com/lightec-xyz/btc_provers/utils/client"
	txinchainUtil "github.com/lightec-xyz/btc_provers/utils/txinchain"
	"github.com/lightec-xyz/daemon/circuits"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/rpc/beacon"
	btcrpc "github.com/lightec-xyz/daemon/rpc/bitcoin"
	ethrpc "github.com/lightec-xyz/daemon/rpc/ethereum"
	"github.com/lightec-xyz/daemon/store"
	proverType "github.com/lightec-xyz/provers/circuits/types"
	apiclient "github.com/lightec-xyz/provers/utils/api-client"
	redeemUtils "github.com/lightec-xyz/provers/utils/redeem-tx"
	txineth2Utils "github.com/lightec-xyz/provers/utils/tx-in-eth2"
)

type Prepared struct {
	filestore        *FileStorage
	chainStore       *ChainStore
	proverClient     btcproverClient.IClient
	btcClient        *btcrpc.Client
	ethClient        *ethrpc.Client
	apiClient        *apiclient.Client
	beaconClient     *beacon.Client
	genesisPeriod    uint64
	genesisSlot      uint64
	minerAddr        string
	btcGenesisHeight uint64 // startIndex
}

func (p *Prepared) GetBtcBaseRequest(start, end uint64) (*rpc.BtcBaseRequest, bool, error) {
	logger.Debug("base:%v~%v", start, end)
	data, err := baselevelUtil.GetBaseLevelProofData(p.proverClient, uint32(end-1))
	if err != nil {
		logger.Error("get base level proof data error: %v~%v %v", start, end, err)
		return nil, false, err
	}
	baseRequest := rpc.BtcBaseRequest{
		Data: data,
	}
	return &baseRequest, true, nil
}

func (p *Prepared) GetBtcMiddleRequest(start, end uint64) (*rpc.BtcMiddleRequest, bool, error) {
	logger.Debug("middle:%v~%v", start, end)

	var proofs []rpc.Proof
	for index := start; index <= end-common.BtcBaseDistance; index = index + common.BtcBaseDistance {
		baseStartIndex := index
		baseEndIndex := index + common.BtcBaseDistance
		baseProof, ok, err := p.filestore.GetBtcBaseProof(baseStartIndex, baseEndIndex)
		if err != nil {
			logger.Error("get base level proof data error: %v~%v %v", baseStartIndex, baseEndIndex, err)
			return nil, false, err
		}
		if ok {
			proofs = append(proofs, rpc.Proof{
				Proof:   baseProof.Proof,
				Witness: baseProof.Witness,
			})
		}
	}
	needProofLen := int(end-start) / common.BtcBaseDistance
	if len(proofs) != needProofLen {
		logger.Warn("get middle base proof not match: %v~%v, %v %v", start, end, needProofLen, len(proofs))
		return nil, false, nil
	}
	data, err := blockchainUtil.GetMidLevelProofData(p.proverClient, uint32(end-1))
	if err != nil {
		logger.Error("get base level proof data error: %v %v", start, err)
		return nil, false, err
	}
	logger.Debug("middle:%v~%v,baseLen: %v", start, end, len(proofs))
	baseRequest := rpc.BtcMiddleRequest{
		Data:   data,
		Proofs: proofs,
	}
	return &baseRequest, true, nil
}

func (p *Prepared) GetBtcUpperRequest(start, end uint64) (*rpc.BtcUpperRequest, bool, error) {
	logger.Debug("upper:%v~%v", start, end)
	var proofs []rpc.Proof
	for index := start; index <= end-common.BtcMiddleDistance; index = index + common.BtcMiddleDistance {
		middleStartIndex := index
		middleEndIndex := index + common.BtcMiddleDistance
		middleProof, ok, err := p.filestore.GetBtcMiddleProof(middleStartIndex, middleEndIndex)
		if err != nil {
			logger.Error("get base level proof data error: %v~%v %v", middleStartIndex, middleEndIndex, err)
			return nil, false, err
		}
		if ok {
			proofs = append(proofs, rpc.Proof{
				Proof:   middleProof.Proof,
				Witness: middleProof.Witness,
			})
		}
	}
	needProofLen := int(end-start) / common.BtcMiddleDistance
	if len(proofs) != needProofLen {
		logger.Warn("get upper middle level proof not match: %v~%v, %v  %v", start, end, needProofLen, len(proofs))
		return nil, false, nil
	}
	logger.Debug("upper:%v~%v,  middleLen: %v", start, end, len(proofs))
	data, err := blockchainUtil.GetUpperLevelProofData(p.proverClient, uint32(end-1))
	if err != nil {
		logger.Error("get upper level proof data error: %v~%v %v", start, end, err)
		return nil, false, err
	}
	baseRequest := rpc.BtcUpperRequest{
		Data:   data,
		Proofs: proofs,
	}
	return &baseRequest, true, nil
}

func (p *Prepared) getProveType(step uint64) string {
	if step == common.BtcUpperDistance || step == common.BtcBaseDistance {
		return circuits.BtcBlockChain
	}
	return circuits.BtcChainHybrid

}

func (p *Prepared) GetBtcDuperRecursiveRequest(start, end uint64) (*rpc.BtcDuperRecursiveRequest, bool, error) {
	var first, second *StoreProof
	var ok bool
	var err error
	isBlockChain := true
	var blockProofData *blockchainUtil.BlockChainProofData
	var hybridChainData *blockchainUtil.HybridProofData
	var firstType, secondType string
	var firstStep, secondStep uint64
	step := end - start
	if step == common.BtcUpperDistance {
		firstProof, exists, errC := p.filestore.FindBtcChainProof(start)
		if errC != nil {
			logger.Error("get btc genesis proof data error: %v %v", 0, errC)
			return nil, false, errC
		}
		if !exists {
			return nil, false, nil
		}
		if firstProof.Start == p.btcGenesisHeight {
			firstType = circuits.BtcChainUpper
		} else {
			firstType = p.getProveType(firstProof.Step)
		}

		firstStep = firstProof.Step
		first = &firstProof.StoreProof

		secondType = circuits.BtcChainUpper
		secondStep = common.BtcUpperDistance
		second, ok, err = p.filestore.GetBtcUpperProof(start, end)
		if err != nil {
			logger.Error("get btc genesis proof data error: %v %v", 0, err)
			return nil, false, err
		}
		if !ok {
			return nil, false, nil
		}
	} else if step == common.BtcBaseDistance {
		firstProof, exists, errC := p.filestore.FindBtcChainProof(start)
		if errC != nil {
			logger.Error("get btc genesis proof data error: %v %v", 0, errC)
			return nil, false, errC
		}
		if !exists {
			return nil, false, nil
		}
		firstType = p.getProveType(firstProof.Step)
		firstStep = firstProof.Step
		first = &firstProof.StoreProof

		secondType = circuits.BtcChainBase
		secondStep = common.BtcBaseDistance
		second, ok, err = p.filestore.GetBtcBaseProof(start, end)
		if err != nil {
			logger.Error("get btc genesis proof data error: %v %v", 0, err)
			return nil, false, err
		}
		if !ok {
			return nil, false, nil
		}
	} else {
		isBlockChain = false
		firstProof, exists, errC := p.filestore.FindBtcChainProof(start)
		if errC != nil {
			logger.Error("get btc genesis proof data error: %v %v", 0, errC)
			return nil, false, errC
		}
		if !exists {
			return nil, false, nil
		}
		firstStep = firstProof.Step
		firstType = p.getProveType(firstStep)
		first = &firstProof.StoreProof
	}
	if isBlockChain {
		blockProofData, err = blockchainUtil.GetBlockChainProofData(p.proverClient, uint32(end-1), uint32(p.btcGenesisHeight), uint32(start))
		if err != nil {
			logger.Error("get recursive proof data error: %v", err)
			return nil, false, err
		}
		err = blockProofData.AdjustFirstBlockTime(p.proverClient)
		if err != nil {
			logger.Error("adjust first block time error: %v", err)
			return nil, false, err
		}
	} else {
		hybridChainData, err = blockchainUtil.GetHybridProofData(p.proverClient, uint32(start), uint32(end-1), uint32(p.btcGenesisHeight))
		if err != nil {
			logger.Error("get recursive proof data error: %v", err)
			return nil, false, err
		}
		err = hybridChainData.AdjustFirstBlockTime(p.proverClient)
		if err != nil {
			logger.Error("adjust first block time error: %v", err)
			return nil, false, err
		}
	}
	request := rpc.BtcDuperRecursiveRequest{
		BlockChainData:  blockProofData,
		HybridChainData: hybridChainData,
		First: rpc.Proof{
			Proof:   first.Proof,
			Witness: first.Witness,
		},
		Start:      start,
		End:        end,
		FirstStep:  firstStep,
		SecondStep: secondStep,
		FirstType:  firstType,
		SecondType: secondType,
	}
	if second != nil {
		request.Second = rpc.Proof{
			Proof:   second.Proof,
			Witness: second.Witness,
		}
	}
	return &request, true, nil
}

func (p *Prepared) GetTxInEth2Request(txHash string, getSlotByNumber func(uint64) (uint64, error)) (*rpc.TxInEth2ProveRequest, bool, error) {
	txData, err := txineth2Utils.GetTxInEth2ProofData(p.ethClient.Client, p.apiClient, getSlotByNumber, ethcommon.HexToHash(txHash))
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
	finalizedSlot, ok, err := p.filestore.GetTxFinalizedSlot(index)
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
	var middleHeader []*proverType.BeaconHeader
	err = common.ParseObj(beaconBlockHeaders[1:], &middleHeader)
	if err != nil {
		logger.Error("deep copy error %v", err)
		return nil, false, err
	}
	return &rpc.BlockHeaderRequest{
		Data: &proverType.BeaconHeaderChain{
			BeginSlot:           beginSlot,
			EndSlot:             endSlot,
			BeginRoot:           hex.EncodeToString(beginRoot),
			EndRoot:             hex.EncodeToString(endRoot),
			MiddleBeaconHeaders: middleHeader,
		},
		Index: index,
	}, true, nil
}

func (p *Prepared) GetDutyRequest(period uint64) (*rpc.SyncCommDutyRequest, bool, error) {
	genesisPeriod := p.filestore.GetGenesisPeriod()
	if period < genesisPeriod+1 {
		return nil, false, fmt.Errorf(" recursive less than %v", genesisPeriod+1)
	}
	genesisId, ok, err := p.GetSyncCommitRootId(genesisPeriod)
	if err != nil {
		logger.Error("get genesis id error: %v", err)
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v FIndex genesis commitId no find", genesisPeriod)
		return nil, false, nil
	}
	relayId, ok, err := p.GetSyncCommitRootId(period)
	if err != nil {
		logger.Error("get relay id error: %v", err)
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v FIndex relay commitId no find", period)
		return nil, false, nil
	}
	var firstProof *StoreProof
	var exists bool
	var choice string
	if period == genesisPeriod+1 {
		firstProof, exists, err = p.filestore.GetUnitProof(p.genesisPeriod)
		choice = circuits.SyncCommitteeGenesis
	} else {
		prePeriod := period - 1
		firstProof, exists, err = p.filestore.GetRecursiveProof(prePeriod)
		choice = circuits.SyncCommitteeRecursive
	}
	if err != nil {
		logger.Error("get first proof error: %v", err)
		return nil, false, err
	}
	if !exists {
		return nil, false, nil
	}
	endPeriod := period + 1
	endId, ok, err := p.GetSyncCommitRootId(endPeriod)
	if err != nil {
		logger.Error("get end id error: %v", err)
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	secondProof, exists, err := p.filestore.GetUnitProof(period)
	if err != nil {
		logger.Error("get second proof error: %v", err)
		return nil, false, err
	}
	if !exists {
		logger.Warn("no find %v unit proof Responses, send new proof request", period)
		return nil, false, nil
	}
	outerProof, exists, err := p.filestore.GetOuterProof(endPeriod)
	if err != nil {
		logger.Error("get outer proof error: %v", err)
		return nil, false, err
	}
	if !exists {
		return nil, false, nil
	}
	update, exists, err := p.GetSyncCommitUpdate(endPeriod)
	if err != nil {
		logger.Error("get update error: %v", err)
		return nil, false, err
	}
	if !exists {
		return nil, false, nil
	}
	return &rpc.SyncCommDutyRequest{
		Choice: choice,
		Period: period,
		FirstProof: rpc.Proof{
			Proof:   firstProof.Proof,
			Witness: firstProof.Witness},
		SecondProof: rpc.Proof{
			Proof:   secondProof.Proof,
			Witness: secondProof.Witness,
		},
		Outer: rpc.Proof{
			Proof:   outerProof.Proof,
			Witness: outerProof.Witness,
		},
		BeginId: hex.EncodeToString(genesisId),
		RelayId: hex.EncodeToString(relayId),
		EndId:   hex.EncodeToString(endId),
		ScIndex: int(period), //todo
		Update:  update.SyncCommitteeUpdate,
	}, true, nil
}

func (p *Prepared) getSlotByNumber(number uint64) (uint64, error) {
	slot, ok, err := p.chainStore.ReadSlotByHeight(number)
	if err != nil {
		logger.Error("get slot error: %v %v", number, err)
		return 0, err
	}
	if !ok {
		return 0, fmt.Errorf("no find %v slot", number)
	}
	return slot, nil
}

func (p *Prepared) GetSyncCommitRootId(period uint64) ([]byte, bool, error) {
	update, ok, err := p.GetSyncCommitUpdate(period)
	if err != nil {
		logger.Error("get %v FIndex update error: %v", period, err)
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	syncCommitRoot, err := circuits.SyncCommitRoot(update.CurrentSyncCommittee)
	if err != nil {
		logger.Error("get commit root: %v error: %v", period, err)
		return nil, false, err
	}
	return syncCommitRoot, true, nil
}

func (p *Prepared) GetSyncComInnerRequest(period, index uint64) (*rpc.SyncCommInnerRequest, bool, error) {
	syncCommittee, exists, err := p.GetSyncCommittee(period)
	if err != nil {
		logger.Error("get %v syncCommittee: %v", period, err)
		return nil, false, err
	}
	if !exists {
		return nil, false, nil
	}
	return &rpc.SyncCommInnerRequest{
		Data:    syncCommittee.SyncCommittee,
		Index:   index,
		Period:  period,
		Version: syncCommittee.Version,
	}, true, nil

}

func (p *Prepared) GetSyncOuterRequest(period uint64) (*rpc.SyncCommOuterRequest, bool, error) {
	syncCommittee, exists, err := p.GetSyncCommittee(period)
	if err != nil {
		logger.Error("get %v syncCommittee: %v", period, err)
		return nil, false, err
	}
	if !exists {
		return nil, false, nil
	}
	var innerProofs []rpc.Proof
	for index := 0; index < common.SyncInnerNum; index++ {
		proof, exists, err := p.filestore.GetSyncInnerProof(period, uint64(index))
		if err != nil {
			logger.Error("get %v inner proof: %v", period, err)
			return nil, false, err
		}
		if !exists {
			return nil, false, nil
		}
		innerProofs = append(innerProofs, rpc.Proof{
			Proof:   proof.Proof,
			Witness: proof.Witness,
		})
	}
	if len(innerProofs) != common.SyncInnerNum {
		logger.Warn("sync outer period: %v, inner proof num: %v", period, len(innerProofs))
		return nil, false, nil
	}
	return &rpc.SyncCommOuterRequest{
		Period:      period,
		Data:        syncCommittee.SyncCommittee,
		Version:     syncCommittee.Version,
		InnerProofs: innerProofs,
	}, true, nil
}

func (p *Prepared) GetSyncComUnitRequest(period uint64) (*rpc.SyncCommUnitsRequest, bool, error) {
	update, ok, err := p.GetSyncCommitUpdate(period)
	if err != nil {
		logger.Error("get %v FIndex update error: %v", period, err)
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	outerProof, exists, err := p.filestore.GetOuterProof(period)
	if err != nil {
		logger.Error("get %v outer proof error: %v", period, err)
		return nil, false, err
	}
	if !exists {
		return nil, false, nil
	}
	commUnitsRequest := rpc.SyncCommUnitsRequest{
		Data:    update,
		Index:   period,
		Version: update.Version,
		Outer: rpc.Proof{
			Proof:   outerProof.Proof,
			Witness: outerProof.Witness,
		},
	}
	return &commUnitsRequest, true, nil

}

func (p *Prepared) GetBtcBulkRequest(start, end, prefix uint64) (*rpc.BtcBulkRequest, error) {
	proofData, err := blockdepthUtil.GetBlockBulkProofData(p.proverClient, uint32(start), uint32(end))
	if err != nil {
		logger.Error("btc bulk data error: %v", err)
		return nil, err
	}
	return &rpc.BtcBulkRequest{
		Data: proofData,
	}, nil

}

func (p *Prepared) GetSyncCommittee(period uint64) (*WrapSyncCommittee, bool, error) {
	if period == p.genesisPeriod {
		var bootstrap common.BootstrapResponse
		exists, err := p.filestore.GetBootStrapBySlot(p.genesisSlot, &bootstrap)
		if err != nil {
			logger.Error("get bootstrap error: %v", err)
			return nil, false, err
		}
		if !exists {
			return nil, false, nil
		}
		return &WrapSyncCommittee{
			SyncCommittee: &proverType.SyncCommittee{
				PubKeys:         bootstrap.Data.CurrentSyncCommittee.Pubkeys,
				AggregatePubKey: bootstrap.Data.CurrentSyncCommittee.AggregatePubkey,
			},
			Version: bootstrap.Version,
		}, true, nil

	} else {
		var lightClientUpdate common.LightClientUpdateResponse
		exists, err := p.filestore.GetUpdate(period-1, &lightClientUpdate)
		if err != nil {
			logger.Error("get %v index update error: %v", period, err)
			return nil, false, err
		}
		if !exists {
			return nil, false, nil
		}
		return &WrapSyncCommittee{
			SyncCommittee: &proverType.SyncCommittee{
				PubKeys:         lightClientUpdate.Data.NextSyncCommittee.Pubkeys,
				AggregatePubKey: lightClientUpdate.Data.NextSyncCommittee.AggregatePubkey,
			},
			Version: lightClientUpdate.Version,
		}, true, nil
	}

}

func (p *Prepared) GetSyncCommitUpdate(period uint64) (*rpc.WrapSyncCommitteeUpdate, bool, error) {
	var currentPeriodUpdate common.LightClientUpdateResponse
	exists, err := p.filestore.GetUpdate(period, &currentPeriodUpdate)
	if err != nil {
		logger.Error("get %v index update error: %v", period, err)
		return nil, false, err
	}
	if !exists {
		//logger.Warn("no find %v index update Responses", period)
		return nil, false, nil
	}
	update, err := dbUpdateToZkUpdate(&currentPeriodUpdate)
	if err != nil {
		logger.Error("parse obj error: %v %v", period, err)
		return nil, false, err
	}
	if p.genesisPeriod == period {
		var bootstrap common.BootstrapResponse
		genesisExists, err := p.filestore.GetBootStrapBySlot(p.genesisSlot, &bootstrap)
		if err != nil {
			logger.Error("get genesis update error: %v %v", period, err)
			return nil, false, err
		}
		if !genesisExists {
			logger.Warn("no find genesis update Responses,%v", period)
			return nil, false, nil
		}
		update.CurrentSyncCommittee = &proverType.SyncCommittee{
			PubKeys:         bootstrap.Data.CurrentSyncCommittee.Pubkeys,
			AggregatePubKey: bootstrap.Data.CurrentSyncCommittee.AggregatePubkey,
		}
		update.CurrentSyncCommitteeBranch = bootstrap.Data.CurrentSyncCommitteeBranch
	} else {
		prePeriod := period - 1
		if prePeriod < p.filestore.GetGenesisPeriod() {
			logger.Error("should never happen: %v", prePeriod)
			return nil, false, nil
		}
		var preUpdateData common.LightClientUpdateResponse
		preUpdateExists, err := p.filestore.GetUpdate(prePeriod, &preUpdateData)
		if err != nil {
			logger.Error("get %v index update error: %v", prePeriod, err)
			return nil, false, err
		}
		if !preUpdateExists {
			logger.Warn("get unit Responses,no find %v FIndex update Responses", prePeriod)
			return nil, false, nil
		}
		update.CurrentSyncCommittee = &proverType.SyncCommittee{
			PubKeys:         preUpdateData.Data.NextSyncCommittee.Pubkeys,
			AggregatePubKey: preUpdateData.Data.NextSyncCommittee.AggregatePubkey,
		}
		update.CurrentSyncCommitteeBranch = preUpdateData.Data.NextSyncCommitteeBranch
	}
	ok, err := update.SyncCommitteeUpdate.Verify()
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

func (p *Prepared) GetBhfUpdateRequest(finalizedSlot uint64) (*rpc.BlockHeaderFinalityRequest, bool, error) {
	logger.Debug("get bhf update data: %v", finalizedSlot)
	var currentFinalityUpdate common.LightClientFinalityUpdateEvent
	exists, err := p.filestore.GetFinalityUpdate(finalizedSlot, &currentFinalityUpdate)
	if err != nil {
		logger.Error("get finality update error: %v %v", finalizedSlot, err)
		return nil, false, err
	}
	if !exists {
		logger.Warn("no find finality update: %v", finalizedSlot)
		return nil, false, nil
	}
	attestedSlot, err := strconv.ParseUint(currentFinalityUpdate.Data.AttestedHeader.Slot, 10, 64)
	if err != nil {
		logger.Error("parse big error %v %v", currentFinalityUpdate.Data.AttestedHeader.Slot, err)
		return nil, false, err
	}
	period := attestedSlot / common.SlotPerPeriod
	syncCommittee, exists, err := p.GetSyncCommittee(period)
	if err != nil {
		logger.Error("get %v syncCommittee: %v", period, err)
		return nil, false, err
	}
	if !exists {
		return nil, false, nil
	}
	proversFinalityUpdate := dbFinalityUpdateToZkFinalityUpdate(&currentFinalityUpdate)
	request := rpc.BlockHeaderFinalityRequest{
		Index:          finalizedSlot,
		FinalityUpdate: proversFinalityUpdate,
		SyncCommittee:  syncCommittee.SyncCommittee,
	}
	return &request, true, nil
}

func (p *Prepared) GetRedeemRequest(txHash string) (*rpc.RedeemRequest, bool, error) {

	txSlot, ok, err := p.chainStore.ReadSlotByHash(txHash)
	if err != nil {
		logger.Error("get txSlot error: %v %v", err, txHash)
		return nil, false, err
	}
	if !ok {
		logger.Warn("no find  tx %v beacon slot", txHash)
		return nil, false, nil
	}

	txProof, ok, err := p.filestore.GetTxProof(txHash)
	if err != nil {
		logger.Error("get tx proof error: %v", err)
		return nil, false, err
	}
	if !ok {
		logger.Debug("proof request data not prepared: %v", txHash)
		return nil, false, nil
	}
	// todo
	finalizedSlot, ok, err := p.filestore.GetTxFinalizedSlot(txSlot)
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
	genesisRoot, ok, err := p.GetSyncCommitRootId(p.genesisPeriod)
	if err != nil {
		logger.Error("get genesis root error: %v", err)
		return nil, false, err
	}
	if !ok {
		logger.Warn("no find genesis root %v", p.genesisPeriod)
		return nil, false, nil
	}

	var finalityUpdate common.LightClientFinalityUpdateEvent
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
	if !ok {
		logger.Warn("no find current root %v", period)
		return nil, false, nil
	}

	dutyProof, exists, err := p.filestore.GetDutyProof(period - 1)
	if err != nil {
		logger.Error("get recursive proof error: %v", err)
		return nil, false, err
	}
	if !exists {
		logger.Warn("no find recursive proof %v", period)
		return nil, false, nil
	}

	//todo slot
	beaconHeaders, err := p.beaconClient.RetrieveBeaconHeaders(txSlot, finalizedSlot)
	if err != nil {
		logger.Error("get beacon headers error: %v", err)
		return nil, false, err
	}
	if len(beaconHeaders) == 0 {
		return nil, false, fmt.Errorf("%v-%v beaconHeaer length is 0", txSlot, finalizedSlot)
	}
	receipt, err := p.ethClient.TransactionReceipt(context.Background(), ethcommon.HexToHash(txHash))
	if err != nil {
		logger.Error("get eth tx receipt error:%v %v", txHash, err)
		return nil, false, err
	}
	btcTxId, minerReward, _, sigHashes, _, err := redeemUtils.DecodeRedeemReceipt(receipt)
	if err != nil {
		logger.Error("get tx receipt proof data error:%v %v", txHash, err)
		return nil, false, err
	}

	logger.Debug("Redeem btcTxid: %x,minerReward: %x,sigHashes: %x", btcTxId, minerReward, sigHashes)

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
		Duty: rpc.Proof{
			Proof:   dutyProof.Proof,
			Witness: dutyProof.Witness,
		},
		GenesisScRoot:    hex.EncodeToString(genesisRoot),
		CurrentSCSSZRoot: hex.EncodeToString(currentRoot),
		SigHashes:        common.BytesArrayToHex(sigHashes),
		NbBeaconHeaders:  len(beaconHeaders) - 1,
		MinerReward:      hex.EncodeToString(minerReward[:]),
		TxId:             hex.EncodeToString(btcTxId[:]),
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

func (p *Prepared) GetBtcDepthRecursiveRequest(prefix, start, end uint64, isCp bool) (*rpc.BtcDepthRecursiveRequest, bool, error) {
	var first *StoreProof
	var ok bool
	var err error
	var step uint64
	var isRecursive bool
	if start-prefix == common.BtcTxUnitMaxDepth || (start-prefix == common.BtcCpMinDepth && isCp) { // todo
		first, ok, err = p.filestore.GetBtcBulkProof(prefix, start)
		if err != nil {
			logger.Error("get btc chain first proof data error: %v %v", start, err)
			return nil, false, err
		}
		if !ok {
			return nil, false, nil
		}
		step = start - prefix
	} else {
		firstProof, exists, err := p.filestore.FindDepthProof(prefix, start)
		if err != nil {
			logger.Error("get btc chain first proof data error: %v %v", start, err)
			return nil, false, err
		}
		if !exists {
			return nil, false, nil
		}
		first = &firstProof.StoreProof
		step = firstProof.Step
		isRecursive = true
	}
	proofData, err := blockdepthUtil.GetRecursiveBulksProofData(p.proverClient, uint32(prefix), uint32(start), uint32(end))
	if err != nil {
		logger.Error("btc bulk data error: %v", err)
		return nil, false, err
	}
	request := rpc.BtcDepthRecursiveRequest{
		Data: proofData,
		First: rpc.Proof{
			Proof:   first.Proof,
			Witness: first.Witness,
		},
		Genesis:     prefix,
		Start:       start,
		End:         end,
		PreStep:     step,
		IsRecursive: isRecursive,
	}
	return &request, true, nil
}

func (p *Prepared) getTxDepthProof(depthHeight, latestHeight uint64) (*StoreProof, uint64, bool, error) {
	logger.Debug("get depth: %v  %v", depthHeight, latestHeight)
	step := latestHeight - depthHeight
	if step >= common.BtcTxMinDepth && step <= common.BtcTxUnitMaxDepth {
		storeProof, exists, err := p.filestore.GetBtcBulkProof(depthHeight, latestHeight)
		if err != nil {
			logger.Error("get depth proof error: %v", err)
			return nil, 0, false, err
		}
		if !exists {
			return nil, 0, false, nil
		}
		return storeProof, latestHeight - depthHeight, exists, nil
	} else if step > common.BtcTxUnitMaxDepth {
		storageProof, exists, err := p.filestore.FindDepthProof(depthHeight, latestHeight)
		if err != nil {
			logger.Error("find depth proof error: %v", err)
			return nil, 0, false, err
		}
		if !exists {
			return nil, 0, false, nil
		}
		return &storageProof.StoreProof, storageProof.Step, exists, nil
	} else {
		return nil, 0, false, fmt.Errorf("never should happen:%v %v", depthHeight, latestHeight)
	}
}

func (p *Prepared) GetBtcDepositRequest(hash string) (*rpc.BtcDepositRequest, bool, error) {
	dbTx, ok, err := p.chainStore.ReadBtcTx(hash)
	if err != nil {
		logger.Error("read btc tx error: %v", err)
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	blockChainProof, exists, err := p.filestore.FindBtcChainProof(dbTx.LatestHeight + 1) //todo
	if err != nil {
		logger.Error("get btc blockHash chain proof data error: %v %v", 0, err)
		return nil, false, err
	}
	if !exists {
		return nil, false, nil
	}
	logger.Debug("btc tx: %v,chain proof: %v", hash, dbTx.LatestHeight+1)

	chainProveType := p.getProveType(blockChainProof.Step)

	txRecursive := dbTx.LatestHeight-dbTx.Height > common.BtcTxUnitMaxDepth

	txDepthProof, txStep, ok, err := p.getTxDepthProof(dbTx.Height, dbTx.LatestHeight)
	if err != nil {
		logger.Error("get btc tx depth proof data error: %v %v", 0, err)
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}

	cpRecursive := dbTx.LatestHeight-dbTx.CheckPointHeight > common.BtcCpMinDepth
	cpDepthProof, cpStep, ok, err := p.getDepthProof(common.BtcCpMinDepth, dbTx.CheckPointHeight, dbTx.LatestHeight)
	if err != nil {
		logger.Error("get btc cp depth proof data error: %v %v", 0, err)
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}

	timestampProof, ok, err := p.filestore.GetBtcTimestampProof(dbTx.Height, dbTx.LatestHeight)
	if err != nil {
		logger.Error("get btc timestamp proof data error: %v %v %v", dbTx.Height, dbTx.LatestHeight, err)
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}

	icpSignature, ok, err := p.chainStore.ReadIcpSignature(dbTx.LatestHeight)
	if err != nil {
		logger.Error("read dfinity sign error: %v", err)
		return nil, false, err
	}
	if !ok {
		logger.Warn("no find: %v icp %v signature", hash, dbTx.LatestHeight)
		// no work,just placeholder
		icpSignature.Hash = "6aeb6ec6f0fbc707b91a3bec690ae6536fe0abaa1994ef24c3463eb20494785d"
		icpSignature.Signature = "3f8e02c743e76a4bd655873a428db4fa2c46ac658854ba38f8be0fbbf9af9b2b6b377aaaaf231b6b890a5ee3c15a558f1ccc18dae0c844b6f06343b88a8d12e3"
	}
	blockHash, err := p.btcClient.GetBlockHash(int64(dbTx.Height))
	if err != nil {
		logger.Error("get blockHash by number error: %v", err)
		return nil, false, err
	}

	logger.Debug("txHash: %v, blockHash: %v,genesisHeight: %v, txHeight: %v,latestHeight: %v, cpHeight: %v",
		hash, blockHash, p.btcGenesisHeight, dbTx.Height, dbTx.LatestHeight, dbTx.CheckPointHeight)
	data, err := txinchainUtil.GetTxInChainProofData(p.proverClient, hash, blockHash, uint32(dbTx.LatestHeight), uint32(dbTx.CheckPointHeight),
		uint32(p.btcGenesisHeight))
	if err != nil {
		logger.Error("get tx in chain proof data error: %v", err)
		return nil, false, err
	}
	smoothedTimestampProofData, err := blockdepthUtil.GetSmoothedTimestampProofData(p.proverClient, uint32(dbTx.LatestHeight))
	if err != nil {
		logger.Error("get timestamp proof data error: %v", err)
		return nil, false, err
	}

	cptimeData, err := blockdepthUtil.GetCpTimestampProofData(p.proverClient, uint32(dbTx.Height))
	if err != nil {
		logger.Error("%v", err)
		return nil, false, err
	}

	sigVerifyData, err := blockdepthUtil.GetSigVerifProofData(
		ethcommon.FromHex(icpSignature.Hash),
		ethcommon.FromHex(icpSignature.Signature),
		ethcommon.FromHex(TestnetIcpPublicKey))
	if err != nil {
		logger.Error("get sig verif proof data error: %v", err)
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
		SigVerify: rpc.Proof{
			Proof:   timestampProof.Proof,
			Witness: timestampProof.Witness,
		},
		ProverAddr:        p.minerAddr,
		ChainType:         chainProveType,
		ChainStep:         blockChainProof.Step,
		TxDepthStep:       txStep,
		CpDepthStep:       cpStep,
		TxRecursive:       txRecursive,
		CpRecursive:       cpRecursive,
		CpFlag:            cptimeData.Flag,
		SmoothedTimestamp: smoothedTimestampProofData.Timestamp,
		SigVerifyData:     sigVerifyData,
	}
	return &request, true, nil
}

func (p *Prepared) getDepthProof(genesisCount, depthHeight, latestHeight uint64) (*StoreProof, uint64, bool, error) {
	logger.Debug("get depth: %v  %v", depthHeight, latestHeight)
	if latestHeight-depthHeight <= genesisCount {
		storeProof, exists, err := p.filestore.GetBtcBulkProof(depthHeight, latestHeight)
		if err != nil {
			logger.Error("get depth proof error: %v", err)
			return nil, 0, false, err
		}
		if !exists {
			return nil, 0, false, nil
		}
		return storeProof, latestHeight - depthHeight, exists, nil
	} else {
		storageProof, exists, err := p.filestore.FindDepthProof(depthHeight, latestHeight)
		if err != nil {
			logger.Error("find depth proof error: %v", err)
			return nil, 0, false, err
		}
		if !exists {
			return nil, 0, false, nil
		}
		return &storageProof.StoreProof, storageProof.Step, exists, nil
	}
}

func (p *Prepared) GetBtcChangeRequest(hash string) (*rpc.BtcChangeRequest, bool, error) {
	depositRequest, ok, err := p.GetBtcDepositRequest(hash)
	if err != nil {
		logger.Error("get btc deposit request error: %v", err)
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	destHash, err := p.chainStore.ReadDestHash(hash)
	if err != nil {
		logger.Error("read dest hash error: %v", err)
		return nil, false, err
	}
	backendRedeemProof, ok, err := p.filestore.GetBackendRedeemProof(destHash)
	if err != nil {
		logger.Error("get backend Redeem proof data error: %v %v", 0, err)
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	receipt, err := p.ethClient.TransactionReceipt(context.Background(), ethcommon.HexToHash(destHash))
	if err != nil {
		logger.Error("get transaction receipt error: %v", err)
		return nil, false, err
	}
	_, rewardBytes, _, _, _, err := redeemUtils.DecodeRedeemReceipt(receipt)
	if err != nil {
		logger.Error("decode Redeem log error:%v %v", hash, err)
		return nil, false, err
	}
	request := rpc.BtcChangeRequest{
		Data:       depositRequest.Data,
		BlockChain: depositRequest.BlockChain,
		TxDepth:    depositRequest.TxDepth,
		CpDepth:    depositRequest.CpDepth,
		SigVerify:  depositRequest.SigVerify,
		Redeem: rpc.Proof{
			Proof:   backendRedeemProof.Proof,
			Witness: backendRedeemProof.Witness,
		},
		ProverAddr:        depositRequest.ProverAddr,
		ChainStep:         depositRequest.ChainStep,
		TxDepthStep:       depositRequest.TxDepthStep,
		CpDepthStep:       depositRequest.CpDepthStep,
		TxRecursive:       depositRequest.TxRecursive,
		CpRecursive:       depositRequest.CpRecursive,
		ChainType:         depositRequest.ChainType,
		MinerReward:       hex.EncodeToString(rewardBytes[:]),
		CpFlag:            depositRequest.CpFlag,
		SmoothedTimestamp: depositRequest.SmoothedTimestamp,
		SigVerifyData:     depositRequest.SigVerifyData,
	}
	return &request, true, nil
}

func (p *Prepared) GetBtcTimestampRequest(fIndex uint64, sIndex uint64) (*rpc.BtcTimestampRequest, bool, error) {
	cpData, err := blockdepthUtil.GetCpTimestampProofData(p.proverClient, uint32(fIndex))
	if err != nil {
		logger.Error("get timestamp proof data error: %v", err)
		return nil, false, err
	}
	smoothedTimestampProofData, err := blockdepthUtil.GetSmoothedTimestampProofData(p.proverClient, uint32(sIndex))
	if err != nil {
		logger.Error("get timestamp proof data error: %v", err)
		return nil, false, err
	}
	request := &rpc.BtcTimestampRequest{
		CpTime:     cpData,
		SmoothData: smoothedTimestampProofData,
	}
	return request, true, nil

}

func NewPreparedData(filestore *FileStorage, store store.IStore, genesisSlot, btcGenesisHeight uint64, proverClient btcproverClient.IClient, btcClient *btcrpc.Client,
	ethClient *ethrpc.Client, apiClient *apiclient.Client, beaconClient *beacon.Client, minerAddr string) (*Prepared, error) {
	return &Prepared{
		filestore:        filestore,
		chainStore:       NewChainStore(store),
		proverClient:     proverClient,
		btcClient:        btcClient,
		ethClient:        ethClient,
		apiClient:        apiClient,
		beaconClient:     beaconClient,
		genesisSlot:      genesisSlot,
		genesisPeriod:    genesisSlot / common.SlotPerPeriod,
		btcGenesisHeight: btcGenesisHeight,
		minerAddr:        minerAddr,
	}, nil
}

func dbFinalityUpdateToZkFinalityUpdate(update *common.LightClientFinalityUpdateEvent) *proverType.FinalityUpdate {
	return &proverType.FinalityUpdate{
		Version: update.Version,
		AttestedHeader: &proverType.BeaconHeader{
			Slot:          update.Data.AttestedHeader.Slot,
			ProposerIndex: update.Data.AttestedHeader.ProposerIndex,
			ParentRoot:    update.Data.AttestedHeader.ParentRoot,
			StateRoot:     update.Data.AttestedHeader.StateRoot,
			BodyRoot:      update.Data.AttestedHeader.BodyRoot,
		},
		FinalizedHeader: &proverType.BeaconHeader{
			Slot:          update.Data.FinalizedHeader.Slot,
			ProposerIndex: update.Data.FinalizedHeader.ProposerIndex,
			ParentRoot:    update.Data.FinalizedHeader.ParentRoot,
			StateRoot:     update.Data.FinalizedHeader.StateRoot,
			BodyRoot:      update.Data.FinalizedHeader.BodyRoot,
		},
		//CurrentSyncCommitteeBranch: nil,
		FinalityBranch: update.Data.FinalityBranch,
		SyncAggregate: &proverType.SyncAggregate{
			SyncCommitteeBits:      update.Data.SyncAggregate.SyncCommitteeBits,
			SyncCommitteeSignature: update.Data.SyncAggregate.SyncCommitteeSignature,
		},
		SignatureSlot: update.Data.SignatureSlot,
	}

}
func dbUpdateToZkUpdate(update *common.LightClientUpdateResponse) (rpc.WrapSyncCommitteeUpdate, error) {
	zkUpdate := rpc.WrapSyncCommitteeUpdate{
		SyncCommitteeUpdate: &proverType.SyncCommitteeUpdate{
			AttestedHeader: &proverType.BeaconHeader{
				Slot:          update.Data.AttestedHeader.Slot,
				ProposerIndex: update.Data.AttestedHeader.ProposerIndex,
				ParentRoot:    update.Data.AttestedHeader.ParentRoot,
				StateRoot:     update.Data.AttestedHeader.StateRoot,
				BodyRoot:      update.Data.AttestedHeader.BodyRoot,
			},
			CurrentSyncCommittee: nil,
			FinalizedHeader: &proverType.BeaconHeader{
				Slot:          update.Data.FinalizedHeader.Slot,
				ProposerIndex: update.Data.FinalizedHeader.ProposerIndex,
				ParentRoot:    update.Data.FinalizedHeader.ParentRoot,
				StateRoot:     update.Data.FinalizedHeader.StateRoot,
				BodyRoot:      update.Data.FinalizedHeader.BodyRoot,
			},
			Version: update.Version,
			SyncAggregate: &proverType.SyncAggregate{
				SyncCommitteeBits:      update.Data.SyncAggregate.SyncCommitteeBits,
				SyncCommitteeSignature: update.Data.SyncAggregate.SyncCommitteeSignature,
			},
			FinalityBranch: update.Data.FinalityBranch,
			NextSyncCommittee: &proverType.SyncCommittee{
				PubKeys:         update.Data.NextSyncCommittee.Pubkeys,
				AggregatePubKey: update.Data.NextSyncCommittee.AggregatePubkey,
			},
			NextSyncCommitteeBranch: update.Data.NextSyncCommitteeBranch,
			SignatureSlot:           update.Data.SignatureSlot,
		},
	}
	return zkUpdate, nil
}
