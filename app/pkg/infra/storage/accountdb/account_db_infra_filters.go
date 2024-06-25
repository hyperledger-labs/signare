package accountdb

import (
	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence"
)

// AccountDBFilter to filter lists of resources from the database
type AccountDBFilter struct {
	// AccountDB is the data struct of the resource in the database
	AccountDB
	// FilterGroup is a collection of filters
	FilterGroup *persistence.FilterGroup `valid:"optional"`
}

// AppendFilter Append filter.
func (filter *AccountDBFilter) AppendFilter(theFilter persistence.Filter) {
	if filter.FilterGroup == nil {
		filter.FilterGroup = &persistence.FilterGroup{
			Filters: make([]persistence.Filter, 0),
		}
	}
	filter.FilterGroup.Filters = append(filter.FilterGroup.Filters, theFilter)
}
