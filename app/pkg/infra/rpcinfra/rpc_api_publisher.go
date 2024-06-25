package rpcinfra

import (
	"errors"
)

// JSON-RPC supported methods by the signare
const (
	generateAccountMethod = "eth_generateAccount"
	removeAccountMethod   = "eth_removeAccount"
	listAccountsMethod    = "eth_accounts"
	signTransactionMethod = "eth_signTransaction"
)

// JSONRPCAPIPublisherOptions options to create a JSONRPCAPIRoutesPublished.
type JSONRPCAPIPublisherOptions struct {
	RPCRouter RPCRouter
	Handler   JSONRPCAPIHandler
}

// JSONRPCAPIRoutesPublished type for the JSON-RPC API published routes.
type JSONRPCAPIRoutesPublished int

// ProvideJSONRPCMethods creates a new JSONRPCAPIRoutesPublished.
func ProvideJSONRPCMethods(options JSONRPCAPIPublisherOptions) (JSONRPCAPIRoutesPublished, error) {
	if options.RPCRouter == nil {
		return 0, errors.New("mandatory 'RPCRouter' not provided")
	}
	if options.Handler == nil {
		return 0, errors.New("mandatory 'Handler' not provided")
	}

	// Register RPC handlers
	var err error
	err = options.RPCRouter.RegisterRPCHandlerFunc(generateAccountMethod, options.Handler.HandleGenerateAccount)
	if err != nil {
		return 0, err
	}
	err = options.RPCRouter.RegisterRPCHandlerFunc(removeAccountMethod, options.Handler.HandleRemoveAccount)
	if err != nil {
		return 0, err
	}
	err = options.RPCRouter.RegisterRPCHandlerFunc(listAccountsMethod, options.Handler.HandleListAccounts)
	if err != nil {
		return 0, err
	}
	err = options.RPCRouter.RegisterRPCHandlerFunc(signTransactionMethod, options.Handler.HandleSignTX)
	if err != nil {
		return 0, err
	}

	// HTTP Handler
	options.RPCRouter.Router().HandleFunc("/", options.RPCRouter.HandleRPCRequest).Methods("POST").Name("rpc.method")

	return 0, nil
}
