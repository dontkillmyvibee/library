package models

import (
	"time"
	"uuid"
)

type Book struct {
	ID          uuid.UUID
	Title       string
	Description string
	IsAvailable bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

func NewBook(title, description string) Book {
	now := time.Now()
	return Book{
		ID:          uuid.NewV4(),
		Title:       title,
		Description: description,
		IsAvailable: true,
		CreatedAt:   now,
		UpdatedAt:   now,
		DeletedAt:   nil,
	}
}
