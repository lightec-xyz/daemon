package common

import (
	"fmt"
	"testing"
)

func Test_Reverse(t *testing.T) {
	var res = []byte{1, 2, 3, 4, 5}
	fmt.Println(ReverseBytes(res))
}

func TestBtcChainGenesisIndex(t *testing.T) {
	endIndex := 83929
	startIndex := ((endIndex-10)/2016 - 2) * 2016
	fmt.Printf("chainStartIndex: %d \n", startIndex)
	cpStartIndex := endIndex - 200
	fmt.Printf("checkPointStartIndex: %d \n", cpStartIndex)

}
