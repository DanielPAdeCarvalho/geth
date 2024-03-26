package main

import (
	"encoding/json"
	"net/http"
	"truswallet/blockchain"
)

// Define a global Parser instance
var ethereumParser blockchain.Parser = &blockchain.EthereumParser{}

func getCurrentBlock(w http.ResponseWriter, r *http.Request) {
	block := ethereumParser.GetCurrentBlock()
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

	http.ListenAndServe(":8080", nil)
}
