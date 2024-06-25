// Package address defines the custom Ethereum address data type of the signare.
package address

import (
	"encoding/hex"
	"encoding/json"
	"strings"

	"golang.org/x/crypto/sha3"

	signererrors "github.com/hyperledger-labs/signare/app/pkg/internal/errors"
)

const addressLength = 20

// Address represents an ethereum account address
type Address [addressLength]byte

// ZeroAddress is an empty address
var ZeroAddress = Address{}

// NewFromHexString creates a new Address from the provided hex string.
// If an error is returned, the Address will have a ZeroAddress value.
func NewFromHexString(address string) (Address, error) {
	a := Address{}
	isValidAddress := isHexAddress(address)
	if !isValidAddress {
		return ZeroAddress, signererrors.InvalidArgument().SetHumanReadableMessage("invalid address: '%s'", address)
	}
	bytesAddress, err := hex.DecodeString(removeHexPrefix(address))
	if err != nil {
		return ZeroAddress, err
	}

	copy(a[:], bytesAddress)
	return a, nil
}

// MustNewFromHexString returns an Address.
// If the address input is invalid, a ZeroAddress will be returned. Suitable for tests.
func MustNewFromHexString(address string) Address {
	a, _ := NewFromHexString(address)
	return a
}

// NewFromRawBytes creates a new address from raw unencoded bytes. If b is larger than len(a), b will be cropped from the left
func NewFromRawBytes(b []byte) (*Address, error) {
	var a [addressLength]byte
	if len(b) > len(a) {
		b = b[len(b)-addressLength:]
	} else if len(b) < len(a) {
		return nil, signererrors.InvalidArgument().WithMessage("not an eth address")
	}
	copy(a[addressLength-len(b):], b)
	address := Address{}
	copy(address[:], a[:])

	return &address, nil
}

// IsEmpty checks if the Address is empty
func (a Address) IsEmpty() bool {
	return a == ZeroAddress
}

// String returns an EIP55-compliant hex encoded string representation of the Address.
// String implements fmt.Stringer
func (a Address) String() string {
	buf := a.prefixedHex()
	sha := sha3.NewLegacyKeccak256()
	_, _ = sha.Write(buf[2:])
	hash := sha.Sum(nil)
	for i := 2; i < len(buf); i++ {
		hashByte := hash[(i-2)/2]
		if i%2 == 0 {
			hashByte = hashByte >> 4
		} else {
			hashByte &= 0xf
		}
		if buf[i] > '9' && hashByte > 7 {
			buf[i] -= 32
		}
	}
	return string(buf[:])
}

// MarshalJSON implements json.Marshaller
func (a Address) MarshalJSON() ([]byte, error) {
	ethAddressStr := a.String()
	data, err := json.Marshal(ethAddressStr)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// UnmarshalJSON implements json.Unmarshaler
func (a *Address) UnmarshalJSON(data []byte) error {
	var addrString string
	err := json.Unmarshal(data, &addrString)
	if err != nil {
		return err
	}

	addr, err := NewFromHexString(addrString)
	if err != nil {
		return err
	}

	*a = addr
	return nil
}

// prefixedHex returns a prefixed hex encoded string representation of the address
func (a Address) prefixedHex() []byte {
	var buf [addressLength*2 + 2]byte
	copy(buf[:2], "0x")
	copy(buf[2:], hex.EncodeToString(a[:]))
	return buf[:]
}

// isHexAddress verifies whether a string can represent a valid hex-encoded Address.
func isHexAddress(s string) bool {
	if isHexPrefixed(s) {
		s = s[2:]
	}
	return len(s) == 2*addressLength && isHex(s)
}

// isHexPrefixed validates if a str begins with '0x' or '0X'.
func isHexPrefixed(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

// isHex validates whether each byte is valid hexadecimal string.
func isHex(str string) bool {
	if len(str)%2 != 0 {
		return false
	}
	for _, c := range []byte(str) {
		if !isHexCharacter(c) {
			return false
		}
	}
	return true
}

// isHexCharacter returns bool of c being a valid hexadecimal.
func isHexCharacter(c byte) bool {
	return ('0' <= c && c <= '9') || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
}

// removeHexPrefix removes '0x' of the string if it is present
func removeHexPrefix(str string) string {
	if str[:2] == "0X" {
		return strings.TrimPrefix(str, "0X")
	}
	return strings.TrimPrefix(str, "0x")
}
