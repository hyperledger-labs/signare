package hsmmodule

import "github.com/hyperledger-labs/signare/app/pkg/entities"

// HSMModule defines the HSMModule resource.
type HSMModule struct {
	// StandardResourceMeta defines the identifier of the HSMModule resource.
	entities.StandardResourceMeta
	// InternalResourceID uniquely identifies a HSMModule by a single ID.
	entities.InternalResourceID
	// Description of this CreateHSMModuleOutput.
	Description *string
	// Configuration defines the configuration of the CreateHSMModuleOutput.
	Configuration HSMModuleConfiguration
	// Kind defines the type of the CreateHSMModuleOutput.
	Kind ModuleKind
}

// HSMModuleConfiguration configuration of a module.
type HSMModuleConfiguration struct {
	// SoftHSMConfiguration configuration of a SoftHSM module.
	SoftHSMConfiguration *SoftHSMConfiguration
}

// SoftHSMConfiguration configuration of a SoftHSM module.
type SoftHSMConfiguration struct{}

// HSMModulesCollection defines a collection of HSMModule resources.
type HSMModulesCollection struct {
	// Items CreateHSMModuleOutput in collection.
	Items []HSMModule
	// StandardCollectionPage is the page data of the collection.
	entities.StandardCollectionPage
}

// CreateHSMModuleInput configures the creation of a HSMModule.
type CreateHSMModuleInput struct {
	// ID defines the identifier of the CreateHSMModuleOutput resource.
	ID *string
	// Description of this CreateHSMModuleOutput.
	Description *string
	// Configuration defines the configuration of the CreateHSMModuleOutput.
	Configuration HSMModuleConfiguration
	// ModuleKind defines the type of the CreateHSMModuleOutput.
	ModuleKind ModuleKind `valid:"in(SoftHSM)"`
}

// CreateHSMModuleOutput defines the output of the CreateHSMModule method.
type CreateHSMModuleOutput struct {
	// HSMModule defines the HSMModule resource.
	HSMModule
}

// ListHSMModulesInput defines all possible options to list CreateHSMModuleOutput resources.
type ListHSMModulesInput struct {
	// PageLimit maximum amount of CreateHSMModuleOutput in list output.
	PageLimit int `valid:"natural"`
	// PageOffset amount of CreateHSMModuleOutput elapsed in list output.
	PageOffset int `valid:"natural"`
	// OrderBy whether to order by last update date.
	OrderBy string
	// OrderDirection the direction in which the list will be ordered base on the attribute selected in OrderBy.
	OrderDirection string
}

// ListHSMModulesOutput defines the output of the ListHSMModules method.
type ListHSMModulesOutput struct {
	// HSMModulesCollection defines a collection of HSMModule resources.
	HSMModulesCollection
}

// GetHSMModuleInput configures getting a HSMModule.
type GetHSMModuleInput struct {
	// StandardID defines the identifier of the CreateHSMModuleOutput resource.
	entities.StandardID
}

// GetHSMModuleOutput defines the output of the GetHSMModule method.
type GetHSMModuleOutput struct {
	// HSMModule defines the HSMModule resource.
	HSMModule
}

// EditHSMModuleInput configures the update of a HSMModule.
type EditHSMModuleInput struct {
	// HSMModule defines the HSMModule resource.
	HSMModule
}

// EditHSMModuleOutput defines the output of the EditHSMModule method.
type EditHSMModuleOutput struct {
	// HSMModule defines the HSMModule resource.
	HSMModule
}

// DeleteHSMModuleInput configures the update of a HSMModule
type DeleteHSMModuleInput struct {
	// StandardID defines the identifier of the CreateHSMModuleOutput resource.
	entities.StandardID
}

// DeleteHSMModuleOutput defines the output of the DeleteHSMModule method.
type DeleteHSMModuleOutput struct {
	// HSMModule defines the HSMModule resource.
	HSMModule
}
