// Package httpin provides the implementation of the HTTP input adapters.
package httpin

import (
	"context"
	"errors"
	"fmt"

	"github.com/hyperledger-labs/signare/app/pkg/entities"
	generatedhttpinfra "github.com/hyperledger-labs/signare/app/pkg/infra/generated/httpinfra"
	"github.com/hyperledger-labs/signare/app/pkg/infra/httpinfra"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/admin"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/application"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/hsmmodule"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/hsmslot"
	"github.com/hyperledger-labs/signare/app/pkg/utils"
)

const (
	defaultApplicationListLimit int = 30
	maxListApplicationLimit     int = 100
	defaultAdminUserListLimit   int = 30
	maxListAdminUserLimit       int = 100
)

var _ generatedhttpinfra.AdminAPIAdapter = new(DefaultAdminAPIAdapter)

/************************/
/*    Applications     */
/**********************/

func (adapter *DefaultAdminAPIAdapter) AdaptAdminApplicationsCreate(ctx context.Context, data generatedhttpinfra.AdminApplicationsCreateRequest) (*generatedhttpinfra.AdminApplicationsCreateResponseWrapper, *httpinfra.HTTPError) {
	chainID, err := entities.NewInt256FromString(*data.ApplicationCreation.Spec.ChainId)
	if err != nil {
		return nil, httpinfra.NewHTTPErrorFromError(ctx, err, httpinfra.StatusInvalidArgument)
	}

	input := application.CreateApplicationInput{}
	input.ChainID = *chainID

	if data.ApplicationCreation.Meta != nil && data.ApplicationCreation.Meta.Id != nil {
		input.ID = data.ApplicationCreation.Meta.Id
	}
	if data.ApplicationCreation.Spec != nil && data.ApplicationCreation.Spec.Description != nil {
		input.Description = data.ApplicationCreation.Spec.Description
	}
	out, err := adapter.applicationUseCase.CreateApplication(ctx, input)
	if err != nil {
		return nil, httpinfra.NewHTTPErrorFromUseCaseError(ctx, err)
	}

	response := generatedhttpinfra.AdminApplicationsCreateResponseWrapper{
		ApplicationDetail: mapApplicationOut(out.Application),
		ResponseInfo: httpinfra.ResponseInfo{
			ResponseType: httpinfra.ResponseTypeCreated,
		},
	}
	return &response, nil
}

func (adapter *DefaultAdminAPIAdapter) AdaptAdminApplicationsDescribe(ctx context.Context, data generatedhttpinfra.AdminApplicationsDescribeRequest) (*generatedhttpinfra.AdminApplicationsDescribeResponseWrapper, *httpinfra.HTTPError) {
	input := application.GetApplicationInput{
		StandardID: entities.StandardID{
			ID: data.ApplicationId,
		},
	}
	out, err := adapter.applicationUseCase.GetApplication(ctx, input)
	if err != nil {
		return nil, httpinfra.NewHTTPErrorFromUseCaseError(ctx, err)
	}

	response := generatedhttpinfra.AdminApplicationsDescribeResponseWrapper{
		ApplicationDetail: mapApplicationOut(out.Application),
		ResponseInfo: httpinfra.ResponseInfo{
			ResponseType: httpinfra.ResponseTypeOk,
		},
	}
	return &response, nil
}

func (adapter *DefaultAdminAPIAdapter) AdaptAdminApplicationsEdit(ctx context.Context, data generatedhttpinfra.AdminApplicationsEditRequest) (*generatedhttpinfra.AdminApplicationsEditResponseWrapper, *httpinfra.HTTPError) {
	chainID, err := entities.NewInt256FromString(*data.ApplicationUpdate.Spec.ChainId)
	if err != nil {
		return nil, httpinfra.NewHTTPErrorFromError(ctx, err, httpinfra.StatusInvalidArgument)
	}

	input := application.EditApplicationInput{
		ID:              data.ApplicationId,
		ResourceVersion: *data.ApplicationUpdate.Meta.ResourceVersion,
		ChainID:         chainID,
		Description:     data.ApplicationUpdate.Spec.Description,
	}
	out, err := adapter.applicationUseCase.EditApplication(ctx, input)
	if err != nil {
		return nil, httpinfra.NewHTTPErrorFromUseCaseError(ctx, err)
	}

	response := generatedhttpinfra.AdminApplicationsEditResponseWrapper{
		ApplicationDetail: mapApplicationOut(out.Application),
		ResponseInfo: httpinfra.ResponseInfo{
			ResponseType: httpinfra.ResponseTypeOk,
		},
	}
	return &response, nil
}

