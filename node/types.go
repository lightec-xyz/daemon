package node

import (
	"bytes"
	"fmt"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc/ethereum"
	"github.com/lightec-xyz/daemon/store"
	"strconv"
	"sync"
)

type ProofStatus int

const (
	ProofDefault ProofStatus = iota
	ProofPending
	ProofSuccess
	ProofFailed
)

type BitcoinTx struct {
	EthAddr string
	Amount  int64 // btc

	EthTxHash string
	Height    int64
	BlockHash string
	TxId      string
	Utxos     []Utxo
	TxType    int
}

type EthereumTx struct {
	Height    int64
	BlockHash string
	Inputs    []Utxo
	Outputs   []TxOut
	TxId      string
}

func (rt *EthereumTx) String() string {
	var buf bytes.Buffer
	buf.WriteString("inputs:[")
	for _, vin := range rt.Inputs {
		buf.WriteString(vin.TxId)
		buf.WriteString(":")
		buf.WriteString(strconv.Itoa(int(vin.Index)))
		buf.WriteString(",")
	}
	buf.WriteString("]")
	buf.WriteString("outputs:[")
	for _, out := range rt.Outputs {
		buf.WriteString(fmt.Sprintf("%x", out.PkScript))
		buf.WriteString(":")
		buf.WriteString(fmt.Sprintf("%v", out.Value))
		buf.WriteString(",")
	}
	buf.WriteString("]")
	return buf.String()

}

type Utxo struct {
	TxId  string `json:"txId"`
	Index uint32 `json:"index"`
}

type TxOut struct {
	Value    int64
	PkScript []byte
}

type TxProof struct {
	Height    int64       `json:"height"`
	BlockHash string      `json:"blockHash"`
	TxId      string      `json:"txId"`
	ProofType string      `json:"type"`
	Proof     string      `json:"proof"`
	Msg       string      `json:"msg"`
	Status    ProofStatus `json:"status"`
}

// todo
type ProofRequest struct {
	// redeem
	Inputs  []Utxo  `json:"inputs"`
	Outputs []TxOut `json:"outputs"`
	BtcTxId string  `json:"btcTxId"`

	// deposit
	Utxos   []Utxo
	Amount  int64  `json:"amount"`
	EthAddr string `json:"ethAddr"`

	Height    int64  `json:"height"`
	BlockHash string `json:"blockHash"`
	TxId      string `json:"txId"`
	ProofType string `json:"type"`
	Proof     string `json:"proof"`
	Msg       string `json:"msg"`
}

func (req *ProofRequest) String() string {
	if req.ProofType == Deposit {
		return fmt.Sprintf("txType:%v,txid: %v, utxos:%v, amount:%v, ethAddr:%v", req.ProofType, req.TxId, req.Utxos, req.Amount, req.EthAddr)
	} else if req.ProofType == Redeem {
		return fmt.Sprintf("txType:%v,txid:%v, utxos:%v, outputs: %v", req.ProofType, req.TxId, formatUtxo(req.Inputs), formatOut(req.Outputs))
	}
	return ""
}

// todo
type ProofResponse struct {
	// redeem
	Inputs  []Utxo  `json:"inputs"`
	Outputs []TxOut `json:"outputs"`
	BtcTxId string  `json:"btcTxId"`

	// deposit
	Utxos   []Utxo
	Amount  int64  `json:"amount"`
	EthAddr string `json:"ethAddr"`

	Height    int64  `json:"height"`
	BlockHash string `json:"blockHash"`
	TxId      string `json:"txId"`
	ProofType string `json:"type"`
	Proof     string `json:"proof"`
	Msg       string `json:"msg"`
	Status    int    `json:"status"`
}

func (resp *ProofResponse) String() string {
	if resp.ProofType == Deposit {
		return fmt.Sprintf("txType:%v, utxos:%v, amount:%v, ethAddr:%v,statrus: %v", resp.ProofType, resp.Utxos, resp.Amount, resp.EthAddr, resp.Status)
	} else if resp.ProofType == Redeem {
		return fmt.Sprintf("txType:%v, utxos:%v, outputs: %v,status:%v", resp.ProofType, formatUtxo(resp.Inputs), formatOut(resp.Outputs), resp.Status)
	}
	return ""
}

func formatUtxo(utxos []Utxo) string {
	var buf bytes.Buffer
	for _, vin := range utxos {
		buf.WriteString(vin.TxId)
		buf.WriteString(":")
		buf.WriteString(strconv.Itoa(int(vin.Index)))
		buf.WriteString(",")
	}
	return buf.String()
}
func formatOut(outputs []TxOut) string {
	var buf bytes.Buffer
	for _, out := range outputs {
		buf.WriteString(fmt.Sprintf("%x", out.PkScript))
		buf.WriteString(":")
		buf.WriteString(fmt.Sprintf("%v", out.Value))
		buf.WriteString(",")
	}
	return buf.String()
}

type NonceManager struct {
	sync.Mutex
}

func NewNonceManager() *NonceManager {
	return &NonceManager{}
}

func (nm *NonceManager) GetNonce(address string, ethClient *ethereum.Client, store store.IStore) (uint64, error) {
	//todo
	nm.Lock()
	defer nm.Unlock()
	if ethClient == nil {
		return 0, fmt.Errorf("ethClient is nil")
	}
	chainNonce, err := ethClient.GetNonce(address)
	if err != nil {
		logger.Error("nonce manager get nonce error: %v %v", address, err)
		return 0, err
	}
	hasObj, err := store.HasObj(address)
	if err != nil {
		return 0, err
	}
	if hasObj {
		var nonce uint64
		err := store.GetObj(address, &nonce)
		if err != nil {
			return 0, err
		}
		if chainNonce <= nonce {
			return nonce + 1, nil
		}
		return chainNonce, nil

	} else {
		return chainNonce, nil
	}

}
