package handlers

import (
	"errors"
	"net/http"
	"uuid"

	"github.com/dontkillmyvibee/library.git/local_storage"
	"github.com/dontkillmyvibee/library.git/schemas"
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
	readerID, err := uuid.Parse(r.PathValue("readerID"))
	if err != nil {
		errDTO := schemas.NewError(
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
			ErrInvalidUUID.Error(),
		)
		http.Error(w, errDTO.ToJSONString(), http.StatusBadRequest)
		return
	}

	bookID, err := uuid.Parse(r.PathValue("bookID"))
	if err != nil {
		errDTO := schemas.NewError(
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
			ErrInvalidUUID.Error(),
		)
		http.Error(w, errDTO.ToJSONString(), http.StatusBadRequest)
		return
	}

	if err := h.localStorage.ReaderTakeBook(readerID, bookID); err != nil {
		if errors.Is(err, local_storage.ErrReaderNotFound) || errors.Is(err, local_storage.ErrBookNotFound) {
			errDTO := schemas.NewError(http.StatusNotFound, http.StatusText(http.StatusNotFound), err.Error())
			http.Error(w, errDTO.ToJSONString(), http.StatusNotFound)
			return
		}

		if errors.Is(err, local_storage.ErrBookUnavailable) || errors.Is(err, local_storage.ErrReaderHasBook) {
			errDTO := schemas.NewError(http.StatusConflict, http.StatusText(http.StatusConflict), err.Error())
			http.Error(w, errDTO.ToJSONString(), http.StatusConflict)
			return
		}

		errDTO := schemas.NewError(
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			"internal server error",
		)
		http.Error(w, errDTO.ToJSONString(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *HTTPBookReaderHandlers) ReturnBook(w http.ResponseWriter, r *http.Request) {
	readerID, err := uuid.Parse(r.PathValue("readerID"))
	if err != nil {
		errDTO := schemas.NewError(
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
			ErrInvalidUUID.Error(),
		)
		http.Error(w, errDTO.ToJSONString(), http.StatusBadRequest)
		return
	}

	bookID, err := uuid.Parse(r.PathValue("bookID"))
	if err != nil {
		errDTO := schemas.NewError(
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
			ErrInvalidUUID.Error(),
		)
		http.Error(w, errDTO.ToJSONString(), http.StatusBadRequest)
		return
	}

	if err := h.localStorage.ReaderReturnBook(readerID, bookID); err != nil {
		if errors.Is(err, local_storage.ErrBookReaderNotFound) ||
			errors.Is(err, local_storage.ErrBookNotFound) ||
			errors.Is(err, local_storage.ErrReaderNotFound) {
			errDTO := schemas.NewError(http.StatusNotFound, http.StatusText(http.StatusNotFound), err.Error())
			http.Error(w, errDTO.ToJSONString(), http.StatusNotFound)
			return
		}

		errDTO := schemas.NewError(
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			"internal server error",
		)
		http.Error(w, errDTO.ToJSONString(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
