package pdp

import (
	"context"

	"github.com/hyperledger-labs/signare/app/pkg/entities/address"
)

// AccountsPolicyInformationPort is a port to adapt requests related to the Account resource
type AccountsPolicyInformationPort interface {
	// GetAccount returns data of a given Account
	GetAccount(ctx context.Context, input GetAccountInput) (*GetAccountOutput, error)
}

// GetAccountInput are the attributes needed to get an Account
type GetAccountInput struct {
	AccountID
}

// GetAccountOutput is the result of getting an account
type GetAccountOutput struct {
}

// AccountID defines the identifier of the Account resource.
type AccountID struct {
	// Address defines the address of the Account resource.
	Address address.Address `valid:"address"`
	// UserID defines the identifier of the Account resource.
	UserID string `valid:"required"`
	// ApplicationID defines the identifier of the Application of the Account resource.
	ApplicationID string `valid:"required"`
}
