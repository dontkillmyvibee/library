package models

import (
	"time"
	"uuid"
)

type Reader struct {
	ID         uuid.UUID
	FirstName  string
	LastName   string
	MiddleName string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}

func NewReader(firstName, lastName, middleName string) Reader {
	now := time.Now()
	return Reader{
		ID:         uuid.NewV4(),
		FirstName:  firstName,
		LastName:   lastName,
		MiddleName: middleName,
		CreatedAt:  now,
		UpdatedAt:  now,
		DeletedAt:  nil,
	}
}

type UpdateReaderData struct {
	FirstName  string
	LastName   string
	MiddleName string
}
