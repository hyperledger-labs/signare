package user_test

import (
	"context"
	"testing"

	"github.com/hyperledger-labs/signare/app/pkg/entities/address"
	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/application"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/hsmmodule"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/hsmslot"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/user"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var (
	ctx       = context.Background()
	addresses = []string{
		"0xbB3fbc5Ab07D7866b62860aB540E5F581B8710E1",
		"0xbB3fbc5Ab07D7866b62860aB540E5F581B8710A1",
		"0xbB3fbc5Ab07D7866b62860aB540E5F581B8710B1",
		"0xbB3fbc5Ab07D7866b62860aB540E5F581B8710C1",
		"0xbB3fbc5Ab07D7866b62860aB540E5F581B8710D1",
		"0xbB3fbc5Ab17D7866b62860aB540E5F581B8710E1",
		"0xbB3fbc5Ab07D7866b62860aB540E5F581B8710F1",
		"0xbB3fbc5Ab07D7866b62860aB540E5F581B8710F0",
		"0xbB3fbc5Ab07D7866b62860aB540E5F581B8710F2",
		"0xbB3fbc5Ab07D7866b62860aB540E5F581B8710F3",
		"0xbB3fbc5Ab07D7866b62860aB540E5F581B8710F4",
		"0xbB3fbc5Ab07D7866b62860aB540E5F581B8710F5",
		"0xbB3fbc5Ab07D7866b62860aB540E5F581B8710F6",
		"0xbB3fbc5Ab07D7866b62860aB540E5F581B8710F7",
		"0xbB3fbc5Ab07D7866b62860aB540E5F581B8710F8",
		"0xbB3fbc5Ab07D7866b62860aB540E5F581B8710F9",
		"0xcB3fbc5Ab07D7866b62860aB540E5F581B8710F9",
		"0xdB3fbc5Ab07D7866b62860aB540E5F581B8710F9",
		"0xeB3fbc5Ab07D7866b62860aB540E5F581B8710F9",
		"0xfB3fbc5Ab07D7866b62860aB540E5F581B8710F9",
		"0xfB3fbc1Ab07D7866b62860aB540E5F581B8710F9",
		"0xfB3fbc2Ab07D7866b62860aB540E5F581B8710F9",
		"0xfB3fbc3Ab07D7866b62860aB540E5F581B8710F9",
	}
)

func TestDefaultUseCase_CreateAccount(t *testing.T) {
	applicationID := uuid.NewString()
	createApplicationInput := application.CreateApplicationInput{
		ID:      &applicationID,
		ChainID: *chainID,
	}
	createApplicationOutput, createApplicationErr := app.ApplicationUseCase.CreateApplication(ctx, createApplicationInput)
	require.NoError(t, createApplicationErr)
	require.NotNil(t, createApplicationOutput)

	userID := uuid.NewString()
	createUserInput := user.CreateUserInput{
		ID:            &userID,
		ApplicationID: applicationID,
		Roles:         []string{"application-admin"},
	}
	createUserOutput, createUserErr := app.UserUseCase.CreateUser(ctx, createUserInput)
	require.NoError(t, createUserErr)
	require.NotNil(t, createUserOutput)
	t.Run("failure: invalid arguments", func(t *testing.T) {
		input := user.CreateAccountInput{
			AccountID: user.AccountID{
				Address:       address.MustNewFromHexString("invalid-address"),
				UserID:        "test-user-id",
				ApplicationID: applicationID,
			},
		}
		acc, err := app.AccountUseCase.CreateAccount(ctx, input)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
		require.Nil(t, acc)

		input = user.CreateAccountInput{
			AccountID: user.AccountID{
				Address:       address.MustNewFromHexString("0xbB3fbc5Ab07D7866b62860aB540E5F581B8710E1"),
				UserID:        "",
				ApplicationID: applicationID,
			},
		}
		acc, err = app.AccountUseCase.CreateAccount(ctx, input)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
		require.Nil(t, acc)

		input = user.CreateAccountInput{
			AccountID: user.AccountID{
				Address:       address.MustNewFromHexString("0xbB3fbc5Ab07D7866b62860aB540E5F581B8710E1"),
				UserID:        "test-user-id",
				ApplicationID: "",
			},
		}
		acc, err = app.AccountUseCase.CreateAccount(ctx, input)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
		require.Nil(t, acc)
	})

	t.Run("failure: already exists", func(t *testing.T) {
		input := user.CreateAccountInput{
			AccountID: user.AccountID{
				Address:       address.MustNewFromHexString("0xbB3fbc5Ab07D7866b62860aB540E5F581B8710E1"),
				UserID:        userID,
				ApplicationID: applicationID,
			},
		}
		acc, err := app.AccountUseCase.CreateAccount(ctx, input)
		require.NoError(t, err)
		require.NotNil(t, acc)

		// Duplicated Entry
		acc, err = app.AccountUseCase.CreateAccount(ctx, input)
		require.Error(t, err)
		require.True(t, errors.IsAlreadyExists(err))
		require.Nil(t, acc)
	})

	t.Run("success", func(t *testing.T) {
		input := user.CreateAccountInput{
			AccountID: user.AccountID{
				Address:       address.MustNewFromHexString("0xbB3fbc5Ab07D7866b62860aB540E5F581B8710A1"),
				UserID:        userID,
				ApplicationID: applicationID,
			},
		}
		acc, err := app.AccountUseCase.CreateAccount(ctx, input)
		require.NoError(t, err)
		require.NotNil(t, acc)
		require.NotNil(t, acc.InternalResourceID)
		require.Equal(t, input.Address.String(), acc.Address.String())
		require.Equal(t, input.UserID, acc.UserID)
		require.Equal(t, input.ApplicationID, acc.ApplicationID)
	})
}

func TestDefaultUseCase_ListAccounts(t *testing.T) {
	applicationOneID := uuid.NewString()
	createApplicationOneInput := application.CreateApplicationInput{
		ID:      &applicationOneID,
		ChainID: *chainID,
	}
	createApplicationOneOutput, createApplicationOneErr := app.ApplicationUseCase.CreateApplication(ctx, createApplicationOneInput)
	require.NoError(t, createApplicationOneErr)
	require.NotNil(t, createApplicationOneOutput)

	userOneID := uuid.NewString()
	createUserOneInput := user.CreateUserInput{
		ID:            &userOneID,
		ApplicationID: applicationOneID,
		Roles:         []string{"application-admin"},
	}
	createUserOneOutput, createUserOneErr := app.UserUseCase.CreateUser(ctx, createUserOneInput)
	require.NoError(t, createUserOneErr)
	require.NotNil(t, createUserOneOutput)

	applicationTwoID := uuid.NewString()
	createApplicationTwoInput := application.CreateApplicationInput{
		ID:      &applicationTwoID,
		ChainID: *chainID,
	}
	createApplicationTwoOutput, createApplicationTwoErr := app.ApplicationUseCase.CreateApplication(ctx, createApplicationTwoInput)
	require.NoError(t, createApplicationTwoErr)
	require.NotNil(t, createApplicationTwoOutput)

	userTwoID := uuid.NewString()
	createUserTwoInput := user.CreateUserInput{
		ID:            &userTwoID,
		ApplicationID: applicationTwoID,
		Roles:         []string{"application-admin"},
	}
	createUserTwoOutput, createUserTwoErr := app.UserUseCase.CreateUser(ctx, createUserTwoInput)
	require.NoError(t, createUserTwoErr)
	require.NotNil(t, createUserTwoOutput)

	accountsToCreate := 20
	// Accounts for application-1
	for i := 0; i < accountsToCreate; i++ {
		input := user.CreateAccountInput{
			AccountID: user.AccountID{
				Address:       address.MustNewFromHexString(addresses[i]),
				UserID:        userOneID,
				ApplicationID: applicationOneID,
			},
		}
		_, err := app.AccountUseCase.CreateAccount(ctx, input)
		require.NoError(t, err)
	}
	// Accounts for application-2
	for i := 0; i < accountsToCreate; i++ {
		input := user.CreateAccountInput{
			AccountID: user.AccountID{
				Address:       address.MustNewFromHexString(addresses[i]),
				UserID:        userTwoID,
				ApplicationID: applicationTwoID,
			},
		}
		_, err := app.AccountUseCase.CreateAccount(ctx, input)
		require.NoError(t, err)
	}

	t.Run("failure: invalid input arguments", func(t *testing.T) {
		output, err := app.AccountUseCase.ListAccounts(ctx, user.ListAccountsInput{
			ApplicationID: "",
		})
		require.Nil(t, output)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
	})

	t.Run("success: list all accounts", func(t *testing.T) {
		output, err := app.AccountUseCase.ListAccounts(ctx, user.ListAccountsInput{
			ApplicationID: applicationOneID,
		})
		require.NoError(t, err)
		require.Len(t, output.Items, accountsToCreate)
		for _, acc := range output.Items {
			require.NotNil(t, acc)
			require.Equal(t, acc.ApplicationID, applicationOneID)
		}
	})

	t.Run("success: list all accounts for a specific user", func(t *testing.T) {
		output, err := app.AccountUseCase.ListAccounts(ctx, user.ListAccountsInput{
			ApplicationID: applicationOneID,
			UserID:        &userOneID,
		})
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(output.Items), accountsToCreate)
		require.NotNil(t, output)
		require.Equal(t, output.Items[0].ApplicationID, applicationOneID)
	})

	t.Run("success: list all accounts for a specific slot", func(t *testing.T) {
		output, err := app.AccountUseCase.ListAccounts(ctx, user.ListAccountsInput{
			ApplicationID: applicationTwoID,
		})
		require.NoError(t, err)
		require.Len(t, output.Items, accountsToCreate)
		for _, acc := range output.Items {
			require.NotNil(t, acc)
			require.Equal(t, acc.ApplicationID, applicationTwoID)
		}
	})
}

