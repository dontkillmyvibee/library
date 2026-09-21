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

type HTTPBookHandlers struct {
	localStorage *local_storage.Storage
}

func NewHTTPBookHandlers(localStorage *local_storage.Storage) *HTTPBookHandlers {
	return &HTTPBookHandlers{
		localStorage: localStorage,
	}
}

func (h *HTTPBookHandlers) CreateBook(w http.ResponseWriter, r *http.Request) {
	var createBookRequest schemas.CreateBookRequestSchema

	if !helpers.DecodeJSONHelper(w, r, &createBookRequest) {
		return
	}

	createBookRequest.Normalize()
	if err := createBookRequest.Validate(); err != nil {
		errDTO := schemas.NewError(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), err.Error())
		http.Error(w, errDTO.ToJSONString(), http.StatusBadRequest)
		return
	}

	book := models.NewBook(*createBookRequest.Title, *createBookRequest.Description)
	if _, err := h.localStorage.AddBook(book, *createBookRequest.AuthorIDs); err != nil {
		if errors.Is(err, local_storage.ErrBookAlreadyExists) ||
			errors.Is(err, local_storage.ErrDuplicateAuthor) ||
			errors.Is(err, local_storage.ErrBookAuthorAlreadyExist) {
			errDTO := schemas.NewError(http.StatusConflict, http.StatusText(http.StatusConflict), err.Error())
			http.Error(w, errDTO.ToJSONString(), http.StatusConflict)
			return
		}
		if errors.Is(err, local_storage.ErrAuthorNotFound) {
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

	authors := h.localStorage.GetAuthorsByBookID(book.ID)

	response := schemas.CreateBookResponseSchema{
		ID:          book.ID,
		Title:       book.Title,
		Description: book.Description,
		Authors:     make([]schemas.Authors, 0, len(authors)),
		IsAvailable: book.IsAvailable,
		CreatedAt:   book.CreatedAt,
		UpdatedAt:   book.UpdatedAt,
	}

	for _, author := range authors {
		response.Authors = append(response.Authors, schemas.Authors{
			ID:         author.ID,
			FirstName:  author.FirstName,
			LastName:   author.LastName,
			MiddleName: author.MiddleName,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Println("err:", err)
	}
}

func (h *HTTPBookHandlers) GetBook(w http.ResponseWriter, r *http.Request) {
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

	book, err := h.localStorage.GetBook(id)
	if err != nil {
		if errors.Is(err, local_storage.ErrBookNotFound) {
			errDTO := schemas.NewError(
				http.StatusNotFound,
				http.StatusText(http.StatusNotFound),
				err.Error(),
			)
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

	authors := h.localStorage.GetAuthorsByBookID(book.ID)

	response := schemas.CreateBookResponseSchema{
		ID:          book.ID,
		Title:       book.Title,
		Description: book.Description,
		Authors:     make([]schemas.Authors, 0, len(authors)),
		IsAvailable: book.IsAvailable,
		CreatedAt:   book.CreatedAt,
		UpdatedAt:   book.UpdatedAt,
	}

	for _, author := range authors {
		response.Authors = append(response.Authors, schemas.Authors{
			ID:         author.ID,
			FirstName:  author.FirstName,
			LastName:   author.LastName,
			MiddleName: author.MiddleName,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// это надо логировать по идеи
		fmt.Println("err:", err)
	}
}

func (h *HTTPBookHandlers) GetFilteredBooks(w http.ResponseWriter, r *http.Request) {
	available := r.URL.Query().Get("available")

	var books map[uuid.UUID]models.Book

	switch available {
	case "":
		books = h.localStorage.GetAllBooks()
	case "true":
		books = h.localStorage.GetAllAvailableBooks()
	case "false":
		books = h.localStorage.GetAllUnavailableBooks()
	default:
		errDTO := schemas.NewError(
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
			ErrInvalidQueryParameter.Error(),
		)
		http.Error(w, errDTO.ToJSONString(), http.StatusBadRequest)
		return
	}

	response := make(schemas.GetAllBooksResponseSchema, 0, len(books))

	for _, book := range books {
		authors := h.localStorage.GetAuthorsByBookID(book.ID)

		responseAuthors := make([]schemas.Authors, 0, len(authors))

		for _, author := range authors {
			responseAuthors = append(responseAuthors, schemas.Authors{
				ID:         author.ID,
				FirstName:  author.FirstName,
				LastName:   author.LastName,
				MiddleName: author.MiddleName,
			})
		}

		response = append(response, schemas.GetBookResponseSchema{
			ID:          book.ID,
			Title:       book.Title,
			Description: book.Description,
			Authors:     responseAuthors,
			IsAvailable: book.IsAvailable,
			CreatedAt:   book.CreatedAt,
			UpdatedAt:   book.UpdatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Println("err:", err)
	}
}

func (h *HTTPBookHandlers) DeleteBook(w http.ResponseWriter, r *http.Request) {
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

	if err := h.localStorage.DeleteBook(id); err != nil {
		if errors.Is(err, local_storage.ErrBookNotFound) {
			errDTO := schemas.NewError(
				http.StatusNotFound,
				http.StatusText(http.StatusNotFound),
				local_storage.ErrBookNotFound.Error(),
			)
			http.Error(w, errDTO.ToJSONString(), http.StatusNotFound)
			return
		}

		if errors.Is(err, local_storage.ErrBookAlreadyDeleted) {
			errDTO := schemas.NewError(
				http.StatusConflict,
				http.StatusText(http.StatusConflict),
				local_storage.ErrBookAlreadyDeleted.Error(),
			)
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
