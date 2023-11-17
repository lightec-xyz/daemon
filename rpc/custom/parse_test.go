package custom

import (
	"testing"
)

func TestParse_Parse(t *testing.T) {
	parse, err := NewParse("/Users/red/lworkspace/lightec/btcmainnetdata/headers.txt")
	if err != nil {
		t.Error(err)
	}
	count := 700000
	for index := 0; index < count; index++ {
		blockHash, err := parse.GetBlockHashByHeight(int64(index))
		if err != nil {
			t.Error(err)
		}
		header, err := parse.GetHeaderByHash(blockHash)
		if err != nil {
			t.Error(err)
		}
		header1, err := parse.GetHeaderByHeight(int64(index))
		if err != nil {
			t.Error(err)
		}
		t.Logf("%v,%v,%v,%v \n", index, blockHash, header, header1)

	}
}

//00000000000000000000121c1b4143d3c1ef82789a3e68c4073c0cdb776c3a8b
