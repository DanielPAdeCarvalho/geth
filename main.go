package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"truswallet/blockchain"
)

func main() {
	rpcClient := blockchain.NewJSONRPCClient("https://cloudflare-eth.com")

	address := "0x71C7656EC7ab88b098defB751B7401B5f6d8976F"
	result, err := rpcClient.Call("eth_getBalance", []interface{}{address, "latest"})
	if err != nil {
		log.Fatalf("Failed to get balance: %v", err)
	}

	var balanceHex string
	err = json.Unmarshal(result, &balanceHex)
	if err != nil {
		log.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Convert hexadecimal balance to a big.Int
	balanceWei := new(big.Int)
	balanceWei, ok := balanceWei.SetString(balanceHex[2:], 16) // Remove the "0x" prefix
	if !ok {
		log.Fatal("Failed to convert balance to big.Int")
	}

	// Convert Wei to ETH
	weiToEth := new(big.Float).Quo(new(big.Float).SetInt(balanceWei), big.NewFloat(1e18))

	fmt.Printf("Balance of address %s: %s ETH\n", address, weiToEth.Text('f', 18)) // 'f' for decimal point, 18 digits after decimal
}
