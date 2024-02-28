package node

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/rpc/beacon"
	"github.com/prysmaticlabs/prysm/v5/api/server/structs"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type DownloadStatus int32

const (
	DownloadStatus_NONE        DownloadStatus = 0
	DownloadStatus_Downloading DownloadStatus = 1
	DownloadStatus_Done        DownloadStatus = 2
)

type Status int32

const (
	Status_NONE            Status = 0 //apply to ProofG, ProofU, ProofR,  init state
	Status_ProofGenerating Status = 1 //apply to ProofG, ProofU, ProofR, generate proof
	Status_ProofGenerated  Status = 2 //apply to ProofG, ProofU, ProofR, proof generated
)

type ProofG struct {
	period          uint64
	status          Status
	bootStrapStatus DownloadStatus
}

type ProofU struct {
	period           uint64
	isGenesis        bool
	status           Status
	preUpdateStatus  DownloadStatus //if isGenesis is true, preUpdateStatus indicating bootstrap_xx.json status
	currUpdateStatus DownloadStatus
}

type ProofR struct {
	period uint64
	status Status
	proofU ProofU
}

type BeaconAgent struct {
	updatesDir        string
	unitProofDir      string
	recursiveProofDir string

	beaconClient      *beacon.Client
	proofClient       rpc.ISyncCommitteeProof
	datadir           string
	genesisSyncPeriod uint64
	targetSyncPeriod  atomic.Uint64

	proofG              ProofG
	proofR              []*ProofR
	lastGeneratedProofG atomic.Int64 //indicate which sync period has been generated

	newPeriodArriveCh   chan uint64
	genesisProofReqCh   chan uint64 //notify to generate genesis proof
	unitProofReqCh      chan uint64 //notify to generate unit proof
	recursiveProofReqCh chan uint64

	downloadingUpdate        *sync.Map                                  //TODO(keep), for future use,
	downloadedUpdateChan     chan *structs.LightClientUpdateWithVersion //TODO(keep), for future use,
	generatingGenesisProof   *sync.Map                                  //TODO(keep), for future use,
	genesisProofRespCh       chan *rpc.SyncCommitteeProofResponse       //TODO(keep), for future use,
	generatingUnitProof      *sync.Map                                  //TODO(keep), for future use,
	unitProofRespCh          chan *rpc.SyncCommitteeProofResponse       //TODO(keep), for future use,
	generatingRecursiveProof *sync.Map                                  //TODO(keep), for future use,
	recursiveProofRespCh     chan *rpc.SyncCommitteeProofResponse       //TODO(keep), for future use,

	exitCh chan struct{}
	wg     sync.WaitGroup
}