func TestDefaultUseCase_GetAccount(t *testing.T) {
	applicationID := uuid.NewString()
	createApplicationInput := application.CreateApplicationInput{
		ID:      &applicationID,
		ChainID: *chainID,
	}
	createApplicationOutput, createApplicationErr := app.ApplicationUseCase.CreateApplication(ctx, createApplicationInput)
	require.NoError(t, createApplicationErr)
	require.NotNil(t, createApplicationOutput)

	userID := uuid.NewString()
	createUserInput := user.CreateUserInput{
		ID:            &userID,
		ApplicationID: applicationID,
		Roles:         []string{"application-admin"},
	}
	createUserOutput, createUserErr := app.UserUseCase.CreateUser(ctx, createUserInput)
	require.NoError(t, createUserErr)
	require.NotNil(t, createUserOutput)

	t.Run("failure: invalid input arguments", func(t *testing.T) {
		input := user.GetAccountInput{
			AccountID: user.AccountID{
				Address:       address.Address{},
				UserID:        userID,
				ApplicationID: applicationID,
			},
		}
		output, err := app.AccountUseCase.GetAccount(ctx, input)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
		require.Nil(t, output)

		input = user.GetAccountInput{
			AccountID: user.AccountID{
				Address:       address.MustNewFromHexString(addresses[0]),
				UserID:        "",
				ApplicationID: applicationID,
			},
		}
		output, err = app.AccountUseCase.GetAccount(ctx, input)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
		require.Nil(t, output)

		input = user.GetAccountInput{
			AccountID: user.AccountID{
				Address:       address.MustNewFromHexString(addresses[0]),
				UserID:        userID,
				ApplicationID: "",
			},
		}
		output, err = app.AccountUseCase.GetAccount(ctx, input)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
		require.Nil(t, output)
	})

	t.Run("failure: application not found", func(t *testing.T) {
		input := user.GetAccountInput{
			AccountID: user.AccountID{
				Address:       address.MustNewFromHexString(addresses[0]),
				UserID:        userID,
				ApplicationID: applicationID,
			},
		}
		output, err := app.AccountUseCase.GetAccount(ctx, input)
		require.Error(t, err)
		require.True(t, errors.IsNotFound(err))
		require.Nil(t, output)
	})

	t.Run("success", func(t *testing.T) {
		// Create account before retrieving it
		createAccountInput := user.CreateAccountInput{
			AccountID: user.AccountID{
				Address:       address.MustNewFromHexString(addresses[0]),
				UserID:        userID,
				ApplicationID: applicationID,
			},
		}
		_, err := app.AccountUseCase.CreateAccount(ctx, createAccountInput)
		require.NoError(t, err)

		createdAcc, err := app.AccountUseCase.GetAccount(ctx, user.GetAccountInput(createAccountInput))
		require.NoError(t, err)
		require.NotNil(t, createdAcc)
		require.Equal(t, createdAcc.UserID, createAccountInput.UserID)
		require.Equal(t, createdAcc.ApplicationID, createAccountInput.ApplicationID)
		require.Equal(t, createdAcc.Address.String(), createAccountInput.Address.String())
	})
}

