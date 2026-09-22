package schemas

import (
	"strings"
	"time"
	"unicode/utf8"
	"uuid"
)

type CreateAuthorRequestSchema struct {
	FirstName  *string `json:"first_name"`
	LastName   *string `json:"last_name"`
	MiddleName *string `json:"middle_name"`
}

type UpdateAuthorRequestSchema struct {
	FirstName  *string `json:"first_name"`
	LastName   *string `json:"last_name"`
	MiddleName *string `json:"middle_name"`
}

type CreateAuthorResponseSchema struct {
	ID         uuid.UUID `json:"id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	MiddleName string    `json:"middle_name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type GetAuthorResponseSchema CreateAuthorResponseSchema

type UpdateAuthorResponseSchema CreateAuthorResponseSchema

type GetAllAuthorsResponseSchema []GetAuthorResponseSchema

func (req *CreateAuthorRequestSchema) Validate() error {
	if req.FirstName == nil {
		return ErrMissingRequiredFieldsFirstName
	}

	if req.LastName == nil {
		return ErrMissingRequiredFieldsLastName
	}

	if req.MiddleName == nil {
		return ErrMissingRequiredFieldsMiddleName
	}

	if n := utf8.RuneCountInString(*req.FirstName); n < 1 || n > 100 {
		return ErrValidationFailedFirstName
	}

	if n := utf8.RuneCountInString(*req.LastName); n < 1 || n > 100 {
		return ErrValidationFailedLastName
	}

	if n := utf8.RuneCountInString(*req.MiddleName); n < 1 || n > 100 {
		return ErrValidationFailedMiddleName
	}

	return nil
}

func (req *CreateAuthorRequestSchema) Normalize() {
	if req.FirstName != nil {
		*req.FirstName = strings.TrimSpace(*req.FirstName)
	}

	if req.LastName != nil {
		*req.LastName = strings.TrimSpace(*req.LastName)
	}

	if req.MiddleName != nil {
		*req.MiddleName = strings.TrimSpace(*req.MiddleName)
	}
}

func (req *UpdateAuthorRequestSchema) Validate() error {
	if req.FirstName == nil {
		return ErrMissingRequiredFieldsFirstName
	}

	if req.LastName == nil {
		return ErrMissingRequiredFieldsLastName
	}

	if req.MiddleName == nil {
		return ErrMissingRequiredFieldsMiddleName
	}

	if n := utf8.RuneCountInString(*req.FirstName); n < 1 || n > 100 {
		return ErrValidationFailedFirstName
	}

	if n := utf8.RuneCountInString(*req.LastName); n < 1 || n > 100 {
		return ErrValidationFailedLastName
	}

	if n := utf8.RuneCountInString(*req.MiddleName); n < 1 || n > 100 {
		return ErrValidationFailedMiddleName
	}

	return nil
}

func (req *UpdateAuthorRequestSchema) Normalize() {
	if req.FirstName != nil {
		*req.FirstName = strings.TrimSpace(*req.FirstName)
	}

	if req.LastName != nil {
		*req.LastName = strings.TrimSpace(*req.LastName)
	}

	if req.MiddleName != nil {
		*req.MiddleName = strings.TrimSpace(*req.MiddleName)
	}
}
