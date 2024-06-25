package pep

import "context"

// AccountUserPolicyDecisionPointPort is a port to adapt account usage authorization checks
type AccountUserPolicyDecisionPointPort interface {
	// AuthorizeAccountUser checks if a user is authorized to use an account, returns an error if it doesn't
	AuthorizeAccountUser(ctx context.Context, input AuthorizeAccountUserInput) (*AuthorizeAccountUserOutput, error)
}
