package user_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hyperledger-labs/signare/app/pkg/adapters/storage/postgres/userdbout"
	"github.com/hyperledger-labs/signare/app/pkg/commons/time"
	"github.com/hyperledger-labs/signare/app/pkg/commons/validators"
	"github.com/hyperledger-labs/signare/app/pkg/entities"
	"github.com/hyperledger-labs/signare/app/pkg/entities/address"
	"github.com/hyperledger-labs/signare/app/pkg/graph"
	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/application"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/hsmconnection"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/hsmconnector"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/hsmmodule"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/hsmslot"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/user"
	"github.com/hyperledger-labs/signare/app/test/dbtesthelper"
	"github.com/hyperledger-labs/signare/app/test/signaturemanagertesthelper"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var (
	slotID string

	chainID          = entities.NewInt256FromInt(44844)
	slotPin          = signaturemanagertesthelper.SlotPin
	hsmLoadedAddress = signaturemanagertesthelper.ImportedKeyAddress

	app graph.GraphShared
)

func TestMain(m *testing.M) {
	initializedSlotID, _, err := signaturemanagertesthelper.InitializeSoftHSMSlot()
	if err != nil {
		panic(err)
	}
	slotID = *initializedSlotID

	testApp, err := dbtesthelper.InitializeApp()
	if err != nil {
		panic(err)
	}
	app = *testApp

	validators.SetValidators()
	os.Exit(m.Run())
}

func TestProvideDefaultUseCase(t *testing.T) {
	t.Run("nil storage", func(t *testing.T) {
		userUseCase, err := user.ProvideDefaultUseCase(user.DefaultUserUseCaseOptions{
			Storage:            nil,
			ApplicationUseCase: &application.DefaultUseCase{},
		})
		require.Error(t, err)
		require.Nil(t, userUseCase)
	})

	t.Run("nil application use case", func(t *testing.T) {
		userUseCase, err := user.ProvideDefaultUseCase(user.DefaultUserUseCaseOptions{
			Storage:            &userdbout.Repository{},
			ApplicationUseCase: nil,
		})
		require.Error(t, err)
		require.Nil(t, userUseCase)
	})
}

func TestDefaultUseCase_CreateUser(t *testing.T) {
	ctx := context.Background()

	applicationID := uuid.NewString()
	description := "application for CreateUser test"
	createApplicationInput := application.CreateApplicationInput{
		ID:          &applicationID,
		ChainID:     *chainID,
		Description: &description,
	}
	createdApplication, createApplicationErr := app.ApplicationUseCase.CreateApplication(ctx, createApplicationInput)
	require.NoError(t, createApplicationErr)
	require.NotNil(t, createdApplication)

	t.Run("failure: invalid input arguments", func(t *testing.T) {
		input := user.CreateUserInput{
			ApplicationID: "",
			Roles: []string{
				"application-admin",
			},
		}
		output, err := app.UserUseCase.CreateUser(ctx, input)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
		require.Nil(t, output)

		input = user.CreateUserInput{
			ApplicationID: createdApplication.ID,
			Roles:         []string{},
		}
		output, err = app.UserUseCase.CreateUser(ctx, input)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
		require.Nil(t, output)
	})

	t.Run("success: if ID is not provided a random one is generated", func(t *testing.T) {
		input := user.CreateUserInput{
			ApplicationID: createdApplication.ID,
			Roles: []string{
				"application-admin",
			},
		}
		output, err := app.UserUseCase.CreateUser(ctx, input)
		require.NoError(t, err)
		require.NotNil(t, output)
		require.NotEmpty(t, output.ID)
	})

	t.Run("success", func(t *testing.T) {
		userID := uuid.New().String()
		input := user.CreateUserInput{
			ID:            &userID,
			ApplicationID: createdApplication.ID,
			Roles: []string{
				"application-admin",
			},
		}
		output, err := app.UserUseCase.CreateUser(ctx, input)
		require.NoError(t, err)
		require.NotNil(t, output)
		require.NotEmpty(t, output.InternalResourceID)
		require.Equal(t, userID, output.ID)
	})

	t.Run("failure: invalid role", func(t *testing.T) {
		userID := uuid.New().String()
		input := user.CreateUserInput{
			ID:            &userID,
			ApplicationID: createdApplication.ID,
			Roles: []string{
				"invalid-role",
			},
		}
		output, err := app.UserUseCase.CreateUser(ctx, input)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
		require.Nil(t, output)
	})

	t.Run("failure: already exists", func(t *testing.T) {
		userID := uuid.New().String()
		input := user.CreateUserInput{
			ID:            &userID,
			ApplicationID: createdApplication.ID,
			Roles: []string{
				"application-admin",
			},
		}
		_, err := app.UserUseCase.CreateUser(ctx, input)
		require.NoError(t, err)

		// Duplicated entry
		output, err := app.UserUseCase.CreateUser(ctx, input)
		require.Error(t, err)
		require.True(t, errors.IsAlreadyExists(err))
		require.Nil(t, output)
	})
}