func (adapter *DefaultAdminAPIAdapter) AdaptAdminApplicationsList(ctx context.Context, data generatedhttpinfra.AdminApplicationsListRequest) (*generatedhttpinfra.AdminApplicationsListResponseWrapper, *httpinfra.HTTPError) {
	var input application.ListApplicationsInput

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

	outputData, err := adapter.applicationUseCase.ListApplications(ctx, input)
	if err != nil {
		return nil, httpinfra.NewHTTPErrorFromUseCaseError(ctx, err)
	}

	adaptedUseCaseCollection := make([]generatedhttpinfra.ApplicationDetail, len(outputData.Items))
	for i, outItem := range outputData.Items {
		item := outItem
		adaptedUseCaseCollection[i] = mapApplicationOut(item)
	}
	offset := int32(outputData.Offset)
	limit := int32(outputData.Limit)
	response := generatedhttpinfra.AdminApplicationsListResponseWrapper{
		ApplicationCollection: generatedhttpinfra.ApplicationCollection{
			Items:     &adaptedUseCaseCollection,
			Offset:    &offset,
			Limit:     &limit,
			MoreItems: &outputData.MoreItems,
		},
		ResponseInfo: httpinfra.ResponseInfo{
			ResponseType: httpinfra.ResponseTypeOk,
		},
	}
	return &response, nil
}

func (adapter *DefaultAdminAPIAdapter) AdaptAdminApplicationsRemove(ctx context.Context, data generatedhttpinfra.AdminApplicationsRemoveRequest) (*generatedhttpinfra.AdminApplicationsRemoveResponseWrapper, *httpinfra.HTTPError) {
	input := application.DeleteApplicationInput{
		StandardID: entities.StandardID{
			ID: data.ApplicationId,
		},
	}
	out, err := adapter.applicationUseCase.DeleteApplication(ctx, input)
	if err != nil {
		return nil, httpinfra.NewHTTPErrorFromUseCaseError(ctx, err)
	}

	response := generatedhttpinfra.AdminApplicationsRemoveResponseWrapper{
		ApplicationDetail: mapApplicationOut(out.Application),
		ResponseInfo: httpinfra.ResponseInfo{
			ResponseType: httpinfra.ResponseTypeOk,
		},
	}
	return &response, nil
}

func mapApplicationOut(in application.Application) generatedhttpinfra.ApplicationDetail {
	creationDate := in.CreationDate.String()
	lastUpdate := in.LastUpdate.String()
	chainID := in.ChainID.String()
	return generatedhttpinfra.ApplicationDetail{
		Meta: &generatedhttpinfra.ResourceMetaDetail{
			Id:              &in.ID,
			ResourceVersion: &in.ResourceVersion,
			CreationDate:    &creationDate,
			LastUpdate:      &lastUpdate,
		},
		Spec: &generatedhttpinfra.ApplicationDetailSpec{
			ChainId:     &chainID,
			Description: in.Description,
		},
	}
}

/*******************/
/*    Modules     */
/*****************/

func (adapter *DefaultAdminAPIAdapter) AdaptAdminModulesCreate(ctx context.Context, request generatedhttpinfra.AdminModulesCreateRequest) (*generatedhttpinfra.AdminModulesCreateResponseWrapper, *httpinfra.HTTPError) {
	var input hsmmodule.CreateHSMModuleInput
	if request.ModuleCreation.Meta != nil && request.ModuleCreation.Meta.Id != nil {
		input.ID = request.ModuleCreation.Meta.Id
	}

	if request.ModuleCreation.Spec != nil {
		if request.ModuleCreation.Spec.Description != nil {
			input.Description = request.ModuleCreation.Spec.Description
		}

		if request.ModuleCreation.Spec.Configuration != nil {
			moduleKind, err := mapUseCaseHSMKindFrom(request.ModuleCreation.Spec.Configuration.HsmKind)
			if err != nil {
				return nil, err
			}
			input.ModuleKind = *moduleKind
			if request.ModuleCreation.Spec.Configuration.HsmKind == generatedhttpinfra.HsmKindSofthsm {
				input.Configuration.SoftHSMConfiguration = &hsmmodule.SoftHSMConfiguration{}
			}
		}
	}
	out, err := adapter.hsmUseCase.CreateHSMModule(ctx, input)
	if err != nil {
		return nil, httpinfra.NewHTTPErrorFromUseCaseError(ctx, err)
	}

	module, err := mapModule(out.HSMModule)
	if err != nil {
		return nil, httpinfra.NewHTTPErrorFromError(ctx, err, httpinfra.StatusInternal)
	}

	response := generatedhttpinfra.AdminModulesCreateResponseWrapper{
		ModuleDetail: *module,
		ResponseInfo: httpinfra.ResponseInfo{
			ResponseType: httpinfra.ResponseTypeCreated,
		},
	}
	return &response, nil
}

