package pdp

import (
	"context"
	"reflect"
)

// ActionsPolicyInformationPointPort is a port to adapt requests related to the Actions.
type ActionsPolicyInformationPointPort interface {
	// ListActions list the Actions assigned to the given roles.
	ListActions(ctx context.Context, input ListActionsInput) (*ListActionsOutput, error)
}

// ListActionsInput are the attributes needed to list actions related to a set of roles.
type ListActionsInput struct {
	// Roles is the set of roles to get actions for. It can be one or more roles.
	Roles []string
}

// ListActionsOutput is the result of listing actions related to a set of roles.
type ListActionsOutput struct {
	// Actions the actions authorized for the roles.
	Actions Actions
}

// Actions is the result of listing actions related to a set of roles.
type Actions struct {
	// actionSet hashset of actions. It is a hashset as it is more efficient for lookups than a slice.
	actionSet map[string]string
}

// Merge actions with the ones provided as an argument.
func (a *Actions) Merge(actions Actions) {
	for action := range actions.actionSet {
		a.actionSet[action] = ""
	}
}

// Equal checks if the object has the same actions as the one provided as input argument.
func (a *Actions) Equal(actions Actions) bool {
	return reflect.DeepEqual(a.actionSet, actions.actionSet)
}

// NewActions creates a new Actions pointer.
func NewActions(actions []string) *Actions {
	actionSet := make(map[string]string)
	for _, action := range actions {
		actionSet[action] = ""
	}
	return &Actions{
		actionSet: actionSet,
	}
}

// Contains checks whether the provided action is present.
func (a Actions) Contains(actionID string) bool {
	_, ok := a.actionSet[actionID]
	return ok
}