func NewBeaconAgent(cfg NodeConfig, beaconClient *beacon.Client, proofClient rpc.ISyncCommitteeProof) (*BeaconAgent, error) {
	genesisSyncPeriod := uint64(cfg.BeaconConfig.AltairForkEpoch) / 256
	targetSyncPeriod, err := beaconClient.GetLatestSyncPeriod()
	if err != nil {
		logger.Error("get latest sync period error:%v", err)
		return nil, err
	}

	//1. check targetSyncPeriod >= genesisSyncPeriod+1
	if targetSyncPeriod < genesisSyncPeriod+1 {
		logger.Error("targetSyncPeriod < genesisSyncPeriod+1")
		return nil, fmt.Errorf("targetSyncPeriod < genesisSyncPeriod+1")
	}

	//2. build proofG
	pg := ProofG{
		period:          genesisSyncPeriod,
		status:          Status_NONE,
		bootStrapStatus: DownloadStatus_NONE,
	}

	//3. build proofR
	//proofR value store in [0, genesisSyncPeriod] are empty, proofR value store in [genesisSyncPeriod+1, targetSyncPeriod]
	proofR := make([]*ProofR, targetSyncPeriod)
	for i := uint64(genesisSyncPeriod) + 1; i < targetSyncPeriod; i++ {
		isGenesis := i == genesisSyncPeriod+1
		pu := ProofU{
			period:           i - 1,
			isGenesis:        isGenesis,
			status:           Status_NONE,
			preUpdateStatus:  DownloadStatus_NONE,
			currUpdateStatus: DownloadStatus_NONE,
		}

		pr := &ProofR{
			period: i,
			status: Status_NONE,
			proofU: pu,
		}
		proofR[i] = pr
	}

	agent := &BeaconAgent{
		updatesDir:               cfg.DataDir + "/updates_data",
		unitProofDir:             cfg.DataDir + "/uint_proofs",
		recursiveProofDir:        cfg.DataDir + "/recursive_proofs",
		beaconClient:             beaconClient,
		proofClient:              proofClient,
		datadir:                  cfg.DataDir,
		genesisSyncPeriod:        genesisSyncPeriod,
		proofG:                   pg,
		proofR:                   proofR,
		newPeriodArriveCh:        make(chan uint64, 10),
		genesisProofReqCh:        make(chan uint64, 10),
		unitProofReqCh:           make(chan uint64, 10),
		recursiveProofReqCh:      make(chan uint64, 10),
		downloadingUpdate:        &sync.Map{},
		generatingUnitProof:      &sync.Map{},
		generatingRecursiveProof: &sync.Map{},
		exitCh:                   make(chan struct{}),
	}
	agent.targetSyncPeriod.Store(targetSyncPeriod)
	agent.lastGeneratedProofG.Store(int64(genesisSyncPeriod) - 1)

	return agent, nil
}

func (b *BeaconAgent) Init() error {
	logger.Info("beacon agent init now")

	updatesDir := b.updatesDir
	unitProofDir := b.unitProofDir
	recursiveProofDir := b.recursiveProofDir

	b.recursiveProofDir = b.datadir + "/recursive_proofs"

	//build Pg
	f1 := recursiveProofDir + fmt.Sprintf("/genesis_proof_%v.json", b.genesisSyncPeriod)
	_, err := os.Stat(f1)
	if err == nil {
		b.proofG.status = Status_ProofGenerated
	} else {
		f2 := updatesDir + fmt.Sprintf("/bootstrap_%v.json", b.genesisSyncPeriod)
		_, err = os.Stat(f2)
		if err == nil {
			b.proofG.bootStrapStatus = DownloadStatus_Done
		}
	}

	//build Prs
	for i := uint64(b.genesisSyncPeriod) + 1; i < b.targetSyncPeriod.Load(); i++ {
		pr := b.proofR[i]
		f1 = recursiveProofDir + fmt.Sprintf("/recursive_proof_%v.json", i)
		_, err = os.Stat(f1)
		if err == nil {
			pr.status = Status_ProofGenerated
		} else {
			f2 := unitProofDir + fmt.Sprintf("/uint_proof_%v.json", pr.proofU.period)
			_, err = os.Stat(f2)
			if err == nil {
				pr.proofU.status = Status_ProofGenerated
			} else {
				//check whether current update files exist
				f3 := updatesDir + fmt.Sprintf("/update_%v.json", pr.proofU.period)
				_, err = os.Stat(f3)
				if err == nil {
					pr.proofU.currUpdateStatus = DownloadStatus_Done
				}

				//check whether previous update files exist
				if pr.proofU.isGenesis {
					f3 = updatesDir + fmt.Sprintf("/bootstrap_%v.json", b.genesisSyncPeriod)
				} else {
					f3 = updatesDir + fmt.Sprintf("/update_%v.json", pr.proofU.period-1)
				}
				_, err = os.Stat(f3)
				if err == nil {
					pr.proofU.preUpdateStatus = DownloadStatus_Done
				}
			}
		}
	}

	//download bootstrap_x.json before other update_x.json
	if b.proofG.bootStrapStatus == DownloadStatus_NONE {
		bootStrap, err := b.beaconClient.GetBootstrap(uint64(b.genesisSyncPeriod) * 32)
		if err != nil {
			logger.Error("get bootstrap error:%v", err)
			return err
		}

		f1 = updatesDir + fmt.Sprintf("/bootstrap_%v.json", b.genesisSyncPeriod)
		data, err := json.Marshal(bootStrap)
		if err != nil {
			logger.Error("marshal bootstrap error:%v", err)
			return err
		}
		err = os.WriteFile(f1, data, 0644)
		if err != nil {
			logger.Error("write genesis file error:%v", err)
			return err
		}
		b.proofG.bootStrapStatus = DownloadStatus_Done
		b.proofR[b.genesisSyncPeriod+1].proofU.preUpdateStatus = DownloadStatus_Done
	}

	b.wg.Add(5)
	go b.fetchUpdates()
	go b.genGenesisProof()
	go b.genUnitProof()
	go b.genRecursiveProof()
	go b.scanSyncPeriod()

	return nil
}

