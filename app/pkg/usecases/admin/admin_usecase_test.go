package admin_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hyperledger-labs/signare/app/pkg/adapters/storage/postgres/admindbout"
	"github.com/hyperledger-labs/signare/app/pkg/commons/validators"
	"github.com/hyperledger-labs/signare/app/pkg/entities"
	"github.com/hyperledger-labs/signare/app/pkg/graph"
	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/admin"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/authorization/role"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/referentialintegrity"
	"github.com/hyperledger-labs/signare/app/test/dbtesthelper"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
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

func TestProvideDefaultUseCase(t *testing.T) {
	t.Run("nil admin storage", func(t *testing.T) {
		options := admin.DefaultUseCaseOptions{
			AdminStorage:                nil,
			RoleUseCase:                 &role.DefaultRoleUseCase{},
			ReferentialIntegrityUseCase: &referentialintegrity.DefaultUseCase{},
		}
		defaultUseCase, err := admin.ProvideDefaultUseCase(options)
		require.Error(t, err)
		require.Nil(t, defaultUseCase)
	})

	t.Run("nil role use case", func(t *testing.T) {
		options := admin.DefaultUseCaseOptions{
			AdminStorage:                &admindbout.Repository{},
			RoleUseCase:                 nil,
			ReferentialIntegrityUseCase: &referentialintegrity.DefaultUseCase{},
		}
		defaultUseCase, err := admin.ProvideDefaultUseCase(options)
		require.Error(t, err)
		require.Nil(t, defaultUseCase)
	})

	t.Run("nil referential integrity use case", func(t *testing.T) {
		options := admin.DefaultUseCaseOptions{
			AdminStorage:                &admindbout.Repository{},
			RoleUseCase:                 &role.DefaultRoleUseCase{},
			ReferentialIntegrityUseCase: nil,
		}
		defaultUseCase, err := admin.ProvideDefaultUseCase(options)
		require.Error(t, err)
		require.Nil(t, defaultUseCase)
	})

	t.Run("success", func(t *testing.T) {
		options := admin.DefaultUseCaseOptions{
			AdminStorage:                &admindbout.Repository{},
			RoleUseCase:                 &role.DefaultRoleUseCase{},
			ReferentialIntegrityUseCase: &referentialintegrity.DefaultUseCase{},
		}
		defaultUseCase, err := admin.ProvideDefaultUseCase(options)
		require.NoError(t, err)
		require.NotNil(t, defaultUseCase)
	})
}

func TestDefaultUseCase_CreateAdmin(t *testing.T) {
	ctx := context.Background()
	description := "admin for testing"
	adminID := entities.StandardID{
		ID: "admin-test",
	}
	t.Run("success", func(t *testing.T) {
		createAdminInput := admin.CreateAdminInput{
			StandardID:  adminID,
			Description: &description,
		}
		createAdminOutput, createAdminErr := app.AdminUseCase.CreateAdmin(ctx, createAdminInput)
		require.Nil(t, createAdminErr)
		require.NotEmpty(t, createAdminOutput.InternalResourceID)
		require.NotNil(t, createAdminOutput)
	})

	t.Run("failure: admin already exists", func(t *testing.T) {
		createAdminInput := admin.CreateAdminInput{
			StandardID:  adminID,
			Description: &description,
		}
		createAdminOutput, createAdminErr := app.AdminUseCase.CreateAdmin(ctx, createAdminInput)
		require.NotNil(t, createAdminErr)
		require.True(t, errors.IsAlreadyExists(createAdminErr))
		require.Nil(t, createAdminOutput)
	})
}

