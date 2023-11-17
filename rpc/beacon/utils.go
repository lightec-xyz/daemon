package beacon

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	apiclient "github.com/lightec-xyz/provers/utils/api-client"
	"github.com/prysmaticlabs/prysm/v5/api/client/beacon"
	"github.com/prysmaticlabs/prysm/v5/encoding/ssz/detect"
	"strconv"
)

type Eth1MapToEth2 struct {
	BlockNumber uint64 `json:"blockNumber"`
	BlockHash   string `json:"blockHash"`
	BlockSlot   uint64 `json:"blockSlot"`
	BlockRoot   string `json:"blockRoot"`
}

func GetFinalizedHeadSlot(cl *apiclient.Client) (uint64, error) {
	headerResp, err := cl.GetBlockHeader(beacon.IdFinalized)
	if err != nil {
		return 0, err
	}

	if headerResp == nil {
		return 0, errors.New("headerResp is nil")
	}

	headSlot, err := strconv.ParseUint(headerResp.Data.Header.Message.Slot, 10, 64)
	if err != nil {
		return 0, err
	}

	return headSlot, nil
}

func GetEth1MapToEth2(cl *apiclient.Client, slot uint64) (*Eth1MapToEth2, error) {
	blockId := beacon.StateOrBlockId(fmt.Sprintf("%d", slot))

	bb, err := cl.GetBlindedBlock(blockId)
	if err != nil {
		return nil, err
	}

	vu, err := detect.FromBlock(bb)
	if err != nil {
		return nil, err
	}

	blk, err := vu.UnmarshalBlindedBeaconBlock(bb)
	if err != nil {
		return nil, err
	}

	blockRoot, err := blk.Block().HashTreeRoot()
	if err != nil {
		return nil, err
	}

	executionData, err := blk.Block().Body().Execution()
	if err != nil {
		return nil, err
	}

	return &Eth1MapToEth2{
		BlockNumber: executionData.BlockNumber(),
		BlockHash:   hexutil.Encode(executionData.BlockHash()),
		BlockSlot:   uint64(blk.Block().Slot()),
		BlockRoot:   hexutil.Encode(blockRoot[:]),
	}, nil
}
