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

	useDatabase := true

	var userRepo repository.UserRepository

	if useDatabase {
		connString := "postgres://postgres:postgres@localhost:5432/learning_db?sslmode=disable"

		dbStore, err := repository.NewPostgresUserStore(connString)
		if err != nil {
			log.Fatal("Unable to connect to database", err)
		}

		defer dbStore.Close()

		userRepo = dbStore
		fmt.Println("Using PostgrsSQL store")
	} else {
		userRepo = repository.NewUserStore()
		fmt.Println("Using in-memory store")
	}
	userService := services.NewUserService(userRepo)
	UserHandler := handlers.NewUserHandler(userService)

	mux := http.NewServeMux()
	mux.HandleFunc("/", UserHandler.GetAllUsers)

	err := http.ListenAndServe(":8888", mux)
	if err != nil {
		log.Fatalf("Unable to start the server")
	}

	fmt.Println("Server started on port 8888")
}
