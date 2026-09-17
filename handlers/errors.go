package handlers

import "errors"

var (
	ErrInvalidUUID           = errors.New("invalid UUID")
	ErrInvalidQueryParameter = errors.New("invalid query parameter")
)
