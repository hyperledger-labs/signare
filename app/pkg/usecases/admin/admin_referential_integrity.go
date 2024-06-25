package admin

import (
	"context"
	"fmt"

	"github.com/hyperledger-labs/signare/app/pkg/entities"
	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/referentialintegrity"
)

func (u *DefaultUseCase) removeAllDependencies(ctx context.Context, adminStandardID entities.StandardID) error {
	getAdminInput := GetAdminInput{
		StandardID: adminStandardID,
	}
	getAdminOutput, getAdminErr := u.GetAdmin(ctx, getAdminInput)
	if getAdminErr != nil {
		return getAdminErr
	}

	var listMyChildrenEntriesInput referentialintegrity.ListMyChildrenEntriesInput
	listMyChildrenEntriesInput.ParentResourceID = string(getAdminOutput.InternalResourceID)
	listMyChildrenEntriesInput.ParentResourceKind = referentialintegrity.KindAdmin
	children, listMyChildrenEntriesErr := u.referentialIntegrityUseCase.ListMyChildrenEntries(ctx, listMyChildrenEntriesInput)
	if listMyChildrenEntriesErr != nil {
		return listMyChildrenEntriesErr
	}

	if len(children.Items) != 0 {
		var referencingResourcesString string
		for _, entry := range children.Items {
			referencingResourcesString += fmt.Sprintf("[id=%s,thisKind=%s]", entry.ResourceID, entry.ResourceKind)
		}
		msg := fmt.Sprintf("admin '%s' can't be removed, there are elements depending on it", adminStandardID.ID)
		return errors.PreconditionFailed().WithMessage(msg).SetHumanReadableMessage(msg)
	}

	var deleteInput referentialintegrity.DeleteMyEntriesIfAnyInput
	deleteInput.ResourceID = string(getAdminOutput.InternalResourceID)
	deleteInput.ResourceKind = referentialintegrity.KindAdmin
	deleteMyEntriesIfAnyErr := u.referentialIntegrityUseCase.DeleteMyEntriesIfAny(ctx, deleteInput)
	if listMyChildrenEntriesErr != nil {
		return deleteMyEntriesIfAnyErr
	}

	return nil

}
