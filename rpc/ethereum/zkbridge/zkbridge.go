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

// IBtcTxVerifierPublicWitnessParams is an auto generated low-level Go binding around an user-defined struct.
type IBtcTxVerifierPublicWitnessParams struct {
	Checkpoint        [32]byte
	CpDepth           uint32
	TxDepth           uint32
	TxBlockHash       [32]byte
	TxTimestamp       uint32
	ZkpMiner          common.Address
	Flag              *big.Int
	SmoothedTimestamp uint32
}

// ZkBTCBridgeOperator is an auto generated low-level Go binding around an user-defined struct.
type ZkBTCBridgeOperator struct {
	ChangeLockScript []byte
	MultiSigScript   []byte
	UtxoManager      common.Address
}

// ZkbridgeMetaData contains all meta data concerning the Zkbridge contract.
var ZkbridgeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIFeePool\",\"name\":\"_feeAccount\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_changeLockScript\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"_multiSigScript\",\"type\":\"bytes\"},{\"internalType\":\"contractUTXOManagerInterface\",\"name\":\"_utxoAddress\",\"type\":\"address\"},{\"internalType\":\"contractIBtcTxVerifier\",\"name\":\"_txVerifier\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"depositAmount\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"minDepositAmount\",\"type\":\"uint64\"}],\"name\":\"InvalidDepositAmount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"lockScriptLength\",\"type\":\"uint256\"}],\"name\":\"InvalidLockScriptLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"proofData\",\"type\":\"bytes\"}],\"name\":\"InvalidRedeemOrProof\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"lockScript\",\"type\":\"bytes\"}],\"name\":\"LockScriptIsChange\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes\",\"name\":\"sourceLockingScript\",\"type\":\"bytes\"},{\"indexed\":true,\"internalType\":\"bytes\",\"name\":\"destinationLockingScript\",\"type\":\"bytes\"}],\"name\":\"BeginMigration\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"minerReward\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32[]\",\"name\":\"sigHashs\",\"type\":\"bytes32[]\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"rawTx\",\"type\":\"bytes\"}],\"name\":\"CreateRedeemUnsignedTx\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"old\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"now\",\"type\":\"uint64\"}],\"name\":\"DustThresholdsUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"}],\"name\":\"RetrievedData\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"old\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"now\",\"type\":\"address\"}],\"name\":\"UpdateReserveInterface\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MIN_DEPOSIT_AMOUNT\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"changeLockScript\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"multiSigScript\",\"type\":\"bytes\"},{\"internalType\":\"contractUTXOManagerInterface\",\"name\":\"utxoManager\",\"type\":\"address\"}],\"internalType\":\"structzkBTCBridge.Operator\",\"name\":\"destinationOperator\",\"type\":\"tuple\"}],\"name\":\"beginMigration\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"rawBtcTx\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"checkpoint\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"cpDepth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"txDepth\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"txBlockHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"txTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"zkpMiner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"flag\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"smoothedTimestamp\",\"type\":\"uint32\"}],\"internalType\":\"structIBtcTxVerifier.PublicWitnessParams\",\"name\":\"params\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"proofData\",\"type\":\"bytes\"}],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dustThresholds\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feeAccount\",\"outputs\":[{\"internalType\":\"contractIFeePool\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLatestOperator\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"changeLockScript\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"multiSigScript\",\"type\":\"bytes\"},{\"internalType\":\"contractUTXOManagerInterface\",\"name\":\"utxoManager\",\"type\":\"address\"}],\"internalType\":\"structzkBTCBridge.Operator\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"script\",\"type\":\"bytes\"}],\"name\":\"getOperator\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"changeLockScript\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"multiSigScript\",\"type\":\"bytes\"},{\"internalType\":\"contractUTXOManagerInterface\",\"name\":\"utxoManager\",\"type\":\"address\"}],\"internalType\":\"structzkBTCBridge.Operator\",\"name\":\"op\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"txBlockHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"amount\",\"type\":\"uint64\"}],\"name\":\"getRaiseIf\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"getRoleMember\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleMemberCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleMembers\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"redeemAmount\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"btcMinerFee\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"sourceScript\",\"type\":\"bytes\"}],\"name\":\"migrate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"migrationList\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"changeLockScript\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"multiSigScript\",\"type\":\"bytes\"},{\"internalType\":\"contractUTXOManagerInterface\",\"name\":\"utxoManager\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"provenDatas\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"associatedAmount\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"retrieved\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"redeemAmount\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"btcMinerFee\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"recipientLockScript\",\"type\":\"bytes\"}],\"name\":\"redeem\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reserveInterface\",\"outputs\":[{\"internalType\":\"contractReserveInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"}],\"name\":\"retrieve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"suggestPrice\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"txCounts\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"txVerifier\",\"outputs\":[{\"internalType\":\"contractIBtcTxVerifier\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"newThresholds\",\"type\":\"uint64\"}],\"name\":\"updateDustThresholds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"checkpoint\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"cpDepth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"txDepth\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"txBlockHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"txTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"zkpMiner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"flag\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"smoothedTimestamp\",\"type\":\"uint32\"}],\"internalType\":\"structIBtcTxVerifier.PublicWitnessParams\",\"name\":\"params\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"minerReward\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"proofData\",\"type\":\"bytes\"}],\"name\":\"updateRedeem\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractReserveInterface\",\"name\":\"newInterface\",\"type\":\"address\"}],\"name\":\"updateReserveInterface\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
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
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ZkbridgeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ZkbridgeCallerSession struct {
	Contract *ZkbridgeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
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

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Zkbridge *ZkbridgeCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Zkbridge.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Zkbridge *ZkbridgeSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Zkbridge.Contract.DEFAULTADMINROLE(&_Zkbridge.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Zkbridge *ZkbridgeCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Zkbridge.Contract.DEFAULTADMINROLE(&_Zkbridge.CallOpts)
}

// MINDEPOSITAMOUNT is a free data retrieval call binding the contract method 0x1ea30fef.
//
// Solidity: function MIN_DEPOSIT_AMOUNT() view returns(uint64)
func (_Zkbridge *ZkbridgeCaller) MINDEPOSITAMOUNT(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Zkbridge.contract.Call(opts, &out, "MIN_DEPOSIT_AMOUNT")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// MINDEPOSITAMOUNT is a free data retrieval call binding the contract method 0x1ea30fef.
//
// Solidity: function MIN_DEPOSIT_AMOUNT() view returns(uint64)
func (_Zkbridge *ZkbridgeSession) MINDEPOSITAMOUNT() (uint64, error) {
	return _Zkbridge.Contract.MINDEPOSITAMOUNT(&_Zkbridge.CallOpts)
}

// MINDEPOSITAMOUNT is a free data retrieval call binding the contract method 0x1ea30fef.
//
// Solidity: function MIN_DEPOSIT_AMOUNT() view returns(uint64)
func (_Zkbridge *ZkbridgeCallerSession) MINDEPOSITAMOUNT() (uint64, error) {
	return _Zkbridge.Contract.MINDEPOSITAMOUNT(&_Zkbridge.CallOpts)
}

// DustThresholds is a free data retrieval call binding the contract method 0x9609f7d3.
//
// Solidity: function dustThresholds() view returns(uint64)
func (_Zkbridge *ZkbridgeCaller) DustThresholds(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Zkbridge.contract.Call(opts, &out, "dustThresholds")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// DustThresholds is a free data retrieval call binding the contract method 0x9609f7d3.
//
// Solidity: function dustThresholds() view returns(uint64)
func (_Zkbridge *ZkbridgeSession) DustThresholds() (uint64, error) {
	return _Zkbridge.Contract.DustThresholds(&_Zkbridge.CallOpts)
}

// DustThresholds is a free data retrieval call binding the contract method 0x9609f7d3.
//
// Solidity: function dustThresholds() view returns(uint64)
func (_Zkbridge *ZkbridgeCallerSession) DustThresholds() (uint64, error) {
	return _Zkbridge.Contract.DustThresholds(&_Zkbridge.CallOpts)
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

// GetLatestOperator is a free data retrieval call binding the contract method 0x88716b28.
//
// Solidity: function getLatestOperator() view returns((bytes,bytes,address))
func (_Zkbridge *ZkbridgeCaller) GetLatestOperator(opts *bind.CallOpts) (ZkBTCBridgeOperator, error) {
	var out []interface{}
	err := _Zkbridge.contract.Call(opts, &out, "getLatestOperator")

	if err != nil {
		return *new(ZkBTCBridgeOperator), err
	}

	out0 := *abi.ConvertType(out[0], new(ZkBTCBridgeOperator)).(*ZkBTCBridgeOperator)

	return out0, err

}

// GetLatestOperator is a free data retrieval call binding the contract method 0x88716b28.
//
// Solidity: function getLatestOperator() view returns((bytes,bytes,address))
func (_Zkbridge *ZkbridgeSession) GetLatestOperator() (ZkBTCBridgeOperator, error) {
	return _Zkbridge.Contract.GetLatestOperator(&_Zkbridge.CallOpts)
}

// GetLatestOperator is a free data retrieval call binding the contract method 0x88716b28.
//
// Solidity: function getLatestOperator() view returns((bytes,bytes,address))
func (_Zkbridge *ZkbridgeCallerSession) GetLatestOperator() (ZkBTCBridgeOperator, error) {
	return _Zkbridge.Contract.GetLatestOperator(&_Zkbridge.CallOpts)
}

// GetOperator is a free data retrieval call binding the contract method 0x9eaffa96.
//
// Solidity: function getOperator(bytes script) view returns((bytes,bytes,address) op)
func (_Zkbridge *ZkbridgeCaller) GetOperator(opts *bind.CallOpts, script []byte) (ZkBTCBridgeOperator, error) {
	var out []interface{}
	err := _Zkbridge.contract.Call(opts, &out, "getOperator", script)

	if err != nil {
		return *new(ZkBTCBridgeOperator), err
	}

	out0 := *abi.ConvertType(out[0], new(ZkBTCBridgeOperator)).(*ZkBTCBridgeOperator)

	return out0, err

}

// GetOperator is a free data retrieval call binding the contract method 0x9eaffa96.
//
// Solidity: function getOperator(bytes script) view returns((bytes,bytes,address) op)
func (_Zkbridge *ZkbridgeSession) GetOperator(script []byte) (ZkBTCBridgeOperator, error) {
	return _Zkbridge.Contract.GetOperator(&_Zkbridge.CallOpts, script)
}

// GetOperator is a free data retrieval call binding the contract method 0x9eaffa96.
//
// Solidity: function getOperator(bytes script) view returns((bytes,bytes,address) op)
func (_Zkbridge *ZkbridgeCallerSession) GetOperator(script []byte) (ZkBTCBridgeOperator, error) {
	return _Zkbridge.Contract.GetOperator(&_Zkbridge.CallOpts, script)
}

// GetRaiseIf is a free data retrieval call binding the contract method 0x964dc667.
//
// Solidity: function getRaiseIf(bytes32 txBlockHash, uint64 amount) view returns(bool)
func (_Zkbridge *ZkbridgeCaller) GetRaiseIf(opts *bind.CallOpts, txBlockHash [32]byte, amount uint64) (bool, error) {
	var out []interface{}
	err := _Zkbridge.contract.Call(opts, &out, "getRaiseIf", txBlockHash, amount)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetRaiseIf is a free data retrieval call binding the contract method 0x964dc667.
//
// Solidity: function getRaiseIf(bytes32 txBlockHash, uint64 amount) view returns(bool)
func (_Zkbridge *ZkbridgeSession) GetRaiseIf(txBlockHash [32]byte, amount uint64) (bool, error) {
	return _Zkbridge.Contract.GetRaiseIf(&_Zkbridge.CallOpts, txBlockHash, amount)
}

// GetRaiseIf is a free data retrieval call binding the contract method 0x964dc667.
//
// Solidity: function getRaiseIf(bytes32 txBlockHash, uint64 amount) view returns(bool)
func (_Zkbridge *ZkbridgeCallerSession) GetRaiseIf(txBlockHash [32]byte, amount uint64) (bool, error) {
	return _Zkbridge.Contract.GetRaiseIf(&_Zkbridge.CallOpts, txBlockHash, amount)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Zkbridge *ZkbridgeCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _Zkbridge.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Zkbridge *ZkbridgeSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Zkbridge.Contract.GetRoleAdmin(&_Zkbridge.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Zkbridge *ZkbridgeCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Zkbridge.Contract.GetRoleAdmin(&_Zkbridge.CallOpts, role)
}

// GetRoleMember is a free data retrieval call binding the contract method 0x9010d07c.
//
// Solidity: function getRoleMember(bytes32 role, uint256 index) view returns(address)
func (_Zkbridge *ZkbridgeCaller) GetRoleMember(opts *bind.CallOpts, role [32]byte, index *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Zkbridge.contract.Call(opts, &out, "getRoleMember", role, index)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetRoleMember is a free data retrieval call binding the contract method 0x9010d07c.
//
// Solidity: function getRoleMember(bytes32 role, uint256 index) view returns(address)
func (_Zkbridge *ZkbridgeSession) GetRoleMember(role [32]byte, index *big.Int) (common.Address, error) {
	return _Zkbridge.Contract.GetRoleMember(&_Zkbridge.CallOpts, role, index)
}

// GetRoleMember is a free data retrieval call binding the contract method 0x9010d07c.
//
// Solidity: function getRoleMember(bytes32 role, uint256 index) view returns(address)
func (_Zkbridge *ZkbridgeCallerSession) GetRoleMember(role [32]byte, index *big.Int) (common.Address, error) {
	return _Zkbridge.Contract.GetRoleMember(&_Zkbridge.CallOpts, role, index)
}

// GetRoleMemberCount is a free data retrieval call binding the contract method 0xca15c873.
//
// Solidity: function getRoleMemberCount(bytes32 role) view returns(uint256)
func (_Zkbridge *ZkbridgeCaller) GetRoleMemberCount(opts *bind.CallOpts, role [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _Zkbridge.contract.Call(opts, &out, "getRoleMemberCount", role)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRoleMemberCount is a free data retrieval call binding the contract method 0xca15c873.
//
// Solidity: function getRoleMemberCount(bytes32 role) view returns(uint256)
func (_Zkbridge *ZkbridgeSession) GetRoleMemberCount(role [32]byte) (*big.Int, error) {
	return _Zkbridge.Contract.GetRoleMemberCount(&_Zkbridge.CallOpts, role)
}

// GetRoleMemberCount is a free data retrieval call binding the contract method 0xca15c873.
//
// Solidity: function getRoleMemberCount(bytes32 role) view returns(uint256)
func (_Zkbridge *ZkbridgeCallerSession) GetRoleMemberCount(role [32]byte) (*big.Int, error) {
	return _Zkbridge.Contract.GetRoleMemberCount(&_Zkbridge.CallOpts, role)
}

// GetRoleMembers is a free data retrieval call binding the contract method 0xa3246ad3.
//
// Solidity: function getRoleMembers(bytes32 role) view returns(address[])
func (_Zkbridge *ZkbridgeCaller) GetRoleMembers(opts *bind.CallOpts, role [32]byte) ([]common.Address, error) {
	var out []interface{}
	err := _Zkbridge.contract.Call(opts, &out, "getRoleMembers", role)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetRoleMembers is a free data retrieval call binding the contract method 0xa3246ad3.
//
// Solidity: function getRoleMembers(bytes32 role) view returns(address[])
func (_Zkbridge *ZkbridgeSession) GetRoleMembers(role [32]byte) ([]common.Address, error) {
	return _Zkbridge.Contract.GetRoleMembers(&_Zkbridge.CallOpts, role)
}

// GetRoleMembers is a free data retrieval call binding the contract method 0xa3246ad3.
//
// Solidity: function getRoleMembers(bytes32 role) view returns(address[])
func (_Zkbridge *ZkbridgeCallerSession) GetRoleMembers(role [32]byte) ([]common.Address, error) {
	return _Zkbridge.Contract.GetRoleMembers(&_Zkbridge.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Zkbridge *ZkbridgeCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _Zkbridge.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Zkbridge *ZkbridgeSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Zkbridge.Contract.HasRole(&_Zkbridge.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Zkbridge *ZkbridgeCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Zkbridge.Contract.HasRole(&_Zkbridge.CallOpts, role, account)
}

// MigrationList is a free data retrieval call binding the contract method 0xca1e161d.
//
// Solidity: function migrationList(uint256 ) view returns(bytes changeLockScript, bytes multiSigScript, address utxoManager)
func (_Zkbridge *ZkbridgeCaller) MigrationList(opts *bind.CallOpts, arg0 *big.Int) (struct {
	ChangeLockScript []byte
	MultiSigScript   []byte
	UtxoManager      common.Address
}, error) {
	var out []interface{}
	err := _Zkbridge.contract.Call(opts, &out, "migrationList", arg0)

	outstruct := new(struct {
		ChangeLockScript []byte
		MultiSigScript   []byte
		UtxoManager      common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ChangeLockScript = *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	outstruct.MultiSigScript = *abi.ConvertType(out[1], new([]byte)).(*[]byte)
	outstruct.UtxoManager = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)

	return *outstruct, err

}

// MigrationList is a free data retrieval call binding the contract method 0xca1e161d.
//
// Solidity: function migrationList(uint256 ) view returns(bytes changeLockScript, bytes multiSigScript, address utxoManager)
func (_Zkbridge *ZkbridgeSession) MigrationList(arg0 *big.Int) (struct {
	ChangeLockScript []byte
	MultiSigScript   []byte
	UtxoManager      common.Address
}, error) {
	return _Zkbridge.Contract.MigrationList(&_Zkbridge.CallOpts, arg0)
}

// MigrationList is a free data retrieval call binding the contract method 0xca1e161d.
//
// Solidity: function migrationList(uint256 ) view returns(bytes changeLockScript, bytes multiSigScript, address utxoManager)
func (_Zkbridge *ZkbridgeCallerSession) MigrationList(arg0 *big.Int) (struct {
	ChangeLockScript []byte
	MultiSigScript   []byte
	UtxoManager      common.Address
}, error) {
	return _Zkbridge.Contract.MigrationList(&_Zkbridge.CallOpts, arg0)
}

// ProvenDatas is a free data retrieval call binding the contract method 0xa67dca44.
//
// Solidity: function provenDatas(bytes32 ) view returns(uint32 index, bytes32 blockHash, uint64 associatedAmount, bytes data, bool retrieved)
func (_Zkbridge *ZkbridgeCaller) ProvenDatas(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Index            uint32
	BlockHash        [32]byte
	AssociatedAmount uint64
	Data             []byte
	Retrieved        bool
}, error) {
	var out []interface{}
	err := _Zkbridge.contract.Call(opts, &out, "provenDatas", arg0)

	outstruct := new(struct {
		Index            uint32
		BlockHash        [32]byte
		AssociatedAmount uint64
		Data             []byte
		Retrieved        bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Index = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.BlockHash = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.AssociatedAmount = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	outstruct.Data = *abi.ConvertType(out[3], new([]byte)).(*[]byte)
	outstruct.Retrieved = *abi.ConvertType(out[4], new(bool)).(*bool)

	return *outstruct, err

}

// ProvenDatas is a free data retrieval call binding the contract method 0xa67dca44.
//
// Solidity: function provenDatas(bytes32 ) view returns(uint32 index, bytes32 blockHash, uint64 associatedAmount, bytes data, bool retrieved)
func (_Zkbridge *ZkbridgeSession) ProvenDatas(arg0 [32]byte) (struct {
	Index            uint32
	BlockHash        [32]byte
	AssociatedAmount uint64
	Data             []byte
	Retrieved        bool
}, error) {
	return _Zkbridge.Contract.ProvenDatas(&_Zkbridge.CallOpts, arg0)
}

// ProvenDatas is a free data retrieval call binding the contract method 0xa67dca44.
//
// Solidity: function provenDatas(bytes32 ) view returns(uint32 index, bytes32 blockHash, uint64 associatedAmount, bytes data, bool retrieved)
func (_Zkbridge *ZkbridgeCallerSession) ProvenDatas(arg0 [32]byte) (struct {
	Index            uint32
	BlockHash        [32]byte
	AssociatedAmount uint64
	Data             []byte
	Retrieved        bool
}, error) {
	return _Zkbridge.Contract.ProvenDatas(&_Zkbridge.CallOpts, arg0)
}

// ReserveInterface is a free data retrieval call binding the contract method 0xe6ea61ba.
//
// Solidity: function reserveInterface() view returns(address)
func (_Zkbridge *ZkbridgeCaller) ReserveInterface(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Zkbridge.contract.Call(opts, &out, "reserveInterface")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ReserveInterface is a free data retrieval call binding the contract method 0xe6ea61ba.
//
// Solidity: function reserveInterface() view returns(address)
func (_Zkbridge *ZkbridgeSession) ReserveInterface() (common.Address, error) {
	return _Zkbridge.Contract.ReserveInterface(&_Zkbridge.CallOpts)
}

// ReserveInterface is a free data retrieval call binding the contract method 0xe6ea61ba.
//
// Solidity: function reserveInterface() view returns(address)
func (_Zkbridge *ZkbridgeCallerSession) ReserveInterface() (common.Address, error) {
	return _Zkbridge.Contract.ReserveInterface(&_Zkbridge.CallOpts)
}

// SuggestPrice is a free data retrieval call binding the contract method 0x7e354628.
//
// Solidity: function suggestPrice() view returns(uint64)
func (_Zkbridge *ZkbridgeCaller) SuggestPrice(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Zkbridge.contract.Call(opts, &out, "suggestPrice")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// SuggestPrice is a free data retrieval call binding the contract method 0x7e354628.
//
// Solidity: function suggestPrice() view returns(uint64)
func (_Zkbridge *ZkbridgeSession) SuggestPrice() (uint64, error) {
	return _Zkbridge.Contract.SuggestPrice(&_Zkbridge.CallOpts)
}

// SuggestPrice is a free data retrieval call binding the contract method 0x7e354628.
//
// Solidity: function suggestPrice() view returns(uint64)
func (_Zkbridge *ZkbridgeCallerSession) SuggestPrice() (uint64, error) {
	return _Zkbridge.Contract.SuggestPrice(&_Zkbridge.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Zkbridge *ZkbridgeCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _Zkbridge.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Zkbridge *ZkbridgeSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Zkbridge.Contract.SupportsInterface(&_Zkbridge.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Zkbridge *ZkbridgeCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Zkbridge.Contract.SupportsInterface(&_Zkbridge.CallOpts, interfaceId)
}

// TxCounts is a free data retrieval call binding the contract method 0xc48c2de1.
//
// Solidity: function txCounts(bytes32 ) view returns(uint32)
func (_Zkbridge *ZkbridgeCaller) TxCounts(opts *bind.CallOpts, arg0 [32]byte) (uint32, error) {
	var out []interface{}
	err := _Zkbridge.contract.Call(opts, &out, "txCounts", arg0)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// TxCounts is a free data retrieval call binding the contract method 0xc48c2de1.
//
// Solidity: function txCounts(bytes32 ) view returns(uint32)
func (_Zkbridge *ZkbridgeSession) TxCounts(arg0 [32]byte) (uint32, error) {
	return _Zkbridge.Contract.TxCounts(&_Zkbridge.CallOpts, arg0)
}

// TxCounts is a free data retrieval call binding the contract method 0xc48c2de1.
//
// Solidity: function txCounts(bytes32 ) view returns(uint32)
func (_Zkbridge *ZkbridgeCallerSession) TxCounts(arg0 [32]byte) (uint32, error) {
	return _Zkbridge.Contract.TxCounts(&_Zkbridge.CallOpts, arg0)
}

// TxVerifier is a free data retrieval call binding the contract method 0x23a45c93.
//
// Solidity: function txVerifier() view returns(address)
func (_Zkbridge *ZkbridgeCaller) TxVerifier(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Zkbridge.contract.Call(opts, &out, "txVerifier")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TxVerifier is a free data retrieval call binding the contract method 0x23a45c93.
//
// Solidity: function txVerifier() view returns(address)
func (_Zkbridge *ZkbridgeSession) TxVerifier() (common.Address, error) {
	return _Zkbridge.Contract.TxVerifier(&_Zkbridge.CallOpts)
}

// TxVerifier is a free data retrieval call binding the contract method 0x23a45c93.
//
// Solidity: function txVerifier() view returns(address)
func (_Zkbridge *ZkbridgeCallerSession) TxVerifier() (common.Address, error) {
	return _Zkbridge.Contract.TxVerifier(&_Zkbridge.CallOpts)
}

// BeginMigration is a paid mutator transaction binding the contract method 0x6e629937.
//
// Solidity: function beginMigration((bytes,bytes,address) destinationOperator) returns()
func (_Zkbridge *ZkbridgeTransactor) BeginMigration(opts *bind.TransactOpts, destinationOperator ZkBTCBridgeOperator) (*types.Transaction, error) {
	return _Zkbridge.contract.Transact(opts, "beginMigration", destinationOperator)
}

// BeginMigration is a paid mutator transaction binding the contract method 0x6e629937.
//
// Solidity: function beginMigration((bytes,bytes,address) destinationOperator) returns()
func (_Zkbridge *ZkbridgeSession) BeginMigration(destinationOperator ZkBTCBridgeOperator) (*types.Transaction, error) {
	return _Zkbridge.Contract.BeginMigration(&_Zkbridge.TransactOpts, destinationOperator)
}

// BeginMigration is a paid mutator transaction binding the contract method 0x6e629937.
//
// Solidity: function beginMigration((bytes,bytes,address) destinationOperator) returns()
func (_Zkbridge *ZkbridgeTransactorSession) BeginMigration(destinationOperator ZkBTCBridgeOperator) (*types.Transaction, error) {
	return _Zkbridge.Contract.BeginMigration(&_Zkbridge.TransactOpts, destinationOperator)
}

// Deposit is a paid mutator transaction binding the contract method 0xee631bfb.
//
// Solidity: function deposit(bytes rawBtcTx, (bytes32,uint32,uint32,bytes32,uint32,address,uint256,uint32) params, bytes proofData) returns()
func (_Zkbridge *ZkbridgeTransactor) Deposit(opts *bind.TransactOpts, rawBtcTx []byte, params IBtcTxVerifierPublicWitnessParams, proofData []byte) (*types.Transaction, error) {
	return _Zkbridge.contract.Transact(opts, "deposit", rawBtcTx, params, proofData)
}

// Deposit is a paid mutator transaction binding the contract method 0xee631bfb.
//
// Solidity: function deposit(bytes rawBtcTx, (bytes32,uint32,uint32,bytes32,uint32,address,uint256,uint32) params, bytes proofData) returns()
func (_Zkbridge *ZkbridgeSession) Deposit(rawBtcTx []byte, params IBtcTxVerifierPublicWitnessParams, proofData []byte) (*types.Transaction, error) {
	return _Zkbridge.Contract.Deposit(&_Zkbridge.TransactOpts, rawBtcTx, params, proofData)
}

// Deposit is a paid mutator transaction binding the contract method 0xee631bfb.
//
// Solidity: function deposit(bytes rawBtcTx, (bytes32,uint32,uint32,bytes32,uint32,address,uint256,uint32) params, bytes proofData) returns()
func (_Zkbridge *ZkbridgeTransactorSession) Deposit(rawBtcTx []byte, params IBtcTxVerifierPublicWitnessParams, proofData []byte) (*types.Transaction, error) {
	return _Zkbridge.Contract.Deposit(&_Zkbridge.TransactOpts, rawBtcTx, params, proofData)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Zkbridge *ZkbridgeTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Zkbridge.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Zkbridge *ZkbridgeSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Zkbridge.Contract.GrantRole(&_Zkbridge.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Zkbridge *ZkbridgeTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Zkbridge.Contract.GrantRole(&_Zkbridge.TransactOpts, role, account)
}

// Migrate is a paid mutator transaction binding the contract method 0x34f2fe70.
//
// Solidity: function migrate(uint64 redeemAmount, uint64 btcMinerFee, bytes sourceScript) returns()
func (_Zkbridge *ZkbridgeTransactor) Migrate(opts *bind.TransactOpts, redeemAmount uint64, btcMinerFee uint64, sourceScript []byte) (*types.Transaction, error) {
	return _Zkbridge.contract.Transact(opts, "migrate", redeemAmount, btcMinerFee, sourceScript)
}

// Migrate is a paid mutator transaction binding the contract method 0x34f2fe70.
//
// Solidity: function migrate(uint64 redeemAmount, uint64 btcMinerFee, bytes sourceScript) returns()
func (_Zkbridge *ZkbridgeSession) Migrate(redeemAmount uint64, btcMinerFee uint64, sourceScript []byte) (*types.Transaction, error) {
	return _Zkbridge.Contract.Migrate(&_Zkbridge.TransactOpts, redeemAmount, btcMinerFee, sourceScript)
}

// Migrate is a paid mutator transaction binding the contract method 0x34f2fe70.
//
// Solidity: function migrate(uint64 redeemAmount, uint64 btcMinerFee, bytes sourceScript) returns()
func (_Zkbridge *ZkbridgeTransactorSession) Migrate(redeemAmount uint64, btcMinerFee uint64, sourceScript []byte) (*types.Transaction, error) {
	return _Zkbridge.Contract.Migrate(&_Zkbridge.TransactOpts, redeemAmount, btcMinerFee, sourceScript)
}

// Redeem is a paid mutator transaction binding the contract method 0x74f4897c.
//
// Solidity: function redeem(uint64 redeemAmount, uint64 btcMinerFee, bytes recipientLockScript) returns()
func (_Zkbridge *ZkbridgeTransactor) Redeem(opts *bind.TransactOpts, redeemAmount uint64, btcMinerFee uint64, recipientLockScript []byte) (*types.Transaction, error) {
	return _Zkbridge.contract.Transact(opts, "redeem", redeemAmount, btcMinerFee, recipientLockScript)
}

// Redeem is a paid mutator transaction binding the contract method 0x74f4897c.
//
// Solidity: function redeem(uint64 redeemAmount, uint64 btcMinerFee, bytes recipientLockScript) returns()
func (_Zkbridge *ZkbridgeSession) Redeem(redeemAmount uint64, btcMinerFee uint64, recipientLockScript []byte) (*types.Transaction, error) {
	return _Zkbridge.Contract.Redeem(&_Zkbridge.TransactOpts, redeemAmount, btcMinerFee, recipientLockScript)
}

// Redeem is a paid mutator transaction binding the contract method 0x74f4897c.
//
// Solidity: function redeem(uint64 redeemAmount, uint64 btcMinerFee, bytes recipientLockScript) returns()
func (_Zkbridge *ZkbridgeTransactorSession) Redeem(redeemAmount uint64, btcMinerFee uint64, recipientLockScript []byte) (*types.Transaction, error) {
	return _Zkbridge.Contract.Redeem(&_Zkbridge.TransactOpts, redeemAmount, btcMinerFee, recipientLockScript)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_Zkbridge *ZkbridgeTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _Zkbridge.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_Zkbridge *ZkbridgeSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _Zkbridge.Contract.RenounceRole(&_Zkbridge.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_Zkbridge *ZkbridgeTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _Zkbridge.Contract.RenounceRole(&_Zkbridge.TransactOpts, role, callerConfirmation)
}

// Retrieve is a paid mutator transaction binding the contract method 0xd1ed74ad.
//
// Solidity: function retrieve(bytes32 txid) returns()
func (_Zkbridge *ZkbridgeTransactor) Retrieve(opts *bind.TransactOpts, txid [32]byte) (*types.Transaction, error) {
	return _Zkbridge.contract.Transact(opts, "retrieve", txid)
}

// Retrieve is a paid mutator transaction binding the contract method 0xd1ed74ad.
//
// Solidity: function retrieve(bytes32 txid) returns()
func (_Zkbridge *ZkbridgeSession) Retrieve(txid [32]byte) (*types.Transaction, error) {
	return _Zkbridge.Contract.Retrieve(&_Zkbridge.TransactOpts, txid)
}

// Retrieve is a paid mutator transaction binding the contract method 0xd1ed74ad.
//
// Solidity: function retrieve(bytes32 txid) returns()
func (_Zkbridge *ZkbridgeTransactorSession) Retrieve(txid [32]byte) (*types.Transaction, error) {
	return _Zkbridge.Contract.Retrieve(&_Zkbridge.TransactOpts, txid)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Zkbridge *ZkbridgeTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Zkbridge.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Zkbridge *ZkbridgeSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Zkbridge.Contract.RevokeRole(&_Zkbridge.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Zkbridge *ZkbridgeTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Zkbridge.Contract.RevokeRole(&_Zkbridge.TransactOpts, role, account)
}

// UpdateDustThresholds is a paid mutator transaction binding the contract method 0x4cfe5c16.
//
// Solidity: function updateDustThresholds(uint64 newThresholds) returns()
func (_Zkbridge *ZkbridgeTransactor) UpdateDustThresholds(opts *bind.TransactOpts, newThresholds uint64) (*types.Transaction, error) {
	return _Zkbridge.contract.Transact(opts, "updateDustThresholds", newThresholds)
}

// UpdateDustThresholds is a paid mutator transaction binding the contract method 0x4cfe5c16.
//
// Solidity: function updateDustThresholds(uint64 newThresholds) returns()
func (_Zkbridge *ZkbridgeSession) UpdateDustThresholds(newThresholds uint64) (*types.Transaction, error) {
	return _Zkbridge.Contract.UpdateDustThresholds(&_Zkbridge.TransactOpts, newThresholds)
}

// UpdateDustThresholds is a paid mutator transaction binding the contract method 0x4cfe5c16.
//
// Solidity: function updateDustThresholds(uint64 newThresholds) returns()
func (_Zkbridge *ZkbridgeTransactorSession) UpdateDustThresholds(newThresholds uint64) (*types.Transaction, error) {
	return _Zkbridge.Contract.UpdateDustThresholds(&_Zkbridge.TransactOpts, newThresholds)
}

// UpdateRedeem is a paid mutator transaction binding the contract method 0x40dbb323.
//
// Solidity: function updateRedeem((bytes32,uint32,uint32,bytes32,uint32,address,uint256,uint32) params, bytes32 txid, uint256 minerReward, bytes proofData) returns()
func (_Zkbridge *ZkbridgeTransactor) UpdateRedeem(opts *bind.TransactOpts, params IBtcTxVerifierPublicWitnessParams, txid [32]byte, minerReward *big.Int, proofData []byte) (*types.Transaction, error) {
	return _Zkbridge.contract.Transact(opts, "updateRedeem", params, txid, minerReward, proofData)
}

// UpdateRedeem is a paid mutator transaction binding the contract method 0x40dbb323.
//
// Solidity: function updateRedeem((bytes32,uint32,uint32,bytes32,uint32,address,uint256,uint32) params, bytes32 txid, uint256 minerReward, bytes proofData) returns()
func (_Zkbridge *ZkbridgeSession) UpdateRedeem(params IBtcTxVerifierPublicWitnessParams, txid [32]byte, minerReward *big.Int, proofData []byte) (*types.Transaction, error) {
	return _Zkbridge.Contract.UpdateRedeem(&_Zkbridge.TransactOpts, params, txid, minerReward, proofData)
}

// UpdateRedeem is a paid mutator transaction binding the contract method 0x40dbb323.
//
// Solidity: function updateRedeem((bytes32,uint32,uint32,bytes32,uint32,address,uint256,uint32) params, bytes32 txid, uint256 minerReward, bytes proofData) returns()
func (_Zkbridge *ZkbridgeTransactorSession) UpdateRedeem(params IBtcTxVerifierPublicWitnessParams, txid [32]byte, minerReward *big.Int, proofData []byte) (*types.Transaction, error) {
	return _Zkbridge.Contract.UpdateRedeem(&_Zkbridge.TransactOpts, params, txid, minerReward, proofData)
}

// UpdateReserveInterface is a paid mutator transaction binding the contract method 0x4f2f3df7.
//
// Solidity: function updateReserveInterface(address newInterface) returns()
func (_Zkbridge *ZkbridgeTransactor) UpdateReserveInterface(opts *bind.TransactOpts, newInterface common.Address) (*types.Transaction, error) {
	return _Zkbridge.contract.Transact(opts, "updateReserveInterface", newInterface)
}

// UpdateReserveInterface is a paid mutator transaction binding the contract method 0x4f2f3df7.
//
// Solidity: function updateReserveInterface(address newInterface) returns()
func (_Zkbridge *ZkbridgeSession) UpdateReserveInterface(newInterface common.Address) (*types.Transaction, error) {
	return _Zkbridge.Contract.UpdateReserveInterface(&_Zkbridge.TransactOpts, newInterface)
}

// UpdateReserveInterface is a paid mutator transaction binding the contract method 0x4f2f3df7.
//
// Solidity: function updateReserveInterface(address newInterface) returns()
func (_Zkbridge *ZkbridgeTransactorSession) UpdateReserveInterface(newInterface common.Address) (*types.Transaction, error) {
	return _Zkbridge.Contract.UpdateReserveInterface(&_Zkbridge.TransactOpts, newInterface)
}

// ZkbridgeBeginMigrationIterator is returned from FilterBeginMigration and is used to iterate over the raw logs and unpacked data for BeginMigration events raised by the Zkbridge contract.
type ZkbridgeBeginMigrationIterator struct {
	Event *ZkbridgeBeginMigration // Event containing the contract specifics and raw log

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
func (it *ZkbridgeBeginMigrationIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZkbridgeBeginMigration)
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
		it.Event = new(ZkbridgeBeginMigration)
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
func (it *ZkbridgeBeginMigrationIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZkbridgeBeginMigrationIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZkbridgeBeginMigration represents a BeginMigration event raised by the Zkbridge contract.
type ZkbridgeBeginMigration struct {
	SourceLockingScript      common.Hash
	DestinationLockingScript common.Hash
	Raw                      types.Log // Blockchain specific contextual infos
}

// FilterBeginMigration is a free log retrieval operation binding the contract event 0x451ad61039ac2e35a6f3a9c3e5b21d5b928dd02ca1be6fcc5e99c08a71323103.
//
// Solidity: event BeginMigration(bytes indexed sourceLockingScript, bytes indexed destinationLockingScript)
func (_Zkbridge *ZkbridgeFilterer) FilterBeginMigration(opts *bind.FilterOpts, sourceLockingScript [][]byte, destinationLockingScript [][]byte) (*ZkbridgeBeginMigrationIterator, error) {

	var sourceLockingScriptRule []interface{}
	for _, sourceLockingScriptItem := range sourceLockingScript {
		sourceLockingScriptRule = append(sourceLockingScriptRule, sourceLockingScriptItem)
	}
	var destinationLockingScriptRule []interface{}
	for _, destinationLockingScriptItem := range destinationLockingScript {
		destinationLockingScriptRule = append(destinationLockingScriptRule, destinationLockingScriptItem)
	}

	logs, sub, err := _Zkbridge.contract.FilterLogs(opts, "BeginMigration", sourceLockingScriptRule, destinationLockingScriptRule)
	if err != nil {
		return nil, err
	}
	return &ZkbridgeBeginMigrationIterator{contract: _Zkbridge.contract, event: "BeginMigration", logs: logs, sub: sub}, nil
}

// WatchBeginMigration is a free log subscription operation binding the contract event 0x451ad61039ac2e35a6f3a9c3e5b21d5b928dd02ca1be6fcc5e99c08a71323103.
//
// Solidity: event BeginMigration(bytes indexed sourceLockingScript, bytes indexed destinationLockingScript)
func (_Zkbridge *ZkbridgeFilterer) WatchBeginMigration(opts *bind.WatchOpts, sink chan<- *ZkbridgeBeginMigration, sourceLockingScript [][]byte, destinationLockingScript [][]byte) (event.Subscription, error) {

	var sourceLockingScriptRule []interface{}
	for _, sourceLockingScriptItem := range sourceLockingScript {
		sourceLockingScriptRule = append(sourceLockingScriptRule, sourceLockingScriptItem)
	}
	var destinationLockingScriptRule []interface{}
	for _, destinationLockingScriptItem := range destinationLockingScript {
		destinationLockingScriptRule = append(destinationLockingScriptRule, destinationLockingScriptItem)
	}

	logs, sub, err := _Zkbridge.contract.WatchLogs(opts, "BeginMigration", sourceLockingScriptRule, destinationLockingScriptRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZkbridgeBeginMigration)
				if err := _Zkbridge.contract.UnpackLog(event, "BeginMigration", log); err != nil {
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

// ParseBeginMigration is a log parse operation binding the contract event 0x451ad61039ac2e35a6f3a9c3e5b21d5b928dd02ca1be6fcc5e99c08a71323103.
//
// Solidity: event BeginMigration(bytes indexed sourceLockingScript, bytes indexed destinationLockingScript)
func (_Zkbridge *ZkbridgeFilterer) ParseBeginMigration(log types.Log) (*ZkbridgeBeginMigration, error) {
	event := new(ZkbridgeBeginMigration)
	if err := _Zkbridge.contract.UnpackLog(event, "BeginMigration", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
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
	Txid        [32]byte
	MinerReward *big.Int
	SigHashs    [][32]byte
	RawTx       []byte
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterCreateRedeemUnsignedTx is a free log retrieval operation binding the contract event 0x379299efe6911678ce0f23cfce13a7c61a5b2c1723f583f9217b6ee0887b3ef4.
//
// Solidity: event CreateRedeemUnsignedTx(bytes32 indexed txid, uint256 indexed minerReward, bytes32[] sigHashs, bytes rawTx)
func (_Zkbridge *ZkbridgeFilterer) FilterCreateRedeemUnsignedTx(opts *bind.FilterOpts, txid [][32]byte, minerReward []*big.Int) (*ZkbridgeCreateRedeemUnsignedTxIterator, error) {

	var txidRule []interface{}
	for _, txidItem := range txid {
		txidRule = append(txidRule, txidItem)
	}
	var minerRewardRule []interface{}
	for _, minerRewardItem := range minerReward {
		minerRewardRule = append(minerRewardRule, minerRewardItem)
	}

	logs, sub, err := _Zkbridge.contract.FilterLogs(opts, "CreateRedeemUnsignedTx", txidRule, minerRewardRule)
	if err != nil {
		return nil, err
	}
	return &ZkbridgeCreateRedeemUnsignedTxIterator{contract: _Zkbridge.contract, event: "CreateRedeemUnsignedTx", logs: logs, sub: sub}, nil
}

// WatchCreateRedeemUnsignedTx is a free log subscription operation binding the contract event 0x379299efe6911678ce0f23cfce13a7c61a5b2c1723f583f9217b6ee0887b3ef4.
//
// Solidity: event CreateRedeemUnsignedTx(bytes32 indexed txid, uint256 indexed minerReward, bytes32[] sigHashs, bytes rawTx)
func (_Zkbridge *ZkbridgeFilterer) WatchCreateRedeemUnsignedTx(opts *bind.WatchOpts, sink chan<- *ZkbridgeCreateRedeemUnsignedTx, txid [][32]byte, minerReward []*big.Int) (event.Subscription, error) {

	var txidRule []interface{}
	for _, txidItem := range txid {
		txidRule = append(txidRule, txidItem)
	}
	var minerRewardRule []interface{}
	for _, minerRewardItem := range minerReward {
		minerRewardRule = append(minerRewardRule, minerRewardItem)
	}

	logs, sub, err := _Zkbridge.contract.WatchLogs(opts, "CreateRedeemUnsignedTx", txidRule, minerRewardRule)
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

// ParseCreateRedeemUnsignedTx is a log parse operation binding the contract event 0x379299efe6911678ce0f23cfce13a7c61a5b2c1723f583f9217b6ee0887b3ef4.
//
// Solidity: event CreateRedeemUnsignedTx(bytes32 indexed txid, uint256 indexed minerReward, bytes32[] sigHashs, bytes rawTx)
func (_Zkbridge *ZkbridgeFilterer) ParseCreateRedeemUnsignedTx(log types.Log) (*ZkbridgeCreateRedeemUnsignedTx, error) {
	event := new(ZkbridgeCreateRedeemUnsignedTx)
	if err := _Zkbridge.contract.UnpackLog(event, "CreateRedeemUnsignedTx", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZkbridgeDustThresholdsUpdatedIterator is returned from FilterDustThresholdsUpdated and is used to iterate over the raw logs and unpacked data for DustThresholdsUpdated events raised by the Zkbridge contract.
type ZkbridgeDustThresholdsUpdatedIterator struct {
	Event *ZkbridgeDustThresholdsUpdated // Event containing the contract specifics and raw log

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
func (it *ZkbridgeDustThresholdsUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZkbridgeDustThresholdsUpdated)
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
		it.Event = new(ZkbridgeDustThresholdsUpdated)
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
func (it *ZkbridgeDustThresholdsUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZkbridgeDustThresholdsUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZkbridgeDustThresholdsUpdated represents a DustThresholdsUpdated event raised by the Zkbridge contract.
type ZkbridgeDustThresholdsUpdated struct {
	Old uint64
	Now uint64
	Raw types.Log // Blockchain specific contextual infos
}

// FilterDustThresholdsUpdated is a free log retrieval operation binding the contract event 0x5ff000e3aef2750b281f3b6d0694bd55a08f256595ee736af44cdfabf978f88d.
//
// Solidity: event DustThresholdsUpdated(uint64 old, uint64 now)
func (_Zkbridge *ZkbridgeFilterer) FilterDustThresholdsUpdated(opts *bind.FilterOpts) (*ZkbridgeDustThresholdsUpdatedIterator, error) {

	logs, sub, err := _Zkbridge.contract.FilterLogs(opts, "DustThresholdsUpdated")
	if err != nil {
		return nil, err
	}
	return &ZkbridgeDustThresholdsUpdatedIterator{contract: _Zkbridge.contract, event: "DustThresholdsUpdated", logs: logs, sub: sub}, nil
}

// WatchDustThresholdsUpdated is a free log subscription operation binding the contract event 0x5ff000e3aef2750b281f3b6d0694bd55a08f256595ee736af44cdfabf978f88d.
//
// Solidity: event DustThresholdsUpdated(uint64 old, uint64 now)
func (_Zkbridge *ZkbridgeFilterer) WatchDustThresholdsUpdated(opts *bind.WatchOpts, sink chan<- *ZkbridgeDustThresholdsUpdated) (event.Subscription, error) {

	logs, sub, err := _Zkbridge.contract.WatchLogs(opts, "DustThresholdsUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZkbridgeDustThresholdsUpdated)
				if err := _Zkbridge.contract.UnpackLog(event, "DustThresholdsUpdated", log); err != nil {
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

// ParseDustThresholdsUpdated is a log parse operation binding the contract event 0x5ff000e3aef2750b281f3b6d0694bd55a08f256595ee736af44cdfabf978f88d.
//
// Solidity: event DustThresholdsUpdated(uint64 old, uint64 now)
func (_Zkbridge *ZkbridgeFilterer) ParseDustThresholdsUpdated(log types.Log) (*ZkbridgeDustThresholdsUpdated, error) {
	event := new(ZkbridgeDustThresholdsUpdated)
	if err := _Zkbridge.contract.UnpackLog(event, "DustThresholdsUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZkbridgeRetrievedDataIterator is returned from FilterRetrievedData and is used to iterate over the raw logs and unpacked data for RetrievedData events raised by the Zkbridge contract.
type ZkbridgeRetrievedDataIterator struct {
	Event *ZkbridgeRetrievedData // Event containing the contract specifics and raw log

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
func (it *ZkbridgeRetrievedDataIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZkbridgeRetrievedData)
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
		it.Event = new(ZkbridgeRetrievedData)
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
func (it *ZkbridgeRetrievedDataIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZkbridgeRetrievedDataIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZkbridgeRetrievedData represents a RetrievedData event raised by the Zkbridge contract.
type ZkbridgeRetrievedData struct {
	Txid [32]byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterRetrievedData is a free log retrieval operation binding the contract event 0x7d78d94c375a6d97295caa73300d4a88e21294e9eeacec2aac8b34add3bc5fd0.
//
// Solidity: event RetrievedData(bytes32 txid)
func (_Zkbridge *ZkbridgeFilterer) FilterRetrievedData(opts *bind.FilterOpts) (*ZkbridgeRetrievedDataIterator, error) {

	logs, sub, err := _Zkbridge.contract.FilterLogs(opts, "RetrievedData")
	if err != nil {
		return nil, err
	}
	return &ZkbridgeRetrievedDataIterator{contract: _Zkbridge.contract, event: "RetrievedData", logs: logs, sub: sub}, nil
}

// WatchRetrievedData is a free log subscription operation binding the contract event 0x7d78d94c375a6d97295caa73300d4a88e21294e9eeacec2aac8b34add3bc5fd0.
//
// Solidity: event RetrievedData(bytes32 txid)
func (_Zkbridge *ZkbridgeFilterer) WatchRetrievedData(opts *bind.WatchOpts, sink chan<- *ZkbridgeRetrievedData) (event.Subscription, error) {

	logs, sub, err := _Zkbridge.contract.WatchLogs(opts, "RetrievedData")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZkbridgeRetrievedData)
				if err := _Zkbridge.contract.UnpackLog(event, "RetrievedData", log); err != nil {
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

// ParseRetrievedData is a log parse operation binding the contract event 0x7d78d94c375a6d97295caa73300d4a88e21294e9eeacec2aac8b34add3bc5fd0.
//
// Solidity: event RetrievedData(bytes32 txid)
func (_Zkbridge *ZkbridgeFilterer) ParseRetrievedData(log types.Log) (*ZkbridgeRetrievedData, error) {
	event := new(ZkbridgeRetrievedData)
	if err := _Zkbridge.contract.UnpackLog(event, "RetrievedData", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZkbridgeRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the Zkbridge contract.
type ZkbridgeRoleAdminChangedIterator struct {
	Event *ZkbridgeRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *ZkbridgeRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZkbridgeRoleAdminChanged)
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
		it.Event = new(ZkbridgeRoleAdminChanged)
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
func (it *ZkbridgeRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZkbridgeRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZkbridgeRoleAdminChanged represents a RoleAdminChanged event raised by the Zkbridge contract.
type ZkbridgeRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Zkbridge *ZkbridgeFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*ZkbridgeRoleAdminChangedIterator, error) {

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

	logs, sub, err := _Zkbridge.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &ZkbridgeRoleAdminChangedIterator{contract: _Zkbridge.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Zkbridge *ZkbridgeFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *ZkbridgeRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _Zkbridge.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZkbridgeRoleAdminChanged)
				if err := _Zkbridge.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_Zkbridge *ZkbridgeFilterer) ParseRoleAdminChanged(log types.Log) (*ZkbridgeRoleAdminChanged, error) {
	event := new(ZkbridgeRoleAdminChanged)
	if err := _Zkbridge.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZkbridgeRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the Zkbridge contract.
type ZkbridgeRoleGrantedIterator struct {
	Event *ZkbridgeRoleGranted // Event containing the contract specifics and raw log

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
func (it *ZkbridgeRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZkbridgeRoleGranted)
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
		it.Event = new(ZkbridgeRoleGranted)
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
func (it *ZkbridgeRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZkbridgeRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZkbridgeRoleGranted represents a RoleGranted event raised by the Zkbridge contract.
type ZkbridgeRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Zkbridge *ZkbridgeFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*ZkbridgeRoleGrantedIterator, error) {

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

	logs, sub, err := _Zkbridge.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &ZkbridgeRoleGrantedIterator{contract: _Zkbridge.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Zkbridge *ZkbridgeFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *ZkbridgeRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _Zkbridge.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZkbridgeRoleGranted)
				if err := _Zkbridge.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_Zkbridge *ZkbridgeFilterer) ParseRoleGranted(log types.Log) (*ZkbridgeRoleGranted, error) {
	event := new(ZkbridgeRoleGranted)
	if err := _Zkbridge.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZkbridgeRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the Zkbridge contract.
type ZkbridgeRoleRevokedIterator struct {
	Event *ZkbridgeRoleRevoked // Event containing the contract specifics and raw log

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
func (it *ZkbridgeRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZkbridgeRoleRevoked)
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
		it.Event = new(ZkbridgeRoleRevoked)
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
func (it *ZkbridgeRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZkbridgeRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZkbridgeRoleRevoked represents a RoleRevoked event raised by the Zkbridge contract.
type ZkbridgeRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Zkbridge *ZkbridgeFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*ZkbridgeRoleRevokedIterator, error) {

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

	logs, sub, err := _Zkbridge.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &ZkbridgeRoleRevokedIterator{contract: _Zkbridge.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Zkbridge *ZkbridgeFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *ZkbridgeRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _Zkbridge.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZkbridgeRoleRevoked)
				if err := _Zkbridge.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_Zkbridge *ZkbridgeFilterer) ParseRoleRevoked(log types.Log) (*ZkbridgeRoleRevoked, error) {
	event := new(ZkbridgeRoleRevoked)
	if err := _Zkbridge.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZkbridgeUpdateReserveInterfaceIterator is returned from FilterUpdateReserveInterface and is used to iterate over the raw logs and unpacked data for UpdateReserveInterface events raised by the Zkbridge contract.
type ZkbridgeUpdateReserveInterfaceIterator struct {
	Event *ZkbridgeUpdateReserveInterface // Event containing the contract specifics and raw log

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
func (it *ZkbridgeUpdateReserveInterfaceIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZkbridgeUpdateReserveInterface)
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
		it.Event = new(ZkbridgeUpdateReserveInterface)
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
func (it *ZkbridgeUpdateReserveInterfaceIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZkbridgeUpdateReserveInterfaceIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZkbridgeUpdateReserveInterface represents a UpdateReserveInterface event raised by the Zkbridge contract.
type ZkbridgeUpdateReserveInterface struct {
	Old common.Address
	Now common.Address
	Raw types.Log // Blockchain specific contextual infos
}

// FilterUpdateReserveInterface is a free log retrieval operation binding the contract event 0xb2b5c6f91bbee84404bdb28b8bb2483ebf2872b6255f0bb366a8a228ed7da536.
//
// Solidity: event UpdateReserveInterface(address old, address now)
func (_Zkbridge *ZkbridgeFilterer) FilterUpdateReserveInterface(opts *bind.FilterOpts) (*ZkbridgeUpdateReserveInterfaceIterator, error) {

	logs, sub, err := _Zkbridge.contract.FilterLogs(opts, "UpdateReserveInterface")
	if err != nil {
		return nil, err
	}
	return &ZkbridgeUpdateReserveInterfaceIterator{contract: _Zkbridge.contract, event: "UpdateReserveInterface", logs: logs, sub: sub}, nil
}

// WatchUpdateReserveInterface is a free log subscription operation binding the contract event 0xb2b5c6f91bbee84404bdb28b8bb2483ebf2872b6255f0bb366a8a228ed7da536.
//
// Solidity: event UpdateReserveInterface(address old, address now)
func (_Zkbridge *ZkbridgeFilterer) WatchUpdateReserveInterface(opts *bind.WatchOpts, sink chan<- *ZkbridgeUpdateReserveInterface) (event.Subscription, error) {

	logs, sub, err := _Zkbridge.contract.WatchLogs(opts, "UpdateReserveInterface")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZkbridgeUpdateReserveInterface)
				if err := _Zkbridge.contract.UnpackLog(event, "UpdateReserveInterface", log); err != nil {
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

// ParseUpdateReserveInterface is a log parse operation binding the contract event 0xb2b5c6f91bbee84404bdb28b8bb2483ebf2872b6255f0bb366a8a228ed7da536.
//
// Solidity: event UpdateReserveInterface(address old, address now)
func (_Zkbridge *ZkbridgeFilterer) ParseUpdateReserveInterface(log types.Log) (*ZkbridgeUpdateReserveInterface, error) {
	event := new(ZkbridgeUpdateReserveInterface)
	if err := _Zkbridge.contract.UnpackLog(event, "UpdateReserveInterface", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
