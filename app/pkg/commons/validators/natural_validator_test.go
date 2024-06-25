package validators

import (
	"testing"

	"github.com/asaskevich/govalidator"
	"github.com/stretchr/testify/assert"
)

type NaturalTestType struct {
	Amount int `valid:"required,natural"`
}

func TestNaturalValidator(t *testing.T) {
	setNaturalValidator()

	tests := []struct {
		name string
		have NaturalTestType
		want bool
	}{
		{
			name: "empty",
			have: NaturalTestType{},
			want: false,
		},
		{
			name: "negative amount",
			have: NaturalTestType{Amount: -10},
			want: false,
		},
		{
			name: "zero amount",
			have: NaturalTestType{Amount: 0},
			want: false,
		},
		{
			name: "positive amount",
			have: NaturalTestType{Amount: 10},
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
