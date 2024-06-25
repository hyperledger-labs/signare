package userdb

import (
	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence"
)

// UserDBFilter to filter lists of resources from the database
type UserDBFilter struct {
	// UserDB is the data struct of the resource in the database
	UserDB
	// Order is the order of the list based on an attribute
	Order *persistence.Order `valid:"optional"`
	// FilterGroup is a collection of filters
	FilterGroup *persistence.FilterGroup `valid:"optional"`
	// Pagination is the page info of the list
	Pagination *persistence.Pagination `valid:"optional"`
}

// AppendFilter Append filter.
func (filter *UserDBFilter) AppendFilter(theFilter persistence.Filter) {
	if filter.FilterGroup == nil {
		filter.FilterGroup = &persistence.FilterGroup{
			Filters: make([]persistence.Filter, 0),
		}
	}
	filter.FilterGroup.Filters = append(filter.FilterGroup.Filters, theFilter)
}

// Paged creates a pagination filter.
func (filter *UserDBFilter) Paged(limit, offset int) *UserDBFilter {
	filter.Pagination = &persistence.Pagination{
		Limit:  limit,
		Offset: offset,
	}
	return filter
}

// Sort creates a sorting filter.
func (filter *UserDBFilter) Sort(orderBy string, orderDirection persistence.OrderDirection) *UserDBFilter {
	filter.Order = &persistence.Order{
		By:        persistence.OrderByOption(orderBy),
		Direction: orderDirection,
	}
	return filter
}
