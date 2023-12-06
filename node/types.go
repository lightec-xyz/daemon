package node

import (
	"bytes"
	"fmt"
	"strconv"
)

type DepositTx struct {
	TxId    string
	TxIndex int
	EthAddr string
	Amount  string
}

type RedeemTx struct {
	Inputs  []TxIn
	Outputs []TxOut
	TxIndex uint32
	TxId    string
}

func (rt *RedeemTx) String() string {
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

type TxIn struct {
	TxId  string
	Index uint32
}

type TxOut struct {
	Value    int64
	PkScript []byte
}

type TxProof struct {
	PTxId  string      `json:"pTxId"`
	Proof  string      `json:"proof"`
	Status ProofStatus `json:"status"`
	TxId   string      `json:"txId"`
	PType  string      `json:"type"`
	ToAddr string      `json:"toAddr"`
	Amount string      `json:"amount"`
	Msg    string      `json:"msg"`
}

// todo

type ProofRequest struct {
	// redeem
	Inputs  []TxIn  `json:"inputs"`
	Outputs []TxOut `json:"outputs"`

	// deposit
	Amount  string `json:"amount"`
	EthAddr string `json:"ethAddr"`
	Vout    int    `json:"index"`

	TxId  string `json:"txId"`
	PType string `json:"type"`
	Proof string `json:"proof"`
	Msg   string `json:"msg"`
}

// todo
type ProofResponse struct {
	// redeem
	Inputs  []TxIn  `json:"inputs"`
	Outputs []TxOut `json:"outputs"`

	// deposit
	Amount  string `json:"amount"`
	Vout    int    `json:"index"`
	EthAddr string `json:"ethAddr"`

	TxId  string `json:"txId"`
	PType string `json:"type"`
	Proof string `json:"proof"`
	Msg   string `json:"msg"`
}

func (resp *ProofResponse) String() string {
	var buf bytes.Buffer
	buf.WriteString("inputs:[")
	for _, vin := range resp.Inputs {
		buf.WriteString(vin.TxId)
		buf.WriteString(":")
		buf.WriteString(strconv.Itoa(int(vin.Index)))
		buf.WriteString(",")
	}
	buf.WriteString("]")
	buf.WriteString("outputs:[")
	for _, out := range resp.Outputs {
		buf.WriteString(fmt.Sprintf("%x", out.PkScript))
		buf.WriteString(":")
		buf.WriteString(fmt.Sprintf("%v", out.Value))
		buf.WriteString(",")
	}
	buf.WriteString("]")
	return buf.String()

}
