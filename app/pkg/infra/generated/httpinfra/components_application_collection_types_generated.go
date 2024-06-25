// Code generated by Signare OpenAPI generator. DO NOT EDIT

package httpinfra

import (
	"github.com/hyperledger-labs/signare/app/pkg/infra/httpinfra"
)

type ApplicationCollection struct {
	// The size of the collection's page
	Limit *int32 `json:"limit"`
	// The entry of the table on which the collection starts
	Offset *int32 `json:"offset"`
	// True if there are more pages to collect from the database
	MoreItems *bool `json:"moreItems"`
	// collection of applications.
	Items *[]ApplicationDetail `json:"items"`
}

// ValidateWith check whether ApplicationCollection is valid
func (data ApplicationCollection) ValidateWith() (*httpinfra.ValidationResult, *httpinfra.HTTPError) {
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
	if data.Items == nil {
		httpError := httpinfra.NewHTTPError(httpinfra.StatusInvalidArgument)
		httpError.SetMessage("error validating field [items]")
		return nil, httpError
	}
	for _, item := range *data.Items {
		item = item
		itemValidated, err := item.ValidateWith()
		if err != nil {
			httpError := httpinfra.NewHTTPError(httpinfra.StatusInvalidArgument)
			httpError.SetMessage("error validating field [Items]")
			return nil, httpError
		}
		if !itemValidated.Valid {
			return itemValidated, nil
		}
	}
	return &httpinfra.ValidationResult{
		Valid: true,
	}, nil
}

// SetDefaults sets default values as defined in the API spec
func (data *ApplicationCollection) SetDefaults() {
}
