package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"truswallet/blockchain"
	"truswallet/client"
	"truswallet/storage"
)

// Define a global Parser instance
var ethereumParser blockchain.Parser

// Initialize JSONRPCClient and Storage
func init() {
	ethClient := client.NewJSONRPC("https://cloudflare-eth.com/")
	ethStorage := storage.NewInMemoryStorage()

	ethereumParser = blockchain.NewEthereumParser(ethClient, ethStorage)
}

func getCurrentBlock(w http.ResponseWriter, r *http.Request) {
	block, err := ethereumParser.GetCurrentBlock()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]int{"currentBlock": block})
}

func subscribeAddress(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	success := ethereumParser.Subscribe(address)
	json.NewEncoder(w).Encode(map[string]bool{"subscribed": success})
}

func getTransactions(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	transactions := ethereumParser.GetTransactions(address)
	json.NewEncoder(w).Encode(transactions)
}

func main() {
	http.HandleFunc("/currentBlock", getCurrentBlock)
	http.HandleFunc("/subscribe", subscribeAddress)
	http.HandleFunc("/transactions", getTransactions)
	fmt.Println("Server is listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
