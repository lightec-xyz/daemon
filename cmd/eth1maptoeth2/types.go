package main

import (
	"errors"
	"strconv"

	"github.com/ethereum/go-ethereum/common/hexutil"
	apiclient "github.com/lightec-xyz/provers/utils/api-client"
	"github.com/prysmaticlabs/prysm/v5/api/client/beacon"
	ethpb "github.com/prysmaticlabs/prysm/v5/proto/prysm/v1alpha1"
)

type Eth1MapToEth2 struct {
	BlockNumber uint64 `json:"blockNumber"`
	BlockHash   string `json:"blockHash"`
	BlockSlot   uint64 `json:"blockSlot"`
	BlockRoot   string `json:"blockRoot"`
}

func GetHeadSlot(cl *apiclient.Client) (int, error) {
	headerResp, err := cl.GetBlockHeader(beacon.IdHead)
	if err != nil {
		return 0, err
	}

	if headerResp == nil {
		return 0, errors.New("headerResp is nil")
	}

	headSlot, err := strconv.Atoi(headerResp.Data.Header.Message.Slot)
	if err != nil {
		return 0, err
	}

	return headSlot, nil
}

func GetEth1MapToEth2(cl *apiclient.Client, slot int) (*Eth1MapToEth2, error) {
	blockId := beacon.StateOrBlockId(strconv.Itoa(slot))

	bb, err := cl.GetBlindedBlock(blockId)
	if err != nil {
		return nil, err
	}

	blk := &ethpb.SignedBlindedBeaconBlockBellatrix{}
	err = blk.UnmarshalSSZ(bb)
	if err != nil {
		return nil, err
	}

	blockRoot, err := blk.Block.HashTreeRoot()
	if err != nil {
		return nil, err
	}

	return &Eth1MapToEth2{
		BlockNumber: blk.Block.Body.ExecutionPayloadHeader.BlockNumber,
		BlockHash:   hexutil.Encode(blk.Block.Body.ExecutionPayloadHeader.BlockHash),
		BlockSlot:   uint64(blk.Block.Slot),
		BlockRoot:   hexutil.Encode(blockRoot[:]),
	}, nil
}
