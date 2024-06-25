package hsmconnector_test

import (
	"context"
	"encoding/hex"
	"math/big"
	"os"
	"testing"

	"github.com/hyperledger-labs/signare/app/pkg/commons/validators"
	"github.com/hyperledger-labs/signare/app/pkg/entities"
	"github.com/hyperledger-labs/signare/app/pkg/entities/address"
	"github.com/hyperledger-labs/signare/app/pkg/graph"
	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/application"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/hsmconnector"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/hsmmodule"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/hsmslot"
	"github.com/hyperledger-labs/signare/app/test/dbtesthelper"
	"github.com/hyperledger-labs/signare/app/test/signaturemanagertesthelper"

	"github.com/stretchr/testify/require"
)

var (
	app    graph.GraphShared
	slotID string
	ctx    context.Context

	applicationID  = "my-app"
	moduleID       = "module-id"
	chainID        = entities.NewInt256FromInt(44844)
	invalidChainID = entities.NewInt256FromInt(0)
	validAddress   = address.MustNewFromHexString("0x970e8128ab834e8eac17ab8e3812f010678cf791")
	slotPin        = signaturemanagertesthelper.SlotPin
)

func TestMain(m *testing.M) {
	initializedSlotID, _, err := signaturemanagertesthelper.InitializeSoftHSMSlot()
	if err != nil {
		panic(err)
	}
	slotID = *initializedSlotID

	a, err := dbtesthelper.InitializeApp()
	if err != nil {
		panic(err)
	}
	app = *a

	ctx = context.Background()
	err = provisionTest(ctx)
	if err != nil {
		panic(err)
	}

	validators.SetValidators()
	os.Exit(m.Run())
}

func TestProvideDefaultUseCase(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		options := hsmconnector.DefaultUseCaseOptions{
			DigitalSignatureManagerFactory: app.DigitalSignatureManagerFactory,
		}
		defaultUseCase, err := hsmconnector.ProvideDefaultHSMConnector(options)
		require.Nil(t, err)
		require.NotNil(t, defaultUseCase)
	})
	t.Run("nil digitalSignatureManagerFactory", func(t *testing.T) {
		options := hsmconnector.DefaultUseCaseOptions{
			DigitalSignatureManagerFactory: nil,
		}
		defaultUseCase, err := hsmconnector.ProvideDefaultHSMConnector(options)
		require.Error(t, err)
		require.Nil(t, defaultUseCase)
	})

}

