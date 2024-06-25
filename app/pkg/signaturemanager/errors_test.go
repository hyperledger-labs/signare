package signaturemanager_test

import (
	"testing"

	"github.com/hyperledger-labs/signare/app/pkg/signaturemanager"

	"github.com/stretchr/testify/assert"
)

func TestError_IsError(t *testing.T) {
	var err error
	err = signaturemanager.NewLibFailedError()
	assert.True(t, signaturemanager.IsLibFailedFailedError(err))

	err = signaturemanager.NewInvalidSlotError()
	assert.True(t, signaturemanager.IsInvalidSlotError(err))

	err = signaturemanager.NewKeyGenerationError()
	assert.True(t, signaturemanager.IsKeyGenerationError(err))

	err = signaturemanager.NewInternalError()
	assert.True(t, signaturemanager.IsInternalError(err))

	err = signaturemanager.NewNotFoundError()
	assert.True(t, signaturemanager.IsNotFoundError(err))

	err = signaturemanager.NewInvalidArgumentError()
	assert.True(t, signaturemanager.IsInvalidArgumentError(err))
}

func TestError_Description(t *testing.T) {
	t.Run("empty description", func(t *testing.T) {
		expectedError := "key generation failed"

		err := signaturemanager.NewKeyGenerationError()
		assert.Contains(t, err.Error(), expectedError)
	})

	t.Run("with description", func(t *testing.T) {
		expectedDescription := "test description"
		expectedError := "key generation failed: " + expectedDescription

		err := signaturemanager.NewKeyGenerationError().WithMessage(expectedDescription)
		assert.Contains(t, err.Error(), expectedError)
	})

}
