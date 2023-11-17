// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package zkbridge_signer

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

// ZkbridgeSignerMetaData contains all meta data concerning the ZkbridgeSigner contract.
var ZkbridgeSignerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"count\",\"type\":\"uint8\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"txId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"proofData\",\"type\":\"bytes\"}],\"name\":\"InvalidRedeemOrProof\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"addKeypairOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_new\",\"type\":\"address\"}],\"name\":\"addOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPublicKeys\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"\",\"type\":\"bytes[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isOperator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_new\",\"type\":\"address\"}],\"name\":\"removeOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"txId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"minerReward\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"sigHashs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"currentScRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"proofData\",\"type\":\"bytes\"}],\"name\":\"signBtcTx\",\"outputs\":[{\"internalType\":\"bytes[][]\",\"name\":\"\",\"type\":\"bytes[][]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"signerAddresses\",\"outputs\":[{\"internalType\":\"contractBTCKeypair\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractBTCKeypair[]\",\"name\":\"pairs\",\"type\":\"address[]\"}],\"name\":\"updateKeypairs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIEthTxVerifier\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"updateVerifierAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"verifierAddress\",\"outputs\":[{\"internalType\":\"contractIEthTxVerifier\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// ZkbridgeSignerABI is the input ABI used to generate the binding from.
// Deprecated: Use ZkbridgeSignerMetaData.ABI instead.
var ZkbridgeSignerABI = ZkbridgeSignerMetaData.ABI

// ZkbridgeSigner is an auto generated Go binding around an Ethereum contract.
type ZkbridgeSigner struct {
	ZkbridgeSignerCaller     // Read-only binding to the contract
	ZkbridgeSignerTransactor // Write-only binding to the contract
	ZkbridgeSignerFilterer   // Log filterer for contract events
}

