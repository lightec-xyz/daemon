// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package zkbridge_verify

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

// ZkbridgeVerifyMetaData contains all meta data concerning the ZkbridgeVerify contract.
var ZkbridgeVerifyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_new\",\"type\":\"address\"}],\"name\":\"addOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isOperator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_new\",\"type\":\"address\"}],\"name\":\"removeOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"signer_public_key\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"verifierAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"updateVerifierAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"txRaw\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"receiptRaw\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"proofData\",\"type\":\"bytes\"}],\"name\":\"signBtcTx\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"\",\"type\":\"bytes[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// ZkbridgeVerifyABI is the input ABI used to generate the binding from.
// Deprecated: Use ZkbridgeVerifyMetaData.ABI instead.
var ZkbridgeVerifyABI = ZkbridgeVerifyMetaData.ABI

// ZkbridgeVerify is an auto generated Go binding around an Ethereum contract.
type ZkbridgeVerify struct {
	ZkbridgeVerifyCaller     // Read-only binding to the contract
	ZkbridgeVerifyTransactor // Write-only binding to the contract
	ZkbridgeVerifyFilterer   // Log filterer for contract events
}

// ZkbridgeVerifyCaller is an auto generated read-only Go binding around an Ethereum contract.
type ZkbridgeVerifyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZkbridgeVerifyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ZkbridgeVerifyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZkbridgeVerifyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ZkbridgeVerifyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZkbridgeVerifySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ZkbridgeVerifySession struct {
	Contract     *ZkbridgeVerify   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ZkbridgeVerifyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ZkbridgeVerifyCallerSession struct {
	Contract *ZkbridgeVerifyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// ZkbridgeVerifyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ZkbridgeVerifyTransactorSession struct {
	Contract     *ZkbridgeVerifyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// ZkbridgeVerifyRaw is an auto generated low-level Go binding around an Ethereum contract.
type ZkbridgeVerifyRaw struct {
	Contract *ZkbridgeVerify // Generic contract binding to access the raw methods on
}

// ZkbridgeVerifyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ZkbridgeVerifyCallerRaw struct {
	Contract *ZkbridgeVerifyCaller // Generic read-only contract binding to access the raw methods on
}

// ZkbridgeVerifyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ZkbridgeVerifyTransactorRaw struct {
	Contract *ZkbridgeVerifyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewZkbridgeVerify creates a new instance of ZkbridgeVerify, bound to a specific deployed contract.
func NewZkbridgeVerify(address common.Address, backend bind.ContractBackend) (*ZkbridgeVerify, error) {
	contract, err := bindZkbridgeVerify(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ZkbridgeVerify{ZkbridgeVerifyCaller: ZkbridgeVerifyCaller{contract: contract}, ZkbridgeVerifyTransactor: ZkbridgeVerifyTransactor{contract: contract}, ZkbridgeVerifyFilterer: ZkbridgeVerifyFilterer{contract: contract}}, nil
}

// NewZkbridgeVerifyCaller creates a new read-only instance of ZkbridgeVerify, bound to a specific deployed contract.
func NewZkbridgeVerifyCaller(address common.Address, caller bind.ContractCaller) (*ZkbridgeVerifyCaller, error) {
	contract, err := bindZkbridgeVerify(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ZkbridgeVerifyCaller{contract: contract}, nil
}

// NewZkbridgeVerifyTransactor creates a new write-only instance of ZkbridgeVerify, bound to a specific deployed contract.
func NewZkbridgeVerifyTransactor(address common.Address, transactor bind.ContractTransactor) (*ZkbridgeVerifyTransactor, error) {
	contract, err := bindZkbridgeVerify(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ZkbridgeVerifyTransactor{contract: contract}, nil
}

// NewZkbridgeVerifyFilterer creates a new log filterer instance of ZkbridgeVerify, bound to a specific deployed contract.
func NewZkbridgeVerifyFilterer(address common.Address, filterer bind.ContractFilterer) (*ZkbridgeVerifyFilterer, error) {
	contract, err := bindZkbridgeVerify(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ZkbridgeVerifyFilterer{contract: contract}, nil
}

// bindZkbridgeVerify binds a generic wrapper to an already deployed contract.
func bindZkbridgeVerify(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ZkbridgeVerifyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ZkbridgeVerify *ZkbridgeVerifyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ZkbridgeVerify.Contract.ZkbridgeVerifyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ZkbridgeVerify *ZkbridgeVerifyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZkbridgeVerify.Contract.ZkbridgeVerifyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ZkbridgeVerify *ZkbridgeVerifyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ZkbridgeVerify.Contract.ZkbridgeVerifyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ZkbridgeVerify *ZkbridgeVerifyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ZkbridgeVerify.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ZkbridgeVerify *ZkbridgeVerifyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZkbridgeVerify.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ZkbridgeVerify *ZkbridgeVerifyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ZkbridgeVerify.Contract.contract.Transact(opts, method, params...)
}

// IsOperator is a free data retrieval call binding the contract method 0x6d70f7ae.
//
// Solidity: function isOperator(address addr) view returns(bool)
func (_ZkbridgeVerify *ZkbridgeVerifyCaller) IsOperator(opts *bind.CallOpts, addr common.Address) (bool, error) {
	var out []interface{}
	err := _ZkbridgeVerify.contract.Call(opts, &out, "isOperator", addr)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOperator is a free data retrieval call binding the contract method 0x6d70f7ae.
//
// Solidity: function isOperator(address addr) view returns(bool)
func (_ZkbridgeVerify *ZkbridgeVerifySession) IsOperator(addr common.Address) (bool, error) {
	return _ZkbridgeVerify.Contract.IsOperator(&_ZkbridgeVerify.CallOpts, addr)
}

// IsOperator is a free data retrieval call binding the contract method 0x6d70f7ae.
//
// Solidity: function isOperator(address addr) view returns(bool)
func (_ZkbridgeVerify *ZkbridgeVerifyCallerSession) IsOperator(addr common.Address) (bool, error) {
	return _ZkbridgeVerify.Contract.IsOperator(&_ZkbridgeVerify.CallOpts, addr)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ZkbridgeVerify *ZkbridgeVerifyCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ZkbridgeVerify.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ZkbridgeVerify *ZkbridgeVerifySession) Owner() (common.Address, error) {
	return _ZkbridgeVerify.Contract.Owner(&_ZkbridgeVerify.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ZkbridgeVerify *ZkbridgeVerifyCallerSession) Owner() (common.Address, error) {
	return _ZkbridgeVerify.Contract.Owner(&_ZkbridgeVerify.CallOpts)
}

// SignBtcTx is a free data retrieval call binding the contract method 0xe43d4471.
//
// Solidity: function signBtcTx(bytes txRaw, bytes receiptRaw, bytes proofData) view returns(bytes[])
func (_ZkbridgeVerify *ZkbridgeVerifyCaller) SignBtcTx(opts *bind.CallOpts, txRaw []byte, receiptRaw []byte, proofData []byte) ([][]byte, error) {
	var out []interface{}
	err := _ZkbridgeVerify.contract.Call(opts, &out, "signBtcTx", txRaw, receiptRaw, proofData)

	if err != nil {
		return *new([][]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][]byte)).(*[][]byte)

	return out0, err

}

// SignBtcTx is a free data retrieval call binding the contract method 0xe43d4471.
//
// Solidity: function signBtcTx(bytes txRaw, bytes receiptRaw, bytes proofData) view returns(bytes[])
func (_ZkbridgeVerify *ZkbridgeVerifySession) SignBtcTx(txRaw []byte, receiptRaw []byte, proofData []byte) ([][]byte, error) {
	return _ZkbridgeVerify.Contract.SignBtcTx(&_ZkbridgeVerify.CallOpts, txRaw, receiptRaw, proofData)
}

// SignBtcTx is a free data retrieval call binding the contract method 0xe43d4471.
//
// Solidity: function signBtcTx(bytes txRaw, bytes receiptRaw, bytes proofData) view returns(bytes[])
func (_ZkbridgeVerify *ZkbridgeVerifyCallerSession) SignBtcTx(txRaw []byte, receiptRaw []byte, proofData []byte) ([][]byte, error) {
	return _ZkbridgeVerify.Contract.SignBtcTx(&_ZkbridgeVerify.CallOpts, txRaw, receiptRaw, proofData)
}

// SignerPublicKey is a free data retrieval call binding the contract method 0xc093278a.
//
// Solidity: function signer_public_key() view returns(bytes)
func (_ZkbridgeVerify *ZkbridgeVerifyCaller) SignerPublicKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _ZkbridgeVerify.contract.Call(opts, &out, "signer_public_key")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// SignerPublicKey is a free data retrieval call binding the contract method 0xc093278a.
//
// Solidity: function signer_public_key() view returns(bytes)
func (_ZkbridgeVerify *ZkbridgeVerifySession) SignerPublicKey() ([]byte, error) {
	return _ZkbridgeVerify.Contract.SignerPublicKey(&_ZkbridgeVerify.CallOpts)
}

// SignerPublicKey is a free data retrieval call binding the contract method 0xc093278a.
//
// Solidity: function signer_public_key() view returns(bytes)
func (_ZkbridgeVerify *ZkbridgeVerifyCallerSession) SignerPublicKey() ([]byte, error) {
	return _ZkbridgeVerify.Contract.SignerPublicKey(&_ZkbridgeVerify.CallOpts)
}

// VerifierAddress is a free data retrieval call binding the contract method 0x18bdffbb.
//
// Solidity: function verifierAddress() view returns(address)
func (_ZkbridgeVerify *ZkbridgeVerifyCaller) VerifierAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ZkbridgeVerify.contract.Call(opts, &out, "verifierAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// VerifierAddress is a free data retrieval call binding the contract method 0x18bdffbb.
//
// Solidity: function verifierAddress() view returns(address)
func (_ZkbridgeVerify *ZkbridgeVerifySession) VerifierAddress() (common.Address, error) {
	return _ZkbridgeVerify.Contract.VerifierAddress(&_ZkbridgeVerify.CallOpts)
}

// VerifierAddress is a free data retrieval call binding the contract method 0x18bdffbb.
//
// Solidity: function verifierAddress() view returns(address)
func (_ZkbridgeVerify *ZkbridgeVerifyCallerSession) VerifierAddress() (common.Address, error) {
	return _ZkbridgeVerify.Contract.VerifierAddress(&_ZkbridgeVerify.CallOpts)
}

// AddOperator is a paid mutator transaction binding the contract method 0x9870d7fe.
//
// Solidity: function addOperator(address _new) returns()
func (_ZkbridgeVerify *ZkbridgeVerifyTransactor) AddOperator(opts *bind.TransactOpts, _new common.Address) (*types.Transaction, error) {
	return _ZkbridgeVerify.contract.Transact(opts, "addOperator", _new)
}

// AddOperator is a paid mutator transaction binding the contract method 0x9870d7fe.
//
// Solidity: function addOperator(address _new) returns()
func (_ZkbridgeVerify *ZkbridgeVerifySession) AddOperator(_new common.Address) (*types.Transaction, error) {
	return _ZkbridgeVerify.Contract.AddOperator(&_ZkbridgeVerify.TransactOpts, _new)
}

// AddOperator is a paid mutator transaction binding the contract method 0x9870d7fe.
//
// Solidity: function addOperator(address _new) returns()
func (_ZkbridgeVerify *ZkbridgeVerifyTransactorSession) AddOperator(_new common.Address) (*types.Transaction, error) {
	return _ZkbridgeVerify.Contract.AddOperator(&_ZkbridgeVerify.TransactOpts, _new)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0xac8a584a.
//
// Solidity: function removeOperator(address _new) returns()
func (_ZkbridgeVerify *ZkbridgeVerifyTransactor) RemoveOperator(opts *bind.TransactOpts, _new common.Address) (*types.Transaction, error) {
	return _ZkbridgeVerify.contract.Transact(opts, "removeOperator", _new)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0xac8a584a.
//
// Solidity: function removeOperator(address _new) returns()
func (_ZkbridgeVerify *ZkbridgeVerifySession) RemoveOperator(_new common.Address) (*types.Transaction, error) {
	return _ZkbridgeVerify.Contract.RemoveOperator(&_ZkbridgeVerify.TransactOpts, _new)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0xac8a584a.
//
// Solidity: function removeOperator(address _new) returns()
func (_ZkbridgeVerify *ZkbridgeVerifyTransactorSession) RemoveOperator(_new common.Address) (*types.Transaction, error) {
	return _ZkbridgeVerify.Contract.RemoveOperator(&_ZkbridgeVerify.TransactOpts, _new)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ZkbridgeVerify *ZkbridgeVerifyTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZkbridgeVerify.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ZkbridgeVerify *ZkbridgeVerifySession) RenounceOwnership() (*types.Transaction, error) {
	return _ZkbridgeVerify.Contract.RenounceOwnership(&_ZkbridgeVerify.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ZkbridgeVerify *ZkbridgeVerifyTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ZkbridgeVerify.Contract.RenounceOwnership(&_ZkbridgeVerify.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ZkbridgeVerify *ZkbridgeVerifyTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ZkbridgeVerify.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ZkbridgeVerify *ZkbridgeVerifySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ZkbridgeVerify.Contract.TransferOwnership(&_ZkbridgeVerify.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ZkbridgeVerify *ZkbridgeVerifyTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ZkbridgeVerify.Contract.TransferOwnership(&_ZkbridgeVerify.TransactOpts, newOwner)
}

// UpdateVerifierAddress is a paid mutator transaction binding the contract method 0x736f1618.
//
// Solidity: function updateVerifierAddress(address addr) returns()
func (_ZkbridgeVerify *ZkbridgeVerifyTransactor) UpdateVerifierAddress(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _ZkbridgeVerify.contract.Transact(opts, "updateVerifierAddress", addr)
}

// UpdateVerifierAddress is a paid mutator transaction binding the contract method 0x736f1618.
//
// Solidity: function updateVerifierAddress(address addr) returns()
func (_ZkbridgeVerify *ZkbridgeVerifySession) UpdateVerifierAddress(addr common.Address) (*types.Transaction, error) {
	return _ZkbridgeVerify.Contract.UpdateVerifierAddress(&_ZkbridgeVerify.TransactOpts, addr)
}

// UpdateVerifierAddress is a paid mutator transaction binding the contract method 0x736f1618.
//
// Solidity: function updateVerifierAddress(address addr) returns()
func (_ZkbridgeVerify *ZkbridgeVerifyTransactorSession) UpdateVerifierAddress(addr common.Address) (*types.Transaction, error) {
	return _ZkbridgeVerify.Contract.UpdateVerifierAddress(&_ZkbridgeVerify.TransactOpts, addr)
}

// ZkbridgeVerifyOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ZkbridgeVerify contract.
type ZkbridgeVerifyOwnershipTransferredIterator struct {
	Event *ZkbridgeVerifyOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ZkbridgeVerifyOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZkbridgeVerifyOwnershipTransferred)
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
		it.Event = new(ZkbridgeVerifyOwnershipTransferred)
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
func (it *ZkbridgeVerifyOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZkbridgeVerifyOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZkbridgeVerifyOwnershipTransferred represents a OwnershipTransferred event raised by the ZkbridgeVerify contract.
type ZkbridgeVerifyOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ZkbridgeVerify *ZkbridgeVerifyFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ZkbridgeVerifyOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ZkbridgeVerify.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ZkbridgeVerifyOwnershipTransferredIterator{contract: _ZkbridgeVerify.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ZkbridgeVerify *ZkbridgeVerifyFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ZkbridgeVerifyOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ZkbridgeVerify.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZkbridgeVerifyOwnershipTransferred)
				if err := _ZkbridgeVerify.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_ZkbridgeVerify *ZkbridgeVerifyFilterer) ParseOwnershipTransferred(log types.Log) (*ZkbridgeVerifyOwnershipTransferred, error) {
	event := new(ZkbridgeVerifyOwnershipTransferred)
	if err := _ZkbridgeVerify.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
