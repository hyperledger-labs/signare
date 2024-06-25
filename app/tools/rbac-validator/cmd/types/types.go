package types

type Role struct {
	ID          string   `yaml:"id"`
	Description string   `yaml:"description"`
	Permissions []string `yaml:"permissions"`
}

type RoleCollection struct {
	Roles []Role `yaml:"roles"`
}

type Permission struct {
	ID          string   `yaml:"id"`
	Description string   `yaml:"description"`
	Actions     []string `yaml:"actions"`
}

type PermissionCollection struct {
	Permissions []Permission `yaml:"permissions"`
}

type ActionCollection struct {
	Actions []string `yaml:"actions"`
}

// MergeActions concatenates returns an ActionCollection after concatenating actions from the two inputs
func MergeActions(actionsOne, actionsTwo ActionCollection, actionsToExclude []string) ActionCollection {
	merged := append(actionsOne.Actions, actionsTwo.Actions...)
	var result []string
	for _, item := range merged {
		if !sliceContainsString(item, actionsToExclude) {
			result = append(result, item)
		}
	}

	return ActionCollection{
		Actions: result,
	}
}

func sliceContainsString(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
