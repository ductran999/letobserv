package consts

const (
	ProductionEnv = "prod"
	DeveloperEnv  = "dev"
	StagingEnv    = "staging"

	LogDebugLevel = "debug"
	LogInfoLevel  = "info"
)

const (
	// 5xx - Server errors.

	InternalServerError = "INTERNAL_SERVER_ERROR"

	// 4xx - Client errors.

	BadRequest       = "BAD_REQUEST"
	Unauthorized     = "UNAUTHORIZED"
	Forbidden        = "FORBIDDEN"
	NotFound         = "NOT_FOUND"
	MethodNotAllowed = "METHOD_NOT_ALLOWED"
	Conflict         = "CONFLICT"
	TooManyRequests  = "TOO_MANY_REQUESTS"
)

const (
	HealthyState = "heathy"
)
