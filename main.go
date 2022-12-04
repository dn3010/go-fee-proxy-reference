package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	RawURL = "https://porcini.au.rootnet.app"
	// RawURL = "http://localhost:9933"
	AlicePrivateKey = "cb6df9de1efca7a3998a8ead4e02159d5fa99c3e0d4fd6432667390bb4726854"
	BobPublicAddress = "0x25451A4de12dcCc2D166922fA938E900fCc4ED24"
	FeeProxyAddress = "0x00000000000000000000000000000000000004bb"
	SyloTokenAddress = "0xCCcCCcCC00000C64000000000000000000000000"
	SyloAssetId = 3172
	XRPAssetId = 2
)

var OneEth = new(big.Int).SetInt64(int64(1e18))

func main() {
	if err := run(); err != nil {
		log.Panicf("%v", err)
	}
}

func run() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	// make temp data directory
	tempStore, err := ioutil.TempDir("", "*-keystore")
	if err != nil {
		return fmt.Errorf("could not create test directory for keystore: %v", err)
	}
	// remove temp data directory
	defer func() {
		if err := os.RemoveAll(tempStore); err != nil {
			log.Printf("could not remove test directory for keystore: %v", err)
		}
	}()

	ks := keystore.NewKeyStore(tempStore, keystore.LightScryptN, keystore.LightScryptP)

	// private key is just for testing :)
	sk, err := ethcrypto.HexToECDSA(AlicePrivateKey)
	if err != nil {
		return fmt.Errorf("failed to derive private key: %v", err)
	}

	acct := accounts.Account{Address: ethcrypto.PubkeyToAddress(sk.PublicKey)}

	if _, err = ks.ImportECDSA(sk, ""); err != nil {
		return fmt.Errorf("failed to import private key: %v", err)
	}

	if err = ks.Unlock(acct, ""); err != nil {
		return fmt.Errorf("failed to unlock account: %v", err)
	}

	evmClient, err := ethclient.Dial(RawURL)
	if err != nil {
		return fmt.Errorf("could not create ethereum virtual client: %v", err)
	}
	chainID, err := evmClient.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("could not get chain id: %v", err)
	}
	gasPrice, err := evmClient.SuggestGasPrice(ctx)
	if err != nil {
		return fmt.Errorf("could not suggest gas price: %v", err)
	}

	initialXRPBalance, err := evmClient.BalanceAt(ctx, acct.Address, nil)
	if err != nil {
		return fmt.Errorf("failed to retrieve xrp balance: %v", err)
	}
	initialXRPBalanceNative := new(big.Int).Div(initialXRPBalance, math.BigPow(10, 12))
	initialXRPBalanceScaled := new(big.Float).Quo(new(big.Float).SetInt(initialXRPBalance), new(big.Float).SetInt(OneEth))
	log.Printf("initial Balance (XRP): %v = %v (actual)", initialXRPBalanceNative.String(), initialXRPBalanceScaled.String())

	tokenAddress := common.HexToAddress(SyloTokenAddress)
	token, err := NewSyloToken(tokenAddress, evmClient)
	if err != nil {
		return fmt.Errorf("failed to bind sylo token contract: %v", err)
	}

	initialSyloBalance, err := token.BalanceOf(nil, acct.Address)
	if err != nil {
		return fmt.Errorf("failed to retrieve sylo balance: %v", err)
	}
	initialSyloBalanceScaled := new(big.Float).Quo(new(big.Float).SetInt(initialSyloBalance), new(big.Float).SetInt(OneEth))
	log.Printf("initial Balance (SYLO): %v = %v (actual)", initialSyloBalance.String(), initialSyloBalanceScaled.String())

	feeProxyAddress := common.HexToAddress(FeeProxyAddress)
	feeProxy, err := NewFeeProxy(feeProxyAddress, evmClient)
	if err != nil {
		return fmt.Errorf("failed to bind fee proxy contract: %v", err)
	}

	// set receiver and amount sent for transfer tx data
	receiver, amount := common.HexToAddress(BobPublicAddress), big.NewInt(1)
	transferData, err := packTxData(SyloTokenMetaData, "transfer", receiver, amount)
	if err != nil {
		return fmt.Errorf("could not derive input bytes: %w", err)
	}
	
	maxFeePayment := new(big.Int).Mul(big.NewInt(2), OneEth) // 2 SYLO, TODO: Use dex rpc

	// Estimate Gas Limit for fee proxy transaction
	feeProxyData, err := packTxData(FeeProxyMetaData, "callWithFeePreferences", tokenAddress, maxFeePayment, tokenAddress, transferData)
	if err != nil {
		return fmt.Errorf("could not derive input bytes: %w", err)
	}

	// Get account nonce
	nonce, err := evmClient.PendingNonceAt(ctx, acct.Address)
	if err != nil {
		return fmt.Errorf("could not get account nonce: %v", err)
	}
	// nonce = nonce + 3
	log.Printf("Account Nonce: %v", nonce)
	
	opts, err := bind.NewKeyStoreTransactorWithChainID(ks, acct, chainID)
	if err != nil {
		return fmt.Errorf("could not create keystore transactor: %v", err)
	}
	opts.Nonce = big.NewInt(int64(nonce))
	msg := ethereum.CallMsg{
		From: opts.From,
		To:   &feeProxyAddress,
		Data: feeProxyData,
		Value: big.NewInt(0),
	}
	gasLimit, err := evmClient.EstimateGas(ctx, msg)
	if err != nil {
		return fmt.Errorf("could not estimate gas for fee proxy: %v", err)
	}

	xrpGasCost := new(big.Int).Mul(gasPrice, big.NewInt(int64(gasLimit)))
	xrpGasCostNative := new(big.Int).Div(xrpGasCost, math.BigPow(10, 12))
	xrpGasCostScaled := new(big.Float).Quo(new(big.Float).SetInt(xrpGasCost), new(big.Float).SetInt(OneEth))
	log.Printf("estimated gas cost (XRP): %v = %v (actual)", xrpGasCostNative, xrpGasCostScaled)

	syloCost, err := getSyloCost(xrpGasCostNative.Uint64(), []uint32{SyloAssetId, XRPAssetId})
	if err != nil {
		return fmt.Errorf("could not get sylo cost: %v", err)
	}

	syloCostEther := syloCost.Result.Ok[0]
	syloCostScaled := new(big.Float).Quo(new(big.Float).SetInt(&syloCostEther), new(big.Float).SetInt(OneEth))
	log.Printf("estimated gas (SYLO): %s = %s (actual)", syloCostEther.String(), syloCostScaled.String())

	opts.GasLimit = gasLimit
	opts.GasPrice = gasPrice

	feeProxySession := &FeeProxySession{
		Contract:     feeProxy,
		TransactOpts: *opts,
	}
	log.Printf("Sending Fee Proxy Transaction from %v for token=%v, target=%v", acct.Address.Hex(), tokenAddress.Hex(), tokenAddress.Hex())

	tx, err := feeProxySession.CallWithFeePreferences(tokenAddress, maxFeePayment, tokenAddress, transferData)
	if err != nil {
		return fmt.Errorf("failed to send fee proxy transaction: %v", err)
	}

	log.Printf("Sent Fee Proxy transaction: %v", tx.Hash())

	log.Printf("Waiting for tx receipt...")

	// poll for transaction receipt
	receipt, err := bind.WaitMined(ctx, evmClient, tx)
	if err != nil {
		return fmt.Errorf("failed to wait for transaction receipt: %v", err)
	}
	// log.Printf("Successfully received tx receipt. Gas used=%v", receipt.GasUsed)

	actualCost := new(big.Int).Mul(gasPrice, big.NewInt(int64(receipt.GasUsed)))
	xrpGasCostActual := new(big.Int).Div(actualCost, math.BigPow(10, 12))
	log.Printf("actual gas (XRP): %v", xrpGasCostActual)

	// XRP balance difference: 0
	finalXRPBalance, err := evmClient.BalanceAt(ctx, acct.Address, nil)
	if err != nil {
		return fmt.Errorf("failed to retrieve xrp balance: %v", err)
	}

	finalSyloBalance, err := token.BalanceOf(nil, acct.Address)
	if err != nil {
		return fmt.Errorf("failed to retrieve sylo balance: %v", err)
	}

	// XRP balance difference
	xrpBalanceDiff := new(big.Int).Sub(finalXRPBalance, initialXRPBalance)
	log.Printf("XRP balance difference: %v", xrpBalanceDiff.String())

	// Sylo balance difference
	syloBalanceDiff := new(big.Int).Sub(finalSyloBalance, initialSyloBalance)
	syloCostScaledDiff := new(big.Float).Quo(new(big.Float).SetInt(syloBalanceDiff), new(big.Float).SetInt(OneEth))
	log.Printf("Sylo balance difference: %v = %v (actual)", syloBalanceDiff.String(), syloCostScaledDiff.String())

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