func (adapter *DefaultAdminAPIAdapter) AdaptAdminModulesDescribe(ctx context.Context, request generatedhttpinfra.AdminModulesDescribeRequest) (*generatedhttpinfra.AdminModulesDescribeResponseWrapper, *httpinfra.HTTPError) {
	var hsmModuleStandardID entities.StandardID
	hsmModuleStandardID.ID = request.ModuleId
	input := hsmmodule.GetHSMModuleInput{
		StandardID: hsmModuleStandardID,
	}
	out, err := adapter.hsmUseCase.GetHSMModule(ctx, input)
	if err != nil {
		return nil, httpinfra.NewHTTPErrorFromUseCaseError(ctx, err)
	}

	module, err := mapModule(out.HSMModule)
	if err != nil {
		return nil, httpinfra.NewHTTPErrorFromError(ctx, err, httpinfra.StatusInternal)
	}

	response := generatedhttpinfra.AdminModulesDescribeResponseWrapper{
		ModuleDetail: *module,
		ResponseInfo: httpinfra.ResponseInfo{
			ResponseType: httpinfra.ResponseTypeCreated,
		},
	}
	return &response, nil
}

func (adapter *DefaultAdminAPIAdapter) AdaptAdminModulesEdit(ctx context.Context, request generatedhttpinfra.AdminModulesEditRequest) (*generatedhttpinfra.AdminModulesEditResponseWrapper, *httpinfra.HTTPError) {
	input := hsmmodule.EditHSMModuleInput{
		HSMModule: hsmmodule.HSMModule{
			StandardResourceMeta: entities.StandardResourceMeta{
				StandardResource: entities.StandardResource{
					StandardID: entities.StandardID{
						ID: request.ModuleId,
					},
				},
			},
		},
	}

	if request.ModuleUpdate.Meta != nil && request.ModuleUpdate.Meta.ResourceVersion != nil {
		input.ResourceVersion = *request.ModuleUpdate.Meta.ResourceVersion
	}

	if request.ModuleUpdate.Spec != nil {
		if request.ModuleUpdate.Spec.Description != nil {
			input.Description = request.ModuleUpdate.Spec.Description
		}
		if request.ModuleUpdate.Spec.Configuration != nil {
			moduleKind, err := mapUseCaseHSMKindFrom(request.ModuleUpdate.Spec.Configuration.HsmKind)
			if err != nil {
				return nil, err
			}
			input.Kind = *moduleKind
			if request.ModuleUpdate.Spec.Configuration.HsmKind == generatedhttpinfra.HsmKindSofthsm {
				input.Configuration.SoftHSMConfiguration = &hsmmodule.SoftHSMConfiguration{}
			}
		}
	}

	out, err := adapter.hsmUseCase.EditHSMModule(ctx, input)
	if err != nil {
		return nil, httpinfra.NewHTTPErrorFromUseCaseError(ctx, err)
	}

	module, err := mapModule(out.HSMModule)
	if err != nil {
		return nil, httpinfra.NewHTTPErrorFromError(ctx, err, httpinfra.StatusInternal)
	}

	response := generatedhttpinfra.AdminModulesEditResponseWrapper{
		ModuleDetail: *module,
		ResponseInfo: httpinfra.ResponseInfo{
			ResponseType: httpinfra.ResponseTypeCreated,
		},
	}
	return &response, nil
}

