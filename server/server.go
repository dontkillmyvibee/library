package server

import (
	"errors"
	"net/http"

	"github.com/dontkillmyvibee/library.git/handlers"
)

type HTTPServer struct {
	bookHandlers   *handlers.HTTPBookHandlers
	authorHandlers *handlers.HTTPAuthorHandlers
}

func NewHTTPServer(bookHandlers *handlers.HTTPBookHandlers, authorHandlers *handlers.HTTPAuthorHandlers) *HTTPServer {
	return &HTTPServer{
		bookHandlers:   bookHandlers,
		authorHandlers: authorHandlers,
	}
}

func (s *HTTPServer) StartServer() error {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /books", s.bookHandlers.GetFilteredBooks)
	mux.HandleFunc("GET /books/{id}", s.bookHandlers.GetBook)
	mux.HandleFunc("POST /books", s.bookHandlers.CreateBook)
	mux.HandleFunc("DELETE /books/{id}", s.bookHandlers.DeleteBook)

	mux.HandleFunc("GET /authors", s.authorHandlers.GetAllAuthors)
	mux.HandleFunc("GET /authors/{id}", s.authorHandlers.GetAuthor)
	mux.HandleFunc("POST /authors", s.authorHandlers.CreateAuthor)
	mux.HandleFunc("DELETE /authors/{id}", s.authorHandlers.DeleteAuthor)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
	return nil
}
