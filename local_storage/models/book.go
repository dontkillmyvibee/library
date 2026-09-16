package models

import (
	"time"
	"uuid"
)

type Book struct {
	ID          uuid.UUID
	Title       string
	Description string
	AuthorNames []string
	IsAvailable bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

func NewBook(title, description string, authorNames []string) Book {
	now := time.Now()
	tmp := make([]string, 0, len(authorNames))
	tmp = append(tmp, authorNames...)
	return Book{
		ID:          uuid.NewV4(),
		Title:       title,
		Description: description,
		AuthorNames: tmp,
		IsAvailable: true,
		CreatedAt:   now,
		UpdatedAt:   now,
		DeletedAt:   nil,
	}
}
