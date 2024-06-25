package utils_test

import (
	"math/big"
	"testing"

	"github.com/hyperledger-labs/signare/app/pkg/utils"

	"github.com/stretchr/testify/assert"
)

func Test_DefaultIntValue(t *testing.T) {
	tests := []struct {
		name         string
		have         int
		defaultValue int
		expected     int
	}{
		{
			name:         "Zero case",
			have:         0,
			defaultValue: 1,
			expected:     1,
		}, {
			name:         "Normal Case",
			have:         10,
			defaultValue: 1,
			expected:     10,
		},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			result := utils.DefaultIntValue(tt.have, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_DefaultMaxValue(t *testing.T) {
	tests := []struct {
		name     string
		have     int
		maxValue int
		expected int
	}{
		{
			name:     "Inferior case",
			have:     0,
			maxValue: 1,
			expected: 0,
		}, {
			name:     "Equal Case",
			have:     1,
			maxValue: 1,
			expected: 1,
		}, {
			name:     "Superior Case",
			have:     2,
			maxValue: 1,
			expected: 1,
		},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			result := defaultIntMaxValue(tt.have, tt.maxValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_DefaultPtrBigIntValue(t *testing.T) {
	tests := []struct {
		name         string
		have         *big.Int
		defaultValue *big.Int
		expected     *big.Int
	}{
		{
			name:         "Zero case",
			have:         big.NewInt(0),
			defaultValue: big.NewInt(1),
			expected:     big.NewInt(0),
		},
		{
			name:         "Nil case",
			have:         nil,
			defaultValue: big.NewInt(1),
			expected:     big.NewInt(1),
		},
		{
			name:         "Normal Case",
			have:         big.NewInt(10),
			defaultValue: big.NewInt(1),
			expected:     big.NewInt(10),
		},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			result := defaultPtrBigIntValue(tt.have, tt.defaultValue)
			assert.Equal(t, tt.expected.String(), result.String())
		})
	}
}

func Test_PtrStringDefaultValue(t *testing.T) {
	defaultValue := "defaultValue"
	testValue1 := "testValue1"

	tests := []struct {
		name     string
		value    *string
		expected string
	}{
		{
			name:     "Nil value",
			value:    nil,
			expected: defaultValue,
		}, {
			name:     "Test value 1",
			value:    &testValue1,
			expected: testValue1,
		},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			result := utils.PtrStringDefaultValue(tt.value, defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_PtrIntDefaultValue(t *testing.T) {
	defaultValue := 123
	testValue1 := 456

	tests := []struct {
		name     string
		value    *int
		expected int
	}{
		{
			name:     "Nil value",
			value:    nil,
			expected: defaultValue,
		}, {
			name:     "Test value",
			value:    &testValue1,
			expected: testValue1,
		},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			result := utils.PtrIntDefaultValue(tt.value, defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_DefaultString(t *testing.T) {
	defaultValue := "defaultValue"

	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{
			name:     "Equal",
			value:    "testValue-1",
			expected: "testValue-1",
		}, {
			name:     "Equal",
			value:    " testValue-2",
			expected: " testValue-2",
		}, {
			name:     "Test value",
			value:    " ",
			expected: defaultValue,
		},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			result := utils.DefaultString(tt.value, defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMaxOf(t *testing.T) {
	tests := []struct {
		name           string
		firstValue     int
		secondValue    int
		expectedResult int
	}{
		{
			name:           "both values equal",
			firstValue:     1,
			secondValue:    1,
			expectedResult: 1,
		}, {
			name:           "first greater",
			firstValue:     10,
			secondValue:    5,
			expectedResult: 10,
		}, {
			name:           "second greater",
			firstValue:     6,
			secondValue:    18,
			expectedResult: 18,
		}, {
			name:           "negative values",
			firstValue:     -1,
			secondValue:    -10,
			expectedResult: -1,
		},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			result := utils.MaxOf(tt.firstValue, tt.secondValue)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func defaultIntMaxValue(realValue int, maxValue int) int {
	if realValue > maxValue {
		return maxValue
	}
	return realValue
}

func defaultPtrBigIntValue(realValue *big.Int, defaultValue *big.Int) *big.Int {
	if realValue == nil {
		return defaultValue
	}
	return realValue
}
