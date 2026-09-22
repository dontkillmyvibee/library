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
	readers     map[uuid.UUID]models.Reader
	bookReaders map[models.BookReader]struct{}
}

func NewStorage() *Storage {
	return &Storage{
		books:       make(map[uuid.UUID]models.Book),
		authors:     make(map[uuid.UUID]models.Author),
		bookAuthors: make(map[models.BookAuthor]struct{}),
		readers:     make(map[uuid.UUID]models.Reader),
		bookReaders: make(map[models.BookReader]struct{}),
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
			s.DeleteBookAuthorByBookID(newBook.ID)
			return models.Book{}, err
		}
	}

	if err := s.Save(); err != nil {
		delete(s.books, newBook.ID)
		s.DeleteBookAuthorByBookID(newBook.ID)

		return models.Book{}, err
	}

	return newBook, nil
}

func (s *Storage) UpdateBook(bookID uuid.UUID, data models.UpdateBookData) (models.Book, error) {
	book, ok := s.books[bookID]
	if !ok || book.DeletedAt != nil {
		return models.Book{}, ErrBookNotFound
	}

	uniqueAuthorIDs := make(map[uuid.UUID]struct{}, len(data.AuthorIDs))

	for _, authorID := range data.AuthorIDs {
		if _, exists := uniqueAuthorIDs[authorID]; exists {
			return models.Book{}, ErrDuplicateAuthor
		}

		author, ok := s.authors[authorID]
		if !ok || author.DeletedAt != nil {
			return models.Book{}, ErrAuthorNotFound
		}

		uniqueAuthorIDs[authorID] = struct{}{}
	}

	for id, existingBook := range s.books {
		if id == bookID || existingBook.DeletedAt != nil {
			continue
		}

		if existingBook.Title == data.Title {
			return models.Book{}, ErrBookAlreadyExists
		}
	}

	oldBook := book
	oldBookAuthors := s.GetBookAuthorByBookID(bookID)

	book.Title = data.Title
	book.Description = data.Description
	book.UpdatedAt = time.Now()

	s.books[bookID] = book

	for bookAuthor := range oldBookAuthors {
		delete(s.bookAuthors, bookAuthor)
	}

	for authorID := range uniqueAuthorIDs {
		s.bookAuthors[models.BookAuthor{
			BookID:   bookID,
			AuthorID: authorID,
		}] = struct{}{}
	}

	if err := s.Save(); err != nil {
		s.books[bookID] = oldBook

		for bookAuthor := range s.GetBookAuthorByBookID(bookID) {
			delete(s.bookAuthors, bookAuthor)
		}

		for bookAuthor := range oldBookAuthors {
			s.bookAuthors[bookAuthor] = struct{}{}
		}

		return models.Book{}, err
	}

	return book, nil
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
	s.DeleteBookAuthorByBookID(id)

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

func (s *Storage) UpdateAuthor(id uuid.UUID, data models.UpdateAuthorData) (models.Author, error) {
	author, ok := s.authors[id]
	if !ok || author.DeletedAt != nil {
		return models.Author{}, ErrAuthorNotFound
	}

	for existingID, existingAuthor := range s.authors {
		if existingID == id || existingAuthor.DeletedAt != nil {
			continue
		}

		if data.LastName == existingAuthor.LastName &&
			data.MiddleName == existingAuthor.MiddleName &&
			data.FirstName == existingAuthor.FirstName {
			return models.Author{}, ErrAuthorAlreadyExists
		}
	}

	oldAuthor := author

	author.FirstName = data.FirstName
	author.LastName = data.LastName
	author.MiddleName = data.MiddleName
	author.UpdatedAt = time.Now()

	s.authors[id] = author

	if err := s.Save(); err != nil {
		s.authors[id] = oldAuthor
		return models.Author{}, err
	}

	return author, nil
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

	bookAuthors := s.GetBookAuthorByAuthorID(id)

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
	s.DeleteBookAuthorByAuthorID(id)

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
		return ErrBookAuthorAlreadyExists
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

func (s *Storage) DeleteBookAuthorByAuthorID(authorID uuid.UUID) {
	for bookAuthor := range s.bookAuthors {
		if bookAuthor.AuthorID == authorID {
			delete(s.bookAuthors, bookAuthor)
		}
	}
}

func (s *Storage) DeleteBookAuthorByBookID(bookID uuid.UUID) {
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

func (s *Storage) GetBookAuthorByAuthorID(authorID uuid.UUID) map[models.BookAuthor]struct{} {
	tmp := make(map[models.BookAuthor]struct{})

	for bookAuthor := range s.bookAuthors {
		if bookAuthor.AuthorID == authorID {
			tmp[bookAuthor] = struct{}{}
		}
	}

	return tmp
}

func (s *Storage) AddReader(reader models.Reader) error {
	if _, ok := s.readers[reader.ID]; ok {
		return ErrReaderAlreadyExists
	}

	for _, existingReader := range s.readers {
		if existingReader.FirstName == reader.FirstName &&
			existingReader.LastName == reader.LastName &&
			existingReader.MiddleName == reader.MiddleName {
			return ErrReaderAlreadyExists
		}
	}

	s.readers[reader.ID] = reader

	if err := s.Save(); err != nil {
		delete(s.readers, reader.ID)
		return err
	}

	return nil
}

func (s *Storage) UpdateReader(id uuid.UUID, data models.UpdateReaderData) (models.Reader, error) {
	reader, ok := s.readers[id]
	if !ok || reader.DeletedAt != nil {
		return models.Reader{}, ErrReaderNotFound
	}

	for existingID, existingReader := range s.readers {
		if existingID == id || existingReader.DeletedAt != nil {
			continue
		}

		if existingReader.FirstName == data.FirstName &&
			existingReader.LastName == data.LastName &&
			existingReader.MiddleName == data.MiddleName {
			return models.Reader{}, ErrReaderAlreadyExists
		}
	}

	oldReader := reader

	reader.FirstName = data.FirstName
	reader.LastName = data.LastName
	reader.MiddleName = data.MiddleName
	reader.UpdatedAt = time.Now()

	s.readers[id] = reader

	if err := s.Save(); err != nil {
		s.readers[id] = oldReader
		return models.Reader{}, err
	}

	return reader, nil
}

func (s *Storage) DeleteReader(id uuid.UUID) error {
	reader, ok := s.readers[id]
	if !ok {
		return ErrReaderNotFound
	}

	if reader.DeletedAt != nil {
		return ErrReaderAlreadyDeleted
	}

	bookReaders := s.GetBookReaderByReaderID(id)

	if len(bookReaders) > 0 {
		return ErrReaderHasBook
	}

	oldReader := s.readers[id]

	now := time.Now()
	reader.DeletedAt = &now
	reader.UpdatedAt = now

	s.readers[id] = reader

	if err := s.Save(); err != nil {
		s.readers[id] = oldReader
		return err
	}

	return nil
}

func (s *Storage) GetReader(id uuid.UUID) (models.Reader, error) {
	reader, ok := s.readers[id]
	if !ok {
		return models.Reader{}, ErrReaderNotFound
	}

	if reader.DeletedAt != nil {
		return models.Reader{}, ErrReaderAlreadyDeleted
	}

	return reader, nil
}

func (s *Storage) GetAllReaders() map[uuid.UUID]models.Reader {
	tmp := make(map[uuid.UUID]models.Reader, len(s.readers))

	for _, reader := range s.readers {
		if reader.DeletedAt == nil {
			tmp[reader.ID] = reader
		}
	}

	return tmp
}

func (s *Storage) ReaderTakeBook(readerID, bookID uuid.UUID) error {
	reader, ok := s.readers[readerID]
	if !ok || reader.DeletedAt != nil {
		return ErrReaderNotFound
	}

	book, ok := s.books[bookID]
	if !ok || book.DeletedAt != nil {
		return ErrBookNotFound
	}

	if !book.IsAvailable {
		return ErrBookUnavailable
	}

	if bookReaders := s.GetBookReaderByReaderID(readerID); len(bookReaders) > 0 {
		return ErrReaderHasBook
	}

	bookReader := models.NewBookReader(bookID, readerID)

	if err := s.AddBookReader(bookReader); err != nil {
		return err
	}

	oldBook := book

	book.IsAvailable = false
	book.UpdatedAt = time.Now()
	s.books[bookID] = book

	if err := s.Save(); err != nil {
		delete(s.bookReaders, bookReader)
		s.books[bookID] = oldBook

		return err
	}

	return nil
}

func (s *Storage) ReaderReturnBook(readerID, bookID uuid.UUID) error {
	reader, ok := s.readers[readerID]
	if !ok || reader.DeletedAt != nil {
		return ErrReaderNotFound
	}

	bookReader := models.NewBookReader(bookID, readerID)

	if _, ok := s.bookReaders[bookReader]; !ok {
		return ErrBookReaderNotFound
	}

	book, ok := s.books[bookID]
	if !ok || book.DeletedAt != nil {
		return ErrBookNotFound
	}

	oldBook := book

	delete(s.bookReaders, bookReader)
	book.IsAvailable = true
	book.UpdatedAt = time.Now()

	s.books[bookID] = book

	if err := s.Save(); err != nil {
		s.bookReaders[bookReader] = struct{}{}
		s.books[bookID] = oldBook

		return err
	}

	return nil
}

func (s *Storage) AddBookReader(bookReader models.BookReader) error {
	if _, ok := s.bookReaders[bookReader]; ok {
		return ErrBookReaderAlreadyExists
	}

	s.bookReaders[bookReader] = struct{}{}
	return nil
}

func (s *Storage) GetBookReaderByReaderID(readerID uuid.UUID) map[models.BookReader]struct{} {
	tmp := make(map[models.BookReader]struct{})

	for bookReader := range s.bookReaders {
		if bookReader.ReaderID == readerID {
			tmp[bookReader] = struct{}{}
		}
	}

	return tmp
}
