package hsmconnector

import (
	"math/big"

	"github.com/hyperledger-labs/signare/app/pkg/commons/rlp"
	"github.com/hyperledger-labs/signare/app/pkg/entities"
	"github.com/hyperledger-labs/signare/app/pkg/entities/address"
	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"
)

// PKCS11Library path to the library to connect to a PKCS11 compatible HSM.
type PKCS11Library string

// ModuleKind HSM type.
type ModuleKind string

const (
	SoftHSMModuleKind ModuleKind = "SoftHSM"
)

// CreateInput input data to create a new instance using the factory.
type CreateInput struct {
	ModuleKind ModuleKind
}

// PKCS11ConnectionDetails HSM connection details.
type PKCS11ConnectionDetails struct {
	// Configuration details for the HSM.
	Configuration PKCS11ConnectionDetailsConfiguration
	// Slot to be accessed.
	Slot string
	// Pin that grants access to the slot.
	Pin string
}

// CloseInput input to close all the signature manager resources.
type CloseInput struct {
}

// CloseOutput input to close all the signature manager resources.
type CloseOutput struct {
}

// PKCS11ConnectionDetailsConfiguration configuration details for the HSM.
type PKCS11ConnectionDetailsConfiguration struct{}

// GenerateAddressInput for account generation requests.
type GenerateAddressInput struct {
	// SlotConnectionData configuration to connect to a slot.
	SlotConnectionData
}

// SlotConnectionData configuration to connect to a slot.
type SlotConnectionData struct {
	// Slot to be accessed.
	Slot string `valid:"required"`
	// Pin that grants access to the slot.
	Pin string `valid:"required"`
	// ModuleKind of the Hardware Security Module.
	ModuleKind ModuleKind `valid:"in(SoftHSM)"`
	// ChainID id of the chain.
	ChainID entities.Int256 `valid:"required"`
}

// GenerateAddressOutput for account generation responses.
type GenerateAddressOutput struct {
	// Address an Ethereum account to interact with the network.
	Address address.Address `json:"address"`
}

// RemoveAddressInput for account removal requests.
type RemoveAddressInput struct {
	// SlotConnectionData configuration to connect to a slot.
	SlotConnectionData
	// Address an Ethereum account to interact with the network.
	Address address.Address `valid:"address"`
}

// RemoveAddressOutput for address removal responses.
type RemoveAddressOutput struct {
	// Address an Ethereum account to interact with the network.
	Address address.Address `json:"address"`
}

// ListAddressesInput for account listing requests.
type ListAddressesInput struct {
	// SlotConnectionData configuration to connect to a slot.
	SlotConnectionData
}

// ListAddressesOutput for account listing responses.
type ListAddressesOutput struct {
	// Items is an array of Ethereum accounts to interact with the network.
	Items []address.Address `json:"items"`
}

// SignTxInput for transaction signing requests.
type SignTxInput struct {
	// SlotConnectionData configuration to connect to a slot.
	SlotConnectionData
	// From address.
	From address.Address `valid:"address"`
	// To address.
	To *address.Address `valid:"optional"`
	// Gas amount to use for transaction execution.
	Gas *entities.HexUInt64 `valid:"optional"`
	// GasPrice to use for each paid gas.
	GasPrice *entities.HexInt256 `valid:"optional"`
	// Value amount sent with this transaction.
	Value *entities.HexInt256 `valid:"optional"`
	// Data arguments packed according to JSON RPC standard.
	Data entities.HexBytes // it can be empty (byte array of length 0) in eth-transfers
	// Nonce integer to identify request.
	Nonce entities.HexUInt64
}

// SignTxOutput for transaction signing responses.
type SignTxOutput struct {
	// SignedTx an encrypted transaction with the corresponding private key of the Ethereum account.
	SignedTx string
	// Transaction represents an Ethereum transaction.
	Transaction EthereumTransaction
}

// CloseAllInput input to close all the signature manager resources.
type CloseAllInput struct {
}

// CloseAllOutput input to close all the signature manager resources.
type CloseAllOutput struct {
}

// IsAliveInput input to check the availability of the HSM slot.
type IsAliveInput struct {
	// Slot to be accessed.
	Slot string `valid:"required"`
	// Pin that grants access to the slot.
	Pin string `valid:"required"`
	// ModuleKind of the Hardware Security Module.
	ModuleKind ModuleKind `valid:"in(SoftHSM)"`
}

// IsAliveOutput whether the slot is available.
type IsAliveOutput struct {
	//IsAlive is true if the slot is reachable.
	IsAlive bool
}

// ResetInput input to reset the connection with the HSM library.
type ResetInput struct {
	// ModuleKind is the kind of the module that will be reset.
	ModuleKind ModuleKind
}

// ResetOutput output from the reset operation.
type ResetOutput struct {
}

// EthereumTransaction represents an Ethereum transaction.
type EthereumTransaction struct {
	// From address.
	From address.Address
	// To address.
	To *address.Address
	// Gas amount to use for transaction execution.
	Gas entities.HexUInt64
	// GasPrice to use for each paid gas.
	GasPrice entities.HexInt256
	// Value amount sent with this transaction.
	Value *entities.HexInt256
	// Data arguments packed according to json rpc standard.
	Data entities.HexBytes
	// Nonce integer to identify request.
	Nonce entities.HexUInt64
	// ChainID id of the blockchain network where the transaction is sent to.
	ChainID entities.HexInt256
	// Signature Ethereum transaction signature.
	Signature *EthereumTransactionSignature
}

