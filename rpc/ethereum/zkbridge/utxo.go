// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package zkbridge

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// BtcTxLibMsgTx is an auto generated low-level Go binding around an user-defined struct.
type BtcTxLibMsgTx struct {
	Version  [4]byte
	TxIns    []BtcTxLibTxIn
	TxOuts   []BtcTxLibTxOut
	LockTime [4]byte
}

// BtcTxLibTxIn is an auto generated low-level Go binding around an user-defined struct.
type BtcTxLibTxIn struct {
	Hash     [32]byte
	Index    uint32
	Script   []byte
	Witness  [][]byte
	Sequence [4]byte
}

// BtcTxLibTxOut is an auto generated low-level Go binding around an user-defined struct.
type BtcTxLibTxOut struct {
	Value    uint64
	PkScript []byte
}

// UTXOManagerUTXO is an auto generated low-level Go binding around an user-defined struct.
type UTXOManagerUTXO struct {
	Txid              [32]byte
	Index             uint32
	Amount            uint64
	IsChangeConfirmed bool
}

// UtxoMetaData contains all meta data concerning the Utxo contract.
var UtxoMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"}],\"name\":\"ChangeIsNotExisted\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"}],\"name\":\"ChangeIsUpdated\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"totalAmount\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"amount\",\"type\":\"uint64\"}],\"name\":\"InsufficientTotalAmount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"amount\",\"type\":\"uint64\"}],\"name\":\"InvalidAmount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"}],\"name\":\"UTXOIsExisted\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"limit\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"utxoCount\",\"type\":\"uint8\"}],\"name\":\"UtxoExceedLimit\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"amount\",\"type\":\"uint64\"}],\"name\":\"UTXOAdd\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"UTXOAddBatchSucess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"amount\",\"type\":\"uint64\"}],\"name\":\"UTXOUpdate\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_new\",\"type\":\"address\"}],\"name\":\"addOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isOperator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_new\",\"type\":\"address\"}],\"name\":\"removeOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalAvailableAmount\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"utxos\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"amount\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"isChangeConfirmed\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"}],\"name\":\"utxoOf\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"amount\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"isChangeConfirmed\",\"type\":\"bool\"}],\"internalType\":\"structUTXOManager.UTXO\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getUtxosAvailableKeys\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"amount\",\"type\":\"uint64\"}],\"name\":\"addUTXO\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"amount\",\"type\":\"uint64\"}],\"name\":\"addChange\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"}],\"name\":\"updateChange\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"targetAmount\",\"type\":\"uint64\"}],\"name\":\"spentUTXOs\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"amount\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"isChangeConfirmed\",\"type\":\"bool\"}],\"internalType\":\"structUTXOManager.UTXO[]\",\"name\":\"foundUTXOs\",\"type\":\"tuple[]\"},{\"internalType\":\"uint64\",\"name\":\"totalAmount\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"targetAmount\",\"type\":\"uint64\"}],\"name\":\"findUTXOs\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"amount\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"isChangeConfirmed\",\"type\":\"bool\"}],\"internalType\":\"structUTXOManager.UTXO[]\",\"name\":\"findedUTXOs\",\"type\":\"tuple[]\"},{\"internalType\":\"uint64\",\"name\":\"totalAmount\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"amount\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"isChangeConfirmed\",\"type\":\"bool\"}],\"internalType\":\"structUTXOManager.UTXO[]\",\"name\":\"inputs\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"value\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"pkScript\",\"type\":\"bytes\"}],\"internalType\":\"structBtcTxLib.TxOut[]\",\"name\":\"outputs\",\"type\":\"tuple[]\"}],\"name\":\"createUnsignedTx\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes4\",\"name\":\"version\",\"type\":\"bytes4\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"script\",\"type\":\"bytes\"},{\"internalType\":\"bytes[]\",\"name\":\"witness\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes4\",\"name\":\"sequence\",\"type\":\"bytes4\"}],\"internalType\":\"structBtcTxLib.TxIn[]\",\"name\":\"txIns\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"value\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"pkScript\",\"type\":\"bytes\"}],\"internalType\":\"structBtcTxLib.TxOut[]\",\"name\":\"txOuts\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes4\",\"name\":\"lockTime\",\"type\":\"bytes4\"}],\"internalType\":\"structBtcTxLib.MsgTx\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes4\",\"name\":\"version\",\"type\":\"bytes4\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"script\",\"type\":\"bytes\"},{\"internalType\":\"bytes[]\",\"name\":\"witness\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes4\",\"name\":\"sequence\",\"type\":\"bytes4\"}],\"internalType\":\"structBtcTxLib.TxIn[]\",\"name\":\"txIns\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"value\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"pkScript\",\"type\":\"bytes\"}],\"internalType\":\"structBtcTxLib.TxOut[]\",\"name\":\"txOuts\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes4\",\"name\":\"lockTime\",\"type\":\"bytes4\"}],\"internalType\":\"structBtcTxLib.MsgTx\",\"name\":\"msgTx\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"multiSigScript\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"amount\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"isChangeConfirmed\",\"type\":\"bool\"}],\"internalType\":\"structUTXOManager.UTXO[]\",\"name\":\"inputs\",\"type\":\"tuple[]\"}],\"name\":\"calcSigHashs\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
}

