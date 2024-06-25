package contextdefinition

const (
	DefaultUserHeader        = "X-Auth-UserId"
	DefaultApplicationHeader = "X-Auth-ApplicationId"
)

// AuthHeadersConfiguration configure the auth headers of the requests
type AuthHeadersConfiguration struct {
	// UserRequestHeader configure the auth header key of the user
	UserRequestHeader string
	// ApplicationRequestHeader configure the auth header key of the application
	ApplicationRequestHeader string
}
