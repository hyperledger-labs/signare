package pep

import (
	"encoding/json"

	"github.com/hyperledger-labs/signare/app/pkg/entities/address"
)

// AuthorizeUserInput attributes to authorize a user
type AuthorizeUserInput struct {
	// UserID the ID of the user
	UserID string
	// ApplicationID the ID of the application
	ApplicationID *string
	// ActionID the ID of the action
	ActionID string
}

// AuthorizeUserOutput the result of the user authorization
type AuthorizeUserOutput struct {
}

// AuthorizeAccountUserInput attributes to authorize a user to use an account
type AuthorizeAccountUserInput struct {
	// UserID the ID of the user
	UserID string
	// ApplicationID the ID of the application
	ApplicationID string
	// Address related to the account to be authorized
	Address address.Address
}

// AuthorizeAccountUserOutput the result of the account usage authorization
type AuthorizeAccountUserOutput struct {
}

// AuthorizeAccountRPCBody represents a JSON-RPC request object as defined in: https://www.jsonrpc.org/specification
type AuthorizeAccountRPCBody struct {
	// Method defines the name of the method to be invoked.
	Method string `json:"method"`
	// ID defines a unique identifier established by the client.
	ID any `json:"id"`
	// Params defines a structured value that holds the parameter values to be used during the invocation of the method.
	// src: https://www.jsonrpc.org/specification#:~:text=4.2%20Parameter%20Structures
	Params json.RawMessage `json:"params"`
}

// AuthorizeAccountRPCParams is used to unmarshal the parameters of a request that has a From field
type AuthorizeAccountRPCParams struct {
	// From is an Ethereum account.
	From string `json:"from"`
}