func TestDefaultUseCase_ListAdmins(t *testing.T) {
	ctx := context.Background()
	adminsToCreate := 40
	description := "admin for testing"
	createAdminInput := admin.CreateAdminInput{
		StandardID:  entities.StandardID{},
		Description: &description,
	}
	for i := 0; i < adminsToCreate; i++ {
		createAdminInput.ID = fmt.Sprintf("test-admin-%d", i)
		createAdminOutput, createAdminErr := app.AdminUseCase.CreateAdmin(ctx, createAdminInput)
		require.Nil(t, createAdminErr)
		require.NotNil(t, createAdminOutput)
	}

	t.Run("failure: invalid input arguments", func(t *testing.T) {
		listAdminsInput := admin.ListAdminsInput{
			PageLimit:  -50,
			PageOffset: -30,
		}
		listAdminsOutput, listAdminsErr := app.AdminUseCase.ListAdmins(ctx, listAdminsInput)
		require.NotNil(t, listAdminsErr)
		require.Nil(t, listAdminsOutput)
	})

	t.Run("success: list all applications", func(t *testing.T) {
		listAdminsInput := admin.ListAdminsInput{}
		listAdminsOutput, listAdminsErr := app.AdminUseCase.ListAdmins(ctx, listAdminsInput)
		require.NoError(t, listAdminsErr)
		require.True(t, len(listAdminsOutput.Items) >= adminsToCreate)
	})

	t.Run("success: list DESC with limit", func(t *testing.T) {
		desiredLimit := 5
		listAdminsInput := admin.ListAdminsInput{
			OrderBy:        "creationDate",
			OrderDirection: "desc",
			PageLimit:      desiredLimit,
		}
		listAdminsOutput, listAdminsErr := app.AdminUseCase.ListAdmins(ctx, listAdminsInput)
		require.NoError(t, listAdminsErr)
		require.Len(t, listAdminsOutput.Items, desiredLimit)
		require.True(t, listAdminsOutput.MoreItems)
		// Assert order
		for i := 1; i < len(listAdminsOutput.Items); i++ {
			require.GreaterOrEqual(t, listAdminsOutput.Items[i-1].CreationDate.ToInt64(), listAdminsOutput.Items[i].CreationDate.ToInt64())
		}
	})

	t.Run("success: list ASC with limit", func(t *testing.T) {
		desiredLimit := 5
		listAdminsInput := admin.ListAdminsInput{
			OrderBy:        "creationDate",
			OrderDirection: "asc",
			PageLimit:      desiredLimit,
		}
		listAdminsOutput, listAdminsErr := app.AdminUseCase.ListAdmins(ctx, listAdminsInput)
		require.NoError(t, listAdminsErr)
		require.Len(t, listAdminsOutput.Items, desiredLimit)
		require.True(t, listAdminsOutput.MoreItems)
		// Assert order
		for i := 1; i < len(listAdminsOutput.Items); i++ {
			require.Less(t, listAdminsOutput.Items[i-1].LastUpdate.ToInt64(), listAdminsOutput.Items[i].LastUpdate.ToInt64())
		}
	})
}

func TestDefaultUseCase_GetAdmin(t *testing.T) {
	ctx := context.Background()
	t.Run("failure: invalid ID", func(t *testing.T) {
		getAdminInput := admin.GetAdminInput{
			StandardID: entities.StandardID{
				ID: "",
			},
		}
		getAdminOutput, getAdminErr := app.AdminUseCase.GetAdmin(ctx, getAdminInput)
		require.Error(t, getAdminErr)
		require.True(t, errors.IsInvalidArgument(getAdminErr))
		require.Nil(t, getAdminOutput)

		getAdminInput = admin.GetAdminInput{
			StandardID: entities.StandardID{
				ID: "a-very-long-id-0000000000000000000000000000000000000000000000000000000000000",
			},
		}
		getAdminOutput, getAdminErr = app.AdminUseCase.GetAdmin(ctx, getAdminInput)
		require.Error(t, getAdminErr)
		require.True(t, errors.IsInvalidArgument(getAdminErr))
		require.Nil(t, getAdminOutput)
	})

	t.Run("failure: application not found", func(t *testing.T) {
		getAdminInput := admin.GetAdminInput{
			StandardID: entities.StandardID{
				ID: "this-id-does-not-belong-to-any-admin",
			},
		}
		output, err := app.AdminUseCase.GetAdmin(ctx, getAdminInput)
		require.Error(t, err)
		require.True(t, errors.IsNotFound(err))
		require.Nil(t, output)
	})

	t.Run("success", func(t *testing.T) {
		// Create application before retrieving it
		adminID := entities.StandardID{
			ID: uuid.New().String(),
		}
		createAdminInput := admin.CreateAdminInput{
			StandardID: adminID,
		}
		_, err := app.AdminUseCase.CreateAdmin(ctx, createAdminInput)
		require.NoError(t, err)

		getAdminInput := admin.GetAdminInput{
			StandardID: adminID,
		}
		admin, err := app.AdminUseCase.GetAdmin(ctx, getAdminInput)
		require.NoError(t, err)
		require.NotNil(t, admin)
		require.Equal(t, adminID.ID, admin.ID)
		require.NotEmpty(t, admin.InternalResourceID)
	})
}

