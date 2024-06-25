package address_test

import (
	"testing"

	"github.com/hyperledger-labs/signare/app/pkg/entities/address"

	"github.com/stretchr/testify/assert"
)

var (
	validAddrString = "970e8128ab834e8eac17ab8e3812f010678cf791"

	expectedAddressEIP55 = "0x970E8128AB834E8EAC17Ab8E3812F010678CF791"
	zeroAddressString    = "0x0000000000000000000000000000000000000000"
)

func TestNewFromHexString(t *testing.T) {
	var err error
	var got address.Address

	t.Run("invalid addresses", func(t *testing.T) {
		got, err = address.NewFromHexString("invalid address format")
		assert.Error(t, err)
		assert.Equal(t, address.ZeroAddress, got)
		assert.Equal(t, zeroAddressString, got.String())

		// valid length, invalid hex characters
		got, err = address.NewFromHexString("0xz70z8128zb834e8enc17az8e3812k010678zf791")
		assert.Error(t, err)
		assert.Equal(t, address.ZeroAddress, got)
		assert.Equal(t, zeroAddressString, got.String())
		got, err = address.NewFromHexString("z70z8128zb834e8enc17az8e3812k010678zf791")
		assert.Error(t, err)
		assert.Equal(t, address.ZeroAddress, got)
		assert.Equal(t, zeroAddressString, got.String())

		got, err = address.NewFromHexString("")
		assert.Error(t, err)
		assert.Equal(t, address.ZeroAddress, got)
		assert.Equal(t, zeroAddressString, got.String())
	})

	t.Run("valid addresses", func(t *testing.T) {
		got, err = address.NewFromHexString(validAddrString)
		assert.NoError(t, err)
		assert.Equal(t, expectedAddressEIP55, got.String())

		got, err = address.NewFromHexString("0x" + validAddrString)
		assert.NoError(t, err)
		assert.Equal(t, expectedAddressEIP55, got.String())

		got, err = address.NewFromHexString("0X" + validAddrString)
		assert.NoError(t, err)
		assert.Equal(t, expectedAddressEIP55, got.String())
	})
}