// UtxoABI is the input ABI used to generate the binding from.
// Deprecated: Use UtxoMetaData.ABI instead.
var UtxoABI = UtxoMetaData.ABI

// Utxo is an auto generated Go binding around an Ethereum contract.
type Utxo struct {
	UtxoCaller     // Read-only binding to the contract
	UtxoTransactor // Write-only binding to the contract
	UtxoFilterer   // Log filterer for contract events
}

// UtxoCaller is an auto generated read-only Go binding around an Ethereum contract.
type UtxoCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UtxoTransactor is an auto generated write-only Go binding around an Ethereum contract.
type UtxoTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UtxoFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type UtxoFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UtxoSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type UtxoSession struct {
	Contract     *Utxo             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// UtxoCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type UtxoCallerSession struct {
	Contract *UtxoCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// UtxoTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type UtxoTransactorSession struct {
	Contract     *UtxoTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// UtxoRaw is an auto generated low-level Go binding around an Ethereum contract.
type UtxoRaw struct {
	Contract *Utxo // Generic contract binding to access the raw methods on
}

// UtxoCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type UtxoCallerRaw struct {
	Contract *UtxoCaller // Generic read-only contract binding to access the raw methods on
}

// UtxoTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type UtxoTransactorRaw struct {
	Contract *UtxoTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUtxo creates a new instance of Utxo, bound to a specific deployed contract.
