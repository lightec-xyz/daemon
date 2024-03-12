package node

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/prysmaticlabs/prysm/v5/api/server/structs"
	"strconv"
	"sync"
)

type ZkProofType int

const (
	DepositTxType ZkProofType = iota + 1
	RedeemTxType
	VerifyTxType
	SyncComGenesisType
	SyncComUnitType
	SyncComRecursiveType
)

func (zkpr *ZkProofType) String() string {
	switch *zkpr {
	case DepositTxType:
		return "DepositTxType"
	case RedeemTxType:
		return "RedeemTxType"
	case VerifyTxType:
		return "VerifyTxType"
	case SyncComGenesisType:
		return "SyncComGenesisType"
	case SyncComUnitType:
		return "SyncComUnitType"
	case SyncComRecursiveType:
		return "SyncComRecursiveType"
	}
	return ""
}

type ZkProofRequest struct {
	reqType ZkProofType // 0: genesis proof, 1: unit proof, 2: recursive proof
	period  uint64
	data    interface{} // current request data
}

func (r *ZkProofRequest) String() string {
	return fmt.Sprintf("ZkProofRequest{reqType:%v,period:%v,data:%v}", r.reqType, r.period, r.data)

}

type ZkProofResponse struct {
	zkProofType ZkProofType // 0: genesis proof, 1: unit proof, 2: recursive proof
	period      uint64
	data        interface{}
	Status      ProofStatus
	body        []byte
	proof       []byte
}

func toDepositZkProofRequest(list []ProofRequest) ([]ZkProofRequest, error) {
	var result []ZkProofRequest
	for _, item := range list {
		result = append(result, ZkProofRequest{
			reqType: DepositTxType,
			data: DepositProofParam{
				Body: item,
			},
		})
	}
	return result, nil
}

func toRedeemZkProofRequest(list []ProofRequest) ([]ZkProofRequest, error) {
	var result []ZkProofRequest
	for _, item := range list {
		result = append(result, ZkProofRequest{
			reqType: RedeemTxType,
			data: RedeemProofParam{
				Body: item,
			},
		})
	}
	return result, nil
}

func (zkResp *ZkProofResponse) String() string {
	return fmt.Sprintf("zkProofType:%v period:%v proof:%v", zkResp.zkProofType, zkResp.period, zkResp.proof)
}

func (zkRep *ZkProofResponse) ParseDepositProof() (DepositProof, error) {
	if zkRep.zkProofType != DepositTxType {
		return DepositProof{}, fmt.Errorf("not deposit proof")
	}
	depositProof := DepositProof{}
	err := json.Unmarshal(zkRep.body, &depositProof)
	if err != nil {
		return DepositProof{}, err
	}
	return depositProof, nil
}

func (zkRep *ZkProofResponse) ParseRedeemProof() (RedeemProof, error) {
	if zkRep.zkProofType != RedeemTxType {
		return RedeemProof{}, fmt.Errorf("not redeem proof")
	}
	redeemProof := RedeemProof{}
	err := json.Unmarshal(zkRep.body, &redeemProof)
	if err != nil {
		return RedeemProof{}, err
	}
	return redeemProof, nil
}

type DepositProofParam struct {
	Version string
	Body    interface{}
}

type RedeemProofParam struct {
	Version string
	Body    interface{}
}

type VerifyProofParam struct {
	Version string
	Body    interface{}
}

type GenesisProofParam struct {
	Version string
	data    structs.LightClientBootstrapResponse
}

type UnitProofParam struct {
	Version   string
	update    []structs.LightClientUpdateWithVersion // todo
	preUpdate []structs.LightClientUpdateWithVersion
	genesis   structs.LightClientBootstrapResponse
	isGenesis bool
}

type RecursiveProofParam struct {
	Version           string
	unitProof         string
	preRecursiveProof string
	isGenesis         bool
}

type DepositProof struct {
	// redeem
	Inputs  []Utxo  `json:"inputs"`
	Outputs []TxOut `json:"outputs"`
	BtcTxId string  `json:"btcTxId"`

	// deposit
	Utxos   []Utxo
	Amount  int64  `json:"amount"`
	EthAddr string `json:"ethAddr"`

	// other
	Height    int64       `json:"height"`
	BlockHash string      `json:"blockHash"`
	TxId      string      `json:"txId"`
	ProofType ProofType   `json:"type"`
	Proof     string      `json:"proof"`
	Msg       string      `json:"msg"`
	Status    ProofStatus `json:"status"`
}

func (dp *DepositProof) String() string {
	return fmt.Sprintf("inputs:%v outputs:%v btcTxId:%v amount:%v ethAddr:%v height:%v blockHash:%v txId:%v type:%v proof:%v msg:%v status:%v",
		dp.Inputs, dp.Outputs, dp.BtcTxId, dp.Amount, dp.EthAddr, dp.Height, dp.BlockHash, dp.TxId, dp.ProofType, dp.Proof, dp.Msg, dp.Status)
}

