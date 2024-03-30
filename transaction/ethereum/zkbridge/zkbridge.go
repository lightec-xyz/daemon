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

// UTXOManagerInterfaceTxOut is an auto generated low-level Go binding around an user-defined struct.
type UTXOManagerInterfaceTxOut struct {
	Amount       uint64
	ScriptPubKey []byte
}

// ZkbridgeMetaData contains all meta data concerning the Zkbridge contract.
var ZkbridgeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_feeAccount\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_changeLockScript\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"_multiSigScript\",\"type\":\"bytes\"},{\"internalType\":\"contractUTXOManagerInterface\",\"name\":\"_utxoAddress\",\"type\":\"address\"},{\"internalType\":\"contractLITInterface\",\"name\":\"_litAddress\",\"type\":\"address\"},{\"internalType\":\"contractzkBTCInterface\",\"name\":\"_zkBTCAddress\",\"type\":\"address\"},{\"internalType\":\"contractEconomicVariationInterface\",\"name\":\"_variationAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"DepositAccountIsBridge\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"proofData\",\"type\":\"bytes\"}],\"name\":\"InvalidChangeProof\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"depositAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minDepositAmount\",\"type\":\"uint256\"}],\"name\":\"InvalidDepositAmount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"changeLockScript\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"proofData\",\"type\":\"bytes\"}],\"name\":\"InvalidDepositProof\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"lockScriptLength\",\"type\":\"uint256\"}],\"name\":\"InvalidLockScriptLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"lockScript\",\"type\":\"bytes\"}],\"name\":\"LockScriptIsChange\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"rawTx\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes32[]\",\"name\":\"sigHashs\",\"type\":\"bytes32[]\"}],\"name\":\"CreateRedeemUnsignedTx\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_new\",\"type\":\"address\"}],\"name\":\"addOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"changeLockScript\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"economicVariation\",\"outputs\":[{\"internalType\":\"contractEconomicVariationInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feeAccount\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isOperator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"litToken\",\"outputs\":[{\"internalType\":\"contractLITInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minDepositAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"multiSigScript\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_new\",\"type\":\"address\"}],\"name\":\"removeOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalDepositAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"utxoManager\",\"outputs\":[{\"internalType\":\"contractUTXOManagerInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"zkBTCToken\",\"outputs\":[{\"internalType\":\"contractzkBTCInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractEconomicVariationInterface\",\"name\":\"_newAddress\",\"type\":\"address\"}],\"name\":\"updateEconomicAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractUTXOManagerInterface\",\"name\":\"_utxoAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_changeLockScript\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"_multiSigScript\",\"type\":\"bytes\"}],\"name\":\"updateUtxoAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"receiveAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"proofData\",\"type\":\"bytes\"}],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"targetAmount\",\"type\":\"uint256\"}],\"name\":\"estimateTxWeight\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"receiveLockScript\",\"type\":\"bytes\"}],\"name\":\"checkReceiveLockScript\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"redeemAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"btcMinerFee\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"receiveLockScript\",\"type\":\"bytes\"}],\"name\":\"redeem\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"userAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"changeAmount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"receiveLockScript\",\"type\":\"bytes\"}],\"name\":\"getTxOuts\",\"outputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"amount\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"scriptPubKey\",\"type\":\"bytes\"}],\"internalType\":\"structUTXOManagerInterface.TxOut[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"proofData\",\"type\":\"bytes\"}],\"name\":\"updateChange\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"proofData\",\"type\":\"bytes\"}],\"name\":\"verifyDepositProof\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"proofData\",\"type\":\"bytes\"}],\"name\":\"verifyChangeProof\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"getBridgeDepositToll\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"userAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"feeAmount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"getBridgeRedeemToll\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"userAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"feeAmount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"feeAmount\",\"type\":\"uint256\"}],\"name\":\"getDepositLITMintAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"feeAmount\",\"type\":\"uint256\"}],\"name\":\"getRedeemLITMintAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// ZkbridgeABI is the input ABI used to generate the binding from.
// Deprecated: Use ZkbridgeMetaData.ABI instead.
var ZkbridgeABI = ZkbridgeMetaData.ABI

// Zkbridge is an auto generated Go binding around an Ethereum contract.
type Zkbridge struct {
	ZkbridgeCaller     // Read-only binding to the contract
	ZkbridgeTransactor // Write-only binding to the contract
	ZkbridgeFilterer   // Log filterer for contract events
}

// ZkbridgeCaller is an auto generated read-only Go binding around an Ethereum contract.
type ZkbridgeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZkbridgeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ZkbridgeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZkbridgeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ZkbridgeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZkbridgeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ZkbridgeSession struct {
	Contract     *Zkbridge         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ZkbridgeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ZkbridgeCallerSession struct {
	Contract *ZkbridgeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // call options to use throughout this session
}

// ZkbridgeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ZkbridgeTransactorSession struct {
	Contract     *ZkbridgeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ZkbridgeRaw is an auto generated low-level Go binding around an Ethereum contract.
type ZkbridgeRaw struct {
	Contract *Zkbridge // Generic contract binding to access the raw methods on
}

// ZkbridgeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ZkbridgeCallerRaw struct {
	Contract *ZkbridgeCaller // Generic read-only contract binding to access the raw methods on
}

// ZkbridgeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ZkbridgeTransactorRaw struct {
	Contract *ZkbridgeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewZkbridge creates a new instance of Zkbridge, bound to a specific deployed contract.
