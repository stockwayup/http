package dictionary

const (
	BadRequestTitle   = "Bad request"
	InvalidAttribute  = "Invalid Attribute"
	UnauthorizedTitle = "Unauthorized"
	NotFoundTitle     = "Not found"
	ForbiddenTitle    = "Forbidden"

	InvalidJSONBodyCode = "1"
	InvalidJSONBodyDesc = "Request body contains invalid json"

	InvalidAccessTokenCode = "401"
	InvalidAccessTokenDesc = "Access token or refresh token missing or invalid"

	NotFoundCode = "404"
	NotFoundDesc = "Resource not found"

	ForbiddenCode = "403"
	ForbiddenDesc = "Access denied"

	TimeoutCode  = "408"
	TimeoutTitle = "Request timeout"
)
