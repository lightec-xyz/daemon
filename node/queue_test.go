package node

import (
	"github.com/lightec-xyz/daemon/common"
	"testing"
)

func TestArrayQueue(t *testing.T) {
	queue := NewArrayQueue(sortRequest)
	for period := 396; period <= 400; period++ {
		for index := 0; index <= 7; index++ {
			queue.Push(common.NewProofRequest(common.SyncComInnerType, nil, uint64(period), uint64(index), 0, ""))
		}
	}
	for queue.Len() > 0 {
		request, ok := queue.Pop()
		if ok {
			t.Log(request.ProofId())
		}
	}
}
