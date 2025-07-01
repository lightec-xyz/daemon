package common

import (
	"fmt"
	"testing"
	"time"
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
func TestStrToTime(t *testing.T) {
	str := "2025-07-01 03:58:55"
	ti, err := time.Parse("2006-01-02 15:04:05", str)
	if err != nil {
		fmt.Println("parse error:", err)
		return
	}
	t.Log(ti.Unix())
}
