package persistence

import (
	"context"
	"database/sql"
	"fmt"
)

type Storage interface {
	// QueryAll obtains all element from statement with arguments provided. It returns an error if it fails
	QueryAll(ctx context.Context, stmtID string, args interface{}, dst interface{}) error
	// ExecuteStmt executes statement with arguments provided. It returns an error if it fails
	ExecuteStmt(ctx context.Context, stmtID string, args interface{}) error
	// ExecuteStmtWithStorageResult executes statement with arguments provided. It returns information on stmt or an error if it fails
	ExecuteStmtWithStorageResult(ctx context.Context, stmtID string, args interface{}) (*ExecuteStmtWithStorageResultOutput, error)
	// BeginTransaction starts a Storage operation that is transactional. It returns an error if it fails
	BeginTransaction(ctx context.Context) (context.Context, error)
	// RollbackTransaction rollbacks to previous state a Storage operation that is transactional. It returns an error if it fails
	RollbackTransaction(ctx context.Context) error
	// CommitTransaction confirms a Storage operation that is transactional. It returns an error if it fails
	CommitTransaction(ctx context.Context) error
	// GetTransaction returns the transaction from context (in case of success).
	GetTransaction(ctx context.Context) (any, error)
	// AddConfig registers SQL statements
	AddConfig(config StorageConfig) error
}

// OrderByOption defines the field to apply order
type OrderByOption string

// OrderDirection defines the order direction
type OrderDirection string

const (
	// Asc ascendant direction order
	Asc OrderDirection = "asc"
	// Desc descendant direction order
	Desc OrderDirection = "desc"
	// CreationDate creation date order
	CreationDate OrderByOption = "creation_date"
	// LastUpdate update date order
	LastUpdate OrderByOption = "last_update"
)

// Order defines order on Storage data
type Order struct {
	By        OrderByOption  `valid:"required"`
	Direction OrderDirection `valid:"required"`
}

// FilterGroup defines the group of filters Storage data
type FilterGroup struct {
	Filters []Filter
}

// Filter defines the functionality to convert filters into statements
type Filter interface {
	// ToSQLStmt SQL statement for the filter
	ToSQLStmt() string
}

// EqualFilter to filter by the specified field with equal value
type EqualFilter struct {
	By string
}

// ToSQLStmt SQL statement for the filter
func (f EqualFilter) ToSQLStmt() string {
	return fmt.Sprintf("%s=:%s", f.By, f.By)
}

// BetweenFilter to filter by the specified field between two values
type BetweenFilter struct {
	By       string
	MinValue string
	MaxValue string
}

// ToSQLStmt SQL statement for the filter
func (f BetweenFilter) ToSQLStmt() string {
	return fmt.Sprintf("%s>=%s AND %s<=%s", f.By, f.MinValue, f.By, f.MaxValue)
}

// LessFilter to filter by the specified field with less value
type LessFilter struct {
	By string
}

// ToSQLStmt SQL statement for the filter
func (f LessFilter) ToSQLStmt() string {
	return fmt.Sprintf("%s<:%s", f.By, f.By)
}

// LessOrEqualFilter to filter by the specified field with less or equal value
type LessOrEqualFilter struct {
	By string
}

// ToSQLStmt SQL statement for the filter
func (f LessOrEqualFilter) ToSQLStmt() string {
	return fmt.Sprintf("%s<=:%s", f.By, f.By)
}

// GreaterFilter to filter by the specified field with greater value
type GreaterFilter struct {
	By string
}

// ToSQLStmt SQL statement for the filter
func (f GreaterFilter) ToSQLStmt() string {
	return fmt.Sprintf("%s>:%s", f.By, f.By)
}

// GreaterOrEqualFilter to filter by the specified field with greater or equal value
type GreaterOrEqualFilter struct {
	By string
}

// ToSQLStmt SQL statement for the filter
func (f GreaterOrEqualFilter) ToSQLStmt() string {
	return fmt.Sprintf("%s>=:%s", f.By, f.By)
}

// ListEqualFilter to filter by the specified field with a list of possible value
type ListEqualFilter struct {
	By     string
	Values []string
}

// ToSQLStmt SQL statement for the filter
func (f ListEqualFilter) ToSQLStmt() string {
	filterDefinitionFirstIndex := "%s='%s'"
	filterDefinitionOtherIndexes := " OR %s='%s'"
	var filterBuild string
	if f.Values == nil || len(f.Values) == 0 {
		// It is assumed that if no filter is provided, then we should return nothing
		return "FALSE"
	}
	if len(f.By) > 1 {
		for index := range f.Values {
			if index == 0 {
				filterBuild = fmt.Sprintf(filterDefinitionFirstIndex, f.By, f.Values[index])
			} else {
				filterBuild += fmt.Sprintf(filterDefinitionOtherIndexes, f.By, f.Values[index])
			}
		}
	}
	return "(" + filterBuild + ")"
}

// FilterOption defines a filter option
type FilterOption string

// Pagination from the Storage data with limit/offset navigation
type Pagination struct {
	// Limit maximum number of entries to return
	Limit int `valid:"required"`
	// Offset zero-based offset of the first item in the collection to return
	Offset int `valid:"required"`
}

// ExecuteStmtWithStorageResultOutput wraps the result of an statement
type ExecuteStmtWithStorageResultOutput struct {
	Result sql.Result
}
