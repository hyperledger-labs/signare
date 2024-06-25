// Package pip defines the Policy Information Point.
package pip

import (
	"context"

	"github.com/hyperledger-labs/signare/app/pkg/usecases/authorization/pdp"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/user"
)

var _ pdp.AccountsPolicyInformationPort = (*DefaultAccountsPIPAdapter)(nil)

// GetAccount returns data of a given Account
func (d DefaultAccountsPIPAdapter) GetAccount(ctx context.Context, input pdp.GetAccountInput) (*pdp.GetAccountOutput, error) {
	getAccountInput := user.GetAccountInput{
		AccountID: user.AccountID(input.AccountID),
	}

	_, getAccountErr := d.accountUseCase.GetAccount(ctx, getAccountInput)
	if getAccountErr != nil {
		return nil, getAccountErr
	}

	return &pdp.GetAccountOutput{}, nil
}

// DefaultAccountsPIPAdapterOptions are the set of fields to create an DefaultAccountsPIPAdapter
type DefaultAccountsPIPAdapterOptions struct {
	// AccountUseCase defines the management of the Account resource
	AccountUseCase user.AccountUseCase
}

// DefaultAccountsPIPAdapter is a port to adapt requests related to the Account resource
type DefaultAccountsPIPAdapter struct {
	accountUseCase user.AccountUseCase
}

// ProvideDefaultAccountsPIPAdapter provides an instance of an DefaultAccountsPIPAdapter
func ProvideDefaultAccountsPIPAdapter(options DefaultAccountsPIPAdapterOptions) (*DefaultAccountsPIPAdapter, error) {
	return &DefaultAccountsPIPAdapter{
		accountUseCase: options.AccountUseCase,
	}, nil
}
