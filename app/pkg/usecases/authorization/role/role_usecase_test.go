package role_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hyperledger-labs/signare/app/pkg/adapters/storage/infile/roleinfile"
	"github.com/hyperledger-labs/signare/app/pkg/commons/validators"
	"github.com/hyperledger-labs/signare/app/pkg/graph"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/authorization/role"
	"github.com/hyperledger-labs/signare/app/test/dbtesthelper"
)

var app graph.GraphShared

func TestMain(m *testing.M) {
	application, err := dbtesthelper.InitializeApp()
	if err != nil {
		panic(err)
	}
	app = *application
	validators.SetValidators()
	os.Exit(m.Run())
}

func TestProvideDefaultRoleUseCase(t *testing.T) {
	t.Run("nil storage", func(t *testing.T) {
		options := role.DefaultRoleUseCaseOptions{}
		defaultUseCase, err := role.ProvideDefaultRoleUseCase(options)
		require.Error(t, err)
		require.Nil(t, defaultUseCase)
	})

	t.Run("success", func(t *testing.T) {
		options := role.DefaultRoleUseCaseOptions{
			RoleStorage: &roleinfile.DefaultRoleStorageInFile{},
		}
		defaultUseCase, err := role.ProvideDefaultRoleUseCase(options)
		require.NoError(t, err)
		require.NotNil(t, defaultUseCase)
	})
}

func TestProvideDefaultRoleUseCase_CreateHSM(t *testing.T) {
	ctx := context.Background()
	t.Run("success", func(t *testing.T) {
		getSupportedRolesInput := role.GetSupportedRolesInput{}
		getSupportedRolesOutput, getSupportedRolesErr := app.RoleUseCase.GetSupportedRoles(ctx, getSupportedRolesInput)
		require.Nil(t, getSupportedRolesErr)
		require.NotNil(t, getSupportedRolesOutput)
	})
}
