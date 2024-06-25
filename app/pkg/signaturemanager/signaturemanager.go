// Package signaturemanager defines the management of signature managers
package signaturemanager

import (
	"context"
)

// DigitalSignatureManager defines an interface for interacting with a signature manager. It utilizes the concept of addresses to identify key pairs, where an address consists of the 20 last characters of the public key.
type DigitalSignatureManager interface {
	// GenerateKey generates a new public and private key returning the corresponding address.
	GenerateKey(ctx context.Context, input GenerateKeyInput) (*GenerateKeyOutput, error)
	// RemoveKey removes the public and private key identified by the provided address. It returns an error if it fails or if the key pair doesn't exist.
	RemoveKey(ctx context.Context, input RemoveKeyInput) (*RemoveKeyOutput, error)
	// ListKeys retrieves all stored keys as a list of addresses.
	ListKeys(ctx context.Context, input ListKeysInput) (*ListKeysOutput, error)
	// Sign signs a set of bytes with the private key identified by the provided address.
	Sign(ctx context.Context, input SignInput) (*SignOutput, error)
	// Close closes the connection and cleans up open resources.
	Close(ctx context.Context, input CloseInput) (*CloseOutput, error)
	// Open opens the connection to a digital signature manager provider.
	Open(ctx context.Context, input OpenInput) (*OpenOutput, error)
	// IsAlive checks if a given slot healthiness in a digital signature manager, returns true if it's healthy
	IsAlive(ctx context.Context, input IsAliveInput) (*IsAliveOutput, error)
}
