package hsmmodule_test

import (
	"context"
	"os"
	"testing"

	"github.com/hyperledger-labs/signare/app/pkg/adapters/storage/postgres/hsmdbout"
	"github.com/hyperledger-labs/signare/app/pkg/commons/validators"
	"github.com/hyperledger-labs/signare/app/pkg/entities"
	"github.com/hyperledger-labs/signare/app/pkg/graph"
	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/application"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/hsmmodule"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/hsmslot"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/referentialintegrity"
	"github.com/hyperledger-labs/signare/app/test/dbtesthelper"
	"github.com/hyperledger-labs/signare/app/test/signaturemanagertesthelper"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var (
	slotID string

	slotPin = signaturemanagertesthelper.SlotPin
	app     graph.GraphShared
)

func TestMain(m *testing.M) {
	slot, _, err := signaturemanagertesthelper.InitializeSoftHSMSlot()
	if err != nil {
		panic(err)
	}
	slotID = *slot

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
		options := hsmmodule.DefaultUseCaseOptions{
			HSMModuleStorage:            nil,
			ReferentialIntegrityUseCase: &referentialintegrity.DefaultUseCase{},
		}
		defaultUseCase, err := hsmmodule.ProvideDefaultHSMModuleUseCase(options)
		require.Error(t, err)
		require.Nil(t, defaultUseCase)
	})

	t.Run("nil referential integrity use case", func(t *testing.T) {
		options := hsmmodule.DefaultUseCaseOptions{
			HSMModuleStorage:            &hsmdbout.Repository{},
			ReferentialIntegrityUseCase: nil,
		}
		defaultUseCase, err := hsmmodule.ProvideDefaultHSMModuleUseCase(options)
		require.Error(t, err)
		require.Nil(t, defaultUseCase)
	})

	t.Run("success", func(t *testing.T) {
		options := hsmmodule.DefaultUseCaseOptions{
			HSMModuleStorage:            &hsmdbout.Repository{},
			ReferentialIntegrityUseCase: &referentialintegrity.DefaultUseCase{},
		}
		defaultUseCase, err := hsmmodule.ProvideDefaultHSMModuleUseCase(options)
		require.NoError(t, err)
		require.NotNil(t, defaultUseCase)
	})
}

func TestDefaultUseCase_CreateHSM(t *testing.T) {
	ctx := context.Background()
	description := "HSM for testing"
	hsmID := "HSM-01"
	t.Run("success", func(t *testing.T) {
		createHSMInput := hsmmodule.CreateHSMModuleInput{
			ID:          &hsmID,
			Description: &description,
			Configuration: hsmmodule.HSMModuleConfiguration{
				SoftHSMConfiguration: &hsmmodule.SoftHSMConfiguration{},
			},
			ModuleKind: hsmmodule.SoftHSMModuleKind,
		}
		createdHSM, createHSMErr := app.HSMModuleUseCase.CreateHSMModule(ctx, createHSMInput)
		require.Nil(t, createHSMErr)
		require.NotEmpty(t, createdHSM.InternalResourceID)
		require.NotNil(t, createdHSM)
	})

	t.Run("failure: HSM already exists", func(t *testing.T) {
		createHSMInput := hsmmodule.CreateHSMModuleInput{
			ID:          &hsmID,
			Description: &description,
			Configuration: hsmmodule.HSMModuleConfiguration{
				SoftHSMConfiguration: &hsmmodule.SoftHSMConfiguration{},
			},
			ModuleKind: hsmmodule.SoftHSMModuleKind,
		}
		createHSMOutput, createHSMErr := app.HSMModuleUseCase.CreateHSMModule(ctx, createHSMInput)
		require.NotNil(t, createHSMErr)
		require.True(t, errors.IsAlreadyExists(createHSMErr))
		require.Nil(t, createHSMOutput)
	})
}

