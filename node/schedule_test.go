package node

import (
	"github.com/lightec-xyz/daemon/common"
	"testing"
)

func TestDemo01(t *testing.T) {
	start := 2015
	startIndex := (start / common.BtcUpperDistance) * common.BtcUpperDistance
	endIndex := startIndex + common.BtcUpperDistance
	t.Logf("startIndex: %v, endIndex: %v", startIndex, endIndex)

}