func (b *BeaconAgent) fetchAndStoreUpdate(i uint64) error {
	//if the update_i.json already exist, treat it as successfully download
	f := b.updatesDir + fmt.Sprintf("/update_%v.json", i)
	_, err := os.Stat(f)
	if err == nil {
		log.Error("update_%v.json already exist", i)
		return nil
	}

	updates, err := b.beaconClient.GetLightClientUpdates(i, 1)
	if err != nil {
		logger.Error("get light client updates error:%v", err)
		return err
	}

	data, err := json.Marshal(updates[0])
	if err != nil {
		logger.Error("marshal update error:%v", err)
		return err
	}

	err = os.WriteFile(f, data, 0644)
	if err != nil {
		logger.Error("write update file error:%v", err)
		return err
	}
	return nil
}

// fetch update_[i-1].json  and update_[i-2].json for proofR[i]
func (b *BeaconAgent) doFetchUpdate(i uint64) error {
	if i == b.genesisSyncPeriod {
		//the bootstrap_x.json should exist now, because we download it in Init()
		//TODO(keep): exception should be consider here
		if b.proofG.status == Status_NONE && b.proofG.bootStrapStatus == DownloadStatus_Done {
			//notify genGenesisProof routine work
			b.genesisProofReqCh <- b.genesisSyncPeriod
		}
		return nil
	}

	//not the genesis period
	pr := b.proofR[i]
	if pr.status == Status_ProofGenerated || pr.status == Status_ProofGenerating {
		//already in generating proof, or generated, do nothing
		return nil
	}

	var err error
	//check pr.proofU.preUpdateStatus
	if pr.proofU.preUpdateStatus == DownloadStatus_NONE {
		//download
		err = b.fetchAndStoreUpdate(pr.proofU.period - 1)
		if err != nil {
			return err
		}
		pr.proofU.preUpdateStatus = DownloadStatus_Done
	}

	//check pr.proofU.currUpdateStatus
	if pr.proofU.currUpdateStatus == DownloadStatus_NONE {
		//download
		err = b.fetchAndStoreUpdate(pr.proofU.period)
		if err != nil {
			return err
		}
		pr.proofU.currUpdateStatus = DownloadStatus_Done
	}

	//after download 2 updates, notify genUnitProof to generate Pu[i-1]
	b.unitProofReqCh <- pr.proofU.period

	return nil

}

func (b *BeaconAgent) fetchUpdates() {
	defer b.wg.Done()

	duration := 60 * time.Second
	ticker := time.NewTicker(duration)
	defer ticker.Stop()
	for {
		select {
		case <-b.exitCh:
			logger.Info("beacon agent exit ..., shut down fetch updates")
			return

		case <-ticker.C:
			for i := uint64(b.lastGeneratedProofG.Load() + 1); i < b.targetSyncPeriod.Load(); i++ {
				// download at most 10 files in one round
				err := b.doFetchUpdate(i)
				log.Error("fail to fetch update_%v, err:%v", i, err)
			}

		case i := <-b.newPeriodArriveCh:
			err := b.doFetchUpdate(i)
			if err != nil {
				log.Error("fail to fetch update_%v, err:%v", i, err)
			}
		}
	}

}