func TestDefaultUseCase_ListUsers(t *testing.T) {
	ctx := context.Background()

	application1ID := uuid.NewString()
	description := "application for CreateUser test"
	createApplication1Input := application.CreateApplicationInput{
		ID:          &application1ID,
		ChainID:     *chainID,
		Description: &description,
	}
	createdApplication1, createApplication1Err := app.ApplicationUseCase.CreateApplication(ctx, createApplication1Input)
	require.NoError(t, createApplication1Err)
	require.NotNil(t, createdApplication1)

	usersToCreate := 20
	// Users for application-1
	for i := 1; i <= usersToCreate; i++ {
		userID := fmt.Sprint("user-", i)
		testDesc := "my-description"
		input := user.CreateUserInput{
			ID:            &userID,
			ApplicationID: createdApplication1.ID,
			Description:   &testDesc,
			Roles: []string{
				"application-admin",
			},
		}
		_, err := app.UserUseCase.CreateUser(ctx, input)
		require.NoError(t, err)
	}

	application2ID := uuid.NewString()
	createApplication2Input := application.CreateApplicationInput{
		ID:          &application2ID,
		ChainID:     *chainID,
		Description: &description,
	}
	createdApplication2, createApplication2Err := app.ApplicationUseCase.CreateApplication(ctx, createApplication2Input)
	require.NoError(t, createApplication2Err)
	require.NotNil(t, createdApplication2)

	// Users for application-2
	for i := 1; i <= usersToCreate; i++ {
		userID := fmt.Sprint("user-", i)
		testDesc := "my-description"
		input := user.CreateUserInput{
			ID:            &userID,
			ApplicationID: createdApplication2.ID,
			Description:   &testDesc,
			Roles: []string{
				"application-admin",
			},
		}
		_, err := app.UserUseCase.CreateUser(ctx, input)
		require.NoError(t, err)
	}

	t.Run("failure: invalid arguments", func(t *testing.T) {
		input := user.ListUsersInput{
			ApplicationID: "",
		}
		output, err := app.UserUseCase.ListUsers(ctx, input)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
		require.Nil(t, output)

		input = user.ListUsersInput{
			ApplicationID: "test-app",
			PageLimit:     -50,
			PageOffset:    -30,
		}
		output, err = app.UserUseCase.ListUsers(ctx, input)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
		require.Nil(t, output)
	})

	t.Run("success: list all users without limit", func(t *testing.T) {
		input := user.ListUsersInput{
			ApplicationID: createdApplication1.ID,
		}
		output, err := app.UserUseCase.ListUsers(ctx, input)
		require.NoError(t, err)
		require.Len(t, output.Items, usersToCreate)
	})

	t.Run("success: list DESC with limit", func(t *testing.T) {
		desiredLimit := 5
		input := user.ListUsersInput{
			ApplicationID:  createdApplication2.ID,
			OrderBy:        "creationDate",
			OrderDirection: "desc",
			PageLimit:      desiredLimit,
		}
		output, err := app.UserUseCase.ListUsers(ctx, input)
		require.NoError(t, err)
		require.Len(t, output.Items, desiredLimit)
		require.True(t, output.MoreItems)
		require.Equal(t, "user-20", output.Items[0].ID)
		require.Equal(t, "user-19", output.Items[1].ID)
		require.Equal(t, "user-18", output.Items[2].ID)
		require.Equal(t, "user-17", output.Items[3].ID)
		require.Equal(t, "user-16", output.Items[4].ID)
		// Assert order
		for i := 1; i < len(output.Items); i++ {
			require.Greater(t, output.Items[i-1].CreationDate.ToInt64(), output.Items[i].CreationDate.ToInt64())
		}
	})

	t.Run("success: list ASC with limit", func(t *testing.T) {
		desiredLimit := 5
		input := user.ListUsersInput{
			ApplicationID:  createdApplication2.ID,
			OrderBy:        "lastUpdate",
			OrderDirection: "asc",
			PageLimit:      desiredLimit,
		}
		output, err := app.UserUseCase.ListUsers(ctx, input)
		require.NoError(t, err)
		require.Len(t, output.Items, desiredLimit)
		require.True(t, output.MoreItems)
		require.Equal(t, "user-1", output.Items[0].ID)
		require.Equal(t, "user-2", output.Items[1].ID)
		require.Equal(t, "user-3", output.Items[2].ID)
		require.Equal(t, "user-4", output.Items[3].ID)
		require.Equal(t, "user-5", output.Items[4].ID)
		// Assert order
		for i := 1; i < len(output.Items); i++ {
			require.Less(t, output.Items[i-1].LastUpdate.ToInt64(), output.Items[i].LastUpdate.ToInt64())
		}
	})
}

