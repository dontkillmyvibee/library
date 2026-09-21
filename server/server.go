package server

import (
	"errors"
	"net/http"

	"github.com/dontkillmyvibee/library.git/handlers"
)

type HTTPServer struct {
	bookHandlers       *handlers.HTTPBookHandlers
	authorHandlers     *handlers.HTTPAuthorHandlers
	readerHandlers     *handlers.HTTPReaderHandlers
	bookReaderHandlers *handlers.HTTPBookReaderHandlers
}

func NewHTTPServer(
	bookHandlers *handlers.HTTPBookHandlers,
	authorHandlers *handlers.HTTPAuthorHandlers,
	readerHandlers *handlers.HTTPReaderHandlers,
	bookReaderHandlers *handlers.HTTPBookReaderHandlers,
) *HTTPServer {
	return &HTTPServer{
		bookHandlers:       bookHandlers,
		authorHandlers:     authorHandlers,
		readerHandlers:     readerHandlers,
		bookReaderHandlers: bookReaderHandlers,
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

	mux.HandleFunc("GET /readers", s.readerHandlers.GetAllReaders)
	mux.HandleFunc("GET /readers/{id}", s.readerHandlers.GetReader)
	mux.HandleFunc("POST /readers", s.readerHandlers.CreateReader)
	mux.HandleFunc("DELETE /readers/{id}", s.readerHandlers.DeleteReader)

	mux.HandleFunc("POST /readers/{readerID}/books/{bookID}", s.bookReaderHandlers.TakeBook)
	mux.HandleFunc("DELETE /readers/{readerID}/books/{bookID}", s.bookReaderHandlers.ReturnBook)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
	return nil
}
