package p2p

func NewP2pMinerMsg(addr string, power int64, timestamp int64) *Msg {
	minerType := Msg_Miner
	msg := &Msg{
		Type: &minerType,
		Miner: &Miner{
			MinerAddr: &addr,
			Power:     &power,
		},
		Timestamp: &timestamp,
	}
	return msg
}
