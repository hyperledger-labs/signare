package hsmconnector

import (
	"fmt"
	"math/big"

	"github.com/hyperledger-labs/signare/app/pkg/entities"
	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"

	curves "github.com/btcsuite/btcd/btcec/v2"
	"golang.org/x/crypto/sha3"
)

func generateEthereumTransactionSignature(signature []byte, chainID entities.HexInt256) *EthereumTransactionSignature {
	r := new(big.Int).SetBytes(signature[1:33])
	s := new(big.Int).SetBytes(signature[33:signatureLength])
	v := new(big.Int).SetBytes(signature[0:1])

	if chainID.Int.Sign() != 0 {
		ethV := int64(signature[0]) - minSignatureOffsetBitcoin // since we used the bitcoin library, the V value is 27 or 28. However, Ethereum expects either a 0 or a 1, so we substract 27.
		// calculate the V value based on https://github.com/ethereum/EIPs/blob/master/EIPS/eip-155.md#specification
		v = big.NewInt(ethV + 35)
		mul := new(big.Int).Mul(chainID.BigInt(), big.NewInt(2))
		v.Add(v, mul)
	}

	return &EthereumTransactionSignature{
		V: entities.Int256{
			Int: *v,
		},
		R: entities.Int256{
			Int: *r,
		},
		S: entities.Int256{
			Int: *s,
		},
	}
}

// signatureToLowS ensures that the signature has a low S value as Ethereum requires in EIP-2 https://github.com/ethereum/EIPs/blob/master/EIPS/eip-2.md.
func signatureToLowS(sig []byte) []byte {
	rVal := new(big.Int)
	sVal := new(big.Int)
	rVal.SetBytes(sig[0 : len(sig)/2])
	sVal.SetBytes(sig[len(sig)/2:])
	if sVal.Cmp(halfOrder()) == 1 {
		sVal.Sub(curves.S256().Params().N, sVal)
	}
	rBytes := rVal.Bytes()
	sBytes := sVal.Bytes()
	ret := make([]byte, len(sig))
	rOffset := len(sig)/2 - len(rBytes)
	sOffset := len(sig)/2 - len(sBytes)
	copy(ret[rOffset:len(sig)/2], rBytes)
	copy(ret[len(sig)/2+sOffset:], sBytes)
	return ret
}

// unmarshalECDSAKey converts bytes to a secp256k1 public key.
func unmarshalECDSAKey(pubKeyBytes []byte) (*curves.PublicKey, error) {
	pk, err := curves.ParsePubKey(pubKeyBytes)
	if err != nil {
		return nil, errors.Internal().WithMessage(fmt.Sprintf("unable to parse public key. Error: %v", err))
	}

	return pk, nil
}

// halfOrder returns half the order of the secp256k1 curve.
func halfOrder() *big.Int {
	return new(big.Int).Rsh(curves.S256().Params().N, 1)
}

// hashKeccak256 returns the keccak256 hash of the input data
func hashKeccak256(data []byte) ([]byte, error) {
	d := sha3.NewLegacyKeccak256()
	for i := range data {
		_, err := d.Write(data[i : i+1])
		if err != nil {
			return nil, err
		}
	}
	return d.Sum(nil), nil
}
