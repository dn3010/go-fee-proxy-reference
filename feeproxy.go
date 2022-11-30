// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package main

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
)

// FeeProxyMetaData contains all meta data concerning the FeeProxy contract.
var FeeProxyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint128\",\"name\":\"maxPayment\",\"type\":\"uint128\"},{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"}],\"name\":\"callWithFeePreferences\",\"outputs\":[],\"type\":\"function\"}]",
}

// FeeProxyABI is the input ABI used to generate the binding from.
// Deprecated: Use FeeProxyMetaData.ABI instead.
var FeeProxyABI = FeeProxyMetaData.ABI

// FeeProxy is an auto generated Go binding around an Ethereum contract.
type FeeProxy struct {
	FeeProxyCaller     // Read-only binding to the contract
	FeeProxyTransactor // Write-only binding to the contract
	FeeProxyFilterer   // Log filterer for contract events
}

// FeeProxyCaller is an auto generated read-only Go binding around an Ethereum contract.
type FeeProxyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeeProxyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FeeProxyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeeProxyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FeeProxyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeeProxySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FeeProxySession struct {
	Contract     *FeeProxy         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FeeProxyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FeeProxyCallerSession struct {
	Contract *FeeProxyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// FeeProxyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FeeProxyTransactorSession struct {
	Contract     *FeeProxyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// FeeProxyRaw is an auto generated low-level Go binding around an Ethereum contract.
type FeeProxyRaw struct {
	Contract *FeeProxy // Generic contract binding to access the raw methods on
}

// FeeProxyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FeeProxyCallerRaw struct {
	Contract *FeeProxyCaller // Generic read-only contract binding to access the raw methods on
}

// FeeProxyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FeeProxyTransactorRaw struct {
	Contract *FeeProxyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFeeProxy creates a new instance of FeeProxy, bound to a specific deployed contract.
func NewFeeProxy(address common.Address, backend bind.ContractBackend) (*FeeProxy, error) {
	contract, err := bindFeeProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FeeProxy{FeeProxyCaller: FeeProxyCaller{contract: contract}, FeeProxyTransactor: FeeProxyTransactor{contract: contract}, FeeProxyFilterer: FeeProxyFilterer{contract: contract}}, nil
}

// NewFeeProxyCaller creates a new read-only instance of FeeProxy, bound to a specific deployed contract.
func NewFeeProxyCaller(address common.Address, caller bind.ContractCaller) (*FeeProxyCaller, error) {
	contract, err := bindFeeProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FeeProxyCaller{contract: contract}, nil
}

// NewFeeProxyTransactor creates a new write-only instance of FeeProxy, bound to a specific deployed contract.
func NewFeeProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*FeeProxyTransactor, error) {
	contract, err := bindFeeProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FeeProxyTransactor{contract: contract}, nil
}

// NewFeeProxyFilterer creates a new log filterer instance of FeeProxy, bound to a specific deployed contract.
func NewFeeProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*FeeProxyFilterer, error) {
	contract, err := bindFeeProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FeeProxyFilterer{contract: contract}, nil
}

// bindFeeProxy binds a generic wrapper to an already deployed contract.
func bindFeeProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(FeeProxyABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FeeProxy *FeeProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FeeProxy.Contract.FeeProxyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FeeProxy *FeeProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FeeProxy.Contract.FeeProxyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FeeProxy *FeeProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FeeProxy.Contract.FeeProxyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FeeProxy *FeeProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FeeProxy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FeeProxy *FeeProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FeeProxy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FeeProxy *FeeProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FeeProxy.Contract.contract.Transact(opts, method, params...)
}

// CallWithFeePreferences is a paid mutator transaction binding the contract method 0x255a3432.
//
// Solidity: function callWithFeePreferences(address asset, uint128 maxPayment, address target, bytes input) returns()
func (_FeeProxy *FeeProxyTransactor) CallWithFeePreferences(opts *bind.TransactOpts, asset common.Address, maxPayment *big.Int, target common.Address, input []byte) (*types.Transaction, error) {
	return _FeeProxy.contract.Transact(opts, "callWithFeePreferences", asset, maxPayment, target, input)
}

// CallWithFeePreferences is a paid mutator transaction binding the contract method 0x255a3432.
//
// Solidity: function callWithFeePreferences(address asset, uint128 maxPayment, address target, bytes input) returns()
func (_FeeProxy *FeeProxySession) CallWithFeePreferences(asset common.Address, maxPayment *big.Int, target common.Address, input []byte) (*types.Transaction, error) {
	return _FeeProxy.Contract.CallWithFeePreferences(&_FeeProxy.TransactOpts, asset, maxPayment, target, input)
}

// CallWithFeePreferences is a paid mutator transaction binding the contract method 0x255a3432.
//
// Solidity: function callWithFeePreferences(address asset, uint128 maxPayment, address target, bytes input) returns()
func (_FeeProxy *FeeProxyTransactorSession) CallWithFeePreferences(asset common.Address, maxPayment *big.Int, target common.Address, input []byte) (*types.Transaction, error) {
	return _FeeProxy.Contract.CallWithFeePreferences(&_FeeProxy.TransactOpts, asset, maxPayment, target, input)
}
