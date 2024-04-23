package beacon

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/prysmaticlabs/go-bitfield"
	"github.com/prysmaticlabs/prysm/v5/api/server/structs"
	"github.com/prysmaticlabs/prysm/v5/beacon-chain/core/signing"
	"github.com/prysmaticlabs/prysm/v5/beacon-chain/state"
	"github.com/prysmaticlabs/prysm/v5/config/params"
	"github.com/prysmaticlabs/prysm/v5/container/trie"
	"github.com/prysmaticlabs/prysm/v5/crypto/bls"
	ethpb "github.com/prysmaticlabs/prysm/v5/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/v5/time/slots"
	"io"
	"math/big"
	"net/http"
)

var rpcURL = "http://37.120.151.183:8970"

func retrieveNonEmptySlotInPeriod(ctx context.Context, period uint) (uint, string, error) {
	resp := &structs.GetBlockHeaderResponse{}
	slot := period*8192 + 2048
	for {
		uri := rpcURL + fmt.Sprintf("/eth/v1/beacon/headers?slot=%v", slot)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
		if err != nil {
			slot++
			continue
		}

		r, err := http.DefaultClient.Do(req)
		defer func() {
			err = r.Body.Close()
		}()
		if err != nil {
			slot++
			continue
		}

		if r.StatusCode != http.StatusOK {
			slot++
			continue
		}
		data, err := io.ReadAll(r.Body)
		if err != nil {
			slot++
			continue
		}

		err = json.Unmarshal(data, resp)
		if err != nil {
			slot++
			continue
		}

		if slot >= (period+1)*8192 {
			return 0, "", fmt.Errorf("failed to get a non empty slot in period %v", period)
		}
		//success to get the slot's root
		break
	}

	return slot, resp.Data.Root, nil
}

// root is the parent root of the next non empty slot
func retrieveNextNonEmptySlot(ctx context.Context, root uint) (uint, string, error) {
	resp := &structs.GetBlockHeaderResponse{}
	//make sure current exist
	uri := rpcURL + fmt.Sprintf("/eth/v1/beacon/headers?parent_root=%v", root)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return 0, "", fmt.Errorf("failt to next non empty slot, failed: %s", err)
	}

	r, err := http.DefaultClient.Do(req)
	defer func() {
		err = r.Body.Close()
	}()
	if err != nil {
		return 0, "", fmt.Errorf("failt to next non empty slot, failed: %s", err)
	}

	if r.StatusCode != http.StatusOK {
		return 0, "", fmt.Errorf("failt to next non empty slot, failed: %s", err)
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return 0, "", fmt.Errorf("failt to next non empty slot, failed: %s", err)
	}

	err = json.Unmarshal(data, resp)
	if err != nil {
		return 0, "", fmt.Errorf("failt to next non empty slot, failed: %s", err)
	}

	slot, ok := big.NewInt(0).SetString(resp.Data.Header.Message.Slot, 10)
	if !ok {
		return 0, "", fmt.Errorf("fail to decode next non empty slot's slot number")
	}

	return uint(slot.Uint64()), resp.Data.Root, nil
}

func RetrieveSyncCommittee(ctx context.Context, period uint) (*structs.LightClientBootstrap, error) {
	slot, root, err := retrieveNonEmptySlotInPeriod(ctx, period)
	if err != nil {
		return nil, err
	}

	fmt.Printf("slot:%v\n", slot)
	uri := rpcURL + fmt.Sprintf("/eth/v1/beacon/light_client/bootstrap/%s", root)
	fmt.Printf("Requesting sync committee, uri:%v\n", uri)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = r.Body.Close()
	}()
	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("requesting light client bootstrap, bad status code %d", r.StatusCode)
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("requesting light client bootstrap,failed: %s", err)
	}

	resp := &structs.LightClientBootstrapResponse{}
	err = json.Unmarshal(data, resp)
	if err != nil {
		return nil, fmt.Errorf("unmarshal light client bootstrap response, failed: %s", err)
	}
	return resp.Data, nil
}

func RetrieveLightClientUpdateByRange(ctx context.Context, startPeriod uint, count uint) ([]*structs.LightClientUpdate, error) {
	uri := rpcURL + fmt.Sprintf("http://127.0.0.1:8970/eth/v1/beacon/light_client/updates/start_period=%d&count=%d", startPeriod, count)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting light client updates %s\n", err)
	}
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = r.Body.Close()
	}()
	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("requesting light client updates, bad status code %d", r.StatusCode)
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("requesting light client updates,failed: %s", err)
	}

	updates := []structs.LightClientUpdateWithVersion{}
	err = json.Unmarshal(data, &updates)
	if err != nil {
		return nil, fmt.Errorf("unmarshal light client bootstrap response, failed: %s", err)
	}

	res := make([]*structs.LightClientUpdate, len(updates))
	for i, update := range updates {
		res[i] = update.Data
	}
	return res, nil
}

