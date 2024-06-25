// Code generated by Signare OpenAPI generator. DO NOT EDIT

package httpinfra

import (
	"github.com/hyperledger-labs/signare/app/pkg/infra/httpinfra"
)

type UserUpdateSpec struct {
	Roles *[]string `json:"roles"`
	// Description of the resource.
	Description *string `json:"description,omitempty"`
}

// ValidateWith check whether UserUpdateSpec is valid
func (data UserUpdateSpec) ValidateWith() (*httpinfra.ValidationResult, *httpinfra.HTTPError) {
	if data.Roles == nil {
		httpError := httpinfra.NewHTTPError(httpinfra.StatusInvalidArgument)
		httpError.SetMessage("error validating field [roles]")
		return nil, httpError
	}
	for _, item := range *data.Roles {
		item = item
	}
	if data.Description != nil {
		if len(*data.Description) > 256 {
			httpError := httpinfra.NewHTTPError(httpinfra.StatusInvalidArgument)
			httpError.SetMessage("field [description] exceeds max length of 256")
			return nil, httpError
		}
	}
	return &httpinfra.ValidationResult{
		Valid: true,
	}, nil
}

// SetDefaults sets default values as defined in the API spec
func (data *UserUpdateSpec) SetDefaults() {
}
