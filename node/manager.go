package node

import (
	"encoding/hex"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/node/p2p"
	"github.com/lightec-xyz/daemon/rpc/bitcoin"
	"github.com/lightec-xyz/daemon/rpc/dfinity"
	"github.com/lightec-xyz/daemon/rpc/ethereum"
	"github.com/lightec-xyz/daemon/store"
	"time"
)

type manager struct {
	fileStore      *FileStorage
	libP2p         *p2p.LibP2p
	scheduler      *Scheduler
	chainStore     *ChainStore
	btcProofResp   chan *common.ProofResponse
	ethProofResp   chan *common.ProofResponse
	syncCommitResp chan *common.ProofResponse
	ethNotify      chan *Notify
	btcNotify      chan *Notify
	beaconNotify   chan *Notify
	minerPower     *MinerPower
	btcForks       []*ChainFork
	appStartTime   time.Time
}

func NewManager(minerAddr string, libP2p *p2p.LibP2p, icpClient *dfinity.Client, btcClient *bitcoin.Client, ethClient *ethereum.Client, prep *Prepared,
	btcProofResp, ethProofResp, syncCommitteeProofResp chan *common.ProofResponse, store store.IStore, fileStore *FileStorage,
	btcNotify, ethNotify, beaconNotify chan *Notify) (IManager, error) {
	scheduler, err := NewScheduler(fileStore, store, prep, icpClient, btcClient, ethClient)
	if err != nil {
		logger.Error("new scheduler error:%v", err)
		return nil, err
	}
	return &manager{
		libP2p:         libP2p,
		chainStore:     NewChainStore(store),
		fileStore:      fileStore,
		btcProofResp:   btcProofResp,
		ethProofResp:   ethProofResp,
		syncCommitResp: syncCommitteeProofResp,
		scheduler:      scheduler,
		btcNotify:      btcNotify,
		ethNotify:      ethNotify,
		beaconNotify:   beaconNotify,
		appStartTime:   time.Now(),
		minerPower:     NewMinerPower(minerAddr, 0, time.Now()),
	}, nil
}

func (m *manager) Init() error {
	btcChainFork, err := m.chainStore.ReadChainForks(common.BitcoinChain.String())
	if err != nil {
		logger.Error("read btc chain fork error:%v", err)
		return err
	}
	m.btcForks = append(m.btcForks, btcChainFork...)
	err = m.chainStore.WriteMiner(m.minerPower.Address)
	if err != nil {
		logger.Error("write miner error:%v", err)
		return err
	}
	err = m.chainStore.WriteMinerPower(m.minerPower.Address, m.minerPower.Power, uint64(time.Now().Unix()))
	if err != nil {
		logger.Error("write miner power error:%v", err)
		return err
	}
	return nil
}

func (m *manager) MinerPower() error {
	//todo  temporarily a simple method to calculate the power per hour,
	avgConstantPerHour := int64(m.minerPower.AvgConstantPerHour())
	timestamp := time.Now().Unix()
	minerMsg := p2p.NewP2pMinerMsg(m.minerPower.Address, avgConstantPerHour, timestamp)
	err := m.chainStore.WriteMinerPower(m.minerPower.Address, m.minerPower.Power, uint64(timestamp))
	if err != nil {
		logger.Error("write miner power error:%v", err)
		//return err
	}
	if m.libP2p != nil {
		err := m.libP2p.Broadcast(minerMsg)
		if err != nil {
			logger.Error("broadcast minerType power error:%v", err)
			return err
		}
	}
	return nil
}

func (m *manager) LibP2pMessage(msg *p2p.Msg) error {
	logger.Debug("libp2p message:%v", msg.String())
	switch msg.GetType() {
	case p2p.Msg_Hello:
		err := m.chainStore.WriteMiner(msg.GetHello().GetAddress())
		if err != nil {
			logger.Error("write miner error:%v", err)
			return err
		}
		return nil
	case p2p.Msg_Miner:
		miner := msg.GetMiner()
		err := m.chainStore.WriteMinerPower(miner.GetMinerAddr(),
			uint64(miner.GetPower()), uint64(*msg.Timestamp))
		if err != nil {
			logger.Error("write miner power error:%v", err)
			return err
		}
		return nil
	default:
		logger.Error("unknown message type:%v", msg)
		return nil
	}
}

