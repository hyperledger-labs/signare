package requestcontext_test

import (
	"context"
	"testing"

	"github.com/hyperledger-labs/signare/app/pkg/entities"
	"github.com/hyperledger-labs/signare/app/pkg/infra/requestcontext"

	"github.com/stretchr/testify/require"
)

const (
	testUser = "Tom"
	testApp  = "App"
)

func TestRequestContext_ValidUserAndApplication(t *testing.T) {
	ctx := context.WithValue(context.Background(), requestcontext.UserContextKey, testUser)
	user, err := requestcontext.UserFromContext(ctx)
	require.NoError(t, err)
	require.Equal(t, *user, testUser)

	ctx = context.WithValue(ctx, requestcontext.ApplicationContextKey, testApp)
	app, err := requestcontext.ApplicationFromContext(ctx)
	require.NoError(t, err)
	require.Equal(t, *app, testApp)
}

func TestRequestContext_InValidUserAndApplication(t *testing.T) {
	invalidHeader := entities.ContextKey("invalid-header")
	ctx := context.WithValue(context.Background(), invalidHeader, testUser)
	user, err := requestcontext.UserFromContext(ctx)
	require.Error(t, err)
	require.Nil(t, user)

	ctx = context.WithValue(ctx, invalidHeader, testApp)
	app, err := requestcontext.ApplicationFromContext(ctx)
	require.Error(t, err)
	require.Nil(t, app)
}