// example: {"jsonrpc":"2.0","result":{"Ok":[1,1000000000000]},"id":1}
type SyloCostResponse struct {
	JSONrpc string `json:"jsonrpc"`
	Result  struct {
		Ok []big.Int `json:"Ok"`
	} `json:"result"`
	ID uint64 `json:"id"`
}

/*
 * getSyloCost performs a POST request to the dex pallet RPC to retrieve input tokens required for a swap
 * example: curl -X POST -H "Content-Type: application/json" -d '{"id":1, "jsonrpc":"2.0", "method": "dex_getAmountsIn", "params": [348615, [3172, 2]]}' http://localhost:9933
 * path max length is 2 (due to the way we parse the response)
*/
func getSyloCost(input uint64, path []uint32) (*SyloCostResponse, error) {
	// create request body
	reqBodyStr := fmt.Sprintf(`{"id":1, "jsonrpc":"2.0", "method": "dex_getAmountsIn", "params": [%d, [%d, %d]]}`, input, path[0], path[1])
	resp, err := http.Post(RawURL, "application/json", bytes.NewBuffer([]byte(reqBodyStr)))
	if err != nil {
		return nil, fmt.Errorf("could not perform POST request: %v", err)
	}
	defer resp.Body.Close()

	// print body as string
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %v", err)
	}

	// parse response
	var syloCostResponse SyloCostResponse
	err = json.Unmarshal(bodyBytes, &syloCostResponse)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response: %v", err)
	}

	return &syloCostResponse, nil
}