func (m *manager) GetProofRequest(proofTypes []common.ProofType) (*common.ProofRequest, bool, error) {
	if m.scheduler.queueManager.RequestLen() == 0 {
		return nil, false, nil
	}
	var request *common.ProofRequest
	var ok bool
	if len(proofTypes) == 0 {
		request, ok = m.scheduler.queueManager.PopRequest()
	} else {
		request, ok = m.scheduler.queueManager.PopFnRequest(func(request *common.ProofRequest) bool {
			for _, req := range proofTypes {
				if request.ProofType == req {
					return true
				}
			}
			return false
		})
	}
	if !ok {
		logger.Warn("no find match proof task")
		return nil, false, nil
	}
	storeKey := NewStoreKey(request.ProofType, request.Hash, request.Prefix, request.FIndex, request.SIndex)
	proofId := storeKey.ProofId()
	exists, err := m.fileStore.CheckProof(storeKey)
	if err != nil {
		logger.Error("check Proof error:%v %v", proofId, err)
		return nil, false, err
	}
	if exists {
		logger.Debug("proof request exists: %v", proofId)
		m.scheduler.removeRequest(proofId)
		return nil, false, nil
	}
	request.SetStartTime(time.Now())
	m.scheduler.addRequestToPending(request)
	err = m.updateProofStatus(request.Hash, "", request.ProofType, common.ProofGenerating)
	if err != nil {
		logger.Error("update Proof status error:%v %v", proofId, err)
	}
	return request, true, nil
}

func (m *manager) ReceiveProofs(res *common.SubmitProof) error {
	if res.Status {
		go m.storeProof(res.Responses)
	} else {
		for _, req := range res.Requests {
			m.scheduler.removeRequest(req.ProofId())
		}
	}
	return nil
}

func (m *manager) CheckState() error {
	logger.Debug("check pending req now")
	pendingRequest := m.scheduler.PendingRequest()
	for _, req := range pendingRequest {
		if req == nil {
			continue
		}
		if req.StartTime.IsZero() {
			logger.Error("never should happen,req start time is zero: %v", req.ProofId())
			m.scheduler.removeRequest(req.ProofId())
			continue
		}
		timeout := time.Now().Sub(req.StartTime) >= req.ProofType.ProveTime()
		if timeout {
			storeKey := NewStoreKey(req.ProofType, req.Hash, req.Prefix, req.FIndex, req.SIndex)
			proofId := storeKey.ProofId()
			exists, err := m.fileStore.CheckProof(storeKey)
			if err != nil {
				logger.Error("check proof error:%v %v", proofId, err)
				continue
			}
			m.scheduler.removeRequest(proofId)
			if !exists {
				logger.Debug("%v timeout,add proof queue again", proofId)
			}
		}
		return nil
	}
	return nil
}

func (m *manager) storeProof(responses []*common.ProofResponse) {
	for _, item := range responses {
		//todo
		if common.IsBtcProofType(item.ProofType) && item.ReqCreateTime.Before(m.appStartTime) {
			continue
		}
		m.minerPower.AddConstant(item.ProofType.ConstraintQuantity())
		forked := m.checkForkedProof(item)
		storeKey := NewStoreKey(item.ProofType, item.Hash, item.Prefix, item.FIndex, item.SIndex)
		proofId := storeKey.ProofId()
		if forked {
			logger.Warn("success receive proof,%v but it's forked proof,try again", proofId)
			m.scheduler.removeRequest(proofId)
			continue
		}
		// first btcUpper proof,store it to btcDuperProof
		if item.ProofType == common.BtcUpperType && item.FIndex == m.fileStore.btcGenesisHeight {
			key := NewDoubleStoreKey(common.BtcDuperRecursiveType, item.FIndex, item.SIndex)
			err := m.fileStore.StoreProof(key, item.Proof, item.Witness)
			if err != nil {
				logger.Error("store proof error: %v %v", proofId, err)
			}
		}
		err := m.fileStore.StoreProof(storeKey, item.Proof, item.Witness)
		if err != nil {
			logger.Error("store proof error: %v %v", proofId, err)
			continue
		}
		logger.Debug("store zk proof: %v", proofId)
		m.notify(item)
		m.scheduler.removeRequest(proofId)
		chanResponse := m.getChanResponse(item.ProofType)
		if chanResponse != nil {
			chanResponse <- item
		}
		logger.Info("delete pending request:%v", proofId)
		err = m.updateProofStatus(item.Hash, hex.EncodeToString(item.Proof), item.ProofType, common.ProofSuccess)
		if err != nil {
			logger.Error("update Proof status error:%v", err)
		}
	}
}
func (m *manager) checkForkedProof(resp *common.ProofResponse) bool {
	if common.IsBtcProofType(resp.ProofType) {
		if m.btcForks == nil {
			return false
		}
		for _, forkInfo := range m.btcForks {
			if resp.SIndex >= forkInfo.ForkHeight && resp.ReqCreateTime.UnixNano() <= forkInfo.Timestamp {
				logger.Warn("find forked response proofId: %v,height:%v,time:%v", resp.ProofId(), resp.SIndex, resp.ReqCreateTime)
				return true
			}
		}
	}
	return false
}

