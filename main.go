package main

import (
	"fmt"

	"github.com/dontkillmyvibee/library.git/handlers"
	"github.com/dontkillmyvibee/library.git/local_storage"
	server2 "github.com/dontkillmyvibee/library.git/server"
)

func main() {
	localStorage := local_storage.NewStorage()
	if err := localStorage.Load(); err != nil {
		//тут тоже надо норм логировать
		fmt.Println(err.Error())
	}
	bookHandler := handlers.NewHTTPBookHandlers(localStorage)
	authorHandler := handlers.NewHTTPAuthorHandlers(localStorage)
	readerHandler := handlers.NewHTTPReaderHandlers(localStorage)
	bookReaderHandler := handlers.NewHTTPBookReaderHandlers(localStorage)
	server := server2.NewHTTPServer(bookHandler, authorHandler, readerHandler, bookReaderHandler)

	if err := server.StartServer(); err != nil {
		fmt.Println("err", err)
	}
}
