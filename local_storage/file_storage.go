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
	BookAuthors []models.BookAuthor         `json:"book_authors"`
}

func (s *Storage) Save() error {
	data := fileStorage{
		Books:       s.books,
		Authors:     s.authors,
		BookAuthors: make([]models.BookAuthor, 0, len(s.bookAuthors)),
	}

	for bookAuthor := range s.bookAuthors {
		data.BookAuthors = append(data.BookAuthors, bookAuthor)
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

	s.books = data.Books
	s.authors = data.Authors
	s.bookAuthors = make(map[models.BookAuthor]struct{}, len(data.BookAuthors))

	for _, bookAuthor := range data.BookAuthors {
		s.bookAuthors[bookAuthor] = struct{}{}
	}

	return nil
}