func RetrieveLightClientFinalityUpdate(ctx context.Context) (*structs.LightClientFinalityUpdate, error) {
	uri := rpcURL + fmt.Sprintf("/eth/v1/beacon/light_client/finality_update")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting light client updates %s\n", err)
	}
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = r.Body.Close()
	}()
	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("requesting light client updates, bad status code %d", r.StatusCode)
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("requesting light client updates,failed: %s", err)
	}

	update := structs.LightClientFinalityUpdate{}
	err = json.Unmarshal(data, &update)
	if err != nil {
		return nil, fmt.Errorf("unmarshal light client bootstrap response, failed: %s", err)
	}

	return &update, nil
}

type LightClientUpdateInfo struct {
	Version                 string                     `json:"version"`
	AttestedHeader          *structs.BeaconBlockHeader `json:"attested_header"`
	CurrentSyncCommittee    *structs.SyncCommittee     `json:"current_sync_committee,omitempty"`     //current_sync_committee
	SyncAggregate           *structs.SyncAggregate     `json:"sync_aggregate"`                       //sync_aggregate for attested_header, signed by current_sync_committee
	FinalizedHeader         *structs.BeaconBlockHeader `json:"finalized_header,omitempty"`           //finalized_header in attested_header.state_root
	FinalityBranch          []string                   `json:"finality_branch,omitempty"`            // finality_branch in attested_header.state_root
	NextSyncCommittee       *structs.SyncCommittee     `json:"next_sync_committee,omitempty"`        //next_sync_committee in finalized_header.state_root
	NextSyncCommitteeBranch []string                   `json:"next_sync_committee_branch,omitempty"` //next_sync_committee branch in finalized_header.state_root
	SignatureSlot           string                     `json:"signature_slot"`
}

func buildCurrentSyncCommitteeMerkleProof(ctx context.Context, s state.BeaconState) (*structs.SyncCommittee, []string, error) {
	stateRoot, err := s.HashTreeRoot(ctx)
	if err != nil {
		return nil, nil, err
	}

	// verify current_sync_committee against finalized_header.state_root
	syncCommittee, err := s.CurrentSyncCommittee()
	if err != nil {
		return nil, nil, err
	}
	ps, err := s.CurrentSyncCommitteeProof(ctx)
	if err != nil {
		return nil, nil, err
	}

	root, err := syncCommittee.HashTreeRoot()
	if err != nil {
		return nil, nil, err
	}

	gIndex := uint64(22)
	valid := trie.VerifyMerkleProof(stateRoot[:], root[:], gIndex, ps)
	if !valid {
		return nil, nil, fmt.Errorf("fail to verify current sync committee proof, failed: %s", err)
	}

	proof := make([]string, len(ps))
	for i, v := range ps {
		proof[i] = hex.EncodeToString(v)
	}

	pubkeys := make([]string, len(syncCommittee.Pubkeys))
	for i, v := range syncCommittee.Pubkeys {
		pubkeys[i] = hex.EncodeToString(v)
	}

	return &structs.SyncCommittee{
		Pubkeys:         pubkeys,
		AggregatePubkey: hex.EncodeToString(syncCommittee.AggregatePubkey),
	}, proof, nil
}

func buildNextSyncCommitteeMerkleProof(ctx context.Context, s state.BeaconState) (*structs.SyncCommittee, []string, error) {
	stateRoot, err := s.HashTreeRoot(ctx)
	if err != nil {
		return nil, nil, err
	}

	// verify current_sync_committee against finalized_header.state_root
	syncCommittee, err := s.NextSyncCommittee()
	if err != nil {
		return nil, nil, err
	}
	ps, err := s.NextSyncCommitteeProof(ctx)
	if err != nil {
		return nil, nil, err
	}

	root, err := syncCommittee.HashTreeRoot()
	if err != nil {
		return nil, nil, err
	}

	gIndex := uint64(23)
	valid := trie.VerifyMerkleProof(stateRoot[:], root[:], gIndex, ps)
	if !valid {
		return nil, nil, fmt.Errorf("fail to verify current sync committee proof, failed: %s", err)
	}

	proof := make([]string, len(ps))
	for i, v := range ps {
		proof[i] = hex.EncodeToString(v)
	}

	pubkeys := make([]string, len(syncCommittee.Pubkeys))
	for i, v := range syncCommittee.Pubkeys {
		pubkeys[i] = hex.EncodeToString(v)
	}

	return &structs.SyncCommittee{
		Pubkeys:         pubkeys,
		AggregatePubkey: hex.EncodeToString(syncCommittee.AggregatePubkey),
	}, proof, nil
}

