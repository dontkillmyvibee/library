package models

import (
	"time"
	"uuid"
)

type Author struct {
	ID         uuid.UUID
	FirstName  string
	LastName   string
	MiddleName string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}

func NewAuthor(firstName, lastName, middleName string) Author {
	now := time.Now()
	return Author{
		ID:         uuid.NewV4(),
		FirstName:  firstName,
		LastName:   lastName,
		MiddleName: middleName,
		CreatedAt:  now,
		UpdatedAt:  now,
		DeletedAt:  nil,
	}
}
