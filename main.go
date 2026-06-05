package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/okoraretega/doc_stream_server/handlers"
	"github.com/okoraretega/doc_stream_server/repository"
	"github.com/okoraretega/doc_stream_server/services"
)

func main() {

	userStore := repository.NewUserStore()
	userService := services.NewUserService(userStore)
	UserHandler := handlers.NewUserHandler(userService)

	mux := http.NewServeMux()
	mux.HandleFunc("/", UserHandler.GetAllUsers)

	err := http.ListenAndServe(":8888", mux)
	if err != nil {
		log.Fatalf("Unable to start the server")
	}

	fmt.Println("Server started on port 8888")
}
