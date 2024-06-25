// Code generated by Signare OpenAPI generator. DO NOT EDIT

package httpinfra

import (
	"github.com/hyperledger-labs/signare/app/pkg/infra/httpinfra"
)

type SlotUpdatePinSpec struct {
	// PIN that provides access to the slot number inside the HSM.
	Pin *string `json:"pin"`
}

// ValidateWith check whether SlotUpdatePinSpec is valid
func (data SlotUpdatePinSpec) ValidateWith() (*httpinfra.ValidationResult, *httpinfra.HTTPError) {
	if data.Pin == nil {
		httpError := httpinfra.NewHTTPError(httpinfra.StatusInvalidArgument)
		httpError.SetMessage("error validating field [pin]")
		return nil, httpError
	}
	return &httpinfra.ValidationResult{
		Valid: true,
	}, nil
}

// SetDefaults sets default values as defined in the API spec
func (data *SlotUpdatePinSpec) SetDefaults() {
}
