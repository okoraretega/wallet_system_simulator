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
	var walletStore repository.WalletRepository

	if useDatabase {
		connString := "postgres://postgres:postgres@localhost:5432/bank_db?sslmode=disable"

		dbStore, err := repository.NewPostgresStore(connString)
		if err != nil {
			log.Fatal("Unable to connect to database", err)
		}

		defer dbStore.Close()

		userRepo = dbStore
		walletStore = dbStore
		fmt.Println("Using PostgrsSQL store")
	} else {
		userRepo = repository.NewUserStore()
		fmt.Println("Using in-memory store")
	}

	mux := http.NewServeMux()

	// Wallets
	walletService := services.NewWalletService(walletStore)
	walletHandler := handlers.NewWalletHandler(walletService)

	mux.HandleFunc("/wallets", walletHandler.GetAllWalets)
	mux.HandleFunc("/wallets/", walletHandler.GetWalletByUserId)

	// Users
	userService := services.NewUserService(userRepo)
	UserHandler := handlers.NewUserHandler(userService)

	mux.HandleFunc("/users/", UserHandler.HandleUsers)
	mux.HandleFunc("/create", UserHandler.CreateUser)
	mux.HandleFunc("/users", UserHandler.GetAllUsers)

	err := http.ListenAndServe(":8888", mux)
	if err != nil {
		log.Fatalf("Unable to start the server")
	}

	fmt.Println("Server started on port 8888")
}
