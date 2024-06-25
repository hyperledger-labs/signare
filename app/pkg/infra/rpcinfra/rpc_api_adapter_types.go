package rpcinfra

import (
	"errors"
	"fmt"
)

// GenerateAccountRequestParams request definition
type GenerateAccountRequestParams struct {
	// ApplicationID requesting the Ethereum account generation.
	ApplicationID string
}

// RemoveAccountRequestParams request definition
type RemoveAccountRequestParams struct {
	// ApplicationID requesting the Ethereum account removal.
	ApplicationID string
	// Address is the Ethereum account to be removed.
	Address string `json:"address"`
}

func (p *RemoveAccountRequestParams) SetParamsFrom(params []any) error {
	if len(params) != 1 {
		return fmt.Errorf("only one object is expected")
	}
	paramMap := params[0].(map[string]any)
	addressParam, ok := paramMap["address"]
	if !ok {
		return errors.New("missing required field [address]")
	}
	address, ok := addressParam.(string)
	if !ok {
		return errors.New("[address] must be of type string")
	}
	p.Address = address
	return nil
}

func (p *RemoveAccountRequestParams) ValidateParams() error {
	if len(p.Address) == 0 {
		return errors.New("[address] cannot be nil")
	}
	return nil
}

// ListAccountsRequestParams request definition
type ListAccountsRequestParams struct {
	ApplicationID string
}

// SignTXRequestParams request definition
type SignTXRequestParams struct {
	ApplicationID string
	// From address
	From string `json:"from"`
	// To address
	To *string `json:"to"`
	// Gas amount to use for transaction execution
	Gas *string `json:"gas"`
	// GasPrice to use for each paid gas
	GasPrice *string `json:"gasPrice"`
	// Value amount sent with this transaction
	Value *string `json:"value"`
	// Data arguments packed according to json rpc standard
	Data string `json:"data"`
	// Nonce integer to identify request
	Nonce string `json:"nonce"`
}

func (p *SignTXRequestParams) SetParamsFrom(params []any) error {
	if len(params) != 1 {
		return fmt.Errorf("only one object is expected")
	}
	paramMap := params[0].(map[string]any)

	// Required fields
	fromParam, ok := paramMap["from"]
	if !ok {
		return errors.New("missing required field [from]")
	}
	from, ok := fromParam.(string)
	if !ok {
		return errors.New("[from] must be of type string")
	}
	p.From = from

	dataParam, ok := paramMap["data"]
	if !ok {
		return errors.New("missing required field [data]")
	}
	data, ok := dataParam.(string)
	if !ok {
		return errors.New("[data] must be of type string")
	}
	p.Data = data

	nonceParam, ok := paramMap["nonce"]
	if !ok {
		return errors.New("missing required field [nonce]")
	}
	nonce, ok := nonceParam.(string)
	if !ok {
		return errors.New("[nonce] must be of type string")
	}
	p.Nonce = nonce

	// Optional fields
	var to, gas, gasPrice, value string

	toParam, ok := paramMap["to"]
	if ok {
		to, ok = toParam.(string)
		if !ok {
			return errors.New("[to] must be of type string")
		}
		p.To = &to
	}

	gasParam, ok := paramMap["gas"]
	if ok {
		gas, ok = gasParam.(string)
		if !ok {
			return errors.New("[gas] must be of type string")
		}
		p.Gas = &gas
	}

	gasPriceParam, ok := paramMap["gasPrice"]
	if ok {
		gasPrice, ok = gasPriceParam.(string)
		if !ok {
			return errors.New("[gasPrice] must be of type string")
		}
		p.GasPrice = &gasPrice
	}

	valueParam, ok := paramMap["value"]
	if ok {
		value, ok = valueParam.(string)
		if !ok {
			return errors.New("[value] must be of type string")
		}
		p.Value = &value
	}
	return nil
}

func (p *SignTXRequestParams) ValidateParams() error {
	if len(p.From) == 0 {
		return errors.New("[from] cannot be nil")
	}
	if len(p.Nonce) == 0 {
		return errors.New("[nonce] cannot be nil")
	}
	return nil
}
