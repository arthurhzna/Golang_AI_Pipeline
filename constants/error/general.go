package error

import "errors"

var (
	ErrInternalServerError = errors.New("internal server error")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrForbidden           = errors.New("forbidden")
	ErrInvalidAPIKey       = errors.New("invalid or missing API key")
)

var GeneralErrors = []error{
	ErrInternalServerError,
	ErrUnauthorized,
	ErrForbidden,
}