func TestDefaultUseCase_GetUser(t *testing.T) {
	ctx := context.Background()

	applicationID := uuid.NewString()
	description := "application for CreateUser test"
	createApplicationInput := application.CreateApplicationInput{
		ID:          &applicationID,
		ChainID:     *chainID,
		Description: &description,
	}
	createdApplication, createApplicationErr := app.ApplicationUseCase.CreateApplication(ctx, createApplicationInput)
	require.NoError(t, createApplicationErr)
	require.NotNil(t, createdApplication)

	t.Run("failure: invalid arguments", func(t *testing.T) {
		input := user.GetUserInput{
			ApplicationStandardID: entities.ApplicationStandardID{
				ID:            "",
				ApplicationID: createdApplication.ID,
			},
		}
		output, err := app.UserUseCase.GetUser(ctx, input)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
		require.Nil(t, output)

		input = user.GetUserInput{
			ApplicationStandardID: entities.ApplicationStandardID{
				ID:            "my-user",
				ApplicationID: "",
			},
		}
		output, err = app.UserUseCase.GetUser(ctx, input)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
		require.Nil(t, output)
	})

	t.Run("failure: user not found", func(t *testing.T) {
		input := user.GetUserInput{
			ApplicationStandardID: entities.ApplicationStandardID{
				ID:            "my-id",
				ApplicationID: createdApplication.ID,
			},
		}
		output, err := app.UserUseCase.GetUser(ctx, input)
		require.Error(t, err)
		require.True(t, errors.IsNotFound(err))
		require.Nil(t, output)
	})

	t.Run("success", func(t *testing.T) {
		// Create user before retrieving it
		userID := uuid.New().String()
		testDesc := "my-description"
		createInput := user.CreateUserInput{
			ID:            &userID,
			ApplicationID: createdApplication.ID,
			Description:   &testDesc,
			Roles: []string{
				"application-admin",
			},
		}
		_, err := app.UserUseCase.CreateUser(ctx, createInput)
		require.NoError(t, err)

		getInput := user.GetUserInput{
			ApplicationStandardID: entities.ApplicationStandardID{
				ID:            userID,
				ApplicationID: createdApplication.ID,
			},
		}
		output, err := app.UserUseCase.GetUser(ctx, getInput)
		require.NoError(t, err)
		require.NotNil(t, output)
		require.Equal(t, createInput.ApplicationID, output.ApplicationID)
		require.Equal(t, createInput.Description, output.Description)
		require.Len(t, output.Roles, 1)
		require.Equal(t, createInput.Roles[0], output.Roles[0])
		require.Zero(t, len(output.Accounts))
		require.NotEmpty(t, output.InternalResourceID)
	})
}

