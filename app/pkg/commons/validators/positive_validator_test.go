package validators

import (
	"testing"

	"github.com/asaskevich/govalidator"
	"github.com/stretchr/testify/assert"
)

type PositiveTestType struct {
	Amount int `valid:"required,positive"`
}

func TestPositiveValidator(t *testing.T) {
	setPositiveValidator()

	tests := []struct {
		name string
		have PositiveTestType
		want bool
	}{
		{
			name: "empty",
			have: PositiveTestType{},
			want: false,
		},
		{
			name: "negative amount",
			have: PositiveTestType{Amount: -10},
			want: false,
		},
		{
			name: "zero amount",
			have: PositiveTestType{Amount: 0},
			want: false,
		},
		{
			name: "positive amount",
			have: PositiveTestType{Amount: 10},
			want: true,
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