func TestDefaultUseCase_EditAdmin(t *testing.T) {
	ctx := context.Background()
	t.Run("failure: invalid input arguments", func(t *testing.T) {
		editAdminInput := admin.EditAdminInput{
			AdminEditable: admin.AdminEditable{
				StandardResourceMeta: entities.StandardResourceMeta{
					StandardResource: entities.StandardResource{
						StandardID: entities.StandardID{
							ID: "",
						},
					},
					ResourceVersion: "my-resource-version",
				},
			},
		}
		editAdminOutput, err := app.AdminUseCase.EditAdmin(ctx, editAdminInput)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
		require.Nil(t, editAdminOutput)

		editAdminInput = admin.EditAdminInput{
			AdminEditable: admin.AdminEditable{
				StandardResourceMeta: entities.StandardResourceMeta{
					StandardResource: entities.StandardResource{
						StandardID: entities.StandardID{
							ID: "",
						},
					},
					ResourceVersion: "",
				},
			},
		}
		editAdminOutput, err = app.AdminUseCase.EditAdmin(ctx, editAdminInput)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
		require.Nil(t, editAdminOutput)
	})

	t.Run("failure: nonexistent admin", func(t *testing.T) {
		editAdminInput := admin.EditAdminInput{
			AdminEditable: admin.AdminEditable{
				StandardResourceMeta: entities.StandardResourceMeta{
					StandardResource: entities.StandardResource{
						StandardID: entities.StandardID{
							ID: "my-id",
						},
					},
					ResourceVersion: "my-resource-version",
				},
			},
		}
		editAdminOutput, err := app.AdminUseCase.EditAdmin(ctx, editAdminInput)
		require.Error(t, err)
		require.True(t, errors.IsNotFound(err))
		require.Nil(t, editAdminOutput)
	})

	t.Run("failure: invalid resource version", func(t *testing.T) {
		// Create a valid admin
		validAdminID := entities.StandardID{
			ID: uuid.New().String(),
		}
		createAdminInput := admin.CreateAdminInput{
			StandardID: validAdminID,
		}
		_, createAdminErr := app.AdminUseCase.CreateAdmin(ctx, createAdminInput)
		require.NoError(t, createAdminErr)

		// Edit the application with an invalid resource version
		editAdminInput := admin.EditAdminInput{
			AdminEditable: admin.AdminEditable{
				StandardResourceMeta: entities.StandardResourceMeta{
					StandardResource: entities.StandardResource{
						StandardID: validAdminID,
					},
					ResourceVersion: "invalid-resource-version",
				},
			},
		}
		editAdminOutput, editAdminErr := app.AdminUseCase.EditAdmin(ctx, editAdminInput)
		require.Error(t, editAdminErr)
		require.True(t, errors.IsNotFound(editAdminErr))
		require.Nil(t, editAdminOutput)
	})

	t.Run("success", func(t *testing.T) {
		// Create a valid application
		validAdminID := entities.StandardID{
			ID: uuid.New().String(),
		}
		createAdminInput := admin.CreateAdminInput{
			StandardID: validAdminID,
		}
		createAdminOutput, createAdminErr := app.AdminUseCase.CreateAdmin(ctx, createAdminInput)
		require.NoError(t, createAdminErr)

		// Edit the application
		newDescription := "this is a new description"
		editAdminInput := admin.EditAdminInput{
			AdminEditable: admin.AdminEditable{
				StandardResourceMeta: entities.StandardResourceMeta{
					StandardResource: entities.StandardResource{
						StandardID: createAdminOutput.StandardID,
					},
					ResourceVersion: createAdminOutput.ResourceVersion,
				},
				Description: &newDescription,
			},
		}
		editAdminOutput, err := app.AdminUseCase.EditAdmin(ctx, editAdminInput)
		require.NoError(t, err)
		require.NotNil(t, editAdminOutput)
		require.Equal(t, newDescription, *editAdminOutput.Description)

		// Get application from storage
		getAdminInput := admin.GetAdminInput{
			StandardID: validAdminID,
		}
		getAdminOutput, err := app.AdminUseCase.GetAdmin(ctx, getAdminInput)
		require.NoError(t, err)
		require.NotNil(t, getAdminOutput)
		require.Equal(t, newDescription, *getAdminOutput.Description)
	})
}

