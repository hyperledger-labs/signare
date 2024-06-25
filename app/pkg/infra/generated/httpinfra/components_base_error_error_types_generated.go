// Code generated by Signare OpenAPI generator. DO NOT EDIT

package httpinfra

import (
	"github.com/hyperledger-labs/signare/app/pkg/infra/httpinfra"
)

type BaseErrorError struct {
	// Response code. It is always aligned with http response code
	Code *int32 `json:"code"`
	// Error category. Next list shows all possible statuses (and related codes):   * 400    `INVALID_ARGUMENT` Client specified an invalid argument. Check error details message for more information.   * 400    `ALREADY_EXISTS`  The resource that a client tried to create already exists.   * 401    `UNAUTHENTICATED`  Request not authenticated due to missing, invalid, or expired OAuth token.   * 403    `PERMISSION_DENIED`  Client does not have sufficient permission. This can happen because the OAuth token does not have the right scopes, the client doesn't have permission, or the API has not been enabled for the client project.   * 404    `NOT_FOUND`  A specified resource is not found, or the request is rejected by undisclosed reasons, such as whitelisting.   * 412    `PRECONDITION_FAILED`  Request can not be executed in the current system state.   * 429    `TOO_MANY_REQ` Too many incoming request at this moment.   * 500    `INTERNAL` Internal server error. Typically a server bug.   * 501    `NOT_IMPLEMENTED`  API method not implemented by the server.   * 502    `BAD_GATEWAY`  The server, while acting as a gateway or proxy, received an invalid response from the upstream server.   * 503    `UNAVAILABLE`  Service unavailable. Typically the server is down.   * 504    `TIMEOUT`  Upstream request deadline exceeded.
	Status  *string                `json:"status"`
	Details *BaseErrorErrorDetails `json:"details"`
}

// ValidateWith check whether BaseErrorError is valid
func (data BaseErrorError) ValidateWith() (*httpinfra.ValidationResult, *httpinfra.HTTPError) {
	if data.Code == nil {
		httpError := httpinfra.NewHTTPError(httpinfra.StatusInvalidArgument)
		httpError.SetMessage("error validating field [code]")
		return nil, httpError
	}
	if data.Status == nil {
		httpError := httpinfra.NewHTTPError(httpinfra.StatusInvalidArgument)
		httpError.SetMessage("error validating field [status]")
		return nil, httpError
	}
	if data.Details == nil {
		httpError := httpinfra.NewHTTPError(httpinfra.StatusInvalidArgument)
		httpError.SetMessage("error validating field [details]")
		return nil, httpError
	}
	validatedDetails, errDetails := data.Details.ValidateWith()
	if errDetails != nil {
		httpError := httpinfra.NewHTTPError(httpinfra.StatusInvalidArgument)
		httpError.SetMessage("error validating field [details]")
		return nil, httpError
	}
	if validatedDetails != nil && !validatedDetails.Valid {
		return validatedDetails, nil
	}
	return &httpinfra.ValidationResult{
		Valid: true,
	}, nil
}

// SetDefaults sets default values as defined in the API spec
func (data *BaseErrorError) SetDefaults() {
	data.Details.SetDefaults()
}