func TestDefaultUseCase_GenerateAddress(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		generateAddressInput := hsmconnector.GenerateAddressInput{
			SlotConnectionData: hsmconnector.SlotConnectionData{
				Slot:       slotID,
				Pin:        slotPin,
				ModuleKind: hsmconnector.SoftHSMModuleKind,
				ChainID:    *chainID,
			},
		}
		generateAddressOutput, generateAddressErr := app.HSMConnector.GenerateAddress(ctx, generateAddressInput)
		require.Nil(t, generateAddressErr)
		require.NotNil(t, generateAddressOutput)

		// Clean up the created resource
		removeAddressInput := hsmconnector.RemoveAddressInput{
			SlotConnectionData: hsmconnector.SlotConnectionData{
				Slot:       slotID,
				Pin:        slotPin,
				ModuleKind: hsmconnector.SoftHSMModuleKind,
				ChainID:    *chainID,
			},
			Address: generateAddressOutput.Address,
		}
		removeAddressOutput, removeAddressErr := app.HSMConnector.RemoveAddress(ctx, removeAddressInput)
		require.Nil(t, removeAddressErr)
		require.NotNil(t, removeAddressOutput)
	})

	t.Run("failure: invalid input arguments", func(t *testing.T) {
		generateAddressInput := hsmconnector.GenerateAddressInput{
			SlotConnectionData: hsmconnector.SlotConnectionData{
				Slot:       "",
				Pin:        slotPin,
				ModuleKind: hsmconnector.SoftHSMModuleKind,
				ChainID:    *chainID,
			},
		}
		generateAddressOutput, generateAddressErr := app.HSMConnector.GenerateAddress(ctx, generateAddressInput)
		require.Error(t, generateAddressErr)
		require.True(t, errors.IsInvalidArgument(generateAddressErr))
		require.Nil(t, generateAddressOutput)

		generateAddressInput = hsmconnector.GenerateAddressInput{
			SlotConnectionData: hsmconnector.SlotConnectionData{
				Slot:       slotID,
				Pin:        "",
				ModuleKind: hsmconnector.SoftHSMModuleKind,
				ChainID:    *chainID,
			},
		}
		generateAddressOutput, generateAddressErr = app.HSMConnector.GenerateAddress(ctx, generateAddressInput)
		require.Error(t, generateAddressErr)
		require.True(t, errors.IsInvalidArgument(generateAddressErr))
		require.Nil(t, generateAddressOutput)

		generateAddressInput = hsmconnector.GenerateAddressInput{
			SlotConnectionData: hsmconnector.SlotConnectionData{
				Slot:       slotID,
				Pin:        slotPin,
				ModuleKind: "invalid module kind",
				ChainID:    *chainID,
			},
		}
		generateAddressOutput, generateAddressErr = app.HSMConnector.GenerateAddress(ctx, generateAddressInput)
		require.Error(t, generateAddressErr)
		require.True(t, errors.IsInvalidArgument(generateAddressErr))
		require.Nil(t, generateAddressOutput)

		generateAddressInput = hsmconnector.GenerateAddressInput{
			SlotConnectionData: hsmconnector.SlotConnectionData{
				Slot:       slotID,
				Pin:        slotPin,
				ModuleKind: hsmconnector.SoftHSMModuleKind,
				ChainID:    *invalidChainID,
			},
		}
		generateAddressOutput, generateAddressErr = app.HSMConnector.GenerateAddress(ctx, generateAddressInput)
		require.Error(t, generateAddressErr)
		require.True(t, errors.IsInvalidArgument(generateAddressErr))
		require.Nil(t, generateAddressOutput)
	})
}

