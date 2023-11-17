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

// BtcTxLibTxOut is an auto generated low-level Go binding around an user-defined struct.
type BtcTxLibTxOut struct {
	Value    uint64
	PkScript []byte
}

// UTXOManagerLocator is an auto generated low-level Go binding around an user-defined struct.
type UTXOManagerLocator struct {
	ShuntType uint8
	Index     uint64
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
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"}],\"name\":\"ChangeAlreadyUpdated\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"}],\"name\":\"ChangeNotExisting\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"totalAmount\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"amount\",\"type\":\"uint64\"}],\"name\":\"TotalAmountInsufficient\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"}],\"name\":\"UTXOExisting\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"amount\",\"type\":\"uint64\"}],\"name\":\"ChangeUTXOUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"amount\",\"type\":\"uint64\"}],\"name\":\"UTXOAdded\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"NORMAL\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PETTY\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"TINY\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"amount\",\"type\":\"uint64\"}],\"name\":\"addChange\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"amount\",\"type\":\"uint64\"}],\"name\":\"addUTXO\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"multiSigScript\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"amount\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"isChangeConfirmed\",\"type\":\"bool\"}],\"internalType\":\"structUTXOManager.UTXO[]\",\"name\":\"inputs\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"value\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"pkScript\",\"type\":\"bytes\"}],\"internalType\":\"structBtcTxLib.TxOut[]\",\"name\":\"outputs\",\"type\":\"tuple[]\"}],\"name\":\"calcSigHashs\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"dolphinAvailableKeys\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"targetAmount\",\"type\":\"uint64\"}],\"name\":\"findUTXOs\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"amount\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"isChangeConfirmed\",\"type\":\"bool\"}],\"internalType\":\"structUTXOManager.UTXO[]\",\"name\":\"\",\"type\":\"tuple[]\"},{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"shuntType\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"index\",\"type\":\"uint64\"}],\"internalType\":\"structUTXOManager.Locator[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"fragmentAvailableKeys\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"getRoleMember\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleMemberCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleMembers\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"shuntType\",\"type\":\"uint8\"}],\"name\":\"mapToArray\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"targetAmount\",\"type\":\"uint64\"}],\"name\":\"spendUTXOs\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"amount\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"isChangeConfirmed\",\"type\":\"bool\"}],\"internalType\":\"structUTXOManager.UTXO[]\",\"name\":\"\",\"type\":\"tuple[]\"},{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalAvailableAmount\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"}],\"name\":\"updateChange\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"}],\"name\":\"utxoOf\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"amount\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"isChangeConfirmed\",\"type\":\"bool\"}],\"internalType\":\"structUTXOManager.UTXO\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"utxos\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"amount\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"isChangeConfirmed\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"walrusAvailableKeys\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"whaleAvailableKeys\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
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

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Utxo *UtxoCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Utxo.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Utxo *UtxoSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Utxo.Contract.DEFAULTADMINROLE(&_Utxo.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Utxo *UtxoCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Utxo.Contract.DEFAULTADMINROLE(&_Utxo.CallOpts)
}

