// Code generated by Signare OpenAPI generator. DO NOT EDIT

package httpinfra

import (
	"github.com/hyperledger-labs/signare/app/pkg/infra/httpinfra"
)

// CollectionPage - Page data of the query
type CollectionPage struct {
	// The size of the collection's page
	Limit *int32 `json:"limit"`
	// The entry of the table on which the collection starts
	Offset *int32 `json:"offset"`
	// True if there are more pages to collect from the database
	MoreItems *bool `json:"moreItems"`
}

// ValidateWith check whether CollectionPage is valid
func (data CollectionPage) ValidateWith() (*httpinfra.ValidationResult, *httpinfra.HTTPError) {
	if data.Limit == nil {
		httpError := httpinfra.NewHTTPError(httpinfra.StatusInvalidArgument)
		httpError.SetMessage("error validating field [limit]")
		return nil, httpError
	}
	if data.Offset == nil {
		httpError := httpinfra.NewHTTPError(httpinfra.StatusInvalidArgument)
		httpError.SetMessage("error validating field [offset]")
		return nil, httpError
	}
	if data.MoreItems == nil {
		httpError := httpinfra.NewHTTPError(httpinfra.StatusInvalidArgument)
		httpError.SetMessage("error validating field [moreItems]")
		return nil, httpError
	}
	return &httpinfra.ValidationResult{
		Valid: true,
	}, nil
}

// SetDefaults sets default values as defined in the API spec
func (data *CollectionPage) SetDefaults() {
}
