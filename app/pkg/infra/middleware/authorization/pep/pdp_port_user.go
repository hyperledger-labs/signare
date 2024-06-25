package pep

import "context"

// UserPolicyDecisionPointPort is a port to adapt authorization checks
type UserPolicyDecisionPointPort interface {
	// AuthorizeUser checks if a user is authorized to perform an action, returns an error if it doesn't
	AuthorizeUser(ctx context.Context, input AuthorizeUserInput) (*AuthorizeUserOutput, error)
}
