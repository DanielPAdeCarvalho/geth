package blockchain

import (
	"testing"
)

func TestFetchBlockTransactions(t *testing.T) {
	// Example block data with transactions
	blockData := `{"transactions":[{"hash":"0xhash1","from":"0xfrom1","to":"0xsubscribed1","value":"100","gasPrice":"10","gasUsed":21000},{"hash":"0xhash2","from":"0xsubscribed2","to":"0xto2","value":"200","gasPrice":"20","gasUsed":21000}]}`
	mockClient := &MockJSONRPCClient{
		Response: []byte(blockData),
	}
	parser := NewEthereumParser(mockClient)

	// Subscribe to addresses to test both inbound and outbound transactions
	subscribedAddress1 := "0xsubscribed1"
	subscribedAddress2 := "0xsubscribed2"
	parser.Subscribe(subscribedAddress1)
	parser.Subscribe(subscribedAddress2)

	// Execute FetchBlockTransactions
	blockNumber := 12345
	parser.FetchBlockTransactions(blockNumber)

	// Assertions
	if len(parser.transactions[subscribedAddress1].Inbound) != 1 || len(parser.transactions[subscribedAddress2].Outbound) != 1 {
		t.Errorf("Transactions were not processed correctly")
	}

}
