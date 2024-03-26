package blockchain

import (
	"encoding/json"
	"fmt"
	"log"
)

// FetchBlockTransactions fetches transactions from a specific block and updates the transactions map.
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

	e.mutex.Lock()
	defer e.mutex.Unlock()

	// Ensure this loop is within a method of EthereumParser where e.mutex is defined.
	for _, tx := range block.Transactions {
		// Process outbound transactions
		if _, ok := e.subscribedAddresses[tx.From]; ok {
			addrTxs := e.transactions[tx.From]
			addrTxs.Outbound = append(addrTxs.Outbound, tx)
			e.transactions[tx.From] = addrTxs
		}

		// Process inbound transactions
		if _, ok := e.subscribedAddresses[tx.To]; ok {
			addrTxs := e.transactions[tx.To]
			addrTxs.Inbound = append(addrTxs.Inbound, tx)
			e.transactions[tx.To] = addrTxs
		}
	}

}
