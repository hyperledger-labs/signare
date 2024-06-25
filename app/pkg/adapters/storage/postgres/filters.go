package postgres

import "github.com/hyperledger-labs/signare/app/pkg/commons/persistence"

// NewEqualFilter creates a EqualFilter instance with the given options
func NewEqualFilter(by string) persistence.EqualFilter {
	return persistence.EqualFilter{
		By: by,
	}
}

// NewBetweenFilter creates a BetweenFilter instance with the given options
func NewBetweenFilter(by string, minValue string, maxValue string) persistence.BetweenFilter {
	return persistence.BetweenFilter{
		By:       by,
		MinValue: minValue,
		MaxValue: maxValue,
	}
}

// NewLessFilter creates a LessFilter instance with the given options
func NewLessFilter(by string) persistence.LessFilter {
	return persistence.LessFilter{
		By: by,
	}
}

// NewLessOrEqualFilter creates a LessOrEqualFilter instance with the given options
func NewLessOrEqualFilter(by string) persistence.LessOrEqualFilter {
	return persistence.LessOrEqualFilter{
		By: by,
	}
}

// NewGreaterFilter creates a GreaterFilter instance with the given options
func NewGreaterFilter(by string) persistence.GreaterFilter {
	return persistence.GreaterFilter{
		By: by,
	}
}

// NewGreaterOrEqualFilter creates a GreaterOrEqualFilter instance with the given options
func NewGreaterOrEqualFilter(by string) persistence.GreaterOrEqualFilter {
	return persistence.GreaterOrEqualFilter{
		By: by,
	}
}

// NewListEqualFilter creates a ListEqualFilter instance with the given options
func NewListEqualFilter(by string, values []string) persistence.ListEqualFilter {
	return persistence.ListEqualFilter{
		By:     by,
		Values: values,
	}
}
