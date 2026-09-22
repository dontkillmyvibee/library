package handlers

import (
	"errors"
	"net/http"

	"github.com/dontkillmyvibee/library.git/handlers/helpers"
	"github.com/dontkillmyvibee/library.git/local_storage"
	"github.com/dontkillmyvibee/library.git/local_storage/models"
	"github.com/dontkillmyvibee/library.git/schemas"
)

type HTTPAuthorHandlers struct {
	localStorage *local_storage.Storage
}

func NewHTTPAuthorHandlers(localStorage *local_storage.Storage) *HTTPAuthorHandlers {
	return &HTTPAuthorHandlers{
		localStorage: localStorage,
	}
}

func (h *HTTPAuthorHandlers) CreateAuthor(w http.ResponseWriter, r *http.Request) {
	var createAuthorRequest schemas.CreateAuthorRequestSchema

	if !helpers.DecodeJSONHelper(w, r, &createAuthorRequest) {
		return
	}

	createAuthorRequest.Normalize()

	if err := createAuthorRequest.Validate(); err != nil {
		helpers.InitError(w, http.StatusBadRequest, err)
		return
	}

	author := models.NewAuthor(
		*createAuthorRequest.FirstName,
		*createAuthorRequest.LastName,
		*createAuthorRequest.MiddleName,
	)
	if err := h.localStorage.AddAuthor(author); err != nil {
		if errors.Is(err, local_storage.ErrAuthorAlreadyExists) {
			helpers.InitError(w, http.StatusBadRequest, err)
			return
		}

		helpers.InternalServerError(w)
		return
	}

	response := schemas.CreateAuthorResponseSchema{
		ID:         author.ID,
		FirstName:  author.FirstName,
		LastName:   author.LastName,
		MiddleName: author.MiddleName,
		CreatedAt:  author.CreatedAt,
		UpdatedAt:  author.UpdatedAt,
	}

	helpers.EncodeJSONHelper(w, http.StatusOK, response)
}

func (h *HTTPAuthorHandlers) GetAuthor(w http.ResponseWriter, r *http.Request) {
	id, ok := helpers.ParseUUIDPath(w, r, "id")
	if !ok {
		return
	}

	author, err := h.localStorage.GetAuthor(id)
	if err != nil {
		if errors.Is(err, local_storage.ErrAuthorNotFound) || errors.Is(err, local_storage.ErrAuthorAlreadyDeleted) {
			helpers.InitError(w, http.StatusNotFound, err)
			return
		}

		helpers.InternalServerError(w)
		return
	}

	response := schemas.GetAuthorResponseSchema{
		ID:         author.ID,
		FirstName:  author.FirstName,
		LastName:   author.LastName,
		MiddleName: author.MiddleName,
		CreatedAt:  author.CreatedAt,
		UpdatedAt:  author.UpdatedAt,
	}

	helpers.EncodeJSONHelper(w, http.StatusOK, response)
}

func (h *HTTPAuthorHandlers) GetAllAuthors(w http.ResponseWriter, r *http.Request) {
	authors := h.localStorage.GetAllAuthors()

	response := make(schemas.GetAllAuthorsResponseSchema, 0, len(authors))

	for _, author := range authors {
		response = append(response, schemas.GetAuthorResponseSchema{
			ID:         author.ID,
			FirstName:  author.FirstName,
			LastName:   author.LastName,
			MiddleName: author.MiddleName,
			CreatedAt:  author.CreatedAt,
			UpdatedAt:  author.UpdatedAt,
		})
	}

	helpers.EncodeJSONHelper(w, http.StatusOK, response)
}

func (h *HTTPAuthorHandlers) DeleteAuthor(w http.ResponseWriter, r *http.Request) {
	id, ok := helpers.ParseUUIDPath(w, r, "id")
	if !ok {
		return
	}

	if err := h.localStorage.DeleteAuthor(id); err != nil {
		if errors.Is(err, local_storage.ErrAuthorNotFound) || errors.Is(err, local_storage.ErrAuthorAlreadyDeleted) {
			helpers.InitError(w, http.StatusNotFound, err)
			return
		}

		if errors.Is(err, local_storage.ErrLastBookAuthor) {
			helpers.InitError(w, http.StatusConflict, err)
			return
		}

		helpers.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
}
