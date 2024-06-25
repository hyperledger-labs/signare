// Package rpcin defines the implementation of the input adapters for the JSON RPC infra.
package rpcin

import (
	"context"
	"errors"
	"fmt"

	"github.com/hyperledger-labs/signare/app/pkg/entities"
	"github.com/hyperledger-labs/signare/app/pkg/entities/address"
	"github.com/hyperledger-labs/signare/app/pkg/infra/rpcinfra"
	"github.com/hyperledger-labs/signare/app/pkg/infra/rpcinfra/rpcerrors"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/hsmconnection"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/hsmconnector"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/user"
)

var _ rpcinfra.JSONRPCAPIAdapter = new(DefaultAPIAdapter)

func (adapter *DefaultAPIAdapter) AdaptGenerateAccount(ctx context.Context, data rpcinfra.GenerateAccountRequestParams) (*string, *rpcerrors.RPCError) {
	input := hsmconnection.ByApplicationInput{
		ApplicationID: data.ApplicationID,
	}
	hsmConnection, err := adapter.hsmConnectionResolver.ByApplication(ctx, input)
	if err != nil {
		return nil, adaptError(err)
	}

	generateAddressInput := hsmconnector.GenerateAddressInput{
		SlotConnectionData: hsmconnector.SlotConnectionData{
			Pin:        hsmConnection.Pin,
			Slot:       hsmConnection.Slot,
			ModuleKind: hsmconnector.ModuleKind(hsmConnection.ModuleKind),
			ChainID:    hsmConnection.ChainID,
		},
	}
	out, err := adapter.hsmConnector.GenerateAddress(ctx, generateAddressInput)
	if err != nil {
		return nil, adaptError(err)
	}
	response := out.Address.String()
	return &response, nil
}

func (adapter *DefaultAPIAdapter) AdaptRemoveAccount(ctx context.Context, data rpcinfra.RemoveAccountRequestParams) (*string, *rpcerrors.RPCError) {
	addr, err := address.NewFromHexString(data.Address)
	if err != nil {
		return nil, rpcerrors.NewInvalidParamsFromErr(err)
	}
	input := user.DeleteAllAccountsForAddressInput{
		Address:       addr,
		ApplicationID: data.ApplicationID,
	}
	_, deleteErr := adapter.accountUseCase.DeleteAllAccountsForAddress(ctx, input)
	if deleteErr != nil {
		return nil, adaptError(deleteErr)
	}
	response := addr.String()
	return &response, nil
}

func (adapter *DefaultAPIAdapter) AdaptListAccounts(ctx context.Context, data rpcinfra.ListAccountsRequestParams) ([]string, *rpcerrors.RPCError) {
	input := hsmconnection.ByApplicationInput{
		ApplicationID: data.ApplicationID,
	}
	hsmConnection, err := adapter.hsmConnectionResolver.ByApplication(ctx, input)
	if err != nil {
		return nil, adaptError(err)
	}

	listAddressesInput := hsmconnector.ListAddressesInput{
		SlotConnectionData: hsmconnector.SlotConnectionData{
			Pin:        hsmConnection.Pin,
			Slot:       hsmConnection.Slot,
			ModuleKind: hsmconnector.ModuleKind(hsmConnection.ModuleKind),
			ChainID:    hsmConnection.ChainID,
		},
	}
	out, err := adapter.hsmConnector.ListAddresses(ctx, listAddressesInput)
	if err != nil {
		return nil, adaptError(err)
	}
	response := make([]string, len(out.Items))
	for i, addr := range out.Items {
		response[i] = addr.String()
	}
	return response, nil
}

