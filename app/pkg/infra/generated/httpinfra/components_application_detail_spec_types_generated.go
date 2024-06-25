// Code generated by Signare OpenAPI generator. DO NOT EDIT

package httpinfra

import (
	"github.com/hyperledger-labs/signare/app/pkg/infra/httpinfra"
)

type ApplicationDetailSpec struct {
	// The chain identifier with which the application interacts. It must be a valid integer.
	ChainId *string `json:"chainId"`
	// Description of the resource.
	Description *string `json:"description"`
}

// ValidateWith check whether ApplicationDetailSpec is valid
func (data ApplicationDetailSpec) ValidateWith() (*httpinfra.ValidationResult, *httpinfra.HTTPError) {
	if data.ChainId == nil {
		httpError := httpinfra.NewHTTPError(httpinfra.StatusInvalidArgument)
		httpError.SetMessage("error validating field [chainId]")
		return nil, httpError
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
func (data *ApplicationDetailSpec) SetDefaults() {
}