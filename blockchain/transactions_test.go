package blockchain

import (
	"bytes"
	"math/big"
	"testing"
	"truswallet/storage"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

// Mock or prepare the EthereumParser with an instance of InMemoryStorage and a mock JSONRPCClient.
func TestFetchBlockTransactionsUsingStorage(t *testing.T) {
	// Mock JSONRPCClient setup to return a predefined block with transactions
	blockNumber := uint64(12345)
	subscribedAddress1 := "0xsubscribed1"
	subscribedAddress2 := "0xsubscribed2"

	// Prepare transactions for the test block
	tx1 := types.NewTransaction(1, common.HexToAddress(subscribedAddress1), big.NewInt(100), 21000, big.NewInt(10), nil)
	tx2 := types.NewTransaction(2, common.HexToAddress("0xto2"), big.NewInt(200), 21000, big.NewInt(20), nil)

	// Encode transactions to RLP for the mock response
	var buf bytes.Buffer
	err := rlp.Encode(&buf, []*types.Transaction{tx1, tx2})
	if err != nil {
		t.Fatalf("Failed to RLP encode transactions: %v", err)
	}

	mockClient := &MockJSONRPCClient{
		// Simulate the response for eth_getBlockByNumber
		Response: buf.Bytes(),
	}

	storage := storage.NewInMemoryStorage() // Assuming this is your storage implementation
	parser := NewEthereumParser(mockClient, storage)
	parser.Subscribe(subscribedAddress1)
	parser.Subscribe(subscribedAddress2)

	// Execute FetchBlockTransactions to simulate fetching and storing
	parser.FetchBlockTransactions(int(blockNumber))

	// Verify that the transactions related to the subscribed addresses are correctly saved
	txs1, err := storage.GetTransactionsByAddress(subscribedAddress1)
	if err != nil || len(txs1) != 1 {
		t.Errorf("Expected to find 1 transaction for address %s, found %d", subscribedAddress1, len(txs1))
	}

	// Assuming you also track transactions from the sender's perspective, you could also check for outbound transactions for subscribedAddress2
	txs2, err := storage.GetTransactionsByAddress(subscribedAddress2)
	if err != nil || len(txs2) != 0 { // Adjust this based on your actual logic for tracking outbound transactions
		t.Errorf("Expected to find 0 transactions for address %s, found %d", subscribedAddress2, len(txs2))
	}
}
