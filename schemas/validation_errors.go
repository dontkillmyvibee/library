package schemas

import "errors"

var (
	ErrMissingRequiredFieldsTitle       = errors.New("missing required fields: title")
	ErrMissingRequiredFieldsDescription = errors.New("missing required fields: description")
	ErrMissingRequiredFieldsAuthorNames = errors.New("missing required fields: author_names")
	ErrValidationFailedTitle            = errors.New("validation failed: title")
	ErrValidationFailedDescription      = errors.New("validation failed: description")
	ErrValidationFailedAuthorNames      = errors.New("validation failed: author_names")
)