func TestDefaultUseCase_ListHSMs(t *testing.T) {
	ctx := context.Background()
	hsmsToCreate := 40
	description := "HSM for testing"
	createHSMInput := hsmmodule.CreateHSMModuleInput{
		Description: &description,
		Configuration: hsmmodule.HSMModuleConfiguration{
			SoftHSMConfiguration: &hsmmodule.SoftHSMConfiguration{},
		},
		ModuleKind: hsmmodule.SoftHSMModuleKind,
	}
	for i := 0; i < hsmsToCreate; i++ {
		createHSMOutput, createHSMErr := app.HSMModuleUseCase.CreateHSMModule(ctx, createHSMInput)
		require.Nil(t, createHSMErr)
		require.NotNil(t, createHSMOutput)
	}

	t.Run("failure: invalid input arguments", func(t *testing.T) {
		listHSMsInput := hsmmodule.ListHSMModulesInput{
			PageLimit:  -50,
			PageOffset: -30,
		}
		listHSMsOutput, listHSMsErr := app.HSMModuleUseCase.ListHSMModules(ctx, listHSMsInput)
		require.NotNil(t, listHSMsErr)
		require.Nil(t, listHSMsOutput)
	})

	t.Run("success: list all HSMs", func(t *testing.T) {
		listHSMsInput := hsmmodule.ListHSMModulesInput{}
		listHSMsOutput, listHSMsErr := app.HSMModuleUseCase.ListHSMModules(ctx, listHSMsInput)
		require.NoError(t, listHSMsErr)
		require.True(t, len(listHSMsOutput.Items) >= hsmsToCreate)
	})

	t.Run("success: list DESC with limit", func(t *testing.T) {
		desiredLimit := 5
		listHSMsInput := hsmmodule.ListHSMModulesInput{
			OrderBy:        "creationDate",
			OrderDirection: "desc",
			PageLimit:      desiredLimit,
		}
		listHSMsOutput, listHSMsErr := app.HSMModuleUseCase.ListHSMModules(ctx, listHSMsInput)
		require.NoError(t, listHSMsErr)
		require.Len(t, listHSMsOutput.Items, desiredLimit)
		require.True(t, listHSMsOutput.MoreItems)
		// Assert order
		for i := 1; i < len(listHSMsOutput.Items); i++ {
			require.GreaterOrEqual(t, listHSMsOutput.Items[i-1].CreationDate.ToInt64(), listHSMsOutput.Items[i].CreationDate.ToInt64())
		}
	})

	t.Run("success: list ASC with limit", func(t *testing.T) {
		desiredLimit := 5
		listHSMsInput := hsmmodule.ListHSMModulesInput{
			OrderBy:        "creationDate",
			OrderDirection: "asc",
			PageLimit:      desiredLimit,
		}
		listHSMsOutput, listHSMsErr := app.HSMModuleUseCase.ListHSMModules(ctx, listHSMsInput)
		require.NoError(t, listHSMsErr)
		require.Len(t, listHSMsOutput.Items, desiredLimit)
		require.True(t, listHSMsOutput.MoreItems)
		// Assert order
		for i := 1; i < len(listHSMsOutput.Items); i++ {
			require.Less(t, listHSMsOutput.Items[i-1].LastUpdate.ToInt64(), listHSMsOutput.Items[i].LastUpdate.ToInt64())
		}
	})
}

func TestDefaultUseCase_GetHSM(t *testing.T) {
	ctx := context.Background()
	t.Run("failure: invalid ID", func(t *testing.T) {
		getHSMInput := hsmmodule.GetHSMModuleInput{
			StandardID: entities.StandardID{
				ID: "",
			},
		}
		getHSMOutput, getHSMErr := app.HSMModuleUseCase.GetHSMModule(ctx, getHSMInput)
		require.Error(t, getHSMErr)
		require.True(t, errors.IsInvalidArgument(getHSMErr))
		require.Nil(t, getHSMOutput)

		getHSMInput = hsmmodule.GetHSMModuleInput{
			StandardID: entities.StandardID{
				ID: "a-very-long-id-0000000000000000000000000000000000000000000000000000000000000",
			},
		}
		getHSMOutput, getHSMErr = app.HSMModuleUseCase.GetHSMModule(ctx, getHSMInput)
		require.Error(t, getHSMErr)
		require.True(t, errors.IsInvalidArgument(getHSMErr))
		require.Nil(t, getHSMOutput)
	})

	t.Run("failure: HSM not found", func(t *testing.T) {
		getHSMInput := hsmmodule.GetHSMModuleInput{
			StandardID: entities.StandardID{
				ID: "this-id-does-not-belong-to-any-HSM",
			},
		}
		getHSMOutput, getHSMErr := app.HSMModuleUseCase.GetHSMModule(ctx, getHSMInput)
		require.Error(t, getHSMErr)
		require.True(t, errors.IsNotFound(getHSMErr))
		require.Nil(t, getHSMOutput)
	})

	t.Run("success", func(t *testing.T) {
		// Create HSM before retrieving it
		hsmID := entities.StandardID{
			ID: uuid.New().String(),
		}
		description := "HSM for testing"
		createHSMInput := hsmmodule.CreateHSMModuleInput{
			ID:          &hsmID.ID,
			Description: &description,
			Configuration: hsmmodule.HSMModuleConfiguration{
				SoftHSMConfiguration: &hsmmodule.SoftHSMConfiguration{},
			},
			ModuleKind: hsmmodule.SoftHSMModuleKind,
		}
		_, createHSMErr := app.HSMModuleUseCase.CreateHSMModule(ctx, createHSMInput)
		require.NoError(t, createHSMErr)

		getHSMInput := hsmmodule.GetHSMModuleInput{
			StandardID: hsmID,
		}
		getHSMOutput, err := app.HSMModuleUseCase.GetHSMModule(ctx, getHSMInput)
		require.NoError(t, err)
		require.NotNil(t, getHSMOutput)
		require.Equal(t, hsmID.ID, getHSMOutput.ID)
		require.NotEmpty(t, getHSMOutput.InternalResourceID)
	})
}

