package blockchain

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

// FetchBlockTransactions fetches transactions from a specific block and updates the storage.
func (e *EthereumParser) FetchBlockTransactions(blockNumber int) {
	blockNumberHex := fmt.Sprintf("0x%x", blockNumber)
	result, err := e.client.Call("eth_getBlockByNumber", []interface{}{blockNumberHex, true})
	if err != nil {
		log.Printf("Failed to fetch block %d: %v", blockNumber, err)
		return
	}

	// Parse the result to get transactions
	var block struct {
		Transactions []json.RawMessage `json:"transactions"`
	}
	if err := json.Unmarshal(result, &block); err != nil {
		log.Printf("Failed to unmarshal block: %v", err)
		return
	}

	// Convert transactions to go-ethereum's types.Transaction format
	var transactions []types.Transaction
	for _, txData := range block.Transactions {
		var tx types.Transaction
		if err := rlp.DecodeBytes(hexutil.MustDecode(string(txData)), &tx); err != nil {
			log.Printf("Failed to decode transaction: %v", err)
			continue
		}
		transactions = append(transactions, tx)
	}

	// Save the transactions using InMemoryStorage
	if err := e.storage.SaveTransactions(uint64(blockNumber), transactions); err != nil {
		log.Printf("Failed to save transactions for block %d: %v", blockNumber, err)
	}
}
