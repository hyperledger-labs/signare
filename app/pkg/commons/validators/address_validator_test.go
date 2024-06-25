package validators

import (
	"testing"

	"github.com/hyperledger-labs/signare/app/pkg/entities/address"

	"github.com/asaskevich/govalidator"
	"github.com/stretchr/testify/assert"
)

type RequiredAddressTestType struct {
	Address any `valid:"address"`
}

func TestRequiredAddressValidator(t *testing.T) {
	setAddressValidator()

	tests := []struct {
		name string
		have RequiredAddressTestType
		want bool
	}{
		{
			name: "valid address",
			have: RequiredAddressTestType{
				Address: address.MustNewFromHexString("0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaed"),
			},
			want: true,
		},
		{
			name: "valid slice of addresses",
			have: RequiredAddressTestType{
				Address: []address.Address{
					address.MustNewFromHexString("0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaed"),
					address.MustNewFromHexString("0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaed"),
					address.MustNewFromHexString("0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaed"),
				},
			},
			want: true,
		},
		{
			name: "empty address",
			have: RequiredAddressTestType{
				Address: address.ZeroAddress,
			},
			want: false,
		},
		{
			name: "invalid address type",
			have: RequiredAddressTestType{
				// This is a string, not an address.Address
				Address: "0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaed",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ok, _ := govalidator.ValidateStruct(tt.have)
			assert.Equal(t, tt.want, ok)
		})
	}
}
