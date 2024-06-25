package admin

import (
	"github.com/hyperledger-labs/signare/app/pkg/entities"
)

// Admin defines the administrator user in the application.
type Admin struct {
	entities.StandardResourceMeta
	// InternalResourceID uniquely identifies an Admin by a single ID.
	entities.InternalResourceID
	// Description of this Admin.
	Description *string
	// Roles of this Admin for the RBAC.
	Roles []string
}

// AdminCollection defines a collection of Admin resources.
type AdminCollection struct {
	// Items Admin in collection
	Items []Admin
	// StandardCollectionPage is the page data of the collection
	entities.StandardCollectionPage
}

// AdminEditable defines the editable fields of the Admin.
type AdminEditable struct {
	entities.StandardResourceMeta
	// Description of this Admin.
	Description *string
}

// CreateAdminInput configures the creation of an Admin
type CreateAdminInput struct {
	// StandardID defines the identifier of the Admin resource.
	entities.StandardID
	// Description of this Admin.
	Description *string
}

// CreateAdminOutput output from Admin users creation.
type CreateAdminOutput struct {
	Admin
}

// GetAdminInput input to get a specific Admin.
type GetAdminInput struct {
	// StandardID defines the identifier of the Admin resource.
	entities.StandardID
}

// GetAdminOutput the requested Admin resource.
type GetAdminOutput struct {
	Admin
}

// EditAdminInput configures the update of an Admin.
type EditAdminInput struct {
	AdminEditable
}

// EditAdminOutput the edited Admin resource.
type EditAdminOutput struct {
	Admin
}

// DeleteAdminInput identifies the Admin resource to be removed.
type DeleteAdminInput struct {
	// StandardID defines the identifier of the Admin resource.
	entities.StandardID
}

// DeleteAdminOutput the deleted Admin resource.
type DeleteAdminOutput struct {
	Admin
}

// ListAdminsInput defines all possible options to list Admin resources.
type ListAdminsInput struct {
	// PageLimit maximum amount of Admin in list output.
	PageLimit int `valid:"natural"`
	// PageOffset amount of Admin elapsed in list output.
	PageOffset int `valid:"natural"`
	// OrderBy whether to order by last update date.
	OrderBy string
	// OrderDirection the direction in which the list will be ordered base on the attribute selected in OrderBy.
	OrderDirection string
}

// ListAdminsOutput collection of Admin resources.
type ListAdminsOutput struct {
	AdminCollection
}
