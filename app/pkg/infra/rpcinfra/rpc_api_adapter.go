package rpcinfra

import (
	"context"

	"github.com/hyperledger-labs/signare/app/pkg/infra/rpcinfra/rpcerrors"
)

// JSONRPCAPIAdapter adapts the set of operations that are supported by the RPC protocol
type JSONRPCAPIAdapter interface {
	// AdaptGenerateAccount adapts the generation of an Ethereum account.
	AdaptGenerateAccount(ctx context.Context, data GenerateAccountRequestParams) (*string, *rpcerrors.RPCError)
	// AdaptRemoveAccount adapts the removal of an Ethereum account.
	AdaptRemoveAccount(ctx context.Context, data RemoveAccountRequestParams) (*string, *rpcerrors.RPCError)
	// AdaptListAccounts adapts the listing of all the Ethereum accounts in an Application.
	AdaptListAccounts(ctx context.Context, data ListAccountsRequestParams) ([]string, *rpcerrors.RPCError)
	// AdaptSignTx adapts the signature of a transaction with an Ethereum account.
	AdaptSignTx(ctx context.Context, data SignTXRequestParams) (*string, *rpcerrors.RPCError)
}
