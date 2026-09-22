package errors

import "errors"

var (
	ErrInvalidUUID           = errors.New("invalid UUID")
	ErrInvalidQueryParameter = errors.New("invalid query parameter")
	ErrInternalServer        = errors.New("internal server error")
)
