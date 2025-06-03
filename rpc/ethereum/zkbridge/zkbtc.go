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

// ZkbtcMetaData contains all meta data concerning the Zkbtc contract.
var ZkbtcMetaData = &bind.MetaData{
	ABI: "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"},{\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}]",
}

// ZkbtcABI is the input ABI used to generate the binding from.
// Deprecated: Use ZkbtcMetaData.ABI instead.
var ZkbtcABI = ZkbtcMetaData.ABI

// Zkbtc is an auto generated Go binding around an Ethereum contract.
type Zkbtc struct {
	ZkbtcCaller     // Read-only binding to the contract
	ZkbtcTransactor // Write-only binding to the contract
	ZkbtcFilterer   // Log filterer for contract events
}

// ZkbtcCaller is an auto generated read-only Go binding around an Ethereum contract.
type ZkbtcCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZkbtcTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ZkbtcTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZkbtcFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ZkbtcFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZkbtcSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ZkbtcSession struct {
	Contract     *Zkbtc            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ZkbtcCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ZkbtcCallerSession struct {
	Contract *ZkbtcCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ZkbtcTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ZkbtcTransactorSession struct {
	Contract     *ZkbtcTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ZkbtcRaw is an auto generated low-level Go binding around an Ethereum contract.
type ZkbtcRaw struct {
	Contract *Zkbtc // Generic contract binding to access the raw methods on
}

// ZkbtcCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ZkbtcCallerRaw struct {
	Contract *ZkbtcCaller // Generic read-only contract binding to access the raw methods on
}

// ZkbtcTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ZkbtcTransactorRaw struct {
	Contract *ZkbtcTransactor // Generic write-only contract binding to access the raw methods on
}

// NewZkbtc creates a new instance of Zkbtc, bound to a specific deployed contract.
func NewZkbtc(address common.Address, backend bind.ContractBackend) (*Zkbtc, error) {
	contract, err := bindZkbtc(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Zkbtc{ZkbtcCaller: ZkbtcCaller{contract: contract}, ZkbtcTransactor: ZkbtcTransactor{contract: contract}, ZkbtcFilterer: ZkbtcFilterer{contract: contract}}, nil
}

// NewZkbtcCaller creates a new read-only instance of Zkbtc, bound to a specific deployed contract.
func NewZkbtcCaller(address common.Address, caller bind.ContractCaller) (*ZkbtcCaller, error) {
	contract, err := bindZkbtc(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ZkbtcCaller{contract: contract}, nil
}

// NewZkbtcTransactor creates a new write-only instance of Zkbtc, bound to a specific deployed contract.
func NewZkbtcTransactor(address common.Address, transactor bind.ContractTransactor) (*ZkbtcTransactor, error) {
	contract, err := bindZkbtc(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ZkbtcTransactor{contract: contract}, nil
}

// NewZkbtcFilterer creates a new log filterer instance of Zkbtc, bound to a specific deployed contract.
func NewZkbtcFilterer(address common.Address, filterer bind.ContractFilterer) (*ZkbtcFilterer, error) {
	contract, err := bindZkbtc(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ZkbtcFilterer{contract: contract}, nil
}

// bindZkbtc binds a generic wrapper to an already deployed contract.
func bindZkbtc(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ZkbtcMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Zkbtc *ZkbtcRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Zkbtc.Contract.ZkbtcCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Zkbtc *ZkbtcRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Zkbtc.Contract.ZkbtcTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Zkbtc *ZkbtcRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Zkbtc.Contract.ZkbtcTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Zkbtc *ZkbtcCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Zkbtc.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Zkbtc *ZkbtcTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Zkbtc.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Zkbtc *ZkbtcTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Zkbtc.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address _owner, address _spender) view returns(uint256)
func (_Zkbtc *ZkbtcCaller) Allowance(opts *bind.CallOpts, _owner common.Address, _spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Zkbtc.contract.Call(opts, &out, "allowance", _owner, _spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address _owner, address _spender) view returns(uint256)
func (_Zkbtc *ZkbtcSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _Zkbtc.Contract.Allowance(&_Zkbtc.CallOpts, _owner, _spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address _owner, address _spender) view returns(uint256)
func (_Zkbtc *ZkbtcCallerSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _Zkbtc.Contract.Allowance(&_Zkbtc.CallOpts, _owner, _spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _owner) view returns(uint256 balance)
func (_Zkbtc *ZkbtcCaller) BalanceOf(opts *bind.CallOpts, _owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Zkbtc.contract.Call(opts, &out, "balanceOf", _owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _owner) view returns(uint256 balance)
func (_Zkbtc *ZkbtcSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _Zkbtc.Contract.BalanceOf(&_Zkbtc.CallOpts, _owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _owner) view returns(uint256 balance)
func (_Zkbtc *ZkbtcCallerSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _Zkbtc.Contract.BalanceOf(&_Zkbtc.CallOpts, _owner)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_Zkbtc *ZkbtcCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Zkbtc.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_Zkbtc *ZkbtcSession) Decimals() (uint8, error) {
	return _Zkbtc.Contract.Decimals(&_Zkbtc.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_Zkbtc *ZkbtcCallerSession) Decimals() (uint8, error) {
	return _Zkbtc.Contract.Decimals(&_Zkbtc.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Zkbtc *ZkbtcCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Zkbtc.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Zkbtc *ZkbtcSession) Name() (string, error) {
	return _Zkbtc.Contract.Name(&_Zkbtc.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Zkbtc *ZkbtcCallerSession) Name() (string, error) {
	return _Zkbtc.Contract.Name(&_Zkbtc.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Zkbtc *ZkbtcCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Zkbtc.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Zkbtc *ZkbtcSession) Symbol() (string, error) {
	return _Zkbtc.Contract.Symbol(&_Zkbtc.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Zkbtc *ZkbtcCallerSession) Symbol() (string, error) {
	return _Zkbtc.Contract.Symbol(&_Zkbtc.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Zkbtc *ZkbtcCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Zkbtc.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Zkbtc *ZkbtcSession) TotalSupply() (*big.Int, error) {
	return _Zkbtc.Contract.TotalSupply(&_Zkbtc.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Zkbtc *ZkbtcCallerSession) TotalSupply() (*big.Int, error) {
	return _Zkbtc.Contract.TotalSupply(&_Zkbtc.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address _spender, uint256 _value) returns(bool)
func (_Zkbtc *ZkbtcTransactor) Approve(opts *bind.TransactOpts, _spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Zkbtc.contract.Transact(opts, "approve", _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address _spender, uint256 _value) returns(bool)
func (_Zkbtc *ZkbtcSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Zkbtc.Contract.Approve(&_Zkbtc.TransactOpts, _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address _spender, uint256 _value) returns(bool)
func (_Zkbtc *ZkbtcTransactorSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Zkbtc.Contract.Approve(&_Zkbtc.TransactOpts, _spender, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _to, uint256 _value) returns(bool)
func (_Zkbtc *ZkbtcTransactor) Transfer(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Zkbtc.contract.Transact(opts, "transfer", _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _to, uint256 _value) returns(bool)
func (_Zkbtc *ZkbtcSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Zkbtc.Contract.Transfer(&_Zkbtc.TransactOpts, _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _to, uint256 _value) returns(bool)
func (_Zkbtc *ZkbtcTransactorSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Zkbtc.Contract.Transfer(&_Zkbtc.TransactOpts, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address _from, address _to, uint256 _value) returns(bool)
func (_Zkbtc *ZkbtcTransactor) TransferFrom(opts *bind.TransactOpts, _from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Zkbtc.contract.Transact(opts, "transferFrom", _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address _from, address _to, uint256 _value) returns(bool)
func (_Zkbtc *ZkbtcSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Zkbtc.Contract.TransferFrom(&_Zkbtc.TransactOpts, _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address _from, address _to, uint256 _value) returns(bool)
func (_Zkbtc *ZkbtcTransactorSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Zkbtc.Contract.TransferFrom(&_Zkbtc.TransactOpts, _from, _to, _value)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Zkbtc *ZkbtcTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _Zkbtc.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Zkbtc *ZkbtcSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Zkbtc.Contract.Fallback(&_Zkbtc.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Zkbtc *ZkbtcTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Zkbtc.Contract.Fallback(&_Zkbtc.TransactOpts, calldata)
}

// ZkbtcApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the Zkbtc contract.
type ZkbtcApprovalIterator struct {
	Event *ZkbtcApproval // Event containing the contract specifics and raw log

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
func (it *ZkbtcApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZkbtcApproval)
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
		it.Event = new(ZkbtcApproval)
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
func (it *ZkbtcApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZkbtcApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZkbtcApproval represents a Approval event raised by the Zkbtc contract.
type ZkbtcApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Zkbtc *ZkbtcFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*ZkbtcApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Zkbtc.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &ZkbtcApprovalIterator{contract: _Zkbtc.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Zkbtc *ZkbtcFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ZkbtcApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Zkbtc.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZkbtcApproval)
				if err := _Zkbtc.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Zkbtc *ZkbtcFilterer) ParseApproval(log types.Log) (*ZkbtcApproval, error) {
	event := new(ZkbtcApproval)
	if err := _Zkbtc.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZkbtcTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the Zkbtc contract.
type ZkbtcTransferIterator struct {
	Event *ZkbtcTransfer // Event containing the contract specifics and raw log

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
func (it *ZkbtcTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZkbtcTransfer)
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
		it.Event = new(ZkbtcTransfer)
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
func (it *ZkbtcTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZkbtcTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZkbtcTransfer represents a Transfer event raised by the Zkbtc contract.
type ZkbtcTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Zkbtc *ZkbtcFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ZkbtcTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Zkbtc.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ZkbtcTransferIterator{contract: _Zkbtc.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Zkbtc *ZkbtcFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ZkbtcTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Zkbtc.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZkbtcTransfer)
				if err := _Zkbtc.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Zkbtc *ZkbtcFilterer) ParseTransfer(log types.Log) (*ZkbtcTransfer, error) {
	event := new(ZkbtcTransfer)
	if err := _Zkbtc.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
