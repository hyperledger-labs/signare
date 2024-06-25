package hsmslot

import (
	"context"
	"fmt"

	"github.com/hyperledger-labs/signare/app/pkg/entities"
	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/application"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/hsmmodule"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/referentialintegrity"
)

func (u *DefaultUseCase) addApplicationDependency(ctx context.Context, data HSMSlot) error {
	applicationStandardID := application.GetApplicationInput{
		StandardID: entities.StandardID{
			ID: data.ApplicationID,
		},
	}
	getApplicationOutput, getApplicationErr := u.applicationUseCase.GetApplication(ctx, applicationStandardID)
	if getApplicationErr != nil {
		if errors.IsNotFound(getApplicationErr) {
			msg := fmt.Sprintf("can't create HSM slot becauase application '%s' does not exist", data.ApplicationID)
			return errors.PreconditionFailed().WithMessage(msg).SetHumanReadableMessage(msg)
		}
		return getApplicationErr
	}

	var referentialIntegrityCreateEntryInput referentialintegrity.CreateEntryInput
	referentialIntegrityCreateEntryInput.ResourceID = string(data.InternalResourceID)
	referentialIntegrityCreateEntryInput.ResourceKind = referentialintegrity.KindHSMSlot
	referentialIntegrityCreateEntryInput.ParentResourceID = string(getApplicationOutput.InternalResourceID)
	referentialIntegrityCreateEntryInput.ParentResourceKind = referentialintegrity.KindApplication

	_, createEntryErr := u.referentialIntegrityUseCase.CreateEntry(ctx, referentialIntegrityCreateEntryInput)
	if createEntryErr != nil && !errors.IsAlreadyExists(createEntryErr) {
		return createEntryErr
	}
	return nil
}

func (u *DefaultUseCase) addHSMModuleDependency(ctx context.Context, data HSMSlot) error {
	hsmModuleStandardID := hsmmodule.GetHSMModuleInput{
		StandardID: entities.StandardID{
			ID: data.HSMModuleID,
		},
	}
	getHSMModuleOutput, getHSMModuleErr := u.hsmModuleUseCase.GetHSMModule(ctx, hsmModuleStandardID)
	if getHSMModuleErr != nil {
		if errors.IsNotFound(getHSMModuleErr) {
			msg := fmt.Sprintf("can't create HSM slot becauase HSM module '%s' does not exist", data.HSMModuleID)
			return errors.PreconditionFailed().WithMessage(msg).SetHumanReadableMessage(msg)
		}
		return getHSMModuleErr
	}

	var referentialIntegrityCreateEntryInput referentialintegrity.CreateEntryInput
	referentialIntegrityCreateEntryInput.ResourceID = string(data.InternalResourceID)
	referentialIntegrityCreateEntryInput.ResourceKind = referentialintegrity.KindHSMSlot
	referentialIntegrityCreateEntryInput.ParentResourceID = string(getHSMModuleOutput.InternalResourceID)
	referentialIntegrityCreateEntryInput.ParentResourceKind = referentialintegrity.KindHSMModule

	_, createEntryErr := u.referentialIntegrityUseCase.CreateEntry(ctx, referentialIntegrityCreateEntryInput)
	if createEntryErr != nil && !errors.IsAlreadyExists(createEntryErr) {
		return createEntryErr
	}
	return nil
}

func (u *DefaultUseCase) removeAllDependencies(ctx context.Context, hsmSlotStandardID entities.StandardID) error {
	getHSMSlotInput := GetHSMSlotInput{
		StandardID: hsmSlotStandardID,
	}
	getHSMSlotOutput, getHSMSlotErr := u.GetHSMSlot(ctx, getHSMSlotInput)
	if getHSMSlotErr != nil {
		return getHSMSlotErr
	}

	var listMyChildrenEntriesInput referentialintegrity.ListMyChildrenEntriesInput
	listMyChildrenEntriesInput.ParentResourceID = string(getHSMSlotOutput.InternalResourceID)
	listMyChildrenEntriesInput.ParentResourceKind = referentialintegrity.KindHSMSlot
	children, listMyChildrenEntriesErr := u.referentialIntegrityUseCase.ListMyChildrenEntries(ctx, listMyChildrenEntriesInput)
	if listMyChildrenEntriesErr != nil {
		return listMyChildrenEntriesErr
	}

	if len(children.Items) != 0 {
		var referencingResourcesString string
		for _, entry := range children.Items {
			referencingResourcesString += fmt.Sprintf("[id=%s,thisKind=%s]", entry.ResourceID, entry.ResourceKind)
		}
		msg := fmt.Sprintf("HSM slot '%s' can't be removed, there are elements depending on it", hsmSlotStandardID.ID)
		return errors.PreconditionFailed().WithMessage(msg).SetHumanReadableMessage(msg)
	}

	var deleteInput referentialintegrity.DeleteMyEntriesIfAnyInput
	deleteInput.ResourceID = string(getHSMSlotOutput.InternalResourceID)
	deleteInput.ResourceKind = referentialintegrity.KindHSMSlot
	deleteMyEntriesIfAnyErr := u.referentialIntegrityUseCase.DeleteMyEntriesIfAny(ctx, deleteInput)
	if listMyChildrenEntriesErr != nil {
		return deleteMyEntriesIfAnyErr
	}

	return nil

}
