package blockchain

import (
	"context"
	"math/big"

	"truswallet/store" // Import your store package

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ParseBlock parses transactions from a given block number, filters transactions from subscribed addresses,
// and saves them using the provided storage.
func ParseBlock(client *ethclient.Client, blockNumber uint64, subscribedAddresses map[string]bool, storage store.Storage) error {
	block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(blockNumber)))
	if err != nil {
		return err
	}

	var relevantTxs []types.Transaction
	for _, tx := range block.Transactions() {
		// Ensure the transaction has a to address (i.e., is not a contract creation transaction)
		if tx.To() != nil && subscribedAddresses[tx.To().Hex()] {
			relevantTxs = append(relevantTxs, *tx)
		}
	}

	// Save the filtered transactions to storage
	if len(relevantTxs) > 0 {
		err = storage.SaveTransactions(blockNumber, relevantTxs)
		if err != nil {
			return err
		}
	}

	return nil
}
