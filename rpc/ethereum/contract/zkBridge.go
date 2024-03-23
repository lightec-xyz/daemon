// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package zkBridge

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

// UTXOManagerInterfaceTxOut is an auto generated low-level Go binding around an user-defined struct.
type UTXOManagerInterfaceTxOut struct {
	Amount       uint64
	ScriptPubKey []byte
}

// ZkBridgeMetaData contains all meta data concerning the ZkBridge contract.
var ZkBridgeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_feeAccount\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_changeLockScript\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"_multiSigScript\",\"type\":\"bytes\"},{\"internalType\":\"contractUTXOManagerInterface\",\"name\":\"_utxoAddress\",\"type\":\"address\"},{\"internalType\":\"contractLITInterface\",\"name\":\"_litAddress\",\"type\":\"address\"},{\"internalType\":\"contractzkBTCInterface\",\"name\":\"_zkBTCAddress\",\"type\":\"address\"},{\"internalType\":\"contractEconomicVariationInterface\",\"name\":\"_variationAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"DepositAccountIsBridge\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"proofData\",\"type\":\"bytes\"}],\"name\":\"InvalidChangeProof\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"depositAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minDepositAmount\",\"type\":\"uint256\"}],\"name\":\"InvalidDepositAmount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"changeLockScript\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"proofData\",\"type\":\"bytes\"}],\"name\":\"InvalidDepositProof\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"lockScriptLength\",\"type\":\"uint256\"}],\"name\":\"InvalidLockScriptLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"lockScript\",\"type\":\"bytes\"}],\"name\":\"LockScriptIsChange\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"rawTx\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes32[]\",\"name\":\"sigHashs\",\"type\":\"bytes32[]\"}],\"name\":\"CreateRedeemUnsignedTx\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_new\",\"type\":\"address\"}],\"name\":\"addOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"changeLockScript\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"economicVariation\",\"outputs\":[{\"internalType\":\"contractEconomicVariationInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feeAccount\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isOperator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"litToken\",\"outputs\":[{\"internalType\":\"contractLITInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minDepositAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"multiSigScript\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_new\",\"type\":\"address\"}],\"name\":\"removeOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalDepositAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"utxoManager\",\"outputs\":[{\"internalType\":\"contractUTXOManagerInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"zkBTCToken\",\"outputs\":[{\"internalType\":\"contractzkBTCInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractEconomicVariationInterface\",\"name\":\"_newAddress\",\"type\":\"address\"}],\"name\":\"updateEconomicAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractUTXOManagerInterface\",\"name\":\"_utxoAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_changeLockScript\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"_multiSigScript\",\"type\":\"bytes\"}],\"name\":\"updateUtxoAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"receiveAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"proofData\",\"type\":\"bytes\"}],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"targetAmount\",\"type\":\"uint256\"}],\"name\":\"estimateTxWeight\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"receiveLockScript\",\"type\":\"bytes\"}],\"name\":\"checkReceiveLockScript\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"redeemAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"btcMinerFee\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"receiveLockScript\",\"type\":\"bytes\"}],\"name\":\"redeem\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"userAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"changeAmount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"receiveLockScript\",\"type\":\"bytes\"}],\"name\":\"getTxOuts\",\"outputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"amount\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"scriptPubKey\",\"type\":\"bytes\"}],\"internalType\":\"structUTXOManagerInterface.TxOut[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"proofData\",\"type\":\"bytes\"}],\"name\":\"updateChange\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"proofData\",\"type\":\"bytes\"}],\"name\":\"verifyDepositProof\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"proofData\",\"type\":\"bytes\"}],\"name\":\"verifyChangeProof\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"getBridgeDepositToll\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"userAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"feeAmount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"getBridgeRedeemToll\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"userAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"feeAmount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"feeAmount\",\"type\":\"uint256\"}],\"name\":\"getDepositLITMintAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"feeAmount\",\"type\":\"uint256\"}],\"name\":\"getRedeemLITMintAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// ZkBridgeABI is the input ABI used to generate the binding from.
// Deprecated: Use ZkBridgeMetaData.ABI instead.
var ZkBridgeABI = ZkBridgeMetaData.ABI

// ZkBridge is an auto generated Go binding around an Ethereum contract.
type ZkBridge struct {
	ZkBridgeCaller     // Read-only binding to the contract
	ZkBridgeTransactor // Write-only binding to the contract
	ZkBridgeFilterer   // Log filterer for contract events
}

