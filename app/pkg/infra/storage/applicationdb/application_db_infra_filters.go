package applicationdb

import (
	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence"
)

// ApplicationDBFilter to filter lists of resources from the database
type ApplicationDBFilter struct {
	// ApplicationDB is the data struct of the resource in the database
	ApplicationDB
	// Order is the order of the list based on an attribute
	Order *persistence.Order `valid:"optional"`
	// FilterGroup is a collection of filters
	FilterGroup *persistence.FilterGroup `valid:"optional"`
	// Pagination is the page info of the list
	Pagination *persistence.Pagination `valid:"optional"`
}

// Paged creates a pagination filter.
func (filter *ApplicationDBFilter) Paged(limit, offset int) *ApplicationDBFilter {
	filter.Pagination = &persistence.Pagination{
		Limit:  limit,
		Offset: offset,
	}
	return filter
}

// Sort creates a sorting filter.
func (filter *ApplicationDBFilter) Sort(orderBy string, orderDirection persistence.OrderDirection) *ApplicationDBFilter {
	filter.Order = &persistence.Order{
		By:        persistence.OrderByOption(orderBy),
		Direction: orderDirection,
	}
	return filter
}
