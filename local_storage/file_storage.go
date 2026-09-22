package local_storage

import (
	"encoding/json"
	"errors"
	"os"
	"uuid"

	"github.com/dontkillmyvibee/library.git/local_storage/models"
)

const fileStoragePath = "local_storage/file_storage.json"

type fileStorage struct {
	Books       map[uuid.UUID]models.Book   `json:"books"`
	Authors     map[uuid.UUID]models.Author `json:"authors"`
	Readers     map[uuid.UUID]models.Reader `json:"readers"`
	BookAuthors []models.BookAuthor         `json:"book_authors"`
	BookReaders []models.BookReader         `json:"book_readers"`
}

func (s *Storage) Save() error {
	data := fileStorage{
		Books:       s.books,
		Authors:     s.authors,
		Readers:     s.readers,
		BookAuthors: make([]models.BookAuthor, 0, len(s.bookAuthors)),
		BookReaders: make([]models.BookReader, 0, len(s.bookReaders)),
	}

	for bookAuthor := range s.bookAuthors {
		data.BookAuthors = append(data.BookAuthors, bookAuthor)
	}

	for bookReader := range s.bookReaders {
		data.BookReaders = append(data.BookReaders, bookReader)
	}

	file, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	return os.WriteFile(fileStoragePath, file, 0644)
}

func (s *Storage) Load() error {
	file, err := os.ReadFile(fileStoragePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return err
	}

	var data fileStorage

	if err := json.Unmarshal(file, &data); err != nil {
		return err
	}

	if data.Books == nil {
		data.Books = make(map[uuid.UUID]models.Book)
	}

	if data.Authors == nil {
		data.Authors = make(map[uuid.UUID]models.Author)
	}

	if data.Readers == nil {
		data.Readers = make(map[uuid.UUID]models.Reader)
	}

	s.books = data.Books
	s.authors = data.Authors
	s.readers = data.Readers

	s.bookAuthors = make(map[models.BookAuthor]struct{}, len(data.BookAuthors))

	for _, bookAuthor := range data.BookAuthors {
		s.bookAuthors[bookAuthor] = struct{}{}
	}

	s.bookReaders = make(map[models.BookReader]struct{}, len(data.BookReaders))

	for _, bookReader := range data.BookReaders {
		s.bookReaders[bookReader] = struct{}{}
	}

	return nil
}
