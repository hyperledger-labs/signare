package referentialintegrity_test

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/hyperledger-labs/signare/app/pkg/adapters/storage/postgres/referentialintegritydbout"
	"github.com/hyperledger-labs/signare/app/pkg/commons/validators"
	"github.com/hyperledger-labs/signare/app/pkg/entities"
	"github.com/hyperledger-labs/signare/app/pkg/graph"
	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/referentialintegrity"
	"github.com/hyperledger-labs/signare/app/test/dbtesthelper"
)

var app graph.GraphShared

var (
	resourceKind       = referentialintegrity.ResourceKind("user")
	parentResourceKind = referentialintegrity.ResourceKind("application")
)

func TestMain(m *testing.M) {
	application, err := dbtesthelper.InitializeApp()
	if err != nil {
		panic(err)
	}
	app = *application
	validators.SetValidators()
	os.Exit(m.Run())
}

func TestProvideDefaultUseCase(t *testing.T) {
	t.Run("nil storage", func(t *testing.T) {
		useCase, err := referentialintegrity.ProvideDefaultUseCase(referentialintegrity.DefaultUseCaseOptions{
			Storage: nil,
		})
		require.Error(t, err)
		require.Nil(t, useCase)
	})

	t.Run("success", func(t *testing.T) {
		useCase, err := referentialintegrity.ProvideDefaultUseCase(referentialintegrity.DefaultUseCaseOptions{
			Storage: &referentialintegritydbout.Repository{},
		})
		require.NoError(t, err)
		require.NotNil(t, useCase)
	})
}

func TestDefaultUseCase_CreateEntry(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		input := referentialintegrity.CreateEntryInput{
			ResourceID:         uuid.NewString(),
			ResourceKind:       resourceKind,
			ParentResourceID:   uuid.NewString(),
			ParentResourceKind: parentResourceKind,
		}
		output, err := app.ReferentialIntegrityUseCase.CreateEntry(ctx, input)
		require.NoError(t, err)
		require.NotNil(t, output)
	})

	t.Run("invalid argument: resource ID empty", func(t *testing.T) {
		input := referentialintegrity.CreateEntryInput{
			ResourceID:         "",
			ResourceKind:       resourceKind,
			ParentResourceID:   uuid.NewString(),
			ParentResourceKind: parentResourceKind,
		}
		output, err := app.ReferentialIntegrityUseCase.CreateEntry(ctx, input)
		require.Error(t, err)
		require.Nil(t, output)
	})

	t.Run("invalid argument: resource kind empty", func(t *testing.T) {
		input := referentialintegrity.CreateEntryInput{
			ResourceID:         uuid.NewString(),
			ResourceKind:       "",
			ParentResourceID:   uuid.NewString(),
			ParentResourceKind: parentResourceKind,
		}
		output, err := app.ReferentialIntegrityUseCase.CreateEntry(ctx, input)
		require.Error(t, err)
		require.Nil(t, output)
	})

	t.Run("invalid argument: parent resource ID empty", func(t *testing.T) {
		input := referentialintegrity.CreateEntryInput{
			ResourceID:         uuid.NewString(),
			ResourceKind:       resourceKind,
			ParentResourceID:   "",
			ParentResourceKind: parentResourceKind,
		}
		output, err := app.ReferentialIntegrityUseCase.CreateEntry(ctx, input)
		require.Error(t, err)
		require.Nil(t, output)
	})

	t.Run("invalid argument: parent resource kind empty", func(t *testing.T) {
		input := referentialintegrity.CreateEntryInput{
			ResourceID:         uuid.NewString(),
			ResourceKind:       resourceKind,
			ParentResourceID:   uuid.NewString(),
			ParentResourceKind: "",
		}
		output, err := app.ReferentialIntegrityUseCase.CreateEntry(ctx, input)
		require.Error(t, err)
		require.Nil(t, output)
	})
}

