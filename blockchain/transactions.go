package blockchain

import (
	"encoding/json"
	"fmt"
	"log"
)

// FetchBlockTransactions fetches transactions from a specific block and updates the transactions map.
// This is a simplified example. In practice, you'd need to handle pagination or batch requests for blocks with many transactions.
func (e *EthereumParser) FetchBlockTransactions(blockNumber int) {
	blockNumberHex := fmt.Sprintf("0x%x", blockNumber)
	result, err := e.client.Call("eth_getBlockByNumber", []interface{}{blockNumberHex, true})
	if err != nil {
		log.Printf("Failed to fetch block %d: %v", blockNumber, err)
		return
	}

	var block struct {
		Transactions []Transaction `json:"transactions"`
	}
	if err := json.Unmarshal(result, &block); err != nil {
		log.Printf("Failed to unmarshal block: %v", err)
		return
	}

	for _, tx := range block.Transactions {
		if _, ok := e.subscribedAddresses[tx.From]; ok || e.subscribedAddresses[tx.To] {
			e.transactions[tx.To] = append(e.transactions[tx.To], tx)
		}
	}
}
