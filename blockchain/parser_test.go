package blockchain

import (
	"encoding/json"
	"testing"
)

// MockJSONRPCClient is a mock of the JSONRPCClient for testing purposes.
type MockJSONRPCClient struct {
	Response json.RawMessage
	Error    error
}

func (m *MockJSONRPCClient) Call(method string, params []interface{}) (json.RawMessage, error) {
	return m.Response, m.Error
}

func TestGetCurrentBlock(t *testing.T) {
	// Setup a successful response
	mockClient := &MockJSONRPCClient{
		Response: json.RawMessage(`"0x5BAD55"`),
	}
	parser := NewEthereumParser(mockClient)

	blockNumber, err := parser.GetCurrentBlock()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedBlockNumber := 6008149 // Example block number corresponding to "0x5BAD55"
	if blockNumber != expectedBlockNumber {
		t.Errorf("Expected block number %d, got %d", expectedBlockNumber, blockNumber)
	}
}

func TestSubscribe(t *testing.T) {
	parser := NewEthereumParser(nil) // Client is not used in Subscribe

	address := "0x123"
	subscribed := parser.Subscribe(address)
	if !subscribed {
		t.Errorf("Expected true, got %v", subscribed)
	}

	// Test subscribing the same address again should return false
	subscribedAgain := parser.Subscribe(address)
	if subscribedAgain {
		t.Errorf("Expected false for subscribing the same address, got %v", subscribedAgain)
	}
}

func TestGetTransactions(t *testing.T) {
	// Setup a mock JSONRPC client (not used in this test, but required for initialization)
	mockClient := &MockJSONRPCClient{}

	// Initialize the parser with the mock client
	parser := NewEthereumParser(mockClient)

	// Manually subscribe an address and add transactions for testing
	testAddress := "0xTestAddress"
	parser.Subscribe(testAddress)

	inboundTx := Transaction{Hash: "0xInbound", From: "0xFrom", To: testAddress, Value: "100"}
	outboundTx := Transaction{Hash: "0xOutbound", From: testAddress, To: "0xTo", Value: "50"}

	// Directly manipulate the parser's state for the test (bypassing the FetchBlockTransactions method)
	parser.transactions[testAddress] = AddressTransactions{
		Inbound:  []Transaction{inboundTx},
		Outbound: []Transaction{outboundTx},
	}

	// Execution: Retrieve transactions for the test address
	transactions := parser.GetTransactions(testAddress)

	// Verification: Check if the returned transactions match the inserted ones
	if len(transactions) != 2 {
		t.Fatalf("Expected 2 transactions, got %d", len(transactions))
	}

	// More detailed checks can be added to verify each transaction's fields
	if transactions[0].Hash != inboundTx.Hash || transactions[1].Hash != outboundTx.Hash {
		t.Errorf("Transactions do not match the expected transactions")
	}
}
