package hsmslot

import "github.com/hyperledger-labs/signare/app/pkg/entities"

// HSMSlot defines the HSMSlot resource.
type HSMSlot struct {
	entities.StandardResourceMeta
	// InternalResourceID uniquely identifies an HSMSlot by a single ID.
	entities.InternalResourceID
	// ApplicationID defines the identifier of the Application of the HSMSlot.
	ApplicationID string `valid:"required"`
	// HSMModuleID defines the identifier of the module of the HSMSlot.
	HSMModuleID string `valid:"required"`
	// Slot defines the logical container on the HSM.
	Slot string `valid:"required"`
	// Pin defines the alphanumeric code used for authentication in the HSM.
	Pin string `valid:"required"`
}

// CreateHSMSlotInput configures the creation of an HSMSlot.
type CreateHSMSlotInput struct {
	// ID defines the identifier of the User resource.
	ID *string `valid:"optional"`
	// ApplicationID defines the identifier of the Application of the HSMSlot.
	ApplicationID string `valid:"required"`
	// HSMModuleID defines the identifier of the module of the HSMSlot.
	HSMModuleID string `valid:"required"`
	// Slot defines the logical container on the HSM.
	Slot string `valid:"required"`
	// Pin defines the alphanumeric code used for authentication in the HSM.
	Pin string `valid:"required"`
}

// CreateHSMSlotOutput defines the output of creating an HSMSlot.
type CreateHSMSlotOutput struct {
	HSMSlot
}

// GetHSMSlotInput defines the input for getting an HSMSlot.
type GetHSMSlotInput struct {
	entities.StandardID
}

// GetHSMSlotOutput defines the output of getting an HSMSlot.
type GetHSMSlotOutput struct {
	HSMSlot
}

// GetHSMSlotByApplicationInput defines the input for getting an application's HSMSlot.
type GetHSMSlotByApplicationInput struct {
	ApplicationID entities.StandardID `valid:"required"`
}

// GetHSMSlotByApplicationOutput defines the output of getting an application's HSMSlot.
type GetHSMSlotByApplicationOutput struct {
	HSMSlot
}

// EditPinInput configures the update of an HSMSlot's Pin.
type EditPinInput struct {
	entities.StandardID
	// ResourceVersion resource version for resource locking.
	ResourceVersion string `valid:"required"`
	// Pin defines the alphanumeric code used for authentication in the HSM.
	Pin string `valid:"required"`
	// HSMModuleID defines the alphanumeric code used for authentication in the HSM.
	HSMModuleID string `valid:"required"`
}

// EditPinOutput defines the output of editing an HSMSlot's Pin.
type EditPinOutput struct {
	HSMSlot
}

// DeleteHSMSlotInput configures the deletion of an HSMSlot.
type DeleteHSMSlotInput struct {
	entities.StandardID
}

// DeleteHSMSlotOutput defines the output of deleting an HSMSlot.
type DeleteHSMSlotOutput struct {
	HSMSlot
}

// ListHSMSlotsByApplicationInput defines all possible options to list HSMSlot resources for a specific Application.
type ListHSMSlotsByApplicationInput struct {
	// ApplicationID defines the identifier of the Application of the HSMSlot resource.
	ApplicationID entities.StandardID
	// PageLimit maximum amount of HSMSlot in list output.
	PageLimit int `valid:"natural"`
	// PageOffset amount of HSMSlot elapsed in list output.
	PageOffset int `valid:"natural"`
	// OrderBy whether to order by last update date.
	OrderBy string `valid:"optional"`
	// OrderDirection the direction in which the list will be ordered base on the attribute selected in OrderBy.
	OrderDirection string `valid:"optional"`
}

// ListHSMSlotsByApplicationOutput defines the output of listing HSMSlots for a specific Application.
type ListHSMSlotsByApplicationOutput struct {
	HSMSlotCollection
}

// ListHSMSlotsByHSMModuleInput defines all possible options to list HSMSlot resources for a specific HSMModuleID.
type ListHSMSlotsByHSMModuleInput struct {
	// HSMModuleID defines the identifier of the HSMModuleID of the HSMSlot resource.
	HSMModuleID entities.StandardID
	// PageLimit maximum amount of HSMSlot in list output.
	PageLimit int `valid:"natural"`
	// PageOffset amount of HSMSlot elapsed in list output.
	PageOffset int `valid:"natural"`
	// OrderBy whether to order by last update date.
	OrderBy string `valid:"optional"`
	// OrderDirection the direction in which the list will be ordered base on the attribute selected in OrderBy.
	OrderDirection string `valid:"optional"`
	// ApplicationID the application that will be used to filter the slots of the list
	ApplicationID *string `valid:"optional"`
}

// ListHSMSlotsByHSMModuleOutput defines the output of listing HSMSlots for a specific HSMModuleID.
type ListHSMSlotsByHSMModuleOutput struct {
	HSMSlotCollection
}

// HSMSlotCollection defines a collection of HSMSlot resources.
type HSMSlotCollection struct {
	// Items HSMSlot in collection.
	Items []HSMSlot
	// StandardCollectionPage is the page data of the collection.
	entities.StandardCollectionPage
}
