package blockchain

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sync"
	"truswallet/client"
)

// Parser defines an interface for parsing blockchain data.
type Parser interface {
	// last parsed block
	GetCurrentBlock() (int, error)

	// add address to observer
	Subscribe(address string) bool

	// list of inbound or outbound transactions for an address
	GetTransactions(address string) []Transaction
}

type AddressTransactions struct {
	Inbound  []Transaction
	Outbound []Transaction
}

type Transaction struct {
	Hash     string // Transaction hash
	From     string // Sender address
	To       string // Recipient address
	Value    string // Amount transferred
	GasPrice string // Gas price
	GasUsed  uint64 // Gas used
}

type EthereumParser struct {
	subscribedAddresses map[string]bool
	// transactions stores a list of transactions for each subscribed address.
	// The key is the address, and the value is a slice of Transactions.
	transactions map[string]AddressTransactions
	client       client.JSONRPCClient // JSON-RPC client to interact with the blockchain node
	mutex        sync.RWMutex         // Protects subscribedAddresses and transactions
}

func NewEthereumParser(client client.JSONRPCClient) *EthereumParser {
	return &EthereumParser{
		subscribedAddresses: make(map[string]bool),
		transactions:        make(map[string]AddressTransactions),
		client:              client,
	}
}

// GetCurrentBlock calls the JSON-RPC eth_blockNumber method to get the latest block number.
func (p *EthereumParser) GetCurrentBlock() (int, error) {
	result, err := p.client.Call("eth_blockNumber", nil)
	if err != nil {
		return 0, fmt.Errorf("failed to get current block number: %w", err)
	}

	var blockNumberHex string
	if err := json.Unmarshal(result, &blockNumberHex); err != nil {
		return 0, fmt.Errorf("failed to unmarshal block number: %w", err)
	}

	// Convert hex string to big.Int
	blockNumber, ok := new(big.Int).SetString(blockNumberHex[2:], 16)
	if !ok {
		return 0, fmt.Errorf("failed to parse block number from %s", blockNumberHex)
	}

	return int(blockNumber.Int64()), nil // Convert to int for the interface
}

// Subscribe adds an Ethereum address to the list of subscribed addresses.
// Returns true if the address was added, false if it was already subscribed.
func (p *EthereumParser) Subscribe(address string) bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	// Check if the address is already subscribed
	if _, exists := p.subscribedAddresses[address]; exists {
		return false
	}

	// Subscribe the address
	p.subscribedAddresses[address] = true

	return true
}

func (p *EthereumParser) GetTransactions(address string) []Transaction {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	var transactions []Transaction

	if addrTxs, exists := p.transactions[address]; exists {
		// Combine inbound and outbound transactions
		transactions = append(transactions, addrTxs.Inbound...)
		transactions = append(transactions, addrTxs.Outbound...)
	}

	return transactions
}