func TestDefaultUseCase_GetEntry(t *testing.T) {
	ctx := context.Background()

	t.Run("failure: invalid ID", func(t *testing.T) {
		input := referentialintegrity.GetEntryInput{
			StandardID: entities.StandardID{
				ID: "",
			},
		}
		output, err := app.ReferentialIntegrityUseCase.GetEntry(ctx, input)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
		require.Nil(t, output)
	})

	t.Run("invalid argument: resource ID empty", func(t *testing.T) {
		input := referentialintegrity.GetEntryInput{
			StandardID: entities.StandardID{
				ID: "not-found",
			},
		}
		output, err := app.ReferentialIntegrityUseCase.GetEntry(ctx, input)
		require.Error(t, err)
		require.True(t, errors.IsNotFound(err))
		require.Nil(t, output)
	})

	t.Run("success", func(t *testing.T) {
		inputCreate := referentialintegrity.CreateEntryInput{
			ResourceID:         uuid.NewString(),
			ResourceKind:       resourceKind,
			ParentResourceID:   uuid.NewString(),
			ParentResourceKind: parentResourceKind,
		}
		outputCreate, err := app.ReferentialIntegrityUseCase.CreateEntry(ctx, inputCreate)
		require.NoError(t, err)

		input := referentialintegrity.GetEntryInput{
			StandardID: outputCreate.StandardID,
		}
		output, err := app.ReferentialIntegrityUseCase.GetEntry(ctx, input)
		require.NoError(t, err)
		require.NotNil(t, output)
	})
}

func TestDefaultUseCase_ListEntries(t *testing.T) {
	ctx := context.Background()
	notUsedKind := referentialintegrity.ResourceKind("not-used")

	deterministicResourceID := uuid.NewString()
	deterministicParentID := uuid.NewString()
	inputCreate := referentialintegrity.CreateEntryInput{
		ResourceID:         deterministicResourceID,
		ResourceKind:       resourceKind,
		ParentResourceID:   deterministicParentID,
		ParentResourceKind: parentResourceKind,
	}
	_, errCreate := app.ReferentialIntegrityUseCase.CreateEntry(ctx, inputCreate)
	require.NoError(t, errCreate)

	entriesToCreate := 5
	for i := 0; i < entriesToCreate; i++ {
		input := referentialintegrity.CreateEntryInput{
			ResourceID:         uuid.NewString(),
			ResourceKind:       resourceKind,
			ParentResourceID:   uuid.NewString(),
			ParentResourceKind: parentResourceKind,
		}
		_, err := app.ReferentialIntegrityUseCase.CreateEntry(ctx, input)
		require.NoError(t, err)
	}

	t.Run("success: no elements for non existing resource", func(t *testing.T) {
		input := referentialintegrity.ListEntriesInput{
			Resource: &referentialintegrity.Resource{
				ResourceID:   uuid.NewString(),
				ResourceKind: notUsedKind,
			},
		}
		output, err := app.ReferentialIntegrityUseCase.ListEntries(ctx, input)
		require.NoError(t, err)
		require.NotNil(t, output)
		require.Len(t, output.Items, 0)
	})

	t.Run("success: no elements for non existing parent resource", func(t *testing.T) {
		input := referentialintegrity.ListEntriesInput{
			Parent: &referentialintegrity.Resource{
				ResourceID:   uuid.NewString(),
				ResourceKind: notUsedKind,
			},
		}
		output, err := app.ReferentialIntegrityUseCase.ListEntries(ctx, input)
		require.NoError(t, err)
		require.NotNil(t, output)
		require.Len(t, output.Items, 0)
	})

	t.Run("success: filter by resource", func(t *testing.T) {
		input := referentialintegrity.ListEntriesInput{
			Resource: &referentialintegrity.Resource{
				ResourceID:   deterministicResourceID,
				ResourceKind: resourceKind,
			},
		}
		output, err := app.ReferentialIntegrityUseCase.ListEntries(ctx, input)
		require.NoError(t, err)
		require.NotNil(t, output)
		require.Len(t, output.Items, 1)
	})

	t.Run("success: filter by parent resource", func(t *testing.T) {
		input := referentialintegrity.ListEntriesInput{
			Parent: &referentialintegrity.Resource{
				ResourceID:   deterministicParentID,
				ResourceKind: parentResourceKind,
			},
		}
		output, err := app.ReferentialIntegrityUseCase.ListEntries(ctx, input)
		require.NoError(t, err)
		require.NotNil(t, output)
		require.Len(t, output.Items, 1)
	})

	t.Run("error: invalid page limit", func(t *testing.T) {
		input := referentialintegrity.ListEntriesInput{
			PageLimit: -1,
		}
		output, err := app.ReferentialIntegrityUseCase.ListEntries(ctx, input)
		require.Error(t, err)
		require.Nil(t, output)
	})

	t.Run("error: invalid page offset", func(t *testing.T) {
		input := referentialintegrity.ListEntriesInput{
			PageOffset: -1,
		}
		output, err := app.ReferentialIntegrityUseCase.ListEntries(ctx, input)
		require.Error(t, err)
		require.Nil(t, output)
	})
}

