package common

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/lightec-xyz/provers/circuits/utils"
	"github.com/lightec-xyz/provers/common"
	"github.com/prysmaticlabs/go-bitfield"
	"github.com/prysmaticlabs/prysm/v5/api/server/structs"
	"github.com/prysmaticlabs/prysm/v5/beacon-chain/core/signing"
	"github.com/prysmaticlabs/prysm/v5/container/trie"
	"github.com/prysmaticlabs/prysm/v5/crypto/bls"
)

const (
	SignatureThreshold = 2
)

var (
	NotEnoughSigner  = fmt.Errorf("not enough signer")
	InvalidSignature = fmt.Errorf("invalid signature")
)

type SyncCommitteeUpdate struct {
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

func verifySyncCommitteeUpdate(update *SyncCommitteeUpdate) (bool, error) {
	domain, err := utils.BuildDomainBytes(update.Version)
	if err != nil {
		return false, err
	}

	nextSyncCommitteeMerkleIndex := common.NextSyncCommitteeIndex + 1<<common.NextSyncCommitteeDepth
	finalityMerkleIndex := common.FinalizedHeaderIndex + 1<<common.FinalizedHeaderDepth

	if update.Version == common.Electra {
		nextSyncCommitteeMerkleIndex = common.NextSyncCommitteeIndexElectra + 1<<common.NextSyncCommitteeDepthElectra
		finalityMerkleIndex = common.FinalizedHeaderIndexElectra + 1<<common.FinalizedHeaderDepthElectra
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
	valid := trie.VerifyMerkleProof(attestedHeader.GetStateRoot(), nextSyncCommitteeRoot[:], uint64(nextSyncCommitteeMerkleIndex), nextSyncCommitteeBranch)
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
	valid = trie.VerifyMerkleProof(attestedHeader.GetStateRoot(), finalizedHeaderRoot[:], uint64(finalityMerkleIndex), finalityBranch)
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

	signingRoot, err := signing.ComputeSigningRoot(attestedHeader, domain[:])
	return sig.FastAggregateVerify(pubkeys, signingRoot), nil
}

func VerifyLightClientUpdate(update interface{}) (bool, error) {
	var innerUpdate SyncCommitteeUpdate
	err := ParseObj(update, &innerUpdate)
	if err != nil {
		return false, err
	}
	return verifySyncCommitteeUpdate(&innerUpdate)
}

// VerifyFinalityUpdateSignature check the finalityUpdate signature is valid or not,
// issue: FinalityUpdate_2189440 and FinalityUpdate_2193440
func VerifyFinalityUpdateSignature(update *LightClientFinalityUpdateEvent, currentSyncCommittee *SyncCommittee) (bool, error) {
	var domain []byte
	switch update.Version {
	case "bellatrix":
		domain, _ = decodeHex("0700000069b7d97441dbd33e5ee5b4cb8fc8b08d6a58a7274b6e6daf19ef4ca7")
	case "capella":
		domain, _ = decodeHex("0700000017e2dad36f1d3595152042a9ad23430197557e2e7e82bc7f7fc72972")
	case "deneb":
		domain, _ = decodeHex("0700000069ae0e9900d509b38350c53915fccde15c6ef44214aa1b5bdec34d3a")
	case "electra":
		domain, _ = decodeHex("07000000019e21ada5c73dd2b07fd515e7cd6d5f1eeb22e1fc0cfcfac4e03667")
	default:
		panic("unknown version")
	}

	attestedHeader, err := update.Data.AttestedHeader.ToConsensus()
	if err != nil {
		return false, err
	}

	finalizedHeader, err := update.Data.FinalizedHeader.ToConsensus()
	if err != nil {
		return false, err
	}

	finalizedHeaderRoot, err := finalizedHeader.HashTreeRoot()
	if err != nil {
		return false, err
	}

	sc, err := currentSyncCommittee.ToConsensus()
	if err != nil {
		return false, err
	}

	finalityMerkleIndex := common.FinalizedHeaderIndex + 1<<common.FinalizedHeaderDepth
	if update.Version == common.Electra {
		finalityMerkleIndex = common.FinalizedHeaderIndexElectra + 1<<common.FinalizedHeaderDepthElectra
	}

	finalityBranch := make([][]byte, len(update.Data.FinalityBranch))
	for i, v := range update.Data.FinalityBranch {
		finalityBranch[i] = make([]byte, 32)
		finalityBranch[i], err = decodeHex(v)
		if err != nil {
			return false, err
		}
	}
	valid := trie.VerifyMerkleProof(attestedHeader.GetStateRoot(), finalizedHeaderRoot[:], uint64(finalityMerkleIndex), finalityBranch)
	if !valid {
		return false, nil
	}

	var pubkeys []bls.PublicKey
	aggregateBytes, err := decodeHex(update.Data.SyncAggregate.SyncCommitteeBits)
	if err != nil {
		return false, err
	}

	nbSigner := 0
	aggregateBits := bitfield.Bitvector512(aggregateBytes)
	for i := uint64(0); i < aggregateBits.Len(); i++ {
		if aggregateBits.BitAt(i) {
			pubKey, err := bls.PublicKeyFromBytes(sc.Pubkeys[i])
			if err != nil {
				return false, err
			}
			pubkeys = append(pubkeys, pubKey)
			nbSigner++
		}
	}

	sigBytes, err := decodeHex(update.Data.SyncAggregate.SyncCommitteeSignature)
	if err != nil {
		return false, err
	}
	sig, err := bls.SignatureFromBytes(sigBytes)
	if err != nil {
		return false, err
	}

	signingRoot, err := signing.ComputeSigningRoot(attestedHeader, domain)
	valid = sig.FastAggregateVerify(pubkeys, signingRoot)
	if !valid {
		return false, InvalidSignature
	}
	if nbSigner < SignatureThreshold {
		return false, NotEnoughSigner
	}
	return true, nil
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

func Md5(data []byte) []byte {
	ret := md5.Sum(data)
	return ret[:]
}
func HexMd5(data []byte) string {
	return hex.EncodeToString(Md5(data))
}