func TestDefaultUseCase_EditUser(t *testing.T) {
	ctx := context.Background()

	applicationID := uuid.NewString()
	description := "application for CreateUser test"
	createApplicationInput := application.CreateApplicationInput{
		ID:          &applicationID,
		ChainID:     *chainID,
		Description: &description,
	}
	createdApplication, createApplicationErr := app.ApplicationUseCase.CreateApplication(ctx, createApplicationInput)
	require.NoError(t, createApplicationErr)
	require.NotNil(t, createdApplication)

	t.Run("failure: invalid arguments", func(t *testing.T) {
		input := user.EditUserInput{
			ApplicationStandardID: entities.ApplicationStandardID{
				ID:            "my-id",
				ApplicationID: "",
			},
			ResourceVersion: "",
			Roles:           []string{"application-admin"},
		}
		output, err := app.UserUseCase.EditUser(ctx, input)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
		require.Nil(t, output)

		input = user.EditUserInput{
			ApplicationStandardID: entities.ApplicationStandardID{
				ID:            "",
				ApplicationID: createdApplication.ID,
			},
			ResourceVersion: "",
		}
		output, err = app.UserUseCase.EditUser(ctx, input)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
		require.Nil(t, output)

		input = user.EditUserInput{
			ApplicationStandardID: entities.ApplicationStandardID{
				ID:            "",
				ApplicationID: "",
			},
			ResourceVersion: "the-resource-version",
			Roles:           []string{"application-admin"},
		}
		output, err = app.UserUseCase.EditUser(ctx, input)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
		require.Nil(t, output)
	})

	t.Run("failure: user not found", func(t *testing.T) {
		input := user.EditUserInput{
			ApplicationStandardID: entities.ApplicationStandardID{
				ID:            uuid.New().String(),
				ApplicationID: createdApplication.ID,
			},
			ResourceVersion: "the-resource-version",
			Roles:           []string{"application-admin"},
		}
		output, err := app.UserUseCase.EditUser(ctx, input)
		require.Error(t, err)
		require.True(t, errors.IsNotFound(err))
		require.Nil(t, output)
	})

	t.Run("failure: invalid resource version", func(t *testing.T) {
		// Create the User first
		userID := uuid.New().String()
		testDesc := "my-description"
		createInput := user.CreateUserInput{
			ID:            &userID,
			ApplicationID: createdApplication.ID,
			Description:   &testDesc,
			Roles: []string{
				"application-admin",
			},
		}
		_, err := app.UserUseCase.CreateUser(ctx, createInput)
		require.NoError(t, err)

		input := user.EditUserInput{
			ApplicationStandardID: entities.ApplicationStandardID{
				ID:            *createInput.ID,
				ApplicationID: createInput.ApplicationID,
			},
			ResourceVersion: "invalid-resource-version",
			Roles:           []string{"application-admin"},
		}
		output, err := app.UserUseCase.EditUser(ctx, input)
		require.Error(t, err)
		require.True(t, errors.IsNotFound(err))
		require.Nil(t, output)
	})

	t.Run("failure: role is mandatory", func(t *testing.T) {
		// Create the User first
		userID := uuid.New().String()
		testDesc := "my-description"
		createInput := user.CreateUserInput{
			ID:            &userID,
			ApplicationID: createdApplication.ID,
			Description:   &testDesc,
			Roles: []string{
				"application-admin",
			},
		}
		_, err := app.UserUseCase.CreateUser(ctx, createInput)
		require.NoError(t, err)

		input := user.EditUserInput{
			ApplicationStandardID: entities.ApplicationStandardID{
				ID:            *createInput.ID,
				ApplicationID: createInput.ApplicationID,
			},
			ResourceVersion: "invalid-resource-version",
		}
		output, err := app.UserUseCase.EditUser(ctx, input)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
		require.Nil(t, output)
	})

	t.Run("failure: invalid role", func(t *testing.T) {
		// Create the User first
		userID := uuid.New().String()
		testDesc := "my-description"
		createInput := user.CreateUserInput{
			ID:            &userID,
			ApplicationID: createdApplication.ID,
			Description:   &testDesc,
			Roles: []string{
				"application-admin",
			},
		}
		_, err := app.UserUseCase.CreateUser(ctx, createInput)
		require.NoError(t, err)

		input := user.EditUserInput{
			ApplicationStandardID: entities.ApplicationStandardID{
				ID:            *createInput.ID,
				ApplicationID: createInput.ApplicationID,
			},
			ResourceVersion: "invalid-resource-version",
			Roles:           []string{"not-supported-role"},
		}
		output, err := app.UserUseCase.EditUser(ctx, input)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
		require.Nil(t, output)
	})

	t.Run("success", func(t *testing.T) {
		// Create the User first
		userID := uuid.New().String()
		testDesc := "my-description"
		createInput := user.CreateUserInput{
			ID:            &userID,
			ApplicationID: createdApplication.ID,
			Description:   &testDesc,
			Roles: []string{
				"application-admin",
			},
		}
		createdUser, err := app.UserUseCase.CreateUser(ctx, createInput)
		require.NoError(t, err)

		// Edit the created user
		newDescription := "this is the new description"
		newRoles := []string{
			"application-admin",
			"transaction-signer",
		}
		input := user.EditUserInput{
			ApplicationStandardID: entities.ApplicationStandardID{
				ID:            createdUser.ID,
				ApplicationID: createdUser.ApplicationID,
			},
			ResourceVersion: createdUser.ResourceVersion,
			Description:     &newDescription,
			Roles:           newRoles,
		}
		editedUser, err := app.UserUseCase.EditUser(ctx, input)
		require.NoError(t, err)
		require.NotNil(t, editedUser)
		require.Equal(t, newDescription, *editedUser.Description)
		require.Equal(t, newRoles, editedUser.Roles)
		require.NotEqual(t, createdUser.ResourceVersion, editedUser.ResourceVersion)
		require.Equal(t, createdUser.CreationDate, editedUser.CreationDate)
		require.NotEqual(t, createdUser.LastUpdate, editedUser.LastUpdate)
	})
}