func verifyFinalityBranch(update *structs.LightClientUpdate) (bool, error) {
	finalizedHeader, err := update.FinalizedHeader.ToConsensus()
	if err != nil {
		return false, err
	}
	finalizedHeaderRoot, err := finalizedHeader.HashTreeRoot()
	if err != nil {
		return false, err
	}

	attestedStateRoot, err := decodeHex(update.AttestedHeader.StateRoot)
	if err != nil {
		return false, err
	}

	branch := make([][]byte, len(update.FinalityBranch))
	for i := 0; i < len(update.FinalityBranch); i++ {
		branch[i] = make([]byte, 32)
		b, err := decodeHex(update.FinalityBranch[i])
		if err != nil {
			return false, err
		}
		branch[i] = b
	}

	//verfiy finalized_header  against attested_header.state_root
	valid := trie.VerifyMerkleProof(attestedStateRoot[:], finalizedHeaderRoot[:], 105, branch)
	return valid, nil
}

func verifySyncCommitteeSignature(syncCommittee *ethpb.SyncCommittee, update *structs.LightClientUpdate, finalized state.BeaconState) (bool, error) {
	var pubkeys []bls.PublicKey
	aggregateBytes, err := decodeHex(update.SyncAggregate.SyncCommitteeBits)
	if err != nil {
		return false, err
	}

	aggregateBits := bitfield.Bitvector512(aggregateBytes)
	for i := uint64(0); i < aggregateBits.Len(); i++ {
		if aggregateBits.BitAt(i) {
			pubKey, err := bls.PublicKeyFromBytes(syncCommittee.Pubkeys[i])
			if err != nil {
				return false, err
			}
			pubkeys = append(pubkeys, pubKey)
		}
	}

	sigBytes, err := decodeHex(update.SyncAggregate.SyncCommitteeSignature)
	if err != nil {
		return false, err
	}
	sig, err := bls.SignatureFromBytes(sigBytes)
	if err != nil {
		return false, err
	}

	header, err := update.AttestedHeader.ToConsensus()
	if err != nil {
		return false, err
	}

	ps := slots.PrevSlot(header.Slot)
	domain, err := signing.Domain(finalized.Fork(), slots.ToEpoch(ps), params.BeaconConfig().DomainSyncCommittee, finalized.GenesisValidatorsRoot())
	if err != nil {
		return false, err
	}
	fmt.Printf("domain: %v\n", hex.EncodeToString(domain))

	signingRoot, err := signing.ComputeSigningRoot(header, domain)
	return sig.FastAggregateVerify(pubkeys, signingRoot), nil
}

func buildLightClientUpdateWithProof(ctx context.Context, update *structs.LightClientUpdate, finalized state.BeaconState) (*LightClientUpdateInfo, error) {
	attestedSlot, ok := big.NewInt(0).SetString(update.AttestedHeader.Slot, 10)
	if !ok {
		return nil, fmt.Errorf("failed to parse attested slot, failed: %s", update.AttestedHeader.Slot)
	}

	finalizedSlot, ok := big.NewInt(0).SetString(update.FinalizedHeader.Slot, 10)
	if !ok {
		return nil, fmt.Errorf("failed to parse finalized slot, failed: %s", update.FinalizedHeader.Slot)
	}

	if attestedSlot.Uint64()/8192 != finalizedSlot.Uint64()/8192 {
		return nil, fmt.Errorf("attested slot and finalized slot are not in same sync committee period")
	}

	// verify finalized_header.state_root against update.finalized_header.state_root
	finalizedHeaderStateRoot, err := decodeHex(update.FinalizedHeader.StateRoot)
	if err != nil {
		return nil, err
	}

	valid, err := verifyFinalityBranch(update)
	if err != nil {
		return nil, err
	}

	if !valid {
		return nil, fmt.Errorf("finality branch is not valid")
	}

	calculatedFinalizedHeaderStateRoot, err := finalized.HashTreeRoot(ctx)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(finalizedHeaderStateRoot, calculatedFinalizedHeaderStateRoot[:]) {
		return nil, fmt.Errorf("finalized header state root is not match")
	}

	currentSyncCommittee, _, err := buildCurrentSyncCommitteeMerkleProof(ctx, finalized)
	if err != nil {
		return nil, err
	}

	nextSyncCommittee, _, err := buildNextSyncCommitteeMerkleProof(ctx, finalized)
	if err != nil {
		return nil, err
	}

	aggregatePubkey1, err := decodeHex(nextSyncCommittee.AggregatePubkey)
	if err != nil {
		return nil, err
	}

	aggregatePubkey2, err := decodeHex(update.NextSyncCommittee.AggregatePubkey)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(aggregatePubkey1, aggregatePubkey2) {
		return nil, fmt.Errorf("aggregate pubkey is not match")
	}

	committee, err := finalized.CurrentSyncCommittee()
	if err != nil {
		return nil, err
	}

	//check signature
	valid, err = verifySyncCommitteeSignature(committee, update, finalized)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, fmt.Errorf("fail to verify sync committee signature, failed: %s", err)
	}

	return &LightClientUpdateInfo{
		AttestedHeader:          update.AttestedHeader,
		SyncAggregate:           update.SyncAggregate,
		FinalityBranch:          update.FinalityBranch,
		CurrentSyncCommittee:    currentSyncCommittee,
		FinalizedHeader:         update.FinalizedHeader,
		NextSyncCommittee:       update.NextSyncCommittee,
		NextSyncCommitteeBranch: update.NextSyncCommitteeBranch,
		SignatureSlot:           update.SignatureSlot,
	}, nil
}

