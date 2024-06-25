// Package roleinfile defines the implementation of the output adapter to read roles from YAML files.
package roleinfile

import (
	"context"
	"io/fs"
	"path"

	"gopkg.in/yaml.v3"

	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/authorization/role"
)

const defaultRolesFileName = "roles.yaml"

var _ role.RoleStorage = new(DefaultRoleStorageInFile)

// ListRoles fetches a collection of Role based on the ListRolesInput
func (d DefaultRoleStorageInFile) ListRoles(_ context.Context, _ role.ListRolesInput) (*role.ListRolesOutput, error) {
	roles := make([]role.Role, 0)
	for _, r := range d.rolesInfo.Roles {
		newRole := role.Role{
			ID: r.ID,
		}
		roles = append(roles, newRole)
	}

	return &role.ListRolesOutput{
		Roles: roles,
	}, nil
}

// DefaultRoleStorageInFileOptions are the set of fields to create an DefaultRoleStorageInFile
type DefaultRoleStorageInFileOptions struct {
	FileSystem fs.FS
	BasePath   string
}

// DefaultRoleStorageInFile is a port to adapt requests related to the Roles
type DefaultRoleStorageInFile struct {
	rolesInfo RolesInfo
}

// ProvideDefaultRoleStorageInFile provides an instance of an DefaultRoleStorageInFile
func ProvideDefaultRoleStorageInFile(options DefaultRoleStorageInFileOptions) (*DefaultRoleStorageInFile, error) {
	if options.FileSystem == nil {
		return nil, errors.Internal().WithMessage("mandatory 'FileSystem' not provided")
	}
	// As roles, permissions and actions are static, we can read them at start up time
	rolesInfo, err := loadRoles(options.FileSystem, options.BasePath)
	if err != nil {
		return nil, err
	}

	return &DefaultRoleStorageInFile{
		rolesInfo: *rolesInfo,
	}, nil
}

func loadRoles(fileSystem fs.FS, basePath string) (*RolesInfo, error) {
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

	return &roles, nil
}