func TestDefaultUseCase_DeleteEntry(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		inputCreate := referentialintegrity.CreateEntryInput{
			ResourceID:         uuid.NewString(),
			ResourceKind:       resourceKind,
			ParentResourceID:   uuid.NewString(),
			ParentResourceKind: parentResourceKind,
		}
		outputCreate, err := app.ReferentialIntegrityUseCase.CreateEntry(ctx, inputCreate)
		require.NoError(t, err)

		inputDelete := referentialintegrity.DeleteEntryInput{
			StandardID: outputCreate.StandardID,
		}
		output, err := app.ReferentialIntegrityUseCase.DeleteEntry(ctx, inputDelete)
		require.NoError(t, err)
		require.NotNil(t, output)
	})

	t.Run("invalid argument: resource ID empty", func(t *testing.T) {
		inputCreate := referentialintegrity.CreateEntryInput{
			ResourceID:         uuid.NewString(),
			ResourceKind:       resourceKind,
			ParentResourceID:   uuid.NewString(),
			ParentResourceKind: parentResourceKind,
		}
		_, err := app.ReferentialIntegrityUseCase.CreateEntry(ctx, inputCreate)
		require.NoError(t, err)

		inputDelete := referentialintegrity.DeleteEntryInput{
			StandardID: entities.StandardID{
				ID: "",
			},
		}
		output, err := app.ReferentialIntegrityUseCase.DeleteEntry(ctx, inputDelete)
		require.Error(t, err)
		require.Nil(t, output)
		require.True(t, errors.IsInvalidArgument(err))
	})

	t.Run("not found", func(t *testing.T) {
		inputCreate := referentialintegrity.CreateEntryInput{
			ResourceID:         uuid.NewString(),
			ResourceKind:       resourceKind,
			ParentResourceID:   uuid.NewString(),
			ParentResourceKind: parentResourceKind,
		}
		_, err := app.ReferentialIntegrityUseCase.CreateEntry(ctx, inputCreate)
		require.NoError(t, err)

		inputDelete := referentialintegrity.DeleteEntryInput{
			StandardID: entities.StandardID{
				ID: "not-found",
			},
		}
		output, err := app.ReferentialIntegrityUseCase.DeleteEntry(ctx, inputDelete)
		require.Error(t, err)
		require.Nil(t, output)
		require.True(t, errors.IsNotFound(err))
	})
}

func TestDefaultUseCase_ListMyChildrenEntries(t *testing.T) {
	ctx := context.Background()
	parentID := uuid.NewString()

	entriesToCreate := 5
	for i := 0; i < entriesToCreate; i++ {
		input := referentialintegrity.CreateEntryInput{
			ResourceID:         uuid.NewString(),
			ResourceKind:       resourceKind,
			ParentResourceID:   parentID,
			ParentResourceKind: parentResourceKind,
		}
		_, err := app.ReferentialIntegrityUseCase.CreateEntry(ctx, input)
		require.NoError(t, err)
	}

	t.Run("success: no elements for unused resource kind", func(t *testing.T) {
		input := referentialintegrity.ListMyChildrenEntriesInput{
			ParentResourceID:   uuid.NewString(),
			ParentResourceKind: parentResourceKind,
		}
		output, err := app.ReferentialIntegrityUseCase.ListMyChildrenEntries(ctx, input)
		require.NoError(t, err)
		require.NotNil(t, output)
		require.Len(t, output.Items, 0)
	})

	t.Run("success: found children", func(t *testing.T) {
		input := referentialintegrity.ListMyChildrenEntriesInput{
			ParentResourceID:   parentID,
			ParentResourceKind: parentResourceKind,
		}
		output, err := app.ReferentialIntegrityUseCase.ListMyChildrenEntries(ctx, input)
		require.NoError(t, err)
		require.NotNil(t, output)
		require.True(t, len(output.Items) >= entriesToCreate)
	})

	t.Run("error: invalid parent resource ID", func(t *testing.T) {
		input := referentialintegrity.ListMyChildrenEntriesInput{
			ParentResourceID:   "",
			ParentResourceKind: parentResourceKind,
		}
		output, err := app.ReferentialIntegrityUseCase.ListMyChildrenEntries(ctx, input)
		require.Error(t, err)
		require.Nil(t, output)
	})

	t.Run("error: invalid parent resource kind", func(t *testing.T) {
		input := referentialintegrity.ListMyChildrenEntriesInput{
			ParentResourceID:   parentID,
			ParentResourceKind: "",
		}
		output, err := app.ReferentialIntegrityUseCase.ListMyChildrenEntries(ctx, input)
		require.Error(t, err)
		require.Nil(t, output)
	})
}