// ZkBridgeCaller is an auto generated read-only Go binding around an Ethereum contract.
type ZkBridgeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZkBridgeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ZkBridgeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZkBridgeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ZkBridgeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZkBridgeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ZkBridgeSession struct {
	Contract     *ZkBridge         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ZkBridgeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ZkBridgeCallerSession struct {
	Contract *ZkBridgeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// ZkBridgeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ZkBridgeTransactorSession struct {
	Contract     *ZkBridgeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ZkBridgeRaw is an auto generated low-level Go binding around an Ethereum contract.
type ZkBridgeRaw struct {
	Contract *ZkBridge // Generic contract binding to access the raw methods on
}

// ZkBridgeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ZkBridgeCallerRaw struct {
	Contract *ZkBridgeCaller // Generic read-only contract binding to access the raw methods on
}

// ZkBridgeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ZkBridgeTransactorRaw struct {
	Contract *ZkBridgeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewZkBridge creates a new instance of ZkBridge, bound to a specific deployed contract.
func NewZkBridge(address common.Address, backend bind.ContractBackend) (*ZkBridge, error) {
	contract, err := bindZkBridge(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ZkBridge{ZkBridgeCaller: ZkBridgeCaller{contract: contract}, ZkBridgeTransactor: ZkBridgeTransactor{contract: contract}, ZkBridgeFilterer: ZkBridgeFilterer{contract: contract}}, nil
}

// NewZkBridgeCaller creates a new read-only instance of ZkBridge, bound to a specific deployed contract.
func NewZkBridgeCaller(address common.Address, caller bind.ContractCaller) (*ZkBridgeCaller, error) {
	contract, err := bindZkBridge(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ZkBridgeCaller{contract: contract}, nil
}

// NewZkBridgeTransactor creates a new write-only instance of ZkBridge, bound to a specific deployed contract.
func NewZkBridgeTransactor(address common.Address, transactor bind.ContractTransactor) (*ZkBridgeTransactor, error) {
	contract, err := bindZkBridge(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ZkBridgeTransactor{contract: contract}, nil
}

// NewZkBridgeFilterer creates a new log filterer instance of ZkBridge, bound to a specific deployed contract.
func NewZkBridgeFilterer(address common.Address, filterer bind.ContractFilterer) (*ZkBridgeFilterer, error) {
	contract, err := bindZkBridge(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ZkBridgeFilterer{contract: contract}, nil
}

// bindZkBridge binds a generic wrapper to an already deployed contract.
func bindZkBridge(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ZkBridgeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ZkBridge *ZkBridgeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ZkBridge.Contract.ZkBridgeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ZkBridge *ZkBridgeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZkBridge.Contract.ZkBridgeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ZkBridge *ZkBridgeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ZkBridge.Contract.ZkBridgeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ZkBridge *ZkBridgeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ZkBridge.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ZkBridge *ZkBridgeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZkBridge.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ZkBridge *ZkBridgeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ZkBridge.Contract.contract.Transact(opts, method, params...)
}

// ChangeLockScript is a free data retrieval call binding the contract method 0x8e0a98af.
//
// Solidity: function changeLockScript() view returns(bytes)
func (_ZkBridge *ZkBridgeCaller) ChangeLockScript(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _ZkBridge.contract.Call(opts, &out, "changeLockScript")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// ChangeLockScript is a free data retrieval call binding the contract method 0x8e0a98af.
//
// Solidity: function changeLockScript() view returns(bytes)
func (_ZkBridge *ZkBridgeSession) ChangeLockScript() ([]byte, error) {
	return _ZkBridge.Contract.ChangeLockScript(&_ZkBridge.CallOpts)
}

// ChangeLockScript is a free data retrieval call binding the contract method 0x8e0a98af.
//
// Solidity: function changeLockScript() view returns(bytes)
func (_ZkBridge *ZkBridgeCallerSession) ChangeLockScript() ([]byte, error) {
	return _ZkBridge.Contract.ChangeLockScript(&_ZkBridge.CallOpts)
}

// CheckReceiveLockScript is a free data retrieval call binding the contract method 0x56af8603.
//
// Solidity: function checkReceiveLockScript(bytes receiveLockScript) view returns(bool)
func (_ZkBridge *ZkBridgeCaller) CheckReceiveLockScript(opts *bind.CallOpts, receiveLockScript []byte) (bool, error) {
	var out []interface{}
	err := _ZkBridge.contract.Call(opts, &out, "checkReceiveLockScript", receiveLockScript)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CheckReceiveLockScript is a free data retrieval call binding the contract method 0x56af8603.
//
// Solidity: function checkReceiveLockScript(bytes receiveLockScript) view returns(bool)
func (_ZkBridge *ZkBridgeSession) CheckReceiveLockScript(receiveLockScript []byte) (bool, error) {
	return _ZkBridge.Contract.CheckReceiveLockScript(&_ZkBridge.CallOpts, receiveLockScript)
}

// CheckReceiveLockScript is a free data retrieval call binding the contract method 0x56af8603.
//
// Solidity: function checkReceiveLockScript(bytes receiveLockScript) view returns(bool)
func (_ZkBridge *ZkBridgeCallerSession) CheckReceiveLockScript(receiveLockScript []byte) (bool, error) {
	return _ZkBridge.Contract.CheckReceiveLockScript(&_ZkBridge.CallOpts, receiveLockScript)
}

// EconomicVariation is a free data retrieval call binding the contract method 0xa78fab16.
//
// Solidity: function economicVariation() view returns(address)
func (_ZkBridge *ZkBridgeCaller) EconomicVariation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ZkBridge.contract.Call(opts, &out, "economicVariation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EconomicVariation is a free data retrieval call binding the contract method 0xa78fab16.
//
// Solidity: function economicVariation() view returns(address)
func (_ZkBridge *ZkBridgeSession) EconomicVariation() (common.Address, error) {
	return _ZkBridge.Contract.EconomicVariation(&_ZkBridge.CallOpts)
}

// EconomicVariation is a free data retrieval call binding the contract method 0xa78fab16.
//
// Solidity: function economicVariation() view returns(address)
func (_ZkBridge *ZkBridgeCallerSession) EconomicVariation() (common.Address, error) {
	return _ZkBridge.Contract.EconomicVariation(&_ZkBridge.CallOpts)
}

// EstimateTxWeight is a free data retrieval call binding the contract method 0x820e6fc1.
//
// Solidity: function estimateTxWeight(uint256 targetAmount) view returns(uint256)
func (_ZkBridge *ZkBridgeCaller) EstimateTxWeight(opts *bind.CallOpts, targetAmount *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _ZkBridge.contract.Call(opts, &out, "estimateTxWeight", targetAmount)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EstimateTxWeight is a free data retrieval call binding the contract method 0x820e6fc1.
//
// Solidity: function estimateTxWeight(uint256 targetAmount) view returns(uint256)
func (_ZkBridge *ZkBridgeSession) EstimateTxWeight(targetAmount *big.Int) (*big.Int, error) {
	return _ZkBridge.Contract.EstimateTxWeight(&_ZkBridge.CallOpts, targetAmount)
}

// EstimateTxWeight is a free data retrieval call binding the contract method 0x820e6fc1.
//
// Solidity: function estimateTxWeight(uint256 targetAmount) view returns(uint256)
func (_ZkBridge *ZkBridgeCallerSession) EstimateTxWeight(targetAmount *big.Int) (*big.Int, error) {
	return _ZkBridge.Contract.EstimateTxWeight(&_ZkBridge.CallOpts, targetAmount)
}

// FeeAccount is a free data retrieval call binding the contract method 0x65e17c9d.
//
// Solidity: function feeAccount() view returns(address)
func (_ZkBridge *ZkBridgeCaller) FeeAccount(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ZkBridge.contract.Call(opts, &out, "feeAccount")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FeeAccount is a free data retrieval call binding the contract method 0x65e17c9d.
//
// Solidity: function feeAccount() view returns(address)
func (_ZkBridge *ZkBridgeSession) FeeAccount() (common.Address, error) {
	return _ZkBridge.Contract.FeeAccount(&_ZkBridge.CallOpts)
}

// FeeAccount is a free data retrieval call binding the contract method 0x65e17c9d.
//
// Solidity: function feeAccount() view returns(address)
func (_ZkBridge *ZkBridgeCallerSession) FeeAccount() (common.Address, error) {
	return _ZkBridge.Contract.FeeAccount(&_ZkBridge.CallOpts)
}

// GetBridgeDepositToll is a free data retrieval call binding the contract method 0x3382d33e.
//
// Solidity: function getBridgeDepositToll(uint256 amount) view returns(uint256 userAmount, uint256 feeAmount)
func (_ZkBridge *ZkBridgeCaller) GetBridgeDepositToll(opts *bind.CallOpts, amount *big.Int) (struct {
	UserAmount *big.Int
	FeeAmount  *big.Int
}, error) {
	var out []interface{}
	err := _ZkBridge.contract.Call(opts, &out, "getBridgeDepositToll", amount)

	outstruct := new(struct {
		UserAmount *big.Int
		FeeAmount  *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.UserAmount = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.FeeAmount = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetBridgeDepositToll is a free data retrieval call binding the contract method 0x3382d33e.
//
// Solidity: function getBridgeDepositToll(uint256 amount) view returns(uint256 userAmount, uint256 feeAmount)
func (_ZkBridge *ZkBridgeSession) GetBridgeDepositToll(amount *big.Int) (struct {
	UserAmount *big.Int
	FeeAmount  *big.Int
}, error) {
	return _ZkBridge.Contract.GetBridgeDepositToll(&_ZkBridge.CallOpts, amount)
}

// GetBridgeDepositToll is a free data retrieval call binding the contract method 0x3382d33e.
//
// Solidity: function getBridgeDepositToll(uint256 amount) view returns(uint256 userAmount, uint256 feeAmount)
func (_ZkBridge *ZkBridgeCallerSession) GetBridgeDepositToll(amount *big.Int) (struct {
	UserAmount *big.Int
	FeeAmount  *big.Int
}, error) {
	return _ZkBridge.Contract.GetBridgeDepositToll(&_ZkBridge.CallOpts, amount)
}

// GetBridgeRedeemToll is a free data retrieval call binding the contract method 0x44272ac0.
//
// Solidity: function getBridgeRedeemToll(uint256 amount) view returns(uint256 userAmount, uint256 feeAmount)
func (_ZkBridge *ZkBridgeCaller) GetBridgeRedeemToll(opts *bind.CallOpts, amount *big.Int) (struct {
	UserAmount *big.Int
	FeeAmount  *big.Int
}, error) {
	var out []interface{}
	err := _ZkBridge.contract.Call(opts, &out, "getBridgeRedeemToll", amount)

	outstruct := new(struct {
		UserAmount *big.Int
		FeeAmount  *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.UserAmount = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.FeeAmount = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetBridgeRedeemToll is a free data retrieval call binding the contract method 0x44272ac0.
//
// Solidity: function getBridgeRedeemToll(uint256 amount) view returns(uint256 userAmount, uint256 feeAmount)
func (_ZkBridge *ZkBridgeSession) GetBridgeRedeemToll(amount *big.Int) (struct {
	UserAmount *big.Int
	FeeAmount  *big.Int
}, error) {
	return _ZkBridge.Contract.GetBridgeRedeemToll(&_ZkBridge.CallOpts, amount)
}

// GetBridgeRedeemToll is a free data retrieval call binding the contract method 0x44272ac0.
//
// Solidity: function getBridgeRedeemToll(uint256 amount) view returns(uint256 userAmount, uint256 feeAmount)
func (_ZkBridge *ZkBridgeCallerSession) GetBridgeRedeemToll(amount *big.Int) (struct {
	UserAmount *big.Int
	FeeAmount  *big.Int
}, error) {
	return _ZkBridge.Contract.GetBridgeRedeemToll(&_ZkBridge.CallOpts, amount)
}

// GetDepositLITMintAmount is a free data retrieval call binding the contract method 0x2acd03e9.
//
// Solidity: function getDepositLITMintAmount(uint256 feeAmount) view returns(uint256)
func (_ZkBridge *ZkBridgeCaller) GetDepositLITMintAmount(opts *bind.CallOpts, feeAmount *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _ZkBridge.contract.Call(opts, &out, "getDepositLITMintAmount", feeAmount)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetDepositLITMintAmount is a free data retrieval call binding the contract method 0x2acd03e9.
//
// Solidity: function getDepositLITMintAmount(uint256 feeAmount) view returns(uint256)
func (_ZkBridge *ZkBridgeSession) GetDepositLITMintAmount(feeAmount *big.Int) (*big.Int, error) {
	return _ZkBridge.Contract.GetDepositLITMintAmount(&_ZkBridge.CallOpts, feeAmount)
}

// GetDepositLITMintAmount is a free data retrieval call binding the contract method 0x2acd03e9.
//
// Solidity: function getDepositLITMintAmount(uint256 feeAmount) view returns(uint256)
func (_ZkBridge *ZkBridgeCallerSession) GetDepositLITMintAmount(feeAmount *big.Int) (*big.Int, error) {
	return _ZkBridge.Contract.GetDepositLITMintAmount(&_ZkBridge.CallOpts, feeAmount)
}

// GetRedeemLITMintAmount is a free data retrieval call binding the contract method 0x5750e398.
//
// Solidity: function getRedeemLITMintAmount(uint256 feeAmount) view returns(uint256)
func (_ZkBridge *ZkBridgeCaller) GetRedeemLITMintAmount(opts *bind.CallOpts, feeAmount *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _ZkBridge.contract.Call(opts, &out, "getRedeemLITMintAmount", feeAmount)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRedeemLITMintAmount is a free data retrieval call binding the contract method 0x5750e398.
//
// Solidity: function getRedeemLITMintAmount(uint256 feeAmount) view returns(uint256)
func (_ZkBridge *ZkBridgeSession) GetRedeemLITMintAmount(feeAmount *big.Int) (*big.Int, error) {
	return _ZkBridge.Contract.GetRedeemLITMintAmount(&_ZkBridge.CallOpts, feeAmount)
}

// GetRedeemLITMintAmount is a free data retrieval call binding the contract method 0x5750e398.
//
// Solidity: function getRedeemLITMintAmount(uint256 feeAmount) view returns(uint256)
func (_ZkBridge *ZkBridgeCallerSession) GetRedeemLITMintAmount(feeAmount *big.Int) (*big.Int, error) {
	return _ZkBridge.Contract.GetRedeemLITMintAmount(&_ZkBridge.CallOpts, feeAmount)
}

// GetTxOuts is a free data retrieval call binding the contract method 0x41a1ca42.
//
// Solidity: function getTxOuts(uint256 userAmount, uint256 changeAmount, bytes receiveLockScript) view returns((uint64,bytes)[])
func (_ZkBridge *ZkBridgeCaller) GetTxOuts(opts *bind.CallOpts, userAmount *big.Int, changeAmount *big.Int, receiveLockScript []byte) ([]UTXOManagerInterfaceTxOut, error) {
	var out []interface{}
	err := _ZkBridge.contract.Call(opts, &out, "getTxOuts", userAmount, changeAmount, receiveLockScript)

	if err != nil {
		return *new([]UTXOManagerInterfaceTxOut), err
	}

	out0 := *abi.ConvertType(out[0], new([]UTXOManagerInterfaceTxOut)).(*[]UTXOManagerInterfaceTxOut)

	return out0, err

}

// GetTxOuts is a free data retrieval call binding the contract method 0x41a1ca42.
//
// Solidity: function getTxOuts(uint256 userAmount, uint256 changeAmount, bytes receiveLockScript) view returns((uint64,bytes)[])
func (_ZkBridge *ZkBridgeSession) GetTxOuts(userAmount *big.Int, changeAmount *big.Int, receiveLockScript []byte) ([]UTXOManagerInterfaceTxOut, error) {
	return _ZkBridge.Contract.GetTxOuts(&_ZkBridge.CallOpts, userAmount, changeAmount, receiveLockScript)
}

// GetTxOuts is a free data retrieval call binding the contract method 0x41a1ca42.
//
// Solidity: function getTxOuts(uint256 userAmount, uint256 changeAmount, bytes receiveLockScript) view returns((uint64,bytes)[])
func (_ZkBridge *ZkBridgeCallerSession) GetTxOuts(userAmount *big.Int, changeAmount *big.Int, receiveLockScript []byte) ([]UTXOManagerInterfaceTxOut, error) {
	return _ZkBridge.Contract.GetTxOuts(&_ZkBridge.CallOpts, userAmount, changeAmount, receiveLockScript)
}

// IsOperator is a free data retrieval call binding the contract method 0x6d70f7ae.
//
// Solidity: function isOperator(address addr) view returns(bool)
func (_ZkBridge *ZkBridgeCaller) IsOperator(opts *bind.CallOpts, addr common.Address) (bool, error) {
	var out []interface{}
	err := _ZkBridge.contract.Call(opts, &out, "isOperator", addr)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOperator is a free data retrieval call binding the contract method 0x6d70f7ae.
//
// Solidity: function isOperator(address addr) view returns(bool)
func (_ZkBridge *ZkBridgeSession) IsOperator(addr common.Address) (bool, error) {
	return _ZkBridge.Contract.IsOperator(&_ZkBridge.CallOpts, addr)
}

// IsOperator is a free data retrieval call binding the contract method 0x6d70f7ae.
//
// Solidity: function isOperator(address addr) view returns(bool)
func (_ZkBridge *ZkBridgeCallerSession) IsOperator(addr common.Address) (bool, error) {
	return _ZkBridge.Contract.IsOperator(&_ZkBridge.CallOpts, addr)
}

// LitToken is a free data retrieval call binding the contract method 0x1c0c4564.
//
// Solidity: function litToken() view returns(address)
func (_ZkBridge *ZkBridgeCaller) LitToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ZkBridge.contract.Call(opts, &out, "litToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LitToken is a free data retrieval call binding the contract method 0x1c0c4564.
//
// Solidity: function litToken() view returns(address)
func (_ZkBridge *ZkBridgeSession) LitToken() (common.Address, error) {
	return _ZkBridge.Contract.LitToken(&_ZkBridge.CallOpts)
}

// LitToken is a free data retrieval call binding the contract method 0x1c0c4564.
//
// Solidity: function litToken() view returns(address)
func (_ZkBridge *ZkBridgeCallerSession) LitToken() (common.Address, error) {
	return _ZkBridge.Contract.LitToken(&_ZkBridge.CallOpts)
}

// MinDepositAmount is a free data retrieval call binding the contract method 0x645006ca.
//
// Solidity: function minDepositAmount() view returns(uint256)
func (_ZkBridge *ZkBridgeCaller) MinDepositAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ZkBridge.contract.Call(opts, &out, "minDepositAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinDepositAmount is a free data retrieval call binding the contract method 0x645006ca.
//
// Solidity: function minDepositAmount() view returns(uint256)
func (_ZkBridge *ZkBridgeSession) MinDepositAmount() (*big.Int, error) {
	return _ZkBridge.Contract.MinDepositAmount(&_ZkBridge.CallOpts)
}

// MinDepositAmount is a free data retrieval call binding the contract method 0x645006ca.
//
// Solidity: function minDepositAmount() view returns(uint256)
func (_ZkBridge *ZkBridgeCallerSession) MinDepositAmount() (*big.Int, error) {
	return _ZkBridge.Contract.MinDepositAmount(&_ZkBridge.CallOpts)
}

// MultiSigScript is a free data retrieval call binding the contract method 0x342561b8.
//
// Solidity: function multiSigScript() view returns(bytes)
func (_ZkBridge *ZkBridgeCaller) MultiSigScript(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _ZkBridge.contract.Call(opts, &out, "multiSigScript")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// MultiSigScript is a free data retrieval call binding the contract method 0x342561b8.
//
// Solidity: function multiSigScript() view returns(bytes)
func (_ZkBridge *ZkBridgeSession) MultiSigScript() ([]byte, error) {
	return _ZkBridge.Contract.MultiSigScript(&_ZkBridge.CallOpts)
}

// MultiSigScript is a free data retrieval call binding the contract method 0x342561b8.
//
// Solidity: function multiSigScript() view returns(bytes)
func (_ZkBridge *ZkBridgeCallerSession) MultiSigScript() ([]byte, error) {
	return _ZkBridge.Contract.MultiSigScript(&_ZkBridge.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ZkBridge *ZkBridgeCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ZkBridge.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ZkBridge *ZkBridgeSession) Owner() (common.Address, error) {
	return _ZkBridge.Contract.Owner(&_ZkBridge.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ZkBridge *ZkBridgeCallerSession) Owner() (common.Address, error) {
	return _ZkBridge.Contract.Owner(&_ZkBridge.CallOpts)
}

// TotalDepositAmount is a free data retrieval call binding the contract method 0xc5408d50.
//
// Solidity: function totalDepositAmount() view returns(uint256)
func (_ZkBridge *ZkBridgeCaller) TotalDepositAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ZkBridge.contract.Call(opts, &out, "totalDepositAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalDepositAmount is a free data retrieval call binding the contract method 0xc5408d50.
//
// Solidity: function totalDepositAmount() view returns(uint256)
func (_ZkBridge *ZkBridgeSession) TotalDepositAmount() (*big.Int, error) {
	return _ZkBridge.Contract.TotalDepositAmount(&_ZkBridge.CallOpts)
}

// TotalDepositAmount is a free data retrieval call binding the contract method 0xc5408d50.
//
// Solidity: function totalDepositAmount() view returns(uint256)
func (_ZkBridge *ZkBridgeCallerSession) TotalDepositAmount() (*big.Int, error) {
	return _ZkBridge.Contract.TotalDepositAmount(&_ZkBridge.CallOpts)
}

// UtxoManager is a free data retrieval call binding the contract method 0xab8db1b1.
//
// Solidity: function utxoManager() view returns(address)
func (_ZkBridge *ZkBridgeCaller) UtxoManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ZkBridge.contract.Call(opts, &out, "utxoManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// UtxoManager is a free data retrieval call binding the contract method 0xab8db1b1.
//
// Solidity: function utxoManager() view returns(address)
func (_ZkBridge *ZkBridgeSession) UtxoManager() (common.Address, error) {
	return _ZkBridge.Contract.UtxoManager(&_ZkBridge.CallOpts)
}

// UtxoManager is a free data retrieval call binding the contract method 0xab8db1b1.
//
// Solidity: function utxoManager() view returns(address)
func (_ZkBridge *ZkBridgeCallerSession) UtxoManager() (common.Address, error) {
	return _ZkBridge.Contract.UtxoManager(&_ZkBridge.CallOpts)
}

// VerifyChangeProof is a free data retrieval call binding the contract method 0xde4b7179.
//
// Solidity: function verifyChangeProof(bytes32 txid, bytes proofData) pure returns(bool)
func (_ZkBridge *ZkBridgeCaller) VerifyChangeProof(opts *bind.CallOpts, txid [32]byte, proofData []byte) (bool, error) {
	var out []interface{}
	err := _ZkBridge.contract.Call(opts, &out, "verifyChangeProof", txid, proofData)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VerifyChangeProof is a free data retrieval call binding the contract method 0xde4b7179.
//
// Solidity: function verifyChangeProof(bytes32 txid, bytes proofData) pure returns(bool)
func (_ZkBridge *ZkBridgeSession) VerifyChangeProof(txid [32]byte, proofData []byte) (bool, error) {
	return _ZkBridge.Contract.VerifyChangeProof(&_ZkBridge.CallOpts, txid, proofData)
}

// VerifyChangeProof is a free data retrieval call binding the contract method 0xde4b7179.
//
// Solidity: function verifyChangeProof(bytes32 txid, bytes proofData) pure returns(bool)
func (_ZkBridge *ZkBridgeCallerSession) VerifyChangeProof(txid [32]byte, proofData []byte) (bool, error) {
	return _ZkBridge.Contract.VerifyChangeProof(&_ZkBridge.CallOpts, txid, proofData)
}

// VerifyDepositProof is a free data retrieval call binding the contract method 0x0c8d9769.
//
// Solidity: function verifyDepositProof(bytes32 txid, uint32 index, uint256 amount, address account, bytes proofData) pure returns(bool)
func (_ZkBridge *ZkBridgeCaller) VerifyDepositProof(opts *bind.CallOpts, txid [32]byte, index uint32, amount *big.Int, account common.Address, proofData []byte) (bool, error) {
	var out []interface{}
	err := _ZkBridge.contract.Call(opts, &out, "verifyDepositProof", txid, index, amount, account, proofData)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VerifyDepositProof is a free data retrieval call binding the contract method 0x0c8d9769.
//
// Solidity: function verifyDepositProof(bytes32 txid, uint32 index, uint256 amount, address account, bytes proofData) pure returns(bool)
func (_ZkBridge *ZkBridgeSession) VerifyDepositProof(txid [32]byte, index uint32, amount *big.Int, account common.Address, proofData []byte) (bool, error) {
	return _ZkBridge.Contract.VerifyDepositProof(&_ZkBridge.CallOpts, txid, index, amount, account, proofData)
}

// VerifyDepositProof is a free data retrieval call binding the contract method 0x0c8d9769.
//
// Solidity: function verifyDepositProof(bytes32 txid, uint32 index, uint256 amount, address account, bytes proofData) pure returns(bool)
func (_ZkBridge *ZkBridgeCallerSession) VerifyDepositProof(txid [32]byte, index uint32, amount *big.Int, account common.Address, proofData []byte) (bool, error) {
	return _ZkBridge.Contract.VerifyDepositProof(&_ZkBridge.CallOpts, txid, index, amount, account, proofData)
}

// ZkBTCToken is a free data retrieval call binding the contract method 0x281904e4.
//
// Solidity: function zkBTCToken() view returns(address)
func (_ZkBridge *ZkBridgeCaller) ZkBTCToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ZkBridge.contract.Call(opts, &out, "zkBTCToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ZkBTCToken is a free data retrieval call binding the contract method 0x281904e4.
//
// Solidity: function zkBTCToken() view returns(address)
func (_ZkBridge *ZkBridgeSession) ZkBTCToken() (common.Address, error) {
	return _ZkBridge.Contract.ZkBTCToken(&_ZkBridge.CallOpts)
}

// ZkBTCToken is a free data retrieval call binding the contract method 0x281904e4.
//
// Solidity: function zkBTCToken() view returns(address)
func (_ZkBridge *ZkBridgeCallerSession) ZkBTCToken() (common.Address, error) {
	return _ZkBridge.Contract.ZkBTCToken(&_ZkBridge.CallOpts)
}

// AddOperator is a paid mutator transaction binding the contract method 0x9870d7fe.
//
// Solidity: function addOperator(address _new) returns()
func (_ZkBridge *ZkBridgeTransactor) AddOperator(opts *bind.TransactOpts, _new common.Address) (*types.Transaction, error) {
	return _ZkBridge.contract.Transact(opts, "addOperator", _new)
}

// AddOperator is a paid mutator transaction binding the contract method 0x9870d7fe.
//
// Solidity: function addOperator(address _new) returns()
func (_ZkBridge *ZkBridgeSession) AddOperator(_new common.Address) (*types.Transaction, error) {
	return _ZkBridge.Contract.AddOperator(&_ZkBridge.TransactOpts, _new)
}

// AddOperator is a paid mutator transaction binding the contract method 0x9870d7fe.
//
// Solidity: function addOperator(address _new) returns()
func (_ZkBridge *ZkBridgeTransactorSession) AddOperator(_new common.Address) (*types.Transaction, error) {
	return _ZkBridge.Contract.AddOperator(&_ZkBridge.TransactOpts, _new)
}

// Deposit is a paid mutator transaction binding the contract method 0x11f4a8b5.
//
// Solidity: function deposit(bytes32 txid, uint32 index, uint256 amount, address receiveAddress, bytes proofData) returns()
func (_ZkBridge *ZkBridgeTransactor) Deposit(opts *bind.TransactOpts, txid [32]byte, index uint32, amount *big.Int, receiveAddress common.Address, proofData []byte) (*types.Transaction, error) {
	return _ZkBridge.contract.Transact(opts, "deposit", txid, index, amount, receiveAddress, proofData)
}

// Deposit is a paid mutator transaction binding the contract method 0x11f4a8b5.
//
// Solidity: function deposit(bytes32 txid, uint32 index, uint256 amount, address receiveAddress, bytes proofData) returns()
func (_ZkBridge *ZkBridgeSession) Deposit(txid [32]byte, index uint32, amount *big.Int, receiveAddress common.Address, proofData []byte) (*types.Transaction, error) {
	return _ZkBridge.Contract.Deposit(&_ZkBridge.TransactOpts, txid, index, amount, receiveAddress, proofData)
}

// Deposit is a paid mutator transaction binding the contract method 0x11f4a8b5.
//
// Solidity: function deposit(bytes32 txid, uint32 index, uint256 amount, address receiveAddress, bytes proofData) returns()
func (_ZkBridge *ZkBridgeTransactorSession) Deposit(txid [32]byte, index uint32, amount *big.Int, receiveAddress common.Address, proofData []byte) (*types.Transaction, error) {
	return _ZkBridge.Contract.Deposit(&_ZkBridge.TransactOpts, txid, index, amount, receiveAddress, proofData)
}

// Redeem is a paid mutator transaction binding the contract method 0x0449015a.
//
// Solidity: function redeem(uint256 redeemAmount, uint256 btcMinerFee, bytes receiveLockScript) returns()
func (_ZkBridge *ZkBridgeTransactor) Redeem(opts *bind.TransactOpts, redeemAmount *big.Int, btcMinerFee *big.Int, receiveLockScript []byte) (*types.Transaction, error) {
	return _ZkBridge.contract.Transact(opts, "redeem", redeemAmount, btcMinerFee, receiveLockScript)
}

// Redeem is a paid mutator transaction binding the contract method 0x0449015a.
//
// Solidity: function redeem(uint256 redeemAmount, uint256 btcMinerFee, bytes receiveLockScript) returns()
func (_ZkBridge *ZkBridgeSession) Redeem(redeemAmount *big.Int, btcMinerFee *big.Int, receiveLockScript []byte) (*types.Transaction, error) {
	return _ZkBridge.Contract.Redeem(&_ZkBridge.TransactOpts, redeemAmount, btcMinerFee, receiveLockScript)
}

// Redeem is a paid mutator transaction binding the contract method 0x0449015a.
//
// Solidity: function redeem(uint256 redeemAmount, uint256 btcMinerFee, bytes receiveLockScript) returns()
func (_ZkBridge *ZkBridgeTransactorSession) Redeem(redeemAmount *big.Int, btcMinerFee *big.Int, receiveLockScript []byte) (*types.Transaction, error) {
	return _ZkBridge.Contract.Redeem(&_ZkBridge.TransactOpts, redeemAmount, btcMinerFee, receiveLockScript)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0xac8a584a.
//
// Solidity: function removeOperator(address _new) returns()
func (_ZkBridge *ZkBridgeTransactor) RemoveOperator(opts *bind.TransactOpts, _new common.Address) (*types.Transaction, error) {
	return _ZkBridge.contract.Transact(opts, "removeOperator", _new)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0xac8a584a.
//
// Solidity: function removeOperator(address _new) returns()
func (_ZkBridge *ZkBridgeSession) RemoveOperator(_new common.Address) (*types.Transaction, error) {
	return _ZkBridge.Contract.RemoveOperator(&_ZkBridge.TransactOpts, _new)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0xac8a584a.
//
// Solidity: function removeOperator(address _new) returns()
func (_ZkBridge *ZkBridgeTransactorSession) RemoveOperator(_new common.Address) (*types.Transaction, error) {
	return _ZkBridge.Contract.RemoveOperator(&_ZkBridge.TransactOpts, _new)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ZkBridge *ZkBridgeTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZkBridge.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ZkBridge *ZkBridgeSession) RenounceOwnership() (*types.Transaction, error) {
	return _ZkBridge.Contract.RenounceOwnership(&_ZkBridge.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ZkBridge *ZkBridgeTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ZkBridge.Contract.RenounceOwnership(&_ZkBridge.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ZkBridge *ZkBridgeTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ZkBridge.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ZkBridge *ZkBridgeSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ZkBridge.Contract.TransferOwnership(&_ZkBridge.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ZkBridge *ZkBridgeTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ZkBridge.Contract.TransferOwnership(&_ZkBridge.TransactOpts, newOwner)
}

// UpdateChange is a paid mutator transaction binding the contract method 0x8ce343ab.
//
// Solidity: function updateChange(bytes32 txid, bytes proofData) returns()
func (_ZkBridge *ZkBridgeTransactor) UpdateChange(opts *bind.TransactOpts, txid [32]byte, proofData []byte) (*types.Transaction, error) {
	return _ZkBridge.contract.Transact(opts, "updateChange", txid, proofData)
}

// UpdateChange is a paid mutator transaction binding the contract method 0x8ce343ab.
//
// Solidity: function updateChange(bytes32 txid, bytes proofData) returns()
func (_ZkBridge *ZkBridgeSession) UpdateChange(txid [32]byte, proofData []byte) (*types.Transaction, error) {
	return _ZkBridge.Contract.UpdateChange(&_ZkBridge.TransactOpts, txid, proofData)
}

// UpdateChange is a paid mutator transaction binding the contract method 0x8ce343ab.
//
// Solidity: function updateChange(bytes32 txid, bytes proofData) returns()
func (_ZkBridge *ZkBridgeTransactorSession) UpdateChange(txid [32]byte, proofData []byte) (*types.Transaction, error) {
	return _ZkBridge.Contract.UpdateChange(&_ZkBridge.TransactOpts, txid, proofData)
}

// UpdateEconomicAddress is a paid mutator transaction binding the contract method 0x5e5a769a.
//
// Solidity: function updateEconomicAddress(address _newAddress) returns()
func (_ZkBridge *ZkBridgeTransactor) UpdateEconomicAddress(opts *bind.TransactOpts, _newAddress common.Address) (*types.Transaction, error) {
	return _ZkBridge.contract.Transact(opts, "updateEconomicAddress", _newAddress)
}

// UpdateEconomicAddress is a paid mutator transaction binding the contract method 0x5e5a769a.
//
// Solidity: function updateEconomicAddress(address _newAddress) returns()
func (_ZkBridge *ZkBridgeSession) UpdateEconomicAddress(_newAddress common.Address) (*types.Transaction, error) {
	return _ZkBridge.Contract.UpdateEconomicAddress(&_ZkBridge.TransactOpts, _newAddress)
}

// UpdateEconomicAddress is a paid mutator transaction binding the contract method 0x5e5a769a.
//
// Solidity: function updateEconomicAddress(address _newAddress) returns()
func (_ZkBridge *ZkBridgeTransactorSession) UpdateEconomicAddress(_newAddress common.Address) (*types.Transaction, error) {
	return _ZkBridge.Contract.UpdateEconomicAddress(&_ZkBridge.TransactOpts, _newAddress)
}

// UpdateUtxoAddress is a paid mutator transaction binding the contract method 0x9314c393.
//
// Solidity: function updateUtxoAddress(address _utxoAddress, bytes _changeLockScript, bytes _multiSigScript) returns()
func (_ZkBridge *ZkBridgeTransactor) UpdateUtxoAddress(opts *bind.TransactOpts, _utxoAddress common.Address, _changeLockScript []byte, _multiSigScript []byte) (*types.Transaction, error) {
	return _ZkBridge.contract.Transact(opts, "updateUtxoAddress", _utxoAddress, _changeLockScript, _multiSigScript)
}

// UpdateUtxoAddress is a paid mutator transaction binding the contract method 0x9314c393.
//
// Solidity: function updateUtxoAddress(address _utxoAddress, bytes _changeLockScript, bytes _multiSigScript) returns()
func (_ZkBridge *ZkBridgeSession) UpdateUtxoAddress(_utxoAddress common.Address, _changeLockScript []byte, _multiSigScript []byte) (*types.Transaction, error) {
	return _ZkBridge.Contract.UpdateUtxoAddress(&_ZkBridge.TransactOpts, _utxoAddress, _changeLockScript, _multiSigScript)
}

// UpdateUtxoAddress is a paid mutator transaction binding the contract method 0x9314c393.
//
// Solidity: function updateUtxoAddress(address _utxoAddress, bytes _changeLockScript, bytes _multiSigScript) returns()
func (_ZkBridge *ZkBridgeTransactorSession) UpdateUtxoAddress(_utxoAddress common.Address, _changeLockScript []byte, _multiSigScript []byte) (*types.Transaction, error) {
	return _ZkBridge.Contract.UpdateUtxoAddress(&_ZkBridge.TransactOpts, _utxoAddress, _changeLockScript, _multiSigScript)
}

// ZkBridgeCreateRedeemUnsignedTxIterator is returned from FilterCreateRedeemUnsignedTx and is used to iterate over the raw logs and unpacked data for CreateRedeemUnsignedTx events raised by the ZkBridge contract.
type ZkBridgeCreateRedeemUnsignedTxIterator struct {
	Event *ZkBridgeCreateRedeemUnsignedTx // Event containing the contract specifics and raw log

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
func (it *ZkBridgeCreateRedeemUnsignedTxIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZkBridgeCreateRedeemUnsignedTx)
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
		it.Event = new(ZkBridgeCreateRedeemUnsignedTx)
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
func (it *ZkBridgeCreateRedeemUnsignedTxIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZkBridgeCreateRedeemUnsignedTxIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZkBridgeCreateRedeemUnsignedTx represents a CreateRedeemUnsignedTx event raised by the ZkBridge contract.
type ZkBridgeCreateRedeemUnsignedTx struct {
	Txid     [32]byte
	RawTx    []byte
	SigHashs [][32]byte
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterCreateRedeemUnsignedTx is a free log retrieval operation binding the contract event 0x1e5e2baa6d11cc5bcae8c0d1187d7b9ebf13d6d9b932f7dbbf4e396438845fb8.
//
// Solidity: event CreateRedeemUnsignedTx(bytes32 indexed txid, bytes rawTx, bytes32[] sigHashs)
func (_ZkBridge *ZkBridgeFilterer) FilterCreateRedeemUnsignedTx(opts *bind.FilterOpts, txid [][32]byte) (*ZkBridgeCreateRedeemUnsignedTxIterator, error) {

	var txidRule []interface{}
	for _, txidItem := range txid {
		txidRule = append(txidRule, txidItem)
	}

	logs, sub, err := _ZkBridge.contract.FilterLogs(opts, "CreateRedeemUnsignedTx", txidRule)
	if err != nil {
		return nil, err
	}
	return &ZkBridgeCreateRedeemUnsignedTxIterator{contract: _ZkBridge.contract, event: "CreateRedeemUnsignedTx", logs: logs, sub: sub}, nil
}

// WatchCreateRedeemUnsignedTx is a free log subscription operation binding the contract event 0x1e5e2baa6d11cc5bcae8c0d1187d7b9ebf13d6d9b932f7dbbf4e396438845fb8.
//
// Solidity: event CreateRedeemUnsignedTx(bytes32 indexed txid, bytes rawTx, bytes32[] sigHashs)
func (_ZkBridge *ZkBridgeFilterer) WatchCreateRedeemUnsignedTx(opts *bind.WatchOpts, sink chan<- *ZkBridgeCreateRedeemUnsignedTx, txid [][32]byte) (event.Subscription, error) {

	var txidRule []interface{}
	for _, txidItem := range txid {
		txidRule = append(txidRule, txidItem)
	}

	logs, sub, err := _ZkBridge.contract.WatchLogs(opts, "CreateRedeemUnsignedTx", txidRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZkBridgeCreateRedeemUnsignedTx)
				if err := _ZkBridge.contract.UnpackLog(event, "CreateRedeemUnsignedTx", log); err != nil {
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

// ParseCreateRedeemUnsignedTx is a log parse operation binding the contract event 0x1e5e2baa6d11cc5bcae8c0d1187d7b9ebf13d6d9b932f7dbbf4e396438845fb8.
//
// Solidity: event CreateRedeemUnsignedTx(bytes32 indexed txid, bytes rawTx, bytes32[] sigHashs)
func (_ZkBridge *ZkBridgeFilterer) ParseCreateRedeemUnsignedTx(log types.Log) (*ZkBridgeCreateRedeemUnsignedTx, error) {
	event := new(ZkBridgeCreateRedeemUnsignedTx)
	if err := _ZkBridge.contract.UnpackLog(event, "CreateRedeemUnsignedTx", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZkBridgeOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ZkBridge contract.
type ZkBridgeOwnershipTransferredIterator struct {
	Event *ZkBridgeOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ZkBridgeOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZkBridgeOwnershipTransferred)
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
		it.Event = new(ZkBridgeOwnershipTransferred)
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
func (it *ZkBridgeOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZkBridgeOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZkBridgeOwnershipTransferred represents a OwnershipTransferred event raised by the ZkBridge contract.
type ZkBridgeOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ZkBridge *ZkBridgeFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ZkBridgeOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ZkBridge.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ZkBridgeOwnershipTransferredIterator{contract: _ZkBridge.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ZkBridge *ZkBridgeFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ZkBridgeOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ZkBridge.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZkBridgeOwnershipTransferred)
				if err := _ZkBridge.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_ZkBridge *ZkBridgeFilterer) ParseOwnershipTransferred(log types.Log) (*ZkBridgeOwnershipTransferred, error) {
	event := new(ZkBridgeOwnershipTransferred)
	if err := _ZkBridge.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
