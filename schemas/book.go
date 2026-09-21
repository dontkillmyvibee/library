package schemas

import (
	"strings"
	"time"
	"unicode/utf8"
	"uuid"
)

type Authors struct {
	ID         uuid.UUID `json:"id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	MiddleName string    `json:"middle_name"`
}

type CreateBookRequestSchema struct {
	Title       *string      `json:"title"`
	Description *string      `json:"description"`
	AuthorIDs   *[]uuid.UUID `json:"author_ids"`
}

type GetBookResponseSchema struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Authors     []Authors `json:"authors"`
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

	if req.AuthorIDs == nil {
		return ErrMissingRequiredFieldsAuthorIDs
	}

	if len(*req.AuthorIDs) == 0 {
		return ErrValidationFailedAuthorIDs
	}

	for _, authorID := range *req.AuthorIDs {
		if authorID == uuid.Nil() {
			return ErrValidationFailedAuthorIDs
		}
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
}
