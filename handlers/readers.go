package handlers

import (
	"errors"
	"net/http"

	"github.com/dontkillmyvibee/library.git/handlers/helpers"
	"github.com/dontkillmyvibee/library.git/local_storage"
	"github.com/dontkillmyvibee/library.git/local_storage/models"
	"github.com/dontkillmyvibee/library.git/schemas"
)

type HTTPReaderHandlers struct {
	localStorage *local_storage.Storage
}

func NewHTTPReaderHandlers(localStorage *local_storage.Storage) *HTTPReaderHandlers {
	return &HTTPReaderHandlers{
		localStorage: localStorage,
	}
}

func (h *HTTPReaderHandlers) CreateReader(w http.ResponseWriter, r *http.Request) {
	var createReaderRequest schemas.CreateReaderRequestSchema

	if !helpers.DecodeJSONHelper(w, r, &createReaderRequest) {
		return
	}

	createReaderRequest.Normalize()

	if err := createReaderRequest.Validate(); err != nil {
		helpers.InitError(w, http.StatusBadRequest, err)
		return
	}

	reader := models.NewReader(
		*createReaderRequest.FirstName,
		*createReaderRequest.MiddleName,
		*createReaderRequest.LastName,
	)

	if err := h.localStorage.AddReader(reader); err != nil {
		if errors.Is(err, local_storage.ErrReaderAlreadyExists) {
			helpers.InitError(w, http.StatusConflict, err)
			return
		}

		helpers.InternalServerError(w)
		return
	}

	response := schemas.CreateReaderResponseSchema{
		ID:         reader.ID,
		FirstName:  reader.FirstName,
		LastName:   reader.LastName,
		MiddleName: reader.MiddleName,
		CreatedAt:  reader.CreatedAt,
		UpdatedAt:  reader.UpdatedAt,
	}

	helpers.EncodeJSONHelper(w, http.StatusOK, response)
}

func (h *HTTPReaderHandlers) GetReader(w http.ResponseWriter, r *http.Request) {
	id, ok := helpers.ParseUUIDPath(w, r, "id")
	if !ok {
		return
	}

	reader, err := h.localStorage.GetReader(id)
	if err != nil {
		if errors.Is(err, local_storage.ErrReaderNotFound) || errors.Is(err, local_storage.ErrReaderAlreadyDeleted) {
			helpers.InitError(w, http.StatusNotFound, err)
			return
		}

		helpers.InternalServerError(w)
		return
	}

	response := schemas.GetReaderResponseSchema{
		ID:         reader.ID,
		FirstName:  reader.FirstName,
		LastName:   reader.LastName,
		MiddleName: reader.MiddleName,
		CreatedAt:  reader.CreatedAt,
		UpdatedAt:  reader.UpdatedAt,
	}

	helpers.EncodeJSONHelper(w, http.StatusOK, response)
}

func (h *HTTPReaderHandlers) GetAllReaders(w http.ResponseWriter, r *http.Request) {
	readers := h.localStorage.GetAllReaders()

	response := make(schemas.GetAllReadersResponseSchema, 0, len(readers))

	for _, reader := range readers {
		response = append(response, schemas.GetReaderResponseSchema{
			ID:         reader.ID,
			FirstName:  reader.FirstName,
			LastName:   reader.LastName,
			MiddleName: reader.MiddleName,
			CreatedAt:  reader.CreatedAt,
			UpdatedAt:  reader.UpdatedAt,
		})
	}

	helpers.EncodeJSONHelper(w, http.StatusOK, response)
}

func (h *HTTPReaderHandlers) DeleteReader(w http.ResponseWriter, r *http.Request) {
	id, ok := helpers.ParseUUIDPath(w, r, "id")
	if !ok {
		return
	}

	if err := h.localStorage.DeleteReader(id); err != nil {
		if errors.Is(err, local_storage.ErrReaderNotFound) || errors.Is(err, local_storage.ErrReaderAlreadyDeleted) {
			helpers.InitError(w, http.StatusNotFound, err)
			return
		}

		if errors.Is(err, local_storage.ErrReaderHasBook) {
			helpers.InitError(w, http.StatusConflict, err)
			return
		}

		helpers.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
}
