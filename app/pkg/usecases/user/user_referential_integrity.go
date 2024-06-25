package user

import (
	"context"
	"fmt"

	"github.com/hyperledger-labs/signare/app/pkg/entities"
	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/application"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/referentialintegrity"
)

func (u *DefaultUserUseCase) addUserToApplicationDependency(ctx context.Context, data User) error {
	getApplicationInput := application.GetApplicationInput{
		StandardID: entities.StandardID{
			ID: data.ApplicationID,
		},
	}
	getApplicationOutput, getApplicationErr := u.applicationUseCase.GetApplication(ctx, getApplicationInput)
	if getApplicationErr != nil {
		if errors.IsNotFound(getApplicationErr) {
			msg := fmt.Sprintf("user can't be created because the application '%s' does not exist", data.ApplicationID)
			return errors.PreconditionFailed().WithMessage(msg).SetHumanReadableMessage(msg)
		}
		return getApplicationErr
	}

	var referentialIntegrityCreateEntryInput referentialintegrity.CreateEntryInput
	referentialIntegrityCreateEntryInput.ResourceID = string(data.InternalResourceID)
	referentialIntegrityCreateEntryInput.ResourceKind = referentialintegrity.KindUser
	referentialIntegrityCreateEntryInput.ParentResourceID = string(getApplicationOutput.InternalResourceID)
	referentialIntegrityCreateEntryInput.ParentResourceKind = referentialintegrity.KindApplication

	_, createEntryErr := u.referentialIntegrityUseCase.CreateEntry(ctx, referentialIntegrityCreateEntryInput)
	if createEntryErr != nil && !errors.IsAlreadyExists(createEntryErr) {
		return createEntryErr
	}
	return nil
}

func (u *DefaultUserUseCase) removeAllUserDependencies(ctx context.Context, applicationStandardID entities.ApplicationStandardID) error {
	getUserInput := GetUserInput{
		ApplicationStandardID: applicationStandardID,
	}
	getUserOutput, getUserErr := u.GetUser(ctx, getUserInput)
	if getUserErr != nil {
		return getUserErr
	}

	var listMyChildrenEntriesInput referentialintegrity.ListMyChildrenEntriesInput
	listMyChildrenEntriesInput.ParentResourceID = string(getUserOutput.InternalResourceID)
	listMyChildrenEntriesInput.ParentResourceKind = referentialintegrity.KindUser
	children, listMyChildrenEntriesErr := u.referentialIntegrityUseCase.ListMyChildrenEntries(ctx, listMyChildrenEntriesInput)
	if listMyChildrenEntriesErr != nil {
		return listMyChildrenEntriesErr
	}

	if len(children.Items) != 0 {
		var referencingResourcesString string
		for _, entry := range children.Items {
			referencingResourcesString += fmt.Sprintf("[id=%s,thisKind=%s]", entry.ResourceID, entry.ResourceKind)
		}
		msg := fmt.Sprintf("user '%s:%s' can't be removed, there are elements depending on it", applicationStandardID.ApplicationID, applicationStandardID.ID)
		return errors.PreconditionFailed().WithMessage(msg).SetHumanReadableMessage(msg)
	}

	var deleteInput referentialintegrity.DeleteMyEntriesIfAnyInput
	deleteInput.ResourceID = string(getUserOutput.InternalResourceID)
	deleteInput.ResourceKind = referentialintegrity.KindUser
	deleteMyEntriesIfAnyErr := u.referentialIntegrityUseCase.DeleteMyEntriesIfAny(ctx, deleteInput)
	if listMyChildrenEntriesErr != nil {
		return deleteMyEntriesIfAnyErr
	}

	return nil

}