func (b *BeaconAgent) genGenesisProof() {
	defer b.wg.Done()
	for {
		select {
		case <-b.exitCh:
			logger.Info("beacon agent exit ..., shut down gen genesis proof")
			return

		case i := <-b.genesisProofReqCh:
			f1 := b.updatesDir + fmt.Sprintf("/bootstrap_%v.json", i)
			data1, err := os.ReadFile(f1)
			if err != nil {
				log.Error("read bootstrap_%v.json error:%v", i, err)
				return
			}
			var bootstrap structs.LightClientBootstrapResponse
			err = json.Unmarshal(data1, &bootstrap)
			if err != nil {
				log.Error("unmarshal bootstrap_%v.json error:%v", i, err)
				return
			}

			f2 := b.updatesDir + fmt.Sprintf("/update_%v.json", i)
			data2, err := os.ReadFile(f2)
			if err != nil {
				log.Error("read update_%v.json error:%v", i, err)
				return
			}
			var update structs.LightClientUpdateWithVersion
			err = json.Unmarshal(data2, &update)
			if err != nil {
				log.Error("unmarshal update_%v.json error:%v", i, err)
				return
			}

			req := rpc.GenesisSyncCommitteeProofRequest{
				Version:                    bootstrap.Version,
				AttestedHeader:             *bootstrap.Data.Header,
				CurrentSyncCommittee:       *bootstrap.Data.CurrentSyncCommittee,
				CurrentSyncCommitteeBranch: bootstrap.Data.CurrentSyncCommitteeBranch,
			}

			//TODO(keep), can be optimized with channel&sunc map
			proof, err := b.proofClient.GenGenesisSyncCommitteeProof(req)
			if err != nil {
				log.Error("fail to generate genesis proof, err:%v", err)
				return
			}
			if proof.Period != b.genesisSyncPeriod {
				msg := fmt.Sprintf("proof.Period: %v != b.genesisSyncPeriod:%v", proof.Period, b.genesisSyncPeriod)
				log.Error(msg)
				return
			}
			if proof.ProofType != rpc.SyncCommitteeProofType_Genesis {
				msg := fmt.Sprintf("unexpected proof type:%v", proof.ProofType)
				log.Error(msg)
				return
			}

			if proof.Status == rpc.SyncCommitteeProofGenerateStatus_Done {
				log.Error("unexpected proof status:%v", proof.Status)
				return
			}

			f3 := b.recursiveProofDir + fmt.Sprintf("/genesis_proof_%v.json", b.genesisSyncPeriod)
			data, err := json.Marshal(proof)
			if err != nil {
				log.Error(err.Error())
				return
			}
			err = os.WriteFile(f3, data, 0644)
			if err != nil {
				log.Error(err.Error())
				return
			}

			//notify to generate next recursive proof
			nextPr := b.proofR[b.genesisSyncPeriod+1]
			if nextPr.status == Status_NONE && nextPr.proofU.status == Status_ProofGenerated {
				b.recursiveProofReqCh <- nextPr.period
			}

			/*
				case proof := <-b.genesisProofRespCh:
					if proof.Period != b.genesisSyncPeriod {
						msg := fmt.Sprintf("proof.Period: %v != b.genesisSyncPeriod:%v", proof.Period, b.genesisSyncPeriod)
						log.Error(msg)
						return
					}
					if proof.ProofType != rpc.SyncCommitteeProofType_Genesis {
						msg := fmt.Sprintf("unexpected proof type:%v", proof.ProofType)
						log.Error(msg)
						return
					}

					if proof.Status == rpc.SyncCommitteeProofGenerateStatus_Done {
						log.Error("unexpected proof status:%v", proof.Status)
						return
					}

					f1 := b.recursiveProofDir + fmt.Sprintf("/genesis_proof_%v.json", b.genesisSyncPeriod)
					data, err := json.Marshal(proof)
					if err != nil {
						log.Error(err.Error())
						return
					}
					err = os.WriteFile(f1, data, 0644)
					if err != nil {
						log.Error(err.Error())
						return
					}

					//check the generatingGenesisProof and delete it unconditionally
					_, exist := b.generatingGenesisProof.Load(b.genesisSyncPeriod)
					if !exist {
						msg := fmt.Sprintf("response with no requesd")
						log.Error(msg)
					}
					b.generatingGenesisProof.Delete(b.genesisSyncPeriod)

					//trig to generate next recursive proof
					nextPr := b.proofR[b.genesisSyncPeriod+1]
					if nextPr.status == Status_NONE && nextPr.proofU.status == Status_ProofGenerated {
						b.recursiveProofReqCh <- nextPr.period
					}
			*/

		}
	}

}

