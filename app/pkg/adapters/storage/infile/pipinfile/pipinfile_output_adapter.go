// Package pipinfile defines the implementations of the Policy Information Point output adapters to read information from files.
package pipinfile

import (
	"context"
	"io/fs"
	"path"

	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/authorization/pdp"

	"gopkg.in/yaml.v3"
)

const (
	defaultRolesFileName            = "roles.yaml"
	defaultPermissionsFileName      = "permissions.yaml"
	defaultActionsGeneratedFileName = "actions-generated.yaml"
	defaultActionsManualFileName    = "actions-manual.yaml"
)

var _ pdp.ActionsPolicyInformationPointPort = new(DefaultRBACActionsPolicyInformationPointYAMLOutputAdapter)

// ListActions list the Actions assigned to the given roles
func (d DefaultRBACActionsPolicyInformationPointYAMLOutputAdapter) ListActions(_ context.Context, input pdp.ListActionsInput) (*pdp.ListActionsOutput, error) {
	if len(input.Roles) < 1 {
		return nil, errors.InvalidArgument().WithMessage("no 'Roles' provided")
	}

	actions := pdp.NewActions([]string{})
	for _, role := range input.Roles {
		acts, ok := d.roleActionMap[role]
		if !ok {
			return nil, errors.InvalidArgument().WithMessage("role %s does not exist", role)
		}
		actions.Merge(acts)
	}

	return &pdp.ListActionsOutput{
		Actions: *actions,
	}, nil
}

func loadRolesAndActions(fileSystem fs.FS, basePath string) (*loadRolesResult, error) {
	roleActionMap := make(map[string]pdp.Actions)

	// Read manual actions
	var actionsManual ActionCollection
	{
		manualActionsFilePath := path.Join(basePath, defaultActionsManualFileName)
		actionsManualBytes, err := fs.ReadFile(fileSystem, manualActionsFilePath)
		if err != nil {
			return nil, errors.Internal().WithMessage("could not read manual actions file from %s", manualActionsFilePath)
		}
		err = yaml.Unmarshal(actionsManualBytes, &actionsManual)
		if err != nil {
			return nil, errors.Internal().WithMessage("manual actions from file %s could not be unmarshalled", manualActionsFilePath)
		}
	}

	// Read generated actions
	var actionsGenerated ActionCollection
	{
		generatedActionsFilePath := path.Join(basePath, defaultActionsGeneratedFileName)
		actionsGeneratedBytes, err := fs.ReadFile(fileSystem, generatedActionsFilePath)
		if err != nil {
			return nil, errors.Internal().WithMessage("could not read generated actions file from %s", generatedActionsFilePath)
		}
		err = yaml.Unmarshal(actionsGeneratedBytes, &actionsGenerated)
		if err != nil {
			return nil, errors.Internal().WithMessage("generated actions from file %s could not be unmarshalled", generatedActionsFilePath)
		}
	}

	// Merge manual and generated actions
	actions := mergeActions(actionsManual, actionsGenerated)

	// Read permissions and validate they point to existing actions
	var permissions map[string]ActionCollection
	{
		var permissionsCollection PermissionCollection
		permissionsFilePath := path.Join(basePath, defaultPermissionsFileName)
		permissionsBytes, err := fs.ReadFile(fileSystem, permissionsFilePath)
		if err != nil {
			return nil, errors.Internal().WithMessage("could not read permissions file from %s", permissionsFilePath)
		}
		err = yaml.Unmarshal(permissionsBytes, &permissionsCollection)
		if err != nil {
			return nil, errors.Internal().WithMessage("permissions from file %s could not be unmarshalled", permissionsFilePath)
		}
		err = validatePermissions(permissionsCollection, actions)
		if err != nil {
			return nil, err
		}

		permissions = make(map[string]ActionCollection)
		for _, permission := range permissionsCollection.Permissions {
			permissions[permission.ID] = ActionCollection{Actions: permission.Actions}
		}
	}

	// Read roles
	var roles RolesInfo
	{
		rolesFilePath := path.Join(basePath, defaultRolesFileName)
		rolesBytes, err := fs.ReadFile(fileSystem, rolesFilePath)
		if err != nil {
			return nil, errors.Internal().WithMessage("could not read roles file from %s", rolesFilePath)
		}
		err = yaml.Unmarshal(rolesBytes, &roles)
		if err != nil {
			return nil, errors.Internal().WithMessage("roles from file %s could not be unmarshalled", rolesFilePath)
		}
	}

	for _, role := range roles.Roles {
		permissionGrantedActions := ActionCollection{
			Actions: []string{},
		}
		for _, permission := range role.Permissions {
			actionsForPermission, ok := permissions[permission]
			if !ok {
				return nil, errors.Internal().WithMessage("role '%s' points to a permission '%s' that does not exist", role.ID, permission)
			}
			permissionGrantedActions = mergeActions(permissionGrantedActions, actionsForPermission)
		}
		roleActionMap[role.ID] = *pdp.NewActions(permissionGrantedActions.Actions)
	}

	return &loadRolesResult{
		roleActionMap: roleActionMap,
	}, nil
}

type loadRolesResult struct {
	roleActionMap map[string]pdp.Actions
}

// validatePermissions checks that actions pointed by permissions exist and returns an error if any of them doesn't exist
func validatePermissions(permissions PermissionCollection, actions ActionCollection) error {
	actionMap := make(map[string]string)
	for _, action := range actions.Actions {
		actionMap[action] = ""
	}
	for _, permission := range permissions.Permissions {
		for _, action := range permission.Actions {
			if _, actionExists := actionMap[action]; !actionExists {
				return errors.Internal().WithMessage("permission '%s' points to an action '%s' that does not exist", permission.ID, action)
			}
		}
	}
	return nil
}

// DefaultRBACActionsPolicyInformationPointYAMLOutputAdapterOptions are the set of fields to create an DefaultRBACActionsPolicyInformationPointYAMLOutputAdapter
type DefaultRBACActionsPolicyInformationPointYAMLOutputAdapterOptions struct {
	FileSystem fs.FS
	BasePath   string
}

// DefaultRBACActionsPolicyInformationPointYAMLOutputAdapter is a port to adapt requests related to the Actions
type DefaultRBACActionsPolicyInformationPointYAMLOutputAdapter struct {
	roleActionMap map[string]pdp.Actions
}

// ProvideDefaultRBACActionsPolicyInformationPointYAMLOutputAdapter provides an instance of an DefaultRBACActionsPolicyInformationPointYAMLOutputAdapter
func ProvideDefaultRBACActionsPolicyInformationPointYAMLOutputAdapter(options DefaultRBACActionsPolicyInformationPointYAMLOutputAdapterOptions) (*DefaultRBACActionsPolicyInformationPointYAMLOutputAdapter, error) {
	if options.FileSystem == nil {
		return nil, errors.Internal().WithMessage("mandatory 'FileSystem' not provided")
	}
	// As roles, permissions and actions are static, we can read them at start up time
	loadResult, err := loadRolesAndActions(options.FileSystem, options.BasePath)
	if err != nil {
		return nil, err
	}

	return &DefaultRBACActionsPolicyInformationPointYAMLOutputAdapter{
		roleActionMap: loadResult.roleActionMap,
	}, nil
}
