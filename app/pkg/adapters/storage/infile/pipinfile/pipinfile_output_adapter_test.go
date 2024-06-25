package pipinfile_test

import (
	"context"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/hyperledger-labs/signare/app/pkg/adapters/storage/infile/pipinfile"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/authorization/pdp"

	"github.com/stretchr/testify/require"
)

const (
	defaultRoleTestAdmin = "test-admin"
	defaultRoleTestUser  = "test-user"

	defaultTestFilesPath = "testdata"

	defaultActionOne   = "generated.action.one"
	defaultActionTwo   = "generated.action.two"
	defaultActionThree = "generated.action.three"
	defaultActionFour  = "manual.action.four"
	defaultActionFive  = "manual.action.five"
)

var (
	expectedTestAdminActions = pdp.NewActions([]string{defaultActionOne, defaultActionTwo, defaultActionThree, defaultActionFour, defaultActionFive})
	expectedTestUserActions  = pdp.NewActions([]string{defaultActionFour, defaultActionFive})
)

func NewDefaultRBACPolicyInformationPointYAMLOutputAdapterForTest() (*pipinfile.DefaultRBACActionsPolicyInformationPointYAMLOutputAdapter, error) {
	// get the location of the current file
	_, filename, _, _ := runtime.Caller(0)
	basePath := ""
	fileSystem := os.DirFS(path.Join(path.Dir(filename), defaultTestFilesPath))

	adapter, err := pipinfile.ProvideDefaultRBACActionsPolicyInformationPointYAMLOutputAdapter(pipinfile.DefaultRBACActionsPolicyInformationPointYAMLOutputAdapterOptions{
		FileSystem: fileSystem,
		BasePath:   basePath,
	})
	if err != nil {
		return nil, err
	}

	return adapter, nil
}

func TestDefaultRBACPolicyInformationPointYAMLOutputAdapter_ListActions_Success(t *testing.T) {
	ctx := context.TODO()

	adapter, err := NewDefaultRBACPolicyInformationPointYAMLOutputAdapterForTest()
	require.NoError(t, err)
	require.NotNil(t, adapter)

	listActionsInputAdmin := pdp.ListActionsInput{
		Roles: []string{defaultRoleTestAdmin},
	}
	adminActions, failure := adapter.ListActions(ctx, listActionsInputAdmin)
	require.Nil(t, failure)
	require.True(t, adminActions.Actions.Equal(*expectedTestAdminActions))

	listActionsInputUser := pdp.ListActionsInput{
		Roles: []string{defaultRoleTestUser},
	}
	userActions, failure := adapter.ListActions(ctx, listActionsInputUser)
	require.Nil(t, failure)
	require.True(t, userActions.Actions.Equal(*expectedTestUserActions))
}
