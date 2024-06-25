// Package httpin provides the implementation of the HTTP input adapters.
package httpin

import (
	"context"
	"errors"
	"fmt"

	"github.com/hyperledger-labs/signare/app/pkg/entities"
	"github.com/hyperledger-labs/signare/app/pkg/entities/address"
	generatedhttpinfra "github.com/hyperledger-labs/signare/app/pkg/infra/generated/httpinfra"
	"github.com/hyperledger-labs/signare/app/pkg/infra/httpinfra"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/user"
	"github.com/hyperledger-labs/signare/app/pkg/utils"
)

var _ generatedhttpinfra.ApplicationAPIAdapter = new(DefaultApplicationAPIAdapter)

func (adapter *DefaultApplicationAPIAdapter) AdaptApplicationAccountsCreate(ctx context.Context, data generatedhttpinfra.ApplicationAccountsCreateRequest) (*generatedhttpinfra.ApplicationAccountsCreateResponseWrapper, *httpinfra.HTTPError) {
	addresses := make([]address.Address, len(*data.AccountCreation.Spec.Accounts))
	for i, addr := range *data.AccountCreation.Spec.Accounts {
		a, err := address.NewFromHexString(addr)
		if err != nil {
			httpError := httpinfra.NewHTTPError(httpinfra.StatusInvalidArgument).SetMessage(fmt.Sprintf("address '%s' is not a valid hex address", addr))
			return nil, httpError
		}
		addresses[i] = a
	}

	input := user.EnableAccountsInput{
		UserID:        data.UserId,
		ApplicationID: data.ApplicationId,
		Addresses:     addresses,
	}
	out, err := adapter.userUseCase.EnableAccounts(ctx, input)
	if err != nil {
		return nil, httpinfra.NewHTTPErrorFromUseCaseError(ctx, err)
	}

	return &generatedhttpinfra.ApplicationAccountsCreateResponseWrapper{
		UserDetail: mapUser(out.User),
		ResponseInfo: httpinfra.ResponseInfo{
			ResponseType: httpinfra.ResponseTypeCreated,
		},
	}, nil
}

func (adapter *DefaultApplicationAPIAdapter) AdaptApplicationAccountsRemove(ctx context.Context, data generatedhttpinfra.ApplicationAccountsRemoveRequest) (*generatedhttpinfra.ApplicationAccountsRemoveResponseWrapper, *httpinfra.HTTPError) {
	addr, err := address.NewFromHexString(data.AccountId)
	if err != nil {
		httpError := httpinfra.NewHTTPError(httpinfra.StatusInvalidArgument).SetMessage(fmt.Sprintf("address '%s' is not a valid hex address", data.AccountId))
		return nil, httpError
	}

	input := user.DisableAccountInput{
		UserID:        data.UserId,
		ApplicationID: data.ApplicationId,
		Address:       addr,
	}
	out, err := adapter.userUseCase.DisableAccount(ctx, input)
	if err != nil {
		return nil, httpinfra.NewHTTPErrorFromUseCaseError(ctx, err)
	}

	return &generatedhttpinfra.ApplicationAccountsRemoveResponseWrapper{
		UserDetail: mapUser(out.User),
		ResponseInfo: httpinfra.ResponseInfo{
			ResponseType: httpinfra.ResponseTypeOk,
		},
	}, nil
}

func (adapter *DefaultApplicationAPIAdapter) AdaptApplicationUsersCreate(ctx context.Context, data generatedhttpinfra.ApplicationUsersCreateRequest) (*generatedhttpinfra.ApplicationUsersCreateResponseWrapper, *httpinfra.HTTPError) {
	input := user.CreateUserInput{
		ApplicationID: data.ApplicationId,
	}
	if data.UserCreation.Meta != nil && data.UserCreation.Meta.Id != nil {
		input.ID = data.UserCreation.Meta.Id
	}
	if data.UserCreation.Spec != nil && data.UserCreation.Spec.Roles != nil {
		input.Roles = *data.UserCreation.Spec.Roles
	}
	if data.UserCreation.Spec != nil && data.UserCreation.Spec.Description != nil {
		input.Description = data.UserCreation.Spec.Description
	}

	out, err := adapter.userUseCase.CreateUser(ctx, input)
	if err != nil {
		return nil, httpinfra.NewHTTPErrorFromUseCaseError(ctx, err)
	}

	return &generatedhttpinfra.ApplicationUsersCreateResponseWrapper{
		UserDetail: mapUser(out.User),
		ResponseInfo: httpinfra.ResponseInfo{
			ResponseType: httpinfra.ResponseTypeCreated,
		},
	}, nil
}

func (adapter *DefaultApplicationAPIAdapter) AdaptApplicationUsersDescribe(ctx context.Context, data generatedhttpinfra.ApplicationUsersDescribeRequest) (*generatedhttpinfra.ApplicationUsersDescribeResponseWrapper, *httpinfra.HTTPError) {
	input := user.GetUserInput{
		ApplicationStandardID: entities.ApplicationStandardID{
			ID:            data.UserId,
			ApplicationID: data.ApplicationId,
		},
	}
	out, err := adapter.userUseCase.GetUser(ctx, input)
	if err != nil {
		return nil, httpinfra.NewHTTPErrorFromUseCaseError(ctx, err)
	}

	return &generatedhttpinfra.ApplicationUsersDescribeResponseWrapper{
		UserDetail: mapUser(out.User),
		ResponseInfo: httpinfra.ResponseInfo{
			ResponseType: httpinfra.ResponseTypeOk,
		},
	}, nil
}