func NewZkbridge(address common.Address, backend bind.ContractBackend) (*Zkbridge, error) {
	contract, err := bindZkbridge(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Zkbridge{ZkbridgeCaller: ZkbridgeCaller{contract: contract}, ZkbridgeTransactor: ZkbridgeTransactor{contract: contract}, ZkbridgeFilterer: ZkbridgeFilterer{contract: contract}}, nil
}

// NewZkbridgeCaller creates a new read-only instance of Zkbridge, bound to a specific deployed contract.
func NewZkbridgeCaller(address common.Address, caller bind.ContractCaller) (*ZkbridgeCaller, error) {
	contract, err := bindZkbridge(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ZkbridgeCaller{contract: contract}, nil
}

// NewZkbridgeTransactor creates a new write-only instance of Zkbridge, bound to a specific deployed contract.
func NewZkbridgeTransactor(address common.Address, transactor bind.ContractTransactor) (*ZkbridgeTransactor, error) {
	contract, err := bindZkbridge(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ZkbridgeTransactor{contract: contract}, nil
}

// NewZkbridgeFilterer creates a new log filterer instance of Zkbridge, bound to a specific deployed contract.
func NewZkbridgeFilterer(address common.Address, filterer bind.ContractFilterer) (*ZkbridgeFilterer, error) {
	contract, err := bindZkbridge(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ZkbridgeFilterer{contract: contract}, nil
}

// bindZkbridge binds a generic wrapper to an already deployed contract.
func bindZkbridge(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ZkbridgeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Zkbridge *ZkbridgeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Zkbridge.Contract.ZkbridgeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Zkbridge *ZkbridgeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Zkbridge.Contract.ZkbridgeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Zkbridge *ZkbridgeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Zkbridge.Contract.ZkbridgeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Zkbridge *ZkbridgeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Zkbridge.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Zkbridge *ZkbridgeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Zkbridge.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Zkbridge *ZkbridgeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Zkbridge.Contract.contract.Transact(opts, method, params...)
}

// ChangeLockScript is a free data retrieval call binding the contract method 0x8e0a98af.
//
// Solidity: function changeLockScript() view returns(bytes)
func (_Zkbridge *ZkbridgeCaller) ChangeLockScript(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _Zkbridge.contract.Call(opts, &out, "changeLockScript")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// ChangeLockScript is a free data retrieval call binding the contract method 0x8e0a98af.
//
// Solidity: function changeLockScript() view returns(bytes)
func (_Zkbridge *ZkbridgeSession) ChangeLockScript() ([]byte, error) {
	return _Zkbridge.Contract.ChangeLockScript(&_Zkbridge.CallOpts)
}

// ChangeLockScript is a free data retrieval call binding the contract method 0x8e0a98af.
//
// Solidity: function changeLockScript() view returns(bytes)
func (_Zkbridge *ZkbridgeCallerSession) ChangeLockScript() ([]byte, error) {
	return _Zkbridge.Contract.ChangeLockScript(&_Zkbridge.CallOpts)
}

// CheckReceiveLockScript is a free data retrieval call binding the contract method 0x56af8603.
//
// Solidity: function checkReceiveLockScript(bytes receiveLockScript) view returns(bool)
func (_Zkbridge *ZkbridgeCaller) CheckReceiveLockScript(opts *bind.CallOpts, receiveLockScript []byte) (bool, error) {
	var out []interface{}
	err := _Zkbridge.contract.Call(opts, &out, "checkReceiveLockScript", receiveLockScript)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CheckReceiveLockScript is a free data retrieval call binding the contract method 0x56af8603.
//
// Solidity: function checkReceiveLockScript(bytes receiveLockScript) view returns(bool)
func (_Zkbridge *ZkbridgeSession) CheckReceiveLockScript(receiveLockScript []byte) (bool, error) {
	return _Zkbridge.Contract.CheckReceiveLockScript(&_Zkbridge.CallOpts, receiveLockScript)
}

// CheckReceiveLockScript is a free data retrieval call binding the contract method 0x56af8603.
//
// Solidity: function checkReceiveLockScript(bytes receiveLockScript) view returns(bool)
func (_Zkbridge *ZkbridgeCallerSession) CheckReceiveLockScript(receiveLockScript []byte) (bool, error) {
	return _Zkbridge.Contract.CheckReceiveLockScript(&_Zkbridge.CallOpts, receiveLockScript)
}

// EconomicVariation is a free data retrieval call binding the contract method 0xa78fab16.
//
// Solidity: function economicVariation() view returns(address)
func (_Zkbridge *ZkbridgeCaller) EconomicVariation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Zkbridge.contract.Call(opts, &out, "economicVariation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EconomicVariation is a free data retrieval call binding the contract method 0xa78fab16.
//
// Solidity: function economicVariation() view returns(address)
func (_Zkbridge *ZkbridgeSession) EconomicVariation() (common.Address, error) {
	return _Zkbridge.Contract.EconomicVariation(&_Zkbridge.CallOpts)
}

// EconomicVariation is a free data retrieval call binding the contract method 0xa78fab16.
//
// Solidity: function economicVariation() view returns(address)
func (_Zkbridge *ZkbridgeCallerSession) EconomicVariation() (common.Address, error) {
	return _Zkbridge.Contract.EconomicVariation(&_Zkbridge.CallOpts)
}

// EstimateTxWeight is a free data retrieval call binding the contract method 0x820e6fc1.
//
// Solidity: function estimateTxWeight(uint256 targetAmount) view returns(uint256)
func (_Zkbridge *ZkbridgeCaller) EstimateTxWeight(opts *bind.CallOpts, targetAmount *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Zkbridge.contract.Call(opts, &out, "estimateTxWeight", targetAmount)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EstimateTxWeight is a free data retrieval call binding the contract method 0x820e6fc1.
//
// Solidity: function estimateTxWeight(uint256 targetAmount) view returns(uint256)
func (_Zkbridge *ZkbridgeSession) EstimateTxWeight(targetAmount *big.Int) (*big.Int, error) {
	return _Zkbridge.Contract.EstimateTxWeight(&_Zkbridge.CallOpts, targetAmount)
}

// EstimateTxWeight is a free data retrieval call binding the contract method 0x820e6fc1.
//
// Solidity: function estimateTxWeight(uint256 targetAmount) view returns(uint256)
func (_Zkbridge *ZkbridgeCallerSession) EstimateTxWeight(targetAmount *big.Int) (*big.Int, error) {
	return _Zkbridge.Contract.EstimateTxWeight(&_Zkbridge.CallOpts, targetAmount)
}

// FeeAccount is a free data retrieval call binding the contract method 0x65e17c9d.
//
// Solidity: function feeAccount() view returns(address)
func (_Zkbridge *ZkbridgeCaller) FeeAccount(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Zkbridge.contract.Call(opts, &out, "feeAccount")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FeeAccount is a free data retrieval call binding the contract method 0x65e17c9d.
//
// Solidity: function feeAccount() view returns(address)
func (_Zkbridge *ZkbridgeSession) FeeAccount() (common.Address, error) {
	return _Zkbridge.Contract.FeeAccount(&_Zkbridge.CallOpts)
}

// FeeAccount is a free data retrieval call binding the contract method 0x65e17c9d.
//
// Solidity: function feeAccount() view returns(address)
func (_Zkbridge *ZkbridgeCallerSession) FeeAccount() (common.Address, error) {
	return _Zkbridge.Contract.FeeAccount(&_Zkbridge.CallOpts)
}

// GetBridgeDepositToll is a free data retrieval call binding the contract method 0x3382d33e.
//
// Solidity: function getBridgeDepositToll(uint256 amount) view returns(uint256 userAmount, uint256 feeAmount)
func (_Zkbridge *ZkbridgeCaller) GetBridgeDepositToll(opts *bind.CallOpts, amount *big.Int) (struct {
	UserAmount *big.Int
	FeeAmount  *big.Int
}, error) {
	var out []interface{}
	err := _Zkbridge.contract.Call(opts, &out, "getBridgeDepositToll", amount)

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
func (_Zkbridge *ZkbridgeSession) GetBridgeDepositToll(amount *big.Int) (struct {
	UserAmount *big.Int
	FeeAmount  *big.Int
}, error) {
	return _Zkbridge.Contract.GetBridgeDepositToll(&_Zkbridge.CallOpts, amount)
}

// GetBridgeDepositToll is a free data retrieval call binding the contract method 0x3382d33e.
//
// Solidity: function getBridgeDepositToll(uint256 amount) view returns(uint256 userAmount, uint256 feeAmount)
func (_Zkbridge *ZkbridgeCallerSession) GetBridgeDepositToll(amount *big.Int) (struct {
	UserAmount *big.Int
	FeeAmount  *big.Int
}, error) {
	return _Zkbridge.Contract.GetBridgeDepositToll(&_Zkbridge.CallOpts, amount)
}

// GetBridgeRedeemToll is a free data retrieval call binding the contract method 0x44272ac0.
//
// Solidity: function getBridgeRedeemToll(uint256 amount) view returns(uint256 userAmount, uint256 feeAmount)
func (_Zkbridge *ZkbridgeCaller) GetBridgeRedeemToll(opts *bind.CallOpts, amount *big.Int) (struct {
	UserAmount *big.Int
	FeeAmount  *big.Int
}, error) {
	var out []interface{}
	err := _Zkbridge.contract.Call(opts, &out, "getBridgeRedeemToll", amount)

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
func (_Zkbridge *ZkbridgeSession) GetBridgeRedeemToll(amount *big.Int) (struct {
	UserAmount *big.Int
	FeeAmount  *big.Int
}, error) {
	return _Zkbridge.Contract.GetBridgeRedeemToll(&_Zkbridge.CallOpts, amount)
}

// GetBridgeRedeemToll is a free data retrieval call binding the contract method 0x44272ac0.
//
// Solidity: function getBridgeRedeemToll(uint256 amount) view returns(uint256 userAmount, uint256 feeAmount)
func (_Zkbridge *ZkbridgeCallerSession) GetBridgeRedeemToll(amount *big.Int) (struct {
	UserAmount *big.Int
	FeeAmount  *big.Int
}, error) {
	return _Zkbridge.Contract.GetBridgeRedeemToll(&_Zkbridge.CallOpts, amount)
}

// GetDepositLITMintAmount is a free data retrieval call binding the contract method 0x2acd03e9.
//
// Solidity: function getDepositLITMintAmount(uint256 feeAmount) view returns(uint256)
func (_Zkbridge *ZkbridgeCaller) GetDepositLITMintAmount(opts *bind.CallOpts, feeAmount *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Zkbridge.contract.Call(opts, &out, "getDepositLITMintAmount", feeAmount)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetDepositLITMintAmount is a free data retrieval call binding the contract method 0x2acd03e9.
//
// Solidity: function getDepositLITMintAmount(uint256 feeAmount) view returns(uint256)
func (_Zkbridge *ZkbridgeSession) GetDepositLITMintAmount(feeAmount *big.Int) (*big.Int, error) {
	return _Zkbridge.Contract.GetDepositLITMintAmount(&_Zkbridge.CallOpts, feeAmount)
}

// GetDepositLITMintAmount is a free data retrieval call binding the contract method 0x2acd03e9.
//
// Solidity: function getDepositLITMintAmount(uint256 feeAmount) view returns(uint256)
func (_Zkbridge *ZkbridgeCallerSession) GetDepositLITMintAmount(feeAmount *big.Int) (*big.Int, error) {
	return _Zkbridge.Contract.GetDepositLITMintAmount(&_Zkbridge.CallOpts, feeAmount)
}

// GetRedeemLITMintAmount is a free data retrieval call binding the contract method 0x5750e398.
//
// Solidity: function getRedeemLITMintAmount(uint256 feeAmount) view returns(uint256)
func (_Zkbridge *ZkbridgeCaller) GetRedeemLITMintAmount(opts *bind.CallOpts, feeAmount *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Zkbridge.contract.Call(opts, &out, "getRedeemLITMintAmount", feeAmount)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRedeemLITMintAmount is a free data retrieval call binding the contract method 0x5750e398.
//
// Solidity: function getRedeemLITMintAmount(uint256 feeAmount) view returns(uint256)
func (_Zkbridge *ZkbridgeSession) GetRedeemLITMintAmount(feeAmount *big.Int) (*big.Int, error) {
	return _Zkbridge.Contract.GetRedeemLITMintAmount(&_Zkbridge.CallOpts, feeAmount)
}

// GetRedeemLITMintAmount is a free data retrieval call binding the contract method 0x5750e398.
//
// Solidity: function getRedeemLITMintAmount(uint256 feeAmount) view returns(uint256)
func (_Zkbridge *ZkbridgeCallerSession) GetRedeemLITMintAmount(feeAmount *big.Int) (*big.Int, error) {
	return _Zkbridge.Contract.GetRedeemLITMintAmount(&_Zkbridge.CallOpts, feeAmount)
}

// GetTxOuts is a free data retrieval call binding the contract method 0x41a1ca42.
//
// Solidity: function getTxOuts(uint256 userAmount, uint256 changeAmount, bytes receiveLockScript) view returns((uint64,bytes)[])
func (_Zkbridge *ZkbridgeCaller) GetTxOuts(opts *bind.CallOpts, userAmount *big.Int, changeAmount *big.Int, receiveLockScript []byte) ([]UTXOManagerInterfaceTxOut, error) {
	var out []interface{}
	err := _Zkbridge.contract.Call(opts, &out, "getTxOuts", userAmount, changeAmount, receiveLockScript)

	if err != nil {
		return *new([]UTXOManagerInterfaceTxOut), err
	}

	out0 := *abi.ConvertType(out[0], new([]UTXOManagerInterfaceTxOut)).(*[]UTXOManagerInterfaceTxOut)

	return out0, err

}

// GetTxOuts is a free data retrieval call binding the contract method 0x41a1ca42.
//
// Solidity: function getTxOuts(uint256 userAmount, uint256 changeAmount, bytes receiveLockScript) view returns((uint64,bytes)[])
func (_Zkbridge *ZkbridgeSession) GetTxOuts(userAmount *big.Int, changeAmount *big.Int, receiveLockScript []byte) ([]UTXOManagerInterfaceTxOut, error) {
	return _Zkbridge.Contract.GetTxOuts(&_Zkbridge.CallOpts, userAmount, changeAmount, receiveLockScript)
}

// GetTxOuts is a free data retrieval call binding the contract method 0x41a1ca42.
//
// Solidity: function getTxOuts(uint256 userAmount, uint256 changeAmount, bytes receiveLockScript) view returns((uint64,bytes)[])
func (_Zkbridge *ZkbridgeCallerSession) GetTxOuts(userAmount *big.Int, changeAmount *big.Int, receiveLockScript []byte) ([]UTXOManagerInterfaceTxOut, error) {
	return _Zkbridge.Contract.GetTxOuts(&_Zkbridge.CallOpts, userAmount, changeAmount, receiveLockScript)
}

// IsOperator is a free data retrieval call binding the contract method 0x6d70f7ae.
//
// Solidity: function isOperator(address addr) view returns(bool)
func (_Zkbridge *ZkbridgeCaller) IsOperator(opts *bind.CallOpts, addr common.Address) (bool, error) {
	var out []interface{}
	err := _Zkbridge.contract.Call(opts, &out, "isOperator", addr)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOperator is a free data retrieval call binding the contract method 0x6d70f7ae.
//
// Solidity: function isOperator(address addr) view returns(bool)
func (_Zkbridge *ZkbridgeSession) IsOperator(addr common.Address) (bool, error) {
	return _Zkbridge.Contract.IsOperator(&_Zkbridge.CallOpts, addr)
}

// IsOperator is a free data retrieval call binding the contract method 0x6d70f7ae.
//
// Solidity: function isOperator(address addr) view returns(bool)
func (_Zkbridge *ZkbridgeCallerSession) IsOperator(addr common.Address) (bool, error) {
	return _Zkbridge.Contract.IsOperator(&_Zkbridge.CallOpts, addr)
}

// LitToken is a free data retrieval call binding the contract method 0x1c0c4564.
//
// Solidity: function litToken() view returns(address)
func (_Zkbridge *ZkbridgeCaller) LitToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Zkbridge.contract.Call(opts, &out, "litToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LitToken is a free data retrieval call binding the contract method 0x1c0c4564.
//
// Solidity: function litToken() view returns(address)
func (_Zkbridge *ZkbridgeSession) LitToken() (common.Address, error) {
	return _Zkbridge.Contract.LitToken(&_Zkbridge.CallOpts)
}

// LitToken is a free data retrieval call binding the contract method 0x1c0c4564.
//
// Solidity: function litToken() view returns(address)
func (_Zkbridge *ZkbridgeCallerSession) LitToken() (common.Address, error) {
	return _Zkbridge.Contract.LitToken(&_Zkbridge.CallOpts)
}

// MinDepositAmount is a free data retrieval call binding the contract method 0x645006ca.
//
// Solidity: function minDepositAmount() view returns(uint256)
func (_Zkbridge *ZkbridgeCaller) MinDepositAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Zkbridge.contract.Call(opts, &out, "minDepositAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinDepositAmount is a free data retrieval call binding the contract method 0x645006ca.
//
// Solidity: function minDepositAmount() view returns(uint256)
func (_Zkbridge *ZkbridgeSession) MinDepositAmount() (*big.Int, error) {
	return _Zkbridge.Contract.MinDepositAmount(&_Zkbridge.CallOpts)
}

// MinDepositAmount is a free data retrieval call binding the contract method 0x645006ca.
//
// Solidity: function minDepositAmount() view returns(uint256)
func (_Zkbridge *ZkbridgeCallerSession) MinDepositAmount() (*big.Int, error) {
	return _Zkbridge.Contract.MinDepositAmount(&_Zkbridge.CallOpts)
}

// MultiSigScript is a free data retrieval call binding the contract method 0x342561b8.
//
// Solidity: function multiSigScript() view returns(bytes)
func (_Zkbridge *ZkbridgeCaller) MultiSigScript(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _Zkbridge.contract.Call(opts, &out, "multiSigScript")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// MultiSigScript is a free data retrieval call binding the contract method 0x342561b8.
//
// Solidity: function multiSigScript() view returns(bytes)
func (_Zkbridge *ZkbridgeSession) MultiSigScript() ([]byte, error) {
	return _Zkbridge.Contract.MultiSigScript(&_Zkbridge.CallOpts)
}

// MultiSigScript is a free data retrieval call binding the contract method 0x342561b8.
//
// Solidity: function multiSigScript() view returns(bytes)
func (_Zkbridge *ZkbridgeCallerSession) MultiSigScript() ([]byte, error) {
	return _Zkbridge.Contract.MultiSigScript(&_Zkbridge.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Zkbridge *ZkbridgeCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Zkbridge.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Zkbridge *ZkbridgeSession) Owner() (common.Address, error) {
	return _Zkbridge.Contract.Owner(&_Zkbridge.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Zkbridge *ZkbridgeCallerSession) Owner() (common.Address, error) {
	return _Zkbridge.Contract.Owner(&_Zkbridge.CallOpts)
}

// TotalDepositAmount is a free data retrieval call binding the contract method 0xc5408d50.
//
// Solidity: function totalDepositAmount() view returns(uint256)
func (_Zkbridge *ZkbridgeCaller) TotalDepositAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Zkbridge.contract.Call(opts, &out, "totalDepositAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalDepositAmount is a free data retrieval call binding the contract method 0xc5408d50.
//
// Solidity: function totalDepositAmount() view returns(uint256)
func (_Zkbridge *ZkbridgeSession) TotalDepositAmount() (*big.Int, error) {
	return _Zkbridge.Contract.TotalDepositAmount(&_Zkbridge.CallOpts)
}

// TotalDepositAmount is a free data retrieval call binding the contract method 0xc5408d50.
//
// Solidity: function totalDepositAmount() view returns(uint256)
func (_Zkbridge *ZkbridgeCallerSession) TotalDepositAmount() (*big.Int, error) {
	return _Zkbridge.Contract.TotalDepositAmount(&_Zkbridge.CallOpts)
}

// UtxoManager is a free data retrieval call binding the contract method 0xab8db1b1.
//
// Solidity: function utxoManager() view returns(address)
func (_Zkbridge *ZkbridgeCaller) UtxoManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Zkbridge.contract.Call(opts, &out, "utxoManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// UtxoManager is a free data retrieval call binding the contract method 0xab8db1b1.
//
// Solidity: function utxoManager() view returns(address)
func (_Zkbridge *ZkbridgeSession) UtxoManager() (common.Address, error) {
	return _Zkbridge.Contract.UtxoManager(&_Zkbridge.CallOpts)
}

// UtxoManager is a free data retrieval call binding the contract method 0xab8db1b1.
//
// Solidity: function utxoManager() view returns(address)
func (_Zkbridge *ZkbridgeCallerSession) UtxoManager() (common.Address, error) {
	return _Zkbridge.Contract.UtxoManager(&_Zkbridge.CallOpts)
}

// VerifyChangeProof is a free data retrieval call binding the contract method 0xde4b7179.
//
// Solidity: function verifyChangeProof(bytes32 txid, bytes proofData) pure returns(bool)
func (_Zkbridge *ZkbridgeCaller) VerifyChangeProof(opts *bind.CallOpts, txid [32]byte, proofData []byte) (bool, error) {
	var out []interface{}
	err := _Zkbridge.contract.Call(opts, &out, "verifyChangeProof", txid, proofData)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VerifyChangeProof is a free data retrieval call binding the contract method 0xde4b7179.
//
// Solidity: function verifyChangeProof(bytes32 txid, bytes proofData) pure returns(bool)
func (_Zkbridge *ZkbridgeSession) VerifyChangeProof(txid [32]byte, proofData []byte) (bool, error) {
	return _Zkbridge.Contract.VerifyChangeProof(&_Zkbridge.CallOpts, txid, proofData)
}

// VerifyChangeProof is a free data retrieval call binding the contract method 0xde4b7179.
//
// Solidity: function verifyChangeProof(bytes32 txid, bytes proofData) pure returns(bool)
func (_Zkbridge *ZkbridgeCallerSession) VerifyChangeProof(txid [32]byte, proofData []byte) (bool, error) {
	return _Zkbridge.Contract.VerifyChangeProof(&_Zkbridge.CallOpts, txid, proofData)
}

// VerifyDepositProof is a free data retrieval call binding the contract method 0x0c8d9769.
//
// Solidity: function verifyDepositProof(bytes32 txid, uint32 index, uint256 amount, address account, bytes proofData) pure returns(bool)
func (_Zkbridge *ZkbridgeCaller) VerifyDepositProof(opts *bind.CallOpts, txid [32]byte, index uint32, amount *big.Int, account common.Address, proofData []byte) (bool, error) {
	var out []interface{}
	err := _Zkbridge.contract.Call(opts, &out, "verifyDepositProof", txid, index, amount, account, proofData)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VerifyDepositProof is a free data retrieval call binding the contract method 0x0c8d9769.
//
// Solidity: function verifyDepositProof(bytes32 txid, uint32 index, uint256 amount, address account, bytes proofData) pure returns(bool)
func (_Zkbridge *ZkbridgeSession) VerifyDepositProof(txid [32]byte, index uint32, amount *big.Int, account common.Address, proofData []byte) (bool, error) {
	return _Zkbridge.Contract.VerifyDepositProof(&_Zkbridge.CallOpts, txid, index, amount, account, proofData)
}

// VerifyDepositProof is a free data retrieval call binding the contract method 0x0c8d9769.
//
// Solidity: function verifyDepositProof(bytes32 txid, uint32 index, uint256 amount, address account, bytes proofData) pure returns(bool)
func (_Zkbridge *ZkbridgeCallerSession) VerifyDepositProof(txid [32]byte, index uint32, amount *big.Int, account common.Address, proofData []byte) (bool, error) {
	return _Zkbridge.Contract.VerifyDepositProof(&_Zkbridge.CallOpts, txid, index, amount, account, proofData)
}

// ZkBTCToken is a free data retrieval call binding the contract method 0x281904e4.
//
// Solidity: function zkBTCToken() view returns(address)
func (_Zkbridge *ZkbridgeCaller) ZkBTCToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Zkbridge.contract.Call(opts, &out, "zkBTCToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ZkBTCToken is a free data retrieval call binding the contract method 0x281904e4.
//
// Solidity: function zkBTCToken() view returns(address)
func (_Zkbridge *ZkbridgeSession) ZkBTCToken() (common.Address, error) {
	return _Zkbridge.Contract.ZkBTCToken(&_Zkbridge.CallOpts)
}

// ZkBTCToken is a free data retrieval call binding the contract method 0x281904e4.
//
// Solidity: function zkBTCToken() view returns(address)
func (_Zkbridge *ZkbridgeCallerSession) ZkBTCToken() (common.Address, error) {
	return _Zkbridge.Contract.ZkBTCToken(&_Zkbridge.CallOpts)
}

// AddOperator is a paid mutator transaction binding the contract method 0x9870d7fe.
//
// Solidity: function addOperator(address _new) returns()
func (_Zkbridge *ZkbridgeTransactor) AddOperator(opts *bind.TransactOpts, _new common.Address) (*types.Transaction, error) {
	return _Zkbridge.contract.Transact(opts, "addOperator", _new)
}

// AddOperator is a paid mutator transaction binding the contract method 0x9870d7fe.
//
// Solidity: function addOperator(address _new) returns()
func (_Zkbridge *ZkbridgeSession) AddOperator(_new common.Address) (*types.Transaction, error) {
	return _Zkbridge.Contract.AddOperator(&_Zkbridge.TransactOpts, _new)
}

// AddOperator is a paid mutator transaction binding the contract method 0x9870d7fe.
//
// Solidity: function addOperator(address _new) returns()
func (_Zkbridge *ZkbridgeTransactorSession) AddOperator(_new common.Address) (*types.Transaction, error) {
	return _Zkbridge.Contract.AddOperator(&_Zkbridge.TransactOpts, _new)
}

// Deposit is a paid mutator transaction binding the contract method 0x11f4a8b5.
//
// Solidity: function deposit(bytes32 txid, uint32 index, uint256 amount, address receiveAddress, bytes proofData) returns()
func (_Zkbridge *ZkbridgeTransactor) Deposit(opts *bind.TransactOpts, txid [32]byte, index uint32, amount *big.Int, receiveAddress common.Address, proofData []byte) (*types.Transaction, error) {
	return _Zkbridge.contract.Transact(opts, "deposit", txid, index, amount, receiveAddress, proofData)
}

// Deposit is a paid mutator transaction binding the contract method 0x11f4a8b5.
//
// Solidity: function deposit(bytes32 txid, uint32 index, uint256 amount, address receiveAddress, bytes proofData) returns()
func (_Zkbridge *ZkbridgeSession) Deposit(txid [32]byte, index uint32, amount *big.Int, receiveAddress common.Address, proofData []byte) (*types.Transaction, error) {
	return _Zkbridge.Contract.Deposit(&_Zkbridge.TransactOpts, txid, index, amount, receiveAddress, proofData)
}

// Deposit is a paid mutator transaction binding the contract method 0x11f4a8b5.
//
// Solidity: function deposit(bytes32 txid, uint32 index, uint256 amount, address receiveAddress, bytes proofData) returns()
func (_Zkbridge *ZkbridgeTransactorSession) Deposit(txid [32]byte, index uint32, amount *big.Int, receiveAddress common.Address, proofData []byte) (*types.Transaction, error) {
	return _Zkbridge.Contract.Deposit(&_Zkbridge.TransactOpts, txid, index, amount, receiveAddress, proofData)
}

// Redeem is a paid mutator transaction binding the contract method 0x0449015a.
//
// Solidity: function redeem(uint256 redeemAmount, uint256 btcMinerFee, bytes receiveLockScript) returns()
func (_Zkbridge *ZkbridgeTransactor) Redeem(opts *bind.TransactOpts, redeemAmount *big.Int, btcMinerFee *big.Int, receiveLockScript []byte) (*types.Transaction, error) {
	return _Zkbridge.contract.Transact(opts, "redeem", redeemAmount, btcMinerFee, receiveLockScript)
}

// Redeem is a paid mutator transaction binding the contract method 0x0449015a.
//
// Solidity: function redeem(uint256 redeemAmount, uint256 btcMinerFee, bytes receiveLockScript) returns()
func (_Zkbridge *ZkbridgeSession) Redeem(redeemAmount *big.Int, btcMinerFee *big.Int, receiveLockScript []byte) (*types.Transaction, error) {
	return _Zkbridge.Contract.Redeem(&_Zkbridge.TransactOpts, redeemAmount, btcMinerFee, receiveLockScript)
}

// Redeem is a paid mutator transaction binding the contract method 0x0449015a.
//
// Solidity: function redeem(uint256 redeemAmount, uint256 btcMinerFee, bytes receiveLockScript) returns()
func (_Zkbridge *ZkbridgeTransactorSession) Redeem(redeemAmount *big.Int, btcMinerFee *big.Int, receiveLockScript []byte) (*types.Transaction, error) {
	return _Zkbridge.Contract.Redeem(&_Zkbridge.TransactOpts, redeemAmount, btcMinerFee, receiveLockScript)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0xac8a584a.
//
// Solidity: function removeOperator(address _new) returns()
func (_Zkbridge *ZkbridgeTransactor) RemoveOperator(opts *bind.TransactOpts, _new common.Address) (*types.Transaction, error) {
	return _Zkbridge.contract.Transact(opts, "removeOperator", _new)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0xac8a584a.
//
// Solidity: function removeOperator(address _new) returns()
func (_Zkbridge *ZkbridgeSession) RemoveOperator(_new common.Address) (*types.Transaction, error) {
	return _Zkbridge.Contract.RemoveOperator(&_Zkbridge.TransactOpts, _new)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0xac8a584a.
//
// Solidity: function removeOperator(address _new) returns()
func (_Zkbridge *ZkbridgeTransactorSession) RemoveOperator(_new common.Address) (*types.Transaction, error) {
	return _Zkbridge.Contract.RemoveOperator(&_Zkbridge.TransactOpts, _new)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Zkbridge *ZkbridgeTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Zkbridge.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Zkbridge *ZkbridgeSession) RenounceOwnership() (*types.Transaction, error) {
	return _Zkbridge.Contract.RenounceOwnership(&_Zkbridge.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Zkbridge *ZkbridgeTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Zkbridge.Contract.RenounceOwnership(&_Zkbridge.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Zkbridge *ZkbridgeTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Zkbridge.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Zkbridge *ZkbridgeSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Zkbridge.Contract.TransferOwnership(&_Zkbridge.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Zkbridge *ZkbridgeTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Zkbridge.Contract.TransferOwnership(&_Zkbridge.TransactOpts, newOwner)
}

// UpdateChange is a paid mutator transaction binding the contract method 0x8ce343ab.
//
// Solidity: function updateChange(bytes32 txid, bytes proofData) returns()
func (_Zkbridge *ZkbridgeTransactor) UpdateChange(opts *bind.TransactOpts, txid [32]byte, proofData []byte) (*types.Transaction, error) {
	return _Zkbridge.contract.Transact(opts, "updateChange", txid, proofData)
}

// UpdateChange is a paid mutator transaction binding the contract method 0x8ce343ab.
//
// Solidity: function updateChange(bytes32 txid, bytes proofData) returns()
func (_Zkbridge *ZkbridgeSession) UpdateChange(txid [32]byte, proofData []byte) (*types.Transaction, error) {
	return _Zkbridge.Contract.UpdateChange(&_Zkbridge.TransactOpts, txid, proofData)
}

// UpdateChange is a paid mutator transaction binding the contract method 0x8ce343ab.
//
// Solidity: function updateChange(bytes32 txid, bytes proofData) returns()
func (_Zkbridge *ZkbridgeTransactorSession) UpdateChange(txid [32]byte, proofData []byte) (*types.Transaction, error) {
	return _Zkbridge.Contract.UpdateChange(&_Zkbridge.TransactOpts, txid, proofData)
}

// UpdateEconomicAddress is a paid mutator transaction binding the contract method 0x5e5a769a.
//
// Solidity: function updateEconomicAddress(address _newAddress) returns()
func (_Zkbridge *ZkbridgeTransactor) UpdateEconomicAddress(opts *bind.TransactOpts, _newAddress common.Address) (*types.Transaction, error) {
	return _Zkbridge.contract.Transact(opts, "updateEconomicAddress", _newAddress)
}

// UpdateEconomicAddress is a paid mutator transaction binding the contract method 0x5e5a769a.
//
// Solidity: function updateEconomicAddress(address _newAddress) returns()
func (_Zkbridge *ZkbridgeSession) UpdateEconomicAddress(_newAddress common.Address) (*types.Transaction, error) {
	return _Zkbridge.Contract.UpdateEconomicAddress(&_Zkbridge.TransactOpts, _newAddress)
}

// UpdateEconomicAddress is a paid mutator transaction binding the contract method 0x5e5a769a.
//
// Solidity: function updateEconomicAddress(address _newAddress) returns()
func (_Zkbridge *ZkbridgeTransactorSession) UpdateEconomicAddress(_newAddress common.Address) (*types.Transaction, error) {
	return _Zkbridge.Contract.UpdateEconomicAddress(&_Zkbridge.TransactOpts, _newAddress)
}

// UpdateUtxoAddress is a paid mutator transaction binding the contract method 0x9314c393.
//
// Solidity: function updateUtxoAddress(address _utxoAddress, bytes _changeLockScript, bytes _multiSigScript) returns()
func (_Zkbridge *ZkbridgeTransactor) UpdateUtxoAddress(opts *bind.TransactOpts, _utxoAddress common.Address, _changeLockScript []byte, _multiSigScript []byte) (*types.Transaction, error) {
	return _Zkbridge.contract.Transact(opts, "updateUtxoAddress", _utxoAddress, _changeLockScript, _multiSigScript)
}

// UpdateUtxoAddress is a paid mutator transaction binding the contract method 0x9314c393.
//
// Solidity: function updateUtxoAddress(address _utxoAddress, bytes _changeLockScript, bytes _multiSigScript) returns()
func (_Zkbridge *ZkbridgeSession) UpdateUtxoAddress(_utxoAddress common.Address, _changeLockScript []byte, _multiSigScript []byte) (*types.Transaction, error) {
	return _Zkbridge.Contract.UpdateUtxoAddress(&_Zkbridge.TransactOpts, _utxoAddress, _changeLockScript, _multiSigScript)
}

// UpdateUtxoAddress is a paid mutator transaction binding the contract method 0x9314c393.
//
// Solidity: function updateUtxoAddress(address _utxoAddress, bytes _changeLockScript, bytes _multiSigScript) returns()
func (_Zkbridge *ZkbridgeTransactorSession) UpdateUtxoAddress(_utxoAddress common.Address, _changeLockScript []byte, _multiSigScript []byte) (*types.Transaction, error) {
	return _Zkbridge.Contract.UpdateUtxoAddress(&_Zkbridge.TransactOpts, _utxoAddress, _changeLockScript, _multiSigScript)
}

// ZkbridgeCreateRedeemUnsignedTxIterator is returned from FilterCreateRedeemUnsignedTx and is used to iterate over the raw logs and unpacked data for CreateRedeemUnsignedTx events raised by the Zkbridge contract.
type ZkbridgeCreateRedeemUnsignedTxIterator struct {
	Event *ZkbridgeCreateRedeemUnsignedTx // Event containing the contract specifics and raw log

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
func (it *ZkbridgeCreateRedeemUnsignedTxIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZkbridgeCreateRedeemUnsignedTx)
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
		it.Event = new(ZkbridgeCreateRedeemUnsignedTx)
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
func (it *ZkbridgeCreateRedeemUnsignedTxIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZkbridgeCreateRedeemUnsignedTxIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZkbridgeCreateRedeemUnsignedTx represents a CreateRedeemUnsignedTx event raised by the Zkbridge contract.
type ZkbridgeCreateRedeemUnsignedTx struct {
	Txid     [32]byte
	RawTx    []byte
	SigHashs [][32]byte
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterCreateRedeemUnsignedTx is a free log retrieval operation binding the contract event 0x1e5e2baa6d11cc5bcae8c0d1187d7b9ebf13d6d9b932f7dbbf4e396438845fb8.
//
// Solidity: event CreateRedeemUnsignedTx(bytes32 indexed txid, bytes rawTx, bytes32[] sigHashs)
func (_Zkbridge *ZkbridgeFilterer) FilterCreateRedeemUnsignedTx(opts *bind.FilterOpts, txid [][32]byte) (*ZkbridgeCreateRedeemUnsignedTxIterator, error) {

	var txidRule []interface{}
	for _, txidItem := range txid {
		txidRule = append(txidRule, txidItem)
	}

	logs, sub, err := _Zkbridge.contract.FilterLogs(opts, "CreateRedeemUnsignedTx", txidRule)
	if err != nil {
		return nil, err
	}
	return &ZkbridgeCreateRedeemUnsignedTxIterator{contract: _Zkbridge.contract, event: "CreateRedeemUnsignedTx", logs: logs, sub: sub}, nil
}

// WatchCreateRedeemUnsignedTx is a free log subscription operation binding the contract event 0x1e5e2baa6d11cc5bcae8c0d1187d7b9ebf13d6d9b932f7dbbf4e396438845fb8.
//
// Solidity: event CreateRedeemUnsignedTx(bytes32 indexed txid, bytes rawTx, bytes32[] sigHashs)
func (_Zkbridge *ZkbridgeFilterer) WatchCreateRedeemUnsignedTx(opts *bind.WatchOpts, sink chan<- *ZkbridgeCreateRedeemUnsignedTx, txid [][32]byte) (event.Subscription, error) {

	var txidRule []interface{}
	for _, txidItem := range txid {
		txidRule = append(txidRule, txidItem)
	}

	logs, sub, err := _Zkbridge.contract.WatchLogs(opts, "CreateRedeemUnsignedTx", txidRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZkbridgeCreateRedeemUnsignedTx)
				if err := _Zkbridge.contract.UnpackLog(event, "CreateRedeemUnsignedTx", log); err != nil {
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
func (_Zkbridge *ZkbridgeFilterer) ParseCreateRedeemUnsignedTx(log types.Log) (*ZkbridgeCreateRedeemUnsignedTx, error) {
	event := new(ZkbridgeCreateRedeemUnsignedTx)
	if err := _Zkbridge.contract.UnpackLog(event, "CreateRedeemUnsignedTx", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZkbridgeOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Zkbridge contract.
type ZkbridgeOwnershipTransferredIterator struct {
	Event *ZkbridgeOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ZkbridgeOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZkbridgeOwnershipTransferred)
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
		it.Event = new(ZkbridgeOwnershipTransferred)
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
func (it *ZkbridgeOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZkbridgeOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZkbridgeOwnershipTransferred represents a OwnershipTransferred event raised by the Zkbridge contract.
type ZkbridgeOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Zkbridge *ZkbridgeFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ZkbridgeOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Zkbridge.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ZkbridgeOwnershipTransferredIterator{contract: _Zkbridge.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Zkbridge *ZkbridgeFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ZkbridgeOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Zkbridge.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZkbridgeOwnershipTransferred)
				if err := _Zkbridge.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Zkbridge *ZkbridgeFilterer) ParseOwnershipTransferred(log types.Log) (*ZkbridgeOwnershipTransferred, error) {
	event := new(ZkbridgeOwnershipTransferred)
	if err := _Zkbridge.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
