// Package pepin defines the implementation of the Policy Decision Point input adapters from the Policy Enforcement Point.
package pepin

import (
	"context"
	"errors"

	"github.com/hyperledger-labs/signare/app/pkg/infra/middleware/authorization/pep"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/authorization/pdp"
)

var _ pep.AccountUserPolicyDecisionPointPort = new(DefaultAccountUserPolicyDecisionPointAdapter)

// AuthorizeAccountUser checks if a user is authorized to use an account, returns an error if it doesn't
func (adapter DefaultAccountUserPolicyDecisionPointAdapter) AuthorizeAccountUser(ctx context.Context, input pep.AuthorizeAccountUserInput) (*pep.AuthorizeAccountUserOutput, error) {
	authorizeUserAccountInput := pdp.AuthorizeUserAccountInput{
		AccountID: pdp.AccountID{
			UserID:        input.UserID,
			ApplicationID: input.ApplicationID,
			Address:       input.Address,
		},
	}

	_, authorizeUserAccountErr := adapter.policyDecisionPoint.AuthorizeUserAccount(ctx, authorizeUserAccountInput)
	if authorizeUserAccountErr != nil {
		return nil, authorizeUserAccountErr
	}

	return &pep.AuthorizeAccountUserOutput{}, nil
}

// DefaultAccountUserPolicyDecisionPointAdapterOptions are the set of fields to create an DefaultAccountUserPolicyDecisionPointAdapter
type DefaultAccountUserPolicyDecisionPointAdapterOptions struct {
	// DefaultPolicyDecisionPointUseCase is the business logic to perform user authorization for different actions
	DefaultPolicyDecisionPointUseCase *pdp.DefaultPolicyDecisionPointUseCase
}

// DefaultAccountUserPolicyDecisionPointAdapter is an adapter to adapt account usage authorization checks
type DefaultAccountUserPolicyDecisionPointAdapter struct {
	policyDecisionPoint pdp.DefaultPolicyDecisionPointUseCase
}

// ProvideDefaultAccountUserPolicyDecisionPointAdapter provides an instance of an DefaultAccountUserPolicyDecisionPointAdapter
func ProvideDefaultAccountUserPolicyDecisionPointAdapter(options DefaultAccountUserPolicyDecisionPointAdapterOptions) (*DefaultAccountUserPolicyDecisionPointAdapter, error) {
	if options.DefaultPolicyDecisionPointUseCase == nil {
		return nil, errors.New("mandatory 'DefaultPolicyDecisionPointUseCase' not provided")
	}

	return &DefaultAccountUserPolicyDecisionPointAdapter{
		policyDecisionPoint: *options.DefaultPolicyDecisionPointUseCase,
	}, nil
}
