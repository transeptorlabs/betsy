// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package examples

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

// GlobalCounterMetaData contains all meta data concerning the GlobalCounter contract.
var GlobalCounterMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"currentCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"increment\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"}]",
	Bin: "0x608060405234801561001057600080fd5b506000805560cc806100236000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c8063c732d201146037578063d09de08a146051575b600080fd5b603f60005481565b60405190815260200160405180910390f35b60576059565b005b6001600080828254606991906070565b9091555050565b80820180821115609057634e487b7160e01b600052601160045260246000fd5b9291505056fea2646970667358221220dcb66b3ab7c6b03ef3de6b73ba07c12208d29b28aed9ff419717bf1a4fab160f64736f6c63430008170033",
}

// GlobalCounterABI is the input ABI used to generate the binding from.
// Deprecated: Use GlobalCounterMetaData.ABI instead.
var GlobalCounterABI = GlobalCounterMetaData.ABI

// GlobalCounterBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use GlobalCounterMetaData.Bin instead.
var GlobalCounterBin = GlobalCounterMetaData.Bin

// DeployGlobalCounter deploys a new Ethereum contract, binding an instance of GlobalCounter to it.
func DeployGlobalCounter(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *GlobalCounter, error) {
	parsed, err := GlobalCounterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(GlobalCounterBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &GlobalCounter{GlobalCounterCaller: GlobalCounterCaller{contract: contract}, GlobalCounterTransactor: GlobalCounterTransactor{contract: contract}, GlobalCounterFilterer: GlobalCounterFilterer{contract: contract}}, nil
}

// GlobalCounter is an auto generated Go binding around an Ethereum contract.
type GlobalCounter struct {
	GlobalCounterCaller     // Read-only binding to the contract
	GlobalCounterTransactor // Write-only binding to the contract
	GlobalCounterFilterer   // Log filterer for contract events
}

// GlobalCounterCaller is an auto generated read-only Go binding around an Ethereum contract.
type GlobalCounterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GlobalCounterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GlobalCounterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GlobalCounterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GlobalCounterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GlobalCounterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GlobalCounterSession struct {
	Contract     *GlobalCounter    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// GlobalCounterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GlobalCounterCallerSession struct {
	Contract *GlobalCounterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// GlobalCounterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GlobalCounterTransactorSession struct {
	Contract     *GlobalCounterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// GlobalCounterRaw is an auto generated low-level Go binding around an Ethereum contract.
type GlobalCounterRaw struct {
	Contract *GlobalCounter // Generic contract binding to access the raw methods on
}

// GlobalCounterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GlobalCounterCallerRaw struct {
	Contract *GlobalCounterCaller // Generic read-only contract binding to access the raw methods on
}

// GlobalCounterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GlobalCounterTransactorRaw struct {
	Contract *GlobalCounterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGlobalCounter creates a new instance of GlobalCounter, bound to a specific deployed contract.
func NewGlobalCounter(address common.Address, backend bind.ContractBackend) (*GlobalCounter, error) {
	contract, err := bindGlobalCounter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &GlobalCounter{GlobalCounterCaller: GlobalCounterCaller{contract: contract}, GlobalCounterTransactor: GlobalCounterTransactor{contract: contract}, GlobalCounterFilterer: GlobalCounterFilterer{contract: contract}}, nil
}

// NewGlobalCounterCaller creates a new read-only instance of GlobalCounter, bound to a specific deployed contract.
func NewGlobalCounterCaller(address common.Address, caller bind.ContractCaller) (*GlobalCounterCaller, error) {
	contract, err := bindGlobalCounter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GlobalCounterCaller{contract: contract}, nil
}

// NewGlobalCounterTransactor creates a new write-only instance of GlobalCounter, bound to a specific deployed contract.
func NewGlobalCounterTransactor(address common.Address, transactor bind.ContractTransactor) (*GlobalCounterTransactor, error) {
	contract, err := bindGlobalCounter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GlobalCounterTransactor{contract: contract}, nil
}

// NewGlobalCounterFilterer creates a new log filterer instance of GlobalCounter, bound to a specific deployed contract.
func NewGlobalCounterFilterer(address common.Address, filterer bind.ContractFilterer) (*GlobalCounterFilterer, error) {
	contract, err := bindGlobalCounter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GlobalCounterFilterer{contract: contract}, nil
}

// bindGlobalCounter binds a generic wrapper to an already deployed contract.
func bindGlobalCounter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := GlobalCounterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GlobalCounter *GlobalCounterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GlobalCounter.Contract.GlobalCounterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GlobalCounter *GlobalCounterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GlobalCounter.Contract.GlobalCounterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GlobalCounter *GlobalCounterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GlobalCounter.Contract.GlobalCounterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GlobalCounter *GlobalCounterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GlobalCounter.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GlobalCounter *GlobalCounterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GlobalCounter.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GlobalCounter *GlobalCounterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GlobalCounter.Contract.contract.Transact(opts, method, params...)
}

// CurrentCount is a free data retrieval call binding the contract method 0xc732d201.
//
// Solidity: function currentCount() view returns(uint256)
func (_GlobalCounter *GlobalCounterCaller) CurrentCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GlobalCounter.contract.Call(opts, &out, "currentCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CurrentCount is a free data retrieval call binding the contract method 0xc732d201.
//
// Solidity: function currentCount() view returns(uint256)
func (_GlobalCounter *GlobalCounterSession) CurrentCount() (*big.Int, error) {
	return _GlobalCounter.Contract.CurrentCount(&_GlobalCounter.CallOpts)
}

// CurrentCount is a free data retrieval call binding the contract method 0xc732d201.
//
// Solidity: function currentCount() view returns(uint256)
func (_GlobalCounter *GlobalCounterCallerSession) CurrentCount() (*big.Int, error) {
	return _GlobalCounter.Contract.CurrentCount(&_GlobalCounter.CallOpts)
}

// Increment is a paid mutator transaction binding the contract method 0xd09de08a.
//
// Solidity: function increment() returns()
func (_GlobalCounter *GlobalCounterTransactor) Increment(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GlobalCounter.contract.Transact(opts, "increment")
}

// Increment is a paid mutator transaction binding the contract method 0xd09de08a.
//
// Solidity: function increment() returns()
func (_GlobalCounter *GlobalCounterSession) Increment() (*types.Transaction, error) {
	return _GlobalCounter.Contract.Increment(&_GlobalCounter.TransactOpts)
}

// Increment is a paid mutator transaction binding the contract method 0xd09de08a.
//
// Solidity: function increment() returns()
func (_GlobalCounter *GlobalCounterTransactorSession) Increment() (*types.Transaction, error) {
	return _GlobalCounter.Contract.Increment(&_GlobalCounter.TransactOpts)
}
