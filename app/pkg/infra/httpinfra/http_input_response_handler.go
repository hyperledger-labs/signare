package httpinfra

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/hyperledger-labs/signare/app/pkg/commons/logger"
)

var _ HTTPResponseHandler = (*DefaultHTTPResponseHandler)(nil)

// HandleSuccessResponse handles the HTTP response when the operation succeeded
func (d *DefaultHTTPResponseHandler) HandleSuccessResponse(ctx context.Context, w http.ResponseWriter, responseType ResponseInfo, responseData any) {
	statusCode, ok := d.responseMap[responseType.ResponseType]
	if !ok {
		logger.LogEntry(ctx).Errorf("unknown response type [%s]", responseType.ResponseType)
		statusCode = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if responseData != nil {
		err := json.NewEncoder(w).Encode(responseData)
		if err != nil {
			logger.LogEntry(ctx).Errorf("error encoding json response for data [%v]", responseData)
		}
	}
}

// HandleErrorResponse handles the HTTP response when the operation returned an error
func (d *DefaultHTTPResponseHandler) HandleErrorResponse(ctx context.Context, w http.ResponseWriter, receivedError *HTTPError) {
	// Since this handler function is only used through generated code we do not need to check for nils on the receivedError
	statusCode := receivedError.code

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := receivedError.toErrorResponse(ctx)

	if receivedError.Code() == http.StatusForbidden {
		d.httpMetrics.IncrementForbiddenAccessCounter(ctx)
	}

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		logger.LogEntry(ctx).Errorf("error encoding json response for error [%v]", response)
	}
}

// DefaultHTTPResponseHandler implementation of HTTPResponseHandler
type DefaultHTTPResponseHandler struct {
	responseMap map[ResponseType]int
	httpMetrics HTTPMetrics
}

// DefaultHTTPResponseHandlerOptions attributes to create an HTTPResponseHandler
type DefaultHTTPResponseHandlerOptions struct {
	// HTTPMetrics is a set of metrics to count forbidden access attempts
	HTTPMetrics HTTPMetrics
}

// ProvideDefaultHTTPResponseHandler returns a new DefaultHTTPResponseHandler
func ProvideDefaultHTTPResponseHandler(options DefaultHTTPResponseHandlerOptions) (*DefaultHTTPResponseHandler, error) {
	if options.HTTPMetrics == nil {
		return nil, errors.New("mandatory 'HTTPMetrics' not provided")
	}
	return &DefaultHTTPResponseHandler{
		responseMap: buildResponseMap(),
		httpMetrics: options.HTTPMetrics,
	}, nil
}

func buildResponseMap() map[ResponseType]int {
	responseMap := make(map[ResponseType]int)
	responseMap[ResponseTypeOk] = http.StatusOK
	responseMap[ResponseTypeCreated] = http.StatusCreated
	responseMap[ResponseTypeAccepted] = http.StatusAccepted
	return responseMap
}
