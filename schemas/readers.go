package schemas

import (
	"strings"
	"time"
	"unicode/utf8"
	"uuid"
)

type CreateReaderRequestSchema struct {
	FirstName  *string `json:"first_name"`
	LastName   *string `json:"last_name"`
	MiddleName *string `json:"middle_name"`
}

type UpdateReaderRequestSchema struct {
	FirstName  *string `json:"first_name"`
	LastName   *string `json:"last_name"`
	MiddleName *string `json:"middle_name"`
}

type CreateReaderResponseSchema struct {
	ID         uuid.UUID `json:"id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	MiddleName string    `json:"middle_name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type GetReaderResponseSchema CreateReaderResponseSchema

type UpdateReaderResponseSchema CreateReaderResponseSchema
type GetAllReadersResponseSchema []GetReaderResponseSchema

func (req *CreateReaderRequestSchema) Validate() error {
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

func (req *CreateReaderRequestSchema) Normalize() {
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

func (req *UpdateReaderRequestSchema) Validate() error {
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

func (req *UpdateReaderRequestSchema) Normalize() {
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
