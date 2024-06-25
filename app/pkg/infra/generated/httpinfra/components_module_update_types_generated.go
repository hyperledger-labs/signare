// Code generated by Signare OpenAPI generator. DO NOT EDIT

package httpinfra

import (
	"github.com/hyperledger-labs/signare/app/pkg/infra/httpinfra"
)

type ModuleUpdate struct {
	Meta *ResourceMetaUpdate `json:"meta"`
	Spec *ModuleSpec         `json:"spec"`
}

// ValidateWith check whether ModuleUpdate is valid
func (data ModuleUpdate) ValidateWith() (*httpinfra.ValidationResult, *httpinfra.HTTPError) {
	if data.Meta == nil {
		httpError := httpinfra.NewHTTPError(httpinfra.StatusInvalidArgument)
		httpError.SetMessage("error validating field [meta]")
		return nil, httpError
	}
	validatedMeta, errMeta := data.Meta.ValidateWith()
	if errMeta != nil {
		httpError := httpinfra.NewHTTPError(httpinfra.StatusInvalidArgument)
		httpError.SetMessage("error validating field [meta]")
		return nil, httpError
	}
	if validatedMeta != nil && !validatedMeta.Valid {
		return validatedMeta, nil
	}
	if data.Spec == nil {
		httpError := httpinfra.NewHTTPError(httpinfra.StatusInvalidArgument)
		httpError.SetMessage("error validating field [spec]")
		return nil, httpError
	}
	validatedSpec, errSpec := data.Spec.ValidateWith()
	if errSpec != nil {
		httpError := httpinfra.NewHTTPError(httpinfra.StatusInvalidArgument)
		httpError.SetMessage("error validating field [spec]")
		return nil, httpError
	}
	if validatedSpec != nil && !validatedSpec.Valid {
		return validatedSpec, nil
	}
	return &httpinfra.ValidationResult{
		Valid: true,
	}, nil
}

// SetDefaults sets default values as defined in the API spec
func (data *ModuleUpdate) SetDefaults() {
	data.Meta.SetDefaults()
	data.Spec.SetDefaults()
}
