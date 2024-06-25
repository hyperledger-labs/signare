package rpcinfra

import (
	"context"
	"encoding/json"

	"github.com/hyperledger-labs/signare/app/pkg/infra/httpinfra"
	"github.com/hyperledger-labs/signare/app/pkg/infra/rpcinfra/rpcerrors"

	"github.com/gorilla/mux"
)

const (
	SupportedRPCVersion = "2.0"
)

// RPCHandler defines JSON-RPC method handler functions
type RPCHandler func(ctx context.Context, request RPCRequest) (any, *rpcerrors.RPCError)

// DefaultRPCRouter defines a JSON-RPC router that handles multiple requests.
type DefaultRPCRouter struct {
	// Router for HTTP
	router *mux.Router
	// rpcHandlers registered JSON-RPC method handlers
	rpcHandlers map[string]RPCHandler
	// defaultRPCInfraResponseHandler handles the RPC responses of the server
	defaultRPCInfraResponseHandler httpinfra.HTTPResponseHandler
}

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

// RPCResponse represents a JSON-RPC response object as defined in: https://www.jsonrpc.org/specification
type RPCResponse struct {
	// RPCVersion specifies the version of the JSON-RPC protocol. Must be exactly "2.0".
	RPCVersion string `json:"jsonrpc"`
	// ID contains the client established request id or null.
	ID any `json:"id"`
	// Error contains the error if there was an error triggered during the request.
	Error *rpcerrors.RPCError `json:"error,omitempty"`
	// Result contains the result of the called method.
	// It MUST be defined in a successful response, and it MUST not be defined if there was an error triggered
	// during the request.
	Result any `json:"result,omitempty"`
}

// JSONRPCParams defines methods for processing JSON-RPC request parameters.
type JSONRPCParams interface {
	// SetParamsFrom validates the request parameters against the provided interface and if the validation is correct,
	// completes the JSONRPCParams with the interface values.
	SetParamsFrom([]any) error
	// ValidateParams checks if all the parameters are defined in the JSONRPCParams struct and if the validation fails,
	// it returns an error.
	ValidateParams() error
}

// ProcessParams processes the params data structure from the RPCRequest.
func ProcessParams(reqParams json.RawMessage, rpcParams JSONRPCParams) *rpcerrors.RPCError {
	if err := json.Unmarshal(reqParams, rpcParams); err != nil {
		// If the unmarshall fails, we try to unmarshall it into an interface and set the JSONRPCParams from there.
		posParams := make([]any, 0)
		if err = json.Unmarshal(reqParams, &posParams); err != nil {
			return rpcerrors.NewInvalidParamsFromErr(err)
		}
		if err = rpcParams.SetParamsFrom(posParams); err != nil {
			return rpcerrors.NewInvalidParamsFromErr(err)
		}
	}
	return nil
}
