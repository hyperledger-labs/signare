package application_test

import (
	"context"
	"os"
	"testing"

	"github.com/hyperledger-labs/signare/app/pkg/adapters/storage/postgres/applicationdbout"
	"github.com/hyperledger-labs/signare/app/pkg/commons/validators"
	"github.com/hyperledger-labs/signare/app/pkg/entities"
	"github.com/hyperledger-labs/signare/app/pkg/graph"
	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/application"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/referentialintegrity"
	"github.com/hyperledger-labs/signare/app/test/dbtesthelper"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var app graph.GraphShared

var (
	chain       = entities.NewInt256FromInt(50)
	description = "test description"
)

func TestMain(m *testing.M) {
	a, err := dbtesthelper.InitializeApp()
	if err != nil {
		panic(err)
	}
	app = *a
	validators.SetValidators()
	os.Exit(m.Run())
}

func TestProvideDefaultUseCase(t *testing.T) {
	t.Run("nil storage", func(t *testing.T) {
		useCase, err := application.ProvideDefaultUseCase(application.DefaultUseCaseOptions{
			Storage:                     nil,
			ReferentialIntegrityUseCase: &referentialintegrity.DefaultUseCase{},
		})
		require.Error(t, err)
		require.Nil(t, useCase)
	})

	t.Run("nil referentialIntegrity use case", func(t *testing.T) {
		useCase, err := application.ProvideDefaultUseCase(application.DefaultUseCaseOptions{
			Storage:                     &applicationdbout.Repository{},
			ReferentialIntegrityUseCase: nil,
		})
		require.Error(t, err)
		require.Nil(t, useCase)
	})

	t.Run("success", func(t *testing.T) {
		options := application.DefaultUseCaseOptions{
			Storage:                     &applicationdbout.Repository{},
			ReferentialIntegrityUseCase: &referentialintegrity.DefaultUseCase{},
		}
		useCase, err := application.ProvideDefaultUseCase(options)
		require.NoError(t, err)
		require.NotNil(t, useCase)
	})
}

func TestDefaultUseCase_CreateApplication(t *testing.T) {
	ctx := context.Background()
	t.Run("success: if an Application ID is not defined, a random one must be assigned", func(t *testing.T) {
		input := application.CreateApplicationInput{
			ID:          nil,
			ChainID:     *chain,
			Description: &description,
		}
		output, err := app.ApplicationUseCase.CreateApplication(ctx, input)
		require.NoError(t, err)
		require.NotEmpty(t, output.ID)
		require.Equal(t, *chain, output.ChainID)
		require.Equal(t, description, *output.Description)
	})

	t.Run("success", func(t *testing.T) {
		randomID := uuid.New().String()
		input := application.CreateApplicationInput{
			ID:          &randomID,
			ChainID:     *chain,
			Description: &description,
		}
		output, err := app.ApplicationUseCase.CreateApplication(ctx, input)
		require.NoError(t, err)
		require.NotEmpty(t, output.InternalResourceID)
		require.Equal(t, randomID, output.ID)
		require.Equal(t, *chain, output.ChainID)
		require.Equal(t, description, *output.Description)
	})

	t.Run("failure: chain ID cannot be empty", func(t *testing.T) {
		invalidChainID := entities.NewInt256FromInt(0)
		input := application.CreateApplicationInput{
			ChainID:     *invalidChainID,
			Description: &description,
		}
		output, err := app.ApplicationUseCase.CreateApplication(ctx, input)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
		require.Nil(t, output)
	})

	t.Run("failure: application already exists", func(t *testing.T) {
		id := uuid.New().String()
		// First add
		input := application.CreateApplicationInput{
			ID:          &id,
			ChainID:     *chain,
			Description: &description,
		}
		_, err := app.ApplicationUseCase.CreateApplication(ctx, input)
		require.NoError(t, err)

		// Second add
		_, err = app.ApplicationUseCase.CreateApplication(ctx, input)
		require.Error(t, err)
		require.True(t, errors.IsAlreadyExists(err))
	})
}

func TestDefaultUseCase_ListApplications(t *testing.T) {
	ctx := context.Background()
	appsToCreate := 40
	for i := 0; i < appsToCreate; i++ {
		input := application.CreateApplicationInput{
			ID:          nil,
			ChainID:     *chain,
			Description: &description,
		}
		_, err := app.ApplicationUseCase.CreateApplication(ctx, input)
		require.NoError(t, err)
	}

	t.Run("failure: invalid input arguments", func(t *testing.T) {
		output, err := app.ApplicationUseCase.ListApplications(ctx, application.ListApplicationsInput{
			PageLimit:  -50,
			PageOffset: -30,
		})
		require.Nil(t, output)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
	})

	t.Run("success: list all applications", func(t *testing.T) {
		output, err := app.ApplicationUseCase.ListApplications(ctx, application.ListApplicationsInput{})
		require.NoError(t, err)
		require.True(t, len(output.Items) >= appsToCreate)
		for _, app := range output.Items {
			require.NotNil(t, app)
		}
	})

	t.Run("success: list DESC with limit", func(t *testing.T) {
		desiredLimit := 5
		input := application.ListApplicationsInput{
			OrderBy:        "creationDate",
			OrderDirection: "desc",
			PageLimit:      desiredLimit,
		}
		output, err := app.ApplicationUseCase.ListApplications(ctx, input)
		require.NoError(t, err)
		require.Len(t, output.Items, desiredLimit)
		require.True(t, output.MoreItems)
		// Assert order
		for i := 1; i < len(output.Items); i++ {
			require.Greater(t, output.Items[i-1].CreationDate.ToInt64(), output.Items[i].CreationDate.ToInt64())
		}
	})

	t.Run("success: list ASC with limit", func(t *testing.T) {
		desiredLimit := 5
		input := application.ListApplicationsInput{
			OrderBy:        "lastUpdate",
			OrderDirection: "asc",
			PageLimit:      desiredLimit,
		}
		output, err := app.ApplicationUseCase.ListApplications(ctx, input)
		require.NoError(t, err)
		require.Len(t, output.Items, desiredLimit)
		require.True(t, output.MoreItems)
		// Assert order
		for i := 1; i < len(output.Items); i++ {
			require.Less(t, output.Items[i-1].LastUpdate.ToInt64(), output.Items[i].LastUpdate.ToInt64())
		}
	})
}

func TestDefaultUseCase_GetApplication(t *testing.T) {
	ctx := context.Background()
	t.Run("failure: invalid ID", func(t *testing.T) {
		input := application.GetApplicationInput{
			StandardID: entities.StandardID{
				ID: "",
			},
		}
		output, err := app.ApplicationUseCase.GetApplication(ctx, input)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
		require.Nil(t, output)

		input = application.GetApplicationInput{
			StandardID: entities.StandardID{
				ID: "a-very-long-id-0000000000000000000000000000000000000000000000000000000000000",
			},
		}
		output, err = app.ApplicationUseCase.GetApplication(ctx, input)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
		require.Nil(t, output)
	})

	t.Run("failure: application not found", func(t *testing.T) {
		input := application.GetApplicationInput{
			StandardID: entities.StandardID{
				ID: "this-id-does-not-belong-to-any-application",
			},
		}
		output, err := app.ApplicationUseCase.GetApplication(ctx, input)
		require.Error(t, err)
		require.True(t, errors.IsNotFound(err))
		require.Nil(t, output)
	})

	t.Run("success", func(t *testing.T) {
		// Create application before retrieving it
		appID := uuid.New().String()
		_, err := app.ApplicationUseCase.CreateApplication(ctx, application.CreateApplicationInput{
			ID:          &appID,
			ChainID:     *chain,
			Description: &description,
		})
		require.NoError(t, err)

		app, err := app.ApplicationUseCase.GetApplication(ctx, application.GetApplicationInput{
			StandardID: entities.StandardID{
				ID: appID,
			},
		})
		require.NoError(t, err)
		require.NotNil(t, app)
		require.Equal(t, appID, app.ID)
		require.NotEmpty(t, app.InternalResourceID)
	})
}

func TestDefaultUseCase_EditApplication(t *testing.T) {
	ctx := context.Background()
	t.Run("failure: invalid input arguments", func(t *testing.T) {
		editInput := application.EditApplicationInput{
			ID:              "",
			ResourceVersion: "my-resource-version",
		}
		output, err := app.ApplicationUseCase.EditApplication(ctx, editInput)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
		require.Nil(t, output)

		editInput = application.EditApplicationInput{
			ID:              "my-id",
			ResourceVersion: "",
		}
		output, err = app.ApplicationUseCase.EditApplication(ctx, editInput)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
		require.Nil(t, output)
	})

	t.Run("failure: nonexistent application", func(t *testing.T) {
		editInput := application.EditApplicationInput{
			ID:              "my-id",
			ResourceVersion: "my-resource-version",
		}
		output, err := app.ApplicationUseCase.EditApplication(ctx, editInput)
		require.Error(t, err)
		require.True(t, errors.IsNotFound(err))
		require.Nil(t, output)
	})

	t.Run("failure: invalid resource version", func(t *testing.T) {
		// Create a valid application
		validAppID := uuid.New().String()
		_, err := app.ApplicationUseCase.CreateApplication(ctx, application.CreateApplicationInput{
			ID:          &validAppID,
			ChainID:     *chain,
			Description: &description,
		})
		require.NoError(t, err)

		// Edit the application with an invalid resource version
		editInput := application.EditApplicationInput{
			ID:              validAppID,
			ResourceVersion: "invalid-resource-version",
		}
		output, err := app.ApplicationUseCase.EditApplication(ctx, editInput)
		require.Error(t, err)
		require.True(t, errors.IsNotFound(err))
		require.Nil(t, output)
	})

	t.Run("success", func(t *testing.T) {
		// Create a valid application
		validAppID := uuid.New().String()
		createdApp, err := app.ApplicationUseCase.CreateApplication(ctx, application.CreateApplicationInput{
			ID:          &validAppID,
			ChainID:     *chain,
			Description: &description,
		})
		require.NoError(t, err)

		// Edit the application
		newDescription := "this is a new description"
		editInput := application.EditApplicationInput{
			ID:              createdApp.ID,
			ResourceVersion: createdApp.ResourceVersion,
			Description:     &newDescription,
		}
		editedApp, err := app.ApplicationUseCase.EditApplication(ctx, editInput)
		require.NoError(t, err)
		require.NotNil(t, editedApp)
		require.Equal(t, createdApp.CreationDate, editedApp.CreationDate)
		require.NotEqual(t, createdApp.LastUpdate, editedApp.LastUpdate)
		require.Equal(t, newDescription, *editedApp.Description)
		require.Equal(t, createdApp.ID, editedApp.ID)
		require.NotEqual(t, createdApp.ResourceVersion, editedApp.ResourceVersion)
	})
}

func TestDefaultUseCase_DeleteApplication(t *testing.T) {
	ctx := context.Background()
	t.Run("failure: invalid input arguments", func(t *testing.T) {
		deleteInput := application.DeleteApplicationInput{
			StandardID: entities.StandardID{
				ID: "",
			},
		}
		output, err := app.ApplicationUseCase.DeleteApplication(ctx, deleteInput)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
		require.Nil(t, output)
	})

	t.Run("failure: nonexistent application", func(t *testing.T) {
		deleteInput := application.DeleteApplicationInput{
			StandardID: entities.StandardID{
				ID: "my-id",
			},
		}
		output, err := app.ApplicationUseCase.DeleteApplication(ctx, deleteInput)
		require.Error(t, err)
		require.True(t, errors.IsNotFound(err))
		require.Nil(t, output)
	})

	t.Run("success", func(t *testing.T) {
		// Create a valid application
		validAppID := uuid.New().String()
		createdApp, err := app.ApplicationUseCase.CreateApplication(ctx, application.CreateApplicationInput{
			ID:          &validAppID,
			ChainID:     *chain,
			Description: &description,
		})
		require.NoError(t, err)

		// Delete the application
		deleteInput := application.DeleteApplicationInput{
			StandardID: entities.StandardID{
				ID: createdApp.ID,
			},
		}
		deletedApp, err := app.ApplicationUseCase.DeleteApplication(ctx, deleteInput)
		require.NoError(t, err)
		require.NotNil(t, deletedApp)
		require.Equal(t, createdApp.Application, deletedApp.Application)

		// Retrieve deleted application
		getOutput, err := app.ApplicationUseCase.GetApplication(ctx, application.GetApplicationInput{
			StandardID: entities.StandardID{
				ID: validAppID,
			},
		})
		require.Error(t, err)
		require.True(t, errors.IsNotFound(err))
		require.Nil(t, getOutput)
	})
}
