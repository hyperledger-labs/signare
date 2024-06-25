package pdp

// AuthorizeUserInput are the attributes to check a user permissions
type AuthorizeUserAccountInput struct {
	AccountID
}

// AuthorizeUserAccountOutput is the result of an account authorization.
type AuthorizeUserAccountOutput struct{}

// AuthorizeUserInput are the attributes to check a user permissions.
type AuthorizeUserInput struct {
	// UserID is the ID of the user.
	UserID string
	// UserID is the ID of the application.
	ApplicationID *string
	// UserID is the ID of the action.
	ActionID string
}

// AuthorizeUserOutput is the result of a user authorization.
type AuthorizeUserOutput struct{}