func (m *manager) getChanResponse(reqType common.ProofType) chan *common.ProofResponse {
	switch reqType {
	case common.BtcUpdateCpType, common.BtcDepositType, common.BtcChangeType:
		return m.btcProofResp
	case common.RedeemTxType:
		return m.ethProofResp
	default:
		return nil
	}
}

func (m *manager) Close() error {
	logger.Debug("manager start  cache cache now ...")
	pendingRequest := m.scheduler.PendingRequest()
	for _, req := range pendingRequest {
		proofId := req.ProofId()
		logger.Debug("write pending request to db :%v", proofId)
		err := m.chainStore.WritePendingRequest(proofId, req)
		if err != nil {
			logger.Error("write pending request error:%v %v", proofId, err)
			return err
		}
		return nil
	}
	return nil

}

func (m *manager) updateProofStatus(hash, proof string, proofType common.ProofType, status common.ProofStatus) error {
	if proofType == common.BtcDepositType || proofType == common.RedeemTxType || proofType == common.BtcChangeType {
		err := m.chainStore.UpdateProof(hash, proof, proofType, status)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *manager) notify(req *common.ProofResponse) {
	// maybe should schedule by proofs a proofType,it`s will better
	switch req.ProofType {
	case common.TxInEth2Type, common.BeaconHeaderType, common.BeaconHeaderFinalityType:
		select {
		case m.ethNotify <- &Notify{}:
		default:
		}
	case common.BtcBaseType, common.BtcTimestampType, common.BtcMiddleType, common.BtcUpperType, common.BtcBulkType, common.BtcDuperRecursiveType, common.BtcDepthRecursiveType:
		select {
		case m.btcNotify <- &Notify{}:
		default:
		}
	case common.SyncComInnerType, common.SyncComOuterType, common.SyncComUnitType:
		select {
		case m.beaconNotify <- &Notify{}:
		default:
		}

	}
}
func (m *manager) ChainFork(info *ChainFork) error {
	logger.Debug("manager chain fork:%v,forkHeight:%v,time:%v", info.Chain.String(), info.ForkHeight, info.Timestamp)
	unLocks := m.scheduler.Locks()
	defer unLocks()
	info.Timestamp = time.Now().UnixNano()
	err := m.chainStore.WriteChainFork(info.Chain.String(), info)
	if err != nil {
		logger.Error("write chainFork: %v", info)
		return err
	}
	switch info.Chain {
	case common.BitcoinChain:
		m.btcForks = append(m.btcForks, info)
		err := m.fileStore.RemoveBtcProof(info.ForkHeight)
		if err != nil {
			logger.Error("remove btc proof error:%v %v", info, err)
			return err
		}
		err = m.scheduler.btcStateRollback(info.ForkHeight)
		if err != nil {
			logger.Error("btc scheduler roll back error:%v %v", info, err)
			return err
		}
		return nil
	case common.EthereumChain:
		//nothing to do
		return nil
	default:
		logger.Error("unknown chain:%v", info.Chain)
		return nil
	}
	return nil
}

func (m *manager) AddP2pPeer(addr string) error {
	err := m.libP2p.AddPeer(addr)
	if err != nil {
		logger.Error("add p2p peer error:%v", err)
		return err
	}
	return nil
}
func (m *manager) UpdateBtcCp() error {
	err := m.scheduler.updateBtcCp()
	if err != nil {
		logger.Error("update btc check point error:%v", err)
		return err
	}
	return nil
}

func (m *manager) CheckPreBtcState() error {
	err := m.scheduler.CheckPreBtcState()
	if err != nil {
		logger.Error("check pre btc scheduler error:%v", err)
		return err
	}
	return nil
}

func (m *manager) CheckBtcState() error {
	err := m.scheduler.CheckBtcState()
	if err != nil {
		logger.Error("check btc scheduler error:%v", err)
		return err
	}
	return nil
}

func (m *manager) CheckEthState() error {
	err := m.scheduler.CheckEthState()
	if err != nil {
		logger.Error("check eth scheduler error:%v", err)
		return err
	}
	return nil
}
func (m *manager) CheckBeaconState() error {
	err := m.scheduler.CheckBeaconState()
	if err != nil {
		logger.Error("check beacon scheduler error:%v", err)
		return err
	}
	return nil
}
func (m *manager) PendingProofRequest() []*common.ProofRequest {
	return m.scheduler.PendingProofRequest()
}
func (m *manager) BlockSignature() error {
	return m.scheduler.BlockSignature()
}
func (m *manager) EthNotify() chan *Notify {
	return m.ethNotify
}
func (m *manager) BtcNotify() chan *Notify {
	return m.btcNotify
}
func (m *manager) BeaconNotify() chan *Notify {
	return m.beaconNotify
}