func (adapter *DefaultAdminAPIAdapter) AdaptAdminModulesList(ctx context.Context, request generatedhttpinfra.AdminModulesListRequest) (*generatedhttpinfra.AdminModulesListResponseWrapper, *httpinfra.HTTPError) {
	var input hsmmodule.ListHSMModulesInput

	var limitInput int
	if request.Limit != nil {
		limitInput = int(*request.Limit)
	}
	var offsetInput int
	if request.Offset != nil {
		offsetInput = int(*request.Offset)
	}
	pageLimit := utils.MaxValue(utils.DefaultIntValue(limitInput, defaultAdminUserListLimit), maxListAdminUserLimit)
	input.PageLimit = pageLimit
	input.PageOffset = offsetInput
	input.OrderBy = request.OrderBy
	input.OrderDirection = request.OrderDirection

	outputData, err := adapter.hsmUseCase.ListHSMModules(ctx, input)
	if err != nil {
		return nil, httpinfra.NewHTTPErrorFromUseCaseError(ctx, err)
	}

	adaptedUseCaseCollection := make([]generatedhttpinfra.ModuleDetail, len(outputData.Items))
	for i, outItem := range outputData.Items {
		item := outItem
		m, mapErr := mapModule(item)
		if mapErr != nil {
			return nil, httpinfra.NewHTTPError(httpinfra.StatusInternal)
		}
		adaptedUseCaseCollection[i] = *m
	}
	offset := int32(outputData.Offset)
	limit := int32(outputData.Limit)
	response := generatedhttpinfra.AdminModulesListResponseWrapper{
		ModuleCollection: generatedhttpinfra.ModuleCollection{
			Items:     &adaptedUseCaseCollection,
			Offset:    &offset,
			Limit:     &limit,
			MoreItems: &outputData.MoreItems,
		},
		ResponseInfo: httpinfra.ResponseInfo{
			ResponseType: httpinfra.ResponseTypeOk,
		},
	}
	return &response, nil
}

func (adapter *DefaultAdminAPIAdapter) AdaptAdminModulesRemove(ctx context.Context, request generatedhttpinfra.AdminModulesRemoveRequest) (*generatedhttpinfra.AdminModulesRemoveResponseWrapper, *httpinfra.HTTPError) {
	input := hsmmodule.DeleteHSMModuleInput{
		StandardID: entities.StandardID{
			ID: request.ModuleId,
		},
	}
	out, err := adapter.hsmUseCase.DeleteHSMModule(ctx, input)
	if err != nil {
		return nil, httpinfra.NewHTTPErrorFromUseCaseError(ctx, err)
	}

	module, err := mapModule(out.HSMModule)
	if err != nil {
		return nil, httpinfra.NewHTTPErrorFromError(ctx, err, httpinfra.StatusInternal)
	}

	response := generatedhttpinfra.AdminModulesRemoveResponseWrapper{
		ModuleDetail: *module,
		ResponseInfo: httpinfra.ResponseInfo{
			ResponseType: httpinfra.ResponseTypeCreated,
		},
	}
	return &response, nil
}

/*******************/
/*     Slots      */
/*****************/

func (adapter *DefaultAdminAPIAdapter) AdaptAdminSlotsCreate(ctx context.Context, data generatedhttpinfra.AdminSlotsCreateRequest) (*generatedhttpinfra.AdminSlotsCreateResponseWrapper, *httpinfra.HTTPError) {
	input := hsmslot.CreateHSMSlotInput{
		HSMModuleID: data.ModuleId,
	}

	if data.SlotCreation.Meta != nil && data.SlotCreation.Meta.Id != nil {
		input.ID = data.SlotCreation.Meta.Id
	}
	if data.SlotCreation.Spec != nil && data.SlotCreation.Spec.ApplicationId != nil {
		input.ApplicationID = *data.SlotCreation.Spec.ApplicationId
	}
	if data.SlotCreation.Spec != nil && data.SlotCreation.Spec.Slot != nil {
		input.Slot = *data.SlotCreation.Spec.Slot
	}
	if data.SlotCreation.Spec != nil && data.SlotCreation.Spec.Pin != nil {
		input.Pin = *data.SlotCreation.Spec.Pin
	}

	out, err := adapter.hsmSlotUseCase.CreateHSMSlot(ctx, input)
	if err != nil {
		return nil, httpinfra.NewHTTPErrorFromUseCaseError(ctx, err)
	}

	return &generatedhttpinfra.AdminSlotsCreateResponseWrapper{
		SlotDetail: mapSlot(out.HSMSlot),
		ResponseInfo: httpinfra.ResponseInfo{
			ResponseType: httpinfra.ResponseTypeCreated,
		},
	}, nil
}

func (adapter *DefaultAdminAPIAdapter) AdaptAdminSlotsDescribe(ctx context.Context, data generatedhttpinfra.AdminSlotsDescribeRequest) (*generatedhttpinfra.AdminSlotsDescribeResponseWrapper, *httpinfra.HTTPError) {
	input := hsmslot.GetHSMSlotInput{
		StandardID: entities.StandardID{
			ID: data.SlotId,
		},
	}

	out, err := adapter.hsmSlotUseCase.GetHSMSlot(ctx, input)
	if err != nil {
		return nil, httpinfra.NewHTTPErrorFromUseCaseError(ctx, err)
	}

	return &generatedhttpinfra.AdminSlotsDescribeResponseWrapper{
		SlotDetail: mapSlot(out.HSMSlot),
		ResponseInfo: httpinfra.ResponseInfo{
			ResponseType: httpinfra.ResponseTypeOk,
		},
	}, nil
}