func TestDefaultUseCase_RemoveAddress(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		createAddressInput := hsmconnector.GenerateAddressInput{
			SlotConnectionData: hsmconnector.SlotConnectionData{
				Slot:       slotID,
				Pin:        slotPin,
				ModuleKind: hsmconnector.SoftHSMModuleKind,
				ChainID:    *chainID,
			},
		}
		createAddressOutput, createAddressErr := app.HSMConnector.GenerateAddress(ctx, createAddressInput)
		require.Nil(t, createAddressErr)
		require.NotNil(t, createAddressOutput)

		removeAddressInput := hsmconnector.RemoveAddressInput{
			Address: createAddressOutput.Address,
			SlotConnectionData: hsmconnector.SlotConnectionData{
				Slot:       slotID,
				Pin:        slotPin,
				ModuleKind: hsmconnector.SoftHSMModuleKind,
				ChainID:    *chainID,
			},
		}
		removeAddressOutput, removeAddressErr := app.HSMConnector.RemoveAddress(ctx, removeAddressInput)
		require.Nil(t, removeAddressErr)
		require.NotNil(t, removeAddressOutput)
	})

	t.Run("failure: invalid input arguments", func(t *testing.T) {
		removeAddressInput := hsmconnector.RemoveAddressInput{
			SlotConnectionData: hsmconnector.SlotConnectionData{
				Slot:       "",
				Pin:        slotPin,
				ModuleKind: hsmconnector.SoftHSMModuleKind,
				ChainID:    *chainID,
			},
			Address: validAddress,
		}
		removeAddressOutput, removeAddressErr := app.HSMConnector.RemoveAddress(ctx, removeAddressInput)
		require.Error(t, removeAddressErr)
		require.True(t, errors.IsInvalidArgument(removeAddressErr))
		require.Nil(t, removeAddressOutput)

		removeAddressInput = hsmconnector.RemoveAddressInput{
			SlotConnectionData: hsmconnector.SlotConnectionData{
				Slot:       slotID,
				Pin:        "",
				ModuleKind: hsmconnector.SoftHSMModuleKind,
				ChainID:    *chainID,
			},
			Address: validAddress,
		}
		removeAddressOutput, removeAddressErr = app.HSMConnector.RemoveAddress(ctx, removeAddressInput)
		require.Error(t, removeAddressErr)
		require.True(t, errors.IsInvalidArgument(removeAddressErr))
		require.Nil(t, removeAddressOutput)

		removeAddressInput = hsmconnector.RemoveAddressInput{
			SlotConnectionData: hsmconnector.SlotConnectionData{
				Slot:       slotID,
				Pin:        slotPin,
				ModuleKind: "invalid type",
				ChainID:    *chainID,
			},
			Address: validAddress,
		}
		removeAddressOutput, removeAddressErr = app.HSMConnector.RemoveAddress(ctx, removeAddressInput)
		require.Error(t, removeAddressErr)
		require.True(t, errors.IsInvalidArgument(removeAddressErr))
		require.Nil(t, removeAddressOutput)

		removeAddressInput = hsmconnector.RemoveAddressInput{
			SlotConnectionData: hsmconnector.SlotConnectionData{
				Slot:       slotID,
				Pin:        slotPin,
				ModuleKind: hsmconnector.SoftHSMModuleKind,
				ChainID:    *invalidChainID,
			},
			Address: validAddress,
		}
		removeAddressOutput, removeAddressErr = app.HSMConnector.RemoveAddress(ctx, removeAddressInput)
		require.Error(t, removeAddressErr)
		require.True(t, errors.IsInvalidArgument(removeAddressErr))
		require.Nil(t, removeAddressOutput)

		removeAddressInput = hsmconnector.RemoveAddressInput{
			SlotConnectionData: hsmconnector.SlotConnectionData{
				Slot:       slotID,
				Pin:        slotPin,
				ModuleKind: hsmconnector.SoftHSMModuleKind,
				ChainID:    *chainID,
			},
			Address: address.ZeroAddress,
		}
		removeAddressOutput, removeAddressErr = app.HSMConnector.RemoveAddress(ctx, removeAddressInput)
		require.Error(t, removeAddressErr)
		require.True(t, errors.IsInvalidArgument(removeAddressErr))
		require.Nil(t, removeAddressOutput)
	})

	t.Run("failure: address doesn't exist", func(t *testing.T) {
		removeAddressInput := hsmconnector.RemoveAddressInput{
			Address: validAddress,
			SlotConnectionData: hsmconnector.SlotConnectionData{
				Slot:       slotID,
				Pin:        slotPin,
				ModuleKind: hsmconnector.SoftHSMModuleKind,
				ChainID:    *chainID,
			},
		}
		removeAddressOutput, removeAddressErr := app.HSMConnector.RemoveAddress(ctx, removeAddressInput)
		require.Error(t, removeAddressErr)
		require.True(t, errors.IsNotFound(removeAddressErr))
		require.Nil(t, removeAddressOutput)
	})
}