func TestDefaultUseCase_DeleteUser(t *testing.T) {
	ctx := context.Background()

	applicationID := uuid.NewString()
	description := "application for CreateUser test"
	createApplicationInput := application.CreateApplicationInput{
		ID:          &applicationID,
		ChainID:     *chainID,
		Description: &description,
	}
	createdApplication, createApplicationErr := app.ApplicationUseCase.CreateApplication(ctx, createApplicationInput)
	require.NoError(t, createApplicationErr)
	require.NotNil(t, createdApplication)

	t.Run("failure: invalid input arguments", func(t *testing.T) {
		deleteInput := user.DeleteUserInput{
			ApplicationStandardID: entities.ApplicationStandardID{
				ID:            "my-id",
				ApplicationID: "",
			},
		}
		output, err := app.UserUseCase.DeleteUser(ctx, deleteInput)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
		require.Nil(t, output)

		deleteInput = user.DeleteUserInput{
			ApplicationStandardID: entities.ApplicationStandardID{
				ID:            "",
				ApplicationID: createdApplication.ID,
			},
		}
		output, err = app.UserUseCase.DeleteUser(ctx, deleteInput)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
		require.Nil(t, output)
	})

	t.Run("failure: nonexistent user", func(t *testing.T) {
		deleteInput := user.DeleteUserInput{
			ApplicationStandardID: entities.ApplicationStandardID{
				ID:            uuid.New().String(),
				ApplicationID: createdApplication.ID,
			},
		}
		output, err := app.UserUseCase.DeleteUser(ctx, deleteInput)
		require.Error(t, err)
		require.True(t, errors.IsNotFound(err))
		require.Nil(t, output)
	})

	t.Run("success", func(t *testing.T) {
		// Create a valid user
		userID := uuid.New().String()
		createInput := user.CreateUserInput{
			ID:            &userID,
			ApplicationID: createdApplication.ID,
			Roles: []string{
				"application-admin",
			},
		}
		createdUser, err := app.UserUseCase.CreateUser(ctx, createInput)
		require.NoError(t, err)

		// Delete the user
		deleteInput := user.DeleteUserInput{
			ApplicationStandardID: entities.ApplicationStandardID{
				ID:            userID,
				ApplicationID: createdApplication.ID,
			},
		}
		deletedUser, err := app.UserUseCase.DeleteUser(ctx, deleteInput)
		require.NoError(t, err)
		require.NotNil(t, deletedUser)
		require.Equal(t, createdUser.ID, deletedUser.ID)
		require.Equal(t, createdUser.InternalResourceID, deletedUser.InternalResourceID)
		require.Equal(t, createdUser.ApplicationID, deletedUser.ApplicationID)
		require.Equal(t, createdUser.ResourceVersion, deletedUser.ResourceVersion)
		require.Equal(t, createdUser.Roles, deletedUser.Roles)
		require.Equal(t, createdUser.Description, deletedUser.Description)
		require.Len(t, deletedUser.Accounts, 0)

		// Retrieve deleted user
		getOutput, err := app.UserUseCase.GetUser(ctx, user.GetUserInput{
			ApplicationStandardID: entities.ApplicationStandardID{
				ID:            userID,
				ApplicationID: createdApplication.ID,
			},
		})
		require.Error(t, err)
		require.True(t, errors.IsNotFound(err))
		require.Nil(t, getOutput)
	})
}

