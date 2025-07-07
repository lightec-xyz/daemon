package node

import "testing"

func TestBtcToSat(t *testing.T) {
	sat := BtcToSat(21.00000011645645645)
	t.Log(sat)

}

func TestUrlToken(t *testing.T) {
	token := getUrlToken("https://localhost:9003/eth_beacon/8c933202fbe8dbe6d63377a319b6020f4a4c35bb4424f6368f630b676b4fcc33")
	t.Log(token)

	token = getUrlToken("http://127.0.0.1:9003")
	t.Log(token)
}
