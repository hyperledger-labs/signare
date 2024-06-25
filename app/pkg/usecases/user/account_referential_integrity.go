package user

import (
	"context"
	"fmt"

	"github.com/hyperledger-labs/signare/app/pkg/entities"
	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/referentialintegrity"
)

func (u *DefaultUserUseCase) addAccountToUserDependency(ctx context.Context, data Account) error {
	getUserInput := GetUserInput{
		ApplicationStandardID: entities.ApplicationStandardID{
			ID:            data.UserID,
			ApplicationID: data.ApplicationID,
		},
	}
	getUserOutput, getUserErr := u.GetUser(ctx, getUserInput)
	if getUserErr != nil {
		if errors.IsNotFound(getUserErr) {
			msg := fmt.Sprintf("can't create account '%s' because user '%s' does not exist", data.Address, data.UserID)
			return errors.PreconditionFailed().WithMessage(msg).SetHumanReadableMessage(msg)
		}
		return getUserErr
	}

	var referentialIntegrityCreateEntryInput referentialintegrity.CreateEntryInput
	referentialIntegrityCreateEntryInput.ResourceID = string(data.InternalResourceID)
	referentialIntegrityCreateEntryInput.ResourceKind = referentialintegrity.KindAccount
	referentialIntegrityCreateEntryInput.ParentResourceID = string(getUserOutput.InternalResourceID)
	referentialIntegrityCreateEntryInput.ParentResourceKind = referentialintegrity.KindUser

	_, createEntryErr := u.referentialIntegrityUseCase.CreateEntry(ctx, referentialIntegrityCreateEntryInput)
	if createEntryErr != nil && !errors.IsAlreadyExists(createEntryErr) {
		return createEntryErr
	}
	return nil
}

func (u *DefaultUserUseCase) removeAllAccountDependencies(ctx context.Context, accountID AccountID) error {
	getAccountInput := GetAccountInput{
		AccountID: accountID,
	}
	getAccountOutput, getAccountErr := u.GetAccount(ctx, getAccountInput)
	if getAccountErr != nil {
		return getAccountErr
	}
	var listMyChildrenEntriesInput referentialintegrity.ListMyChildrenEntriesInput
	listMyChildrenEntriesInput.ParentResourceID = string(getAccountOutput.InternalResourceID)
	listMyChildrenEntriesInput.ParentResourceKind = referentialintegrity.KindAccount
	children, listMyChildrenEntriesErr := u.referentialIntegrityUseCase.ListMyChildrenEntries(ctx, listMyChildrenEntriesInput)
	if listMyChildrenEntriesErr != nil {
		return listMyChildrenEntriesErr
	}

	if len(children.Items) != 0 {
		var referencingResourcesString string
		for _, entry := range children.Items {
			referencingResourcesString += fmt.Sprintf("[id=%s,thisKind=%s]", entry.ResourceID, entry.ResourceKind)
		}
		msg := fmt.Sprintf("acount '%s:%s:%s' can't be removed, there are elements depending on it", accountID.ApplicationID, accountID.UserID, accountID.Address)
		return errors.PreconditionFailed().WithMessage(msg).SetHumanReadableMessage(msg)
	}

	var deleteInput referentialintegrity.DeleteMyEntriesIfAnyInput
	deleteInput.ResourceID = string(getAccountOutput.InternalResourceID)
	deleteInput.ResourceKind = referentialintegrity.KindAccount
	deleteMyEntriesIfAnyErr := u.referentialIntegrityUseCase.DeleteMyEntriesIfAny(ctx, deleteInput)
	if deleteMyEntriesIfAnyErr != nil {
		return deleteMyEntriesIfAnyErr
	}

	return nil

}
