// Code generated by Signare OpenAPI generator. DO NOT EDIT

package httpinfra

import (
	"github.com/hyperledger-labs/signare/app/pkg/infra/httpinfra"
)

type UserDetailSpec struct {
	Roles    *[]string `json:"roles"`
	Accounts *[]string `json:"accounts"`
	// Description of the resource.
	Description *string `json:"description"`
}

// ValidateWith check whether UserDetailSpec is valid
func (data UserDetailSpec) ValidateWith() (*httpinfra.ValidationResult, *httpinfra.HTTPError) {
	if data.Roles == nil {
		httpError := httpinfra.NewHTTPError(httpinfra.StatusInvalidArgument)
		httpError.SetMessage("error validating field [roles]")
		return nil, httpError
	}
	for _, item := range *data.Roles {
		item = item
	}
	if data.Accounts == nil {
		httpError := httpinfra.NewHTTPError(httpinfra.StatusInvalidArgument)
		httpError.SetMessage("error validating field [accounts]")
		return nil, httpError
	}
	for _, item := range *data.Accounts {
		item = item
	}
	if data.Description == nil {
		httpError := httpinfra.NewHTTPError(httpinfra.StatusInvalidArgument)
		httpError.SetMessage("error validating field [description]")
		return nil, httpError
	}
	return &httpinfra.ValidationResult{
		Valid: true,
	}, nil
}

// SetDefaults sets default values as defined in the API spec
func (data *UserDetailSpec) SetDefaults() {
}