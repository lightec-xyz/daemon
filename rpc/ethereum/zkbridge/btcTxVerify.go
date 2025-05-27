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


// BtcTxVerifierPublicWitnessParams is an auto generated low-level Go binding around an user-defined struct.
type BtcTxVerifierPublicWitnessParams struct {
	Checkpoint        [32]byte
	CpDepth           uint32
	TxDepth           uint32
	TxBlockHash       [32]byte
	TxTimestamp       uint32
	ZkpMiner          common.Address
	Flag              *big.Int
	SmoothedTimestamp uint32
}

// BtcTxVerifyMetaData contains all meta data concerning the BtcTxVerify contract.
var BtcTxVerifyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIBtcZkpVerifier\",\"name\":\"_depositVerifier\",\"type\":\"address\"},{\"internalType\":\"contractIBtcZkpVerifier\",\"name\":\"_redeemVerifier\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_genesisBlockHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_checkpointBlockHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"_checkpointTimestamp\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"signed\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"switchState\",\"type\":\"bool\"}],\"name\":\"UnsignedProtectionStateError\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"oldCp\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"newCp\",\"type\":\"bytes32\"}],\"name\":\"CheckPointRotated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"old\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"now\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"reservedHours\",\"type\":\"uint32\"}],\"name\":\"DepositVerifierUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"old\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"now\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"reservedHours\",\"type\":\"uint32\"}],\"name\":\"RedeemVerifierUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"now\",\"type\":\"bool\"}],\"name\":\"UnsignedProtectionChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"oldDepth\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"newDepth\",\"type\":\"uint8\"}],\"name\":\"UpdatedExtraDepth\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ADD_CP_DEPTH_THRESHOLDS\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"AVERTING_THRESHOLD\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"BRIDGE_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MIN_CP_ADD_INTERVAL\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MIN_CP_DEPTH\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MIN_CP_MATURITY\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ROTATION_THRESHOLD\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"checkpoint\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"cpDepth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"txDepth\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"txBlockHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"txTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"zkpMiner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"flag\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"smoothedTimestamp\",\"type\":\"uint32\"}],\"internalType\":\"structBtcTxVerifier.PublicWitnessParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"name\":\"addNewCP\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"checkpoint\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"cpDepth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"txDepth\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"txBlockHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"txTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"zkpMiner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"flag\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"smoothedTimestamp\",\"type\":\"uint32\"}],\"internalType\":\"structBtcTxVerifier.PublicWitnessParams\",\"name\":\"params\",\"type\":\"tuple\"},{\"internalType\":\"uint32\",\"name\":\"allowance\",\"type\":\"uint32\"}],\"name\":\"checkCpDepth\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"cpCandidates\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"timestamp\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"cpLatestAddedTime\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"rawBtcTx\",\"type\":\"bytes\"}],\"name\":\"decodeTx\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"amount\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"receiveAddress\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"txID\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"toScript\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"opData\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"chainFlag\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"depositCompatibleDeadline\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"depositVerifier\",\"outputs\":[{\"internalType\":\"contractIBtcZkpVerifier\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"enableUnsignedProtection\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"genesisBlockHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"checkpoint\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"cpDepth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"txDepth\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"txBlockHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"txTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"zkpMiner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"flag\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"smoothedTimestamp\",\"type\":\"uint32\"}],\"internalType\":\"structBtcTxVerifier.PublicWitnessParams\",\"name\":\"params\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"}],\"name\":\"getDepositPublicWitness\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"amount\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"raiseIf\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"blockSigned\",\"type\":\"bool\"}],\"name\":\"getDepthByAmount\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"checkpoint\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"cpDepth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"txDepth\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"txBlockHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"txTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"zkpMiner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"flag\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"smoothedTimestamp\",\"type\":\"uint32\"}],\"internalType\":\"structBtcTxVerifier.PublicWitnessParams\",\"name\":\"params\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"minerReward\",\"type\":\"uint256\"}],\"name\":\"getRedeemPublicWitness\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"getRoleMember\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleMemberCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleMembers\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"userAmount\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"changeAmount\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"receiveLockScript\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"changeLockScript\",\"type\":\"bytes\"}],\"name\":\"getTxOuts\",\"outputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"value\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"pkScript\",\"type\":\"bytes\"}],\"internalType\":\"structBtcTxLib.TxOut[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"}],\"name\":\"isCandidateExist\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"juniorCP\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"timestamp\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"previousDepositVerifier\",\"outputs\":[{\"internalType\":\"contractIBtcZkpVerifier\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"previousRedeemVerifier\",\"outputs\":[{\"internalType\":\"contractIBtcZkpVerifier\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"redeemCompatibleDeadline\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"redeemVerifier\",\"outputs\":[{\"internalType\":\"contractIBtcZkpVerifier\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"seniorCP\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"timestamp\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"suggestedCP\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tryRotateCP\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIBtcZkpVerifier\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"reserveHours\",\"type\":\"uint32\"}],\"name\":\"updateDepositVerifier\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"depth\",\"type\":\"uint8\"}],\"name\":\"updateExtraDepth\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIBtcZkpVerifier\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"reserveHours\",\"type\":\"uint32\"}],\"name\":\"updateRedeemVerifier\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"flag\",\"type\":\"bool\"}],\"name\":\"updateUnsignedProtection\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"checkpoint\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"cpDepth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"txDepth\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"txBlockHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"txTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"zkpMiner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"flag\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"smoothedTimestamp\",\"type\":\"uint32\"}],\"internalType\":\"structBtcTxVerifier.PublicWitnessParams\",\"name\":\"params\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"proofData\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"raiseIf\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"amount\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"}],\"name\":\"verifyDepositTx\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"checkpoint\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"cpDepth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"txDepth\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"txBlockHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"txTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"zkpMiner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"flag\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"smoothedTimestamp\",\"type\":\"uint32\"}],\"internalType\":\"structBtcTxVerifier.PublicWitnessParams\",\"name\":\"params\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"txid\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"changeAmount\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"minerReward\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"proofData\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"raiseIf\",\"type\":\"bool\"}],\"name\":\"verifyRedeemTx\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// BtcTxVerifyABI is the input ABI used to generate the binding from.
// Deprecated: Use BtcTxVerifyMetaData.ABI instead.
var BtcTxVerifyABI = BtcTxVerifyMetaData.ABI

// BtcTxVerify is an auto generated Go binding around an Ethereum contract.
type BtcTxVerify struct {
	BtcTxVerifyCaller     // Read-only binding to the contract
	BtcTxVerifyTransactor // Write-only binding to the contract
	BtcTxVerifyFilterer   // Log filterer for contract events
}

