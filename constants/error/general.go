package error

import "errors"

var (
	ErrInternalServerError = errors.New("internal server error")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrInvalidUploadFile   = errors.New("invalid upload file")
	ErrSizeTooBig          = errors.New("size too big")
	ErrForbidden           = errors.New("forbidden")
)

var GeneralErrors = []error{
	ErrInternalServerError,
	ErrUnauthorized,
	ErrForbidden,
}
