package local_storage

import "errors"

var (
	ErrBookAlreadyExist   = errors.New("book already exist")
	ErrBookNotFound       = errors.New("book not found")
	ErrBookAlreadyDeleted = errors.New("book already deleted")
)
