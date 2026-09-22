package handlers

import (
	"errors"
	"net/http"

	"github.com/dontkillmyvibee/library.git/handlers/helpers"
	"github.com/dontkillmyvibee/library.git/local_storage"
)

type HTTPBookReaderHandlers struct {
	localStorage *local_storage.Storage
}

func NewHTTPBookReaderHandlers(localStorage *local_storage.Storage) *HTTPBookReaderHandlers {
	return &HTTPBookReaderHandlers{
		localStorage: localStorage,
	}
}

func (h *HTTPBookReaderHandlers) TakeBook(w http.ResponseWriter, r *http.Request) {
	readerID, ok := helpers.ParseUUIDPath(w, r, "readerID")
	if !ok {
		return
	}

	bookID, ok := helpers.ParseUUIDPath(w, r, "bookID")
	if !ok {
		return
	}

	if err := h.localStorage.ReaderTakeBook(readerID, bookID); err != nil {
		if errors.Is(err, local_storage.ErrReaderNotFound) || errors.Is(err, local_storage.ErrBookNotFound) {
			helpers.InitError(w, http.StatusNotFound, err)
			return
		}

		if errors.Is(err, local_storage.ErrBookUnavailable) || errors.Is(err, local_storage.ErrReaderHasBook) {
			helpers.InitError(w, http.StatusConflict, err)
			return
		}

		helpers.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *HTTPBookReaderHandlers) ReturnBook(w http.ResponseWriter, r *http.Request) {
	readerID, ok := helpers.ParseUUIDPath(w, r, "readerID")
	if !ok {
		return
	}

	bookID, ok := helpers.ParseUUIDPath(w, r, "bookID")
	if !ok {
		return
	}

	if err := h.localStorage.ReaderReturnBook(readerID, bookID); err != nil {
		if errors.Is(err, local_storage.ErrBookReaderNotFound) ||
			errors.Is(err, local_storage.ErrBookNotFound) ||
			errors.Is(err, local_storage.ErrReaderNotFound) {
			helpers.InitError(w, http.StatusNotFound, err)
			return
		}

		helpers.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
}
