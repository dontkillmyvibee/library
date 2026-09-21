package local_storage

import "errors"

var (
	ErrBookAlreadyExists       = errors.New("book already exists")
	ErrBookNotFound            = errors.New("book not found")
	ErrBookAlreadyDeleted      = errors.New("book already deleted")
	ErrBookUnavailable         = errors.New("book unavailable")
	ErrAuthorAlreadyExists     = errors.New("author already exists")
	ErrAuthorNotFound          = errors.New("author not found")
	ErrAuthorAlreadyDeleted    = errors.New("author already deleted")
	ErrBookAuthorAlreadyExists = errors.New("book author already exists")
	ErrBookAuthorNotFound      = errors.New("book author not found")
	ErrDuplicateAuthor         = errors.New("duplicate author")
	ErrLastBookAuthor          = errors.New("last book author")
	ErrReaderAlreadyExists     = errors.New("reader already exists")
	ErrReaderNotFound          = errors.New("reader not found")
	ErrReaderAlreadyDeleted    = errors.New("reader already deleted")
	ErrBookReaderAlreadyExists = errors.New("book reader already exists")
	ErrReaderHasBook           = errors.New("reader has a book")
	ErrBookReaderNotFound      = errors.New("book reader not found")
)