func TestDefaultUseCase_AddUserAccounts(t *testing.T) {
	ctx := context.Background()

	applicationID := uuid.NewString()
	description := "application for CreateUser test"
	createApplicationInput := application.CreateApplicationInput{
		ID:          &applicationID,
		ChainID:     *chainID,
		Description: &description,
	}
	createApplicationOutput, createApplicationErr := app.ApplicationUseCase.CreateApplication(ctx, createApplicationInput)
	require.NoError(t, createApplicationErr)
	require.NotNil(t, createApplicationOutput)
	currentTestApplication := createApplicationOutput

	hsm := createHSM(ctx, t)
	// Create a slot within the HSM
	createHSMSlotInput := hsmslot.CreateHSMSlotInput{
		ApplicationID: currentTestApplication.ID,
		HSMModuleID:   hsm.ID,
		Slot:          slotID,
		Pin:           slotPin,
	}
	createHSMSlotOutput, createHSMSlotErr := app.HSMSlotUseCase.CreateHSMSlot(ctx, createHSMSlotInput)
	require.Nil(t, createHSMSlotErr)
	require.NotNil(t, createHSMSlotOutput)

	t.Run("failure: invalid input arguments", func(t *testing.T) {
		input := user.EnableAccountsInput{
			UserID:        "",
			ApplicationID: currentTestApplication.ID,
			Addresses: []address.Address{
				address.MustNewFromHexString("0xDc611d30c81e723D0A78BE33f5aF3974c108f5cf"),
			},
		}
		output, err := app.UserUseCase.EnableAccounts(ctx, input)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
		require.Nil(t, output)

		input = user.EnableAccountsInput{
			UserID:        "my-user-id",
			ApplicationID: "",
			Addresses: []address.Address{
				address.MustNewFromHexString("0xDc611d30c81e723D0A78BE33f5aF3974c108f5cf"),
			},
		}
		output, err = app.UserUseCase.EnableAccounts(ctx, input)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
		require.Nil(t, output)

		input = user.EnableAccountsInput{
			UserID:        "my-user-id",
			ApplicationID: currentTestApplication.ID,
			Addresses: []address.Address{
				address.MustNewFromHexString("invalid address"),
			},
		}
		output, err = app.UserUseCase.EnableAccounts(ctx, input)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
		require.Nil(t, output)
	})

	t.Run("failure: user not found", func(t *testing.T) {
		input := user.EnableAccountsInput{
			UserID:        uuid.New().String(),
			ApplicationID: currentTestApplication.ID,
			Addresses: []address.Address{
				address.MustNewFromHexString("0xDc611d30c81e723D0A78BE33f5aF3974c108f5cf"),
			},
		}
		output, err := app.UserUseCase.EnableAccounts(ctx, input)
		require.Error(t, err)
		require.True(t, errors.IsNotFound(err))
		require.Nil(t, output)
	})

	t.Run("failure: user's application not found", func(t *testing.T) {
		// Create a valid User
		userID := uuid.New().String()
		createUserInput := user.CreateUserInput{
			ID:            &userID,
			ApplicationID: currentTestApplication.ID,
			Roles: []string{
				"application-admin",
			},
		}
		_, err := app.UserUseCase.CreateUser(ctx, createUserInput)
		require.NoError(t, err)

		// Add user account
		input := user.EnableAccountsInput{
			UserID:        userID,
			ApplicationID: "non-existent-application",
			Addresses: []address.Address{
				address.MustNewFromHexString("0xDc611d30c81e723D0A78BE33f5aF3974c108f5cf"),
			},
		}
		output, err := app.UserUseCase.EnableAccounts(ctx, input)
		require.Error(t, err)
		require.True(t, errors.IsNotFound(err))
		require.Nil(t, output)
	})

	t.Run("success: adding duplicated accounts results in no error but no accounts persisted", func(t *testing.T) {
		// Create a valid User
		userID := uuid.New().String()
		createUserInput := user.CreateUserInput{
			ID:            &userID,
			ApplicationID: currentTestApplication.ID,
			Roles: []string{
				"application-admin",
			},
		}
		_, err := app.UserUseCase.CreateUser(ctx, createUserInput)
		require.NoError(t, err)

		// Add invalid (duplicated) accounts
		byApplicationInput := hsmconnection.ByApplicationInput{
			ApplicationID: currentTestApplication.ID,
		}
		hsmConnection, byApplicationErr := app.HSMConnectionResolver.ByApplication(ctx, byApplicationInput)
		require.NoError(t, byApplicationErr)
		require.NotNil(t, hsmConnection)

		generateAddressInput := hsmconnector.GenerateAddressInput{
			SlotConnectionData: hsmconnector.SlotConnectionData{
				Slot:       hsmConnection.Slot,
				Pin:        hsmConnection.Pin,
				ChainID:    hsmConnection.ChainID,
				ModuleKind: hsmconnector.ModuleKind(hsmConnection.ModuleKind),
			},
		}
		generateAddressOneOutput, generateAddressOneErr := app.HSMConnector.GenerateAddress(ctx, generateAddressInput)
		require.NoError(t, generateAddressOneErr)
		require.NotNil(t, generateAddressOneOutput)

		input := user.EnableAccountsInput{
			UserID:        userID,
			ApplicationID: applicationID,
			Addresses: []address.Address{
				generateAddressOneOutput.Address,
				generateAddressOneOutput.Address,
				generateAddressOneOutput.Address,
			},
		}
		output, err := app.UserUseCase.EnableAccounts(ctx, input)
		require.NoError(t, err)
		require.NotNil(t, output)
		require.Len(t, output.Accounts, 1)
	})

	t.Run("success", func(t *testing.T) {
		// Create a valid User
		userID := uuid.New().String()
		createUserInput := user.CreateUserInput{
			ID:            &userID,
			ApplicationID: currentTestApplication.ID,
			Roles: []string{
				"application-admin",
			},
		}
		_, err := app.UserUseCase.CreateUser(ctx, createUserInput)
		require.NoError(t, err)

		// Add user accounts
		byApplicationInput := hsmconnection.ByApplicationInput{
			ApplicationID: currentTestApplication.ID,
		}
		hsmConnection, byApplicationErr := app.HSMConnectionResolver.ByApplication(ctx, byApplicationInput)
		require.NoError(t, byApplicationErr)
		require.NotNil(t, hsmConnection)

		generateAddressInput := hsmconnector.GenerateAddressInput{
			SlotConnectionData: hsmconnector.SlotConnectionData{
				Slot:       hsmConnection.Slot,
				Pin:        hsmConnection.Pin,
				ChainID:    hsmConnection.ChainID,
				ModuleKind: hsmconnector.ModuleKind(hsmConnection.ModuleKind),
			},
		}
		generateAddressOneOutput, generateAddressOneErr := app.HSMConnector.GenerateAddress(ctx, generateAddressInput)
		require.NoError(t, generateAddressOneErr)
		require.NotNil(t, generateAddressOneOutput)
		input := user.EnableAccountsInput{
			UserID:        userID,
			ApplicationID: currentTestApplication.ID,
			Addresses: []address.Address{
				generateAddressOneOutput.Address,
			},
		}
		output, err := app.UserUseCase.EnableAccounts(ctx, input)
		require.NoError(t, err)
		require.NotNil(t, output)
		require.Equal(t, userID, output.ID)
		require.Equal(t, applicationID, output.ApplicationID)
		require.Len(t, output.Accounts, 1)
		require.Equal(t, input.Addresses[0].String(), output.Accounts[0].Address.String())
	})
}

