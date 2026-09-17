package schemas

import (
	"slices"
	"strings"
	"time"
	"unicode/utf8"
	"uuid"
)

type CreateBookRequestSchema struct {
	Title       *string   `json:"title"`
	Description *string   `json:"description"`
	AuthorNames *[]string `json:"author_names"`
}

type GetBookResponseSchema struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	AuthorNames []string  `json:"author_names"`
	IsAvailable bool      `json:"is_available"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateBookResponseSchema GetBookResponseSchema

type GetAllBooksResponseSchema []GetBookResponseSchema

func (req *CreateBookRequestSchema) Validate() error {
	if req.Title == nil {
		return ErrMissingRequiredFieldsTitle
	}

	if n := utf8.RuneCountInString(*req.Title); n < 1 || n > 100 {
		return ErrValidationFailedTitle
	}

	if req.Description == nil {
		return ErrMissingRequiredFieldsDescription
	}

	if n := utf8.RuneCountInString(*req.Description); n < 1 || n > 300 {
		return ErrValidationFailedDescription
	}

	if req.AuthorNames == nil {
		return ErrMissingRequiredFieldsAuthorNames
	}

	if len(*req.AuthorNames) == 0 {
		return ErrValidationFailedAuthorNames
	}

	if slices.Contains(*req.AuthorNames, "") {
		return ErrValidationFailedAuthorNames
	}

	return nil
}

func (req *CreateBookRequestSchema) Normalize() {
	if req.Title != nil {
		*req.Title = strings.TrimSpace(*req.Title)
	}

	if req.Description != nil {
		*req.Description = strings.TrimSpace(*req.Description)
	}

	if req.AuthorNames != nil {
		for i := range *req.AuthorNames {
			(*req.AuthorNames)[i] = strings.TrimSpace((*req.AuthorNames)[i])
		}
	}
}