func TestDefaultUseCase_DeleteMyEntriesIfAny(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		id := uuid.NewString()
		inputCreate := referentialintegrity.CreateEntryInput{
			ResourceID:         id,
			ResourceKind:       resourceKind,
			ParentResourceID:   uuid.NewString(),
			ParentResourceKind: parentResourceKind,
		}
		_, err := app.ReferentialIntegrityUseCase.CreateEntry(ctx, inputCreate)
		require.NoError(t, err)

		inputCreate = referentialintegrity.CreateEntryInput{
			ResourceID:         id,
			ResourceKind:       resourceKind,
			ParentResourceID:   uuid.NewString(),
			ParentResourceKind: parentResourceKind,
		}
		_, err = app.ReferentialIntegrityUseCase.CreateEntry(ctx, inputCreate)
		require.NoError(t, err)

		inputDelete := referentialintegrity.DeleteMyEntriesIfAnyInput{
			ResourceID:   id,
			ResourceKind: resourceKind,
		}
		err = app.ReferentialIntegrityUseCase.DeleteMyEntriesIfAny(ctx, inputDelete)
		require.NoError(t, err)
	})

	t.Run("invalid argument: resource ID empty", func(t *testing.T) {
		inputDelete := referentialintegrity.DeleteMyEntriesIfAnyInput{
			ResourceID:   "",
			ResourceKind: resourceKind,
		}
		err := app.ReferentialIntegrityUseCase.DeleteMyEntriesIfAny(ctx, inputDelete)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
	})

	t.Run("invalid argument: resource kind empty", func(t *testing.T) {
		inputDelete := referentialintegrity.DeleteMyEntriesIfAnyInput{
			ResourceID:   uuid.NewString(),
			ResourceKind: "",
		}
		err := app.ReferentialIntegrityUseCase.DeleteMyEntriesIfAny(ctx, inputDelete)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
	})
}

func TestDefaultUseCase_GetEntryByResourceAndParent(t *testing.T) {
	ctx := context.Background()

	t.Run("failure: invalid resource ID", func(t *testing.T) {
		input := referentialintegrity.GetEntryByResourceAndParentInput{
			ResourceID:         "",
			ResourceKind:       resourceKind,
			ParentResourceID:   uuid.NewString(),
			ParentResourceKind: parentResourceKind,
		}
		output, err := app.ReferentialIntegrityUseCase.GetEntryByResourceAndParent(ctx, input)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
		require.Nil(t, output)
	})

	t.Run("invalid argument: resource kind empty", func(t *testing.T) {
		input := referentialintegrity.GetEntryByResourceAndParentInput{
			ResourceID:         uuid.NewString(),
			ResourceKind:       "",
			ParentResourceID:   uuid.NewString(),
			ParentResourceKind: parentResourceKind,
		}
		output, err := app.ReferentialIntegrityUseCase.GetEntryByResourceAndParent(ctx, input)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
		require.Nil(t, output)
	})

	t.Run("invalid argument: parent resource ID empty", func(t *testing.T) {
		input := referentialintegrity.GetEntryByResourceAndParentInput{
			ResourceID:         uuid.NewString(),
			ResourceKind:       resourceKind,
			ParentResourceID:   "",
			ParentResourceKind: parentResourceKind,
		}
		output, err := app.ReferentialIntegrityUseCase.GetEntryByResourceAndParent(ctx, input)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
		require.Nil(t, output)
	})

	t.Run("invalid argument: parent resource kind empty", func(t *testing.T) {
		input := referentialintegrity.GetEntryByResourceAndParentInput{
			ResourceID:         uuid.NewString(),
			ResourceKind:       resourceKind,
			ParentResourceID:   uuid.NewString(),
			ParentResourceKind: "",
		}
		output, err := app.ReferentialIntegrityUseCase.GetEntryByResourceAndParent(ctx, input)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
		require.Nil(t, output)
	})

	t.Run("success", func(t *testing.T) {
		inputCreate := referentialintegrity.CreateEntryInput{
			ResourceID:         uuid.NewString(),
			ResourceKind:       resourceKind,
			ParentResourceID:   uuid.NewString(),
			ParentResourceKind: parentResourceKind,
		}
		outputCreate, err := app.ReferentialIntegrityUseCase.CreateEntry(ctx, inputCreate)
		require.NoError(t, err)

		input := referentialintegrity.GetEntryByResourceAndParentInput{
			ResourceID:         outputCreate.ResourceID,
			ResourceKind:       resourceKind,
			ParentResourceID:   outputCreate.ParentResourceID,
			ParentResourceKind: parentResourceKind,
		}
		output, err := app.ReferentialIntegrityUseCase.GetEntryByResourceAndParent(ctx, input)
		require.NoError(t, err)
		require.NotNil(t, output)
	})
}
