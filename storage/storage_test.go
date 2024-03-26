package storage

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

// Helper function to create a mock transaction.
func mockTransaction(to common.Address) *types.Transaction {
	return types.NewTransaction(0, to, nil, 0, nil, nil)
}

func TestSaveAndGetTransactions(t *testing.T) {
	storage := NewInMemoryStorage()

	// Mock address and transactions
	address := common.HexToAddress("0xABC")
	tx1 := mockTransaction(address)
	tx2 := mockTransaction(address)

	// Save transactions
	err := storage.SaveTransactions(1, []types.Transaction{*tx1, *tx2})
	assert.NoError(t, err)

	// Retrieve transactions by address
	txs, err := storage.GetTransactionsByAddress(address.Hex())
	assert.NoError(t, err)
	assert.Len(t, txs, 2)

	// Check if the transactions match
	assert.Equal(t, *tx1, txs[0])
	assert.Equal(t, *tx2, txs[1])
}
