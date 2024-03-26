package blockchain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type JSONRPCClient struct {
	url string
}

func NewJSONRPCClient(url string) *JSONRPCClient {
	return &JSONRPCClient{url: url}
}

type jsonRPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

type jsonRPCResponse struct {
	ID      int             `json:"id"`
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   *jsonRPCError   `json:"error"`
}

type jsonRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error implements the error interface for jsonRPCError.
func (e *jsonRPCError) Error() string {
	return fmt.Sprintf("RPC Error %d: %s", e.Code, e.Message)
}

func (c *JSONRPCClient) Call(method string, params []interface{}) (json.RawMessage, error) {
	request := jsonRPCRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(c.url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response jsonRPCResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, &jsonRPCError{Code: response.Error.Code, Message: response.Error.Message}
	}

	return response.Result, nil
}
