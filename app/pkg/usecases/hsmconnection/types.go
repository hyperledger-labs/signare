package hsmconnection

import (
	"github.com/hyperledger-labs/signare/app/pkg/entities"
)

// ByApplicationInput input to get an HSMConnector given an application.
type ByApplicationInput struct {
	// ApplicationID identifier of the application.
	ApplicationID string
}

// HSMConnection HSM connection details.
type HSMConnection struct {
	// Slot the HSM slot.
	Slot string
	// Pin the pin of the slot.
	Pin string
	// ModuleKind type of the HSM module.
	ModuleKind string
	// ChainID application's chain ID.
	ChainID entities.Int256
}
