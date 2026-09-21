package local_storage

import "errors"

var (
	ErrBookAlreadyExists      = errors.New("book already exists")
	ErrBookNotFound           = errors.New("book not found")
	ErrBookAlreadyDeleted     = errors.New("book already deleted")
	ErrAuthorAlreadyExists    = errors.New("author already exists")
	ErrAuthorNotFound         = errors.New("author not found")
	ErrAuthorAlreadyDeleted   = errors.New("author already deleted")
	ErrBookAuthorAlreadyExist = errors.New("book author already exists")
	ErrBookAuthorNotFound     = errors.New("book author not found")
	ErrDuplicateAuthor        = errors.New("duplicate author")
	ErrLastBookAuthor         = errors.New("last book author")
)
