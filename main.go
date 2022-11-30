package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func run() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	// make temp data directory
	tempStore, err := ioutil.TempDir("", "*-keystore")

	if err != nil {
		return fmt.Errorf("could not create test directory for keystore: %v", err)
	}
	ks := keystore.NewKeyStore(tempStore, keystore.LightScryptN, keystore.LightScryptP)

	evmClient, err := ethclient.Dial("https://porcini.au.rootnet.app")
	if err != nil {
		return fmt.Errorf("could not create ethereum virtual client: %v", err)
	}

	// private key is just for testing :)
	sk, err := ethcrypto.HexToECDSA("cb6df9de1efca7a3998a8ead4e02159d5fa99c3e0d4fd6432667390bb4726854")
	if err != nil {
		return fmt.Errorf("failed to derive private key: %v", err)
	}
	acct := accounts.Account{Address: ethcrypto.PubkeyToAddress(sk.PublicKey)}

	log.Printf("Sending transaction from: %v", acct.Address.Hex())

	_, err = ks.ImportECDSA(sk, "")
	if err != nil {
		return fmt.Errorf("failed to import private key: %v", err)
	}
	err = ks.Unlock(acct, "")
	if err != nil {
		return fmt.Errorf("failed to unlock account: %v", err)
	}

	chainID, err := evmClient.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("could not get chain id: %v", err)
	}

	opts, err := bind.NewKeyStoreTransactorWithChainID(ks, acct, chainID)
	if err != nil {
		return fmt.Errorf("could not create keystore transactor: %v", err)
	}

	tokenAddress := common.HexToAddress("0xCCcCCcCC00000C64000000000000000000000000")
	feeProxyAddress := common.HexToAddress("0x00000000000000000000000000000000000004bb")

	token, err := NewSyloToken(tokenAddress, evmClient)
	if err != nil {
		return fmt.Errorf("failed to bind sylo token contract: %v", err)
	}

	syloBalance, err := token.BalanceOf(nil, acct.Address)
	if err != nil {
		return fmt.Errorf("failed to retrieve sylo balance: %v", err)
	}

	xrpBalance, err := evmClient.BalanceAt(ctx, acct.Address, nil)
	if err != nil {
		return fmt.Errorf("failed to retrieve xrp balance: %v", err)
	}

	log.Printf("Account XRP Balance: %v", xrpBalance.String())
	log.Printf("Account Sylo Balance: %v", syloBalance.String())

	feeProxy, err := NewFeeProxy(feeProxyAddress, evmClient)
	if err != nil {
		return fmt.Errorf("failed to bind fee proxy contract: %v", err)
	}

	// receiver of transfer
	receiver := common.HexToAddress("0x25451A4de12dcCc2D166922fA938E900fCc4ED24")
	// transfer amount
	amount := big.NewInt(1)

	transferData, err := packTxData(SyloTokenMetaData, "transfer", receiver, amount)
	if err != nil {
		return fmt.Errorf("could not derive input bytes: %w", err)
	}

	ETH := new(big.Int).SetInt64(int64(1e18))
	maxFeePayment := new(big.Int).Mul(big.NewInt(100000), ETH) // 10000 SYLO, TODO: Use dex rpc

	// Estimate Gas Limit for fee proxy transaction
	feeProxyData, err := packTxData(FeeProxyMetaData, "callWithFeePreferences", tokenAddress, tokenAddress, transferData)
	msg := ethereum.CallMsg{
		From: opts.From,
		To:   &feeProxyAddress,
		Data: feeProxyData,
	}
	gasLimit, err := evmClient.EstimateGas(ctx, msg)
	if err != nil {
		return fmt.Errorf("could not estimate gas for fee proxy: %v", err)
	}

	log.Printf("Estimated gas limit for fee proxy transaction %v \n", gasLimit)

	opts.GasLimit = gasLimit
	opts.GasLimit = 250000 // The estimation is 21000, but for some reason needs to be a larger number else we get InvalidTransaction:Custom(3)

	feeProxySession := &FeeProxySession{
		Contract:     feeProxy,
		TransactOpts: *opts,
	}

	log.Printf("Sending Fee Proxy Transaction for token=%v, target=%v", tokenAddress.Hex(), tokenAddress.Hex())

	tx, err := feeProxySession.CallWithFeePreferences(tokenAddress, maxFeePayment, tokenAddress, transferData)
	if err != nil {
		return fmt.Errorf("failed to send fee proxy transaction: %v", err)
	}

	log.Printf("Sent Fee Proxy transaction: %v", tx.Hash())

	log.Printf("Waiting for tx receipt...")

	receipt, err := evmClient.TransactionReceipt(ctx, tx.Hash())
	if err != nil {
		return fmt.Errorf("failed to wait for tx receipt: %v", err)
	}

	log.Printf("Successfully received tx receipt. Gas used=%v", receipt.GasUsed)

	return nil
}

func packTxData(contractMetadata *bind.MetaData, method string, params ...interface{}) ([]byte, error) {
	abi, err := contractMetadata.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("could not get contract abi: %v", err)
	}
	input, err := abi.Pack(method, params...)
	if err != nil {
		return nil, fmt.Errorf("could not pack method (%s): %v", method, err)
	}
	return input, nil
}

func main() {
	err := run()

	if err != nil {
		log.Panicf("%v", err)
	}
}
