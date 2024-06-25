package rpccontextdefinition

import "encoding/json"

// RPCRequest represents a JSON-RPC request object as defined in: https://www.jsonrpc.org/specification
type RPCRequest struct {
	// RPCVersion specifies the version of the JSON-RPC protocol. Must be exactly "2.0".
	RPCVersion string `json:"jsonrpc"`
	// ID defines a unique identifier established by the client.
	ID any `json:"id"`
	// Method defines the name of the method to be invoked.
	Method string `json:"method"`
	// Params defines a structured value that holds the parameter values to be used during the invocation of the method.
	Params json.RawMessage `json:"params"`
}
