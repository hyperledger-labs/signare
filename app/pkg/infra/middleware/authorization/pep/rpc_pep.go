package pep

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/hyperledger-labs/signare/app/pkg/commons/logger"
	"github.com/hyperledger-labs/signare/app/pkg/entities/address"
	"github.com/hyperledger-labs/signare/app/pkg/infra/httpinfra"
	"github.com/hyperledger-labs/signare/app/pkg/infra/requestcontext"
	"github.com/hyperledger-labs/signare/app/pkg/utils"
)

// AuthorizeAccount checks if a user is authorized to use an account if it's performing an 'eth_signTransaction' action
func (policyEnforcementPoint *RPCPolicyEnforcementPoint) AuthorizeAccount(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		user, err := requestcontext.UserFromContext(ctx)
		if err != nil {
			policyEnforcementPoint.responseHandler.HandleErrorResponse(ctx, w, httpinfra.NewHTTPErrorFromError(ctx, err, httpinfra.StatusPermissionDenied))
			return
		}
		application, err := requestcontext.ApplicationFromContext(ctx)
		if err != nil {
			policyEnforcementPoint.responseHandler.HandleErrorResponse(ctx, w, httpinfra.NewHTTPErrorFromError(ctx, err, httpinfra.StatusPermissionDenied))
			return
		}
		actionID, err := requestcontext.ActionFromContext(ctx)
		if err != nil {
			policyEnforcementPoint.responseHandler.HandleErrorResponse(ctx, w, httpinfra.NewHTTPErrorFromError(ctx, err, httpinfra.StatusPermissionDenied))
			return
		}

		if !strings.Contains(*actionID, "eth_signTransaction") {
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		var authorizeAccountRPCBody AuthorizeAccountRPCBody
		err = utils.ReadAndResetCloser(&r.Body, &authorizeAccountRPCBody)
		if err != nil {
			policyEnforcementPoint.responseHandler.HandleErrorResponse(r.Context(), w, httpinfra.NewHTTPErrorFromError(ctx, err, httpinfra.StatusInvalidArgument))
			return
		}

		addr, err := getAddressFromParamsArray(ctx, authorizeAccountRPCBody)
		if err != nil {
			addr, err = getAddressFromParamsObject(ctx, authorizeAccountRPCBody)
			if err != nil {
				policyEnforcementPoint.responseHandler.HandleErrorResponse(r.Context(), w, httpinfra.NewHTTPErrorFromError(ctx, err, httpinfra.StatusInvalidArgument))
				return
			}
		}

		authorizeAccountInput := AuthorizeAccountUserInput{
			UserID:        *user,
			ApplicationID: *application,
			Address:       *addr,
		}
		_, err = policyEnforcementPoint.accountUserPolicyDecisionPointAdapter.AuthorizeAccountUser(ctx, authorizeAccountInput)
		if err != nil {
			logger.LogEntry(ctx).Errorf("user [%s] is not authorized to use request's account [%s]", authorizeAccountInput.UserID, authorizeAccountInput.Address.String())
			policyEnforcementPoint.responseHandler.HandleErrorResponse(ctx, w, httpinfra.NewHTTPErrorFromError(ctx, err, httpinfra.StatusPermissionDenied))
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getAddressFromParamsArray(ctx context.Context, params AuthorizeAccountRPCBody) (*address.Address, error) {
	var rpcAddress []AuthorizeAccountRPCParams
	err := json.Unmarshal(params.Params, &rpcAddress)
	if err != nil {
		return nil, err
	}
	addr, err := address.NewFromHexString(rpcAddress[0].From)
	if err != nil {
		logger.LogEntry(ctx).Errorf("invalid [from] address: %s", rpcAddress[0].From)
		return nil, err
	}
	return &addr, nil
}

func getAddressFromParamsObject(ctx context.Context, params AuthorizeAccountRPCBody) (*address.Address, error) {
	var rpcAddress AuthorizeAccountRPCParams
	err := json.Unmarshal(params.Params, &rpcAddress)
	if err != nil {
		return nil, err
	}
	addr, err := address.NewFromHexString(rpcAddress.From)
	if err != nil {
		logger.LogEntry(ctx).Errorf("invalid [from] address: %s", rpcAddress.From)
		return nil, err
	}
	return &addr, nil
}

// RPCPolicyEnforcementPointOptions are the set of fields to create an RPCPolicyEnforcementPoint
type RPCPolicyEnforcementPointOptions struct {
	// ResponseHandler exposes functionality to handle HTTP responses
	ResponseHandler httpinfra.HTTPResponseHandler
	// UserPolicyDecisionPointPort is a port to adapt authorization checks
	UserPolicyDecisionPointAdapter UserPolicyDecisionPointPort
	// AccountUserPolicyDecisionPointPort is a port to adapt account usage authorization checks
	AccountUserPolicyDecisionPointAdapter AccountUserPolicyDecisionPointPort
}

// RPCPolicyEnforcementPoint checks if the user is authorized to perform a given action
type RPCPolicyEnforcementPoint struct {
	responseHandler                       httpinfra.HTTPResponseHandler
	userPolicyDecisionPointAdapter        UserPolicyDecisionPointPort
	accountUserPolicyDecisionPointAdapter AccountUserPolicyDecisionPointPort
}

// ProvideRPCPolicyEnforcementPoint provides an instance of an RPCPolicyEnforcementPoint
func ProvideRPCPolicyEnforcementPoint(options RPCPolicyEnforcementPointOptions) (*RPCPolicyEnforcementPoint, error) {
	if options.ResponseHandler == nil {
		return nil, errors.New("mandatory 'ResponseHandler' not provided")
	}
	if options.UserPolicyDecisionPointAdapter == nil {
		return nil, errors.New("mandatory 'UserPolicyDecisionPointAdapter' not provided")
	}
	if options.AccountUserPolicyDecisionPointAdapter == nil {
		return nil, errors.New("mandatory 'AccountUserPolicyDecisionPointAdapter' not provided")
	}

	return &RPCPolicyEnforcementPoint{
		responseHandler:                       options.ResponseHandler,
		userPolicyDecisionPointAdapter:        options.UserPolicyDecisionPointAdapter,
		accountUserPolicyDecisionPointAdapter: options.AccountUserPolicyDecisionPointAdapter,
	}, nil
}
