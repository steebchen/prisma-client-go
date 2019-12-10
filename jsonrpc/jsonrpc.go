// Package jsonrpc enables communication with Prisma
package jsonrpc

import (
	"encoding/json"
)

// Request sets a generic JSONRPC request, which wraps information and params.
type Request struct {
	// JSONRPC describes the version of the JSON RPC protocol. Defaults to `2.0`.
	JSONRPC string `json:"jsonrpc"`
	// ID identifies a unique request.
	ID int `json:"id"`
	// Method describes the intention of the request.
	Method string `json:"method"`
	// Params contains the payload of the request. Usually parsed into a specific struct for further processing.
	Params json.RawMessage `json:"params"`
}

// Response sets a generic JSONRPC response, which wraps information and a result.
type Response struct {
	// JSONRPC describes the version of the JSON RPC protocol. Defaults to `2.0`.
	JSONRPC string `json:"jsonrpc"`
	// ID identifies a unique request.
	ID int `json:"id"`
	// Result contains the payload of the response.
	Result interface{} `json:"result"`
}

// NewResponse forms a new JSON RPC response to reply to the Prisma CLI commands
func NewResponse(id int, result interface{}) Response {
	return Response{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
}