// EthereumTransactionSignature represents an Ethereum transaction signature.
type EthereumTransactionSignature struct {
	V entities.Int256
	R entities.Int256
	S entities.Int256
}

// RLPEncode RLP encodes the Ethereum transaction (including its signature) according to EIP-155. This function fails if the transaction doesn't have a signature yet.
// As a summary, the result is rlp(nonce, gasPrice, gas, to, value, data, V, R, S)
func (tx EthereumTransaction) RLPEncode() (*entities.HexBytes, error) {
	if tx.Signature == nil {
		return nil, errors.Internal().WithMessage("tx doesn't have a signature so it can't be RLP encoded")
	}
	nonce, err := entities.NewHexBytesFromString(hexStringEvenLength(tx.Nonce.String()))
	if err != nil {
		return nil, errors.Internal().WithMessage("could not convert 'nonce' to hex bytes")
	}
	nonceBytes := nonce.Bytes()
	if tx.Nonce.Uint64() == 0 {
		nonceBytes = []byte{}
	}

	var gasPrice *big.Int
	if tx.GasPrice.BigInt().Sign() != 0 {
		gasPrice = tx.GasPrice.BigInt()
	}

	gas, err := entities.NewHexBytesFromString(hexStringEvenLength(tx.Gas.String()))
	if err != nil {
		return nil, errors.Internal().WithMessage("could not convert 'gas' to hex bytes")
	}

	var toBytes []byte
	if tx.To != nil {
		hexBytes, newErr := entities.NewHexBytesFromString(tx.To.String())
		if newErr != nil {
			return nil, errors.Internal().WithMessage("could not convert 'to' to hex bytes")
		}
		toBytes = hexBytes.Bytes()
	} else {
		toBytes = []byte{}
	}

	var value *big.Int
	if tx.Value != nil && tx.Value.BigInt().Sign() != 0 {
		value = tx.Value.BigInt()
	}

	var data []byte
	if len(tx.Data.Bytes()) > 0 {
		hexBytes, newErr := entities.NewHexBytesFromString(tx.Data.String())
		if newErr != nil {
			return nil, errors.Internal().WithMessage("could not convert 'data' to hex bytes")
		}
		data = hexBytes.Bytes()
	} else {
		data = []byte{}
	}
	dataToEncode := []interface{}{
		&nonceBytes,
		gasPrice,
		gas.Bytes(),
		toBytes,
		value,
		data,
		tx.Signature.V.BigInt(),
		tx.Signature.R.BigInt(),
		tx.Signature.S.BigInt(),
	}

	rlpEncode, err := rlp.Encode(dataToEncode)
	if err != nil {
		return nil, errors.InternalFromErr(err).WithMessage("failed to RLP encode the payload to sign")
	}

	return entities.NewHexBytes(rlpEncode), nil
}

// Hash calculates the Ethereum transaction hash.
func (tx EthereumTransaction) Hash() (*entities.HexBytes, error) {
	nonce, err := entities.NewHexBytesFromString(hexStringEvenLength(tx.Nonce.String()))
	if err != nil {
		return nil, errors.Internal().WithMessage("could not convert 'nonce' to hex bytes")
	}
	nonceBytes := nonce.Bytes()
	if tx.Nonce.Uint64() == 0 {
		nonceBytes = []byte{}
	}

	var gasPrice *big.Int
	if tx.GasPrice.BigInt().Sign() != 0 {
		gasPrice = tx.GasPrice.BigInt()
	}

	gas, err := entities.NewHexBytesFromString(hexStringEvenLength(tx.Gas.String()))
	if err != nil {
		return nil, errors.Internal().WithMessage("could not convert 'gas' to hex bytes")
	}

	var toBytes []byte
	if tx.To != nil {
		hexBytes, newErr := entities.NewHexBytesFromString(tx.To.String())
		if newErr != nil {
			return nil, errors.Internal().WithMessage("could not convert 'to' to hex bytes")
		}
		toBytes = hexBytes.Bytes()
	} else {
		toBytes = []byte{}
	}

	var value *big.Int
	if tx.Value != nil && tx.Value.BigInt().Sign() != 0 {
		value = tx.Value.BigInt()
	}

	var data []byte
	if len(tx.Data.Bytes()) > 0 {
		hexBytes, newErr := entities.NewHexBytesFromString(tx.Data.String())
		if newErr != nil {
			return nil, errors.Internal().WithMessage("could not convert 'data' to hex bytes")
		}
		data = hexBytes.Bytes()
	} else {
		data = []byte{}
	}
	chainID, err := entities.NewHexBytesFromString(hexStringEvenLength(tx.ChainID.String()))
	if err != nil {
		return nil, errors.InternalFromErr(err).WithMessage("failed to calculate the HexBytes from chain ID")
	}

	dataToEncode := []interface{}{
		&nonceBytes,
		gasPrice,
		gas.Bytes(),
		toBytes,
		value,
		data,
		chainID.Bytes(),
		uint(0),
		uint(0),
	}
	// 1. RLP encode of the data
	rlpEncode, err := rlp.Encode(dataToEncode)
	if err != nil {
		return nil, errors.InternalFromErr(err).WithMessage("failed to RLP encode the payload to sign")
	}

	// 2. Keccak256 of the RLP encoded data
	hash, err := hashKeccak256(rlpEncode)
	if err != nil {
		return nil, errors.InternalFromErr(err).WithMessage("failed to calculate the Keccak256 of the payload to sign")
	}

	return entities.NewHexBytes(hash), nil
}

func hexStringEvenLength(input string) string {
	result := input
	if len(input)%2 != 0 {
		result = "0x0" + input[2:]
	}
	return result
}
