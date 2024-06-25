package rpcinfra

import (
	"context"
	"errors"

	"github.com/hyperledger-labs/signare/app/pkg/infra/requestcontext"
	"github.com/hyperledger-labs/signare/app/pkg/infra/rpcinfra/rpcerrors"
)

// JSONRPCAPIHandler handles the set of operations that are supported by the RPC protocol
type JSONRPCAPIHandler interface {
	// HandleGenerateAccount handles the generation of an Ethereum account.
	HandleGenerateAccount(ctx context.Context, r RPCRequest) (any, *rpcerrors.RPCError)
	// HandleRemoveAccount handles the removal of an Ethereum account.
	HandleRemoveAccount(ctx context.Context, r RPCRequest) (any, *rpcerrors.RPCError)
	// HandleListAccounts handles the listing of all the Ethereum accounts in an Application.
	HandleListAccounts(ctx context.Context, r RPCRequest) (any, *rpcerrors.RPCError)
	// HandleSignTX handles the signature of a transaction with an Ethereum account.
	HandleSignTX(ctx context.Context, r RPCRequest) (any, *rpcerrors.RPCError)
}

func (handler DefaultJSONRPCAPIHandler) HandleGenerateAccount(ctx context.Context, r RPCRequest) (any, *rpcerrors.RPCError) {
	reqParams := GenerateAccountRequestParams{}
	applicationID, err := requestcontext.ApplicationFromContext(ctx)
	if err != nil {
		return nil, rpcerrors.NewInternalFromErr(err)
	}
	reqParams.ApplicationID = *applicationID

	out, rpcErr := handler.adapter.AdaptGenerateAccount(ctx, reqParams)
	if rpcErr != nil {
		return nil, rpcErr
	}
	return &RPCResponse{
		RPCVersion: SupportedRPCVersion,
		ID:         r.ID,
		Result:     out,
	}, nil
}

func (handler DefaultJSONRPCAPIHandler) HandleRemoveAccount(ctx context.Context, r RPCRequest) (any, *rpcerrors.RPCError) {
	reqParams := RemoveAccountRequestParams{}
	if err := ProcessParams(r.Params, &reqParams); err != nil {
		return nil, err
	}
	err := reqParams.ValidateParams()
	if err != nil {
		return nil, rpcerrors.NewInvalidParamsFromErr(err)
	}

	applicationID, err := requestcontext.ApplicationFromContext(ctx)
	if err != nil {
		return nil, rpcerrors.NewInternalFromErr(err)
	}
	reqParams.ApplicationID = *applicationID

	out, rpcErr := handler.adapter.AdaptRemoveAccount(ctx, reqParams)
	if rpcErr != nil {
		return nil, rpcErr
	}
	return &RPCResponse{
		RPCVersion: SupportedRPCVersion,
		ID:         r.ID,
		Result:     out,
	}, nil
}

func (handler DefaultJSONRPCAPIHandler) HandleListAccounts(ctx context.Context, r RPCRequest) (any, *rpcerrors.RPCError) {
	reqParams := ListAccountsRequestParams{}
	applicationID, err := requestcontext.ApplicationFromContext(ctx)
	if err != nil {
		return nil, rpcerrors.NewInternalFromErr(err)
	}
	reqParams.ApplicationID = *applicationID

	out, rpcErr := handler.adapter.AdaptListAccounts(ctx, reqParams)
	if rpcErr != nil {
		return nil, rpcErr
	}
	return &RPCResponse{
		RPCVersion: SupportedRPCVersion,
		ID:         r.ID,
		Result:     out,
	}, nil
}

func (handler DefaultJSONRPCAPIHandler) HandleSignTX(ctx context.Context, r RPCRequest) (any, *rpcerrors.RPCError) {
	reqParams := SignTXRequestParams{}
	if err := ProcessParams(r.Params, &reqParams); err != nil {
		return nil, err
	}
	err := reqParams.ValidateParams()
	if err != nil {
		return nil, rpcerrors.NewInvalidParamsFromErr(err)
	}

	applicationID, err := requestcontext.ApplicationFromContext(ctx)
	if err != nil {
		return nil, rpcerrors.NewInternalFromErr(err)
	}
	reqParams.ApplicationID = *applicationID

	out, rpcErr := handler.adapter.AdaptSignTx(ctx, reqParams)
	if rpcErr != nil {
		return nil, rpcErr
	}
	return &RPCResponse{
		RPCVersion: SupportedRPCVersion,
		ID:         r.ID,
		Result:     out,
	}, nil
}

// DefaultJSONRPCAPIHandlerOptions are the attributes to build a DefaultJSONRPCAPIHandler
type DefaultJSONRPCAPIHandlerOptions struct {
	// Adapter  adapts the set of operations that are supported by the RPC protocol
	Adapter JSONRPCAPIAdapter
}

// DefaultJSONRPCAPIHandler is the default JSONRPCAPIHandler
type DefaultJSONRPCAPIHandler struct {
	adapter JSONRPCAPIAdapter
}

// NewDefaultJSONRPCAPIHandler creates a new DefaultJSONRPCAPIHandler from the provided options
func NewDefaultJSONRPCAPIHandler(options DefaultJSONRPCAPIHandlerOptions) (*DefaultJSONRPCAPIHandler, error) {
	if options.Adapter == nil {
		return nil, errors.New("mandatory 'Adapter' not provided")
	}

	return &DefaultJSONRPCAPIHandler{
		adapter: options.Adapter,
	}, nil
}

var _ JSONRPCAPIHandler = (*DefaultJSONRPCAPIHandler)(nil)
