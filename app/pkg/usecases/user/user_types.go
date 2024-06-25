package user

import (
	"github.com/hyperledger-labs/signare/app/pkg/entities"
	"github.com/hyperledger-labs/signare/app/pkg/entities/address"
)

const (
	defaultOrderDirection = entities.OrderDesc
)

// User defines the administrator user in the application.
type User struct {
	entities.ApplicationStandardResourceMeta
	// InternalResourceID uniquely identifies a User by a single ID.
	entities.InternalResourceID
	// Roles of this User.
	Roles []string
	// Description of this User.
	Description *string
	// Accounts assigned to the User.
	Accounts []Account
}

// UserCollection defines a collection of User resources.
type UserCollection struct {
	// Items User in collection.
	Items []User
	// StandardCollectionPage is the page data of the collection.
	entities.StandardCollectionPage
}

// CreateUserInput configures the creation of a User.
type CreateUserInput struct {
	// ID defines the identifier of the User resource.
	ID *string `valid:"optional"`
	// ApplicationID defines the identifier of the Application of the User resource.
	ApplicationID string `valid:"required"`
	// Description of this User.
	Description *string `valid:"optional"`
	// Roles of this User.
	Roles []string `valid:"required"`
}

// CreateUserOutput defines the output of the creation of a User.
type CreateUserOutput struct {
	// User defines the administrator user in the application.
	User
}

// EditUserInput configures the update of a User.
type EditUserInput struct {
	// ApplicationStandardID defines the identifier of the resource.
	entities.ApplicationStandardID
	// ResourceVersion resource version for resource locking.
	ResourceVersion string `valid:"required"`
	// Description of this User.
	Description *string `valid:"optional"`
	// Roles of this User.
	Roles []string `valid:"optional"`
}

// EditUserOutput defines the output for editing a User.
type EditUserOutput struct {
	// User defines the administrator user in the application.
	User
}

// ListUsersInput defines all possible options to list User resources.
type ListUsersInput struct {
	// ApplicationID defines the identifier of the Application of the User resource.
	ApplicationID string `valid:"required"`
	// PageLimit maximum amount of Application in list output.
	PageLimit int `valid:"natural"`
	// PageOffset amount of Application elapsed in list output.
	PageOffset int `valid:"natural"`
	// OrderBy whether to order by last update date.
	OrderBy string
	// OrderDirection the direction of the OrderBy.
	OrderDirection string
}

// ListUsersOutput defines the output of listing Users.
type ListUsersOutput struct {
	// UserCollection defines a collection of User resources.
	UserCollection
}

// GetUserInput defines the input for getting a User.
type GetUserInput struct {
	// ApplicationStandardID defines the identifier of the resource.
	entities.ApplicationStandardID
}

// GetUserOutput defines the output of getting a User.
type GetUserOutput struct {
	// User defines the administrator user in the application.
	User
}

// DeleteUserInput configures the deletion of a User.
type DeleteUserInput struct {
	// ApplicationStandardID defines the identifier of the resource.
	entities.ApplicationStandardID
}

// DeleteUserOutput defines the output of deleting a User.
type DeleteUserOutput struct {
	// User defines the administrator user in the application.
	User
}

// EnableAccountsInput configures the creation of accounts for a User.
type EnableAccountsInput struct {
	// UserID defines the identifier of the User resource.
	UserID string `valid:"required"`
	// ApplicationID defines the identifier of the Application of the User resource.
	ApplicationID string `valid:"required"`
	// Addresses to be assigned to the user.
	Addresses []address.Address `valid:"address"`
}

// EnableAccountsOutput defines the output of creating accounts for a User.
type EnableAccountsOutput struct {
	// User defines the administrator user in the application.
	User
}

// DisableAccountInput configures the removal of an account from a User.
type DisableAccountInput struct {
	// UserID defines the identifier of the User resource.
	UserID string `valid:"required"`
	// ApplicationID defines the identifier of the Application of the User resource.
	ApplicationID string `valid:"required"`
	// Address to be removed from the user's accounts.
	Address address.Address `valid:"address"`
}

// DisableAccountOutput defines the output of removing an account from a User.
type DisableAccountOutput struct {
	// User defines the administrator user in the application.
	User
}
