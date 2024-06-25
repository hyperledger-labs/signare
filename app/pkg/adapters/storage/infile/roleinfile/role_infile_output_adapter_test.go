package roleinfile_test

import (
	"context"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/hyperledger-labs/signare/app/pkg/adapters/storage/infile/roleinfile"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/authorization/role"

	"github.com/stretchr/testify/require"
)

const (
	defaultTestFilesPath = "testdata"
)

func NewDefaultRoleStorageInFileYAMLOutputAdapterForTest() (*roleinfile.DefaultRoleStorageInFile, error) {
	// get the location of the current file
	_, filename, _, _ := runtime.Caller(0)
	basePath := ""
	fileSystem := os.DirFS(path.Join(path.Dir(filename), defaultTestFilesPath))

	adapter, err := roleinfile.ProvideDefaultRoleStorageInFile(roleinfile.DefaultRoleStorageInFileOptions{
		FileSystem: fileSystem,
		BasePath:   basePath,
	})
	if err != nil {
		return nil, err
	}

	return adapter, nil
}

func TestDefaultRoleStorageInFileYAMLOutputAdapter_ListActions_Success(t *testing.T) {
	ctx := context.TODO()

	adapter, err := NewDefaultRoleStorageInFileYAMLOutputAdapterForTest()
	require.NoError(t, err)
	require.NotNil(t, adapter)

	listRolesInput := role.ListRolesInput{}
	listRolesOutput, listRolesErr := adapter.ListRoles(ctx, listRolesInput)
	require.NoError(t, listRolesErr)
	require.NotNil(t, listRolesOutput)

	expectedRoles := map[string]bool{
		"test-admin": false,
		"test-user":  false,
	}

	for _, role := range listRolesOutput.Roles {
		_, ok := expectedRoles[role.ID]
		require.True(t, ok)

		expectedRoles[role.ID] = true
	}

	for _, role := range listRolesOutput.Roles {
		found := expectedRoles[role.ID]
		require.True(t, found)
	}
}
