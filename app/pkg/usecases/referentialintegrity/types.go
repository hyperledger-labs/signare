package referentialintegrity

import (
	"github.com/hyperledger-labs/signare/app/pkg/entities"
)

const (
	defaultOrderDirection = entities.OrderDesc
)

// ReferentialIntegrityEntry represents an entry for managing referential integrity between resources.
type ReferentialIntegrityEntry struct {
	entities.StandardResource
	// ResourceID is the unique identifier of the resource that references to another resource.
	ResourceID string
	// ResourceKind is the kind of the resource that references to another resource.
	ResourceKind ResourceKind
	// ParentResourceID is the unique identifier of the resource being referenced.
	ParentResourceID string
	// ParentResourceKind is the kind of the resource being referenced.
	ParentResourceKind ResourceKind
}

// ReferentialIntegrityEntryCollection defines a collection of ReferentialIntegrityEntry resources.
type ReferentialIntegrityEntryCollection struct {
	// Items is a collection of ReferentialIntegrityEntry.
	Items []ReferentialIntegrityEntry
	// StandardCollectionPage is the page data of the collection.
	entities.StandardCollectionPage
}

// CreateEntryInput configures the creation of a ReferentialIntegrityEntry.
type CreateEntryInput struct {
	// ResourceID is the unique identifier of the resource that references another resource.
	ResourceID string `valid:"required"`
	// ResourceKind is the kind of the resource that references another resource.
	ResourceKind ResourceKind `valid:"required"`
	// ParentResourceID is the unique identifier of the resource being referenced.
	ParentResourceID string `valid:"required"`
	// ParentResourceKind is the kind of the resource being referenced.
	ParentResourceKind ResourceKind `valid:"required"`
}

// CreateEntryOutput defines the output of listing ReferentialIntegrityEntry resources.
type CreateEntryOutput struct {
	ReferentialIntegrityEntry
}

// GetEntryInput defines the input for getting a ReferentialIntegrityEntry.
type GetEntryInput struct {
	entities.StandardID
}

// GetEntryOutput defines the output of getting a ReferentialIntegrityEntry.
type GetEntryOutput struct {
	ReferentialIntegrityEntry
}

// ListEntriesInput defines the options to list ReferentialIntegrityEntry resources.
type ListEntriesInput struct {
	// Resource to filter referential integrity entries for.
	Resource *Resource `valid:"optional"`
	// Parent resource to filter referential integrity entries for.
	Parent *Resource `valid:"optional"`
	// PageLimit maximum amount of Application in list output.
	PageLimit int `valid:"natural"`
	// PageOffset amount of Application elapsed in list output.
	PageOffset int `valid:"natural"`
	// OrderBy whether to order by last update date.
	OrderBy string
	// OrderDirection the direction of the OrderBy.
	OrderDirection string
}

type Resource struct {
	// ResourceID to filter referential integrity entries for.
	ResourceID string `valid:"required"`
	// ResourceKind to filter referential integrity entries for.
	ResourceKind ResourceKind `valid:"required"`
}

// ListEntriesOutput defines the output of listing all ReferentialIntegrityEntry resources.
type ListEntriesOutput struct {
	ReferentialIntegrityEntryCollection
}

// DeleteEntryInput configures the deletion of a ReferentialIntegrityEntry.
type DeleteEntryInput struct {
	entities.StandardID
}

// DeleteEntryOutput defines the output of deleting a ReferentialIntegrityEntry.
type DeleteEntryOutput struct {
	ReferentialIntegrityEntry
}

// ListMyChildrenEntriesInput defines the options to list ReferentialIntegrityEntry resources
// that are referencing the specified resource.
type ListMyChildrenEntriesInput struct {
	// ParentResourceID is the unique identifier of the resource being referenced.
	ParentResourceID string `valid:"required"`
	// ParentResourceKind is the kind of the resource being referenced.
	ParentResourceKind ResourceKind `valid:"required"`
}

// ListMyChildrenEntriesOutput defines the output of listing all ReferentialIntegrityEntry resources
// that are referencing the specified resource.
type ListMyChildrenEntriesOutput struct {
	ReferentialIntegrityEntryCollection
}

// DeleteMyEntriesIfAnyInput configures the deletion of all ReferentialIntegrityEntry that matches
// the specified ResourceID and ResourceKind.
type DeleteMyEntriesIfAnyInput struct {
	// ResourceID is the unique identifier of the resource.
	ResourceID string `valid:"required"`
	// ResourceKind is the kind of the resource.
	ResourceKind ResourceKind `valid:"required"`
}

// GetEntryByResourceAndParentInput defines the input for getting a ReferentialIntegrityEntry
// given the specified parameters.
type GetEntryByResourceAndParentInput struct {
	// ResourceID is the unique identifier of the resource that references another resource.
	ResourceID string `valid:"required"`
	// ResourceKind is the kind of the resource that references another resource.
	ResourceKind ResourceKind `valid:"required"`
	// ParentResourceID is the unique identifier of the resource being referenced.
	ParentResourceID string `valid:"required"`
	// ParentResourceKind is the kind of the resource being referenced.
	ParentResourceKind ResourceKind `valid:"required"`
}

// GetEntryByResourceAndParentOutput defines the output of getting a ReferentialIntegrityEntry
// by the specified parameters.
type GetEntryByResourceAndParentOutput struct {
	ReferentialIntegrityEntry
}