// Can be parrelled
func (b *BeaconAgent) genUnitProof() {
	defer b.wg.Done()
	for {
		select {
		case <-b.exitCh:
			logger.Info("beacon agent exit ..., shut down gen unit proof")
			return
		case i := <-b.unitProofReqCh:
			f := b.unitProofDir + fmt.Sprintf("/unit_proof_%v.json", i)
			if _, err := os.Stat(f); err == nil {
				log.Error("unit proof %v already exists", i)
				continue
			}

			f1 := b.updatesDir + fmt.Sprintf("/update_%v.json", i-1)
			data1, err := os.ReadFile(f1)
			if err != nil {
				log.Error("read upadte_%v.json error:%v", i, err)
				continue
			}
			var update1 structs.LightClientUpdateWithVersion
			err = json.Unmarshal(data1, &update1)
			if err != nil {
				log.Error("unmarshal update_%v.json error:%v", i, err)
				continue
			}

			f2 := b.updatesDir + fmt.Sprintf("/update_%v.json", i)
			data2, err := os.ReadFile(f2)
			if err != nil {
				log.Error("read update_%v.json error:%v", i, err)
				continue
			}
			var update2 structs.LightClientUpdateWithVersion
			err = json.Unmarshal(data2, &update2)
			if err != nil {
				log.Error("unmarshal update_%v.json error:%v", i, err)
				continue
			}

			req := rpc.UnitSyncCommitteeProofRequest{
				Version:                 update2.Version,
				AttestedHeader:          *update2.Data.AttestedHeader,
				CurrentSyncCommittee:    *update1.Data.NextSyncCommittee,
				SyncAggregate:           *update2.Data.SyncAggregate,
				NextSyncCommittee:       *update2.Data.NextSyncCommittee,
				NextSyncCommitteeBranch: update2.Data.NextSyncCommitteeBranch,
			}

			proof, err := b.proofClient.GenUnitSyncCommitteeProof(req)
			if err != nil {
				log.Error("gen unit proof error:%v", err)
				continue
			}

			f3 := b.unitProofDir + fmt.Sprintf("/unit_proof_%v.json", proof.Period)
			data3, err := json.Marshal(proof)
			if err != nil {
				log.Error(err.Error())
				continue
			}
			err = os.WriteFile(f3, data3, 0644)
			if err != nil {
				log.Error(err.Error())
				continue
			}

			b.proofR[proof.Period+1].proofU.status = Status_ProofGenerated

			//check the generatingUnitProof and delete it unconditionally
			if proof.Period == b.genesisSyncPeriod {
				if b.proofG.status == Status_ProofGenerated {
					b.recursiveProofReqCh <- proof.Period
				}
			} else {
				prePr := b.proofR[proof.Period-1]
				if prePr.status == Status_ProofGenerated {
					b.recursiveProofReqCh <- proof.Period
				}

			}
		}
	}
}