func (adapter *DefaultAPIAdapter) AdaptSignTx(ctx context.Context, data rpcinfra.SignTXRequestParams) (*string, *rpcerrors.RPCError) {
	byApplicationInput := hsmconnection.ByApplicationInput{
		ApplicationID: data.ApplicationID,
	}
	hsmConnection, err := adapter.hsmConnectionResolver.ByApplication(ctx, byApplicationInput)
	if err != nil {
		return nil, adaptError(err)
	}

	signTxInput := hsmconnector.SignTxInput{
		SlotConnectionData: hsmconnector.SlotConnectionData{
			Pin:        hsmConnection.Pin,
			Slot:       hsmConnection.Slot,
			ModuleKind: hsmconnector.ModuleKind(hsmConnection.ModuleKind),
			ChainID:    hsmConnection.ChainID,
		},
	}
	if len(data.Data) == 0 {
		emptyBytes := entities.NewHexBytes([]byte{})
		signTxInput.Data = *emptyBytes
	} else {
		inputData, encodeDataErr := entities.NewHexBytesFromString(data.Data)
		if encodeDataErr != nil {
			return nil, rpcerrors.NewInvalidParamsFromErr(fmt.Errorf("invalid [data]: %w", encodeDataErr))
		}
		signTxInput.Data = inputData
	}

	nonce, err := entities.NewHexUInt64FromString(data.Nonce)
	if err != nil {
		return nil, rpcerrors.NewInvalidParamsFromErr(fmt.Errorf("invalid [nonce]: %w", err))
	}
	signTxInput.Nonce = nonce

	from, err := address.NewFromHexString(data.From)
	if err != nil {
		return nil, rpcerrors.NewInvalidParamsFromErr(fmt.Errorf("invalid [from]: %w", err))
	}
	signTxInput.From = from

	if data.To != nil {
		to, errTo := address.NewFromHexString(*data.To)
		if errTo != nil {
			return nil, rpcerrors.NewInvalidParamsFromErr(fmt.Errorf("invalid [to]: %w", errTo))
		}
		signTxInput.To = &to
	}

	if data.Gas != nil {
		gas, errGas := entities.NewHexUInt64FromString(*data.Gas)
		if errGas != nil {
			return nil, rpcerrors.NewInvalidParamsFromErr(fmt.Errorf("invalid [gas]: %w", errGas))
		}
		signTxInput.Gas = &gas
	}

	if data.GasPrice != nil {
		gasPrice, errGasPrice := entities.NewHexInt256FromString(*data.GasPrice)
		if errGasPrice != nil {
			return nil, rpcerrors.NewInvalidParamsFromErr(fmt.Errorf("invalid [gasPrice]: %w", errGasPrice))
		}
		signTxInput.GasPrice = gasPrice
	}

	if data.Value != nil && len(*data.Value) > 0 {
		value, errValue := entities.NewHexInt256FromString(*data.Value)
		if errValue != nil {
			return nil, rpcerrors.NewInvalidParamsFromErr(fmt.Errorf("invalid [value]: %w", errValue))
		}
		signTxInput.Value = value
	}

	out, err := adapter.hsmConnector.SignTx(ctx, signTxInput)
	if err != nil {
		return nil, adaptError(err)
	}
	response := out.SignedTx
	return &response, nil
}

// DefaultAPIAdapter implements JSONRPCAPIAdapter.
type DefaultAPIAdapter struct {
	accountUseCase        user.AccountUseCase
	hsmConnectionResolver hsmconnection.Resolver
	hsmConnector          hsmconnector.HSMConnector
}

// DefaultAPIAdapterOptions options to create a new DefaultAPIAdapter.
type DefaultAPIAdapterOptions struct {
	AccountUseCase        user.AccountUseCase
	HSMConnectionResolver hsmconnection.Resolver
	HSMConnector          hsmconnector.HSMConnector
}

// NewDefaultAPIAdapter creates a new DefaultAPIAdapter instance.
func NewDefaultAPIAdapter(options DefaultAPIAdapterOptions) (*DefaultAPIAdapter, error) {
	if options.AccountUseCase == nil {
		return nil, errors.New("mandatory 'AccountUseCase' not provided")
	}
	if options.HSMConnectionResolver == nil {
		return nil, errors.New("mandatory 'Resolver' not provided")
	}
	if options.HSMConnector == nil {
		return nil, errors.New("mandatory 'HSMConnector' not provided")
	}

	return &DefaultAPIAdapter{
		accountUseCase:        options.AccountUseCase,
		hsmConnectionResolver: options.HSMConnectionResolver,
		hsmConnector:          options.HSMConnector,
	}, nil
}