func TestDefaultUseCase_ListAddress(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		createAddressInput := hsmconnector.GenerateAddressInput{
			SlotConnectionData: hsmconnector.SlotConnectionData{
				Slot:       slotID,
				Pin:        slotPin,
				ModuleKind: hsmconnector.SoftHSMModuleKind,
				ChainID:    *chainID,
			},
		}
		createAddressOutputOne, createAddressOneErr := app.HSMConnector.GenerateAddress(ctx, createAddressInput)
		require.Nil(t, createAddressOneErr)
		require.NotNil(t, createAddressOutputOne)

		createAddressOutputTwo, createAddressTwoErr := app.HSMConnector.GenerateAddress(ctx, createAddressInput)
		require.Nil(t, createAddressTwoErr)
		require.NotNil(t, createAddressOutputTwo)
		listAddressInput := hsmconnector.ListAddressesInput{
			SlotConnectionData: hsmconnector.SlotConnectionData{
				Slot:       slotID,
				Pin:        slotPin,
				ModuleKind: hsmconnector.SoftHSMModuleKind,
				ChainID:    *chainID,
			},
		}

		listAddressOutput, listAddressErr := app.HSMConnector.ListAddresses(ctx, listAddressInput)
		require.Nil(t, listAddressErr)
		require.NotNil(t, listAddressOutput)
		require.Positive(t, len(listAddressOutput.Items))

		defaultLoadedAddressIncluded := false
		addressOneIncluded := false
		addressTwoIncluded := false
		for _, addr := range listAddressOutput.Items {
			if addr.String() == createAddressOutputOne.Address.String() {
				addressOneIncluded = true
			}
			if addr.String() == createAddressOutputTwo.Address.String() {
				addressTwoIncluded = true
			}
			if addr.String() == signaturemanagertesthelper.ImportedKeyAddress {
				defaultLoadedAddressIncluded = true
			}
		}
		require.True(t, addressOneIncluded)
		require.True(t, addressTwoIncluded)
		require.True(t, defaultLoadedAddressIncluded)
	})

	t.Run("failure: invalid input arguments", func(t *testing.T) {
		listAddressInput := hsmconnector.ListAddressesInput{
			SlotConnectionData: hsmconnector.SlotConnectionData{
				Slot:       "",
				Pin:        slotPin,
				ModuleKind: hsmconnector.SoftHSMModuleKind,
				ChainID:    *chainID,
			},
		}
		listAddressOutput, listAddressErr := app.HSMConnector.ListAddresses(ctx, listAddressInput)
		require.Error(t, listAddressErr)
		require.True(t, errors.IsInvalidArgument(listAddressErr))
		require.Nil(t, listAddressOutput)

		listAddressInput = hsmconnector.ListAddressesInput{
			SlotConnectionData: hsmconnector.SlotConnectionData{
				Slot:       slotID,
				Pin:        "",
				ModuleKind: hsmconnector.SoftHSMModuleKind,
				ChainID:    *chainID,
			},
		}
		listAddressOutput, listAddressErr = app.HSMConnector.ListAddresses(ctx, listAddressInput)
		require.Error(t, listAddressErr)
		require.True(t, errors.IsInvalidArgument(listAddressErr))
		require.Nil(t, listAddressOutput)

		listAddressInput = hsmconnector.ListAddressesInput{
			SlotConnectionData: hsmconnector.SlotConnectionData{
				Slot:       slotID,
				Pin:        slotPin,
				ModuleKind: "invalid module kind",
				ChainID:    *chainID,
			},
		}
		listAddressOutput, listAddressErr = app.HSMConnector.ListAddresses(ctx, listAddressInput)
		require.Error(t, listAddressErr)
		require.True(t, errors.IsInvalidArgument(listAddressErr))
		require.Nil(t, listAddressOutput)

		listAddressInput = hsmconnector.ListAddressesInput{
			SlotConnectionData: hsmconnector.SlotConnectionData{
				Slot:       slotID,
				Pin:        slotPin,
				ModuleKind: hsmconnector.SoftHSMModuleKind,
				ChainID:    *invalidChainID,
			},
		}
		listAddressOutput, listAddressErr = app.HSMConnector.ListAddresses(ctx, listAddressInput)
		require.Error(t, listAddressErr)
		require.True(t, errors.IsInvalidArgument(listAddressErr))
		require.Nil(t, listAddressOutput)
	})
}

