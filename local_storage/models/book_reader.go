package models

import "uuid"

type BookReader struct {
	BookID   uuid.UUID
	ReaderID uuid.UUID
}

func NewBookReader(bookID, readerID uuid.UUID) BookReader {
	return BookReader{
		BookID:   bookID,
		ReaderID: readerID,
	}
}