func (adapter *DefaultAdminAPIAdapter) AdaptAdminSlotsList(ctx context.Context, data generatedhttpinfra.AdminSlotsListRequest) (*generatedhttpinfra.AdminSlotsListResponseWrapper, *httpinfra.HTTPError) {
	input := hsmslot.ListHSMSlotsByHSMModuleInput{
		HSMModuleID: entities.StandardID{ID: data.ModuleId},
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

	if len(data.ApplicationId) > 0 {
		input.ApplicationID = &data.ApplicationId
	}

	out, err := adapter.hsmSlotUseCase.ListHSMSlotsByHSMModule(ctx, input)
	if err != nil {
		return nil, httpinfra.NewHTTPErrorFromUseCaseError(ctx, err)
	}

	adaptedItems := make([]generatedhttpinfra.SlotDetail, len(out.Items))
	for i, item := range out.Items {
		adaptedItems[i] = mapSlot(item)
	}

	offset := int32(out.Offset)
	limit := int32(out.Limit)
	return &generatedhttpinfra.AdminSlotsListResponseWrapper{
		SlotCollection: generatedhttpinfra.SlotCollection{
			Items:     &adaptedItems,
			Limit:     &limit,
			Offset:    &offset,
			MoreItems: &out.MoreItems,
		},
		ResponseInfo: httpinfra.ResponseInfo{
			ResponseType: httpinfra.ResponseTypeOk,
		},
	}, nil
}

func (adapter *DefaultAdminAPIAdapter) AdaptAdminSlotsRemove(ctx context.Context, data generatedhttpinfra.AdminSlotsRemoveRequest) (*generatedhttpinfra.AdminSlotsRemoveResponseWrapper, *httpinfra.HTTPError) {
	input := hsmslot.DeleteHSMSlotInput{
		StandardID: entities.StandardID{
			ID: data.SlotId,
		},
	}

	out, err := adapter.hsmSlotUseCase.DeleteHSMSlot(ctx, input)
	if err != nil {
		return nil, httpinfra.NewHTTPErrorFromUseCaseError(ctx, err)
	}

	return &generatedhttpinfra.AdminSlotsRemoveResponseWrapper{
		SlotDetail: mapSlot(out.HSMSlot),
		ResponseInfo: httpinfra.ResponseInfo{
			ResponseType: httpinfra.ResponseTypeOk,
		},
	}, nil
}

func (adapter *DefaultAdminAPIAdapter) AdaptAdminSlotsUpdatePin(ctx context.Context, data generatedhttpinfra.AdminSlotsUpdatePinRequest) (*generatedhttpinfra.AdminSlotsUpdatePinResponseWrapper, *httpinfra.HTTPError) {
	input := hsmslot.EditPinInput{
		StandardID: entities.StandardID{
			ID: data.SlotId,
		},
		HSMModuleID: data.ModuleId,
	}
	if data.SlotUpdatePin.Meta != nil && data.SlotUpdatePin.Meta.ResourceVersion != nil {
		input.ResourceVersion = *data.SlotUpdatePin.Meta.ResourceVersion
	}
	if data.SlotUpdatePin.Spec != nil && data.SlotUpdatePin.Spec.Pin != nil {
		input.Pin = *data.SlotUpdatePin.Spec.Pin
	}

	out, err := adapter.hsmSlotUseCase.EditPin(ctx, input)
	if err != nil {
		return nil, httpinfra.NewHTTPErrorFromUseCaseError(ctx, err)
	}

	return &generatedhttpinfra.AdminSlotsUpdatePinResponseWrapper{
		SlotDetail: mapSlot(out.HSMSlot),
		ResponseInfo: httpinfra.ResponseInfo{
			ResponseType: httpinfra.ResponseTypeOk,
		},
	}, nil
}

func (adapter *DefaultAdminAPIAdapter) AdaptAdminUsersCreate(ctx context.Context, request generatedhttpinfra.AdminUsersCreateRequest) (*generatedhttpinfra.AdminUsersCreateResponseWrapper, *httpinfra.HTTPError) {
	var adminStandardID entities.StandardID
	if request.AdminUserCreation.Meta != nil && request.AdminUserCreation.Meta.Id != nil {
		adminStandardID.ID = *request.AdminUserCreation.Meta.Id
	}
	var description *string
	if request.AdminUserCreation.Spec != nil && request.AdminUserCreation.Spec.Description != nil {
		description = request.AdminUserCreation.Spec.Description
	}
	input := admin.CreateAdminInput{
		StandardID:  adminStandardID,
		Description: description,
	}
	out, err := adapter.adminUseCase.CreateAdmin(ctx, input)
	if err != nil {
		return nil, httpinfra.NewHTTPErrorFromUseCaseError(ctx, err)
	}

	creationDate := out.CreationDate.String()
	lastUpdate := out.LastUpdate.String()
	response := generatedhttpinfra.AdminUsersCreateResponseWrapper{
		AdminUserDetail: generatedhttpinfra.AdminUserDetail{
			Meta: &generatedhttpinfra.ResourceMetaDetail{
				Id:              &out.ID,
				ResourceVersion: &out.ResourceVersion,
				CreationDate:    &creationDate,
				LastUpdate:      &lastUpdate,
			},
			Spec: &generatedhttpinfra.AdminUserDetailSpec{
				Roles:       &out.Roles,
				Description: out.Description,
			},
		},
		ResponseInfo: httpinfra.ResponseInfo{
			ResponseType: httpinfra.ResponseTypeCreated,
		},
	}
	return &response, nil
}

func (adapter *DefaultAdminAPIAdapter) AdaptAdminUsersDescribe(ctx context.Context, request generatedhttpinfra.AdminUsersDescribeRequest) (*generatedhttpinfra.AdminUsersDescribeResponseWrapper, *httpinfra.HTTPError) {
	var adminStandardID entities.StandardID
	adminStandardID.ID = request.AdminUserId
	input := admin.GetAdminInput{
		StandardID: adminStandardID,
	}
	out, err := adapter.adminUseCase.GetAdmin(ctx, input)
	if err != nil {
		return nil, httpinfra.NewHTTPErrorFromUseCaseError(ctx, err)
	}

	creationDate := out.CreationDate.String()
	lastUpdate := out.LastUpdate.String()
	response := generatedhttpinfra.AdminUsersDescribeResponseWrapper{
		AdminUserDetail: generatedhttpinfra.AdminUserDetail{
			Meta: &generatedhttpinfra.ResourceMetaDetail{
				Id:              &out.ID,
				ResourceVersion: &out.ResourceVersion,
				CreationDate:    &creationDate,
				LastUpdate:      &lastUpdate,
			},
			Spec: &generatedhttpinfra.AdminUserDetailSpec{
				Roles:       &out.Roles,
				Description: out.Description,
			},
		},
		ResponseInfo: httpinfra.ResponseInfo{
			ResponseType: httpinfra.ResponseTypeOk,
		},
	}
	return &response, nil
}

func (adapter *DefaultAdminAPIAdapter) AdaptAdminUsersEdit(ctx context.Context, request generatedhttpinfra.AdminUsersEditRequest) (*generatedhttpinfra.AdminUsersEditResponseWrapper, *httpinfra.HTTPError) {
	var description *string
	if request.AdminUserUpdate.Spec != nil && request.AdminUserUpdate.Spec.Description != nil {
		description = request.AdminUserUpdate.Spec.Description
	}

	input := admin.EditAdminInput{
		AdminEditable: admin.AdminEditable{
			StandardResourceMeta: entities.StandardResourceMeta{
				StandardResource: entities.StandardResource{
					StandardID: entities.StandardID{
						ID: request.AdminUserId,
					},
				},
				ResourceVersion: *request.AdminUserUpdate.Meta.ResourceVersion,
			},
			Description: description,
		},
	}
	out, err := adapter.adminUseCase.EditAdmin(ctx, input)
	if err != nil {
		return nil, httpinfra.NewHTTPErrorFromUseCaseError(ctx, err)
	}

	creationDate := out.CreationDate.String()
	lastUpdate := out.LastUpdate.String()
	response := generatedhttpinfra.AdminUsersEditResponseWrapper{
		AdminUserDetail: generatedhttpinfra.AdminUserDetail{
			Meta: &generatedhttpinfra.ResourceMetaDetail{
				Id:              &out.ID,
				ResourceVersion: &out.ResourceVersion,
				CreationDate:    &creationDate,
				LastUpdate:      &lastUpdate,
			},
			Spec: &generatedhttpinfra.AdminUserDetailSpec{
				Roles:       &out.Roles,
				Description: out.Description,
			},
		},
		ResponseInfo: httpinfra.ResponseInfo{
			ResponseType: httpinfra.ResponseTypeOk,
		},
	}
	return &response, nil
}

func (adapter *DefaultAdminAPIAdapter) AdaptAdminUsersList(ctx context.Context, request generatedhttpinfra.AdminUsersListRequest) (*generatedhttpinfra.AdminUsersListResponseWrapper, *httpinfra.HTTPError) {
	var input admin.ListAdminsInput

	var limitInput int
	if request.Limit != nil {
		limitInput = int(*request.Limit)
	}
	var offsetInput int
	if request.Offset != nil {
		offsetInput = int(*request.Offset)
	}
	pageLimit := utils.MaxValue(utils.DefaultIntValue(limitInput, defaultAdminUserListLimit), maxListAdminUserLimit)
	input.PageLimit = pageLimit
	input.PageOffset = offsetInput
	input.OrderBy = request.OrderBy
	input.OrderDirection = request.OrderDirection

	outputData, err := adapter.adminUseCase.ListAdmins(ctx, input)
	if err != nil {
		return nil, httpinfra.NewHTTPErrorFromUseCaseError(ctx, err)
	}

	adaptedUseCaseCollection := make([]generatedhttpinfra.AdminUserDetail, len(outputData.Items))
	for i, outItem := range outputData.Items {
		item := outItem
		creationDate := item.CreationDate.String()
		lastUpdate := item.LastUpdate.String()
		meta := generatedhttpinfra.ResourceMetaDetail{
			Id:              &item.ID,
			ResourceVersion: &item.ResourceVersion,
			CreationDate:    &creationDate,
			LastUpdate:      &lastUpdate,
		}
		spec := generatedhttpinfra.AdminUserDetailSpec{
			Roles:       &item.Roles,
			Description: item.Description,
		}
		detail := generatedhttpinfra.AdminUserDetail{
			Meta: &meta,
			Spec: &spec,
		}
		if item.Description != nil {
			detail.Spec.Description = item.Description
		}

		adaptedUseCaseCollection[i] = detail
	}
	offset := int32(outputData.Offset)
	limit := int32(outputData.Limit)
	response := generatedhttpinfra.AdminUsersListResponseWrapper{
		AdminUserCollection: generatedhttpinfra.AdminUserCollection{
			Items:     &adaptedUseCaseCollection,
			Offset:    &offset,
			Limit:     &limit,
			MoreItems: &outputData.MoreItems,
		},
		ResponseInfo: httpinfra.ResponseInfo{
			ResponseType: httpinfra.ResponseTypeOk,
		},
	}
	return &response, nil
}

func (adapter *DefaultAdminAPIAdapter) AdaptAdminUsersRemove(ctx context.Context, request generatedhttpinfra.AdminUsersRemoveRequest) (*generatedhttpinfra.AdminUsersRemoveResponseWrapper, *httpinfra.HTTPError) {
	input := admin.DeleteAdminInput{
		StandardID: entities.StandardID{
			ID: request.AdminUserId,
		},
	}
	out, err := adapter.adminUseCase.DeleteAdmin(ctx, input)
	if err != nil {
		return nil, httpinfra.NewHTTPErrorFromUseCaseError(ctx, err)
	}

	creationDate := out.CreationDate.String()
	lastUpdate := out.LastUpdate.String()
	response := generatedhttpinfra.AdminUsersRemoveResponseWrapper{
		AdminUserDetail: generatedhttpinfra.AdminUserDetail{
			Meta: &generatedhttpinfra.ResourceMetaDetail{
				Id:              &out.ID,
				ResourceVersion: &out.ResourceVersion,
				CreationDate:    &creationDate,
				LastUpdate:      &lastUpdate,
			},
			Spec: &generatedhttpinfra.AdminUserDetailSpec{
				Roles:       &out.Roles,
				Description: out.Description,
			},
		},
		ResponseInfo: httpinfra.ResponseInfo{
			ResponseType: httpinfra.ResponseTypeOk,
		},
	}
	return &response, nil
}

func mapSlot(slot hsmslot.HSMSlot) generatedhttpinfra.SlotDetail {
	creationDate := slot.CreationDate.String()
	lastUpdate := slot.LastUpdate.String()

	return generatedhttpinfra.SlotDetail{
		Meta: &generatedhttpinfra.ResourceMetaDetail{
			Id:              &slot.ID,
			ResourceVersion: &slot.ResourceVersion,
			CreationDate:    &creationDate,
			LastUpdate:      &lastUpdate,
		},
		Spec: &generatedhttpinfra.SlotDetailSpec{
			HardwareSecurityModuleId: &slot.HSMModuleID,
			ApplicationId:            &slot.ApplicationID,
			Slot:                     &slot.Slot,
		},
	}
}

func mapModule(module hsmmodule.HSMModule) (*generatedhttpinfra.ModuleDetail, error) {
	creationDate := module.CreationDate.String()
	lastUpdate := module.LastUpdate.String()
	var optionalDescription httpinfra.Optional[string]
	if module.Description != nil {
		optionalDescription.SetValue(module.Description)
	}

	infraHSMKind, err := mapModuleCreationSpecConfigurationTypeFrom(module.Kind)
	if err != nil {
		return nil, err
	}
	kind := string(*infraHSMKind)
	return &generatedhttpinfra.ModuleDetail{
		Meta: &generatedhttpinfra.ResourceMetaDetail{
			Id:              &module.ID,
			ResourceVersion: &module.ResourceVersion,
			CreationDate:    &creationDate,
			LastUpdate:      &lastUpdate,
		},
		Spec: &generatedhttpinfra.ModuleSpec{
			Configuration: &generatedhttpinfra.ModuleSpecConfiguration{
				HsmKind: *infraHSMKind,
				SoftHsm: &generatedhttpinfra.SoftHsm{
					HsmKind: &kind,
				},
			},
			Description: module.Description,
		},
	}, nil
}

func mapModuleCreationSpecConfigurationTypeFrom(configurationType hsmmodule.ModuleKind) (*generatedhttpinfra.ModuleSpecConfigurationHsmKind, *httpinfra.HTTPError) {
	if configurationType == hsmmodule.SoftHSMModuleKind {
		t := generatedhttpinfra.HsmKindSofthsm
		return &t, nil
	}
	return nil, httpinfra.NewHTTPError(httpinfra.StatusInternal).SetMessage(fmt.Sprintf("cannot map invalid module kind [%s]", configurationType))
}

func mapUseCaseHSMKindFrom(configurationKind generatedhttpinfra.ModuleSpecConfigurationHsmKind) (*hsmmodule.ModuleKind, *httpinfra.HTTPError) {
	if configurationKind == generatedhttpinfra.HsmKindSofthsm {
		t := hsmmodule.SoftHSMModuleKind
		return &t, nil
	}
	return nil, httpinfra.NewHTTPError(httpinfra.StatusInternal).SetMessage(fmt.Sprintf("can't map '%s' to usecase HSM type", configurationKind))
}

// DefaultAdminAPIAdapter implements AdminAPIAdapter.
type DefaultAdminAPIAdapter struct {
	applicationUseCase application.ApplicationUseCase
	adminUseCase       admin.AdminUseCase
	hsmUseCase         hsmmodule.HSMModuleUseCase
	hsmSlotUseCase     hsmslot.HSMSlotUseCase
}

// DefaultAdminAPIAdapterOptions options to create a new DefaultAdminAPIAdapter.
type DefaultAdminAPIAdapterOptions struct {
	ApplicationUseCase application.ApplicationUseCase
	AdminUseCase       admin.AdminUseCase
	HSMUseCase         hsmmodule.HSMModuleUseCase
	HSMSlotUseCase     hsmslot.HSMSlotUseCase
}

// ProvideDefaultAdminAPIAdapter creates a new DefaultAdminAPIAdapter instance.
func ProvideDefaultAdminAPIAdapter(options DefaultAdminAPIAdapterOptions) (*DefaultAdminAPIAdapter, error) {
	if options.ApplicationUseCase == nil {
		return nil, errors.New("mandatory 'ApplicationUseCase' was not provided")
	}
	if options.AdminUseCase == nil {
		return nil, errors.New("mandatory 'AdminUseCase' was not provided")
	}
	if options.HSMSlotUseCase == nil {
		return nil, errors.New("mandatory 'HSMSlotUseCase' was not provided")
	}
	if options.HSMUseCase == nil {
		return nil, errors.New("mandatory 'HSMUseCase' was not provided")
	}

	return &DefaultAdminAPIAdapter{
		applicationUseCase: options.ApplicationUseCase,
		adminUseCase:       options.AdminUseCase,
		hsmUseCase:         options.HSMUseCase,
		hsmSlotUseCase:     options.HSMSlotUseCase,
	}, nil
}