// BtcTxVerifyCaller is an auto generated read-only Go binding around an Ethereum contract.
type BtcTxVerifyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BtcTxVerifyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BtcTxVerifyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BtcTxVerifyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BtcTxVerifyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BtcTxVerifySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BtcTxVerifySession struct {
	Contract     *BtcTxVerify      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BtcTxVerifyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BtcTxVerifyCallerSession struct {
	Contract *BtcTxVerifyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// BtcTxVerifyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BtcTxVerifyTransactorSession struct {
	Contract     *BtcTxVerifyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// BtcTxVerifyRaw is an auto generated low-level Go binding around an Ethereum contract.
type BtcTxVerifyRaw struct {
	Contract *BtcTxVerify // Generic contract binding to access the raw methods on
}

// BtcTxVerifyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BtcTxVerifyCallerRaw struct {
	Contract *BtcTxVerifyCaller // Generic read-only contract binding to access the raw methods on
}

// BtcTxVerifyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BtcTxVerifyTransactorRaw struct {
	Contract *BtcTxVerifyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBtcTxVerify creates a new instance of BtcTxVerify, bound to a specific deployed contract.
func NewBtcTxVerify(address common.Address, backend bind.ContractBackend) (*BtcTxVerify, error) {
	contract, err := bindBtcTxVerify(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BtcTxVerify{BtcTxVerifyCaller: BtcTxVerifyCaller{contract: contract}, BtcTxVerifyTransactor: BtcTxVerifyTransactor{contract: contract}, BtcTxVerifyFilterer: BtcTxVerifyFilterer{contract: contract}}, nil
}

// NewBtcTxVerifyCaller creates a new read-only instance of BtcTxVerify, bound to a specific deployed contract.
func NewBtcTxVerifyCaller(address common.Address, caller bind.ContractCaller) (*BtcTxVerifyCaller, error) {
	contract, err := bindBtcTxVerify(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BtcTxVerifyCaller{contract: contract}, nil
}

// NewBtcTxVerifyTransactor creates a new write-only instance of BtcTxVerify, bound to a specific deployed contract.
func NewBtcTxVerifyTransactor(address common.Address, transactor bind.ContractTransactor) (*BtcTxVerifyTransactor, error) {
	contract, err := bindBtcTxVerify(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BtcTxVerifyTransactor{contract: contract}, nil
}

// NewBtcTxVerifyFilterer creates a new log filterer instance of BtcTxVerify, bound to a specific deployed contract.
func NewBtcTxVerifyFilterer(address common.Address, filterer bind.ContractFilterer) (*BtcTxVerifyFilterer, error) {
	contract, err := bindBtcTxVerify(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BtcTxVerifyFilterer{contract: contract}, nil
}

// bindBtcTxVerify binds a generic wrapper to an already deployed contract.
func bindBtcTxVerify(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BtcTxVerifyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BtcTxVerify *BtcTxVerifyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BtcTxVerify.Contract.BtcTxVerifyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BtcTxVerify *BtcTxVerifyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BtcTxVerify.Contract.BtcTxVerifyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BtcTxVerify *BtcTxVerifyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BtcTxVerify.Contract.BtcTxVerifyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BtcTxVerify *BtcTxVerifyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BtcTxVerify.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BtcTxVerify *BtcTxVerifyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BtcTxVerify.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BtcTxVerify *BtcTxVerifyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BtcTxVerify.Contract.contract.Transact(opts, method, params...)
}

// ADDCPDEPTHTHRESHOLDS is a free data retrieval call binding the contract method 0xe97cc078.
//
// Solidity: function ADD_CP_DEPTH_THRESHOLDS() view returns(uint8)
func (_BtcTxVerify *BtcTxVerifyCaller) ADDCPDEPTHTHRESHOLDS(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _BtcTxVerify.contract.Call(opts, &out, "ADD_CP_DEPTH_THRESHOLDS")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// ADDCPDEPTHTHRESHOLDS is a free data retrieval call binding the contract method 0xe97cc078.
//
// Solidity: function ADD_CP_DEPTH_THRESHOLDS() view returns(uint8)
func (_BtcTxVerify *BtcTxVerifySession) ADDCPDEPTHTHRESHOLDS() (uint8, error) {
	return _BtcTxVerify.Contract.ADDCPDEPTHTHRESHOLDS(&_BtcTxVerify.CallOpts)
}

// ADDCPDEPTHTHRESHOLDS is a free data retrieval call binding the contract method 0xe97cc078.
//
// Solidity: function ADD_CP_DEPTH_THRESHOLDS() view returns(uint8)
func (_BtcTxVerify *BtcTxVerifyCallerSession) ADDCPDEPTHTHRESHOLDS() (uint8, error) {
	return _BtcTxVerify.Contract.ADDCPDEPTHTHRESHOLDS(&_BtcTxVerify.CallOpts)
}

// AVERTINGTHRESHOLD is a free data retrieval call binding the contract method 0x65d528e9.
//
// Solidity: function AVERTING_THRESHOLD() view returns(uint32)
func (_BtcTxVerify *BtcTxVerifyCaller) AVERTINGTHRESHOLD(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _BtcTxVerify.contract.Call(opts, &out, "AVERTING_THRESHOLD")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// AVERTINGTHRESHOLD is a free data retrieval call binding the contract method 0x65d528e9.
//
// Solidity: function AVERTING_THRESHOLD() view returns(uint32)
func (_BtcTxVerify *BtcTxVerifySession) AVERTINGTHRESHOLD() (uint32, error) {
	return _BtcTxVerify.Contract.AVERTINGTHRESHOLD(&_BtcTxVerify.CallOpts)
}

// AVERTINGTHRESHOLD is a free data retrieval call binding the contract method 0x65d528e9.
//
// Solidity: function AVERTING_THRESHOLD() view returns(uint32)
func (_BtcTxVerify *BtcTxVerifyCallerSession) AVERTINGTHRESHOLD() (uint32, error) {
	return _BtcTxVerify.Contract.AVERTINGTHRESHOLD(&_BtcTxVerify.CallOpts)
}

// BRIDGEROLE is a free data retrieval call binding the contract method 0xb5bfddea.
//
// Solidity: function BRIDGE_ROLE() view returns(bytes32)
func (_BtcTxVerify *BtcTxVerifyCaller) BRIDGEROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _BtcTxVerify.contract.Call(opts, &out, "BRIDGE_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// BRIDGEROLE is a free data retrieval call binding the contract method 0xb5bfddea.
//
// Solidity: function BRIDGE_ROLE() view returns(bytes32)
func (_BtcTxVerify *BtcTxVerifySession) BRIDGEROLE() ([32]byte, error) {
	return _BtcTxVerify.Contract.BRIDGEROLE(&_BtcTxVerify.CallOpts)
}

// BRIDGEROLE is a free data retrieval call binding the contract method 0xb5bfddea.
//
// Solidity: function BRIDGE_ROLE() view returns(bytes32)
func (_BtcTxVerify *BtcTxVerifyCallerSession) BRIDGEROLE() ([32]byte, error) {
	return _BtcTxVerify.Contract.BRIDGEROLE(&_BtcTxVerify.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_BtcTxVerify *BtcTxVerifyCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _BtcTxVerify.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_BtcTxVerify *BtcTxVerifySession) DEFAULTADMINROLE() ([32]byte, error) {
	return _BtcTxVerify.Contract.DEFAULTADMINROLE(&_BtcTxVerify.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_BtcTxVerify *BtcTxVerifyCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _BtcTxVerify.Contract.DEFAULTADMINROLE(&_BtcTxVerify.CallOpts)
}

// MINCPADDINTERVAL is a free data retrieval call binding the contract method 0x6be273a1.
//
// Solidity: function MIN_CP_ADD_INTERVAL() view returns(uint64)
func (_BtcTxVerify *BtcTxVerifyCaller) MINCPADDINTERVAL(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _BtcTxVerify.contract.Call(opts, &out, "MIN_CP_ADD_INTERVAL")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// MINCPADDINTERVAL is a free data retrieval call binding the contract method 0x6be273a1.
//
// Solidity: function MIN_CP_ADD_INTERVAL() view returns(uint64)
func (_BtcTxVerify *BtcTxVerifySession) MINCPADDINTERVAL() (uint64, error) {
	return _BtcTxVerify.Contract.MINCPADDINTERVAL(&_BtcTxVerify.CallOpts)
}

// MINCPADDINTERVAL is a free data retrieval call binding the contract method 0x6be273a1.
//
// Solidity: function MIN_CP_ADD_INTERVAL() view returns(uint64)
func (_BtcTxVerify *BtcTxVerifyCallerSession) MINCPADDINTERVAL() (uint64, error) {
	return _BtcTxVerify.Contract.MINCPADDINTERVAL(&_BtcTxVerify.CallOpts)
}

// MINCPDEPTH is a free data retrieval call binding the contract method 0xd52562fb.
//
// Solidity: function MIN_CP_DEPTH() view returns(uint8)
func (_BtcTxVerify *BtcTxVerifyCaller) MINCPDEPTH(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _BtcTxVerify.contract.Call(opts, &out, "MIN_CP_DEPTH")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// MINCPDEPTH is a free data retrieval call binding the contract method 0xd52562fb.
//
// Solidity: function MIN_CP_DEPTH() view returns(uint8)
func (_BtcTxVerify *BtcTxVerifySession) MINCPDEPTH() (uint8, error) {
	return _BtcTxVerify.Contract.MINCPDEPTH(&_BtcTxVerify.CallOpts)
}

// MINCPDEPTH is a free data retrieval call binding the contract method 0xd52562fb.
//
// Solidity: function MIN_CP_DEPTH() view returns(uint8)
func (_BtcTxVerify *BtcTxVerifyCallerSession) MINCPDEPTH() (uint8, error) {
	return _BtcTxVerify.Contract.MINCPDEPTH(&_BtcTxVerify.CallOpts)
}

// MINCPMATURITY is a free data retrieval call binding the contract method 0x279f7a66.
//
// Solidity: function MIN_CP_MATURITY() view returns(uint32)
func (_BtcTxVerify *BtcTxVerifyCaller) MINCPMATURITY(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _BtcTxVerify.contract.Call(opts, &out, "MIN_CP_MATURITY")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// MINCPMATURITY is a free data retrieval call binding the contract method 0x279f7a66.
//
// Solidity: function MIN_CP_MATURITY() view returns(uint32)
func (_BtcTxVerify *BtcTxVerifySession) MINCPMATURITY() (uint32, error) {
	return _BtcTxVerify.Contract.MINCPMATURITY(&_BtcTxVerify.CallOpts)
}

// MINCPMATURITY is a free data retrieval call binding the contract method 0x279f7a66.
//
// Solidity: function MIN_CP_MATURITY() view returns(uint32)
func (_BtcTxVerify *BtcTxVerifyCallerSession) MINCPMATURITY() (uint32, error) {
	return _BtcTxVerify.Contract.MINCPMATURITY(&_BtcTxVerify.CallOpts)
}

// ROTATIONTHRESHOLD is a free data retrieval call binding the contract method 0xce840ac5.
//
// Solidity: function ROTATION_THRESHOLD() view returns(uint32)
func (_BtcTxVerify *BtcTxVerifyCaller) ROTATIONTHRESHOLD(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _BtcTxVerify.contract.Call(opts, &out, "ROTATION_THRESHOLD")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// ROTATIONTHRESHOLD is a free data retrieval call binding the contract method 0xce840ac5.
//
// Solidity: function ROTATION_THRESHOLD() view returns(uint32)
func (_BtcTxVerify *BtcTxVerifySession) ROTATIONTHRESHOLD() (uint32, error) {
	return _BtcTxVerify.Contract.ROTATIONTHRESHOLD(&_BtcTxVerify.CallOpts)
}

// ROTATIONTHRESHOLD is a free data retrieval call binding the contract method 0xce840ac5.
//
// Solidity: function ROTATION_THRESHOLD() view returns(uint32)
func (_BtcTxVerify *BtcTxVerifyCallerSession) ROTATIONTHRESHOLD() (uint32, error) {
	return _BtcTxVerify.Contract.ROTATIONTHRESHOLD(&_BtcTxVerify.CallOpts)
}

// CheckCpDepth is a free data retrieval call binding the contract method 0xd4b405ff.
//
// Solidity: function checkCpDepth((bytes32,uint32,uint32,bytes32,uint32,address,uint256,uint32) params, uint32 allowance) view returns(bool)
func (_BtcTxVerify *BtcTxVerifyCaller) CheckCpDepth(opts *bind.CallOpts, params BtcTxVerifierPublicWitnessParams, allowance uint32) (bool, error) {
	var out []interface{}
	err := _BtcTxVerify.contract.Call(opts, &out, "checkCpDepth", params, allowance)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CheckCpDepth is a free data retrieval call binding the contract method 0xd4b405ff.
//
// Solidity: function checkCpDepth((bytes32,uint32,uint32,bytes32,uint32,address,uint256,uint32) params, uint32 allowance) view returns(bool)
func (_BtcTxVerify *BtcTxVerifySession) CheckCpDepth(params BtcTxVerifierPublicWitnessParams, allowance uint32) (bool, error) {
	return _BtcTxVerify.Contract.CheckCpDepth(&_BtcTxVerify.CallOpts, params, allowance)
}

// CheckCpDepth is a free data retrieval call binding the contract method 0xd4b405ff.
//
// Solidity: function checkCpDepth((bytes32,uint32,uint32,bytes32,uint32,address,uint256,uint32) params, uint32 allowance) view returns(bool)
func (_BtcTxVerify *BtcTxVerifyCallerSession) CheckCpDepth(params BtcTxVerifierPublicWitnessParams, allowance uint32) (bool, error) {
	return _BtcTxVerify.Contract.CheckCpDepth(&_BtcTxVerify.CallOpts, params, allowance)
}

// CpCandidates is a free data retrieval call binding the contract method 0x594cbe7f.
//
// Solidity: function cpCandidates(uint256 ) view returns(bytes32 blockHash, uint64 timestamp)
func (_BtcTxVerify *BtcTxVerifyCaller) CpCandidates(opts *bind.CallOpts, arg0 *big.Int) (struct {
	BlockHash [32]byte
	Timestamp uint64
}, error) {
	var out []interface{}
	err := _BtcTxVerify.contract.Call(opts, &out, "cpCandidates", arg0)

	outstruct := new(struct {
		BlockHash [32]byte
		Timestamp uint64
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.BlockHash = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.Timestamp = *abi.ConvertType(out[1], new(uint64)).(*uint64)

	return *outstruct, err

}

// CpCandidates is a free data retrieval call binding the contract method 0x594cbe7f.
//
// Solidity: function cpCandidates(uint256 ) view returns(bytes32 blockHash, uint64 timestamp)
func (_BtcTxVerify *BtcTxVerifySession) CpCandidates(arg0 *big.Int) (struct {
	BlockHash [32]byte
	Timestamp uint64
}, error) {
	return _BtcTxVerify.Contract.CpCandidates(&_BtcTxVerify.CallOpts, arg0)
}

// CpCandidates is a free data retrieval call binding the contract method 0x594cbe7f.
//
// Solidity: function cpCandidates(uint256 ) view returns(bytes32 blockHash, uint64 timestamp)
func (_BtcTxVerify *BtcTxVerifyCallerSession) CpCandidates(arg0 *big.Int) (struct {
	BlockHash [32]byte
	Timestamp uint64
}, error) {
	return _BtcTxVerify.Contract.CpCandidates(&_BtcTxVerify.CallOpts, arg0)
}

// CpLatestAddedTime is a free data retrieval call binding the contract method 0x1e821765.
//
// Solidity: function cpLatestAddedTime() view returns(uint64)
func (_BtcTxVerify *BtcTxVerifyCaller) CpLatestAddedTime(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _BtcTxVerify.contract.Call(opts, &out, "cpLatestAddedTime")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// CpLatestAddedTime is a free data retrieval call binding the contract method 0x1e821765.
//
// Solidity: function cpLatestAddedTime() view returns(uint64)
func (_BtcTxVerify *BtcTxVerifySession) CpLatestAddedTime() (uint64, error) {
	return _BtcTxVerify.Contract.CpLatestAddedTime(&_BtcTxVerify.CallOpts)
}

// CpLatestAddedTime is a free data retrieval call binding the contract method 0x1e821765.
//
// Solidity: function cpLatestAddedTime() view returns(uint64)
func (_BtcTxVerify *BtcTxVerifyCallerSession) CpLatestAddedTime() (uint64, error) {
	return _BtcTxVerify.Contract.CpLatestAddedTime(&_BtcTxVerify.CallOpts)
}

// DecodeTx is a free data retrieval call binding the contract method 0xdae029d3.
//
// Solidity: function decodeTx(bytes rawBtcTx) pure returns(uint32 index, uint64 amount, address receiveAddress, bytes32 txID, bytes toScript, bytes opData, bool chainFlag)
func (_BtcTxVerify *BtcTxVerifyCaller) DecodeTx(opts *bind.CallOpts, rawBtcTx []byte) (struct {
	Index          uint32
	Amount         uint64
	ReceiveAddress common.Address
	TxID           [32]byte
	ToScript       []byte
	OpData         []byte
	ChainFlag      bool
}, error) {
	var out []interface{}
	err := _BtcTxVerify.contract.Call(opts, &out, "decodeTx", rawBtcTx)

	outstruct := new(struct {
		Index          uint32
		Amount         uint64
		ReceiveAddress common.Address
		TxID           [32]byte
		ToScript       []byte
		OpData         []byte
		ChainFlag      bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Index = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.Amount = *abi.ConvertType(out[1], new(uint64)).(*uint64)
	outstruct.ReceiveAddress = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.TxID = *abi.ConvertType(out[3], new([32]byte)).(*[32]byte)
	outstruct.ToScript = *abi.ConvertType(out[4], new([]byte)).(*[]byte)
	outstruct.OpData = *abi.ConvertType(out[5], new([]byte)).(*[]byte)
	outstruct.ChainFlag = *abi.ConvertType(out[6], new(bool)).(*bool)

	return *outstruct, err

}

// DecodeTx is a free data retrieval call binding the contract method 0xdae029d3.
//
// Solidity: function decodeTx(bytes rawBtcTx) pure returns(uint32 index, uint64 amount, address receiveAddress, bytes32 txID, bytes toScript, bytes opData, bool chainFlag)
func (_BtcTxVerify *BtcTxVerifySession) DecodeTx(rawBtcTx []byte) (struct {
	Index          uint32
	Amount         uint64
	ReceiveAddress common.Address
	TxID           [32]byte
	ToScript       []byte
	OpData         []byte
	ChainFlag      bool
}, error) {
	return _BtcTxVerify.Contract.DecodeTx(&_BtcTxVerify.CallOpts, rawBtcTx)
}

// DecodeTx is a free data retrieval call binding the contract method 0xdae029d3.
//
// Solidity: function decodeTx(bytes rawBtcTx) pure returns(uint32 index, uint64 amount, address receiveAddress, bytes32 txID, bytes toScript, bytes opData, bool chainFlag)
func (_BtcTxVerify *BtcTxVerifyCallerSession) DecodeTx(rawBtcTx []byte) (struct {
	Index          uint32
	Amount         uint64
	ReceiveAddress common.Address
	TxID           [32]byte
	ToScript       []byte
	OpData         []byte
	ChainFlag      bool
}, error) {
	return _BtcTxVerify.Contract.DecodeTx(&_BtcTxVerify.CallOpts, rawBtcTx)
}

// DepositCompatibleDeadline is a free data retrieval call binding the contract method 0xbd6e8e94.
//
// Solidity: function depositCompatibleDeadline() view returns(uint256)
func (_BtcTxVerify *BtcTxVerifyCaller) DepositCompatibleDeadline(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BtcTxVerify.contract.Call(opts, &out, "depositCompatibleDeadline")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DepositCompatibleDeadline is a free data retrieval call binding the contract method 0xbd6e8e94.
//
// Solidity: function depositCompatibleDeadline() view returns(uint256)
func (_BtcTxVerify *BtcTxVerifySession) DepositCompatibleDeadline() (*big.Int, error) {
	return _BtcTxVerify.Contract.DepositCompatibleDeadline(&_BtcTxVerify.CallOpts)
}

// DepositCompatibleDeadline is a free data retrieval call binding the contract method 0xbd6e8e94.
//
// Solidity: function depositCompatibleDeadline() view returns(uint256)
func (_BtcTxVerify *BtcTxVerifyCallerSession) DepositCompatibleDeadline() (*big.Int, error) {
	return _BtcTxVerify.Contract.DepositCompatibleDeadline(&_BtcTxVerify.CallOpts)
}

// DepositVerifier is a free data retrieval call binding the contract method 0x26120c88.
//
// Solidity: function depositVerifier() view returns(address)
func (_BtcTxVerify *BtcTxVerifyCaller) DepositVerifier(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BtcTxVerify.contract.Call(opts, &out, "depositVerifier")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DepositVerifier is a free data retrieval call binding the contract method 0x26120c88.
//
// Solidity: function depositVerifier() view returns(address)
func (_BtcTxVerify *BtcTxVerifySession) DepositVerifier() (common.Address, error) {
	return _BtcTxVerify.Contract.DepositVerifier(&_BtcTxVerify.CallOpts)
}

// DepositVerifier is a free data retrieval call binding the contract method 0x26120c88.
//
// Solidity: function depositVerifier() view returns(address)
func (_BtcTxVerify *BtcTxVerifyCallerSession) DepositVerifier() (common.Address, error) {
	return _BtcTxVerify.Contract.DepositVerifier(&_BtcTxVerify.CallOpts)
}

// EnableUnsignedProtection is a free data retrieval call binding the contract method 0x80142d74.
//
// Solidity: function enableUnsignedProtection() view returns(bool)
func (_BtcTxVerify *BtcTxVerifyCaller) EnableUnsignedProtection(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _BtcTxVerify.contract.Call(opts, &out, "enableUnsignedProtection")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// EnableUnsignedProtection is a free data retrieval call binding the contract method 0x80142d74.
//
// Solidity: function enableUnsignedProtection() view returns(bool)
func (_BtcTxVerify *BtcTxVerifySession) EnableUnsignedProtection() (bool, error) {
	return _BtcTxVerify.Contract.EnableUnsignedProtection(&_BtcTxVerify.CallOpts)
}

// EnableUnsignedProtection is a free data retrieval call binding the contract method 0x80142d74.
//
// Solidity: function enableUnsignedProtection() view returns(bool)
func (_BtcTxVerify *BtcTxVerifyCallerSession) EnableUnsignedProtection() (bool, error) {
	return _BtcTxVerify.Contract.EnableUnsignedProtection(&_BtcTxVerify.CallOpts)
}

// GenesisBlockHash is a free data retrieval call binding the contract method 0x28e24b3d.
//
// Solidity: function genesisBlockHash() view returns(bytes32)
func (_BtcTxVerify *BtcTxVerifyCaller) GenesisBlockHash(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _BtcTxVerify.contract.Call(opts, &out, "genesisBlockHash")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GenesisBlockHash is a free data retrieval call binding the contract method 0x28e24b3d.
//
// Solidity: function genesisBlockHash() view returns(bytes32)
func (_BtcTxVerify *BtcTxVerifySession) GenesisBlockHash() ([32]byte, error) {
	return _BtcTxVerify.Contract.GenesisBlockHash(&_BtcTxVerify.CallOpts)
}

// GenesisBlockHash is a free data retrieval call binding the contract method 0x28e24b3d.
//
// Solidity: function genesisBlockHash() view returns(bytes32)
func (_BtcTxVerify *BtcTxVerifyCallerSession) GenesisBlockHash() ([32]byte, error) {
	return _BtcTxVerify.Contract.GenesisBlockHash(&_BtcTxVerify.CallOpts)
}

// GetDepositPublicWitness is a free data retrieval call binding the contract method 0x0957c471.
//
// Solidity: function getDepositPublicWitness((bytes32,uint32,uint32,bytes32,uint32,address,uint256,uint32) params, bytes32 txid) view returns(uint256[])
func (_BtcTxVerify *BtcTxVerifyCaller) GetDepositPublicWitness(opts *bind.CallOpts, params BtcTxVerifierPublicWitnessParams, txid [32]byte) ([]*big.Int, error) {
	var out []interface{}
	err := _BtcTxVerify.contract.Call(opts, &out, "getDepositPublicWitness", params, txid)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetDepositPublicWitness is a free data retrieval call binding the contract method 0x0957c471.
//
// Solidity: function getDepositPublicWitness((bytes32,uint32,uint32,bytes32,uint32,address,uint256,uint32) params, bytes32 txid) view returns(uint256[])
func (_BtcTxVerify *BtcTxVerifySession) GetDepositPublicWitness(params BtcTxVerifierPublicWitnessParams, txid [32]byte) ([]*big.Int, error) {
	return _BtcTxVerify.Contract.GetDepositPublicWitness(&_BtcTxVerify.CallOpts, params, txid)
}

// GetDepositPublicWitness is a free data retrieval call binding the contract method 0x0957c471.
//
// Solidity: function getDepositPublicWitness((bytes32,uint32,uint32,bytes32,uint32,address,uint256,uint32) params, bytes32 txid) view returns(uint256[])
func (_BtcTxVerify *BtcTxVerifyCallerSession) GetDepositPublicWitness(params BtcTxVerifierPublicWitnessParams, txid [32]byte) ([]*big.Int, error) {
	return _BtcTxVerify.Contract.GetDepositPublicWitness(&_BtcTxVerify.CallOpts, params, txid)
}

// GetDepthByAmount is a free data retrieval call binding the contract method 0x3fba8f27.
//
// Solidity: function getDepthByAmount(uint64 amount, bool raiseIf, bool blockSigned) view returns(uint32, uint32)
func (_BtcTxVerify *BtcTxVerifyCaller) GetDepthByAmount(opts *bind.CallOpts, amount uint64, raiseIf bool, blockSigned bool) (uint32, uint32, error) {
	var out []interface{}
	err := _BtcTxVerify.contract.Call(opts, &out, "getDepthByAmount", amount, raiseIf, blockSigned)

	if err != nil {
		return *new(uint32), *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)
	out1 := *abi.ConvertType(out[1], new(uint32)).(*uint32)

	return out0, out1, err

}

// GetDepthByAmount is a free data retrieval call binding the contract method 0x3fba8f27.
//
// Solidity: function getDepthByAmount(uint64 amount, bool raiseIf, bool blockSigned) view returns(uint32, uint32)
func (_BtcTxVerify *BtcTxVerifySession) GetDepthByAmount(amount uint64, raiseIf bool, blockSigned bool) (uint32, uint32, error) {
	return _BtcTxVerify.Contract.GetDepthByAmount(&_BtcTxVerify.CallOpts, amount, raiseIf, blockSigned)
}

// GetDepthByAmount is a free data retrieval call binding the contract method 0x3fba8f27.
//
// Solidity: function getDepthByAmount(uint64 amount, bool raiseIf, bool blockSigned) view returns(uint32, uint32)
func (_BtcTxVerify *BtcTxVerifyCallerSession) GetDepthByAmount(amount uint64, raiseIf bool, blockSigned bool) (uint32, uint32, error) {
	return _BtcTxVerify.Contract.GetDepthByAmount(&_BtcTxVerify.CallOpts, amount, raiseIf, blockSigned)
}

// GetRedeemPublicWitness is a free data retrieval call binding the contract method 0x80581666.
//
// Solidity: function getRedeemPublicWitness((bytes32,uint32,uint32,bytes32,uint32,address,uint256,uint32) params, bytes32 txid, uint256 minerReward) view returns(uint256[])
func (_BtcTxVerify *BtcTxVerifyCaller) GetRedeemPublicWitness(opts *bind.CallOpts, params BtcTxVerifierPublicWitnessParams, txid [32]byte, minerReward *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _BtcTxVerify.contract.Call(opts, &out, "getRedeemPublicWitness", params, txid, minerReward)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetRedeemPublicWitness is a free data retrieval call binding the contract method 0x80581666.
//
// Solidity: function getRedeemPublicWitness((bytes32,uint32,uint32,bytes32,uint32,address,uint256,uint32) params, bytes32 txid, uint256 minerReward) view returns(uint256[])
func (_BtcTxVerify *BtcTxVerifySession) GetRedeemPublicWitness(params BtcTxVerifierPublicWitnessParams, txid [32]byte, minerReward *big.Int) ([]*big.Int, error) {
	return _BtcTxVerify.Contract.GetRedeemPublicWitness(&_BtcTxVerify.CallOpts, params, txid, minerReward)
}

// GetRedeemPublicWitness is a free data retrieval call binding the contract method 0x80581666.
//
// Solidity: function getRedeemPublicWitness((bytes32,uint32,uint32,bytes32,uint32,address,uint256,uint32) params, bytes32 txid, uint256 minerReward) view returns(uint256[])
func (_BtcTxVerify *BtcTxVerifyCallerSession) GetRedeemPublicWitness(params BtcTxVerifierPublicWitnessParams, txid [32]byte, minerReward *big.Int) ([]*big.Int, error) {
	return _BtcTxVerify.Contract.GetRedeemPublicWitness(&_BtcTxVerify.CallOpts, params, txid, minerReward)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_BtcTxVerify *BtcTxVerifyCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _BtcTxVerify.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_BtcTxVerify *BtcTxVerifySession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _BtcTxVerify.Contract.GetRoleAdmin(&_BtcTxVerify.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_BtcTxVerify *BtcTxVerifyCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _BtcTxVerify.Contract.GetRoleAdmin(&_BtcTxVerify.CallOpts, role)
}

// GetRoleMember is a free data retrieval call binding the contract method 0x9010d07c.
//
// Solidity: function getRoleMember(bytes32 role, uint256 index) view returns(address)
func (_BtcTxVerify *BtcTxVerifyCaller) GetRoleMember(opts *bind.CallOpts, role [32]byte, index *big.Int) (common.Address, error) {
	var out []interface{}
	err := _BtcTxVerify.contract.Call(opts, &out, "getRoleMember", role, index)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetRoleMember is a free data retrieval call binding the contract method 0x9010d07c.
//
// Solidity: function getRoleMember(bytes32 role, uint256 index) view returns(address)
func (_BtcTxVerify *BtcTxVerifySession) GetRoleMember(role [32]byte, index *big.Int) (common.Address, error) {
	return _BtcTxVerify.Contract.GetRoleMember(&_BtcTxVerify.CallOpts, role, index)
}

// GetRoleMember is a free data retrieval call binding the contract method 0x9010d07c.
//
// Solidity: function getRoleMember(bytes32 role, uint256 index) view returns(address)
func (_BtcTxVerify *BtcTxVerifyCallerSession) GetRoleMember(role [32]byte, index *big.Int) (common.Address, error) {
	return _BtcTxVerify.Contract.GetRoleMember(&_BtcTxVerify.CallOpts, role, index)
}

// GetRoleMemberCount is a free data retrieval call binding the contract method 0xca15c873.
//
// Solidity: function getRoleMemberCount(bytes32 role) view returns(uint256)
func (_BtcTxVerify *BtcTxVerifyCaller) GetRoleMemberCount(opts *bind.CallOpts, role [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _BtcTxVerify.contract.Call(opts, &out, "getRoleMemberCount", role)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRoleMemberCount is a free data retrieval call binding the contract method 0xca15c873.
//
// Solidity: function getRoleMemberCount(bytes32 role) view returns(uint256)
func (_BtcTxVerify *BtcTxVerifySession) GetRoleMemberCount(role [32]byte) (*big.Int, error) {
	return _BtcTxVerify.Contract.GetRoleMemberCount(&_BtcTxVerify.CallOpts, role)
}

// GetRoleMemberCount is a free data retrieval call binding the contract method 0xca15c873.
//
// Solidity: function getRoleMemberCount(bytes32 role) view returns(uint256)
func (_BtcTxVerify *BtcTxVerifyCallerSession) GetRoleMemberCount(role [32]byte) (*big.Int, error) {
	return _BtcTxVerify.Contract.GetRoleMemberCount(&_BtcTxVerify.CallOpts, role)
}

// GetRoleMembers is a free data retrieval call binding the contract method 0xa3246ad3.
//
// Solidity: function getRoleMembers(bytes32 role) view returns(address[])
func (_BtcTxVerify *BtcTxVerifyCaller) GetRoleMembers(opts *bind.CallOpts, role [32]byte) ([]common.Address, error) {
	var out []interface{}
	err := _BtcTxVerify.contract.Call(opts, &out, "getRoleMembers", role)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetRoleMembers is a free data retrieval call binding the contract method 0xa3246ad3.
//
// Solidity: function getRoleMembers(bytes32 role) view returns(address[])
func (_BtcTxVerify *BtcTxVerifySession) GetRoleMembers(role [32]byte) ([]common.Address, error) {
	return _BtcTxVerify.Contract.GetRoleMembers(&_BtcTxVerify.CallOpts, role)
}

// GetRoleMembers is a free data retrieval call binding the contract method 0xa3246ad3.
//
// Solidity: function getRoleMembers(bytes32 role) view returns(address[])
func (_BtcTxVerify *BtcTxVerifyCallerSession) GetRoleMembers(role [32]byte) ([]common.Address, error) {
	return _BtcTxVerify.Contract.GetRoleMembers(&_BtcTxVerify.CallOpts, role)
}

// GetTxOuts is a free data retrieval call binding the contract method 0xf71b99de.
//
// Solidity: function getTxOuts(uint64 userAmount, uint64 changeAmount, bytes receiveLockScript, bytes changeLockScript) pure returns((uint64,bytes)[])
func (_BtcTxVerify *BtcTxVerifyCaller) GetTxOuts(opts *bind.CallOpts, userAmount uint64, changeAmount uint64, receiveLockScript []byte, changeLockScript []byte) ([]BtcTxLibTxOut, error) {
	var out []interface{}
	err := _BtcTxVerify.contract.Call(opts, &out, "getTxOuts", userAmount, changeAmount, receiveLockScript, changeLockScript)

	if err != nil {
		return *new([]BtcTxLibTxOut), err
	}

	out0 := *abi.ConvertType(out[0], new([]BtcTxLibTxOut)).(*[]BtcTxLibTxOut)

	return out0, err

}

// GetTxOuts is a free data retrieval call binding the contract method 0xf71b99de.
//
// Solidity: function getTxOuts(uint64 userAmount, uint64 changeAmount, bytes receiveLockScript, bytes changeLockScript) pure returns((uint64,bytes)[])
func (_BtcTxVerify *BtcTxVerifySession) GetTxOuts(userAmount uint64, changeAmount uint64, receiveLockScript []byte, changeLockScript []byte) ([]BtcTxLibTxOut, error) {
	return _BtcTxVerify.Contract.GetTxOuts(&_BtcTxVerify.CallOpts, userAmount, changeAmount, receiveLockScript, changeLockScript)
}

// GetTxOuts is a free data retrieval call binding the contract method 0xf71b99de.
//
// Solidity: function getTxOuts(uint64 userAmount, uint64 changeAmount, bytes receiveLockScript, bytes changeLockScript) pure returns((uint64,bytes)[])
func (_BtcTxVerify *BtcTxVerifyCallerSession) GetTxOuts(userAmount uint64, changeAmount uint64, receiveLockScript []byte, changeLockScript []byte) ([]BtcTxLibTxOut, error) {
	return _BtcTxVerify.Contract.GetTxOuts(&_BtcTxVerify.CallOpts, userAmount, changeAmount, receiveLockScript, changeLockScript)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_BtcTxVerify *BtcTxVerifyCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _BtcTxVerify.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_BtcTxVerify *BtcTxVerifySession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _BtcTxVerify.Contract.HasRole(&_BtcTxVerify.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_BtcTxVerify *BtcTxVerifyCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _BtcTxVerify.Contract.HasRole(&_BtcTxVerify.CallOpts, role, account)
}

// IsCandidateExist is a free data retrieval call binding the contract method 0x526f083e.
//
// Solidity: function isCandidateExist(bytes32 blockHash) view returns(bool)
func (_BtcTxVerify *BtcTxVerifyCaller) IsCandidateExist(opts *bind.CallOpts, blockHash [32]byte) (bool, error) {
	var out []interface{}
	err := _BtcTxVerify.contract.Call(opts, &out, "isCandidateExist", blockHash)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsCandidateExist is a free data retrieval call binding the contract method 0x526f083e.
//
// Solidity: function isCandidateExist(bytes32 blockHash) view returns(bool)
func (_BtcTxVerify *BtcTxVerifySession) IsCandidateExist(blockHash [32]byte) (bool, error) {
	return _BtcTxVerify.Contract.IsCandidateExist(&_BtcTxVerify.CallOpts, blockHash)
}

// IsCandidateExist is a free data retrieval call binding the contract method 0x526f083e.
//
// Solidity: function isCandidateExist(bytes32 blockHash) view returns(bool)
func (_BtcTxVerify *BtcTxVerifyCallerSession) IsCandidateExist(blockHash [32]byte) (bool, error) {
	return _BtcTxVerify.Contract.IsCandidateExist(&_BtcTxVerify.CallOpts, blockHash)
}

// JuniorCP is a free data retrieval call binding the contract method 0x6fe20155.
//
// Solidity: function juniorCP() view returns(bytes32 blockHash, uint64 timestamp)
func (_BtcTxVerify *BtcTxVerifyCaller) JuniorCP(opts *bind.CallOpts) (struct {
	BlockHash [32]byte
	Timestamp uint64
}, error) {
	var out []interface{}
	err := _BtcTxVerify.contract.Call(opts, &out, "juniorCP")

	outstruct := new(struct {
		BlockHash [32]byte
		Timestamp uint64
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.BlockHash = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.Timestamp = *abi.ConvertType(out[1], new(uint64)).(*uint64)

	return *outstruct, err

}

// JuniorCP is a free data retrieval call binding the contract method 0x6fe20155.
//
// Solidity: function juniorCP() view returns(bytes32 blockHash, uint64 timestamp)
func (_BtcTxVerify *BtcTxVerifySession) JuniorCP() (struct {
	BlockHash [32]byte
	Timestamp uint64
}, error) {
	return _BtcTxVerify.Contract.JuniorCP(&_BtcTxVerify.CallOpts)
}

// JuniorCP is a free data retrieval call binding the contract method 0x6fe20155.
//
// Solidity: function juniorCP() view returns(bytes32 blockHash, uint64 timestamp)
func (_BtcTxVerify *BtcTxVerifyCallerSession) JuniorCP() (struct {
	BlockHash [32]byte
	Timestamp uint64
}, error) {
	return _BtcTxVerify.Contract.JuniorCP(&_BtcTxVerify.CallOpts)
}

// PreviousDepositVerifier is a free data retrieval call binding the contract method 0xd60ff340.
//
// Solidity: function previousDepositVerifier() view returns(address)
func (_BtcTxVerify *BtcTxVerifyCaller) PreviousDepositVerifier(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BtcTxVerify.contract.Call(opts, &out, "previousDepositVerifier")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PreviousDepositVerifier is a free data retrieval call binding the contract method 0xd60ff340.
//
// Solidity: function previousDepositVerifier() view returns(address)
func (_BtcTxVerify *BtcTxVerifySession) PreviousDepositVerifier() (common.Address, error) {
	return _BtcTxVerify.Contract.PreviousDepositVerifier(&_BtcTxVerify.CallOpts)
}

// PreviousDepositVerifier is a free data retrieval call binding the contract method 0xd60ff340.
//
// Solidity: function previousDepositVerifier() view returns(address)
func (_BtcTxVerify *BtcTxVerifyCallerSession) PreviousDepositVerifier() (common.Address, error) {
	return _BtcTxVerify.Contract.PreviousDepositVerifier(&_BtcTxVerify.CallOpts)
}

// PreviousRedeemVerifier is a free data retrieval call binding the contract method 0x103aab4f.
//
// Solidity: function previousRedeemVerifier() view returns(address)
func (_BtcTxVerify *BtcTxVerifyCaller) PreviousRedeemVerifier(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BtcTxVerify.contract.Call(opts, &out, "previousRedeemVerifier")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PreviousRedeemVerifier is a free data retrieval call binding the contract method 0x103aab4f.
//
// Solidity: function previousRedeemVerifier() view returns(address)
func (_BtcTxVerify *BtcTxVerifySession) PreviousRedeemVerifier() (common.Address, error) {
	return _BtcTxVerify.Contract.PreviousRedeemVerifier(&_BtcTxVerify.CallOpts)
}

// PreviousRedeemVerifier is a free data retrieval call binding the contract method 0x103aab4f.
//
// Solidity: function previousRedeemVerifier() view returns(address)
func (_BtcTxVerify *BtcTxVerifyCallerSession) PreviousRedeemVerifier() (common.Address, error) {
	return _BtcTxVerify.Contract.PreviousRedeemVerifier(&_BtcTxVerify.CallOpts)
}

// RedeemCompatibleDeadline is a free data retrieval call binding the contract method 0x856cc534.
//
// Solidity: function redeemCompatibleDeadline() view returns(uint256)
func (_BtcTxVerify *BtcTxVerifyCaller) RedeemCompatibleDeadline(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BtcTxVerify.contract.Call(opts, &out, "redeemCompatibleDeadline")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RedeemCompatibleDeadline is a free data retrieval call binding the contract method 0x856cc534.
//
// Solidity: function redeemCompatibleDeadline() view returns(uint256)
func (_BtcTxVerify *BtcTxVerifySession) RedeemCompatibleDeadline() (*big.Int, error) {
	return _BtcTxVerify.Contract.RedeemCompatibleDeadline(&_BtcTxVerify.CallOpts)
}

// RedeemCompatibleDeadline is a free data retrieval call binding the contract method 0x856cc534.
//
// Solidity: function redeemCompatibleDeadline() view returns(uint256)
func (_BtcTxVerify *BtcTxVerifyCallerSession) RedeemCompatibleDeadline() (*big.Int, error) {
	return _BtcTxVerify.Contract.RedeemCompatibleDeadline(&_BtcTxVerify.CallOpts)
}

// RedeemVerifier is a free data retrieval call binding the contract method 0xf1b9f186.
//
// Solidity: function redeemVerifier() view returns(address)
func (_BtcTxVerify *BtcTxVerifyCaller) RedeemVerifier(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BtcTxVerify.contract.Call(opts, &out, "redeemVerifier")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RedeemVerifier is a free data retrieval call binding the contract method 0xf1b9f186.
//
// Solidity: function redeemVerifier() view returns(address)
func (_BtcTxVerify *BtcTxVerifySession) RedeemVerifier() (common.Address, error) {
	return _BtcTxVerify.Contract.RedeemVerifier(&_BtcTxVerify.CallOpts)
}

// RedeemVerifier is a free data retrieval call binding the contract method 0xf1b9f186.
//
// Solidity: function redeemVerifier() view returns(address)
func (_BtcTxVerify *BtcTxVerifyCallerSession) RedeemVerifier() (common.Address, error) {
	return _BtcTxVerify.Contract.RedeemVerifier(&_BtcTxVerify.CallOpts)
}

// SeniorCP is a free data retrieval call binding the contract method 0x940a1cfc.
//
// Solidity: function seniorCP() view returns(bytes32 blockHash, uint64 timestamp)
func (_BtcTxVerify *BtcTxVerifyCaller) SeniorCP(opts *bind.CallOpts) (struct {
	BlockHash [32]byte
	Timestamp uint64
}, error) {
	var out []interface{}
	err := _BtcTxVerify.contract.Call(opts, &out, "seniorCP")

	outstruct := new(struct {
		BlockHash [32]byte
		Timestamp uint64
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.BlockHash = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.Timestamp = *abi.ConvertType(out[1], new(uint64)).(*uint64)

	return *outstruct, err

}

// SeniorCP is a free data retrieval call binding the contract method 0x940a1cfc.
//
// Solidity: function seniorCP() view returns(bytes32 blockHash, uint64 timestamp)
func (_BtcTxVerify *BtcTxVerifySession) SeniorCP() (struct {
	BlockHash [32]byte
	Timestamp uint64
}, error) {
	return _BtcTxVerify.Contract.SeniorCP(&_BtcTxVerify.CallOpts)
}

// SeniorCP is a free data retrieval call binding the contract method 0x940a1cfc.
//
// Solidity: function seniorCP() view returns(bytes32 blockHash, uint64 timestamp)
func (_BtcTxVerify *BtcTxVerifyCallerSession) SeniorCP() (struct {
	BlockHash [32]byte
	Timestamp uint64
}, error) {
	return _BtcTxVerify.Contract.SeniorCP(&_BtcTxVerify.CallOpts)
}

// SuggestedCP is a free data retrieval call binding the contract method 0xb3d9a8b7.
//
// Solidity: function suggestedCP() view returns(bytes32)
func (_BtcTxVerify *BtcTxVerifyCaller) SuggestedCP(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _BtcTxVerify.contract.Call(opts, &out, "suggestedCP")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// SuggestedCP is a free data retrieval call binding the contract method 0xb3d9a8b7.
//
// Solidity: function suggestedCP() view returns(bytes32)
func (_BtcTxVerify *BtcTxVerifySession) SuggestedCP() ([32]byte, error) {
	return _BtcTxVerify.Contract.SuggestedCP(&_BtcTxVerify.CallOpts)
}

// SuggestedCP is a free data retrieval call binding the contract method 0xb3d9a8b7.
//
// Solidity: function suggestedCP() view returns(bytes32)
func (_BtcTxVerify *BtcTxVerifyCallerSession) SuggestedCP() ([32]byte, error) {
	return _BtcTxVerify.Contract.SuggestedCP(&_BtcTxVerify.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_BtcTxVerify *BtcTxVerifyCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _BtcTxVerify.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_BtcTxVerify *BtcTxVerifySession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _BtcTxVerify.Contract.SupportsInterface(&_BtcTxVerify.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_BtcTxVerify *BtcTxVerifyCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _BtcTxVerify.Contract.SupportsInterface(&_BtcTxVerify.CallOpts, interfaceId)
}

// VerifyDepositTx is a free data retrieval call binding the contract method 0xf1012ee2.
//
// Solidity: function verifyDepositTx((bytes32,uint32,uint32,bytes32,uint32,address,uint256,uint32) params, bytes proofData, bool raiseIf, uint64 amount, bytes32 txid) view returns()
func (_BtcTxVerify *BtcTxVerifyCaller) VerifyDepositTx(opts *bind.CallOpts, params BtcTxVerifierPublicWitnessParams, proofData []byte, raiseIf bool, amount uint64, txid [32]byte) error {
	var out []interface{}
	err := _BtcTxVerify.contract.Call(opts, &out, "verifyDepositTx", params, proofData, raiseIf, amount, txid)

	if err != nil {
		return err
	}

	return err

}

// VerifyDepositTx is a free data retrieval call binding the contract method 0xf1012ee2.
//
// Solidity: function verifyDepositTx((bytes32,uint32,uint32,bytes32,uint32,address,uint256,uint32) params, bytes proofData, bool raiseIf, uint64 amount, bytes32 txid) view returns()
func (_BtcTxVerify *BtcTxVerifySession) VerifyDepositTx(params BtcTxVerifierPublicWitnessParams, proofData []byte, raiseIf bool, amount uint64, txid [32]byte) error {
	return _BtcTxVerify.Contract.VerifyDepositTx(&_BtcTxVerify.CallOpts, params, proofData, raiseIf, amount, txid)
}

// VerifyDepositTx is a free data retrieval call binding the contract method 0xf1012ee2.
//
// Solidity: function verifyDepositTx((bytes32,uint32,uint32,bytes32,uint32,address,uint256,uint32) params, bytes proofData, bool raiseIf, uint64 amount, bytes32 txid) view returns()
func (_BtcTxVerify *BtcTxVerifyCallerSession) VerifyDepositTx(params BtcTxVerifierPublicWitnessParams, proofData []byte, raiseIf bool, amount uint64, txid [32]byte) error {
	return _BtcTxVerify.Contract.VerifyDepositTx(&_BtcTxVerify.CallOpts, params, proofData, raiseIf, amount, txid)
}

// VerifyRedeemTx is a free data retrieval call binding the contract method 0x65d933e6.
//
// Solidity: function verifyRedeemTx((bytes32,uint32,uint32,bytes32,uint32,address,uint256,uint32) params, bytes32 txid, uint64 changeAmount, uint256 minerReward, bytes proofData, bool raiseIf) view returns(bool)
func (_BtcTxVerify *BtcTxVerifyCaller) VerifyRedeemTx(opts *bind.CallOpts, params BtcTxVerifierPublicWitnessParams, txid [32]byte, changeAmount uint64, minerReward *big.Int, proofData []byte, raiseIf bool) (bool, error) {
	var out []interface{}
	err := _BtcTxVerify.contract.Call(opts, &out, "verifyRedeemTx", params, txid, changeAmount, minerReward, proofData, raiseIf)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VerifyRedeemTx is a free data retrieval call binding the contract method 0x65d933e6.
//
// Solidity: function verifyRedeemTx((bytes32,uint32,uint32,bytes32,uint32,address,uint256,uint32) params, bytes32 txid, uint64 changeAmount, uint256 minerReward, bytes proofData, bool raiseIf) view returns(bool)
func (_BtcTxVerify *BtcTxVerifySession) VerifyRedeemTx(params BtcTxVerifierPublicWitnessParams, txid [32]byte, changeAmount uint64, minerReward *big.Int, proofData []byte, raiseIf bool) (bool, error) {
	return _BtcTxVerify.Contract.VerifyRedeemTx(&_BtcTxVerify.CallOpts, params, txid, changeAmount, minerReward, proofData, raiseIf)
}

// VerifyRedeemTx is a free data retrieval call binding the contract method 0x65d933e6.
//
// Solidity: function verifyRedeemTx((bytes32,uint32,uint32,bytes32,uint32,address,uint256,uint32) params, bytes32 txid, uint64 changeAmount, uint256 minerReward, bytes proofData, bool raiseIf) view returns(bool)
func (_BtcTxVerify *BtcTxVerifyCallerSession) VerifyRedeemTx(params BtcTxVerifierPublicWitnessParams, txid [32]byte, changeAmount uint64, minerReward *big.Int, proofData []byte, raiseIf bool) (bool, error) {
	return _BtcTxVerify.Contract.VerifyRedeemTx(&_BtcTxVerify.CallOpts, params, txid, changeAmount, minerReward, proofData, raiseIf)
}

// AddNewCP is a paid mutator transaction binding the contract method 0xd17e5fa7.
//
// Solidity: function addNewCP((bytes32,uint32,uint32,bytes32,uint32,address,uint256,uint32) params) returns(bool)
func (_BtcTxVerify *BtcTxVerifyTransactor) AddNewCP(opts *bind.TransactOpts, params BtcTxVerifierPublicWitnessParams) (*types.Transaction, error) {
	return _BtcTxVerify.contract.Transact(opts, "addNewCP", params)
}

// AddNewCP is a paid mutator transaction binding the contract method 0xd17e5fa7.
//
// Solidity: function addNewCP((bytes32,uint32,uint32,bytes32,uint32,address,uint256,uint32) params) returns(bool)
func (_BtcTxVerify *BtcTxVerifySession) AddNewCP(params BtcTxVerifierPublicWitnessParams) (*types.Transaction, error) {
	return _BtcTxVerify.Contract.AddNewCP(&_BtcTxVerify.TransactOpts, params)
}

// AddNewCP is a paid mutator transaction binding the contract method 0xd17e5fa7.
//
// Solidity: function addNewCP((bytes32,uint32,uint32,bytes32,uint32,address,uint256,uint32) params) returns(bool)
func (_BtcTxVerify *BtcTxVerifyTransactorSession) AddNewCP(params BtcTxVerifierPublicWitnessParams) (*types.Transaction, error) {
	return _BtcTxVerify.Contract.AddNewCP(&_BtcTxVerify.TransactOpts, params)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_BtcTxVerify *BtcTxVerifyTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _BtcTxVerify.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_BtcTxVerify *BtcTxVerifySession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _BtcTxVerify.Contract.GrantRole(&_BtcTxVerify.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_BtcTxVerify *BtcTxVerifyTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _BtcTxVerify.Contract.GrantRole(&_BtcTxVerify.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_BtcTxVerify *BtcTxVerifyTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _BtcTxVerify.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_BtcTxVerify *BtcTxVerifySession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _BtcTxVerify.Contract.RenounceRole(&_BtcTxVerify.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_BtcTxVerify *BtcTxVerifyTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _BtcTxVerify.Contract.RenounceRole(&_BtcTxVerify.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_BtcTxVerify *BtcTxVerifyTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _BtcTxVerify.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_BtcTxVerify *BtcTxVerifySession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _BtcTxVerify.Contract.RevokeRole(&_BtcTxVerify.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_BtcTxVerify *BtcTxVerifyTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _BtcTxVerify.Contract.RevokeRole(&_BtcTxVerify.TransactOpts, role, account)
}

// TryRotateCP is a paid mutator transaction binding the contract method 0x52f00804.
//
// Solidity: function tryRotateCP() returns()
func (_BtcTxVerify *BtcTxVerifyTransactor) TryRotateCP(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BtcTxVerify.contract.Transact(opts, "tryRotateCP")
}

// TryRotateCP is a paid mutator transaction binding the contract method 0x52f00804.
//
// Solidity: function tryRotateCP() returns()
func (_BtcTxVerify *BtcTxVerifySession) TryRotateCP() (*types.Transaction, error) {
	return _BtcTxVerify.Contract.TryRotateCP(&_BtcTxVerify.TransactOpts)
}

// TryRotateCP is a paid mutator transaction binding the contract method 0x52f00804.
//
// Solidity: function tryRotateCP() returns()
func (_BtcTxVerify *BtcTxVerifyTransactorSession) TryRotateCP() (*types.Transaction, error) {
	return _BtcTxVerify.Contract.TryRotateCP(&_BtcTxVerify.TransactOpts)
}

// UpdateDepositVerifier is a paid mutator transaction binding the contract method 0x0c9359d7.
//
// Solidity: function updateDepositVerifier(address addr, uint32 reserveHours) returns()
func (_BtcTxVerify *BtcTxVerifyTransactor) UpdateDepositVerifier(opts *bind.TransactOpts, addr common.Address, reserveHours uint32) (*types.Transaction, error) {
	return _BtcTxVerify.contract.Transact(opts, "updateDepositVerifier", addr, reserveHours)
}

// UpdateDepositVerifier is a paid mutator transaction binding the contract method 0x0c9359d7.
//
// Solidity: function updateDepositVerifier(address addr, uint32 reserveHours) returns()
func (_BtcTxVerify *BtcTxVerifySession) UpdateDepositVerifier(addr common.Address, reserveHours uint32) (*types.Transaction, error) {
	return _BtcTxVerify.Contract.UpdateDepositVerifier(&_BtcTxVerify.TransactOpts, addr, reserveHours)
}

// UpdateDepositVerifier is a paid mutator transaction binding the contract method 0x0c9359d7.
//
// Solidity: function updateDepositVerifier(address addr, uint32 reserveHours) returns()
func (_BtcTxVerify *BtcTxVerifyTransactorSession) UpdateDepositVerifier(addr common.Address, reserveHours uint32) (*types.Transaction, error) {
	return _BtcTxVerify.Contract.UpdateDepositVerifier(&_BtcTxVerify.TransactOpts, addr, reserveHours)
}

// UpdateExtraDepth is a paid mutator transaction binding the contract method 0xbbdec9ef.
//
// Solidity: function updateExtraDepth(uint8 depth) returns()
func (_BtcTxVerify *BtcTxVerifyTransactor) UpdateExtraDepth(opts *bind.TransactOpts, depth uint8) (*types.Transaction, error) {
	return _BtcTxVerify.contract.Transact(opts, "updateExtraDepth", depth)
}

// UpdateExtraDepth is a paid mutator transaction binding the contract method 0xbbdec9ef.
//
// Solidity: function updateExtraDepth(uint8 depth) returns()
func (_BtcTxVerify *BtcTxVerifySession) UpdateExtraDepth(depth uint8) (*types.Transaction, error) {
	return _BtcTxVerify.Contract.UpdateExtraDepth(&_BtcTxVerify.TransactOpts, depth)
}

// UpdateExtraDepth is a paid mutator transaction binding the contract method 0xbbdec9ef.
//
// Solidity: function updateExtraDepth(uint8 depth) returns()
func (_BtcTxVerify *BtcTxVerifyTransactorSession) UpdateExtraDepth(depth uint8) (*types.Transaction, error) {
	return _BtcTxVerify.Contract.UpdateExtraDepth(&_BtcTxVerify.TransactOpts, depth)
}

// UpdateRedeemVerifier is a paid mutator transaction binding the contract method 0x4cb9eb8b.
//
// Solidity: function updateRedeemVerifier(address addr, uint32 reserveHours) returns()
func (_BtcTxVerify *BtcTxVerifyTransactor) UpdateRedeemVerifier(opts *bind.TransactOpts, addr common.Address, reserveHours uint32) (*types.Transaction, error) {
	return _BtcTxVerify.contract.Transact(opts, "updateRedeemVerifier", addr, reserveHours)
}

// UpdateRedeemVerifier is a paid mutator transaction binding the contract method 0x4cb9eb8b.
//
// Solidity: function updateRedeemVerifier(address addr, uint32 reserveHours) returns()
func (_BtcTxVerify *BtcTxVerifySession) UpdateRedeemVerifier(addr common.Address, reserveHours uint32) (*types.Transaction, error) {
	return _BtcTxVerify.Contract.UpdateRedeemVerifier(&_BtcTxVerify.TransactOpts, addr, reserveHours)
}

// UpdateRedeemVerifier is a paid mutator transaction binding the contract method 0x4cb9eb8b.
//
// Solidity: function updateRedeemVerifier(address addr, uint32 reserveHours) returns()
func (_BtcTxVerify *BtcTxVerifyTransactorSession) UpdateRedeemVerifier(addr common.Address, reserveHours uint32) (*types.Transaction, error) {
	return _BtcTxVerify.Contract.UpdateRedeemVerifier(&_BtcTxVerify.TransactOpts, addr, reserveHours)
}

// UpdateUnsignedProtection is a paid mutator transaction binding the contract method 0x31519808.
//
// Solidity: function updateUnsignedProtection(bool flag) returns()
func (_BtcTxVerify *BtcTxVerifyTransactor) UpdateUnsignedProtection(opts *bind.TransactOpts, flag bool) (*types.Transaction, error) {
	return _BtcTxVerify.contract.Transact(opts, "updateUnsignedProtection", flag)
}

// UpdateUnsignedProtection is a paid mutator transaction binding the contract method 0x31519808.
//
// Solidity: function updateUnsignedProtection(bool flag) returns()
func (_BtcTxVerify *BtcTxVerifySession) UpdateUnsignedProtection(flag bool) (*types.Transaction, error) {
	return _BtcTxVerify.Contract.UpdateUnsignedProtection(&_BtcTxVerify.TransactOpts, flag)
}

// UpdateUnsignedProtection is a paid mutator transaction binding the contract method 0x31519808.
//
// Solidity: function updateUnsignedProtection(bool flag) returns()
func (_BtcTxVerify *BtcTxVerifyTransactorSession) UpdateUnsignedProtection(flag bool) (*types.Transaction, error) {
	return _BtcTxVerify.Contract.UpdateUnsignedProtection(&_BtcTxVerify.TransactOpts, flag)
}

// BtcTxVerifyCheckPointRotatedIterator is returned from FilterCheckPointRotated and is used to iterate over the raw logs and unpacked data for CheckPointRotated events raised by the BtcTxVerify contract.
type BtcTxVerifyCheckPointRotatedIterator struct {
	Event *BtcTxVerifyCheckPointRotated // Event containing the contract specifics and raw log

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
func (it *BtcTxVerifyCheckPointRotatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BtcTxVerifyCheckPointRotated)
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
		it.Event = new(BtcTxVerifyCheckPointRotated)
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
func (it *BtcTxVerifyCheckPointRotatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BtcTxVerifyCheckPointRotatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BtcTxVerifyCheckPointRotated represents a CheckPointRotated event raised by the BtcTxVerify contract.
type BtcTxVerifyCheckPointRotated struct {
	OldCp [32]byte
	NewCp [32]byte
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterCheckPointRotated is a free log retrieval operation binding the contract event 0x075a8fc0869a4712d75a7c6c6062d9b26e81f09f4945c584dc70ff7a2660c7d2.
//
// Solidity: event CheckPointRotated(bytes32 oldCp, bytes32 newCp)
func (_BtcTxVerify *BtcTxVerifyFilterer) FilterCheckPointRotated(opts *bind.FilterOpts) (*BtcTxVerifyCheckPointRotatedIterator, error) {

	logs, sub, err := _BtcTxVerify.contract.FilterLogs(opts, "CheckPointRotated")
	if err != nil {
		return nil, err
	}
	return &BtcTxVerifyCheckPointRotatedIterator{contract: _BtcTxVerify.contract, event: "CheckPointRotated", logs: logs, sub: sub}, nil
}

// WatchCheckPointRotated is a free log subscription operation binding the contract event 0x075a8fc0869a4712d75a7c6c6062d9b26e81f09f4945c584dc70ff7a2660c7d2.
//
// Solidity: event CheckPointRotated(bytes32 oldCp, bytes32 newCp)
func (_BtcTxVerify *BtcTxVerifyFilterer) WatchCheckPointRotated(opts *bind.WatchOpts, sink chan<- *BtcTxVerifyCheckPointRotated) (event.Subscription, error) {

	logs, sub, err := _BtcTxVerify.contract.WatchLogs(opts, "CheckPointRotated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BtcTxVerifyCheckPointRotated)
				if err := _BtcTxVerify.contract.UnpackLog(event, "CheckPointRotated", log); err != nil {
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

// ParseCheckPointRotated is a log parse operation binding the contract event 0x075a8fc0869a4712d75a7c6c6062d9b26e81f09f4945c584dc70ff7a2660c7d2.
//
// Solidity: event CheckPointRotated(bytes32 oldCp, bytes32 newCp)
func (_BtcTxVerify *BtcTxVerifyFilterer) ParseCheckPointRotated(log types.Log) (*BtcTxVerifyCheckPointRotated, error) {
	event := new(BtcTxVerifyCheckPointRotated)
	if err := _BtcTxVerify.contract.UnpackLog(event, "CheckPointRotated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BtcTxVerifyDepositVerifierUpdatedIterator is returned from FilterDepositVerifierUpdated and is used to iterate over the raw logs and unpacked data for DepositVerifierUpdated events raised by the BtcTxVerify contract.
type BtcTxVerifyDepositVerifierUpdatedIterator struct {
	Event *BtcTxVerifyDepositVerifierUpdated // Event containing the contract specifics and raw log

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
func (it *BtcTxVerifyDepositVerifierUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BtcTxVerifyDepositVerifierUpdated)
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
		it.Event = new(BtcTxVerifyDepositVerifierUpdated)
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
func (it *BtcTxVerifyDepositVerifierUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BtcTxVerifyDepositVerifierUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BtcTxVerifyDepositVerifierUpdated represents a DepositVerifierUpdated event raised by the BtcTxVerify contract.
type BtcTxVerifyDepositVerifierUpdated struct {
	Old           common.Address
	Now           common.Address
	ReservedHours uint32
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterDepositVerifierUpdated is a free log retrieval operation binding the contract event 0xe499267739d8e3f0ef85a9937594ca2fb74d905a3f977189c85e9fc83ec61b5e.
//
// Solidity: event DepositVerifierUpdated(address old, address now, uint32 reservedHours)
func (_BtcTxVerify *BtcTxVerifyFilterer) FilterDepositVerifierUpdated(opts *bind.FilterOpts) (*BtcTxVerifyDepositVerifierUpdatedIterator, error) {

	logs, sub, err := _BtcTxVerify.contract.FilterLogs(opts, "DepositVerifierUpdated")
	if err != nil {
		return nil, err
	}
	return &BtcTxVerifyDepositVerifierUpdatedIterator{contract: _BtcTxVerify.contract, event: "DepositVerifierUpdated", logs: logs, sub: sub}, nil
}

// WatchDepositVerifierUpdated is a free log subscription operation binding the contract event 0xe499267739d8e3f0ef85a9937594ca2fb74d905a3f977189c85e9fc83ec61b5e.
//
// Solidity: event DepositVerifierUpdated(address old, address now, uint32 reservedHours)
func (_BtcTxVerify *BtcTxVerifyFilterer) WatchDepositVerifierUpdated(opts *bind.WatchOpts, sink chan<- *BtcTxVerifyDepositVerifierUpdated) (event.Subscription, error) {

	logs, sub, err := _BtcTxVerify.contract.WatchLogs(opts, "DepositVerifierUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BtcTxVerifyDepositVerifierUpdated)
				if err := _BtcTxVerify.contract.UnpackLog(event, "DepositVerifierUpdated", log); err != nil {
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

// ParseDepositVerifierUpdated is a log parse operation binding the contract event 0xe499267739d8e3f0ef85a9937594ca2fb74d905a3f977189c85e9fc83ec61b5e.
//
// Solidity: event DepositVerifierUpdated(address old, address now, uint32 reservedHours)
func (_BtcTxVerify *BtcTxVerifyFilterer) ParseDepositVerifierUpdated(log types.Log) (*BtcTxVerifyDepositVerifierUpdated, error) {
	event := new(BtcTxVerifyDepositVerifierUpdated)
	if err := _BtcTxVerify.contract.UnpackLog(event, "DepositVerifierUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BtcTxVerifyRedeemVerifierUpdatedIterator is returned from FilterRedeemVerifierUpdated and is used to iterate over the raw logs and unpacked data for RedeemVerifierUpdated events raised by the BtcTxVerify contract.
type BtcTxVerifyRedeemVerifierUpdatedIterator struct {
	Event *BtcTxVerifyRedeemVerifierUpdated // Event containing the contract specifics and raw log

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
func (it *BtcTxVerifyRedeemVerifierUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BtcTxVerifyRedeemVerifierUpdated)
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
		it.Event = new(BtcTxVerifyRedeemVerifierUpdated)
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
func (it *BtcTxVerifyRedeemVerifierUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BtcTxVerifyRedeemVerifierUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BtcTxVerifyRedeemVerifierUpdated represents a RedeemVerifierUpdated event raised by the BtcTxVerify contract.
type BtcTxVerifyRedeemVerifierUpdated struct {
	Old           common.Address
	Now           common.Address
	ReservedHours uint32
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterRedeemVerifierUpdated is a free log retrieval operation binding the contract event 0xe38b65ae1958c5f6a45d5f3974de42e4ac3baa1ed1a15b15a4734d9ee8c0910d.
//
// Solidity: event RedeemVerifierUpdated(address old, address now, uint32 reservedHours)
func (_BtcTxVerify *BtcTxVerifyFilterer) FilterRedeemVerifierUpdated(opts *bind.FilterOpts) (*BtcTxVerifyRedeemVerifierUpdatedIterator, error) {

	logs, sub, err := _BtcTxVerify.contract.FilterLogs(opts, "RedeemVerifierUpdated")
	if err != nil {
		return nil, err
	}
	return &BtcTxVerifyRedeemVerifierUpdatedIterator{contract: _BtcTxVerify.contract, event: "RedeemVerifierUpdated", logs: logs, sub: sub}, nil
}

// WatchRedeemVerifierUpdated is a free log subscription operation binding the contract event 0xe38b65ae1958c5f6a45d5f3974de42e4ac3baa1ed1a15b15a4734d9ee8c0910d.
//
// Solidity: event RedeemVerifierUpdated(address old, address now, uint32 reservedHours)
func (_BtcTxVerify *BtcTxVerifyFilterer) WatchRedeemVerifierUpdated(opts *bind.WatchOpts, sink chan<- *BtcTxVerifyRedeemVerifierUpdated) (event.Subscription, error) {

	logs, sub, err := _BtcTxVerify.contract.WatchLogs(opts, "RedeemVerifierUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BtcTxVerifyRedeemVerifierUpdated)
				if err := _BtcTxVerify.contract.UnpackLog(event, "RedeemVerifierUpdated", log); err != nil {
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

// ParseRedeemVerifierUpdated is a log parse operation binding the contract event 0xe38b65ae1958c5f6a45d5f3974de42e4ac3baa1ed1a15b15a4734d9ee8c0910d.
//
// Solidity: event RedeemVerifierUpdated(address old, address now, uint32 reservedHours)
func (_BtcTxVerify *BtcTxVerifyFilterer) ParseRedeemVerifierUpdated(log types.Log) (*BtcTxVerifyRedeemVerifierUpdated, error) {
	event := new(BtcTxVerifyRedeemVerifierUpdated)
	if err := _BtcTxVerify.contract.UnpackLog(event, "RedeemVerifierUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BtcTxVerifyRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the BtcTxVerify contract.
type BtcTxVerifyRoleAdminChangedIterator struct {
	Event *BtcTxVerifyRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *BtcTxVerifyRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BtcTxVerifyRoleAdminChanged)
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
		it.Event = new(BtcTxVerifyRoleAdminChanged)
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
func (it *BtcTxVerifyRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BtcTxVerifyRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BtcTxVerifyRoleAdminChanged represents a RoleAdminChanged event raised by the BtcTxVerify contract.
type BtcTxVerifyRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_BtcTxVerify *BtcTxVerifyFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*BtcTxVerifyRoleAdminChangedIterator, error) {

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

	logs, sub, err := _BtcTxVerify.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &BtcTxVerifyRoleAdminChangedIterator{contract: _BtcTxVerify.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_BtcTxVerify *BtcTxVerifyFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *BtcTxVerifyRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _BtcTxVerify.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BtcTxVerifyRoleAdminChanged)
				if err := _BtcTxVerify.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_BtcTxVerify *BtcTxVerifyFilterer) ParseRoleAdminChanged(log types.Log) (*BtcTxVerifyRoleAdminChanged, error) {
	event := new(BtcTxVerifyRoleAdminChanged)
	if err := _BtcTxVerify.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BtcTxVerifyRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the BtcTxVerify contract.
type BtcTxVerifyRoleGrantedIterator struct {
	Event *BtcTxVerifyRoleGranted // Event containing the contract specifics and raw log

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
func (it *BtcTxVerifyRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BtcTxVerifyRoleGranted)
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
		it.Event = new(BtcTxVerifyRoleGranted)
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
func (it *BtcTxVerifyRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BtcTxVerifyRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BtcTxVerifyRoleGranted represents a RoleGranted event raised by the BtcTxVerify contract.
type BtcTxVerifyRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_BtcTxVerify *BtcTxVerifyFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*BtcTxVerifyRoleGrantedIterator, error) {

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

	logs, sub, err := _BtcTxVerify.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &BtcTxVerifyRoleGrantedIterator{contract: _BtcTxVerify.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_BtcTxVerify *BtcTxVerifyFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *BtcTxVerifyRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _BtcTxVerify.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BtcTxVerifyRoleGranted)
				if err := _BtcTxVerify.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_BtcTxVerify *BtcTxVerifyFilterer) ParseRoleGranted(log types.Log) (*BtcTxVerifyRoleGranted, error) {
	event := new(BtcTxVerifyRoleGranted)
	if err := _BtcTxVerify.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BtcTxVerifyRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the BtcTxVerify contract.
type BtcTxVerifyRoleRevokedIterator struct {
	Event *BtcTxVerifyRoleRevoked // Event containing the contract specifics and raw log

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
func (it *BtcTxVerifyRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BtcTxVerifyRoleRevoked)
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
		it.Event = new(BtcTxVerifyRoleRevoked)
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
func (it *BtcTxVerifyRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BtcTxVerifyRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BtcTxVerifyRoleRevoked represents a RoleRevoked event raised by the BtcTxVerify contract.
type BtcTxVerifyRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_BtcTxVerify *BtcTxVerifyFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*BtcTxVerifyRoleRevokedIterator, error) {

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

	logs, sub, err := _BtcTxVerify.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &BtcTxVerifyRoleRevokedIterator{contract: _BtcTxVerify.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_BtcTxVerify *BtcTxVerifyFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *BtcTxVerifyRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _BtcTxVerify.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BtcTxVerifyRoleRevoked)
				if err := _BtcTxVerify.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_BtcTxVerify *BtcTxVerifyFilterer) ParseRoleRevoked(log types.Log) (*BtcTxVerifyRoleRevoked, error) {
	event := new(BtcTxVerifyRoleRevoked)
	if err := _BtcTxVerify.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BtcTxVerifyUnsignedProtectionChangedIterator is returned from FilterUnsignedProtectionChanged and is used to iterate over the raw logs and unpacked data for UnsignedProtectionChanged events raised by the BtcTxVerify contract.
type BtcTxVerifyUnsignedProtectionChangedIterator struct {
	Event *BtcTxVerifyUnsignedProtectionChanged // Event containing the contract specifics and raw log

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
func (it *BtcTxVerifyUnsignedProtectionChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BtcTxVerifyUnsignedProtectionChanged)
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
		it.Event = new(BtcTxVerifyUnsignedProtectionChanged)
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
func (it *BtcTxVerifyUnsignedProtectionChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BtcTxVerifyUnsignedProtectionChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BtcTxVerifyUnsignedProtectionChanged represents a UnsignedProtectionChanged event raised by the BtcTxVerify contract.
type BtcTxVerifyUnsignedProtectionChanged struct {
	Now bool
	Raw types.Log // Blockchain specific contextual infos
}

// FilterUnsignedProtectionChanged is a free log retrieval operation binding the contract event 0x493cf75f56474178b2649e63e9a83efe32bb15c3b498c3981efddea44e318b5c.
//
// Solidity: event UnsignedProtectionChanged(bool now)
func (_BtcTxVerify *BtcTxVerifyFilterer) FilterUnsignedProtectionChanged(opts *bind.FilterOpts) (*BtcTxVerifyUnsignedProtectionChangedIterator, error) {

	logs, sub, err := _BtcTxVerify.contract.FilterLogs(opts, "UnsignedProtectionChanged")
	if err != nil {
		return nil, err
	}
	return &BtcTxVerifyUnsignedProtectionChangedIterator{contract: _BtcTxVerify.contract, event: "UnsignedProtectionChanged", logs: logs, sub: sub}, nil
}

// WatchUnsignedProtectionChanged is a free log subscription operation binding the contract event 0x493cf75f56474178b2649e63e9a83efe32bb15c3b498c3981efddea44e318b5c.
//
// Solidity: event UnsignedProtectionChanged(bool now)
func (_BtcTxVerify *BtcTxVerifyFilterer) WatchUnsignedProtectionChanged(opts *bind.WatchOpts, sink chan<- *BtcTxVerifyUnsignedProtectionChanged) (event.Subscription, error) {

	logs, sub, err := _BtcTxVerify.contract.WatchLogs(opts, "UnsignedProtectionChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BtcTxVerifyUnsignedProtectionChanged)
				if err := _BtcTxVerify.contract.UnpackLog(event, "UnsignedProtectionChanged", log); err != nil {
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

// ParseUnsignedProtectionChanged is a log parse operation binding the contract event 0x493cf75f56474178b2649e63e9a83efe32bb15c3b498c3981efddea44e318b5c.
//
// Solidity: event UnsignedProtectionChanged(bool now)
func (_BtcTxVerify *BtcTxVerifyFilterer) ParseUnsignedProtectionChanged(log types.Log) (*BtcTxVerifyUnsignedProtectionChanged, error) {
	event := new(BtcTxVerifyUnsignedProtectionChanged)
	if err := _BtcTxVerify.contract.UnpackLog(event, "UnsignedProtectionChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BtcTxVerifyUpdatedExtraDepthIterator is returned from FilterUpdatedExtraDepth and is used to iterate over the raw logs and unpacked data for UpdatedExtraDepth events raised by the BtcTxVerify contract.
type BtcTxVerifyUpdatedExtraDepthIterator struct {
	Event *BtcTxVerifyUpdatedExtraDepth // Event containing the contract specifics and raw log

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
func (it *BtcTxVerifyUpdatedExtraDepthIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BtcTxVerifyUpdatedExtraDepth)
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
		it.Event = new(BtcTxVerifyUpdatedExtraDepth)
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
func (it *BtcTxVerifyUpdatedExtraDepthIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BtcTxVerifyUpdatedExtraDepthIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BtcTxVerifyUpdatedExtraDepth represents a UpdatedExtraDepth event raised by the BtcTxVerify contract.
type BtcTxVerifyUpdatedExtraDepth struct {
	OldDepth uint8
	NewDepth uint8
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterUpdatedExtraDepth is a free log retrieval operation binding the contract event 0x3ae1f9c67c8e41f945172b52a888ffec38d938786f2338406f04afe05e21869b.
//
// Solidity: event UpdatedExtraDepth(uint8 oldDepth, uint8 newDepth)
func (_BtcTxVerify *BtcTxVerifyFilterer) FilterUpdatedExtraDepth(opts *bind.FilterOpts) (*BtcTxVerifyUpdatedExtraDepthIterator, error) {

	logs, sub, err := _BtcTxVerify.contract.FilterLogs(opts, "UpdatedExtraDepth")
	if err != nil {
		return nil, err
	}
	return &BtcTxVerifyUpdatedExtraDepthIterator{contract: _BtcTxVerify.contract, event: "UpdatedExtraDepth", logs: logs, sub: sub}, nil
}

// WatchUpdatedExtraDepth is a free log subscription operation binding the contract event 0x3ae1f9c67c8e41f945172b52a888ffec38d938786f2338406f04afe05e21869b.
//
// Solidity: event UpdatedExtraDepth(uint8 oldDepth, uint8 newDepth)
func (_BtcTxVerify *BtcTxVerifyFilterer) WatchUpdatedExtraDepth(opts *bind.WatchOpts, sink chan<- *BtcTxVerifyUpdatedExtraDepth) (event.Subscription, error) {

	logs, sub, err := _BtcTxVerify.contract.WatchLogs(opts, "UpdatedExtraDepth")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BtcTxVerifyUpdatedExtraDepth)
				if err := _BtcTxVerify.contract.UnpackLog(event, "UpdatedExtraDepth", log); err != nil {
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

// ParseUpdatedExtraDepth is a log parse operation binding the contract event 0x3ae1f9c67c8e41f945172b52a888ffec38d938786f2338406f04afe05e21869b.
//
// Solidity: event UpdatedExtraDepth(uint8 oldDepth, uint8 newDepth)
func (_BtcTxVerify *BtcTxVerifyFilterer) ParseUpdatedExtraDepth(log types.Log) (*BtcTxVerifyUpdatedExtraDepth, error) {
	event := new(BtcTxVerifyUpdatedExtraDepth)
	if err := _BtcTxVerify.contract.UnpackLog(event, "UpdatedExtraDepth", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