type RedeemProof struct {
	// redeem
	Inputs  []Utxo  `json:"inputs"`
	Outputs []TxOut `json:"outputs"`
	BtcTxId string  `json:"btcTxId"`

	// deposit
	Utxos   []Utxo
	Amount  int64  `json:"amount"`
	EthAddr string `json:"ethAddr"`

	// other
	Height    int64       `json:"height"`
	BlockHash string      `json:"blockHash"`
	TxId      string      `json:"txId"`
	ProofType ProofType   `json:"type"`
	Proof     string      `json:"proof"`
	Msg       string      `json:"msg"`
	Status    ProofStatus `json:"status"`
}

func (dp *RedeemProof) String() string {
	return fmt.Sprintf("inputs:%v outputs:%v btcTxId:%v amount:%v ethAddr:%v height:%v blockHash:%v txId:%v type:%v proof:%v msg:%v status:%v",
		dp.Inputs, dp.Outputs, dp.BtcTxId, dp.Amount, dp.EthAddr, dp.Height, dp.BlockHash, dp.TxId, dp.ProofType, dp.Proof, dp.Msg, dp.Status)
}

type DownloadStatus int32

const (
	None        DownloadStatus = 0
	Downloading DownloadStatus = 1
	Done        DownloadStatus = 2
	Fail        DownloadStatus = 3
)

type FetchType int

const (
	GenesisUpdateType FetchType = iota + 1
	PeriodUpdateType
)

func (ft FetchType) String() string {
	switch ft {
	case GenesisUpdateType:
		return "GenesisUpdateType"
	case PeriodUpdateType:
		return "PeriodUpdateType"
	default:
		return "unknown"
	}
}

type FetchRequest struct {
	UpdateType FetchType
	Status     DownloadStatus
	period     uint64
}

type FetchDataResponse struct {
	period     uint64
	UpdateType FetchType
}

// ____________________

type ProofStatus int

const (
	ProofDefault ProofStatus = iota
	ProofPending
	ProofSuccess
	ProofFailed
)

type ProofType int

const (
	Deposit ProofType = iota + 1
	Redeem
	Verify
)

type TxType = int

const (
	DepositTx TxType = iota + 1
	RedeemTx
)

type ChainType = int

const (
	Bitcoin ChainType = iota + 1
	Ethereum
)

type Transaction struct {
	TxHash   string
	DestHash string
	Height   int64

	BtcTxId string

	Amount  int64
	EthAddr string
	Utxo    []Utxo

	Inputs  []Utxo
	Outputs []TxOut

	TxType    TxType
	ChainType ChainType
}

type BitcoinTx struct {
	EthAddr string
	Amount  int64 // btc

	EthTxHash string
	Height    int64
	BlockHash string
	TxId      string
	Utxos     []Utxo
	TxType    ProofType
}

type EthereumTx struct {
	Height    int64
	BlockHash string
	Inputs    []Utxo
	Outputs   []TxOut

	Amount  int64
	BtcTxId string
	Vout    int

	TxHash string
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

type Proof struct {
	TxId      string      `json:"txId"`
	ProofType ProofType   `json:"type"`
	Proof     string      `json:"proof"`
	Status    ProofStatus `json:"status"`
}

type ProofRequest struct {
	// redeem
	Inputs  []Utxo  `json:"inputs"`
	Outputs []TxOut `json:"outputs"`
	BtcTxId string  `json:"btcTxId"`

	// deposit
	Utxos   []Utxo
	Amount  int64  `json:"amount"`
	EthAddr string `json:"ethAddr"`

	// other
	Height    int64     `json:"height"`
	BlockHash string    `json:"blockHash"`
	TxHash    string    `json:"txId"`
	ProofType ProofType `json:"type"`
	Proof     string    `json:"proof"`
	Msg       string    `json:"msg"`
}

func (req *ProofRequest) Type() string {
	//TODO implement me
	panic("implement me")
}

func (req *ProofRequest) String() string {
	if req.ProofType == Deposit {
		return fmt.Sprintf("txType:%v,txid: %v, utxos:%v, amount:%v, ethAddr:%v", req.ProofType, req.TxHash, req.Utxos, req.Amount, req.EthAddr)
	} else if req.ProofType == Redeem {
		return fmt.Sprintf("txType:%v,txid:%v, utxos:%v, outputs: %v", req.ProofType, req.TxHash, formatUtxo(req.Inputs), formatOut(req.Outputs))
	}
	return ""
}

type ProofResponse struct {
	// redeem
	Inputs  []Utxo  `json:"inputs"`
	Outputs []TxOut `json:"outputs"`
	BtcTxId string  `json:"btcTxId"`

	// deposit
	Utxos   []Utxo
	Amount  int64  `json:"amount"`
	EthAddr string `json:"ethAddr"`

	// other
	Height    int64       `json:"height"`
	BlockHash string      `json:"blockHash"`
	TxId      string      `json:"txId"`
	ProofType ProofType   `json:"type"`
	Proof     string      `json:"proof"`
	Msg       string      `json:"msg"`
	Status    ProofStatus `json:"status"`
}

func (resp *ProofResponse) Type() string {
	//TODO implement me
	panic("implement me")
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
