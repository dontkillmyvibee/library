package server

import (
	"errors"
	"net/http"

	"github.com/dontkillmyvibee/library.git/handlers"
)

type HTTPServer struct {
	bookHandlers *handlers.HTTPBookHandlers
}

func NewHTTPServer(bookHandlers *handlers.HTTPBookHandlers) *HTTPServer {
	return &HTTPServer{
		bookHandlers: bookHandlers,
	}
}

func (s *HTTPServer) StartServer() error {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /books", s.bookHandlers.GetFilteredBooks)
	mux.HandleFunc("GET /books/{id}", s.bookHandlers.GetBook)
	mux.HandleFunc("POST /books", s.bookHandlers.CreateBook)
	mux.HandleFunc("DELETE /books/{id}", s.bookHandlers.DeleteBook)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
	return nil
}