func (b *BeaconAgent) genRecursiveProof() {
	defer b.wg.Done()
	for {
		select {
		case <-b.exitCh:
			logger.Info("beacon agent exit ..., shut down gen recursive proof")
			return
		case i := <-b.recursiveProofReqCh:
			f := b.recursiveProofDir + fmt.Sprintf("/recursive_proof_%v.json", i)
			if _, err := os.Stat(f); err == nil {
				log.Error("recursive proof %v already exists", i)
				continue
			}

			f1 := b.unitProofDir + fmt.Sprintf("/unit_proof_%v.json", i-1)
			data1, err := os.ReadFile(f1)
			if err != nil {
				log.Error("read unit proof %v error:%v", i-1, err)
				continue
			}
			var proofU rpc.SyncCommitteeProofResponse
			err = json.Unmarshal(data1, &proofU)
			if err != nil {
				log.Error("unmarshal unit proof %v error:%v", i-1, err)
				continue
			}

			var preProof string
			if i == b.genesisSyncPeriod+1 {
				if b.proofG.status == Status_ProofGenerated {
					f2 := b.recursiveProofDir + fmt.Sprintf("/genesis_proof_%v.json", i)
					data2, err := os.ReadFile(f2)
					if err != nil {
						log.Error("read genesis proof %v error:%v", i, err)
						continue
					}
					var proof2 rpc.SyncCommitteeProofResponse
					err = json.Unmarshal(data2, &proof2)
					if err != nil {
						log.Error("unmarshal genesis proof %v error:%v", i, err)
						continue
					}
					preProof = proof2.Proof
				} else {
					//Pg is not generated
					continue
				}
			} else {
				prePr := b.proofR[i-1]
				if prePr.status == Status_ProofGenerated {
					f2 := b.recursiveProofDir + fmt.Sprintf("/recursive_proof_%v.json", i-1)
					data2, err := os.ReadFile(f2)
					if err != nil {
						log.Error("read recursive proof %v error:%v", i-1, err)
						continue
					}
					var proof2 rpc.SyncCommitteeProofResponse
					err = json.Unmarshal(data2, &proof2)
					if err != nil {
						log.Error("unmarshal recursive proof %v error:%v", i-1, err)
						continue
					}
				} else {
					//Previous Pr is not generated
					continue
				}
			}

			req := rpc.RecursiveSyncCommitteeProofRequest{
				Version:          proofU.Version,
				PreProofGOrPoofR: preProof,
				ProofU:           proofU.Proof,
			}

			proof, err := b.proofClient.GenRecursiveSyncCommitteeProof(req)
			if err != nil {
				log.Error("gen recursive proof error:%v", err)
				continue
			}

			data, err := json.Marshal(proof)
			if err != nil {
				log.Error(err.Error())
				continue
			}
			err = os.WriteFile(f, data, 0644)
			if err != nil {
				log.Error(err.Error())
				continue
			}

			b.proofR[i].status = Status_ProofGenerated
			if b.proofR[i+1].status == Status_NONE && b.proofR[i+1].proofU.status == Status_ProofGenerated {
				b.recursiveProofReqCh <- (i + 1)
			}

		}
	}
}

func (b *BeaconAgent) doScanSyncPeriod() {
	period, err := b.beaconClient.GetLatestSyncPeriod()
	if err != nil {
		logger.Error("get latest sync period error:%v, will try later", err)
		return
	}
	targetSyncPeriod := b.targetSyncPeriod.Load()
	if period < b.targetSyncPeriod.Load() {
		logger.Error("latest get sync period %v < previous get sync period %v", period, targetSyncPeriod)
		return
	}

	//notify fetch updates
	if period > targetSyncPeriod {
		b.targetSyncPeriod.Store(period)
		for i := targetSyncPeriod + 1; i <= period; i++ {
			b.newPeriodArriveCh <- i
		}
	}
}

func (b *BeaconAgent) scanSyncPeriod() {
	defer b.wg.Done()
	b.doScanSyncPeriod()

	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-b.exitCh:
			logger.Info("beacon agent exit ..., shut down gen recursive proof")
			b.wg.Done()
			return
		case <-ticker.C:
			b.doScanSyncPeriod()
		}
	}
}

func (b *BeaconAgent) Stop() {
	close(b.exitCh)
	b.wg.Wait()
}
