package blockchain

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"sync"
	"truswallet/client"
	"truswallet/storage"

	"github.com/ethereum/go-ethereum/core/types"
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
	Hash     string   // Transaction hash
	To       string   // Recipient address
	Value    *big.Int // Amount transferred
	GasPrice *big.Int // Gas price
	GasUsed  uint64   // Gas used
}

type EthereumParser struct {
	storage storage.Storage
	client  client.JSONRPCClient // JSON-RPC client to interact with the blockchain node
	mutex   sync.RWMutex         // Protects subscribedAddresses and transactions
}

func NewEthereumParser(client client.JSONRPCClient, storage storage.Storage) *EthereumParser {
	return &EthereumParser{
		client:  client,
		storage: storage,
		mutex:   sync.RWMutex{},
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
	// Locking is handled within the storage package if needed
	if _, exists := storage.SubscribedAddresses[address]; exists {
		return false
	}
	storage.Subscribe(address)
	return true
}

// GetTransactions retrieves transactions for a given address from the storage.
func (p *EthereumParser) GetTransactions(address string) []Transaction {
	// Fetch transactions from storage
	txs, err := p.storage.GetTransactionsByAddress(address)
	if err != nil {
		log.Printf("Failed to get transactions for address %s: %v", address, err)
		return nil
	}

	// Convert types.Transaction (from go-ethereum) to custom Transaction type.
	var transactions []Transaction
	for _, tx := range txs {
		customTx := ConvertToCustomTransactionType(tx)
		transactions = append(transactions, customTx)
	}

	return transactions
}

func ConvertToCustomTransactionType(tx types.Transaction) Transaction {
	// tx.To() can be nil for contract creation transactions
	to := ""
	if tx.To() != nil {
		to = tx.To().Hex()
	}

	value := tx.Value()
	gasPrice := tx.GasPrice()
	gasUsed := tx.Gas()

	return Transaction{
		Hash:     tx.Hash().Hex(),
		To:       to,
		Value:    value,
		GasPrice: gasPrice,
		GasUsed:  gasUsed,
	}
}
