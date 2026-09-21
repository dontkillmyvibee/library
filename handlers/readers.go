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
		errDTO := schemas.NewError(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), err.Error())
		http.Error(w, errDTO.ToJSONString(), http.StatusBadRequest)
		return
	}

	reader := models.NewReader(
		*createReaderRequest.FirstName,
		*createReaderRequest.MiddleName,
		*createReaderRequest.LastName,
	)

	if err := h.localStorage.AddReader(reader); err != nil {
		if errors.Is(err, local_storage.ErrReaderAlreadyExists) {
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

	response := schemas.CreateReaderResponseSchema{
		ID:         reader.ID,
		FirstName:  reader.FirstName,
		LastName:   reader.LastName,
		MiddleName: reader.MiddleName,
		CreatedAt:  reader.CreatedAt,
		UpdatedAt:  reader.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Println("err:", err)
	}
}

func (h *HTTPReaderHandlers) UpdateReader(w http.ResponseWriter, r *http.Request) {
	var updateReaderRequest schemas.UpdateReaderRequestSchema

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

	if !helpers.DecodeJSONHelper(w, r, &updateReaderRequest) {
		return
	}

	updateReaderRequest.Normalize()

	if err := updateReaderRequest.Validate(); err != nil {
		errDTO := schemas.NewError(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), err.Error())
		http.Error(w, errDTO.ToJSONString(), http.StatusBadRequest)
		return
	}

	updateData := models.UpdateReaderData{
		FirstName:  *updateReaderRequest.FirstName,
		LastName:   *updateReaderRequest.LastName,
		MiddleName: *updateReaderRequest.MiddleName,
	}

	reader, err := h.localStorage.UpdateReader(id, updateData)
	if err != nil {
		if errors.Is(err, local_storage.ErrReaderNotFound) {
			errDTO := schemas.NewError(http.StatusNotFound, http.StatusText(http.StatusNotFound), err.Error())
			http.Error(w, errDTO.ToJSONString(), http.StatusNotFound)
			return
		}

		if errors.Is(err, local_storage.ErrReaderAlreadyExists) {
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

	response := schemas.UpdateReaderResponseSchema{
		ID:         reader.ID,
		FirstName:  reader.FirstName,
		LastName:   reader.LastName,
		MiddleName: reader.MiddleName,
		CreatedAt:  reader.CreatedAt,
		UpdatedAt:  reader.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Println("err:", err)
	}
}

func (h *HTTPReaderHandlers) GetReader(w http.ResponseWriter, r *http.Request) {
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

	reader, err := h.localStorage.GetReader(id)
	if err != nil {
		if errors.Is(err, local_storage.ErrReaderNotFound) || errors.Is(err, local_storage.ErrReaderAlreadyDeleted) {
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

	response := schemas.GetReaderResponseSchema{
		ID:         reader.ID,
		FirstName:  reader.FirstName,
		LastName:   reader.LastName,
		MiddleName: reader.MiddleName,
		CreatedAt:  reader.CreatedAt,
		UpdatedAt:  reader.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Println("err:", err)
	}
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Println("err:", err)
	}
}

func (h *HTTPReaderHandlers) DeleteReader(w http.ResponseWriter, r *http.Request) {
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

	if err := h.localStorage.DeleteReader(id); err != nil {
		if errors.Is(err, local_storage.ErrReaderNotFound) || errors.Is(err, local_storage.ErrReaderAlreadyDeleted) {
			errDTO := schemas.NewError(http.StatusNotFound, http.StatusText(http.StatusNotFound), err.Error())
			http.Error(w, errDTO.ToJSONString(), http.StatusNotFound)
			return
		}

		if errors.Is(err, local_storage.ErrReaderHasBook) {
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
