package user

import (
	"github.com/hyperledger-labs/signare/app/pkg/entities"
	"github.com/hyperledger-labs/signare/app/pkg/entities/address"
)

// AccountID defines the identifier of the Account resource.
type AccountID struct {
	// Address defines the address of the Account resource.
	Address address.Address `valid:"address"`
	// UserID defines the identifier of the Account resource.
	UserID string `valid:"required"`
	// ApplicationID defines the identifier of the Application of the Account resource.
	ApplicationID string `valid:"required"`
}

// Account defines the Account resource.
type Account struct {
	// AccountID defines the identifier of the Account resource.
	AccountID
	// InternalResourceID uniquely identifies an Account by a single ID.
	entities.InternalResourceID
	// TimeStamp of the Account resource.
	entities.Timestamps
}

// AccountCollection defines a collection of Account resources.
type AccountCollection struct {
	// Items is a collection of Accounts.
	Items []Account
	// StandardCollectionPage is the page data of the collection.
	entities.StandardCollectionPage
}

// CreateAccountInput configures the creation of an Account.
type CreateAccountInput struct {
	// AccountID defines the identifier of the Account resource.
	AccountID
}

// CreateAccountOutput defines the output of the creation of an Account.
type CreateAccountOutput struct {
	// Account defines the Account resource.
	Account
}

// ListAccountsInput defines all possible options to list Account resources.
type ListAccountsInput struct {
	// ApplicationID to filter Accounts for.
	ApplicationID string `valid:"required"`
	// UserID to filter Accounts for.
	UserID *string `valid:"optional"`
	// Address to filter Accounts for.
	Address *address.Address `valid:"optional"`
}

// ListAccountsOutput defines the output of listing Accounts.
type ListAccountsOutput struct {
	// AccountCollection defines a collection of Account resources.
	AccountCollection
}

// GetAccountInput defines the input for getting an Account.
type GetAccountInput struct {
	// AccountID defines the identifier of the Account resource.
	AccountID
}

// GetAccountOutput defines the output of getting an Account.
type GetAccountOutput struct {
	// Account defines the Account resource.
	Account
}

// DeleteAccountInput configures the deletion of an Account.
type DeleteAccountInput struct {
	// AccountID defines the identifier of the Account resource.
	AccountID
}

// DeleteAccountOutput defines the output of deleting an Account.
type DeleteAccountOutput struct {
	// Account defines the Account resource.
	Account
}

// DeleteAllAccountsForAddressInput configures the deletion of all Accounts with a given address for an Application.
type DeleteAllAccountsForAddressInput struct {
	// Address defines the address of the Account resource.
	Address address.Address `valid:"address"`
	// ApplicationID defines the identifier of the Application of the Account resource.
	ApplicationID string `valid:"required"`
}

// DeleteAllAccountsForAddressOutput defines the output of deleting all Account with a given address for an Application.
type DeleteAllAccountsForAddressOutput struct {
	// Items is a collection of Accounts.
	Items []Account
}