func NewUtxo(address common.Address, backend bind.ContractBackend) (*Utxo, error) {
	contract, err := bindUtxo(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Utxo{UtxoCaller: UtxoCaller{contract: contract}, UtxoTransactor: UtxoTransactor{contract: contract}, UtxoFilterer: UtxoFilterer{contract: contract}}, nil
}

// NewUtxoCaller creates a new read-only instance of Utxo, bound to a specific deployed contract.
func NewUtxoCaller(address common.Address, caller bind.ContractCaller) (*UtxoCaller, error) {
	contract, err := bindUtxo(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UtxoCaller{contract: contract}, nil
}

// NewUtxoTransactor creates a new write-only instance of Utxo, bound to a specific deployed contract.
func NewUtxoTransactor(address common.Address, transactor bind.ContractTransactor) (*UtxoTransactor, error) {
	contract, err := bindUtxo(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UtxoTransactor{contract: contract}, nil
}

// NewUtxoFilterer creates a new log filterer instance of Utxo, bound to a specific deployed contract.
func NewUtxoFilterer(address common.Address, filterer bind.ContractFilterer) (*UtxoFilterer, error) {
	contract, err := bindUtxo(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UtxoFilterer{contract: contract}, nil
}

// bindUtxo binds a generic wrapper to an already deployed contract.
func bindUtxo(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := UtxoMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Utxo *UtxoRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Utxo.Contract.UtxoCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Utxo *UtxoRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Utxo.Contract.UtxoTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Utxo *UtxoRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Utxo.Contract.UtxoTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Utxo *UtxoCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Utxo.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Utxo *UtxoTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Utxo.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Utxo *UtxoTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Utxo.Contract.contract.Transact(opts, method, params...)
}

// CalcSigHashs is a free data retrieval call binding the contract method 0xc021634c.
//
// Solidity: function calcSigHashs((bytes4,(bytes32,uint32,bytes,bytes[],bytes4)[],(uint64,bytes)[],bytes4) msgTx, bytes multiSigScript, (bytes32,uint32,uint64,bool)[] inputs) pure returns(bytes32[])
func (_Utxo *UtxoCaller) CalcSigHashs(opts *bind.CallOpts, msgTx BtcTxLibMsgTx, multiSigScript []byte, inputs []UTXOManagerUTXO) ([][32]byte, error) {
	var out []interface{}
	err := _Utxo.contract.Call(opts, &out, "calcSigHashs", msgTx, multiSigScript, inputs)

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

// CalcSigHashs is a free data retrieval call binding the contract method 0xc021634c.
//
// Solidity: function calcSigHashs((bytes4,(bytes32,uint32,bytes,bytes[],bytes4)[],(uint64,bytes)[],bytes4) msgTx, bytes multiSigScript, (bytes32,uint32,uint64,bool)[] inputs) pure returns(bytes32[])
func (_Utxo *UtxoSession) CalcSigHashs(msgTx BtcTxLibMsgTx, multiSigScript []byte, inputs []UTXOManagerUTXO) ([][32]byte, error) {
	return _Utxo.Contract.CalcSigHashs(&_Utxo.CallOpts, msgTx, multiSigScript, inputs)
}

// CalcSigHashs is a free data retrieval call binding the contract method 0xc021634c.
//
// Solidity: function calcSigHashs((bytes4,(bytes32,uint32,bytes,bytes[],bytes4)[],(uint64,bytes)[],bytes4) msgTx, bytes multiSigScript, (bytes32,uint32,uint64,bool)[] inputs) pure returns(bytes32[])
func (_Utxo *UtxoCallerSession) CalcSigHashs(msgTx BtcTxLibMsgTx, multiSigScript []byte, inputs []UTXOManagerUTXO) ([][32]byte, error) {
	return _Utxo.Contract.CalcSigHashs(&_Utxo.CallOpts, msgTx, multiSigScript, inputs)
}

// CreateUnsignedTx is a free data retrieval call binding the contract method 0x3516b44b.
//
// Solidity: function createUnsignedTx((bytes32,uint32,uint64,bool)[] inputs, (uint64,bytes)[] outputs) pure returns((bytes4,(bytes32,uint32,bytes,bytes[],bytes4)[],(uint64,bytes)[],bytes4))
func (_Utxo *UtxoCaller) CreateUnsignedTx(opts *bind.CallOpts, inputs []UTXOManagerUTXO, outputs []BtcTxLibTxOut) (BtcTxLibMsgTx, error) {
	var out []interface{}
	err := _Utxo.contract.Call(opts, &out, "createUnsignedTx", inputs, outputs)

	if err != nil {
		return *new(BtcTxLibMsgTx), err
	}

	out0 := *abi.ConvertType(out[0], new(BtcTxLibMsgTx)).(*BtcTxLibMsgTx)

	return out0, err

}

// CreateUnsignedTx is a free data retrieval call binding the contract method 0x3516b44b.
//
// Solidity: function createUnsignedTx((bytes32,uint32,uint64,bool)[] inputs, (uint64,bytes)[] outputs) pure returns((bytes4,(bytes32,uint32,bytes,bytes[],bytes4)[],(uint64,bytes)[],bytes4))
func (_Utxo *UtxoSession) CreateUnsignedTx(inputs []UTXOManagerUTXO, outputs []BtcTxLibTxOut) (BtcTxLibMsgTx, error) {
	return _Utxo.Contract.CreateUnsignedTx(&_Utxo.CallOpts, inputs, outputs)
}

// CreateUnsignedTx is a free data retrieval call binding the contract method 0x3516b44b.
//
// Solidity: function createUnsignedTx((bytes32,uint32,uint64,bool)[] inputs, (uint64,bytes)[] outputs) pure returns((bytes4,(bytes32,uint32,bytes,bytes[],bytes4)[],(uint64,bytes)[],bytes4))
func (_Utxo *UtxoCallerSession) CreateUnsignedTx(inputs []UTXOManagerUTXO, outputs []BtcTxLibTxOut) (BtcTxLibMsgTx, error) {
	return _Utxo.Contract.CreateUnsignedTx(&_Utxo.CallOpts, inputs, outputs)
}

// FindUTXOs is a free data retrieval call binding the contract method 0x9186c857.
//
// Solidity: function findUTXOs(uint64 targetAmount) view returns((bytes32,uint32,uint64,bool)[] findedUTXOs, uint64 totalAmount)
func (_Utxo *UtxoCaller) FindUTXOs(opts *bind.CallOpts, targetAmount uint64) (struct {
	FindedUTXOs []UTXOManagerUTXO
	TotalAmount uint64
}, error) {
	var out []interface{}
	err := _Utxo.contract.Call(opts, &out, "findUTXOs", targetAmount)

	outstruct := new(struct {
		FindedUTXOs []UTXOManagerUTXO
		TotalAmount uint64
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.FindedUTXOs = *abi.ConvertType(out[0], new([]UTXOManagerUTXO)).(*[]UTXOManagerUTXO)
	outstruct.TotalAmount = *abi.ConvertType(out[1], new(uint64)).(*uint64)

	return *outstruct, err

}

// FindUTXOs is a free data retrieval call binding the contract method 0x9186c857.
//
// Solidity: function findUTXOs(uint64 targetAmount) view returns((bytes32,uint32,uint64,bool)[] findedUTXOs, uint64 totalAmount)
func (_Utxo *UtxoSession) FindUTXOs(targetAmount uint64) (struct {
	FindedUTXOs []UTXOManagerUTXO
	TotalAmount uint64
}, error) {
	return _Utxo.Contract.FindUTXOs(&_Utxo.CallOpts, targetAmount)
}

// FindUTXOs is a free data retrieval call binding the contract method 0x9186c857.
//
// Solidity: function findUTXOs(uint64 targetAmount) view returns((bytes32,uint32,uint64,bool)[] findedUTXOs, uint64 totalAmount)
func (_Utxo *UtxoCallerSession) FindUTXOs(targetAmount uint64) (struct {
	FindedUTXOs []UTXOManagerUTXO
	TotalAmount uint64
}, error) {
	return _Utxo.Contract.FindUTXOs(&_Utxo.CallOpts, targetAmount)
}

// GetUtxosAvailableKeys is a free data retrieval call binding the contract method 0x67a066aa.
//
// Solidity: function getUtxosAvailableKeys() view returns(bytes32[])
func (_Utxo *UtxoCaller) GetUtxosAvailableKeys(opts *bind.CallOpts) ([][32]byte, error) {
	var out []interface{}
	err := _Utxo.contract.Call(opts, &out, "getUtxosAvailableKeys")

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

// GetUtxosAvailableKeys is a free data retrieval call binding the contract method 0x67a066aa.
//
// Solidity: function getUtxosAvailableKeys() view returns(bytes32[])
func (_Utxo *UtxoSession) GetUtxosAvailableKeys() ([][32]byte, error) {
	return _Utxo.Contract.GetUtxosAvailableKeys(&_Utxo.CallOpts)
}

// GetUtxosAvailableKeys is a free data retrieval call binding the contract method 0x67a066aa.
//
// Solidity: function getUtxosAvailableKeys() view returns(bytes32[])
func (_Utxo *UtxoCallerSession) GetUtxosAvailableKeys() ([][32]byte, error) {
	return _Utxo.Contract.GetUtxosAvailableKeys(&_Utxo.CallOpts)
}

// IsOperator is a free data retrieval call binding the contract method 0x6d70f7ae.
//
// Solidity: function isOperator(address addr) view returns(bool)
func (_Utxo *UtxoCaller) IsOperator(opts *bind.CallOpts, addr common.Address) (bool, error) {
	var out []interface{}
	err := _Utxo.contract.Call(opts, &out, "isOperator", addr)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOperator is a free data retrieval call binding the contract method 0x6d70f7ae.
//
// Solidity: function isOperator(address addr) view returns(bool)
func (_Utxo *UtxoSession) IsOperator(addr common.Address) (bool, error) {
	return _Utxo.Contract.IsOperator(&_Utxo.CallOpts, addr)
}

// IsOperator is a free data retrieval call binding the contract method 0x6d70f7ae.
//
// Solidity: function isOperator(address addr) view returns(bool)
func (_Utxo *UtxoCallerSession) IsOperator(addr common.Address) (bool, error) {
	return _Utxo.Contract.IsOperator(&_Utxo.CallOpts, addr)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Utxo *UtxoCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Utxo.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Utxo *UtxoSession) Owner() (common.Address, error) {
	return _Utxo.Contract.Owner(&_Utxo.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Utxo *UtxoCallerSession) Owner() (common.Address, error) {
	return _Utxo.Contract.Owner(&_Utxo.CallOpts)
}

// TotalAvailableAmount is a free data retrieval call binding the contract method 0x7a1406fb.
//
// Solidity: function totalAvailableAmount() view returns(uint64)
func (_Utxo *UtxoCaller) TotalAvailableAmount(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Utxo.contract.Call(opts, &out, "totalAvailableAmount")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// TotalAvailableAmount is a free data retrieval call binding the contract method 0x7a1406fb.
//
// Solidity: function totalAvailableAmount() view returns(uint64)
func (_Utxo *UtxoSession) TotalAvailableAmount() (uint64, error) {
	return _Utxo.Contract.TotalAvailableAmount(&_Utxo.CallOpts)
}

// TotalAvailableAmount is a free data retrieval call binding the contract method 0x7a1406fb.
//
// Solidity: function totalAvailableAmount() view returns(uint64)
func (_Utxo *UtxoCallerSession) TotalAvailableAmount() (uint64, error) {
	return _Utxo.Contract.TotalAvailableAmount(&_Utxo.CallOpts)
}

// UtxoOf is a free data retrieval call binding the contract method 0xd45c66e4.
//
// Solidity: function utxoOf(bytes32 txid) view returns((bytes32,uint32,uint64,bool))
func (_Utxo *UtxoCaller) UtxoOf(opts *bind.CallOpts, txid [32]byte) (UTXOManagerUTXO, error) {
	var out []interface{}
	err := _Utxo.contract.Call(opts, &out, "utxoOf", txid)

	if err != nil {
		return *new(UTXOManagerUTXO), err
	}

	out0 := *abi.ConvertType(out[0], new(UTXOManagerUTXO)).(*UTXOManagerUTXO)

	return out0, err

}

// UtxoOf is a free data retrieval call binding the contract method 0xd45c66e4.
//
// Solidity: function utxoOf(bytes32 txid) view returns((bytes32,uint32,uint64,bool))
func (_Utxo *UtxoSession) UtxoOf(txid [32]byte) (UTXOManagerUTXO, error) {
	return _Utxo.Contract.UtxoOf(&_Utxo.CallOpts, txid)
}

// UtxoOf is a free data retrieval call binding the contract method 0xd45c66e4.
//
// Solidity: function utxoOf(bytes32 txid) view returns((bytes32,uint32,uint64,bool))
func (_Utxo *UtxoCallerSession) UtxoOf(txid [32]byte) (UTXOManagerUTXO, error) {
	return _Utxo.Contract.UtxoOf(&_Utxo.CallOpts, txid)
}

// Utxos is a free data retrieval call binding the contract method 0x3cf18ab9.
//
// Solidity: function utxos(bytes32 ) view returns(bytes32 txid, uint32 index, uint64 amount, bool isChangeConfirmed)
func (_Utxo *UtxoCaller) Utxos(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Txid              [32]byte
	Index             uint32
	Amount            uint64
	IsChangeConfirmed bool
}, error) {
	var out []interface{}
	err := _Utxo.contract.Call(opts, &out, "utxos", arg0)

	outstruct := new(struct {
		Txid              [32]byte
		Index             uint32
		Amount            uint64
		IsChangeConfirmed bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Txid = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.Index = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.Amount = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	outstruct.IsChangeConfirmed = *abi.ConvertType(out[3], new(bool)).(*bool)

	return *outstruct, err

}

// Utxos is a free data retrieval call binding the contract method 0x3cf18ab9.
//
// Solidity: function utxos(bytes32 ) view returns(bytes32 txid, uint32 index, uint64 amount, bool isChangeConfirmed)
func (_Utxo *UtxoSession) Utxos(arg0 [32]byte) (struct {
	Txid              [32]byte
	Index             uint32
	Amount            uint64
	IsChangeConfirmed bool
}, error) {
	return _Utxo.Contract.Utxos(&_Utxo.CallOpts, arg0)
}

// Utxos is a free data retrieval call binding the contract method 0x3cf18ab9.
//
// Solidity: function utxos(bytes32 ) view returns(bytes32 txid, uint32 index, uint64 amount, bool isChangeConfirmed)
func (_Utxo *UtxoCallerSession) Utxos(arg0 [32]byte) (struct {
	Txid              [32]byte
	Index             uint32
	Amount            uint64
	IsChangeConfirmed bool
}, error) {
	return _Utxo.Contract.Utxos(&_Utxo.CallOpts, arg0)
}

// AddChange is a paid mutator transaction binding the contract method 0xb1d53798.
//
// Solidity: function addChange(bytes32 txid, uint32 index, uint64 amount) returns()
func (_Utxo *UtxoTransactor) AddChange(opts *bind.TransactOpts, txid [32]byte, index uint32, amount uint64) (*types.Transaction, error) {
	return _Utxo.contract.Transact(opts, "addChange", txid, index, amount)
}

// AddChange is a paid mutator transaction binding the contract method 0xb1d53798.
//
// Solidity: function addChange(bytes32 txid, uint32 index, uint64 amount) returns()
func (_Utxo *UtxoSession) AddChange(txid [32]byte, index uint32, amount uint64) (*types.Transaction, error) {
	return _Utxo.Contract.AddChange(&_Utxo.TransactOpts, txid, index, amount)
}

// AddChange is a paid mutator transaction binding the contract method 0xb1d53798.
//
// Solidity: function addChange(bytes32 txid, uint32 index, uint64 amount) returns()
func (_Utxo *UtxoTransactorSession) AddChange(txid [32]byte, index uint32, amount uint64) (*types.Transaction, error) {
	return _Utxo.Contract.AddChange(&_Utxo.TransactOpts, txid, index, amount)
}

// AddOperator is a paid mutator transaction binding the contract method 0x9870d7fe.
//
// Solidity: function addOperator(address _new) returns()
func (_Utxo *UtxoTransactor) AddOperator(opts *bind.TransactOpts, _new common.Address) (*types.Transaction, error) {
	return _Utxo.contract.Transact(opts, "addOperator", _new)
}

// AddOperator is a paid mutator transaction binding the contract method 0x9870d7fe.
//
// Solidity: function addOperator(address _new) returns()
func (_Utxo *UtxoSession) AddOperator(_new common.Address) (*types.Transaction, error) {
	return _Utxo.Contract.AddOperator(&_Utxo.TransactOpts, _new)
}

// AddOperator is a paid mutator transaction binding the contract method 0x9870d7fe.
//
// Solidity: function addOperator(address _new) returns()
func (_Utxo *UtxoTransactorSession) AddOperator(_new common.Address) (*types.Transaction, error) {
	return _Utxo.Contract.AddOperator(&_Utxo.TransactOpts, _new)
}

// AddUTXO is a paid mutator transaction binding the contract method 0xe8128197.
//
// Solidity: function addUTXO(bytes32 txid, uint32 index, uint64 amount) returns()
func (_Utxo *UtxoTransactor) AddUTXO(opts *bind.TransactOpts, txid [32]byte, index uint32, amount uint64) (*types.Transaction, error) {
	return _Utxo.contract.Transact(opts, "addUTXO", txid, index, amount)
}

// AddUTXO is a paid mutator transaction binding the contract method 0xe8128197.
//
// Solidity: function addUTXO(bytes32 txid, uint32 index, uint64 amount) returns()
func (_Utxo *UtxoSession) AddUTXO(txid [32]byte, index uint32, amount uint64) (*types.Transaction, error) {
	return _Utxo.Contract.AddUTXO(&_Utxo.TransactOpts, txid, index, amount)
}

// AddUTXO is a paid mutator transaction binding the contract method 0xe8128197.
//
// Solidity: function addUTXO(bytes32 txid, uint32 index, uint64 amount) returns()
func (_Utxo *UtxoTransactorSession) AddUTXO(txid [32]byte, index uint32, amount uint64) (*types.Transaction, error) {
	return _Utxo.Contract.AddUTXO(&_Utxo.TransactOpts, txid, index, amount)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0xac8a584a.
//
// Solidity: function removeOperator(address _new) returns()
func (_Utxo *UtxoTransactor) RemoveOperator(opts *bind.TransactOpts, _new common.Address) (*types.Transaction, error) {
	return _Utxo.contract.Transact(opts, "removeOperator", _new)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0xac8a584a.
//
// Solidity: function removeOperator(address _new) returns()
func (_Utxo *UtxoSession) RemoveOperator(_new common.Address) (*types.Transaction, error) {
	return _Utxo.Contract.RemoveOperator(&_Utxo.TransactOpts, _new)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0xac8a584a.
//
// Solidity: function removeOperator(address _new) returns()
func (_Utxo *UtxoTransactorSession) RemoveOperator(_new common.Address) (*types.Transaction, error) {
	return _Utxo.Contract.RemoveOperator(&_Utxo.TransactOpts, _new)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Utxo *UtxoTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Utxo.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Utxo *UtxoSession) RenounceOwnership() (*types.Transaction, error) {
	return _Utxo.Contract.RenounceOwnership(&_Utxo.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Utxo *UtxoTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Utxo.Contract.RenounceOwnership(&_Utxo.TransactOpts)
}

// SpentUTXOs is a paid mutator transaction binding the contract method 0xcd33c014.
//
// Solidity: function spentUTXOs(uint64 targetAmount) returns((bytes32,uint32,uint64,bool)[] foundUTXOs, uint64 totalAmount)
func (_Utxo *UtxoTransactor) SpentUTXOs(opts *bind.TransactOpts, targetAmount uint64) (*types.Transaction, error) {
	return _Utxo.contract.Transact(opts, "spentUTXOs", targetAmount)
}

// SpentUTXOs is a paid mutator transaction binding the contract method 0xcd33c014.
//
// Solidity: function spentUTXOs(uint64 targetAmount) returns((bytes32,uint32,uint64,bool)[] foundUTXOs, uint64 totalAmount)
func (_Utxo *UtxoSession) SpentUTXOs(targetAmount uint64) (*types.Transaction, error) {
	return _Utxo.Contract.SpentUTXOs(&_Utxo.TransactOpts, targetAmount)
}

// SpentUTXOs is a paid mutator transaction binding the contract method 0xcd33c014.
//
// Solidity: function spentUTXOs(uint64 targetAmount) returns((bytes32,uint32,uint64,bool)[] foundUTXOs, uint64 totalAmount)
func (_Utxo *UtxoTransactorSession) SpentUTXOs(targetAmount uint64) (*types.Transaction, error) {
	return _Utxo.Contract.SpentUTXOs(&_Utxo.TransactOpts, targetAmount)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Utxo *UtxoTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Utxo.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Utxo *UtxoSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Utxo.Contract.TransferOwnership(&_Utxo.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Utxo *UtxoTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Utxo.Contract.TransferOwnership(&_Utxo.TransactOpts, newOwner)
}

// UpdateChange is a paid mutator transaction binding the contract method 0x9d23d0fa.
//
// Solidity: function updateChange(bytes32 txid) returns()
func (_Utxo *UtxoTransactor) UpdateChange(opts *bind.TransactOpts, txid [32]byte) (*types.Transaction, error) {
	return _Utxo.contract.Transact(opts, "updateChange", txid)
}

// UpdateChange is a paid mutator transaction binding the contract method 0x9d23d0fa.
//
// Solidity: function updateChange(bytes32 txid) returns()
func (_Utxo *UtxoSession) UpdateChange(txid [32]byte) (*types.Transaction, error) {
	return _Utxo.Contract.UpdateChange(&_Utxo.TransactOpts, txid)
}

// UpdateChange is a paid mutator transaction binding the contract method 0x9d23d0fa.
//
// Solidity: function updateChange(bytes32 txid) returns()
func (_Utxo *UtxoTransactorSession) UpdateChange(txid [32]byte) (*types.Transaction, error) {
	return _Utxo.Contract.UpdateChange(&_Utxo.TransactOpts, txid)
}

// UtxoOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Utxo contract.
type UtxoOwnershipTransferredIterator struct {
	Event *UtxoOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *UtxoOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UtxoOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(UtxoOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *UtxoOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UtxoOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UtxoOwnershipTransferred represents a OwnershipTransferred event raised by the Utxo contract.
type UtxoOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Utxo *UtxoFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*UtxoOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Utxo.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &UtxoOwnershipTransferredIterator{contract: _Utxo.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Utxo *UtxoFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *UtxoOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Utxo.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UtxoOwnershipTransferred)
				if err := _Utxo.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Utxo *UtxoFilterer) ParseOwnershipTransferred(log types.Log) (*UtxoOwnershipTransferred, error) {
	event := new(UtxoOwnershipTransferred)
	if err := _Utxo.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UtxoUTXOAddIterator is returned from FilterUTXOAdd and is used to iterate over the raw logs and unpacked data for UTXOAdd events raised by the Utxo contract.
type UtxoUTXOAddIterator struct {
	Event *UtxoUTXOAdd // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *UtxoUTXOAddIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UtxoUTXOAdd)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(UtxoUTXOAdd)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *UtxoUTXOAddIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UtxoUTXOAddIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UtxoUTXOAdd represents a UTXOAdd event raised by the Utxo contract.
type UtxoUTXOAdd struct {
	Txid   [32]byte
	Index  uint32
	Amount uint64
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterUTXOAdd is a free log retrieval operation binding the contract event 0xbfb6a0aa850eff6109c854ffb48321dcf37f02d6c7a44c46987a5ddf3419fc07.
//
// Solidity: event UTXOAdd(bytes32 indexed txid, uint32 indexed index, uint64 amount)
func (_Utxo *UtxoFilterer) FilterUTXOAdd(opts *bind.FilterOpts, txid [][32]byte, index []uint32) (*UtxoUTXOAddIterator, error) {

	var txidRule []interface{}
	for _, txidItem := range txid {
		txidRule = append(txidRule, txidItem)
	}
	var indexRule []interface{}
	for _, indexItem := range index {
		indexRule = append(indexRule, indexItem)
	}

	logs, sub, err := _Utxo.contract.FilterLogs(opts, "UTXOAdd", txidRule, indexRule)
	if err != nil {
		return nil, err
	}
	return &UtxoUTXOAddIterator{contract: _Utxo.contract, event: "UTXOAdd", logs: logs, sub: sub}, nil
}

// WatchUTXOAdd is a free log subscription operation binding the contract event 0xbfb6a0aa850eff6109c854ffb48321dcf37f02d6c7a44c46987a5ddf3419fc07.
//
// Solidity: event UTXOAdd(bytes32 indexed txid, uint32 indexed index, uint64 amount)
func (_Utxo *UtxoFilterer) WatchUTXOAdd(opts *bind.WatchOpts, sink chan<- *UtxoUTXOAdd, txid [][32]byte, index []uint32) (event.Subscription, error) {

	var txidRule []interface{}
	for _, txidItem := range txid {
		txidRule = append(txidRule, txidItem)
	}
	var indexRule []interface{}
	for _, indexItem := range index {
		indexRule = append(indexRule, indexItem)
	}

	logs, sub, err := _Utxo.contract.WatchLogs(opts, "UTXOAdd", txidRule, indexRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UtxoUTXOAdd)
				if err := _Utxo.contract.UnpackLog(event, "UTXOAdd", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUTXOAdd is a log parse operation binding the contract event 0xbfb6a0aa850eff6109c854ffb48321dcf37f02d6c7a44c46987a5ddf3419fc07.
//
// Solidity: event UTXOAdd(bytes32 indexed txid, uint32 indexed index, uint64 amount)
func (_Utxo *UtxoFilterer) ParseUTXOAdd(log types.Log) (*UtxoUTXOAdd, error) {
	event := new(UtxoUTXOAdd)
	if err := _Utxo.contract.UnpackLog(event, "UTXOAdd", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UtxoUTXOAddBatchSucessIterator is returned from FilterUTXOAddBatchSucess and is used to iterate over the raw logs and unpacked data for UTXOAddBatchSucess events raised by the Utxo contract.
type UtxoUTXOAddBatchSucessIterator struct {
	Event *UtxoUTXOAddBatchSucess // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *UtxoUTXOAddBatchSucessIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UtxoUTXOAddBatchSucess)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(UtxoUTXOAddBatchSucess)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *UtxoUTXOAddBatchSucessIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UtxoUTXOAddBatchSucessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UtxoUTXOAddBatchSucess represents a UTXOAddBatchSucess event raised by the Utxo contract.
type UtxoUTXOAddBatchSucess struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterUTXOAddBatchSucess is a free log retrieval operation binding the contract event 0xfb9d208236e4647154ad06e593e7792387aca07600c97d6c8072759b3a3e316a.
//
// Solidity: event UTXOAddBatchSucess()
func (_Utxo *UtxoFilterer) FilterUTXOAddBatchSucess(opts *bind.FilterOpts) (*UtxoUTXOAddBatchSucessIterator, error) {

	logs, sub, err := _Utxo.contract.FilterLogs(opts, "UTXOAddBatchSucess")
	if err != nil {
		return nil, err
	}
	return &UtxoUTXOAddBatchSucessIterator{contract: _Utxo.contract, event: "UTXOAddBatchSucess", logs: logs, sub: sub}, nil
}

// WatchUTXOAddBatchSucess is a free log subscription operation binding the contract event 0xfb9d208236e4647154ad06e593e7792387aca07600c97d6c8072759b3a3e316a.
//
// Solidity: event UTXOAddBatchSucess()
func (_Utxo *UtxoFilterer) WatchUTXOAddBatchSucess(opts *bind.WatchOpts, sink chan<- *UtxoUTXOAddBatchSucess) (event.Subscription, error) {

	logs, sub, err := _Utxo.contract.WatchLogs(opts, "UTXOAddBatchSucess")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UtxoUTXOAddBatchSucess)
				if err := _Utxo.contract.UnpackLog(event, "UTXOAddBatchSucess", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUTXOAddBatchSucess is a log parse operation binding the contract event 0xfb9d208236e4647154ad06e593e7792387aca07600c97d6c8072759b3a3e316a.
//
// Solidity: event UTXOAddBatchSucess()
func (_Utxo *UtxoFilterer) ParseUTXOAddBatchSucess(log types.Log) (*UtxoUTXOAddBatchSucess, error) {
	event := new(UtxoUTXOAddBatchSucess)
	if err := _Utxo.contract.UnpackLog(event, "UTXOAddBatchSucess", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UtxoUTXOUpdateIterator is returned from FilterUTXOUpdate and is used to iterate over the raw logs and unpacked data for UTXOUpdate events raised by the Utxo contract.
type UtxoUTXOUpdateIterator struct {
	Event *UtxoUTXOUpdate // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *UtxoUTXOUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UtxoUTXOUpdate)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(UtxoUTXOUpdate)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *UtxoUTXOUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UtxoUTXOUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UtxoUTXOUpdate represents a UTXOUpdate event raised by the Utxo contract.
type UtxoUTXOUpdate struct {
	Txid   [32]byte
	Index  uint32
	Amount uint64
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterUTXOUpdate is a free log retrieval operation binding the contract event 0xa62639770d75b3f144cd4b62df93eb44cf61fdaac90d461d1f68389a2d461ebe.
//
// Solidity: event UTXOUpdate(bytes32 indexed txid, uint32 indexed index, uint64 amount)
func (_Utxo *UtxoFilterer) FilterUTXOUpdate(opts *bind.FilterOpts, txid [][32]byte, index []uint32) (*UtxoUTXOUpdateIterator, error) {

	var txidRule []interface{}
	for _, txidItem := range txid {
		txidRule = append(txidRule, txidItem)
	}
	var indexRule []interface{}
	for _, indexItem := range index {
		indexRule = append(indexRule, indexItem)
	}

	logs, sub, err := _Utxo.contract.FilterLogs(opts, "UTXOUpdate", txidRule, indexRule)
	if err != nil {
		return nil, err
	}
	return &UtxoUTXOUpdateIterator{contract: _Utxo.contract, event: "UTXOUpdate", logs: logs, sub: sub}, nil
}

// WatchUTXOUpdate is a free log subscription operation binding the contract event 0xa62639770d75b3f144cd4b62df93eb44cf61fdaac90d461d1f68389a2d461ebe.
//
// Solidity: event UTXOUpdate(bytes32 indexed txid, uint32 indexed index, uint64 amount)
func (_Utxo *UtxoFilterer) WatchUTXOUpdate(opts *bind.WatchOpts, sink chan<- *UtxoUTXOUpdate, txid [][32]byte, index []uint32) (event.Subscription, error) {

	var txidRule []interface{}
	for _, txidItem := range txid {
		txidRule = append(txidRule, txidItem)
	}
	var indexRule []interface{}
	for _, indexItem := range index {
		indexRule = append(indexRule, indexItem)
	}

	logs, sub, err := _Utxo.contract.WatchLogs(opts, "UTXOUpdate", txidRule, indexRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UtxoUTXOUpdate)
				if err := _Utxo.contract.UnpackLog(event, "UTXOUpdate", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUTXOUpdate is a log parse operation binding the contract event 0xa62639770d75b3f144cd4b62df93eb44cf61fdaac90d461d1f68389a2d461ebe.
//
// Solidity: event UTXOUpdate(bytes32 indexed txid, uint32 indexed index, uint64 amount)
func (_Utxo *UtxoFilterer) ParseUTXOUpdate(log types.Log) (*UtxoUTXOUpdate, error) {
	event := new(UtxoUTXOUpdate)
	if err := _Utxo.contract.UnpackLog(event, "UTXOUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
