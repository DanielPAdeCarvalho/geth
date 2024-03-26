package store

import (
	"sync"

	"github.com/ethereum/go-ethereum/core/types"
)

// Storage interface defines the methods to be implemented by any storage backend.
type Storage interface {
	SaveTransactions(blockNumber uint64, transactions []types.Transaction) error
	GetTransactionsByAddress(address string) ([]types.Transaction, error)
}

// Ensure the in-memory implementation satisfies the Storage interface.
var _ Storage = (*InMemoryStorage)(nil)

// InMemoryStorage is an in-memory storage for transactions.
type InMemoryStorage struct {
	mu           sync.Mutex                      // protects the following
	transactions map[uint64][]types.Transaction  // maps block number to transactions
	addressIndex map[string][]*types.Transaction // maps address to transactions
}

// NewInMemoryStorage creates a new instance of InMemoryStorage.
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		transactions: make(map[uint64][]types.Transaction),
		addressIndex: make(map[string][]*types.Transaction),
	}
}

// SaveTransactions saves transactions to the in-memory storage.
func (s *InMemoryStorage) SaveTransactions(blockNumber uint64, transactions []types.Transaction) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, tx := range transactions {
		// Add the transaction to the block's list of transactions.
		s.transactions[blockNumber] = append(s.transactions[blockNumber], tx)

		// Index the transaction by the recipient address.
		addr := tx.To().Hex()
		s.addressIndex[addr] = append(s.addressIndex[addr], &tx)
	}
	return nil
}

// GetTransactionsByAddress returns all transactions for a given address.
func (s *InMemoryStorage) GetTransactionsByAddress(address string) ([]types.Transaction, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if txs, ok := s.addressIndex[address]; ok {
		// Make a copy to avoid returning a reference to the internal slice.
		copy := make([]types.Transaction, len(txs))
		for i, txPtr := range txs {
			copy[i] = *txPtr
		}
		return copy, nil
	}
	return nil, nil // No transactions found for the address.
}
