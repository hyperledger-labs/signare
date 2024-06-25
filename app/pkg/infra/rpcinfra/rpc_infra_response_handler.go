package rpcinfra

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/hyperledger-labs/signare/app/pkg/commons/logger"
	"github.com/hyperledger-labs/signare/app/pkg/infra/httpinfra"
	"github.com/hyperledger-labs/signare/app/pkg/infra/requestcontext"
	"github.com/hyperledger-labs/signare/app/pkg/infra/rpcinfra/rpcerrors"
	signererrors "github.com/hyperledger-labs/signare/app/pkg/internal/errors"
)

var _ httpinfra.HTTPResponseHandler = (*DefaultRPCInfraResponseHandler)(nil)

// HandleSuccessResponse handles the HTTP response when the operation succeeded
func (d DefaultRPCInfraResponseHandler) HandleSuccessResponse(ctx context.Context, w http.ResponseWriter, _ httpinfra.ResponseInfo, responseData any) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(responseData)
	if err != nil {
		logger.LogEntry(ctx).Errorf("error encoding json response for data [%v]", responseData)
	}
}

// HandleErrorResponse handles the HTTP response when the operation returned an error
func (d DefaultRPCInfraResponseHandler) HandleErrorResponse(ctx context.Context, w http.ResponseWriter, receivedError *httpinfra.HTTPError) {
	w.Header().Set("Content-Type", "application/json")

	// Typically, in JSON-RPC APIs, the default status code is 200 OK as error details are conveyed within the response body.
	// Non-200 status codes are utilized to signal issues with the HTTP transport itself rather than the semantics of the method call.
	// Only HTTP middlewares will set Code in an HTTPError.
	if receivedError.Code() == http.StatusBadRequest {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	rpcError, ok := rpcerrors.CastAsRPCError(receivedError.OriginalError())
	if !ok {
		rpcError = rpcerrors.NewInternalFromErr(receivedError)
	}

	if receivedError.Code() == http.StatusForbidden {
		d.httpMetrics.IncrementForbiddenAccessCounter(ctx)
		rpcError = rpcerrors.NewUnauthorizedFromErr(receivedError)
	}

	requestID, requestIDFromContextErr := requestcontext.RPCRequestIDFromContext(ctx)
	if requestIDFromContextErr != nil {
		logger.LogEntry(ctx).Warn("error obtaining the RPC request ID from request's context")
	}

	logRPCError(ctx, rpcError)
	response := &RPCResponse{
		RPCVersion: SupportedRPCVersion,
		Error:      rpcError,
		Result:     nil,
		ID:         requestID,
	}

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		logger.LogEntry(ctx).Errorf("error encoding json response for data [%v]", response)
	}
}

// logRPCError logs a RPCError with the correct log level.
func logRPCError(ctx context.Context, error *rpcerrors.RPCError) {
	if error.Code == rpcerrors.InternalErrorCode {
		if error.WrappedErr != nil {
			useCaseErr, ok := signererrors.CastAsUseCaseError(error.WrappedErr)
			if ok {
				logger.LogEntry(ctx).Error(useCaseErr.GetStack())
				return
			}
		}
		logger.LogEntry(ctx).Errorf("%+v", error)
		return
	}
	logger.LogEntry(ctx).Debugf("%+v", error)
}

// DefaultRPCInfraResponseHandler handles the RPC responses of the server
type DefaultRPCInfraResponseHandler struct {
	httpMetrics httpinfra.HTTPMetrics
}

// DefaultRPCInfraResponseHandlerOptions attributes to create an DefaultRPCInfraResponseHandler
type DefaultRPCInfraResponseHandlerOptions struct {
	// HTTPMetrics is a set of metrics to count forbidden access attempts
	HTTPMetrics httpinfra.HTTPMetrics
}

// ProvideDefaultRPCInfraResponseHandler creates a new DefaultRPCInfraResponseHandler with the provided options
func ProvideDefaultRPCInfraResponseHandler(options DefaultRPCInfraResponseHandlerOptions) (*DefaultRPCInfraResponseHandler, error) {
	if options.HTTPMetrics == nil {
		return nil, errors.New("mandatory 'HTTPMetrics' not provided")
	}

	return &DefaultRPCInfraResponseHandler{
		httpMetrics: options.HTTPMetrics,
	}, nil
}
