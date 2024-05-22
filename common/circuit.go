package common

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/prysmaticlabs/go-bitfield"
	"github.com/prysmaticlabs/prysm/v5/api/server/structs"
	"github.com/prysmaticlabs/prysm/v5/beacon-chain/core/signing"
	"github.com/prysmaticlabs/prysm/v5/container/trie"
	"github.com/prysmaticlabs/prysm/v5/crypto/bls"
	"os"
)

func CheckZkParametersMd5(zkDir string, list []*Parameters) error {
	for _, item := range list {
		path := zkDir + "/" + item.FileName
		fileBytes, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read zk file error: %v %v", path, err)
		}
		fileMd5 := HexMd5(fileBytes)
		if fileMd5 != item.Hash {
			return fmt.Errorf("check zk md5 not match: %v %v %v", path, fileMd5, item.Hash)
		}
	}
	return nil
}

// todo

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

// todo
func VerifyLightClientUpdate(update interface{}) (bool, error) {
	var innerUpdate LightClientUpdateInfo
	err := ParseObj(update, &innerUpdate)
	if err != nil {
		return false, err
	}
	return verifyLightClientUpdateInfo(&innerUpdate)
}

func verifyLightClientUpdateInfo(update *LightClientUpdateInfo) (bool, error) {
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

type Parameters struct {
	FileName string `json:"file"`
	Hash     string `json:"md5"`
}

func Md5(data []byte) []byte {
	ret := md5.Sum(data)
	return ret[:]
}
func HexMd5(data []byte) string {
	return hex.EncodeToString(Md5(data))
}

func FileMd5(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(Md5(data)), nil
}

// todo ,get from server
const ParametersStr = `[
  {
    "file": "beacon_header_finality.ccs",
    "md5": "40b39c50429cc9e6d5176f29b5f973ba"
  },
  {
    "file": "beacon_header_finality.pk",
    "md5": "401c0cde4cbca094cd1c3c83db1780e3"
  },
  {
    "file": "beacon_header_finality.vk",
    "md5": "f115daf39cfde8fd4e9243c04bf1db36"
  },
  {
    "file": "beacon_header_inner_0.ccs",
    "md5": "a83199a36a3049495cfa987192420825"
  },
  {
    "file": "beacon_header_inner_0.pk",
    "md5": "39183badc980a57db5426dc66a9d7f21"
  },
  {
    "file": "beacon_header_inner_0.vk",
    "md5": "a55d2f6525e526f84bd529ab6678c2a8"
  },
  {
    "file": "beacon_header_inner_1.ccs",
    "md5": "f84464c51ac419080bc2e067ea87363a"
  },
  {
    "file": "beacon_header_inner_1.pk",
    "md5": "90591533b7b58db177523361bc18b7f6"
  },
  {
    "file": "beacon_header_inner_1.vk",
    "md5": "ae36aeb36d5677eb3f29a9e4e41aa4fc"
  },
  {
    "file": "beacon_header_inner_10.ccs",
    "md5": "7e8b8f31e34a93ace2d91940f6154417"
  },
  {
    "file": "beacon_header_inner_10.pk",
    "md5": "bff7e033d40ab035480e6b425f4391b6"
  },
  {
    "file": "beacon_header_inner_10.vk",
    "md5": "e2d85bb25172a361e7cf5378ceeb6af3"
  },
  {
    "file": "beacon_header_inner_11.ccs",
    "md5": "c89b7e2300905535684039feb60d7b61"
  },
  {
    "file": "beacon_header_inner_11.pk",
    "md5": "7a855e9bb8fef4261ecbf0ae7177d852"
  },
  {
    "file": "beacon_header_inner_11.vk",
    "md5": "910883b057fef423b6bdd44d7c870861"
  },
  {
    "file": "beacon_header_inner_12.ccs",
    "md5": "5fa431fb082eb369157c73bfd39b0933"
  },
  {
    "file": "beacon_header_inner_12.pk",
    "md5": "4d9767bb0a5a0dba39f0f5cd5b622d55"
  },
  {
    "file": "beacon_header_inner_12.vk",
    "md5": "421fa4e7810d6f90d006c6c20d025e38"
  },
  {
    "file": "beacon_header_inner_13.ccs",
    "md5": "a569130d8bc3a0c9554809a4c4b178b8"
  },
  {
    "file": "beacon_header_inner_13.pk",
    "md5": "5846f9764f1ba150674bad4b7e759315"
  },
  {
    "file": "beacon_header_inner_13.vk",
    "md5": "bebea568d1607b7213b069b2404d0bb0"
  },
  {
    "file": "beacon_header_inner_14.ccs",
    "md5": "8e1422713e90a80ffb7f87a173384b7e"
  },
  {
    "file": "beacon_header_inner_14.pk",
    "md5": "1afc4e30c9d1e8364aa5b8eb5ddfefbf"
  },
  {
    "file": "beacon_header_inner_14.vk",
    "md5": "dd84ba6c4e5d0e3dec6498dbd3c6db7a"
  },
  {
    "file": "beacon_header_inner_15.ccs",
    "md5": "082db188ffe0fc99145340f1725cd2d9"
  },
  {
    "file": "beacon_header_inner_15.pk",
    "md5": "490d5af07f85498c2e733adce39f6b14"
  },
  {
    "file": "beacon_header_inner_15.vk",
    "md5": "91db09ea064b42b67626e65b89453aa1"
  },
  {
    "file": "beacon_header_inner_16.ccs",
    "md5": "e0b2e710cc61896af68262e6d7e6707e"
  },
  {
    "file": "beacon_header_inner_16.pk",
    "md5": "fd5e361296dfaaacedc19d17f7aff342"
  },
  {
    "file": "beacon_header_inner_16.vk",
    "md5": "feaa76d45462c759950c393151d65114"
  },
  {
    "file": "beacon_header_inner_17.ccs",
    "md5": "84c7c4af497fd665f1d162749da9c03b"
  },
  {
    "file": "beacon_header_inner_17.pk",
    "md5": "4d0cce8d70ac2b7ae417bb148790cd19"
  },
  {
    "file": "beacon_header_inner_17.vk",
    "md5": "e43dae772ab4ba6251012139818513fb"
  },
  {
    "file": "beacon_header_inner_18.ccs",
    "md5": "fb1645f1415fb36a59c0ba405aea96c0"
  },
  {
    "file": "beacon_header_inner_18.pk",
    "md5": "5ba18021711357b23520f63b9c83a961"
  },
  {
    "file": "beacon_header_inner_18.vk",
    "md5": "c52aa34ee1c1765c1911b734205851c1"
  },
  {
    "file": "beacon_header_inner_19.ccs",
    "md5": "01bd3588de8e6f6895bab9b19623959c"
  },
  {
    "file": "beacon_header_inner_19.pk",
    "md5": "306e06783348c11973e21f72b42e7920"
  },
  {
    "file": "beacon_header_inner_19.vk",
    "md5": "75a48b514d838cb229089bed8c11351b"
  },
  {
    "file": "beacon_header_inner_2.ccs",
    "md5": "025e4b013bda6c11215235f2978ffda9"
  },
  {
    "file": "beacon_header_inner_2.pk",
    "md5": "cb5c869246ead7a8d8b780e182df29d3"
  },
  {
    "file": "beacon_header_inner_2.vk",
    "md5": "b25698d0de2d0cb73b3d315305d413d1"
  },
  {
    "file": "beacon_header_inner_20.ccs",
    "md5": "8a592c775bb82c05ae6863a939314e05"
  },
  {
    "file": "beacon_header_inner_20.pk",
    "md5": "425caea54ff0a23c456148fad67333a2"
  },
  {
    "file": "beacon_header_inner_20.vk",
    "md5": "91c03b6cb4013332b4ae14b766adb65e"
  },
  {
    "file": "beacon_header_inner_21.ccs",
    "md5": "63e5ebc564b722f412743f86cf72a161"
  },
  {
    "file": "beacon_header_inner_21.pk",
    "md5": "55f83aa7d779648ea454e08633c3cacf"
  },
  {
    "file": "beacon_header_inner_21.vk",
    "md5": "35f0f4ada2d3e587371ce867011d8275"
  },
  {
    "file": "beacon_header_inner_22.ccs",
    "md5": "25ca5f6670229f131f141495f14eddaf"
  },
  {
    "file": "beacon_header_inner_22.pk",
    "md5": "578a1fca8405035aecd34eeee54288f1"
  },
  {
    "file": "beacon_header_inner_22.vk",
    "md5": "77718590991b25b875ee743a507cef52"
  },
  {
    "file": "beacon_header_inner_23.ccs",
    "md5": "0994506e6c285ec9d02e636770d56006"
  },
  {
    "file": "beacon_header_inner_23.pk",
    "md5": "9def34f5059b30fa4df1a87dd6928f3e"
  },
  {
    "file": "beacon_header_inner_23.vk",
    "md5": "411baa57d2f7af67758352273ac9da4a"
  },
  {
    "file": "beacon_header_inner_24.ccs",
    "md5": "aad1fd185d8dc26449dbb7a3c887ed55"
  },
  {
    "file": "beacon_header_inner_24.pk",
    "md5": "f2844f505fc3c3e9648ca783410a9017"
  },
  {
    "file": "beacon_header_inner_24.vk",
    "md5": "5f6a13a69dee4da6a1b6d3adae5e8fc8"
  },
  {
    "file": "beacon_header_inner_25.ccs",
    "md5": "b92710c00645e8994445763e3f793a6d"
  },
  {
    "file": "beacon_header_inner_25.pk",
    "md5": "28e1a412fce11a9b5ab8e083e7cc0fc4"
  },
  {
    "file": "beacon_header_inner_25.vk",
    "md5": "36e2b043de5fd6edfcbcb0b8eff19e29"
  },
  {
    "file": "beacon_header_inner_26.ccs",
    "md5": "a18764d574e77771f240e2e07991375a"
  },
  {
    "file": "beacon_header_inner_26.pk",
    "md5": "be1901c3b6aad972d30463ed05b97230"
  },
  {
    "file": "beacon_header_inner_26.vk",
    "md5": "a7723662da363ad8568ba034ab408cfe"
  },
  {
    "file": "beacon_header_inner_27.ccs",
    "md5": "116da202d3178ab8ec57d87767785596"
  },
  {
    "file": "beacon_header_inner_27.pk",
    "md5": "6c5e4ac3fd1c8da39ed6838360725801"
  },
  {
    "file": "beacon_header_inner_27.vk",
    "md5": "5a820894f47409bfcdc629557572ed39"
  },
  {
    "file": "beacon_header_inner_28.ccs",
    "md5": "53f10c20952358e44aff4edf0ea146c0"
  },
  {
    "file": "beacon_header_inner_28.pk",
    "md5": "07a2b7950593f7fdd250468dea2996a4"
  },
  {
    "file": "beacon_header_inner_28.vk",
    "md5": "5756c8df8619e30b7b30e46e1026055f"
  },
  {
    "file": "beacon_header_inner_29.ccs",
    "md5": "d0401a5597c9424b35a7e001202bb25d"
  },
  {
    "file": "beacon_header_inner_29.pk",
    "md5": "496ab30388d2c8f57ecfd6503b309c80"
  },
  {
    "file": "beacon_header_inner_29.vk",
    "md5": "4ba10f47390863ecf70b7f4b626001cb"
  },
  {
    "file": "beacon_header_inner_3.ccs",
    "md5": "be98b78a43cca2bf1d4488ec85757bb8"
  },
  {
    "file": "beacon_header_inner_3.pk",
    "md5": "9d27f3887ce459b191efb5fb21461cd3"
  },
  {
    "file": "beacon_header_inner_3.vk",
    "md5": "16d62c1c6ef737e1e4d942bab84cf416"
  },
  {
    "file": "beacon_header_inner_30.ccs",
    "md5": "3eeceea0ef651a143eb62fd5cb19407e"
  },
  {
    "file": "beacon_header_inner_30.pk",
    "md5": "2b149cc0890cf10f4b9c803a59547238"
  },
  {
    "file": "beacon_header_inner_30.vk",
    "md5": "e1bf842fd9b95057d7df1199638e0e69"
  },
  {
    "file": "beacon_header_inner_31.ccs",
    "md5": "0820e7f781feea894ad7e7697ed1ef0d"
  },
  {
    "file": "beacon_header_inner_31.pk",
    "md5": "33e7ff747b6aa6a21a58d8abca4c5f11"
  },
  {
    "file": "beacon_header_inner_31.vk",
    "md5": "7b66ac7d4feb00997c8446253aa30b6c"
  },
  {
    "file": "beacon_header_inner_32.ccs",
    "md5": "c9c8fa2c55a5f2022337cb01bd3820e1"
  },
  {
    "file": "beacon_header_inner_32.pk",
    "md5": "c278a27e49a6e33ad7e87654c6a8352d"
  },
  {
    "file": "beacon_header_inner_32.vk",
    "md5": "2c806ce44b50ce98ed318620c121a2a1"
  },
  {
    "file": "beacon_header_inner_33.ccs",
    "md5": "766eb5b8be62fc54d5daae5de7393ff9"
  },
  {
    "file": "beacon_header_inner_33.pk",
    "md5": "cf12c8e404c6667efc908b6eb5f9aeda"
  },
  {
    "file": "beacon_header_inner_33.vk",
    "md5": "d2656db4a25ef76a205a84887ca57902"
  },
  {
    "file": "beacon_header_inner_34.ccs",
    "md5": "fa430adc13f08759f1f487d2584deeaf"
  },
  {
    "file": "beacon_header_inner_34.pk",
    "md5": "37d861c04b062db445020b4a2d05b986"
  },
  {
    "file": "beacon_header_inner_34.vk",
    "md5": "70903cbf3a96a0b693f840e45be667e7"
  },
  {
    "file": "beacon_header_inner_35.ccs",
    "md5": "3b912292bc9f41527717f3893458b418"
  },
  {
    "file": "beacon_header_inner_35.pk",
    "md5": "2811a6f30707a46ff5d7dfd64a3691a4"
  },
  {
    "file": "beacon_header_inner_35.vk",
    "md5": "3c278f63afa25bf8391123a24046c471"
  },
  {
    "file": "beacon_header_inner_36.ccs",
    "md5": "24c142e7b481e6470a015457e1f882e7"
  },
  {
    "file": "beacon_header_inner_36.pk",
    "md5": "455d7c7f6c808d0b04671ddb2558a66a"
  },
  {
    "file": "beacon_header_inner_36.vk",
    "md5": "14ad46d940012f390ecb43f2caea45b3"
  },
  {
    "file": "beacon_header_inner_37.ccs",
    "md5": "fde3c3589b5d220c94e67a3b44e93342"
  },
  {
    "file": "beacon_header_inner_37.pk",
    "md5": "560156c5e827ed17f6bb4225b238c342"
  },
  {
    "file": "beacon_header_inner_37.vk",
    "md5": "004c2fd221d1823c8a540fec715183ba"
  },
  {
    "file": "beacon_header_inner_38.ccs",
    "md5": "a4074e3994d6fc5ce46c642016e373c3"
  },
  {
    "file": "beacon_header_inner_38.pk",
    "md5": "30e0579a9331cc415e0544b8fa0c754d"
  },
  {
    "file": "beacon_header_inner_38.vk",
    "md5": "fa24957613df68370c961947b755aeda"
  },
  {
    "file": "beacon_header_inner_39.ccs",
    "md5": "37eff4fa8901344fe6500610a9c05c27"
  },
  {
    "file": "beacon_header_inner_39.pk",
    "md5": "58fe823d760c722df18b2e9db5798972"
  },
  {
    "file": "beacon_header_inner_39.vk",
    "md5": "e5a66dd650ee514cc1515d697905801b"
  },
  {
    "file": "beacon_header_inner_4.ccs",
    "md5": "28b11755183a85c9a7c9ea037fdbfd35"
  },
  {
    "file": "beacon_header_inner_4.pk",
    "md5": "8bde78daa1fad32e61a18d553ca0da62"
  },
  {
    "file": "beacon_header_inner_4.vk",
    "md5": "e4144683156c92102a22aefec78a3a7d"
  },
  {
    "file": "beacon_header_inner_40.ccs",
    "md5": "241c9026f6f4a29f19efd61989b6cc45"
  },
  {
    "file": "beacon_header_inner_40.pk",
    "md5": "6793aab24e8d991d70727333d1b7b1c9"
  },
  {
    "file": "beacon_header_inner_40.vk",
    "md5": "e7f699e8e9bb76ee07d40624bcc9a3dd"
  },
  {
    "file": "beacon_header_inner_41.ccs",
    "md5": "192c6d39a7921ae0cc27273f8c9d1956"
  },
  {
    "file": "beacon_header_inner_41.pk",
    "md5": "86e1a9ca370c0a936046945b61ccecec"
  },
  {
    "file": "beacon_header_inner_41.vk",
    "md5": "f0435d2049a580746acc94fbdbe6fd81"
  },
  {
    "file": "beacon_header_inner_42.ccs",
    "md5": "c28753f259e7347b1e4686d38dc16f4c"
  },
  {
    "file": "beacon_header_inner_42.pk",
    "md5": "1f6f10cdf360cb1b35845331af7f595c"
  },
  {
    "file": "beacon_header_inner_42.vk",
    "md5": "9bc8c949ea566d9bf5fde17eac6bb2f1"
  },
  {
    "file": "beacon_header_inner_43.ccs",
    "md5": "fb4905e5d3acf45a077335fb6841e6c6"
  },
  {
    "file": "beacon_header_inner_43.pk",
    "md5": "045d70d2b5e027f64e9083e84e130dd2"
  },
  {
    "file": "beacon_header_inner_43.vk",
    "md5": "8d9f37f83d157223b146f92c7e1833f9"
  },
  {
    "file": "beacon_header_inner_44.ccs",
    "md5": "7b74b9105b6f0518392432545a68010f"
  },
  {
    "file": "beacon_header_inner_44.pk",
    "md5": "957f8b380d6bfb2dff01f34908566fa6"
  },
  {
    "file": "beacon_header_inner_44.vk",
    "md5": "f3d9482203ef3b1d46f237c352cb55c0"
  },
  {
    "file": "beacon_header_inner_45.ccs",
    "md5": "dce1c0e269b4cb01bc27c4aa93dcc300"
  },
  {
    "file": "beacon_header_inner_45.pk",
    "md5": "cd20935dbd99bab94ccfd9a53e2bc913"
  },
  {
    "file": "beacon_header_inner_45.vk",
    "md5": "0bffaccf6325bea2af3b8c6dbc822244"
  },
  {
    "file": "beacon_header_inner_46.ccs",
    "md5": "2d777b58f7c70078c4c7cc71b5bf0063"
  },
  {
    "file": "beacon_header_inner_46.pk",
    "md5": "be05e1b58626f3bcd4e16f0288ab1301"
  },
  {
    "file": "beacon_header_inner_46.vk",
    "md5": "818208a06d4b5d518ca1df749a4cff23"
  },
  {
    "file": "beacon_header_inner_47.ccs",
    "md5": "4e4025669b7fdf5adbc3cbbce1c677e8"
  },
  {
    "file": "beacon_header_inner_47.pk",
    "md5": "043e2dbb298d7d871cdc5da0702f7eda"
  },
  {
    "file": "beacon_header_inner_47.vk",
    "md5": "8ab74ccafd908c849bca6d9c66805c66"
  },
  {
    "file": "beacon_header_inner_48.ccs",
    "md5": "c56152204123f0cf603b17e1a76cc57a"
  },
  {
    "file": "beacon_header_inner_48.pk",
    "md5": "9b1d83d6c8c632ac44b2ddba6e74fd14"
  },
  {
    "file": "beacon_header_inner_48.vk",
    "md5": "52d6a3220f848ca4343b4ee3ed17e656"
  },
  {
    "file": "beacon_header_inner_49.ccs",
    "md5": "c9d2a74ffdb2cc77e8741d41c4eebbb7"
  },
  {
    "file": "beacon_header_inner_49.pk",
    "md5": "719e36fce225fb16ba68aa0d32f9e60c"
  },
  {
    "file": "beacon_header_inner_49.vk",
    "md5": "fe00012b82e90130d4f433af776d5943"
  },
  {
    "file": "beacon_header_inner_5.ccs",
    "md5": "74d7b144d18e96aa2c3e4b0c24507f3a"
  },
  {
    "file": "beacon_header_inner_5.pk",
    "md5": "b9a1d4a1e294a510a8998a5de26b77eb"
  },
  {
    "file": "beacon_header_inner_5.vk",
    "md5": "8c8f9607b1cc78867314fa34f5f18fab"
  },
  {
    "file": "beacon_header_inner_50.ccs",
    "md5": "1e7cb392194acd2860c81bc65f2973ad"
  },
  {
    "file": "beacon_header_inner_50.pk",
    "md5": "bab6dc53758d5b463227624259a29428"
  },
  {
    "file": "beacon_header_inner_50.vk",
    "md5": "36ccc66b0b2215480a78e4358d8761e3"
  },
  {
    "file": "beacon_header_inner_51.ccs",
    "md5": "c1dbcc8bc5be580c2610c2fc9eaa30dc"
  },
  {
    "file": "beacon_header_inner_51.pk",
    "md5": "ee7ad8ff9068acd26136f29600b6299d"
  },
  {
    "file": "beacon_header_inner_51.vk",
    "md5": "0a3c37dcbe0c586c41b7a279ef039273"
  },
  {
    "file": "beacon_header_inner_52.ccs",
    "md5": "f0b051543a8956846b396be1a9faa527"
  },
  {
    "file": "beacon_header_inner_52.pk",
    "md5": "ab82147b781b02138a2474fe967b99c5"
  },
  {
    "file": "beacon_header_inner_52.vk",
    "md5": "b471dd92af04dc0eb7f5d9449e61e80f"
  },
  {
    "file": "beacon_header_inner_53.ccs",
    "md5": "7d792af778841b1b125f9a3127bb9526"
  },
  {
    "file": "beacon_header_inner_53.pk",
    "md5": "b77bc90d35843de22b3ea80c5d3724de"
  },
  {
    "file": "beacon_header_inner_53.vk",
    "md5": "f3949db0328ebba559c52c14d9b9b4f6"
  },
  {
    "file": "beacon_header_inner_54.ccs",
    "md5": "502eea2a99fa49a29e544409092268c6"
  },
  {
    "file": "beacon_header_inner_54.pk",
    "md5": "7f9eb085bb79c598e66701526bb26c90"
  },
  {
    "file": "beacon_header_inner_54.vk",
    "md5": "4ac6203aff696294f29c243c5b95ad96"
  },
  {
    "file": "beacon_header_inner_55.ccs",
    "md5": "8fcd4936d20fe4ae02919a645b6ca7b7"
  },
  {
    "file": "beacon_header_inner_55.pk",
    "md5": "5c22572588f7a9acf9bb2ca6dad7e161"
  },
  {
    "file": "beacon_header_inner_55.vk",
    "md5": "c8931e8f19c1e0168ac1ae4469c40da9"
  },
  {
    "file": "beacon_header_inner_56.ccs",
    "md5": "9a15f5c86b8b2e977a54ef90b9dc14f4"
  },
  {
    "file": "beacon_header_inner_56.pk",
    "md5": "9b494438b86f354734396b11f6e372f2"
  },
  {
    "file": "beacon_header_inner_56.vk",
    "md5": "d8502060f67fdeddcde3b43814daa8fa"
  },
  {
    "file": "beacon_header_inner_57.ccs",
    "md5": "b1c270ecfd377cc61b6efa9560ae0f27"
  },
  {
    "file": "beacon_header_inner_57.pk",
    "md5": "ead2eaab5c96832b51ef95b839d8ef92"
  },
  {
    "file": "beacon_header_inner_57.vk",
    "md5": "95886426bedd34bbe3a0f37a185d64ae"
  },
  {
    "file": "beacon_header_inner_58.ccs",
    "md5": "fc26e2ae2a80cd8adffb9db561e70ab1"
  },
  {
    "file": "beacon_header_inner_58.pk",
    "md5": "728ace3004f944d22b1fe6adaa3b6b38"
  },
  {
    "file": "beacon_header_inner_58.vk",
    "md5": "ab82ef86724ce02738a174677cc61ba8"
  },
  {
    "file": "beacon_header_inner_59.ccs",
    "md5": "a101efa80c660bddc20b1adb6151d59f"
  },
  {
    "file": "beacon_header_inner_59.pk",
    "md5": "28f9db881c9df61b22a5cc3326ec27e0"
  },
  {
    "file": "beacon_header_inner_59.vk",
    "md5": "e2bb575d2bb849d54973eb5b8c19b186"
  },
  {
    "file": "beacon_header_inner_6.ccs",
    "md5": "309c5e306155f7d71489ad6d5cbe6026"
  },
  {
    "file": "beacon_header_inner_6.pk",
    "md5": "81718b7106d7f06e396da963c6298b67"
  },
  {
    "file": "beacon_header_inner_6.vk",
    "md5": "2de03f668ae236e9861880db5e4822c8"
  },
  {
    "file": "beacon_header_inner_60.ccs",
    "md5": "df811a3526225542b72f77ef12ed8c31"
  },
  {
    "file": "beacon_header_inner_60.pk",
    "md5": "dc8bb8b5bde7a13098ccedeb7acc5b39"
  },
  {
    "file": "beacon_header_inner_60.vk",
    "md5": "b41696c4521b38d3d5ee67fce1e88b8b"
  },
  {
    "file": "beacon_header_inner_61.ccs",
    "md5": "e2b18c96fd339757a7d3d4651c8009fc"
  },
  {
    "file": "beacon_header_inner_61.pk",
    "md5": "c1cd363dde128fcad94aa60bea8416cd"
  },
  {
    "file": "beacon_header_inner_61.vk",
    "md5": "8d0e14fae77a1be689e3c14b2d09b068"
  },
  {
    "file": "beacon_header_inner_62.ccs",
    "md5": "31b816768dcfd4d1842d9bb69746255f"
  },
  {
    "file": "beacon_header_inner_62.pk",
    "md5": "8b92f7740f2af579bbe1f44a29a9fe7d"
  },
  {
    "file": "beacon_header_inner_62.vk",
    "md5": "3885d4ed6ee0f4dbff55318c0bad24fc"
  },
  {
    "file": "beacon_header_inner_63.ccs",
    "md5": "97ad819c56dd07a84965614dbaf883d9"
  },
  {
    "file": "beacon_header_inner_63.pk",
    "md5": "484742aed7f5312db86c5b2bc05047f8"
  },
  {
    "file": "beacon_header_inner_63.vk",
    "md5": "7171344a1c307cfc02288f83b1d627e6"
  },
  {
    "file": "beacon_header_inner_64.ccs",
    "md5": "28fd8c276657c415b511d74de1d2bb2d"
  },
  {
    "file": "beacon_header_inner_64.pk",
    "md5": "3f5a9fc46fc35c435b54e91f86abb361"
  },
  {
    "file": "beacon_header_inner_64.vk",
    "md5": "b8b589d58f0f17855924ae8490b329fd"
  },
  {
    "file": "beacon_header_inner_7.ccs",
    "md5": "86fdf53da86daa7f7a587ef67e47a623"
  },
  {
    "file": "beacon_header_inner_7.pk",
    "md5": "1e8f3a85a775bb2dfb1380f22951b5a3"
  },
  {
    "file": "beacon_header_inner_7.vk",
    "md5": "566ba726cb773071985510564ed7df17"
  },
  {
    "file": "beacon_header_inner_8.ccs",
    "md5": "92672417cf44574414e797dcb2e07791"
  },
  {
    "file": "beacon_header_inner_8.pk",
    "md5": "e68ff86fd70bf15adbcb0284ecccca6d"
  },
  {
    "file": "beacon_header_inner_8.vk",
    "md5": "79c05362ea35f41ed88b3e24e51297ef"
  },
  {
    "file": "beacon_header_inner_9.ccs",
    "md5": "3be0d2d2c7c0304d37fa95a8aa7e2f4a"
  },
  {
    "file": "beacon_header_inner_9.pk",
    "md5": "99c2e1484611aa8aa5d8a14b7acd2ee5"
  },
  {
    "file": "beacon_header_inner_9.vk",
    "md5": "16f2be46d1b7eac7db5011d6cc91f82a"
  },
  {
    "file": "beacon_header_outer.ccs",
    "md5": "ce7554d017b1db4ee0352b6eb0806cc9"
  },
  {
    "file": "beacon_header_outer.pk",
    "md5": "7d11e67613c2e8625f0d174853a9bf7a"
  },
  {
    "file": "beacon_header_outer.vk",
    "md5": "7acd798544d6f5a2fba0c4ba49ad6766"
  },
  {
    "file": "redeem.ccs",
    "md5": "7eb4fbbe5de7bd6f80151a64b00c24b4"
  },
  {
    "file": "redeem.pk",
    "md5": "65e802e4b46dbed95447ac01a5f3d6d5"
  },
  {
    "file": "redeem.sol",
    "md5": "2442088ac47e8915587fba8cfeb9141f"
  },
  {
    "file": "redeem.vk",
    "md5": "869ccbd24937f085af6d27216617ea5f"
  },
  {
    "file": "sc_genesis.ccs",
    "md5": "ffcc7616ebdc5d0ffcfc4e757084c2ad"
  },
  {
    "file": "sc_genesis.pk",
    "md5": "d81d9d245f9273e8b810acba8a2b785f"
  },
  {
    "file": "sc_genesis.vk",
    "md5": "76b8e9d5324d9847dd0bc7ab9ec64c70"
  },
  {
    "file": "sc_inner.ccs",
    "md5": "5f36bc9eb67d2b98fb6173fb3c822c65"
  },
  {
    "file": "sc_inner.pk",
    "md5": "795c5573691fbbb5f561cc096866a3b0"
  },
  {
    "file": "sc_inner.vk",
    "md5": "83eaaa7f896dd59add2875de6d1ed258"
  },
  {
    "file": "sc_outer.ccs",
    "md5": "9d41e1ae4b81053b2c30c75ba52b4849"
  },
  {
    "file": "sc_outer.pk",
    "md5": "725d0a533c0f8430bdad5aa7b37b3712"
  },
  {
    "file": "sc_outer.vk",
    "md5": "8ecbae071cd6309b00fad6698fed1569"
  },
  {
    "file": "sc_recursive.ccs",
    "md5": "98b53b40742845702c38cb63f0b1ff00"
  },
  {
    "file": "sc_recursive.pk",
    "md5": "bda797ff06a6b29caa92f4db8774fcd4"
  },
  {
    "file": "sc_recursive.vk",
    "md5": "6ce011e96db4e58e3be6c9f14a341115"
  },
  {
    "file": "sc_unit.ccs",
    "md5": "9ac80f79badc046353a05eb2a0e55037"
  },
  {
    "file": "sc_unit.pk",
    "md5": "75fa86bbce1c5914261a50ffd819dc14"
  },
  {
    "file": "sc_unit.vk",
    "md5": "afbafee91d6467b075de1569341c1a83"
  },
  {
    "file": "tx_in_eth2.ccs",
    "md5": "5eef33cbc824d697730151eb5f4f5054"
  },
  {
    "file": "tx_in_eth2.pk",
    "md5": "275d68e497af80faf9b09e8a18e02176"
  },
  {
    "file": "tx_in_eth2.sol",
    "md5": "079b9af78b8fc991351817da4ff5ca62"
  },
  {
    "file": "tx_in_eth2.vk",
    "md5": "c140af35089340b4a20ed69fb387f70c"
  },
  {
    "file": "grandrollup.ccs",
    "md5": "44a7057f1239714400fce8003b84a0fa"
  },
  {
    "file": "grandrollup.pk",
    "md5": "44f9b35c90995286e31b95b02ee9a1ca"
  },
  {
    "file": "grandrollup.vk",
    "md5": "da61a7a997ba0527aabbccbbae935964"
  },
  {
    "file": "grandrollup.sol",
    "md5": "0faa3ff3c2266d3acf75e7ab03316c98"
  }
]`
