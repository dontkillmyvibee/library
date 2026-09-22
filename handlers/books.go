package handlers

import (
	"errors"
	"net/http"
	"uuid"

	errors2 "github.com/dontkillmyvibee/library.git/handlers/errors"
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
			errors.Is(err, local_storage.ErrBookAuthorAlreadyExists) {
			helpers.InitError(w, http.StatusConflict, err)
			return
		}
		if errors.Is(err, local_storage.ErrAuthorNotFound) {
			helpers.InitError(w, http.StatusNotFound, err)
			return
		}

		helpers.InternalServerError(w)
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

	helpers.EncodeJSONHelper(w, http.StatusOK, response)
}

func (h *HTTPBookHandlers) UpdateBook(w http.ResponseWriter, r *http.Request) {
	var updateBookRequest schemas.UpdateBookRequestSchema

	id, ok := helpers.ParseUUIDPath(w, r, "readerID")
	if !ok {
		return
	}

	if !helpers.DecodeJSONHelper(w, r, &updateBookRequest) {
		return
	}

	updateBookRequest.Normalize()

	if err := updateBookRequest.Validate(); err != nil {
		helpers.InitError(w, http.StatusBadRequest, err)
		return
	}

	updateData := models.UpdateBookData{
		Title:       *updateBookRequest.Title,
		Description: *updateBookRequest.Description,
		AuthorIDs:   *updateBookRequest.AuthorIDs,
	}

	book, err := h.localStorage.UpdateBook(id, updateData)
	if err != nil {
		if errors.Is(err, local_storage.ErrBookNotFound) || errors.Is(err, local_storage.ErrAuthorNotFound) {
			helpers.InitError(w, http.StatusNotFound, err)
			return
		}

		if errors.Is(err, local_storage.ErrDuplicateAuthor) || errors.Is(err, local_storage.ErrBookAlreadyExists) {
			helpers.InitError(w, http.StatusConflict, err)
			return
		}

		helpers.InternalServerError(w)
		return
	}

	authors := h.localStorage.GetAuthorsByBookID(book.ID)

	response := schemas.UpdateBookResponseSchema{
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

	helpers.EncodeJSONHelper(w, http.StatusOK, response)
}

func (h *HTTPBookHandlers) GetBook(w http.ResponseWriter, r *http.Request) {
	id, ok := helpers.ParseUUIDPath(w, r, "id")
	if !ok {
		return
	}

	book, err := h.localStorage.GetBook(id)
	if err != nil {
		if errors.Is(err, local_storage.ErrBookNotFound) {
			helpers.InitError(w, http.StatusNotFound, err)
			return
		}

		helpers.InternalServerError(w)
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

	helpers.EncodeJSONHelper(w, http.StatusOK, response)
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
		helpers.InitError(w, http.StatusBadRequest, errors2.ErrInvalidQueryParameter)
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

	helpers.EncodeJSONHelper(w, http.StatusOK, response)
}

func (h *HTTPBookHandlers) DeleteBook(w http.ResponseWriter, r *http.Request) {
	id, ok := helpers.ParseUUIDPath(w, r, "id")
	if !ok {
		return
	}

	if err := h.localStorage.DeleteBook(id); err != nil {
		if errors.Is(err, local_storage.ErrBookNotFound) || errors.Is(err, local_storage.ErrBookAlreadyDeleted) {
			helpers.InitError(w, http.StatusNotFound, err)
			return
		}

		helpers.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
}
