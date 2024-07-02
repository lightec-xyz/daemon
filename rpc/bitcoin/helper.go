package bitcoin

import "github.com/blockcypher/gobcy/v2"

// todo just for test

func BroadcastTx(txHex string) (string, error) {
	bc := gobcy.API{
		Token: "46ef69aa6c2349bc9a38fb5b6ae6080c",
		Coin:  "btc",
		Chain: "test3",
	}
	trans, err := bc.PushTX(txHex)
	if err != nil {
		return "", err
	}
	return trans.Trans.Hash, nil
}
