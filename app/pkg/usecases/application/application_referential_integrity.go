package application

import (
	"context"
	"fmt"

	"github.com/hyperledger-labs/signare/app/pkg/entities"
	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/referentialintegrity"
)

func (u *DefaultUseCase) removeAllDependencies(ctx context.Context, applicationStandardID entities.StandardID) error {
	getApplicationInput := GetApplicationInput{
		StandardID: applicationStandardID,
	}
	getApplication, getApplicationErr := u.GetApplication(ctx, getApplicationInput)
	if getApplicationErr != nil {
		return getApplicationErr
	}

	internalResourceID := getApplication.InternalResourceID
	var listMyChildrenEntriesInput referentialintegrity.ListMyChildrenEntriesInput
	listMyChildrenEntriesInput.ParentResourceID = string(internalResourceID)
	listMyChildrenEntriesInput.ParentResourceKind = referentialintegrity.KindApplication
	children, listMyChildrenEntriesErr := u.referentialIntegrityUseCase.ListMyChildrenEntries(ctx, listMyChildrenEntriesInput)
	if listMyChildrenEntriesErr != nil {
		return listMyChildrenEntriesErr
	}

	if len(children.Items) != 0 {
		var referencingResourcesString string
		for _, entry := range children.Items {
			referencingResourcesString += fmt.Sprintf("[id=%s,thisKind=%s]", entry.ResourceID, entry.ResourceKind)
		}
		msg := fmt.Sprintf("application '%s' can't be removed, there are elements depending on it", applicationStandardID.ID)
		return errors.PreconditionFailed().WithMessage(msg).SetHumanReadableMessage(msg)
	}

	var deleteInput referentialintegrity.DeleteMyEntriesIfAnyInput
	deleteInput.ResourceID = string(internalResourceID)
	deleteInput.ResourceKind = referentialintegrity.KindApplication
	deleteMyEntriesIfAnyErr := u.referentialIntegrityUseCase.DeleteMyEntriesIfAny(ctx, deleteInput)
	if listMyChildrenEntriesErr != nil {
		return deleteMyEntriesIfAnyErr
	}

	return nil

}
