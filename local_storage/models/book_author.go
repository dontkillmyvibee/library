package models

import "uuid"

type BookAuthor struct {
	BookID   uuid.UUID
	AuthorID uuid.UUID
}

func NewBookAuthor(bookID, authorID uuid.UUID) BookAuthor {
	return BookAuthor{
		BookID:   bookID,
		AuthorID: authorID,
	}
}