// ZkbridgeSignerCaller is an auto generated read-only Go binding around an Ethereum contract.
type ZkbridgeSignerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZkbridgeSignerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ZkbridgeSignerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZkbridgeSignerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ZkbridgeSignerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZkbridgeSignerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ZkbridgeSignerSession struct {
	Contract     *ZkbridgeSigner   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ZkbridgeSignerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ZkbridgeSignerCallerSession struct {
	Contract *ZkbridgeSignerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// ZkbridgeSignerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ZkbridgeSignerTransactorSession struct {
	Contract     *ZkbridgeSignerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// ZkbridgeSignerRaw is an auto generated low-level Go binding around an Ethereum contract.
type ZkbridgeSignerRaw struct {
	Contract *ZkbridgeSigner // Generic contract binding to access the raw methods on
}

// ZkbridgeSignerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ZkbridgeSignerCallerRaw struct {
	Contract *ZkbridgeSignerCaller // Generic read-only contract binding to access the raw methods on
}

// ZkbridgeSignerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ZkbridgeSignerTransactorRaw struct {
	Contract *ZkbridgeSignerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewZkbridgeSigner creates a new instance of ZkbridgeSigner, bound to a specific deployed contract.
func NewZkbridgeSigner(address common.Address, backend bind.ContractBackend) (*ZkbridgeSigner, error) {
	contract, err := bindZkbridgeSigner(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ZkbridgeSigner{ZkbridgeSignerCaller: ZkbridgeSignerCaller{contract: contract}, ZkbridgeSignerTransactor: ZkbridgeSignerTransactor{contract: contract}, ZkbridgeSignerFilterer: ZkbridgeSignerFilterer{contract: contract}}, nil
}

// NewZkbridgeSignerCaller creates a new read-only instance of ZkbridgeSigner, bound to a specific deployed contract.
func NewZkbridgeSignerCaller(address common.Address, caller bind.ContractCaller) (*ZkbridgeSignerCaller, error) {
	contract, err := bindZkbridgeSigner(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ZkbridgeSignerCaller{contract: contract}, nil
}

// NewZkbridgeSignerTransactor creates a new write-only instance of ZkbridgeSigner, bound to a specific deployed contract.
func NewZkbridgeSignerTransactor(address common.Address, transactor bind.ContractTransactor) (*ZkbridgeSignerTransactor, error) {
	contract, err := bindZkbridgeSigner(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ZkbridgeSignerTransactor{contract: contract}, nil
}

// NewZkbridgeSignerFilterer creates a new log filterer instance of ZkbridgeSigner, bound to a specific deployed contract.
func NewZkbridgeSignerFilterer(address common.Address, filterer bind.ContractFilterer) (*ZkbridgeSignerFilterer, error) {
	contract, err := bindZkbridgeSigner(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ZkbridgeSignerFilterer{contract: contract}, nil
}

// bindZkbridgeSigner binds a generic wrapper to an already deployed contract.
func bindZkbridgeSigner(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ZkbridgeSignerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ZkbridgeSigner *ZkbridgeSignerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ZkbridgeSigner.Contract.ZkbridgeSignerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ZkbridgeSigner *ZkbridgeSignerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZkbridgeSigner.Contract.ZkbridgeSignerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ZkbridgeSigner *ZkbridgeSignerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ZkbridgeSigner.Contract.ZkbridgeSignerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ZkbridgeSigner *ZkbridgeSignerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ZkbridgeSigner.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ZkbridgeSigner *ZkbridgeSignerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZkbridgeSigner.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ZkbridgeSigner *ZkbridgeSignerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ZkbridgeSigner.Contract.contract.Transact(opts, method, params...)
}

// GetPublicKeys is a free data retrieval call binding the contract method 0x15285fed.
//
// Solidity: function getPublicKeys() view returns(bytes[])
func (_ZkbridgeSigner *ZkbridgeSignerCaller) GetPublicKeys(opts *bind.CallOpts) ([][]byte, error) {
	var out []interface{}
	err := _ZkbridgeSigner.contract.Call(opts, &out, "getPublicKeys")

	if err != nil {
		return *new([][]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][]byte)).(*[][]byte)

	return out0, err

}

// GetPublicKeys is a free data retrieval call binding the contract method 0x15285fed.
//
// Solidity: function getPublicKeys() view returns(bytes[])
func (_ZkbridgeSigner *ZkbridgeSignerSession) GetPublicKeys() ([][]byte, error) {
	return _ZkbridgeSigner.Contract.GetPublicKeys(&_ZkbridgeSigner.CallOpts)
}

// GetPublicKeys is a free data retrieval call binding the contract method 0x15285fed.
//
// Solidity: function getPublicKeys() view returns(bytes[])
func (_ZkbridgeSigner *ZkbridgeSignerCallerSession) GetPublicKeys() ([][]byte, error) {
	return _ZkbridgeSigner.Contract.GetPublicKeys(&_ZkbridgeSigner.CallOpts)
}

// IsOperator is a free data retrieval call binding the contract method 0x6d70f7ae.
//
// Solidity: function isOperator(address addr) view returns(bool)
func (_ZkbridgeSigner *ZkbridgeSignerCaller) IsOperator(opts *bind.CallOpts, addr common.Address) (bool, error) {
	var out []interface{}
	err := _ZkbridgeSigner.contract.Call(opts, &out, "isOperator", addr)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOperator is a free data retrieval call binding the contract method 0x6d70f7ae.
//
// Solidity: function isOperator(address addr) view returns(bool)
func (_ZkbridgeSigner *ZkbridgeSignerSession) IsOperator(addr common.Address) (bool, error) {
	return _ZkbridgeSigner.Contract.IsOperator(&_ZkbridgeSigner.CallOpts, addr)
}

// IsOperator is a free data retrieval call binding the contract method 0x6d70f7ae.
//
// Solidity: function isOperator(address addr) view returns(bool)
func (_ZkbridgeSigner *ZkbridgeSignerCallerSession) IsOperator(addr common.Address) (bool, error) {
	return _ZkbridgeSigner.Contract.IsOperator(&_ZkbridgeSigner.CallOpts, addr)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ZkbridgeSigner *ZkbridgeSignerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ZkbridgeSigner.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ZkbridgeSigner *ZkbridgeSignerSession) Owner() (common.Address, error) {
	return _ZkbridgeSigner.Contract.Owner(&_ZkbridgeSigner.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ZkbridgeSigner *ZkbridgeSignerCallerSession) Owner() (common.Address, error) {
	return _ZkbridgeSigner.Contract.Owner(&_ZkbridgeSigner.CallOpts)
}

// SignBtcTx is a free data retrieval call binding the contract method 0x60825364.
//
// Solidity: function signBtcTx(bytes32 txId, uint256 minerReward, bytes32[] sigHashs, bytes32 currentScRoot, bytes proofData) view returns(bytes[][])
func (_ZkbridgeSigner *ZkbridgeSignerCaller) SignBtcTx(opts *bind.CallOpts, txId [32]byte, minerReward *big.Int, sigHashs [][32]byte, currentScRoot [32]byte, proofData []byte) ([][][]byte, error) {
	var out []interface{}
	err := _ZkbridgeSigner.contract.Call(opts, &out, "signBtcTx", txId, minerReward, sigHashs, currentScRoot, proofData)

	if err != nil {
		return *new([][][]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][][]byte)).(*[][][]byte)

	return out0, err

}

// SignBtcTx is a free data retrieval call binding the contract method 0x60825364.
//
// Solidity: function signBtcTx(bytes32 txId, uint256 minerReward, bytes32[] sigHashs, bytes32 currentScRoot, bytes proofData) view returns(bytes[][])
func (_ZkbridgeSigner *ZkbridgeSignerSession) SignBtcTx(txId [32]byte, minerReward *big.Int, sigHashs [][32]byte, currentScRoot [32]byte, proofData []byte) ([][][]byte, error) {
	return _ZkbridgeSigner.Contract.SignBtcTx(&_ZkbridgeSigner.CallOpts, txId, minerReward, sigHashs, currentScRoot, proofData)
}

// SignBtcTx is a free data retrieval call binding the contract method 0x60825364.
//
// Solidity: function signBtcTx(bytes32 txId, uint256 minerReward, bytes32[] sigHashs, bytes32 currentScRoot, bytes proofData) view returns(bytes[][])
func (_ZkbridgeSigner *ZkbridgeSignerCallerSession) SignBtcTx(txId [32]byte, minerReward *big.Int, sigHashs [][32]byte, currentScRoot [32]byte, proofData []byte) ([][][]byte, error) {
	return _ZkbridgeSigner.Contract.SignBtcTx(&_ZkbridgeSigner.CallOpts, txId, minerReward, sigHashs, currentScRoot, proofData)
}

// SignerAddresses is a free data retrieval call binding the contract method 0xb72217d6.
//
// Solidity: function signerAddresses(uint256 ) view returns(address)
func (_ZkbridgeSigner *ZkbridgeSignerCaller) SignerAddresses(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _ZkbridgeSigner.contract.Call(opts, &out, "signerAddresses", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SignerAddresses is a free data retrieval call binding the contract method 0xb72217d6.
//
// Solidity: function signerAddresses(uint256 ) view returns(address)
func (_ZkbridgeSigner *ZkbridgeSignerSession) SignerAddresses(arg0 *big.Int) (common.Address, error) {
	return _ZkbridgeSigner.Contract.SignerAddresses(&_ZkbridgeSigner.CallOpts, arg0)
}

// SignerAddresses is a free data retrieval call binding the contract method 0xb72217d6.
//
// Solidity: function signerAddresses(uint256 ) view returns(address)
func (_ZkbridgeSigner *ZkbridgeSignerCallerSession) SignerAddresses(arg0 *big.Int) (common.Address, error) {
	return _ZkbridgeSigner.Contract.SignerAddresses(&_ZkbridgeSigner.CallOpts, arg0)
}

// VerifierAddress is a free data retrieval call binding the contract method 0x18bdffbb.
//
// Solidity: function verifierAddress() view returns(address)
func (_ZkbridgeSigner *ZkbridgeSignerCaller) VerifierAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ZkbridgeSigner.contract.Call(opts, &out, "verifierAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// VerifierAddress is a free data retrieval call binding the contract method 0x18bdffbb.
//
// Solidity: function verifierAddress() view returns(address)
func (_ZkbridgeSigner *ZkbridgeSignerSession) VerifierAddress() (common.Address, error) {
	return _ZkbridgeSigner.Contract.VerifierAddress(&_ZkbridgeSigner.CallOpts)
}

// VerifierAddress is a free data retrieval call binding the contract method 0x18bdffbb.
//
// Solidity: function verifierAddress() view returns(address)
func (_ZkbridgeSigner *ZkbridgeSignerCallerSession) VerifierAddress() (common.Address, error) {
	return _ZkbridgeSigner.Contract.VerifierAddress(&_ZkbridgeSigner.CallOpts)
}

// AddKeypairOperator is a paid mutator transaction binding the contract method 0x9793079f.
//
// Solidity: function addKeypairOperator(address operator) returns()
func (_ZkbridgeSigner *ZkbridgeSignerTransactor) AddKeypairOperator(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _ZkbridgeSigner.contract.Transact(opts, "addKeypairOperator", operator)
}

// AddKeypairOperator is a paid mutator transaction binding the contract method 0x9793079f.
//
// Solidity: function addKeypairOperator(address operator) returns()
func (_ZkbridgeSigner *ZkbridgeSignerSession) AddKeypairOperator(operator common.Address) (*types.Transaction, error) {
	return _ZkbridgeSigner.Contract.AddKeypairOperator(&_ZkbridgeSigner.TransactOpts, operator)
}

// AddKeypairOperator is a paid mutator transaction binding the contract method 0x9793079f.
//
// Solidity: function addKeypairOperator(address operator) returns()
func (_ZkbridgeSigner *ZkbridgeSignerTransactorSession) AddKeypairOperator(operator common.Address) (*types.Transaction, error) {
	return _ZkbridgeSigner.Contract.AddKeypairOperator(&_ZkbridgeSigner.TransactOpts, operator)
}

// AddOperator is a paid mutator transaction binding the contract method 0x9870d7fe.
//
// Solidity: function addOperator(address _new) returns()
func (_ZkbridgeSigner *ZkbridgeSignerTransactor) AddOperator(opts *bind.TransactOpts, _new common.Address) (*types.Transaction, error) {
	return _ZkbridgeSigner.contract.Transact(opts, "addOperator", _new)
}

// AddOperator is a paid mutator transaction binding the contract method 0x9870d7fe.
//
// Solidity: function addOperator(address _new) returns()
func (_ZkbridgeSigner *ZkbridgeSignerSession) AddOperator(_new common.Address) (*types.Transaction, error) {
	return _ZkbridgeSigner.Contract.AddOperator(&_ZkbridgeSigner.TransactOpts, _new)
}

// AddOperator is a paid mutator transaction binding the contract method 0x9870d7fe.
//
// Solidity: function addOperator(address _new) returns()
func (_ZkbridgeSigner *ZkbridgeSignerTransactorSession) AddOperator(_new common.Address) (*types.Transaction, error) {
	return _ZkbridgeSigner.Contract.AddOperator(&_ZkbridgeSigner.TransactOpts, _new)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0xac8a584a.
//
// Solidity: function removeOperator(address _new) returns()
func (_ZkbridgeSigner *ZkbridgeSignerTransactor) RemoveOperator(opts *bind.TransactOpts, _new common.Address) (*types.Transaction, error) {
	return _ZkbridgeSigner.contract.Transact(opts, "removeOperator", _new)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0xac8a584a.
//
// Solidity: function removeOperator(address _new) returns()
func (_ZkbridgeSigner *ZkbridgeSignerSession) RemoveOperator(_new common.Address) (*types.Transaction, error) {
	return _ZkbridgeSigner.Contract.RemoveOperator(&_ZkbridgeSigner.TransactOpts, _new)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0xac8a584a.
//
// Solidity: function removeOperator(address _new) returns()
func (_ZkbridgeSigner *ZkbridgeSignerTransactorSession) RemoveOperator(_new common.Address) (*types.Transaction, error) {
	return _ZkbridgeSigner.Contract.RemoveOperator(&_ZkbridgeSigner.TransactOpts, _new)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ZkbridgeSigner *ZkbridgeSignerTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZkbridgeSigner.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ZkbridgeSigner *ZkbridgeSignerSession) RenounceOwnership() (*types.Transaction, error) {
	return _ZkbridgeSigner.Contract.RenounceOwnership(&_ZkbridgeSigner.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ZkbridgeSigner *ZkbridgeSignerTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ZkbridgeSigner.Contract.RenounceOwnership(&_ZkbridgeSigner.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ZkbridgeSigner *ZkbridgeSignerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ZkbridgeSigner.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ZkbridgeSigner *ZkbridgeSignerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ZkbridgeSigner.Contract.TransferOwnership(&_ZkbridgeSigner.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ZkbridgeSigner *ZkbridgeSignerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ZkbridgeSigner.Contract.TransferOwnership(&_ZkbridgeSigner.TransactOpts, newOwner)
}

// UpdateKeypairs is a paid mutator transaction binding the contract method 0xc897c995.
//
// Solidity: function updateKeypairs(address[] pairs) returns()
func (_ZkbridgeSigner *ZkbridgeSignerTransactor) UpdateKeypairs(opts *bind.TransactOpts, pairs []common.Address) (*types.Transaction, error) {
	return _ZkbridgeSigner.contract.Transact(opts, "updateKeypairs", pairs)
}

// UpdateKeypairs is a paid mutator transaction binding the contract method 0xc897c995.
//
// Solidity: function updateKeypairs(address[] pairs) returns()
func (_ZkbridgeSigner *ZkbridgeSignerSession) UpdateKeypairs(pairs []common.Address) (*types.Transaction, error) {
	return _ZkbridgeSigner.Contract.UpdateKeypairs(&_ZkbridgeSigner.TransactOpts, pairs)
}

// UpdateKeypairs is a paid mutator transaction binding the contract method 0xc897c995.
//
// Solidity: function updateKeypairs(address[] pairs) returns()
func (_ZkbridgeSigner *ZkbridgeSignerTransactorSession) UpdateKeypairs(pairs []common.Address) (*types.Transaction, error) {
	return _ZkbridgeSigner.Contract.UpdateKeypairs(&_ZkbridgeSigner.TransactOpts, pairs)
}

// UpdateVerifierAddress is a paid mutator transaction binding the contract method 0x736f1618.
//
// Solidity: function updateVerifierAddress(address addr) returns()
func (_ZkbridgeSigner *ZkbridgeSignerTransactor) UpdateVerifierAddress(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _ZkbridgeSigner.contract.Transact(opts, "updateVerifierAddress", addr)
}

// UpdateVerifierAddress is a paid mutator transaction binding the contract method 0x736f1618.
//
// Solidity: function updateVerifierAddress(address addr) returns()
func (_ZkbridgeSigner *ZkbridgeSignerSession) UpdateVerifierAddress(addr common.Address) (*types.Transaction, error) {
	return _ZkbridgeSigner.Contract.UpdateVerifierAddress(&_ZkbridgeSigner.TransactOpts, addr)
}

// UpdateVerifierAddress is a paid mutator transaction binding the contract method 0x736f1618.
//
// Solidity: function updateVerifierAddress(address addr) returns()
func (_ZkbridgeSigner *ZkbridgeSignerTransactorSession) UpdateVerifierAddress(addr common.Address) (*types.Transaction, error) {
	return _ZkbridgeSigner.Contract.UpdateVerifierAddress(&_ZkbridgeSigner.TransactOpts, addr)
}

// ZkbridgeSignerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ZkbridgeSigner contract.
type ZkbridgeSignerOwnershipTransferredIterator struct {
	Event *ZkbridgeSignerOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ZkbridgeSignerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZkbridgeSignerOwnershipTransferred)
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
		it.Event = new(ZkbridgeSignerOwnershipTransferred)
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
func (it *ZkbridgeSignerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZkbridgeSignerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZkbridgeSignerOwnershipTransferred represents a OwnershipTransferred event raised by the ZkbridgeSigner contract.
type ZkbridgeSignerOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ZkbridgeSigner *ZkbridgeSignerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ZkbridgeSignerOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ZkbridgeSigner.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ZkbridgeSignerOwnershipTransferredIterator{contract: _ZkbridgeSigner.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ZkbridgeSigner *ZkbridgeSignerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ZkbridgeSignerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ZkbridgeSigner.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZkbridgeSignerOwnershipTransferred)
				if err := _ZkbridgeSigner.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_ZkbridgeSigner *ZkbridgeSignerFilterer) ParseOwnershipTransferred(log types.Log) (*ZkbridgeSignerOwnershipTransferred, error) {
	event := new(ZkbridgeSignerOwnershipTransferred)
	if err := _ZkbridgeSigner.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
