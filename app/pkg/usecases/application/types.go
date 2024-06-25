package application

import (
	"github.com/hyperledger-labs/signare/app/pkg/commons/time"
	"github.com/hyperledger-labs/signare/app/pkg/entities"

	"github.com/google/uuid"
)

const (
	defaultOrderDirection = entities.OrderDesc
)

// Application defines an Application.
type Application struct {
	entities.StandardResourceMeta
	// InternalResourceID uniquely identifies an Application by a single ID.
	entities.InternalResourceID
	// ChainID defines the identifier of the chain configured in the Application.
	ChainID entities.Int256
	// Description contains the definition of the Application.
	Description *string
}

// CreateApplicationInput configures the creation of an Application.
type CreateApplicationInput struct {
	// ID defines the identifier of the Application resource.
	ID *string `valid:"optional"`
	// ChainID defines the identifier of the chain configured in the Application.
	ChainID entities.Int256 `valid:"required"`
	// Description contains the definition of the Application.
	Description *string `valid:"optional"`
}

// CreateApplicationOutput defines the output of the creation of an Application.
type CreateApplicationOutput struct {
	Application
}

// EditApplicationInput configures the edit of an Application.
type EditApplicationInput struct {
	// ID defines the identifier of the Application resource.
	ID string `valid:"required"`
	// ResourceVersion resource version for resource locking.
	ResourceVersion string `valid:"required"`
	// ChainID defines the identifier of the chain configured in the Application.
	ChainID *entities.Int256
	// Description contains the definition of the Application.
	Description *string
}

// EditApplicationOutput defines the output of editing of an Application.
type EditApplicationOutput struct {
	Application
}

// GetApplicationInput defines the input for getting an Application.
type GetApplicationInput struct {
	// StandardID defines the identifier of the resource.
	entities.StandardID
}

// GetApplicationOutput defines the output of getting an Application.
type GetApplicationOutput struct {
	Application
}

// DeleteApplicationInput configures the deletion of an Application.
type DeleteApplicationInput struct {
	// StandardID defines the identifier of the resource.
	entities.StandardID
}

// DeleteApplicationOutput defines the output of deleting an Application.
type DeleteApplicationOutput struct {
	Application
}

// ApplicationCollection defines a collection of Application.
type ApplicationCollection struct {
	// Items defines the Application resources in the collection.
	Items []Application
	// StandardCollectionPage is the page data of the collection.
	entities.StandardCollectionPage
}

// ListApplicationsInput defines all possible options to list Application resources.
type ListApplicationsInput struct {
	// PageLimit maximum amount of Application in list output.
	PageLimit int `valid:"natural"`
	// PageOffset amount of Application elapsed in list output.
	PageOffset int `valid:"natural"`
	// OrderBy whether to order by last update date.
	OrderBy string
	// OrderDirection the direction of the OrderBy.
	OrderDirection string
}

// ListApplicationsOutput defines the output of listing Applications.
type ListApplicationsOutput struct {
	ApplicationCollection
}

func createToApplication(input CreateApplicationInput) Application {
	now := time.Now()
	if input.ID == nil {
		randomID := uuid.NewString()
		input.ID = &randomID
	}
	return Application{
		StandardResourceMeta: entities.StandardResourceMeta{
			StandardResource: entities.StandardResource{
				StandardID: entities.StandardID{
					ID: *input.ID,
				},
				Timestamps: entities.Timestamps{
					CreationDate: now,
					LastUpdate:   now,
				},
			},
		},
		ChainID:     input.ChainID,
		Description: input.Description,
	}
}

func mapEditedValues(input EditApplicationInput) Application {
	application := Application{
		StandardResourceMeta: entities.StandardResourceMeta{
			StandardResource: entities.StandardResource{
				StandardID: entities.StandardID{
					ID: input.ID,
				},
				Timestamps: entities.Timestamps{
					LastUpdate: time.Now(),
				},
			},
			ResourceVersion: input.ResourceVersion,
		},
	}
	if input.ChainID != nil {
		application.ChainID = *input.ChainID
	}
	if input.Description != nil {
		application.Description = input.Description
	}
	return application

}