func (adapter *DefaultApplicationAPIAdapter) AdaptApplicationUsersEdit(ctx context.Context, data generatedhttpinfra.ApplicationUsersEditRequest) (*generatedhttpinfra.ApplicationUsersEditResponseWrapper, *httpinfra.HTTPError) {
	input := user.EditUserInput{
		ApplicationStandardID: entities.ApplicationStandardID{
			ID:            data.UserId,
			ApplicationID: data.ApplicationId,
		},
		ResourceVersion: *data.UserUpdate.Meta.ResourceVersion,
	}
	if data.UserUpdate.Spec != nil && data.UserUpdate.Spec.Roles != nil {
		input.Roles = *data.UserUpdate.Spec.Roles
	}
	if data.UserUpdate.Spec != nil && data.UserUpdate.Spec.Description != nil {
		input.Description = data.UserUpdate.Spec.Description
	}
	out, err := adapter.userUseCase.EditUser(ctx, input)
	if err != nil {
		return nil, httpinfra.NewHTTPErrorFromUseCaseError(ctx, err)
	}

	return &generatedhttpinfra.ApplicationUsersEditResponseWrapper{
		UserDetail: mapUser(out.User),
		ResponseInfo: httpinfra.ResponseInfo{
			ResponseType: httpinfra.ResponseTypeOk,
		},
	}, nil
}

func (adapter *DefaultApplicationAPIAdapter) AdaptApplicationUsersList(ctx context.Context, data generatedhttpinfra.ApplicationUsersListRequest) (*generatedhttpinfra.ApplicationUsersListResponseWrapper, *httpinfra.HTTPError) {
	input := user.ListUsersInput{
		ApplicationID: data.ApplicationId,
	}
	var limitInput int
	if data.Limit != nil {
		limitInput = int(*data.Limit)
	}
	var offsetInput int
	if data.Offset != nil {
		offsetInput = int(*data.Offset)
	}
	pageLimit := utils.MaxValue(utils.DefaultIntValue(limitInput, defaultApplicationListLimit), maxListApplicationLimit)
	input.PageLimit = pageLimit
	input.PageOffset = offsetInput
	input.OrderBy = data.OrderBy
	input.OrderDirection = data.OrderDirection

	out, err := adapter.userUseCase.ListUsers(ctx, input)
	if err != nil {
		return nil, httpinfra.NewHTTPErrorFromUseCaseError(ctx, err)
	}

	adaptedItems := make([]generatedhttpinfra.UserDetail, len(out.Items))
	for i, item := range out.Items {
		adaptedItems[i] = mapUser(item)
	}

	offset := int32(out.Offset)
	limit := int32(out.Limit)
	return &generatedhttpinfra.ApplicationUsersListResponseWrapper{
		UserCollection: generatedhttpinfra.UserCollection{
			Limit:     &limit,
			Offset:    &offset,
			MoreItems: &out.MoreItems,
			Items:     &adaptedItems,
		},
		ResponseInfo: httpinfra.ResponseInfo{
			ResponseType: httpinfra.ResponseTypeOk,
		},
	}, nil
}

func (adapter *DefaultApplicationAPIAdapter) AdaptApplicationUsersRemove(ctx context.Context, data generatedhttpinfra.ApplicationUsersRemoveRequest) (*generatedhttpinfra.ApplicationUsersRemoveResponseWrapper, *httpinfra.HTTPError) {
	input := user.DeleteUserInput{
		ApplicationStandardID: entities.ApplicationStandardID{
			ID:            data.UserId,
			ApplicationID: data.ApplicationId,
		},
	}

	out, err := adapter.userUseCase.DeleteUser(ctx, input)
	if err != nil {
		return nil, httpinfra.NewHTTPErrorFromUseCaseError(ctx, err)
	}

	return &generatedhttpinfra.ApplicationUsersRemoveResponseWrapper{
		UserDetail: mapUser(out.User),
		ResponseInfo: httpinfra.ResponseInfo{
			ResponseType: httpinfra.ResponseTypeOk,
		},
	}, nil
}

// DefaultApplicationAPIAdapter implements ApplicationAPIAdapter.
type DefaultApplicationAPIAdapter struct {
	userUseCase user.UserUseCase
}

// DefaultApplicationAPIAdapterOptions options to create a new DefaultApplicationAPIAdapter.
type DefaultApplicationAPIAdapterOptions struct {
	UserUseCase user.UserUseCase
}

// ProvideDefaultApplicationAPIAdapter creates a new DefaultApplicationAPIAdapter instance.
func ProvideDefaultApplicationAPIAdapter(options DefaultApplicationAPIAdapterOptions) (*DefaultApplicationAPIAdapter, error) {
	if options.UserUseCase == nil {
		return nil, errors.New("mandatory 'UserUseCase' was not provided")
	}

	return &DefaultApplicationAPIAdapter{
		userUseCase: options.UserUseCase,
	}, nil
}

func mapUser(user user.User) generatedhttpinfra.UserDetail {
	creationDate := user.CreationDate.String()
	lastUpdate := user.LastUpdate.String()
	accounts := make([]string, len(user.Accounts))
	for i, acc := range user.Accounts {
		accounts[i] = acc.Address.String()
	}

	return generatedhttpinfra.UserDetail{
		Meta: &generatedhttpinfra.ResourceMetaDetail{
			Id:              &user.ID,
			ResourceVersion: &user.ResourceVersion,
			CreationDate:    &creationDate,
			LastUpdate:      &lastUpdate,
		},
		Spec: &generatedhttpinfra.UserDetailSpec{
			Roles:       &user.Roles,
			Accounts:    &accounts,
			Description: user.Description,
		},
	}
}
