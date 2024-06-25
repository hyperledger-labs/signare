package hsmmoduledb

import (
	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence"
)

// HardwareSecurityModuleDBFilter to filter lists of resources from the database
type HardwareSecurityModuleDBFilter struct {
	// HardwareSecurityModuleDB is the data struct of the resource in the database
	HardwareSecurityModuleDB
	// Order is the order of the list based on an attribute
	Order *persistence.Order `valid:"optional"`
	// FilterGroup is a collection of filters
	FilterGroup *persistence.FilterGroup `valid:"optional"`
	// Pagination is the page info of the list
	Pagination *persistence.Pagination `valid:"optional"`
}

// Paged creates a pagination filter.
func (filter *HardwareSecurityModuleDBFilter) Paged(limit, offset int) *HardwareSecurityModuleDBFilter {
	filter.Pagination = &persistence.Pagination{
		Limit:  limit,
		Offset: offset,
	}
	return filter
}

// Sort creates a sorting filter.
func (filter *HardwareSecurityModuleDBFilter) Sort(orderBy string, orderDirection persistence.OrderDirection) *HardwareSecurityModuleDBFilter {
	filter.Order = &persistence.Order{
		By:        persistence.OrderByOption(orderBy),
		Direction: orderDirection,
	}
	return filter
}

// AppendFilter Append filter.
func (filter *HardwareSecurityModuleDBFilter) AppendFilter(theFilter persistence.Filter) {
	if filter.FilterGroup == nil {
		filter.FilterGroup = &persistence.FilterGroup{
			Filters: make([]persistence.Filter, 0),
		}
	}
	filter.FilterGroup.Filters = append(filter.FilterGroup.Filters, theFilter)
}