func TestDefaultUseCase_DeleteAccount(t *testing.T) {
	applicationID := uuid.NewString()
	createApplicationInput := application.CreateApplicationInput{
		ID:      &applicationID,
		ChainID: *chainID,
	}
	createApplicationOutput, createApplicationErr := app.ApplicationUseCase.CreateApplication(ctx, createApplicationInput)
	require.NoError(t, createApplicationErr)
	require.NotNil(t, createApplicationOutput)

	userID := uuid.NewString()
	createUserInput := user.CreateUserInput{
		ID:            &userID,
		ApplicationID: applicationID,
		Roles:         []string{"application-admin"},
	}
	createUserOutput, createUserErr := app.UserUseCase.CreateUser(ctx, createUserInput)
	require.NoError(t, createUserErr)
	require.NotNil(t, createUserOutput)

	t.Run("failure: invalid input arguments", func(t *testing.T) {
		deleteInput := user.DeleteAccountInput{
			AccountID: user.AccountID{
				Address:       address.Address{},
				UserID:        userID,
				ApplicationID: applicationID,
			},
		}
		output, err := app.AccountUseCase.DeleteAccount(ctx, deleteInput)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
		require.Nil(t, output)

		deleteInput = user.DeleteAccountInput{
			AccountID: user.AccountID{
				Address:       address.MustNewFromHexString(addresses[0]),
				UserID:        "",
				ApplicationID: applicationID,
			},
		}
		output, err = app.AccountUseCase.DeleteAccount(ctx, deleteInput)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
		require.Nil(t, output)

		deleteInput = user.DeleteAccountInput{
			AccountID: user.AccountID{
				Address:       address.MustNewFromHexString(addresses[0]),
				UserID:        userID,
				ApplicationID: "",
			},
		}
		output, err = app.AccountUseCase.DeleteAccount(ctx, deleteInput)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
		require.Nil(t, output)
	})

	t.Run("failure: nonexistent account", func(t *testing.T) {
		deleteInput := user.DeleteAccountInput{
			AccountID: user.AccountID{
				Address:       address.MustNewFromHexString(addresses[0]),
				UserID:        userID,
				ApplicationID: applicationID,
			},
		}
		output, err := app.AccountUseCase.DeleteAccount(ctx, deleteInput)
		require.Error(t, err)
		require.True(t, errors.IsNotFound(err))
		require.Nil(t, output)
	})

	t.Run("success", func(t *testing.T) {
		// Create a valid account
		createAccountInput := user.CreateAccountInput{
			AccountID: user.AccountID{
				Address:       address.MustNewFromHexString(addresses[0]),
				UserID:        userID,
				ApplicationID: applicationID,
			},
		}
		createdAcc, err := app.AccountUseCase.CreateAccount(ctx, createAccountInput)
		require.NoError(t, err)

		// Delete the account
		deleteInput := user.DeleteAccountInput(createAccountInput)
		deletedAcc, err := app.AccountUseCase.DeleteAccount(ctx, deleteInput)
		require.NoError(t, err)
		require.NotNil(t, deletedAcc)
		require.Equal(t, createdAcc.Account, deletedAcc.Account)

		// Retrieve deleted account
		getOutput, err := app.AccountUseCase.GetAccount(ctx, user.GetAccountInput(createAccountInput))
		require.Error(t, err)
		require.True(t, errors.IsNotFound(err))
		require.Nil(t, getOutput)
	})
}

