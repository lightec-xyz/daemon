package common

import (
	"fmt"
	"testing"
)

func TestBtcChainGenesisIndex(t *testing.T) {
	endIndex := 79768
	startIndex := ((endIndex-10)/2016 - 2) * 2016
	fmt.Printf("chainStartIndex: %d \n", startIndex)
	cpStartIndex := endIndex - 200
	fmt.Printf("checkPointStartIndex: %d \n", cpStartIndex)

}
