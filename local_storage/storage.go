package local_storage

import (
	"time"
	"uuid"

	"github.com/dontkillmyvibee/library.git/local_storage/models"
)

type Storage struct {
	books       map[uuid.UUID]models.Book
	authors     map[uuid.UUID]models.Author
	bookAuthors map[models.BookAuthor]struct{}
}

func NewStorage() *Storage {
	return &Storage{
		books:       make(map[uuid.UUID]models.Book),
		authors:     make(map[uuid.UUID]models.Author),
		bookAuthors: make(map[models.BookAuthor]struct{}),
	}
}

func (s *Storage) AddBook(newBook models.Book, authorIDs []uuid.UUID) (models.Book, error) {
	uniqueAuthorIDs := make(map[uuid.UUID]struct{}, len(authorIDs))

	for _, authorID := range authorIDs {
		if _, exists := uniqueAuthorIDs[authorID]; exists {
			return models.Book{}, ErrDuplicateAuthor
		}
		uniqueAuthorIDs[authorID] = struct{}{}
	}

	for authorID := range uniqueAuthorIDs {
		author, ok := s.authors[authorID]
		if !ok || author.DeletedAt != nil {
			return models.Book{}, ErrAuthorNotFound
		}
	}

	if _, exists := s.books[newBook.ID]; exists {
		return models.Book{}, ErrBookAlreadyExists
	}

	for _, book := range s.books {
		if book.DeletedAt == nil && book.Title == newBook.Title {
			return models.Book{}, ErrBookAlreadyExists
		}
	}

	s.books[newBook.ID] = newBook

	for authorID := range uniqueAuthorIDs {
		bookModel := models.NewBookAuthor(newBook.ID, authorID)
		if err := s.AddBookAuthor(bookModel); err != nil {
			delete(s.books, newBook.ID)
			s.DeleteBookAuthorByBookId(newBook.ID)
			return models.Book{}, err
		}
	}

	if err := s.Save(); err != nil {
		delete(s.books, newBook.ID)
		s.DeleteBookAuthorByBookId(newBook.ID)

		return models.Book{}, err
	}

	return newBook, nil
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

	oldBook := book
	oldBookAuthors := s.GetBookAuthorByBookID(id)

	now := time.Now()
	book.DeletedAt = &now
	book.UpdatedAt = now

	s.books[id] = book
	s.DeleteBookAuthorByBookId(id)

	if err := s.Save(); err != nil {
		s.books[id] = oldBook

		for bookAuthor := range oldBookAuthors {
			s.bookAuthors[bookAuthor] = struct{}{}
		}

		return err
	}

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

	if err := s.Save(); err != nil {
		delete(s.authors, author.ID)
		return err
	}

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

func (s *Storage) GetAuthorsByBookID(bookID uuid.UUID) []models.Author {
	bookAuthors := s.GetBookAuthorByBookID(bookID)

	authors := make([]models.Author, 0, len(bookAuthors))

	for bookAuthor := range bookAuthors {
		author, ok := s.authors[bookAuthor.AuthorID]
		if !ok || author.DeletedAt != nil {
			continue
		}

		authors = append(authors, author)
	}

	return authors
}

func (s *Storage) DeleteAuthor(id uuid.UUID) error {
	author, ok := s.authors[id]
	if !ok {
		return ErrAuthorNotFound
	}

	if author.DeletedAt != nil {
		return ErrAuthorAlreadyDeleted
	}

	bookAuthors := s.GetBookAuthorByAuthorId(id)

	for bookAuthor := range bookAuthors {
		authors := s.GetAuthorsByBookID(bookAuthor.BookID)

		if len(authors) == 1 {
			return ErrLastBookAuthor
		}
	}

	oldAuthor := author

	now := time.Now()
	author.DeletedAt = &now
	author.UpdatedAt = now

	s.authors[id] = author
	s.DeleteBookAuthorByAuthorId(id)

	if err := s.Save(); err != nil {
		s.authors[id] = oldAuthor

		for bookAuthor := range bookAuthors {
			s.bookAuthors[bookAuthor] = struct{}{}
		}

		return err
	}

	return nil
}

func (s *Storage) AddBookAuthor(bookAuthor models.BookAuthor) error {
	if _, ok := s.bookAuthors[bookAuthor]; ok {
		return ErrBookAuthorAlreadyExist
	}

	s.bookAuthors[bookAuthor] = struct{}{}
	return nil
}

func (s *Storage) DeleteBookAuthor(bookAuthor models.BookAuthor) error {
	if _, ok := s.bookAuthors[bookAuthor]; !ok {
		return ErrBookAuthorNotFound
	}

	delete(s.bookAuthors, bookAuthor)
	return nil
}

func (s *Storage) DeleteBookAuthorByAuthorId(authorID uuid.UUID) {
	for bookAuthor := range s.bookAuthors {
		if bookAuthor.AuthorID == authorID {
			delete(s.bookAuthors, bookAuthor)
		}
	}
}

func (s *Storage) DeleteBookAuthorByBookId(bookID uuid.UUID) {
	for bookAuthor := range s.bookAuthors {
		if bookAuthor.BookID == bookID {
			delete(s.bookAuthors, bookAuthor)
		}
	}
}

func (s *Storage) GetBookAuthorByBookID(bookID uuid.UUID) map[models.BookAuthor]struct{} {
	tmp := make(map[models.BookAuthor]struct{})

	for bookAuthor := range s.bookAuthors {
		if bookAuthor.BookID == bookID {
			tmp[bookAuthor] = struct{}{}
		}
	}

	return tmp
}

func (s *Storage) GetBookAuthorByAuthorId(authorID uuid.UUID) map[models.BookAuthor]struct{} {
	tmp := make(map[models.BookAuthor]struct{})

	for bookAuthor := range s.bookAuthors {
		if bookAuthor.AuthorID == authorID {
			tmp[bookAuthor] = struct{}{}
		}
	}

	return tmp
}
