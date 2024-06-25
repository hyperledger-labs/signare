// Package utils provides a set of utility methods that can be used across the application.
package utils

import (
	"strings"

	"github.com/hyperledger-labs/signare/app/pkg/entities"
)

// DefaultIntValue returns a default integer value if the real value equals zero
func DefaultIntValue(realValue, defaultValue int) int {
	if realValue == 0 {
		return defaultValue
	}
	return realValue
}

// DefaultBigIntString returns an empty string if ptr is nil and the string value of ptr if not
func DefaultBigIntString(ptr *entities.Int256) string {
	if ptr == nil {
		return ""
	}
	return ptr.String()
}

// DefaultUIntDecString returns an empty string if ptr is nil and the base 10 string representation if not
func DefaultUIntDecString(ptr *entities.UInt64) string {
	if ptr == nil {
		return ""
	}
	return ptr.String()
}

// MaxOf returns the maximum of firstValue and secondValue
func MaxOf(firstValue int, secondValue int) int {
	if firstValue > secondValue {
		return firstValue
	}
	return secondValue
}

// MaxValue returns maxValue if realValue is greater than maxValue and realValue if it is lower
func MaxValue(realValue int, maxValue int) int {
	if realValue > maxValue {
		return maxValue
	}
	return realValue
}

// PtrStringDefaultValue returns the default value if ptr is nil and the string value if not
func PtrStringDefaultValue(ptr *string, defaultValue string) string {
	if ptr != nil {
		return *ptr
	}
	return defaultValue
}

// PtrIntDefaultValue returns the default value if ptr is nil and the int value if not
func PtrIntDefaultValue(ptr *int, defaultValue int) int {
	if ptr != nil {
		return *ptr
	}
	return defaultValue
}

// DefaultString returns src if it is not empty after trim (TrimSpace) or defaultValue
func DefaultString(src string, defaultValue string) string {
	trimmedStr := strings.TrimSpace(src)
	if len(trimmedStr) == 0 {
		return defaultValue
	}
	return src
}
