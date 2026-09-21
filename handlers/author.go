package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"uuid"

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
		errDTO := schemas.NewError(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), err.Error())
		http.Error(w, errDTO.ToJSONString(), http.StatusBadRequest)
		return
	}

	author := models.NewAuthor(
		*createAuthorRequest.FirstName,
		*createAuthorRequest.LastName,
		*createAuthorRequest.MiddleName,
	)
	if err := h.localStorage.AddAuthor(author); err != nil {
		if errors.Is(err, local_storage.ErrAuthorAlreadyExists) {
			errDTO := schemas.NewError(
				http.StatusBadRequest,
				http.StatusText(http.StatusBadRequest),
				err.Error(),
			)
			http.Error(w, errDTO.ToJSONString(), http.StatusBadRequest)
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

	response := schemas.CreateAuthorResponseSchema{
		ID:         author.ID,
		FirstName:  author.FirstName,
		LastName:   author.LastName,
		MiddleName: author.MiddleName,
		CreatedAt:  author.CreatedAt,
		UpdatedAt:  author.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Println("err:", err)
	}
}

func (h *HTTPAuthorHandlers) GetAuthor(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		errDTO := schemas.NewError(
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
			ErrInvalidUUID.Error(),
		)
		http.Error(w, errDTO.ToJSONString(), http.StatusBadRequest)
		return
	}

	author, err := h.localStorage.GetAuthor(id)
	if err != nil {
		if errors.Is(err, local_storage.ErrAuthorNotFound) {
			errDTO := schemas.NewError(http.StatusNotFound, http.StatusText(http.StatusNotFound), err.Error())
			http.Error(w, errDTO.ToJSONString(), http.StatusNotFound)
			return
		}

		if errors.Is(err, local_storage.ErrAuthorAlreadyDeleted) {
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

	response := schemas.GetAuthorResponseSchema{
		ID:         author.ID,
		FirstName:  author.FirstName,
		LastName:   author.LastName,
		MiddleName: author.MiddleName,
		CreatedAt:  author.CreatedAt,
		UpdatedAt:  author.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Println("err:", err)
	}
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Println("err:", err)
	}
}

func (h *HTTPAuthorHandlers) UpdateAuthor(w http.ResponseWriter, r *http.Request) {
	var updateAuthorRequest schemas.UpdateAuthorRequestSchema

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		errDTO := schemas.NewError(
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
			ErrInvalidUUID.Error(),
		)
		http.Error(w, errDTO.ToJSONString(), http.StatusBadRequest)
		return
	}

	if !helpers.DecodeJSONHelper(w, r, &updateAuthorRequest) {
		return
	}

	updateAuthorRequest.Normalize()

	if err := updateAuthorRequest.Validate(); err != nil {
		errDTO := schemas.NewError(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), err.Error())
		http.Error(w, errDTO.ToJSONString(), http.StatusBadRequest)
		return
	}

	updateData := models.UpdateAuthorData{
		FirstName:  *updateAuthorRequest.FirstName,
		LastName:   *updateAuthorRequest.LastName,
		MiddleName: *updateAuthorRequest.MiddleName,
	}

	author, err := h.localStorage.UpdateAuthor(id, updateData)
	if err != nil {
		if errors.Is(err, local_storage.ErrAuthorNotFound) {
			errDTO := schemas.NewError(http.StatusNotFound, http.StatusText(http.StatusNotFound), err.Error())
			http.Error(w, errDTO.ToJSONString(), http.StatusNotFound)
			return
		}

		if errors.Is(err, local_storage.ErrAuthorAlreadyExists) {
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

	response := schemas.UpdateAuthorResponseSchema{
		ID:         author.ID,
		FirstName:  author.FirstName,
		LastName:   author.LastName,
		MiddleName: author.MiddleName,
		CreatedAt:  author.CreatedAt,
		UpdatedAt:  author.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Println("err:", err)
	}
}

func (h *HTTPAuthorHandlers) DeleteAuthor(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		errDTO := schemas.NewError(
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
			ErrInvalidUUID.Error(),
		)
		http.Error(w, errDTO.ToJSONString(), http.StatusBadRequest)
		return
	}

	if err := h.localStorage.DeleteAuthor(id); err != nil {
		if errors.Is(err, local_storage.ErrAuthorNotFound) {
			errDTO := schemas.NewError(http.StatusNotFound, http.StatusText(http.StatusNotFound), err.Error())
			http.Error(w, errDTO.ToJSONString(), http.StatusNotFound)
			return
		}

		if errors.Is(err, local_storage.ErrAuthorAlreadyDeleted) {
			errDTO := schemas.NewError(http.StatusNotFound, http.StatusText(http.StatusNotFound), err.Error())
			http.Error(w, errDTO.ToJSONString(), http.StatusNotFound)
			return
		}

		if errors.Is(err, local_storage.ErrLastBookAuthor) {
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
