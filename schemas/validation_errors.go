package schemas

import "errors"

var (
	ErrMissingRequiredFieldsTitle       = errors.New("missing required fields: title")
	ErrMissingRequiredFieldsDescription = errors.New("missing required fields: description")
	ErrMissingRequiredFieldsAuthorNames = errors.New("missing required fields: author_names")
	ErrValidationFailedTitle            = errors.New("validation failed: title")
	ErrValidationFailedDescription      = errors.New("validation failed: description")
	ErrValidationFailedAuthorNames      = errors.New("validation failed: author_names")
	ErrMissingRequiredFieldsFirstName   = errors.New("missing required fields: first_name")
	ErrMissingRequiredFieldsLastName    = errors.New("missing required fields: last_name")
	ErrMissingRequiredFieldsMiddleName  = errors.New("missing required fields: middle_name")
	ErrValidationFailedFirstName        = errors.New("validation failed: first_name")
	ErrValidationFailedLastName         = errors.New("validation failed: last_name")
	ErrValidationFailedMiddleName       = errors.New("validation failed: middle_name")
)
