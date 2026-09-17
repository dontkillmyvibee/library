package main

import (
	"fmt"

	"github.com/dontkillmyvibee/library.git/handlers"
	"github.com/dontkillmyvibee/library.git/local_storage"
	server2 "github.com/dontkillmyvibee/library.git/server"
)

func main() {
	localStorage := local_storage.NewStorage()
	bookHandler := handlers.NewHTTPBookHandlers(localStorage)
	server := server2.NewHTTPServer(bookHandler)

	if err := server.StartServer(); err != nil {
		fmt.Println("err", err)
	}
}