func TestDefaultUseCase_EditHSM(t *testing.T) {
	ctx := context.Background()
	t.Run("failure: invalid input arguments", func(t *testing.T) {
		editHSMInput := hsmmodule.EditHSMModuleInput{
			HSMModule: hsmmodule.HSMModule{
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
		editHSMOutput, editHSMErr := app.HSMModuleUseCase.EditHSMModule(ctx, editHSMInput)
		require.Error(t, editHSMErr)
		require.True(t, errors.IsInvalidArgument(editHSMErr))
		require.Nil(t, editHSMOutput)

		editHSMInput = hsmmodule.EditHSMModuleInput{
			HSMModule: hsmmodule.HSMModule{
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
		editHSMOutput, editHSMErr = app.HSMModuleUseCase.EditHSMModule(ctx, editHSMInput)
		require.Error(t, editHSMErr)
		require.True(t, errors.IsInvalidArgument(editHSMErr))
		require.Nil(t, editHSMOutput)
	})

	t.Run("failure: nonexistent HSM", func(t *testing.T) {
		editHSMInput := hsmmodule.EditHSMModuleInput{
			HSMModule: hsmmodule.HSMModule{
				StandardResourceMeta: entities.StandardResourceMeta{
					StandardResource: entities.StandardResource{
						StandardID: entities.StandardID{
							ID: "my-id",
						},
					},
					ResourceVersion: "my-resource-version",
				},
				Kind: hsmmodule.SoftHSMModuleKind,
				Configuration: hsmmodule.HSMModuleConfiguration{
					SoftHSMConfiguration: &hsmmodule.SoftHSMConfiguration{},
				},
			},
		}
		editHSMOutput, editHSMErr := app.HSMModuleUseCase.EditHSMModule(ctx, editHSMInput)
		require.Error(t, editHSMErr)
		require.True(t, errors.IsNotFound(editHSMErr))
		require.Nil(t, editHSMOutput)
	})

	t.Run("failure: invalid resource version", func(t *testing.T) {
		// Create a valid HSM
		description := "HSM for testing"
		validHSMID := entities.StandardID{
			ID: uuid.New().String(),
		}
		createHSMInput := hsmmodule.CreateHSMModuleInput{
			ID:          &validHSMID.ID,
			Description: &description,
			ModuleKind:  hsmmodule.SoftHSMModuleKind,
			Configuration: hsmmodule.HSMModuleConfiguration{
				SoftHSMConfiguration: &hsmmodule.SoftHSMConfiguration{},
			},
		}
		_, createHSMErr := app.HSMModuleUseCase.CreateHSMModule(ctx, createHSMInput)
		require.NoError(t, createHSMErr)

		// Edit the HSM with an invalid resource version
		editHSMInput := hsmmodule.EditHSMModuleInput{
			HSMModule: hsmmodule.HSMModule{
				StandardResourceMeta: entities.StandardResourceMeta{
					StandardResource: entities.StandardResource{
						StandardID: validHSMID,
					},
					ResourceVersion: "invalid-resource-version",
				},
				Configuration: hsmmodule.HSMModuleConfiguration{
					SoftHSMConfiguration: &hsmmodule.SoftHSMConfiguration{},
				},
				Kind: hsmmodule.SoftHSMModuleKind,
			},
		}
		editHSMOutput, editHSMErr := app.HSMModuleUseCase.EditHSMModule(ctx, editHSMInput)
		require.Error(t, editHSMErr)
		require.True(t, errors.IsNotFound(editHSMErr))
		require.Nil(t, editHSMOutput)
	})

	t.Run("success", func(t *testing.T) {
		// Create a valid HSM
		description := "HSM for testing"
		validHSMID := entities.StandardID{
			ID: uuid.New().String(),
		}
		createHSMInput := hsmmodule.CreateHSMModuleInput{
			ID:          &validHSMID.ID,
			Description: &description,
			Configuration: hsmmodule.HSMModuleConfiguration{
				SoftHSMConfiguration: &hsmmodule.SoftHSMConfiguration{},
			},
			ModuleKind: hsmmodule.SoftHSMModuleKind,
		}
		createHSMOutput, createHSMErr := app.HSMModuleUseCase.CreateHSMModule(ctx, createHSMInput)
		require.NoError(t, createHSMErr)

		// Edit the HSM
		newDescription := "this is a new description"
		editHSMInput := hsmmodule.EditHSMModuleInput{
			HSMModule: hsmmodule.HSMModule{
				StandardResourceMeta: entities.StandardResourceMeta{
					StandardResource: entities.StandardResource{
						StandardID: createHSMOutput.StandardID,
					},
					ResourceVersion: createHSMOutput.ResourceVersion,
				},
				Description: &newDescription,
				Configuration: hsmmodule.HSMModuleConfiguration{
					SoftHSMConfiguration: &hsmmodule.SoftHSMConfiguration{},
				},
				Kind: hsmmodule.SoftHSMModuleKind,
			},
		}
		editHSMOutput, editHSMErr := app.HSMModuleUseCase.EditHSMModule(ctx, editHSMInput)
		require.NoError(t, editHSMErr)
		require.NotNil(t, editHSMOutput)
		require.Equal(t, newDescription, *editHSMOutput.Description)

		// Get HSM from storage
		getHSMInput := hsmmodule.GetHSMModuleInput{
			StandardID: validHSMID,
		}
		getHSMOutput, editHSMErr := app.HSMModuleUseCase.GetHSMModule(ctx, getHSMInput)
		require.NoError(t, editHSMErr)
		require.NotNil(t, getHSMOutput)
		require.Equal(t, newDescription, *getHSMOutput.Description)
	})
}

func TestDefaultUseCase_DeleteHSM(t *testing.T) {
	ctx := context.Background()
	t.Run("failure: invalid input arguments", func(t *testing.T) {
		deleteHSMInput := hsmmodule.DeleteHSMModuleInput{
			StandardID: entities.StandardID{
				ID: "",
			},
		}
		deleteHSMOutput, deleteHSMErr := app.HSMModuleUseCase.DeleteHSMModule(ctx, deleteHSMInput)
		require.Error(t, deleteHSMErr)
		require.True(t, errors.IsInvalidArgument(deleteHSMErr))
		require.Nil(t, deleteHSMOutput)
	})

	t.Run("failure: cannot delete HSM if there are slots in use", func(t *testing.T) {
		// Create a valid HSM
		hsmID := entities.StandardID{
			ID: uuid.New().String(),
		}
		description := "HSM for testing"
		createHSMInput := hsmmodule.CreateHSMModuleInput{
			ID:          &hsmID.ID,
			Description: &description,
			Configuration: hsmmodule.HSMModuleConfiguration{
				SoftHSMConfiguration: &hsmmodule.SoftHSMConfiguration{},
			},
			ModuleKind: hsmmodule.SoftHSMModuleKind,
		}
		_, createHSMErr := app.HSMModuleUseCase.CreateHSMModule(ctx, createHSMInput)
		require.NoError(t, createHSMErr)

		// Create a valid application for the slot
		applicationID := "hsm_module_test-app"
		input := application.CreateApplicationInput{
			ID:          &applicationID,
			ChainID:     *entities.NewInt256FromInt(10),
			Description: &description,
		}
		output, err := app.ApplicationUseCase.CreateApplication(ctx, input)
		require.NoError(t, err)
		require.Equal(t, applicationID, output.ID)

		createHSMSlotInput := hsmslot.CreateHSMSlotInput{
			ApplicationID: applicationID,
			HSMModuleID:   hsmID.ID,
			Slot:          slotID,
			Pin:           slotPin,
		}
		_, createHSMSlotErr := app.HSMSlotUseCase.CreateHSMSlot(ctx, createHSMSlotInput)
		require.NoError(t, createHSMSlotErr)
		deleteHSMInput := hsmmodule.DeleteHSMModuleInput{
			StandardID: hsmID,
		}
		deleteHSMOutput, deleteHSMErr := app.HSMModuleUseCase.DeleteHSMModule(ctx, deleteHSMInput)
		if deleteHSMErr != nil {
			require.Error(t, deleteHSMErr)
			require.True(t, errors.IsPreconditionFailed(deleteHSMErr))
			require.Nil(t, deleteHSMOutput)
		}
	})

	t.Run("failure: HSM not found", func(t *testing.T) {
		// Create a valid HSM
		hsmID := entities.StandardID{
			ID: uuid.New().String(),
		}
		description := "HSM for testing"
		createHSMInput := hsmmodule.CreateHSMModuleInput{
			ID:          &hsmID.ID,
			Description: &description,
			Configuration: hsmmodule.HSMModuleConfiguration{
				SoftHSMConfiguration: &hsmmodule.SoftHSMConfiguration{},
			},
			ModuleKind: hsmmodule.SoftHSMModuleKind,
		}
		_, createHSMErr := app.HSMModuleUseCase.CreateHSMModule(ctx, createHSMInput)
		require.NoError(t, createHSMErr)
		hsmID = entities.StandardID{
			ID: uuid.New().String(),
		}
		createHSMInput = hsmmodule.CreateHSMModuleInput{
			ID:          &hsmID.ID,
			Description: &description,
			Configuration: hsmmodule.HSMModuleConfiguration{
				SoftHSMConfiguration: &hsmmodule.SoftHSMConfiguration{},
			},
			ModuleKind: hsmmodule.SoftHSMModuleKind,
		}
		_, createHSMErr = app.HSMModuleUseCase.CreateHSMModule(ctx, createHSMInput)
		require.NoError(t, createHSMErr)

		deleteHSMInput := hsmmodule.DeleteHSMModuleInput{
			StandardID: entities.StandardID{
				ID: "id-do-not-exist",
			},
		}
		deleteHSMOutput, deleteHSMErr := app.HSMModuleUseCase.DeleteHSMModule(ctx, deleteHSMInput)
		require.Error(t, deleteHSMErr)
		require.True(t, errors.IsNotFound(deleteHSMErr))
		require.Nil(t, deleteHSMOutput)
	})

	t.Run("success", func(t *testing.T) {
		// Create a valid HSM
		hsmID := entities.StandardID{
			ID: uuid.New().String(),
		}
		description := "HSM for testing"
		createHSMInput := hsmmodule.CreateHSMModuleInput{
			ID:          &hsmID.ID,
			Description: &description,
			Configuration: hsmmodule.HSMModuleConfiguration{
				SoftHSMConfiguration: &hsmmodule.SoftHSMConfiguration{},
			},
			ModuleKind: hsmmodule.SoftHSMModuleKind,
		}
		createHSMOutput, createHSMErr := app.HSMModuleUseCase.CreateHSMModule(ctx, createHSMInput)
		require.NoError(t, createHSMErr)

		// Delete the HSM
		deleteHSMInput := hsmmodule.DeleteHSMModuleInput{
			StandardID: entities.StandardID{
				ID: createHSMOutput.ID,
			},
		}
		deleteHSMOutput, deleteHSMErr := app.HSMModuleUseCase.DeleteHSMModule(ctx, deleteHSMInput)
		require.NoError(t, deleteHSMErr)
		require.NotNil(t, deleteHSMOutput)
		require.Equal(t, createHSMOutput.ID, deleteHSMOutput.ID)

		getHSMInput := hsmmodule.GetHSMModuleInput{
			StandardID: hsmID,
		}
		// Retrieve deleted HSM
		getHSMOutput, getHSMErr := app.HSMModuleUseCase.GetHSMModule(ctx, getHSMInput)
		require.Error(t, getHSMErr)
		require.True(t, errors.IsNotFound(getHSMErr))
		require.Nil(t, getHSMOutput)
	})
}
