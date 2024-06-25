package signaturemanager

import (
	"golang.org/x/crypto/sha3"

	"github.com/hyperledger-labs/signare/app/pkg/entities/address"
)

func DeriveAddressFromPublicKey(publicKeyBytes []byte) (*address.Address, error) {
	keccak, err := hashKeccak256(publicKeyBytes[1:])
	if err != nil {
		return nil, err
	}
	addr, err := address.NewFromRawBytes(keccak[12:])
	if err != nil {
		return nil, err
	}
	return addr, nil
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
