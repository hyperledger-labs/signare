// Package pepin defines the implementation of the Policy Decision Point input adapters from the Policy Enforcement Point.
package pepin

import (
	"context"
	"errors"

	"github.com/hyperledger-labs/signare/app/pkg/infra/middleware/authorization/pep"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/authorization/pdp"
)

var _ pep.UserPolicyDecisionPointPort = new(DefaultUserPolicyDecisionPointAdapter)

// AuthorizeUser checks if a user is authorized to perform an action, returns an error if it doesn't
func (adapter DefaultUserPolicyDecisionPointAdapter) AuthorizeUser(ctx context.Context, input pep.AuthorizeUserInput) (*pep.AuthorizeUserOutput, error) {
	authorizeUserAccountInput := pdp.AuthorizeUserInput{
		UserID:        input.UserID,
		ApplicationID: input.ApplicationID,
		ActionID:      input.ActionID,
	}

	_, authorizeUserAccountErr := adapter.policyDecisionPoint.AuthorizeUser(ctx, authorizeUserAccountInput)
	if authorizeUserAccountErr != nil {
		return nil, authorizeUserAccountErr
	}

	return &pep.AuthorizeUserOutput{}, nil
}

// DefaultUserPolicyDecisionPointAdapterOptions are the set of fields to create an DefaultUserPolicyDecisionPointAdapter
type DefaultUserPolicyDecisionPointAdapterOptions struct {
	// DefaultPolicyDecisionPointUseCase is the business logic to perform user authorization for different actions
	DefaultPolicyDecisionPoint *pdp.DefaultPolicyDecisionPointUseCase
}

// DefaultUserPolicyDecisionPointAdapter is an adapter to adapt authorization checks
type DefaultUserPolicyDecisionPointAdapter struct {
	policyDecisionPoint pdp.DefaultPolicyDecisionPointUseCase
}

// ProvideUserPolicyDecisionPointAdapter provides an instance of an DefaultUserPolicyDecisionPointAdapter
func ProvideUserPolicyDecisionPointAdapter(options DefaultUserPolicyDecisionPointAdapterOptions) (*DefaultUserPolicyDecisionPointAdapter, error) {
	if options.DefaultPolicyDecisionPoint == nil {
		return nil, errors.New("mandatory 'DefaultPolicyDecisionPointUseCase' not provided")
	}

	return &DefaultUserPolicyDecisionPointAdapter{
		policyDecisionPoint: *options.DefaultPolicyDecisionPoint,
	}, nil
}
