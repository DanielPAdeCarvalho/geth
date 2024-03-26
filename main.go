package main

import (
	"log"
	"truswallet/blockchain"
	"truswallet/store"
)

func main() {
	client := blockchain.NewClient("https://mainnet.infura.io/v3/YOUR_PROJECT_ID")
	storage := store.NewInMemoryStorage() // Initialize in-memory storage

	// Example: Subscribing to an address
	store.Subscribe("0xAddressHere")

	// TODO: Adjusted call to ParseBlock with storage passed as an argument
	err := blockchain.ParseBlock(client, 12345678, store.SubscribedAddresses, storage)
	if err != nil {
		log.Fatalf("Failed to parse block: %v", err)
	}

	// Optionally, you could retrieve and process transactions for a specific address from storage here...
}