// NORMAL is a free data retrieval call binding the contract method 0xae24595c.
//
// Solidity: function NORMAL() view returns(uint64)
func (_Utxo *UtxoCaller) NORMAL(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Utxo.contract.Call(opts, &out, "NORMAL")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// NORMAL is a free data retrieval call binding the contract method 0xae24595c.
//
// Solidity: function NORMAL() view returns(uint64)
func (_Utxo *UtxoSession) NORMAL() (uint64, error) {
	return _Utxo.Contract.NORMAL(&_Utxo.CallOpts)
}

// NORMAL is a free data retrieval call binding the contract method 0xae24595c.
//
// Solidity: function NORMAL() view returns(uint64)
func (_Utxo *UtxoCallerSession) NORMAL() (uint64, error) {
	return _Utxo.Contract.NORMAL(&_Utxo.CallOpts)
}

// PETTY is a free data retrieval call binding the contract method 0x877e28e6.
//
// Solidity: function PETTY() view returns(uint64)
func (_Utxo *UtxoCaller) PETTY(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Utxo.contract.Call(opts, &out, "PETTY")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// PETTY is a free data retrieval call binding the contract method 0x877e28e6.
//
// Solidity: function PETTY() view returns(uint64)
func (_Utxo *UtxoSession) PETTY() (uint64, error) {
	return _Utxo.Contract.PETTY(&_Utxo.CallOpts)
}

// PETTY is a free data retrieval call binding the contract method 0x877e28e6.
//
// Solidity: function PETTY() view returns(uint64)
func (_Utxo *UtxoCallerSession) PETTY() (uint64, error) {
	return _Utxo.Contract.PETTY(&_Utxo.CallOpts)
}

// TINY is a free data retrieval call binding the contract method 0x29b5a596.
//
// Solidity: function TINY() view returns(uint64)
func (_Utxo *UtxoCaller) TINY(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Utxo.contract.Call(opts, &out, "TINY")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// TINY is a free data retrieval call binding the contract method 0x29b5a596.
//
// Solidity: function TINY() view returns(uint64)
func (_Utxo *UtxoSession) TINY() (uint64, error) {
	return _Utxo.Contract.TINY(&_Utxo.CallOpts)
}

// TINY is a free data retrieval call binding the contract method 0x29b5a596.
//
// Solidity: function TINY() view returns(uint64)
func (_Utxo *UtxoCallerSession) TINY() (uint64, error) {
	return _Utxo.Contract.TINY(&_Utxo.CallOpts)
}

// CalcSigHashs is a free data retrieval call binding the contract method 0xdde51cb6.
//
// Solidity: function calcSigHashs(bytes multiSigScript, (bytes32,uint32,uint64,bool)[] inputs, (uint64,bytes)[] outputs) pure returns(bytes32[], bytes, bytes32)
func (_Utxo *UtxoCaller) CalcSigHashs(opts *bind.CallOpts, multiSigScript []byte, inputs []UTXOManagerUTXO, outputs []BtcTxLibTxOut) ([][32]byte, []byte, [32]byte, error) {
	var out []interface{}
	err := _Utxo.contract.Call(opts, &out, "calcSigHashs", multiSigScript, inputs, outputs)

	if err != nil {
		return *new([][32]byte), *new([]byte), *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)
	out1 := *abi.ConvertType(out[1], new([]byte)).(*[]byte)
	out2 := *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return out0, out1, out2, err

}

// CalcSigHashs is a free data retrieval call binding the contract method 0xdde51cb6.
//
// Solidity: function calcSigHashs(bytes multiSigScript, (bytes32,uint32,uint64,bool)[] inputs, (uint64,bytes)[] outputs) pure returns(bytes32[], bytes, bytes32)
func (_Utxo *UtxoSession) CalcSigHashs(multiSigScript []byte, inputs []UTXOManagerUTXO, outputs []BtcTxLibTxOut) ([][32]byte, []byte, [32]byte, error) {
	return _Utxo.Contract.CalcSigHashs(&_Utxo.CallOpts, multiSigScript, inputs, outputs)
}

// CalcSigHashs is a free data retrieval call binding the contract method 0xdde51cb6.
//
// Solidity: function calcSigHashs(bytes multiSigScript, (bytes32,uint32,uint64,bool)[] inputs, (uint64,bytes)[] outputs) pure returns(bytes32[], bytes, bytes32)
func (_Utxo *UtxoCallerSession) CalcSigHashs(multiSigScript []byte, inputs []UTXOManagerUTXO, outputs []BtcTxLibTxOut) ([][32]byte, []byte, [32]byte, error) {
	return _Utxo.Contract.CalcSigHashs(&_Utxo.CallOpts, multiSigScript, inputs, outputs)
}

// DolphinAvailableKeys is a free data retrieval call binding the contract method 0x884d952a.
//
// Solidity: function dolphinAvailableKeys(uint256 ) view returns(bytes32)
func (_Utxo *UtxoCaller) DolphinAvailableKeys(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _Utxo.contract.Call(opts, &out, "dolphinAvailableKeys", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DolphinAvailableKeys is a free data retrieval call binding the contract method 0x884d952a.
//
// Solidity: function dolphinAvailableKeys(uint256 ) view returns(bytes32)
func (_Utxo *UtxoSession) DolphinAvailableKeys(arg0 *big.Int) ([32]byte, error) {
	return _Utxo.Contract.DolphinAvailableKeys(&_Utxo.CallOpts, arg0)
}

// DolphinAvailableKeys is a free data retrieval call binding the contract method 0x884d952a.
//
// Solidity: function dolphinAvailableKeys(uint256 ) view returns(bytes32)
func (_Utxo *UtxoCallerSession) DolphinAvailableKeys(arg0 *big.Int) ([32]byte, error) {
	return _Utxo.Contract.DolphinAvailableKeys(&_Utxo.CallOpts, arg0)
}

// FindUTXOs is a free data retrieval call binding the contract method 0x9186c857.
//
// Solidity: function findUTXOs(uint64 targetAmount) view returns((bytes32,uint32,uint64,bool)[], uint64, (uint8,uint64)[])
func (_Utxo *UtxoCaller) FindUTXOs(opts *bind.CallOpts, targetAmount uint64) ([]UTXOManagerUTXO, uint64, []UTXOManagerLocator, error) {
	var out []interface{}
	err := _Utxo.contract.Call(opts, &out, "findUTXOs", targetAmount)

	if err != nil {
		return *new([]UTXOManagerUTXO), *new(uint64), *new([]UTXOManagerLocator), err
	}

	out0 := *abi.ConvertType(out[0], new([]UTXOManagerUTXO)).(*[]UTXOManagerUTXO)
	out1 := *abi.ConvertType(out[1], new(uint64)).(*uint64)
	out2 := *abi.ConvertType(out[2], new([]UTXOManagerLocator)).(*[]UTXOManagerLocator)

	return out0, out1, out2, err

}

// FindUTXOs is a free data retrieval call binding the contract method 0x9186c857.
//
// Solidity: function findUTXOs(uint64 targetAmount) view returns((bytes32,uint32,uint64,bool)[], uint64, (uint8,uint64)[])
func (_Utxo *UtxoSession) FindUTXOs(targetAmount uint64) ([]UTXOManagerUTXO, uint64, []UTXOManagerLocator, error) {
	return _Utxo.Contract.FindUTXOs(&_Utxo.CallOpts, targetAmount)
}

// FindUTXOs is a free data retrieval call binding the contract method 0x9186c857.
//
// Solidity: function findUTXOs(uint64 targetAmount) view returns((bytes32,uint32,uint64,bool)[], uint64, (uint8,uint64)[])
func (_Utxo *UtxoCallerSession) FindUTXOs(targetAmount uint64) ([]UTXOManagerUTXO, uint64, []UTXOManagerLocator, error) {
	return _Utxo.Contract.FindUTXOs(&_Utxo.CallOpts, targetAmount)
}

// FragmentAvailableKeys is a free data retrieval call binding the contract method 0x292adb04.
//
// Solidity: function fragmentAvailableKeys(uint256 ) view returns(bytes32)
func (_Utxo *UtxoCaller) FragmentAvailableKeys(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _Utxo.contract.Call(opts, &out, "fragmentAvailableKeys", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// FragmentAvailableKeys is a free data retrieval call binding the contract method 0x292adb04.
//
// Solidity: function fragmentAvailableKeys(uint256 ) view returns(bytes32)
func (_Utxo *UtxoSession) FragmentAvailableKeys(arg0 *big.Int) ([32]byte, error) {
	return _Utxo.Contract.FragmentAvailableKeys(&_Utxo.CallOpts, arg0)
}

// FragmentAvailableKeys is a free data retrieval call binding the contract method 0x292adb04.
//
// Solidity: function fragmentAvailableKeys(uint256 ) view returns(bytes32)
func (_Utxo *UtxoCallerSession) FragmentAvailableKeys(arg0 *big.Int) ([32]byte, error) {
	return _Utxo.Contract.FragmentAvailableKeys(&_Utxo.CallOpts, arg0)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Utxo *UtxoCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _Utxo.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Utxo *UtxoSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Utxo.Contract.GetRoleAdmin(&_Utxo.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Utxo *UtxoCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Utxo.Contract.GetRoleAdmin(&_Utxo.CallOpts, role)
}

// GetRoleMember is a free data retrieval call binding the contract method 0x9010d07c.
//
// Solidity: function getRoleMember(bytes32 role, uint256 index) view returns(address)
func (_Utxo *UtxoCaller) GetRoleMember(opts *bind.CallOpts, role [32]byte, index *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Utxo.contract.Call(opts, &out, "getRoleMember", role, index)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetRoleMember is a free data retrieval call binding the contract method 0x9010d07c.
//
// Solidity: function getRoleMember(bytes32 role, uint256 index) view returns(address)
func (_Utxo *UtxoSession) GetRoleMember(role [32]byte, index *big.Int) (common.Address, error) {
	return _Utxo.Contract.GetRoleMember(&_Utxo.CallOpts, role, index)
}

// GetRoleMember is a free data retrieval call binding the contract method 0x9010d07c.
//
// Solidity: function getRoleMember(bytes32 role, uint256 index) view returns(address)
func (_Utxo *UtxoCallerSession) GetRoleMember(role [32]byte, index *big.Int) (common.Address, error) {
	return _Utxo.Contract.GetRoleMember(&_Utxo.CallOpts, role, index)
}

// GetRoleMemberCount is a free data retrieval call binding the contract method 0xca15c873.
//
// Solidity: function getRoleMemberCount(bytes32 role) view returns(uint256)
func (_Utxo *UtxoCaller) GetRoleMemberCount(opts *bind.CallOpts, role [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _Utxo.contract.Call(opts, &out, "getRoleMemberCount", role)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRoleMemberCount is a free data retrieval call binding the contract method 0xca15c873.
//
// Solidity: function getRoleMemberCount(bytes32 role) view returns(uint256)
func (_Utxo *UtxoSession) GetRoleMemberCount(role [32]byte) (*big.Int, error) {
	return _Utxo.Contract.GetRoleMemberCount(&_Utxo.CallOpts, role)
}

// GetRoleMemberCount is a free data retrieval call binding the contract method 0xca15c873.
//
// Solidity: function getRoleMemberCount(bytes32 role) view returns(uint256)
func (_Utxo *UtxoCallerSession) GetRoleMemberCount(role [32]byte) (*big.Int, error) {
	return _Utxo.Contract.GetRoleMemberCount(&_Utxo.CallOpts, role)
}

// GetRoleMembers is a free data retrieval call binding the contract method 0xa3246ad3.
//
// Solidity: function getRoleMembers(bytes32 role) view returns(address[])
func (_Utxo *UtxoCaller) GetRoleMembers(opts *bind.CallOpts, role [32]byte) ([]common.Address, error) {
	var out []interface{}
	err := _Utxo.contract.Call(opts, &out, "getRoleMembers", role)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetRoleMembers is a free data retrieval call binding the contract method 0xa3246ad3.
//
// Solidity: function getRoleMembers(bytes32 role) view returns(address[])
func (_Utxo *UtxoSession) GetRoleMembers(role [32]byte) ([]common.Address, error) {
	return _Utxo.Contract.GetRoleMembers(&_Utxo.CallOpts, role)
}

// GetRoleMembers is a free data retrieval call binding the contract method 0xa3246ad3.
//
// Solidity: function getRoleMembers(bytes32 role) view returns(address[])
func (_Utxo *UtxoCallerSession) GetRoleMembers(role [32]byte) ([]common.Address, error) {
	return _Utxo.Contract.GetRoleMembers(&_Utxo.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Utxo *UtxoCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _Utxo.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Utxo *UtxoSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Utxo.Contract.HasRole(&_Utxo.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Utxo *UtxoCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Utxo.Contract.HasRole(&_Utxo.CallOpts, role, account)
}

// MapToArray is a free data retrieval call binding the contract method 0xc259561a.
//
// Solidity: function mapToArray(uint8 shuntType) view returns(bytes32[])
func (_Utxo *UtxoCaller) MapToArray(opts *bind.CallOpts, shuntType uint8) ([][32]byte, error) {
	var out []interface{}
	err := _Utxo.contract.Call(opts, &out, "mapToArray", shuntType)

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

// MapToArray is a free data retrieval call binding the contract method 0xc259561a.
//
// Solidity: function mapToArray(uint8 shuntType) view returns(bytes32[])
func (_Utxo *UtxoSession) MapToArray(shuntType uint8) ([][32]byte, error) {
	return _Utxo.Contract.MapToArray(&_Utxo.CallOpts, shuntType)
}

// MapToArray is a free data retrieval call binding the contract method 0xc259561a.
//
// Solidity: function mapToArray(uint8 shuntType) view returns(bytes32[])
func (_Utxo *UtxoCallerSession) MapToArray(shuntType uint8) ([][32]byte, error) {
	return _Utxo.Contract.MapToArray(&_Utxo.CallOpts, shuntType)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Utxo *UtxoCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _Utxo.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Utxo *UtxoSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Utxo.Contract.SupportsInterface(&_Utxo.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Utxo *UtxoCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Utxo.Contract.SupportsInterface(&_Utxo.CallOpts, interfaceId)
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

// WalrusAvailableKeys is a free data retrieval call binding the contract method 0x6d5c47b1.
//
// Solidity: function walrusAvailableKeys(uint256 ) view returns(bytes32)
func (_Utxo *UtxoCaller) WalrusAvailableKeys(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _Utxo.contract.Call(opts, &out, "walrusAvailableKeys", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// WalrusAvailableKeys is a free data retrieval call binding the contract method 0x6d5c47b1.
//
// Solidity: function walrusAvailableKeys(uint256 ) view returns(bytes32)
func (_Utxo *UtxoSession) WalrusAvailableKeys(arg0 *big.Int) ([32]byte, error) {
	return _Utxo.Contract.WalrusAvailableKeys(&_Utxo.CallOpts, arg0)
}

// WalrusAvailableKeys is a free data retrieval call binding the contract method 0x6d5c47b1.
//
// Solidity: function walrusAvailableKeys(uint256 ) view returns(bytes32)
func (_Utxo *UtxoCallerSession) WalrusAvailableKeys(arg0 *big.Int) ([32]byte, error) {
	return _Utxo.Contract.WalrusAvailableKeys(&_Utxo.CallOpts, arg0)
}

// WhaleAvailableKeys is a free data retrieval call binding the contract method 0x724071f5.
//
// Solidity: function whaleAvailableKeys(uint256 ) view returns(bytes32)
func (_Utxo *UtxoCaller) WhaleAvailableKeys(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _Utxo.contract.Call(opts, &out, "whaleAvailableKeys", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// WhaleAvailableKeys is a free data retrieval call binding the contract method 0x724071f5.
//
// Solidity: function whaleAvailableKeys(uint256 ) view returns(bytes32)
func (_Utxo *UtxoSession) WhaleAvailableKeys(arg0 *big.Int) ([32]byte, error) {
	return _Utxo.Contract.WhaleAvailableKeys(&_Utxo.CallOpts, arg0)
}

// WhaleAvailableKeys is a free data retrieval call binding the contract method 0x724071f5.
//
// Solidity: function whaleAvailableKeys(uint256 ) view returns(bytes32)
func (_Utxo *UtxoCallerSession) WhaleAvailableKeys(arg0 *big.Int) ([32]byte, error) {
	return _Utxo.Contract.WhaleAvailableKeys(&_Utxo.CallOpts, arg0)
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

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Utxo *UtxoTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Utxo.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Utxo *UtxoSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Utxo.Contract.GrantRole(&_Utxo.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Utxo *UtxoTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Utxo.Contract.GrantRole(&_Utxo.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_Utxo *UtxoTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _Utxo.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_Utxo *UtxoSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _Utxo.Contract.RenounceRole(&_Utxo.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_Utxo *UtxoTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _Utxo.Contract.RenounceRole(&_Utxo.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Utxo *UtxoTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Utxo.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Utxo *UtxoSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Utxo.Contract.RevokeRole(&_Utxo.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Utxo *UtxoTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Utxo.Contract.RevokeRole(&_Utxo.TransactOpts, role, account)
}

// SpendUTXOs is a paid mutator transaction binding the contract method 0x2d8240b3.
//
// Solidity: function spendUTXOs(uint64 targetAmount) returns((bytes32,uint32,uint64,bool)[], uint64)
func (_Utxo *UtxoTransactor) SpendUTXOs(opts *bind.TransactOpts, targetAmount uint64) (*types.Transaction, error) {
	return _Utxo.contract.Transact(opts, "spendUTXOs", targetAmount)
}

// SpendUTXOs is a paid mutator transaction binding the contract method 0x2d8240b3.
//
// Solidity: function spendUTXOs(uint64 targetAmount) returns((bytes32,uint32,uint64,bool)[], uint64)
func (_Utxo *UtxoSession) SpendUTXOs(targetAmount uint64) (*types.Transaction, error) {
	return _Utxo.Contract.SpendUTXOs(&_Utxo.TransactOpts, targetAmount)
}

// SpendUTXOs is a paid mutator transaction binding the contract method 0x2d8240b3.
//
// Solidity: function spendUTXOs(uint64 targetAmount) returns((bytes32,uint32,uint64,bool)[], uint64)
func (_Utxo *UtxoTransactorSession) SpendUTXOs(targetAmount uint64) (*types.Transaction, error) {
	return _Utxo.Contract.SpendUTXOs(&_Utxo.TransactOpts, targetAmount)
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

// UtxoChangeUTXOUpdatedIterator is returned from FilterChangeUTXOUpdated and is used to iterate over the raw logs and unpacked data for ChangeUTXOUpdated events raised by the Utxo contract.
type UtxoChangeUTXOUpdatedIterator struct {
	Event *UtxoChangeUTXOUpdated // Event containing the contract specifics and raw log

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
func (it *UtxoChangeUTXOUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UtxoChangeUTXOUpdated)
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
		it.Event = new(UtxoChangeUTXOUpdated)
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
func (it *UtxoChangeUTXOUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UtxoChangeUTXOUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UtxoChangeUTXOUpdated represents a ChangeUTXOUpdated event raised by the Utxo contract.
type UtxoChangeUTXOUpdated struct {
	Txid   [32]byte
	Index  uint32
	Amount uint64
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterChangeUTXOUpdated is a free log retrieval operation binding the contract event 0x5d6cc5c33e60ae274f09159956ae8fd20271c63c95b9004445271437b335ed6e.
//
// Solidity: event ChangeUTXOUpdated(bytes32 indexed txid, uint32 indexed index, uint64 amount)
func (_Utxo *UtxoFilterer) FilterChangeUTXOUpdated(opts *bind.FilterOpts, txid [][32]byte, index []uint32) (*UtxoChangeUTXOUpdatedIterator, error) {

	var txidRule []interface{}
	for _, txidItem := range txid {
		txidRule = append(txidRule, txidItem)
	}
	var indexRule []interface{}
	for _, indexItem := range index {
		indexRule = append(indexRule, indexItem)
	}

	logs, sub, err := _Utxo.contract.FilterLogs(opts, "ChangeUTXOUpdated", txidRule, indexRule)
	if err != nil {
		return nil, err
	}
	return &UtxoChangeUTXOUpdatedIterator{contract: _Utxo.contract, event: "ChangeUTXOUpdated", logs: logs, sub: sub}, nil
}

// WatchChangeUTXOUpdated is a free log subscription operation binding the contract event 0x5d6cc5c33e60ae274f09159956ae8fd20271c63c95b9004445271437b335ed6e.
//
// Solidity: event ChangeUTXOUpdated(bytes32 indexed txid, uint32 indexed index, uint64 amount)
func (_Utxo *UtxoFilterer) WatchChangeUTXOUpdated(opts *bind.WatchOpts, sink chan<- *UtxoChangeUTXOUpdated, txid [][32]byte, index []uint32) (event.Subscription, error) {

	var txidRule []interface{}
	for _, txidItem := range txid {
		txidRule = append(txidRule, txidItem)
	}
	var indexRule []interface{}
	for _, indexItem := range index {
		indexRule = append(indexRule, indexItem)
	}

	logs, sub, err := _Utxo.contract.WatchLogs(opts, "ChangeUTXOUpdated", txidRule, indexRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UtxoChangeUTXOUpdated)
				if err := _Utxo.contract.UnpackLog(event, "ChangeUTXOUpdated", log); err != nil {
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

// ParseChangeUTXOUpdated is a log parse operation binding the contract event 0x5d6cc5c33e60ae274f09159956ae8fd20271c63c95b9004445271437b335ed6e.
//
// Solidity: event ChangeUTXOUpdated(bytes32 indexed txid, uint32 indexed index, uint64 amount)
func (_Utxo *UtxoFilterer) ParseChangeUTXOUpdated(log types.Log) (*UtxoChangeUTXOUpdated, error) {
	event := new(UtxoChangeUTXOUpdated)
	if err := _Utxo.contract.UnpackLog(event, "ChangeUTXOUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UtxoRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the Utxo contract.
type UtxoRoleAdminChangedIterator struct {
	Event *UtxoRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *UtxoRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UtxoRoleAdminChanged)
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
		it.Event = new(UtxoRoleAdminChanged)
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
func (it *UtxoRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UtxoRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UtxoRoleAdminChanged represents a RoleAdminChanged event raised by the Utxo contract.
type UtxoRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Utxo *UtxoFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*UtxoRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _Utxo.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &UtxoRoleAdminChangedIterator{contract: _Utxo.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Utxo *UtxoFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *UtxoRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _Utxo.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UtxoRoleAdminChanged)
				if err := _Utxo.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Utxo *UtxoFilterer) ParseRoleAdminChanged(log types.Log) (*UtxoRoleAdminChanged, error) {
	event := new(UtxoRoleAdminChanged)
	if err := _Utxo.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UtxoRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the Utxo contract.
type UtxoRoleGrantedIterator struct {
	Event *UtxoRoleGranted // Event containing the contract specifics and raw log

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
func (it *UtxoRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UtxoRoleGranted)
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
		it.Event = new(UtxoRoleGranted)
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
func (it *UtxoRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UtxoRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UtxoRoleGranted represents a RoleGranted event raised by the Utxo contract.
type UtxoRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Utxo *UtxoFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*UtxoRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Utxo.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &UtxoRoleGrantedIterator{contract: _Utxo.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Utxo *UtxoFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *UtxoRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Utxo.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UtxoRoleGranted)
				if err := _Utxo.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Utxo *UtxoFilterer) ParseRoleGranted(log types.Log) (*UtxoRoleGranted, error) {
	event := new(UtxoRoleGranted)
	if err := _Utxo.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UtxoRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the Utxo contract.
type UtxoRoleRevokedIterator struct {
	Event *UtxoRoleRevoked // Event containing the contract specifics and raw log

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
func (it *UtxoRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UtxoRoleRevoked)
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
		it.Event = new(UtxoRoleRevoked)
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
func (it *UtxoRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UtxoRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UtxoRoleRevoked represents a RoleRevoked event raised by the Utxo contract.
type UtxoRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Utxo *UtxoFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*UtxoRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Utxo.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &UtxoRoleRevokedIterator{contract: _Utxo.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Utxo *UtxoFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *UtxoRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Utxo.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UtxoRoleRevoked)
				if err := _Utxo.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Utxo *UtxoFilterer) ParseRoleRevoked(log types.Log) (*UtxoRoleRevoked, error) {
	event := new(UtxoRoleRevoked)
	if err := _Utxo.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UtxoUTXOAddedIterator is returned from FilterUTXOAdded and is used to iterate over the raw logs and unpacked data for UTXOAdded events raised by the Utxo contract.
type UtxoUTXOAddedIterator struct {
	Event *UtxoUTXOAdded // Event containing the contract specifics and raw log

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
func (it *UtxoUTXOAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UtxoUTXOAdded)
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
		it.Event = new(UtxoUTXOAdded)
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
func (it *UtxoUTXOAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UtxoUTXOAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UtxoUTXOAdded represents a UTXOAdded event raised by the Utxo contract.
type UtxoUTXOAdded struct {
	Txid   [32]byte
	Index  uint32
	Amount uint64
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterUTXOAdded is a free log retrieval operation binding the contract event 0xd063609fea0cb9b8a1b53a4fbf0e659c270b3bc99eab08dcc7f4433b4937e074.
//
// Solidity: event UTXOAdded(bytes32 indexed txid, uint32 indexed index, uint64 amount)
func (_Utxo *UtxoFilterer) FilterUTXOAdded(opts *bind.FilterOpts, txid [][32]byte, index []uint32) (*UtxoUTXOAddedIterator, error) {

	var txidRule []interface{}
	for _, txidItem := range txid {
		txidRule = append(txidRule, txidItem)
	}
	var indexRule []interface{}
	for _, indexItem := range index {
		indexRule = append(indexRule, indexItem)
	}

	logs, sub, err := _Utxo.contract.FilterLogs(opts, "UTXOAdded", txidRule, indexRule)
	if err != nil {
		return nil, err
	}
	return &UtxoUTXOAddedIterator{contract: _Utxo.contract, event: "UTXOAdded", logs: logs, sub: sub}, nil
}

// WatchUTXOAdded is a free log subscription operation binding the contract event 0xd063609fea0cb9b8a1b53a4fbf0e659c270b3bc99eab08dcc7f4433b4937e074.
//
// Solidity: event UTXOAdded(bytes32 indexed txid, uint32 indexed index, uint64 amount)
func (_Utxo *UtxoFilterer) WatchUTXOAdded(opts *bind.WatchOpts, sink chan<- *UtxoUTXOAdded, txid [][32]byte, index []uint32) (event.Subscription, error) {

	var txidRule []interface{}
	for _, txidItem := range txid {
		txidRule = append(txidRule, txidItem)
	}
	var indexRule []interface{}
	for _, indexItem := range index {
		indexRule = append(indexRule, indexItem)
	}

	logs, sub, err := _Utxo.contract.WatchLogs(opts, "UTXOAdded", txidRule, indexRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UtxoUTXOAdded)
				if err := _Utxo.contract.UnpackLog(event, "UTXOAdded", log); err != nil {
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

// ParseUTXOAdded is a log parse operation binding the contract event 0xd063609fea0cb9b8a1b53a4fbf0e659c270b3bc99eab08dcc7f4433b4937e074.
//
// Solidity: event UTXOAdded(bytes32 indexed txid, uint32 indexed index, uint64 amount)
func (_Utxo *UtxoFilterer) ParseUTXOAdded(log types.Log) (*UtxoUTXOAdded, error) {
	event := new(UtxoUTXOAdded)
	if err := _Utxo.contract.UnpackLog(event, "UTXOAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
