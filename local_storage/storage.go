package local_storage

import (
	"slices"
	"time"
	"uuid"

	"github.com/dontkillmyvibee/library.git/local_storage/models"
)

type Storage struct {
	books map[uuid.UUID]models.Book
}

func NewStorage() *Storage {
	return &Storage{
		books: make(map[uuid.UUID]models.Book),
	}
}

func (s *Storage) AddBook(book models.Book) error {
	if _, ok := s.books[book.ID]; ok {
		return ErrBookAlreadyExist
	}

	for _, existingBook := range s.books {
		if existingBook.Title == book.Title &&
			slices.Equal(existingBook.AuthorNames, book.AuthorNames) {
			return ErrBookAlreadyExist
		}
	}

	s.books[book.ID] = book
	return nil
}

func (s *Storage) GetBook(id uuid.UUID) (models.Book, error) {
	book, ok := s.books[id]
	if !ok {
		return models.Book{}, ErrBookNotFound
	}

	if book.DeletedAt != nil {
		return models.Book{}, ErrBookNotFound
	}

	return book, nil
}

func (s *Storage) GetAllBooks() map[uuid.UUID]models.Book {
	tmp := make(map[uuid.UUID]models.Book, len(s.books))

	for id, book := range s.books {
		if book.DeletedAt == nil {
			tmp[id] = book
		}
	}

	return tmp
}

func (s *Storage) GetAllAvailableBooks() map[uuid.UUID]models.Book {
	tmp := make(map[uuid.UUID]models.Book)

	for id, book := range s.books {
		if book.DeletedAt == nil && book.IsAvailable {
			tmp[id] = book
		}
	}

	return tmp
}

func (s *Storage) GetAllUnavailableBooks() map[uuid.UUID]models.Book {
	tmp := make(map[uuid.UUID]models.Book)

	for id, book := range s.books {
		if book.DeletedAt == nil && !book.IsAvailable {
			tmp[id] = book
		}
	}

	return tmp
}

func (s *Storage) DeleteBook(id uuid.UUID) error {
	book, ok := s.books[id]
	if !ok {
		return ErrBookNotFound
	}

	if book.DeletedAt != nil {
		return ErrBookAlreadyDeleted
	}

	now := time.Now()
	book.DeletedAt = &now
	book.UpdatedAt = now

	s.books[id] = book

	return nil
}
