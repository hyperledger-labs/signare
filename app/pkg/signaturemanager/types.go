package signaturemanager

import (
	"github.com/hyperledger-labs/signare/app/pkg/commons/logger"
	"github.com/hyperledger-labs/signare/app/pkg/entities"
	"github.com/hyperledger-labs/signare/app/pkg/entities/address"
)

// GenerateKeyInput for account generation requests.
type GenerateKeyInput struct {
	// Slot the slot to look for the keys
	Slot string
	// Pin the pin to authorize the user
	Pin string
	// Tracer to log what is needed
	Tracer logger.Tracer
}

// GenerateKeyOutput for account generation responses.
type GenerateKeyOutput struct {
	// Address derived from the generated public key.
	Address address.Address `json:"address"`
}

// RemoveKeyInput for account removal requests.
type RemoveKeyInput struct {
	// Slot the slot to look for the keys
	Slot string
	// Pin the pin to authorize the user
	Pin string
	// Tracer to log what is needed
	Tracer logger.Tracer
	// Address identifies the key pair to remove.
	Address address.Address `json:"address"`
}

// RemoveKeyOutput for account removal responses.
type RemoveKeyOutput struct{}

// ListKeysInput for account listing requests.
type ListKeysInput struct {
	// Slot the slot to look for the keys
	Slot string
	// Pin the pin to authorize the user
	Pin string
	// Tracer to log what is needed
	Tracer logger.Tracer
}

// ListKeysOutput for account listing responses.
type ListKeysOutput struct {
	Items []address.Address `json:"items"`
}

// SignInput for transaction signing requests.
type SignInput struct {
	// Slot the slot to look for the keys
	Slot string
	// Pin the pin to authorize the user
	Pin string
	// Tracer to log what is needed
	Tracer logger.Tracer
	// From address identifying the private key to use.
	From address.Address
	// Data to sign.
	Data entities.HexBytes
}

// SignOutput for transaction signing responses.
type SignOutput struct {
	// Signature signed bytes
	Signature []byte
}

// CloseInput input to close connection and clean up resources.
type CloseInput struct {
	// Tracer to log what is needed
	Tracer logger.Tracer
}

// CloseOutput output closing the connection and cleaning up all the resources.
type CloseOutput struct{}

// OpenInput input to open connection
type OpenInput struct {
	// Tracer to log what is needed
	Tracer logger.Tracer
}

// OpenOutput output opening the connection.
type OpenOutput struct{}

// IsAliveInput input to check the healthiness of a slot
type IsAliveInput struct {
	// Slot the slot to look for the keys
	Slot string
	// Pin the pin to authorize the user
	Pin string
	// Tracer to log what is needed
	Tracer logger.Tracer
}

// IsAliveOutput response of the healthiness of a slot
type IsAliveOutput struct {
	// IsAlive is true if the slot is ready to be used
	IsAlive bool
}
