package httpinfra_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hyperledger-labs/signare/app/pkg/infra/httpinfra"
	signererrors "github.com/hyperledger-labs/signare/app/pkg/internal/errors"

	"github.com/stretchr/testify/require"
)

const (
	defaultMessage = "test message"
)

// TestDefaultHTTPResponseHandler_ProvideDefaultHTTPResponseHandler_Correct tests ProvideDefaultHTTPResponseHandler function
func TestDefaultHTTPResponseHandler_ProvideDefaultHTTPResponseHandler_Correct(t *testing.T) {
	options := httpinfra.DefaultHTTPResponseHandlerOptions{
		HTTPMetrics: httpinfra.DefaultHTTPMetrics{},
	}
	handler, err := httpinfra.ProvideDefaultHTTPResponseHandler(options)
	require.Nil(t, err)
	require.NotNil(t, handler)
}

// TestDefaultHTTPResponseHandler_HandleSuccessResponse_Correct tests HandleSuccessResponse function
func TestDefaultHTTPResponseHandler_HandleSuccessResponse_Correct(t *testing.T) {
	options := httpinfra.DefaultHTTPResponseHandlerOptions{
		HTTPMetrics: httpinfra.DefaultHTTPMetrics{},
	}
	handler, err := httpinfra.ProvideDefaultHTTPResponseHandler(options)
	require.Nil(t, err)
	require.NotNil(t, handler)

	responseInfo := httpinfra.ResponseInfo{
		ResponseType: httpinfra.ResponseTypeOk,
	}
	rr := httptest.NewRecorder()
	handler.HandleSuccessResponse(context.Background(), rr, responseInfo, nil)
	require.Equal(t, http.StatusOK, rr.Code)

	responseInfo = httpinfra.ResponseInfo{
		ResponseType: httpinfra.ResponseTypeCreated,
	}
	rr = httptest.NewRecorder()
	handler.HandleSuccessResponse(context.Background(), rr, responseInfo, nil)
	require.Equal(t, http.StatusCreated, rr.Code)
}

// TestDefaultHTTPResponseHandler_HandleErrorResponse_Correct tests HandleErrorResponse function
func TestDefaultHTTPResponseHandler_HandleErrorResponse_Correct(t *testing.T) {
	options := httpinfra.DefaultHTTPResponseHandlerOptions{
		HTTPMetrics: httpinfra.DefaultHTTPMetrics{},
	}
	handler, err := httpinfra.ProvideDefaultHTTPResponseHandler(options)
	require.Nil(t, err)
	require.NotNil(t, handler)

	rr := httptest.NewRecorder()
	receivedError := httpinfra.NewHTTPError(httpinfra.StatusPermissionDenied)
	handler.HandleErrorResponse(context.Background(), rr, receivedError)
	require.Equal(t, http.StatusForbidden, rr.Code)
}

// TestNewHTTPError tests NewHTTPError function
func TestNewHTTPError(t *testing.T) {
	err := httpinfra.NewHTTPError(httpinfra.StatusInvalidArgument).SetMessage(defaultMessage)
	require.Equal(t, defaultMessage, err.Error())
}

// TestNewHTTPErrorFromSignerError tests NewHTTPErrorFromUseCaseError function
func TestNewHTTPErrorFromSignerError(t *testing.T) {
	err := httpinfra.NewHTTPErrorFromUseCaseError(context.Background(), signererrors.Internal().WithMessage(defaultMessage))
	require.Equal(t, "INTERNAL: "+defaultMessage, err.OriginalError().Error())
}
