package entities

import (
	"github.com/hyperledger-labs/signare/app/pkg/commons/time"

	"github.com/google/uuid"
)

const (
	OrderAsc  = "asc"
	OrderDesc = "desc"

	OrderByLastUpdate   = "lastUpdate"
	OrderByCreationDate = "creationDate"
)

// ContextKey used in the signare for storing and retrieving values from context.
type ContextKey string

// InternalResourceID is a hash composed by the unique keys of a resource, and it uniquely identifies a resource by a single ID.
type InternalResourceID string

func NewInternalResourceID() InternalResourceID {
	id := uuid.NewString()

	return InternalResourceID(id)
}

// String implements fmt.Stringer
func (i InternalResourceID) String() string {
	return string(i)
}

// StandardID the Standard identification for every not application-related resource.
type StandardID struct {
	ID string `json:"id" valid:"required,maxstringlength(64)" storage:"id"`
}

// StandardResource the Standard identification for every not application-related resource including when it was created and last updated
type StandardResource struct {
	// StandardID identifier of StandardResource object
	StandardID `json:",inline"`
	// Timestamps of StandardResource object
	Timestamps `json:",inline"`
}

// ApplicationStandardID the Standard identification for every application-related resource.
type ApplicationStandardID struct {
	ID            string `json:"id" valid:"required,maxstringlength(64)" storage:"id"`
	ApplicationID string `json:"application_id"  valid:"required" storage:"application_id"`
}

// ApplicationStandardResource the Standard identification for every application-related resource including when it was created and last updated
type ApplicationStandardResource struct {
	// ApplicationStandardID identifier of ApplicationStandardResource object
	ApplicationStandardID `json:",inline"`
	// Timestamps of ApplicationStandardResource object
	Timestamps `json:",inline"`
}

type Timestamps struct {
	// CreationDate creation date of the resource
	CreationDate time.Timestamp `json:"creationDate"`
	// LastUpdate last update date of the resource
	LastUpdate time.Timestamp `json:"lastUpdate"`
}

// StandardCollectionPage defines the standard data structure included in all the collections items.
type StandardCollectionPage struct {
	// Limit of entries contained in the page
	Limit int `json:"limit"`
	// Offset zero-based offset of the first item in the page
	Offset int `json:"offset"`
	// MoreItems is true if there are entries after the < Offset + Size > entry
	MoreItems bool `json:"moreItems"`
}

func NewUnlimitedQueryStandardCollectionPage(size int) StandardCollectionPage {
	return StandardCollectionPage{
		Limit:     size,
		Offset:    0,
		MoreItems: false,
	}
}

// StandardResourceMeta defines the meta-information of a standard non-application resource.
type StandardResourceMeta struct {
	// entities.StandardResource standard resource data
	StandardResource
	// ResourceVersion resource version for resource locking
	ResourceVersion string `json:"ResourceVersion" storage:"resource_version"`
}

// ApplicationStandardResourceMeta defines the meta-information of an application resource.
type ApplicationStandardResourceMeta struct {
	// entities.ApplicationStandardResource standard application resource data
	ApplicationStandardResource
	// ResourceVersion resource version for resource locking
	ResourceVersion string `json:"ResourceVersion" storage:"resource_version"`
}
