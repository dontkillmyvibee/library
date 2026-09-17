package local_storage

import (
	"slices"
	"time"
	"uuid"

	"github.com/dontkillmyvibee/library.git/local_storage/models"
)

type Storage struct {
	books   map[uuid.UUID]models.Book
	authors map[uuid.UUID]models.Author
}

func NewStorage() *Storage {
	return &Storage{
		books:   make(map[uuid.UUID]models.Book),
		authors: make(map[uuid.UUID]models.Author),
	}
}

func (s *Storage) AddBook(book models.Book) error {
	if _, ok := s.books[book.ID]; ok {
		return ErrBookAlreadyExists
	}

	for _, existingBook := range s.books {
		if existingBook.Title == book.Title &&
			slices.Equal(existingBook.AuthorNames, book.AuthorNames) {
			return ErrBookAlreadyExists
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

func (s *Storage) AddAuthor(author models.Author) error {
	if _, ok := s.authors[author.ID]; ok {
		return ErrAuthorAlreadyExists
	}

	for _, existingAuthor := range s.authors {
		if existingAuthor.FirstName == author.FirstName &&
			existingAuthor.LastName == author.LastName &&
			existingAuthor.MiddleName == author.MiddleName {
			return ErrAuthorAlreadyExists
		}
	}

	s.authors[author.ID] = author

	return nil
}

func (s *Storage) GetAuthor(id uuid.UUID) (models.Author, error) {
	author, ok := s.authors[id]
	if !ok {
		return models.Author{}, ErrAuthorNotFound
	}

	if author.DeletedAt != nil {
		return models.Author{}, ErrAuthorAlreadyDeleted
	}

	return author, nil
}

func (s *Storage) GetAllAuthors() map[uuid.UUID]models.Author {
	tmp := make(map[uuid.UUID]models.Author, len(s.authors))

	for id, author := range s.authors {
		if author.DeletedAt == nil {
			tmp[id] = author
		}
	}

	return tmp
}

func (s *Storage) DeleteAuthor(id uuid.UUID) error {
	author, ok := s.authors[id]
	if !ok {
		return ErrAuthorNotFound
	}

	if author.DeletedAt != nil {
		return ErrAuthorAlreadyDeleted
	}

	now := time.Now()
	author.DeletedAt = &now
	author.UpdatedAt = now

	s.authors[id] = author

	return nil
}