func TestDefaultUseCase_SignTx(t *testing.T) {
	toAddress := address.MustNewFromHexString("0xA4F666f1860D2aCbe49b342C87867754a21dE850")
	gasPrice := big.NewInt(20)
	value := big.NewInt(3)
	nonce := entities.UInt64(1)

	t.Run("success: regular transaction filled values", func(t *testing.T) {
		data := entities.NewHexBytes(hexStringToBytes("0x1f170873")) // simpleMethod()
		signTxInput := hsmconnector.SignTxInput{
			SlotConnectionData: hsmconnector.SlotConnectionData{
				Slot:       slotID,
				Pin:        slotPin,
				ModuleKind: hsmconnector.SoftHSMModuleKind,
				ChainID:    *chainID,
			},
			From: address.MustNewFromHexString(signaturemanagertesthelper.ImportedKeyAddress),
			To:   &toAddress,
			Gas: &entities.HexUInt64{
				UInt64: 1000,
			},
			GasPrice: &entities.HexInt256{
				Int256: entities.Int256{
					Int: *gasPrice,
				},
			},
			Value: &entities.HexInt256{
				Int256: entities.Int256{
					Int: *value,
				},
			},
			Data: *data,
			Nonce: entities.HexUInt64{
				UInt64: nonce,
			},
		}
		signTxOutput, err := app.HSMConnector.SignTx(ctx, signTxInput)
		require.Nil(t, err)
		require.NotNil(t, signTxOutput)
	})
	t.Run("success: regular transaction defaults for optional values", func(t *testing.T) {
		data := entities.NewHexBytes(hexStringToBytes("0x"))
		signTxInput := hsmconnector.SignTxInput{
			SlotConnectionData: hsmconnector.SlotConnectionData{
				Slot:       slotID,
				Pin:        slotPin,
				ModuleKind: hsmconnector.SoftHSMModuleKind,
				ChainID:    *chainID,
			},
			From:     address.MustNewFromHexString(signaturemanagertesthelper.ImportedKeyAddress),
			To:       &toAddress,
			Gas:      nil,
			GasPrice: nil,
			Value:    nil,
			Data:     *data,
			Nonce: entities.HexUInt64{
				UInt64: nonce,
			},
		}
		signTxOutput, err := app.HSMConnector.SignTx(ctx, signTxInput)
		require.Nil(t, err)
		require.NotNil(t, signTxOutput)
	})
	t.Run("success: smart contract deployment (nil to address)", func(t *testing.T) {
		data := entities.NewHexBytes(hexStringToBytes("0x1234"))
		signTxInput := hsmconnector.SignTxInput{
			SlotConnectionData: hsmconnector.SlotConnectionData{
				Slot:       slotID,
				Pin:        slotPin,
				ModuleKind: hsmconnector.SoftHSMModuleKind,
				ChainID:    *chainID,
			},
			From: address.MustNewFromHexString(signaturemanagertesthelper.ImportedKeyAddress),
			To:   nil,
			Gas: &entities.HexUInt64{
				UInt64: 1000,
			},
			GasPrice: &entities.HexInt256{
				Int256: entities.Int256{
					Int: *gasPrice,
				},
			},
			Value: &entities.HexInt256{
				Int256: entities.Int256{
					Int: *value,
				},
			},
			Data: *data,
			Nonce: entities.HexUInt64{
				UInt64: nonce,
			},
		}
		signTxOutput, err := app.HSMConnector.SignTx(ctx, signTxInput)
		require.Nil(t, err)
		require.NotNil(t, signTxOutput)
	})
}

func hexStringToBytes(input string) []byte {
	if len(input) == 0 {
		panic("empty string")
	}
	if !has0xPrefix(input) {
		panic("missing '0x' prefix")
	}
	b, err := hex.DecodeString(input[2:])
	if err != nil {
		panic(err.Error())
	}
	return b
}

func has0xPrefix(input string) bool {
	return len(input) >= 2 && input[0] == '0' && (input[1] == 'x' || input[1] == 'X')
}

func provisionTest(ctx context.Context) error {
	createApplicationInput := application.CreateApplicationInput{
		ID:      &applicationID,
		ChainID: *chainID,
	}
	_, err := app.ApplicationUseCase.CreateApplication(ctx, createApplicationInput)
	if err != nil {
		return err
	}

	createModuleInput := hsmmodule.CreateHSMModuleInput{
		ID: &moduleID,
		Configuration: hsmmodule.HSMModuleConfiguration{
			SoftHSMConfiguration: &hsmmodule.SoftHSMConfiguration{},
		},
		ModuleKind: hsmmodule.SoftHSMModuleKind,
	}
	_, err = app.HSMModuleUseCase.CreateHSMModule(ctx, createModuleInput)
	if err != nil {
		return err
	}

	createSlotInput := hsmslot.CreateHSMSlotInput{
		ApplicationID: applicationID,
		HSMModuleID:   moduleID,
		Slot:          slotID,
		Pin:           slotPin,
	}
	_, err = app.HSMSlotUseCase.CreateHSMSlot(ctx, createSlotInput)
	if err != nil {
		return err
	}
	return nil
}
