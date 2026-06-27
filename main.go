package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/okoraretega/doc_stream_server/handlers"
	"github.com/okoraretega/doc_stream_server/repository"
	"github.com/okoraretega/doc_stream_server/services"
)

var JWT_SECRET = "edwed32eedqEdDEDrwdweD" //This should be in your env

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

	mux.Handle("/wallets", TestMiddleWare(http.HandlerFunc(walletHandler.GetAllWalets)))
	mux.Handle("/wallets/", TestMiddleWare(http.HandlerFunc(walletHandler.GetWalletByUserId)))
	mux.HandleFunc("/wallets/transfer", walletHandler.Transfer)

	// Users
	userService := services.NewUserService(userRepo)
	UserHandler := handlers.NewUserHandler(userService)

	mux.HandleFunc("/users/", UserHandler.HandleUsers)
	mux.HandleFunc("/create", UserHandler.CreateUser)
	mux.HandleFunc("/login", UserHandler.Login)
	mux.Handle("/users", AuthMiddleWare(http.HandlerFunc(UserHandler.GetAllUsers)))

	fmt.Println("Server starting on port 8888")

	err := http.ListenAndServe(":8888", mux)
	if err != nil {
		log.Fatalf("Unable to start the server")
	}
}

func TestMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		fmt.Printf("Hello from the %v\n", path)

		next.ServeHTTP(w, r)
	})
}

func AuthMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		parts := strings.Split(authHeader, " ")

		if len(parts) > 2 || parts[0] != "Bearer" {
			http.Error(w, "Please provide a valid authorization header", http.StatusBadRequest)
			return
		}

		token := parts[1]
		claims, err := ValidateToken(token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), "user", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ValidateToken(tokenStr string) (jwt.Claims, error) {

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("Wrong signature method")
		}

		return []byte(JWT_SECRET), nil
	})

	if err != nil {
		return nil, fmt.Errorf("Failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("No claims found")
	}

	return claims, nil
}
