// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package noop_caller

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

// NoopCallerMetaData contains all meta data concerning the NoopCaller contract.
var NoopCallerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"noop\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"noop_static_call\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234601c57600e6020565b61049261002c823961049290f35b6026565b60405190565b600080fdfe60806040526004361015610013575b610157565b61001e60003561003d565b80635dfc2e4a146100385763a79ad1a50361000e57610122565b610069565b60e01c90565b60405190565b600080fd5b600080fd5b600091031261005e57565b61004e565b60000190565b3461009757610079366004610053565b61008161022c565b610089610043565b8061009381610063565b0390f35b610049565b5190565b60209181520190565b60005b8381106100bd575050906000910152565b8060209183015181850152016100ac565b601f801991011690565b6100f7610100602093610105936100ee8161009c565b938480936100a0565b958691016100a9565b6100ce565b0190565b61011f91602082019160008184039101526100d8565b90565b3461015257610132366004610053565b61014e61013d6103f3565b610145610043565b91829182610109565b0390f35b610049565b600080fd5b60016009609c1b0190565b60018060a01b031690565b90565b61018961018461018e92610167565b610172565b610167565b90565b61019a90610175565b90565b6101a690610191565b90565b6101b290610175565b90565b6101be906101a9565b90565b600080fd5b634e487b7160e01b600052604160045260246000fd5b906101e6906100ce565b810190811067ffffffffffffffff82111761020057604052565b6101c6565b60e01b90565b600091031261021657565b61004e565b610223610043565b3d6000823e3d90fd5b61024461023f61023a61015c565b61019d565b6101b5565b635dfc2e4a90803b156102bc5761026891600091610260610043565b938492610205565b8252818061027860048201610063565b03915afa80156102b75761028a575b50565b6102aa9060003d81116102b0575b6102a281836101dc565b81019061020b565b38610287565b503d610298565b61021b565b6101c1565b606090565b906102d96102d2610043565b92836101dc565b565b67ffffffffffffffff81116102f9576102f56020916100ce565b0190565b6101c6565b9061031061030b836102db565b6102c6565b918252565b3d600014610332576103263d6102fe565b903d6000602084013e5b565b61033a6102c1565b90610330565b60209181520190565b60207f6c65640000000000000000000000000000000000000000000000000000000000917f63616c6c20746f20707265636f6d70696c656420636f6e74726163742066616960008201520152565b6103a46023604092610340565b6103ad81610349565b0190565b6103c79060208101906000818303910152610397565b90565b156103d157565b6103d9610043565b62461bcd60e51b8152806103ef600482016103b1565b0390fd5b6103fb6102c1565b506000806004610436632efe172560e11b610427610417610043565b9384926020840190815201610063565b602082018103825203826101dc565b61043e61015c565b90602081019051915afa610459610453610315565b916103ca565b9056fea2646970667358221220e789f3802bdf70bd0bc015146d5093e7d8b7efa129b3f663af58a671421d9b6964736f6c63430008190033",
}

// NoopCallerABI is the input ABI used to generate the binding from.
// Deprecated: Use NoopCallerMetaData.ABI instead.
var NoopCallerABI = NoopCallerMetaData.ABI

// NoopCallerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use NoopCallerMetaData.Bin instead.
var NoopCallerBin = NoopCallerMetaData.Bin