func TestDefaultUseCase_RemoveUserAccount(t *testing.T) {
	ctx := context.Background()

	applicationID := uuid.NewString()
	description := "application for CreateUser test"
	createApplicationInput := application.CreateApplicationInput{
		ID:          &applicationID,
		ChainID:     *chainID,
		Description: &description,
	}
	createdApplication, createApplicationErr := app.ApplicationUseCase.CreateApplication(ctx, createApplicationInput)
	require.NoError(t, createApplicationErr)
	require.NotNil(t, createdApplication)

	t.Run("failure: invalid input arguments", func(t *testing.T) {
		input := user.DisableAccountInput{
			UserID:        "",
			ApplicationID: createdApplication.ID,
			Address:       address.MustNewFromHexString("0xDc611d30c81e723D0A78BE33f5aF3974c108f5cf"),
		}
		output, err := app.UserUseCase.DisableAccount(ctx, input)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
		require.Nil(t, output)

		input = user.DisableAccountInput{
			UserID:        "my-user-id",
			ApplicationID: "",
			Address:       address.MustNewFromHexString("0xDc611d30c81e723D0A78BE33f5aF3974c108f5cf"),
		}
		output, err = app.UserUseCase.DisableAccount(ctx, input)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
		require.Nil(t, output)

		input = user.DisableAccountInput{
			UserID:        "my-user-id",
			ApplicationID: createdApplication.ID,
			Address:       address.MustNewFromHexString("invalid address"),
		}
		output, err = app.UserUseCase.DisableAccount(ctx, input)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
		require.Nil(t, output)
	})

	t.Run("failure: account not found", func(t *testing.T) {
		input := user.DisableAccountInput{
			UserID:        "my-user-id",
			ApplicationID: createdApplication.ID,
			Address:       address.MustNewFromHexString("0xDc611d30c81e723D0A78BE33f5aF3974c108f5cf"),
		}
		output, err := app.UserUseCase.DisableAccount(ctx, input)
		require.Error(t, err)
		require.True(t, errors.IsNotFound(err))
		require.Nil(t, output)
	})

	t.Run("success", func(t *testing.T) {
		// Create an HSM to create accounts
		hsm := createHSM(ctx, t)
		require.NotNil(t, hsm)

		// Create a slot within the HSM
		createHSMSlotInput := hsmslot.CreateHSMSlotInput{
			ApplicationID: createdApplication.ID,
			HSMModuleID:   hsm.ID,
			Slot:          slotID,
			Pin:           slotPin,
		}
		createHSMSlotOutput, createHSMSlotErr := app.HSMSlotUseCase.CreateHSMSlot(ctx, createHSMSlotInput)
		require.Nil(t, createHSMSlotErr)
		require.NotNil(t, createHSMSlotOutput)

		// Create a valid User
		userID := uuid.New().String()
		createUserInput := user.CreateUserInput{
			ID:            &userID,
			ApplicationID: createdApplication.ID,
			Roles: []string{
				"application-admin",
			},
		}
		_, err := app.UserUseCase.CreateUser(ctx, createUserInput)
		require.NoError(t, err)

		// Add user account
		byApplicationInput := hsmconnection.ByApplicationInput{
			ApplicationID: createdApplication.ID,
		}
		hsmConnection, byApplicationErr := app.HSMConnectionResolver.ByApplication(ctx, byApplicationInput)
		require.NoError(t, byApplicationErr)
		require.NotNil(t, hsmConnection)

		generateAddressInput := hsmconnector.GenerateAddressInput{
			SlotConnectionData: hsmconnector.SlotConnectionData{
				Slot:       hsmConnection.Slot,
				Pin:        hsmConnection.Pin,
				ChainID:    hsmConnection.ChainID,
				ModuleKind: hsmconnector.ModuleKind(hsmConnection.ModuleKind),
			},
		}
		generateAddressOneOutput, generateAddressOneErr := app.HSMConnector.GenerateAddress(ctx, generateAddressInput)
		require.NoError(t, generateAddressOneErr)
		require.NotNil(t, generateAddressOneOutput)
		input := user.EnableAccountsInput{
			UserID:        userID,
			ApplicationID: applicationID,
			Addresses: []address.Address{
				generateAddressOneOutput.Address,
			},
		}
		_, err = app.UserUseCase.EnableAccounts(ctx, input)
		require.NoError(t, err)

		removeAccountInput := user.DisableAccountInput{
			UserID:        userID,
			ApplicationID: applicationID,
			Address:       generateAddressOneOutput.Address,
		}
		output, err := app.UserUseCase.DisableAccount(ctx, removeAccountInput)
		require.NoError(t, err)
		require.NotNil(t, output)
		require.Len(t, output.Accounts, 0)
	})
}
func createHSM(ctx context.Context, t *testing.T) *hsmmodule.HSMModule {
	description := "HSM module for testing"
	hsmModule := hsmmodule.HSMModule{
		StandardResourceMeta: entities.StandardResourceMeta{
			StandardResource: entities.StandardResource{
				StandardID: entities.StandardID{
					ID: uuid.NewString(),
				},
				Timestamps: entities.Timestamps{
					CreationDate: time.Now(),
					LastUpdate:   time.Now(),
				},
			},
			ResourceVersion: uuid.NewString(),
		},
		Description: &description,
		Configuration: hsmmodule.HSMModuleConfiguration{
			SoftHSMConfiguration: &hsmmodule.SoftHSMConfiguration{},
		},
		Kind: hsmmodule.SoftHSMModuleKind,
	}
	createHSMModuleInput := hsmmodule.CreateHSMModuleInput{
		ID:            &hsmModule.ID,
		Description:   hsmModule.Description,
		Configuration: hsmModule.Configuration,
		ModuleKind:    hsmModule.Kind,
	}
	addedModule, err := app.HSMModuleUseCase.CreateHSMModule(ctx, createHSMModuleInput)
	if err != nil {
		require.True(t, errors.IsAlreadyExists(err))
		getHSMModuleInput := hsmmodule.GetHSMModuleInput{
			StandardID: entities.StandardID{
				ID: hsmModule.ID,
			},
		}
		module, errGet := app.HSMModuleUseCase.GetHSMModule(ctx, getHSMModuleInput)
		require.NoError(t, errGet)
		require.NotNil(t, module)
		return &module.HSMModule
	}
	require.NotNil(t, addedModule)
	return &addedModule.HSMModule
}
