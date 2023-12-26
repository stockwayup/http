package response

import "github.com/stockwayup/http/server/http/dictionary"

type Errors struct {
	Errors []Error `json:"errors"`
}

type Error struct {
	Code   string            `json:"code"`
	Title  string            `json:"title"`
	Detail string            `json:"detail,omitempty"`
	Meta   map[string]string `json:"meta,omitempty"`
	Source *Source           `json:"source,omitempty"`
}

type Source struct {
	Pointer string `json:"pointer,omitempty"`
}

//nolint:gochecknoglobals,stylecheck,golint
var UnauthorizedError = Error{
	Code:   dictionary.InvalidAccessTokenCode,
	Title:  dictionary.UnauthorizedTitle,
	Detail: dictionary.InvalidAccessTokenDesc,
}

//nolint:gochecknoglobals,stylecheck,golint
var NotFoundError = Error{
	Code:   dictionary.NotFoundCode,
	Title:  dictionary.NotFoundTitle,
	Detail: dictionary.NotFoundDesc,
}

//nolint:gochecknoglobals,stylecheck,golint
var ForbiddenError = Error{
	Code:   dictionary.ForbiddenCode,
	Title:  dictionary.ForbiddenTitle,
	Detail: dictionary.ForbiddenDesc,
}

//nolint:gochecknoglobals,stylecheck,golint
var TimeoutError = Error{
	Code:  dictionary.TimeoutCode,
	Title: dictionary.TimeoutTitle,
}