func TestDefaultUseCase_DeleteAdmin(t *testing.T) {
	ctx := context.Background()
	t.Run("failure: invalid input arguments", func(t *testing.T) {
		deleteAdminInput := admin.DeleteAdminInput{
			StandardID: entities.StandardID{
				ID: "",
			},
		}
		deleteAdminOutput, deleteAdminErr := app.AdminUseCase.DeleteAdmin(ctx, deleteAdminInput)
		require.Error(t, deleteAdminErr)
		require.True(t, errors.IsInvalidArgument(deleteAdminErr))
		require.Nil(t, deleteAdminOutput)
	})

	t.Run("failure: cannot delete admin if there is only one created", func(t *testing.T) {
		// Create a valid application
		adminID := entities.StandardID{
			ID: uuid.New().String(),
		}
		createAdminInput := admin.CreateAdminInput{
			StandardID: adminID,
		}
		_, createAdminErr := app.AdminUseCase.CreateAdmin(ctx, createAdminInput)
		require.NoError(t, createAdminErr)

		listAdminsInput := admin.ListAdminsInput{}
		listAdminsOutput, listAdminsErr := app.AdminUseCase.ListAdmins(ctx, listAdminsInput)
		require.NoError(t, listAdminsErr)
		require.NotNil(t, listAdminsOutput)

		for _, item := range listAdminsOutput.Items {
			deleteAdminInput := admin.DeleteAdminInput{
				StandardID: entities.StandardID{
					ID: item.ID,
				},
			}
			deleteAdminOutput, deleteAdminErr := app.AdminUseCase.DeleteAdmin(ctx, deleteAdminInput)
			if deleteAdminErr != nil {
				require.Error(t, deleteAdminErr)
				require.True(t, errors.IsPreconditionFailed(deleteAdminErr))
				require.Nil(t, deleteAdminOutput)
			}
		}
	})

	t.Run("failure: admin not found", func(t *testing.T) {
		// Create a valid application
		adminID := entities.StandardID{
			ID: uuid.New().String(),
		}
		createAdminInput := admin.CreateAdminInput{
			StandardID: adminID,
		}
		_, createAdminErr := app.AdminUseCase.CreateAdmin(ctx, createAdminInput)
		require.NoError(t, createAdminErr)
		adminID = entities.StandardID{
			ID: uuid.New().String(),
		}
		createAdminInput = admin.CreateAdminInput{
			StandardID: adminID,
		}
		_, createAdminErr = app.AdminUseCase.CreateAdmin(ctx, createAdminInput)
		require.NoError(t, createAdminErr)

		deleteAdminInput := admin.DeleteAdminInput{
			StandardID: entities.StandardID{
				ID: "id-do-not-exist",
			},
		}
		deleteAdminOutput, deleteAdminErr := app.AdminUseCase.DeleteAdmin(ctx, deleteAdminInput)
		require.Error(t, deleteAdminErr)
		require.True(t, errors.IsNotFound(deleteAdminErr))
		require.Nil(t, deleteAdminOutput)
	})

	t.Run("success", func(t *testing.T) {
		// Create a valid application
		initialAdminID := entities.StandardID{
			ID: uuid.New().String(),
		}
		createInitialAdminInput := admin.CreateAdminInput{
			StandardID: initialAdminID,
		}
		_, createInitialAdminErr := app.AdminUseCase.CreateAdmin(ctx, createInitialAdminInput)

		require.NoError(t, createInitialAdminErr)
		adminID := entities.StandardID{
			ID: uuid.New().String(),
		}
		createAdminInput := admin.CreateAdminInput{
			StandardID: adminID,
		}
		createAdminOutput, createAdminErr := app.AdminUseCase.CreateAdmin(ctx, createAdminInput)
		require.NoError(t, createAdminErr)

		// Delete the application
		deleteAdminInput := admin.DeleteAdminInput{
			StandardID: entities.StandardID{
				ID: createAdminOutput.ID,
			},
		}
		deleteAdminOutput, deleteAdminErr := app.AdminUseCase.DeleteAdmin(ctx, deleteAdminInput)
		require.NoError(t, deleteAdminErr)
		require.NotNil(t, deleteAdminOutput)
		require.Equal(t, createAdminOutput.ID, deleteAdminOutput.ID)

		getAdminInput := admin.GetAdminInput{
			StandardID: adminID,
		}
		// Retrieve deleted application
		getAdminOutput, getAdminErr := app.AdminUseCase.GetAdmin(ctx, getAdminInput)
		require.Error(t, getAdminErr)
		require.True(t, errors.IsNotFound(getAdminErr))
		require.Nil(t, getAdminOutput)
	})
}