func VerifyLightClientUpdateInfo(update *LightClientUpdateInfo) (bool, error) {
	/*
		holesky, Bellatrix: 0700000069b7d97441dbd33e5ee5b4cb8fc8b08d6a58a7274b6e6daf19ef4ca7
		holesky,Capella: 0700000017e2dad36f1d3595152042a9ad23430197557e2e7e82bc7f7fc72972
	*/
	var domain []byte
	switch update.Version {
	case "bellatrix":
		domain, _ = decodeHex("0700000069b7d97441dbd33e5ee5b4cb8fc8b08d6a58a7274b6e6daf19ef4ca7")
	case "capella":
		domain, _ = decodeHex("0700000017e2dad36f1d3595152042a9ad23430197557e2e7e82bc7f7fc72972")
	case "deneb":
		domain, _ = decodeHex("0700000069ae0e9900d509b38350c53915fccde15c6ef44214aa1b5bdec34d3a")
	default:
		panic("unknown version")
	}

	attestedHeader, err := update.AttestedHeader.ToConsensus()
	if err != nil {
		return false, err
	}

	finalizedHeader, err := update.FinalizedHeader.ToConsensus()
	if err != nil {
		return false, err
	}

	finalizedHeaderRoot, err := finalizedHeader.HashTreeRoot()
	if err != nil {
		return false, err
	}

	currentSyncCommittee, err := update.CurrentSyncCommittee.ToConsensus()
	if err != nil {
		return false, err
	}

	nextSyncCommittee, err := update.NextSyncCommittee.ToConsensus()
	if err != nil {
		return false, err
	}
	nextSyncCommitteeRoot, err := nextSyncCommittee.HashTreeRoot()
	if err != nil {
		return false, err
	}

	nextSyncCommitteeBranch := make([][]byte, len(update.NextSyncCommitteeBranch))
	for i, v := range update.NextSyncCommitteeBranch {
		nextSyncCommitteeBranch[i] = make([]byte, 32)
		nextSyncCommitteeBranch[i], err = decodeHex(v)
		if err != nil {
			return false, err
		}
	}
	valid := trie.VerifyMerkleProof(attestedHeader.GetStateRoot(), nextSyncCommitteeRoot[:], 23, nextSyncCommitteeBranch)
	if !valid {
		return false, nil
	}

	finalityBranch := make([][]byte, len(update.FinalityBranch))
	for i, v := range update.FinalityBranch {
		finalityBranch[i] = make([]byte, 32)
		finalityBranch[i], err = decodeHex(v)
		if err != nil {
			return false, err
		}
	}
	valid = trie.VerifyMerkleProof(attestedHeader.GetStateRoot(), finalizedHeaderRoot[:], 105, finalityBranch)
	if !valid {
		return false, nil
	}

	var pubkeys []bls.PublicKey
	aggregateBytes, err := decodeHex(update.SyncAggregate.SyncCommitteeBits)
	if err != nil {
		return false, err
	}

	aggregateBits := bitfield.Bitvector512(aggregateBytes)
	for i := uint64(0); i < aggregateBits.Len(); i++ {
		if aggregateBits.BitAt(i) {
			pubKey, err := bls.PublicKeyFromBytes(currentSyncCommittee.Pubkeys[i])
			if err != nil {
				return false, err
			}
			pubkeys = append(pubkeys, pubKey)
		}
	}

	sigBytes, err := decodeHex(update.SyncAggregate.SyncCommitteeSignature)
	if err != nil {
		return false, err
	}
	sig, err := bls.SignatureFromBytes(sigBytes)
	if err != nil {
		return false, err
	}

	signingRoot, err := signing.ComputeSigningRoot(attestedHeader, domain)
	return sig.FastAggregateVerify(pubkeys, signingRoot), nil
}

func decodeHex(hexString string) ([]byte, error) {
	if hexString[0:2] == "0x" {
		hexString = hexString[2:]
	}
	decoded, err := hex.DecodeString(hexString)
	if err != nil {
		return nil, err
	}
	return decoded, nil
}
