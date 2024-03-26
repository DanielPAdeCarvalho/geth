package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type BlockNumberResponse struct {
	BlockNumber string `json:"blockNumber"`
}

func TestJSONRPCCall(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Successful JSON-RPC response
		resp := jsonRPCResponse{
			JSONRPC: "2.0",
			Result:  json.RawMessage(`{"blockNumber":"0x10d4f"}`),
		}

		err := json.NewEncoder(w).Encode(resp)
		if err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	// Create a new JSONRPC client pointing to the mock server
	client := NewJSONRPC(server.URL)

	// Test Call
	result, err := client.Call("eth_blockNumber", nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var blockNumberResponse BlockNumberResponse
	if err := json.Unmarshal(result, &blockNumberResponse); err != nil {
		t.Fatalf("Expected to unmarshal response without error, got %v", err)
	}

	expectedBlockNumber := "0x10d4f"
	if blockNumberResponse.BlockNumber != expectedBlockNumber {
		t.Errorf("Expected block number %s, got %s", expectedBlockNumber, blockNumberResponse.BlockNumber)
	}
}
