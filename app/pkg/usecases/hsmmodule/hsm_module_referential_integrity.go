package hsmmodule

import (
	"context"
	"fmt"

	"github.com/hyperledger-labs/signare/app/pkg/entities"
	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/referentialintegrity"
)

func (u *DefaultUseCase) removeAllDependencies(ctx context.Context, hsmModuleStandardID entities.StandardID) error {
	getHSMModuleInput := GetHSMModuleInput{
		StandardID: hsmModuleStandardID,
	}
	getHSMModuleOutput, getHSMModuleErr := u.GetHSMModule(ctx, getHSMModuleInput)
	if getHSMModuleErr != nil {
		return getHSMModuleErr
	}

	var listMyChildrenEntriesInput referentialintegrity.ListMyChildrenEntriesInput
	listMyChildrenEntriesInput.ParentResourceID = string(getHSMModuleOutput.InternalResourceID)
	listMyChildrenEntriesInput.ParentResourceKind = referentialintegrity.KindHSMModule
	children, listMyChildrenEntriesErr := u.referentialIntegrityUseCase.ListMyChildrenEntries(ctx, listMyChildrenEntriesInput)
	if listMyChildrenEntriesErr != nil {
		return listMyChildrenEntriesErr
	}

	if len(children.Items) != 0 {
		var referencingResourcesString string
		for _, entry := range children.Items {
			referencingResourcesString += fmt.Sprintf("[id=%s,thisKind=%s]", entry.ResourceID, entry.ResourceKind)
		}
		msg := fmt.Sprintf("HSM '%s' can't be removed, there are elements depending on it", hsmModuleStandardID.ID)
		return errors.PreconditionFailed().WithMessage(msg).SetHumanReadableMessage(msg)
	}

	var deleteInput referentialintegrity.DeleteMyEntriesIfAnyInput
	deleteInput.ResourceID = string(getHSMModuleOutput.InternalResourceID)
	deleteInput.ResourceKind = referentialintegrity.KindHSMModule
	deleteMyEntriesIfAnyErr := u.referentialIntegrityUseCase.DeleteMyEntriesIfAny(ctx, deleteInput)
	if listMyChildrenEntriesErr != nil {
		return deleteMyEntriesIfAnyErr
	}

	return nil

}
