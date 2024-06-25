package errors_test

import (
	"fmt"
	"testing"

	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"

	"github.com/stretchr/testify/assert"
)

var errExternal = fmt.Errorf("third party error")

func TestAlreadyExists(t *testing.T) {
	t.Run("AlreadyExists error", func(t *testing.T) {
		err := errors.AlreadyExists()
		assert.True(t, errors.IsAlreadyExists(err))
		assert.Equal(t, errors.ErrAlreadyExists, err.Type())
		assert.Equal(t, string(errors.ErrAlreadyExists), err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})

	t.Run("AlreadyExists error ", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s: operation with AdminID %v failed", errors.ErrAlreadyExists, 10)

		err := errors.AlreadyExists().WithMessage("operation with AdminID %v failed", 10)
		assert.True(t, errors.IsAlreadyExists(err))
		assert.Equal(t, errors.ErrAlreadyExists, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})

	t.Run("AlreadyExists error from external error", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s (wrapped error: third party error)", errors.ErrAlreadyExists)

		err := errors.AlreadyExistsFromErr(errExternal)
		assert.True(t, errors.IsAlreadyExists(err))
		assert.Equal(t, errors.ErrAlreadyExists, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})

	t.Run("AlreadyExists error with error message from external error", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s: unexpected error (wrapped error: third party error)", errors.ErrAlreadyExists)

		err := errors.AlreadyExistsFromErr(errExternal).WithMessage("unexpected error")
		assert.True(t, errors.IsAlreadyExists(err))
		assert.Equal(t, errors.ErrAlreadyExists, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})

	t.Run("AlreadyExists error with error message and human readable message", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s: operation with AdminID %v failed", errors.ErrAlreadyExists, 10)
		expectedHumanReadableMessage := "Please try again later"

		err := errors.AlreadyExists().WithMessage("operation with AdminID %v failed", 10).SetHumanReadableMessage("Please try again later")
		assert.True(t, errors.IsAlreadyExists(err))
		assert.Equal(t, errors.ErrAlreadyExists, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Equal(t, &expectedHumanReadableMessage, err.HumanReadableMessage())
	})

	t.Run("AlreadyExists error with error message and human readable message from external error", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s: operation with AdminID %v failed (wrapped error: third party error)", errors.ErrAlreadyExists, 10)
		expectedHumanReadableMessage := "Please try again later"

		err := errors.AlreadyExistsFromErr(errExternal).WithMessage("operation with AdminID %v failed", 10).SetHumanReadableMessage("Please try again later")
		assert.True(t, errors.IsAlreadyExists(err))
		assert.Equal(t, errors.ErrAlreadyExists, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Equal(t, &expectedHumanReadableMessage, err.HumanReadableMessage())
	})

	t.Run("AlreadyExists error wrapping another signare error", func(t *testing.T) {
		errNotFound := errors.NotFound().WithMessage("no records found")
		expectedErrMsg := fmt.Sprintf("%s (wrapped error: %s: no records found)", errors.ErrAlreadyExists, errors.ErrNotFound)

		err := errors.AlreadyExistsFromErr(errNotFound)
		assert.False(t, errors.IsNotFound(err))
		assert.True(t, errors.IsAlreadyExists(err))
		assert.Equal(t, errors.ErrAlreadyExists, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})
	assert.False(t, errors.IsAlreadyExists(errExternal))
}

func TestBadGateway(t *testing.T) {
	t.Run("BadGateway error", func(t *testing.T) {
		err := errors.BadGateway()
		assert.True(t, errors.IsBadGateway(err))
		assert.Equal(t, errors.ErrBadGateway, err.Type())
		assert.Equal(t, string(errors.ErrBadGateway), err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})

	t.Run("BadGateway error ", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s: operation with AdminID %v failed", errors.ErrBadGateway, 10)

		err := errors.BadGateway().WithMessage("operation with AdminID %v failed", 10)
		assert.True(t, errors.IsBadGateway(err))
		assert.Equal(t, errors.ErrBadGateway, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})

	t.Run("BadGateway error from external error", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s (wrapped error: third party error)", errors.ErrBadGateway)

		err := errors.BadGatewayFromErr(errExternal)
		assert.True(t, errors.IsBadGateway(err))
		assert.Equal(t, errors.ErrBadGateway, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})

	t.Run("BadGateway error with error message from external error", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s: unexpected error (wrapped error: third party error)", errors.ErrBadGateway)

		err := errors.BadGatewayFromErr(errExternal).WithMessage("unexpected error")
		assert.True(t, errors.IsBadGateway(err))
		assert.Equal(t, errors.ErrBadGateway, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})

	t.Run("BadGateway error with error message and human readable message", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s: operation with AdminID %v failed", errors.ErrBadGateway, 10)
		expectedHumanReadableMessage := "Please try again later"

		err := errors.BadGateway().WithMessage("operation with AdminID %v failed", 10).SetHumanReadableMessage("Please try again later")
		assert.True(t, errors.IsBadGateway(err))
		assert.Equal(t, errors.ErrBadGateway, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Equal(t, &expectedHumanReadableMessage, err.HumanReadableMessage())
	})

	t.Run("BadGateway error with error message and human readable message from external error", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s: operation with AdminID %v failed (wrapped error: third party error)", errors.ErrBadGateway, 10)
		expectedHumanReadableMessage := "Please try again later"

		err := errors.BadGatewayFromErr(errExternal).WithMessage("operation with AdminID %v failed", 10).SetHumanReadableMessage("Please try again later")
		assert.True(t, errors.IsBadGateway(err))
		assert.Equal(t, errors.ErrBadGateway, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Equal(t, &expectedHumanReadableMessage, err.HumanReadableMessage())
	})

	t.Run("BadGateway error wrapping another signare error", func(t *testing.T) {
		errNotFound := errors.NotFound().WithMessage("no records found")
		expectedErrMsg := fmt.Sprintf("%s (wrapped error: %s: no records found)", errors.ErrBadGateway, errors.ErrNotFound)

		err := errors.BadGatewayFromErr(errNotFound)
		assert.False(t, errors.IsNotFound(err))
		assert.True(t, errors.IsBadGateway(err))
		assert.Equal(t, errors.ErrBadGateway, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})
	assert.False(t, errors.IsBadGateway(errExternal))
}

func TestInternal(t *testing.T) {
	t.Run("Internal error", func(t *testing.T) {
		err := errors.Internal()
		assert.True(t, errors.IsInternal(err))
		assert.Equal(t, errors.ErrInternal, err.Type())
		assert.Equal(t, string(errors.ErrInternal), err.Error())
	})

	t.Run("Internal error ", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s: operation with AdminID %v failed", errors.ErrInternal, 10)

		err := errors.Internal().WithMessage("operation with AdminID %v failed", 10)
		assert.True(t, errors.IsInternal(err))
		assert.Equal(t, errors.ErrInternal, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
	})

	t.Run("Internal error from external error", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s (wrapped error: third party error)", errors.ErrInternal)

		err := errors.InternalFromErr(errExternal)
		assert.True(t, errors.IsInternal(err))
		assert.Equal(t, errors.ErrInternal, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
	})

	t.Run("Internal error with error message from external error", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s: unexpected error (wrapped error: third party error)", errors.ErrInternal)

		err := errors.InternalFromErr(errExternal).WithMessage("unexpected error")
		assert.True(t, errors.IsInternal(err))
		assert.Equal(t, errors.ErrInternal, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
	})

	t.Run("Internal error wrapping another signare error", func(t *testing.T) {
		errNotFound := errors.NotFound().WithMessage("no records found")
		expectedErrMsg := fmt.Sprintf("%s (wrapped error: %s: no records found)", errors.ErrInternal, errors.ErrNotFound)

		err := errors.InternalFromErr(errNotFound)
		assert.False(t, errors.IsNotFound(err))
		assert.True(t, errors.IsInternal(err))
		assert.Equal(t, errors.ErrInternal, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
	})
	assert.False(t, errors.IsInternal(errExternal))
}

func TestInvalidArgument(t *testing.T) {
	t.Run("InvalidArgument error", func(t *testing.T) {
		err := errors.InvalidArgument()
		assert.True(t, errors.IsInvalidArgument(err))
		assert.Equal(t, errors.ErrInvalidArgument, err.Type())
		assert.Equal(t, string(errors.ErrInvalidArgument), err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})

	t.Run("InvalidArgument error ", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s: operation with AdminID %v failed", errors.ErrInvalidArgument, 10)

		err := errors.InvalidArgument().WithMessage("operation with AdminID %v failed", 10)
		assert.True(t, errors.IsInvalidArgument(err))
		assert.Equal(t, errors.ErrInvalidArgument, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})

	t.Run("InvalidArgument error from external error", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s (wrapped error: third party error)", errors.ErrInvalidArgument)

		err := errors.InvalidArgumentFromErr(errExternal)
		assert.True(t, errors.IsInvalidArgument(err))
		assert.Equal(t, errors.ErrInvalidArgument, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})

	t.Run("InvalidArgument error with error message from external error", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s: unexpected error (wrapped error: third party error)", errors.ErrInvalidArgument)

		err := errors.InvalidArgumentFromErr(errExternal).WithMessage("unexpected error")
		assert.True(t, errors.IsInvalidArgument(err))
		assert.Equal(t, errors.ErrInvalidArgument, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})

	t.Run("InvalidArgument error with error message and human readable message", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s: operation with AdminID %v failed", errors.ErrInvalidArgument, 10)
		expectedHumanReadableMessage := "Please try again later"

		err := errors.InvalidArgument().WithMessage("operation with AdminID %v failed", 10).SetHumanReadableMessage("Please try again later")
		assert.True(t, errors.IsInvalidArgument(err))
		assert.Equal(t, errors.ErrInvalidArgument, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Equal(t, &expectedHumanReadableMessage, err.HumanReadableMessage())
	})

	t.Run("InvalidArgument error with error message and human readable message from external error", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s: operation with AdminID %v failed (wrapped error: third party error)", errors.ErrInvalidArgument, 10)
		expectedHumanReadableMessage := "Please try again later"

		err := errors.InvalidArgumentFromErr(errExternal).WithMessage("operation with AdminID %v failed", 10).SetHumanReadableMessage("Please try again later")
		assert.True(t, errors.IsInvalidArgument(err))
		assert.Equal(t, errors.ErrInvalidArgument, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Equal(t, &expectedHumanReadableMessage, err.HumanReadableMessage())
	})

	t.Run("InvalidArgument error wrapping another signare error", func(t *testing.T) {
		errNotFound := errors.NotFound().WithMessage("no records found")
		expectedErrMsg := fmt.Sprintf("%s (wrapped error: %s: no records found)", errors.ErrInvalidArgument, errors.ErrNotFound)

		err := errors.InvalidArgumentFromErr(errNotFound)
		assert.False(t, errors.IsNotFound(err))
		assert.True(t, errors.IsInvalidArgument(err))
		assert.Equal(t, errors.ErrInvalidArgument, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})
	assert.False(t, errors.IsInvalidArgument(errExternal))
}

func TestNotFound(t *testing.T) {
	t.Run("NotFound error", func(t *testing.T) {
		err := errors.NotFound()
		assert.True(t, errors.IsNotFound(err))
		assert.Equal(t, errors.ErrNotFound, err.Type())
		assert.Equal(t, string(errors.ErrNotFound), err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})

	t.Run("NotFound error ", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s: operation with AdminID %v failed", errors.ErrNotFound, 10)

		err := errors.NotFound().WithMessage("operation with AdminID %v failed", 10)
		assert.True(t, errors.IsNotFound(err))
		assert.Equal(t, errors.ErrNotFound, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})

	t.Run("NotFound error from external error", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s (wrapped error: third party error)", errors.ErrNotFound)

		err := errors.NotFoundFromErr(errExternal)
		assert.True(t, errors.IsNotFound(err))
		assert.Equal(t, errors.ErrNotFound, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})

	t.Run("NotFound error with error message from external error", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s: unexpected error (wrapped error: third party error)", errors.ErrNotFound)

		err := errors.NotFoundFromErr(errExternal).WithMessage("unexpected error")
		assert.True(t, errors.IsNotFound(err))
		assert.Equal(t, errors.ErrNotFound, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})

	t.Run("NotFound error with error message and human readable message", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s: operation with AdminID %v failed", errors.ErrNotFound, 10)
		expectedHumanReadableMessage := "Please try again later"

		err := errors.NotFound().WithMessage("operation with AdminID %v failed", 10).SetHumanReadableMessage("Please try again later")
		assert.True(t, errors.IsNotFound(err))
		assert.Equal(t, errors.ErrNotFound, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Equal(t, &expectedHumanReadableMessage, err.HumanReadableMessage())
	})

	t.Run("NotFound error with error message and human readable message from external error", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s: operation with AdminID %v failed (wrapped error: third party error)", errors.ErrNotFound, 10)
		expectedHumanReadableMessage := "Please try again later"

		err := errors.NotFoundFromErr(errExternal).WithMessage("operation with AdminID %v failed", 10).SetHumanReadableMessage("Please try again later")
		assert.True(t, errors.IsNotFound(err))
		assert.Equal(t, errors.ErrNotFound, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Equal(t, &expectedHumanReadableMessage, err.HumanReadableMessage())
	})

	t.Run("NotFound error wrapping another signare error", func(t *testing.T) {
		errNotFound := errors.Unauthenticated().WithMessage("request not authorized")
		expectedErrMsg := fmt.Sprintf("%s (wrapped error: %s: request not authorized)", errors.ErrNotFound, errors.ErrUnauthenticated)

		err := errors.NotFoundFromErr(errNotFound)
		assert.False(t, errors.IsUnauthenticated(err))
		assert.True(t, errors.IsNotFound(err))
		assert.Equal(t, errors.ErrNotFound, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})
	assert.False(t, errors.IsNotFound(errExternal))
}

func TestNotImplemented(t *testing.T) {
	t.Run("NotImplemented error", func(t *testing.T) {
		err := errors.NotImplemented()
		assert.True(t, errors.IsNotImplemented(err))
		assert.Equal(t, errors.ErrNotImplemented, err.Type())
		assert.Equal(t, string(errors.ErrNotImplemented), err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})

	t.Run("NotImplemented error ", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s: operation with AdminID %v failed", errors.ErrNotImplemented, 10)

		err := errors.NotImplemented().WithMessage("operation with AdminID %v failed", 10)
		assert.True(t, errors.IsNotImplemented(err))
		assert.Equal(t, errors.ErrNotImplemented, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})

	t.Run("NotImplemented error from external error", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s (wrapped error: third party error)", errors.ErrNotImplemented)

		err := errors.NotImplementedFromErr(errExternal)
		assert.True(t, errors.IsNotImplemented(err))
		assert.Equal(t, errors.ErrNotImplemented, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})

	t.Run("NotImplemented error with error message from external error", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s: unexpected error (wrapped error: third party error)", errors.ErrNotImplemented)

		err := errors.NotImplementedFromErr(errExternal).WithMessage("unexpected error")
		assert.True(t, errors.IsNotImplemented(err))
		assert.Equal(t, errors.ErrNotImplemented, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})

	t.Run("NotImplemented error with error message and human readable message", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s: operation with AdminID %v failed", errors.ErrNotImplemented, 10)
		expectedHumanReadableMessage := "Please try again later"

		err := errors.NotImplemented().WithMessage("operation with AdminID %v failed", 10).SetHumanReadableMessage("Please try again later")
		assert.True(t, errors.IsNotImplemented(err))
		assert.Equal(t, errors.ErrNotImplemented, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Equal(t, &expectedHumanReadableMessage, err.HumanReadableMessage())
	})

	t.Run("NotImplemented error with error message and human readable message from external error", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s: operation with AdminID %v failed (wrapped error: third party error)", errors.ErrNotImplemented, 10)
		expectedHumanReadableMessage := "Please try again later"

		err := errors.NotImplementedFromErr(errExternal).WithMessage("operation with AdminID %v failed", 10).SetHumanReadableMessage("Please try again later")
		assert.True(t, errors.IsNotImplemented(err))
		assert.Equal(t, errors.ErrNotImplemented, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Equal(t, &expectedHumanReadableMessage, err.HumanReadableMessage())
	})

	t.Run("NotImplemented error wrapping another signare error", func(t *testing.T) {
		errNotFound := errors.NotFound().WithMessage("no records found")
		expectedErrMsg := fmt.Sprintf("%s (wrapped error: %s: no records found)", errors.ErrNotImplemented, errors.ErrNotFound)

		err := errors.NotImplementedFromErr(errNotFound)
		assert.False(t, errors.IsNotFound(err))
		assert.True(t, errors.IsNotImplemented(err))
		assert.Equal(t, errors.ErrNotImplemented, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})
	assert.False(t, errors.IsNotImplemented(errExternal))
}

func TestPermissionDenied(t *testing.T) {
	t.Run("PermissionDenied error", func(t *testing.T) {
		err := errors.PermissionDenied()
		assert.True(t, errors.IsPermissionDenied(err))
		assert.Equal(t, errors.ErrPermissionDenied, err.Type())
		assert.Equal(t, string(errors.ErrPermissionDenied), err.Error())
	})

	t.Run("PermissionDenied error ", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s: operation with AdminID %v failed", errors.ErrPermissionDenied, 10)

		err := errors.PermissionDenied().WithMessage("operation with AdminID %v failed", 10)
		assert.True(t, errors.IsPermissionDenied(err))
		assert.Equal(t, errors.ErrPermissionDenied, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
	})

	t.Run("PermissionDenied error from external error", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s (wrapped error: third party error)", errors.ErrPermissionDenied)

		err := errors.PermissionDeniedFromErr(errExternal)
		assert.True(t, errors.IsPermissionDenied(err))
		assert.Equal(t, errors.ErrPermissionDenied, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
	})

	t.Run("PermissionDenied error with error message from external error", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s: unexpected error (wrapped error: third party error)", errors.ErrPermissionDenied)

		err := errors.PermissionDeniedFromErr(errExternal).WithMessage("unexpected error")
		assert.True(t, errors.IsPermissionDenied(err))
		assert.Equal(t, errors.ErrPermissionDenied, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
	})

	t.Run("PermissionDenied error wrapping another signare error", func(t *testing.T) {
		errNotFound := errors.NotFound().WithMessage("no records found")
		expectedErrMsg := fmt.Sprintf("%s (wrapped error: %s: no records found)", errors.ErrPermissionDenied, errors.ErrNotFound)

		err := errors.PermissionDeniedFromErr(errNotFound)
		assert.False(t, errors.IsNotFound(err))
		assert.True(t, errors.IsPermissionDenied(err))
		assert.Equal(t, errors.ErrPermissionDenied, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
	})
	assert.False(t, errors.IsPermissionDenied(errExternal))
}

func TestPreconditionFailed(t *testing.T) {
	t.Run("PreconditionFailed error", func(t *testing.T) {
		err := errors.PreconditionFailed()
		assert.True(t, errors.IsPreconditionFailed(err))
		assert.Equal(t, errors.ErrPreconditionFailed, err.Type())
		assert.Equal(t, string(errors.ErrPreconditionFailed), err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})

	t.Run("PreconditionFailed error ", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s: operation with AdminID %v failed", errors.ErrPreconditionFailed, 10)

		err := errors.PreconditionFailed().WithMessage("operation with AdminID %v failed", 10)
		assert.True(t, errors.IsPreconditionFailed(err))
		assert.Equal(t, errors.ErrPreconditionFailed, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})

	t.Run("PreconditionFailed error from external error", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s (wrapped error: third party error)", errors.ErrPreconditionFailed)

		err := errors.PreconditionFailedFromErr(errExternal)
		assert.True(t, errors.IsPreconditionFailed(err))
		assert.Equal(t, errors.ErrPreconditionFailed, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})

	t.Run("PreconditionFailed error with error message from external error", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s: unexpected error (wrapped error: third party error)", errors.ErrPreconditionFailed)

		err := errors.PreconditionFailedFromErr(errExternal).WithMessage("unexpected error")
		assert.True(t, errors.IsPreconditionFailed(err))
		assert.Equal(t, errors.ErrPreconditionFailed, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})

	t.Run("PreconditionFailed error with error message and human readable message", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s: operation with AdminID %v failed", errors.ErrPreconditionFailed, 10)
		expectedHumanReadableMessage := "Please try again later"

		err := errors.PreconditionFailed().WithMessage("operation with AdminID %v failed", 10).SetHumanReadableMessage("Please try again later")
		assert.True(t, errors.IsPreconditionFailed(err))
		assert.Equal(t, errors.ErrPreconditionFailed, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Equal(t, &expectedHumanReadableMessage, err.HumanReadableMessage())
	})

	t.Run("PreconditionFailed error with error message and human readable message from external error", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s: operation with AdminID %v failed (wrapped error: third party error)", errors.ErrPreconditionFailed, 10)
		expectedHumanReadableMessage := "Please try again later"

		err := errors.PreconditionFailedFromErr(errExternal).WithMessage("operation with AdminID %v failed", 10).SetHumanReadableMessage("Please try again later")
		assert.True(t, errors.IsPreconditionFailed(err))
		assert.Equal(t, errors.ErrPreconditionFailed, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Equal(t, &expectedHumanReadableMessage, err.HumanReadableMessage())
	})

	t.Run("PreconditionFailed error wrapping another signare error", func(t *testing.T) {
		errNotFound := errors.NotFound().WithMessage("no records found")
		expectedErrMsg := fmt.Sprintf("%s (wrapped error: %s: no records found)", errors.ErrPreconditionFailed, errors.ErrNotFound)

		err := errors.PreconditionFailedFromErr(errNotFound)
		assert.False(t, errors.IsNotFound(err))
		assert.True(t, errors.IsPreconditionFailed(err))
		assert.Equal(t, errors.ErrPreconditionFailed, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})
	assert.False(t, errors.IsPreconditionFailed(errExternal))
}

func TestTimeout(t *testing.T) {
	t.Run("Timeout error", func(t *testing.T) {
		err := errors.Timeout()
		assert.True(t, errors.IsTimeout(err))
		assert.Equal(t, errors.ErrTimeout, err.Type())
		assert.Equal(t, string(errors.ErrTimeout), err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})

	t.Run("Timeout error ", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s: operation with AdminID %v failed", errors.ErrTimeout, 10)

		err := errors.Timeout().WithMessage("operation with AdminID %v failed", 10)
		assert.True(t, errors.IsTimeout(err))
		assert.Equal(t, errors.ErrTimeout, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})

	t.Run("Timeout error from external error", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s (wrapped error: third party error)", errors.ErrTimeout)

		err := errors.TimeoutFromErr(errExternal)
		assert.True(t, errors.IsTimeout(err))
		assert.Equal(t, errors.ErrTimeout, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})

	t.Run("Timeout error with error message from external error", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s: unexpected error (wrapped error: third party error)", errors.ErrTimeout)

		err := errors.TimeoutFromErr(errExternal).WithMessage("unexpected error")
		assert.True(t, errors.IsTimeout(err))
		assert.Equal(t, errors.ErrTimeout, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})

	t.Run("Timeout error with error message and human readable message", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s: operation with AdminID %v failed", errors.ErrTimeout, 10)
		expectedHumanReadableMessage := "Please try again later"

		err := errors.Timeout().WithMessage("operation with AdminID %v failed", 10).SetHumanReadableMessage("Please try again later")
		assert.True(t, errors.IsTimeout(err))
		assert.Equal(t, errors.ErrTimeout, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Equal(t, &expectedHumanReadableMessage, err.HumanReadableMessage())
	})

	t.Run("Timeout error with error message and human readable message from external error", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s: operation with AdminID %v failed (wrapped error: third party error)", errors.ErrTimeout, 10)
		expectedHumanReadableMessage := "Please try again later"

		err := errors.TimeoutFromErr(errExternal).WithMessage("operation with AdminID %v failed", 10).SetHumanReadableMessage("Please try again later")
		assert.True(t, errors.IsTimeout(err))
		assert.Equal(t, errors.ErrTimeout, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Equal(t, &expectedHumanReadableMessage, err.HumanReadableMessage())
	})

	t.Run("Timeout error wrapping another signare error", func(t *testing.T) {
		errNotFound := errors.NotFound().WithMessage("no records found")
		expectedErrMsg := fmt.Sprintf("%s (wrapped error: %s: no records found)", errors.ErrTimeout, errors.ErrNotFound)

		err := errors.TimeoutFromErr(errNotFound)
		assert.False(t, errors.IsNotFound(err))
		assert.True(t, errors.IsTimeout(err))
		assert.Equal(t, errors.ErrTimeout, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})
	assert.False(t, errors.IsTimeout(errExternal))
}

func TestTooManyReq(t *testing.T) {
	t.Run("TooManyReq error", func(t *testing.T) {
		err := errors.TooManyReq()
		assert.True(t, errors.IsTooManyReq(err))
		assert.Equal(t, errors.ErrTooManyReq, err.Type())
		assert.Equal(t, string(errors.ErrTooManyReq), err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})

	t.Run("TooManyReq error ", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s: operation with AdminID %v failed", errors.ErrTooManyReq, 10)

		err := errors.TooManyReq().WithMessage("operation with AdminID %v failed", 10)
		assert.True(t, errors.IsTooManyReq(err))
		assert.Equal(t, errors.ErrTooManyReq, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})

	t.Run("TooManyReq error from external error", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s (wrapped error: third party error)", errors.ErrTooManyReq)

		err := errors.TooManyReqFromErr(errExternal)
		assert.True(t, errors.IsTooManyReq(err))
		assert.Equal(t, errors.ErrTooManyReq, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})

	t.Run("TooManyReq error with error message from external error", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s: unexpected error (wrapped error: third party error)", errors.ErrTooManyReq)

		err := errors.TooManyReqFromErr(errExternal).WithMessage("unexpected error")
		assert.True(t, errors.IsTooManyReq(err))
		assert.Equal(t, errors.ErrTooManyReq, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})

	t.Run("TooManyReq error with error message and human readable message", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s: operation with AdminID %v failed", errors.ErrTooManyReq, 10)
		expectedHumanReadableMessage := "Please try again later"

		err := errors.TooManyReq().WithMessage("operation with AdminID %v failed", 10).SetHumanReadableMessage("Please try again later")
		assert.True(t, errors.IsTooManyReq(err))
		assert.Equal(t, errors.ErrTooManyReq, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Equal(t, &expectedHumanReadableMessage, err.HumanReadableMessage())
	})

	t.Run("TooManyReq error with error message and human readable message from external error", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s: operation with AdminID %v failed (wrapped error: third party error)", errors.ErrTooManyReq, 10)
		expectedHumanReadableMessage := "Please try again later"

		err := errors.TooManyReqFromErr(errExternal).WithMessage("operation with AdminID %v failed", 10).SetHumanReadableMessage("Please try again later")
		assert.True(t, errors.IsTooManyReq(err))
		assert.Equal(t, errors.ErrTooManyReq, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Equal(t, &expectedHumanReadableMessage, err.HumanReadableMessage())
	})

	t.Run("TooManyReq error wrapping another signare error", func(t *testing.T) {
		errNotFound := errors.NotFound().WithMessage("no records found")
		expectedErrMsg := fmt.Sprintf("%s (wrapped error: %s: no records found)", errors.ErrTooManyReq, errors.ErrNotFound)

		err := errors.TooManyReqFromErr(errNotFound)
		assert.False(t, errors.IsNotFound(err))
		assert.True(t, errors.IsTooManyReq(err))
		assert.Equal(t, errors.ErrTooManyReq, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})
	assert.False(t, errors.IsTooManyReq(errExternal))
}

func TestUnauthenticated(t *testing.T) {
	t.Run("Unauthenticated error", func(t *testing.T) {
		err := errors.Unauthenticated()
		assert.True(t, errors.IsUnauthenticated(err))
		assert.Equal(t, errors.ErrUnauthenticated, err.Type())
		assert.Equal(t, string(errors.ErrUnauthenticated), err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})

	t.Run("Unauthenticated error ", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s: operation with AdminID %v failed", errors.ErrUnauthenticated, 10)

		err := errors.Unauthenticated().WithMessage("operation with AdminID %v failed", 10)
		assert.True(t, errors.IsUnauthenticated(err))
		assert.Equal(t, errors.ErrUnauthenticated, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})

	t.Run("Unauthenticated error from external error", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s (wrapped error: third party error)", errors.ErrUnauthenticated)

		err := errors.UnauthenticatedFromErr(errExternal)
		assert.True(t, errors.IsUnauthenticated(err))
		assert.Equal(t, errors.ErrUnauthenticated, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})

	t.Run("Unauthenticated error with error message from external error", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s: unexpected error (wrapped error: third party error)", errors.ErrUnauthenticated)

		err := errors.UnauthenticatedFromErr(errExternal).WithMessage("unexpected error")
		assert.True(t, errors.IsUnauthenticated(err))
		assert.Equal(t, errors.ErrUnauthenticated, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})

	t.Run("Unauthenticated error with error message and human readable message", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s: operation with AdminID %v failed", errors.ErrUnauthenticated, 10)
		expectedHumanReadableMessage := "Please try again later"

		err := errors.Unauthenticated().WithMessage("operation with AdminID %v failed", 10).SetHumanReadableMessage("Please try again later")
		assert.True(t, errors.IsUnauthenticated(err))
		assert.Equal(t, errors.ErrUnauthenticated, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Equal(t, &expectedHumanReadableMessage, err.HumanReadableMessage())
	})

	t.Run("Unauthenticated error with error message and human readable message from external error", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s: operation with AdminID %v failed (wrapped error: third party error)", errors.ErrUnauthenticated, 10)
		expectedHumanReadableMessage := "Please try again later"

		err := errors.UnauthenticatedFromErr(errExternal).WithMessage("operation with AdminID %v failed", 10).SetHumanReadableMessage("Please try again later")
		assert.True(t, errors.IsUnauthenticated(err))
		assert.Equal(t, errors.ErrUnauthenticated, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Equal(t, &expectedHumanReadableMessage, err.HumanReadableMessage())
	})

	t.Run("Unauthenticated error wrapping another signare error", func(t *testing.T) {
		errNotFound := errors.NotFound().WithMessage("no records found")
		expectedErrMsg := fmt.Sprintf("%s (wrapped error: %s: no records found)", errors.ErrUnauthenticated, errors.ErrNotFound)

		err := errors.UnauthenticatedFromErr(errNotFound)
		assert.False(t, errors.IsNotFound(err))
		assert.True(t, errors.IsUnauthenticated(err))
		assert.Equal(t, errors.ErrUnauthenticated, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})
	assert.False(t, errors.IsUnauthenticated(errExternal))
}

func TestUnavailable(t *testing.T) {
	t.Run("Unavailable error", func(t *testing.T) {
		err := errors.Unavailable()
		assert.True(t, errors.IsUnavailable(err))
		assert.Equal(t, errors.ErrUnavailable, err.Type())
		assert.Equal(t, string(errors.ErrUnavailable), err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})

	t.Run("Unavailable error ", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s: operation with AdminID %v failed", errors.ErrUnavailable, 10)

		err := errors.Unavailable().WithMessage("operation with AdminID %v failed", 10)
		assert.True(t, errors.IsUnavailable(err))
		assert.Equal(t, errors.ErrUnavailable, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})

	t.Run("Unavailable error from external error", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s (wrapped error: third party error)", errors.ErrUnavailable)

		err := errors.UnavailableFromErr(errExternal)
		assert.True(t, errors.IsUnavailable(err))
		assert.Equal(t, errors.ErrUnavailable, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})

	t.Run("Unavailable error with error message from external error", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s: unexpected error (wrapped error: third party error)", errors.ErrUnavailable)

		err := errors.UnavailableFromErr(errExternal).WithMessage("unexpected error")
		assert.True(t, errors.IsUnavailable(err))
		assert.Equal(t, errors.ErrUnavailable, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})

	t.Run("Unavailable error with error message and human readable message", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s: operation with AdminID %v failed", errors.ErrUnavailable, 10)
		expectedHumanReadableMessage := "Please try again later"

		err := errors.Unavailable().WithMessage("operation with AdminID %v failed", 10).SetHumanReadableMessage("Please try again later")
		assert.True(t, errors.IsUnavailable(err))
		assert.Equal(t, errors.ErrUnavailable, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Equal(t, &expectedHumanReadableMessage, err.HumanReadableMessage())
	})

	t.Run("Unavailable error with error message and human readable message from external error", func(t *testing.T) {
		expectedErrMsg := fmt.Sprintf("%s: operation with AdminID %v failed (wrapped error: third party error)", errors.ErrUnavailable, 10)
		expectedHumanReadableMessage := "Please try again later"

		err := errors.UnavailableFromErr(errExternal).WithMessage("operation with AdminID %v failed", 10).SetHumanReadableMessage("Please try again later")
		assert.True(t, errors.IsUnavailable(err))
		assert.Equal(t, errors.ErrUnavailable, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Equal(t, &expectedHumanReadableMessage, err.HumanReadableMessage())
	})

	t.Run("Unavailable error wrapping another signare error", func(t *testing.T) {
		errNotFound := errors.NotFound().WithMessage("no records found")
		expectedErrMsg := fmt.Sprintf("%s (wrapped error: %s: no records found)", errors.ErrUnavailable, errors.ErrNotFound)

		err := errors.UnavailableFromErr(errNotFound)
		assert.False(t, errors.IsNotFound(err))
		assert.True(t, errors.IsUnavailable(err))
		assert.Equal(t, errors.ErrUnavailable, err.Type())
		assert.Equal(t, expectedErrMsg, err.Error())
		assert.Nil(t, err.HumanReadableMessage())
	})
	assert.False(t, errors.IsUnavailable(errExternal))
}

func TestCastAsSignerError(t *testing.T) {
	castedErr, ok := errors.CastAsUseCaseError(errExternal)
	assert.False(t, ok)
	assert.Nil(t, castedErr)

	alreadyExists := errors.AlreadyExists()
	castedErr, ok = errors.CastAsUseCaseError(alreadyExists)
	assert.True(t, ok)
	assert.Equal(t, alreadyExists, castedErr)

	badGateway := errors.BadGateway()
	castedErr, ok = errors.CastAsUseCaseError(badGateway)
	assert.True(t, ok)
	assert.Equal(t, badGateway, castedErr)

	internal := errors.Internal()
	castedErr, ok = errors.CastAsUseCaseError(internal)
	assert.True(t, ok)
	assert.Equal(t, internal, castedErr)

	invalidArgument := errors.InvalidArgument()
	castedErr, ok = errors.CastAsUseCaseError(invalidArgument)
	assert.True(t, ok)
	assert.Equal(t, invalidArgument, castedErr)

	notFound := errors.NotFound()
	castedErr, ok = errors.CastAsUseCaseError(notFound)
	assert.True(t, ok)
	assert.Equal(t, notFound, castedErr)

	notImplemented := errors.NotImplemented()
	castedErr, ok = errors.CastAsUseCaseError(notImplemented)
	assert.True(t, ok)
	assert.Equal(t, notImplemented, castedErr)

	permissionDenied := errors.PermissionDenied()
	castedErr, ok = errors.CastAsUseCaseError(permissionDenied)
	assert.True(t, ok)
	assert.Equal(t, permissionDenied, castedErr)

	preconditionFailed := errors.PreconditionFailed()
	castedErr, ok = errors.CastAsUseCaseError(preconditionFailed)
	assert.True(t, ok)
	assert.Equal(t, preconditionFailed, castedErr)

	timeout := errors.Timeout()
	castedErr, ok = errors.CastAsUseCaseError(timeout)
	assert.True(t, ok)
	assert.Equal(t, timeout, castedErr)

	tooManyReq := errors.TooManyReq()
	castedErr, ok = errors.CastAsUseCaseError(tooManyReq)
	assert.True(t, ok)
	assert.Equal(t, tooManyReq, castedErr)

	unauthenticated := errors.Unauthenticated()
	castedErr, ok = errors.CastAsUseCaseError(unauthenticated)
	assert.True(t, ok)
	assert.Equal(t, unauthenticated, castedErr)

	unavailable := errors.Unavailable()
	castedErr, ok = errors.CastAsUseCaseError(unavailable)
	assert.True(t, ok)
	assert.Equal(t, unavailable, castedErr)
}

func TestStackTrace(t *testing.T) {
	// This line corresponds to the location where the original error is created: inside errorParentFunc()
	// If this test file has lines added or removed this value must change accordingly.
	originalErrLine := 932
	err := wrapperErrorFunc(t)
	assert.Error(t, err)
	wrappedErr := errors.InternalFromErr(err)
	expectedErrorMessage := "Error: INTERNAL (wrapped error: PRECONDITION_FAILED: records can not be empty (wrapped error: NOT_FOUND (wrapped error: third party error)))"
	expectedFileAndLine := fmt.Sprintf("pkg/internal/errors/errors_test.go:%v", originalErrLine)
	expectedOriginalErrorStackTrace := `Original Error Stack Trace:
	at github.com/hyperledger-labs/signare/app/pkg/internal/errors_test.TestStackTrace`

	stackTrace := wrappedErr.GetStack()
	assert.Contains(t, stackTrace, expectedErrorMessage)
	assert.Contains(t, stackTrace, expectedFileAndLine)
	assert.Contains(t, stackTrace, expectedOriginalErrorStackTrace)
}

func wrapperErrorFunc(t *testing.T) error {
	t.Helper()
	err := errorParentFunc(t)
	assert.Error(t, err)
	return errors.PreconditionFailedFromErr(err).WithMessage("records can not be empty")
}

func errorParentFunc(t *testing.T) error {
	t.Helper()
	return errors.NotFoundFromErr(errExternal) // This is the line number of the original error being created.
}