func TestDefaultUseCase_DeleteAllAccountsForAddress(t *testing.T) {
	applicationID := uuid.NewString()
	createApplicationInput := application.CreateApplicationInput{
		ID:      &applicationID,
		ChainID: *chainID,
	}
	createApplicationOutput, createApplicationErr := app.ApplicationUseCase.CreateApplication(ctx, createApplicationInput)
	require.NoError(t, createApplicationErr)
	require.NotNil(t, createApplicationOutput)

	userID := uuid.NewString()
	createUserInput := user.CreateUserInput{
		ID:            &userID,
		ApplicationID: applicationID,
		Roles:         []string{"application-admin"},
	}
	createUserOutput, createUserErr := app.UserUseCase.CreateUser(ctx, createUserInput)
	require.NoError(t, createUserErr)
	require.NotNil(t, createUserOutput)

	hsmModuleID := uuid.NewString()
	createHSMModuleInput := hsmmodule.CreateHSMModuleInput{
		ID: &hsmModuleID,
		Configuration: hsmmodule.HSMModuleConfiguration{
			SoftHSMConfiguration: &hsmmodule.SoftHSMConfiguration{},
		},
		ModuleKind: hsmmodule.SoftHSMModuleKind,
	}
	createHSMModuleOutput, createHSMModuleErr := app.HSMModuleUseCase.CreateHSMModule(ctx, createHSMModuleInput)
	require.NoError(t, createHSMModuleErr)
	require.NotNil(t, createHSMModuleOutput)

	createHSMSlotInput := hsmslot.CreateHSMSlotInput{
		ApplicationID: applicationID,
		HSMModuleID:   hsmModuleID,
		Slot:          slotID,
		Pin:           slotPin,
	}
	createHSMSlotOutput, createHSMSlotErr := app.HSMSlotUseCase.CreateHSMSlot(ctx, createHSMSlotInput)
	require.NoError(t, createHSMSlotErr)
	require.NotNil(t, createHSMSlotOutput)

	t.Run("failure: invalid input arguments", func(t *testing.T) {
		deleteInput := user.DeleteAllAccountsForAddressInput{
			Address:       address.Address{},
			ApplicationID: applicationID,
		}
		output, err := app.AccountUseCase.DeleteAllAccountsForAddress(ctx, deleteInput)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
		require.Nil(t, output)

		deleteInput = user.DeleteAllAccountsForAddressInput{
			Address:       address.MustNewFromHexString(addresses[0]),
			ApplicationID: "",
		}
		output, err = app.AccountUseCase.DeleteAllAccountsForAddress(ctx, deleteInput)
		require.Error(t, err)
		require.True(t, errors.IsInvalidArgument(err))
		require.Nil(t, output)
	})

	t.Run("success", func(t *testing.T) {
		// Create a valid account
		createAccountInputOne := user.CreateAccountInput{
			AccountID: user.AccountID{
				Address:       address.MustNewFromHexString(hsmLoadedAddress),
				UserID:        userID,
				ApplicationID: applicationID,
			},
		}
		_, err := app.AccountUseCase.CreateAccount(ctx, createAccountInputOne)
		require.NoError(t, err)

		// Create another valid account
		userTwoID := uuid.NewString()
		createUserTwoInput := user.CreateUserInput{
			ID:            &userTwoID,
			ApplicationID: applicationID,
			Roles:         []string{"application-admin"},
		}
		createUserTwoOutput, createUserTwoErr := app.UserUseCase.CreateUser(ctx, createUserTwoInput)
		require.NoError(t, createUserTwoErr)
		require.NotNil(t, createUserTwoOutput)
		createAccountTwoInput := user.CreateAccountInput{
			AccountID: user.AccountID{
				Address:       address.MustNewFromHexString(hsmLoadedAddress),
				UserID:        userTwoID,
				ApplicationID: applicationID,
			},
		}
		_, createAccountTwoErr := app.AccountUseCase.CreateAccount(ctx, createAccountTwoInput)
		require.NoError(t, createAccountTwoErr)

		// Delete both accounts
		deleteInput := user.DeleteAllAccountsForAddressInput{
			Address:       address.MustNewFromHexString(hsmLoadedAddress),
			ApplicationID: applicationID,
		}
		deletedAccounts, err := app.AccountUseCase.DeleteAllAccountsForAddress(ctx, deleteInput)
		require.NoError(t, err)
		require.NotNil(t, deletedAccounts)
		require.Len(t, deletedAccounts.Items, 2)

		// Retrieve deleted accounts
		getOutput, err := app.AccountUseCase.GetAccount(ctx, user.GetAccountInput(createAccountInputOne))
		require.Error(t, err)
		require.True(t, errors.IsNotFound(err))
		require.Nil(t, getOutput)

		getOutput, err = app.AccountUseCase.GetAccount(ctx, user.GetAccountInput(createAccountTwoInput))
		require.Error(t, err)
		require.True(t, errors.IsNotFound(err))
		require.Nil(t, getOutput)
	})
}
