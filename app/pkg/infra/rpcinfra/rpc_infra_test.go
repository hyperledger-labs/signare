package rpcinfra_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hyperledger-labs/signare/app/pkg/infra/rpcinfra"
	"github.com/hyperledger-labs/signare/app/pkg/infra/rpcinfra/rpcerrors"

	"github.com/stretchr/testify/require"
)

const (
	method = "rpc_method"
)

func TestRPCInfra_RegisterRPCHandlerFunc(t *testing.T) {
	defaultRPCRouterOptions := rpcinfra.DefaultRPCRouterOptions{}
	rpcRouter := rpcinfra.ProvideDefaultRPCRouter(defaultRPCRouterOptions)
	require.NotNil(t, rpcRouter)

	err := rpcRouter.RegisterRPCHandlerFunc(method, FooFunc)
	require.Nil(t, err)
	out, rpcError := rpcRouter.RPCHandler(method)
	require.Nil(t, rpcError)
	require.NotNil(t, out)
}

func TestRPCInfra_StartRPCServerWithRegisteredRPCHandler(t *testing.T) {
	defaultRPCRouterOptions := rpcinfra.DefaultRPCRouterOptions{}
	rpcRouter := rpcinfra.ProvideDefaultRPCRouter(defaultRPCRouterOptions)
	require.NotNil(t, rpcRouter)

	err := rpcRouter.RegisterRPCHandlerFunc(method, FooFunc)
	require.Nil(t, err)
	out, rpcError := rpcRouter.RPCHandler(method)
	require.Nil(t, rpcError)
	require.NotNil(t, out)

	httpRouter := rpcRouter.Router()
	httpRouter.HandleFunc("/", rpcRouter.HandleRPCRequest).Methods("POST")
}

type fooParams struct {
	Foo *string `json:"foo"`
}

func (fp *fooParams) SetParamsFrom(params []any) error {
	if len(params) != 1 {
		return fmt.Errorf("only one string is required")
	}
	f := params[0].(string)
	fp.Foo = &f
	return nil
}

func (fp *fooParams) ValidateParams() error {
	return nil
}

func FooFunc(_ context.Context, request rpcinfra.RPCRequest) (any, *rpcerrors.RPCError) {
	methodParams := new(fooParams)
	if err := rpcinfra.ProcessParams(request.Params, methodParams); err != nil {
		return nil, err
	}

	if methodParams.Foo == nil {
		return nil, rpcerrors.NewInvalidParams()
	}
	return methodParams.Foo, nil
}