// DeployNoopCaller deploys a new Ethereum contract, binding an instance of NoopCaller to it.
func DeployNoopCaller(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *NoopCaller, error) {
	parsed, err := NoopCallerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(NoopCallerBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &NoopCaller{NoopCallerCaller: NoopCallerCaller{contract: contract}, NoopCallerTransactor: NoopCallerTransactor{contract: contract}, NoopCallerFilterer: NoopCallerFilterer{contract: contract}}, nil
}

// NoopCaller is an auto generated Go binding around an Ethereum contract.
type NoopCaller struct {
	NoopCallerCaller     // Read-only binding to the contract
	NoopCallerTransactor // Write-only binding to the contract
	NoopCallerFilterer   // Log filterer for contract events
}

// NoopCallerCaller is an auto generated read-only Go binding around an Ethereum contract.
type NoopCallerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NoopCallerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type NoopCallerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NoopCallerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type NoopCallerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NoopCallerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type NoopCallerSession struct {
	Contract     *NoopCaller       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// NoopCallerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type NoopCallerCallerSession struct {
	Contract *NoopCallerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// NoopCallerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type NoopCallerTransactorSession struct {
	Contract     *NoopCallerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// NoopCallerRaw is an auto generated low-level Go binding around an Ethereum contract.
type NoopCallerRaw struct {
	Contract *NoopCaller // Generic contract binding to access the raw methods on
}

// NoopCallerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type NoopCallerCallerRaw struct {
	Contract *NoopCallerCaller // Generic read-only contract binding to access the raw methods on
}

// NoopCallerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type NoopCallerTransactorRaw struct {
	Contract *NoopCallerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewNoopCaller creates a new instance of NoopCaller, bound to a specific deployed contract.
func NewNoopCaller(address common.Address, backend bind.ContractBackend) (*NoopCaller, error) {
	contract, err := bindNoopCaller(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NoopCaller{NoopCallerCaller: NoopCallerCaller{contract: contract}, NoopCallerTransactor: NoopCallerTransactor{contract: contract}, NoopCallerFilterer: NoopCallerFilterer{contract: contract}}, nil
}

// NewNoopCallerCaller creates a new read-only instance of NoopCaller, bound to a specific deployed contract.
func NewNoopCallerCaller(address common.Address, caller bind.ContractCaller) (*NoopCallerCaller, error) {
	contract, err := bindNoopCaller(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NoopCallerCaller{contract: contract}, nil
}

// NewNoopCallerTransactor creates a new write-only instance of NoopCaller, bound to a specific deployed contract.
func NewNoopCallerTransactor(address common.Address, transactor bind.ContractTransactor) (*NoopCallerTransactor, error) {
	contract, err := bindNoopCaller(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NoopCallerTransactor{contract: contract}, nil
}

// NewNoopCallerFilterer creates a new log filterer instance of NoopCaller, bound to a specific deployed contract.
func NewNoopCallerFilterer(address common.Address, filterer bind.ContractFilterer) (*NoopCallerFilterer, error) {
	contract, err := bindNoopCaller(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NoopCallerFilterer{contract: contract}, nil
}

// bindNoopCaller binds a generic wrapper to an already deployed contract.
func bindNoopCaller(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := NoopCallerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NoopCaller *NoopCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NoopCaller.Contract.NoopCallerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NoopCaller *NoopCallerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NoopCaller.Contract.NoopCallerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NoopCaller *NoopCallerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NoopCaller.Contract.NoopCallerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NoopCaller *NoopCallerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NoopCaller.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NoopCaller *NoopCallerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NoopCaller.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NoopCaller *NoopCallerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NoopCaller.Contract.contract.Transact(opts, method, params...)
}

// Noop is a free data retrieval call binding the contract method 0x5dfc2e4a.
//
// Solidity: function noop() view returns()
func (_NoopCaller *NoopCallerCaller) Noop(opts *bind.CallOpts) error {
	var out []interface{}
	err := _NoopCaller.contract.Call(opts, &out, "noop")

	if err != nil {
		return err
	}

	return err

}

// Noop is a free data retrieval call binding the contract method 0x5dfc2e4a.
//
// Solidity: function noop() view returns()
func (_NoopCaller *NoopCallerSession) Noop() error {
	return _NoopCaller.Contract.Noop(&_NoopCaller.CallOpts)
}

// Noop is a free data retrieval call binding the contract method 0x5dfc2e4a.
//
// Solidity: function noop() view returns()
func (_NoopCaller *NoopCallerCallerSession) Noop() error {
	return _NoopCaller.Contract.Noop(&_NoopCaller.CallOpts)
}

// NoopStaticCall is a free data retrieval call binding the contract method 0xa79ad1a5.
//
// Solidity: function noop_static_call() view returns(bytes)
func (_NoopCaller *NoopCallerCaller) NoopStaticCall(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _NoopCaller.contract.Call(opts, &out, "noop_static_call")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// NoopStaticCall is a free data retrieval call binding the contract method 0xa79ad1a5.
//
// Solidity: function noop_static_call() view returns(bytes)
func (_NoopCaller *NoopCallerSession) NoopStaticCall() ([]byte, error) {
	return _NoopCaller.Contract.NoopStaticCall(&_NoopCaller.CallOpts)
}

// NoopStaticCall is a free data retrieval call binding the contract method 0xa79ad1a5.
//
// Solidity: function noop_static_call() view returns(bytes)
func (_NoopCaller *NoopCallerCallerSession) NoopStaticCall() ([]byte, error) {
	return _NoopCaller.Contract.NoopStaticCall(&_NoopCaller.CallOpts)
}
